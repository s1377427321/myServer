package common

import (
	"sync"
	"fmt"
	"data"
	"github.com/gomodule/redigo/redis"
	"lib/db"
)

type SysPrizePool struct {
	sync.RWMutex
}

var (
	sysPrizePool *SysPrizePool
)

func GetSysPrizePool() *SysPrizePool {
	if single == nil {
		sysPrizePool = &SysPrizePool{}
	}

	return sysPrizePool
}


// 增加奖池
func (this *SysPrizePool) AddPrizePool(AppID, Type, SessionID int32, prize int64) {
	this.Lock()
	defer this.Unlock()

	key := fmt.Sprintf(data.Sys_Prize_Pool_Key, AppID, Type, SessionID)

	val, _ := redis.Int64(db.Do("GET", key))

	db.Send("SET", key, prize+val)
}

// 设置奖池
func (this *SysPrizePool) SetPrizePool(AppID, Type, SessionID int32, prize int64) {
	this.Lock()
	defer this.Unlock()

	key := fmt.Sprintf(data.Sys_Prize_Pool_Key, AppID, Type, SessionID)

	db.Send("SET", key, prize)
}

// 获取奖池
func (this *SysPrizePool) GetPrizePool(AppID, Type, SessionID int32) int64 {
	this.Lock()
	defer this.Unlock()

	key := fmt.Sprintf(data.Sys_Prize_Pool_Key, AppID, Type, SessionID)

	val, _ := redis.Int64(db.Do("GET", key))

	return val
}