package llrp

import (
	"fmt"
)

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
func ROSpecStartTrigger(typeof int, params ...[]interface{}) []interface{} {
	var (
		l = 5
	)

	for _, k := range params {
		l += calcLen(k)
	}
	r := []interface{}{
		uint16(P_ROSpecStartTrigger),
		uint16(l),
		//uint8(typeof),
	}
	for _, k := range params {
		r = append(r, k...)
	}
	return r
}

// option period trigger for rospec start/stop trigger
func PeriodicTriggerValue(UTCTime uint64, offset uint32, period uint32) []interface{} {
	l := 4 + 8 + 4 + 4
	return []interface{}{
		uint16(P_PeriodicTriggerValue),
		uint16(l),
		UTCTime,
		offset,
		period,
	}
}

// option gpi trigger for rospec start/stop trigger
func GPITriggerValue(GPIPortNum uint16, GPIEvent bool, Timeout uint32) []interface{} {
	l := 4 + 2 + 1 + 4
	ev := uint8(0) // disable
	if GPIEvent {
		ev = 1
	}
	return []interface{}{
		uint16(P_GPITriggerValue),
		uint16(l),
		GPIPortNum,
		ev,
		Timeout,
	}
}

func ROSpecStopTrigger(typeof int, DurationTrigger uint32, params ...[]interface{}) []interface{} {
	var (
		l = 4 + 1 + 4
	)
	for _, k := range params {
		l += calcLen(k)
	}
	r := []interface{}{
		uint16(P_ROSpecStopTrigger),
		uint16(l),
		uint8(typeof),
		DurationTrigger,
	}
	for _, k := range params {
		r = append(r, k...)
	}
	return r

}
func RoBoundSpec(startTriggerType, stopTriggerType, DurationTrigger int) []interface{} {
	bound := []interface{}{
		uint16(P_ROBoundarySpec),
		uint16(18),
		uint16(P_ROSpecStartTrigger),
		uint16(5),
		uint8(startTriggerType),
		uint16(P_ROSpecStopTrigger),
		uint16(9),
		uint8(stopTriggerType),
		uint32(DurationTrigger),
	}
	return bound
}

func AISpecStopTrigger(typeAIspec, duration int) []interface{} {
	return []interface{}{
		uint16(P_AISpecStopTrigger),
		uint16(9),
		uint8(typeAIspec),
		uint32(duration),
	}
}
func AntennaConfiguration(id int, params ...[]interface{}) []interface{} {
	var (
		l = 6
	)
	for _, k := range params {
		l += calcLen(k)
	}
	r := []interface{}{
		uint16(P_AntennaConfiguration),
		uint16(l),
		uint16(id),
	}
	for _, k := range params {
		r = append(r, k...)
	}
	return r
}
func RFReceiver(sensitivity int) []interface{} {
	return []interface{}{
		uint16(P_RFReceiver),
		uint16(6),
		uint16(sensitivity),
	}
}
func RFTransmitter(hopTableid, channelIndex, power int) []interface{} {
	return []interface{}{
		uint16(P_RFTransmitter),
		uint16(10),
		uint16(hopTableid),
		uint16(channelIndex),
		uint16(power),
	}
}
func C1G2InventoryCommand(stateTagaware int, params ...[]interface{}) []interface{} {
	var (
		l = 5 // type(2) + len(2) + tag state aware(1)
	)
	for _, k := range params {
		l += calcLen(k)
	}
	r := []interface{}{
		uint16(P_C1G2InventoryCommand),
		uint16(l),
		uint8(stateTagaware),
	}
	for _, k := range params {
		r = append(r, k...)
	}
	return r
}

func C1G2RFControl(mode, tari int) []interface{} {
	return []interface{}{
		uint16(P_C1G2RFControl),
		uint16(8),
		uint16(mode),
		uint16(tari),
	}
}
func C1G2SingulationControl(session, population, tranzittime int) []interface{} {
	return []interface{}{
		uint16(P_C1G2SingulationControl),
		uint16(11),
		uint8(session),
		uint16(population),
		uint32(tranzittime),
	}

}
func CustomParameter(params ...interface{}) []interface{} {
	l := calcLen(params) + 4
	r := []interface{}{
		uint16(P_Custom),
		uint16(l),
	}
	return append(r, params...)
}
func InventoryParameterSpec(id, protocolId int, params ...[]interface{}) []interface{} {
	var (
		l = 7
	)
	for _, k := range params {
		l += calcLen(k)
	}
	r := []interface{}{
		uint16(P_InventoryParameterSpec),
		uint16(l),
		uint16(id),
		uint8(protocolId),
	}
	for _, k := range params {
		r = append(r, k...)
	}
	return r
}
func AISpec(count_antenna int, params ...[]interface{}) []interface{} {
	var (
		len_ = (count_antenna * 2) + 6
	)
	for _, k := range params {
		len_ += calcLen(k)
	}
	ai := []interface{}{
		uint16(P_AISpec),
		uint16(len_),
		uint16(count_antenna),
	}
	for i := 0; i < count_antenna; i++ {
		ai = append(ai, uint16(i+1))
	}
	for _, k := range params {
		ai = append(ai, k...)
	}
	return ai
}

func RoSpec(rospecId, priority, currentState int, params ...[]interface{}) []interface{} {
	var (
		len_ = 10
	)
	for _, k := range params {
		len_ += calcLen(k)
	}
	spec := []interface{}{
		uint16(P_ROSpec),
		uint16(len_),
		uint32(rospecId),
		uint8(priority),
		uint8(currentState),
	}
	for _, k := range params {
		spec = append(spec, k...)
	}
	return spec
}
