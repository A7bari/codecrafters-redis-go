package main

import (
	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"github.com/codecrafters-io/redis-starter-go/app/structures"
)

var handlers = map[string]func([]*resp.RESP) *resp.RESP{
	"GET": structures.Get,
	"SET": structures.Set,
}
