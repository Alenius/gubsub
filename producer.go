package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

func getInput() string {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("-> ")
		text, _ := reader.ReadString('\n')
		// convert CRLF to LF
		text = strings.Replace(text, "\n", "", -1)
		return text
	}
}

func startProducer(config gsConfig) {
	conn, err := net.Dial("tcp", "127.0.0.1:8081")
	checkError(err)

	sendConfig(conn, config)

	for {
		msg := getInput()

		isExitMsg := strings.Compare(msg, "__EXIT")
		if isExitMsg == 0 {
			return
		}
		gsMsg := gsMsg{}.Create(msg, config.Topic)
		writeGsMsg(gsMsg, conn)
	}
}

func closeProducer(conn net.Conn) {
	log.Println("Closing producer")
	conn.Close()
}
