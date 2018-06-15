package llrp

// This parameter is used to enable or disable notification of a single Reader event type.
func EventNotificationStateParam(EventType int, NotificationState bool) []interface{} {
	return commonSpec(
		P_EventNotificationState,
		[]interface{}{
			uint16(EventType),
			convertBooleanUint8(NotificationState),
		},
	)
}

/*
This parameter is used by the Client to enable or disable notification of one or more Reader events. Notification of buffer overflow events and connection events (attempt/close) are mandatory, and not configurable.
*/
// params = List of <EventNotificationState Parameter>
func ReaderEventNotificationSpec(params ...[]interface{}) []interface{} {
	if len(params) == 0 {
		// default
		default_ := map[int]bool{
			0: true,
			1: true,
			2: true,
			3: true,
			4: true,
			5: false,
			6: true,
			7: false,
			8: true,
		}
		for v, k := range default_ {
			params = append(
				params,
				EventNotificationStateParam(v, k),
			)
		}
	}
	return commonSpec(
		P_ReaderEventNotificationSpec,
		nil,
		params...,
	)
}

/*
This Parameter carries the Reader inventory and RF survey reporting definition for the antenna. This parameter describes the contents of the report sent by the Reader and defines the events that cause the report to be sent.
*/
func RoReportSpec(ROReportTrigger, N int, params ...[]interface{}) []interface{} {

	return commonSpec(
		P_ROReportSpec,
		[]interface{}{
			uint8(ROReportTrigger),
			uint16(N),
		},
		params...,
	)
}

// N = 0 is unlimited
// ROReportTrigger type
const (
	C_ROReportTrigger_None = iota
	C_ROReportTrigger_Upon_N_TagReportData_or_EndOfAISpec
	C_ROReportTrigger_Upon_N_TagReportData_or_EndOfROSpec
	C_ROReportTrigger_Upon_N_Seconds_or_EndAISpec_EndROSpec
	C_ROReportTrigger_Upon_N_Seconds_or_EndROSpec
	C_ROReportTrigger_Upon
	C_ROReportTrigger_Upon_N_milliseconds_or_EndROSpec_EndRFSurveyspec
	C_ROReportTrigger_Upon_N_milliseconds_or_EndROSpec
)

/*
This Parameter carries the Reader inventory and RF survey reporting definition for the antenna. This parameter describes the contents of the report sent by the Reader and defines the events that cause the report to be sent.

The ROReportTrigger field defines the events that cause the report to be sent.

The TagReportContentSelector parameter defines the desired contents of the report. The ROReportTrigger defines the event that causes the report to be sent by the Reader to the Client.

Custom extensions to this parameter are intended to specify summary data to be reported as an extension to the RO_ACCESS_REPORT message (see section 14.1.2).
*/
// params = <TagReportContentSelector Parameter> , List of <Custom Parameter>
func ROReportSpec(ROReportTrigger, N int, params ...[]interface{}) []interface{} {
	return commonSpec(
		P_ROReportSpec,
		[]interface{}{
			uint8(ROReportTrigger),
			uint16(N),
		},
		params...,
	)
}

// This parameter is used to configure the contents that are of interest in TagReportData. If enabled, the field is reported along with the tag data in the TagReportData.
func TagReportContentSelector(
	EnableROSpecID,
	EnableSpecIndex,
	EnableInventoryParameterSpecID,
	EnableAntennaID,
	EnableChannelIndex,
	EnablePeakRSSI,
	EnableFirstSeenTimestamp,
	EnableLastSeenTimestamp,
	EnableTagSeenCount,
	EnableAccessSpecID bool,
	params ...[]interface{},
) []interface{} {
	return commonSpec(
		P_TagReportContentSelector,
		[]interface{}{
			convert16uintbit(
				EnableROSpecID,
				EnableSpecIndex,
				EnableInventoryParameterSpecID,
				EnableAntennaID,
				EnableChannelIndex,
				EnablePeakRSSI,
				EnableFirstSeenTimestamp,
				EnableLastSeenTimestamp,
				EnableTagSeenCount,
				EnableAccessSpecID,
			),
		},
		params...,
	)
}

// AccessReportTrigger
/*
	0 : Whenever ROReport is generated for the RO that triggered the execution of this AccessSpec.
	1 : End of AccessSpec (immediately upon completionof the access operation)
*/
const (
	C_AccessReportTrigger_GenTriggerd = iota
	C_AccessReportTrigger_EndOfAcessSpec
)

// This parameter sets up the triggers for the Reader to send the access results to the Client. In addition, the Client can enable or disable reporting of ROSpec details in the access results.
func AccessReportSpec(AccessReportTrigger int) []interface{} {
	return commonSpec(
		P_AccessReportSpec,
		[]interface{}{
			uint8(AccessReportTrigger),
		},
	)
}
