package controller

import (
	"github.com/gavinklfong/go-rest-api-demo/service"
	"github.com/gin-gonic/gin"
)

type TradeDealController struct {
	forexTradeDealService *service.ForexTradeDealService
}

func NewTradeDealController(forexTradeDealService *service.ForexTradeDealService) *TradeDealController {
	return &TradeDealController{forexTradeDealService}
}

func (e *TradeDealController) SubmitTradeDeal(c *gin.Context) {

}

func (e *TradeDealController) GetTradeDeal(c *gin.Context) {

}
