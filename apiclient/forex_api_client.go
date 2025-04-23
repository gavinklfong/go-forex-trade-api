package apiclient

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/gavinklfong/go-forex-trade-api/apiclient/model"
	retryablehttp "github.com/hashicorp/go-retryablehttp"
)

type ForexApiClientImpl struct {
	url string
}

func NewForexApiClient(url string) ForexApiClient {
	return &ForexApiClientImpl{url: url}
}

func newClient() (*http.Client, error) {
	retryWaitMin, err := time.ParseDuration("500ms")
	if err != nil {
		return nil, err
	}

	retryWaitMax, err := time.ParseDuration("5s")
	if err != nil {
		return nil, err
	}

	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = 3
	retryClient.RetryWaitMin = retryWaitMin
	retryClient.RetryWaitMax = retryWaitMax
	retryClient.CheckRetry = retryPolicy
	retryClient.Backoff = backOffPolicy

	return retryClient.StandardClient(), nil
}

func retryPolicy(ctx context.Context, resp *http.Response, err error) (bool, error) {
	shouldRetry, err := retryablehttp.DefaultRetryPolicy(ctx, resp, err)
	if err != nil {
		return false, err
	}

	if shouldRetry {
		return shouldRetry, nil
	}

	if resp.StatusCode == http.StatusBadRequest {
		return true, nil
	}

	return false, nil
}

func backOffPolicy(min, max time.Duration, attemptNum int, resp *http.Response) time.Duration {
	if min != max {
		return retryablehttp.DefaultBackoff(min, max, attemptNum, resp)
	} else {
		return min
	}

}

func (c *ForexApiClientImpl) GetRateByCurrencyPair(base, counter string) (*model.ForexRateResponse, error) {
	requestURL := fmt.Sprintf("%s/rates/%s-%s", c.url, base, counter)
	slog.Info(fmt.Sprintf("GET %s", requestURL))

	httpClient, err := newClient()
	if err != nil {
		return nil, err
	}

	res, err := httpClient.Get(requestURL)
	if err != nil {
		slog.Error(fmt.Sprintf("error making http request: %s\n", err))
		return nil, err
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		slog.Error(fmt.Sprintf("error reading response body: %s\n", err))
		return nil, err
	}

	var rate model.ForexRateResponse
	err = json.Unmarshal(resBody, &rate)
	if err != nil {
		slog.Error(fmt.Sprintf("error parsing response body: %s", err))
		return nil, err
	}

	return &rate, nil
}

func (c *ForexApiClientImpl) GetRateByBaseCurrency(base string) (*model.ForexRateResponse, error) {
	requestURL := fmt.Sprintf("%s/rates/%s", c.url, base)
	slog.Info(fmt.Sprintf("GET %s", requestURL))
	res, err := http.Get(requestURL)
	if err != nil {
		slog.Error(fmt.Sprintf("error making http request: %s\n", err))
		return nil, err
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		slog.Error(fmt.Sprintf("error reading response body: %s\n", err))
		return nil, err
	}

	var rate model.ForexRateResponse
	err = json.Unmarshal(resBody, &rate)
	if err != nil {
		slog.Error(fmt.Sprintf("error parsing response body: %s", err))
		return nil, err
	}

	return &rate, nil
}
