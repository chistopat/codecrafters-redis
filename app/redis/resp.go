package redis

import (
	"bufio"
	"fmt"
)

type Resp struct {
	s *bufio.Scanner
}

func NewResp(in *bufio.Reader) Resp {
	scanner := bufio.NewScanner(in)
	scanner.Split(ScanCRLF)
	return Resp{
		s: scanner,
	}
}

func (r *Resp) ParseArray() ([]string, error) {
	n, err := r.GetArrayLen(r.NextToken())
	if err != nil {
		return nil, err
	}
	results := make([]string, 0, n)
	for n != 0 {
		token := r.NextToken()
		if NeedSkip(token) {
			continue
		}
		results = append(results, token)
	}
	return results, nil
}

func (r *Resp) GetArrayLen(token string) (int, error) {
	if len(token) == 2 && token[0] == '*' {
		return int(token[1]-'0') * 2, nil
	}
	return -1, fmt.Errorf("invalid array token: %s", token)
}

func (r *Resp) NextToken() string {
	r.s.Scan()
	return r.s.Text()
}

func NeedSkip(token string) bool {
	if len(token) > 0 && token[0] == '$' {
		return true
	}
	return false
}
