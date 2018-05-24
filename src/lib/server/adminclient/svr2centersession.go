package adminclient

import (
	"sync"
	"lib/socket"
	"github.com/astaxie/beego"
	"errors"
	"protocol"
	"code.google.com/p/goprotobuf/proto"
	"lib/route"
	"reflect"
	"lib/server"
	"os"
	"lib/event"
)

type Svr2CenterSession struct {
	Lock              sync.RWMutex
	Uid               uint32
	Sid               uint32
	ServiceType_      int32
	ServerAddr_       string
	AppId             int
	connection        *socket.Connection
	GateWayConnection *socket.Connection

	Dispatch *event.Dispatcher

	//游戏中房间号生成的ID范围
	CodeSectionMin_ int32
	CodeSectionMax_ int32
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

//初始化
func (this *Svr2CenterSession) OnInit(c *socket.Connection) {
	//beego.Info("GateClient->OnInit....")
	this.connection = c
}

func (this *Svr2CenterSession) OnClose() {
}

func (this *Svr2CenterSession) Write(msg interface{}, head *socket.PackHead) error{
	if this.connection == nil {
		beego.Error("Inside Svr2CenterSession.Write(), connection is nil!")
		return errors.New("Inside Svr2CenterSession.Write(), connection is nil!")
	}

	return this.connection.Write(head, msg)
}

func (this *Svr2CenterSession) Close() {
	//beego.Info("TmpSession===>Close()", this.Sid)
	this.connection.Stop(false)
}

func (this *Svr2CenterSession) OnConnect(isOk bool) {
	if isOk {
		this.RegistServer()
	}
}

func (this *Svr2CenterSession) RegistServer() {
	regist := &protocol.NonCenterSvrRegist{
		ServerId:       proto.Int32(int32(this.Sid)),
		ServerType:     proto.Int32(this.ServiceType_),
		ServerAddress:  proto.String(this.ServerAddr_),
		AppId:          proto.Int32(int32(this.AppId)),
		CodeSectionMin: proto.Int32(int32(this.CodeSectionMin_)),
		CodeSectionMax: proto.Int32(int32(this.CodeSectionMax_)),
	}

	head := &socket.PackHead{Sid: uint32(this.Sid), Uid: 0, Cmd: uint32(protocol.CMD_AMPQ_APP_SERVER_REGIST)}

	err := this.Write(regist, head)
	if err != nil {
		beego.Error("Send register request to center failed:", err)
	} else {
		beego.Warn("Send register request to center....")
	}
}

//回包关键接口函数
func (this *Svr2CenterSession) OnProcessPack(ph *socket.PackHead) {
	if route.Exist(ph.Cmd) {
		data := ph.Body
		hd := route.GetFunc(ph.Cmd)
		t := route.GetProto(ph.Cmd)
		v := reflect.New(t)
		if err := proto.Unmarshal(data, v.Interface().(proto.Message)); err == nil {
			seesion := &server.TmpSession{
				Uid: ph.Uid,
				Sid: ph.Sid,
			}
			seesion.OnInit(this.connection)

			result := hd.Call([]reflect.Value{reflect.ValueOf(seesion), v})
			if result[1].IsNil() && !result[0].IsNil() {
				this.connection.Write(ph, result[0].Interface())
			}
		}
		return
	}

	if ph.Cmd == uint32(protocol.CMD_AMPQ_APP_SERVER_REGIST_RESULT) {
		appServerRegistResult := &protocol.NonCenterSvrRegistResult{}
		err := proto.Unmarshal(ph.Body, appServerRegistResult)
		if err != nil || appServerRegistResult.GetResult() != int32(0) {
			//beego.Error("服务注册失败>>>", string(ph.Body))
			os.Exit(0) // 退出进程
		} else {
			//beego.Error("服务注册成功>>>", string(ph.Body))
		}
		return
	}

	if this.Sid != ph.Sid || (this.Sid/100 != ph.Sid && ph.Sid < 100) {
		//beego.Info("不是发给本服务器的包 cmd>>>", ph.Cmd, " 应给>>>", ph.Sid, " 我是>>>", this.Sid)
		return
	}

	//失败回包
	if ph.Cmd == uint32(protocol.CMD_AMPQ_FAIL) {
		//beego.Error("admin服务器发包失败，原因>>>", string(ph.Body))
		return
	}

	this.Dispatch.Dispatch("OnProcessPack", ph)
	return


}