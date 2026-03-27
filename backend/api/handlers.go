package api

import (
	"net/http"
	"strconv"
	"time"

	"werewolf-game/backend/models"

	"github.com/gin-gonic/gin"
)

// createGame 创建游戏
func createGame(c *gin.Context) {
	// 获取请求参数
	var req struct {
		HostID uint `json:"host_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// 返回模拟数据
	c.JSON(http.StatusOK, gin.H{"game": models.Game{
		ID:        1,
		Status:    "waiting",
		GameCode:  "GAME123",
		HostID:    req.HostID,
		CreatedAt: time.Now(),
	}})
}

// joinGame 加入游戏
func joinGame(c *gin.Context) {
	// 获取请求参数
	var req struct {
		GameCode string `json:"game_code" binding:"required"`
		UserID   uint   `json:"user_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// 返回模拟数据
	c.JSON(http.StatusOK, gin.H{"game": models.Game{
		ID:        1,
		Status:    "waiting",
		GameCode:  req.GameCode,
		HostID:    1,
		CreatedAt: time.Now(),
	}})
}

// startGame 开始游戏
func startGame(c *gin.Context) {
	// 获取请求参数
	var req struct {
		GameID uint `json:"game_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// 返回模拟数据
	c.JSON(http.StatusOK, gin.H{"message": "Game started successfully"})
}

// getGameInfo 获取游戏信息
func getGameInfo(c *gin.Context) {
	// 获取游戏ID
	gameIDStr := c.Param("game_id")
	_, err := strconv.ParseUint(gameIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid game ID"})
		return
	}

	// 返回模拟数据
	c.JSON(http.StatusOK, gin.H{"game": models.Game{
		ID:        1,
		Status:    "playing",
		GameCode:  "GAME123",
		HostID:    1,
		CreatedAt: time.Now(),
	}})
}

// getGamePlayers 获取游戏玩家
func getGamePlayers(c *gin.Context) {
	// 获取游戏ID
	gameIDStr := c.Param("game_id")
	_, err := strconv.ParseUint(gameIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid game ID"})
		return
	}

	// 返回模拟数据
	c.JSON(http.StatusOK, gin.H{"players": []models.GamePlayer{
		{
			ID:            1,
			GameID:        1,
			UserID:        uintPtr(1),
			Role:          "werewolf",
			Status:        "alive",
			IsRealPerson:  true,
			Number:        1,
			CreatedAt:     time.Now(),
		},
		{
			ID:            2,
			GameID:        1,
			UserID:        nil,
			AIDigitalPersonID: uintPtr(1),
			Role:          "villager",
			Status:        "alive",
			IsRealPerson:  false,
			Number:        2,
			CreatedAt:     time.Now(),
		},
		{
			ID:            3,
			GameID:        1,
			UserID:        nil,
			AIDigitalPersonID: uintPtr(2),
			Role:          "seer",
			Status:        "alive",
			IsRealPerson:  false,
			Number:        3,
			CreatedAt:     time.Now(),
		},
		{
			ID:            4,
			GameID:        1,
			UserID:        nil,
			AIDigitalPersonID: uintPtr(3),
			Role:          "witch",
			Status:        "alive",
			IsRealPerson:  false,
			Number:        4,
			CreatedAt:     time.Now(),
		},
		{
			ID:            5,
			GameID:        1,
			UserID:        nil,
			AIDigitalPersonID: uintPtr(4),
			Role:          "hunter",
			Status:        "alive",
			IsRealPerson:  false,
			Number:        5,
			CreatedAt:     time.Now(),
		},
		{
			ID:            6,
			GameID:        1,
			UserID:        nil,
			AIDigitalPersonID: uintPtr(5),
			Role:          "villager",
			Status:        "alive",
			IsRealPerson:  false,
			Number:        6,
			CreatedAt:     time.Now(),
		},
	}})
}

// getPlayerInfo 获取玩家信息
func getPlayerInfo(c *gin.Context) {
	// 获取玩家ID
	playerIDStr := c.Param("player_id")
	_, err := strconv.ParseUint(playerIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid player ID"})
		return
	}

	// 返回模拟数据
	c.JSON(http.StatusOK, gin.H{"player": models.GamePlayer{
		ID:            1,
		GameID:        1,
		UserID:        uintPtr(1),
		Role:          "werewolf",
		Status:        "alive",
		IsRealPerson:  true,
		Number:        1,
		CreatedAt:     time.Now(),
	}})
}

// nightAction 夜晚行动
func nightAction(c *gin.Context) {
	// 获取请求参数
	var req struct {
		GameID     uint   `json:"game_id" binding:"required"`
		PlayerID   uint   `json:"player_id" binding:"required"`
		ActionType string `json:"action_type" binding:"required"`
		TargetID   uint   `json:"target_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// 返回模拟数据
	c.JSON(http.StatusOK, gin.H{"message": "Night action executed successfully"})
}

// dayAction 白天行动
func dayAction(c *gin.Context) {
	// 获取请求参数
	var req struct {
		GameID     uint   `json:"game_id" binding:"required"`
		PlayerID   uint   `json:"player_id" binding:"required"`
		ActionType string `json:"action_type" binding:"required"`
		Content    string `json:"content"`
		TargetID   *uint  `json:"target_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// 返回模拟数据
	c.JSON(http.StatusOK, gin.H{"message": "Day action executed successfully"})
}

// voteAction 投票行动
func voteAction(c *gin.Context) {
	// 获取请求参数
	var req struct {
		GameID     uint   `json:"game_id" binding:"required"`
		PlayerID   uint   `json:"player_id" binding:"required"`
		TargetID   uint   `json:"target_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// 返回模拟数据
	c.JSON(http.StatusOK, gin.H{"message": "Vote action executed successfully"})
}

// endPhase 结束阶段
func endPhase(c *gin.Context) {
	// 获取请求参数
	var req struct {
		GameID uint `json:"game_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// 返回模拟数据
	c.JSON(http.StatusOK, gin.H{"message": "Phase ended successfully"})
}

// getCurrentPhase 获取当前阶段
func getCurrentPhase(c *gin.Context) {
	// 获取游戏ID
	gameIDStr := c.Param("game_id")
	_, err := strconv.ParseUint(gameIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid game ID"})
		return
	}

	// 返回模拟数据
	c.JSON(http.StatusOK, gin.H{"phase": models.GamePhase{
		ID:        1,
		GameID:    1,
		Phase:     "night",
		Round:     1,
		StartTime: time.Now(),
	}})
}

// guessRealPerson 真人猜测
func guessRealPerson(c *gin.Context) {
	// 获取请求参数
	var req struct {
		GameID       uint `json:"game_id" binding:"required"`
		UserID       uint `json:"user_id" binding:"required"`
		TargetID     uint `json:"target_id" binding:"required"`
		IsRealPerson bool `json:"is_real_person"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// 返回模拟数据
	c.JSON(http.StatusOK, gin.H{"guess": models.RealPersonGuess{
		ID:           1,
		GameID:       req.GameID,
		UserID:       req.UserID,
		TargetID:     req.TargetID,
		IsRealPerson: req.IsRealPerson,
		IsCorrect:    true,
		CreatedAt:    time.Now(),
	}})
}

// getGuessResult 获取猜测结果
func getGuessResult(c *gin.Context) {
	// 获取游戏ID
	gameIDStr := c.Param("game_id")
	_, err := strconv.ParseUint(gameIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid game ID"})
		return
	}

	// 返回模拟数据
	c.JSON(http.StatusOK, gin.H{"guesses": []models.RealPersonGuess{
		{
			ID:           1,
			GameID:       1,
			UserID:       1,
			TargetID:     2,
			IsRealPerson: false,
			IsCorrect:    true,
			CreatedAt:    time.Now(),
		},
		{
			ID:           2,
			GameID:       1,
			UserID:       1,
			TargetID:     3,
			IsRealPerson: false,
			IsCorrect:    true,
			CreatedAt:    time.Now(),
		},
	}})
}

// registerUser 注册用户
func registerUser(c *gin.Context) {
	// 获取请求参数
	var req struct {
		Username string `json:"username" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// 返回模拟数据
	c.JSON(http.StatusOK, gin.H{"user": models.User{
		ID:           1,
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: req.Password,
		Role:         "user",
		CreatedAt:    time.Now(),
	}})
}

// loginUser 登录用户
func loginUser(c *gin.Context) {
	// 获取请求参数
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// 返回模拟数据
	c.JSON(http.StatusOK, gin.H{
		"user": models.User{
			ID:           1,
			Username:     req.Username,
			Email:        "test@example.com",
			PasswordHash: req.Password,
			Role:         "user",
			CreatedAt:    time.Now(),
		},
		"token": "dummy_token",
	})
}

// getUserInfo 获取用户信息
func getUserInfo(c *gin.Context) {
	// 获取用户ID
	userIDStr := c.Param("user_id")
	_, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// 返回模拟数据
	c.JSON(http.StatusOK, gin.H{"user": models.User{
		ID:           1,
		Username:     "testuser",
		Email:        "test@example.com",
		Role:         "user",
		CreatedAt:    time.Now(),
	}})
}

// 辅助函数：返回uint指针
func uintPtr(u uint) *uint {
	return &u
}
