package admin

import "sync"

type GameRevenue struct {
	sync.RWMutex

	GameRevenueInfos []*GameRevenueInfo
}

type GameRevenueInfo struct {
	AppId     int32
	AppType   uint32
	ServerID  uint32
	RoomType  uint32
	SessionID int32

	SysWinGold   int64 // 系统输赢值
	RobotWinGold int64 // 机器人输赢值
	InputGold    int64 // 投入
	OutputGold   int64 // 产出
}

var server_GameRevenue *GameRevenue

func GetGameRevenue() *GameRevenue {
	return server_GameRevenue
}


func (this *GameRevenue) GetRevenueValue(AppId int32, RoomType uint32, SessionID int32, lock bool) (int64, int64, int64, int64) {
	if lock {
		this.Lock()
		this.Unlock()
	}

	var SysWinGold int64
	var RobotWinGold int64
	var InputGold int64
	var OutputGold int64
	for _, v := range this.GameRevenueInfos {
		if (SessionID != 0 && v.AppId == AppId && v.RoomType == RoomType && v.SessionID == SessionID) ||
			(SessionID == 0 && v.AppId == AppId && v.RoomType == RoomType) {
			SysWinGold += v.SysWinGold
			RobotWinGold += v.RobotWinGold
			InputGold += v.InputGold
			OutputGold += v.OutputGold
		}
	}
	return SysWinGold, RobotWinGold, InputGold, OutputGold
}