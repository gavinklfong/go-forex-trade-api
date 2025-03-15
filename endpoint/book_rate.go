package endpoint

import (
	"net/http"

	"github.com/gavinklfong/go-rest-api-demo/model"
	"github.com/gavinklfong/go-rest-api-demo/service"
	"github.com/gin-gonic/gin"
)

type BookRateEndpoint struct {
	r *service.ForexRateService
}

func NewBookRateEndpoint(rateService *service.ForexRateService) *BookRateEndpoint {
	return &BookRateEndpoint{r: rateService}
}

func (e *BookRateEndpoint) BookRate(c *gin.Context) {
	var request model.ForexRateBookingRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	booking := e.r.BookRate(&request)

	c.JSON(http.StatusOK, booking)
}
