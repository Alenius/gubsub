package main

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/google/uuid"
)

type gsConfig struct {
	Id         uuid.UUID  `json:"id"`
	ClientType ClientType `json:"type"`
	Topic      string     `json:"topic"`
}

func (gsConfig) Create(client_type ClientType, topic string) (gsConfig, error) {
	id := uuid.New()

	validation_err := failInvalidClientType(&client_type)

	if validation_err != nil {
		return gsConfig{}, validation_err
	}

	return gsConfig{Id: id, ClientType: client_type, Topic: topic}, nil
}

func (gsConfig) CreateFromFile(configFilePath string) (gsConfig, error) {
	id := uuid.New()
	file, err := os.Open(configFilePath)
	checkError(err)
	defer file.Close()

	configBytes, _ := ioutil.ReadAll(file)

	base_config := gsConfig{}
	json.Unmarshal(configBytes, &base_config)

	return gsConfig{Id: id, ClientType: base_config.ClientType, Topic: base_config.Topic}, nil
}
