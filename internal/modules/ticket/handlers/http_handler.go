package handlers

import (
	"ticket-service/internal/modules/ticket"
	"ticket-service/internal/modules/ticket/models/request"
	"ticket-service/internal/pkg/errors"
	"ticket-service/internal/pkg/helpers"
	"ticket-service/internal/pkg/log"
	"ticket-service/internal/pkg/redis"

	middlewares "ticket-service/configs/middleware"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type TicketHttpHandler struct {
	TicketUsecaseQuery ticket.UsecaseQuery
	Logger             log.Logger
	Validator          *validator.Validate
}

func InitTicketHttpHandler(app *fiber.App, tuq ticket.UsecaseQuery, log log.Logger, redisClient redis.Collections) {
	handler := &TicketHttpHandler{
		TicketUsecaseQuery: tuq,
		Logger:             log,
		Validator:          validator.New(),
	}
	middlewares := middlewares.NewMiddlewares(redisClient)
	route := app.Group("/api/tickets")

	route.Get("/v1/list", middlewares.VerifyBearer(), handler.GetTickets)
	// route.Get("/v1/list-available", middlewares.VerifyBearer(), handler.GetAvailableTicket)
	route.Get("/v1/online", middlewares.VerifyBearer(), handler.GetOnlineTicket)
}

func (t TicketHttpHandler) GetTickets(c *fiber.Ctx) error {
	req := new(request.TicketReq)
	if err := c.QueryParser(req); err != nil {
		return helpers.RespError(c, t.Logger, errors.BadRequest("bad request"))
	}

	if err := t.Validator.Struct(req); err != nil {
		return helpers.RespError(c, t.Logger, errors.BadRequest(err.Error()))
	}
	resp, err := t.TicketUsecaseQuery.FindTickets(c.Context(), *req)
	if err != nil {
		return helpers.RespCustomError(c, t.Logger, err)
	}
	return helpers.RespSuccess(c, t.Logger, resp, "Get ticket success")
}

func (t TicketHttpHandler) GetOnlineTicket(c *fiber.Ctx) error {
	req := new(request.TicketReq)
	if err := c.QueryParser(req); err != nil {
		return helpers.RespError(c, t.Logger, errors.BadRequest("bad request"))
	}

	if err := t.Validator.Struct(req); err != nil {
		return helpers.RespError(c, t.Logger, errors.BadRequest(err.Error()))
	}
	resp, err := t.TicketUsecaseQuery.FindOnlineTicket(c.Context(), *req)
	if err != nil {
		return helpers.RespCustomError(c, t.Logger, err)
	}
	return helpers.RespSuccess(c, t.Logger, resp, "Get online ticket success")
}

// func (t TicketHttpHandler) GetAvailableTicket(c *fiber.Ctx) error {
// 	resp, err := t.TicketUsecaseQuery.FindAvailableTicket(c.Context())
// 	if err != nil {
// 		return helpers.RespCustomError(c, t.Logger, err)
// 	}
// 	return helpers.RespSuccess(c, t.Logger, resp, "Get available ticket success")
// }
