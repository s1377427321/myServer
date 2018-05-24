package common

import (
	"encoding/json"
	"lib/db"
	"strconv"
	"data"
)

//设置等级概率
func SetRateLevel(level, Type int, s string) {
	if Type == 1 {
		m := make(map[int]int)
		err := json.Unmarshal([]byte(s), &m)
		if err != nil {
			return
		}
		if len(m) != 13 {
			return
		}
		//设置
		db.Send("HSET", data.Probability_Level_Key, strconv.Itoa(level), s)
	} else {
		//解除
		db.Send("HDEL", data.Probability_Level_Key, strconv.Itoa(level))
	}
}

//设置玩家等级
func SetPlayerLevel(uid int64, level, Type int) {
	if Type == 1 {
		//设置
		db.Send("HSET", data.Login_User_Key, strconv.Itoa(int(uid)), level)
	} else {
		//解除
		db.Send("HDEL", data.Login_User_Key, strconv.Itoa(int(uid)))
	}
}

//玩家id列表s:"[1,2,3,4]"
func SetPlayersValue(s string, value int64, Type int) {
	uids := make([]int64, 0)
	err := json.Unmarshal([]byte(s), &uids)
	if err != nil {
		return
	}
	for _, v := range uids {
		SetPlayerValue(v, value, Type)
	}
}

//玩家输赢值
func SetPlayerValue(uid int64, num int64, Type int) {
	if Type == 1 {
		//设置
		db.Send("HSET", data.WinloseValue_Key, strconv.Itoa(int(uid)), num)
	} else {
		//解除
		db.Send("HDEL", data.WinloseValue_Key, strconv.Itoa(int(uid)))
	}
}
