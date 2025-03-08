package router

import (
	"fmt"
	"net/http"

	"github.com/gavinklfong/go-rest-api-demo/model"
	"github.com/gin-gonic/gin"
)

var engine *gin.Engine

func SetupRouter() *gin.Engine {
	r := gin.Default()
	doSetup(r)
	return r
}

func doSetup(engine *gin.Engine) {
	engine.POST("/rates/book", func(c *gin.Context) {
		var request model.ForexRateBookingRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		fmt.Printf("booking request: %+v\n", request)

		c.JSON(http.StatusOK, gin.H{"status": "OK"})
	})
}
