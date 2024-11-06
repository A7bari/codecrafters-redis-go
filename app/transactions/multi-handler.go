package transactions

import "github.com/codecrafters-io/redis-starter-go/app/resp"

func Multi(params []resp.RESP) []byte {
	return resp.String("OK").Marshal()
}
