package respConnection

import (
	"fmt"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/config"
	"github.com/codecrafters-io/redis-starter-go/app/handlers"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

func (master *RespConn) Handshack() error {

	master.Write(resp.Command("PING").Marshal())
	master.Read()

	master.Write(resp.Command("REPLCONF", "listening-port", config.Get().Port).Marshal())
	master.Read()

	master.Write(resp.Command("REPLCONF", "capa", "psync2").Marshal())
	master.Read()

	master.Write(resp.Command("PSYNC", "?", "-1").Marshal())
	master.Read()    // +FULLRESYNC 8371b4fb1155b71f4a04d3e1bc3e18c4a990aeeb 0
	master.ReadRDB() // +RDBFIILE

	return nil
}

func (master *RespConn) handleMaster(args []resp.RESP) {
	command := strings.ToUpper(args[0].Bulk)
	handler := handlers.GetHandler(command)

	data := handler(args[1:])
	if command == "REPLCONF" && strings.ToUpper(args[1].Bulk) == "GETACK" {
		master.Write(data)
	}
}

func (master *RespConn) ListenOnMaster(errChan chan error) {
	for {
		value, err := master.Read()
		if err != nil {
			errChan <- err
			continue
		}

		if value.Type == "array" && len(value.Array) > 0 {
			master.handleMaster(value.Array)
		} else {
			errChan <- fmt.Errorf("invalid command")
		}

		config.IncOffset(
			len(value.Marshal()),
		)
	}
}
