package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	dbPass := "abc123"
	os.Setenv("ENV", "sit")          // env specific configuration
	os.Setenv("DB_PASSWORD", dbPass) // environment variable

	err := LoadConfig()

	assert.Nil(t, err)
	assert.Equal(t, 8080, AppConfig.ServerPort)
	assert.Equal(t, "http://forex.sit.example.com:9090", AppConfig.ForexRateApiUrl)
	assert.Equal(t, dbPass, AppConfig.DbPassword)
}
