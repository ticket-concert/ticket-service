// query_test.go
package usecases_test

import (
	"context"
	"testing"

	"ticket-service/internal/modules/ticket"
	ticketEntity "ticket-service/internal/modules/ticket/models/entity"
	ticketRequest "ticket-service/internal/modules/ticket/models/request"
	uc "ticket-service/internal/modules/ticket/usecases"
	"ticket-service/internal/pkg/errors"
	"ticket-service/internal/pkg/helpers"
	mockcert "ticket-service/mocks/modules/ticket"
	mockkafka "ticket-service/mocks/pkg/kafka"
	mocklog "ticket-service/mocks/pkg/log"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type QueryUsecaseTestSuite struct {
	suite.Suite
	mockTicketRepositoryQuery *mockcert.MongodbRepositoryQuery
	mockKafkaProducer         *mockkafka.Producer
	mockLogger                *mocklog.Logger
	usecase                   ticket.UsecaseQuery
	ctx                       context.Context
}

func (suite *QueryUsecaseTestSuite) SetupTest() {
	suite.mockTicketRepositoryQuery = &mockcert.MongodbRepositoryQuery{}
	suite.mockKafkaProducer = &mockkafka.Producer{}
	suite.mockLogger = &mocklog.Logger{}
	suite.ctx = context.Background()
	suite.usecase = uc.NewQueryUsecase(
		suite.mockTicketRepositoryQuery,
		suite.mockKafkaProducer,
		suite.mockLogger,
	)
}
func TestQueryUsecaseTestSuite(t *testing.T) {
	suite.Run(t, new(QueryUsecaseTestSuite))
}

func (suite *QueryUsecaseTestSuite) TestFindTicketSuccess() {
	// Arrange
	payload := ticketRequest.TicketReq{
		CountryCode: "code",
		EventId:     "id",
	}

	mockTicketQueryResponse := helpers.Result{
		Data: &[]ticketEntity.Ticket{
			{
				TicketId:       "id",
				EventId:        "id",
				TicketType:     "type",
				TicketPrice:    50,
				TotalRemaining: 5,
				Tag:            "tag",
			},
		},
		Error: nil,
	}
	suite.mockTicketRepositoryQuery.On("FindOfflineTicketByCountry", mock.Anything, payload).Return(mockChannel(mockTicketQueryResponse))
	suite.mockTicketRepositoryQuery.On("FindTicketByLowestPrice", mock.Anything, mock.Anything).Return(mockChannel(mockTicketQueryResponse))
	suite.mockTicketRepositoryQuery.On("FindOfflineTicketByCountryCode", mock.Anything, mock.Anything).Return(mockChannel(mockTicketQueryResponse))
	suite.mockKafkaProducer.On("Publish", mock.Anything, mock.Anything, mock.Anything)
	suite.mockLogger.On("Info", mock.Anything, mock.Anything, mock.Anything)

	// Act
	result, err := suite.usecase.FindTickets(suite.ctx, payload)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
}

func (suite *QueryUsecaseTestSuite) TestFindTicketErr() {
	// Arrange
	payload := ticketRequest.TicketReq{
		CountryCode: "code",
		EventId:     "id",
	}

	mockTicketQueryResponse := helpers.Result{
		Data: &[]ticketEntity.Ticket{
			{
				TicketId:       "id",
				EventId:        "id",
				TicketType:     "type",
				TicketPrice:    50,
				TotalRemaining: 5,
				Tag:            "tag",
			},
		},
		Error: errors.BadRequest("error"),
	}
	suite.mockTicketRepositoryQuery.On("FindOfflineTicketByCountry", mock.Anything, payload).Return(mockChannel(mockTicketQueryResponse))
	suite.mockTicketRepositoryQuery.On("FindTicketByLowestPrice", mock.Anything, mock.Anything).Return(mockChannel(mockTicketQueryResponse))
	suite.mockTicketRepositoryQuery.On("FindOfflineTicketByCountryCode", mock.Anything, mock.Anything).Return(mockChannel(mockTicketQueryResponse))
	suite.mockLogger.On("Error", mock.Anything, mock.Anything, mock.Anything)

	// Act
	_, err := suite.usecase.FindTickets(suite.ctx, payload)

	// Assert
	assert.Error(suite.T(), err)

	mockTicketQueryResponse = helpers.Result{
		Data:  nil,
		Error: nil,
	}
	suite.mockTicketRepositoryQuery.On("FindOfflineTicketByCountry", mock.Anything, payload).Return(mockChannel(mockTicketQueryResponse))
	suite.mockTicketRepositoryQuery.On("FindTicketByLowestPrice", mock.Anything, mock.Anything).Return(mockChannel(mockTicketQueryResponse))
	suite.mockTicketRepositoryQuery.On("FindOfflineTicketByCountryCode", mock.Anything, mock.Anything).Return(mockChannel(mockTicketQueryResponse))
	suite.mockKafkaProducer.On("Publish", mock.Anything, mock.Anything, mock.Anything)
	suite.mockLogger.On("Error", mock.Anything, mock.Anything, mock.Anything)

	// Act
	_, err = suite.usecase.FindTickets(suite.ctx, payload)

	// Assert
	assert.Error(suite.T(), err)
}

func (suite *QueryUsecaseTestSuite) TestFindTicketErrParse() {
	// Arrange
	payload := ticketRequest.TicketReq{
		CountryCode: "code",
		EventId:     "id",
	}

	mockTicketQueryResponse := helpers.Result{
		Data: &ticketEntity.Country{
			Name: "name",
			Code: "code",
		},
		Error: nil,
	}
	suite.mockTicketRepositoryQuery.On("FindOfflineTicketByCountry", mock.Anything, payload).Return(mockChannel(mockTicketQueryResponse))
	suite.mockTicketRepositoryQuery.On("FindTicketByLowestPrice", mock.Anything, mock.Anything).Return(mockChannel(mockTicketQueryResponse))
	suite.mockTicketRepositoryQuery.On("FindOfflineTicketByCountryCode", mock.Anything, mock.Anything).Return(mockChannel(mockTicketQueryResponse))
	suite.mockKafkaProducer.On("Publish", mock.Anything, mock.Anything, mock.Anything)
	suite.mockLogger.On("Error", mock.Anything, mock.Anything, mock.Anything)

	// Act
	_, err := suite.usecase.FindTickets(suite.ctx, payload)

	// Assert
	assert.Error(suite.T(), err)
}

func (suite *QueryUsecaseTestSuite) TestFindTicketSuggestion() {
	// Arrange
	payload := ticketRequest.TicketReq{
		CountryCode: "code",
		EventId:     "id",
	}

	mockTicketQueryResponse := helpers.Result{
		Data: &[]ticketEntity.Ticket{
			{
				TicketId:       "id",
				EventId:        "id",
				TicketType:     "type",
				TicketPrice:    50,
				TotalRemaining: 0,
				Tag:            "tag",
			},
			{
				TicketId:       "id",
				EventId:        "id",
				TicketType:     "type",
				TicketPrice:    50,
				TotalRemaining: 0,
				Tag:            "tag",
			},
			{
				TicketId:       "id",
				EventId:        "id",
				TicketType:     "type",
				TicketPrice:    50,
				TotalRemaining: 0,
				Tag:            "tag",
			},
			{
				TicketId:       "id",
				EventId:        "id",
				TicketType:     "type",
				TicketPrice:    50,
				TotalRemaining: 0,
				Tag:            "tag",
			},
		},
		Error: nil,
	}
	suite.mockTicketRepositoryQuery.On("FindOfflineTicketByCountry", mock.Anything, payload).Return(mockChannel(mockTicketQueryResponse))
	suite.mockTicketRepositoryQuery.On("FindTicketByLowestPrice", mock.Anything, mock.Anything).Return(mockChannel(mockTicketQueryResponse))
	suite.mockTicketRepositoryQuery.On("FindOfflineTicketByCountryCode", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockTicketQueryResponse))
	suite.mockKafkaProducer.On("Publish", mock.Anything, mock.Anything, mock.Anything)
	suite.mockLogger.On("Info", mock.Anything, mock.Anything, mock.Anything)

	// Act
	result, err := suite.usecase.FindTickets(suite.ctx, payload)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
}

