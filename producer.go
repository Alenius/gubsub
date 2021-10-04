package main

import (
	"bufio"
	"encoding/json"
	"net"
)

func startProducer() {
	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	checkError(err)

	sendConfig(conn, ClientType(ProducerClient))

	for {
		raw_config, err := gsConfig{}.Create(ClientType(ProducerClient))

		checkError(err)
		serializedMsg, err := json.Marshal(raw_config)
		msgWithNewline := append(serializedMsg, []byte{'\n'}...)
		_, err = conn.Write(msgWithNewline)
		msg, err := bufio.NewReader(conn).ReadBytes('\n')
		checkTcpMsgError(err)

	}
}
