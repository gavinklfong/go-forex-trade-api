package service

import (
	"testing"
	"time"

	apimodel "github.com/gavinklfong/go-forex-trade-api/apiclient/model"
	mockapiclient "github.com/gavinklfong/go-forex-trade-api/mocks/apiclient"
	mockservice "github.com/gavinklfong/go-forex-trade-api/mocks/service"
	"github.com/gavinklfong/go-forex-trade-api/model"
	"github.com/stretchr/testify/suite"

	mockdao "github.com/gavinklfong/go-forex-trade-api/mocks/dao"
)

type ForexRateServiceTestSuite struct {
	suite.Suite
	forexRateService    ForexRateService
	mockForexApiClient  *mockapiclient.MockForexApiClient
	mockForexPricingDao *mockdao.MockForexPricingDao
	mockTimeProvider    *mockservice.MockTimeProvider
}

func TestForexRateServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ForexRateServiceTestSuite))
}

func (s *ForexRateServiceTestSuite) SetupTest() {
	s.mockForexApiClient = mockapiclient.NewMockForexApiClient(s.T())
	s.mockForexPricingDao = mockdao.NewMockForexPricingDao(s.T())
	s.mockTimeProvider = mockservice.NewMockTimeProvider(s.T())

	s.forexRateService = NewForexRateService(s.mockForexApiClient,
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
