package main

import (
	"bufio"
	"fmt"
	"net"
)

func startClient() {
	conn, err := net.Dial("tcp", "127.0.0.1:8080")

	checkError(err)

	status, err := bufio.NewReader(conn).ReadString('\n')
	checkError(err)

	fmt.Println(status)

}
