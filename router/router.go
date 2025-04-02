package router

import (
	"github.com/gavinklfong/go-forex-trade-api/controller"
	"github.com/gin-gonic/gin"
)

type Router struct {
	e                   *gin.Engine
	getRateController   *controller.GetRateController
	bookRateController  *controller.BookRateController
	tradeDealController *controller.TradeDealController
}

func NewRouter(getRateController *controller.GetRateController,
	bookRateController *controller.BookRateController,
	tradeDealController *controller.TradeDealController) *Router {
	r := &Router{e: gin.Default(),
		getRateController:   getRateController,
		bookRateController:  bookRateController,
		tradeDealController: tradeDealController}
	r.doSetup()
	return r
}

func (r *Router) doSetup() {

	r.e.GET("/rates/:baseCurrency/:counterCurrency", r.getRateController.GetRateByCurrencyPair)
	r.e.GET("/rates/:baseCurrency", r.getRateController.GetRateByBaseCurrency)
	r.e.GET("/rates/latest", r.getRateController.GetAllRates)

	r.e.POST("/rates/book", r.bookRateController.BookRate)

	r.e.GET("/deals", r.tradeDealController.GetTradeDeal)
	r.e.POST("/deals", r.tradeDealController.SubmitTradeDeal)
}

func (r *Router) Run(s string) {
	r.e.Run(s)
}
