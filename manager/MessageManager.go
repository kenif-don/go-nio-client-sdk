package manager

import (
	"fmt"
	"go-nio-client-sdk/model"
	"go-nio-client-sdk/process"
	"time"

	"github.com/go-netty/go-netty"
)

type MessageManager struct {
	// qos定时器 对外不可见
	qosTicker *time.Ticker
	// 用来装qos的消息map key需要用来做唯一判断
	qosMessageDTO map[string]*model.QosMsg
	LogicProcess  process.IIMProcess
	Channel       netty.Channel
	preHeartTime  int64
}

func New(Channel netty.Channel, process process.IIMProcess) *MessageManager {
	return &MessageManager{
		LogicProcess:  process,
		Channel:       Channel,
		qosMessageDTO: make(map[string]*model.QosMsg),
	}
}
func (_self *MessageManager) HandlerAck(protocol *model.Protocol) {
	v := _self.qosMessageDTO[protocol.No]
	//删除Qos中的数据
	delete(_self.qosMessageDTO, protocol.No)
	_self.LogicProcess.SendOk(v.Protocol)
}

// SendLogout 发起退出请求
func (_self *MessageManager) SendLogout() {
	_self.Send(&model.Protocol{Type: model.ChannelLoginOut})
	_self.LogicProcess.Logout()
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
	//ACK为100 且 No不为空 就将消息放入Qos
	if protocol.Ack == 100 && protocol.No != "" {
		//判断qos中是否已存在此消息 存在 那么此消息就不发 交给Qos即可
		if _self.qosMessageDTO[protocol.No] != nil {
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
}
func (_self *MessageManager) BaseSend(protocol *model.Protocol) {
	//转换为json对象
	//json, err := util.Obj2Str(protocol)
	//println(json)
	err := _self.Channel.Write(protocol)
	if err != nil {
		fmt.Printf("【IM】IM发送消息失败！ %s\n", err.Error())
	}
}

// StartupQos 启动Qos
func (_self *MessageManager) StartupQos() {
	fmt.Printf("【IM】启动Qos\n")
	_self.qosTicker = time.NewTicker(time.Second * 2)
	go func() {
		for {
			select {
			case <-_self.qosTicker.C:
				for _, msg := range _self.qosMessageDTO {
					//当前发送时间必须比上次发送时间至少间隔QOS_DELAY
					curTime := time.Now()
					if curTime.Unix() < msg.PreSendTimeStamp.Unix()-1000 {
						continue
					}
					////次数超限--意味着失败
					//if msg.Frequency > 30 {
					//	delete(_self.qosMessageDTO, k)
					//	_self.LogicProcess.SendFailedCallback(msg.Protocol)
					//	continue
					//}
					////记录当前发送时间
					//msg.Frequency++
					msg.PreSendTimeStamp = curTime
					_self.BaseSend(msg.Protocol)
					//_self.LogicProcess.SendOkCallback(msg.Protocol)
				}
			}
		}
	}()
}

// StopQos 停止Qos
func (_self *MessageManager) StopQos() {
	if _self.qosTicker != nil {
		fmt.Printf("【IM】停止Qos\n")
		_self.qosTicker.Stop()
	}
}

// SendHeartbeat 发起心跳包
func (_self *MessageManager) SendHeartbeat() {
	//因为读超时和写超时都会发送心跳 如果同时读写都超时了 这里会发送两次 所以通过时间来做一次限制
	now := time.Now().UnixMilli()
	//因为超时设置的5秒 所以这里使用略低于这个值的判断即可
	if now-_self.preHeartTime <= 4500 {
		return
	}
	_self.preHeartTime = now
	_self.BaseSend(model.NewHeartbeatPack())
}
