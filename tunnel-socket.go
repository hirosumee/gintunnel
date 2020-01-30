package gintunnel

import "net"

type TunnelConn struct {
	net.Conn
	host string
}
