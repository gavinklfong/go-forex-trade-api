package apiclient

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/gavinklfong/go-forex-trade-api/apiclient/model"
)

type ForexApiClientImpl struct {
	url string
}

func NewForexApiClient(url string) ForexApiClient {
	return &ForexApiClientImpl{url: url}
}

func (c *ForexApiClientImpl) GetRateByCurrencyPair(base, counter string) (*model.ForexRateResponse, error) {
	requestURL := fmt.Sprintf("%s/rates/%s-%s", c.url, base, counter)
	slog.Info("GET %s", requestURL)
	res, err := http.Get(requestURL)
	if err != nil {
		slog.Error("error making http request: %s\n", err)
		return nil, err
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		slog.Error("error reading response body: %s\n", err)
		return nil, err
	}

	var rate model.ForexRateResponse
	err = json.Unmarshal(resBody, &rate)
	if err != nil {
		slog.Error("error parsing response body: %s", err)
		return nil, err
	}

	return &rate, nil
}

func (c *ForexApiClientImpl) GetRateByBaseCurrency(base string) (*model.ForexRateResponse, error) {
	requestURL := fmt.Sprintf("%s/rates/%s", c.url, base)
	slog.Info(fmt.Sprintf("GET %s", requestURL))
	res, err := http.Get(requestURL)
	if err != nil {
		slog.Error("error making http request: %s\n", err)
		return nil, err
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		slog.Error("error reading response body: %s\n", err)
		return nil, err
	}

	var rate model.ForexRateResponse
	err = json.Unmarshal(resBody, &rate)
	if err != nil {
		slog.Error("error parsing response body: %s", err)
		return nil, err
	}

	return &rate, nil
}
