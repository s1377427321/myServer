package admin

import (
	"sync"
	"lib/utils"
	"server_console/conf"
	"lib/public_config"
	"errors"
	"strings"
	"lib/socket"
)

// 房间信息
type RoomInfo struct {
	RoomID          uint32
	CreateUID       uint32 // 创建者UID
	ServerID        int32
	GameID          int32
	GameType        int32
	AppID           int32
	Players         []uint32
	ModelKind       uint32 //二进制位表示的,创建房间模式类型(钻石公：1<<1,钻石私：1<<2,房卡：1<<3)
	RoomCreateInfo  string //房间创建信息
	RealPlayerCount int32  // 真实玩家数量
}

// 游戏玩家人数
type GamePlayers struct {
	AppID           int32
	GameID          int32
	GameType        int32
	PlayersCount    int32
	RealPlayerCount int32 // 真实玩家数量
}


type CenterSvr struct {
	sync.RWMutex
	number        uint32
	wg            utils.WaitGroupWrapper
	servers map[uint32][]*Center4SvrSession
	rooms       map[uint32][]*RoomInfo   //key=appid
	gamePlayers map[int32][]*GamePlayers //key=appid 类型不一样?
}

var centerSvr *CenterSvr

func InitCenterSvr() {
	centerSvr = &CenterSvr{
		servers:     make(map[uint32][]*Center4SvrSession),
		rooms:       make(map[uint32][]*RoomInfo),
		gamePlayers: make(map[int32][]*GamePlayers),
	}
}

func StartListenSvrs() {
	svrInfo := public_config.GetCfgCenterServer(conf.CenterCfg.CfgServer.ServerID)
	if svrInfo == nil {
		panic(errors.New("no center-serverId configured"))
	}

	addr := ":" + strings.Split(svrInfo.TcpListenAddr, ":")[1]

	go socket.StartTcpListen(centerSvr, addr)
}

