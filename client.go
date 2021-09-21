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

		status, err := bufio.NewReader(conn).ReadBytes(' ')

		if err != nil {
			if err.Error() != "EOF" {
				log.Println("Error", err)
			}
		}

		if len(status) > 0 {
			gs_msg := gs_msg{}
			err = json.Unmarshal(status, &gs_msg)
			fmt.Println("res", string(status))
			fmt.Println("Status", status)
			// fmt.Printf("gs msg %+v \n", gs_msg)
			fmt.Println("err", err)
		}
	}

}
