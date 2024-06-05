package handler

import (
	"fmt"
	"net"
	"time"

	"github.com/go-netty/go-netty"
)

type AllIdleHandler struct {
	Timeout time.Duration
}

func (h *AllIdleHandler) HandleActive(ctx netty.ActiveContext) {
	if conn, ok := ctx.Channel().Transport().(net.Conn); ok {
		conn.SetDeadline(time.Now().Add(h.Timeout))
	}
	ctx.HandleActive()
}
func (h *AllIdleHandler) HandleRead(ctx netty.InboundContext, message netty.Message) {
	if conn, ok := ctx.Channel().Transport().(net.Conn); ok {
		conn.SetDeadline(time.Now().Add(h.Timeout))
	}
	fmt.Printf("收到消息，重置超时时间\n")
	ctx.HandleRead(message)
}
func (h *AllIdleHandler) HandleWrite(ctx netty.OutboundContext, message netty.Message) {
	if conn, ok := ctx.Channel().Transport().(net.Conn); ok {
		conn.SetDeadline(time.Now().Add(h.Timeout))
	}
	fmt.Printf("发出消息，重置超时时间\n")
	ctx.HandleWrite(message)
}
func (h *AllIdleHandler) HandleInactive(ctx netty.InactiveContext, ex netty.Exception) {
	ctx.HandleInactive(ex)
}
