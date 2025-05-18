package server

import "time"

type Config struct {
	ServerAddress   string
	ShutdownTimeout time.Duration
}
