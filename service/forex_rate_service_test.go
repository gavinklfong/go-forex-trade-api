package service

import (
	"fmt"
	"testing"
	"time"

	apimodel "github.com/gavinklfong/go-forex-trade-api/apiclient/model"
	mockapiclient "github.com/gavinklfong/go-forex-trade-api/mocks/apiclient"
	mockservice "github.com/gavinklfong/go-forex-trade-api/mocks/service"
	"github.com/gavinklfong/go-forex-trade-api/model"
	"github.com/stretchr/testify/suite"
	"github.com/xyproto/randomstring"

	mockdao "github.com/gavinklfong/go-forex-trade-api/mocks/dao"
)

type ForexRateServiceTestSuite struct {
	suite.Suite
	forexRateService    ForexRateService
	mockForexApiClient  *mockapiclient.MockForexApiClient
	mockForexPricingDao *mockdao.MockForexPricingDao
	mockForexRateDao    *mockdao.MockForexRateDao
	mockTimeProvider    *mockservice.MockTimeProvider
}

func TestForexRateServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ForexRateServiceTestSuite))
}

func (s *ForexRateServiceTestSuite) SetupTest() {
	s.mockForexApiClient = mockapiclient.NewMockForexApiClient(s.T())
	s.mockForexPricingDao = mockdao.NewMockForexPricingDao(s.T())
	s.mockForexRateDao = mockdao.NewMockForexRateDao(s.T())
	s.mockTimeProvider = mockservice.NewMockTimeProvider(s.T())

	s.forexRateService = NewForexRateService(s.mockForexApiClient,
		s.mockForexRateDao,
		s.mockForexPricingDao,
		s.mockTimeProvider)
}

func (s *ForexRateServiceTestSuite) TestGetRateByCurrencyPair() {

	// Given
	rate := float32(1.29)
	buyPip := float32(3)
	sellPip := float32(-6)
	now := time.Date(2025, 4, 1, 14, 30, 0, 0, time.UTC)

	// mock forex api response
	rates := make(map[string]float32)
	rates["USD"] = rate
	mockApiResp := apimodel.ForexRateResponse{"GBP-USD", time.Now().Truncate(24 * time.Hour), "GBP", rates}
	s.mockForexApiClient.EXPECT().GetRateByCurrencyPair("GBP", "USD").Return(&mockApiResp, nil).Once()

	// mock forex pricing
	mockForexPricing := model.ForexPricing{"GBP", "USD", buyPip, sellPip}
	s.mockForexPricingDao.EXPECT().GetPricingByCurrencyPair("GBP", "USD").Return(&mockForexPricing).Once()

	// mock current time
	s.mockTimeProvider.EXPECT().Now().Return(now)

	// When
	forexRate, err := s.forexRateService.GetRateByCurrencyPair("GBP", "USD")

	// Then
	expectedForexRate := &model.ForexRate{
		Timestamp:       now,
		BaseCurrency:    "GBP",
		CounterCurrency: "USD",
		BuyRate:         rate + buyPip/10000,
		SellRate:        rate + sellPip/10000,
		Spread:          buyPip - sellPip,
	}

	s.Nil(err)
	s.Equal(expectedForexRate, forexRate)

	s.mockForexApiClient.AssertExpectations(s.T())
	s.mockForexPricingDao.AssertExpectations(s.T())
}

