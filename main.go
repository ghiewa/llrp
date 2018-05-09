package main

import (
	. "./llrp"
	"fmt"
	//"io"
	"net"
	"reflect"
	"time"
)

var (
	state_toggle bool
	conn         net.Conn
	messageId    = 100
	err          error

	ip   = "192.168.33.16"
	port = "5084"
)

const (
	BufferSize = 512
)

func test_set_gpi() {
	// no. 168
	conn.Write(SET_READER_CONFIG(messageId, false,
		GPIPortCurrentState_Param(1, 0, false),
		GPIPortCurrentState_Param(2, 0, false),
		GPIPortCurrentState_Param(3, 0, false),
		GPIPortCurrentState_Param(4, 0, false),
	))
}
func test_set_region() {
	// no. 170
	conn.Write(SET_READER_CONFIG(messageId, true,
		CustomParameter(
			uint32(25882),
			uint32(22),
			uint16(14),
		)))
}
func test_set_gpo() {
	// no. 179
	conn.Write(SET_READER_CONFIG(messageId, false,
		GPOWriteData_Param(1, true),
		GPOWriteData_Param(2, true),
		GPOWriteData_Param(3, true),
		GPOWriteData_Param(4, true),
	))
}
func test_set_event_notice_spec() {
	// no. 182
	conn.Write(SET_READER_CONFIG(messageId, false,
		ReaderEventNotification(
			true, true, true,
			true, true, false,
			true, false, true,
		)),
	)

}
func test_resettofactory() {
	conn.Write(SET_READER_CONFIG(messageId, true))
}
func test_delro() {
	spec := 0
	conn.Write(DEL_ROSPEC(messageId, spec))
}
func test_extension() {
	messageType := M_CUSTOM_MESSAGE
	vendor := 25882
	subtype := 21
	reserve := 0
	conn.Write(CustomPack(
		messageType,
		messageId,
		[]interface{}{
			uint32(vendor),
			uint8(subtype),
			uint32(reserve),
		},
	))

}
func test_delacc() {
	spec := 0
	messageType := M_DELETE_ACCESSSPEC
	conn.Write(CustomPack(messageType, messageId,
		[]interface{}{
			uint32(spec),
		},
	))
}
func test_addro_v2() {
	messageId = 201
	conn.Write(ADD_ROSPEC(messageId,
		RoSpec(1234, 0, 0,
			RoBoundSpec(1, 0, 0),
			AISpec(2,
				AISpecStopTrigger(0, 0),
				InventoryParameterSpec(1234, 1,
					AntennaConfiguration(1,
						RFTransmitter(1, 0, 81),
						C1G2InventoryCommand(0,
							C1G2RFControl(1000, 0),
							C1G2SingulationControl(0x80, 32, 0),
							CustomParameter(uint32(25882), uint32(23), uint16(2)),
							CustomParameter(uint32(25882), uint32(26), uint16(0), uint16(0), uint16(0)),
							CustomParameter(uint32(25882), uint32(28), uint16(0), uint16(0), uint16(0)),
						),
					),
					AntennaConfiguration(2,
						RFTransmitter(1, 0, 81),
						C1G2InventoryCommand(0,
							C1G2RFControl(1000, 0),
							C1G2SingulationControl(0x80, 32, 0),
							CustomParameter(uint32(25882), uint32(23), uint16(2)),
							CustomParameter(uint32(25882), uint32(26), uint16(0), uint16(0), uint16(0)),
							CustomParameter(uint32(25882), uint32(28), uint16(0), uint16(0), uint16(0)),
						),
					),
				),
			),
			RoReportSpec(2, 1,
				TagReportContentSelector(0x1e40),
			),
		),
	))

}
func test_addro() {
	conn.Write(ADD_ROSPEC(messageId,
		RoSpec(1234, 0, 0,
			RoBoundSpec(1, 0, 0),
			AISpec(2,
				AISpecStopTrigger(0, 0),
				InventoryParameterSpec(1234, 1,
					AntennaConfiguration(1,
						RFReceiver(2),
						RFTransmitter(1, 0, 1),
						C1G2InventoryCommand(0,
							C1G2RFControl(1000, 0),
							C1G2SingulationControl(0x80, 32, 0),
							CustomParameter(uint32(25882), uint32(23), uint16(2)),
							CustomParameter(uint32(25882), uint32(26), uint16(0), uint16(0), uint16(0)),
							CustomParameter(uint32(25882), uint32(28), uint16(1), uint16(2000), uint16(250)),
						),
					),
					AntennaConfiguration(2,
						RFReceiver(2),
						RFTransmitter(1, 0, 81),
						C1G2InventoryCommand(0,
							C1G2RFControl(1000, 0),
							C1G2SingulationControl(0x80, 32, 0),
							CustomParameter(uint32(25882), uint32(23), uint16(2)),
							CustomParameter(uint32(25882), uint32(26), uint16(0), uint16(0), uint16(0)),
							CustomParameter(uint32(25882), uint32(28), uint16(1), uint16(2000), uint16(250)),
						),
					),
				),
			),
			RoReportSpec(2, 1,
				TagReportContentSelector(0x1e40),
			),
		),
	))
}
func response_test(buf []byte, len_ int) {
	reports := Response(buf, len_)
	for _, k := range reports {
		switch k.(type) {
		case *ROAccessReportResponse:
			kk := k.(*ROAccessReportResponse)
			if kk.Data != nil {
				fmt.Printf("\n[RO][%d][%s]", kk.MsgId, kk.Data.EPC_96)
			} else {
				fmt.Printf("\n[RO]")
			}
		case *DELETE_ROSPEC_RESPONSE:
			kk := k.(*DELETE_ROSPEC_RESPONSE)
			if kk.Status != nil {
				fmt.Printf("\n[DELRO] success = %v\n", kk.Status.Success)
			}
		case *DELETE_ACCESSSPEC_RESPONSE:
			kk := k.(*DELETE_ACCESSSPEC_RESPONSE)
			if kk.Status != nil {
				fmt.Printf("\n[DELACC] success = %v\n", kk.Status.Success)
			}
		case *ADD_ROSPEC_RESPONSE:
			kk := k.(*ADD_ROSPEC_RESPONSE)
			if kk.Status != nil {
				fmt.Printf("\n[ADD_ROSPEC] success = %v ,%s\n", kk.Status.Success, kk.Status.ErrMsg)
			}
		case *GetConfigResponse:
			kk := k.(*GetConfigResponse)
			fmt.Printf("\n[GET][%d] : %+v", kk.MsgId, kk.Status)
			if kk.GPI != nil {
				fmt.Printf("\ngpi=")
				for _, kkk := range kk.GPI {
					fmt.Printf("[%d=%d],", kkk.Number, kkk.State)
				}
			}
		case *SetConfigResponse:
			kk := k.(*SetConfigResponse)
			if kk.Status == nil {
				continue
			}
			fmt.Printf("\n[SET][%d] success=%v", kk.MsgId, kk.Status.Success)
		case *CUSTOM_MESSAGE_RESPONSE:
			kk := k.(*CUSTOM_MESSAGE_RESPONSE)
			if kk.Status != nil {
				fmt.Printf("\n[CUSTOM][%d] success=%v", kk.MsgId, kk.Status.Success)
			}
		case *ENABLE_ROSPEC_RESPONSE:
			kk := k.(*ENABLE_ROSPEC_RESPONSE)
			if kk.Status != nil {
				fmt.Printf("\n[ENA_RO] Success=%v", kk.Status.Success)
			}
		case *EventNotificationResponse:
			fmt.Printf("\n[EVT]")
		case *ERROR_MESSAGE:
			kk := k.(*ERROR_MESSAGE)
			if kk.Status != nil {
				fmt.Printf("\n[ERROR] code=%d ,msg=%s", kk.Status.StatusCode, kk.Status.ErrMsg)
			}
		case *MsgLoss:
			kk := k.(*MsgLoss)
			fmt.Printf("\n[MSG_DAMAGE] len=%d ", kk.Len)
		default:
			panic(fmt.Sprintf("Can't handle type %v", reflect.TypeOf(k)))
		}
	}
}
func test_enable_ro() {
	messageId = 202
	conn.Write(
		ENABLE_ROSPEC(messageId, 1234),
	)

}
func onloop() {
	// will loop send gpi request [bypass ack send] & get response both gpi response & ro report
	var (
		send_gpi_req         = time.NewTicker(time.Millisecond * 4000).C
		test_toggle_gpo      = time.NewTicker(time.Millisecond * 15000).C
		time_delay           = time.Second * 15
		reconnect            = time.NewTicker(time.Hour * 24).C
		be_offline           = make(chan bool)
		msg                  = make(chan bool)
		buf                  = make([]byte, BufferSize)
		reqLen               = 0
		cmd                  = make(chan bool)
		command              string
		args                 int
		something_wrong_here      = make(chan bool)
		online               bool = true
	)
	test_resettofactory()        // 1
	test_delro()                 // 2
	test_delacc()                // 3
	test_extension()             // 4
	test_set_gpi()               // 5
	test_set_region()            // 6
	test_set_gpo()               // 7
	test_set_event_notice_spec() // 8
	//test_addro()                 // 9
	test_addro_v2()
	test_enable_ro()
	//test_set_gpo()
	for {
		if online {
			go func() {

				reqLen, err = conn.Read(buf)
				if err != nil {
					be_offline <- true
				}
				if reqLen < BufferSize {
					msg <- true
				} else {
					something_wrong_here <- true
				}
			}()
		}
		select {
		//case <-send_keep_alive:
		case <-msg:
			response_test(buf, reqLen)
		case <-something_wrong_here:
			fmt.Printf("\n***************detecting error : %d , Will reconnecting", reqLen)
			panic(fmt.Sprintf("buf full %d", reqLen))
			be_offline <- true
		case <-cmd:
			if !online {
				continue
			}
			// test gpo
			set_gpo := false
			switch command {
			case "on":
				set_gpo = true
			case "off":
			default:
				break
				panic("Wrong cmd")
			}

			if args == 0 {
				// set all gpo
				conn.Write(SET_READER_CONFIG(messageId, false,
					GPOWriteData_Param(1, set_gpo),
					GPOWriteData_Param(2, set_gpo),
					GPOWriteData_Param(3, set_gpo),
					GPOWriteData_Param(4, set_gpo),
				))
			} else {
				conn.Write(SET_READER_CONFIG(messageId, false,
					GPOWriteData_Param(args, set_gpo),
				))
			}
		case gpo := <-gpos:
			if !online {
				continue
			}
			conn.Write(SET_READER_CONFIG(messageId, false,
				GPOWriteData_Param(gpo.Port, gpo.State),
			))
		case <-test_toggle_gpo:
			if online {
				web_req()
			}
		case <-send_gpi_req:
			if !online {
				continue
			}
			fmt.Printf("\n-----SEND GPI_REQ")
			conn.Write(GET_READER_CONFIG_V1311(
				messageId,
				0, //
				V_1311_GPIPortCurrentState,
				0, // all gpi
				0, // ignore
			))
		case <-reconnect:
			if online {
				continue
			}
			conn, err = net.Dial("tcp", ip+":"+port)
			if err != nil {
				fmt.Errorf(err.Error())
				continue
			}
			online = true
			reconnect = time.NewTicker(time.Hour * 24).C
		case <-be_offline:
			online = false
			fmt.Printf("\nClosing connection = %v", conn.Close())
			reconnect = time.NewTicker(time_delay).C
		}
	}
}

type gpoCmd struct {
	Port  int
	State bool
}

var (
	gpos = make(chan *gpoCmd)
)

func set_gpo_cmd(port int, state bool) {
	go func() {
		gpos <- &gpoCmd{
			Port:  port,
			State: state,
		}
	}()
}

func init() {
	conn, err = net.Dial("tcp", ip+":"+port)
	if err != nil {
		fmt.Errorf(err.Error())
		return
	}

}

func web_req() {
	fmt.Printf("\n--------***************************************** Web request Simulate")
	state_toggle = !state_toggle
	set_gpo_cmd(1, state_toggle)
}

func main() {
	start := time.Now()
	defer func() {
		fmt.Printf("%s", time.Now().Sub(start).String)
	}()
	onloop()

}
