package usecases

import (
	"context"
	"encoding/json"
	"fmt"
	"ticket-service/internal/modules/ticket"
	"ticket-service/internal/modules/ticket/models/entity"
	"ticket-service/internal/modules/ticket/models/request"
	"ticket-service/internal/modules/ticket/models/response"
	"ticket-service/internal/pkg/errors"
	"ticket-service/internal/pkg/log"
	"time"

	kafkaConfluent "ticket-service/internal/pkg/kafka/confluent"

	"go.elastic.co/apm"
)

type queryUsecase struct {
	ticketRepositoryQuery ticket.MongodbRepositoryQuery
	kafkaProducer         kafkaConfluent.Producer
	logger                log.Logger
}

func NewQueryUsecase(tmq ticket.MongodbRepositoryQuery, kp kafkaConfluent.Producer, log log.Logger) ticket.UsecaseQuery {
	return queryUsecase{
		ticketRepositoryQuery: tmq,
		kafkaProducer:         kp,
		logger:                log,
	}
}

func (q queryUsecase) FindTickets(origCtx context.Context, payload request.TicketReq) (*response.TicketResp, error) {
	domain := "addressUsecase-FindTickets"
	span, ctx := apm.StartSpanOptions(origCtx, domain, "function", apm.SpanOptions{
		Start:  time.Now(),
		Parent: apm.TraceContext{},
	})
	defer span.End()

	resp := <-q.ticketRepositoryQuery.FindOfflineTicketByCountry(ctx, payload)
	if resp.Error != nil {
		msg := "Error query ticket"
		q.logger.Error(ctx, msg, fmt.Sprintf("%+v", resp.Error))
		return nil, resp.Error
	}

	if resp.Data == nil {
		msg := "Ticket Not Found"
		q.logger.Error(ctx, msg, fmt.Sprintf("%+v", payload))
		return nil, errors.NotFound("ticket not found")
	}

	availableTicket, ok := resp.Data.(*[]entity.Ticket)
	if !ok {
		return nil, errors.InternalServerError("cannot parsing data")
	}

	var result response.TicketResp
	var tag string
	emptyCounter := 0
	var collectionData = make([]response.Ticket, 0)
	for _, value := range *availableTicket {
		isSold := false
		if value.TotalRemaining == 0 {
			isSold = true
			emptyCounter = emptyCounter + 1
		}
		collectionData = append(collectionData, response.Ticket{
			TicketType:    value.TicketType,
			TicketPrice:   fmt.Sprintf("$%d", value.TicketPrice),
			ContinentName: value.ContinentName,
			ContinentCode: value.ContinentCode,
			CountryName:   value.Country.Name,
			CountryCode:   value.Country.Code,
			IsSold:        isSold,
		})
		tag = value.Tag
	}
	result.Tickets = collectionData

	if emptyCounter >= 4 {
		res := <-q.ticketRepositoryQuery.FindTicketByLowestPrice(ctx, tag)
		if res.Error != nil {
			msg := "Error query ticket"
			q.logger.Error(ctx, msg, fmt.Sprintf("%+v", res.Error))
			return nil, res.Error
		}

		if res.Data == nil {
			msg := "Ticket Not Found"
			q.logger.Error(ctx, msg, fmt.Sprintf("%+v", payload))
			return nil, errors.NotFound("ticket not found")
		}

		suggestionTicket, ok := res.Data.(*[]entity.Ticket)
		if !ok {
			return nil, errors.InternalServerError("cannot parsing data")
		}

		if len(*suggestionTicket) > 0 {
			resp := <-q.ticketRepositoryQuery.FindOfflineTicketByCountryCode(ctx, (*suggestionTicket)[0].Country.Code, tag)
			if resp.Error != nil {
				msg := "Error query ticket"
				q.logger.Error(ctx, msg, fmt.Sprintf("%+v", resp.Error))
				return nil, resp.Error
			}

			if resp.Data == nil {
				msg := "Ticket Not Found"
				q.logger.Error(ctx, msg, fmt.Sprintf("%+v", payload))
				return nil, errors.NotFound("ticket not found")
			}

			availableTicket, ok := resp.Data.(*[]entity.Ticket)
			if !ok {
				return nil, errors.InternalServerError("cannot parsing data")
			}

			var collectionData = make([]response.SuggestionTicket, 0)
			for _, value := range *availableTicket {
				isSold := false
				if value.TotalRemaining == 0 {
					isSold = true
				}
				discountPrice := value.TicketPrice * 80 / 100
				collectionData = append(collectionData, response.SuggestionTicket{
					TicketType:          value.TicketType,
					NormalTicketPrice:   fmt.Sprintf("$%d", value.TicketPrice),
					DiscountTicketPrice: fmt.Sprintf("$%d", discountPrice),
					Discount:            "20%",
					ContinentName:       value.ContinentName,
					ContinentCode:       value.ContinentCode,
					CountryName:         value.Country.Name,
					CountryCode:         value.Country.Code,
					IsSold:              isSold,
				})
			}
			result.Suggestion = collectionData
		}

		updateOnlineTicket := request.CreateOnlineTicketReq{
			Tag:         tag,
			CountryCode: payload.CountryCode,
		}
		marshaledKafkaData, _ := json.Marshal(updateOnlineTicket)
		topic := "concert-update-online-bank-ticket"
		q.kafkaProducer.Publish(topic, marshaledKafkaData, nil)
		q.logger.Info(ctx, fmt.Sprintf("Send kafka update online bank ticket, tag : %s", tag), fmt.Sprintf("%+v", updateOnlineTicket))
	}

	return &result, nil

}

