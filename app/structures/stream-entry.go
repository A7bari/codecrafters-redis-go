package structures

import (
	"fmt"
	"strconv"
	"strings"
)

type Entry struct {
	Key   string
	Pairs map[string]string
}

func NewEntry(timestamp int64, seq int, pairs map[string]string) Entry {
	key := fmt.Sprintf("%d-%d", timestamp, seq)
	return Entry{Key: key, Pairs: pairs}
}

func (e Entry) Timestamp() int64 {
	ids := strings.Split(e.Key, "-")
	timestamp, _ := strconv.ParseInt(ids[0], 10, 64)
	return timestamp
}

func (e Entry) Seq() int {
	ids := strings.Split(e.Key, "-")
	seq, _ := strconv.Atoi(ids[1])
	return seq
}
