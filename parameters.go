package llrp

import (
	"fmt"
)

func gPIPortCurrentState_Param(port_number, status int, config bool) []interface{} {
	_config := uint8(0)
	if config {
		_config += 1
	}
	return []interface{}{
		uint16(P_GPIPortCurrentState),
		uint16(8),
		uint16(port_number),
		_config,
		uint8(status),
	}
}

func gPOWriteData_Param(port_number int, data bool) []interface{} {
	_data := uint8(0)
	if data {
		_data += 0x80
	}
	return []interface{}{
		uint16(P_GPOWriteData),
		uint16(7),
		uint16(port_number),
		uint8(_data),
	}
}

// required 9 field
// no. 182
func ReaderEventNotification(config ...bool) []interface{} {
	r := []interface{}{
		uint16(P_ReaderEventNotificationSpec),
		uint16(67),
	}
	if len(config) != 9 {
		panic(fmt.Sprintf("Reader event config notify need 9 field"))
		return nil
	}
	for i := 0; i < 9; i++ {
		val := uint8(0x80)
		if !config[i] {
			val = 0
		}
		r = append(r,
			[]interface{}{
				uint16(P_EventNotificationState),
				uint16(7),
				uint16(i),
				val,
			}...,
		)
	}
	return r
}

// Periodic trigger is specified using UTC time, offset and period
func PeriodicTriggerValueParam(offset, period uint32, utc []interface{}) []interface{} {
	len_ := calcLen(utc)
	r := []interface{}{
		uint16(P_PeriodicTriggerValue),
		uint16(8 + len_),
		offset,
		period,
	}
	r = append(r, utc...)
	return r
}

// UTCTimestamp Parameter
func UTCTimestampParam(microsecond uint64) []interface{} {
	return []interface{}{
		uint16(P_UTCTimeStamp),
		uint16(12),
		microsecond,
	}
}

