package llrp

import (
	"bufio"
	"bytes"
	"errors"
	"net"
	"sync"
	"time"
)

// Errors
var (
	ErrConnectionClosed     = errors.New("llrp: connection closed")
	ErrSecureConnRequired   = errors.New("llrp: secure connection required")
	ErrSecureConnWanted     = errors.New("llrp: secure connection not available")
	ErrBadSubscription      = errors.New("llrp: invalid subscription")
	ErrTypeSubscription     = errors.New("llrp: invalid subscription type")
	ErrBadSubject           = errors.New("llrp: invalid subject")
	ErrSlowConsumer         = errors.New("llrp: slow consumer, messages dropped")
	ErrTimeout              = errors.New("llrp: timeout")
	ErrBadTimeout           = errors.New("llrp: timeout invalid")
	ErrNoServers            = errors.New("llrp: no readers available for connection")
	ErrJsonParse            = errors.New("llrp: connect message, json parse error")
	ErrChanArg              = errors.New("llrp: argument needs to be a channel type")
	ErrMaxPayload           = errors.New("llrp: maximum payload exceeded")
	ErrMaxMessages          = errors.New("llrp: maximum messages delivered")
	ErrSyncSubRequired      = errors.New("llrp: illegal call on an async subscription")
	ErrMultipleTLSConfigs   = errors.New("llrp: multiple tls.Configs not allowed")
	ErrNoInfoReceived       = errors.New("llrp: protocol exception, INFO not received")
	ErrReconnectBufExceeded = errors.New("llrp: outbound buffer limit exceeded")
	ErrInvalidConnection    = errors.New("llrp: invalid connection")
	ErrInvalidMsg           = errors.New("llrp: invalid message or message nil")
	ErrInvalidArg           = errors.New("llrp: invalid argument")
	ErrInvalidContext       = errors.New("llrp: invalid context")
)

// Pending Limits
const (
	DefaultSubPendingMsgsLimit  = 65536
	DefaultSubPendingBytesLimit = 65536 * 1024
)
const (
	BufferSize           = 512
	port_default         = "5084"
	DefaultMaxReconnect  = 60
	DefaultReconnectWait = 2 * time.Second
	DefaultTimeout       = 2 * time.Second
	DefaultPingInterval  = 2 * time.Minute
)

const (
	_CRLF_  = "\r\n"
	_EMPTY_ = ""
	_SPC_   = " "
)

type NetworkIssue struct {
	From *Subscription
}

// Status represents the state of the connection.
type Status int

const (
	DISCONNECTED = Status(iota)
	CONNECTED
	CLOSED
	RECONNECTING
	CONNECTING
)

// ConnHandler is used for asynchronous events such as
// disconnected and closed connections.
type ConnHandler func(*RConn)

// ErrHandler is used to process asynchronous errors encountered
// while processing inbound messages.
type ErrHandler func(*Conn, *Subscription, error)

// asyncCB is used to preserve order for async callbacks.
type asyncCB func()

// Option is a function on the options for a connection.
type Option func(*Options) error

// CustomDialer can be used to specify any dialer, not necessarily
// a *net.Dialer.
type CustomDialer interface {
	Dial(network, address string) (net.Conn, error)
}

// Default Constants
const (
	DefaultPort             = 5084
	DefaultMaxPingOut       = 2
	DefaultMaxChanLen       = 8192            // 8k
	DefaultReconnectBufSize = 8 * 1024 * 1024 // 8MB
	RequestChanLen          = 8
)

