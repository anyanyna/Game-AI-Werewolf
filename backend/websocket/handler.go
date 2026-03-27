package websocket

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// Message WebSocket消息结构
type Message struct {
	Type    string          `json:"type"`
	GameID  uint            `json:"game_id,omitempty"`
	PlayerID uint           `json:"player_id,omitempty"`
	Data    json.RawMessage `json:"data,omitempty"`
}

// UpgradeConnection WebSocket连接升级
var UpgradeConnection = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // 允许所有来源的连接
	},
}

// Handler WebSocket处理器
type Handler struct {
	manager *ClientManager
}

// NewHandler 创建WebSocket处理器
func NewHandler(manager *ClientManager) *Handler {
	return &Handler{
		manager: manager,
	}
}

// HandleConnection 处理WebSocket连接
func (h *Handler) HandleConnection(c *gin.Context) {
	// 升级HTTP连接为WebSocket连接
	conn, err := UpgradeConnection.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v\n", err)
		return
	}

	// 获取游戏ID和玩家ID
	gameIDStr := c.Query("game_id")
	playerIDStr := c.Query("player_id")
	userIDStr := c.Query("user_id")

	// 解析参数
	var gameIDUint uint
	var playerIDUint uint
	var userIDUint *uint

	// 解析游戏ID
	if gameIDStr != "" {
		_, err := fmt.Sscanf(gameIDStr, "%d", &gameIDUint)
		if err != nil {
			log.Printf("Invalid game_id: %v\n", err)
		}
	}

	// 解析玩家ID
	if playerIDStr != "" {
		_, err := fmt.Sscanf(playerIDStr, "%d", &playerIDUint)
		if err != nil {
			log.Printf("Invalid player_id: %v\n", err)
		}
	}

	// 解析用户ID
	if userIDStr != "" {
		var uid uint
		_, err := fmt.Sscanf(userIDStr, "%d", &uid)
		if err != nil {
			log.Printf("Invalid user_id: %v\n", err)
		} else {
			userIDUint = &uid
		}
	}

	// 创建客户端
	client := &Client{
		Conn:      conn,
		GameID:    gameIDUint,
		PlayerID:  playerIDUint,
		UserID:    userIDUint,
		Send:      make(chan []byte, 256),
		IsRealPerson: userIDUint != nil,
	}

	// 注册客户端
	h.manager.Register <- client

	// 启动客户端读写协程
	go h.readPump(client)
	go h.writePump(client)
}

// readPump 读取客户端消息
func (h *Handler) readPump(client *Client) {
	defer func() {
		h.manager.Unregister <- client
		client.Conn.Close()
	}()

	// 设置读取超时
	client.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	client.Conn.SetPongHandler(func(string) error {
		client.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, message, err := client.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket read error: %v\n", err)
			}
			break
		}

		// 处理消息
		h.handleMessage(client, message)
	}
}

// writePump 向客户端写入消息
func (h *Handler) writePump(client *Client) {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		client.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-client.Send:
			client.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				// 通道已关闭
				client.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := client.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// 添加队列中的所有消息
			n := len(client.Send)
			for i := 0; i < n; i++ {
				w.Write([]byte("\n"))
				w.Write(<-client.Send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			client.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := client.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// handleMessage 处理客户端消息
func (h *Handler) handleMessage(client *Client, message []byte) {
	// 解析消息
	var msg Message
	if err := json.Unmarshal(message, &msg); err != nil {
		log.Printf("Failed to unmarshal message: %v\n", err)
		return
	}

	// 根据消息类型处理
	switch msg.Type {
	case "join_game":
		// 处理加入游戏
		h.handleJoinGame(client, msg)
	case "night_action":
		// 处理夜晚行动
		h.handleNightAction(client, msg)
	case "day_action":
		// 处理白天行动
		h.handleDayAction(client, msg)
	case "end_phase":
		// 处理结束阶段
		h.handleEndPhase(client, msg)
	case "guess_real_person":
		// 处理真人猜测
		h.handleGuessRealPerson(client, msg)
	default:
		log.Printf("Unknown message type: %s\n", msg.Type)
	}
}

// handleJoinGame 处理加入游戏
func (h *Handler) handleJoinGame(client *Client, msg Message) {
	// 这里应该添加加入游戏的逻辑
}

// handleNightAction 处理夜晚行动
func (h *Handler) handleNightAction(client *Client, msg Message) {
	// 这里应该添加夜晚行动的逻辑
}

// handleDayAction 处理白天行动
func (h *Handler) handleDayAction(client *Client, msg Message) {
	// 这里应该添加白天行动的逻辑
}

// handleEndPhase 处理结束阶段
func (h *Handler) handleEndPhase(client *Client, msg Message) {
	// 这里应该添加结束阶段的逻辑
}

// handleGuessRealPerson 处理真人猜测
func (h *Handler) handleGuessRealPerson(client *Client, msg Message) {
	// 这里应该添加真人猜测的逻辑
}
