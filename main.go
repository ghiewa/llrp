package main

import (
	. "./llrp"
	"bufio"
	log "github.com/sirupsen/logrus"
	"os"
	"reflect"
	"time"
)

var (
	card_evt   bool = true
	count           = 0
	limit_card      = 100
	toggle          = true
	host       *Conn
	am         bool
)

func handler(msg *Msg) {
	for _, k := range msg.Reports {
		ip := msg.From.Ip
		switch k.(type) {
		case *NetworkIssue:
			kk := k.(*NetworkIssue)
			switch kk.Type {
			case NETW_LOSS:
				log.Warningf("Network loss on [%s] ", ip)
			case NETW_CONNECTED:
				log.Infof("Network connected [%s]", ip)
			default:
				log.Warningf("Network unknow state %s %s ", msg.From.Id, ip)
			}
		case *EventNotificationResponse:
			kk := k.(*EventNotificationResponse)
			log.Infof("[EVT] %+v", kk)
		case *ROAccessReportResponse:
			if card_evt || am {
				//log.Warnf("--- Form %s", msg.From.Id)
				kk := k.(*ROAccessReportResponse)
				if kk.Data != nil {
					log.Infof("[RO][%s][%s][%s]", kk.Data.EPC_96, ip, msg.From.Id)
				}
				count++
				if count%11 == 0 {
					host.Lock()
					defer host.Unlock()
					// toggle io on reader_01
					err := host.GPOset(123, "random_reader_id_01", toggle, toggle, toggle, toggle)
					toggle = !toggle
					if err != nil {
						log.Errorf("error on gpo test set")
					}
				}
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
			if kk.Status != nil {
				log.Infof("[ENA_RO] Success=%v", kk.Status.Success)
			} else {
				log.Infof("[ENA_RO] %v", kk.Status)
			}
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
	opt.Timeout = time.Minute * 2
	opt.MaxReconnect = 10000
	opt.ReconnectWait = time.Minute * 2

	host = opt.NewConn()
	log.SetOutput(os.Stdout)
	//log.SetLevel(log.DebugLevel)
	var (
		valid        bool
		timeout      = uint32(10000) // milliseconds
		port_trigger = uint16(1)
	)
	readers := []*SPReaderInfo{
		&SPReaderInfo{
			Id:   "random_reader_id_00",
			Host: "192.168.33.16:5084",
			InitCommand: [][]byte{
				ResetFactoryOpt(),
				DelROSpecOpt(),
				DelAccOption(),
				ExtensionOption(),
				SetRegion(),
				SetEventSpecOption(),
				//AddROSpecOption(),
				AddROSpecCustom(
					// set trigger option - gpi
					RoBoundSpecCustom(
						//GPITriggerValue option = 3
						ROSpecStartTrigger(3,
							GPITriggerValue(port_trigger, true, timeout),
						),
						ROSpecStopTrigger(
							1, // stop by duration trigger
							timeout,
						),
					),
					GetDefaultAISpec(),
					GetRoReportSpec(),
				),

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
	var err error
	scanner := bufio.NewScanner(os.Stdin)
	for {
		log.Infof("\n***\tPlease enter command\nlist - list of readers\nd - disable card event log\ne - enable card event log \nio - control gpo/get gpi state\nam - long run to test card logs")
		scanner.Scan()
		text := scanner.Text()
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
			log.Infof("sample command please enter reader id number(00 - 10)")
			scanner.Scan()
			text = scanner.Text()
			id := text
			log.Infof("reader [%s] selected", id)
			reader_id := "random_reader_id_" + id
			log.Infof("Please set command [gp(o) | gp(i)]")
			scanner.Scan()
			cmd := scanner.Text()
			switch cmd {
			case "i":
				log.Infof("[GPI] get gpi port on reader [%s]", reader_id)
				err = host.GPIget(333, reader_id)
			case "o":
				log.Infof("[GPO] Please set state [ o(n) | of(f) ]")
				state := false
				scanner.Scan()
				switch scanner.Text() {
				case "n":
					state = true
				case "f":
				default:
					valid = false
					continue
				}
				log.Infof("[GPO] Please set port to command [%s] [ all | [1-4]]", reader_id)
				scanner.Scan()
				switch scanner.Text() {
				case "all":
					// GPOset(random_message_id,reader_id , ... order_state_port 1->4)
					err = host.GPOset(123, reader_id, state, state, state, state)
				case "1":
					// GPOsetp(random_message_id,reader_id , port,set_state)
					err = host.GPOsetp(222, reader_id, 1, state)
				case "2":
					err = host.GPOsetp(222, reader_id, 2, state)
				case "3":
					err = host.GPOsetp(222, reader_id, 3, state)
				case "4":
					err = host.GPOsetp(222, reader_id, 4, state)
				default:
					valid = false
					continue
				}
			default:
				valid = false
				continue
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
	// close connection
	host.Close()

}
