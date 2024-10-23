package handlers

import (
	"encoding/hex"
	"fmt"
	"strconv"
	"time"

	"github.com/codecrafters-io/redis-starter-go/app/config"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

func replconf(params []resp.RESP) []byte {
	if params[0].Bulk == "GETACK" {
		return resp.Command("REPLCONF", "ACK", strconv.Itoa(config.Get().Offset)).Marshal()
	}

	return resp.String("OK").Marshal()
}

func psync(params []resp.RESP) []byte {
	valid := len(params) > 1 && params[0].Bulk == "?" && params[1].Bulk == "-1"

	if valid {
		msg := resp.String(
			fmt.Sprintf("FULLRESYNC %s 0", config.Get().MasterReplid),
		).Marshal()

		dbfile := getRdbFile()
		msg = append(msg, []byte(fmt.Sprintf("$%d\r\n", len(dbfile)))...)
		msg = append(msg, dbfile...)
		return msg
	}

	return resp.Error("Uncompleted command").Marshal()
}

func getRdbFile() []byte {
	file := "524544495330303131fa0972656469732d76657205372e322e30fa0a72656469732d62697473c040fa056374696d65c26d08bc65fa08757365642d6d656dc2b0c41000fa08616f662d62617365c000fff06e3bfec0ff5aa2"
	data, err := hex.DecodeString(file)
	if err != nil {
		return nil
	}
	return data
}

func wait(params []resp.RESP) []byte {
	count, _ := strconv.Atoi(params[0].Bulk)
	timeout, _ := strconv.Atoi(params[1].Bulk)
	cha := make(chan bool, count)
	ack := 0
	for i := 0; i < len(config.Get().Replicas); i++ {
		rep := config.Get().Replicas[i]

		if rep.GetOffset() > 0 {
			go func(replica *config.Node) {
				size, err := replica.Write(resp.Command("REPLCONF", "GETACK", "*").Marshal())
				if err != nil {
					fmt.Println("err REPLCONF: lost connection " + err.Error())
				}
				replica.AddOffset(size)
				offset := <-replica.AckChan
				fmt.Println("WAIT COMM REPLCONF: ", offset)
				cha <- true
			}(rep)
		} else {
			ack++
		}
	}

loop:
	for ack < count {
		select {
		case <-cha:
			ack++
			fmt.Println("ack: ", ack)
		case <-time.After(time.Duration(timeout) * time.Millisecond):
			break loop
		}
	}

	return resp.Integer(ack).Marshal()

}
