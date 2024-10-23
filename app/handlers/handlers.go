package handlers

import (
	"fmt"
	"net"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/config"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"github.com/codecrafters-io/redis-starter-go/app/structures"
)

type CommandHandler func([]resp.RESP) []byte

var handlers = map[string]func([]resp.RESP) []byte{
	"GET":      structures.Get,
	"SET":      structures.Set,
	"CONFIG":   config.GetConfigHandler,
	"PING":     ping,
	"ECHO":     echo,
	"KEYS":     structures.Keys,
	"INFO":     info,
	"PSYNC":    psync,
	"REPLCONF": replconf,
	"WAIT":     wait,
}

func Handle(conn net.Conn, args []resp.RESP) error {
	command := strings.ToUpper(args[0].Bulk)
	handler, ok := handlers[command]
	if !ok {
		handler = notfound
	}

	if command == "REPLCONF" && strings.ToUpper(args[1].Bulk) == "ACK" {
		fmt.Println("warning: recieved replica ack from client handler ")
	}

	conn.Write(handler(args[1:]))

	if command == "PSYNC" {
		config.AddReplicat(conn)
		return fmt.Errorf("PSYNC")
	}

	// Propagate the command to all replicas
	if isWriteCommand(command) {
		for _, replica := range config.Get().Replicas {
			go func(replica config.Node) {
				writtenSize, _ := replica.Write(resp.Array(args...).Marshal())
				config.SetReplOffset(replica, writtenSize)
			}(replica)
		}
	}

	return nil
}

func HandleMaster(conn net.Conn, args []resp.RESP) {
	command := strings.ToUpper(args[0].Bulk)
	handler, ok := handlers[command]
	if !ok {
		handler = notfound
	}

	data := handler(args[1:])
	if command == "REPLCONF" && strings.ToUpper(args[1].Bulk) == "GETACK" {
		conn.Write(data)
	}
}

func ping(params []resp.RESP) []byte {
	return resp.String("PONG").Marshal()
}

func echo(params []resp.RESP) []byte {
	return resp.String(params[0].Bulk).Marshal()
}

func notfound(params []resp.RESP) []byte {
	return resp.Error("Command not found").Marshal()
}

func isWriteCommand(command string) bool {
	return command == "SET"
}
