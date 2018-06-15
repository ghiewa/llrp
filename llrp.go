package llrp

func bundle(messageType int, messageId int, config []interface{}, params ...[]interface{}) []byte {
	var (
		length = 10
	)
	length += calcLen(config)
	for _, k := range params {
		length += calcLen(k)
	}
	var (
		data = []interface{}{
			uint16(messageType + 1024),
			uint32(length),
			uint32(messageId),
		}
	)
	data = append(data, config...)
	for _, k := range params {
		data = append(data, k...)
	}
	return pack(data)
}
