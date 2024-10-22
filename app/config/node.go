package config

import (
	"bufio"
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

type Node struct {
	Conn   net.Conn
	Reader *resp.RespReader
	Offset int
}

func NewNode(conn net.Conn) Node {
	return Node{
		Conn:   conn,
		Reader: resp.NewRespReader(bufio.NewReader(conn)),
		Offset: 0,
	}
}

func (r Node) Close() {
	r.Conn.Close()
}

func (r Node) Read() (resp.RESP, error) {
	return r.Reader.Read()
}

func (r Node) ReadRDB() (resp.RESP, error) {
	return r.Reader.ReadRDB()
}

func (r Node) Write(data []byte) (int, error) {
	return r.Conn.Write(data)
}
