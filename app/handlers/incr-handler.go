package handlers

import (
	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"github.com/codecrafters-io/redis-starter-go/app/structures"
)

func Incr(params []resp.RESP) []byte {
	if len(params) != 1 {
		return resp.Error("ERR wrong number of arguments for 'incr' command").Marshal()
	}

	val := params[0].Bulk

	intVal, err := structures.Incr(val)
	if err != nil {
		return resp.Error("ERR value is not an integer or out of range").Marshal()
	}

	return resp.Integer(intVal).Marshal()
}
