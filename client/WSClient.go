package client

import (
	"encoding/binary"
	"errors"
	"go-nio-client-sdk/handler"
	"go-nio-client-sdk/process"
	"go-nio-client-sdk/util"
	"time"

	"github.com/go-netty/go-netty"
	"github.com/go-netty/go-netty-transport/websocket"
	"github.com/go-netty/go-netty/codec/format"
	"github.com/go-netty/go-netty/codec/frame"
	"github.com/go-netty/go-netty/transport"
	"github.com/go-netty/go-netty/transport/tcp"
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
func (_self *WSClient) getTransport(tp string) netty.Option {
	var p transport.Factory
	switch tp {
	case "TCP", "tcp":
		p = tcp.New()
		break
	case "ws", "WS", "websocket":
		p = websocket.New()
		break
	}
	return netty.WithTransport(p)
}

// Startup p 协议 tcp.New()或者 websocket.New()
func (_self *WSClient) Startup(process process.IIMProcess, tp string) error {
	_self.handler = handler.NewClientHandler(_self.Reconnect, process)
	client := func(channel netty.Channel) {
		pipeline := channel.Pipeline()
		switch tp {
		case "TCP", "tcp":
			pipeline.AddLast(frame.LengthFieldCodec(binary.LittleEndian, 1024*1024*10+4, 0, 4, 0, 4))
			break
		case "ws", "WS", "websocket":
			pipeline.AddLast(frame.PacketCodec(1024 * 1024 * 10))
			break
		}
		pipeline.
			//读写超时
			AddLast(&handler.AllIdleHandler{Timeout: 4 * time.Second}).
			AddLast(format.JSONCodec(true, false)).
			AddLast(_self.handler)
	}
	var bootstrap = netty.NewBootstrap(netty.WithClientInitializer(client), _self.getTransport(tp))
	channel, err := bootstrap.Connect(_self.Url)
	_self.Channel = channel
	if err != nil {
		return err
	}
	go func() {
		select {
		case <-channel.Context().Done():
			util.Err("【IM】连接异常断开 重连1？" + channel.Context().Err().Error())
		case <-bootstrap.Context().Done():
			util.Err("【IM】连接异常断开 重连2？" + channel.Context().Err().Error())
		}
	}()
	return nil
}
func (_self *WSClient) Reconnect(tp string) {
	println("【IM】重连中...")
	//如果通道在线 先关闭
	if _self.Channel != nil && _self.Channel.IsActive() {
		_self.Channel.Close(errors.New("【IM】IM客户端正常关闭"))
	}
	//停止Qos
	_self.handler.GetMessageManager().StopQos()
	//再重新启动
	err := _self.Startup(_self.handler.GetMessageManager().LogicProcess, tp)
	if err != nil {
		//延迟重连
		time.Sleep(time.Second * 5)
		_self.Reconnect(tp)
	}
}
func (_self *WSClient) OpenLog() {
	util.OpenLog()
}
