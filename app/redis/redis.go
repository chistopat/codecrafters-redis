package redis

import "bufio"

type MyRedis struct{}

func NewMyRedis() *MyRedis {
	return &MyRedis{}
}

func (r *MyRedis) OnConnect(in *bufio.Reader, out *bufio.Writer) {
	for {
		_, _ = in.ReadString('\n')
		_, _ = out.Write([]byte("+PONG\r\n"))
		out.Flush()
	}
}
