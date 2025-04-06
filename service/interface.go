package service

import "github.com/gavinklfong/go-forex-trade-api/model"

type ForexRateService interface {
	GetRateByCurrencyPair(baseCurrency, counterCurrency string) (*model.ForexRate, error)
	GetRatesByBaseCurrency(baseCurrency string) ([]*model.ForexRate, error)
	BookRate(request *model.ForexRateBookingRequest) *model.ForexRateBooking
}

type ForexTradeDealService interface {
}
