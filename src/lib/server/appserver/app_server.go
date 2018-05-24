package appserver

import (
	"lib/utils"
	"sync"
	"lib/socket"
)

type NonCenterSvr struct {
	wg         utils.WaitGroupWrapper
	svr4SvrSession map[int32]*Svr4SvrSession
	sync.RWMutex

	UseRunTime bool
}

/*
type Server interface {
	OnCreateConnection(c Conner)
	OnClose()
}

 */

func (this *NonCenterSvr) OnClose() {

}

func (this *NonCenterSvr) OnCreateConnection(c socket.Conner) {
	session := CreateSvrSession(this)
	c.SetSession(session)
	session.connection = c.(*socket.Connection)
}



func (this *NonCenterSvr) RemoveGateway(serverId int32) {
	this.Lock()
	defer this.Unlock()

	if _, ok := this.svr4SvrSession[serverId]; ok {
		delete(this.svr4SvrSession, serverId)
	}
}

func (this *NonCenterSvr) RegistGateway(serverId int32, session *Svr4SvrSession) {
	this.Lock()
	defer this.Unlock()

	if _, ok := this.svr4SvrSession[serverId]; !ok {
		this.svr4SvrSession[serverId] = session
	}
}