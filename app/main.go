package main

import (
	"bufio"
	"fmt"
	"os"
)

const PORT = 6379

type MyRedis struct{}

func NewMyRedis() *MyRedis {
	return &MyRedis{}
}

func (r *MyRedis) OnConnect(in *bufio.Reader, out *bufio.Writer) {
	for {
		_, _ = in.ReadString('\n')
		_, _ = out.Write([]byte("+PING\r\n"))
		out.Flush()
	}
}

func main() {
	fmt.Println("My simple redis started!")
	cache := NewMyRedis()
	server, err := NewNetworkServer("0.0.0.0", PORT, cache.OnConnect)

	if err != nil {
		fmt.Println("Failed to start network server")
		os.Exit(1)
	}
	server.Run()
}
