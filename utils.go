package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"log"
	"net"
	"os"
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

func writeGsMsg(gs_msg gsMsg, conn net.Conn) {
	serializedMsg, err := json.Marshal(gs_msg)
	msgWithNewline := append(serializedMsg, []byte{'\n'}...)
	checkError(err)

	_, err = conn.Write(msgWithNewline)
	checkError(err)
}

func readGsMsg(conn net.Conn) (gsMsg, error) {
	msg, err := bufio.NewReader(conn).ReadBytes('\n')

	if err != nil {
		return gsMsg{}, err
	}

	gs_msg := gsMsg{}
	err = json.Unmarshal(msg, &gs_msg)
	checkError(err)
	return gs_msg, nil
}

func readGsConfig(conn net.Conn) (gsConfig, error) {
	msg, err := bufio.NewReader(conn).ReadBytes('\n')

	if err != nil {
		return gsConfig{}, err
	}

	gs_config := gsConfig{}
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

func FileExists(fileName string) (bool, error) {
	_, err := os.Stat(fileName)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return false, err
}
