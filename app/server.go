package main

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"os"
)

func main() {
	fmt.Println("My simple redis started!")

	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		go Serve(conn)
	}
}

func Serve(conn net.Conn) {
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			fmt.Println("Alarm! can't close connection: ", err.Error())
			os.Exit(1)
		}
	}(conn)

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)
	Handle(reader, writer)
}

func Handle(in *bufio.Reader, out *bufio.Writer) {
	scanner := bufio.NewScanner(in)
	scanner.Split(ScanCRLF)
	for scanner.Scan() {
		token := scanner.Text()
		//row = strings.TrimSpace(row)
		//if row != Ping {
		//	continue
		//}
		fmt.Println(token)
		_, err := out.Write([]byte("+PONG\r\n"))
		if err != nil {
			return
		}
		out.Flush()
	}
}

// https://stackoverflow.com/questions/37530451/golang-bufio-read-multiline-until-crlf-r-n-delimiter
// dropCR drops a terminal \r from the data.
func dropCR(data []byte) []byte {
	if len(data) > 0 && data[len(data)-1] == '\r' {
		return data[0 : len(data)-1]
	}
	return data
}

func ScanCRLF(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if i := bytes.Index(data, []byte{'\r', '\n'}); i >= 0 {
		// We have a full newline-terminated line.
		return i + 2, dropCR(data[0:i]), nil
	}
	// If we're at EOF, we have a final, non-terminated line. Return it.
	if atEOF {
		return len(data), dropCR(data), nil
	}
	// Request more data.
	return 0, nil, nil
}
