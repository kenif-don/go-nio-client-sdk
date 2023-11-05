package manager

import (
	"github.com/go-netty/go-netty"
	"github.com/go-netty/go-netty/utils"
	"im-sdk/model"
	"im-sdk/process"
	"im-sdk/util"
	"time"
)

type MessageManager struct {
	// qos定时器 对外不可见
	qosTicker *time.Ticker
	// 心跳定时器 对外不可见
	heartbeatTicker *time.Ticker
	// 用来装qos的消息map key需要用来做唯一判断
	qosMessageDTO map[string]*model.QosMsg
	LogicProcess  process.IIMProcess
	Channel       netty.Channel
}

func New(Channel netty.Channel, process process.IIMProcess) *MessageManager {
	return &MessageManager{
		LogicProcess: process,
		Channel:      Channel,
	}
}
func (_self *MessageManager) HandlerAck(protocol *model.Protocol) {
	v := _self.qosMessageDTO[protocol.No]
	//删除Qos中的数据
	delete(_self.qosMessageDTO, protocol.No)
	_self.LogicProcess.SendOk(v.Protocol)
}

// SendLogin 发起登录请求
func (_self *MessageManager) SendLogin(loginInfo *model.LoginInfo) {
	_self.Send(model.NewLoginInfoPack(loginInfo))
}

// SendAck 发送应答包
func (_self *MessageManager) SendAck(protocol *model.Protocol) {
	_self.BaseSend(model.NewAckPack(protocol.No))
}

// Send 通用的发送请求函数
func (_self *MessageManager) Send(protocol *model.Protocol) {
	//判断是否在线 不在线就重连？
	if !_self.Channel.IsActive() {
		util.Out("【IM】IM未连接，重连中...")
		_self.LogicProcess.SendFailedCallback(protocol)
		return
	}
	//ACK为100 且 No不为空 就将消息放入Qos
	if protocol.Ack == 100 && protocol.No != "" {
		//判断qos中是否已存在此消息 存在 那么此消息就不发 交给Qos即可
		if _self.qosMessageDTO[protocol.No].Protocol.No != "" {
			util.Out("【IM】Qos中已存在ID[%s]的消息,直接交由Qos管理，不再发送\n", protocol.No)
			return
		}
		//放入Qos
		_self.qosMessageDTO[protocol.No] = &model.QosMsg{
			Protocol:         protocol,
			PreSendTimeStamp: time.Now(),
			Frequency:        1,
		}
	}
	//发送
	_self.BaseSend(protocol)
	util.Out("【IM】消息发送成功！")
}
func (_self *MessageManager) BaseSend(protocol *model.Protocol) {
	err := _self.Channel.Write(protocol)
	utils.Assert(err)
}

// StartupQos 启动Qos
func (_self *MessageManager) StartupQos() {
	_self.qosTicker = time.NewTicker(time.Second * 2)
	select {
	case <-_self.qosTicker.C:
		for k, msg := range _self.qosMessageDTO {
			//当前发送时间必须比上次发送时间至少间隔QOS_DELAY
			curTime := time.Now()
			if curTime.Unix()-msg.PreSendTimeStamp.Unix() < 2000 {
				continue
			}
			//次数超限--意味着失败
			if msg.Frequency > 15 {
				delete(_self.qosMessageDTO, k)
				_self.LogicProcess.SendFailedCallback(msg.Protocol)
				continue
			}
			//记录当前发送时间
			msg.Frequency++
			msg.PreSendTimeStamp = curTime
			_self.BaseSend(msg.Protocol)
			_self.LogicProcess.SendOkCallback(msg.Protocol)
		}
	}
}

// StopQos 停止Qos
func (_self *MessageManager) StopQos() {
	if _self.qosTicker != nil {
		_self.qosTicker.Stop()
	}
}

// StartupHeartbeat 启动Qos
func (_self *MessageManager) StartupHeartbeat() {
	_self.heartbeatTicker = time.NewTicker(time.Second * 30)
	select {
	case <-_self.heartbeatTicker.C:
		_self.BaseSend(model.NewHeartbeatPack())
	}
}

// StopHeartbeat 停止心跳
func (_self *MessageManager) StopHeartbeat() {
	if _self.heartbeatTicker != nil {
		_self.heartbeatTicker.Stop()
	}
}
