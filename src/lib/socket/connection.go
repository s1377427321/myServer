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
	"time"
	"fmt"
)

const (
	StateInit = iota
	StateConnected
	StateManualClosed
	StateClose
	StateConnecting
	StateHalt

)

const maxSendBufferSize int = 10240

type OutConns struct {
	connMap map[string]*Connection	// key:peer address
	sync.Mutex
}

var(
	outConns *OutConns
	connId uint32
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

func CreateConnection(addr string, is_server bool, is_center_server bool, session ConnSession) *Connection {
	conn := &Connection{
		is_server:        is_server,
		is_center_server: is_center_server,
		addr_:            addr,
		hdBuf: make([]byte, PackHeadLength),
	}

	if session != nil {
		conn.SetSession(session)
		atomic.StoreInt32(&conn.state, StateInit)
		conn.Id = AllocConnId()

		Add2OutConnMap(conn)
	}

	return conn
}

func Add2OutConnMap(conn *Connection) {
	outConns.Lock()
	defer outConns.Unlock()

	k := GenConnKey(conn)

	_, exists := outConns.connMap[k]; if exists {
		beego.Warn("The key:", k, "has ever exists!")
	} else {
		outConns.connMap[k] = conn
	}
}

func GenConnKey(conn *Connection) string{
	return fmt.Sprintf("%s%_%d", conn.addr_, conn.Id)
}

func AllocConnId() uint32 {
	return atomic.AddUint32(&connId, 1)
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

		this.DispatchMsg(pack)
	}
}

func (this *Connection)DispatchMsg(pack *PackHead) {
	if pack == nil {
		beego.Error("Inside DispatchMsg, pack is nil!")
		return
	}
	if this.session == nil {
		beego.Error("Inside DispatchMsg, this.session is nil!")
		return
	}

	this.session.OnProcessPack(pack)

}

func (this *Connection) Write(ph *PackHead, msg interface{}) error{
	_, data, err := SerializePackWithPB(ph, msg, maxSendBufferSize)

	if err != nil {
		beego.Error("Inside Connection.Write(), SerializePackWithPB failed!")
		return err
	}

	if err = this.WriteData(data); err != nil {
		beego.Error(err)
		return err
	}

	return  nil
}

func (this *Connection) WriteData(msg []byte) error {
	var retErr error

	defer func() {
		if retErr != nil && this.rwc != nil {
			(*this.rwc).Close()
		}
	}()


	if this.rwc == nil {
		retErr = errors.New("Inside Connection.WriteData(), rwc is nil!")
		beego.Error(retErr)
		return retErr
	}

	retErr = (*this.rwc).SetWriteDeadline(time.Now().Add(500 * time.Millisecond))
	if retErr != nil {
		beego.Error(retErr)
		return retErr
	}

	var n int
	n, retErr = (*this.rwc).Write(msg)
	if retErr != nil {
		beego.Error(retErr)
		return retErr
	}

	if n != len(msg) {
		retErr = errors.New("Only partial data written!")
		beego.Error(retErr)
		return retErr
	}

	return nil
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


//is_manual_close 主动关闭，如果是主动关闭，则不重连。
func (this *Connection) Stop(is_manual_close bool) {

	//this.Lock()
	//defer this.Unlock()

	//beego.Info("Stop,remoteAddr", this.remoteAddr, "state:", this.state, "is_manual_close:", is_manual_close)
	if is_manual_close {
		if !atomic.CompareAndSwapInt32(&this.state, StateConnected, StateManualClosed) {
			//beego.Info("手动关闭连接原子更新失败。。。逻辑bug,state:", this.state)
			return
		}
	} else {
		if !atomic.CompareAndSwapInt32(&this.state, StateConnected, StateInit) {
			//beego.Info("非主动关闭连接，原子更新失败。。。逻辑bug,state:", this.state)
			return
		}
	}
	//close(this.CloseChan)
	if this.rwc != nil {
		(*this.rwc).Close()
	}
	if this.session != nil && !is_manual_close {
		go this.session.OnClose()
	}
	//this.SetIsConn(false)
	//this.IsConnected = false
}