func (suite *QueryUsecaseTestSuite) TestFindTicketSuggestionErrParse() {
	// Arrange
	payload := ticketRequest.TicketReq{
		CountryCode: "code",
		EventId:     "id",
	}
	mockTicketLowestPrice := helpers.Result{
		Data:  nil,
		Error: errors.BadRequest("error"),
	}
	suite.mockTicketRepositoryQuery.On("FindOfflineTicketByCountry", mock.Anything, payload).Return(mockChannel(getMockTicketSold()))
	suite.mockTicketRepositoryQuery.On("FindTicketByLowestPrice", mock.Anything, mock.Anything).Return(mockChannel(mockTicketLowestPrice))
	suite.mockTicketRepositoryQuery.On("FindOfflineTicketByCountryCode", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(getMockTicketSold()))
	suite.mockKafkaProducer.On("Publish", mock.Anything, mock.Anything, mock.Anything)
	suite.mockLogger.On("Error", mock.Anything, mock.Anything, mock.Anything)

	// Act
	_, err := suite.usecase.FindTickets(suite.ctx, payload)

	// Assert
	assert.Error(suite.T(), err)
}

func (suite *QueryUsecaseTestSuite) TestFindTicketSuggestionErrNil() {
	// Arrange
	payload := ticketRequest.TicketReq{
		CountryCode: "code",
		EventId:     "id",
	}

	nilTicketLowestPrice := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	suite.mockTicketRepositoryQuery.On("FindOfflineTicketByCountry", mock.Anything, payload).Return(mockChannel(getMockTicketSold()))
	suite.mockTicketRepositoryQuery.On("FindTicketByLowestPrice", mock.Anything, mock.Anything).Return(mockChannel(nilTicketLowestPrice))
	suite.mockTicketRepositoryQuery.On("FindOfflineTicketByCountryCode", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(getMockTicketSold()))
	suite.mockKafkaProducer.On("Publish", mock.Anything, mock.Anything, mock.Anything)
	suite.mockLogger.On("Error", mock.Anything, mock.Anything, mock.Anything)

	// Act
	_, err := suite.usecase.FindTickets(suite.ctx, payload)
	suite.T().Log(err)

	// Assert
	assert.Error(suite.T(), err)
}

