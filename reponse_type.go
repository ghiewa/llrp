package llrp

type LLRPStatus struct {
	Success    bool
	StatusCode uint16
	ErrMsg     string
}
type ERROR_MESSAGE struct {
	MsgId  uint32
	Status *LLRPStatus
}
type ROAccessReportResponse struct {
	MsgId uint32
	Data  *TagReportData
}
type SetConfigResponse struct {
	MsgId  uint32
	Status *LLRPStatus
}
type CUSTOM_MESSAGE_RESPONSE struct {
	MsgId   uint32
	Vendor  uint32
	SubType uint8
	Status  *LLRPStatus
}
type ENABLE_ROSPEC_RESPONSE struct {
	MsgId  uint32
	Status *LLRPStatus
}
type EventNotificationResponse struct {
	MsgId uint32
	Data  *EvtData
}

//  ReaderEventNotificationData
type EvtData struct {
	TimestampUTC      uint64
	Hopping           []*HoppingEventParameter
	GpiEvt            []*GpiEventParam
	ROSpEvt           []*ROSpecEventParameter
	ReportBuf         []*ReportBufferLevelWarningEventParameter
	ReaderException   []*ReaderExceptionEventParameter
	RFSurvey          []*RFSurveyEventParameter
	AISpec            []*AISpecEventParameter
	Antenna           []*AntennaEventParameter
	ConnectionAttempt []*ConnectionAttemptEventParameter
}

const (
	C_AntennaEventParameter_Antenna_Disconnected = iota
	C_AntennaEventParameter_Antenna_Connected
)

type AntennaEventParameter struct {
	EventType uint8
	AntennaID uint16
}

const (
	C_ConnectionAttemptEvent_Status_Success = iota
	C_ConnectionAttemptEvent_Status_Failed_Reader_Conn_Exists
	C_ConnectionAttemptEvent_Status_Failed_Client_Conn_Exists
	C_ConnectionAttemptEvent_Status_Another_Conn_Attempted
)

type ConnectionAttemptEventParameter struct {
	Status uint16
}
type HoppingEventParameter struct {
	HopTableID       uint16
	NextChannelIndex uint16
}
type GpiEventParam struct {
	PortNumber uint16
	Evt        bool
}

const (
	C_ROSpecEvent_Type_Start_ROSpec = iota
	C_ROSpecEvent_Type_End_ROSpec
	C_ROSpecEvent_Type_Preemption_ROSpec
)

// This parameter carries the ROSpec event details. The EventType could be start or end of the ROSpec.
type ROSpecEventParameter struct {
	EventType          uint8
	ROSpecID           uint32
	PreemptingROSpecID uint32
}
type ReportBufferLevelWarningEventParameter struct {
	ReportBufferPercentageFull uint8
}
type ReportBufferOverflowErrorEvent struct {
}

type ReaderExceptionEventParameter struct {
	ROSpecID                 []*ROSpecIDParameter
	SpecIndex                []*SpecIndexParameter
	InventoryParameterSpecID []*InventoryParameterSpecIDParameter
	AntennaID                []*AntennaIDParameter
	AccessSpecID             []*AccessSpecIDParameter
	OpSpecID                 []*OpSpecIDParameter
	Message                  string
	//CustomExtensionPointList []*CustomParameterResp
}
type ROSpecIDParameter struct {
	ROSpecID uint32
}
type SpecIndexParameter struct {
	SpecIndex uint16
}
type InventoryParameterSpecIDParameter struct {
	InventoryParameterSpecId uint16
}
type AntennaIDParameter struct {
	AntennaID uint16
}
type PeakRSSIParameter struct {
	PeakRSSI uint8
}
type ChannelIndexParameter struct {
	ChannelIndex uint16
}
type AccessSpecIDParameter struct {
}

type OpSpecIDParameter struct {
	OpSpecId uint16
}
type RFSurveyEventParameter struct {
	EventType uint8
	ROSpecID  uint32
	SpecIndex uint16
}
type AISpecEventParameter struct {
	EventType uint8
	ROSpecID  uint32
	SpecIndex uint16
	// air prtocols
}

//-------------------
type DELETE_ROSPEC_RESPONSE struct {
	MsgId  uint32
	Status *LLRPStatus
}
type DELETE_ACCESSSPEC_RESPONSE struct {
	MsgId  uint32
	Status *LLRPStatus
}
type ADD_ROSPEC_RESPONSE struct {
	MsgId  uint32
	Status *LLRPStatus
}
type DISABLE_ROSPEC_RESPONSE struct {
	MsgId  uint32
	Status *LLRPStatus
}
type START_ROSPEC_RESPONSE struct {
	MsgId  uint32
	Status *LLRPStatus
}
type STOP_ROSPEC_RESPONSE struct {
	MsgId  uint32
	Status *LLRPStatus
}
type CLOSE_CONNECTION_RESPONSE struct {
	MsgId  uint32
	Status *LLRPStatus
}
type GetConfigResponse struct {
	MsgId  uint32
	Status *LLRPStatus
	Id     *IdentificationParam
	AnnPty []*AntennaProperty
	//AnnCon    []*AntennaConfig
	EvtSpec *ReaderEventNoticationSpec
	//ROSpec    *ROReportSpec
	//AccSpec   *AccessReportSpec
	//CfgState  *LLRPConfigState
	Ka        *KeepaliveResponse
	GPI       []*GPICurrentState
	GPO       []*GPOWriteData
	EvtReport *EventReport
	//Customs   []*CustomParamer
}
type EventReport struct {
	HoldEventReportsUponReconnect bool
}

type IdentificationParam struct {
	Type int
	Id   string
}
type AntennaProperty struct {
	Connected bool
	Id        uint16
	Gain      uint16
}
type ReaderEventNoticationSpec struct {
	State []*EventNotificationState
}
type EventNotificationState struct {
	EventType uint16
	State     int
}
type KeepaliveResponse struct {
	Type         int
	IntervalTime uint32
}

type GPICurrentState struct {
	Number uint16
	Config bool
	State  int
}
type GPOWriteData struct {
	Number uint16
	Data   bool
}
type TagReportData struct {
	EPC_96                string
	AntennaID             uint16
	PeakRSSI              int
	ChannelIndex          uint16
	FirstSeenTimestampUTC uint64
	AccessSpecId          uint32
}
type MsgLoss struct {
	Len int
}
