package admin

import (
	"sync"
	"lib/utils"
	"lib/public_config"
	"server_console/conf"
	"errors"
	"strings"
	"lib/socket"
	"protocol"
	"data"
	"github.com/astaxie/beego"
	"code.google.com/p/goprotobuf/proto"
	"strconv"
	"fmt"
	"net/url"
	"lib/http_agent"
	"lib/common"
	"github.com/gomodule/redigo/redis"
	"lib/db"
	"encoding/json"
	"time"
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

func  StartOnlineLog() {
	return

	go func() {
		for range time.NewTicker(time.Duration(600) * time.Second).C {
			centerSvr.Lock()

			for _, its := range centerSvr.gamePlayers {
				for _, v := range its {
					d := data.PfGameOnlineLog{
						Appid:     v.AppID,
						Sessionid: v.GameID,
						GameType:  v.GameType,
						Count:     v.PlayersCount,
					}

					d.Save()
				}
			}

			centerSvr.Unlock()
		}
	}()
}


// 更新java平台在线人数
func (this *CenterSvr) UpdateOnlineUsers() {
	this.RLock()
	defer this.RUnlock()

	var count int32 = 0
	for _, its := range this.gamePlayers {
		for _, it := range its {
			//count += it.PlayersCount
			count += it.RealPlayerCount
		}
	}

	path := fmt.Sprintf("/platform/updateonlineusers?onlinegameusersnum=%d", count)
	go http_agent.AgentRequest(public_config.GetCfgPlateform().PaltformAddr, path)
}

// 给java平台注册中心服务器信息
func (this *CenterSvr) RegistCenterToPlatform() {
	//centerstr := "{\"addr\":\"" + public_config.GetCfgSvrAdmin().HttpServer + "\"}"

	svrInfo := public_config.GetCfgCenterServer(conf.CenterCfg.CfgServer.ServerID)
	if svrInfo == nil {
		return
	}
	//centerstr := public_config.GetCfgSvrAdmin().HttpServer
	centerstr := svrInfo.HttpListenAddr
	path := fmt.Sprintf("/platform/updatecenter?data=%s", centerstr)
	go http_agent.AgentRequest(public_config.GetCfgPlateform().PaltformAddr, path)
}


/*
type Server interface {
	OnCreateConnection(c Conner)
	OnClose()
}

 */

func StartListenSvrs() {
	svrInfo := public_config.GetCfgCenterServer(conf.CenterCfg.CfgServer.ServerID)
	if svrInfo == nil {
		panic(errors.New("no center-serverId configured"))
	}

	addr := ":" + strings.Split(svrInfo.TcpListenAddr, ":")[1]

	go socket.StartTcpListen(centerSvr, addr)
}

func GetMiscSvr() *CenterSvr {
	return centerSvr
}

func (this *CenterSvr) OnCreateConnection(c socket.Conner) {
	session := NewSvrSession(this)
	c.SetSession(session)
	session.connection = c.(*socket.Connection)
}

func (this *CenterSvr) OnClose() {

}

//新建客户端session
func NewSvrSession(svr *CenterSvr) *Center4SvrSession {
	return &Center4SvrSession{server: svr}
}

// 删除注册的服务器
func (this *CenterSvr) RemoveServer(session *Center4SvrSession) {
	this.Lock()
	defer this.Unlock()

	//排除对应得session id
	if v, ok := this.servers[session.service_type]; ok && v != nil {
		items := make([]*Center4SvrSession, 0, len(v))
		for _, item := range v {
			if item.serverID != session.serverID {
				items = append(items, item)
			}
		}
		this.servers[session.service_type] = items
	}

	// 如果是游戏服则清除相应的房间信息
	if session.service_type == uint32(protocol.ServiceType_ServiceType_Game) {
		roomPlayerCount := 0
		// 暂时每款游戏只有一个游戏服
		if its, ok := this.rooms[uint32(session.AppId)]; ok {
			for _, room := range its {
				//if room.ServerID == int32(session.serverID) {
				roomPlayerCount += len(room.Players)

				// 玩家下线
				for _, v := range room.Players {
					data.UserOutGame(v)
				}

				//}
			}
			delete(this.rooms, uint32(session.AppId))
		}

		for _, its := range this.gamePlayers {
			for _, it := range its {
				if it.AppID == session.AppId {
					it.PlayersCount = it.PlayersCount - int32(roomPlayerCount)
					if it.PlayersCount < 0 {
						it.PlayersCount = 0
					}
				}
			}
		}
	}
}

// 注册服务器
func (this *CenterSvr) RegistServer(session *Center4SvrSession) {
	beego.Info("Inside CenterSvr.RegistServer(), svrType:", session.service_type, "svrId:", session.serverID)
	this.Lock()
	defer this.Unlock()

	var isFindServer bool = false
	var isFindSection bool = false

	if v, ok := this.servers[session.service_type]; ok && v != nil {
		for _, item := range v {
			if item.AppId != session.AppId {
				continue
			}
			if item.serverID == session.serverID {
				isFindServer = true
				//break
			}
			if item.serverID == session.serverID || (session.CodeSectionMin <= item.CodeSectionMax && session.CodeSectionMin >= item.CodeSectionMin) || (session.CodeSectionMax <= item.CodeSectionMax && session.CodeSectionMax >= item.CodeSectionMin) {
				isFindSection = true
				//break
			}
		}
	}

	registResult := &protocol.NonCenterSvrRegistResult{
		Result: proto.Int32(0),
	}

	if isFindSection && session.service_type == uint32(protocol.ServiceType_ServiceType_Game) {
		registResult.Result = proto.Int32(2)
		head := &socket.PackHead{Sid: 0, Uid: 0, Cmd: uint32(protocol.CMD_AMPQ_APP_SERVER_REGIST_RESULT)}
		session.Write(registResult, head)
		return
	}

	if !isFindServer {
		this.servers[session.service_type] = append(this.servers[session.service_type], session)
	} else if isFindServer && session.service_type != uint32(protocol.ServiceType_ServiceType_Game) {
		if v, ok := this.servers[session.service_type]; ok && v != nil {
			for _, item := range v {
				if item.serverID == session.serverID {
					item = session
					break
				}
			}
		}
	}

	// 如果是网关
	if session.service_type == uint32(protocol.ServiceType_ServiceType_Gateway) {
		this.NotifyAllServersToSession(session, uint32(protocol.ServiceType_ServiceType_Login), 1)
		this.NotifyAllServersToSession(session, uint32(protocol.ServiceType_ServiceType_Game), 1)
	} else if session.service_type == uint32(protocol.ServiceType_ServiceType_Game) || session.service_type == uint32(protocol.ServiceType_ServiceType_Login) { // 如果是游戏服务器服 或者是 登陆服
		this.NotifyAllServersToSession(session, uint32(protocol.ServiceType_ServiceType_DBA), 1)
	} else if session.service_type == uint32(protocol.ServiceType_ServiceType_DBA) { // 如果是数据服
		this.NotifyServerToAllServiceType(session, 1, uint32(protocol.ServiceType_ServiceType_Game))
		this.NotifyServerToAllServiceType(session, 1, uint32(protocol.ServiceType_ServiceType_Login))
	}

	if session.service_type != uint32(protocol.ServiceType_ServiceType_DBA) && session.service_type != uint32(protocol.ServiceType_ServiceType_Gateway) { // 如果是游戏服务器或者是大厅服务器 则通知网关
		this.NotifyServerToAllServiceType(session, 1, uint32(protocol.ServiceType_ServiceType_Gateway))
	}


	// 注册成功
	head := &socket.PackHead{Sid: 0, Uid: 0, Cmd: uint32(protocol.CMD_AMPQ_APP_SERVER_REGIST_RESULT)}
	session.Write(registResult, head)


	// 如果是网关
	if session.service_type == uint32(protocol.ServiceType_ServiceType_Gateway) {
		this.NoticeGatewayToPlatform(false)
	}
}


// 通知所有服务器状态信息给指定网关
func (this *CenterSvr) NotifyAllServersToSession(session *Center4SvrSession, service_type uint32, state int32) {

	if v, ok := this.servers[service_type]; ok && v != nil {
		for _, item := range v {
			this.NotifyServerToSession(session, item, state)
		}
	}
}

// 通知所有网关服务器状态信息
func (this *CenterSvr) NotifyServerToAllServiceType(server *Center4SvrSession, state int32, serviceType uint32) {

	if v, ok := this.servers[serviceType]; ok && v != nil {
		for _, item := range v {
			this.NotifyServerToSession(item, server, state)
		}
	}
}

// 通知单个网关服务器状态信息
func (this *CenterSvr) NotifyServerToSession(session *Center4SvrSession, server *Center4SvrSession, state int32) {
	serverNoticy := &protocol.NonCenterSvrRegistNotify{
		ServerId:       proto.Int32(int32(server.serverID)),
		ServerType:     proto.Int32(int32(server.service_type)),
		ServerAddress:  proto.String(server.ServerAddress),
		ServerState:    proto.Int32(state),
		AppId:          proto.Int32(int32(server.AppId)),
		CodeSectionMin: proto.Int32(int32(server.CodeSectionMin)),
		CodeSectionMax: proto.Int32(int32(server.CodeSectionMax)),
	}

	head := &socket.PackHead{Sid: 0, Uid: 0, Cmd: uint32(protocol.CMD_AMPQ_APP_SERVER_REGIST_NOTIFY)}
	session.Write(serverNoticy, head)
}


// 给java平台通知网关信息
func (this *CenterSvr) NoticeGatewayToPlatform(lock bool) {

	if lock {
		this.RLock()
		defer this.RUnlock()
	}
	if items, ok := this.servers[uint32(protocol.ServiceType_ServiceType_Gateway)]; ok && items != nil {

		var gateways string

		gateways += "["

		for s, item := range items {

			gateways += "{"

			gateways += "\"addr\":\"" + item.ServerAddress + "\","
			gateways += "\"usercount\":" + strconv.Itoa(int(item.UserCount))

			gateways += "}"

			if s < len(items)-1 {
				gateways += ","
			}

		}

		gateways += "]"

		path := fmt.Sprintf("/platform/updategateway?data=%s", url.QueryEscape(gateways))
		http_agent.AgentRequest(public_config.GetCfgPlateform().PaltformAddr, path)

	} else {
		path := fmt.Sprintf("/platform/updategateway?data=%s", url.QueryEscape("[]"))
		http_agent.AgentRequest(public_config.GetCfgPlateform().PaltformAddr, path)
	}

}


// 更新console房间信息
func (this *CenterSvr) UpdateRoomInfo(roomInfo *protocol.RoomInfoUpdate, AppId int32) {
	this.Lock()
	defer this.Unlock()
	if int32(roomInfo.GetOperate()) == int32(1) { // 更新
		// 更新人数
		curRoomPlayerCount := 0
		var curRealPlayerCount int32 = 0
		isFind := false

		// 更新房间信息
		var players []uint32
		for _, it := range roomInfo.Players {
			players = append(players, it)
		}
		if its, ok := this.rooms[uint32(AppId)]; ok {
			for _, room := range its {
				if room.RoomID == roomInfo.GetRoomId() {
					curRoomPlayerCount = len(room.Players)
					curRealPlayerCount = room.RealPlayerCount
					isFind = true

					room.Players = players
					room.ServerID = roomInfo.GetServerID()
					room.CreateUID = roomInfo.GetCreateUID()
					room.RoomID = roomInfo.GetRoomId()
					room.GameID = roomInfo.GetGameID()
					room.AppID = AppId
					room.GameType = int32(roomInfo.GetRoomType())
					room.RealPlayerCount = roomInfo.GetRealPlayerCount()
					room.ModelKind = roomInfo.GetModelKind()
					room.RoomCreateInfo = roomInfo.GetRoomCreateInfo()
					break
				}
			}
		}
		if !isFind {
			room := &RoomInfo{
				Players:         players,
				ServerID:        int32(roomInfo.GetServerID()),
				CreateUID:       uint32(roomInfo.GetCreateUID()),
				RoomID:          uint32(roomInfo.GetRoomId()),
				AppID:           AppId,
				GameID:          int32(roomInfo.GetGameID()),
				GameType:        int32(roomInfo.GetRoomType()),
				RealPlayerCount: int32(roomInfo.GetRealPlayerCount()),
				ModelKind:       roomInfo.GetModelKind(),
				RoomCreateInfo:  roomInfo.GetRoomCreateInfo(),
			}
			if _, ok := this.rooms[uint32(AppId)]; !ok {
				m := make([]*RoomInfo, 0)
				m = append(m, room)
				this.rooms[uint32(AppId)] = m
			} else {
				this.rooms[uint32(AppId)] = append(this.rooms[uint32(AppId)], room)
			}
		}

		gamePlayer := &GamePlayers{AppID: AppId, GameID: int32(roomInfo.GetGameID()), GameType: int32(roomInfo.GetRoomType()), PlayersCount: int32(len(roomInfo.Players))}
		isFind = false
		if its, ok := this.gamePlayers[AppId]; ok {
			for _, it := range its {
				if it.GameID == int32(roomInfo.GetGameID()) && it.GameType == int32(roomInfo.GetRoomType()) {
					it.PlayersCount = it.PlayersCount + int32(len(roomInfo.Players)-curRoomPlayerCount)
					it.RealPlayerCount = it.RealPlayerCount + int32(roomInfo.GetRealPlayerCount()-curRealPlayerCount)
					if it.PlayersCount < 0 {
						it.PlayersCount = 0
					}
					if it.RealPlayerCount < 0 {
						it.RealPlayerCount = 0
					}
					isFind = true
					break
				}
			}
		}
		if !isFind {
			this.gamePlayers[AppId] = append(this.gamePlayers[AppId], gamePlayer)
		}

	} else if int32(roomInfo.GetOperate()) == int32(2) { // 删除

		if its, ok := this.rooms[uint32(AppId)]; ok {
			rooms := make([]*RoomInfo, 0)
			for _, room := range its {
				if room.RoomID == uint32(roomInfo.GetRoomId()) {

					// 更新人数
					if its, ok1 := this.gamePlayers[AppId]; ok1 {
						for _, it := range its {
							if it.GameID == int32(roomInfo.GetGameID()) && it.GameType == int32(roomInfo.GetRoomType()) {
								it.PlayersCount = it.PlayersCount - int32(len(room.Players))
								it.RealPlayerCount = it.RealPlayerCount - roomInfo.GetRealPlayerCount()
								if it.PlayersCount < 0 {
									it.PlayersCount = 0
								}
								if it.RealPlayerCount < 0 {
									it.RealPlayerCount = 0
								}
								break
							}
						}
					}

					room.Players = make([]uint32, 0)
					room.ServerID = roomInfo.GetServerID()
					room.RoomID = roomInfo.GetRoomId()
					room.GameID = roomInfo.GetGameID()
					room.AppID = AppId
				} else {
					rooms = append(rooms, room)
				}
			}

			this.rooms[uint32(AppId)] = rooms

		}

	}
}


func (this *CenterSvr) DetectCurrentRoom(uid uint32) *RoomInfo {
	this.Lock()
	defer this.Unlock()


	// 改为从redis中获取
	gsInviteCode := common.GetPlayerGameInviteCode(uid)
	if gsInviteCode != 0 {
		return &RoomInfo{
			RoomID: uint32(gsInviteCode),
		}
	}

	return nil
}

// 获取游戏玩家信息
func (this *CenterSvr) GetRoomPlayers(AppId, Type int32) *protocol.SGetRoomPlayers {
	this.RLock()
	defer this.RUnlock()

	sGetRoomPlayers := &protocol.SGetRoomPlayers{}

	if its, ok := this.gamePlayers[AppId]; ok {
		for _, it := range its {
			if it.PlayersCount > 0 && it.GameType == Type {
				sGetRoomPlayers.ID = append(sGetRoomPlayers.ID, it.GameID)
				sGetRoomPlayers.PlayerCount = append(sGetRoomPlayers.PlayerCount, it.PlayersCount)
			}
		}
	}
	return sGetRoomPlayers
}


// 改变玩家数据信息
func (this *CenterSvr) ChangeUserData(msg *protocol.SChangeUserData) {
	this.RLock()
	defer this.RUnlock()

	var index int32 = 0
	var sindex int32 = 0
	var slen int32 = 0
	if items, itok := this.servers[uint32(protocol.ServiceType_ServiceType_Login)]; itok {
		slen = int32(len(items))
	}

	if slen > 0 {
		sindex = utils.RandInt32Section(0, slen)
	}

	// 发送给大厅 再转发给具体玩家
	if items, itok := this.servers[uint32(protocol.ServiceType_ServiceType_Login)]; itok && items != nil {
		for _, item := range items {
			if index == sindex {
				head := &socket.PackHead{Sid: 0, Uid: 0, Cmd: uint32(protocol.CMD_CHANGE_USER_DATA)}
				item.Write(msg, head)
				break
			}
			index++
		}
	}
}


// 获取大厅的总人数
func (this *CenterSvr) GetPlatformUserCount() int32 {
	this.RLock()
	defer this.RUnlock()

	var playerCount int32 = 0
	if items, ok := this.servers[uint32(protocol.ServiceType_ServiceType_Gateway)]; ok && items != nil {
		for _, item := range items {
			//playerCount += int32(item.UserCount)
			playerCount += int32(item.TotalUserCount)
		}
	}
	return playerCount
}

func (this *CenterSvr) UpdateRoomPlayerInfo(info *protocol.PlayerInfoUpdate) {
	this.Lock()
	defer this.Unlock()

	var uid uint32 = info.GetUid()

	if info.GetOperate() == 0 { // 刪除
		var isFind bool = false
		for _, items := range this.rooms {
			for _, room := range items {
				isFind = false
				for _, it := range room.Players {
					if it == uid {
						isFind = true
						break
					}
				}
				if isFind {
					var players []uint32 = make([]uint32, 0)
					for _, it := range room.Players {
						if it != uid {
							players = append(players, it)
						}
					}
					room.Players = players

					for _, its := range this.gamePlayers {
						for _, it1 := range its {
							if it1.AppID == room.AppID && it1.GameType == room.GameType && it1.GameID == room.GameID {
								it1.PlayersCount = it1.PlayersCount - 1
								if it1.PlayersCount < 0 {
									it1.PlayersCount = 0
								}
							}
						}
					}

					//break
				}
			}
			if isFind {
				//break
			}
		}
	}
}

// 获取房间列表
func (this *CenterSvr) GetRoomList(AppId, Type, CurPage int32, ModelKind uint32) *protocol.SGetAppRoomList {
	this.Lock()
	defer this.Unlock()

	return this.GetRoomListByRedis(AppId, Type, CurPage, ModelKind)
}

func (this *CenterSvr) GetRoomListByRedis(AppId, Type, CurPage int32, ModelKind uint32) *protocol.SGetAppRoomList {

	var pageSize int32 = 20 // 每次发20个房间信息
	var startIndex int32 = pageSize * CurPage
	//var endIndex int32 = startIndex + pageSize

	appRoomList := &protocol.SGetAppRoomList{
		AppID:   proto.Int32(AppId),
		Type:    proto.Int32(Type),
		CurPage: proto.Int32(CurPage),
	}

	roomkey := fmt.Sprintf(data.Game_Room_List_Key, AppId)
	var index int32 = 0
	for {

		// 从redis中获取房间列表
		values, _ := redis.Values(db.Do("lrange", roomkey, startIndex, startIndex + pageSize))
		if len(values) == 0 {
			break
		}
		/*var roomkeys string = ""
		for _, v := range values {
			roomkeys += string(v.([]byte)) + ""
		}*/

		// 从redis中获取房间信息
		roominfos, err := redis.Values(db.Do("MGET", values[0:]...))
		if len(roominfos) == 0 {
			break
		}
		if err != nil {
			beego.Error(err.Error())
		}
		for _, v := range roominfos {
			if v == nil {
				continue
			}
			rinfo := &common.GameRoomInfo{}
			if err := json.Unmarshal(v.([]byte), rinfo); err == nil {
				// 如果客户端有选显示的模式ModelKind (所有：0[默认],普通模式：1<<4,疯狂模式：1<<5)
				if ModelKind != 0 {
					var find bool = false
					for i := uint32(2); i <= uint32(32); i++ {
						if (ModelKind&(1<<i) != 0) && (rinfo.ModelKind&(1<<i) != 0) {
							find = true
							break
						}
					}
					if !find {
						continue
					}
				}

				index++
				appRoomList.RoomInfos = append(appRoomList.RoomInfos, &protocol.AppRoomInfo{
					InviteCode:     proto.Uint32(rinfo.RoomId),
					ModelKind:      proto.Uint32(rinfo.ModelKind),
					RoomCreateInfo: proto.String(rinfo.RoomCreateInfo),
					PlayerCount:    proto.Int32(int32(len(rinfo.Players))),
				})
				if index >= pageSize {
					break
				}
			}

		}

		if index >= pageSize {
			break
		}

		startIndex += pageSize
	}

	return appRoomList
}




// 推送消息给所有玩家
func (this *CenterSvr) PushMessageToAllUsers(msg *protocol.SPushMessage) {
	this.RLock()
	defer this.RUnlock()

	// 发送给网关 由网关转发给玩家
	if items, itok := this.servers[uint32(protocol.ServiceType_ServiceType_Gateway)]; itok && items != nil {
		for _, item := range items {
			head := &socket.PackHead{Sid: 0, Uid: 0, Cmd: uint32(protocol.CMD_PUSH_MESSAGE)}
			item.Write(msg, head)
		}
	}
}

// 推送消息给指定玩家
func (this *CenterSvr) PushMessageToUser(msg *protocol.SPushMessage) {
	this.RLock()
	defer this.RUnlock()

	var index int32 = 0
	var sindex int32 = 0
	var slen int32 = 0
	if items, itok := this.servers[uint32(protocol.ServiceType_ServiceType_Login)]; itok {
		slen = int32(len(items))
	}

	if slen > 0 {
		sindex = utils.RandInt32Section(0, slen)
	}

	// 发送给大厅 再转发给具体玩家
	if items, itok := this.servers[uint32(protocol.ServiceType_ServiceType_Login)]; itok && items != nil {
		for _, item := range items {
			if sindex == index {
				head := &socket.PackHead{Sid: 0, Uid: 0, Cmd: uint32(protocol.CMD_PUSH_MESSAGE)}
				item.Write(msg, head)
				break //
			}
			index++
		}
	}
}

func (this *CenterSvr) GetSvrSession(svrType, svrId uint32) *Center4SvrSession {
	this.RLock()
	defer this.RUnlock()

	if v, ok := this.servers[svrType]; ok && v != nil {
		for _, item := range v {
			if svrId == item.serverID {
				return item
			}
		}
	}

	return nil
}

// 更新游戏配置信息
func (this *CenterSvr) UpdateGameConfig(appID, appType, tType, uid int32, data string) {
	this.RLock()
	defer this.RUnlock()

	updateGameConfig := &protocol.UpdateGameConfig{
		AppId:   proto.Int32(appID),
		AppType: proto.Int32(appType),
		Type:    proto.Int32(tType), //更新类型 0：基本配置  1：概率配置  2：指定玩家概率  3: 牌型赔率
		Uid:     proto.Uint32(uint32(uid)),
		Data:    proto.String(data),
	}

	if v, ok := this.servers[uint32(protocol.ServiceType_ServiceType_Game)]; ok && v != nil {
		for _, item := range v {
			if item.AppId != appID {
				continue
			}
			head := &socket.PackHead{Sid: 0, Uid: 0, Cmd: uint32(protocol.CMD_AMPQ_UPDATE_GAME_CONFIG)}
			item.Write(updateGameConfig, head)
		}
	}
}

// 修改游戏服务器状态
func (this *CenterSvr) ChangeGameServerStatus(appid, status int32) {
	this.RLock()
	defer this.RUnlock()

	if appid == 0 && status == 0 {
		common.GetGameManger().StopGame()
	} else if appid == 0 {
		common.GetGameManger().RunGame()
	}

	msg := &protocol.ChangeGameServer{
		Status: proto.Int32(status),
	}

	if v, ok := this.servers[uint32(protocol.ServiceType_ServiceType_Game)]; ok && v != nil {
		for _, item := range v {
			if appid != 0 && item.AppId != appid {
				continue
			}
			head := &socket.PackHead{Sid: 0, Uid: 0, Cmd: uint32(protocol.CMD_CHANGE_GAME_SERVER_STATUS)}
			item.Write(msg, head)
		}
	}
}


// 踢出玩家
func (this *CenterSvr) KickPlayer(uid uint32, Type int32) {
	this.RLock()
	defer this.RUnlock()

	if items, ok := this.servers[uint32(protocol.ServiceType_ServiceType_Gateway)]; ok && items != nil {

		msg := &protocol.SKickPlayerOffLine{
			Uid:  proto.Uint32(uid),
			Type: proto.Int32(Type),
		}

		head := &socket.PackHead{Sid: 0, Uid: uid, Cmd: uint32(protocol.CMD_KICK_PLAYER_OFFLINE)}

		for _, item := range items {

			item.Write(msg, head)

		}
	}

}

// 刷新游戏配置文件
func (this *CenterSvr) RefGameConfig(Type int32) {
	this.RLock()
	defer this.RUnlock()

	msg := &protocol.RefGameConfig{
		Type: proto.Int32(Type),
	}

	// 游戏服务器
	if v, ok := this.servers[uint32(protocol.ServiceType_ServiceType_Game)]; ok && v != nil {
		for _, item := range v {
			head := &socket.PackHead{Sid: 0, Uid: 0, Cmd: uint32(protocol.CMD_REF_GAME_CONFIG)}
			item.Write(msg, head)
		}
	}

	// 数据库统计服务器
	if v, ok := this.servers[uint32(protocol.ServiceType_ServiceType_DBACensus)]; ok && v != nil {
		for _, item := range v {
			head := &socket.PackHead{Sid: 0, Uid: 0, Cmd: uint32(protocol.CMD_REF_GAME_CONFIG)}
			item.Write(msg, head)
		}
	}
}


func StartReportSvrInfo2Platform() {
	go func() {
		for range time.NewTicker(time.Duration(60) * time.Second).C {
			ReportSvrInfo2Platform()
		}
	}()
}

func ReportSvrInfo2Platform() {
	centerSvr.NoticeGatewayToPlatform(true)
	centerSvr.RegistCenterToPlatform()
	centerSvr.UpdateOnlineUsers()
}

// 修改玩家数据
func ChangeUserData(msg *protocol.SChangeUserData) {
	centerSvr.ChangeUserData(msg)
}

// 通知gs游戏内充值
func NoticeToGs(svrId uint32, cmd uint32, msg proto.Message) {
	session := centerSvr.GetSvrSession(uint32(protocol.ServiceType_ServiceType_Game), svrId)

	if session == nil {
		beego.Error("GetSvrSession failed, svrId:", svrId)
		return
	}

	head := &socket.PackHead{Sid: 0, Uid: 0, Cmd: cmd}
	session.Write(msg, head)
}

// java平台推送消息
func PushMessage(msg *protocol.SPushMessage) {
	if msg.GetType() == 1 || msg.GetType() == 3 { // 系统消息
		centerSvr.PushMessageToAllUsers(msg)
	} else if msg.GetType() == 2 && msg.GetUid() != 0 { // 玩家消息
		centerSvr.PushMessageToUser(msg)
	}
}

// 更新游戏配置表
func UpdateGameConfig(appID, appType, tType, uid int32, data string) {
	centerSvr.UpdateGameConfig(appID, appType, tType, uid, data)
}

// 停止游戏
func ChangeGameServerStatus(appid, status int32) {
	centerSvr.ChangeGameServerStatus(appid, status)
}

// 停止游戏
func RefGameConfig(Type int32) {
	centerSvr.RefGameConfig(Type)
}