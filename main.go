package main

import (
	"log"
	"log/slog"

	"github.com/gavinklfong/go-rest-api-demo/apiclient"
	"github.com/gavinklfong/go-rest-api-demo/config"
	"github.com/gavinklfong/go-rest-api-demo/router"
	"go.uber.org/dig"
)

var c *dig.Container
var r *router.Router
var target *apiclient.ForexApiClient

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

	// r.Run(fmt.Sprintf(":%s", config.AppConfig.ServerPort))

}

func initComponent() error {
	return c.Invoke(func(t *apiclient.ForexApiClient) {
		if t == nil {
			slog.Error("target is nil")
		}
	},
	)
}

func initRouter() error {
	return c.Invoke(
		func(router *router.Router) {
			r = router
		},
	)
}

func provideComponents() error {
	c := dig.New()
	// err := c.Provide(config.InitializeDBConnection)
	// if err != nil {
	// 	slog.Error("application initialization failed: %v", err)
	// 	return err
	// }

	// err = c.Provide(func(db *sql.DB) (*dao.CustomerDao, *dao.ForexRateDao, *dao.ForexTradeDealDao) {
	// 	return dao.NewCustomerDao(db), dao.NewForexRateDao(db), dao.NewForexTradeDealDao(db)
	// })
	// if err != nil {
	// 	slog.Error("application initialization failed: %v", err)
	// 	return err
	// }

	err := c.Provide(func() *apiclient.ForexApiClient {
		return apiclient.NewForexApiClient(config.AppConfig.ForexRateApiUrl)
	})
	if err != nil {
		slog.Error("application initialization failed: %v", err)
		return err
	}

	// err = c.Provide(service.NewForexRateService)
	// err = c.Provide(service.NewForexTradeDealService)
	// err = c.Provide(endpoint.NewBookRateEndpoint)
	// err = c.Provide(endpoint.NewGetRateEndpoint)
	// err = c.Provide(endpoint.NewTradeDealEndpoint)
	// err = c.Provide(router.NewRouter)

	// providers := [...]interface{}{
	// 	service.NewForexRateService,
	// 	service.NewForexTradeDealService,
	// 	endpoint.NewBookRateEndpoint,
	// 	endpoint.NewGetRateEndpoint,
	// 	endpoint.NewTradeDealEndpoint,
	// 	router.NewRouter,
	// }

	// for _, provider := range providers {
	// 	err = c.Provide(provider)
	// 	if err != nil {
	// 		slog.Error("application initialization failed: %v", err)
	// 		return err
	// 	}
	// }

	return nil

}
