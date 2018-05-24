package common



type GameRoomInfo struct {
	RoomId           uint32
	CreateUID        uint32
	RoomType         uint32
	ServerID         int32
	GameID           int32
	Players          []uint32
	Operate          int32
	ModelKind        uint32
	RoomCreateInfo   string
	RealPlayerCount  int32
}