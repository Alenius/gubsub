package main

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/google/uuid"
)

type gsConfig struct {
	Id    uuid.UUID      `json:"id"`
	Type  ConnectionType `json:"type"`
	Topic string         `json:"topic"`
}

type gsConsumerConfig struct {
	gsConfig
	StartTimestamp string `json:"startTimestamp"`
}

func (gsConfig) Create(connection_type ConnectionType, topic string) (gsConfig, error) {
	id := uuid.New()

	validation_err := failInvalidConnectionType(&connection_type)

	if validation_err != nil {
		return gsConfig{}, validation_err
	}

	return gsConfig{Id: id, Type: connection_type, Topic: topic}, nil
}

func (gsConfig) CreateFromFile(configFilePath string) (gsConfig, error) {
	id := uuid.New()
	file, err := os.Open(configFilePath)
	checkError(err)
	defer file.Close()

	configBytes, _ := ioutil.ReadAll(file)

	base_config := gsConfig{}
	json.Unmarshal(configBytes, &base_config)

	return gsConfig{Id: id, Type: base_config.Type, Topic: base_config.Topic}, nil
}
