package service

import (
	"fmt"
	"log/slog"
	"math/rand/v2"
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
	forexPricingDao dao.ForexPricingDao
	forexApiClient  apiclient.ForexApiClient
	timeProvider    TimeProvider
}

func NewForexRateService(forexApiClient apiclient.ForexApiClient, forexPricingDao dao.ForexPricingDao, timeProvider TimeProvider) ForexRateService {
	return &ForexRateServiceImpl{forexPricingDao, forexApiClient, timeProvider}
}

func (s *ForexRateServiceImpl) GetRateByCurrencyPair(baseCurrency, counterCurrency string) (*model.ForexRate, error) {

	if !isValidCurrency(baseCurrency) {
		return nil, fmt.Errorf("unsupported base currency %s", baseCurrency)
	}

	if !isValidCurrency(counterCurrency) {
		return nil, fmt.Errorf("unsupported counter currency %s", counterCurrency)
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

	if !isValidCurrency(baseCurrency) {
		return nil, fmt.Errorf("unsupported base currency %s", baseCurrency)
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

func (s *ForexRateServiceImpl) BookRate(request *model.ForexRateBookingRequest) (*model.ForexRateBooking, error) {
	// Validate currencies
	if !isValidCurrency(request.BaseCurrency) {
		return nil, fmt.Errorf("unsupported base currency %s", request.BaseCurrency)
	}

	if !isValidCurrency(request.CounterCurrency) {
		return nil, fmt.Errorf("unsupported counter currency %s", request.CounterCurrency)
	}

	// Get current rate from API
	forexRateResponse, err := s.forexApiClient.GetRateByCurrencyPair(request.BaseCurrency, request.CounterCurrency)
	if err != nil {
		slog.Error(fmt.Sprintf("forex api returned error: %v", err))
		return nil, err
	}

	rate, exist := forexRateResponse.Rates[request.CounterCurrency]
	if !exist {
		return nil, fmt.Errorf("forex rate not found for %s/%s", request.BaseCurrency, request.CounterCurrency)
	}

	// Get pricing to apply
	pricing := s.forexPricingDao.GetPricingByCurrencyPair(request.BaseCurrency, request.CounterCurrency)
	if pricing == nil {
		return nil, fmt.Errorf("pricing entry does not exist for %s/%s", request.BaseCurrency, request.CounterCurrency)
	}

	// Calculate the final rate based on trade action
	finalRate := rate
	if request.TradeAction == "BUY" {
		finalRate = rate + pricing.BuyPip/10000
	} else if request.TradeAction == "SELL" {
		finalRate = rate + pricing.SellPip/10000
	} else {
		return nil, fmt.Errorf("invalid trade action: %s", request.TradeAction)
	}

	now := s.timeProvider.Now().UTC()
	
	return &model.ForexRateBooking{
		ForexRateBookingRequest: model.ForexRateBookingRequest{
			BaseCurrency:       request.BaseCurrency,
			CounterCurrency:    request.CounterCurrency,
			BaseCurrencyAmount: request.BaseCurrencyAmount,
			TradeAction:        request.TradeAction,
			CustomerId:         request.CustomerId,
		},
		Timestamp:  now,
		Rate:       finalRate,
		BookingRef: randomstring.String(8),
		ExpiryTime: now.Add(time.Minute * 10),
	}, nil
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

func isValidCurrency(currency string) bool {
	for _, item := range CURRENCIES {
		if strings.Compare(item, currency) == 0 {
			return true
		}
	}

	return false
}
