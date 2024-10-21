package main

import (
	"bufio"
	"fmt"
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/config"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

func handshack() error {
	link := config.Get("master_host") + ":" + config.Get("master_port")
	conn, err := net.Dial("tcp", link)
	if err != nil {
		return err
	}
	// defer conn.Close()

	reader := resp.NewRespReader(bufio.NewReader(conn))

	send(conn, "PING")
	value, _ := reader.Read()
	fmt.Printf("Received response from master: %v\n", value.Bulk)

	send(conn, "REPLCONF", "listening-port", config.Get("port"))
	value, _ = reader.Read()
	fmt.Printf("Received response from master: %v\n", value.Bulk)

	send(conn, "REPLCONF", "capa", "psync2")
	value, _ = reader.Read()
	fmt.Printf("Received response from master: %v\n", value.Bulk)

	return nil
}

func send(conn net.Conn, commands ...string) error {
	fmt.Printf("Sending command to master: %v\n", commands)

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
