package service

import (
	"fmt"
	"testing"
	"time"

	mockdao "github.com/gavinklfong/go-forex-trade-api/mocks/dao"
	mockservice "github.com/gavinklfong/go-forex-trade-api/mocks/service"
	"github.com/gavinklfong/go-forex-trade-api/model"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type ForexTradeDealServiceTestSuite struct {
	suite.Suite
	forexTradeDealService ForexTradeDealService
	mockTradeDealDao      *mockdao.MockForexTradeDealDao
	mockRateDao           *mockdao.MockForexRateDao
	mockTimeProvider      *mockservice.MockTimeProvider
}

func TestForexTradeDealServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ForexTradeDealServiceTestSuite))
}

func (s *ForexTradeDealServiceTestSuite) SetupTest() {
	s.mockTradeDealDao = mockdao.NewMockForexTradeDealDao(s.T())
	s.mockRateDao = mockdao.NewMockForexRateDao(s.T())
	s.mockTimeProvider = mockservice.NewMockTimeProvider(s.T())

	s.forexTradeDealService = NewForexTradeDealService(
		s.mockTradeDealDao,
		s.mockRateDao,
		s.mockTimeProvider)
}

func (s *ForexTradeDealServiceTestSuite) TestSubmitTradeDealSuccess() {
	// Given
	baseCurrency := "GBP"
	counterCurrency := "USD"
	baseCurrencyAmount := float32(1000.00)
	rate := float32(1.2950)
	bookingRef := "ABC12345"
	customerID := "customer123"
	tradeAction := "BUY"
	
	now := time.Date(2025, 4, 1, 14, 30, 0, 0, time.UTC)
	expiryTime := now.Add(10 * time.Minute)
	
	// Mock booking found
	booking := &model.ForexRateBooking{
		ForexRateBookingRequest: model.ForexRateBookingRequest{
			BaseCurrency:       baseCurrency,
			CounterCurrency:    counterCurrency,
			BaseCurrencyAmount: baseCurrencyAmount,
			TradeAction:        tradeAction,
		},
		Timestamp:  now.Add(-5 * time.Minute),
		Rate:       rate,
		BookingRef: bookingRef,
		ExpiryTime: expiryTime,
		CustomerID: customerID,
	}
	s.mockRateDao.EXPECT().FindByBookingRef(bookingRef).Return(booking, nil).Once()
	
	// Mock time provider
	s.mockTimeProvider.EXPECT().Now().Return(now).Once()
	
	// Mock insert deal - capture the deal to verify it
	s.mockTradeDealDao.EXPECT().Insert(mock.AnythingOfType("*model.ForexTradeDeal")).
		Run(func(deal *model.ForexTradeDeal) {
			s.Equal(baseCurrency, deal.BaseCurrency)
			s.Equal(counterCurrency, deal.CounterCurrency)
			s.Equal(rate, deal.Rate)
			s.Equal(tradeAction, deal.TradeAction)
			s.Equal(fmt.Sprintf("%.2f", baseCurrencyAmount), deal.BaseCurrencyAmount)
			s.Equal(customerID, deal.CustomerID)
			s.Equal(now, deal.Timestamp)
			s.NotEmpty(deal.ID)
			s.NotEmpty(deal.Ref)
		}).Return(int64(1), nil).Once()

	// When
	deal, err := s.forexTradeDealService.SubmitTradeDeal(
		baseCurrency, counterCurrency, baseCurrencyAmount, rate, bookingRef)

	// Then
	s.Nil(err)
	s.NotNil(deal)
	s.Equal(baseCurrency, deal.BaseCurrency)
	s.Equal(counterCurrency, deal.CounterCurrency)
	s.Equal(rate, deal.Rate)
	s.Equal(tradeAction, deal.TradeAction)
	s.Equal(fmt.Sprintf("%.2f", baseCurrencyAmount), deal.BaseCurrencyAmount)
	s.Equal(customerID, deal.CustomerID)

	s.mockRateDao.AssertExpectations(s.T())
	s.mockTradeDealDao.AssertExpectations(s.T())
	s.mockTimeProvider.AssertExpectations(s.T())
}

