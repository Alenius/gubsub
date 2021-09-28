package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
)

func sendConfig(conn net.Conn) {
	raw_config, err := gs_config{}.Create(ClientType(Subscriber))

	checkError(err)
	serializedMsg, err := json.Marshal(raw_config)
	log.Println(serializedMsg)

	checkError(err)

	_, err = conn.Write(serializedMsg)
	checkError(err)
	_, err = conn.Write([]byte{'\n'})
	checkError(err)
}

func startClient() {
	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	checkError(err)

	sendConfig(conn)

	for {
		msg, err := bufio.NewReader(conn).ReadBytes('\n')
		checkTcpMsgError(err)

		if len(msg) > 0 {
			gs_msg := gs_msg{}
			err = json.Unmarshal(msg, &gs_msg)
			checkError(err)
			fmt.Println("res", string(msg))

		}
	}
}