func (s *ForexRateServiceTestSuite) TestGetRateByBaseCurrency() {

	// Given
	usdRate := float32(1.29)
	usdBuyPip := float32(3)
	usdSellPip := float32(-6)

	audRate := float32(2.12)
	audBuyPip := float32(9)
	audSellPip := float32(7)

	now := time.Date(2025, 4, 1, 14, 30, 0, 0, time.UTC)

	// mock forex api response
	rates := make(map[string]float32)
	rates["USD"] = usdRate
	rates["AUD"] = audRate
	mockApiResp := apimodel.ForexRateResponse{"GBP", time.Now().Truncate(24 * time.Hour), "GBP", rates}
	s.mockForexApiClient.EXPECT().GetRateByBaseCurrency("GBP").Return(&mockApiResp, nil).Once()

	// mock forex pricing
	mockUSDForexPricing := model.ForexPricing{"GBP", "USD", usdBuyPip, usdSellPip}
	mockAUDForexPricing := model.ForexPricing{"GBP", "AUD", audBuyPip, audSellPip}
	s.mockForexPricingDao.EXPECT().GetPricingByCurrencyPair("GBP", "USD").Return(&mockUSDForexPricing).Once()
	s.mockForexPricingDao.EXPECT().GetPricingByCurrencyPair("GBP", "AUD").Return(&mockAUDForexPricing).Once()

	// mock current time
	s.mockTimeProvider.EXPECT().Now().Return(now)

	// When
	forexRates, err := s.forexRateService.GetRatesByBaseCurrency("GBP")

	// Then
	expectedUSDForexRate := &model.ForexRate{
		Timestamp:       now,
		BaseCurrency:    "GBP",
		CounterCurrency: "USD",
		BuyRate:         usdRate + usdBuyPip/10000,
		SellRate:        usdRate + usdSellPip/10000,
		Spread:          usdBuyPip - usdSellPip,
	}

	expectedAUDForexRate := &model.ForexRate{
		Timestamp:       now,
		BaseCurrency:    "GBP",
		CounterCurrency: "AUD",
		BuyRate:         audRate + audBuyPip/10000,
		SellRate:        audRate + audSellPip/10000,
		Spread:          audBuyPip - audSellPip,
	}

	s.Nil(err)
	s.Equal(2, len(forexRates))
	s.Contains(forexRates, expectedUSDForexRate)
	s.Contains(forexRates, expectedAUDForexRate)

	s.mockForexApiClient.AssertExpectations(s.T())
	s.mockForexPricingDao.AssertExpectations(s.T())
}

func (s *ForexRateServiceTestSuite) TestGetRateByCurrencyPairWithUnknownBaseCurrency() {

	forexRates, err := s.forexRateService.GetRateByCurrencyPair("ZZZ", "USD")
	s.Nil(forexRates)
	s.NotNil(err)
}

func (s *ForexRateServiceTestSuite) TestGetRateByCurrencyPairWithUnknownCounterCurrency() {

	forexRates, err := s.forexRateService.GetRateByCurrencyPair("USD", "ZZZ")
	s.Nil(forexRates)
	s.NotNil(err)
}

func (s *ForexRateServiceTestSuite) TestGetRateByBaseCurrencyWithUnknownCurrency() {

	forexRates, err := s.forexRateService.GetRatesByBaseCurrency("ZZZ")
	s.Nil(forexRates)
	s.NotNil(err)
}

func (s *ForexRateServiceTestSuite) TestBookRateWithBuyAction() {
	// Given
	baseCurrency := "GBP"
	counterCurrency := "USD"
	rate := float32(1.29)
	buyPip := float32(3)
	sellPip := float32(-6)
	now := time.Date(2025, 4, 1, 14, 30, 0, 0, time.UTC)
	expiryTime := now.Add(10 * time.Minute)

	// Create booking request
	request := &model.ForexRateBookingRequest{
		BaseCurrency:       baseCurrency,
		CounterCurrency:    counterCurrency,
		BaseCurrencyAmount: 1000,
		TradeAction:        "BUY",
		CustomerId:         123,
	}

	// Mock forex API response
	rates := make(map[string]float32)
	rates[counterCurrency] = rate
	mockApiResp := apimodel.ForexRateResponse{baseCurrency + "-" + counterCurrency, now, baseCurrency, rates}
	s.mockForexApiClient.EXPECT().GetRateByCurrencyPair(baseCurrency, counterCurrency).Return(&mockApiResp, nil).Once()

	// Mock forex pricing
	mockForexPricing := model.ForexPricing{baseCurrency, counterCurrency, buyPip, sellPip}
	s.mockForexPricingDao.EXPECT().GetPricingByCurrencyPair(baseCurrency, counterCurrency).Return(&mockForexPricing).Once()

	// Mock time provider
	s.mockTimeProvider.EXPECT().Now().Return(now).Once()

	// When
	booking, err := s.forexRateService.BookRate(request)

	// Then
	s.Nil(err)
	s.NotNil(booking)
	s.Equal(baseCurrency, booking.BaseCurrency)
	s.Equal(counterCurrency, booking.CounterCurrency)
	s.Equal(float32(1000), booking.BaseCurrencyAmount)
	s.Equal("BUY", booking.TradeAction)
	s.Equal(int32(123), booking.CustomerId)
	s.Equal(now, booking.Timestamp)
	s.Equal(rate+buyPip/10000, booking.Rate)
	s.Equal(expiryTime, booking.ExpiryTime)
	s.NotEmpty(booking.BookingRef)
	s.Len(booking.BookingRef, 8)

	s.mockForexApiClient.AssertExpectations(s.T())
	s.mockForexPricingDao.AssertExpectations(s.T())
	s.mockTimeProvider.AssertExpectations(s.T())
}

