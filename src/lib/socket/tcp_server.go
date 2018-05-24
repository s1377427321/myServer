package socket

import (
	"net"
	"github.com/astaxie/beego"
)

type Server interface {
	OnCreateConnection(c Conner)
	OnClose()
}

//启动tcp server
func StartTcpListen(svr Server, address string) {
	defer svr.OnClose()

	addr, e := net.ResolveTCPAddr("tcp", address)

	if e != nil {
		panic(e.Error())
	}

	l, e := net.ListenTCP("tcp", addr)
	if e != nil {
		panic(e.Error())
	}
	defer l.Close()

	for {
		rw, e := l.Accept()
		if e != nil {
			if ne, ok := e.(net.Error); ok && ne.Temporary() {
				continue
			}

			beego.Error("[TCPServer] accept error: %v", e)
			return
		}

		// 设置TCP参数
		if rw.(*net.TCPConn) != nil {
			rw.(*net.TCPConn).SetKeepAlive(false)
			rw.(*net.TCPConn).SetLinger(0)
		}

		c := CreateConnection(address, true, false, nil)
		svr.OnCreateConnection(c)
		go c.StartConnection(&rw)
	}

}


