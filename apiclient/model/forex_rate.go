package model

import (
	"encoding/json"
	"time"
)

type ForexRateResponse struct {
	ID    string
	Date  time.Time
	Base  string
	Rates map[string]float32
}

func (f *ForexRateResponse) UnmarshalJSON(data []byte) error {
	type Alias ForexRateResponse
	aux := &struct {
		Date string
		*Alias
	}{
		Alias: (*Alias)(f),
	}

	var err error

	if err = json.Unmarshal(data, &aux); err != nil {
		return err
	}

	f.Date, err = time.Parse("2006-01-02", aux.Date)
	if err != nil {
		return err
	}

	return nil
}

func (f *ForexRateResponse) MarshalJSON() ([]byte, error) {
	type Alias ForexRateResponse
	return json.Marshal(&struct {
		Date string `json:"date"`
		*Alias
	}{
		Date:  f.Date.Format("2006-01-02"),
		Alias: (*Alias)(f),
	})
}
