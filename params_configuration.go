package llrp

/*
This parameter, LLRPConfigurationStateValue, is a 32-bit value which represents a
Reader’s entire LLRP configuration state including: LLRP configuration parameters,
vendor extension configuration parameters, ROSpecs, and AccessSpecs. A Reader
SHALL change this value only:
• Upon successful execution of any of the following messages:
	o ADD_ROSPEC
	o DELETE_ROSPEC
	o ADD_ACCESSSPEC
	o DELETE_ACCESSSPEC
	o SET_READER_CONFIG
	o Any CUSTOM_MESSAGE command that alters the reader’s internal configuration.
• Upon an automatically deleted AccessSpec due to OperationCountValue number of operations (Section 12.2.1.1).
A Reader SHALL not change this value when the CurrentState of a ROSpec or AccessSpec changes.
The mechanism used to compute the LLRP configuration state value is implementation
dependent. However, a good implementation will insure that there’s a high probability
that the value will change when the Reader’s configuration state changes.
It is expected that a Client will configure the Reader and then request the Reader’s
configuration state value. The Client will then save this state value. If this value does not
change between two requests for it, then a Client may assume that the above components
of the LLRP configuration have also not changed.
*/
func LLRPConfigurationStateValue(llrpConfigurationStateValue int) []interface{} {
	return commonSpec(
		P_LLRPConfigurationStateValue,
		[]interface{}{
			uint32(llrpConfigurationStateValue),
		},
	)
}

// IDType used by Identification Parameter
const (
	C_Identification_IDType_MAC = iota
	C_Identification_IDType_EPC
)

// Reader ID: Byte array. If IDType=0, the MAC address SHALL be encoded as EUI-64.[EUI64]
func Identification(IDType int, ReaderId string) []interface{} {
	b := []uint8(ReaderId)
	inf := []interface{}{
		uint8(IDType),
		uint16(len(b)),
	}
	for _, k := range b {
		inf = append(inf, k)
	}

	return commonSpec(
		P_Identification,
		inf,
	)
}

// This parameter carries the data pertinent to perform the write to a general purpose output port.
/*
	GPO Port Number : Unsigned Short Integer. 0 is invalid.
	GPO Data: Boolean. The state to output on the specified GPO port.
*/
func GPOWriteDataFunc(GPOPortNumber int, GPOData bool) []interface{} {
	data := uint8(0)
	if GPOData {
		data += 0x80
	}
	return commonSpec(
		P_GPOWriteData,
		[]interface{}{
			uint16(GPOPortNumber),
			data,
		},
	)
}

// This parameter carries the specification for the keepalive message generation by the Reader. This includes the definition of the periodic trigger to send the keepalive message
// PeriodicTriggerValue: Integer. Time interval in milliseconds. This field is ignored when KeepaliveTriggerType is not 1.
func KeepaliveSpec(PeriodicTriggerValue int) []interface{} {

	typeof := uint8(0)
	if PeriodicTriggerValue > 0 {
		typeof += 1
	}
	return commonSpec(
		P_KeepaliveSpec,
		[]interface{}{
			typeof,
			uint32(PeriodicTriggerValue),
		},
	)
}

/*
	This parameter carries a single antenna's properties. The properties include the gain and the connectivity status of the antenna.The antenna gain is the composite gain and includes the loss of the associated cable from the Reader to the antenna. The gain is represented in dBi*100 to allow fractional dBi representation.
*/
func AntennaProperties(AntennaID, AntennaGain int, AntennaConnected bool) []interface{} {
	and := uint8(0)
	if AntennaConnected {
		and += 0x80
	}
	return commonSpec(
		P_AntennaProperties,
		[]interface{}{
			and,
			uint16(AntennaID),
			uint16(AntennaGain),
		},
	)
}

