package postgres

import "time"

type Config struct {
	Host     string        `yaml:"host" env:"DB_HOST" envDefault:"localhost"`
	Port     string        `yaml:"port" env:"DB_PORT" envDefault:"5432"`
	User     string        `yaml:"user" env:"DB_USER" envDefault:"postgres"`
	Password string        `yaml:"password" env:"DB_PASSWORD" envDefault:"postgres"`
	Name     string        `yaml:"name" env:"DB_NAME" envDefault:"postgres"`
	Timeout  time.Duration `yaml:"timeout" env:"DB_TIMEOUT"`
}
