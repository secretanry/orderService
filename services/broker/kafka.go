package broker

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	kafka_lib "github.com/segmentio/kafka-go"

	"wb-L0/modules/config"
	"wb-L0/modules/kafka"
)

type KafkaBroker struct {
	kafkaConn *kafka.Kafka
	address   string
}

func NewKafkaBroker(kafkaInstance *kafka.Kafka) Broker {
	address := config.GetConfig().KafkaUrl
	if address == "" {
		address = "localhost:9092"
	}
	return &KafkaBroker{
		kafkaConn: kafkaInstance,
		address:   address,
	}
}

func (b *KafkaBroker) StartConsuming(ctx context.Context) chan Message {
	resultChan := make(chan Message, 1)

	go func() {
		defer close(resultChan)
		reader := b.kafkaConn.Reader

		var (
			retryCount    int
			maxRetries    = 10
			retryDelay    = 2 * time.Second
			maxRetryDelay = 30 * time.Second
		)

		for {
			select {
			case <-ctx.Done():
				return
			default:
				msg, err := reader.FetchMessage(ctx)
				if err != nil {
					if ctx.Err() != nil {
						return
					}

					if isCoordinatorError(err) {
						if retryCount < maxRetries {
							retryCount++
							log.Printf("Coordinator unavailable (retry %d/%d), waiting %s",
								retryCount, maxRetries, retryDelay)

							select {
							case <-time.After(retryDelay):
								retryDelay = minDuration(retryDelay*2, maxRetryDelay)
								continue
							case <-ctx.Done():
								return
							}
						} else {
							resultChan <- Message{Err: fmt.Errorf("coordinator unavailable after %d retries", maxRetries)}
							return
						}
					}
					if !isTimedOutError(err) {
						resultChan <- Message{Err: err}
					}
					continue
				}

				retryCount = 0
				retryDelay = 2 * time.Second

				msgCopy := msg
				resultChan <- Message{
					Value: msgCopy.Value,
					Ack: func() error {
						return b.commitWithRetry(ctx, msgCopy)
					},
					Nack: func() error {
						return reader.SetOffset(msgCopy.Offset)
					},
				}
			}
		}
	}()

	return resultChan
}

func isCoordinatorError(err error) bool {
	if kafkaErr, ok := err.(kafka_lib.Error); ok {
		return errors.Is(kafkaErr, kafka_lib.GroupCoordinatorNotAvailable) ||
			errors.Is(kafkaErr, kafka_lib.NotCoordinatorForGroup)
	}
	return false
}

func isTimedOutError(err error) bool {
	if kafkaErr, ok := err.(kafka_lib.Error); ok {
		return errors.Is(kafkaErr, kafka_lib.RequestTimedOut)
	}
	return false
}

func (b *KafkaBroker) commitWithRetry(ctx context.Context, msg kafka_lib.Message) error {
	const maxCommitRetries = 5
	var lastErr error

	for i := 0; i < maxCommitRetries; i++ {
		err := b.kafkaConn.Reader.CommitMessages(ctx, msg)
		if err == nil {
			return nil
		}

		if !isCoordinatorError(err) {
			return err
		}

		lastErr = err
		select {
		case <-time.After(time.Duration(i+1) * time.Second):
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	return fmt.Errorf("commit failed after %d retries: %w", maxCommitRetries, lastErr)
}

func minDuration(a, b time.Duration) time.Duration {
	if a < b {
		return a
	}
	return b
}

// HealthCheck performs a health check on Kafka
func (b *KafkaBroker) HealthCheck(ctx context.Context) error {
	// Try to connect to Kafka and check if it's reachable
	conn, err := kafka_lib.DialContext(ctx, "tcp", b.address)
	if err != nil {
		return fmt.Errorf("failed to connect to Kafka: %w", err)
	}
	defer conn.Close()
	// Check if we can list topics (basic connectivity test)
	_, err = conn.ReadPartitions()
	if err != nil {
		return fmt.Errorf("failed to read partitions: %w", err)
	}
	return nil
}
