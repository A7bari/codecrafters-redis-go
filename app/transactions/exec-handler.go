package transactions

import "github.com/codecrafters-io/redis-starter-go/app/resp"

func Exec(params []resp.RESP) []byte {
	return resp.Error("ERR EXEC without MULTI").Marshal()
}
