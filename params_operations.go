package llrp

/*
This parameter carries the information of the Reader inventory and survey operation.
*/
// params = ROBoundarySpec Parameter , SpecParameter (1-n) , ROReportSpec Parameter
func ROSpec(ROSpecID, Priority, CurrentState int, params ...[]interface{}) []interface{} {
	return commonSpec(
		P_ROSpec,
		[]interface{}{
			uint32(ROSpecID),
			uint8(Priority),
			uint8(CurrentState),
		},
		params...,
	)

}

// params = AISpecStopTrigger , InventoryParameter Spec and Custom Parameter
func AISpec(AntennaCount int, AntennaIDn []int, params ...[]interface{}) []interface{} {
	inf := []interface{}{
		uint16(AntennaCount),
	}
	for _, k := range AntennaIDn {
		inf = append(inf, uint16(k))
	}
	return commonSpec(
		P_AISpec,
		inf,
		params...,
	)
}

const (
	C_AISpecStopTrigger_NULL = iota
	C_AISpecStopTrigger_DURATION
	C_AISpecStopTrigger_GPI_WITH_TIMEOUT
	C_AISpecStopTrigger_TAG_OBSERVATION
)

// params = GPITriggerValue Parameter , TagObservationTrigger Parameter
func AISpecStopTrigger(AISpecStopTriggerType, DurationTrigger int, params ...[]interface{}) []interface{} {
	return commonSpec(
		P_AISpecStopTrigger,
		[]interface{}{
			uint8(AISpecStopTriggerType),
			uint32(DurationTrigger),
		},
		params...,
	)
}

/*
	TriggerType: Integer
	value : 	 modulation
		0	Upon seeing N tag observations, or timeout. The definition of an "observation" is vendor specific.
		1	Upon seeing no more new tag observations for T ms, or timeout. The definition of an "observation" is
		vendor specific.
		2	N attempts to see all tags in the FOV, or timeout.
		3	Upon seeing N unique tag observations, or timeout.
		4	Upon seeing no more new unique tag observations for T ms, or timeout.
	-----
	NumberOfTags: Unsigned Short Integer. This field SHALL be ignored when
	TriggerType != 0 and TriggerType != 3.
	NumberOfAttempts; Unsigned Short Integer. This field SHALL be ignored when
	TriggerType != 2.
	T : Unsigned Short Integer. Idle time between tag responses in milliseconds. This field
	SHALL be ignored when TriggerType != 1 and TriggerType != 4.
	Timeout : Unsigned Integer; Trigger timeout value in milliseconds. If set to zero, it
	indicates that there is no timeout.
*/
func TagObservationTrigger(TriggerType, NumberOfTags, NumberOfAttempts, T, Timeout int) []interface{} {
	return commonSpec(
		P_TagObservationTrigger,
		[]interface{}{
			uint8(TriggerType),
			uint16(NumberOfTags),
			uint16(T),
			uint32(Timeout),
		},
	)
}

// Operational parameters for an inventory using a single air protocol.
// params = AntennaConfigurationParameter , Custom Parameter
func InventoryParameterSpec(InventoryParameterSpecID, ProtocolID int, params ...[]interface{}) []interface{} {
	return commonSpec(
		P_InventoryParameterSpec,
		[]interface{}{
			uint16(InventoryParameterSpecID),
			uint8(ProtocolID),
		},
		params...,
	)
}

// Details of a RF Survey operation
// params = RFSurveySpecStopTrigger , Custom Parameter
func RFSurveySpec(AntennaID, StartFrequency, EndFrequency int, params ...[]interface{}) []interface{} {
	return commonSpec(
		P_RFSurveySpec,
		[]interface{}{
			uint16(AntennaID),
			uint32(StartFrequency),
			uint32(EndFrequency),
		},
		params...,
	)
}

// RFSurveySpecStopTrigger Type
const (
	C_RFSurveySpecStopTrigger_NULL = iota
	C_RFSurveySpecStopTrigger_Duration
	C_RFSurveySpecStopTrigger_N_Iteration
)

/*
Duration: Unsigned Integer; The maximum duration of the RFSurvey operation
specified in milliseconds. This field SHALL be ignored when StopTriggerType != 1.
When StopTriggerType = 1, the value SHALL be greater than zero.
N: Unsigned Integer. The maximum number of iterations through the specified
frequency range. This field SHALL be ignored when StopTriggerType != 2. When
StopTriggerType = 2, the value SHALL be greater than zero.
*/
func RFSurveySpecStopTrigger(StopTriggerType, Duration, N int) []interface{} {
	return commonSpec(
		P_RFSurveySpecStopTrigger,
		[]interface{}{
			uint8(StopTriggerType),
			uint32(Duration),
			uint32(N),
		},
	)
}

