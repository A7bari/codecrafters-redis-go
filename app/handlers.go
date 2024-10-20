package main

import (
	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"github.com/codecrafters-io/redis-starter-go/app/structures"
)

var Handlers = map[string]func([]resp.RESP) resp.RESP{
	"GET":    structures.Get,
	"SET":    structures.Set,
	"CONFIG": GetConfig,
	"PING":   ping,
	"ECHO":   echo,
	"KEYS":   structures.Keys,
	"INFO":   info,
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

func info(params []resp.RESP) resp.RESP {
	if len(params) < 1 {
		return resp.RESP{
			Type: "error",
			Bulk: "ERR wrong number of arguments for 'info' command",
		}
	}

	if params[0].Bulk == "replication" {
		return resp.RESP{
			Type: "bulk",
			Bulk: "role:master",
		}
	}

	return resp.RESP{
		Type: "nil",
	}
}
