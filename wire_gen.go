// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/gavinklfong/go-rest-api-demo/apiclient"
	"github.com/gavinklfong/go-rest-api-demo/config"
	"github.com/gavinklfong/go-rest-api-demo/dao"
	"github.com/gavinklfong/go-rest-api-demo/endpoint"
	"github.com/gavinklfong/go-rest-api-demo/router"
	"github.com/gavinklfong/go-rest-api-demo/service"
)

// Injectors from wire.go:

func InitializeRouter() (*router.Router, error) {
	db, err := config.InitializeDBConnection()
	if err != nil {
		return nil, err
	}
	customerDao := dao.NewCustomerDao(db)
	forexRateDao := dao.NewForexRateDao(db)
	forexApiClient := apiclient.NewForexApiClient()
	forexRateService := service.NewForexRateService(customerDao, forexRateDao, forexApiClient)
	getRateEndpoint := endpoint.NewGetRateEndpoint(forexRateService)
	bookRateEndpoint := endpoint.NewBookRateEndpoint(forexRateService)
	forexTradeDealDao := dao.NewForexTradeDealDao(db)
	forexTradeDealService := service.NewForexTradeDealService(forexTradeDealDao, forexRateService)
	tradeDealEndpoint := endpoint.NewTradeDealEndpoint(forexTradeDealService)
	routerRouter := router.NewRouter(getRateEndpoint, bookRateEndpoint, tradeDealEndpoint)
	return routerRouter, nil
}
