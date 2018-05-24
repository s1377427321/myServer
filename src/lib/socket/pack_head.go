package socket

import (
	"github.com/golang/protobuf/proto"
	"code.google.com/p/go.tools/go/types"
	"errors"
	"bytes"
	"encoding/binary"
)

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

func SerializePackHead(ph *PackHead) (result int, buf []byte) {
	ph.Length = uint32(len(ph.Body))
	bytebuf := new(bytes.Buffer)
	binary.Write(bytebuf, binary.BigEndian, PackFlag)
	binary.Write(bytebuf, binary.BigEndian, ph.SequenceID)
	binary.Write(bytebuf, binary.BigEndian, ph.Cmd)
	binary.Write(bytebuf, binary.BigEndian, ph.Uid)
	binary.Write(bytebuf, binary.BigEndian, ph.Sid)
	binary.Write(bytebuf, binary.BigEndian, ph.Length)
	binary.Write(bytebuf, binary.BigEndian, ph.Reserve)
	buf = make([]byte, PackHeadLength+ph.Length)
	copy(buf, bytebuf.Bytes())
	copy(buf[PackHeadLength:], ph.Body)
	result = len(buf)
	return
}

func SerializePackWithPB(ph *PackHead, msg interface{}, max_size int /*, head_buff, data_buff []byte*/) (length int, data []byte, err error) {
	switch v:=msg.(type) {
	case []byte:
		ph.Body = v
		ph.Length = uint32(len(ph.Body))
		length, data = SerializePackHead(ph)
		return
	case proto.Message:
		if ph.Body, err = proto.Marshal(v); err == nil {
			ph.Length = uint32(len(ph.Body))
			length, data = SerializePackHead(ph)
			return
		} else {
			//beego.Error("proto marshal cmd: %d sid: %d uid: %d error: %v",
			//	ph.Cmd, ph.Sid, ph.Uid, err)
			return 0, nil, err
		}
	case types.Nil:
		ph.Body = nil
		ph.Length = 0
		length, data = SerializePackHead(ph)
		return
	default:
		return 0, nil, errors.New("tcp_conn: error msg type")

	}
	return 0, nil, errors.New("unkonw error while serialize proto")
}

