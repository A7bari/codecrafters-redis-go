package main

import (
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

func handshack(replica string) error {

	conn, err := net.Dial("tcp", replica)
	if err != nil {
		return err
	}
	defer conn.Close()

	// send ping
	ping, _ := resp.RESP{
		Type: "array",
		Array: []resp.RESP{
			{
				Type: "bulk",
				Bulk: "PING",
			},
		},
	}.Marshal()

	conn.Write(ping)
	return nil
}
