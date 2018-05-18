package socket


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


