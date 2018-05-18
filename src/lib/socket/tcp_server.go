package socket

type Server interface {
	OnCreateConnection(c Conner)
	OnClose()
}

//启动tcp server
func StartTcpListen(svr Server, address string) {
	defer svr.OnClose()



}
