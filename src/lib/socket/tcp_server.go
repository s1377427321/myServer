package socket


//启动tcp server
func StartTcpListen(svr Server, address string) {
	defer svr.OnClose()



}
