package controller

import (
	"net/http"

	"github.com/gavinklfong/go-forex-trade-api/service"
	"github.com/gin-gonic/gin"
)

const DEFAULT_BASE_CURRENCY = "GBP"

type GetRateByBaseCurrencyRequest struct {
	BaseCurrency string `uri:"baseCurrency" binding:"required,string"`
}

type GetRateByCurrencyPairRequest struct {
	BaseCurrency    string `uri:"baseCurrency" binding:"required,string"`
	CounterCurrency string `uri:"counterCurrency" binding:"required,string"`
}

type GetRateController struct {
	r service.ForexRateService
}

func NewGetRateController(ForexRateService service.ForexRateService) *GetRateController {
	return &GetRateController{r: ForexRateService}
}

func (e *GetRateController) GetRateByBaseCurrency(c *gin.Context) {
	var request GetRateByBaseCurrencyRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	rates, err := e.r.GetRatesByBaseCurrency(request.BaseCurrency)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, rates)
}

func (e *GetRateController) GetRateByCurrencyPair(c *gin.Context) {
	var request GetRateByCurrencyPairRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	rate, err := e.r.GetRateByCurrencyPair(request.BaseCurrency, request.CounterCurrency)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, rate)
}

func (e *GetRateController) GetDefaultRates(c *gin.Context) {
	rates, err := e.r.GetRatesByBaseCurrency(DEFAULT_BASE_CURRENCY)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, rates)
}
