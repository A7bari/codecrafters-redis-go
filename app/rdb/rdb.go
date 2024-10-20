package rdb

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/codecrafters-io/redis-starter-go/app/structures"
)

type RDB struct {
	dir        string
	dbfilename string
	reader     *bufio.Reader
}

func newRDB(dir, dbfilename string) (*RDB, error) {
	// Open the file
	file, err := os.Open(dir + "/" + dbfilename)
	if err != nil {
		return nil, err
	}

	return &RDB{
		dir:        dir,
		dbfilename: dbfilename,
		reader:     bufio.NewReader(file),
	}, nil
}

func ReadFromRDB(dir, dbfilename string) (structures.RedisMap, error) {
	rdb, err := newRDB(dir, dbfilename)
	if err != nil {
		return nil, err
	}

	return rdb.readKeys()
}

func (r *RDB) readKeys() (structures.RedisMap, error) {
	// Read the header
	header := make([]byte, 9)
	r.reader.Read(header)
	if string(header) != "REDIS0004" {
		return nil, fmt.Errorf("invalid RDB file")
	}

	for {
		// Read the type
		typ, err := r.reader.ReadByte()
		if err != nil {
			break
		}

		switch typ {
		case 0xFE:
			fmt.Println("start reading database info...")
			return r.startDbRead()
		case 0xFF:
			return nil, nil
		}
	}
	return nil, nil
}

func (r *RDB) startDbRead() (structures.RedisMap, error) {
	// read db index
	_, err := r.readSizeEncoded()
	if err != nil {
		return nil, err
	}

	redisMap := structures.RedisMap{}
	currentKey := ""
	for {
		it, err := r.reader.ReadByte()
		if err != nil {
			return nil, err
		}

		switch it {
		case 0xFB: // db info
			_, err := r.readSizeEncoded()
			if err != nil {
				return nil, err
			}

			_, err = r.readSizeEncoded()
			if err != nil {
				return nil, err
			}

		case 0x00: // type string
			key, err := r.readString()
			if err != nil {
				return nil, err
			}

			currentKey = key

			value, err := r.readString()
			if err != nil {
				return nil, err
			}

			redisMap[key] = structures.MapValue{
				Value:  value,
				Expiry: time.Time{},
			}
		case 0xFC: // the curent key has expiry in ms
			timestampBytes := make([]byte, 8)
			r.reader.Read(timestampBytes)
			timestamp := binary.LittleEndian.Uint64(timestampBytes)

			value := redisMap[currentKey]
			value.Expiry = time.Unix(0, int64(timestamp)*int64(time.Millisecond))
			redisMap[currentKey] = value
		case 0xFD: // the curent key has expiry in s
			timestampBytes := make([]byte, 4)
			r.reader.Read(timestampBytes)
			timestamp := binary.LittleEndian.Uint64(timestampBytes)

			value := redisMap[currentKey]
			value.Expiry = time.Unix(0, int64(timestamp)*int64(time.Second))
			redisMap[currentKey] = value
		case 0xFF:
			return redisMap, nil
		}
	}
}

func (r *RDB) readSizeEncoded() (int, error) {
	first, err := r.reader.ReadByte()
	if err != nil {
		return 0, fmt.Errorf("invalid RDB file")
	}

	switch first >> 6 {
	case 0b00:
		return int(first), nil
	case 0b01:
		nextByte, err := r.reader.ReadByte()
		if err != nil {
			return 0, fmt.Errorf("invalid RDB file")
		}
		return int(first&0b00111111)<<8 | int(nextByte), nil
	case 0b10:
		size := make([]byte, 4)
		r.reader.Read(size)

		return int(binary.BigEndian.Uint32(size)), nil
	case 0b11:
		format := int(first&0b00111111) + 1
		sizeBytes := make([]byte, format)
		r.reader.Read(sizeBytes)

		strSize := string(sizeBytes)
		return strconv.Atoi(strSize)
	}

	return 0, nil
}

func (r *RDB) readString() (string, error) {
	size, err := r.readSizeEncoded()
	if err != nil {
		return "", err
	}

	buf := make([]byte, size)
	r.reader.Read(buf)
	return string(buf), nil
}
