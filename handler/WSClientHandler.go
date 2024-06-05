package handler

import (
	"fmt"
	"go-nio-client-sdk/manager"
	"go-nio-client-sdk/model"
	"go-nio-client-sdk/process"
	"go-nio-client-sdk/util"
	"strings"

	"github.com/go-netty/go-netty"
)

var wsClientHandler = &WSClientHandler{}

type WSClientHandler struct {
	messageManager *manager.MessageManager
	process        process.IIMProcess
	reconnect      Operation
}
type Operation func()

func GetClientHandler() *WSClientHandler {
	return wsClientHandler
}
func NewClientHandler(reconnect Operation, process process.IIMProcess) *WSClientHandler {
	wsClientHandler.process = process
	wsClientHandler.reconnect = reconnect
	return wsClientHandler
}
func (_self *WSClientHandler) GetMessageManager() *manager.MessageManager {
	return _self.messageManager
}

// HandleActive 客户端链接
func (_self *WSClientHandler) HandleActive(ctx netty.ActiveContext) {
	ctx.HandleActive()
	println("【IM】与服务器连接成功")
	_self.messageManager = manager.New(ctx.Channel(), _self.process)
	//启动qos
	_self.messageManager.StartupQos()
	_self.process.Connected()
}

// HandleRead 接收到服务器消息
// 组装未通用类型
// 交由CoreHandler处理
func (_self *WSClientHandler) HandleRead(ctx netty.InboundContext, message netty.Message) {
	if res, ok := message.(map[string]interface{}); ok {
		protocol := model.NewProtocol()
		err := util.Map2Obj(res, protocol)
		if err != nil {
			_self.messageManager.LogicProcess.Exception(nil, err)
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
			//自己发送的qos需要删除
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
func (_self *WSClientHandler) HandleEvent(ctx netty.EventContext, event netty.Event) {
	if _, ok := event.(netty.ReadIdleEvent); ok {
		//心跳
		_self.messageManager.SendHeartbeat()
	} else if _, ok := event.(netty.WriteIdleEvent); ok {
		//心跳
		_self.messageManager.SendHeartbeat()
	}
	ctx.HandleEvent(event)
}

// HandleException 处理异常
func (_self *WSClientHandler) HandleException(ctx netty.ExceptionContext, e netty.Exception) {
	if strings.Contains(e.Error(), "An existing connection was forcibly closed by the remote host") ||
		strings.Contains(e.Error(), "unexpected EOF") ||
		strings.Contains(e.Error(), " An established connection was aborted by the software in your host machine.") {
		//重连
		_self.reconnect()
	} else {
		_self.messageManager.LogicProcess.Exception(ctx, e)
	}
}

// HandleInactive 断开链接
func (_self *WSClientHandler) HandleInactive(ctx netty.InactiveContext, ex netty.Exception) {
	fmt.Printf("【IM】链接断开,异常信息:%s \n", ex.Error())
}
