package llrp

import (
	"fmt"
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
		GetDefaultAISpec(),
		GetRoReportSpec(),
	)
	lenn := len(b)
	fmt.Printf("\n-len-\n%d", lenn)
	fmt.Printf("\n-raw-\n% x", b)
	fmt.Printf("\n-ro start-\n% x", b[33:45])
	fmt.Printf("\n")
}

func TestO(t *testing.T) {
	//t.Run("dars", default_add_ro_spec)
	t.Run("ars", add_ro_spec)
}
