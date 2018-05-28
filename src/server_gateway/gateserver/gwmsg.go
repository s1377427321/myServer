package gateserver

import (
	"lib/route"
	"protocol"
)

func RegisterMsgHandler() {
	route.RegisterFunc(protocol.CMD_AMPQ_APP_SERVER_REGIST_NOTIFY, appServerRegistNotify, protocol.NonCenterSvrRegistNotify{})

	route.RegisterFunc(protocol.CMD_DETECT_CURRENT_ROOM, detectCurrentRoom, protocol.SDetectCurrentRoom{})
	route.RegisterFunc(protocol.CMD_GET_ROOM_PLAYERS, getRoomPlayers, protocol.SGetRoomPlayers{})
	route.RegisterFunc(protocol.CMD_PUSH_MESSAGE, pushMessage, protocol.SPushMessage{})
	route.RegisterFunc(protocol.CMD_GET__PLATFORM_PLAYER_COUNT, getPlatformPlayers, protocol.SGetPlatformPlayers{})
	route.RegisterFunc(protocol.CMD_KICK_PLAYER_OFFLINE, kickPlayerOffLine, protocol.SKickPlayerOffLine{})
	route.RegisterFunc(protocol.CMD_GET_App_ROOMLIST, getAppRoomList, protocol.SGetAppRoomList{})

}