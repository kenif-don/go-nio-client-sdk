package client

import (
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
}

func New(url string) *WSClient {
	return &WSClient{
		Url: url,
	}
}
func (_self *WSClient) Startup(process process.IIMProcess) error {
	client := func(channel netty.Channel) {
		channel.Pipeline().
			AddLast(frame.PacketCodec(1024)).
			AddLast(format.JSONCodec(true, false)).
			AddLast(handler.NewClientHandler(process))
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
