package dao

import "github.com/gavinklfong/go-forex-trade-api/model"

type CustomerDao interface {
	Insert(customer *model.Customer) (int64, error)
	FindByID(id string) (*model.Customer, error)
	FindByTier(tier int) (result []*model.Customer, err error)
}

type ForexRateDao interface {
	Insert(booking *model.ForexRateBooking) (int64, error)
	FindByID(id string) (*model.ForexRateBooking, error)
}

type ForexTradeDealDao interface {
	Insert(deal *model.ForexTradeDeal) (int64, error)
	FindByID(id string) (*model.ForexTradeDeal, error)
}

type ForexPricingDao interface {
	GetPricingByCurrencyPair(base, counter string) *model.ForexPricing
}
