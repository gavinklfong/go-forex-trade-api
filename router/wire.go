package router

import "github.com/google/wire"

var Providers = wire.NewSet(NewRouter)
