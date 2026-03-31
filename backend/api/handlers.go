package api

import (
	"net/http"
	"strconv"
	"strings"

	"werewolf-game/backend/game"
	"werewolf-game/backend/models"
	"werewolf-game/backend/websocket"

	"github.com/gin-gonic/gin"
)

// SetupRoutes 设置API路由
func SetupRoutes(router *gin.Engine, wsHandler *websocket.Handler, gameService *game.GameService) {
	// 使用中间件注入gameService
	router.Use(func(c *gin.Context) {
		c.Set("gameService", gameService)
		c.Next()
	})

	// API路由组
	api := router.Group("/api")
	{
		// 游戏相关API
		game := api.Group("/game")
		{
			game.POST("/create", createGame)
			game.POST("/join", joinGame)
			game.POST("/start", startGame)
			game.GET("/info/:game_id", getGameInfo)
			game.GET("/players/:game_id", getGamePlayers)
		}

		// 玩家相关API
		player := api.Group("/player")
		{
			player.GET("/info/:player_id", getPlayerInfo)
			player.POST("/action/night", nightAction)
			player.POST("/action/day", dayAction)
			player.POST("/action/vote", voteAction)
		}

		// 游戏阶段相关API
		phase := api.Group("/phase")
		{
			phase.POST("/end", endPhase)
			phase.GET("/current/:game_id", getCurrentPhase)
		}

		// 游戏日志相关API
		logs := api.Group("/logs")
		{
			logs.GET("/game/:game_id", getGameLogs)
			logs.GET("/player/:game_id/:player_id", getPlayerLogs)
		}

		// 真人猜测相关API
		guess := api.Group("/guess")
		{
			guess.POST("/real_person", guessRealPerson)
			guess.GET("/result/:game_id", getGuessResult)
		}

		// 用户相关API
		user := api.Group("/user")
		{
			user.POST("/register", registerUser)
			user.POST("/login", loginUser)
			user.GET("/info/:user_id", getUserInfo)
		}
	}

	// WebSocket路由
	router.GET("/ws", wsHandler.HandleConnection)
}

// getGameService 从上下文获取gameService
func getGameService(c *gin.Context) *game.GameService {
	if s, exists := c.Get("gameService"); exists {
		return s.(*game.GameService)
	}
	return nil
}

