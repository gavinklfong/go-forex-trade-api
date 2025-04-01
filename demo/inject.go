package demo

import (
	"log"
	"log/slog"

	"go.uber.org/dig"
)

var c *dig.Container
var a *ServiceA

func InitializeServiceA() (*ServiceA, error) {

	c = dig.New()

	err := c.Provide(NewComponentY)
	if err != nil {
		slog.Error("New Component Y failed: %v", err)
		return nil, err
	}

	err = c.Provide(NewComponentZ)
	if err != nil {
		slog.Error("New Component Z failed: %v", err)
		return nil, err
	}

	err = c.Provide(NewServiceA)
	if err != nil {
		slog.Error("New Service A failed: %v", err)
		return nil, err
	}

	err = c.Invoke(func(serviceA *ServiceA) {
		a = serviceA
	})
	if err != nil {
		log.Panicf("fail to initialize target: %v", err)
	}

	return a, nil

}

func initComponent() error {
	return c.Invoke(func(serviceA *ServiceA) {
		a = serviceA
	},
	)
}
