package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/spf13/viper"
)

const (
	configName         = "config"
	configOverrideName = "config.override"
)

type Config struct {
	ValkeyConnString string `mapstructure:"valkey_conn_string"`
	OAPIPath         string `mapstructure:"oapi_path"`
	CryptoKey        string `mapstructure:"crypto_key"`
}

var (
	cfg        *Config
	createOnce sync.Once
)

func createConfig() error {
	var err error

	createOnce.Do(func() {
		viper.AddConfigPath(".")
		viper.SetConfigName(configName)

		err = viper.ReadInConfig()
		if err != nil {
			return
		}

		err = viper.BindEnv("crypto_key", "CRYPTO_KEY")
		if err != nil {
			return
		}

		cfg = &Config{}
		err = viper.Unmarshal(cfg)
	})
	if err != nil {
		return fmt.Errorf("unable to read config: %w", err)
	}

	return nil
}

func Overload() (*Config, error) {
	err := createConfig()
	if err != nil {
		return nil, fmt.Errorf("unable to create: %w", err)
	}

	_, err = os.Stat(filepath.Base(filepath.Join(".", configOverrideName+".yaml")))
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("unable to check config override existence: %w", err)
	}

	if err == nil {
		viper.AddConfigPath(".")
		viper.SetConfigName(configOverrideName)

		err = viper.MergeInConfig()
		if err != nil {
			return nil, fmt.Errorf("unable to merge: %w", err)
		}
	}

	err = viper.Unmarshal(cfg)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal: %w", err)
	}

	return cfg, nil
}
