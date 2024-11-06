package respConnection

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"

	"github.com/codecrafters-io/redis-starter-go/app/handlers"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

type RespConn struct {
	Conn     net.Conn
	Reader   *resp.RespReader
	offset   int
	id       string
	mu       sync.Mutex
	AckChans []chan int
	TxQueue  []*resp.RESP
}

func NewRespConn(conn net.Conn) *RespConn {
	fmt.Println("New connection from: ", conn.RemoteAddr().String())
	return &RespConn{
		Conn:     conn,
		Reader:   resp.NewRespReader(bufio.NewReader(conn)),
		offset:   0,
		id:       conn.RemoteAddr().String(),
		mu:       sync.Mutex{},
		AckChans: make([]chan int, 0),
		TxQueue:  nil,
	}
}

func (r *RespConn) Close() {
	r.Conn.Close()
}

func (r *RespConn) Id() string {
	return r.id
}

func (r *RespConn) Listen() {
	for {
		value, err := r.Reader.Read()
		if err != nil {
			break
		}

		if value.Type != "array" || len(value.Array) < 1 {
			break
		}
		r.handleClient(value.Array)
	}

	r.Close()
}

func (r *RespConn) handleClient(args []resp.RESP) error {
	command := strings.ToUpper(args[0].Bulk)

	// handle replication commands
	if command == "REPLCONF" && strings.ToUpper(args[1].Bulk) == "ACK" {
		offset, _ := strconv.Atoi(args[2].Bulk)
		go r.AckRecieved(offset)
		return nil
	}

	if command == "WAIT" {
		r.Write(Wait(args[1:]))
		return nil
	}

	// handle tx commands
	handler := r.GetTxHandler(command)
	if handler == nil {
		// if no tx handler is found, use the default handler
		handler = handlers.GetHandler(command)
	}

	r.Conn.Write(handler(args[1:]))

	if command == "PSYNC" {
		GetReplicaManager().AddReplica(r)
		return nil
	}

	// Propagate the command to all replicas
	if isWriteCommand(command) {
		GetReplicaManager().PropagateCommand(args)
	}

	return nil
}

func isWriteCommand(command string) bool {
	return command == "SET" || command == "DEL"
}

func (r *RespConn) AddOffset(offset int) {
	r.mu.Lock()
	r.offset += offset
	r.mu.Unlock()
}

func (r *RespConn) GetOffset() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.offset
}

func (r *RespConn) Read() (resp.RESP, error) {
	return r.Reader.Read()
}

func (r *RespConn) Write(data []byte) (int, error) {
	return r.Conn.Write(data)
}

func (r *RespConn) ReadRDB() (resp.RESP, error) {
	return r.Reader.ReadRDB()
}

func Wait(params []resp.RESP) []byte {
	fmt.Println("Received WAIT command: ", params)
	count, _ := strconv.Atoi(params[0].Bulk)
	timeout, _ := strconv.Atoi(params[1].Bulk)
	acks := GetReplicaManager().SendAck(timeout, count)

	return resp.Integer(acks).Marshal()
}
