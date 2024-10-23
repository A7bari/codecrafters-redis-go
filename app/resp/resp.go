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
	Array   []RESP
}

func NewRespReader(r *bufio.Reader) *RespReader {
	return &RespReader{reader: r}
}

func Command(cmd string, args ...string) RESP {
	res := make([]RESP, len(args)+1)
	res[0] = Bulk(cmd)
	for i, arg := range args {
		res[i+1] = Bulk(arg)
	}
	return Array(res...)
}

func Error(msg string) RESP {
	return RESP{
		Type: "error",
		Bulk: msg,
	}
}

func Integer(i int) RESP {
	return RESP{
		Type:    "integer",
		Integer: i,
	}
}

func String(s string) RESP {
	return RESP{
		Type: "string",
		Bulk: s,
	}
}

func Array(a ...RESP) RESP {
	return RESP{
		Type:  "array",
		Array: a,
	}
}

func Nil() RESP {
	return RESP{
		Type: "nil",
	}
}

func Bulk(b string) RESP {
	return RESP{
		Type: "bulk",
		Bulk: b,
	}
}

func (r *RespReader) Read() (RESP, error) {
	typ, err := r.reader.ReadByte()
	if err != nil {
		return RESP{}, err
	}

	switch typ {
	case ArrayByte:
		return r.readArray()
	case BulkByte:
		return r.readBulk()
	case StringByte:
		return r.readString()
	case IntegerByte:
		i, err := r.readInt()
		if err != nil {
			return RESP{}, err
		}
		return Integer(i), nil

	}
	return RESP{}, fmt.Errorf("unsupported type: %c", typ)
}

func (r *RespReader) ReadRDB() (RESP, error) {
	typ, _ := r.reader.ReadByte()
	if typ != BulkByte {
		return RESP{}, fmt.Errorf("expected '$' as first byte for RDB")
	}
	// read size
	size, err := r.readInt()

	if err != nil {
		return RESP{}, err
	}

	buf := make([]byte, size)
	r.reader.Read(buf)

	return RESP{
		Type: "rdb",
		Bulk: string(buf),
	}, nil
}

func (r *RespReader) readArray() (RESP, error) {
	size, err := r.readInt()
	if err != nil {
		return RESP{}, err
	}

	resuts := make([]RESP, size)
	for i := 0; i < size; i++ {
		resuts[i], err = r.Read()
		if err != nil {
			return RESP{}, err
		}
	}
	return RESP{
		Type:  "array",
		Array: resuts,
	}, nil
}

func (r *RespReader) readBulk() (RESP, error) {
	_, err := r.readInt()
	if err != nil {
		return RESP{}, err
	}

	buf, err := r.readLine()
	if err != nil {
		return RESP{}, err
	}

	return RESP{
		Type: "bulk",
		Bulk: string(buf),
	}, nil
}

func (r *RespReader) readString() (RESP, error) {
	buf, err := r.readLine()
	if err != nil {
		return RESP{}, err
	}

	return String(string(buf)), nil
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

func (r RESP) Marshal() []byte {
	switch r.Type {
	case "bulk":
		msg := fmt.Sprintf("$%d\r\n%s\r\n", len(r.Bulk), r.Bulk)
		return []byte(msg)
	case "integer":
		msg := fmt.Sprintf(":%d\r\n", r.Integer)
		return []byte(msg)
	case "string":
		msg := fmt.Sprintf("+%s\r\n", r.Bulk)
		return []byte(msg)
	case "array":
		buf := []byte(fmt.Sprintf("*%d\r\n", len(r.Array)))
		for _, item := range r.Array {
			buf = append(buf, item.Marshal()...)
		}

		return buf
	case "error":
		msg := fmt.Sprintf("-%s\r\n", r.Bulk)
		return []byte(msg)
	case "nil":
		return []byte("$-1\r\n")
	case "rdb":
		return []byte(fmt.Sprintf("$%d\r\n%s", len(r.Bulk), r.Bulk))
	}

	return nil
}
