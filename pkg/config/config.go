package config

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"reflect"
)

func MustLoad[Config any]() *Config {
	var cfg Config

	if reflect.TypeOf(cfg).Kind() != reflect.Struct {
		panic("Config must be a struct")
	}

	configPath := fetchConfigPath()
	if configPath == "" {
		return mustLoadEnv[Config]()
	}

	return mustLoadPath[Config](configPath)
}

func mustLoadPath[Config any](configPath string) *Config {
	// check config file
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("Config file does not exist: " + configPath)
	}
	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("Cannot read config: " + err.Error())
	}

	return &cfg
}

func mustLoadEnv[Config any]() *Config {
	var cfg Config

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		panic("Cannot read env variables: " + err.Error())
	}

	return &cfg
}

func fetchConfigPath() string {
	var value string

	flag.StringVar(&value, "config", "", "Path to config file")
	flag.Parse()

	if value == "" {
		value = os.Getenv("CONFIG_PATH")
	}

	return value
}
