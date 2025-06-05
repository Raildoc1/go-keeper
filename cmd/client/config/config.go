package config

import (
	"os"
	"time"
)

type Config struct {
	ServerAddress    string
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
		ServerAddress:    serverAddress,
		ShutdownTimeout:  defaultShutdownTimeout,
		LocalStoragePath: "./storage.str",
	}
}
