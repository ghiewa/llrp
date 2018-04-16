package llrp

import (
	//"bytes"
	"encoding/binary"
	"fmt"
	//"strconv"
)

type LLRPStatus struct {
	Success    bool
	StatusCode uint16
	ErrMsg     string
}
type ERROR_MESSAGE struct {
	MsgId  uint32
	Status *LLRPStatus
}
type ROAccessReportResponse struct {
	MsgId uint32
	Data  *TagReportData
}
type SetConfigResponse struct {
	MsgId  uint32
	Status *LLRPStatus
}
type CUSTOM_MESSAGE_RESPONSE struct {
	MsgId   uint32
	Vendor  uint32
	SubType uint8
	Status  *LLRPStatus
}
type ENABLE_ROSPEC_RESPONSE struct {
	MsgId  uint32
	Status *LLRPStatus
}
type EventNotificationResponse struct {
	MsgId uint32
	Data  *EvtData
}
type EvtData struct {
	TimestampUTC uint64
	GpiEvt       []*GpiEventParam
}
type GpiEventParam struct {
	PortNumber uint16
	Evt        bool
}
type DELETE_ROSPEC_RESPONSE struct {
	MsgId  uint32
	Status *LLRPStatus
}
type DELETE_ACCESSSPEC_RESPONSE struct {
	MsgId  uint32
	Status *LLRPStatus
}
type ADD_ROSPEC_RESPONSE struct {
	MsgId  uint32
	Status *LLRPStatus
}
type GetConfigResponse struct {
	MsgId  uint32
	Status *LLRPStatus
	Id     *Identification
	AnnPty []*AntennaProperty
	//AnnCon    []*AntennaConfig
	EvtSpec *ReaderEventNoticationSpec
	//ROSpec    *ROReportSpec
	//AccSpec   *AccessReportSpec
	//CfgState  *LLRPConfigState
	Ka        *KeepaliveSpec
	GPI       []*GPICurrentState
	GPO       []*GPOWriteData
	EvtReport *EventReport
	//Customs   []*CustomParamer
}
type EventReport struct {
	HoldEventReportsUponReconnect bool
}

type Identification struct {
	Type int
	Id   string
}
type AntennaProperty struct {
	Connected bool
	Id        uint16
	Gain      uint16
}
type ReaderEventNoticationSpec struct {
	State []*EventNotificationState
}
type EventNotificationState struct {
	EventType uint16
	State     int
}
type KeepaliveSpec struct {
	Type         int
	IntervalTime uint32
}

type GPICurrentState struct {
	Number uint16
	Config bool
	State  int
}
type GPOWriteData struct {
	Number uint16
	Data   bool
}
type TagReportData struct {
	EPC_96                string
	AntennaID             uint16
	PeakRSSI              int
	ChannelIndex          uint16
	FirstSeenTimestampUTC uint64
	AccessSpecId          uint32
}
type MsgLoss struct {
	Len int
}

