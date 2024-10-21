package handlers

import (
	"fmt"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/config"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

func info(params []resp.RESP) []byte {
	if len(params) < 1 {
		return resp.Error("ERR wrong number of arguments for 'info' command").Marshal()
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
		return resp.Bulk(msg).Marshal()
	}

	return resp.Nil().Marshal()
}