func (s *ForexTradeDealServiceTestSuite) TestSubmitTradeDealWithNonExistentBooking() {
	// Given
	bookingRef := "NONEXIST"
	s.mockRateDao.EXPECT().FindByBookingRef(bookingRef).Return(nil, nil).Once()

	// When
	deal, err := s.forexTradeDealService.SubmitTradeDeal(
		"GBP", "USD", 1000.00, 1.2950, bookingRef)

	// Then
	s.Nil(deal)
	s.NotNil(err)
	s.Contains(err.Error(), "rate booking not found")
	
	s.mockRateDao.AssertExpectations(s.T())
}

func (s *ForexTradeDealServiceTestSuite) TestSubmitTradeDealWithExpiredBooking() {
	// Given
	bookingRef := "EXPIRED"
	baseCurrency := "GBP"
	counterCurrency := "USD"
	
	now := time.Date(2025, 4, 1, 14, 30, 0, 0, time.UTC)
	bookingTime := now.Add(-20 * time.Minute)
	expiryTime := bookingTime.Add(10 * time.Minute) // Expired 10 minutes ago
	
	// Mock booking found but expired
	booking := &model.ForexRateBooking{
		ForexRateBookingRequest: model.ForexRateBookingRequest{
			BaseCurrency:    baseCurrency,
			CounterCurrency: counterCurrency,
		},
		Timestamp:  bookingTime,
		ExpiryTime: expiryTime,
		BookingRef: bookingRef,
	}
	s.mockRateDao.EXPECT().FindByBookingRef(bookingRef).Return(booking, nil).Once()
	s.mockTimeProvider.EXPECT().Now().Return(now).Once()

	// When
	deal, err := s.forexTradeDealService.SubmitTradeDeal(
		baseCurrency, counterCurrency, 1000.00, 1.2950, bookingRef)

	// Then
	s.Nil(deal)
	s.NotNil(err)
	s.Contains(err.Error(), "rate booking has expired")
	
	s.mockRateDao.AssertExpectations(s.T())
	s.mockTimeProvider.AssertExpectations(s.T())
}

func (s *ForexTradeDealServiceTestSuite) TestSubmitTradeDealWithCurrencyMismatch() {
	// Given
	bookingRef := "MISMATCH"
	bookingBaseCurrency := "GBP"
	bookingCounterCurrency := "USD"
	
	requestBaseCurrency := "EUR" // Different from booking
	requestCounterCurrency := "USD"
	
	now := time.Date(2025, 4, 1, 14, 30, 0, 0, time.UTC)
	expiryTime := now.Add(5 * time.Minute) // Still valid
	
	// Mock booking found but with different currencies
	booking := &model.ForexRateBooking{
		ForexRateBookingRequest: model.ForexRateBookingRequest{
			BaseCurrency:    bookingBaseCurrency,
			CounterCurrency: bookingCounterCurrency,
		},
		Timestamp:  now.Add(-5 * time.Minute),
		ExpiryTime: expiryTime,
		BookingRef: bookingRef,
	}
	s.mockRateDao.EXPECT().FindByBookingRef(bookingRef).Return(booking, nil).Once()
	s.mockTimeProvider.EXPECT().Now().Return(now).Once()

	// When
	deal, err := s.forexTradeDealService.SubmitTradeDeal(
		requestBaseCurrency, requestCounterCurrency, 1000.00, 1.2950, bookingRef)

	// Then
	s.Nil(deal)
	s.NotNil(err)
	s.Contains(err.Error(), "currency pair mismatch")
	
	s.mockRateDao.AssertExpectations(s.T())
	s.mockTimeProvider.AssertExpectations(s.T())
}

func (s *ForexTradeDealServiceTestSuite) TestSubmitTradeDealWithRateMismatch() {
	// Given
	bookingRef := "RATEMISMATCH"
	baseCurrency := "GBP"
	counterCurrency := "USD"
	bookingRate := float32(1.2950)
	requestRate := float32(1.3000) // Different rate
	
	now := time.Date(2025, 4, 1, 14, 30, 0, 0, time.UTC)
	expiryTime := now.Add(5 * time.Minute) // Still valid
	
	// Mock booking found but with different rate
	booking := &model.ForexRateBooking{
		ForexRateBookingRequest: model.ForexRateBookingRequest{
			BaseCurrency:    baseCurrency,
			CounterCurrency: counterCurrency,
		},
		Timestamp:  now.Add(-5 * time.Minute),
		ExpiryTime: expiryTime,
		BookingRef: bookingRef,
		Rate:       bookingRate,
	}
	s.mockRateDao.EXPECT().FindByBookingRef(bookingRef).Return(booking, nil).Once()
	s.mockTimeProvider.EXPECT().Now().Return(now).Once()

	// When
	deal, err := s.forexTradeDealService.SubmitTradeDeal(
		baseCurrency, counterCurrency, 1000.00, requestRate, bookingRef)

	// Then
	s.Nil(deal)
	s.NotNil(err)
	s.Contains(err.Error(), "rate mismatch")
	
	s.mockRateDao.AssertExpectations(s.T())
	s.mockTimeProvider.AssertExpectations(s.T())
}

