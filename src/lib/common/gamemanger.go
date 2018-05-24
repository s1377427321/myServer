package common

import (
	"sync"
	"lib/server/adminclient"
	"lib/server/appserver"
)

var (
	single *GameManger
)


type GameManger struct {
	sync.RWMutex
	svr2CenterSession *adminclient.Svr2CenterSession
	nonCenterSvr   *appserver.NonCenterSvr

	IsStop bool // 是否停止游戏
}

func GetGameManger() *GameManger {
	if single == nil {
		single = new(GameManger)
		single.RunGame()
	}

	return single
}


// 开始游戏
func (this *GameManger) RunGame() {
	this.Lock()
	defer this.Unlock()

	this.IsStop = false
}

// 停止游戏
func (this *GameManger) StopGame() {
	this.Lock()
	defer this.Unlock()

	this.IsStop = true
}

