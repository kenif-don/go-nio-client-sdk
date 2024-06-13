package handler

import (
	"go-nio-client-sdk/manager"
	"go-nio-client-sdk/model"
	"go-nio-client-sdk/process"
	"go-nio-client-sdk/util"
	"strings"
	"time"

	"github.com/go-netty/go-netty"
)

var wsClientHandler = &WSClientHandler{}

type WSClientHandler struct {
	messageManager  *manager.MessageManager
	process         process.IIMProcess
	reconnectTicker *time.Ticker //重连定时器
	reconnect       Operation
	preGetMsgTime   int64 // 上一次收到消息的时间
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
	//先停 再启动qos
	_self.messageManager.StopQos()
	_self.messageManager.StartupQos()
	//在内部会自动停止 这里只需启动
	_self.startReConnect()
	_self.process.Connected()
}

// HandleRead 接收到服务器消息
// 组装未通用类型
// 交由CoreHandler处理
func (_self *WSClientHandler) HandleRead(ctx netty.InboundContext, message netty.Message) {
	//记录接收到服务器响应时间
	_self.preGetMsgTime = time.Now().Unix()
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
			//自己发送的qos需要删除
			if protocol.Ack == 1 {
				_self.messageManager.HandlerAck(protocol)
			}
			break
		case model.ChannelHeart:
			println("收到服务器心跳")
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
	_self.handlerException(e)
}

// HandleInactive 断开链接
func (_self *WSClientHandler) HandleInactive(ctx netty.InactiveContext, e netty.Exception) {
	_self.handlerException(e)
}
func (_self *WSClientHandler) handlerException(e netty.Exception) {
	if strings.Contains(e.Error(), "An existing connection was forcibly closed by the remote host") ||
		strings.Contains(e.Error(), "unexpected EOF") ||
		strings.Contains(e.Error(), " An established connection was aborted by the software in your host machine.") ||
		strings.Contains(e.Error(), "wsarecv: A connection attempt failed because the connected party did not properly respond after a period of time") {
		return
	} else {
		_self.messageManager.LogicProcess.Exception(e)
	}
}

// startReConnect 重连机制 仅此一处具备重连功能
func (_self *WSClientHandler) startReConnect() {
	_self.reconnectTicker = time.NewTicker(time.Second)
	go func() {
		for {
			select {
			case <-_self.reconnectTicker.C:
				//避免刚发心跳服务器秒回 导致需要等5秒才能收到下一条回复 所以这里要大于5秒
				if time.Now().Unix() > (_self.preGetMsgTime + 6) {
					_self.reconnect()
					//开始重连后 取消定时器 不然会反复调重连方法
					_self.stopReConnect()
				}
			}
		}
	}()
}

func (_self *WSClientHandler) stopReConnect() {
	if _self.reconnectTicker != nil {
		_self.reconnectTicker.Stop()
	}
}
