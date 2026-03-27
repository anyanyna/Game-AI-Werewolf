package api

import (
	"werewolf-game/backend/websocket"

	"github.com/gin-gonic/gin"
)

// SetupRoutes 设置API路由
func SetupRoutes(router *gin.Engine, wsHandler *websocket.Handler) {
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


