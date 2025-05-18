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

	defaultShutdownTimeout = 5 * time.Second
)

var defaultRetryAttempts = []time.Duration{time.Second, 3 * time.Second, 5 * time.Second}

func Load() (*Config, error) {

	dbConnectionString := ""

	if valStr, ok := os.LookupEnv(dbConnectionStringEnv); ok {
		dbConnectionString = valStr
	}

	return &Config{
		Server: server.Config{
			ServerAddress:   ":5001",
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
