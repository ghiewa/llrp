package llrp

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"testing"
)

func default_add_ro_spec(t *testing.T) {
	//fmt.Printf("% x", AddROSpecOption())
}
func add_ro_spec(t *testing.T) {
	var (
		timeout      = uint32(10000) // milliseconds
		port_trigger = uint16(1)
	)
	log.SetLevel(log.DebugLevel)
	b := AddROSpecCustom(
		// set trigger option - gpi
		RoBoundSpecCustom(
			//GPITriggerValue option = 3
			ROSpecStartTrigger(3,
				GPITriggerValue(port_trigger, true, timeout),
			),
			// stop by duration trigger
			ROSpecStopTrigger(1,
				timeout,
			),
		),
		//GetDefaultAISpec(),
		//GetRoReportSpec(),
	)
	fmt.Printf("\n-raw-\n% x", b)
	fmt.Printf("\n-ro start-\n% x", b[34:45])
	fmt.Printf("\n-pack-\n% x", pack(GPITriggerValue(port_trigger, true, timeout)))
	fmt.Printf("\n")

}
func check_keepalive(t *testing.T) {

}
func check_set_event(t *testing.T) {
	br := []byte{
		0x04, 0x03, 0x00, 0x00, 0x00, 0x4e, 0x00, 0x00, 0x00, 0x64, 0x00, 0x00, 0xf4, 0x00, 0x43, 0x00,
		0xf5, 0x00, 0x07, 0x00, 0x00, 0x80, 0x00, 0xf5, 0x00, 0x07, 0x00, 0x01, 0x80, 0x00, 0xf5, 0x00,
		0x07, 0x00, 0x02, 0x80, 0x00, 0xf5, 0x00, 0x07, 0x00, 0x03, 0x80, 0x00, 0xf5, 0x00, 0x07, 0x00,
		0x04, 0x80, 0x00, 0xf5, 0x00, 0x07, 0x00, 0x05, 0x00, 0x00, 0xf5, 0x00, 0x07, 0x00, 0x06, 0x80,
		0x00, 0xf5, 0x00, 0x07, 0x00, 0x07, 0x00, 0x00, 0xf5, 0x00, 0x07, 0x00, 0x08, 0x80,
	}
	b := SET_READER_CONFIG(
		100, // message id
		false,
		ReaderEventNotificationSpec(),
	)
	log.Infof("check_set_event\n % x", b)
	if len(b) == len(br) {
		log.Infof("check_set_event ok\n % x", br)
	} else {
		log.Errorf("check_set_event notok\n % x", br)
	}
}
func TestO(t *testing.T) {
	//t.Run("dars", default_add_ro_spec)
	t.Run("keep_alive", check_keepalive)
	t.Run("event", check_set_event)
	t.Run("add_ro", add_ro_spec)
}
