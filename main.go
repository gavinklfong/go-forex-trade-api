package main

import (
	"database/sql"
	"fmt"
	"log"
	"log/slog"

	"github.com/gavinklfong/go-forex-trade-api/apiclient"
	"github.com/gavinklfong/go-forex-trade-api/config"
	"github.com/gavinklfong/go-forex-trade-api/controller"
	"github.com/gavinklfong/go-forex-trade-api/dao"
	"github.com/gavinklfong/go-forex-trade-api/router"
	"github.com/gavinklfong/go-forex-trade-api/service"
	"go.uber.org/dig"
)

var c *dig.Container
var r *router.Router

func main() {

	err := config.LoadConfig()
	if err != nil {
		log.Panicf("fail to load configuration: %v", err)
	}

	err = provideComponents()
	if err != nil {
		log.Panicf("fail to set up providers: %v", err)
	}

	err = initComponent()
	if err != nil {
		log.Panicf("fail to initialize target: %v", err)
	}

	r.Run(fmt.Sprintf(":%d", config.AppConfig.ServerPort))

	config.CloseDBConnection()
}

func initComponent() error {
	return c.Invoke(func(t *router.Router) {
		r = t
	},
	)
}

func provideComponents() error {
	c = dig.New()
	err := c.Provide(config.InitializeDBConnection)
	if err != nil {
		slog.Error("application initialization failed: %v", err)
		return err
	}

	err = c.Provide(func(db *sql.DB) (*dao.CustomerDao, *dao.ForexRateDao, *dao.ForexTradeDealDao) {
		return dao.NewCustomerDao(db), dao.NewForexRateDao(db), dao.NewForexTradeDealDao(db)
	})
	if err != nil {
		slog.Error("application initialization failed: %v", err)
		return err
	}

	err = c.Provide(func() *apiclient.ForexApiClient {
		return apiclient.NewForexApiClient(config.AppConfig.ForexRateApiUrl)
	})
	if err != nil {
		slog.Error("application initialization failed: %v", err)
		return err
	}

	providers := [...]interface{}{
		service.NewForexRateService,
		service.NewForexTradeDealService,
		controller.NewBookRateController,
		controller.NewGetRateController,
		controller.NewTradeDealController,
		router.NewRouter,
	}

	for _, provider := range providers {
		err = c.Provide(provider)
		if err != nil {
			slog.Error("application initialization failed: %v", err)
			return err
		}
	}

	return nil

}
