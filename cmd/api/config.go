package main

import (
	"time"

	"github.com/anvidev/goenv"
)

type config struct {
	server   serverConfig
	database databaseConfig
}

type serverConfig struct {
	env          string
	addr         string
	readTimeout  time.Duration
	writeTimeout time.Duration
	idleTimeout  time.Duration
}

type databaseConfig struct {
	token string
	url   string
}

func loadConfig() config {
	return config{
		server: serverConfig{
			env:          goenv.String("SERVER_ENV", "development"),
			addr:         goenv.String("SERVER_ADDR", ":9090"),
			readTimeout:  goenv.Duration("SERVER_READ_TIMEOUT", time.Second*10),
			writeTimeout: goenv.Duration("SERVER_WRITE_TIMEOUT", time.Second*30),
			idleTimeout:  goenv.Duration("SERVER_IDLE_TIMEOUT", time.Minute),
		},
		database: databaseConfig{
			token: goenv.MustString("TURSO_AUTH_TOKEN"),
			url:   goenv.MustString("TURSO_DATABASE_URL"),
		},
	}
}
