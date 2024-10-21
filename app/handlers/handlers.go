package handlers

import (
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/config"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"github.com/codecrafters-io/redis-starter-go/app/structures"
)

type CommandHandler func([]resp.RESP) resp.RESP

var handlers = map[string]CommandHandler{
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

func GetHandler(command string) CommandHandler {
	handler, ok := handlers[strings.ToUpper(command)]
	if !ok {
		return unfound
	}
	return handler
}

func ping(params []resp.RESP) resp.RESP {

	return resp.RESP{
		Type: "string",
		Bulk: "PONG",
	}
}

func echo(params []resp.RESP) resp.RESP {
	return resp.RESP{
		Type: "bulk",
		Bulk: params[0].Bulk,
	}
}

func unfound(params []resp.RESP) resp.RESP {
	return resp.Error("Command not found")
}
