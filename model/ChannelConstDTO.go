package model

// 协议常量
const (
	ChannelLogin      int = 0 //客户端登录
	ChannelOne2oneMsg int = 1 //单聊消息交互
	ChannelAck        int = 2 //消息回执
	ChannelHeart      int = 3 //心跳类型
	ChannelGroupMsg   int = 8 //群聊消息交互
)
