package model

import (
	"encoding/json"
	"time"
)

type ForexRateBookingRequest struct {
	BaseCurrency       string
	CounterCurrency    string
	BaseCurrencyAmount float32
	TradeAction        string
	CustomerId         int32
}

type ForexRateBookingResponse struct {
	*ForexRateBookingRequest
	Timestamp  time.Time
	Rate       float32
	BookingRef string
	ExpiryTime time.Time
}

func (f *ForexRateBookingResponse) UnmarshalJSON(data []byte) error {
	type Alias ForexRateBookingResponse
	aux := &struct {
		Timestamp  string
		ExpiryTime string
		*Alias
	}{
		Alias: (*Alias)(f),
	}

	var err error

	if err = json.Unmarshal(data, &aux); err != nil {
		return err
	}

	f.Timestamp, err = time.Parse("2006-01-02T15:04:05", aux.Timestamp)
	if err != nil {
		return err
	}

	f.ExpiryTime, err = time.Parse("2006-01-02T15:04:05", aux.ExpiryTime)
	if err != nil {
		return err
	}

	return nil
}

func (f *ForexRateBookingResponse) MarshalJSON() ([]byte, error) {
	type Alias ForexRateBookingResponse
	return json.Marshal(&struct {
		Timestamp  string `json:"timestamp"`
		ExpiryTime string `json:"expiryTime"`
		*Alias
	}{
		Timestamp:  f.Timestamp.Format("2006-01-02T15:04:05"),
		ExpiryTime: f.ExpiryTime.Format("2006-01-02T15:04:05"),
		Alias:      (*Alias)(f),
	})
}
