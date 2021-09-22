package main

import (
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

func handleConnection(conn net.Conn) {
	log.Println("Client connected")

	defer closeConnection(conn)

	raw_msg := gs_msg{}.Create("Hello world, this is a longer string")
	serializedMsg, err := json.Marshal(raw_msg)

	checkError(err)

	size, err := conn.Write(serializedMsg)
	checkError(err)
	fmt.Println("size", size)
}

func closeConnection(conn net.Conn) {
	log.Println("Closing connection")
	conn.Close()
}
