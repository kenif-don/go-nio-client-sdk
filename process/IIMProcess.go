package process

import (
	"go-nio-client-sdk/model"

	"github.com/go-netty/go-netty"
)

type IIMProcess interface {
	//Connected 与服务器链接成功
	Connected()
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
	//SendOk qos中的消息发送成功 服务器成功返回 消息真正的发送成功了
	SendOk(protocol *model.Protocol)
	//Exception 链接发生异常
	Exception(ctx netty.ExceptionContext, e netty.Exception)
	//Logout 退出登录回调 可以在里面做重连 只有手动调用MessageManager->SendLogout才有这个回调
	Logout()
}
