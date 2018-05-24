package dbagtent

import (
	"reflect"
	"sync"
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

