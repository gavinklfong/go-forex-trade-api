package model

import (
	"encoding/json"
	"time"
)

type ForexRateResponse struct {
	Timestamp       time.Time
	BaseCurrency    string
	CounterCurrency string
	BuyRate         float32
	SellRate        float32
	Spread          float32
}

func (f *ForexRateResponse) UnmarshalJSON(data []byte) error {
	type Alias ForexRateResponse
	aux := &struct {
		Timestamp string
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

	return nil
}

func (f *ForexRateResponse) MarshalJSON() ([]byte, error) {
	type Alias ForexRateResponse
	return json.Marshal(&struct {
		Timestamp string `json:"timestamp"`
		*Alias
	}{
		Timestamp: f.Timestamp.Format("2006-01-02T15:04:05"),
		Alias:     (*Alias)(f),
	})
}
