package main

import (
	"fmt"
	"go-nio-client-sdk/client"
	"go-nio-client-sdk/handler"
	"go-nio-client-sdk/model"

	"github.com/go-netty/go-netty"
)

type IMProcess struct {
}

// OnConnecting 链接中 可以在这里做链接状态更新
func (_self *IMProcess) OnConnecting() {

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
	fmt.Printf("登录成功 %v \n", protocol)
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
func (_self *IMProcess) Disconnect() {
	ct.Reconnect()
}

var ct *client.Client

func main() {
	ct = client.New("ws", "ws://127.0.0.1:1003", &IMProcess{})
	e := ct.Startup()
	if e != nil {
		panic(e)
		return
	}
	select {}
}
