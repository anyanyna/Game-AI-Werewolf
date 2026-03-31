package database

import (
	"fmt"
	"log"
	"os"

	"werewolf-game/backend/models"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB 全局数据库连接
var DB *gorm.DB

// 测试模式标志
var testDBPath string

// SetTestDB 设置测试数据库路径
func SetTestDB(path string) {
	testDBPath = path
}

// InitDB 初始化数据库连接
func InitDB() error {
	// 创建数据库文件
	dbPath := testDBPath
	if dbPath == "" {
		dbPath = "./werewolf.db"
	}

	// 检查文件是否存在
	_, err := os.Stat(dbPath)
	if os.IsNotExist(err) {
		log.Printf("Creating new database: %s", dbPath)
	} else {
		log.Printf("Using existing database: %s", dbPath)
	}

	// 连接数据库
	logLevel := logger.Info
	if testDBPath != "" {
		logLevel = logger.Silent // 测试模式静默日志
	}
	DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		return fmt.Errorf("failed to connect database: %w", err)
	}

	log.Println("Database connected successfully")
	return nil
}

// GetDB 获取数据库连接
func GetDB() *gorm.DB {
	return DB
}

// MigrateDB 迁移数据库表结构
func MigrateDB() error {
	if DB == nil {
		return fmt.Errorf("database not initialized")
	}

	// 自动迁移表结构
	err := DB.AutoMigrate(
		&models.User{},
		&models.Game{},
		&models.GamePlayer{},
		&models.AIDigitalPerson{},
		&models.GamePhase{},
		&models.NightAction{},
		&models.DayAction{},
		&models.RealPersonGuess{},
	)
	if err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	log.Println("Database migration completed")
	return nil
}

// CloseDB 关闭数据库连接
func CloseDB() error {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}
