package structures

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Stream struct {
	Entries       map[int64][]Entry
	size          int
	lastTimestamp int64
}

func NewStream() *Stream {
	return &Stream{
		Entries:       map[int64][]Entry{},
		size:          0,
		lastTimestamp: -1,
	}
}

func (s *Stream) Add(key string, pairs map[string]string) error {
	tmstmp, strSeq, err := parseKey(key)
	if err != nil {
		return err
	}

	timestamp, seq, err := s.formatKey(tmstmp, strSeq)
	if err != nil {
		return err
	}

	err = s.validateKey(timestamp, seq)

	if err != nil {
		return err
	}

	if _, ok := s.Entries[timestamp]; !ok {
		s.Entries[timestamp] = []Entry{}
	}

	s.Entries[timestamp] = append(s.Entries[timestamp], NewEntry(timestamp, seq, pairs))

	s.size++
	if timestamp > s.lastTimestamp {
		s.lastTimestamp = timestamp
	}

	return nil
}

func (s *Stream) Get(key string) (map[string]string, bool) {
	timestamp, _, err := parseKey(key)
	if err != nil {
		return nil, false
	}

	for _, e := range s.Entries[timestamp] {
		if e.Key == key {
			return e.Pairs, true
		}
	}
	return nil, false
}

func (s *Stream) validateKey(timestamp int64, seq int) error {

	if timestamp < 0 || (timestamp == 0 && seq <= 0) {
		return fmt.Errorf("ERR The ID specified in XADD must be greater than 0-0")
	}

	if s.size == 0 {
		return nil
	}

	if _, ok := s.Entries[timestamp]; !ok {
		if (timestamp < s.lastTimestamp) || (timestamp == s.lastTimestamp && seq <= s.LastSeq(timestamp)) {
			return fmt.Errorf("ERR The ID specified in XADD is equal or smaller than the target stream top item")
		}
	}

	return nil
}

func (s *Stream) formatKey(timestamp int64, strSeq string) (int64, int, error) {
	if timestamp < 0 {
		timestamp := time.Now().UnixMilli()
		return timestamp, 0, nil
	}

	if strSeq == "*" {
		lastSeq := s.LastSeq(timestamp)
		if timestamp == 0 && lastSeq == -1 {
			lastSeq = 0
		}

		return timestamp, lastSeq + 1, nil
	}

	seq, err := strconv.Atoi(strSeq)
	if err != nil {
		return 0, 0, fmt.Errorf("ERR invalid stream ID format: seq must be an integer")
	}

	return timestamp, seq, nil
}

func parseKey(key string) (int64, string, error) {
	if key == "*" {
		return -1, "*", nil
	}

	ids := strings.Split(key, "-")
	if len(ids) != 2 {
		return 0, "", fmt.Errorf("ERR invalid stream ID format: <timestamp>-<seq>")
	}

	timestamp, err := strconv.ParseInt(ids[0], 10, 64)
	if err != nil {
		return 0, "", fmt.Errorf("ERR invalid stream ID format: timestamp must be an integer")
	}

	return timestamp, ids[1], nil
}

func (s *Stream) LastSeq(timestamp int64) int {
	entries, ok := s.Entries[timestamp]
	if !ok {
		return -1
	}

	return entries[len(entries)-1].Seq()
}

func (s *Stream) Len() int {
	return s.size
}
