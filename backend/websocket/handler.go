package websocket

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"werewolf-game/backend/game"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// Message WebSocket消息结构
type Message struct {
	Type     string          `json:"type"`
	GameID   uint            `json:"game_id,omitempty"`
	PlayerID uint            `json:"player_id,omitempty"`
	UserID   uint            `json:"user_id,omitempty"`
	Data     json.RawMessage `json:"data,omitempty"`
}

// MessageData 消息数据体
type MessageData struct {
	ActionType   string `json:"action_type,omitempty"`
	TargetID     uint   `json:"target_id,omitempty"`
	Content      string `json:"content,omitempty"`
	IsRealPerson bool   `json:"is_real_person,omitempty"`
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
	manager     *ClientManager
	gameService *game.GameService
}

// NewHandler 创建WebSocket处理器
func NewHandler(manager *ClientManager, gameService *game.GameService) *Handler {
	return &Handler{
		manager:     manager,
		gameService: gameService,
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
	var userIDUint uint

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
		_, err := fmt.Sscanf(userIDStr, "%d", &userIDUint)
		if err != nil {
			log.Printf("Invalid user_id: %v\n", err)
		}
	}

	// 创建客户端
	client := &Client{
		Conn:         conn,
		GameID:       gameIDUint,
		PlayerID:     playerIDUint,
		UserID:       &userIDUint,
		Send:         make(chan []byte, 256),
		IsRealPerson: userIDUint != 0,
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
		h.handleJoinGame(client, msg)
	case "night_action":
		h.handleNightAction(client, msg)
	case "day_action":
		h.handleDayAction(client, msg)
	case "vote":
		h.handleVote(client, msg)
	case "end_phase":
		h.handleEndPhase(client, msg)
	case "guess_real_person":
		h.handleGuessRealPerson(client, msg)
	case "get_players":
		h.handleGetPlayers(client, msg)
	case "get_phase":
		h.handleGetPhase(client, msg)
	default:
		log.Printf("Unknown message type: %s\n", msg.Type)
	}
}

// sendMessage 发送消息给客户端
func (h *Handler) sendMessage(client *Client, msgType string, data interface{}) {
	msg := map[string]interface{}{
		"type": msgType,
	}
	if data != nil {
		msg["data"] = data
	}

	msgBytes, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Failed to marshal message: %v\n", err)
		return
	}

	select {
	case client.Send <- msgBytes:
	default:
		log.Printf("Failed to send message, channel full")
	}
}

// broadcastToGame 广播消息给游戏所有玩家
func (h *Handler) broadcastToGame(gameID uint, msgType string, data interface{}) {
	msg := map[string]interface{}{
		"type": msgType,
	}
	if data != nil {
		msg["data"] = data
	}

	msgBytes, _ := json.Marshal(msg)
	h.manager.BroadcastToGame(gameID, msgBytes)
}

// handleJoinGame 处理加入游戏
func (h *Handler) handleJoinGame(client *Client, msg Message) {
	gameCode := string(msg.Data)
	if gameCode == "" {
		h.sendMessage(client, "error", map[string]string{"message": "Game code required"})
		return
	}

	// 如果提供了用户ID，尝试加入游戏
	if msg.UserID != 0 {
		game, player, err := h.gameService.JoinGame(gameCode, msg.UserID)
		if err != nil {
			h.sendMessage(client, "error", map[string]string{"message": err.Error()})
			return
		}

		client.GameID = game.ID
		client.PlayerID = player.ID
		client.UserID = &msg.UserID

		// 发送游戏信息
		h.sendMessage(client, "game_joined", map[string]interface{}{
			"game":   game,
			"player": player,
		})

		// 广播给其他玩家
		h.broadcastToGame(game.ID, "player_joined", map[string]interface{}{
			"player": player,
		})
	}
}

// handleNightAction 处理夜晚行动
func (h *Handler) handleNightAction(client *Client, msg Message) {
	if client.GameID == 0 {
		h.sendMessage(client, "error", map[string]string{"message": "Not in a game"})
		return
	}

	var data MessageData
	if err := json.Unmarshal(msg.Data, &data); err != nil {
		h.sendMessage(client, "error", map[string]string{"message": "Invalid data"})
		return
	}

	err := h.gameService.ProcessNightAction(client.GameID, client.PlayerID, data.ActionType, data.TargetID)
	if err != nil {
		h.sendMessage(client, "error", map[string]string{"message": err.Error()})
		return
	}

	h.sendMessage(client, "night_action_ok", map[string]interface{}{
		"action_type": data.ActionType,
		"target_id":   data.TargetID,
	})
}