func (suite *QueryUsecaseTestSuite) TestFindTicketErrMarshal() {
	// Arrange
	payload := ticketRequest.TicketReq{
		CountryCode: "code",
		EventId:     "id",
	}

	nilTicketLowestPrice := helpers.Result{
		Data: &ticketEntity.Country{
			Code: "code",
			Name: "name",
		},
		Error: nil,
	}

	suite.mockTicketRepositoryQuery.On("FindOfflineTicketByCountry", mock.Anything, payload).Return(mockChannel(getMockTicketSold()))
	suite.mockTicketRepositoryQuery.On("FindTicketByLowestPrice", mock.Anything, mock.Anything).Return(mockChannel(nilTicketLowestPrice))
	suite.mockTicketRepositoryQuery.On("FindOfflineTicketByCountryCode", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(getMockTicketSold()))
	suite.mockKafkaProducer.On("Publish", mock.Anything, mock.Anything, mock.Anything)
	suite.mockLogger.On("Error", mock.Anything, mock.Anything, mock.Anything)

	// Act
	_, err := suite.usecase.FindTickets(suite.ctx, payload)
	suite.T().Log(err)

	// Assert
	assert.Error(suite.T(), err)
}

func (suite *QueryUsecaseTestSuite) TestFindTicketOfflineErr() {
	// Arrange
	payload := ticketRequest.TicketReq{
		CountryCode: "code",
		EventId:     "id",
	}

	mockOfflineTicket := helpers.Result{
		Data:  nil,
		Error: errors.BadRequest("error"),
	}

	suite.mockTicketRepositoryQuery.On("FindOfflineTicketByCountry", mock.Anything, payload).Return(mockChannel(getMockTicketSold()))
	suite.mockTicketRepositoryQuery.On("FindTicketByLowestPrice", mock.Anything, mock.Anything).Return(mockChannel(getMockTicketSold()))
	suite.mockTicketRepositoryQuery.On("FindOfflineTicketByCountryCode", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockOfflineTicket))
	suite.mockKafkaProducer.On("Publish", mock.Anything, mock.Anything, mock.Anything)
	suite.mockLogger.On("Error", mock.Anything, mock.Anything, mock.Anything)

	// Act
	_, err := suite.usecase.FindTickets(suite.ctx, payload)
	suite.T().Log(err)

	// Assert
	assert.Error(suite.T(), err)
}

