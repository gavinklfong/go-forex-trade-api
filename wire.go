//go:build wireinject
// +build wireinject

package main

import (
	"github.com/gavinklfong/go-rest-api-demo/endpoint"
	"github.com/gavinklfong/go-rest-api-demo/router"
	"github.com/gavinklfong/go-rest-api-demo/service"
	"github.com/google/wire"
)

func InitializeRouter() *router.Router {
	wire.Build(router.NewRouter,
		endpoint.NewGetRateEndpoint, endpoint.NewBookRateEndpoint,
		service.NewRateService, service.NewDealService)
	return &router.Router{}
}