// handleDayAction 处理白天行动
func (h *Handler) handleDayAction(client *Client, msg Message) {
	if client.GameID == 0 {
		h.sendMessage(client, "error", map[string]string{"message": "Not in a game"})
		return
	}

	var data MessageData
	if err := json.Unmarshal(msg.Data, &data); err != nil {
		h.sendMessage(client, "error", map[string]string{"message": "Invalid data"})
		return
	}

	err := h.gameService.ProcessDayAction(client.GameID, client.PlayerID, data.ActionType, data.Content, nil)
	if err != nil {
		h.sendMessage(client, "error", map[string]string{"message": err.Error()})
		return
	}

	h.sendMessage(client, "day_action_ok", nil)
	h.broadcastToGame(client.GameID, "player_spoke", map[string]interface{}{
		"player_id": client.PlayerID,
		"content":   data.Content,
	})
}

// handleVote 处理投票
func (h *Handler) handleVote(client *Client, msg Message) {
	if client.GameID == 0 {
		h.sendMessage(client, "error", map[string]string{"message": "Not in a game"})
		return
	}

	var data MessageData
	if err := json.Unmarshal(msg.Data, &data); err != nil {
		h.sendMessage(client, "error", map[string]string{"message": "Invalid data"})
		return
	}

	err := h.gameService.ProcessDayAction(client.GameID, client.PlayerID, "vote", "", &data.TargetID)
	if err != nil {
		h.sendMessage(client, "error", map[string]string{"message": err.Error()})
		return
	}

	h.sendMessage(client, "vote_ok", map[string]interface{}{
		"target_id": data.TargetID,
	})

	h.broadcastToGame(client.GameID, "player_voted", map[string]interface{}{
		"player_id": client.PlayerID,
		"target_id": data.TargetID,
	})
}

// handleEndPhase 处理结束阶段
func (h *Handler) handleEndPhase(client *Client, msg Message) {
	if client.GameID == 0 {
		h.sendMessage(client, "error", map[string]string{"message": "Not in a game"})
		return
	}

	phase, err := h.gameService.EndPhase(client.GameID)
	if err != nil {
		h.sendMessage(client, "error", map[string]string{"message": err.Error()})
		return
	}

	// 广播阶段变更
	h.broadcastToGame(client.GameID, "phase_changed", map[string]interface{}{
		"phase": phase,
	})

	h.sendMessage(client, "phase_ended", map[string]interface{}{
		"phase": phase,
	})
}

// handleGuessRealPerson 处理真人猜测
func (h *Handler) handleGuessRealPerson(client *Client, msg Message) {
	if client.GameID == 0 {
		h.sendMessage(client, "error", map[string]string{"message": "Not in a game"})
		return
	}

	var data MessageData
	if err := json.Unmarshal(msg.Data, &data); err != nil {
		h.sendMessage(client, "error", map[string]string{"message": "Invalid data"})
		return
	}

	guess, err := h.gameService.GuessRealPerson(client.GameID, *client.UserID, data.TargetID, data.IsRealPerson)
	if err != nil {
		h.sendMessage(client, "error", map[string]string{"message": err.Error()})
		return
	}

	h.sendMessage(client, "guess_result", map[string]interface{}{
		"guess": guess,
	})
}

// handleGetPlayers 获取玩家列表
func (h *Handler) handleGetPlayers(client *Client, msg Message) {
	if client.GameID == 0 {
		h.sendMessage(client, "error", map[string]string{"message": "Not in a game"})
		return
	}

	players, err := h.gameService.GetPlayers(client.GameID)
	if err != nil {
		h.sendMessage(client, "error", map[string]string{"message": err.Error()})
		return
	}

	h.sendMessage(client, "players", map[string]interface{}{
		"players": players,
	})
}

// handleGetPhase 获取当前阶段
func (h *Handler) handleGetPhase(client *Client, msg Message) {
	if client.GameID == 0 {
		h.sendMessage(client, "error", map[string]string{"message": "Not in a game"})
		return
	}

	phase, err := h.gameService.GetCurrentPhase(client.GameID)
	if err != nil {
		h.sendMessage(client, "error", map[string]string{"message": err.Error()})
		return
	}

	h.sendMessage(client, "phase", map[string]interface{}{
		"phase": phase,
	})
}
