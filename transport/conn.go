package transport

import (
	"errors"
	"net"
	"os"

	"github.com/VolantMQ/volantmq/auth"
	"github.com/VolantMQ/volantmq/metrics"
)

// Conn is wrapper to net.Conn
// implemented to encapsulate bytes statistic
type Conn interface {
	net.Conn
	GlobalMessageChannel() chan *Message
	GlobalEventChannel() chan *Event
}

type conn struct {
	net.Conn
	globalMessageChannel chan *Message
	globalEventChannel   chan *Event

	stat metrics.Bytes
}

var _ Conn = (*conn)(nil)

// Handler ...
type Handler interface {
	OnConnection(Conn, *auth.Manager) error
}

func newConn(cn net.Conn, globalMessageChannel chan *Message, globalEventChannel chan *Event, stat metrics.Bytes) *conn {
	c := &conn{
		Conn:                 cn,
		globalMessageChannel: globalMessageChannel,
		globalEventChannel:   globalEventChannel,
		stat:                 stat,
	}

	return c
}

func (c *conn) GlobalMessageChannel() chan *Message {
	return c.globalMessageChannel
}
func (c *conn) GlobalEventChannel() chan *Event {
	return c.globalEventChannel
}

// Read ...
func (c *conn) Read(b []byte) (int, error) {
	n, err := c.Conn.Read(b)

	c.stat.OnRecv(n)

	return n, err
}

// Write ...
func (c *conn) Write(b []byte) (int, error) {
	n, err := c.Conn.Write(b)
	c.stat.OnSent(n)

	return n, err
}

// File ...
func (c *conn) File() (*os.File, error) {
	if t, ok := c.Conn.(*net.TCPConn); ok {
		return t.File()
	}

	return nil, errors.New("not implemented")
}
