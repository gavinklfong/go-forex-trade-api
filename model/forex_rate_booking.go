package model

import (
	"encoding/json"
	"time"
)

const TIMESTAMP_FORMAT = "2006-01-02T15:04:05"

type ForexRateBookingRequest struct {
	BaseCurrency       string
	CounterCurrency    string
	BaseCurrencyAmount float32
	TradeAction        string
	CustomerId         int32
}

type ForexRateBooking struct {
	ForexRateBookingRequest
	Timestamp  time.Time
	Rate       float32
	BookingRef string
	ExpiryTime time.Time
}

func (f *ForexRateBooking) UnmarshalJSON(data []byte) error {
	type Alias ForexRateBooking
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

	f.Timestamp, err = time.Parse(TIMESTAMP_FORMAT, aux.Timestamp)
	if err != nil {
		return err
	}

	f.ExpiryTime, err = time.Parse(TIMESTAMP_FORMAT, aux.ExpiryTime)
	if err != nil {
		return err
	}

	return nil
}

func (f *ForexRateBooking) MarshalJSON() ([]byte, error) {
	type Alias ForexRateBooking
	return json.Marshal(&struct {
		Timestamp  string `json:"timestamp"`
		ExpiryTime string `json:"expiryTime"`
		*Alias
	}{
		Timestamp:  f.Timestamp.Format(TIMESTAMP_FORMAT),
		ExpiryTime: f.ExpiryTime.Format(TIMESTAMP_FORMAT),
		Alias:      (*Alias)(f),
	})
}
