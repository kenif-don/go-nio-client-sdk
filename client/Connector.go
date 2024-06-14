package client

import (
	"sync"
)

// Connector manages the reconnection logic for the client.
type Connector struct {
	client    *Client
	mu        sync.Mutex
	isRunning bool
	reconnect chan bool
}

func NewConnector(client *Client) *Connector {
	return &Connector{
		client:    client,
		reconnect: make(chan bool, 1),
	}
}

func (rm *Connector) Start() {
	rm.mu.Lock()
	if rm.isRunning {
		rm.mu.Unlock()
		return
	}
	rm.isRunning = true
	rm.mu.Unlock()

	go rm.reconnectLoop()
}

func (rm *Connector) reconnectLoop() {
	for {
		<-rm.reconnect
		rm.client.connect()
	}
}

func (rm *Connector) TriggerReconnect() {
	select {
	case rm.reconnect <- true:
	default:
	}
}
