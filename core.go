package llrp

import (
	"bufio"
	"bytes"
	"net"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

func (nc *Conn) registry(sp *SPReaderInfo) error {
	if sp.Host == "" || sp.Id == "" || nc.readers[sp.Id] != nil {
		return ErrInvalidContext
	}
	// add to Conn
	nc.readers[sp.Id] = sp
	sp.conn = &RConn{
		opts:        &nc.Opts,
		initCommand: sp.InitCommand,
	}

	sp.conn.ach = make(chan asyncCB, asyncCBChanSize)
	if err := sp.conn.connect(sp.Host); err != nil {
		return err
	}

	go sp.conn.asyncDispatch()
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
func (c *RConn) createConn(host string) (err error) {
	c.lastAttempt = time.Now()
	c.conn, err = net.Dial("tcp", host)
	if err != nil {
		c.err = err
		return err
	}
	if c.pending != nil && c.bw != nil {
		// move to pending buffer.
		c.bw.Flush()
	}
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

	if nc.conn == nil || nc.bw == nil {
		return
	}

	for {
		nc.mu.Lock()
		defer nc.mu.Unlock()
		if _, ok := <-nc.fch; !ok {
			return
		}
		if !nc.isConnecting() || nc.isConnected() {
			return
		}
		if nc.bw.Buffered() > 0 {
			if nc.opts.FlusherTimeout > 0 {
				nc.conn.SetWriteDeadline(time.Now().Add(nc.opts.FlusherTimeout))
			}
			if err := nc.bw.Flush(); err != nil {
				if nc.err == nil {
					nc.err = err
				}
			}
			nc.conn.SetWriteDeadline(time.Time{})
		}
	}

}

// connect to reader
func (nc *RConn) connect(host string) error {
	nc.mu.Lock()
	nc.initc = true
	if err := nc.createConn(host); err != nil {
		return err
	}
	err := nc.processConnectInit()
	if err != nil {
		nc.mu.Unlock()
		nc.close(DISCONNECTED, false)
		return err
	}
	defer nc.mu.Unlock()
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
		nc.mu.Lock()
		defer nc.mu.Unlock()
		defer nc.kickFlusher()
		// check error condition
		sub := &Subscription{
			Id:   id,
			mcb:  cb,
			conn: nc,
		}
		sub.pMsgsLimit = DefaultSubPendingMsgsLimit
		sub.pBytesLimit = DefaultSubPendingBytesLimit
		sub.pCond = sync.NewCond(&sub.mu)
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
		nc.mu.Lock()
		conn := nc.conn
		nc.mu.Unlock()
		if conn == nil {
			break
		}
		n, err := conn.Read(b)
		if err != nil {
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
	nc.mu.Lock()
	if nc.isConnecting() || nc.isClosed() || nc.isReconnecting() {
		nc.mu.Unlock()
	}
	if nc.opts.AllowReconnect && nc.status == CONNECTED {
		nc.status = RECONNECTING
		if nc.ptmr != nil {
			nc.ptmr.Stop()
		}
		if nc.conn != nil {
			nc.bw.Flush()
			nc.conn.Close()
			nc.conn = nil
		}
		if nc.pending == nil {
			nc.pending = new(bytes.Buffer)
		}
		nc.pending.Reset()
		nc.bw.Reset(nc.pending)

		go nc.doReconnect()
		nc.mu.Unlock()
		return
	}
	nc.status = DISCONNECTED
	nc.err = err
	nc.mu.Unlock()
	nc.Close()
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
		return
	}
	nc.Reconnects++
	if nc.err = nc.processConnectInit(); nc.err != nil {
		nc.status = RECONNECTING
		nc.mu.Unlock()
		return
	}
	nc.didConnect = true
	nc.reconnects = 0

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

	// process init commands ( reset factory / set gpo off and so on..

	err = nc.bw.Flush()
	if err != nil {
		return err
	}
	err = nc.sendPrefixCommand()
	if err != nil {
		return err
	}

	go nc.spinUpGoRoutines()
	return nil
}
func (nc *RConn) sendPrefixCommand() error {
	nc.mu.Lock()
	defer nc.mu.Unlock()
	for _, k := range nc.initCommand {
		_, err := nc.bw.Write(k)
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

	if nc.ptmr != nil {
		nc.ptmr.Stop()
	}

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
