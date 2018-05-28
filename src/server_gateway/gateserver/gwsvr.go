package gateserver

import (
	"sync"
	"lib/server/adminclient"
	"server_gateway/typeconf"
)

type GateWaySvr struct {
	//clientID uint32
	autoTempClientID uint32 //自增，临时userid。用户连接上来

	sync.Mutex   // 改为指针类型
	Appsession   map[typeconf.UserID]*AppSession	// 客户端连接对应的session
	Gw2Svr *Gw2Svr	// 描述内部服务器的数据结构
	Svr2CenterSession *adminclient.Svr2CenterSession	// 描述连接到center的数据结构

	LoginUserCount      int32
	LoginUserTotalCount int32
}

var Hkgw *GateWaySvr