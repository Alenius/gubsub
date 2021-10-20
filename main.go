package main

import (
	"errors"
	"flag"
	"log"
)

type ConnectionType string

const (
	Consumer ConnectionType = "Consumer"
	Producer ConnectionType = "Producer"
	Broker   ConnectionType = "Broker"
)

type ClientType string

const (
	ConsumerClient ClientType = "Consumer"
	ProducerClient ClientType = "Producer"
	BrokerClient   ClientType = "Broker"
)

func failInvalidClientType(ct *ClientType) error {
	switch *ct {
	case ConsumerClient, ProducerClient:
		return nil
	}

	return errors.New("invalid client type")
}

func readCliFlags() string {
	configFilePath := flag.String("configFilePath", "", "Path to the config file")

	flag.Parse()

	return *configFilePath
}

func createConfig(configFilePath string) gsConfig {
	if configFilePath == "" {
		log.Fatal("No config path provided")
	}

	config, _ := gsConfig{}.CreateFromFile(configFilePath)
	log.Println("Config provided: ", config)

	return config
}

func main() {
	configFilePath := readCliFlags()

	config := createConfig(configFilePath)

	switch config.ClientType {
	case ConsumerClient:
		startConsumer(config)
	case ProducerClient:
		startProducer(config)
	case BrokerClient:
		startBroker()
	default:
		log.Fatal("That config type is not recognized: ", config.ClientType)
	}
}
