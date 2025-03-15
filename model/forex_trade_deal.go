package model

import "time"

type ForexTradeDeal struct {
	ID                 string
	Ref                string
	Timestamp          time.Time
	BaseCurrency       string
	CounterCurrency    string
	Rate               float32
	TradeAction        string
	BaseCurrencyAmount string
	CustomerId         string
}
