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
type EvtData struct {
	TimestampUTC uint64
	GpiEvt       []*GpiEventParam
}
type GpiEventParam struct {
	PortNumber uint16
	Evt        bool
}
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
type GetConfigResponse struct {
	MsgId  uint32
	Status *LLRPStatus
	Id     *Identification
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

type Identification struct {
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
