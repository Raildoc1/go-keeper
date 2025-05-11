package config

import (
	client "go-keeper/internal/client/config"
	logic "go-keeper/internal/client/logic/config"
	"time"
)

type Config struct {
	Client          client.Config
	ShutdownTimeout time.Duration
}

const (
	defaultShutdownTimeout = 5 * time.Second
)

func Load() *Config {
	return &Config{
		Client: client.Config{
			LogicConfig: logic.Config{
				Address: "localhost:5001",
			},
		},
		ShutdownTimeout: defaultShutdownTimeout,
	}
}
