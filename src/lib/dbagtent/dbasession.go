package dbagtent

import (
	"lib/socket"
)




type DBASession struct {
	connection    *socket.Connection
	ServerId      int32
	ServerAddress string
	dbagent       *DBAgent
}
