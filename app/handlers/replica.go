package handlers

import (
	"fmt"

	"github.com/codecrafters-io/redis-starter-go/app/config"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

func replconf(params []resp.RESP) resp.RESP {
	return resp.String("OK")
}

func psync(params []resp.RESP) resp.RESP {
	valid := len(params) > 1 && params[0].Bulk == "?" && params[1].Bulk == "-1"

	if valid {
		msg := fmt.Sprintf("FULLRESYNC %s 0\r\n", config.Get("master_replid"))
		return resp.String(msg)
	}

	return resp.Error("Uncompleted command")
}
