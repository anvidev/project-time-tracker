package main

import (
	"time"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Resend   ResendConfig
}

type ServerConfig struct {
	Version      string        `goenv:"SERVER_VERSION,default=v0.1.0"`
	Env          string        `goenv:"SERVER_ENV,default=development"`
	Addr         string        `goenv:"SERVER_ADDR,default=:9090"`
	ReadTimeout  time.Duration `goenv:"SERVER_READ_TIMEOUT,default=10s"`
	WriteTimeout time.Duration `goenv:"SERVER_WRITE_TIMEOUT,default=30s"`
	IdleTimeout  time.Duration `goenv:"SERVER_IDLE_TIMEOUT,default=1m"`
}

type DatabaseConfig struct {
	Token string `goenv:"TURSO_AUTH_TOKEN,required"`
	URL   string `goenv:"TURSO_DATABASE_URL,required"`
}

type ResendConfig struct {
	From   string `goenv:"RESEND_FROM,default=Tidsregistrering <noreply@nemunivers.app>"` // format: "name <email>"
	ApiKey string `goenv:"RESEND_API_KEY,required"`
}