func (s *ForexRateServiceTestSuite) TestBookRateWithSellAction() {
	// Given
	baseCurrency := "GBP"
	counterCurrency := "USD"
	rate := float32(1.29)
	buyPip := float32(3)
	sellPip := float32(-6)
	now := time.Date(2025, 4, 1, 14, 30, 0, 0, time.UTC)

	// Create booking request
	request := &model.ForexRateBookingRequest{
		BaseCurrency:       baseCurrency,
		CounterCurrency:    counterCurrency,
		BaseCurrencyAmount: 1000,
		TradeAction:        "SELL",
		CustomerId:         123,
	}

	// Create expected booking
	booking := &model.ForexRateBooking{
		ForexRateBookingRequest: model.ForexRateBookingRequest{
			BaseCurrency:       request.BaseCurrency,
			CounterCurrency:    request.CounterCurrency,
			BaseCurrencyAmount: request.BaseCurrencyAmount,
			TradeAction:        request.TradeAction,
			CustomerId:         request.CustomerId,
		},
		Timestamp:  now,
		Rate:       rate + sellPip/10000,
		BookingRef: randomstring.String(8),
		ExpiryTime: now.Add(time.Minute * 10),
	}

	// Mock forex API response
	rates := make(map[string]float32)
	rates[counterCurrency] = rate
	mockApiResp := apimodel.ForexRateResponse{
		ID:    baseCurrency + "-" + counterCurrency,
		Date:  now,
		Base:  baseCurrency,
		Rates: rates}
	s.mockForexApiClient.EXPECT().GetRateByCurrencyPair(baseCurrency, counterCurrency).Return(&mockApiResp, nil).Once()

	// Mock forex pricing
	mockForexPricing := model.ForexPricing{
		BaseCurrency:    baseCurrency,
		CounterCurrency: counterCurrency,
		BuyPip:          buyPip,
		SellPip:         sellPip}
	s.mockForexPricingDao.EXPECT().GetPricingByCurrencyPair(baseCurrency, counterCurrency).Return(&mockForexPricing).Once()

	// Mock forex rate service
	s.mockForexRateDao.EXPECT().Insert().Return(1, nil)

	// Mock time provider
	s.mockTimeProvider.EXPECT().Now().Return(now).Twice()

	// When
	booking, err := s.forexRateService.BookRate(request)

	// Then
	s.Nil(err)
	s.NotNil(booking)
	s.Equal("SELL", booking.TradeAction)
	s.Equal(rate+sellPip/10000, booking.Rate)

	s.mockForexApiClient.AssertExpectations(s.T())
	s.mockForexRateDao.AssertExpectations(s.T())
	s.mockForexPricingDao.AssertExpectations(s.T())
	s.mockTimeProvider.AssertExpectations(s.T())
}

