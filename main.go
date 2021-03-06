package main

import (
	"errors"
	"flag"
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

	log.Println("args", os.Args[1])

	config := createConfig(configFilePath)

	switch config.Type {
	case Consumer:
		startConsumer(config)
	case Producer:
		startProducer(config)
	case Broker:
		startBroker()
	default:
		log.Fatal("That config type is not recognized: ", config.ClientType)
	}
}
