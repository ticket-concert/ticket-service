package queries

import (
	"context"
	"ticket-service/internal/modules/ticket"
	"ticket-service/internal/modules/ticket/models/entity"
	"ticket-service/internal/modules/ticket/models/request"
	"ticket-service/internal/pkg/databases/mongodb"
	wrapper "ticket-service/internal/pkg/helpers"
	"ticket-service/internal/pkg/log"

	"go.mongodb.org/mongo-driver/bson"
)

type queryMongodbRepository struct {
	mongoDb mongodb.Collections
	logger  log.Logger
}

func NewQueryMongodbRepository(mongodb mongodb.Collections, log log.Logger) ticket.MongodbRepositoryQuery {
	return &queryMongodbRepository{
		mongoDb: mongodb,
		logger:  log,
	}
}

func (q queryMongodbRepository) FindOfflineTicketByCountry(ctx context.Context, payload request.TicketReq) <-chan wrapper.Result {
	var tickets []entity.Ticket
	output := make(chan wrapper.Result)

	go func() {
		resp := <-q.mongoDb.FindMany(mongodb.FindMany{
			Result:         &tickets,
			CollectionName: "ticket-detail",
			Filter: bson.M{
				"ticketType":   bson.M{"$ne": "Online"},
				"country.code": payload.CountryCode,
				"eventId":      payload.EventId,
			},
			Sort: &mongodb.Sort{
				FieldName: "ticketPrice",
				By:        mongodb.SortAscending,
			},
		}, ctx)
		output <- resp
		close(output)
	}()

	return output
}

func (q queryMongodbRepository) FindOfflineTicketByCountryCode(ctx context.Context, countryCode string, tag string) <-chan wrapper.Result {
	var tickets []entity.Ticket
	output := make(chan wrapper.Result)

	go func() {
		resp := <-q.mongoDb.FindMany(mongodb.FindMany{
			Result:         &tickets,
			CollectionName: "ticket-detail",
			Filter: bson.M{
				"ticketType":   bson.M{"$ne": "Online"},
				"country.code": countryCode,
				"tag":          tag,
			},
			Sort: &mongodb.Sort{
				FieldName: "ticketPrice",
				By:        mongodb.SortAscending,
			},
		}, ctx)
		output <- resp
		close(output)
	}()

	return output
}

func (q queryMongodbRepository) FindTicketByLowestPrice(ctx context.Context, tag string) <-chan wrapper.Result {
	var tickets []entity.Ticket
	output := make(chan wrapper.Result)

	go func() {
		resp := <-q.mongoDb.FindAllData(mongodb.FindAllData{
			Result:         &tickets,
			CollectionName: "ticket-detail",
			Filter: bson.M{
				"$where":     "this.totalRemaining > 0",
				"ticketType": bson.M{"$ne": "Online"},
				"tag":        tag,
			},
			Sort: &mongodb.Sort{
				FieldName: "ticketPrice",
				By:        mongodb.SortAscending,
			},
			Page: 1,
			Size: 1,
		}, ctx)
		output <- resp
		close(output)
	}()

	return output
}

func (q queryMongodbRepository) FindOnlineTicketByCountry(ctx context.Context, payload request.TicketReq) <-chan wrapper.Result {
	var tickets entity.Ticket
	output := make(chan wrapper.Result)

	go func() {
		resp := <-q.mongoDb.FindOne(mongodb.FindOne{
			Result:         &tickets,
			CollectionName: "ticket-detail",
			Filter: bson.M{
				"ticketType":   "Online",
				"country.code": payload.CountryCode,
				"eventId":      payload.EventId,
			},
		}, ctx)
		output <- resp
		close(output)
	}()

	return output
}

// func (q queryMongodbRepository) FindTotalAvalailableTicket(ctx context.Context) <-chan wrapper.Result {
// 	var ticket []entity.AggregateTotalTicket
// 	output := make(chan wrapper.Result)

// 	go func() {
// 		resp := <-q.mongoDb.Aggregate(mongodb.Aggregate{
// 			Result:         &ticket,
// 			CollectionName: "ticket-detail",
// 			Filter: []bson.M{
// 				{
// 					"$group": bson.M{
// 						"_id":                  "$country.code",
// 						"countryName":          bson.M{"$first": "$country.name"},
// 						"totalAvailableTicket": bson.M{"$sum": "$totalRemaining"},
// 					},
// 				},
// 			},
// 		}, ctx)
// 		output <- resp
// 		close(output)
// 	}()

// 	return output
// }
