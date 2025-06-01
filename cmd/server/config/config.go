package config

import (
	"go-keeper/internal/server"
	"go-keeper/internal/server/data/database"
	"go-keeper/pkg/jwtfactory"
	"os"
	"time"
)

type Config struct {
	Server          server.Config
	DB              database.Config
	JWTConfig       jwtfactory.Config
	ShutdownTimeout time.Duration
}

const (
	dbConnectionStringEnv = "DATABASE_URI"
	serverAddressEnv      = "SERVER_ADDRESS"

	defaultShutdownTimeout    = 5 * time.Second
	serverAddressDefaultValue = ":8080"
)

var defaultRetryAttempts = []time.Duration{time.Second, 3 * time.Second, 5 * time.Second}

func Load() (*Config, error) {

	dbConnectionString := ""
	serverAddress := serverAddressDefaultValue

	if valStr, ok := os.LookupEnv(dbConnectionStringEnv); ok {
		dbConnectionString = valStr
	}

	if valStr, ok := os.LookupEnv(serverAddressEnv); ok {
		serverAddress = valStr
	}

	return &Config{
		Server: server.Config{
			ServerAddress:   serverAddress,
			ShutdownTimeout: defaultShutdownTimeout,
		},
		DB: database.Config{
			ConnectionString:   dbConnectionString,
			RetryAttemptDelays: defaultRetryAttempts,
		},
		JWTConfig: jwtfactory.Config{
			Algorithm:      "HS256",
			Secret:         "secret",
			ExpirationTime: time.Hour,
		},
		ShutdownTimeout: defaultShutdownTimeout,
	}, nil
}
