package main

import (
	"fmt"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/config"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"github.com/codecrafters-io/redis-starter-go/app/structures"
)

var Handlers = map[string]func([]resp.RESP) resp.RESP{
	"GET":    structures.Get,
	"SET":    structures.Set,
	"CONFIG": config.GetConfigHandler,
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

	if strings.ToUpper(params[0].Bulk) == "REPLICATION" {
		role := config.Get("role")
		master_replid := config.Get("master_replid")
		master_repl_offset := config.Get("master_repl_offset")

		msg := fmt.Sprintf(
			"role:%s\nmaster_replid:%s\nmaster_repl_offset:%s",
			role,
			master_replid,
			master_repl_offset,
		)
		return resp.RESP{
			Type: "bulk",
			Bulk: msg,
		}
	}

	return resp.RESP{
		Type: "nil",
	}
}
