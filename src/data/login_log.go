package data

import (
	"time"
	"fmt"
)

type PfGameOnlineLog struct {
	ID        uint64    `xorm:"'id' not null pk autoincr INTEGER"` // 记录自增ID
	HourDate  int       `xorm:"'hour_date' INTEGER"`               // 时间:2017120215  小时为单位
	DayDate   int       `xorm:"'day_date' INTEGER"`                // 时间:20171202
	MonthDate int       `xorm:"'month_date' INTEGER"`              // 时间:201712
	Appid     int32     `xorm:"'appid' INTEGER"`                   // 游戏id(如:二八杠,牛牛)
	Sessionid int32     `xorm:"'sessionid' SMALLINT"`              // 场次id(如:土豪场,富豪场)
	GameType  int32     `xorm:"'game_type' SMALLINT"`              //
	Count     int32     `xorm:"'count' INTEGER"`                   //
	Ctime     time.Time `xorm:"'ctime' DATETIME"`                  // 创建时间
}

func (this *PfGameOnlineLog) Save() (int64, error) {
	sql := fmt.Sprintf("INSERT INTO public.pf_game_online_log " +
		"(hour_date,day_date,month_date,appid,sessionid,game_type,count,ctime) " +
		"VALUES (%d, %d, %d, %d, %d, %d, %d, '%v');",
		HourDate(), DayDate(), MonthDate(), this.Appid, this.Sessionid, this.GameType, this.Count, Unix2Str())

	dbagent.GetDBAgent().Exec(3, sql, nil)

	return 0, nil
}

/**
 * 获取本地当天时间2017040215
 * @return int
 */
func HourDate() int {
	now := time.Now()
	return now.Year()*1000000 + int(now.Month())*10000 + now.Day()*100 + now.Hour()
}

/**
 * 获取本地当天时间20170402
 * @return int
 */
func DayDate() int {
	now := time.Now()
	return now.Year()*10000 + int(now.Month())*100 + now.Day()
}

/**
 * 获取本地当月时间201704
 * @return int
 */
func MonthDate() int {
	now := time.Now()
	return now.Year()*10000 + int(now.Month())
}