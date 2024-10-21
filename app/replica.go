package main

import (
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/config"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

func handshack(link string) error {

	send("PING")

	send("REPLCONF", "listening-port", config.Get("port"))

	send("REPLCONF", "capa", "psync2")

	return nil
}

func send(commands ...string) error {
	link := config.Get("master_host") + ":" + config.Get("master_port")
	conn, err := net.Dial("tcp", link)
	if err != nil {
		return err
	}
	defer conn.Close()

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
