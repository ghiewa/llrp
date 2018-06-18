package llrp

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"time"
)

func (o Options) NewConn() *Conn {
	if o.ReconnectBufSize == 0 {
		o.ReconnectBufSize = DefaultReconnectBufSize
	}
	if o.Timeout == 0 {
		o.Timeout = DefaultTimeout
	}
	return &Conn{
		Opts:    o,
		readers: make(map[string]*SPReaderInfo),
	}
}
func (nc *Conn) Registry(reader *SPReaderInfo) error {
	return nc.registry(reader)
}
func (nc *Conn) GPICheck(reader_id string, gpis []*GPICurrentState) error {
	if re, ok := nc.readers[reader_id]; ok {
		state := re.state.gpi
		for _, v := range gpis {
			if pre_state, ok := state.Ports[v.Number]; ok {
				if !state.init || v.State != pre_state {
					pa := new(IOState)
					pa.ReaderId = reader_id
					pa.Ports = map[uint16]int{
						v.Number: v.State,
					}
					state.init = true
					state.Ports[v.Number] = v.State
					// fire to toggle handler
					if state.cb != nil {
						state.cb(pa)
					}
				}
			}
		}
		return nil
	}
	return fmt.Errorf("Cann't find reader id")
}

func (nc *Conn) GPIToggleMonitor(reader_id string, port_trigger int, t time.Duration, cb HandlerIOState) error {
	if re, ok := nc.readers[reader_id]; ok {
		// filter toggle gpi
		var (
			PortTrigger = uint16(port_trigger)
		)
		re.conn.mu.Lock()
		if re.state == nil {
			re.state = new(GState)
		}
		if re.state.gpi == nil {
			re.state.gpi = new(IOState)
		}
		re.state.gpi.Ports = map[uint16]int{
			PortTrigger: 0,
		}
		re.state.gpi.cb = cb
		re.conn.mu.Unlock()
		go func() {
			for {
				select {
				case <-time.After(t):
					log.Debugf("time to check gpi")
					if re.conn.isClosed() {
						return
					}
					err := re.conn.publish(
						GET_READER_CONFIG(
							rand.Int(),
							0,
							C_GET_READER_CONFIG_GPIPortCurrentState,
							0,
							0,
						),
					)
					if err != nil {
						log.Errorf("gpi error %s", err)
					}

				}
			}
		}()
		return nil
	} else {
		return fmt.Errorf("Cann't find reader id")
	}
}

func (nc *Conn) Lock() {
	nc.mu.Lock()
}
func (nc *Conn) Unlock() {
	nc.mu.Unlock()
}

// List of readers registed
func (nc *Conn) ListReader() map[string]*SPReaderInfo {
	return nc.readers
}
func (nc *Conn) StopROSpec(messageId, ROSpecID int, reader_id string) error {
	if re, ok := nc.readers[reader_id]; ok {
		return re.conn.publish(
			STOP_ROSPEC(
				messageId,
				ROSpecID,
			),
		)
	}
	return fmt.Errorf("Cann't find reader id")

}

func (nc *Conn) StartROSpec(messageId, ROSpecID int, reader_id string) error {
	if re, ok := nc.readers[reader_id]; ok {
		return re.conn.publish(
			START_ROSPEC(
				messageId,
				ROSpecID,
			),
		)
	}
	return fmt.Errorf("Cann't find reader id")

}

func (nc *Conn) Enable_ROSpec(messageId, ROSpecID int, reader_id string) error {

	if re, ok := nc.readers[reader_id]; ok {
		return re.conn.publish(
			ENABLE_ROSPEC(
				messageId,
				ROSpecID,
			),
		)
	}
	return fmt.Errorf("Cann't find reader id")
}

func (nc *Conn) Disabled_ROSpec(messageId, ROSpecID int, reader_id string) error {
	if re, ok := nc.readers[reader_id]; ok {
		return re.conn.publish(
			DISABLE_ROSPEC(
				messageId,
				ROSpecID,
			),
		)
	}
	return fmt.Errorf("Cann't find reader id")
}

// set gpo via reader_id by order params 1-4
func (nc *Conn) GPOset(messageId int, reader_id string, params ...bool) error {
	if re, ok := nc.readers[reader_id]; ok {
		var gpo [][]interface{}
		for i, k := range params {
			gpo = append(
				gpo,
				gPOWriteData_Param(i+1, k),
			)
		}
		return re.conn.publish(
			SET_READER_CONFIG(
				messageId,
				false,
				gpo...,
			),
		)
	}
	return fmt.Errorf("Cann't find reader id")
}