const (
	P_UTCTimeStamp                                = 128
	P_Uptime                                      = 129
	P_GeneralDeviceCapabilities                   = 137
	P_MaximumReceiveSensitivity                   = 363
	P_ReceiveSensitivityTableEntry                = 139
	P_PerAntennaAirProtocol                       = 140
	P_GPIOCapabilities                            = 141
	P_LLRPCapabilities                            = 142
	P_RegulatoryCapabilities                      = 143
	P_UHFBandCapabilities                         = 144
	P_TransmitPowerLevelTableEntry                = 145
	P_FrequencyInformation                        = 146
	P_FrequencyHopTable                           = 147
	P_FixedFrequencyTable                         = 148
	P_PerAntennaReceiveSensitivityRange           = 149
	P_RFSurveyFrequencyCapabilities               = 365
	P_ROSpec                                      = 177
	P_ROBoundarySpec                              = 178
	P_ROSpecStartTrigger                          = 179
	P_PeriodicTriggerValue                        = 180
	P_GPITriggerValue                             = 181
	P_ROSpecStopTrigger                           = 182
	P_AISpec                                      = 183
	P_AISpecStopTrigger                           = 184
	P_TagObservationTrigger                       = 185
	P_InventoryParameterSpec                      = 186
	P_RFSurveySpec                                = 187
	P_RFSurveySpecStopTrigger                     = 188
	P_LoopSpec                                    = 355
	P_AccessSpec                                  = 207
	P_AccessSpecStopTrigger                       = 208
	P_AccessCommand                               = 209
	P_ClientRequestOpSpec                         = 210
	P_ClientRequestResponse                       = 211
	P_LLRPConfigurationStateValue                 = 217
	P_Identification                              = 218
	P_GPOWriteData                                = 219
	P_KeepaliveSpec                               = 220
	P_AntennaProperties                           = 221
	P_AntennaConfiguration                        = 222
	P_RFReceiver                                  = 223
	P_RFTransmitter                               = 224
	P_GPIPortCurrentState                         = 225
	P_EventsAndReports                            = 226
	P_ROReportSpec                                = 237
	P_TagReportContentSelector                    = 238
	P_AccessReportSpec                            = 239
	P_TagReportData                               = 240
	P_EPCData                                     = 241
	P_EPC_96                                      = 13
	P_ROSpecID                                    = 9
	P_SpecIndex                                   = 14
	P_InventoryParameterSpecID                    = 10
	P_AntennaID                                   = 1
	P_PeakRSSI                                    = 6
	P_ChannelIndex                                = 7
	P_FirstSeenTimestampUTC                       = 2
	P_FirstSeenTimestampUptime                    = 3
	P_LastSeenTimestampUTC                        = 4
	P_LastSeenTimestampUptim                      = 5
	P_TagSeenCount                                = 8
	P_ClientRequestOpSpecResult                   = 15
	P_AccessSpecID                                = 16
	P_RFSurveyReportData                          = 242
	P_FrequencyRSSILevelEntry                     = 243
	P_ReaderEventNotificationSpec                 = 244
	P_EventNotificationState                      = 245
	P_ReaderEventNotificationData                 = 246
	P_HoppingEvent                                = 247
	P_GPIEvent                                    = 248
	P_ROSpecEvent                                 = 249
	P_ReportBufferLevelWarningEvent               = 250
	P_ReportBufferOverflowErrorEvent              = 251
	P_ReaderExceptionEvent                        = 252
	P_OpSpecID                                    = 17
	P_RFSurveyEvent                               = 253
	P_AISpecEvent                                 = 254
	P_AntennaEvent                                = 255
	P_ConnectionAttemptEvent                      = 256
	P_ConnectionCloseEvent                        = 257
	P_SpecLoopEvent                               = 356
	P_LLRPStatus                                  = 287
	P_FieldError                                  = 288
	P_ParameterError                              = 289
	P_Custom                                      = 1023
	P_C1G2LLRPCapabilities                        = 327
	P_UHFC1G2RFModeTable                          = 328
	P_UHFC1G2RFModeTableEntry                     = 329
	P_C1G2InventoryCommand                        = 330
	P_C1G2Filter                                  = 331
	P_C1G2TagInventoryMask                        = 332
	P_C1G2TagInventoryStateAwareFilterAction      = 333
	P_C1G2TagInventoryStateUnawareFilterAction    = 334
	P_C1G2RFControl                               = 335
	P_C1G2SingulationControl                      = 336
	P_C1G2TagInventoryStateAwareSingulationAction = 337
	P_C1G2TagSpec                                 = 338
	P_C1G2TargetTag                               = 339
	P_C1G2Read                                    = 341
	P_C1G2Write                                   = 342
	P_C1G2Kill                                    = 343
	P_C1G2Recommission                            = 357
	P_C1G2Lock                                    = 344
	P_C1G2LockPayload                             = 345
	P_C1G2BlockErase                              = 346
	P_C1G2BlockWrite                              = 347
	P_C1G2BlockPermalock                          = 358
	P_C1G2GetBlockPermalockStatus                 = 359
	P_C1G2EPCMemorySelector                       = 348
	P_C1G2PC                                      = 12
	P_C1G2XPCW1                                   = 19
	P_C1G2XPCW2                                   = 20
	P_C1G2CRC                                     = 11
	P_C1G2SingulationDetails                      = 18
	P_C1G2ReadOpSpecResult                        = 349
	P_C1G2WriteOpSpecResult                       = 350
	P_C1G2KillOpSpecResult                        = 351
	P_C1G2RecommissionOpSpecResult                = 360
	P_C1G2LockOpSpecResult                        = 352
	P_C1G2BlockEraseOpSpecResult                  = 353
	P_C1G2BlockWriteOpSpecResult                  = 354
	P_C1G2BlockPermalockOpSpecResult              = 361
	P_C1G2GetBlockPermalockStatusOpSpecResult     = 362
)