func (suite *QueryUsecaseTestSuite) TestFindTicketOfflineErrNil() {
	// Arrange
	payload := ticketRequest.TicketReq{
		CountryCode: "code",
		EventId:     "id",
	}

	mockOfflineTicket := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	suite.mockTicketRepositoryQuery.On("FindOfflineTicketByCountry", mock.Anything, payload).Return(mockChannel(getMockTicketSold()))
	suite.mockTicketRepositoryQuery.On("FindTicketByLowestPrice", mock.Anything, mock.Anything).Return(mockChannel(getMockTicketSold()))
	suite.mockTicketRepositoryQuery.On("FindOfflineTicketByCountryCode", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockOfflineTicket))
	suite.mockKafkaProducer.On("Publish", mock.Anything, mock.Anything, mock.Anything)
	suite.mockLogger.On("Error", mock.Anything, mock.Anything, mock.Anything)

	// Act
	_, err := suite.usecase.FindTickets(suite.ctx, payload)
	suite.T().Log(err)

	// Assert
	assert.Error(suite.T(), err)
}

func (suite *QueryUsecaseTestSuite) TestFindTicketOfflineErrParse() {
	// Arrange
	payload := ticketRequest.TicketReq{
		CountryCode: "code",
		EventId:     "id",
	}

	mockOfflineTicket := helpers.Result{
		Data: &ticketEntity.Country{
			Code: "code",
			Name: "name",
		},
		Error: nil,
	}

	suite.mockTicketRepositoryQuery.On("FindOfflineTicketByCountry", mock.Anything, payload).Return(mockChannel(getMockTicketSold()))
	suite.mockTicketRepositoryQuery.On("FindTicketByLowestPrice", mock.Anything, mock.Anything).Return(mockChannel(getMockTicketSold()))
	suite.mockTicketRepositoryQuery.On("FindOfflineTicketByCountryCode", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockOfflineTicket))
	suite.mockKafkaProducer.On("Publish", mock.Anything, mock.Anything, mock.Anything)
	suite.mockLogger.On("Error", mock.Anything, mock.Anything, mock.Anything)

	// Act
	_, err := suite.usecase.FindTickets(suite.ctx, payload)
	suite.T().Log(err)

	// Assert
	assert.Error(suite.T(), err)
}

// Helper function to create a channel
func mockChannel(result helpers.Result) <-chan helpers.Result {
	responseChan := make(chan helpers.Result)

	go func() {
		responseChan <- result
		close(responseChan)
	}()

	return responseChan
}

func getMockTicketSold() helpers.Result {
	return helpers.Result{
		Data: &[]ticketEntity.Ticket{
			{
				TicketId:       "id",
				EventId:        "id",
				TicketType:     "type",
				TicketPrice:    50,
				TotalRemaining: 0,
				Tag:            "tag",
			},
			{
				TicketId:       "id",
				EventId:        "id",
				TicketType:     "type",
				TicketPrice:    50,
				TotalRemaining: 0,
				Tag:            "tag",
			},
			{
				TicketId:       "id",
				EventId:        "id",
				TicketType:     "type",
				TicketPrice:    50,
				TotalRemaining: 0,
				Tag:            "tag",
			},
			{
				TicketId:       "id",
				EventId:        "id",
				TicketType:     "type",
				TicketPrice:    50,
				TotalRemaining: 0,
				Tag:            "tag",
			},
		},
		Error: nil,
	}
}

