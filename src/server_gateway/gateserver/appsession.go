package gateserver

import (
	"lib/socket"
	"server_gateway/typeconf"
)

type AppSession struct {
	connection   *socket.Connection
	remoteAddr   string
	sid          typeconf.ServerID
	State        uint32 //user.proto-->UserOnlineStatus枚举。
	uid          typeconf.UserID
	server       *GateWaySvr
	AppId        uint32
	GameServerID int32  // 玩家所在游戏ID
	InviteCode   uint32 // 桌子ID

	UserType uint32 // 如果不是机器人 则为1， 否则为0  暂时
}


func NewAppSession(svr *GateWaySvr) *AppSession {
	//beego.Info("NewAppSession")
	return &AppSession{server: svr}
}

