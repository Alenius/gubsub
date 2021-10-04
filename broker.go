package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"syscall"
)

func checkServerError(conn net.Conn, err error) {
	if err != nil {
		log.Println("Critical error", err.Error())
		closeConnection(conn, 1)
	}
}

func startBroker() {
	log.Println("starting broker")

	listener, err := net.Listen("tcp", "127.0.0.1:8080")

	checkError(err)

	for {
		if conn, err := listener.Accept(); err == nil {
			go handleConnection(conn)
		}
	}

}

func readConfig(conn net.Conn) (gsConfig, error) {
	for {
		msg, err := bufio.NewReader(conn).ReadBytes('\n')
		checkTcpMsgError(err)

		if len(msg) > 0 {
			gs_config := gsConfig{}
			err = json.Unmarshal(msg, &gs_config)
			checkError(err)
			fmt.Println("res", string(msg))
			return gs_config, nil
		}
	}
}

func writeToLedger(conn net.Conn, msg string) error {
	path := "ledger.txt"

	fileExists, _ := FileExists(path)

	if !fileExists {
		os.Create(path)
	}

	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0644)
	checkError(err)
	defer f.Close()

	bytes, err := f.Write([]byte(msg))
	checkServerError(conn, err)
	_, err = f.Write([]byte{'\n'})
	checkServerError(conn, err)
	f.Sync()

	log.Println("Wrote to file: ", bytes)
	return err
}

func handleProducerConnection(conn net.Conn) {
	for {
		msg_raw := getInput()

		isExitMsg := strings.Compare(msg_raw, "__EXIT")
		if isExitMsg == 0 {
			return
		}
		readGsMsg(conn)

		msg := gsMsg{}.Create(msg_raw)
		writeToLedger(conn, msg.Stringify())
	}
}

func handleConnection(conn net.Conn) {
	defer closeConnection(conn, 0)

	log.Println("Client connected")

	config, err := readConfig(conn)
	log.Println(config)

	checkError(err)

	switch config.ClientType {
	case ClientType(ProducerClient):
		handleProducerConnection(conn)
	default:
		panic("no good")
	}
}

func closeConnection(conn net.Conn, exitCode int) {
	log.Println("Closing connection")
	close_msg := gsMsg{}.CreateCloseMsg()
	writeGsMsg(close_msg, conn)
	conn.Close()
	syscall.Exit(exitCode)
}