// Instructs the Reader to execute the first Spec in the Set of Specs.
// LoopCount: This value instructs the reader on the number of times to loop through the Set of Specs within the ROSpec.
func LoopSpec(LoopCount int) []interface{} {
	return commonSpec(
		P_LoopSpec,
		[]interface{}{
			uint32(LoopCount),
		},
	)
}

//  ROSpecStartTrigger Parameter, ROSpecStopTrigger Parameter
func ROBoundarySpec(params ...[]interface{}) []interface{} {
	return commonSpec(
		P_ROBoundarySpec,
		nil,
		params...,
	)
}

// PeriodicTriggerValue Parameter , GPITriggerValue Parameter
func ROSpecStartTrigger(ROSpecStartTriggerType int, params ...[]interface{}) []interface{} {
	return commonSpec(
		P_ROSpecStartTrigger,
		[]interface{}{
			uint8(ROSpecStartTriggerType),
		},
		params...,
	)
}

// UTCTimestamp Parameter
func PeriodicTriggerValue(Offset, Period int, params ...interface{}) []interface{} {
	return commonSpec(
		P_PeriodicTriggerValue,
		[]interface{}{
			uint32(Offset),
			uint32(Period),
		},
		params,
	)
}
func GPITriggerValue(GPIPortNum int, GPIEvent bool, Timeout int) []interface{} {
	gp := 0
	if GPIEvent {
		gp = 0x80
	}
	return commonSpec(
		P_GPITriggerValue,
		[]interface{}{
			uint16(GPIPortNum),
			uint8(gp),
			uint32(Timeout),
		},
	)
}

// GPITriggerValue Parameter
func ROSpecStopTrigger(ROSpecStopTriggerType int, DurationTriggerValue int, params ...[]interface{}) []interface{} {
	return commonSpec(
		P_ROSpecStopTrigger,
		[]interface{}{
			uint8(ROSpecStopTriggerType),
			uint32(DurationTriggerValue),
		},
		params...,
	)
}

/*
This parameter defines the C1G2 inventory-specific settings to be used during a particular C1G2 inventory operation. This comprises of C1G2Filter Parameter, C1G2RF Parameter and C1G2Singulation Parameter. It is not necessary that the Filter, RF Control and Singulation Control Parameters be specified in each and every inventory command.They are optional parameters. If not specified, the default values in the Reader are used during the inventory operation. If multiple C1G2Filter parameters are encapsulated by the Client in the C1G2InventoryCommand parameter, the ordering of the filter parameters determine the order of C1G2 air-protocol commands (e.g., Select command) generated by the Reader. C1G2Filter parameters included in the C1G2InventoryCommand parameter
*/
func C1G2InventoryCommand(TagInventoryStateAware bool, params ...[]interface{}) []interface{} {
	return commonSpec(
		P_C1G2InventoryCommand,
		[]interface{}{
			convertBooleanUint8(TagInventoryStateAware),
		},
		params...,
	)
}

/*
This Parameter carries the settings relevant to RF forward and reverse link control in the C1G2 air protocol. This is basically the C1G2 RF Mode and the Tari value to use for the inventory operation.
 ---
 ModeIndex: Unsigned Integer. This is an index into the UHFC1G2RFModeTable.
 Tari: Integer. Value of Tari to use for this mode specified in nsec. This is specified if the mode selected has a Tari range. If the selected mode has a range, and the Tari is set to zero, the Reader implementation picks up any Tari value within the range. If the selected mode has a range, and the specified Tari is out of that range and is not set to zero, an error message is generated.
 Possible Values:
 0 or 6250-25000 nsec
*/
func C1G2RFControl(ModeIndex, Tari int) []interface{} {
	return commonSpec(
		P_C1G2RFControl,
		[]interface{}{
			uint16(ModeIndex),
			uint16(Tari),
		},
	)
}

/*
This C1G2SingulationControl Parameter provides controls particular to the singulation process in the C1G2 air protocol. The singulation process is started using a Query command in the C1G2 protocol. The Query command describes the session number, tag state, the start Q value to use, and the RF link parameters. The RF link parameters are specified using the C1G2RFControl Parameter (see section 16.2.1.2.1.2). This Singulation Parameter specifies the session, tag state and description of the target singulation environment
-----
	Session: Integer. Session number to use for the inventory operation.
	Possible Values: 0-3
	Tag population: Unsigned Short Integer. An estimate of the tag population in view of the RF field of the antenna.
	Tag transit time: Unsigned Integer. An estimate of the time a tag will typically remain in the RF field of the antenna specified in milliseconds.
	TagInventoryStateAwareSingulationAction: <C1G2TagInventoryStateAwareSingulationAction Parameter> (optional)
	params = C1G2TagInventoryStateAwareSingulationAction
*/
func C1G2SingulationControl(Session, TagPopulation, TagTransitTime int, params ...[]interface{}) []interface{} {
	return commonSpec(
		P_C1G2SingulationControl,
		[]interface{}{
			uint8(Session),
			uint16(TagPopulation),
			uint32(TagTransitTime),
		},
		params...,
	)
}
