package model

type LoginInfo struct {
	Id     string `json:"id"`
	Device string `json:"device"`
	Token  string `json:"token"`
}

func NewLoginInfo(id string, device string, token string) *LoginInfo {
	return &LoginInfo{
		Id:     id,
		Device: device,
		Token:  token,
	}
}
