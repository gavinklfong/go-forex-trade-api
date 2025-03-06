package main

import (
	"fmt"
	"net/http"

	"github.com/gavinklfong/go-rest-api-demo/model"
	"github.com/gin-gonic/gin"
)

type Router struct {
	Engine *gin.Engine
}

func SetupRouter() *Router {
	r := &Router{Engine: gin.Default()}
	r.doSetup()
	return r
}

func (r *Router) doSetup() {
	r.Engine.POST("/rates/book", func(c *gin.Context) {
		var request model.ForexRateBookingRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		fmt.Printf("booking request: %+v\n", request)

		c.JSON(http.StatusOK, gin.H{"status": "OK"})
	})
}
