package websocket

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

// Client WebSocket客户端
type Client struct {
	Conn      *websocket.Conn
	GameID    uint
	PlayerID  uint
	UserID    *uint
	Send      chan []byte
	IsRealPerson bool
}

// ClientManager 客户端管理器
type ClientManager struct {
	Clients    map[*Client]bool
	GameClients map[uint][]*Client // 游戏ID到客户端的映射
	Register   chan *Client
	Unregister chan *Client
	Mutex      sync.RWMutex
}

// NewClientManager 创建客户端管理器
func NewClientManager() *ClientManager {
	return &ClientManager{
		Clients:    make(map[*Client]bool),
		GameClients: make(map[uint][]*Client),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}

// Run 运行客户端管理器
func (manager *ClientManager) Run() {
	for {
		select {
		case client := <-manager.Register:
			manager.Mutex.Lock()
			manager.Clients[client] = true
			// 添加到游戏客户端列表
			if client.GameID > 0 {
				manager.GameClients[client.GameID] = append(manager.GameClients[client.GameID], client)
			}
			manager.Mutex.Unlock()
			log.Printf("Client registered: GameID=%d, PlayerID=%d\n", client.GameID, client.PlayerID)

		case client := <-manager.Unregister:
			manager.Mutex.Lock()
			if _, ok := manager.Clients[client]; ok {
				delete(manager.Clients, client)
				close(client.Send)
				// 从游戏客户端列表中移除
				if client.GameID > 0 {
					clients := manager.GameClients[client.GameID]
					for i, c := range clients {
						if c == client {
							manager.GameClients[client.GameID] = append(clients[:i], clients[i+1:]...)
							break
						}
					}
					// 如果游戏客户端列表为空，删除该游戏
					if len(manager.GameClients[client.GameID]) == 0 {
						delete(manager.GameClients, client.GameID)
					}
				}
			}
			manager.Mutex.Unlock()
			log.Printf("Client unregistered: GameID=%d, PlayerID=%d\n", client.GameID, client.PlayerID)
		}
	}
}

// BroadcastToGame 向游戏中的所有客户端广播消息
func (manager *ClientManager) BroadcastToGame(gameID uint, message []byte) {
	manager.Mutex.RLock()
	defer manager.Mutex.RUnlock()

	clients := manager.GameClients[gameID]
	for _, client := range clients {
		select {
		case client.Send <- message:
		default:
			close(client.Send)
			delete(manager.Clients, client)
		}
	}
}

// SendToPlayer 向特定玩家发送消息
func (manager *ClientManager) SendToPlayer(gameID uint, playerID uint, message []byte) {
	manager.Mutex.RLock()
	defer manager.Mutex.RUnlock()

	clients := manager.GameClients[gameID]
	for _, client := range clients {
		if client.PlayerID == playerID {
			select {
			case client.Send <- message:
			default:
				close(client.Send)
				delete(manager.Clients, client)
			}
			break
		}
	}
}
