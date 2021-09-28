package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
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
		log.Println("hey")
		msg, err := bufio.NewReader(conn).ReadBytes('\n')
		log.Println("hoy")
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

func handleConnection(conn net.Conn) {
	defer closeConnection(conn)

	log.Println("Client connected")

	config, err := readConfig(conn)
	log.Println(config)

	checkError(err)

	for {
		msg, _ := json.Marshal("oiqjwd\n")
		conn.Write(msg)
	}

}

func closeConnection(conn net.Conn) {
	log.Println("Closing connection")
	conn.Close()
}
