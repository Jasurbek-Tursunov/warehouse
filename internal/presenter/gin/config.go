package pgin

import (
	"github.com/ilyakaznacheev/cleanenv"
	"time"
)

type Config struct {
	Port    int           `yaml:"port" env:"PORT"`
	Timeout time.Duration `yaml:"timeout" env:"TIMEOUT"`
}

func MustLoadEnvConfig() *Config {
	var cfg Config

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		panic(err)
	}

	return &cfg
}