func (suite *QueryUsecaseTestSuite) TestFindOnlineTicketSuccess() {
	// Arrange
	payload := ticketRequest.TicketReq{
		CountryCode: "code",
		EventId:     "id",
	}

	mockOnlineTicket := helpers.Result{
		Data: &ticketEntity.Ticket{
			TicketId:       "id",
			EventId:        "id",
			TicketType:     "type",
			TicketPrice:    50,
			TotalRemaining: 0,
			Tag:            "tag",
		},
		Error: nil,
	}

	suite.mockTicketRepositoryQuery.On("FindOfflineTicketByCountry", mock.Anything, payload).Return(mockChannel(getMockTicketSold()))
	suite.mockTicketRepositoryQuery.On("FindOnlineTicketByCountry", mock.Anything, mock.Anything).Return(mockChannel(mockOnlineTicket))
	suite.mockLogger.On("Error", mock.Anything, mock.Anything, mock.Anything)

	// Act
	result, err := suite.usecase.FindOnlineTicket(suite.ctx, payload)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
}

func (suite *QueryUsecaseTestSuite) TestFindOnlineTicketErr() {
	// Arrange
	payload := ticketRequest.TicketReq{
		CountryCode: "code",
		EventId:     "id",
	}

	mockOfflineTicket := helpers.Result{
		Data:  nil,
		Error: errors.BadRequest("error"),
	}

	suite.mockTicketRepositoryQuery.On("FindOfflineTicketByCountry", mock.Anything, payload).Return(mockChannel(mockOfflineTicket))
	suite.mockLogger.On("Error", mock.Anything, mock.Anything, mock.Anything)

	// Act
	_, err := suite.usecase.FindOnlineTicket(suite.ctx, payload)

	// Assert
	assert.Error(suite.T(), err)

	mockOfflineTicket2 := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	suite.mockTicketRepositoryQuery.On("FindOfflineTicketByCountry", mock.Anything, payload).Return(mockChannel(mockOfflineTicket2))
	suite.mockLogger.On("Error", mock.Anything, mock.Anything, mock.Anything)

	// Act
	_, err2 := suite.usecase.FindOnlineTicket(suite.ctx, payload)

	// Assert
	assert.Error(suite.T(), err2)

	mockOfflineTicket3 := helpers.Result{
		Data: &ticketEntity.Country{
			Code: "code",
			Name: "Name",
		},
		Error: nil,
	}

	suite.mockTicketRepositoryQuery.On("FindOfflineTicketByCountry", mock.Anything, payload).Return(mockChannel(mockOfflineTicket3))
	suite.mockLogger.On("Error", mock.Anything, mock.Anything, mock.Anything)

	// Act
	_, err3 := suite.usecase.FindOnlineTicket(suite.ctx, payload)

	// Assert
	assert.Error(suite.T(), err3)
}

func (suite *QueryUsecaseTestSuite) TestFindOnlineTicketErrParse() {
	// Arrange
	payload := ticketRequest.TicketReq{
		CountryCode: "code",
		EventId:     "id",
	}

	mockOfflineTicket := helpers.Result{
		Data: &ticketEntity.Country{
			Code: "code",
			Name: "Name",
		},
		Error: nil,
	}

	suite.mockTicketRepositoryQuery.On("FindOfflineTicketByCountry", mock.Anything, payload).Return(mockChannel(mockOfflineTicket))
	suite.mockLogger.On("Error", mock.Anything, mock.Anything, mock.Anything)

	// Act
	_, err := suite.usecase.FindOnlineTicket(suite.ctx, payload)

	// Assert
	assert.Error(suite.T(), err)
}

