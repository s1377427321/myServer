package route

import (
	"reflect"
	"github.com/astaxie/beego"
	."protocol"
)

var funcs = make(map[uint32]*router,50)

type router struct {
	fun reflect.Value
	msg reflect.Type
}

// cmd 	消息id
// f   	处理消息的函数 如: login(session *server.TmpSession, req *protocol.CLogin) (resp proto.Message, err error)
// msg 	消息对应的protobuf请求包类型
func RegisterFunc(cmd CMD, f interface{}, msg interface{}) {
	_, ok := funcs[uint32(cmd)]
	if ok {
		beego.Warn("消息重复注册", cmd, cmd.String())
		return
	}

	if reflect.TypeOf(f).Kind() != reflect.Func {
		beego.Warn("消息处理不是函数", cmd, cmd.String())
		return
	}

	if reflect.TypeOf(msg).Kind() != reflect.Struct {
		beego.Warn("请求消息体不是一个结构体", cmd, cmd.String())
		return
	}
	funcs[uint32(cmd)] = &router{
		fun: reflect.ValueOf(f),
		msg: reflect.TypeOf(msg),
	}
}

func GetFunc(cmd uint32) reflect.Value {
	return funcs[cmd].fun
}

func Exist(cmd uint32) bool {
	_, ok := funcs[cmd]
	return ok
}

func GetProto(cmd uint32) reflect.Type {
	return funcs[cmd].msg
}
