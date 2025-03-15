package apiclient

import (
	"context"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/testcontainers/testcontainers-go"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/wiremock/go-wiremock"
	. "github.com/wiremock/wiremock-testcontainers-go"
)

type ForexApiClientTestSuite struct {
	suite.Suite
	apiClient         *ForexApiClient
	wiremockContainer *WireMockContainer
}

func TestForexApiClientTestSuite(t *testing.T) {
	suite.Run(t, new(ForexApiClientTestSuite))
}

func (suite *ForexApiClientTestSuite) TearDownSuite() {
	ctx := context.Background()
	if err := suite.wiremockContainer.Terminate(ctx); err != nil {
		suite.T().Fatalf("failed to terminate container: %s", err)
	}
}

func (suite *ForexApiClientTestSuite) TearDownTest() {
	// reset wiremock stub
	suite.wiremockContainer.Client.Reset()
}

func (suite *ForexApiClientTestSuite) SetupSuite() {
	log.Println("setting up test")
	ctx := context.Background()

	var err error

	// start up Wiremock testcontainer
	var emptyCustomizers []testcontainers.ContainerCustomizer
	suite.wiremockContainer, err = RunContainer(ctx, emptyCustomizers...)
	if err != nil {
		suite.T().Fatalf("Fail to start wiremock container: %s", err)
	}
	log.Println("Wiremock container started")

	// initialize forex API client
	uri, err := GetURI(ctx, suite.wiremockContainer)
	if err != nil {
		suite.T().Fatalf("Fail to obtain uri from wiremock container: %s", err)
	}

	suite.apiClient = NewForexApiClient(uri)
	log.Println("Forex API client initialized")
}

func (suite *ForexApiClientTestSuite) TestGetLatestRates() {

	// Use the WireMock client to stub a new endpoint manually

	body, err := os.ReadFile("rates.json")
	if err != nil {
		suite.T().Fatalf("fail to read stub response from file: %s", err)
	}

	err = suite.wiremockContainer.Client.StubFor(
		wiremock.Get(wiremock.URLEqualTo("/rates/latest")).
			WillReturnResponse(
				wiremock.NewResponse().
					WithBody(string(body)).
					WithHeader("Content-Type", "application/json").
					WithStatus(http.StatusOK),
			),
	)
	if err != nil {
		suite.T().Fatal(err)
	}

	rates, err := suite.apiClient.GetLatestRates()
	if err != nil {
		suite.T().Fatal(err)
	}

	assert.NotNil(suite.T(), rates)
}

func (suite *ForexApiClientTestSuite) TestGetLatestRate() {

	// Use the WireMock client to stub a new endpoint manually

	body, err := os.ReadFile("rate.json")
	if err != nil {
		suite.T().Fatalf("fail to read stub response from file: %s", err)
	}

	err = suite.wiremockContainer.Client.StubFor(
		wiremock.Get(wiremock.URLEqualTo("/rates/latest/AUD/USD")).
			WillReturnResponse(
				wiremock.NewResponse().
					WithBody(string(body)).
					WithHeader("Content-Type", "application/json").
					WithStatus(http.StatusOK),
			),
	)
	if err != nil {
		suite.T().Fatal(err)
	}

	rate, err := suite.apiClient.GetLatestRate("AUD", "USD")
	if err != nil {
		suite.T().Fatal(err)
	}

	assert.NotNil(suite.T(), rate)
}
