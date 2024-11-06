package respConnection

import (
	"time"

	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

type ReplicaManager struct {
	Replicas map[string]*RespConn
}

var replicaManager *ReplicaManager = &ReplicaManager{
	Replicas: make(map[string]*RespConn),
}

func GetReplicaManager() *ReplicaManager {
	return replicaManager
}

func (r *ReplicaManager) AddReplica(conn *RespConn) {
	r.Replicas[conn.Id()] = conn
}

func (r *ReplicaManager) RemoveReplica(id string) {
	delete(r.Replicas, id)
}

func (r *ReplicaManager) GetReplica(id string) *RespConn {
	return r.Replicas[id]
}

func (r *ReplicaManager) GetReplicas() []*RespConn {
	replicas := make([]*RespConn, 0)
	for _, replica := range r.Replicas {
		replicas = append(replicas, replica)
	}
	return replicas
}

func (r *ReplicaManager) SendAck(timeout int, maxCount int) int {
	ackChan := make(chan int)

	count := 0

	for _, replica := range r.Replicas {
		if replica.GetOffset() > 0 {
			go replica.SendAck(ackChan)
		} else {
			count++
		}
	}

loop:
	for count < maxCount {
		select {
		case <-ackChan:
			count++
		case <-time.After(time.Duration(timeout) * time.Millisecond):
			break loop
		}
	}

	close(ackChan)

	return count
}

func (r *ReplicaManager) ClearAckChans(ackChan chan int) {
	for _, replica := range r.Replicas {
		for i := 0; i < len(replica.AckChans); i++ {
			if replica.AckChans[i] == ackChan {
				replica.AckChans = append(replica.AckChans[:i], replica.AckChans[i+1:]...)
				break
			}
		}
	}
}

func (r *ReplicaManager) PropagateCommand(args []resp.RESP) {
	for _, replica := range r.Replicas {
		writtenSize, _ := replica.Write(resp.Array(args...).Marshal())
		replica.AddOffset(writtenSize)
	}
}
