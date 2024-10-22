package handlers

import (
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
}

func GetHandler(commands resp.RESP) CommandHandler {
	command := strings.ToUpper(commands.Bulk)
	handler, ok := handlers[command]
	if !ok {
		return notfound
	}

	if isWriteCommand(command) {
		return propagate(handler)
	}

	return handler
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

// propagate is a middleware that will be used to propagate the command to all replicas
func propagate(handler CommandHandler) CommandHandler {
	return func(params []resp.RESP) []byte {
		value := handler(params)
		for _, replica := range config.Get().Replicas {
			replica.Write(resp.Array(params...).Marshal())
		}
		return value
	}
}

func isWriteCommand(command string) bool {
	return command == "SET"
}
