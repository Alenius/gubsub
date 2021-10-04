package main

import "github.com/google/uuid"

type gsConfig struct {
	Id         uuid.UUID  `json:"id"`
	ClientType ClientType `json:"type"`
}

func (gsConfig) Create(client_type ClientType) (gsConfig, error) {
	id := uuid.New()

	validation_err := failInvalidClientType(&client_type)

	if validation_err != nil {
		return gsConfig{}, validation_err
	}

	return gsConfig{Id: id, ClientType: client_type}, nil
}
