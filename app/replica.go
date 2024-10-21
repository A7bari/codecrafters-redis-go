package main

import (
	"fmt"
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/config"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

func handshack() error {

	send("PING")

	send("REPLCONF", "listening-port", config.Get("port"))

	send("REPLCONF", "capa", "psync2")

	return nil
}

func send(commands ...string) error {
	fmt.Printf("Sending command to master: %v\n", commands)
	link := config.Get("master_host") + ":" + config.Get("master_port")
	conn, err := net.Dial("tcp", link)
	if err != nil {
		return err
	}
	defer conn.Close()

	// send ping
	comArray := make([]resp.RESP, len(commands))
	for i, command := range commands {
		comArray[i] = resp.Bulk(command)
	}

	msg, err := resp.Array(comArray).Marshal()

	if err != nil {
		return err
	}

	conn.Write(msg)
	return nil
}
