package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
)

type ConnectionType string

const (
	Client ConnectionType = "Client"
	Server ConnectionType = "Server"
)

func failInvalidConnectionType(ct *ConnectionType) error {
	switch *ct {
	case Client, Server:
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

	if connType == Client {
		startClient()
	} else {
		startServer()
	}

}

type gs_msg struct {
	Id        uuid.UUID `json:"id"`
	Msg       string    `json:"msg"`
	Timestamp string    `json:"timestamp"`
}

func (gs_msg) Create(msg string) gs_msg {
	id := uuid.New()

	timestamp := time.Now().UTC().Format(time.RFC3339)

	return gs_msg{Id: id, Timestamp: timestamp, Msg: msg}
}

func (gs_msg gs_msg) GetMsg() string {
	return gs_msg.Msg
}

func (gs_msg gs_msg) Stringify() string {
	return gs_msg.Timestamp + " " + gs_msg.Id.String() + " " + gs_msg.Msg
}

type gs_config struct {
	Id         uuid.UUID  `json:"id"`
	ClientType ClientType `json:"type"`
}

type ClientType string

const (
	Subscriber ClientType = "Subscriber"
	Publisher  ClientType = "Publisher"
)

func failInvalidClientType(ct *ClientType) error {
	switch *ct {
	case Subscriber, Publisher:
		return nil
	}

	return errors.New("invalid client type")
}

func (gs_config) Create(client_type ClientType) (gs_config, error) {
	id := uuid.New()

	validation_err := failInvalidClientType(&client_type)

	if validation_err != nil {
		return gs_config{}, validation_err
	}

	return gs_config{Id: id, ClientType: client_type}, nil
}
