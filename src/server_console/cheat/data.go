package cheat


// 系统消息
type SysNotice struct {
	NoticeType   int    `json:"noticetype"` // 1：系统消息  2：用户消息 3：游戏内充值
	Id           int    `json:"id"`
	Uid          int    `json:"uid"`          // 玩家id
	Title        string `json:"title"`        // 标题
	Body         string `json:"body"`         // 内容
	RoomCard     int    `json:"roomcard"`     // 房卡
	Gold         int64  `json:"gold"`         // 金币
	CarryGold    int64  `json:"carrygold"`    // 携带金币
	SafeGold     int64  `json:"safegold"`     // 保险柜金币
	CarryDiamond int64   `json:"carrydiamond"` // 携带钻石
	Diamond      int64   `json:"diamond"`      // 钻石
	SafeDiamond  int64   `json:"safediamond"`   // 保险柜钻石
	GiveUid       int    `json:"giveuid"`       // 赠送玩家ID
	GiveRoomCard  int    `json:"giveroomcard"`  // 房卡
	GiveGold      int64  `json:"givegold"`      // 赠送玩家现有金币
	GiveCarryGold int64  `json:"givecarrygold"` // 赠送玩家携带金币
	GiveSafeGold  int64  `json:"givesafegold"`  // 赠送玩家保险柜金币
	ValidendTime  string `json:"validendtime"`  // 失效时间
	Frequency     int32  `json:"frequency"`     // 频率
	IsAgent       int32  `json:"isagent"`       // 是否是代理

	ChargeType int `json:chargetype`		// 充值类型，金币/房卡/钻石
	ChargeGold int64	`json:"chargegold"`		// 充值金币
	ChargeRoomCard int64 `json:"chargeroomcard"`		// 充值房卡
	ChargeDiamond int64 `json:"chargediamond"`		// 充值钻石
}