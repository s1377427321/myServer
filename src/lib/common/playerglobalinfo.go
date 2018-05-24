package common

import (
	"fmt"
	"data"
	"github.com/gomodule/redigo/redis"
	"lib/db"
)

func GetPlayerGameInviteCode(uid uint32) int {
	key := fmt.Sprintf(data.Game_User_Info_Key, uid)
	gsInviteCode, rerr := redis.Int(db.Do("HGET", key, data.Game_User_Info_GameInviteCode_Key))
	if rerr != nil {
		gsInviteCode = 0
	}
	return gsInviteCode
}

func GetPlayerGameSvrInfo(uid uint32) int {
	key := fmt.Sprintf(data.Game_User_Info_Key, uid)
	gameSvrID, rerr := redis.Int(db.Do("HGET", key, data.Game_User_Info_GameSvrID_Key))
	if rerr != nil {
		gameSvrID = 0
	}
	return gameSvrID
}



func GetPlayerProbabilityInfo(uid int64) (int, int, int64) {
	key := fmt.Sprintf(data.Login_User_Key, uid)

	Level, rerr := redis.Int(db.Do("HGET", key, data.Probability_Level_Key))
	if rerr != nil {
		Level = 0
	}
	Level1, rerr := redis.Int(db.Do("HGET", key, data.Probability_Level1_Key))
	if rerr != nil {
		Level1 = 0
	}
	WinloseValue, rerr := redis.Int64(db.Do("HGET", key, data.WinloseValue_Key))
	if rerr != nil {
		WinloseValue = 0
	}

	return Level, Level1, WinloseValue
}

// 获取系统信息（百人牛牛，二八杠）
func GetSysProbabilityInfo(appid, apptype, sessionid int) (int, int, int, int, int, int) {
	key := fmt.Sprintf(data.Game_Key, appid, apptype, sessionid)
	//var uid int64 = 0
	PlayerBankerLevel, rerr := redis.Int(db.Do("HGET", key, data.Probability_PlayerBankerLevel_Key))
	if rerr != nil {
		PlayerBankerLevel = 0
	}
	PlayerBankerLevel1, rerr := redis.Int(db.Do("HGET", key, data.Probability_PlayerBankerLevel1_Key))
	if rerr != nil {
		PlayerBankerLevel1 = 0
	}
	SysBankerLevel, rerr := redis.Int(db.Do("HGET", key, data.Probability_SysBankerLevel_Key))
	if rerr != nil {
		SysBankerLevel = 0
	}
	SysBankerLevel1, rerr := redis.Int(db.Do("HGET", key, data.Probability_SysBankerLevel1_Key))
	if rerr != nil {
		SysBankerLevel1 = 0
	}
	PlayerLevel, rerr := redis.Int(db.Do("HGET", key, data.Probability_PlayerLevel_Key))
	if rerr != nil {
		PlayerLevel = 0
	}
	PlayerLevel1, rerr := redis.Int(db.Do("HGET", key, data.Probability_PlayerLevel1_Key))
	if rerr != nil {
		PlayerLevel1 = 0
	}

	return PlayerBankerLevel, PlayerBankerLevel1, SysBankerLevel, SysBankerLevel1, PlayerLevel, PlayerLevel1
}

// 获取玩家输赢金币信息
func GetPlayerWinGold(uid int64) int64 {
	key := fmt.Sprintf(data.Login_User_Key, uid)
	WinGold, rerr := redis.Int64(db.Do("HGET", key, data.WinGold_Key))
	if rerr != nil {
		WinGold = 0
	}

	return WinGold
}