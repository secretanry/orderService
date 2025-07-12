package broker

import (
	"context"
)

var brokerInstance Broker

type Message struct {
	Value []byte
	Err   error
	Ack   func() error
	Nack  func() error
}

type Broker interface {
	StartConsuming(context.Context) chan Message
	HealthCheck(context.Context) error
}

func SetBroker(broker Broker) {
	brokerInstance = broker
}

func GetBroker() Broker {
	return brokerInstance
}
