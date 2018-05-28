package conf

import (
	"github.com/BurntSushi/toml"
	"lib/common"
)

type CfgServer struct {
	ServerID uint32
}

type CenterConfig struct {
	CfgLog common.CfgLog
	CfgServer CfgServer
}

var (
	CenterCfg CenterConfig
)

func ReadConf(file string) {
	_, err := toml.DecodeFile(file, &CenterCfg)
	if err != nil {
		panic(err)
	}
}
