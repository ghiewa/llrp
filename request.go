package llrp

const (
	V_1011_All = iota
	V_1011_General_Device_Capabilities
	V_1011_LLRP_Capabilities
	V_1011_Regulatory_Capabilities
	V_1011_Air_Protocol_LLRP_Capabilities
)

/*
This message is issued by the Reader to the Client. This message can be used by the Client to monitor the LLRP-layer connectivity with the Reader. The Client configures the trigger at the Reader to send the Keepalive message. The configuration is done using the KeepaliveSpec parameter
*/
func SEND_KEEPALIVE(messageId int) []byte {
	return bundle(
		M_KEEPALIVE_ACK,
		messageId,
		nil,
	)
}

/*
This message is issued by the Client to the Reader to get the tag reports. In response to this message, the Reader SHALL return tag reports accumulated. If no reports are available to send as a response to a GET_REPORT message, the Reader MAY return an empty RO_ACCESS_REPORT message.
*/
func GET_REPORT(messageId int) []byte {
	return bundle(
		M_GET_REPORT,
		messageId,
		nil,
	)
}

// Convert custom pack to protocol []byte
func CustomPack(messageType, messageId int, config []interface{}, params ...[]interface{}) []byte {
	return bundle(messageType, messageId, config, params...)
}

/*
This message can be issued by the Client to the Reader after a LLRP connection is established. The Client uses this message to inform the Reader that it can remove its hold on event and report messages. Readers that are configured to hold events and reports on reconnection (See Section 13.2.6.4) respond to this message by returning the tag reports accumulated (same way they respond to GET_REPORT (See Section 13.1.1)).

*/
func ENABLE_EVENTS_AND_REPORTS(messageId int) []byte {
	return bundle(
		M_ENABLE_EVENTS_AND_REPORTS,
		messageId,
		nil,
	)
}

/*
This message is issued by the Client to the Reader. Upon receiving the message, the Reader moves the ROSpec corresponding to the ROSpecID passed in this message from the disabled to the inactive state.
*/
func ENABLE_ROSPEC(messageId, ROSpecID int) []byte {
	return bundle(
		M_ENABLE_ROSPEC,
		messageId,
		[]interface{}{
			uint32(ROSpecID),
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

func ADD_ROSPEC(messageId int, params ...[]interface{}) []byte {
	return bundle(
		M_ADD_ROSPEC,
		messageId,
		nil,
		params...,
	)
}
