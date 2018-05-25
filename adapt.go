package llrp

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

// set gpo via reader_id by order params 1-4
func (nc *Conn) GPOset(messageId int, reader_id string, params ...bool) error {
	var gpo [][]interface{}
	for i, k := range params {
		gpo = append(
			gpo,
			gPOWriteData_Param(i+1, k),
		)
	}
	if re, ok := nc.readers[reader_id]; ok {
		return re.conn.publish(
			SET_READER_CONFIG(
				messageId,
				false,
				gpo...,
			),
		)
	}
	return nil
}

func (nc *Conn) GPIget(messageId int, reader_id string) error {
	if re, ok := nc.readers[reader_id]; ok {
		return re.conn.publish(
			GET_READER_CONFIG_V1311(
				messageId,
				0,
				V_1311_GPIPortCurrentState,
				0,
				0,
			),
		)
	}
	return nil
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
	return nil

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
