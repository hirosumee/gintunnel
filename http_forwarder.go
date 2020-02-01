package gintunnel

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"net"
	"strings"
)

type Forwarder struct {
	register Register
}

func (f *Forwarder) Start() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		logrus.Fatal(err)
		return
	}
	f.listen(listener)

}
func (f *Forwarder) listen(lis net.Listener) {
	for {
		conn, err := lis.Accept()
		if err != nil {
			logrus.Error(err)
			continue
		}
		go f.handleConn(conn)
	}
}
func (f *Forwarder) handleConn(clientConn net.Conn) {
	buf := make([]byte, 2*1024)
	clientConn.Read(buf)
	ok, host := getHostName(string(buf))
	if !ok {
		sendError(clientConn)
		logrus.Error("cant found hostname in header")
		logrus.Error(string(buf))
		return
	}
	svAddr := f.register.getSvAddr(host)
	if svAddr == "" {
		sendError(clientConn)
		logrus.Infof("can't found server address for %s", host)
		return
	}
	svConn, err := net.Dial("tcp", fmt.Sprintf("%v:8082", svAddr))
	if err != nil {
		sendError(clientConn)
		logrus.Error(err)
		return
	}
	svConn.Write(buf)
	go transfer(clientConn, svConn)
	go transfer(svConn, clientConn)
}
func transfer(src io.ReadCloser, dest io.WriteCloser) {
	defer src.Close()
	defer dest.Close()
	io.Copy(dest, src)
}
func getHostName(str string) (bool, string) {
	index := strings.LastIndex(str, "Host:") + 6
	var endLineIndex = -1
	str2 := str[index:]
	for i, v := range str2[:] {
		if v == 10 {
			endLineIndex = i
			break
		}
	}
	if endLineIndex == -1 {
		return false, ""
	}
	return true, strings.Trim(str2[:endLineIndex], "\r\n")
}

func sendError(conn net.Conn) {
	conn.Write([]byte("HTTP/1.1 503 Service Unavailable Error\r\n\r\n"))
	conn.Write([]byte("Try later..."))
	conn.Write([]byte("\r\n"))
	conn.Close()
}
