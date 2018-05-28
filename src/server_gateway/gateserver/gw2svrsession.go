package gateserver

import "lib/socket"

type Gw2SvrSession struct {
	connection    *socket.Connection
	server        *Gw2Svr
	ServerId      int32
	ServerType    int32
	ServerAddress string
	AppId         int32

	//游戏中房间号生成的ID范围
	CodeSectionMin int32
	CodeSectionMax int32
}