// createGame 创建游戏
func createGame(c *gin.Context) {
	var req struct {
		HostID      uint `json:"host_id" binding:"required"`
		PlayerCount int  `json:"player_count"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	playerCount := req.PlayerCount
	if playerCount <= 0 {
		playerCount = 12 // 默认12人局
	}

	gameService := getGameService(c)
	g, err := gameService.CreateGame(req.HostID, playerCount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"game":    g,
		"message": "Game created successfully",
	})
}

// joinGame 加入游戏
func joinGame(c *gin.Context) {
	var req struct {
		GameCode string `json:"game_code" binding:"required"`
		UserID   uint   `json:"user_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	gameService := getGameService(c)
	g, player, err := gameService.JoinGame(req.GameCode, req.UserID)
	if err != nil {
		// 检查游戏是否不存在
		errMsg := err.Error()
		if errMsg == "record not found" || strings.HasPrefix(errMsg, "game not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": "Game not found"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"game":    g,
		"player":  player,
		"message": "Joined game successfully",
	})
}

// startGame 开始游戏
func startGame(c *gin.Context) {
	var req struct {
		GameID uint `json:"game_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	gameService := getGameService(c)
	g, err := gameService.StartGame(req.GameID)
	if err != nil {
		// 检查是否是"游戏不存在"错误
		if err.Error() == "record not found" || err.Error() == "game not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Game not found"})
			return
		}
		// 检查是否游戏已经开始
		if err.Error() == "game already started" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"game":    g,
		"message": "Game started successfully",
	})
}

// getGameInfo 获取游戏信息
func getGameInfo(c *gin.Context) {
	gameIDStr := c.Param("game_id")
	gameID, err := strconv.ParseUint(gameIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid game ID"})
		return
	}

	gameService := getGameService(c)
	g, err := gameService.GetGame(uint(gameID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Game not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"game": g})
}

// getGamePlayers 获取游戏玩家
func getGamePlayers(c *gin.Context) {
	gameIDStr := c.Param("game_id")
	gameID, err := strconv.ParseUint(gameIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid game ID"})
		return
	}

	gameService := getGameService(c)

	// 检查是否有 viewer_player_id 参数
	viewerPlayerIDStr := c.Query("viewer_player_id")
	var players []models.GamePlayer

	if viewerPlayerIDStr != "" {
		viewerPlayerID, err := strconv.ParseUint(viewerPlayerIDStr, 10, 32)
		if err == nil {
			players, err = gameService.GetPlayersForViewer(uint(gameID), uint(viewerPlayerID))
		} else {
			players, err = gameService.GetPlayers(uint(gameID))
		}
	} else {
		players, err = gameService.GetPlayers(uint(gameID))
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"players": players})
}

// getPlayerInfo 获取玩家信息
func getPlayerInfo(c *gin.Context) {
	playerIDStr := c.Param("player_id")
	playerID, err := strconv.ParseUint(playerIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid player ID"})
		return
	}

	gameService := getGameService(c)
	player, err := gameService.GetPlayer(uint(playerID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Player not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"player": player})
}

// nightAction 夜晚行动
func nightAction(c *gin.Context) {
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

	gameService := getGameService(c)
	err := gameService.ProcessNightAction(req.GameID, req.PlayerID, req.ActionType, req.TargetID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Night action executed successfully"})
}

// dayAction 白天行动
func dayAction(c *gin.Context) {
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

	gameService := getGameService(c)
	err := gameService.ProcessDayAction(req.GameID, req.PlayerID, req.ActionType, req.Content, req.TargetID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Day action executed successfully"})
}

// voteAction 投票行动
func voteAction(c *gin.Context) {
	var req struct {
		GameID   uint `json:"game_id" binding:"required"`
		PlayerID uint `json:"player_id" binding:"required"`
		TargetID uint `json:"target_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// 投票也是白天行动的一种
	gameService := getGameService(c)
	err := gameService.ProcessDayAction(req.GameID, req.PlayerID, "vote", "", &req.TargetID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Vote action executed successfully"})
}

// endPhase 结束阶段
func endPhase(c *gin.Context) {
	var req struct {
		GameID uint `json:"game_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	gameService := getGameService(c)
	phase, err := gameService.EndPhase(req.GameID)
	if err != nil {
		// 检查游戏是否不存在
		if err.Error() == "record not found" || err.Error() == "game not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Game not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"phase": phase, "message": "Phase ended successfully"})
}

// getCurrentPhase 获取当前阶段
func getCurrentPhase(c *gin.Context) {
	gameIDStr := c.Param("game_id")
	gameID, err := strconv.ParseUint(gameIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid game ID"})
		return
	}

	gameService := getGameService(c)
	phase, err := gameService.GetCurrentPhase(uint(gameID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Phase not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"phase": phase})
}

// guessRealPerson 真人猜测
func guessRealPerson(c *gin.Context) {
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

	gameService := getGameService(c)
	guess, err := gameService.GuessRealPerson(req.GameID, req.UserID, req.TargetID, req.IsRealPerson)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"guess": guess})
}

// getGuessResult 获取猜测结果
func getGuessResult(c *gin.Context) {
	gameIDStr := c.Param("game_id")
	gameID, err := strconv.ParseUint(gameIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid game ID"})
		return
	}

	gameService := getGameService(c)
	guesses, err := gameService.GetGuesses(uint(gameID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"guesses": guesses})
}

// getGameLogs 获取游戏公共日志
func getGameLogs(c *gin.Context) {
	gameIDStr := c.Param("game_id")
	gameID, err := strconv.ParseUint(gameIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid game ID"})
		return
	}

	gameService := getGameService(c)
	logs, err := gameService.GetGameLogs(uint(gameID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, logs)
}

// getPlayerLogs 获取玩家私人日志
func getPlayerLogs(c *gin.Context) {
	gameIDStr := c.Param("game_id")
	gameID, err := strconv.ParseUint(gameIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid game ID"})
		return
	}

	playerIDStr := c.Param("player_id")
	playerID, err := strconv.ParseUint(playerIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid player ID"})
		return
	}

	gameService := getGameService(c)
	logs, err := gameService.GetPlayerLogs(uint(gameID), uint(playerID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, logs)
}

// registerUser 注册用户
func registerUser(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	gameService := getGameService(c)
	user, err := gameService.RegisterUser(req.Username, req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user":    user,
		"message": "User registered successfully",
	})
}

// loginUser 登录用户
func loginUser(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	gameService := getGameService(c)
	user, err := gameService.VerifyLogin(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user":  user,
		"token": "dummy_token",
	})
}

// getUserInfo 获取用户信息
func getUserInfo(c *gin.Context) {
	userIDStr := c.Param("user_id")
	_, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// TODO: 实现真实的用户信息获取
	c.JSON(http.StatusOK, gin.H{
		"user": models.User{
			ID:       1,
			Username: "testuser",
			Email:    "test@example.com",
			Role:     "user",
		},
	})
}
