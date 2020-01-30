package gintunnel

import (
	"bufio"
	"github.com/sirupsen/logrus"
	"io"
	"net"
	"strings"
	"time"
)

type Register struct {
	fm ForwardMap
}

func (r *Register) Start() {
	listener, err := net.Listen("tcp", ":8081")
	if err != nil {
		logrus.Fatal(err)
		return
	}
	r.listen(listener)

}
func (r *Register) listen(listener net.Listener) {
	for {
		conn, err := (listener).Accept()
		if err != nil {
			logrus.Error(err)
			continue
		}
		t := NewRegisterConn(conn)
		go r.handleConnection(&t)
	}
}
func (r *Register) handleConnection(conn *RegisterConn) {
	logrus.Infof("Serving client %s", conn.RemoteAddr())
	addr := conn.RemoteAddr().String()
	reader := bufio.NewReader(conn)
	go pingPong(conn)
	defer r.closeConn(conn)
	for {
		if conn.closed.get() {
			return
		}
		temp, err := reader.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				logrus.Error(err)
			}
			break
		}
		raw := strings.TrimSpace(temp)
		cmd := strings.Split(raw, " ")[0]
		data := strings.TrimSpace(raw[len(cmd):])
		switch cmd {
		case "REG":
			{
				i := strings.LastIndex(addr, ":")
				ok := r.fm.set(data, addr[:i])
				if ok {
					//assign hostname to conn
					logrus.Infof("registered %s with %s", data, addr[:i])
					conn.hostname = data
					conn.Write([]byte("REG-RES success\r\n"))
				} else {
					conn.Write([]byte("REG-RES used\r\n"))
				}
				break
			}
		case "PONG":
			{
				conn.ping.set(false)
				break
			}
		}
	}
}
func (r *Register) closeConn(conn *RegisterConn) {
	if conn.hostname != "" {
		r.fm.remove(conn.hostname)
	}
	conn.Close()
	logrus.Info("closed connection")
}
func (r *Register) getSvAddr(host string) string {
	return r.fm.get(host)
}

func pingPong(conn *RegisterConn) {
	for {
		if conn.closed.get() {
			return
		}
		if conn.ping.get() {
			//set closed
			conn.closed.set(true)
			return
		}
		conn.ping.set(true)
		conn.Write([]byte("PING\r\n"))
		time.Sleep(2 * time.Second)
	}
}
