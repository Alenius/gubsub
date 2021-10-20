package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
)

func sendConfig(conn net.Conn, config gsConfig) {

	serializedMsg, err := json.Marshal(config)
	msgWithNewline := append(serializedMsg, []byte{'\n'}...)

	checkError(err)

	_, err = conn.Write(msgWithNewline)
	checkError(err)
}

func checkIfCloseMsg(gs_msg gsMsg) bool {
	return gs_msg.Type == "CLOSE"
}

func startConsumer(config gsConfig) {
	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	checkError(err)

	sendConfig(conn, config)

	for {
		msg, err := bufio.NewReader(conn).ReadBytes('\n')
		checkTcpMsgError(err)

		if len(msg) > 0 {
			gs_msg := gsMsg{}
			err = json.Unmarshal(msg, &gs_msg)
			checkError(err)

			shouldCloseConn := checkIfCloseMsg(gs_msg)

			if shouldCloseConn {
				conn.Close()
				log.Println("Closing client")
				return
			}

			fmt.Println("res", string(msg))

		}
	}
}
