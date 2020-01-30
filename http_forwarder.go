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
	buf := make([]byte, 100)
	clientConn.Read(buf)
	ok, host := getHostName(string(buf))
	if !ok {
		logrus.Info("cant found hostname in header")
		return
	}
	svAddr := f.register.getSvAddr(host)
	if svAddr == "" {
		logrus.Infof("can't found server address for %s", host)
		return
	}
	svConn, err := net.Dial("tcp", fmt.Sprintf("%v:8082", svAddr))
	if err != nil {
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
