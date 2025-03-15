package apiclient

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gavinklfong/go-rest-api-demo/apiclient/model"
)

type ForexApiClient struct {
	url string
}

func NewForexApiClient(url string) *ForexApiClient {
	return &ForexApiClient{url: url}
}

// func NewForexApiClient() *ForexApiClient {
// 	return &ForexApiClient{url: config.AppConfig.ForexRateApiUrl}
// }

func (c *ForexApiClient) GetLatestRates() ([]model.ForexRateResponse, error) {
	requestURL := fmt.Sprintf("%s/rates/latest", c.url)
	res, err := http.Get(requestURL)
	if err != nil {
		log.Fatalf("error making http request: %s\n", err)
		return nil, err
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("error reading response body: %s\n", err)
		return nil, err
	}

	var rates []model.ForexRateResponse
	err = json.Unmarshal(resBody, &rates)
	if err != nil {
		log.Fatalf("error parsing response body: %s", err)
		return nil, err
	}

	return rates, nil
}

func (c *ForexApiClient) GetLatestRate(base, counter string) (*model.ForexRateResponse, error) {
	requestURL := fmt.Sprintf("%s/rates/latest/%s/%s", c.url, base, counter)
	res, err := http.Get(requestURL)
	if err != nil {
		log.Fatalf("error making http request: %s\n", err)
		return nil, err
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("error reading response body: %s\n", err)
		return nil, err
	}

	var rate model.ForexRateResponse
	err = json.Unmarshal(resBody, &rate)
	if err != nil {
		log.Fatalf("error parsing response body: %s", err)
		return nil, err
	}

	return &rate, nil
}

// func (c *ForexApiClient) BookRate(booking *model.ForexRateBookingRequest) *model.ForexRateBookingResponse {

// }
