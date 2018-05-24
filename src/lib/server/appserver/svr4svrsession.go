package appserver

import (
	"lib/socket"
	"github.com/astaxie/beego"
	"protocol"
	"runtime/debug"
	"os"
	"lib/route"
	"reflect"
	"lib/server"
	"code.google.com/p/goprotobuf/proto"
)

//连接
type Svr4SvrSession struct {
	connection *socket.Connection //链接
	sessionID  uint32
	server     *NonCenterSvr //app服务器
	serverID   int32
}

/*
type ConnSession interface {
	OnInit(*Connection)      //初始化操作，比如心跳的设置...
	OnProcessPack(*PackHead) //处理消息
	OnClose()                //消除所有对Sessioner的引用,心跳...
	Write(msg interface{}, ph *PackHead) error
	Close()
	OnConnect(isOk bool)
}

 */

func (this *Svr4SvrSession) Close() {
	this.connection.SetIsConn(false)
	this.connection.Stop(true)
	if this.serverID != 0 {
		this.server.RemoveGateway(this.serverID)
	}
}

//初始化
func (this *Svr4SvrSession) OnInit(c *socket.Connection) {
	this.connection = c
}

//新建Svr4SvrSession
func CreateSvrSession(svr *NonCenterSvr) *Svr4SvrSession {
	return &Svr4SvrSession{server: svr}
}

func (this *Svr4SvrSession) OnProcessPack(ph *socket.PackHead) {
	if this.server != nil && this.server.UseRunTime {
		go this.ProcessPack(ph)
	} else {
		this.ProcessPack(ph)
	}
}

func (this *Svr4SvrSession) ProcessPack(ph *socket.PackHead) {

	//beego.Info("服务器收到数据cmd:", ph.Cmd, "uid:", ph.Uid)
	defer func() {
		if err := recover(); err != nil {
			beego.Error("ProcessPack ", err)
			beego.Error("Process call failed CMD : %d %s %v %d", ph.Cmd, protocol.CMD(ph.Cmd).String(), err, ph.Length)
			beego.Error(string(debug.Stack()))

			os.Exit(0)
		}
	}()

	//beego.Info("收到协议号：", ph.Cmd)
	data := ph.Body

	// 网关注册
	if ph.Cmd == uint32(protocol.CMD_GATEWAY_REGIST_TO_INTERNAL) {

		regist := &protocol.GateWayRegist{}
		err := proto.Unmarshal(ph.Body, regist)
		if err == nil {
			this.serverID = int32(*regist.ServerId)
			this.server.RegistGateway(int32(*regist.ServerId), this)
		} else {
			beego.Error("protocol  unmarshal fail: ", err)
		}
		return
	}

	if route.Exist(ph.Cmd) {
		//beego.Info("OnProcessPack==>", ph.Cmd)
		hd := route.GetFunc(ph.Cmd)
		t := route.GetProto(ph.Cmd)
		v := reflect.New(t)
		if err := proto.Unmarshal(data, v.Interface().(proto.Message)); err == nil {
			session := &server.TmpSession{
				Uid:      ph.Uid,
				Sid:      ph.Sid,
				ServerID: this.serverID,
			}
			session.OnInit(this.connection)

			result := hd.Call([]reflect.Value{reflect.ValueOf(session), v})
			if result[1].IsNil() && !result[0].IsNil() {
				this.connection.Write(ph, result[0].Interface())
			}
		} else {
			beego.Error("protocol  unmarshal fail: ", err)
		}
	} else {
		beego.Error("not found command %d %s", ph.Cmd, protocol.CMD(ph.Cmd).String())
	}
}


func (this *Svr4SvrSession) OnClose() {

}

func (this *Svr4SvrSession) Write(msg interface{}, ph *socket.PackHead) error {
	return this.connection.Write(ph, msg)
}

func (this *Svr4SvrSession) Close() {
	this.connection.SetIsConn(false)
	this.connection.Stop(true)
	if this.serverID != 0 {
		this.server.RemoveGateway(this.serverID)
	}
}

func (this *Svr4SvrSession) OnConnect(isOk bool) {

}

