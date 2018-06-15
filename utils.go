package llrp

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type paramFunc func(messageType int, params ...interface{}) []interface{}

func convertBooleanUint8(b bool) uint8 {
	if b {
		return 0x80
	} else {
		return 0
	}
}

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
func convert16uintbit(param ...bool) uint16 {
	var res uint16
	for _, k := range param {
		res = res << 1
		if k {
			res += 1
		}
	}
	l := 16 - len(param)
	for l > 0 {
		res = res << 1
		l--
	}
	return res
}
func calcLen(r []interface{}) int {
	len_ := 0
	for _, k := range r {
		switch k.(type) {
		case uint8:
			len_ += 1
		case uint16:
			len_ += 2
		case uint32:
			len_ += 4
		case uint64:
			len_ += 8
		case string:
			kk := k.(string)
			len_ += len(kk)
		default:
			panic(fmt.Sprintf("Can't find type in calc %s", k))
		}
	}
	return len_
}
