package main

import (
	"errors"
	"fmt"
	"log"
	"os"
)

type ConnectionType string

const (
	Consumer ConnectionType = "Consumer"
	Producer ConnectionType = "Producer"
	Broker   ConnectionType = "Broker"
)

func failInvalidConnectionType(ct *ConnectionType) error {
	switch *ct {
	case Consumer, Producer, Broker:
		return nil
	}

	return errors.New("invalid connection type")
}

func main() {
	connTypeUnchecked := os.Getenv("TYPE")
	connType := ConnectionType(connTypeUnchecked)
	typeError := failInvalidConnectionType(&connType)

	if typeError != nil {
		log.Println(fmt.Sprintf("Invalid connection type: %s. Must be 'Server' or 'Client'", connTypeUnchecked))
		return
	}

	if connType == Consumer {
		startConsumer()
	} else if connType == Producer {
		startProducer()
	} else {
		startBroker()
	}
}

type ClientType string

const (
	ConsumerClient ClientType = "Consumer"
	ProducerClient ClientType = "Producer"
)

func failInvalidClientType(ct *ClientType) error {
	switch *ct {
	case ConsumerClient, ProducerClient:
		return nil
	}

	return errors.New("invalid client type")
}
