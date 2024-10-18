package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

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

		command := strings.ToUpper(resp.Array[0].Bulk)

		if command == "ECHO" {
			respBuf, _ := resp.Array[1].Marshal()
			conn.Write(respBuf)
		} else {
			handler, ok := handlers[command]
			if !ok {
				fmt.Println("Unknown command: ", command)
				break
			}
			res, _ := handler(resp.Array[1:]).Marshal()
			conn.Write(res)
		}
	}
}
