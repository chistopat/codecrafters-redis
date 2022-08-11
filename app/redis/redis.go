package redis

import (
	"bufio"
	"fmt"
	"log"
	"strings"
)

type MyRedis struct {
	storage map[string]string
}

func NewMyRedis() *MyRedis {
	return &MyRedis{
		storage: map[string]string{},
	}
}

func (r *MyRedis) OnConnect(in *bufio.Reader, out *bufio.Writer) {
	for {
		command, err := r.ParseInput(in)
		if err != nil {
			log.Printf("%v\n", err)
			continue
		}
		log.Println(command)
		response := r.Invoke(command)
		r.Send(response, out)
	}
}

func (r *MyRedis) ParseInput(in *bufio.Reader) ([]string, error) {
	parser := NewResp(in)
	return parser.ParseArray()
}

func (r *MyRedis) Send(data []byte, out *bufio.Writer) {
	out.Write(data)
	out.Flush()
}

func (r *MyRedis) Invoke(command []string) []byte {
	method := command[0]
	switch strings.ToUpper(method) {
	case "PING":
		return r.Ping()
	case "ECHO":
		return r.Echo([]byte(command[1]))
	case "SET":
		return r.Set(command[1], command[2])
	case "GET":
		return r.Get(command[1])
	default:
		return r.Ping()
	}
}

func (r *MyRedis) Echo(data []byte) []byte {
	return []byte(fmt.Sprintf("+%s\r\n", data))
}

func (r *MyRedis) Ping() []byte {
	return []byte(fmt.Sprintf("+PONG\r\n"))
}

func (r *MyRedis) Set(key string, value string) []byte {
	r.storage[key] = value
	return []byte("+OK\r\n")
}

func (r *MyRedis) Get(key string) []byte {
	if value, ok := r.storage[key]; !ok {
		return []byte("(nil)\r\n")
	} else {
		return []byte(fmt.Sprintf("%s\r\n", value))
	}
}
