package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

// Ensures gofmt doesn't remove the "net" and "os" imports in stage 1 (feel free to remove this!)
var _ = net.Listen
var _ = os.Exit

func main() {

	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			break
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	reader := resp.NewRespReader(bufio.NewReader(conn))
	for {
		resp, err := reader.Read()
		if err != nil {
			fmt.Println("Error reading from connection: ", err.Error())
			break
		}

		if resp.Type != "array" {
			fmt.Println("expected to recieve an array")
			break
		}

		command := resp.Array[0].Bulk

		if strings.ToUpper(command) == "ECHO" {
			respBuf, _ := resp.Array[1].Marshal()
			conn.Write(respBuf)
		} else {
			conn.Write([]byte("+PONG\r\n"))
		}
	}
}
