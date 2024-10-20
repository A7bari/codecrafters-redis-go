package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/config"
	"github.com/codecrafters-io/redis-starter-go/app/handlers"
	"github.com/codecrafters-io/redis-starter-go/app/rdb"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"github.com/codecrafters-io/redis-starter-go/app/structures"
)

func main() {
	// read configs flag
	dir := flag.String("dir", "", "Directory to serve static files from")
	dbfilename := flag.String("dbfilename", "dump.rdb", "Filename to save the DB to")
	port := flag.String("port", "6379", "Port to listen on")
	replicaof := flag.String("replicaof", "", "Replicate to another Redis server")
	flag.Parse()

	config.Set("dir", *dir)
	config.Set("dbfilename", *dbfilename)
	config.Set("port", *port)
	config.Set("replicaof", *replicaof)
	if *replicaof != "" {
		config.Set("role", "slave")
	} else {
		config.Set("role", "master")
	}
	config.Set("master_replid", "8371b4fb1155b71f4a04d3e1bc3e18c4a990aeeb")
	config.Set("master_repl_offset", "0")

	initializeMapStore(*dir, *dbfilename)

	l, err := net.Listen("tcp", "0.0.0.0:"+*port)
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
	defer conn.Close()
	fmt.Println("Received connection: ", conn.RemoteAddr().String())

	reader := resp.NewRespReader(bufio.NewReader(conn))
	for {
		resp, err := reader.Read()
		if err != nil {
			fmt.Println("Error reading from connection: ", err.Error())
			break
		}

		if resp.Type != "array" || len(resp.Array) < 1 {
			fmt.Println("Invalid command")
			break
		}

		command := strings.ToUpper(resp.Array[0].Bulk)

		handler, ok := handlers.GetHandler(command)
		if !ok {
			fmt.Println("Unknown command: ", command)
			break
		}
		res, err := handler(resp.Array[1:]).Marshal()
		if err != nil {
			fmt.Println("Error marshaling response for command:", command, "-", err.Error())
			continue
		}
		conn.Write(res)
	}
}

func initializeMapStore(dir string, dbfilename string) {
	//load rdb
	mapStore, err := rdb.ReadFromRDB(dir, dbfilename)
	if err != nil {
		fmt.Println("Error loading RDB: ", err.Error())
	} else {
		structures.LoadKeys(mapStore)
	}
}
