package router

import (
	"github.com/gavinklfong/go-rest-api-demo/endpoint"
	"github.com/gin-gonic/gin"
)

type Router struct {
	e                 *gin.Engine
	getRateEndpoint   *endpoint.GetRateEndpoint
	bookRateEndpoint  *endpoint.BookRateEndpoint
	tradeDealEndpoint *endpoint.TradeDealEndpoint
}

func NewRouter(getRateEndpoint *endpoint.GetRateEndpoint,
	bookRateEndpoint *endpoint.BookRateEndpoint,
	tradeDealEndpoint *endpoint.TradeDealEndpoint) *Router {
	r := &Router{e: gin.Default(),
		getRateEndpoint:   getRateEndpoint,
		bookRateEndpoint:  bookRateEndpoint,
		tradeDealEndpoint: tradeDealEndpoint}
	r.doSetup()
	return r
}

func (r *Router) doSetup() {

	r.e.GET("/rates/:baseCurrency}/:counterCurrency", r.getRateEndpoint.GetRateByCurrencyPair)
	r.e.GET("/rates/:baseCurrency", r.getRateEndpoint.GetRateByBaseCurrency)
	r.e.GET("/rates/latest", r.getRateEndpoint.GetAllRates)

	r.e.POST("/rates/book", r.bookRateEndpoint.BookRate)

	r.e.GET("/deals", r.tradeDealEndpoint.GetTradeDeal)
	r.e.POST("/deals", r.tradeDealEndpoint.SubmitTradeDeal)
}

func (r *Router) Run(s string) {
	r.e.Run(s)
}
