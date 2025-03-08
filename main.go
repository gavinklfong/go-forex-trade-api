package main

import (
	"fmt"

	"github.com/gavinklfong/go-rest-api-demo/config"
	"github.com/gavinklfong/go-rest-api-demo/endpoint"
	"github.com/gavinklfong/go-rest-api-demo/router"
	"github.com/gavinklfong/go-rest-api-demo/service"
)

func main() {

	config.LoadConfig()
	rateService := service.NewRateService()
	getRateEndpoint := endpoint.NewGetRateEndpoint(rateService)
	bookRateEndpoint := endpoint.NewBookRateEndpoint(rateService)
	r := router.NewRouter(getRateEndpoint, bookRateEndpoint)

	// Listen and serve on 0.0.0.0:8080
	r.Run(fmt.Sprintf(":%d", config.AppConfig.ServerPort))
}
