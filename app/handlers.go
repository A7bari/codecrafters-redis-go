package main

import (
	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"github.com/codecrafters-io/redis-starter-go/app/structures"
)

func GetHandler(command string) (func([]resp.RESP) resp.RESP, bool) {
	switch command {
	case "GET":
		return structures.Get, true
	case "SET":
		return structures.Set, true
	case "PING":
		return ping, true
	case "ECHO":
		return echo, true
	}
	return nil, false
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
