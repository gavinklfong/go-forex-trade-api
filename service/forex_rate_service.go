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
	customerDao     dao.CustomerDao
	forexRateDao    dao.ForexRateDao
	forexPricingDao dao.ForexPricingDao
	forexApiClient  apiclient.ForexApiClient
}

func NewForexRateService(customerDao dao.CustomerDao, forexRateDao dao.ForexRateDao,
	forexApiClient apiclient.ForexApiClient, forexPricingDao dao.ForexPricingDao) ForexRateService {
	return &ForexRateServiceImpl{customerDao, forexRateDao, forexPricingDao, forexApiClient}
}

func (s *ForexRateServiceImpl) GetRateByCurrencyPair(baseCurrency, counterCurrency string) (*model.ForexRate, error) {
	forexRate, err := s.forexApiClient.GetRateByCurrencyPair(baseCurrency, counterCurrency)
	if err != nil {
		slog.Error(fmt.Sprintf("forex api returned error: %v", err))
		return nil, err
	}

	pricing := s.forexPricingDao.GetPricingByCurrencyPair(baseCurrency, counterCurrency)
	if pricing == nil {
		return nil, fmt.Errorf("pricing entry does not exist for %s/%s", baseCurrency, counterCurrency)
	}

	rate, exist := forexRate.Rates[counterCurrency]
	if !exist {
		return nil, fmt.Errorf("Forex rate not found for %s/%s", baseCurrency, counterCurrency)
	}

	return &model.ForexRate{
		Timestamp:       forexRate.Date,
		BaseCurrency:    baseCurrency,
		CounterCurrency: counterCurrency,
		BuyRate:         rate + pricing.BuyPip/10000,
		SellRate:        rate + pricing.SellPip/10000,
		Spread:          pricing.GetSpread(),
	}, nil
}

func (s *ForexRateServiceImpl) GetRatesByBaseCurrency(baseCurrency string) ([]*model.ForexRate, error) {
	forexRate, err := s.forexApiClient.GetRateByBaseCurrency(baseCurrency)
	if err != nil {
		slog.Error(fmt.Sprintf("forex api returned error: %v", err))
		return nil, err
	}

	_, v := range forexRate.Rates {
		
	}

	var result []*model.ForexRate

	for _, counterCurrency := range CURRENCIES {
		if strings.Compare(counterCurrency, baseCurrency) == 0 {
			continue
		}
		result = append(result, buildForexRate(baseCurrency, counterCurrency))
	}

	return result
}

func (s *ForexRateServiceImpl) BookRate(request *model.ForexRateBookingRequest) *model.ForexRateBooking {

	return &model.ForexRateBooking{
		ForexRateBookingRequest: model.ForexRateBookingRequest{
			BaseCurrency:       request.BaseCurrency,
			CounterCurrency:    request.CounterCurrency,
			BaseCurrencyAmount: request.BaseCurrencyAmount,
			TradeAction:        request.TradeAction,
			CustomerId:         request.CustomerId,
		},
		Timestamp:  time.Now(),
		Rate:       rand.Float32(),
		BookingRef: randomstring.String(8),
		ExpiryTime: time.Now().Add(time.Second * 30),
	}
}

func buildForexRate(baseCurrency, counterCurrency string) *model.ForexRate {
	return &model.ForexRate{
		Timestamp:       time.Now(),
		BaseCurrency:    baseCurrency,
		CounterCurrency: counterCurrency,
		BuyRate:         rand.Float32(),
		SellRate:        rand.Float32(),
		Spread:          rand.Float32(),
	}
}
