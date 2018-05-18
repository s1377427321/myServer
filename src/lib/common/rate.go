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