package main

import (
	"fmt"
	"net"
	"os"

	"github.com/codecrafters-io/redis-starter-go/app/config"
	"github.com/codecrafters-io/redis-starter-go/app/rdb"
	respConnection "github.com/codecrafters-io/redis-starter-go/app/resp-connection"
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

		client := respConnection.NewRespConn(conn)
		go client.Listen()
	}
}

func setup() error {
	//initile the map if rdb file is found
	initializeMapStore()
	// if handle the replica if it is a slave
	if config.Get().Role == "slave" {
		masterHost := config.Get().MasterHost
		masterPort := config.Get().MasterPort
		masterConn, err := net.Dial("tcp", masterHost+":"+masterPort)
		if err != nil {
			fmt.Println("Error connecting to master: ", err.Error())
			return err
		}

		master := respConnection.NewRespConn(masterConn)
		master.Handshack()
		errChan := make(chan error)
		go master.ListenOnMaster(errChan)

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
