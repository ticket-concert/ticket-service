package queries_test

import (
	"context"
	"testing"
	"ticket-service/internal/modules/ticket"
	"ticket-service/internal/modules/ticket/models/request"
	mongoRQ "ticket-service/internal/modules/ticket/repositories/queries"
	"ticket-service/internal/pkg/helpers"
	mocks "ticket-service/mocks/pkg/databases/mongodb"
	mocklog "ticket-service/mocks/pkg/log"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type CommandTestSuite struct {
	suite.Suite
	mockMongodb *mocks.Collections
	mockLogger  *mocklog.Logger
	repository  ticket.MongodbRepositoryQuery
	ctx         context.Context
}

func (suite *CommandTestSuite) SetupTest() {
	suite.mockMongodb = new(mocks.Collections)
	suite.mockLogger = &mocklog.Logger{}
	suite.repository = mongoRQ.NewQueryMongodbRepository(
		suite.mockMongodb,
		suite.mockLogger,
	)
	suite.ctx = context.Background()
}

func TestCommandTestSuite(t *testing.T) {
	suite.Run(t, new(CommandTestSuite))
}

func (suite *CommandTestSuite) TestFindOfflineTicketByCountry() {

	request := request.TicketReq{}
	// Mock FindOne
	expectedResult := make(chan helpers.Result)
	suite.mockMongodb.On("FindMany", mock.Anything, mock.Anything).Return((<-chan helpers.Result)(expectedResult))

	// Act
	result := suite.repository.FindOfflineTicketByCountry(suite.ctx, request)
	// Asset
	assert.NotNil(suite.T(), result, "Expected a result")

	// Simulate receiving a result from the channel
	go func() {
		expectedResult <- helpers.Result{Data: "result not nil", Error: nil}
		close(expectedResult)
	}()

	// Wait for the goroutine to complete
	<-result

	// Assert FindOne
	suite.mockMongodb.AssertCalled(suite.T(), "FindMany", mock.Anything, mock.Anything)
}

func (suite *CommandTestSuite) TestFindOfflineTicketByCountryCode() {
	// Mock FindOne
	expectedResult := make(chan helpers.Result)
	suite.mockMongodb.On("FindMany", mock.Anything, mock.Anything).Return((<-chan helpers.Result)(expectedResult))

	// Act
	result := suite.repository.FindOfflineTicketByCountryCode(suite.ctx, mock.Anything, mock.Anything)
	// Asset
	assert.NotNil(suite.T(), result, "Expected a result")

	// Simulate receiving a result from the channel
	go func() {
		expectedResult <- helpers.Result{Data: "result not nil", Error: nil}
		close(expectedResult)
	}()

	// Wait for the goroutine to complete
	<-result

	// Assert FindOne
	suite.mockMongodb.AssertCalled(suite.T(), "FindMany", mock.Anything, mock.Anything)
}

func (suite *CommandTestSuite) TestFindTicketByLowestPrice() {
	// Mock FindOne
	expectedResult := make(chan helpers.Result)
	suite.mockMongodb.On("FindAllData", mock.Anything, mock.Anything).Return((<-chan helpers.Result)(expectedResult))

	// Act
	result := suite.repository.FindTicketByLowestPrice(suite.ctx, mock.Anything)
	// Asset
	assert.NotNil(suite.T(), result, "Expected a result")

	// Simulate receiving a result from the channel
	go func() {
		expectedResult <- helpers.Result{Data: "result not nil", Error: nil}
		close(expectedResult)
	}()

	// Wait for the goroutine to complete
	<-result

	// Assert FindOne
	suite.mockMongodb.AssertCalled(suite.T(), "FindAllData", mock.Anything, mock.Anything)
}

func (suite *CommandTestSuite) TestFindOnlineTicketByCountry() {
	// Mock FindOne
	expectedResult := make(chan helpers.Result)
	suite.mockMongodb.On("FindOne", mock.Anything, mock.Anything).Return((<-chan helpers.Result)(expectedResult))

	request := request.TicketReq{}
	// Act
	result := suite.repository.FindOnlineTicketByCountry(suite.ctx, request)
	// Asset
	assert.NotNil(suite.T(), result, "Expected a result")

	// Simulate receiving a result from the channel
	go func() {
		expectedResult <- helpers.Result{Data: "result not nil", Error: nil}
		close(expectedResult)
	}()

	// Wait for the goroutine to complete
	<-result

	// Assert FindOne
	suite.mockMongodb.AssertCalled(suite.T(), "FindOne", mock.Anything, mock.Anything)
}
