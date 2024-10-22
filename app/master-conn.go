package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/config"
	"github.com/codecrafters-io/redis-starter-go/app/handlers"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

type Master struct {
	masterAddress string
	conn          net.Conn
	reader        *resp.RespReader
}

func NewMaster() (*Master, error) {
	adress := config.Get("master_host") + ":" + config.Get("master_port")
	conn, err := net.Dial("tcp", adress)
	if err != nil {
		return nil, err
	}
	return &Master{
		masterAddress: adress,
		conn:          conn,
		reader:        resp.NewRespReader(bufio.NewReader(conn)),
	}, nil
}

func (r *Master) Handshack() error {
	r.send("PING")
	value, _ := r.reader.Read()

	r.send("REPLCONF", "listening-port", config.Get("port"))
	value, _ = r.reader.Read()

	r.send("REPLCONF", "capa", "psync2")
	value, _ = r.reader.Read()

	r.send("PSYNC", "?", "-1")
	value, _ = r.reader.Read()    // +FULLRESYNC 8371b4fb1155b71f4a04d3e1bc3e18c4a990aeeb 0
	value, _ = r.reader.ReadRDB() // +RDBFIILE
	fmt.Printf("Received response from master: %v\n", value.Bulk)

	return nil
}

func (r *Master) Listen(errChan chan error) {
	go r.listening(errChan)
}

func (r *Master) listening(errChan chan error) {
	for {
		value, err := r.reader.Read()
		if err != nil {
			errChan <- err
			continue
		}

		if value.Type == "array" && len(value.Array) > 0 {
			command := strings.ToUpper(value.Array[0].Bulk)
			handler := handlers.GetHandler(command)
			r.conn.Write(handler(value.Array[1:]))
		} else {
			errChan <- fmt.Errorf("invalid command")
		}
	}
}

func (r *Master) Close() {
	r.conn.Close()
}

func (r *Master) send(commands ...string) error {
	comArray := make([]resp.RESP, len(commands))
	for i, command := range commands {
		comArray[i] = resp.Bulk(command)
	}

	msg := resp.Array(comArray...).Marshal()

	r.conn.Write(msg)
	return nil
}
