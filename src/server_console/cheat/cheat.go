package cheat

import (
	"lib/public_config"
	"server_console/conf"
	"github.com/astaxie/beego"
	"os"
	"strings"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"encoding/json"
	"net/http"
	"fmt"
	"protocol"
	"code.google.com/p/goprotobuf/proto"
	"lib/common"
	"server_console/admin"
	"strconv"
)

type H map[string]interface{}

func StartCheat() {
	svrInfo := public_config.GetCfgCenterServer(conf.CenterCfg.CfgServer.ServerID)
	if svrInfo == nil {
		beego.Error("centerServer listen http fail!")
		os.Exit(0)
	}

	ports := ":" + strings.Split(svrInfo.HttpListenAddr, ":")[1]

	Run(ports)
}

func Run(port string) {
	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.Gzip())
	e.Use(middleware.Secure())
	e.Use(middleware.BodyLimit("1M"))

	e.Static("/", "./assets")

	e.GET("/notice", notice)
	e.GET("/updatePlatform", updatePlatform)
	e.GET("/updateConfig", updateGameConfig)
	e.GET("/kickplayer", kickPlayer)
	e.GET("/getrevenueinfo", getRevenueInfo)
	e.GET("/updateprobabilityinfo", updateProbabilityInfo)
	e.GET("/getprobabilityinfo", getProbabilityInfo)
	e.GET("/setprobability", setProbability)
	e.GET("/getplayerwingold", getPlayerWinGold)
	e.GET("/changegameserverstatus", changeGameServerStatus)
	e.GET("/refgameconfig", refGameConfig)
	e.GET("/getsysprizepool", getSysPrizePool)
	e.GET("/setsysprizepool", setSysPrizePool)

	e.Start(port)

}

// 设置奖池
func setSysPrizePool(c echo.Context) error {
	appid, _ := strconv.Atoi(c.QueryParam("appid"))
	apptype, _ := strconv.Atoi(c.QueryParam("apptype"))
	sessionid, _ := strconv.Atoi(c.QueryParam("sessionid"))
	prize, _ := strconv.Atoi(c.QueryParam("prize"))

	// 百人场的
	if (appid == 102 || appid == 109 || appid == 107) && sessionid > 1 {
		sessionid--
	}

	common.GetSysPrizePool().SetPrizePool(int32(appid), int32(apptype), int32(sessionid), int64(prize))

	return c.JSON(http.StatusOK, H{"msg": "成功"})
}


// 获取奖池
func getSysPrizePool(c echo.Context) error {
	appid, _ := strconv.Atoi(c.QueryParam("appid"))
	apptype, _ := strconv.Atoi(c.QueryParam("apptype"))
	sessionid, _ := strconv.Atoi(c.QueryParam("sessionid"))

	// 百人场的
	if (appid == 102 || appid == 109 || appid == 107) && sessionid > 1 {
		sessionid--
	}

	sysPrizePool := common.GetSysPrizePool().GetPrizePool(int32(appid), int32(apptype), int32(sessionid))

	var dh H = make(H, 0)
	dh["SysPrizePool"] = sysPrizePool

	return c.JSON(http.StatusOK, dh)
}

// 刷新配置文件
func refGameConfig(c echo.Context) error {
	Type, _ := strconv.Atoi(c.QueryParam("type"))
	admin.RefGameConfig(int32(Type))
	return c.JSON(http.StatusOK, H{"msg": "成功"})
}

// 改变游戏服务器状态
func changeGameServerStatus(c echo.Context) error {
	status, _ := strconv.Atoi(c.QueryParam("status"))
	appid, _ := strconv.Atoi(c.QueryParam("appid"))
	admin.ChangeGameServerStatus(int32(appid), int32(status))
	return c.JSON(http.StatusOK, H{"msg": "成功"})
}

// 获取玩家的输赢情况
//http://localhost:7102/getplayerwingold?uid=1
func getPlayerWinGold(c echo.Context) error {
	uid, _ := strconv.Atoi(c.QueryParam("uid"))

	var dh H = make(H, 0)
	dh["uid"] = uid
	dh["WinGold"] = common.GetPlayerWinGold(int64(uid))

	return c.JSON(http.StatusOK, dh)
}

// 设置概率配置
//http://localhost:7102/setprobability?id=1&type=1&data=
func setProbability(c echo.Context) error {

	id, _ := strconv.Atoi(c.QueryParam("id"))
	Type, _ := strconv.Atoi(c.QueryParam("type"))
	data := c.QueryParam("data")

	common.SetRateLevel(id, Type, data)

	return c.JSON(http.StatusOK, H{"msg": "成功"})
}

