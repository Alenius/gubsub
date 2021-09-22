package main

import (
	"errors"
	"fmt"
	"log"
	"os"
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

}




type gs_msg struct {
	Id  string `json:"id"`
	Msg string `json:"msg"`
}
