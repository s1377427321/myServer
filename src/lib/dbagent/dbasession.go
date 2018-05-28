package dbagent

import (
	"lib/socket"
	"protocol"
	"code.google.com/p/goprotobuf/proto"
)




type DBASession struct {
	connection    *socket.Connection
	ServerId      int32
	ServerAddress string
	dbagent       *DBAgent
}


func (this *DBASession) OnInit(c *socket.Connection) {
	this.connection = c
}

func (this *DBASession) OnProcessPack(ph *socket.PackHead) {
	if ph.Cmd == uint32(protocol.CMD_DBA_REQ_MSG) && dbAgent != nil {

		dbaAckMsg := &protocol.DBAServerAckMsg{}
		err := proto.Unmarshal(ph.Body, dbaAckMsg)
		if err == nil {
			dbAgent.HandleMsg(dbaAckMsg)
		}
	}
}

func (this *DBASession) OnConnect(isOk bool) {
	if dbAgent != nil {
		dbAgent.AddSession(this)
	}
}

func (this *DBASession) OnClose() {
	if dbAgent != nil {
		dbAgent.RemoveSession(this)
	}
}

func (this *DBASession) Close() {
	this.connection.SetIsConn(false)
	this.connection.Stop(true)
}

func (this *DBASession) Write(msg interface{}, head *socket.PackHead) error {
	return this.connection.Write(head, msg)
}

func (this *DBASession) Exec(eventID int64, Type int32, sqlStr string) bool {

	msg := &protocol.DBAServerReqMsg{
		EventID: proto.Int64(eventID),
		Type:    proto.Int32(Type),
		SqlStr:  proto.String(sqlStr),
	}

	head := &socket.PackHead{Sid: 0, Uid: 0, Cmd: uint32(protocol.CMD_DBA_REQ_MSG)}
	this.connection.Write(head, msg)

	return true
}