// 获取概率配置
//http://localhost:7102/getprobabilityinfo?appid=0&apptype=0&sessionid=0&uid=0
func getProbabilityInfo(c echo.Context) error {
	uid, _ := strconv.Atoi(c.QueryParam("uid"))
	appid, _ := strconv.Atoi(c.QueryParam("appid"))
	apptype, _ := strconv.Atoi(c.QueryParam("apptype"))
	sessionid, _ := strconv.Atoi(c.QueryParam("sessionid"))

	var dh H = make(H, 0)

	if uid != 0 {
		Level, Level1, WinloseValue := common.GetPlayerProbabilityInfo(int64(uid))
		dh["Level"] = Level
		dh["Level1"] = Level1
		dh["WinloseValue"] = WinloseValue
	} else {
		playerbankerlevel, playerbankerlevel1, sysbankerlevel, sysbankerlevel1, playerlevel, playerlevel1 := common.GetSysProbabilityInfo(appid, apptype, sessionid)

		dh["appid"] = appid
		dh["apptype"] = apptype
		dh["sessionid"] = sessionid
		dh["playerbankerlevel"] = playerbankerlevel
		dh["playerbankerlevel1"] = playerbankerlevel1
		dh["sysbankerlevel"] = sysbankerlevel
		dh["sysbankerlevel1"] = sysbankerlevel1
		dh["playerlevel"] = playerlevel
		dh["playerlevel1"] = playerlevel1

	}
	return c.JSON(http.StatusOK, dh)
}

// 玩家概率配置
//http://localhost:7102/updateprobabilityinfo?uid=1&level=0&level1=0&playerbankerlevel=0&playerbankerlevel1=0&sysbankerlevel=0&sysbankerlevel1=0&playerlevel=0&playerlevel1=0&winlosevalue=0&type=0&appid=0&apptype=0&sessionid=0
//http://localhost:7102/updateprobabilityinfo?uid=1&lv=0&type=0
func updateProbabilityInfo(c echo.Context) error {

	uid, _ := strconv.Atoi(c.QueryParam("uid"))
	Type, _ := strconv.Atoi(c.QueryParam("type"))
	lv, _ := strconv.Atoi(c.QueryParam("lv"))
	WinloseValue, _ := strconv.Atoi(c.QueryParam("winlosevalue"))

	//Level, _ := strconv.Atoi(c.QueryParam("level"))
	//Level1, _ := strconv.Atoi(c.QueryParam("level1"))
	if uid != 0 {
		//common.SetPlayerLevel(int64(uid), (5*(Level-1))+Level1, Type)
		common.SetPlayerLevel(int64(uid), lv, Type)
		//common.SetPlayerWinloseValue(int64(uid), int64(WinloseValue), Type)

		s := fmt.Sprintf("[%d]", uid)
		common.SetPlayersValue(s, int64(WinloseValue), Type)
	}

	/*Level, _ := strconv.Atoi(c.QueryParam("level"))
	Level1, _ := strconv.Atoi(c.QueryParam("level1"))
	WinloseValue, _ := strconv.Atoi(c.QueryParam("winlosevalue"))

	appid, _ := strconv.Atoi(c.QueryParam("appid"))
	apptype, _ := strconv.Atoi(c.QueryParam("apptype"))
	sessionid, _ := strconv.Atoi(c.QueryParam("sessionid"))
	playerbankerlevel, _ := strconv.Atoi(c.QueryParam("playerbankerlevel"))
	playerbankerlevel1, _ := strconv.Atoi(c.QueryParam("playerbankerlevel1"))
	sysbankerlevel, _ := strconv.Atoi(c.QueryParam("sysbankerlevel"))
	sysbankerlevel1, _ := strconv.Atoi(c.QueryParam("sysbankerlevel1"))
	playerlevel, _ := strconv.Atoi(c.QueryParam("playerlevel"))
	playerlevel1, _ := strconv.Atoi(c.QueryParam("playerlevel1"))

	if uid != 0 {
		common.SetPlayerProbabilityInfo(int64(uid), Level, Level1, int64(WinloseValue), Type)

		common.SetPlayerLevel(int64(uid), lv, Type)
	} else {
		common.SetSysProbabilityInfo(appid, apptype, sessionid, playerbankerlevel, playerbankerlevel1, sysbankerlevel, sysbankerlevel1, playerlevel, playerlevel1, Type)
	}*/
	return c.JSON(http.StatusOK, H{"msg": "成功"})
}


// 获取盈利信息（系统 机器人）
//http://localhost:7102/getrevenueinfo?appid=1&roomtype=0&sessionid=0
func getRevenueInfo(c echo.Context) error {
	appid, _ := strconv.Atoi(c.QueryParam("appid"))
	roomtype, _ := strconv.Atoi(c.QueryParam("roomtype"))
	sessionid, _ := strconv.Atoi(c.QueryParam("sessionid"))

	SysWinGold, RobotWinGold, InputGold, OutputGold := admin.GetGameRevenue().GetRevenueValue(int32(appid), uint32(roomtype), int32(sessionid), true)

	var dh H = make(H, 0)
	dh["SysWinGold"] = SysWinGold
	dh["RobotWinGold"] = RobotWinGold
	dh["InputGold"] = InputGold
	dh["OutputGold"] = OutputGold

	return c.JSON(http.StatusOK, dh)
}


