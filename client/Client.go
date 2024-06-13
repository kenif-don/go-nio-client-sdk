package client

import (
	"encoding/binary"
	"go-nio-client-sdk/handler"
	"go-nio-client-sdk/process"
	"time"

	"github.com/go-netty/go-netty"
	"github.com/go-netty/go-netty-transport/websocket"
	"github.com/go-netty/go-netty/codec/format"
	"github.com/go-netty/go-netty/codec/frame"
	"github.com/go-netty/go-netty/transport"
	"github.com/go-netty/go-netty/transport/tcp"
)

type Client struct {
	// Tp 协议类型
	Tp          netty.Option
	Url         string
	nettyClient netty.Bootstrap
	Channel     netty.Channel
	process     process.IIMProcess
}

func New(tp, url string, process process.IIMProcess) *Client {
	client := &Client{
		Url:         url,
		Tp:          getTransport(tp),
		nettyClient: netty.NewBootstrap(),
		process:     process,
	}
	client.nettyClient = netty.NewBootstrap(netty.WithClientInitializer(func(channel netty.Channel) {
		pipeline := channel.Pipeline()
		switch tp {
		case "TCP", "tcp":
			pipeline.AddLast(frame.LengthFieldCodec(binary.LittleEndian, 1024*1024*10+4, 0, 4, 0, 4))
			break
		case "ws", "WS", "websocket":
			pipeline.AddLast(frame.PacketCodec(1024 * 1024))
			break
		}
		pipeline.
			//读写超时
			AddLast(netty.ReadIdleHandler(5 * time.Second)).
			AddLast(netty.WriteIdleHandler(5 * time.Second)).
			AddLast(format.JSONCodec(true, false)).
			AddLast(handler.NewClientHandler(client.Reconnect, client.process))
	}), client.Tp)

	return client
}
func getTransport(tp string) netty.Option {
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

func (_self *Client) Startup() error {
	_self.process.OnConnecting()
	channel, err := _self.nettyClient.Connect(_self.Url)
	if err != nil {
		_self.Channel = nil
		return err
	}
	_self.Channel = channel
	return nil
}
func (_self *Client) Reconnect() {
	//如果通道在线 先关闭
	if _self.Channel != nil && _self.Channel.IsActive() {
		_self.Channel.Close(nil)
	}
	println("【IM】重连中...")
	//再重新启动
	err := _self.Startup()
	if err != nil {
		//延迟重连
		time.Sleep(time.Second * 2)
		_self.Reconnect()
	}
}