func (q queryUsecase) FindOnlineTicket(origCtx context.Context, payload request.TicketReq) (*response.Ticket, error) {
	domain := "addressUsecase-FindTickets"
	span, ctx := apm.StartSpanOptions(origCtx, domain, "function", apm.SpanOptions{
		Start:  time.Now(),
		Parent: apm.TraceContext{},
	})
	defer span.End()

	offlineTicketData := <-q.ticketRepositoryQuery.FindOfflineTicketByCountry(ctx, payload)
	if offlineTicketData.Error != nil {
		msg := "Error query ticket"
		q.logger.Error(ctx, msg, fmt.Sprintf("%+v", offlineTicketData.Error))
		return nil, offlineTicketData.Error
	}

	if offlineTicketData.Data == nil {
		msg := "Offline ticket Not Found"
		q.logger.Error(ctx, msg, fmt.Sprintf("%+v", payload))
		return nil, errors.NotFound("offline ticket not found")
	}

	offlineTicket, ok := offlineTicketData.Data.(*[]entity.Ticket)
	if !ok {
		return nil, errors.InternalServerError("cannot parsing data")
	}

	emptyCounter := 0
	for _, ticket := range *offlineTicket {
		if ticket.TotalRemaining == 0 {
			emptyCounter = emptyCounter + 1
		}
	}

	if emptyCounter < 4 && len(*offlineTicket) > 0 {
		return nil, errors.BadRequest("Offline ticket still available")
	}

	resp := <-q.ticketRepositoryQuery.FindOnlineTicketByCountry(ctx, payload)
	if resp.Error != nil {
		msg := "Error query ticket"
		q.logger.Error(ctx, msg, fmt.Sprintf("%+v", resp.Error))
		return nil, resp.Error
	}

	if resp.Data == nil {
		msg := "Ticket Not Found"
		q.logger.Error(ctx, msg, fmt.Sprintf("%+v", payload))
		return nil, errors.NotFound("ticket not found")
	}

	availableTicket, ok := resp.Data.(*entity.Ticket)
	if !ok {
		return nil, errors.InternalServerError("cannot parsing data")
	}

	var result response.Ticket
	result.CountryCode = availableTicket.Country.Code
	result.CountryName = availableTicket.Country.Name
	result.ContinentCode = availableTicket.ContinentCode
	result.ContinentName = availableTicket.ContinentName
	result.TicketPrice = fmt.Sprintf("$%d", availableTicket.TicketPrice)
	result.TicketType = availableTicket.TicketType
	if availableTicket.TotalRemaining == 0 {
		result.IsSold = true
	}

	return &result, nil

}

// func (q queryUsecase) FindAvailableTicket(origCtx context.Context) ([]response.TicketCountry, error) {
// 	domain := "addressUsecase-FindAvailableTicket"
// 	span, ctx := apm.StartSpanOptions(origCtx, domain, "function", apm.SpanOptions{
// 		Start:  time.Now(),
// 		Parent: apm.TraceContext{},
// 	})
// 	defer span.End()

// 	resp := <-q.ticketRepositoryQuery.FindTotalAvalailableTicket(ctx)
// 	if resp.Error != nil {
// 		msg := "Error query ticket"
// 		q.logger.Error(ctx, msg, fmt.Sprintf("%+v", resp.Error))
// 		return nil, resp.Error
// 	}

// 	if resp.Data == nil {
// 		msg := "Ticket Not Found"
// 		q.logger.Error(ctx, msg, resp)
// 		return nil, errors.NotFound("ticket not found")
// 	}

// 	availableTicket, ok := resp.Data.(*[]entity.AggregateTotalTicket)
// 	if !ok {
// 		return nil, errors.InternalServerError("cannot parsing data")
// 	}

// 	var collectionData = make([]response.TicketCountry, 0)
// 	for _, value := range *availableTicket {
// 		isSold := false
// 		if value.TotalAvailableTicket == 0 {
// 			isSold = true
// 		}
// 		collectionData = append(collectionData, response.TicketCountry{
// 			CountryName: value.CountryName,
// 			CountryCode: value.Id,
// 			IsSold:      isSold,
// 		})
// 	}
// 	return collectionData, nil
// }
