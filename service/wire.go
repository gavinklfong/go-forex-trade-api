package service

import "github.com/google/wire"

var Providers = wire.NewSet(NewForexRateService, NewForexTradeDealService)
