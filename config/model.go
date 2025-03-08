package config

type Config struct {
	ServerPort           int    `config:"server.port"`
	RateBookingDuration  int    `config:"app.rate-booking-duration"`
	DefaultBaseCurrency  string `config:"app.default-base-currency"`
	DefaultAdditionalPip int    `config:"app.default-additional-pip"`
	ForexRateApiUrl      string `config:"app.forex-rate-api-url"`
	DbUrl                string `config:"db.url"`
	DbUser               string `config:"db.user"`
	DbPassword           string `config:"db.password"`
}
