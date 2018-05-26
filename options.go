package llrp

import (
//"bytes"
//"github.com/spaolacci/murmur3"
)

func (sp *SPReaderInfo) processCommandOptions() {
	nc := sp.conn
	nc.mu.Lock()
	defer nc.mu.Unlock()
}

func ResetFactoryOpt() []byte {
	return SET_READER_CONFIG(
		1,
		true,
	)
}

func DelROSpecOpt() []byte {
	return DEL_ROSPEC(
		2,
		0,
	)
}

func DelAccOption() []byte {
	return CustomPack(
		M_DELETE_ACCESSSPEC,
		3,
		[]interface{}{
			uint32(0),
		},
	)
}

// vendor , subtype , reserve
func ExtensionOption(params ...int) []byte {
	param := []int{
		25882,
		21,
		0,
	}
	for i, k := range param {
		param[i] = k
	}
	return CustomPack(
		M_CUSTOM_MESSAGE,
		4,
		[]interface{}{
			uint32(param[0]),
			uint8(param[1]),
			uint32(param[2]),
		},
	)
}

func SetRegion(params ...int) []byte {
	param := []int{
		25882,
		22,
		14,
	}
	for i, k := range param {
		param[i] = k
	}
	return SET_READER_CONFIG(
		5,
		true,
		CustomParameter(
			uint32(param[0]),
			uint32(param[1]),
			uint16(param[2]),
		),
	)
}

// if len(params) != 9 will be default value
func SetEventSpecOption(params ...bool) []byte {
	if len(params) != 9 {
		params = []bool{
			true, true, true,
			true, true, false,
			true, false, true,
		}
	}

	return SET_READER_CONFIG(
		6,
		false,
		ReaderEventNotification(
			params...,
		),
	)
}
func AddROSpecCustom(spec ...[]interface{}) []byte {
	return ADD_ROSPEC(
		11,
		RoSpec(1234, 0, 0,
			spec...,
		),
	)
}

func AddROSpecOption(params ...int) []byte {
	return ADD_ROSPEC(
		7,
		RoSpec(1234, 0, 0,
			// start_trigger_type, stop_trgger_type, duration_trigger
			// set 3 to gpi detect , 2 to event to GPI trigger and 4000 is gpi timeout
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
	)
}
func EnableROSpecOption(params ...int) []byte {
	if len(params) == 0 {
		params = append(params, 1234)
	}
	return ENABLE_ROSPEC(
		8,
		params[0],
	)
}