// 推送通知
//http://localhost:7102/notice?noticetype=1&id=0&uid=169087&title=test2&gold=10&giveuid=169087&giveusergold=111
//http://localhost:7102/notice?data={"noticetype":2,"id":0,"uid":0,"giveuid":0,"title":"test","gold":0,"carrygold ":100,"safegold":100,"givegold":0,"givecarrygold ":100,"givesafegold":100}
func notice(c echo.Context) error {
	c.Response().CloseNotify()

	noticeData := &SysNotice{}
	data := []byte(c.QueryParam("data"))

	beego.Info("console notice ", string(data))

	err := json.Unmarshal(data, noticeData)
	if err != nil {
		beego.Error("json.Unmarshal failed ", err, string(data))
		return c.JSON(http.StatusOK, H{"msg": err.Error()})
	}
	if noticeData.NoticeType == 3 {
		fmt.Println("noticeData.NoticeType")
	}
	switch noticeData.NoticeType {
	case 3:
		noticeChargeInGame(noticeData)
	case 2:
	case 1:
	}

	msg := &protocol.SPushMessage{
		Uid:          proto.Int32(int32(noticeData.GiveUid)),
		ID:           proto.Int32(int32(noticeData.Id)),
		Type:         proto.Uint32(uint32(noticeData.NoticeType)),
		Title:        proto.String(noticeData.Title),
		ValidendTime: proto.String(noticeData.ValidendTime),
		Frequency:    proto.Int32(noticeData.Frequency),
		//Gold:      proto.Int64(noticeData.GiveGold),
		//RoomCard:  proto.Int32(int32(noticeData.GiveRoomCard)),
		//CarryGold: proto.Int64(noticeData.GiveCarryGold),
		//SafeGold:  proto.Int64(noticeData.GiveSafeGold),
	}
	if msg.GetType() != 3 {
		msg.Gold = proto.Int64(noticeData.GiveGold)
		msg.RoomCard = proto.Int32(int32(noticeData.GiveRoomCard))
		msg.CarryGold = proto.Int64(noticeData.GiveCarryGold)
		msg.SafeGold = proto.Int64(noticeData.GiveSafeGold)
	}
	go admin.PushMessage(msg)

	// 玩家信息改变(更新赠送玩家的金币)
	if noticeData.Uid != 0 && msg.GetType() == 2 {

		changeUserData := &protocol.SChangeUserData{
			Uid:          proto.Uint32(uint32(noticeData.Uid)),
			Gold:         proto.Int64(int64(noticeData.Gold)),
			RoomCard:     proto.Int32(int32(noticeData.RoomCard)),
			CarryGold:    proto.Int64(noticeData.CarryGold),
			SafeGold:     proto.Int64(noticeData.SafeGold),
			CarryDiamond: proto.Int64(int64(noticeData.CarryDiamond)),
			IsAgent:      proto.Bool(noticeData.IsAgent == 2),
			Diamond: proto.Int64(int64(noticeData.Diamond)),
			SafeDiamond: proto.Int64(int64(noticeData.SafeDiamond)),
		}

		go admin.ChangeUserData(changeUserData)

	}

	return c.JSON(http.StatusOK, H{"msg": "成功"})
}



func noticeChargeInGame(sysNotice *SysNotice) error {
	gsId := common.GetPlayerGameSvrInfo(uint32(sysNotice.Uid))

	if gsId != 0 {  // 如果在游戏中则通知相关游戏
		msg := &protocol.NoticeChargeToGs{}
		msg.Uid = proto.Uint32(uint32(sysNotice.Uid))
		msg.ChargeType = proto.Int32(int32(sysNotice.ChargeType))
		msg.ChargeGold = proto.Int64(sysNotice.ChargeGold)
		msg.ChargeRoomCard = proto.Int64(sysNotice.ChargeRoomCard)
		msg.ChargeDiamond = proto.Int64(sysNotice.ChargeDiamond)

		admin.NoticeToGs(uint32(gsId), uint32(protocol.CMD_NOTICE_CHARGE_TO_GS), msg)
	}

	return nil
}

// 更新java服务器信息
func updatePlatform(c echo.Context) error {
	c.Response().CloseNotify()

	admin.ReportSvrInfo2Platform()

	return c.JSON(http.StatusOK, H{"msg": "成功"})
}

// 更新游戏配置表
func updateGameConfig(c echo.Context) error {
	appid, _ := strconv.Atoi(c.QueryParam("appid"))
	apptype, _ := strconv.Atoi(c.QueryParam("apptype"))
	uid, _ := strconv.Atoi(c.QueryParam("uid"))
	ttype, _ := strconv.Atoi(c.QueryParam("type")) //更新类型 0：基本配置  1：概率配置
	data := c.QueryParam("data")
	go admin.UpdateGameConfig(int32(appid), int32(apptype), int32(ttype), int32(uid), data)
	return c.JSON(http.StatusOK, H{"msg": "成功"})
}

// 踢出玩家
//http://localhost:7102/kickplayer?uid=1&type=0
func kickPlayer(c echo.Context) error {
	uid, _ := strconv.Atoi(c.QueryParam("uid"))
	Type, _ := strconv.Atoi(c.QueryParam("type"))

	go admin.GetMiscSvr().KickPlayer(uint32(uid), int32(Type))

	return c.JSON(http.StatusOK, H{"msg": "成功"})
}