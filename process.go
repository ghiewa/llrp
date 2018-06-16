package llrp

import (
	"encoding/binary"
	log "github.com/sirupsen/logrus"
	//"strconv"
)

func (nc *RConn) sendReport(len_data int, reports ...interface{}) {
	if nc.sub == nil {
		log.Warningf("Subscription doesn't init yet.")
		return
	}
	nc.subsMu.RLock()
	nc.InMsgs++
	nc.InBytes += uint64(len_data)
	sub := nc.sub
	sub.mu.Lock()
	sub.pMsgs++
	if sub.pMsgs > sub.pMsgsMax {
		sub.pBytesMax = sub.pMsgs
	}
	sub.pBytes += len_data
	if sub.pBytes > sub.pBytesMax {
		sub.pBytesMax = sub.pBytes
	}
	m := &Msg{
		From:    nc.sub,
		Reports: reports,
	}
	if sub.pHead == nil {
		sub.pHead = m
		sub.pTail = m
		sub.pCond.Signal()
	} else {
		sub.pTail.next = m
		sub.pTail = m
	}
	sub.mu.Unlock()
	nc.subsMu.RUnlock()
}

// sent process message via subscript events
func (nc *RConn) process(b []byte, len_data int) error {
	// cut header & messageId
	log.Debugf("incomming package %d", len_data)
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
		disable_ro_resp   *DISABLE_ROSPEC_RESPONSE
		start_rospec_resp *START_ROSPEC_RESPONSE
		stop_rospec_resp  *STOP_ROSPEC_RESPONSE
		close_conn_resp   *CLOSE_CONNECTION_RESPONSE
		err_resp          *ERROR_MESSAGE
		keep_alive_resp   *KeepaliveResponse
		dam_res           *MsgLoss
		duticate_cards    = make(map[string]bool)
		reports           []interface{}
		vaild             = true
	)
	for len_data > 0 {
		vaild = true
		var (
			ro_resp *ROAccessReportResponse
		)
		header := binary.BigEndian.Uint16(b[walk:walk+2]) - 1024
		walk += 2
		len_p := int(binary.BigEndian.Uint32(b[walk : walk+4]))
		walk += 4
		switch header {
		case M_KEEPALIVE:
			keep_alive_resp = new(KeepaliveResponse)
			reports = append(reports, keep_alive_resp)
		case M_RO_ACCESS_REPORT:
			ro_resp = new(ROAccessReportResponse)
			ro_resp.MsgId = binary.BigEndian.Uint32(b[walk : walk+4])
			reports = append(reports, ro_resp)
		case M_READER_EVENT_NOTIFICATION:
			evt_resp = new(EventNotificationResponse)
			evt_resp.MsgId = binary.BigEndian.Uint32(b[walk : walk+4])
			reports = append(reports, evt_resp)
		case M_GET_READER_CONFIG_RESPONSE:
			get_resp = new(GetConfigResponse)
			get_resp.MsgId = binary.BigEndian.Uint32(b[walk : walk+4])
			reports = append(reports, get_resp)
		case M_SET_READER_CONFIG_RESPONSE:
			set_resp = new(SetConfigResponse)
			set_resp.MsgId = binary.BigEndian.Uint32(b[walk : walk+4])
			reports = append(reports, set_resp)
		case M_DELETE_ROSPEC_RESPONSE:
			del_ro_spec_resp = new(DELETE_ROSPEC_RESPONSE)
			del_ro_spec_resp.MsgId = binary.BigEndian.Uint32(b[walk : walk+4])
			reports = append(reports, del_ro_spec_resp)
		case M_CLOSE_CONNECTION_RESPONSE:
			close_conn_resp = new(CLOSE_CONNECTION_RESPONSE)
			close_conn_resp.MsgId = binary.BigEndian.Uint32(b[walk : walk+4])
			reports = append(reports, close_conn_resp)
		case M_DELETE_ACCESSSPEC_RESPONSE:
			del_acc_spec_resp = new(DELETE_ACCESSSPEC_RESPONSE)
			del_acc_spec_resp.MsgId = binary.BigEndian.Uint32(b[walk : walk+4])
			reports = append(reports, del_acc_spec_resp)
		case M_ADD_ROSPEC_RESPONSE:
			add_ro_resp = new(ADD_ROSPEC_RESPONSE)
			add_ro_resp.MsgId = binary.BigEndian.Uint32(b[walk : walk+4])
			reports = append(reports, add_ro_resp)
		case M_CUSTOM_MESSAGE:
			log.Warnf("M_CUSTOM_MESSAGE")
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
		case M_DISABLE_ROSPEC_RESPONSE:
			disable_ro_resp = new(DISABLE_ROSPEC_RESPONSE)
			disable_ro_resp.MsgId = binary.BigEndian.Uint32(b[walk : walk+4])
			reports = append(reports, disable_ro_resp)
		case M_START_ROSPEC_RESPONSE:
			start_rospec_resp = new(START_ROSPEC_RESPONSE)
			start_rospec_resp.MsgId = binary.BigEndian.Uint32(b[walk : walk+4])
			reports = append(reports, start_rospec_resp)
		case M_STOP_ROSPEC_RESPONSE:
			stop_rospec_resp = new(STOP_ROSPEC_RESPONSE)
			stop_rospec_resp.MsgId = binary.BigEndian.Uint32(b[walk : walk+4])
			reports = append(reports, stop_rospec_resp)
		case M_ENABLE_ROSPEC_RESPONSE:
			en_ro_resp = new(ENABLE_ROSPEC_RESPONSE)
			en_ro_resp.MsgId = binary.BigEndian.Uint32(b[walk : walk+4])
			reports = append(reports, en_ro_resp)
		case M_ERROR_MESSAGE:
			err_resp = new(ERROR_MESSAGE)
			err_resp.MsgId = binary.BigEndian.Uint32(b[walk : walk+4])
			reports = append(reports, err_resp)
		default:
			log.Errorf("cant handle code %d : %s", header, b)
			dam_res = new(MsgLoss)
			dam_res.Len = len_p
			reports = append(reports, dam_res)
			vaild = false
			break
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
			case P_ReaderEventNotificationData:
				evt_resp.Data, walk = parseEvtNotificationData(b, walk)
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
				case disable_ro_resp != nil:
					disable_ro_resp.Status = status
				case start_rospec_resp != nil:
					start_rospec_resp.Status = status
				case stop_rospec_resp != nil:
					stop_rospec_resp.Status = status
				case close_conn_resp != nil:
					close_conn_resp.Status = status
				case err_resp != nil:
					err_resp.Status = status
				default:
					// not implement yet(deletero , get cap , addro)
					log.Errorf("\nnot implement")
					vaild = false
					continue
				}
			default:
				log.Warnf("not implement %d ", code)
				// not implement yet will find len_ & skip parameter
				len_skip := int(binary.BigEndian.Uint16(b[walk : walk+2]))
				walk += (len_skip - 2)
			}
			walk_pre = walk - walk_pre
			len_pre -= walk_pre
		}
		len_data -= len_p
	}

	if vaild {
		nc.mu.Lock()
		nc.sendReport(len_data, reports...)
		nc.mu.Unlock()
	}
	return nil
}
