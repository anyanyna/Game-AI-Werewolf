package main

import (
	"fmt"
	"log"
	"os"

	"werewolf-game/backend/ai"
	"werewolf-game/backend/api"
	"werewolf-game/backend/database"
	"werewolf-game/backend/game"
	"werewolf-game/backend/websocket"

	"github.com/gin-gonic/gin"
)

func main() {
	// 创建日志文件
	f, err := os.OpenFile("server.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		log.SetOutput(f)
	}

	fmt.Println("Starting Werewolf Game Server...")

	// 初始化数据库
	if err := database.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// 迁移数据库表结构
	if err := database.MigrateDB(); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// 初始化LangSmith追踪
	ai.InitGlobalTracer(ai.GetLangSmithConfig())

	// 初始化AI工厂
	ai.InitGlobalFactory()

	// 初始化游戏服务
	gameService := game.NewGameService(database.GetDB())

	// 初始化WebSocket客户端管理器
	clientManager := websocket.NewClientManager()
	go clientManager.Run()

	// 初始化WebSocket处理器
	wsHandler := websocket.NewHandler(clientManager, gameService)

	// 初始化Gin路由
	router := gin.Default()

	// 添加CORS中间件
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// 设置API路由
	api.SetupRoutes(router, wsHandler, gameService)

	// 启动服务器
	port := 8080
	log.Printf("Server starting on port %d...", port)
	if err := router.Run(fmt.Sprintf(":%d", port)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
