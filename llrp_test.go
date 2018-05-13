package llrp

import (
	//"encoding/binary"
	"bytes"
	//"fmt"
	"testing"
)

func TestConvertNotify(t *testing.T) {
	// no. 1723
	v := []byte{
		0x04, 0x3d, 0x00, 0x00, 0x00, 0x31, 0x2e, 0x5a, 0x3d, 0x68, 0x00, 0xf0, 0x00, 0x27, 0x8d, 0xe2,
		0x00, 0x20, 0x75, 0x69, 0x09, 0x01, 0x46, 0x28, 0x00, 0x08, 0x07, 0x81, 0x00, 0x01, 0x86, 0xda,
		0x87, 0x00, 0x04, 0x82, 0x00, 0x05, 0x67, 0x60, 0x73, 0x19, 0xf0, 0x7a, 0x90, 0x00, 0x00, 0x00,
		0x00,
		0x04, 0x3d, 0x00, 0x00, 0x00, 0x31, 0x2e, 0x5a, 0x3d, 0x69, 0x00, 0xf0, 0x00, 0x27, 0x8d, 0xe2,
		0x00, 0x20, 0x75, 0x69, 0x09, 0x01, 0x46, 0x28, 0x00, 0x08, 0x07, 0x81, 0x00, 0x01, 0x86, 0xda,
		0x87, 0x00, 0x04, 0x82, 0x00, 0x05, 0x67, 0x60, 0x73, 0x1a, 0x6e, 0x6d, 0x90, 0x00, 0x00, 0x00,
		0x00,
		0x04, 0x3d, 0x00, 0x00, 0x00, 0x31, 0x2e, 0x5a, 0x3d, 0x6a, 0x00, 0xf0, 0x00, 0x27, 0x8d, 0xe2,
		0x00, 0x20, 0x75, 0x69, 0x09, 0x01, 0x46, 0x28, 0x00, 0x08, 0x07, 0x81, 0x00, 0x01, 0x86, 0xda,
		0x87, 0x00, 0x04, 0x82, 0x00, 0x05, 0x67, 0x60, 0x73, 0x1a, 0xf1, 0xe9, 0x90, 0x00, 0x00, 0x00,
		0x00,
		0x04, 0x3d, 0x00, 0x00, 0x00, 0x31, 0x2e, 0x5a, 0x3d, 0x6b, 0x00, 0xf0, 0x00, 0x27, 0x8d, 0xe2,
		0x00, 0x20, 0x75, 0x69, 0x09, 0x01, 0x46, 0x28, 0x00, 0x08, 0x07, 0x81, 0x00, 0x01, 0x86, 0xda,
		0x87, 0x00, 0x04, 0x82, 0x00, 0x05, 0x67, 0x60, 0x73, 0x1b, 0x72, 0xc5, 0x90, 0x00, 0x00, 0x00,
		0x00,
		0x04, 0x3d, 0x00, 0x00, 0x00, 0x31, 0x2e, 0x5a, 0x3d, 0x6c, 0x00, 0xf0, 0x00, 0x27, 0x8d, 0xe2,
		0x00, 0x20, 0x75, 0x69, 0x09, 0x01, 0x46, 0x28, 0x00, 0x08, 0x07, 0x81, 0x00, 0x01, 0x86, 0xda,
		0x87, 0x00, 0x04, 0x82, 0x00, 0x05, 0x67, 0x60, 0x73, 0x1c, 0x1d, 0x09, 0x90, 0x00, 0x00, 0x00,
		0x00,
		0x04, 0x3f, 0x00, 0x00, 0x00, 0x22, 0x2e, 0x5a, 0x3d, 0x6d, 0x00, 0xf6, 0x00, 0x18, 0x00, 0x80,
		0x00, 0x0c, 0x00, 0x05, 0x67, 0x60, 0x73, 0x1c, 0x8d, 0x5c, 0x00, 0xf7, 0x00, 0x08, 0x00, 0x01,
		0x00, 0x05,
	}
	res := Response(v, len(v))
	t.Logf("\nlen : %d", len(res))
	for _, k := range res {
		switch k.(type) {
		case *ROAccessReportResponse:
			kk := k.(*ROAccessReportResponse)
			t.Logf("\n[ROAccessReportResponse] %d : %s", kk.MsgId, kk.Data.EPC_96)
		case *EventNotificationResponse:
			kk := k.(*EventNotificationResponse)
			t.Logf("\n[EventNotificationResponse] %d  : %d", kk.MsgId, kk.Data.TimestampUTC)
		default:
			t.Errorf("not found type %s", k)
		}
	}

}
func TestCustomMsg(t *testing.T) {
	// no. 162
	v := []byte{
		0x07, 0xff, 0x00, 0x00, 0x00, 0x13, 0x00, 0x00, 0x00, 0x51, 0x00, 0x00, 0x65, 0x1a, 0x15, 0x00,
		0x00, 0x00, 0x00,
	}
	messageId = 81
	messageType := M_CUSTOM_MESSAGE
	b := CustomPack(messageType, messageId,
		[]interface{}{
			uint32(25882),
			uint8(21),
			uint32(0),
		},
	)
	if !bytes.Equal(v, b) {
		t.Errorf("want:% x", v)
		t.Errorf("resp:% x", b)
	}

}
func TestCustomMsgRes(t *testing.T) {
	// res
	v := []byte{
		0x07, 0xff, 0x00, 0x00, 0x00, 0x17, 0x00, 0x00, 0x00, 0x51, 0x00, 0x00, 0x65, 0x1a, 0x16, 0x01,
		0x1f, 0x00, 0x08, 0x00, 0x00, 0x00, 0x00,
	}
	res := Response(v, 23)
	kk := res[0].(*CUSTOM_MESSAGE_RESPONSE)
	if !(kk.Vendor == 25882 && kk.MsgId == 81 &&
		kk.SubType == 22 && kk.Status.Success) {
		t.Errorf("\n%+v %+v", kk, kk.Status)
	}

}
func TestSendNotifyConfig(t *testing.T) {
	// no. 182
	v := []byte{
		0x04, 0x03, 0x00, 0x00, 0x00, 0x4e, 0x00, 0x00, 0x00, 0x96, 0x00, 0x00, 0xf4, 0x00, 0x43, 0x00,
		0xf5, 0x00, 0x07, 0x00, 0x00, 0x80, 0x00, 0xf5, 0x00, 0x07, 0x00, 0x01, 0x80, 0x00, 0xf5, 0x00,
		0x07, 0x00, 0x02, 0x80, 0x00, 0xf5, 0x00, 0x07, 0x00, 0x03, 0x80, 0x00, 0xf5, 0x00, 0x07, 0x00,
		0x04, 0x80, 0x00, 0xf5, 0x00, 0x07, 0x00, 0x05, 0x00, 0x00, 0xf5, 0x00, 0x07, 0x00, 0x06, 0x80,
		0x00, 0xf5, 0x00, 0x07, 0x00, 0x07, 0x00, 0x00, 0xf5, 0x00, 0x07, 0x00, 0x08, 0x80,
	}
	messageId = 150
	b := SET_READER_CONFIG(messageId, false,
		ReaderEventNotification(
			true, true, true,
			true, true, false,
			true, false, true,
		))

	if !bytes.Equal(v, b) {
		t.Errorf("want:% x", v)
		t.Errorf("resp:% x", b)
	}
}
func TestSetGPO(t *testing.T) {
	// no. 189
	v := []byte{
		// 0  1   	2		3	4		5	6	  7		8		9	10		11	12		13  14	  15
		0x04, 0x03, 0x00, 0x00, 0x00, 0x27, 0x00, 0x00, 0x00, 0x9b, 0x00, 0x00, 0xdb, 0x00, 0x07, 0x00,
		0x01, 0x00, 0x00, 0xdb, 0x00, 0x07, 0x00, 0x02, 0x00, 0x00, 0xdb, 0x00, 0x07, 0x00, 0x03, 0x00,
		0x00, 0xdb, 0x00, 0x07, 0x00, 0x04, 0x00,
	}
	messageId = 155
	b := SET_READER_CONFIG(messageId, false,
		gPOWriteData_Param(1, false),
		gPOWriteData_Param(2, false),
		gPOWriteData_Param(3, false),
		gPOWriteData_Param(4, false),
	)
	if !bytes.Equal(v, b) {
		t.Errorf("want:% x", v)
		t.Errorf("resp:% x", b)
	}
}
func TestSetRegion(t *testing.T) {
	v := []byte{
		0x04, 0x03, 0x00, 0x00, 0x00, 0x19, 0x00, 0x00, 0x00, 0x96, 0x80, 0x03, 0xff, 0x00, 0x0e, 0x00,
		0x00, 0x65, 0x1a, 0x00, 0x00, 0x00, 0x16, 0x00, 0x0e,
	}
	messageId = 150
	b := SET_READER_CONFIG(messageId, true,
		CustomParameter(
			uint32(25882),
			uint32(22),
			uint16(14),
		))
	if !bytes.Equal(v, b) {
		t.Errorf("want:% x", v)
		t.Errorf("resp:% x", b)
	}
}

