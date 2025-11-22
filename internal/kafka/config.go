package kafka

import "github.com/kelseyhightower/envconfig"

// Config holds the kafka variables
type Config struct {
	Broker string `envconfig:"BROKER" required:"true"`
	Topic  string `envconfig:"TOPIC" required:"true"`
}

// EnvConfig loads the Kafka configuration from environment variables
func EnvConfig() (*Config, error) {
	var cfg Config
	err := envconfig.Process("KAFKA", &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
