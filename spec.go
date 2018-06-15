package llrp

func commonSpec(code int, details []interface{}, params ...[]interface{}) []interface{} {
	var (
		l = calcLen(details) + 4
	)
	for _, k := range params {
		l += calcLen(k)
	}
	r := []interface{}{
		uint16(code),
		uint16(l),
	}
	for _, k := range details {
		r = append(r, k)
	}
	for _, k := range params {
		r = append(r, k...)
	}
	return r
}
