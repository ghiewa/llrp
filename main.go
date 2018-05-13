package llrp

import (
	. "./llrp"
	"bufio"
	log "github.com/sirupsen/logrus"
	"os"
	"reflect"
)

var (
	card_evt bool = true
)

func handler(msg *Msg) {
	log.Warnf("--- Form %s", msg.From.Id)
	for _, k := range msg.Reports {
		switch k.(type) {
		case *ROAccessReportResponse:
			if card_evt {
				kk := k.(*ROAccessReportResponse)
				if kk.Data != nil {
					log.Infof("[RO][%d][%s]", kk.MsgId, kk.Data.EPC_96)
				} else {
					log.Infof("\n[RO]")
				}
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
			log.Infof("[ERROR] code=%d ,msg=%s", kk.Status.StatusCode, kk.Status.ErrMsg)
		case *MsgLoss:
			kk := k.(*MsgLoss)
			log.Errorf("[MSG_DAMAGE] len=%d ", kk.Len)
		default:
			log.Errorf("Can't handle type %v", reflect.TypeOf(k))
		}

	}
	// msg.From - reader id
}

func main() {
	log.Info("loop")
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
	log.Info("registry")
	for _, reader := range readers {
		// doReconnected when loss signal
		err := host.Registry(reader)
		if err != nil {
			log.Error(err)
		}
	}
	host.Subscription(handler)
	var text string
	var err error
	for {
		err = nil
		log.Infof("Please enter command\nreader - list of readers\nnce - disable card event log\nce - enable card event log \nio - control gpo/get gpi state")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		text = scanner.Text()
		if text == "q" {
			break
		}
		switch text {
		case "reader":
			log.Infof("List Readers : %v", host.ListReader())
		case "ne":
			// disable card logs
			card_evt = false
		case "ce":
			// enable card logs
			card_evt = true
		case "io":
			log.Infof("sample command please enter number(0-2)")
			text, _ = reader.ReadString('\n')
			switch text {
			case "0":
				// set gpo all open state // 0 = close , 1 = open , 2 = igonre
				// GPOset(id,port_state ...)  - set 4 port open state
				err = host.GPOset(123, "random_reader_id", true, true, true, true)
			case "1":
				// set gpo spectfic port eg. port no.1 will open
				err = host.GPOsetp(222, "random_reader_id", 1, true)
			case "2":
				// get all gpi port
				err = host.GPIget(333, "random_reader_id")
			default:
				log.Infof("not found cmd here")
			}
		default:
			log.Infof("not found cmd here")
		}
		if err != nil {
			log.Warnf("Send command not success. %v", err)
		} else {
			log.Infof("Send command success. ")
		}
	}
	// close connection
	host.Close()
}
