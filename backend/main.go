package main

import (
	"fmt"
	"log"

	"werewolf-game/backend/api"
	"werewolf-game/backend/database"
	"werewolf-game/backend/websocket"

	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("Starting Werewolf Game Server...")

	// 初始化数据库
	if err := database.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// 迁移数据库表结构
	if err := database.MigrateDB(); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// 初始化WebSocket客户端管理器
	clientManager := websocket.NewClientManager()
	go clientManager.Run()

	// 初始化WebSocket处理器
	wsHandler := websocket.NewHandler(clientManager)

	// 初始化Gin路由
	router := gin.Default()

	// 设置API路由
	api.SetupRoutes(router, wsHandler)

	// 启动服务器
	port := 8080
	log.Printf("Server starting on port %d...", port)
	if err := router.Run(fmt.Sprintf(":%d", port)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
