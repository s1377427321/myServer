package server

import (
	"lib/socket"
	"github.com/astaxie/beego"
	"protocol"
	"lib/route"
	"reflect"
	"code.google.com/p/goprotobuf/proto"
)

type TmpSession struct {
	Uid          uint32
	Sid          uint32
	ServiceType_ int32
	AppId        int32
	ServerID     int32
	connection   *socket.Connection
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

func (this *TmpSession) OnInit(c *socket.Connection) {
	//beego.Info("GateClient->OnInit....")
	this.connection = c
}

func (this *TmpSession) OnProcessPack(ph *socket.PackHead) {
	go this.ProcessPack(ph)
}

func (this *TmpSession) ProcessPack(ph *socket.PackHead) {
	//beego.Info("游戏服务器收到数据cmd:", ph.Cmd, "uid:", ph.Uid)
	defer func() {
		if err := recover(); err != nil {
			beego.Error("Process call failed CMD : %d %s %v %d", ph.Cmd, protocol.CMD(ph.Cmd).String(), err, ph.Length)
		}
	}()

	//beego.Info("收到协议号：", ph.Cmd)
	data := ph.Body

	if route.Exist(ph.Cmd) {
		//beego.Info("OnProcessPack==>", ph.Cmd)
		hd := route.GetFunc(ph.Cmd)
		t := route.GetProto(ph.Cmd)
		v := reflect.New(t)
		// v.Interface的得到的value是proto.Message的一个实现，所以可以断言
		if err := proto.Unmarshal(data, v.Interface().(proto.Message)); err == nil {
			session := &TmpSession{
				Uid:          ph.Uid,
				Sid:          ph.Sid,
				connection:   this.connection,
				ServiceType_: this.ServiceType_,
			}

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

func (this *TmpSession) OnClose() {
}

func (this *TmpSession) Write(msg interface{}, head *socket.PackHead) error{
	return this.connection.Write(head, msg)
}

func (this *TmpSession) Close() {
	//beego.Info("TmpSession===>Close()", this.Sid)
	this.connection.Stop(false)
}

func (this *TmpSession) OnConnect(isOk bool) {
	if isOk {
		this.RegistService()
	} else {
		//beego.Info("网关连接失败，退出进程")
		//os.Exit(-1)
	}
}

//注册服务。
func (this *TmpSession) RegistService() {

	regist := &protocol.ServerRegist{
		ServerId:    proto.Int32(int32(this.Sid)),
		ServiceType: proto.Int32(this.ServiceType_),
		AppId:       proto.Int32(this.AppId),
	}
	//beego.Info("服务注册。。。", int32(this.Sid), this.ServiceType_, regist)
	head := &socket.PackHead{Sid: uint32(this.Sid), Uid: 0, Cmd: uint32(protocol.CMD_SERVICEREGIST)}
	this.Write(regist, head)
}



