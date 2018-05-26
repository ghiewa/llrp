package llrp

func GetDefaultAISpec() []interface{} {
	return AISpec(2,
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
	)
}
func GetRoReportSpec() []interface{} {
	return RoReportSpec(2, 1,
		TagReportContentSelector(0x1e40),
	)

}
