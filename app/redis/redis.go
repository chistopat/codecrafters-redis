package redis

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type MyRedis struct {
	storage map[string]string
	timer   map[string]int64
}

func NewMyRedis() *MyRedis {
	return &MyRedis{
		storage: map[string]string{},
	}
}

func (r *MyRedis) OnConnect(in *bufio.Reader, out *bufio.Writer) {
	for {
		command, err := r.ParseInput(in)
		fmt.Println(command)
		if err != nil {
			fmt.Printf("%v\n", err)
			continue
		}
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
		return r.Set(command[1], command[1:]...)
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

func (r *MyRedis) Set(key string, args ...string) []byte {
	r.storage[key] = args[0]
	if len(args) == 3 {
		duration, _ := strconv.Atoi(args[2])
		r.timer[key] = time.Now().Unix() + int64(duration)
	}
	return []byte("+OK\r\n")
}

func (r *MyRedis) Get(key string) []byte {
	r.CheckTTL(key)
	if value, ok := r.storage[key]; !ok {
		return []byte("(nil)\r\n")
	} else {
		return []byte(fmt.Sprintf("+%s\r\n", value))
	}
}

func (r *MyRedis) CheckTTL(key string) {
	if exp, ok := r.timer[key]; ok {
		if exp <= time.Now().Unix() {
			delete(r.timer, key)
			delete(r.storage, key)
		}
	}
}
