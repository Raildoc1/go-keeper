package config

import (
	client "go-keeper/internal/client/config"
	logic "go-keeper/internal/client/logic/config"
	"os"
	"time"
)

type Config struct {
	Client           client.Config
	ShutdownTimeout  time.Duration
	LocalStoragePath string
}

const (
	serverAddressEnv = "SERVER_ADDRESS"

	defaultShutdownTimeout    = 5 * time.Second
	serverAddressDefaultValue = "localhost:8080"
)

func Load() *Config {

	serverAddress := serverAddressDefaultValue

	if valStr, ok := os.LookupEnv(serverAddressEnv); ok {
		serverAddress = valStr
	}

	return &Config{
		Client: client.Config{
			LogicConfig: logic.Config{
				Address: serverAddress,
			},
		},
		ShutdownTimeout:  defaultShutdownTimeout,
		LocalStoragePath: "./storage.str",
	}
}
