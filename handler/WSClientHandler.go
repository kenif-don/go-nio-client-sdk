package handler

import (
	"github.com/go-netty/go-netty"
	"im-sdk/manager"
	"im-sdk/model"
	"im-sdk/process"
	"im-sdk/util"
)

var wsClientHandler = &WSClientHandler{}

type WSClientHandler struct {
	messageManager *manager.MessageManager
	process        process.IIMProcess
}

func GetClientHandler() *WSClientHandler {
	return wsClientHandler
}
func NewClientHandler(process process.IIMProcess) *WSClientHandler {
	wsClientHandler.process = process
	return wsClientHandler
}
func (_self *WSClientHandler) HandleActive(ctx netty.ActiveContext) {
	ctx.HandleActive()
	util.Out("【IM】与服务器连接成功！")
	_self.messageManager = manager.New(ctx.Channel(), _self.process)
	_self.process.Connected()
	//启动心跳
	_self.messageManager.StartupHeartbeat()
	//启动qos
	_self.messageManager.StartupQos()
}
func (_self *WSClientHandler) GetMessageManager() *manager.MessageManager {
	return _self.messageManager
}

// HandleRead 接收到服务器消息
// 组装未通用类型
// 交由CoreHandler处理
func (_self *WSClientHandler) HandleRead(ctx netty.InboundContext, message netty.Message) {
	if res, ok := message.(map[string]interface{}); ok {
		protocol := model.NewProtocol()
		err := util.Map2Obj(res, protocol)
		if err != nil {
			_self.messageManager.LogicProcess.Exception(err)
			return
		}
		//1-自己发出的消息 服务器返回收到的标志 100-别人给自己发送的
		if protocol.Ack == 1 || protocol.Ack == 100 {
			_self.messageManager.SendAck(protocol)
		}
		switch protocol.Type {
		case model.ChannelLogin:
			if protocol.Ack == 200 {
				_self.process.LoginOk(protocol)
			} else {
				_self.process.LoginFail(protocol)
			}
			break
		case model.ChannelOne2oneMsg, model.ChannelGroupMsg:
			if protocol.Ack == 1 {
				_self.messageManager.HandlerAck(protocol)
			}
			break
		}
		//触发接收到消息的回调
		_self.process.ReceivedMessage(protocol)
	}
	ctx.HandleRead(message)
}

func (_self *WSClientHandler) HandleException(ctx netty.ExceptionContext, err netty.Exception) {
	ctx.HandleException(err)
	_self.messageManager.LogicProcess.Exception(err)
}
