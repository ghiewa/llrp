package llrp

import (
	"bufio"
	log "github.com/sirupsen/logrus"
	"os"
	"testing"
)

var (
	card_evt bool
)

func handler(code int, msg *Msg) {
	if card_evt {
		log.Infof(msg.Data)
	}
	// msg.From - reader id

}

func TestM(t *testing.T) {
	t.Run("th", thread)
}
func thread() {
	opt := GetDefaultOptions()
	host := opt.NewConn()
	readers := []*SPReaderInfo{
		&SPReaderInfo{
			Id:   "random_reader_id",
			Host: "192.168.33.16:5084",
			InitCommand: [][]byte{
				ResetFactoryOpt(),
				DelROSpecOpt(),
				DelAccOption(),
				ExtensionOption(),
				SetRegion(),
				SetEventSpecOption(),
				AddROSpecOption(),
				EnableROSpecOption(),
			},
		},
	}
	for _, reader := range readers {
		// doReconnected when loss signal
		err := host.Registry(reader)
		if err != nil {
			t.Error(err)
		}

	}
	host.Subscription(handler)
	host.Start()
	var text string
	var success bool
	for {
		success = false
		log.Infof("Please enter command\nreader - list of readers\nnce - disable card event log\nce - enable card event log \nio - control gpo/get gpi state")
		scan := bufio.NewScanner(os.Stdin)
		scan.Scan()
		text = scan.Text()
		if text == "q" {
			break
		}
		switch text {
		case "reader":
			log.Infof("List Readers : %v", host.ListReader())
			success = true
		case "ne":
			// disable card logs
			card_evt = false
		case "ce":
			// enable card logs
			card_evt = true
		case "io":
			log.Infof("sample command please enter number(0-3)")
			scan.Scan()
			switch scan.Text() {
			case "0":
				// set gpo all open state // 0 = close , 1 = open , 2 = igonre
				// GPOset(id,port_state ...)  - set 4 port open state
				success = host.GPOset("random_reader_id", 1, 1, 1, 1)
			case "1":
				// set gpo spectfic port eg. port no.1 will open
				success = host.GPOsetp("random_reader_id", 1, true)
			case "2":
				// get all gpi port
				success = host.GPIget("random_reader_id")
			case "3":
				// get gpi port only port 1 will be reply
				success = ost.GPIget("random_reader_id", 1)
			default:
				log.Infof("not found cmd here")
			}
		default:
			log.Infof("not found cmd here")
		}
		if !success {
			log.Warnf("Send command not success.")
		}
	}
	// close connection
	host.Close()
}
