package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
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

	consumerListener, err := net.Listen("tcp", "127.0.0.1:8080")
	checkError(err)
	producerListener, err := net.Listen("tcp", "127.0.0.1:8081")
	checkError(err)

	channel := make(chan gsMsg)
	log.Println("channel created")

	for {
		if conn, err := consumerListener.Accept(); err == nil {
			go handleConsumerConnection(conn, channel)
		}

		if conn, err := producerListener.Accept(); err == nil {
			go handleProducerConnection(conn, channel)
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

func handleProducerConnection(conn net.Conn, channel chan gsMsg) {
	log.Println("Producer connected")
	for {
		gsMsg, err := readGsMsg(conn)
		checkError(err)

		writeToLedger(conn, gsMsg.Stringify())
		channel <- gsMsg
	}
}

func handleConsumerConnection(conn net.Conn, channel chan gsMsg) {
	defer closeConnection(conn, 0)

	config, err := readConfig(conn)
	log.Println(config)

	checkError(err)

	switch config.ClientType {
	case ClientType(ProducerClient):
		log.Print("Error, this channel is for consumers")
		closeConnection(conn, 1)
	case ClientType(ConsumerClient):
		break
	default:
		panic("no good")
	}

	for {
		msg := <-channel
		log.Println("channel msg", msg)
		writeGsMsg(msg, conn)
	}
}

func closeConnection(conn net.Conn, exitCode int) {
	log.Println("Closing connection")
	close_msg := gsMsg{}.CreateCloseMsg()
	writeGsMsg(close_msg, conn)
	conn.Close()
	syscall.Exit(exitCode)
}