func (s *ForexRateServiceTestSuite) TestBookRateWithInvalidTradeAction() {
	// Given
	request := &model.ForexRateBookingRequest{
		BaseCurrency:       "GBP",
		CounterCurrency:    "USD",
		BaseCurrencyAmount: 1000,
		TradeAction:        "INVALID",
		CustomerId:         123,
	}

	// When
	booking, err := s.forexRateService.BookRate(request)

	// Then
	s.Nil(booking)
	s.NotNil(err)
	s.Contains(err.Error(), "invalid trade action")
}

func (s *ForexRateServiceTestSuite) TestBookRateWithInvalidBaseCurrency() {
	// Given
	request := &model.ForexRateBookingRequest{
		BaseCurrency:       "XXX",
		CounterCurrency:    "USD",
		BaseCurrencyAmount: 1000,
		TradeAction:        "BUY",
		CustomerId:         123,
	}

	// When
	booking, err := s.forexRateService.BookRate(request)

	// Then
	s.Nil(booking)
	s.NotNil(err)
	s.Contains(err.Error(), "unsupported base currency")
}

func (s *ForexRateServiceTestSuite) TestBookRateWithInvalidCounterCurrency() {
	// Given
	request := &model.ForexRateBookingRequest{
		BaseCurrency:       "GBP",
		CounterCurrency:    "XXX",
		BaseCurrencyAmount: 1000,
		TradeAction:        "BUY",
		CustomerId:         123,
	}

	// When
	booking, err := s.forexRateService.BookRate(request)

	// Then
	s.Nil(booking)
	s.NotNil(err)
	s.Contains(err.Error(), "unsupported counter currency")
}

func (s *ForexRateServiceTestSuite) TestBookRateWithMissingPricing() {
	// Given
	baseCurrency := "GBP"
	counterCurrency := "USD"
	rate := float32(1.29)
	now := time.Date(2025, 4, 1, 14, 30, 0, 0, time.UTC)

	// Create booking request
	request := &model.ForexRateBookingRequest{
		BaseCurrency:       baseCurrency,
		CounterCurrency:    counterCurrency,
		BaseCurrencyAmount: 1000,
		TradeAction:        "BUY",
		CustomerId:         123,
	}

	// Mock forex API response
	rates := make(map[string]float32)
	rates[counterCurrency] = rate
	mockApiResp := apimodel.ForexRateResponse{baseCurrency + "-" + counterCurrency, now, baseCurrency, rates}
	s.mockForexApiClient.EXPECT().GetRateByCurrencyPair(baseCurrency, counterCurrency).Return(&mockApiResp, nil).Once()

	// Mock missing pricing
	s.mockForexPricingDao.EXPECT().GetPricingByCurrencyPair(baseCurrency, counterCurrency).Return(nil).Once()

	// When
	booking, err := s.forexRateService.BookRate(request)

	// Then
	s.Nil(booking)
	s.NotNil(err)
	s.Contains(err.Error(), "pricing entry does not exist")

	s.mockForexApiClient.AssertExpectations(s.T())
	s.mockForexPricingDao.AssertExpectations(s.T())
}

func (s *ForexRateServiceTestSuite) TestBookRateWithApiError() {
	// Given
	baseCurrency := "GBP"
	counterCurrency := "USD"
	apiError := fmt.Errorf("API connection error")

	// Create booking request
	request := &model.ForexRateBookingRequest{
		BaseCurrency:       baseCurrency,
		CounterCurrency:    counterCurrency,
		BaseCurrencyAmount: 1000,
		TradeAction:        "BUY",
		CustomerId:         123,
	}

	// Mock forex API error
	s.mockForexApiClient.EXPECT().GetRateByCurrencyPair(baseCurrency, counterCurrency).Return(nil, apiError).Once()

	// When
	booking, err := s.forexRateService.BookRate(request)

	// Then
	s.Nil(booking)
	s.NotNil(err)
	s.Equal(apiError, err)

	s.mockForexApiClient.AssertExpectations(s.T())
}