func TestConvert2StatusAndGPICurrentState(t *testing.T) {
	// no. 968
	v := []byte{
		0x04, 0x0c, 0x00, 0x00, 0x00, 0x32, 0x00, 0x00, 0x01, 0x2e, 0x01, 0x1f, 0x00, 0x08, 0x00, 0x00,
		0x00, 0x00, 0x00, 0xe1, 0x00, 0x08, 0x00, 0x01, 0x80, 0x00, 0x00, 0xe1, 0x00, 0x08, 0x00, 0x02,
		0x00, 0x00, 0x00, 0xe1, 0x00, 0x08, 0x00, 0x03, 0x00, 0x00, 0x00, 0xe1, 0x00, 0x08, 0x00, 0x04,
		0x00, 0x00}
	res := Response(v, 50)
	config := res[0].(*GetConfigResponse)
	t.Logf("status=%v", config.Status.Success)
	for _, k := range config.GPI {
		t.Logf("gpi[%d]=%v", k.Number, k.State)
	}

}
func TestConvert2RoResByte(t *testing.T) {
	// no.1228
	v := []byte{
		0x04, 0x3d, 0x00, 0x00, 0x00, 0x31, 0x2e, 0x5a, 0x3d, 0x14, 0x00, 0xf0, 0x00, 0x27, 0x8d, 0xe2,
		0x00, 0x20, 0x75, 0x69, 0x09, 0x01, 0x46, 0x28, 0x00, 0x08, 0x07, 0x81, 0x00, 0x01, 0x86, 0xd9,
		0x87, 0x00, 0x01, 0x82, 0x00, 0x05, 0x67, 0x60, 0x72, 0xf1, 0x8b, 0x26, 0x90, 0x00, 0x00, 0x00,
		0x00,
		0x04, 0x3d, 0x00, 0x00, 0x00, 0x31, 0x2e, 0x5a, 0x3d, 0x14, 0x00, 0xf0, 0x00, 0x27, 0x8d, 0xe2,
		0x00, 0x20, 0x75, 0x69, 0x09, 0x01, 0x46, 0x28, 0x00, 0x08, 0x07, 0x81, 0x00, 0x01, 0x86, 0xd9,
		0x87, 0x00, 0x01, 0x82, 0x00, 0x05, 0x67, 0x60, 0x72, 0xf1, 0x8b, 0x26, 0x90, 0x00, 0x00, 0x00,
		0x00,
	}
	res := Response(v, 98)
	for _, k := range res {
		kk := k.(*ROAccessReportResponse)
		t.Logf("Datacard: %v %v", kk.MsgId, kk.Data.EPC_96)
	}
}

