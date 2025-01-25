package gin

import (
	"time"
)

type Config struct {
	Port    int           `yaml:"port" env:"SERVER_PORT"`
	Timeout time.Duration `yaml:"timeout" env:"SERVER_TIMEOUT"`
}
