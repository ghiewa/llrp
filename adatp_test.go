package llrp

import (
	"bufio"
	log "github.com/sirupsen/logrus"
	"os"
	"testing"
)

var (
	card_evt bool = true
)

func handler(msg *Msg) {
	if card_evt {
		log.Infof("evt %v", msg.Reports)
	}
	// msg.From - reader id

}

func TestM(t *testing.T) {
	t.Run("th", loop)
}
func loop(t *testing.T) {
	log.Info("loop")
	opt := GetDefaultOptions()
	host := opt.NewConn()
	log.SetLevel(log.DebugLevel)
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
	log.Info("registry")
	for _, reader := range readers {
		// doReconnected when loss signal
		err := host.Registry(reader)
		if err != nil {
			log.Error(err)
		}

	}

	log.Info("subscribe")
	var text string
	var err error
	for {
		err = nil
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
		case "ne":
			// disable card logs
			card_evt = false
		case "ce":
			// enable card logs
			card_evt = true
		case "io":
			log.Infof("sample command please enter number(0-2)")
			scan.Scan()
			switch scan.Text() {
			case "0":
				// set gpo all open state // 0 = close , 1 = open , 2 = igonre
				// GPOset(id,port_state ...)  - set 4 port open state
				err = host.GPOset(123, "random_reader_id", true, true, true, true)
			case "1":
				// set gpo spectfic port eg. port no.1 will open
				err = host.GPOsetp(222, "random_reader_id", 1, true)
			case "2":
				// get all gpi port
				err = host.GPIget(333, "random_reader_id")
			default:
				log.Infof("not found cmd here")
			}
		default:
			log.Infof("not found cmd here")
		}
		if err != nil {
			log.Warnf("Send command not success. %v", err)
		} else {
			log.Infof("Send command success. ")
		}
	}
	// close connection
	host.Close()
}
