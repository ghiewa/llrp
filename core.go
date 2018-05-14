package llrp

import (
	"bufio"
	"bytes"
	log "github.com/sirupsen/logrus"
	"net"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// register reader to main conn
func (nc *Conn) registry(sp *SPReaderInfo) error {
	if sp.Host == "" || sp.Id == "" || nc.readers[sp.Id] != nil {
		return ErrInvalidContext
	}
	// add to Conn
	nc.readers[sp.Id] = sp
	sp.conn = &RConn{
		opts:        &nc.Opts,
		initCommand: sp.InitCommand,
		host:        sp.Host,
		ip:          sp.Host[:strings.Index(sp.Host, ":")],
		fch:         make(chan struct{}, flushChanSize),
	}

	sp.conn.ach = make(chan asyncCB, asyncCBChanSize)
	go func() {
		if err := sp.conn.connect(); err != nil {
			log.Errorf("Unable to connected :%s , we will reconnect every %v[%d]", sp.Host, sp.conn.opts.ReconnectWait, sp.conn.reconnects)
			sp.conn.processOpErr(err)
		}
	}()
	go sp.conn.asyncDispatch()
	return nil
}

// logic of pushing msg to reader
func (nc *RConn) publish(data []byte) error {
	if nc == nil {
		return ErrInvalidConnection
	}
	if nc.isClosed() || nc.bw == nil {
		log.Errorf("can't publish channal close", nc.host)
		return ErrConnectionClosed
	}
	log.Infof("publishing %v", nc.status)
	nc.mu.Lock()
	defer nc.mu.Unlock()
	if nc.isReconnecting() {
		log.Infof("publish reconnect %v", nc.status)
		nc.bw.Flush()
		log.Infof("publish reconnect flush %v", nc.status)
		if nc.pending.Len() >= nc.opts.ReconnectBufSize {
			return ErrReconnectBufExceeded
		}
	}
	log.Infof("publish %v", nc.status)
	l, err := nc.bw.Write(data)
	if err != nil {
		return err
	}
	nc.OutMsgs++
	nc.OutBytes += uint64(l)
	if len(nc.fch) == 0 {
		nc.kickFlusher()
	}
	return nil
}

// asyncDispatch is responsible for calling any async callbacks
func (nc *RConn) asyncDispatch() {
	nc.mu.Lock()
	ach := nc.ach
	nc.mu.Unlock()
	for {
		if f, ok := <-ach; !ok {
			return
		} else {
			f()
		}
	}
}

// createConn will connect to the reader and wrap the appropriate
// bufio structures. It will do the right thing when an existing
// connection is in place.
func (c *RConn) createConn() (err error) {
	c.lastAttempt = time.Now()
	log.Infof("try to dial %s | %v", c.host, c.opts.Timeout)
	c.conn, err = net.DialTimeout("tcp", c.host, c.opts.Timeout)
	if err != nil {
		log.Errorf("failed to dial on %s", c.host)
		c.err = err
		return err
	}

	log.Infof("dial to %s", c.host)
	if c.pending != nil && c.bw != nil {
		// move to pending buffer.
		log.Infof("fluse bw %s", c.host)
		c.bw.Flush()
	}
	log.Infof("create buffer %s", c.host)
	c.bw = bufio.NewWriterSize(c.conn, defaultBufSize)
	return nil
}

func (nc *RConn) waitForExits(wg *sync.WaitGroup) {
	// flusher off
	select {
	case nc.fch <- struct{}{}:
	default:
	}
	// wait for any previous go routines
	if wg != nil {
		wg.Wait()
	}
}

// spinUpGoRoutines will launch the Go routines responsible for
// reading and writing to the socket. This will be launched via a
// go routine itself to release any locks that may be held.
// We also use a WaitGroup to make sure we only start them on a
// reconnect when the previous ones have exited.
func (nc *RConn) spinUpGoRoutines() {
	nc.waitForExits(nc.wg)
	nc.wg = &sync.WaitGroup{}
	nc.wg.Add(2)
	// spin
	go nc.readLoop(nc.wg)
	go nc.flusher(nc.wg)
}
func (nc *RConn) flusher(wg *sync.WaitGroup) {
	defer wg.Done()
	nc.mu.Lock()
	bw := nc.bw
	conn := nc.conn
	fch := nc.fch
	flusherTimeout := nc.opts.FlusherTimeout
	nc.mu.Unlock()
	if conn == nil || bw == nil {
		return
	}
	for {
		if _, ok := <-fch; !ok {
			return
		}
		nc.mu.Lock()
		if !nc.isConnecting() || nc.isConnected() || bw != nc.bw || conn != nc.conn {
			nc.mu.Unlock()
			return
		}
		if bw.Buffered() > 0 {
			if flusherTimeout > 0 {
				conn.SetWriteDeadline(time.Now().Add(nc.opts.FlusherTimeout))
			}
			if err := bw.Flush(); err != nil {
				if nc.err == nil {
					nc.err = err
				}
			}
			conn.SetWriteDeadline(time.Time{})
		}
		nc.mu.Unlock()
	}

}

// connect to reader
func (nc *RConn) connect() error {
	// create conn
	if err := nc.createConn(); err != nil {
		log.Errorf("Can't connecting")
		return err
	}
	log.Infof("Setup initCommand")
	err := nc.processConnectInit()
	if err != nil {
		log.Errorf("processConnectInit not success : %v", err)
		nc.close(DISCONNECTED, false)
		return err
	}

	nc.didConnect = true
	nc.reconnects = 0
	log.Infof("Establish connection %s", nc.host)
	return nil
}

func (cnc *Conn) subscribe(cb MsgHandler, ch chan *Msg) ([]*Subscription, error) {
	if cb == nil {
		return nil, ErrBadSubscription
	}
	var (
		subs []*Subscription
	)
	for id, ncc := range cnc.readers {
		nc := ncc.conn
		defer nc.kickFlusher()
		// check error condition
		sub := &Subscription{
			Id:   id,
			Ip:   nc.ip,
			mcb:  cb,
			conn: nc,
		}
		sub.pMsgsLimit = DefaultSubPendingMsgsLimit
		sub.pBytesLimit = DefaultSubPendingBytesLimit
		sub.pCond = sync.NewCond(&sub.mu)
		nc.sub = sub
		go nc.waitForMsgs(sub)

		subs = append(subs, sub)
	}
	return subs, nil

}
func (nc *RConn) waitForMsgs(s *Subscription) {
	var (
		closed         bool
		delivered, max uint64
	)
	for {
		s.mu.Lock()
		if s.pHead == nil && !s.closed {
			s.pCond.Wait()
		}
		// pop msg from list
		m := s.pHead
		if m != nil {
			s.pHead = m.next
			if s.pHead == nil {
				s.pTail = nil
			}
			if m.barrier != nil {
				s.mu.Unlock()
				if atomic.AddInt64(&m.barrier.refs, -1) == 0 {
					m.barrier.f()
				}
				continue
			}
			s.pMsgs--
			s.pBytes -= m.len_data
		}
		mcb := s.mcb
		max = s.max
		closed = s.closed
		if !s.closed {
			s.delivered++
			delivered = s.delivered
		}
		s.mu.Unlock()
		if closed {
			break
		}
		// deliver msg
		if m != nil && (max == 0 || delivered <= max) {
			mcb(m)
		}
		if max > 0 && delivered >= max {
			nc.mu.Lock()
			nc.removeSub(s)
			nc.mu.Unlock()
			break
		}
	}
	// check barrier msg
	s.mu.Lock()
	for m := s.pHead; m != nil; m = s.pHead {
		if m.barrier != nil {
			s.mu.Unlock()
			if atomic.AddInt64(&m.barrier.refs, -1) == 0 {
				m.barrier.f()
			}
			s.mu.Lock()
		}
		s.pHead = m.next
	}
	s.mu.Unlock()
}
func (nc *RConn) removeSub(s *Subscription) {
	s.mu.Lock()
	defer s.mu.Unlock()
	// mark as invalid
	s.conn = nil
	s.closed = true
	if s.pCond != nil {
		s.pCond.Broadcast()
	}
}

// readLoop() will sit on the socket reading and processing the
// protocol from the server. It will dispatch appropriately based
// on the op type.
func (nc *RConn) readLoop(wg *sync.WaitGroup) {
	defer wg.Done()
	// Stack based buffer.

	b := make([]byte, defaultBufSize)
	for {
		log.Infof("ReadLoop %+v", nc.status)
		nc.mu.Lock()
		conn := nc.conn
		conn.SetReadDeadline(time.Now().Add(nc.opts.Timeout))
		nc.mu.Unlock()
		log.Infof("ReadLoop  Lock ")
		if conn == nil {
			log.Errorf("conn is nil")
			break
		}
		n, err := conn.Read(b)
		if err != nil {
			log.Errorf("readLoop op error %d", n)
			nc.processOpErr(err)
			break
		}
		// process
		if err := nc.process(b[:n], n); err != nil {
			nc.processOpErr(err)
			break
		}
	}
}
func (nc *RConn) processOpErr(err error) {
	if nc.opts.AllowReconnect {
		nc.status = RECONNECTING
		nc.didConnect = false
		if nc.conn != nil {
			nc.bw.Flush()
			nc.conn.Close()
			nc.conn = nil
		}
		if nc.pending == nil {
			nc.pending = new(bytes.Buffer)
		}
		nc.pending.Reset()
		if nc.bw != nil {
			nc.bw.Reset(nc.pending)
		}
		nc.mu.Unlock()
		go func() {
			for {
				nc.doReconnect()
				n := &NetworkIssue{
					Period:     nc.opts.ReconnectWait,
					Reconnects: nc.reconnects,
					Type:       NETW_LOSS,
				}
				if nc.didConnect {
					nc.status = CONNECTED
					n.Type = NETW_CONNECTED
					nc.sendReport(0, n)
					break
				}
				nc.sendReport(0, n)
			}
		}()
		return
	}

	nc.status = DISCONNECTED
	nc.err = err
	nc.mu.Unlock()
	nc.close(CLOSED, true)
}
func (nc *RConn) doReconnect() {
	nc.mu.Lock()
	wg := nc.wg
	nc.mu.Unlock()

	nc.waitForExits(wg)

	nc.mu.Lock()
	nc.err = nil

	if nc.opts.DisconnectedCB != nil {
		nc.ach <- func() { nc.opts.DisconnectedCB(nc) }
	}

	sleepTime := int64(0)
	// Sleep appropriate amount of time before the
	// connection attempt if connecting to same server
	// we just got disconnected from..
	if time.Since(nc.lastAttempt) < nc.opts.ReconnectWait {
		sleepTime = int64(nc.opts.ReconnectWait - time.Since(nc.lastAttempt))
	}

	nc.mu.Unlock()
	if sleepTime <= 0 {
		runtime.Gosched()
	} else {
		time.Sleep(time.Duration(sleepTime))
	}
	nc.mu.Lock()
	if nc.isClosed() {
		nc.mu.Unlock()
		return
	}
	nc.Reconnects++
	log.Infof("tried to createConn")
	if err := nc.createConn(); err != nil {
		log.Errorf("createConn error")
		nc.reconnects++
		nc.mu.Unlock()
		return
	}
	log.Infof("tried to processConnectInit")
	if nc.err = nc.processConnectInit(); nc.err != nil {
		log.Errorf("processConnectInit error")
		nc.status = RECONNECTING
		nc.mu.Unlock()
		return
	}

	log.Infof("didConnect %s %v", nc.host, nc.status)
	nc.didConnect = true
	nc.reconnects = 0
	nc.mu.Unlock()

}

// flushReconnectPending will push the pending items that were
// gathered while we were in a RECONNECTING state to the socket.
func (nc *RConn) flushReconnectPendingItems() {
	if nc.pending == nil {
		return
	}
	if nc.pending.Len() > 0 {
		nc.bw.Write(nc.pending.Bytes())

	}
}

// Process a connected connection and initialize properly.
func (nc *RConn) processConnectInit() (err error) {
	nc.conn.SetDeadline(time.Now().Add(nc.opts.Timeout))
	defer nc.conn.SetDeadline(time.Time{})

	nc.status = CONNECTING

	log.Infof("1", nc.mu)
	// process init commands ( reset factory / set gpo off and so on..
	err = nc.bw.Flush()
	log.Infof("2")
	if err != nil {
		log.Infof("2e")
		return err
	}

	log.Infof("3", nc.mu)
	err = nc.sendPrefixCommand()
	if err != nil {
		log.Infof("3e")
		return err
	}
	log.Infof("4")
	nc.kickFlusher()
	log.Infof("5")
	go nc.spinUpGoRoutines()
	return nil
}
func (nc *RConn) sendPrefixCommand() error {
	for _, k := range nc.initCommand {
		_, err := nc.bw.Write(k)
		log.Info("7")
		if err != nil {
			return err
		}
	}
	return nil
}

// Low level close call that will do correct cleanup and set
// desired status. Also controls whether user defined callbacks
// will be triggered. The lock should not be held entering this
// function. This function will handle the locking manually.
func (nc *RConn) close(status Status, doCBs bool) {
	nc.mu.Lock()
	if nc.isClosed() {
		nc.status = status
		nc.mu.Unlock()
		return
	}
	nc.status = CLOSED
	// Kick the Go routines so they fall out.
	nc.kickFlusher()
	nc.mu.Unlock()

	nc.mu.Lock()

	// Go ahead and make sure we have flushed the outbound
	if nc.conn != nil {
		nc.bw.Flush()
		defer nc.conn.Close()
	}

	// Perform appropriate callback if needed for a disconnect.
	if doCBs {
		if nc.opts.DisconnectedCB != nil && nc.conn != nil {
			nc.ach <- func() { nc.opts.DisconnectedCB(nc) }
		}
		if nc.opts.ClosedCB != nil {
			nc.ach <- func() { nc.opts.ClosedCB(nc) }
		}
		nc.ach <- nc.closeAsyncFunc()
	}
	nc.status = status
	nc.mu.Unlock()
}

func (nc *RConn) closeAsyncFunc() asyncCB {
	return func() {
		nc.mu.Lock()
		if nc.ach != nil {
			close(nc.ach)
			nc.ach = nil
		}
		nc.mu.Unlock()
	}
}

// kickFlusher will send a bool on a channel to kick the
// flush Go routine to flush data to the reader.
func (nc *RConn) kickFlusher() {
	if nc.bw != nil {
		select {
		case nc.fch <- struct{}{}:
		default:
		}
	}
}
