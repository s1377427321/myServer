package dbagent

import (
	"reflect"
	"sync"
	"github.com/astaxie/beego"
	"lib/utils"
	"protocol"
	"time"
	"lib/socket"
	"lib/route"
	"lib/server"
	"code.google.com/p/goprotobuf/proto"
)

type DBAgentFunc struct {
	fun reflect.Value
}

type DBAgent struct {
	sync.RWMutex
	funcs       map[int64]*DBAgentFunc
	connections map[int32]*DBASession
	eventID     int64
}

var dbAgent *DBAgent

//启动 tcp client
func DBATCPClient(se socket.ConnSession, addr string, timeout time.Duration) (bool, *socket.Connection) {
	s := socket.CreateConnection(addr, true, false, se)
	//s.SetSession(se)
	//s.Connect(timeout)
	return true, s
}

func GetDBAgent() *DBAgent {
	return dbAgent
}

func InitDbAgent() {
	dbAgent = &DBAgent{
		funcs:       make(map[int64]*DBAgentFunc, 0),
		connections: make(map[int32]*DBASession, 0),
		eventID:     100000,
	}

	route.RegisterFunc(protocol.CMD_AMPQ_APP_SERVER_REGIST_NOTIFY, appServerRegistNotify, protocol.NonCenterSvrRegistNotify{})
}

func appServerRegistNotify(session *server.TmpSession, req *protocol.NonCenterSvrRegistNotify) (proto.Message, error) {
	if req.GetServerType() == int32(protocol.ServiceType_ServiceType_DBA) {
		if int32(req.GetServerState()) == 1 { // 注册服务
			dbAgent.RegistServer(req)
		} else {
			dbAgent.RemoveServer(int32(req.GetServerId()))
		}
	}
	return nil, nil
}


func (this *DBAgent) Exec(Type int32, sqlStr string, f interface{}) int64 {
	this.Lock()
	defer this.Unlock()

	beego.Debug("DBAgent Exec", Type, sqlStr)

	this.eventID++

	var index int32 = 0
	var sindex int32 = 0
	var slen int32 = int32(len(this.connections))

	if slen > 0 {
		sindex = utils.RandInt32Section(0, slen)
	}

	for _, item := range this.connections {
		if sindex == index {
			item.Exec(this.eventID, Type, sqlStr)
			break
		}
		index++
	}

	if Type == int32(protocol.DBASqlType_Select) {
		this.funcs[this.eventID] = &DBAgentFunc{
			fun: reflect.ValueOf(f),
		}
	}

	return this.eventID
}

func (this *DBAgent) RegistServer(req *protocol.NonCenterSvrRegistNotify) {

	//this.Lock()
	//defer this.Unlock()
	TmpSession := &DBASession{
		dbagent:       this,
		ServerId:      int32(req.GetServerId()),
		ServerAddress: string(req.GetServerAddress()),
	}
	if ok, _ := DBATCPClient(TmpSession, string(req.GetServerAddress()), time.Second*10); ok {
		this.connections[int32(req.GetServerId())] = TmpSession
	}
}

func (this *DBAgent) RemoveServer(serverID int32) {

	this.Lock()
	defer this.Unlock()

	delete(this.connections, serverID)
}

func (this *DBAgent) AddSession(session *DBASession) {

	this.Lock()
	defer this.Unlock()

	this.connections[session.ServerId] = session
}

func (this *DBAgent) RemoveSession(session *DBASession) {

	this.Lock()
	defer this.Unlock()

	delete(this.connections, session.ServerId)
}

func (this *DBAgent) HandleMsg(ack *protocol.DBAServerAckMsg) {

	this.Lock()
	defer this.Unlock()

	hd := this.funcs[ack.GetEventID()]
	delete(this.funcs, ack.GetEventID())
	if hd != nil {
		hd.fun.Call([]reflect.Value{reflect.ValueOf(ack)})
	}
}