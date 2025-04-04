package apiclient

import "github.com/gavinklfong/go-forex-trade-api/apiclient/model"

type ForexApiClient interface {
	GetRateByBaseCurrency(base string) (*model.ForexRateResponse, error)
	GetRateByCurrencyPair(base, counter string) (*model.ForexRateResponse, error)
}
