package resp

import (
	"bufio"
	"fmt"
	"strconv"
)

const (
	ArrayByte   = byte('*')
	BulkByte    = byte('$')
	IntegerByte = byte(':')
	StringByte  = byte('+')
)

type RespReader struct {
	reader *bufio.Reader
}

type RESP struct {
	Type    string
	Bulk    string
	Integer int
	Array   []*RESP
}

func NewRespReader(r *bufio.Reader) *RespReader {
	return &RespReader{reader: r}
}

func (r *RespReader) Read() (*RESP, error) {
	typ, _ := r.reader.ReadByte()

	switch typ {
	case ArrayByte:
		return r.readArray()
	case BulkByte:
		return r.readBulk()
	}
	return nil, nil
}

func (r *RespReader) readArray() (*RESP, error) {
	size, err := r.readInt()
	if err != nil {
		return nil, err
	}

	resuts := make([]*RESP, size)
	for i := 0; i < size; i++ {
		resuts[i], err = r.Read()
		if err != nil {
			return nil, err
		}
	}
	return &RESP{
		Type:  "array",
		Array: resuts,
	}, nil
}

func (r *RespReader) readBulk() (*RESP, error) {
	_, err := r.readInt()
	if err != nil {
		return nil, err
	}

	buf, err := r.readLine()
	if err != nil {
		return nil, err
	}

	return &RESP{
		Type: "bulk",
		Bulk: string(buf),
	}, nil
}
func (r *RespReader) readLine() ([]byte, error) {
	line, err := r.reader.ReadBytes(byte('\n'))
	if err != nil {
		return line, fmt.Errorf("err reading line")
	}

	return line[:len(line)-2], nil
}

func (r *RespReader) readInt() (int, error) {
	line, err := r.readLine()
	if err != nil {
		return 0, err
	}

	return strconv.Atoi(string(line))
}

func (r *RESP) Marshal() ([]byte, error) {
	switch r.Type {
	case "bulk":
		msg := fmt.Sprintf("$%d\r\n%s\r\n", len(r.Bulk), r.Bulk)
		return []byte(msg), nil
	case "integer":
		msg := fmt.Sprintf(":%d\r\n", r.Integer)
		return []byte(msg), nil
	case "string":
		msg := fmt.Sprintf("+%s\r\n", r.Bulk)
		return []byte(msg), nil
	case "array":
		buf := []byte(fmt.Sprintf("*%d\r\n", len(r.Array)))
		for _, item := range r.Array {
			m, err := item.Marshal()
			if err != nil {
				return nil, err
			}
			buf = append(buf, m...)
		}

		return buf, nil
	}

	return nil, fmt.Errorf("unsopported type")
}
