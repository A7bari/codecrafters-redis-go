package main

import (
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/config"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

func handshack(link string) error {

	conn, err := net.Dial("tcp", link)
	if err != nil {
		return err
	}
	defer conn.Close()

	send(conn, "PING")

	send(conn, "REPLCONF", "listening-port", config.Get("port"))

	send(conn, "REPLCONF", "capa", "eof")

	return nil
}

func send(conn net.Conn, commands ...string) error {
	// send ping
	comArray := make([]resp.RESP, len(commands))
	for i, command := range commands {
		comArray[i] = resp.RESP{
			Type: "bulk",
			Bulk: command,
		}
	}

	msg, err := resp.RESP{
		Type:  "array",
		Array: comArray,
	}.Marshal()

	if err != nil {
		return err
	}

	conn.Write(msg)
	return nil
}
