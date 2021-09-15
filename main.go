package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
)

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

func startServer() {
	log.Println("starting server")

	listener, err := net.Listen("tcp", "127.0.0.1:8080")

	checkError(err)

	for {
		if conn, err := listener.Accept(); err == nil {
			go handleConnection(conn)
		}
	}

}

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

	return errors.New("Invalid connection type")
}

func handleConnection(conn net.Conn) {
	log.Println("Client connected")

	defer conn.Close()

	messageProto := gs_msg{id: "123", msg: "Hello world!"}
	serializedMsg, err := json.Marshal(messageProto)

	checkError(err)

	conn.Write(serializedMsg)
}

type gs_msg struct {
	id  string
	msg string
}
