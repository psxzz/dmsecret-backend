package config

import (
	"fmt"
	"sync"

	"github.com/spf13/viper"
)

type Config struct {
	PGConnString     string `yaml:"pg_conn_string" mapstructure:"pg_conn_string"`
	ValkeyConnString string `yaml:"valkey_conn_string" mapstructure:"valkey_conn_string"`
	OAPIPath         string `yaml:"oapi_path" mapstructure:"oapi_path"`
}

var (
	cfg        *Config
	createOnce sync.Once
)

func Create() (*Config, error) {
	var err error
	createOnce.Do(func() {
		viper.AddConfigPath(".")
		viper.SetConfigName("config")

		err = viper.ReadInConfig()

		cfg = &Config{}
		err = viper.Unmarshal(cfg)
	})

	if err != nil {
		return nil, fmt.Errorf("unable to read config: %w", err)
	}
	return cfg, nil
}
