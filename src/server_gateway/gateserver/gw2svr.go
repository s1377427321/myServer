package gateserver

import (
	"server_gateway/conf"
	"sync"
	"lib/utils"
)

type Gw2Svr struct {
	autoTempClientID uint32 //自增，临时userid。用户连接上来

	sync.Mutex    // 改为指针类型
	gw2SvrSession map[int32]map[int32]*Gw2SvrSession
	wg            utils.WaitGroupWrapper
	ServerId      int
}

func InitGw2Svr() {
	Hkgw.Gw2Svr = &Gw2Svr{
		gw2SvrSession: make(map[int32]map[int32]*Gw2SvrSession),
		ServerId:      conf.GwCfg.CfgGatewayServer.ServerID}
}



