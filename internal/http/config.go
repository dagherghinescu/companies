package http

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

// Config holds setting for the HTTP server.
type Config struct {
	Addr              string        `envconfig:"ADDR" default:":8080"`
	ReadHeaderTimeout time.Duration `envconfig:"READ_HEADER_TIMEOUT" default:"5s"`
	ReadTimeout       time.Duration `envconfig:"READ_TIMEOUT" default:"10s"`
	WriteTimeout      time.Duration `envconfig:"WRITE_TIMEOUT" default:"10s"`
}

// EnvConfig loads config from environment variables into HTTPConfig.
func EnvConfig() (*Config, error) {
	var cfg Config
	if err := envconfig.Process("HTTP", &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
