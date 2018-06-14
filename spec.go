package llrp

import (
//log "github.com/sirupsen/logrus"
)

/*

Value    Definition
----------
	0 	Upon hopping to next channel (e.g., in FCC regulatory region)
	1	GPI event
	2	ROSpec event (start/end/preempt)
	3	Report buffer fill warning
	4	Reader exception event
	5	RFSurvey event (start/end)
	6	AISpec event (end)
	7	AISpec event (end) with singulation details
	8	Antenna event
	9	SpecLoop event
*/
func ReaderEventNotificationSpec(param ...bool) []interface{} {
	len_p := len(param)
	if len_p == 0 {
		// default
		param = []bool{
			true, true, true, true, true,
			false, true, false, true,
		}
		len_p = 9
	}
	l := 9 - len_p
	for l > 0 {
		param = append(param, false)
		l--
	}
	var (
		len_ = 4
		evt  = EventNotificationStateSpec(param)
	)
	len_ += calcLen(evt)
	r := []interface{}{
		uint16(P_ReaderEventNotificationSpec),
		uint16(len_), // len
	}
	for _, k := range evt {
		r = append(r, k)
	}
	return r
}
func EventNotificationStateSpec(param []bool) []interface{} {
	var (
		res []interface{}
	)
	//log.Infof("len %d", len(param))
	for i, k := range param {
		kk := uint8(0)
		if k {
			kk = 0x80
		}
		r := []interface{}{
			uint16(P_EventNotificationState),
			uint16(7), //len
			uint16(i),
			uint8(kk),
		}
		for _, k := range r {
			res = append(res, k)
		}
	}
	return res
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

// If time = 0 ,Reader will not send to client
func KeepaliveSpec(ms_timeinterval int) []interface{} {
	var (
		enable = 0
	)
	if ms_timeinterval > 0 {
		enable = 1
	}

	r := []interface{}{
		uint16(P_KeepaliveSpec),
		uint16(9),
		// 0 : Null – No keepalives SHALL be sent by the Reader
		// 1 : Periodic
		uint8(enable),
		uint32(ms_timeinterval),
	}

	return r
}

func RoReportSpec(trigger, n int, params ...[]interface{}) []interface{} {
	l := 7
	for _, k := range params {
		l += calcLen(k)
	}
	r := []interface{}{
		uint16(P_ROReportSpec),
		uint16(l),
		uint8(trigger),
		uint16(n),
	}
	for _, k := range params {
		r = append(r, k...)
	}
	return r
}

func TagReportContentSelector(enable int) []interface{} {
	return []interface{}{
		uint16(P_TagReportContentSelector),
		uint16(6),
		uint16(enable),
	}
}

// RoBoundSpec custom
// params should be ROSpecStartTrigger and ROSpecStopTrigger
func RoBoundSpecCustom(params ...[]interface{}) []interface{} {
	var (
		l = 4
	)
	for _, k := range params {
		l += calcLen(k)
	}
	r := []interface{}{
		uint16(P_ROBoundarySpec),
		uint16(l),
	}
	for _, k := range params {
		r = append(r, k...)
	}
	return r

}
