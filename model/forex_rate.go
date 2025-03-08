package model

import (
	"encoding/json"
	"time"
)

type ForexRate struct {
	Timestamp       time.Time
	BaseCurrency    string
	CounterCurrency string
	BuyRate         float32
	SellRate        float32
	Spread          float32
}

func (f *ForexRate) UnmarshalJSON(data []byte) error {
	type Alias ForexRate
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

func (f *ForexRate) MarshalJSON() ([]byte, error) {
	type Alias ForexRate
	return json.Marshal(&struct {
		Timestamp string `json:"timestamp"`
		*Alias
	}{
		Timestamp: f.Timestamp.Format("2006-01-02T15:04:05"),
		Alias:     (*Alias)(f),
	})
}
