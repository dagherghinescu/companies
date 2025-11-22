package repository

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

// Config holds Postgres connection settings loaded from environment
type Config struct {
	User     string `envconfig:"USER" required:"true"`
	Password string `envconfig:"PASSWORD" required:"true"`
	Host     string `envconfig:"HOST" required:"true"`
	Port     string `envconfig:"PORT" default:"5432"`
	Name     string `envconfig:"NAME" required:"true"`
	SSLMode  string `envconfig:"SSLMODE" default:"disable"`
}

// EnvConfig loads the Postgres configuration from environment variables
func EnvConfig() (*Config, error) {
	var cfg Config
	err := envconfig.Process("POSTGRES", &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

// DSN builds the connection string for sql.Open
func (c *Config) DSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.User, c.Password, c.Host, c.Port, c.Name, c.SSLMode,
	)
}
