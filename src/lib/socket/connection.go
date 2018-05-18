package socket

import (
	"sync"
	"net"
	"github.com/astaxie/beego"
	"runtime/debug"
	"os"
	"io"
	"encoding/binary"
	"errors"
	"sync/atomic"
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

type Conner interface {
	StartConnection(rwc *net.Conn) bool
	SetSession(ConnSession)
}

func (this *Connection) SetSession(s ConnSession) {
	if s == nil {
		return
	}

	this.session = s
	this.session.OnInit(this)
}

func (this *Connection) SetIsConn(isConnected bool) {
	this.IsConnected = isConnected
	//this.Stop(false)
}

func (this *Connection) initConnection(rwc *net.Conn) {
	this.remoteAddr = (*rwc).RemoteAddr().String()
	this.localAddr = (*rwc).LocalAddr().String()
	this.rwc = rwc
	this.SetIsConn(true)
}

func (this *Connection) StartConnection (rwc *net.Conn) bool {
	this.initConnection(rwc)

	go this.readLoop()

	return true
}

func (this *Connection) readLoop() {
	defer func() {
		if err := recover(); err != nil {
			beego.Error(string(debug.Stack()))
			os.Exit(0)
		}
	}()

	if this.rwc == nil {
		beego.Warn("The connection has not been established!")
		return
	}

	for {
		pack, err := this.readPck()
		if err != nil {
			this.session.OnClose()
			net.Conn(*this.rwc).Close()
			atomic.StoreInt32(&this.state, StateClose)
			beego.Error(err)
			return
		}


	}
}

func (this *Connection) readPck() (*PackHead, error){
	if n, err := io.ReadFull(*this.rwc, this.hdBuf); n < PackHeadLength || err != nil {
		beego.Error(err)
		return nil, err
	}

	hd := &PackHead{}
	hd.PackFlag = binary.BigEndian.Uint32(this.hdBuf[FlagOffset:SeqOffset])
	hd.SequenceID = binary.BigEndian.Uint32(this.hdBuf[SeqOffset:CmdOffset])
	hd.Cmd = binary.BigEndian.Uint32(this.hdBuf[CmdOffset:UidOffset])
	hd.Uid = binary.BigEndian.Uint32(this.hdBuf[UidOffset:SidOffset])
	hd.Sid = binary.BigEndian.Uint32(this.hdBuf[SidOffset:LenOffset])
	hd.Length = binary.BigEndian.Uint32(this.hdBuf[LenOffset:ReserveOffset])
	hd.Reserve = binary.BigEndian.Uint64(this.hdBuf[ReserveOffset:PackHeadLength])

	if hd.Length > MaxReadBufLen {
		beego.Error("Invalid bodyLenth:", hd.Length)
		return nil, errors.New("Invalid bodyLen")
	}

	return hd, nil
}



