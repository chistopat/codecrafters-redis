package redis

import "bufio"

type MyRedis struct {
}

func NewMyRedis() *MyRedis {
	return &MyRedis{}
}

func (r *MyRedis) OnConnect(in *bufio.Reader, out *bufio.Writer) {
	for {
		command, err := r.parseInput(in)
		if err != nil {
			continue
		}
		response := r.Invoke(command)
		r.Send(response, out)
	}
}

func (r *MyRedis) parseInput(in *bufio.Reader) ([]string, error) {
	scanner := bufio.NewScanner(in)
	return nil, nil
}

func (r *MyRedis) Send(data []byte, out *bufio.Writer) {
	_, err := out.Write(data)
	if err != nil {
		err := out.Flush()
		if err != nil {
			return
		}
	}
}

func (r *MyRedis) Invoke(command []string) []byte {
	return nil
}
