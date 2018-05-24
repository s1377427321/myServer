package admin

import (
	"lib/socket"
	"lib/route"
	"reflect"
	"github.com/golang/protobuf/proto"
	"github.com/astaxie/beego"
	"protocol"
	"lib/common"
)

//客户端连接
type Center4SvrSession struct {
	connection    *socket.Connection //链接
	remoteAddr    string             //客户端地址
	service_type  uint32             //客户端类型
	serverID      uint32             //客户端id
	server        *CenterSvr       //admin服务器
	ServerAddress string             //注册服务器ip端口
	AppId         int32

	UserCount      uint32 //玩家人数
	TotalUserCount uint32 //总玩家人数

	//游戏中房间号生成的ID范围
	CodeSectionMin int32
	CodeSectionMax int32
}


/*

type ConnSession interface {
	OnInit(*Connection)      //初始化操作，比如心跳的设置...
	OnProcessPack(*PackHead) //处理消息
	OnClose()                //消除所有对Sessioner的引用,心跳...
	Write(msg interface{}, ph *PackHead) error
	Close()
	OnConnect(isOk bool)
}
 */



//初始化
func (this *Center4SvrSession) OnInit(c *socket.Connection) {
	this.connection = c
}

func (this *Center4SvrSession) OnProcessPack(ph *socket.PackHead) {
	if route.Exist(ph.Cmd) {
		hd := route.GetFunc(ph.Cmd)
		t := route.GetProto(ph.Cmd)
		v := reflect.New(t)
		if err := proto.Unmarshal(ph.Body, v.Interface().(proto.Message)); err == nil {
			result := hd.Call([]reflect.Value{v})
			if result[1].IsNil() && !result[0].IsNil() {
				this.connection.Write(ph, result[0].Interface())
			}
		}else {
			beego.Error("protocol  unmarshal fail: ", err)
		}
	}

	if ph.Cmd == uint32(protocol.CMD_AMPQ_APP_SERVER_REGIST) { //注册服务器信息
		regist := &protocol.NonCenterSvrRegist{}
		err := proto.Unmarshal(ph.Body, regist)
		if err != nil {
			beego.Error("unmarshaling error: ", err)
			//beego.Error("服务注册失败")
		} else {

			this.service_type = uint32(regist.GetServerType())
			this.serverID = uint32(regist.GetServerId())
			this.ServerAddress = regist.GetServerAddress()
			this.AppId = regist.GetAppId()
			this.CodeSectionMin = regist.GetCodeSectionMin()
			this.CodeSectionMax = regist.GetCodeSectionMax()

			this.server.RegistServer(this)
		}
	} else if ph.Cmd == uint32(protocol.CMD_AMPQ_ROOM_INFO_NOTIFY) { //房间信息
		roomInfoUpdate := &protocol.RoomInfoUpdate{}
		err := proto.Unmarshal(ph.Body, roomInfoUpdate)
		if err != nil {
			beego.Error("unmarshaling error: ", err)
			//beego.Error("房间信息失败")
		} else {
			roomInfoUpdate.ServerID = proto.Int32(int32(this.serverID))
			this.server.UpdateRoomInfo(roomInfoUpdate, this.AppId)
		}
	} else if ph.Cmd == uint32(protocol.CMD_DETECT_CURRENT_ROOM) { //查找玩家当前房间
		detectCurrentRoom := &protocol.CDetectCurrentRoom{}
		err := proto.Unmarshal(ph.Body, detectCurrentRoom)
		if err != nil {
			beego.Error("unmarshaling error: ", err)
			//beego.Error("查找房间信息失败")
		} else {
			curRoom := this.server.DetectCurrentRoom(detectCurrentRoom.GetUid())
			DetectCurrentRoom := &protocol.SDetectCurrentRoom{
				InviteCode: proto.Uint32(0),
				Uid:        detectCurrentRoom.Uid,
				IsServer:   detectCurrentRoom.IsServer,
			}

			if common.GetGameManger().IsStop {
				DetectCurrentRoom.Error = proto.Uint32(uint32(protocol.StatusCode_STOP_GAME_SERVERING))
			}

			if curRoom != nil {
				DetectCurrentRoom.InviteCode = proto.Uint32(uint32(curRoom.RoomID))
				DetectCurrentRoom.ServerID = proto.Int32(int32(curRoom.ServerID))
				DetectCurrentRoom.AppID = proto.Int32(int32(curRoom.AppID))
				DetectCurrentRoom.Type = proto.Int32(int32(curRoom.GameType))
			}
			DetectCurrentRoom.Uid = proto.Uint32(ph.Uid)
			head := &socket.PackHead{Sid: ph.Sid, Uid: ph.Uid, Cmd: uint32(protocol.CMD_DETECT_CURRENT_ROOM)}
			this.Write(DetectCurrentRoom, head)
		}
	} else if ph.Cmd == uint32(protocol.CMD_GET_ROOM_PLAYERS) { //查找房间信息
		getRoomPlayers := &protocol.CGetRoomPlayers{}
		err := proto.Unmarshal(ph.Body, getRoomPlayers)
		if err != nil {
			beego.Error("unmarshaling error: ", err)
			//beego.Error("查找房间信息失败")
		} else {
			sGetRoomPlayers := this.server.GetRoomPlayers(getRoomPlayers.GetAppID(), getRoomPlayers.GetType())
			sGetRoomPlayers.Uid = proto.Uint32(ph.Uid)

			head := &socket.PackHead{Sid: ph.Sid, Uid: ph.Uid, Cmd: uint32(protocol.CMD_GET_ROOM_PLAYERS)}
			this.Write(sGetRoomPlayers, head)
		}

	}else if ph.Cmd == uint32(protocol.CMD_CHANGE_USER_DATA) { // 玩家数据改变
		changeUserData := &protocol.SChangeUserData{}
		err := proto.Unmarshal(ph.Body, changeUserData)
		if err != nil {
			beego.Error("unmarshaling error: ", err)
			//beego.Error("玩家数据改变失败")
		} else {
			this.server.ChangeUserData(changeUserData)
		}

	}else if ph.Cmd == uint32(protocol.CMD_AMPQ_APP_UPDATE_SERVER_INFO) { // 更新服务器信息
		appUpdateServerInfo := &protocol.AppUpdateServerInfo{}
		err := proto.Unmarshal(ph.Body, appUpdateServerInfo)
		if err != nil {
			beego.Error("unmarshaling error: ", err)
			//beego.Error("更新服务器信息失败")
		} else {
			this.UserCount = appUpdateServerInfo.GetUserCount()
			this.TotalUserCount = appUpdateServerInfo.GetUserTotalCount()
		}
	}else if ph.Cmd == uint32(protocol.CMD_GET__PLATFORM_PLAYER_COUNT) { // 获取大厅人数
		count := this.server.GetPlatformUserCount()

		msg := &protocol.SGetPlatformPlayers{
			PlayerCount: proto.Int32(count),
			Uid:         proto.Uint32(ph.Uid),
		}

		head := &socket.PackHead{Sid: ph.Sid, Uid: ph.Uid, Cmd: uint32(protocol.CMD_GET__PLATFORM_PLAYER_COUNT)}
		this.Write(msg, head)
	} else if ph.Cmd == uint32(protocol.CMD_UPDATE_ROOM_PLAYER_INFO) { // 更新游戏中玩家信息
		playerInfoUpdate := &protocol.PlayerInfoUpdate{}
		err := proto.Unmarshal(ph.Body, playerInfoUpdate)
		if err != nil {
			beego.Error("unmarshaling error: ", err)
		} else {
			this.server.UpdateRoomPlayerInfo(playerInfoUpdate)
		}
	}else if ph.Cmd == uint32(protocol.CMD_GET_App_ROOMLIST) { // 获取房间列表
		getAppRoomList := &protocol.CGetAppRoomList{}
		err := proto.Unmarshal(ph.Body, getAppRoomList)
		if err != nil {
			beego.Error("unmarshaling error: ", err)
		} else {
			msg := this.server.GetRoomList(getAppRoomList.GetAppID(), getAppRoomList.GetType(), getAppRoomList.GetCurPage(), getAppRoomList.GetModelKind())
			head := &socket.PackHead{Sid: ph.Sid, Uid: ph.Uid, Cmd: uint32(protocol.CMD_GET_App_ROOMLIST)}
			this.Write(msg, head)
		}
	}else  {

	}

}

func (this *Center4SvrSession) OnClose() {
	this.server.RemoveServer(this)
}

func (this *Center4SvrSession) Write(msg interface{}, ph *socket.PackHead) error {
	go this.connection.Write(ph, msg)
	return nil
}

func (this *Center4SvrSession) OnConnect(isOk bool) {

}

func (this *Center4SvrSession) Close() {

}









