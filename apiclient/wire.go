package apiclient

import "github.com/google/wire"

var Providers = wire.NewSet(NewForexApiClient)
