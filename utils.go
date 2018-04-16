package llrp

import (
	"bytes"
	"encoding/binary"
)

type paramFunc func(messageType int, params ...interface{}) []interface{}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func pack(data []interface{}) []byte {
	buf := new(bytes.Buffer)
	for _, v := range data {
		err := binary.Write(buf, binary.BigEndian, v)
		check(err)
	}
	return buf.Bytes()
}