func TestAddROV2(t *testing.T) {
	// no. 100 (v2)
	v := []byte{
		0x04, 0x14, 0x00, 0x00, 0x01, 0x01, 0x00, 0x00, 0x00, 0xc9, 0x00, 0xb1, 0x00, 0xf7, 0x00, 0x00,
		0x04, 0xd2, 0x00, 0x00, 0x00, 0xb2, 0x00, 0x12, 0x00, 0xb3, 0x00, 0x05, 0x01, 0x00, 0xb6, 0x00,
		0x09, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xb7, 0x00, 0xce, 0x00, 0x02, 0x00, 0x01, 0x00, 0x02,
		0x00, 0xb8, 0x00, 0x09, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xba, 0x00, 0xbb, 0x04, 0xd2, 0x01,
		0x00, 0xde, 0x00, 0x5a, 0x00, 0x01, 0x00, 0xe0, 0x00, 0x0a, 0x00, 0x01, 0x00, 0x00, 0x00, 0x51,
		0x01, 0x4a, 0x00, 0x4a, 0x00, 0x01, 0x4f, 0x00, 0x08, 0x03, 0xe8, 0x00, 0x00, 0x01, 0x50, 0x00,
		0x0b, 0x80, 0x00, 0x20, 0x00, 0x00, 0x00, 0x00, 0x03, 0xff, 0x00, 0x0e, 0x00, 0x00, 0x65, 0x1a,
		0x00, 0x00, 0x00, 0x17, 0x00, 0x02, 0x03, 0xff, 0x00, 0x12, 0x00, 0x00, 0x65, 0x1a, 0x00, 0x00,
		0x00, 0x1a, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x03, 0xff, 0x00, 0x12, 0x00, 0x00, 0x65, 0x1a,
		0x00, 0x00, 0x00, 0x1c, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xde, 0x00, 0x5a, 0x00, 0x02,
		0x00, 0xe0, 0x00, 0x0a, 0x00, 0x01, 0x00, 0x00, 0x00, 0x51, 0x01, 0x4a, 0x00, 0x4a, 0x00, 0x01,
		0x4f, 0x00, 0x08, 0x03, 0xe8, 0x00, 0x00, 0x01, 0x50, 0x00, 0x0b, 0x80, 0x00, 0x20, 0x00, 0x00,
		0x00, 0x00, 0x03, 0xff, 0x00, 0x0e, 0x00, 0x00, 0x65, 0x1a, 0x00, 0x00, 0x00, 0x17, 0x00, 0x02,
		0x03, 0xff, 0x00, 0x12, 0x00, 0x00, 0x65, 0x1a, 0x00, 0x00, 0x00, 0x1a, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x03, 0xff, 0x00, 0x12, 0x00, 0x00, 0x65, 0x1a, 0x00, 0x00, 0x00, 0x1c, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0xed, 0x00, 0x0d, 0x02, 0x00, 0x01, 0x00, 0xee, 0x00, 0x06, 0x1e,
		0x40,
	}
	messageId = 201
	b := ADD_ROSPEC(messageId,
		RoSpec(1234, 0, 0,
			RoBoundSpec(1, 0, 0),
			AISpec(2,
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
			),
			RoReportSpec(2, 1,
				TagReportContentSelector(0x1e40),
			),
		),
	)
	//len_ := len(v)/2 + 12
	//if !bytes.Equal(v[len_:len_+10], b[len_:len_+10]) {
	if !bytes.Equal(v, b) {
		t.Errorf("want:% x", v)
		t.Errorf("resp:% x", b)
	}

}
func TestAddRO(t *testing.T) {
	// no .184
	v := []byte{
		0x04, 0x14, 0x00, 0x00, 0x01, 0x0d, 0x00, 0x00, 0x00, 0xc9, 0x00, 0xb1, 0x01, 0x03, 0x00, 0x00,
		0x04, 0xd2, 0x00, 0x00, 0x00, 0xb2, 0x00, 0x12, 0x00, 0xb3, 0x00, 0x05, 0x01, 0x00, 0xb6, 0x00,
		0x09, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xb7, 0x00, 0xda, 0x00, 0x02, 0x00, 0x01, 0x00, 0x02,
		0x00, 0xb8, 0x00, 0x09, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xba, 0x00, 0xc7, 0x04, 0xd2, 0x01,
		0x00, 0xde, 0x00, 0x60, 0x00, 0x01, 0x00, 0xdf, 0x00, 0x06, 0x00, 0x02, 0x00, 0xe0, 0x00, 0x0a,
		0x00, 0x01, 0x00, 0x00, 0x00, 0x01, 0x01, 0x4a, 0x00, 0x4a, 0x00, 0x01, 0x4f, 0x00, 0x08, 0x03,
		0xe8, 0x00, 0x00, 0x01, 0x50, 0x00, 0x0b, 0x80, 0x00, 0x20, 0x00, 0x00, 0x00, 0x00, 0x03, 0xff,
		0x00, 0x0e, 0x00, 0x00, 0x65, 0x1a, 0x00, 0x00, 0x00, 0x17, 0x00, 0x02, 0x03, 0xff, 0x00, 0x12,
		0x00, 0x00, 0x65, 0x1a, 0x00, 0x00, 0x00, 0x1a, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x03, 0xff,
		0x00, 0x12, 0x00, 0x00, 0x65, 0x1a, 0x00, 0x00, 0x00, 0x1c, 0x00, 0x01, 0x07, 0xd0, 0x00, 0xfa,
		0x00, 0xde, 0x00, 0x60, 0x00, 0x02, 0x00, 0xdf, 0x00, 0x06, 0x00, 0x02, 0x00, 0xe0, 0x00, 0x0a,
		0x00, 0x01, 0x00, 0x00, 0x00, 0x51, 0x01, 0x4a, 0x00, 0x4a, 0x00, 0x01, 0x4f, 0x00, 0x08, 0x03,
		0xe8, 0x00, 0x00, 0x01, 0x50, 0x00, 0x0b, 0x80, 0x00, 0x20, 0x00, 0x00, 0x00, 0x00, 0x03, 0xff,
		0x00, 0x0e, 0x00, 0x00, 0x65, 0x1a, 0x00, 0x00, 0x00, 0x17, 0x00, 0x02, 0x03, 0xff, 0x00, 0x12,
		0x00, 0x00, 0x65, 0x1a, 0x00, 0x00, 0x00, 0x1a, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x03, 0xff,
		0x00, 0x12, 0x00, 0x00, 0x65, 0x1a, 0x00, 0x00, 0x00, 0x1c, 0x00, 0x01, 0x07, 0xd0, 0x00, 0xfa,
		0x00, 0xed, 0x00, 0x0d, 0x02, 0x00, 0x01, 0x00, 0xee, 0x00, 0x06, 0x1e, 0x40,
	}
	messageId = 201
	b := ADD_ROSPEC(messageId,
		RoSpec(1234, 0, 0,
			RoBoundSpec(1, 0, 0),
			AISpec(2,
				AISpecStopTrigger(0, 0),
				InventoryParameterSpec(1234, 1,
					AntennaConfiguration(1,
						RFReceiver(2),
						RFTransmitter(1, 0, 1),
						C1G2InventoryCommand(0,
							C1G2RFControl(1000, 0),
							C1G2SingulationControl(0x80, 32, 0),
							CustomParameter(uint32(25882), uint32(23), uint16(2)),
							CustomParameter(uint32(25882), uint32(26), uint16(0), uint16(0), uint16(0)),
							CustomParameter(uint32(25882), uint32(28), uint16(1), uint16(2000), uint16(250)),
						),
					),
					AntennaConfiguration(2,
						RFReceiver(2),
						RFTransmitter(1, 0, 81),
						C1G2InventoryCommand(0,
							C1G2RFControl(1000, 0),
							C1G2SingulationControl(0x80, 32, 0),
							CustomParameter(uint32(25882), uint32(23), uint16(2)),
							CustomParameter(uint32(25882), uint32(26), uint16(0), uint16(0), uint16(0)),
							CustomParameter(uint32(25882), uint32(28), uint16(1), uint16(2000), uint16(250)),
						),
					),
				),
			),
			RoReportSpec(2, 1,
				TagReportContentSelector(0x1e40),
			),
		),
	)
	if !bytes.Equal(v, b) {
		t.Errorf("want:% x", v)
		t.Errorf("resp:% x", b)
	}
}
func TestDeleteRoSpec(t *testing.T) {
	// no. 158
	v := []byte{0x04, 0x15, 0x00, 0x00, 0x00, 0x0e, 0x00, 0x00, 0x00, 0x66, 0x00, 0x00, 0x00, 0x00}
	spec := 0 // all spec
	messageId = 102
	b := DEL_ROSPEC(messageId, spec)
	if !bytes.Equal(v, b) {
		t.Errorf("want:% x", v)
		t.Errorf("resp:% x", b)
	}

}
func TestDeleteAccessSpec(t *testing.T) {
	// no. 160
	v := []byte{0x04, 0x29, 0x00, 0x00, 0x00, 0x0e, 0x00, 0x00, 0x00, 0x67, 0x00, 0x00, 0x00, 0x00}
	spec := 0
	messageType := M_DELETE_ACCESSSPEC
	messageId = 103
	// we tried to used custom pack to reduce time work
	b := CustomPack(messageType, messageId,
		[]interface{}{
			uint32(spec),
		},
	)
	if !bytes.Equal(v, b) {
		t.Errorf("want:% x", v)
		t.Errorf("resp:% x", b)
	}

}

