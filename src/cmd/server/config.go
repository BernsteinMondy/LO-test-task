package main

import (
	"fmt"
	"github.com/caarlos0/env/v11"
)

type Config struct {
	HTTPServer HTTPServer `envPrefix:"HTTP_SERVER_"`
	Logger     Logger     `envPrefix:"LOGGER_"`
}

type HTTPServer struct {
	ListenAddr string `env:"LISTEN_ADDR,notEmpty"`
}

type Logger struct {
	BufferSize int `env:"BUFFER_SIZE,notEmpty"`
}

func loadConfigFromEnv() (Config, error) {
	c, err := env.ParseAs[Config]()
	if err != nil {
		return Config{}, fmt.Errorf("parse environment: %w", err)
	}

	return c, nil
}
