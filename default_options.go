package llrp

func GetDefaultAISpec() []interface{} {
	return AISpec(
		2,           // Antennacount
		[]int{1, 2}, // AntennaIdn
		AISpecStopTrigger(0, 0),
		InventoryParameterSpec(1234, 1,
			AntennaConfiguration(1,
				RFTransmitter(1, 0, 0x0a),
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
				RFTransmitter(1, 0, 0x0a),
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
	)
}
func GetRoReportSpec() []interface{} {
	return RoReportSpec(2, 1,
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
	)

}
