package service

import (
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/gavinklfong/go-forex-trade-api/apiclient"
	"github.com/gavinklfong/go-forex-trade-api/dao"
	"github.com/gavinklfong/go-forex-trade-api/model"
	"github.com/xyproto/randomstring"
)

var CURRENCIES = [...]string{"CAD", "USD", "EUR", "ISK", "PHP",
	"DKK", "HUF", "CZK", "GBP", "RON", "SEK", "IDR", "INR", "BRL", "JPY"}

type ForexRateServiceImpl struct {
	forexRateDao    dao.ForexRateDao
	forexPricingDao dao.ForexPricingDao
	forexApiClient  apiclient.ForexApiClient
	timeProvider    TimeProvider
}

func NewForexRateService(forexApiClient apiclient.ForexApiClient, forexRateDao dao.ForexRateDao, forexPricingDao dao.ForexPricingDao, timeProvider TimeProvider) ForexRateService {
	return &ForexRateServiceImpl{forexRateDao, forexPricingDao, forexApiClient, timeProvider}
}

func (s *ForexRateServiceImpl) GetRateByCurrencyPair(baseCurrency, counterCurrency string) (*model.ForexRate, error) {
	// Validate currencies
	if err := s.validateCurrencyPair(baseCurrency, counterCurrency); err != nil {
		return nil, err
	}

	forexRate, err := s.forexApiClient.GetRateByCurrencyPair(baseCurrency, counterCurrency)
	if err != nil {
		slog.Error(fmt.Sprintf("forex api returned error: %v", err))
		return nil, err
	}

	rate, exist := forexRate.Rates[counterCurrency]
	if !exist {
		return nil, fmt.Errorf("Forex rate not found for %s/%s", baseCurrency, counterCurrency)
	}

	return s.buildForexRate(baseCurrency, counterCurrency, rate)
}

func (s *ForexRateServiceImpl) GetRatesByBaseCurrency(baseCurrency string) ([]*model.ForexRate, error) {
	// Validate base currency
	if err := s.validateCurrency(baseCurrency, "base"); err != nil {
		return nil, err
	}

	rateResp, err := s.forexApiClient.GetRateByBaseCurrency(baseCurrency)
	if err != nil {
		slog.Error(fmt.Sprintf("forex api returned error: %v", err))
		return nil, err
	}

	var forexRates []*model.ForexRate
	for counterCurrency, rate := range rateResp.Rates {
		forexRate, err := s.buildForexRate(baseCurrency, counterCurrency, rate)
		if err != nil {
			slog.Error(fmt.Sprintf("failed to build forex rate: %v", err))
			return nil, err
		}
		forexRates = append(forexRates, forexRate)
	}

	return forexRates, nil
}

func (s *ForexRateServiceImpl) validateTradeAction(action string) error {
	if action != "BUY" && action != "SELL" {
		return fmt.Errorf("invalid trade action: %s (must be BUY or SELL)", action)
	}
	return nil
}

func (s *ForexRateServiceImpl) BookRate(request *model.ForexRateBookingRequest) (*model.ForexRateBooking, error) {

	err := s.validateTradeAction(request.TradeAction)
	if err != nil {
		return nil, err
	}

	rate, err := s.GetRateByCurrencyPair(request.BaseCurrency, request.CounterCurrency)
	if err != nil {
		return nil, fmt.Errorf("failed retrieve rate by currency pair: %s/%s", request.BaseCurrency, request.CounterCurrency)
	}

	var bookingRate float32
	if request.TradeAction == "BUY" {
		bookingRate = rate.BuyRate
	} else {
		bookingRate = rate.SellRate
	}

	now := s.timeProvider.Now().UTC()

	booking := &model.ForexRateBooking{
		ForexRateBookingRequest: model.ForexRateBookingRequest{
			BaseCurrency:       request.BaseCurrency,
			CounterCurrency:    request.CounterCurrency,
			BaseCurrencyAmount: request.BaseCurrencyAmount,
			TradeAction:        request.TradeAction,
			CustomerId:         request.CustomerId,
		},
		Timestamp:  now,
		Rate:       bookingRate,
		BookingRef: randomstring.String(8),
		ExpiryTime: now.Add(time.Minute * 10),
	}

	_, err = s.forexRateDao.Insert(booking)
	if err != nil {
		return nil, fmt.Errorf("failed save rate booking: %v", booking)
	}

	return booking, nil
}

func (s *ForexRateServiceImpl) buildForexRate(baseCurrency, counterCurrency string, rate float32) (*model.ForexRate, error) {

	pricing := s.forexPricingDao.GetPricingByCurrencyPair(baseCurrency, counterCurrency)
	if pricing == nil {
		return nil, fmt.Errorf("pricing entry does not exist for %s/%s", baseCurrency, counterCurrency)
	}

	return &model.ForexRate{
		Timestamp:       s.timeProvider.Now().UTC(),
		BaseCurrency:    baseCurrency,
		CounterCurrency: counterCurrency,
		BuyRate:         rate + pricing.BuyPip/10000,
		SellRate:        rate + pricing.SellPip/10000,
		Spread:          pricing.GetSpread(),
	}, nil
}

func (s *ForexRateServiceImpl) validateCurrency(currency string, currencyType string) error {
	if !isValidCurrency(currency) {
		return fmt.Errorf("unsupported %s currency %s", currencyType, currency)
	}
	return nil
}

func (s *ForexRateServiceImpl) validateCurrencyPair(baseCurrency, counterCurrency string) error {
	if err := s.validateCurrency(baseCurrency, "base"); err != nil {
		return err
	}

	if err := s.validateCurrency(counterCurrency, "counter"); err != nil {
		return err
	}

	return nil
}

func isValidCurrency(currency string) bool {
	for _, item := range CURRENCIES {
		if strings.Compare(item, currency) == 0 {
			return true
		}
	}

	return false
}
