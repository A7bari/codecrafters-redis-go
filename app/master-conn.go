package main

import (
	"fmt"

	"github.com/codecrafters-io/redis-starter-go/app/config"
	"github.com/codecrafters-io/redis-starter-go/app/handlers"
)

func Handshack() error {
	r := config.Get().Master

	r.Send("PING")
	value, _ := r.Read()

	r.Send("REPLCONF", "listening-port", config.Get().Port)
	value, _ = r.Read()

	r.Send("REPLCONF", "capa", "psync2")
	value, _ = r.Read()

	r.Send("PSYNC", "?", "-1")
	value, _ = r.Read()    // +FULLRESYNC 8371b4fb1155b71f4a04d3e1bc3e18c4a990aeeb 0
	value, _ = r.ReadRDB() // +RDBFIILE
	fmt.Printf("Received response from master: %v\n", value.Bulk)

	return nil
}

func ListenOnMaster(errChan chan error) {
	go listening(errChan)
}

func listening(errChan chan error) {
	master := config.Get().Master
	for {
		value, err := master.Read()
		if err != nil {
			errChan <- err
			continue
		}

		if value.Type == "array" && len(value.Array) > 0 {
			handlers.Handle(master.Conn, value.Array)
		} else {
			errChan <- fmt.Errorf("invalid command")
		}
	}
}
