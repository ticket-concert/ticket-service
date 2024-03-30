package main

import (
	"fmt"
	logGo "log"
	"strconv"
	"ticket-service/configs"
	ticketHandler "ticket-service/internal/modules/ticket/handlers"
	ticketRepoQuery "ticket-service/internal/modules/ticket/repositories/queries"
	ticketUsecase "ticket-service/internal/modules/ticket/usecases"
	"ticket-service/internal/pkg/apm"
	"ticket-service/internal/pkg/databases/mongodb"
	graceful "ticket-service/internal/pkg/gs"
	"ticket-service/internal/pkg/helpers"
	kafkaConfluent "ticket-service/internal/pkg/kafka/confluent"
	"ticket-service/internal/pkg/log"
	"ticket-service/internal/pkg/redis"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"go.elastic.co/apm/module/apmfiber"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

// @BasePath	/
func main() {
	// Init Config
	configs.InitConfig()

	// Init Elastic APM Config or DD Apm
	if configs.GetConfig().Datadog.DatadogEnabled == "true" {
		tracer.Start(
			tracer.WithAgentAddr(fmt.Sprintf("%s:%s", configs.GetConfig().Datadog.DatadogHost, configs.GetConfig().Datadog.DatadogPort)),
			tracer.WithEnv(configs.GetConfig().Datadog.DatadogEnv),
			tracer.WithService(configs.GetConfig().Datadog.DatadogService),
		)
		defer tracer.Stop()
	} else {
		// temp handling until we move all to Datadog
		apm.InitConnection(configs.GetConfig().APMElastic.APMUrl, configs.GetConfig().APMElastic.APMSecretToken)
	}
	// Init MongoDB Connection
	mongo := mongodb.MongoImpl{}
	mongo.SetCollections(&mongo)
	mongo.InitConnection(configs.GetConfig().MongoDB.MongoMasterDBUrl, configs.GetConfig().MongoDB.MongoSlaveDBUrl)
	// Init Logger
	logZap := log.SetupLogger(configs.GetConfig().ServiceName)
	log.Init(logZap)

	// Init BlacklistedEmail
	helpers.InitReadBlackListEmail()

	// Init Kafka Config
	kafkaConfluent.InitKafkaConfig(configs.GetConfig().Kafka.KafkaUrl, configs.GetConfig().Kafka.KafkaUsername, configs.GetConfig().Kafka.KafkaPassword)

	// Init instance fiber
	app := fiber.New(fiber.Config{
		BodyLimit: 30 * 1024 * 1024,
	})
	app.Use(apmfiber.Middleware(apmfiber.WithTracer(apm.GetTracer())))
	app.Use(recover.New())
	app.Use(cors.New())
	app.Use(pprof.New())
	if configs.GetConfig().AppsLimiter {
		app.Use(limiter.New(limiter.Config{
			Max:               100,
			Expiration:        1 * time.Minute,
			LimiterMiddleware: limiter.SlidingWindow{},
		}))
	}
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${latency} ${method} ${path}\n",
		// Format:       `${time} {"router_activity" : [${status},"${latency}","${method}","${path}"], "query_param":${queryParams}, "body_param":${body}}` + "\n",
		TimeInterval: time.Millisecond,
		TimeFormat:   "02-01-2006 15:04:05",
		TimeZone:     "Indonesia/Jakarta",
	}))
	shutdownDelay, _ := strconv.Atoi(configs.GetConfig().ShutDownDelay)
	// graceful shutdown setup
	gs := &graceful.GracefulShutdown{
		Timeout:        5 * time.Second,
		GracefulPeriod: time.Duration(shutdownDelay) * time.Second,
	}
	app.Get("/healthz", gs.LivenessCheck)
	app.Get("/readyz", gs.ReadinessCheck)
	gs.Enable(app)

	setHttp(app, gs)

	//=== listen port ===//
	if err := app.Listen(fmt.Sprintf(":%s", configs.GetConfig().ServicePort)); err != nil {
		logGo.Fatal(err)
	}
}

func setHttp(app *fiber.App, gs *graceful.GracefulShutdown) {
	// Init Redis
	redisClient := redis.InitConnection(configs.GetConfig().Redis.RedisDB, configs.GetConfig().Redis.RedisHost, configs.GetConfig().Redis.RedisPort,
		configs.GetConfig().Redis.RedisPassword, configs.GetConfig().Redis.RedisAppConfig)
	// Init Jwt
	helperImpl := &helpers.JwtImpl{}
	helperImpl.InitConfig(configs.GetConfig().Jwt.JwtPrivateKey, configs.GetConfig().Jwt.JwtPublicKey,
		configs.GetConfig().Jwt.JwtRefreshPrivateKey, configs.GetConfig().Jwt.JwtRefreshPublicKey)

	logger := log.GetLogger()
	mongoMasterClient := mongodb.NewMongoDBLogger(mongodb.GetMasterConn(), mongodb.GetMasterDBName(), logger)
	mongoSlaveClient := mongodb.NewMongoDBLogger(mongodb.GetSlaveConn(), mongodb.GetMasterDBName(), logger)
	kafkaProducer, err := kafkaConfluent.NewProducer(kafkaConfluent.GetConfig().GetKafkaConfig(configs.GetConfig().ServiceName, true), logger)
	if err != nil {
		panic(err)
	}
	gs.Register(
		mongoMasterClient,
		mongoSlaveClient,
		graceful.FnWithError(redisClient.Close),
		kafkaProducer,
	)

	ticketQueryMongodbRepo := ticketRepoQuery.NewQueryMongodbRepository(mongoSlaveClient, logger)
	ticketUsecaseQuery := ticketUsecase.NewQueryUsecase(ticketQueryMongodbRepo, kafkaProducer, logger)

	// set module
	ticketHandler.InitTicketHttpHandler(app, ticketUsecaseQuery, logger, redisClient)

}
