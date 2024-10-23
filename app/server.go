package main

import (
	"bufio"
	"fmt"
	"net"
	"os"

	"github.com/codecrafters-io/redis-starter-go/app/config"
	"github.com/codecrafters-io/redis-starter-go/app/handlers"
	"github.com/codecrafters-io/redis-starter-go/app/rdb"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"github.com/codecrafters-io/redis-starter-go/app/structures"
)

func main() {

	conf := config.Get()

	setup()

	l, err := net.Listen("tcp", "0.0.0.0:"+conf.Port)
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}
	defer l.Close()

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
		value, err := reader.Read()
		if err != nil {
			fmt.Println("Error reading from connection: ", err.Error())
			break
		}

		if value.Type != "array" || len(value.Array) < 1 {
			fmt.Println("Invalid command")
			break
		}

		handlers.Handle(conn, value.Array)
	}

	conn.Close()
}

func setup() error {
	//initile the map if rdb file is found
	initializeMapStore()
	// if handle the replica if itis a slave
	if config.Get().Role == "slave" {
		Handshack()
		errChan := make(chan error)
		ListenOnMaster(errChan)
		go func() {
			err := <-errChan
			fmt.Println("Error reading from master: ", err.Error())
		}()
	}

	return nil
}

func initializeMapStore() {
	mapStore, err := rdb.ReadFromRDB(config.Get().Dir, config.Get().Dbfilename)
	if err != nil {
		fmt.Println("Error loading RDB: ", err.Error())
	} else {
		structures.LoadKeys(mapStore)
	}
}
