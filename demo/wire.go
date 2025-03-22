//go:build wireinject
// +build wireinject

package demo

import "github.com/google/wire"

func initializeServiceA() *ServiceA {
	wire.Build(NewComponentY, NewComponentZ, NewServiceA)
	return &ServiceA{}
}
