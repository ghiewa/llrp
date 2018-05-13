package llrp

import (
	log "github.com/sirupsen/logrus"
	"reflect"
	"testing"
)

var (
	card_evt bool = true
)

func handler(msg *Msg) {
	// msg.From - reader id
	log.Warnf("--- Form %s", msg.From.Id)
	for _, k := range msg.Reports {
		switch k.(type) {
		case *ROAccessReportResponse:
			kk := k.(*ROAccessReportResponse)
			if kk.Data != nil {
				log.Infof("[RO][%d][%s]", kk.MsgId, kk.Data.EPC_96)
			} else {
				log.Infof("\n[RO]")
			}
		case *DELETE_ROSPEC_RESPONSE:
			kk := k.(*DELETE_ROSPEC_RESPONSE)
			if kk.Status != nil {
				log.Infof("[DELRO] success = %v\n", kk.Status.Success)
			}
		case *DELETE_ACCESSSPEC_RESPONSE:
			kk := k.(*DELETE_ACCESSSPEC_RESPONSE)
			if kk.Status != nil {
				log.Infof("[DELACC] success = %v\n", kk.Status.Success)
			}
		case *ADD_ROSPEC_RESPONSE:
			kk := k.(*ADD_ROSPEC_RESPONSE)
			if kk.Status != nil {
				log.Infof("[ADD_ROSPEC] success = %v ,%s\n", kk.Status.Success, kk.Status.ErrMsg)
			}
		case *GetConfigResponse:
			kk := k.(*GetConfigResponse)
			log.Infof("[GET][%d] : %+v", kk.MsgId, kk.Status)
			if kk.GPI != nil {
				log.Infof("\ngpi=")
				for _, kkk := range kk.GPI {
					log.Infof("[%d=%d],", kkk.Number, kkk.State)
				}
			}
		case *SetConfigResponse:
			kk := k.(*SetConfigResponse)
			log.Infof("[SET][%d] success=%v", kk.MsgId, kk.Status.Success)
		case *CUSTOM_MESSAGE_RESPONSE:
			kk := k.(*CUSTOM_MESSAGE_RESPONSE)
			if kk.Status != nil {
				log.Infof("[CUSTOM][%d] success=%v", kk.MsgId, kk.Status.Success)
			}
		case *ENABLE_ROSPEC_RESPONSE:
			kk := k.(*ENABLE_ROSPEC_RESPONSE)
			log.Infof("[ENA_RO] Success=%v", kk.Status.Success)
		case *EventNotificationResponse:
			log.Infof("[EVT]")
		case *ERROR_MESSAGE:
			kk := k.(*ERROR_MESSAGE)
			log.Errorf("[ERROR] code=%d ,msg=%s", kk.Status.StatusCode, kk.Status.ErrMsg)
		case *MsgLoss:
			kk := k.(*MsgLoss)
			log.Errorf("[MSG_DAMAGE] len=%d ", kk.Len)
		default:
			log.Errorf("Can't handle type %v", reflect.TypeOf(k))
		}
	}
}

func TestM(t *testing.T) {
	t.Run("th", loop)
}

func loop(t *testing.T) {
	opt := GetDefaultOptions()
	host := opt.NewConn()
	//log.SetLevel(log.DebugLevel)
	readers := []*SPReaderInfo{
		&SPReaderInfo{
			Id:   "random_reader_id",
			Host: "192.168.33.16:5084",
			InitCommand: [][]byte{
				ResetFactoryOpt(),
				DelROSpecOpt(),
				DelAccOption(),
				ExtensionOption(),
				SetRegion(),
				SetEventSpecOption(),
				AddROSpecOption(),
				EnableROSpecOption(),
			},
		},
	}
	for _, reader := range readers {
		// doReconnected when loss signal
		err := host.Registry(reader)
		if err != nil {
			log.Error(err)
		}
	}
	host.Subscription(handler)
	select {}
	// close connection
	host.Close()
}
