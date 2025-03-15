//go:build wireinject
// +build wireinject

package main

import (
	"github.com/gavinklfong/go-rest-api-demo/apiclient"
	"github.com/gavinklfong/go-rest-api-demo/config"
	"github.com/gavinklfong/go-rest-api-demo/dao"
	"github.com/gavinklfong/go-rest-api-demo/endpoint"
	"github.com/gavinklfong/go-rest-api-demo/router"
	"github.com/gavinklfong/go-rest-api-demo/service"
	"github.com/google/wire"
)

func InitializeRouter() *router.Router {

	// var allProviders = wire.NewSet(apiclient.Providers, dao.Providers, service.Providers, endpoint.Providers, router.Providers)

	// config.LoadConfig()

	wire.Build(config.InitializeDBConnection, apiclient.Providers,
		dao.Providers, service.Providers, endpoint.Providers, router.Providers)

	return &router.Router{}
}

// func injectForexApiClient() apiclient.ForexApiClient {
// 	wire.Build(wire.Value(apiclient.ForexApiClient{url: config.AppConfig.ForexRateApiUrl}))
// 	return apiclient.NewForexApiClient{}
// }
