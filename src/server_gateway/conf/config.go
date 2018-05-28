package conf

import (
	"lib/common"
	"github.com/BurntSushi/toml"
)

type CfgGatewayServer struct {
	AppServerAddress   string //网关，APP端地址
	SvrServerAddress   string //网关，Svr端地址
	ServerID           int    //网关的ID
	ServerType         int    //服务类型
	NetIP              int    //是否上报外网IP
	IPDataPath         string
}

type GatewayConfig struct {
	CfgLog           common.CfgLog
	CfgGatewayServer CfgGatewayServer
}

var (
	GwCfg GatewayConfig
)

func ReadConf(file string) {
	_, err := toml.DecodeFile(file, &GwCfg)
	if err != nil {
		panic(err)
	}
}