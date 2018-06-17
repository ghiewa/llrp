package llrp

import (
//"bytes"
//"github.com/spaolacci/murmur3"
)

func (sp *SPReaderInfo) processCommandOptions() {
	nc := sp.conn
	nc.mu.Lock()
	defer nc.mu.Unlock()
	KeepaliveSpec
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
			param[0],
			param[1],
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

const (
	C_ROSpec_CurrentState_Disabled = iota
	C_ROSpec_CurrentState_Inactive
	C_ROSpec_CurrentState_Active
)

/*
This parameter carries the information of the Reader inventory and survey operation.
ROSpecID: Unsigned Integer; 0 is an illegal ROSpecID for a ROSpec.
Priority: Integer. Lower numbered priority values are given higher priority.
Possible Values: 0-7.
CurrentState: Integer
----------
ROBoundarySpec: ROBoundarySpec Parameter
ListOfSpecs: List of LLRP Parameters
Possible Values:
Each parameter can be either an <AISpec Parameter>, a <RFSurveySpec Parameter>, a <LoopSpec Parameter>, or a Custom Parameter.
ROReportSpec: ROReportSpec Parameter
*/
func AddROSpecCustom(messageId, ROSpecID, Priority, CurrentState int, spec ...[]interface{}) []byte {
	b := ADD_ROSPEC(
		messageId,
		ROSpec(
			ROSpecID,
			Priority,
			CurrentState,
			spec...,
		),
	)
	return b
}

func AddROSpecOptionDefault() []byte {
	return ADD_ROSPEC(
		7, // message id
		ROSpec(1234, 0, 0,
			ROBoundarySpec(
				ROSpecStartTrigger(1),
				ROSpecStopTrigger(0, 0),
			),
			AISpec(
				2,
				[]int{
					1, 2,
				},
				AISpecStopTrigger(0, 0),
				InventoryParameterSpec(1234, 1,
					AntennaConfiguration(1,
						RFTransmitter(1, 0, 10),
						C1G2InventoryCommand(
							false,
							C1G2RFControl(1000, 0),
							C1G2SingulationControl(0x80, 32, 0),
							CustomParameter(25882, 23, uint16(2)),
							CustomParameter(25882, 26, uint16(0), uint16(0), uint16(0)),
							CustomParameter(25882, 28, uint16(0), uint16(0), uint16(0)),
						),
					),
					AntennaConfiguration(2,
						RFTransmitter(1, 0, 81),
						C1G2InventoryCommand(
							false,
							C1G2RFControl(1000, 0),
							C1G2SingulationControl(0x80, 32, 0),
							CustomParameter(25882, 23, uint16(2)),
							CustomParameter(25882, 26, uint16(0), uint16(0), uint16(0)),
							CustomParameter(25882, 28, uint16(0), uint16(0), uint16(0)),
						),
					),
				),
			),
			RoReportSpec(2, 1,
				TagReportContentSelector(
					false,
					false,
					false,
					true,
					true,
					true,
					true,
					false,
					false,
					true,
				),
			),
		),
	)
}
