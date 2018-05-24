package data

import (
	"lib/db"
	"strconv"
)

// 玩家进入游戏
func UserEnterGame(Uid uint32) {
	db.Send("HSET", Login_User_Key+strconv.Itoa(int(Uid)), "gameing", 1)
}

// 玩家退出游戏
func UserOutGame(Uid uint32) {
	db.Send("HSET", Login_User_Key+strconv.Itoa(int(Uid)), "gameing", 0)
}
