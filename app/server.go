package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

const Ping = "PING"

func main() {
	fmt.Println("My simple redis started!")

	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}
	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}

	for {
		message, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			continue
		}
		if string(message) == Ping {
			_, err = conn.Write([]byte("PONG" + "\n"))
			if err != nil {
				continue
			}
		}
	}
}
