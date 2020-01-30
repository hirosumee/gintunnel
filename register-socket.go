package gintunnel

import (
	"net"
	"sync"
)

type RegisterConn struct {
	net.Conn
	ping     PingPong
	closed   Closed
	hostname string
}

func NewRegisterConn(conn net.Conn) RegisterConn {
	return RegisterConn{conn, PingPong{}, Closed{}, ""}
}

type PingPong struct {
	ping bool
	rw   sync.RWMutex
}

func (p *PingPong) get() bool {
	p.rw.RLock()
	defer p.rw.RUnlock()
	return p.ping
}
func (p *PingPong) set(value bool) {
	p.rw.Lock()
	p.ping = value
	p.rw.Unlock()
}

type Closed struct {
	ping bool
	rw   sync.RWMutex
}

func (c *Closed) get() bool {
	c.rw.RLock()
	defer c.rw.RUnlock()
	return c.ping
}

func (c *Closed) set(value bool) {
	c.rw.Lock()
	c.ping = value
	c.rw.Unlock()
}
