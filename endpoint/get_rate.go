package endpoint

import (
	"net/http"

	"github.com/gavinklfong/go-rest-api-demo/service"
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

type GetRateEndpoint struct {
	r *service.ForexRateService
}

func NewGetRateEndpoint(ForexRateService *service.ForexRateService) *GetRateEndpoint {
	return &GetRateEndpoint{r: ForexRateService}
}

func (e *GetRateEndpoint) GetRateByBaseCurrency(c *gin.Context) {
	var request GetRateByBaseCurrencyRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	rates := e.r.GetRatesByBaseCurrency(request.BaseCurrency)

	c.JSON(http.StatusOK, rates)
}

func (e *GetRateEndpoint) GetRateByCurrencyPair(c *gin.Context) {
	var request GetRateByCurrencyPairRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	rate := e.r.GetRateByCurrencyPair(request.BaseCurrency, request.CounterCurrency)

	c.JSON(http.StatusOK, rate)
}

func (e *GetRateEndpoint) GetAllRates(c *gin.Context) {
	rates := e.r.GetRatesByBaseCurrency(DEFAULT_BASE_CURRENCY)
	c.JSON(http.StatusOK, rates)
}
