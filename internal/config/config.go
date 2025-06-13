package config

import (
	"sync"

	"github.com/spf13/viper"
)

type Config struct {
	PG string `yaml:"pg_conn_string" mapstructure:"pg_conn_string"`
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
		return nil, err
	}
	return cfg, nil
}
