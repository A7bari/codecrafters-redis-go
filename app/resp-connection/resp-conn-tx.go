package respConnection

import (
	"github.com/codecrafters-io/redis-starter-go/app/handlers"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

func (client *RespConn) GetTxHandler(command string) handlers.CommandHandler {
	switch command {
	case "MULTI":
		return client.Multi
	case "EXEC":
		return client.Exec
	default:
		return nil
	}

}

func (client *RespConn) Exec(params []resp.RESP) []byte {
	if client.TxQueue != nil {
		client.TxQueue = nil
		return resp.Array().Marshal()
	}

	return resp.Error("ERR EXEC without MULTI").Marshal()
}

func (client *RespConn) Multi(params []resp.RESP) []byte {
	client.TxQueue = make([][]resp.RESP, 0)
	return resp.String("OK").Marshal()
}
