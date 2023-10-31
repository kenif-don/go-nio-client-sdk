package process

import (
	"IM-SDK/model"
)

type IIMProcess interface {
	//SendOkCallback 发送成功的回调
	//仅仅是发出去了 如果是Qos消息 此时还未收到服务器反馈
	//SendOk代表发出Qos消息并接收到了服务器反馈
	SendOkCallback(protocol *model.Protocol)
	//SendFailedCallback 发送失败的回调
	SendFailedCallback(protocol *model.Protocol)
	//LoginOk 登录成功的回调
	LoginOk(protocol *model.Protocol)
	//LoginFail 登录失败的回调
	LoginFail(protocol *model.Protocol)
	//ReceivedMessage 接收到消息
	ReceivedMessage(protocol *model.Protocol)
	//SendOk qos中的消息发送成功 服务器成功返回
	SendOk(protocol *model.Protocol)
}

var LogicProcess = &DefaultIMProcess{}

type DefaultIMProcess struct {
}

func New() *DefaultIMProcess {
	return LogicProcess
}
func (*DefaultIMProcess) SendOkCallback(protocol *model.Protocol) {

}
func (*DefaultIMProcess) SendFailedCallback(protocol *model.Protocol) {

}
func (*DefaultIMProcess) LoginOk(protocol *model.Protocol) {

}
func (*DefaultIMProcess) LoginFail(protocol *model.Protocol) {

}
func (*DefaultIMProcess) ReceivedMessage(protocol *model.Protocol) {

}
func (*DefaultIMProcess) SendOk(protocol *model.Protocol) {

}
