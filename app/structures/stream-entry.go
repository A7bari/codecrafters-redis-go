package structures

import (
	"fmt"
)

type Entry struct {
	Pairs     map[string]string
	seq       int
	timestamp int64
}

func NewEntry(timestamp int64, seq int, pairs map[string]string) Entry {
	return Entry{Pairs: pairs, seq: seq, timestamp: timestamp}
}

func (e Entry) Timestamp() int64 {
	return e.timestamp
}

func (e Entry) Seq() int {
	return e.seq
}

func (e Entry) Key() string {
	return fmt.Sprintf("%d-%d", e.timestamp, e.seq)
}
