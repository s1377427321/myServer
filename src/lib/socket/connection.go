package socket

import (
	"sync"
	"net"
)

type ConnSession interface {
	OnInit(*Connection)      //初始化操作，比如心跳的设置...
	OnProcessPack(*PackHead) //处理消息
	OnClose()                //消除所有对Sessioner的引用,心跳...
	Write(msg interface{}, ph *PackHead) error
	Close()
	OnConnect(isOk bool)
}

type Connection struct {
	sync.RWMutex
	session          ConnSession
	addr_            string
	is_server        bool //是服务端connection类型。
	is_center_server bool //是中心服务端connection类型。
	state            int32
	rwc             *net.Conn
	remoteAddr      string
	localAddr       string
	hdBuf        	[]byte
	IsConnected     bool
	Id uint32	// use it to identify multiple connection to a same peer.
}

