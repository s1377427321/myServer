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


