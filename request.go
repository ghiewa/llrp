package llrp

const (
	V_1011_All = iota
	V_1011_General_Device_Capabilities
	V_1011_LLRP_Capabilities
	V_1011_Regulatory_Capabilities
	V_1011_Air_Protocol_LLRP_Capabilities
)

// Convert custom pack to protocol []byte
func CustomPack(messageType, messageId int, config []interface{}, params ...[]interface{}) []byte {
	return bundle(messageType, messageId, config, params...)
}
func ENABLE_ROSPEC(messageId, id int) []byte {
	return bundle(
		M_ENABLE_ROSPEC,
		messageId,
		[]interface{}{
			uint32(id),
		},
	)
}

func GET_READER_CAPABILITIES_V1011(messageId, v_1011 int) []byte {
	data := []interface{}{
		uint8(v_1011),
	}
	return bundle(
		M_GET_READER_CAPABILITIES,
		messageId,
		data,
	)
}

const (
	V_1311_All = iota
	V_1311_Identification
	V_1311_AntennaProperties
	V_1311_AntennaConfiguration
	V_1311_ROReportSpec
	V_1311_ReaderEventNotificationSpec
	V_1311_AccessReportSpec
	V_1311_LLRPConfigurationStateValue
	V_1311_KeepaliveSpec
	V_1311_GPIPortCurrentState
	V_1311_GPOWriteData
	V_1311_EventsAndReports
)

func GET_READER_CONFIG_V1311(messageId, AntennaID, v_1311, GPIPortNum, GPOPortNum int) []byte {
	config := []interface{}{
		uint16(AntennaID),
		uint8(v_1311),
		uint16(GPIPortNum),
		uint16(GPOPortNum),
	}
	return bundle(
		M_GET_READER_CONFIG,
		messageId,
		config,
	)
}

func SET_READER_CONFIG(messageId int, restore_factory_setting bool, params ...[]interface{}) []byte {
	f := uint8(0)
	if restore_factory_setting {
		f = uint8(0x80)
	}
	config := []interface{}{
		f,
	}
	return bundle(
		M_SET_READER_CONFIG,
		messageId,
		config,
		params...,
	)
}
func DEL_ROSPEC(messageId, spec int) []byte {
	return bundle(
		M_DELETE_ROSPEC,
		messageId,
		[]interface{}{uint32(spec)},
	)
}

func SEND_KEEPALIVE(messageId int) []byte {
	return bundle(
		M_KEEPALIVE,
		0,
		nil,
	)
}

func ADD_ROSPEC(messageId int, params ...[]interface{}) []byte {
	return bundle(
		M_ADD_ROSPEC,
		messageId,
		nil,
		params...,
	)
}
