// user_http_handler_test.go

package handlers_test

import (
	"net/http/httptest"
	"testing"
	"ticket-service/internal/modules/ticket/handlers"
	"ticket-service/internal/modules/ticket/models/response"
	"ticket-service/internal/pkg/errors"
	mockcert "ticket-service/mocks/modules/ticket"
	mocklog "ticket-service/mocks/pkg/log"
	mockredis "ticket-service/mocks/pkg/redis"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/valyala/fasthttp"
)

type ticketHttpHandlerTestSuite struct {
	suite.Suite

	cUQ       *mockcert.UsecaseQuery
	cLog      *mocklog.Logger
	validator *validator.Validate
	handler   *handlers.TicketHttpHandler
	cRedis    *mockredis.Collections
	app       *fiber.App
}

func (suite *ticketHttpHandlerTestSuite) SetupTest() {
	suite.cUQ = new(mockcert.UsecaseQuery)
	suite.cLog = new(mocklog.Logger)
	suite.validator = validator.New()
	suite.cRedis = new(mockredis.Collections)
	suite.handler = &handlers.TicketHttpHandler{
		TicketUsecaseQuery: suite.cUQ,
		Logger:             suite.cLog,
		Validator:          suite.validator,
	}
	suite.app = fiber.New()
	handlers.InitTicketHttpHandler(suite.app, suite.cUQ, suite.cLog, suite.cRedis)
}

func TestUserHttpHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(ticketHttpHandlerTestSuite))
}

func (suite *ticketHttpHandlerTestSuite) TestGetTickets() {

	response := &response.TicketResp{
		Tickets: []response.Ticket{
			{
				TicketType:  "type",
				TicketPrice: "price",
			},
		},
		Suggestion: []response.SuggestionTicket{},
	}
	suite.cUQ.On("FindTickets", mock.Anything, mock.Anything).Return(response, nil)
	suite.cLog.On("Info", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	ctx := suite.app.AcquireCtx(&fasthttp.RequestCtx{})
	req := httptest.NewRequest(fiber.MethodGet, "/v1/list?countryCode=1&eventID=1", nil)
	req.Header.Set("Content-Type", "application/json")
	ctx.Request().SetRequestURI("/v1/list?countryCode=1&eventID=1")
	ctx.Request().Header.SetMethod(fiber.MethodGet)
	ctx.Request().Header.SetContentType("application/json")

	err := suite.handler.GetTickets(ctx)
	assert.Nil(suite.T(), err)
}

func (suite *ticketHttpHandlerTestSuite) TestGetTicketsErrQuery() {

	response := &response.TicketResp{
		Tickets: []response.Ticket{
			{
				TicketType:  "type",
				TicketPrice: "price",
			},
		},
		Suggestion: []response.SuggestionTicket{},
	}
	suite.cUQ.On("FindTickets", mock.Anything, mock.Anything).Return(response, nil)
	suite.cLog.On("Info", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	ctx := suite.app.AcquireCtx(&fasthttp.RequestCtx{})
	req := httptest.NewRequest(fiber.MethodGet, "/v1/list", nil)
	req.Header.Set("Content-Type", "application/json")
	// ctx.Request().SetRequestURI("/v1/list?countryCode=aa&eventID=1")
	ctx.Request().Header.SetMethod(fiber.MethodGet)
	ctx.Request().Header.SetContentType("application/json")

	err := suite.handler.GetTickets(ctx)
	assert.Nil(suite.T(), err)
}

func (suite *ticketHttpHandlerTestSuite) TestGetTicketsErr() {

	suite.cUQ.On("FindTickets", mock.Anything, mock.Anything).Return(nil, errors.BadRequest("error"))
	suite.cLog.On("Info", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	ctx := suite.app.AcquireCtx(&fasthttp.RequestCtx{})
	req := httptest.NewRequest(fiber.MethodGet, "/v1/list?countryCode=1&eventID=1", nil)
	req.Header.Set("Content-Type", "application/json")
	ctx.Request().SetRequestURI("/v1/list?countryCode=1&eventID=1")
	ctx.Request().Header.SetMethod(fiber.MethodGet)
	ctx.Request().Header.SetContentType("application/json")

	err := suite.handler.GetTickets(ctx)
	assert.Nil(suite.T(), err)
}

func (suite *ticketHttpHandlerTestSuite) TestGetOnlineTicket() {

	response := &response.Ticket{
		TicketType:  "type",
		TicketPrice: "price",
	}
	suite.cUQ.On("FindOnlineTicket", mock.Anything, mock.Anything).Return(response, nil)
	suite.cLog.On("Info", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	ctx := suite.app.AcquireCtx(&fasthttp.RequestCtx{})
	req := httptest.NewRequest(fiber.MethodGet, "/v1/online?countryCode=1&eventID=1", nil)
	req.Header.Set("Content-Type", "application/json")
	ctx.Request().SetRequestURI("/v1/online?countryCode=1&eventID=1")
	ctx.Request().Header.SetMethod(fiber.MethodGet)
	ctx.Request().Header.SetContentType("application/json")

	err := suite.handler.GetOnlineTicket(ctx)
	assert.Nil(suite.T(), err)
}

func (suite *ticketHttpHandlerTestSuite) TestGetOnlineTicketErrQuery() {

	response := &response.Ticket{
		TicketType:  "type",
		TicketPrice: "price",
	}
	suite.cUQ.On("FindOnlineTicket", mock.Anything, mock.Anything).Return(response, nil)
	suite.cLog.On("Info", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	ctx := suite.app.AcquireCtx(&fasthttp.RequestCtx{})
	req := httptest.NewRequest(fiber.MethodGet, "/v1/online?countryCode=1&eventID=1", nil)
	req.Header.Set("Content-Type", "application/json")
	// ctx.Request().SetRequestURI("/v1/online?countryCode=1&eventID=1")
	ctx.Request().Header.SetMethod(fiber.MethodGet)
	ctx.Request().Header.SetContentType("application/json")

	err := suite.handler.GetOnlineTicket(ctx)
	assert.Nil(suite.T(), err)
}

func (suite *ticketHttpHandlerTestSuite) TestGetOnlineTicketErr() {
	suite.cUQ.On("FindOnlineTicket", mock.Anything, mock.Anything).Return(nil, errors.BadRequest("error"))
	suite.cLog.On("Info", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	ctx := suite.app.AcquireCtx(&fasthttp.RequestCtx{})
	req := httptest.NewRequest(fiber.MethodGet, "/v1/online?countryCode=1&eventID=1", nil)
	req.Header.Set("Content-Type", "application/json")
	ctx.Request().SetRequestURI("/v1/online?countryCode=1&eventID=1")
	ctx.Request().Header.SetMethod(fiber.MethodGet)
	ctx.Request().Header.SetContentType("application/json")

	err := suite.handler.GetOnlineTicket(ctx)
	assert.Nil(suite.T(), err)
}