func (s *ForexTradeDealServiceTestSuite) TestSubmitTradeDealWithAmountMismatch() {
	// Given
	bookingRef := "AMTMISMATCH"
	baseCurrency := "GBP"
	counterCurrency := "USD"
	bookingAmount := float32(1000.00)
	requestAmount := float32(2000.00) // Different amount
	rate := float32(1.2950)
	
	now := time.Date(2025, 4, 1, 14, 30, 0, 0, time.UTC)
	expiryTime := now.Add(5 * time.Minute) // Still valid
	
	// Mock booking found but with different amount
	booking := &model.ForexRateBooking{
		ForexRateBookingRequest: model.ForexRateBookingRequest{
			BaseCurrency:       baseCurrency,
			CounterCurrency:    counterCurrency,
			BaseCurrencyAmount: bookingAmount,
		},
		Timestamp:  now.Add(-5 * time.Minute),
		ExpiryTime: expiryTime,
		BookingRef: bookingRef,
		Rate:       rate,
	}
	s.mockRateDao.EXPECT().FindByBookingRef(bookingRef).Return(booking, nil).Once()
	s.mockTimeProvider.EXPECT().Now().Return(now).Once()

	// When
	deal, err := s.forexTradeDealService.SubmitTradeDeal(
		baseCurrency, counterCurrency, requestAmount, rate, bookingRef)

	// Then
	s.Nil(deal)
	s.NotNil(err)
	s.Contains(err.Error(), "amount mismatch")
	
	s.mockRateDao.AssertExpectations(s.T())
	s.mockTimeProvider.AssertExpectations(s.T())
}

func (s *ForexTradeDealServiceTestSuite) TestSubmitTradeDealDaoError() {
	// Given
	bookingRef := "DAOFAIL"
	baseCurrency := "GBP"
	counterCurrency := "USD"
	amount := float32(1000.00)
	rate := float32(1.2950)
	customerID := "customer123"
	tradeAction := "BUY"
	
	now := time.Date(2025, 4, 1, 14, 30, 0, 0, time.UTC)
	expiryTime := now.Add(5 * time.Minute)
	
	// Mock valid booking
	booking := &model.ForexRateBooking{
		ForexRateBookingRequest: model.ForexRateBookingRequest{
			BaseCurrency:       baseCurrency,
			CounterCurrency:    counterCurrency,
			BaseCurrencyAmount: amount,
			TradeAction:        tradeAction,
		},
		Timestamp:  now.Add(-5 * time.Minute),
		ExpiryTime: expiryTime,
		BookingRef: bookingRef,
		Rate:       rate,
		CustomerID: customerID,
	}
	s.mockRateDao.EXPECT().FindByBookingRef(bookingRef).Return(booking, nil).Once()
	s.mockTimeProvider.EXPECT().Now().Return(now).Once()
	
	// Mock DAO insert error
	daoError := fmt.Errorf("database connection error")
	s.mockTradeDealDao.EXPECT().Insert(mock.AnythingOfType("*model.ForexTradeDeal")).Return(int64(0), daoError).Once()

	// When
	deal, err := s.forexTradeDealService.SubmitTradeDeal(
		baseCurrency, counterCurrency, amount, rate, bookingRef)

	// Then
	s.Nil(deal)
	s.NotNil(err)
	s.Contains(err.Error(), "failed to save trade deal")
	
	s.mockRateDao.AssertExpectations(s.T())
	s.mockTimeProvider.AssertExpectations(s.T())
	s.mockTradeDealDao.AssertExpectations(s.T())
}