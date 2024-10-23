package handlers

import (
	"fmt"
	"net"
	"strconv"
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
		offset, _ := strconv.Atoi(args[2].Bulk)
		fmt.Println("offset: ", offset)
		config.Replica(conn).AckChan <- offset
	}

	conn.Write(handler(args[1:]))

	if command == "PSYNC" {
		config.AddReplicat(conn)
		return fmt.Errorf("PSYNC")
	}

	// Propagate the command to all replicas
	if isWriteCommand(command) {
		go func() {
			for i := 0; i < len(config.Get().Replicas); i++ {
				replica := config.Get().Replicas[i]
				writtenSize, _ := replica.Write(resp.Array(args...).Marshal())
				// if err != nil {
				// 	// disconnected
				// 	config.RemoveReplica(replica)
				// 	i--
				// 	return
				// }
				replica.AddOffset(writtenSize)
			}
		}()
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
