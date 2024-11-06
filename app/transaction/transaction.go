package transaction

import "github.com/codecrafters-io/redis-starter-go/app/resp"

type TransQueue struct {
	Commands []resp.RESP
}
