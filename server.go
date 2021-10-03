package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
)

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

func readConfig(conn net.Conn) (gs_config, error) {
	for {
		msg, err := bufio.NewReader(conn).ReadBytes('\n')
		checkTcpMsgError(err)

		if len(msg) > 0 {
			gs_config := gs_config{}
			err = json.Unmarshal(msg, &gs_config)
			checkError(err)
			fmt.Println("res", string(msg))
			return gs_config, nil
		}
	}
}

func writeToLedger(msg string) error {
	err := os.WriteFile("ledger.txt", []byte(msg), 0644)
	return err
}

func handleProducerConnection(conn net.Conn) {
	msg := gs_msg{}.Create("hello")
	writeToLedger(msg.Stringify())

}

func handleConnection(conn net.Conn) {
	defer closeConnection(conn)

	log.Println("Client connected")

	config, err := readConfig(conn)
	log.Println(config)

	checkError(err)

	switch config.ClientType {
	case ClientType(Publisher):
		handleProducerConnection(conn)
	default:
		panic("no good")
	}
}

func closeConnection(conn net.Conn) {
	log.Println("Closing connection")
	conn.Close()
}
