package respConnection

import (
	"fmt"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/handlers"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

func (client *RespConn) GetTxHandler(command string) handlers.CommandHandler {
	switch command {
	case "MULTI":
		return client.Multi
	case "EXEC":
		return client.Exec
	case "DISCARD":
		return client.Discard
	default:
		return nil
	}

}

func (client *RespConn) Exec(params []resp.RESP) []byte {
	if client.TxQueue != nil {

		buf := []byte(fmt.Sprintf("*%d\r\n", len(client.TxQueue)))

		for _, agrs := range client.TxQueue {
			handler := handlers.GetHandler(strings.ToUpper(agrs[0].Bulk))
			handlerResponse := handler(agrs[1:])

			buf = append(buf, handlerResponse...)
		}

		client.TxQueue = nil
		return buf
	}

	return resp.Error("ERR EXEC without MULTI").Marshal()
}

func (client *RespConn) Multi(params []resp.RESP) []byte {
	client.TxQueue = make([][]resp.RESP, 0)
	return resp.String("OK").Marshal()
}

func (client *RespConn) Discard(params []resp.RESP) []byte {
	if client.TxQueue != nil {
		client.TxQueue = nil
		return resp.String("OK").Marshal()
	}

	return resp.Error("ERR DISCARD without MULTI").Marshal()
}
