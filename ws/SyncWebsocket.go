package ws

import (
	"github.com/gorilla/websocket"
)

// SyncClient stores a client information
type SyncClient struct {
	Conn       *websocket.Conn
	UserID     uint32
	UserName   string
	ProjectMap map[uint32]struct{}
}

// SyncBroadcast is message struct
type SyncBroadcast struct {
	ProjectID uint32
	State     uint8
	Message   string
	// 0=>错误信息 1=>git信息 2=>rsync信息 3=>运行脚本信息
	DataType uint8
}

// SyncHub is a client struct
type SyncHub struct {
	// Registered clients.
	clients map[*SyncClient]bool

	// Inbound messages from the clients.
	Broadcast chan *SyncBroadcast

	// Register requests from the clients.
	Register chan *SyncClient

	// Unregister requests from clients.
	Unregister chan *SyncClient
}

var instance *SyncHub

// GetSyncHub it will only init once in main.go
func GetSyncHub() *SyncHub {
	if instance == nil {
		instance = &SyncHub{
			Broadcast:  make(chan *SyncBroadcast),
			clients:    make(map[*SyncClient]bool),
			Register:   make(chan *SyncClient),
			Unregister: make(chan *SyncClient),
		}
	}
	return instance
}

// Run goroutine run the sync hub
func (hub *SyncHub) Run() {
	for {
		select {
		case client := <-hub.Register:
			hub.clients[client] = true
		case client := <-hub.Unregister:
			if _, ok := hub.clients[client]; ok {
				delete(hub.clients, client)
			}
		case message := <-hub.Broadcast:
			for client := range hub.clients {
				if _, ok := client.ProjectMap[message.ProjectID]; ok {
					client.Conn.WriteJSON(message)
				}
			}
		}
	}
}