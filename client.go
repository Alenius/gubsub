package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
)

func startClient() {
	conn, err := net.Dial("tcp", "127.0.0.1:8080")

	checkError(err)

	for {

		msg, err := bufio.NewReader(conn).ReadBytes('\n')

		if err != nil {
			if err.Error() != "EOF" {
				log.Println("Error", err)
			}
		}

		if len(msg) > 0 {
			gs_msg := gs_msg{}
			err = json.Unmarshal(msg, &gs_msg)
			checkError(err)
			fmt.Println("res", string(msg))

		}
	}
}
