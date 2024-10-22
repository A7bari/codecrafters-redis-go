package config

import (
	"bufio"
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

type Node struct {
	address string
	Conn    net.Conn
	Reader  *resp.RespReader
}

func NewNode(address string) (*Node, error) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}
	return &Node{
		address: address,
		Conn:    conn,
		Reader:  resp.NewRespReader(bufio.NewReader(conn)),
	}, nil
}

func (r *Node) Close() {
	r.Conn.Close()
}

func (r *Node) Send(commands ...string) error {
	comArray := make([]resp.RESP, len(commands))
	for i, command := range commands {
		comArray[i] = resp.Bulk(command)
	}

	msg := resp.Array(comArray...).Marshal()

	r.Conn.Write(msg)
	return nil
}

func (r *Node) Read() (resp.RESP, error) {
	return r.Reader.Read()
}

func (r *Node) ReadRDB() (resp.RESP, error) {
	return r.Reader.ReadRDB()
}

func (r *Node) Write(data []byte) error {
	r.Conn.Write(data)
	return nil
}
