package data

const (
	// 玩家信息
	Game_User_Info_Key                  string = "Game_User_Info_%d"                // 玩家信息
	Game_User_Info_Data_Key             string = "Game_User_Info_Data"              // 玩家信息(Data)
	Game_User_Info_Uid_Key              string = "Game_User_Info_Uid"               // 玩家信息(uid)
	Game_User_Info_Nickname_Key         string = "Game_User_Info_Nickname"          // 玩家信息(Nickname)
	Game_User_Info_Sex_Key              string = "Game_User_Info_Sex"               // 玩家信息(Sex)
	Game_User_Info_Photo_Key            string = "Game_User_Info_Photo"             // 玩家信息(Photo)
	Game_User_Info_RoomCardNum_Key      string = "Game_User_Info_RoomCardNum"       // 玩家信息(RoomCardNum)
	Game_User_Info_Gold_Key             string = "Game_User_Info_Gold"              // 玩家信息(Gold)
	Game_User_Info_DeviceId_Key         string = "Game_User_Info_DeviceId"          // 玩家信息(DeviceId)
	Game_User_Info_CarryGold_Key        string = "Game_User_Info_CarryGold"         // 玩家信息(CarryGold)
	Game_User_Info_SafeGold_Key         string = "Game_User_Info_SafeGold"          // 玩家信息(SafeGold)
	Game_User_Info_CarryDiamond_Key     string = "Game_User_Info_CarryDiamond"      // 玩家信息(CarryDiamond)
	Game_User_Info_MatchCount_Key       string = "Game_User_Info_MatchCount"        // 玩家信息(MatchCount)
	Game_User_Info_MatchUpdateTime_Key  string = "Game_User_Info_MatchUpdateTime"   // 玩家信息(MatchUpdateTime)
	Game_User_Info_LoginTypes_Key       string = "Game_User_Info_LoginTypes"        // 玩家信息(LoginTypes)
	Game_User_Info_Promotionid_Key      string = "Game_User_Info_Promotionid"       // 玩家信息(Promotionid)
	Game_User_Info_Promotionchannel_Key string = "Game_User_Info_Promotionchannel"  // 玩家信息(Promotionchannel)
	Game_User_Info_GameSvrID_Key        string = "Game_User_Info_GameSvrID"         // 玩家信息(GameSvrID)
	Game_User_Info_GameInviteCode_Key   string = "Game_User_Info_GameInviteCode"    // 玩家信息(GameInviteCode)


	// 房间信息
	Game_Room_List_Key                  string = "Game_Room_List_%d"                 // 房间列表(Game_Room_List_[appid])
	Game_Room_Info_Key                  string = "Game_Room_Info_%d_%d"              // 房间信息(Game_Room_Info_[appid]_[InviteCode])
	Game_Room_Info_Data_Key             string = "Game_Room_Info_Data"               // 房间信息(Data)


	Login_User_Key                  string = "Game_Login_User_%d" // 玩家等级
	Game_Key                        string = "Game_%d_%d_%d"
	Game_User_Probability_Level_Key string = "Game_User_Probability_Level_"
	Probability_Level_Key           string = "probability_level"  // 概率等级
	Probability_Level1_Key          string = "probability_level1" // 概率等级1
	WinloseValue_Key                string = "WinloseValue"       // 输赢值
	WinGold_Key                     string = "WinGold"            // 输赢金币

	Sys_Prize_Pool_Key string = "sysprizepool_%d_%d_%d" // 系统奖池  appid_type_sessionid

	Probability_PlayerBankerLevel_Key  string = "probability_playerbankerlevel"  // 玩家庄家概率等级
	Probability_PlayerBankerLevel1_Key string = "probability_playerbankerlevel1" // 玩家庄家概率等级1
	Probability_SysBankerLevel_Key     string = "probability_sysbankerlevel"     // 系统庄家概率等级
	Probability_SysBankerLevel1_Key    string = "probability_sysbankerlevel1"    // 系统庄家概率等级1
	Probability_PlayerLevel_Key        string = "probability_playerlevel"        // 闲家概率等级
	Probability_PlayerLevel1_Key       string = "probability_playerlevel1"       // 闲家概率等级1
)
