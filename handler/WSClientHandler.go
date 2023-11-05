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

func NewClientHandler(process process.IIMProcess) *WSClientHandler {
	wsClientHandler.process = process
	return wsClientHandler
}
func (_self *WSClientHandler) HandleActive(ctx netty.ActiveContext) {
	ctx.HandleActive()
	util.Out("【IM】与服务器连接成功！")
	_self.messageManager = manager.New(ctx.Channel(), _self.process)
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
		protocol.Type = res["type"].(int)
		protocol.From = res["from"].(string)
		protocol.To = res["to"].(string)
		protocol.Ack = res["ack"].(int)
		protocol.Data = res["data"].(string)
		protocol.No = res["no"].(string)
		protocol.Ext1 = res["ext1"].(string)
		protocol.Ext2 = res["ext2"].(string)
		protocol.Ext3 = res["ext3"].(string)
		protocol.Ext4 = res["ext4"].(int)
		protocol.Ext5 = res["ext5"].(int)

		switch protocol.Type {
		case model.ChannelLogin:
			if protocol.Ack == 200 {
				_self.messageManager.LogicProcess.LoginOk(protocol)
			} else {
				_self.messageManager.LogicProcess.LoginFail(protocol)
			}
			break
		case model.ChannelOne2oneMsg:
		case model.ChannelGroupMsg:
			//1-自己发出的消息 服务器返回收到的标志 100-别人给自己发送的
			if protocol.Ack == 1 || protocol.Ack == 100 {
				_self.messageManager.SendAck(protocol)
			}
			if protocol.Ack == 1 {
				_self.messageManager.HandlerAck(protocol)
			}
			break
		}
		//触发接收到消息的回调
		_self.messageManager.LogicProcess.ReceivedMessage(protocol)
	}
	ctx.HandleRead(message)
}

func (*WSClientHandler) HandleException(ctx netty.ExceptionContext, ex netty.Exception) {
	util.Err("【IM】协议出现异常 异常信息：%s", ex.Error())
	ctx.HandleException(ex)
}