func Response(b []byte, len_data int) (reports []interface{}) {
	// cut header & messageId
	var (
		walk              = 0
		get_resp          *GetConfigResponse
		set_resp          *SetConfigResponse
		evt_resp          *EventNotificationResponse
		del_ro_spec_resp  *DELETE_ROSPEC_RESPONSE
		del_acc_spec_resp *DELETE_ACCESSSPEC_RESPONSE
		add_ro_resp       *ADD_ROSPEC_RESPONSE
		custom_resp       *CUSTOM_MESSAGE_RESPONSE
		en_ro_resp        *ENABLE_ROSPEC_RESPONSE
		err_resp          *ERROR_MESSAGE
		dam_res           *MsgLoss
		duticate_cards    = make(map[string]bool)
	)
	//fmt.Printf("\n----loop")
	for len_data > 0 {
		var (
			ro_resp *ROAccessReportResponse
		)
		//fmt.Printf("\nwalk0 = %d %d", walk, len_data)
		header := binary.BigEndian.Uint16(b[walk:walk+2]) - 1024
		walk += 2
		len_p := int(binary.BigEndian.Uint32(b[walk : walk+4]))
		//fmt.Printf("\nwalk = %d:%d", walk, len_p)
		walk += 4
		switch header {
		case M_RO_ACCESS_REPORT:
			ro_resp = new(ROAccessReportResponse)
			ro_resp.MsgId = binary.BigEndian.Uint32(b[walk : walk+4])
			reports = append(reports, ro_resp)
		case M_READER_EVENT_NOTIFICATION:
			evt_resp = new(EventNotificationResponse)
			evt_resp.MsgId = binary.BigEndian.Uint32(b[walk : walk+4])
			reports = append(reports, evt_resp)
			//fmt.Printf("\nevt")
		case M_GET_READER_CONFIG_RESPONSE:
			get_resp = new(GetConfigResponse)
			get_resp.MsgId = binary.BigEndian.Uint32(b[walk : walk+4])
			reports = append(reports, get_resp)
			//fmt.Printf("\nget")
		case M_SET_READER_CONFIG_RESPONSE:
			set_resp = new(SetConfigResponse)
			set_resp.MsgId = binary.BigEndian.Uint32(b[walk : walk+4])
			reports = append(reports, set_resp)
		case M_DELETE_ROSPEC_RESPONSE:
			del_ro_spec_resp = new(DELETE_ROSPEC_RESPONSE)
			del_ro_spec_resp.MsgId = binary.BigEndian.Uint32(b[walk : walk+4])
			reports = append(reports, del_ro_spec_resp)
		case M_DELETE_ACCESSSPEC_RESPONSE:
			del_acc_spec_resp = new(DELETE_ACCESSSPEC_RESPONSE)
			del_acc_spec_resp.MsgId = binary.BigEndian.Uint32(b[walk : walk+4])
			reports = append(reports, del_acc_spec_resp)
		case M_ADD_ROSPEC_RESPONSE:
			add_ro_resp = new(ADD_ROSPEC_RESPONSE)
			add_ro_resp.MsgId = binary.BigEndian.Uint32(b[walk : walk+4])
			reports = append(reports, add_ro_resp)
		case M_CUSTOM_MESSAGE:
			custom_resp = new(CUSTOM_MESSAGE_RESPONSE)
			custom_resp.MsgId = binary.BigEndian.Uint32(b[walk : walk+4])
			reports = append(reports, custom_resp)
			walk += 4
			custom_resp.Vendor = binary.BigEndian.Uint32(b[walk : walk+4])
			walk += 4
			custom_resp.SubType = uint8(b[walk])
			walk -= 3
			len_p -= 5
			len_data -= 5
		case M_ENABLE_ROSPEC_RESPONSE:
			en_ro_resp = new(ENABLE_ROSPEC_RESPONSE)
			en_ro_resp.MsgId = binary.BigEndian.Uint32(b[walk : walk+4])
			reports = append(reports, en_ro_resp)
		case M_ERROR_MESSAGE:
			err_resp = new(ERROR_MESSAGE)
			err_resp.MsgId = binary.BigEndian.Uint32(b[walk : walk+4])
			reports = append(reports, err_resp)
		default:
			fmt.Printf("\ncant handle code %d : %s", header, b)
			dam_res = new(MsgLoss)
			dam_res.Len = len_p
			reports = append(reports, dam_res)
			return reports
			//panic(fmt.Sprintf("cant handle code %d : %s", header, b))
		}
		walk += 4
		var (
			len_pre = len_p - 10
		)
		for len_pre > 0 {
			walk_pre := walk
			code := binary.BigEndian.Uint16(b[walk : walk+2])
			walk += 2
			switch code {
			case P_TagReportData:
				ro_resp.Data, walk = parseTagData(b, walk)
				if ro_resp.Data != nil {
					if _, same := duticate_cards[ro_resp.Data.EPC_96]; same {
						reports = reports[:len(reports)-1]
					} else {
						duticate_cards[ro_resp.Data.EPC_96] = true
					}
				}
			case P_GPIPortCurrentState:
				var (
					gpi *GPICurrentState
				)
				gpi, walk = parseGPICurrentState(b, walk)
				get_resp.GPI = append(get_resp.GPI, gpi)
			case P_GPOWriteData:
				var (
					gpo *GPOWriteData
				)
				gpo, walk = parseGPOWriteData(b, walk)
				get_resp.GPO = append(get_resp.GPO, gpo)
			case P_LLRPStatus:
				var (
					status *LLRPStatus
				)
				status, walk = parseLLRP(b, walk)
				switch {
				case get_resp != nil:
					get_resp.Status = status
				case set_resp != nil:
					set_resp.Status = status
				case del_ro_spec_resp != nil:
					del_ro_spec_resp.Status = status
				case del_acc_spec_resp != nil:
					del_acc_spec_resp.Status = status
				case add_ro_resp != nil:
					add_ro_resp.Status = status
				case custom_resp != nil:
					custom_resp.Status = status
				case en_ro_resp != nil:
					en_ro_resp.Status = status
				case err_resp != nil:
					err_resp.Status = status
				default:
					// not implement yet(deletero , get cap , addro)
					panic(fmt.Sprintf("\nnot implement"))

				}
			//case P_ReaderEventNotificationData:
			//evt_resp.Data, walk = parseEvtNotificationData(b, walk)
			//fmt.Printf("\nEVT DATA %d", evt_resp.MsgId)
			default:
				// not implement yet will find len_ & skip parameter
				len_skip := int(binary.BigEndian.Uint16(b[walk : walk+2]))
				walk += (len_skip - 2)
				/*fmt.Printf("\nskip %d , %d , %d", code, walk, len_skip)
				if walk == 49 {
					panic(fmt.Sprintf("\n%d \n%d", code, len(b)))
				}
				*/
			}
			walk_pre = walk - walk_pre
			len_pre -= walk_pre
			//fmt.Printf("\nlen %d ,walk %d", len_pre, walk_pre)
		}
		len_data -= len_p
	}
	return
}
