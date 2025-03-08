package config

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

var AppConfig *Config

func LoadConfig() error {

	slog.Info("Initialize Koanf")
	k := koanf.NewWithConf(koanf.Conf{
		Delim: ".",
	})

	// Load YAML config
	slog.Info("Load default configuration")
	if err := k.Load(file.Provider("./application.yaml"), yaml.Parser()); err != nil {
		slog.Error(fmt.Sprintf("error loading config: %v", err))
		return err
	}

	// Load env specific config
	environment := os.Getenv("ENV")
	if len(strings.TrimSpace(environment)) > 0 {
		slog.Info(fmt.Sprintf("Loading env specific configuration - %s", environment))
		// Load Dot Env config
		if err := k.Load(file.Provider(fmt.Sprintf("./application-%s.yaml", environment)), yaml.Parser()); err != nil {
			slog.Error(fmt.Sprintf("error loading config: %v", err))
			return err
		}
	} else {
		slog.Warn("Environment variable 'ENV' is not defined, skip environment specific configuration")
	}

	// Load environment variable
	slog.Info("Load from environment variables")
	err := k.Load(env.Provider("", ".", func(s string) string {
		return strings.Replace(strings.ToLower(s), "_", ".", -1)
	}), nil)
	if err != nil {
		slog.Error(fmt.Sprintf("error loading config: %v", err))
		return err
	}

	AppConfig = &Config{}
	err = k.UnmarshalWithConf("", AppConfig, koanf.UnmarshalConf{Tag: "config", FlatPaths: true})
	if err != nil {
		slog.Error(fmt.Sprintf("error unmarshaling config: %v", err))
		return err
	}

	return nil
}
