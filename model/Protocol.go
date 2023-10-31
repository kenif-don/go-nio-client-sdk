package model

type Protocol struct {
	Type int         `json:"type"`
	From string      `json:"from"`
	To   string      `json:"to"`
	Data interface{} `json:"data"`
	Ack  int         `json:"ack"`
	No   string      `json:"no"`
	Ext1 string      `json:"ext1"`
	Ext2 string      `json:"ext2"`
	Ext3 string      `json:"ext3"`
	Ext4 int         `json:"ext4"`
	Ext5 int         `json:"ext5"`
}

func NewProtocol() *Protocol {
	return &Protocol{}
}

// NewLoginInfoPack 创建一个用于登录的数据包
func NewLoginInfoPack(loginInfo *LoginInfo) *Protocol {
	return &Protocol{
		Type: 0,
		From: loginInfo.Id,
		Data: loginInfo,
	}
}
func NewHeartbeatPack() *Protocol {
	return &Protocol{
		Type: ChannelHeart,
	}
}
func NewAckPack(no string) *Protocol {
	return &Protocol{
		Type: ChannelAck,
		No:   no,
	}
}
