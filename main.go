package main

import (
	"bufio"
	"fmt"
	"net"
)

const (
	SERVER_HOST = "localhost"
	SERVER_PORT = "3000"
	SERVER_TYPE = "tcp"
)

func main() {
	fmt.Println("server started")

	ln, _ := net.Listen("tcp", ":3000")

	conn, _ := ln.Accept()

	for {
		// get message, output
		message, _ := bufio.NewReader(conn).ReadString('\n')
		fmt.Print("Message Received:", string(message))
	}
}
