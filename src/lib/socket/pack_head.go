package socket


//抓包规则：tcp[0:n32]==0x13473&&tcp[8:n32]==0x7d5   (0x7d5 是16进制协议号)
const(
	PackFlag uint32 = 0x13473 // 包头标识
	PackHeadLength = 32
	MaxReadBufLen       = 1024 * 1024 * 4
)

const (
	FlagOffset = 0
	SeqOffset = 4
	CmdOffset = 8
	UidOffset = 12
	SidOffset = 16
	LenOffset = 20
	ReserveOffset = 24
)


type PackHead struct {
	PackFlag   uint32
	SequenceID uint32
	Cmd        uint32 //cmd enum, 命令字枚举
	Uid        uint32
	Sid        uint32 ////服务器ID//[16,20 )
	Length     uint32 //[20, 24)
	Reserve    uint64 //[24,32)
	Body       []byte //[32, 32+ length)
}


