package database

import (
	"log"
)

// DB 模拟数据库连接
var DB interface{}

// InitDB 初始化数据库连接
func InitDB() error {
	// 使用内存数据库模拟
	log.Println("Using in-memory database for development")
	
	// 这里使用一个模拟的数据库连接
	// 在实际生产环境中，应该使用真实的数据库
	DB = struct{}{}

	log.Println("Database connected successfully")
	return nil
}

// GetDB 获取数据库连接
func GetDB() interface{} {
	return DB
}
