package config

import (
	"log"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Address     string        `env:"ADDRESS" env-default:"0.0.0.0:8080"`
	Timeout     time.Duration `env:"TIMEOUT" env-default:"5s"`
	IdleTimeout time.Duration `env:"IDLE_TIMEOUT" env-default:"60s"`
	Email       string        `env:"EMAIL"`
	Password    string        `env:"PASSWORD"`
}

func MustLoad() *Config {
	cfg := Config{}
	err := cleanenv.ReadConfig("./.env", &cfg)
	if err != nil {
		log.Fatalf("error reading config file: %s", err)
	}
	return &cfg
}
