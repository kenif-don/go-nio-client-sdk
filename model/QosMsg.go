package model

import (
	"time"
)

type QosMsg struct {
	Protocol *Protocol `json:"protocol"`
	/**
	 * 上次这个包发送的时间戳
	 * 此字段设计目的：
	 * 定时器中消息每2s发送一次,而且是独立线程运行,可能会出现一种情况就是刚把这个消息包放入qos集合,而qos刚好2s到了
	 * 就立马发送了,但是这个时候正常发送也是才发了,因为丢入qos时，正常也会发送一次
	 * 应该让此消息至少延迟2s再进行第二次发送
	 */
	PreSendTimeStamp time.Time
	/**
	 * 消息的发送次数
	 * 放入队列时会发送一次,所以默认为1
	 */
	Frequency int
}
