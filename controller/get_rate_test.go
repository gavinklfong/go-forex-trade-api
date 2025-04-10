package controller

import (
	"testing"

	"github.com/gavinklfong/go-forex-trade-api/model"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type forexRateServiceMock struct {
	mock.Mock
}

func newForexRateServiceMock() *forexRateServiceMock {
	return &forexRateServiceMock{}
}

func (m *forexRateServiceMock) GetRateByCurrencyPair(baseCurrency, counterCurrency string) (*model.ForexRate, error) {
	args := m.Called
}

type GetRateControllerTestSuite struct {
	suite.Suite
	router *gin.Engine
}

func TestGetRateControllerTestSuite(t *testing.T) {
	suite.Run(t, new(GetRateControllerTestSuite))
}

func (suite *GetRateControllerTestSuite) SetupSuite() {
	controller := NewGetController()

	router := gin.Default()
	router.GET("/rates/:baseCurrency/:counterCurrency", r.getRateController.GetRateByCurrencyPair)
}

func TestGetRateByBaseCurrency(t *testing.T) {

}
