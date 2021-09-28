package main

import (
	"bufio"
	"encoding/json"
	"log"
	"net"
	"syscall"
)

func checkNonCriticalError(err error) {
	if err != nil {
		log.Println("Error", err.Error())
	}
}
func checkError(err error) {
	if err != nil {
		log.Println("Critical error", err.Error())
		syscall.Exit(1)
	}
}

func readGsMsg(conn net.Conn) (gs_msg, error) {
	msg, err := bufio.NewReader(conn).ReadBytes('\n')

	if err != nil {
		return gs_msg{}, err
	}

	gs_msg := gs_msg{}
	err = json.Unmarshal(msg, &gs_msg)
	checkError(err)
	return gs_msg, nil
}

func readGsConfig(conn net.Conn) (gs_config, error) {
	msg, err := bufio.NewReader(conn).ReadBytes('\n')

	if err != nil {
		return gs_config{}, err
	}

	gs_config := gs_config{}
	err = json.Unmarshal(msg, &gs_config)
	checkError(err)
	return gs_config, nil
}

func checkTcpMsgError(err error) {
	if err != nil {
		if err.Error() != "EOF" {
			log.Println("Error", err)
		}
	}
}
