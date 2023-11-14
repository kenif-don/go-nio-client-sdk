package client

import (
	"errors"
	"github.com/go-netty/go-netty"
	"github.com/go-netty/go-netty-transport/websocket"
	"github.com/go-netty/go-netty/codec/format"
	"github.com/go-netty/go-netty/codec/frame"
	"im-sdk/handler"
	"im-sdk/process"
	"im-sdk/util"
)

type WSClient struct {
	Url     string
	Channel netty.Channel
	handler *handler.WSClientHandler
}

func New(url string) *WSClient {
	return &WSClient{
		Url: url,
	}
}
func (_self *WSClient) Startup(process process.IIMProcess) error {
	_self.handler = handler.NewClientHandler(process)
	client := func(channel netty.Channel) {
		channel.Pipeline().
			AddLast(frame.PacketCodec(1024 * 1024 * 10)).
			AddLast(format.JSONCodec(true, false)).
			AddLast(_self.handler)
	}
	var bootstrap = netty.NewBootstrap(netty.WithClientInitializer(client), netty.WithTransport(websocket.New()))
	channel, err := bootstrap.Connect(_self.Url)
	_self.Channel = channel
	if err != nil {
		return err
	}
	go func() {
		select {
		case <-channel.Context().Done():
			util.Err("【IM】连接异常断开 重连1？")
		case <-bootstrap.Context().Done():
			util.Err("【IM】连接异常断开 重连2？")
		}
	}()
	return nil
}
func (_self *WSClient) Reconnect() error {
	//如果通道在线 先关闭
	if _self.Channel != nil && _self.Channel.IsActive() {
		_self.Channel.Close(errors.New("【IM】IM客户端正常关闭"))
	}
	//停止心跳
	_self.handler.GetMessageManager().StopHeartbeat()
	//停止Qos
	_self.handler.GetMessageManager().StopQos()
	//再重新启动
	return _self.Startup(_self.handler.GetMessageManager().LogicProcess)
}
func (_self *WSClient) OpenLog() {
	util.OpenLog()
}
