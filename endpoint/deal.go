package endpoint

import (
	"github.com/gavinklfong/go-rest-api-demo/service"
	"github.com/gin-gonic/gin"
)

type TradeDealEndpoint struct {
	forexTradeDealService *service.ForexTradeDealService
}

func NewTradeDealEndpoint(forexTradeDealService *service.ForexTradeDealService) *TradeDealEndpoint {
	return &TradeDealEndpoint{forexTradeDealService}
}

func (e *TradeDealEndpoint) SubmitTradeDeal(c *gin.Context) {

}

func (e *TradeDealEndpoint) GetTradeDeal(c *gin.Context) {

}
