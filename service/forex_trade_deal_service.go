package service

import (
	"fmt"
	"log/slog"

	"github.com/gavinklfong/go-forex-trade-api/dao"
	"github.com/gavinklfong/go-forex-trade-api/model"
	"github.com/google/uuid"
	"github.com/xyproto/randomstring"
)

type ForexTradeDealServiceImpl struct {
	tradeDealdao dao.ForexTradeDealDao
	rateDao      dao.ForexRateDao
	timeProvider TimeProvider
}

func NewForexTradeDealService(tradeDealdao dao.ForexTradeDealDao, rateDao dao.ForexRateDao, timeProvider TimeProvider) ForexTradeDealService {
	return &ForexTradeDealServiceImpl{tradeDealdao, rateDao, timeProvider}
}

func (s *ForexTradeDealServiceImpl) SubmitTradeDeal(baseCurrency, counterCurrency string,
	baseCurrencyAmount, rate float32, rateBookingRef string) (*model.ForexTradeDeal, error) {

	// Validate booking reference exists
	booking, err := s.rateDao.FindByBookingRef(rateBookingRef)
	if err != nil {
		slog.Error(fmt.Sprintf("error retrieving rate booking: %v", err))
		return nil, fmt.Errorf("failed to retrieve rate booking: %v", err)
	}

	if booking == nil {
		return nil, fmt.Errorf("rate booking not found with reference: %s", rateBookingRef)
	}

	// Validate booking has not expired
	currentTime := s.timeProvider.Now().UTC()
	if currentTime.After(booking.ExpiryTime) {
		return nil, fmt.Errorf("rate booking has expired at %v", booking.ExpiryTime)
	}

	// Validate booking details match the trade request
	if booking.BaseCurrency != baseCurrency || booking.CounterCurrency != counterCurrency {
		return nil, fmt.Errorf("currency pair mismatch: booking is for %s/%s, but trade is for %s/%s",
			booking.BaseCurrency, booking.CounterCurrency, baseCurrency, counterCurrency)
	}

	if booking.Rate != rate {
		return nil, fmt.Errorf("rate mismatch: booking rate is %f, but trade rate is %f", booking.Rate, rate)
	}

	if booking.BaseCurrencyAmount != baseCurrencyAmount {
		return nil, fmt.Errorf("amount mismatch: booking amount is %f, but trade amount is %f",
			booking.BaseCurrencyAmount, baseCurrencyAmount)
	}

	// Create forex trade deal
	deal := &model.ForexTradeDeal{
		ID:                 uuid.New().String(),
		Ref:                randomstring.String(8),
		Timestamp:          currentTime,
		BaseCurrency:       baseCurrency,
		CounterCurrency:    counterCurrency,
		Rate:               rate,
		TradeAction:        booking.TradeAction,
		BaseCurrencyAmount: fmt.Sprintf("%.2f", baseCurrencyAmount),
		CustomerID:         booking.CustomerID,
	}

	// Save the deal
	_, err = s.tradeDealdao.Insert(deal)
	if err != nil {
		slog.Error(fmt.Sprintf("failed to insert trade deal: %v", err))
		return nil, fmt.Errorf("failed to save trade deal: %v", err)
	}

	return deal, nil
}
