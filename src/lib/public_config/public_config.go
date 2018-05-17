package public_config

import (
	"os"
	"github.com/astaxie/beego"
	"github.com/BurntSushi/toml"
)

// redis 配置信息
type CfgRedis struct {
	Address  string //数据库连接地址。
	Password string //数据库连接地址。
}

// 中心服务器信息
type CfgCenterServer struct {
	ServerID       uint32
	TcpListenAddr  string
	HttpListenAddr string
}

type CfgPlateform struct {
	SignKey                    string
	PaltformAddr               string
}

type CfgCommon struct {
	HttpTimeOut                int64  //HTTP请求超时时间 time.Millisecond
	MaxCarryGold               int64  //携带的最大金币数
	MaxCarryDiamond            int64  //携带的最大钻石数
	PersonalRandHeadAddr       string //随机头像地址
	PersonalRandHeadMinSection int32  //随机头像区间
	PersonalRandHeadMaxSection int32
}

type PublicConfig struct{
	CfgRedis              CfgRedis
	CfgCommon             CfgCommon
	CfgPlateform          CfgPlateform
	CfgCenterServer 	  map[string]CfgCenterServer
	
}

var opt *PublicConfig

func ParseToml(file string) error  {
	if _,err := os.Stat(file);os.IsNotExist(err){
		beego.Error(err)
		return nil
	}

	opt = &PublicConfig{}

	_,err :=toml.DecodeFile(file,opt)
	if err != nil{
		beego.Error(err)
		return err
	}

	return  nil
}