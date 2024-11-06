package respConnection

import "github.com/codecrafters-io/redis-starter-go/app/resp"

func (replica *RespConn) SendAck(ack chan int) (int, error) {
	replica.mu.Lock()
	replica.AckChans = append(replica.AckChans, ack)
	replica.mu.Unlock()

	s, err := replica.Conn.Write(
		resp.Command("REPLCONF", "GETACK", "*").Marshal(),
	)
	replica.AddOffset(s)
	return s, err
}

func (replica *RespConn) AckRecieved(offset int) {
	replica.mu.Lock()
	defer replica.mu.Unlock()

	cha := replica.AckChans[0]
	replica.AckChans = replica.AckChans[1:]
	cha <- offset
}
