package dao

import "github.com/google/wire"

var Providers = wire.NewSet(NewForexRateDao, NewForexTradeDealDao, NewCustomerDao)