func (suite *QueryUsecaseTestSuite) TestFindOnlineTicketNoEligible() {
	// Arrange
	payload := ticketRequest.TicketReq{
		CountryCode: "code",
		EventId:     "id",
	}

	mockOfflineTicket := helpers.Result{
		Data: &[]ticketEntity.Ticket{
			{
				TicketId:       "id",
				EventId:        "id",
				TicketType:     "type",
				TicketPrice:    50,
				TotalRemaining: 0,
				Tag:            "tag",
			},
		},
		Error: nil,
	}

	suite.mockTicketRepositoryQuery.On("FindOfflineTicketByCountry", mock.Anything, payload).Return(mockChannel(mockOfflineTicket))
	suite.mockLogger.On("Error", mock.Anything, mock.Anything, mock.Anything)

	// Act
	_, err := suite.usecase.FindOnlineTicket(suite.ctx, payload)

	// Assert
	assert.Error(suite.T(), err)
}

func (suite *QueryUsecaseTestSuite) TestFindOnlineTicketErrCountry() {
	// Arrange
	payload := ticketRequest.TicketReq{
		CountryCode: "code",
		EventId:     "id",
	}

	mockOnlineTicket := helpers.Result{
		Data:  nil,
		Error: errors.BadRequest("err"),
	}

	suite.mockTicketRepositoryQuery.On("FindOfflineTicketByCountry", mock.Anything, payload).Return(mockChannel(getMockTicketSold()))
	suite.mockTicketRepositoryQuery.On("FindOnlineTicketByCountry", mock.Anything, mock.Anything).Return(mockChannel(mockOnlineTicket))
	suite.mockLogger.On("Error", mock.Anything, mock.Anything, mock.Anything)

	// Act
	_, err := suite.usecase.FindOnlineTicket(suite.ctx, payload)

	// Assert
	assert.Error(suite.T(), err)

	mockOnlineTicket2 := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	suite.mockTicketRepositoryQuery.On("FindOfflineTicketByCountry", mock.Anything, payload).Return(mockChannel(getMockTicketSold()))
	suite.mockTicketRepositoryQuery.On("FindOnlineTicketByCountry", mock.Anything, mock.Anything).Return(mockChannel(mockOnlineTicket2))
	suite.mockLogger.On("Error", mock.Anything, mock.Anything, mock.Anything)

	// Act
	_, err2 := suite.usecase.FindOnlineTicket(suite.ctx, payload)

	// Assert
	assert.Error(suite.T(), err2)
}

func (suite *QueryUsecaseTestSuite) TestFindOnlineTicketErrCountryNil() {
	// Arrange
	payload := ticketRequest.TicketReq{
		CountryCode: "code",
		EventId:     "id",
	}

	mockOnlineTicket := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	suite.mockTicketRepositoryQuery.On("FindOfflineTicketByCountry", mock.Anything, payload).Return(mockChannel(getMockTicketSold()))
	suite.mockTicketRepositoryQuery.On("FindOnlineTicketByCountry", mock.Anything, mock.Anything).Return(mockChannel(mockOnlineTicket))
	suite.mockLogger.On("Error", mock.Anything, mock.Anything, mock.Anything)

	// Act
	_, err := suite.usecase.FindOnlineTicket(suite.ctx, payload)

	// Assert
	assert.Error(suite.T(), err)
}

func (suite *QueryUsecaseTestSuite) TestFindOnlineTicketErrCountryParse() {
	// Arrange
	payload := ticketRequest.TicketReq{
		CountryCode: "code",
		EventId:     "id",
	}

	mockOnlineTicket := helpers.Result{
		Data: &ticketEntity.Country{
			Code: "code",
		},
		Error: nil,
	}

	suite.mockTicketRepositoryQuery.On("FindOfflineTicketByCountry", mock.Anything, payload).Return(mockChannel(getMockTicketSold()))
	suite.mockTicketRepositoryQuery.On("FindOnlineTicketByCountry", mock.Anything, mock.Anything).Return(mockChannel(mockOnlineTicket))
	suite.mockLogger.On("Error", mock.Anything, mock.Anything, mock.Anything)

	// Act
	_, err := suite.usecase.FindOnlineTicket(suite.ctx, payload)

	// Assert
	assert.Error(suite.T(), err)
}
