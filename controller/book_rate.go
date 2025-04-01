package controller

import (
	"net/http"

	"github.com/gavinklfong/go-rest-api-demo/model"
	"github.com/gavinklfong/go-rest-api-demo/service"
	"github.com/gin-gonic/gin"
)

type BookRateController struct {
	r *service.ForexRateService
}

func NewBookRateController(rateService *service.ForexRateService) *BookRateController {
	return &BookRateController{r: rateService}
}

func (e *BookRateController) BookRate(c *gin.Context) {
	var request model.ForexRateBookingRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	booking := e.r.BookRate(&request)

	c.JSON(http.StatusOK, booking)
}
