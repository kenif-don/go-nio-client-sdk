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
	connector   *Connector
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
			AddLast(handler.NewClientHandler(client.Startup, client.process))
	}), client.Tp)
	client.connector = NewConnector(client)
	client.connector.Start()
	return client
}
func (_self *Client) Startup() {
	//如果已经在执行启动 这里会直接结束 不会重复启动
	_self.connector.TriggerReconnect()
}
func (_self *Client) connect() {
	println("【IM】开始链接服务器")
	//如果通道在线 先关闭
	if _self.Channel != nil && _self.Channel.IsActive() {
		_self.Channel.Close(nil)
	}
	//再重新启动
	_self.Channel = nil
	_self.process.OnConnecting()
	channel, err := _self.nettyClient.Connect(_self.Url)
	if err == nil {
		_self.Channel = channel
		return

	}
	//延迟重连
	time.Sleep(time.Second * 2)
	//不是递归调用 而是利用连接器
	_self.Startup()
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
