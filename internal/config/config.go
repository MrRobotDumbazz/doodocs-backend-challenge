package config

import (
	"fmt"
	"os"
)

type Config struct {
	Email    string `env:"EMAIL"`
	Password string `env:"PASSWORD"`
}

func MustLoad() (*Config, error) {
	const op = "config.MustLoad"
	email, password := os.Getenv("EMAIL"), os.Getenv("PASSWORD")
	if email == "" || password == "" {
		return nil, fmt.Errorf("%s: nil email or password", op)
	}
	cfg := Config{
		Email:    email,
		Password: password,
	}
	return &cfg, nil
}
