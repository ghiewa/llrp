package llrp

import (
	//"bytes"
	"encoding/binary"
	"fmt"
	log "github.com/sirupsen/logrus"
	//"strconv"
)

func parseEvtNotificationData(b []byte, walk int) (*EvtData, int) {
	// ReaderEventNotificationData
	// ignore
	var (
		evt   = new(EvtData)
		len_p = int(binary.BigEndian.Uint16(b[walk:walk+2])) - 4
	)
	walk += 2
	for len_p > 0 {
		code := int(binary.BigEndian.Uint16(b[walk : walk+2]))
		walk += 2
		len_ := int(binary.BigEndian.Uint16(b[walk : walk+2]))
		walk += 2
		switch code {
		case P_GPIEvent:
			var (
				gpi *GpiEventParam
			)
			log.Infof("---- 1 gpi --state %d", walk)
			gpi, walk = parseGPIEvent(b, walk)
			log.Infof("---- 2 gpi --state %d", walk)
			evt.GpiEvt = append(evt.GpiEvt, gpi)
		case P_UTCTimeStamp:
			evt.TimestampUTC = binary.BigEndian.Uint64(b[walk : walk+8])
			//fmt.Printf("\nget time %d", evt.TimestampUTC)
			walk += 8
		default:
			// ReaderEventNotificationData Parameters
			// will skip
			walk += (len_ - 4)
		}
		//fmt.Printf("\nCODE %d %d %d", code, len_, len_p)
		len_p -= len_

	}
	return evt, walk
}
func parseGPIEvent(b []byte, walk int) (*GpiEventParam, int) {
	var (
		gpi = new(GpiEventParam)
		evt bool
	)
	// skip len
	walk += 2
	gpi.PortNumber = binary.BigEndian.Uint16(b[walk : walk+2])
	walk += 2
	_evt := int(b[walk : walk+1][0])
	if _evt != 0 {
		evt = true
	}
	gpi.Evt = evt
	walk += 1
	return gpi, walk
}

func parseGPOWriteData(b []byte, walk int) (*GPOWriteData, int) {
	var (
		gpo     = new(GPOWriteData)
		config_ = false
	)
	// skip len
	walk += 2
	gpo.Number = binary.BigEndian.Uint16(b[walk : walk+2])
	walk += 2
	config := int(b[walk : walk+1][0])
	if config != 0 {
		config_ = true
	}
	gpo.Data = config_
	walk += 1
	return gpo, walk

}
func parseGPICurrentState(b []byte, walk int) (*GPICurrentState, int) {
	var (
		gpi = new(GPICurrentState)
	)
	//fmt.Printf("\ngpi")
	// skip len
	walk += 2
	gpi.Number = binary.BigEndian.Uint16(b[walk : walk+2])
	walk += 2
	config := int(b[walk : walk+1][0])
	walk += 1
	config_ := false
	if config != 0 {
		config_ = true
	}
	gpi.Config = config_
	gpi.State = int(b[walk : walk+1][0])
	walk += 1
	return gpi, walk
}

func parseTagData(b []byte, walk int) (*TagReportData, int) {
	var (
		tag         = new(TagReportData)
		crumb, step = 0, 0
	)
	len_p := int(binary.BigEndian.Uint16(b[walk:walk+2])) - 4
	walk += 2
	for len_p > 0 {
		code := b[walk : walk+1][0] - 128
		walk += 1
		switch code {
		case P_EPC_96:
			crumb = 12
			step = walk + crumb
			tag.EPC_96 = fmt.Sprintf("%x", (b[walk:step]))
		case P_AntennaID:
			crumb = 2
			step = walk + crumb
			tag.AntennaID = binary.BigEndian.Uint16(b[walk:step])
		case P_PeakRSSI:
			crumb = 1
			step = walk + crumb
			tag.PeakRSSI = int(b[walk:step][0])
		case P_ChannelIndex:
			crumb = 2
			step = walk + crumb
			tag.ChannelIndex = binary.BigEndian.Uint16(b[walk:step])
		case P_FirstSeenTimestampUTC:
			crumb = 8
			step = walk + crumb
			tag.FirstSeenTimestampUTC = binary.BigEndian.Uint64(b[walk:step])
		case P_AccessSpecID:
			crumb = 4
			step = walk + crumb
			tag.AccessSpecId = binary.BigEndian.Uint32(b[walk:step])
		default:
			log.Errorf("Cann't detect header %d on TagReportData", code)
			return tag, walk
		}
		walk += crumb
		len_p -= crumb + 1
	}
	return tag, walk
}

func parseLLRP(b []byte, walk int) (*LLRPStatus, int) {
	var (
		resp *LLRPStatus
	)
	len_ := int(binary.BigEndian.Uint16(b[walk : walk+2]))
	walk += 2
	code := binary.BigEndian.Uint16(b[walk : walk+2])
	walk += 2
	switch code {
	case M_Success:
		resp = &LLRPStatus{
			Success: true,
		}
	default:
		resp = &LLRPStatus{
			Success:    false,
			StatusCode: code,
			ErrMsg:     string(b[walk : walk+len_-6]),
		}
	}
	walk += len_ - 6
	return resp, walk
}
