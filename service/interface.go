package service

import (
	"time"

	"github.com/gavinklfong/go-forex-trade-api/model"
)

type TimeProvider interface {
	Now() time.Time
}

type ForexRateService interface {
	GetRateByCurrencyPair(baseCurrency, counterCurrency string) (*model.ForexRate, error)
	GetRatesByBaseCurrency(baseCurrency string) ([]*model.ForexRate, error)
	BookRate(request *model.ForexRateBookingRequest) *model.ForexRateBooking
}

type ForexTradeDealService interface {
}
