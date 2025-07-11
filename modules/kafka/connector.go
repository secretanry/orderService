package kafka

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/segmentio/kafka-go"

	"wb-L0/modules/config"
)

type Kafka struct {
	Reader *kafka.Reader
}

func (k *Kafka) Init(_ chan error) error {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        []string{config.GetConfig().KafkaUrl},
		Topic:          config.GetConfig().KafkaTopic,
		GroupID:        config.GetConfig().KafkaConsumerGroup,
		MinBytes:       1,
		MaxBytes:       10e6,
		MaxWait:        2 * time.Second,
		CommitInterval: 0,
		MaxAttempts:    5,
		QueueCapacity:  100,
		Logger:         kafka.LoggerFunc(logKafka),
		ErrorLogger:    kafka.LoggerFunc(logKafkaError),
		StartOffset:    kafka.LastOffset,
	})
	k.Reader = reader
	return nil
}

func (k *Kafka) SuccessfulMessage() string {
	return "Kafka successfully initialized"
}

func (k *Kafka) Shutdown(_ context.Context) error {
	err := k.Reader.Close()
	return err
}

func logKafka(msg string, args ...interface{}) {
	if config.GetConfig().RunMode == "debug" &&
		!strings.Contains(msg, "no messages received from kafka within the allocated time") {
		log.Printf("[KAFKA] "+msg, args...)
	}
}

func logKafkaError(msg string, args ...interface{}) {
	log.Printf("[KAFKA-ERROR] "+msg, args...)
}
