package llrp

import (
	"bufio"
	log "github.com/sirupsen/logrus"
	"os"
	"reflect"
	"testing"
	//"time"
)

var (
	card_evt   bool = true
	count           = 0
	limit_card      = 100
	am         bool
)

func handler(msg *Msg) {
	for _, k := range msg.Reports {
		switch k.(type) {
		case *NetworkIssue:
			kk := k.(*NetworkIssue)
			switch kk.Type {
			case NETW_LOSS:
				log.Warningf("Network loss on %s [%d] ", msg.From.Id, kk.Reconnects)
			case NETW_CONNECTED:
				log.Infof("Network connected %s ", msg.From.Id)
			default:
				log.Warningf("Network unknow state %s ", msg.From.Id)
			}
		case *EventNotificationResponse:
			//log.Infof("[EVT]")
		case *ROAccessReportResponse:
			if card_evt || am {
				log.Warnf("--- Form %s", msg.From.Id)
				kk := k.(*ROAccessReportResponse)
				if kk.Data != nil {
					log.Infof("[RO][%d][%s]", kk.MsgId, kk.Data.EPC_96)
				}
				count++
				if count > limit_card {
					count = 0
					card_evt = false
					log.Warningf("----We stop logs card here to do another operation")
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

func TestM(t *testing.T) {
	t.Run("th", loop)
}

func loop(t *testing.T) {
	log.Info("loop")
	opt := GetDefaultOptions()
	host := opt.NewConn()
	//log.SetLevel(log.DebugLevel)
	var valid bool
	readers := []*SPReaderInfo{
		&SPReaderInfo{
			Id:   "random_reader_id1",
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
	var text string
	var err error

	log.Infof("Please enter command\nlist - list of readers\nd - disable card event log\ne - enable card event log \nio - control gpo/get gpi state\nam - long run to test card logs")
	go func() {
		for {
			scanner := bufio.NewScanner(os.Stdin)
			scanner.Scan()
			text = scanner.Text()
			if text == "q" {
				break
			}
			err = nil
			valid = true

			switch text {
			case "am":
				log.Infof("starting non-stop cards log")
				am = true
			case "list":
				log.Infof("List Readers : %v", host.ListReader())
			case "d":
				// disable card logs
				card_evt = false
			case "e":
				// enable card logs
				card_evt = true
			case "io":
				log.Infof("sample command please enter number(0-2)")
				scanner.Scan()
				text = scanner.Text()
				switch text {
				case "0":
					// set gpo all open state // 0 = close , 1 = open , 2 = igonre
					// GPOset(id,port_state ...)  - set 4 port open state
					log.Infof("set gpo all on")
					err = host.GPOset(123, "random_reader_id", true, true, true, true)
				case "1":
					// set gpo spectfic port eg. port no.1 will open
					log.Infof("set gpo port 1 on")
					err = host.GPOsetp(222, "random_reader_id", 1, true)
				case "2":
					log.Infof("set gpo port 1 on")
					err = host.GPOset(123, "random_reader_id", false, true, false, true)
				case "3":
					// get all gpi port
					log.Infof("get gpi")
					err = host.GPIget(333, "random_reader_id")
				default:
					valid = false
				}
			default:
				valid = false
			}
			if valid {
				if err != nil {
					log.Warnf("Send command not success. %v", err)
				} else {
					log.Infof("Send command success. ")
				}
			}
		}
	}()
	select {}
	// close connection
	host.Close()

}