func TestSendFactoryReset(t *testing.T) {
	// no . 155
	v := []byte{0x04, 0x03, 0x00, 0x00, 0x00, 0x0b, 0x00, 0x00, 0x00, 0x65, 0x80}
	messageId = 101
	b := SET_READER_CONFIG(messageId, true)
	if !bytes.Equal(v, b) {
		t.Errorf("want:% x", v)
		t.Errorf("resp:% x", b)
	}
}

func TestSendGPI(t *testing.T) {
	messageId = 156
	v := []byte{0x04, 0x03,
		0, 0, 0, 0x2b,
		0, 0, 0, 0x9c,
		0,
		0, 0xe1, 0, 0x08, 0, 1, 0, 0,
		0, 0xe1, 0, 0x08, 0, 2, 0, 0,
		0, 0xe1, 0, 0x08, 0, 3, 0, 0,
		0, 0xe1, 0, 0x08, 0, 4, 0, 0,
	}
	b := SET_READER_CONFIG(messageId, false,
		gPIPortCurrentState_Param(1, 0, false),
		gPIPortCurrentState_Param(2, 0, false),
		gPIPortCurrentState_Param(3, 0, false),
		gPIPortCurrentState_Param(4, 0, false),
	)
	if !bytes.Equal(v, b) {
		t.Errorf("want:% x", v)
		t.Errorf("resp:% x", b)
	}

}