func (nc *Conn) GetRoReport(messageId int, reader_id string) error {
	if re, ok := nc.readers[reader_id]; ok {
		return re.conn.publish(
			GET_REPORT(messageId),
		)
	}
	return fmt.Errorf("Cann't find reader id")
}

func (nc *Conn) GPIset(messageId int, reader_id string, port int, port_state bool) error {
	if re, ok := nc.readers[reader_id]; ok {
		port_s := 0
		if port_state {
			port_s += 1
		}
		gpi := gPIPortCurrentState_Param(
			port,
			port_s,
			false,
		)
		return re.conn.publish(
			SET_READER_CONFIG(
				messageId,
				false,
				gpi,
			),
		)
	}
	return fmt.Errorf("Cann't find reader id")
}
func (nc *Conn) GPIget(messageId int, reader_id string) error {
	if re, ok := nc.readers[reader_id]; ok {
		return re.conn.publish(
			GET_READER_CONFIG(
				messageId,
				0,
				C_GET_READER_CONFIG_GPIPortCurrentState,
				0,
				0,
			),
		)
	}
	return fmt.Errorf("Cann't find reader id")
}
func (sc *Subscription) Ack(messageId int) error {
	return sc.conn.publish(
		SEND_KEEPALIVE(messageId),
	)
}

func (nc *Conn) GPOsetp(messageId int, reader_id string, number_port int, state bool) error {
	if re, ok := nc.readers[reader_id]; ok {
		return re.conn.publish(
			SET_READER_CONFIG(
				messageId,
				false,
				gPOWriteData_Param(number_port, state),
			),
		)
	}
	return fmt.Errorf("Cann't find reader id")

}

func (nc *Conn) Subscription(cb MsgHandler) ([]*Subscription, error) {
	return nc.subscribe(cb, nil)
}

// Close will close the connection to the server. This call will release
// all blocking calls, such as Flush() and NextMsg()
func (nc *Conn) Close() {
	for _, k := range nc.readers {
		k.conn.close(CLOSED, true)
	}
}

// IsClosed tests if a Conn has been closed.
func (nc *RConn) IsClosed() bool {
	nc.mu.Lock()
	defer nc.mu.Unlock()
	return nc.isClosed()
}

// IsReconnecting tests if a Conn is reconnecting.
func (nc *RConn) IsReconnecting() bool {
	nc.mu.Lock()
	defer nc.mu.Unlock()
	return nc.isReconnecting()
}

// IsConnected tests if a Conn is connected.
func (nc *RConn) IsConnected() bool {
	nc.mu.Lock()
	defer nc.mu.Unlock()
	return nc.isConnected()
}
func GetDefaultOptions() Options {
	return Options{
		AllowReconnect:   true,
		MaxReconnect:     DefaultMaxReconnect,
		ReconnectWait:    DefaultReconnectWait,
		Timeout:          DefaultTimeout,
		ReconnectBufSize: DefaultReconnectBufSize,
	}
}

// Status returns the current state of the connection.
func (nc *RConn) Status() Status {
	nc.mu.Lock()
	defer nc.mu.Unlock()
	return nc.status
}

// Test if Conn has been closed Lock is assumed held.
func (nc *RConn) isClosed() bool {
	return nc.status == CLOSED
}

// Test if Conn is in the process of connecting
func (nc *RConn) isConnecting() bool {
	return nc.status == CONNECTING
}

// Test if Conn is being reconnected.
func (nc *RConn) isReconnecting() bool {
	return nc.status == RECONNECTING
}

// Test if Conn is connected or connecting.
func (nc *RConn) isConnected() bool {
	return nc.status == CONNECTED
}

// Stats will return a race safe copy of the Statistics section for the connection.
func (nc *RConn) Stats() Statistics {
	// Stats are updated either under connection's mu or subsMu mutexes.
	// Lock both to safely get them.
	nc.mu.Lock()
	nc.subsMu.RLock()
	stats := Statistics{
		InMsgs:     nc.InMsgs,
		InBytes:    nc.InBytes,
		OutMsgs:    nc.OutMsgs,
		OutBytes:   nc.OutBytes,
		Reconnects: nc.Reconnects,
	}
	nc.subsMu.RUnlock()
	nc.mu.Unlock()
	return stats
}