// Options can be used to create a customized connection.
type Options struct {
	// command stack when start reader
	InitCommand [][]byte
	// AllowReconnect enables reconnection logic to be used when we
	// encounter a disconnect from the current reader.
	AllowReconnect bool

	// MaxReconnect sets the number of reconnect attempts that will be
	// tried before giving up. If negative, then it will never give up
	// trying to reconnect.
	MaxReconnect int

	// ReconnectWait sets the time to backoff after attempting a reconnect
	// to a reader that we were already connected to previously.
	ReconnectWait time.Duration

	// Timeout sets the timeout for a Dial operation on a connection.
	Timeout time.Duration

	// FlusherTimeout is the maximum time to wait for the flusher loop
	// to be able to finish writing to the underlying connection.
	FlusherTimeout time.Duration

	// ClosedCB sets the closed handler that is called when a client will
	// no longer be connected.
	ClosedCB ConnHandler

	// DisconnectedCB sets the disconnected handler that is called
	// whenever the connection is disconnected.
	DisconnectedCB ConnHandler

	// ReconnectedCB sets the reconnected handler called whenever
	// the connection is successfully reconnected.
	ReconnectedCB ConnHandler

	// DiscoveredServersCB sets the callback that is invoked whenever a new
	// reader has joined the cluster.
	DiscoveredServersCB ConnHandler

	// AsyncErrorCB sets the async error handler (e.g. slow consumer errors)
	AsyncErrorCB ErrHandler

	// ReconnectBufSize is the size of the backing bufio during reconnect.
	// Once this has been exhausted publish operations will return an error.
	ReconnectBufSize int
}

const (
	// Scratch storage for assembling protocol headers
	scratchSize = 512

	// The size of the bufio reader/writer on top of the socket.
	defaultBufSize = 32768

	// The buffered size of the flush "kick" channel
	flushChanSize = 1024

	// Default reader pool size
	srvPoolSize = 4

	// Channel size for the async callback handler.
	asyncCBChanSize = 32

	// NUID size
	nuidSize = 22
)

type RConn struct {
	Statistics
	ip          string
	opts        *Options
	host        string
	didConnect  bool
	reconnects  int
	lastAttempt time.Time
	isImplicit  bool
	mu          sync.Mutex
	wg          *sync.WaitGroup
	conn        net.Conn
	bw          *bufio.Writer
	pending     *bytes.Buffer
	fch         chan struct{}
	subsMu      sync.RWMutex
	ach         chan asyncCB
	scratch     [scratchSize]byte
	status      Status
	initc       bool // true if the connection is performing the initial connect
	err         error
	ptmr        *time.Timer
	pout        int
	sub         *Subscription
	initCommand [][]byte
}

type SPReaderInfo struct {
	Id          string `json:"reader_id"`
	Host        string `json:"host"`
	conn        *RConn
	DidConnect  bool
	Reconnects  int
	LastAttempt time.Time
	isImplicit  bool
	InitCommand [][]byte
}

// A Conn represents a bare connection to a reader.
// It can send and receive []byte payloads.
type Conn struct {
	// Keep all members for which we use atomic at the beginning of the
	// struct and make sure they are all 64bits (or use padding if necessary).
	// atomic.* functions crash on 32bit machines if operand is not aligned
	// at 64bit. See https://github.com/golang/go/issues/599
	// reader
	Opts    Options
	readers map[string]*SPReaderInfo
	inbox   *Msg
}

// SubscriptionType is the type of the Subscription.
type SubscriptionType int

type MsgHandler func(msg *Msg)

type Subscription struct {
	Id         string
	mu         sync.Mutex
	sid        int64
	delivered  uint64
	max        uint64
	conn       *RConn
	mcb        MsgHandler
	closed     bool
	connClosed bool

	// Async linked list
	pHead *Msg
	pTail *Msg
	pCond *sync.Cond

	// Pending stats, async subscriptions, high-speed etc.
	pMsgs       int
	pBytes      int
	pMsgsMax    int
	pBytesMax   int
	pMsgsLimit  int
	pBytesLimit int
	dropped     int
}

// Msg is a structure used by Subscribers and PublishMsg().
type Msg struct {
	From     *Subscription
	Reports  []interface{}
	next     *Msg
	barrier  *barrierInfo
	len_data int
}
type barrierInfo struct {
	refs int64
	f    func()
}

// Tracks various stats received and sent on this connection,
// including counts for messages and bytes.
type Statistics struct {
	InMsgs     uint64
	OutMsgs    uint64
	InBytes    uint64
	OutBytes   uint64
	Reconnects uint64
}
