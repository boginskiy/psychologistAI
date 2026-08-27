package server

import "time"

// Конфигурация сервера
type Config struct {
	Port            string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	IdleTimeout     time.Duration
	ShutdownTimeout time.Duration
}

var ConfigServer Config = Config{
	Port:            ":8080",
	ReadTimeout:     5 * time.Second,
	WriteTimeout:    10 * time.Second,
	IdleTimeout:     120 * time.Second,
	ShutdownTimeout: 5 * time.Second,
}
