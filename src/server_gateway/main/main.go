package main

import (
	"runtime"
	"flag"
	"lib/public_config"
	"lib/common"
	"server_gateway/conf"
	"github.com/astaxie/beego"
	. "lib/socket"
	.  "server_gateway/gateserver"
	"server_gateway/typeconf"
	"lib/server/adminclient"
)

func initGatewaySvr() {
	Hkgw = &GateWaySvr{Appsession: make(map[typeconf.UserID]*AppSession)}

	netAddr := conf.GwCfg.CfgGatewayServer.AppServerAddress

	if conf.GwCfg.CfgGatewayServer.NetIP == 1 {
		netAddr = common.GetNetAddr() // 外网地址
		port := common.SubString(conf.GwCfg.CfgGatewayServer.AppServerAddress,
			common.UnicodeIndex(conf.GwCfg.CfgGatewayServer.AppServerAddress,":"),
			len(conf.GwCfg.CfgGatewayServer.AppServerAddress))

		netAddr += port
	}

	Hkgw.Svr2CenterSession = &adminclient.Svr2CenterSession{
		Sid: uint32(conf.GwCfg.CfgGatewayServer.ServerID),
		ServiceType_: int32(conf.GwCfg.CfgGatewayServer.ServerType),
		ServerAddr_: netAddr}

	centerServers := public_config.GetCfgCenterServers()
	CreateConnection(centerServers["0"].TcpListenAddr, false, true, Hkgw.Svr2CenterSession)


	InitGw2Svr()

}

func InitLog() {
	common.SetLogOption(conf.GwCfg.CfgLog.LogDir,
		conf.GwCfg.CfgLog.LogFileName,
		conf.GwCfg.CfgLog.LogLevel,
		conf.GwCfg.CfgLog.LogToConsole)

	beego.Info("gatewayServer log-options set done!")
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	var pubConfPath string
	flag.StringVar(&pubConfPath, "public_conf", "./public_config.toml", "public_config path")
	flag.Parse()

	public_config.ParseToml(pubConfPath)
	conf.ReadConf("conf/gateway.toml")
	common.LoadIPData(conf.GwCfg.CfgGatewayServer.IPDataPath)

	InitLog()

	MaintainOutConns()

	initGatewaySvr()

}



