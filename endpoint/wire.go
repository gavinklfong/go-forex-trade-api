package endpoint

import "github.com/google/wire"

var Providers = wire.NewSet(NewBookRateEndpoint, NewGetRateEndpoint, NewTradeDealEndpoint)
