package service

import (
	"math/rand/v2"
	"strings"
	"time"

	"github.com/gavinklfong/go-rest-api-demo/apiclient"
	"github.com/gavinklfong/go-rest-api-demo/dao"
	"github.com/gavinklfong/go-rest-api-demo/model"
	"github.com/xyproto/randomstring"
)

var CURRENCIES = [...]string{"CAD", "USD", "EUR", "ISK", "PHP",
	"DKK", "HUF", "CZK", "GBP", "RON", "SEK", "IDR", "INR", "BRL", "JPY"}

type ForexRateService struct {
	customerDao    *dao.CustomerDao
	forexRateDao   *dao.ForexRateDao
	forexApiClient *apiclient.ForexApiClient
}

func NewForexRateService(customerDao *dao.CustomerDao, forexRateDao *dao.ForexRateDao,
	forexApiClient *apiclient.ForexApiClient) *ForexRateService {
	return &ForexRateService{customerDao, forexRateDao, forexApiClient}
}

func (s *ForexRateService) GetRateByCurrencyPair(baseCurrency, counterCurrency string) *model.ForexRate {
	return buildForexRate(baseCurrency, counterCurrency)
}

func (s *ForexRateService) GetRatesByBaseCurrency(baseCurrency string) []*model.ForexRate {

	var result []*model.ForexRate

	for _, counterCurrency := range CURRENCIES {
		if strings.Compare(counterCurrency, baseCurrency) == 0 {
			continue
		}
		result = append(result, buildForexRate(baseCurrency, counterCurrency))
	}

	return result
}

func (s *ForexRateService) BookRate(request *model.ForexRateBookingRequest) *model.ForexRateBooking {
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
