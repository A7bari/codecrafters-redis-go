package main

import (
	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"github.com/codecrafters-io/redis-starter-go/app/structures"
)

var handlers = map[string]func([]*resp.RESP) *resp.RESP{
	"GET":  structures.Get,
	"SET":  structures.Set,
	"PING": ping,
	"ECHO": echo,
}

func ping(params []*resp.RESP) *resp.RESP {
	return &resp.RESP{
		Type: "string",
		Bulk: "PONG",
	}
}

func echo(params []*resp.RESP) *resp.RESP {
	return &resp.RESP{
		Type: "bulk",
		Bulk: params[0].Bulk,
	}
}