/*
	This parameter carries a single antenna's configuration and it specifies the default values for the parameter set that are passed in this parameter block. The scope of the default values is the antenna. The default values are used for parameters during an operation on this antenna if the parameter was unspecified in the spec that describes the operation.
*/
// params = RFReceiver Parameter , RFTransmitter Parameter , AirProtocolInventoryCommandSettings Parameter , Custom Parameter
func AntennaConfiguration(AntennaID int, params ...[]interface{}) []interface{} {
	return commonSpec(
		P_AntennaConfiguration,
		[]interface{}{
			uint16(AntennaID),
		},
		params...,
	)
}

/*
This Parameter carries the RF transmitter information. The Transmit Power defines the transmit power for the antenna expressed as an index into the TransmitPowerTable (section 10.2.4.1.1). The HopTableID is the index of the frequency hop table to be used by the Reader (section 10.2.4.1.2.1) and is used when operating in frequency-hopping regulatory regions. This field is ignored in non-frequency-hopping regulatory regions.
The ChannelIndex is the one-based channel index in the FixedFrequencyTable to use during transmission (section 10.2.4.1.2.2) and is used when operating in non-frequency-hopping regulatory regions. This field is ignored in frequency-hopping regulatory regions.
*/
func RFReceiver(TransmitPower, HopTableID, ChannelIndex int) []interface{} {
	return commonSpec(
		P_RFReceiver,
		[]interface{}{
			uint16(HopTableID),
			uint16(ChannelIndex),
			uint16(TransmitPower),
		},
	)
}

/*
This Parameter carries the current configuration and state of a single GPI port. In a SET_READER_CONFIG message, this parameter is used to enable or disable the GPI port using the GPIConfig field; the GPIState field is ignored by the reader. In a GET_READER_CONFIG message, this parameter reports both the configuration and state of the GPI port. When a ROSpec or AISpec is configured on a GPI-capable reader with GPI start and/or stop triggers, those GPIs must be enabled by the client with a SET_READER_CONFIG message for the triggers to function.
*/
func GPIPortCurrentState(GPIPortNum, GPIState int, GPIConfig bool) []interface{} {
	config := uint8(0)
	if GPIConfig {
		config += 0x80
	}
	return commonSpec(
		P_GPIPortCurrentState,
		[]interface{}{
			uint16(GPIPortNum),
			config,
			uint8(GPIState),
		},
	)
}

/*
This parameter controls the behavior of the Reader when a new LLRP connection is established. In a SET_READER_CONFIG message, this parameter is used to enable or
disable the holding of events and reports upon connection using the HoldEventsAndReportsUponReconnect field. In a GET_READER_CONFIG message,this parameter reports the current configuration. If the ldEventsAndReportsUponReconnect is true, the reader will not deliver any reports or events (except the ConnectionAttemptEvent) to the Client until the Client issues an ENABLE_EVENTS_AND_REPORTS message. Once the ENABLE_EVENTS_AND_REPORTS message is received the reader ceases its hold on events and reports for the duration of the connection.
*/
func EventsAndReports(HoldEventsAndReportsUponReconnect bool) []interface{} {
	return commonSpec(
		P_EventsAndReports,
		[]interface{}{
			convertBooleanUint8(HoldEventsAndReportsUponReconnect),
		},
	)
}

/*
This Parameter carries the RF transmitter information. The Transmit Power defines the transmit power for the antenna expressed as an index into the TransmitPowerTable (section 10.2.4.1.1). The HopTableID is the index of the frequency hop table to be used
by the Reader (section 10.2.4.1.2.1) and is used when operating in frequency-hopping regulatory regions. This field is ignored in non-frequency-hopping regulatory regions. The ChannelIndex is the one-based channel index in the FixedFrequencyTable to use during transmission (section 10.2.4.1.2.2) and is used when operating in non-frequency-hopping regulatory regions. This field is ignored in frequency-hopping regulatory regions.

*/
func RFTransmitter(HopTableID, ChannelIndex, TransmitPower int) []interface{} {

	return commonSpec(
		P_RFTransmitter,
		[]interface{}{
			uint16(HopTableID),
			uint16(ChannelIndex),
			uint16(TransmitPower),
		},
	)
}
