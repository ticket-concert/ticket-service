package ticket

import (
	"context"
	"ticket-service/internal/modules/ticket/models/request"
	"ticket-service/internal/modules/ticket/models/response"
	wrapper "ticket-service/internal/pkg/helpers"
)

type UsecaseQuery interface {
	FindTickets(origCtx context.Context, payload request.TicketReq) (*response.TicketResp, error)
	FindOnlineTicket(origCtx context.Context, payload request.TicketReq) (*response.Ticket, error)
	// FindAvailableTicket(origCtx context.Context) ([]response.TicketCountry, error)
}

type MongodbRepositoryQuery interface {
	FindOfflineTicketByCountry(ctx context.Context, payload request.TicketReq) <-chan wrapper.Result
	FindTicketByLowestPrice(ctx context.Context, tag string) <-chan wrapper.Result
	FindOnlineTicketByCountry(ctx context.Context, payload request.TicketReq) <-chan wrapper.Result
	FindOfflineTicketByCountryCode(ctx context.Context, countryCode string, tag string) <-chan wrapper.Result
	// FindTotalAvalailableTicket(ctx context.Context) <-chan wrapper.Result
}
