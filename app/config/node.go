package config

import (
	"bufio"
	"net"
	"sync"

	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

type Node struct {
	Conn   net.Conn
	Reader *resp.RespReader
	offset int
	id     string
	mu     sync.Mutex
}

func NewNode(conn net.Conn) *Node {
	return &Node{
		Conn:   conn,
		Reader: resp.NewRespReader(bufio.NewReader(conn)),
		offset: 0,
		id:     conn.RemoteAddr().String(),
		mu:     sync.Mutex{},
	}
}

func (r *Node) Close() {
	r.Conn.Close()
}

func (r *Node) AddOffset(offset int) {
	r.mu.Lock()
	r.offset += offset
	r.mu.Unlock()
}

func (r *Node) GetOffset() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.offset
}

func (r *Node) Read() (resp.RESP, error) {
	return r.Reader.Read()
}

func (r *Node) ReadRDB() (resp.RESP, error) {
	return r.Reader.ReadRDB()
}

func (r *Node) Write(data []byte) (int, error) {
	return r.Conn.Write(data)
}
