package admin

import "lib/socket"

//客户端连接
type Center4SvrSession struct {
	connection    *socket.Connection //链接
	remoteAddr    string             //客户端地址
	service_type  uint32             //客户端类型
	serverID      uint32             //客户端id
	server        *CenterSvr       //admin服务器
	ServerAddress string             //注册服务器ip端口
	AppId         int32

	UserCount      uint32 //玩家人数
	TotalUserCount uint32 //总玩家人数

	//游戏中房间号生成的ID范围
	CodeSectionMin int32
	CodeSectionMax int32
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
func (this *Center4SvrSession) OnInit(c *socket.Connection) {
	this.connection = c
}

func (this *Center4SvrSession) OnProcessPack(ph *socket.PackHead) {
	return  nil
}

func (this *Center4SvrSession) OnClose() {
	this.server.RemoveServer(this)
}

func (this *Center4SvrSession) Write(msg interface{}, ph *socket.PackHead) error {
	go this.connection.Write(ph, msg)
	return nil
}

func (this *Center4SvrSession) OnConnect(isOk bool) {

}

func (this *Center4SvrSession) Close() {

}









