package config

import (
	"flag"
	"log"
	"os"
	"strings"

	"github.com/spf13/viper"
)

const CONFIG_PATH = "CONFIG_PATH"

type Option func(*viper.Viper)

func WithEnvPrefix(prefix string) Option {
	return func(v *viper.Viper) {
		v.SetEnvPrefix(prefix)
	}
}

func WithConfigName(name string) Option {
	return func(v *viper.Viper) {
		v.SetConfigName(name)
	}
}

func WithConfigType(configType string) Option {
	return func(v *viper.Viper) {
		v.SetConfigType(configType)
	}
}

func MustLoad[T any](opts ...Option) *T {
	cfg := new(T)
	configPath := fetchConfigPath()

	// Initialize Viper
	v := viper.New()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// Apply options
	for _, opt := range opts {
		opt(v)
	}

	// If a path to the config is specified, load from the file
	// In Kubernetes, all settings come through environment variables
	// from ConfigMap and Secret, so config files are only needed for local development
	if configPath != "" {
		v.SetConfigFile(configPath)

		if err := v.ReadInConfig(); err != nil {
			log.Fatalf("error reading config file: %v", err)
		}
	}

	// Fill the structure
	if err := v.Unmarshal(cfg); err != nil {
		log.Fatalf("error unmarshalling config: %v", err)
	}

	return cfg
}

func fetchConfigPath() string {
	var v string
	flag.StringVar(&v, "config", "", "path to config file")
	flag.Parse()

	if v == "" {
		v = os.Getenv(CONFIG_PATH)
	}

	return v
}
