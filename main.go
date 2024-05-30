package main

import (
	"go-nio-client-sdk/client"
	"go-nio-client-sdk/handler"
	"go-nio-client-sdk/model"
	"go-nio-client-sdk/util"

	"github.com/go-netty/go-netty"
)

type IMProcess struct {
}

func (_self *IMProcess) Connected() {
	//登录
	handler.GetClientHandler().GetMessageManager().SendLogin(&model.LoginInfo{Id: "123", Device: "123", Token: "123"})
}
func (_self *IMProcess) SendOkCallback(protocol *model.Protocol) {

}
func (_self *IMProcess) SendFailedCallback(protocol *model.Protocol) {

}
func (_self *IMProcess) LoginOk(protocol *model.Protocol) {
	util.Out("登录成功 %v", protocol)
}
func (_self *IMProcess) LoginFail(protocol *model.Protocol) {

}
func (_self *IMProcess) ReceivedMessage(protocol *model.Protocol) {

}
func (_self *IMProcess) SendOk(protocol *model.Protocol) {

}
func (_self *IMProcess) Exception(ctx netty.ExceptionContext, e netty.Exception) {
	println("链接异常", e.Error())
}
func (_self *IMProcess) Logout() {

}
func main() {
	ct := client.New("ws://127.0.0.1:1003")
	e := ct.Startup(&IMProcess{}, "ws")
	if e != nil {
		panic(e)
		return
	}
	select {}
}
