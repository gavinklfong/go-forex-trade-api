package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gavinklfong/go-rest-api-demo/model"
	"github.com/stretchr/testify/assert"
)

func TestRouter(t *testing.T) {
	r := SetupRouter()

	request := model.ForexRateBookingRequest{
		BaseCurrency:       "GBP",
		CounterCurrency:    "USD",
		BaseCurrencyAmount: 1200,
		TradeAction:        "BUY",
		CustomerId:         2,
	}

	requestJson, _ := json.Marshal(request)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/rates/book", strings.NewReader(string(requestJson)))
	r.Engine.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	fmt.Printf("response body: %+v/n", w.Body.String())
}
