package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
)

func sendConfig(conn net.Conn) {
	raw_config, err := gs_config{}.Create(ClientType(Publisher))

	checkError(err)
	serializedMsg, err := json.Marshal(raw_config)
	msgWithNewline := append(serializedMsg, []byte{'\n'}...)

	checkError(err)

	_, err = conn.Write(msgWithNewline)
	checkError(err)
}

func checkIfCloseMsg(gs_msg gs_msg) bool {
	return gs_msg.Type == "CLOSE"
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
