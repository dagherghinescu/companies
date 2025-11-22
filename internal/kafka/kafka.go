package kafka

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	kafka "github.com/segmentio/kafka-go"
)

// ProducerInterface defines the methods your app needs
type ProducerInterface interface {
	Publish(ctx context.Context, key string, value any) error
	Close() error
}

// Producer wraps a Kafka writer
type Producer struct {
	writer *kafka.Writer
}

// NewProducer creates a new Kafka producer using KafkaConfig
func NewProducer(cfg *Config) *Producer {
	brokers := strings.Split(cfg.Broker, ",")
	writer := &kafka.Writer{
		Addr:         kafka.TCP(brokers...),
		Topic:        cfg.Topic,
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: kafka.RequireAll,
	}
	return &Producer{writer: writer}
}

// Close shuts down the Kafka producer, releasing all resources.
// It should be called when the producer is no longer needed.
func (p *Producer) Close() error {
	return p.writer.Close()
}

// Publish sends a message to the Kafka topic configured in the producer.
// The `key` is used for message partitioning, and `value` is marshaled to JSON.
// Returns an error if marshaling fails or the message could not be written.
func (p *Producer) Publish(ctx context.Context, key string, value any) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	msg := kafka.Message{
		Key:   []byte(key),
		Value: data,
		Time:  time.Now(),
	}
	return p.writer.WriteMessages(ctx, msg)
}
