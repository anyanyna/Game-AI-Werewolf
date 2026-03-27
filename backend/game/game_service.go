package game

import (
	"log"
	"math/rand"
	"time"

	"werewolf-game/backend/models"
)

// GameService 游戏服务
type GameService struct{}

// NewGameService 创建游戏服务实例
func NewGameService() *GameService {
	return &GameService{}
}

// CreateGame 创建新游戏
func (s *GameService) CreateGame(hostID uint) (*models.Game, error) {
	// 生成游戏代码
	gameCode := s.generateGameCode()

	// 创建游戏实例
	game := &models.Game{
		ID:        1,
		Status:    "waiting",
		GameCode:  gameCode,
		HostID:    hostID,
		CreatedAt: time.Now(),
	}

	log.Printf("Game created: %s, HostID: %d\n", gameCode, hostID)
	return game, nil
}

// JoinGame 加入游戏
func (s *GameService) JoinGame(gameCode string, userID uint) (*models.Game, error) {
	// 模拟加入游戏
	game := &models.Game{
		ID:        1,
		Status:    "waiting",
		GameCode:  gameCode,
		HostID:    1,
		CreatedAt: time.Now(),
	}

	log.Printf("User %d joined game %s\n", userID, gameCode)
	return game, nil
}

// StartGame 开始游戏
func (s *GameService) StartGame(gameID uint) error {
	// 模拟开始游戏
	log.Printf("Game %d started\n", gameID)
	return nil
}

// generateGameCode 生成游戏代码
func (s *GameService) generateGameCode() string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	code := make([]byte, 6)
	for i := range code {
		code[i] = charset[rand.Intn(len(charset))]
	}
	return string(code)
}

// addAIDigitalPersons 添加AI数字人填充
func (s *GameService) addAIDigitalPersons(gameID uint, realPlayerCount int) error {
	// 模拟添加AI数字人
	log.Printf("Added AI digital persons to game %d\n", gameID)
	return nil
}

// assignRoles 分配角色
func (s *GameService) assignRoles(gameID uint) error {
	// 模拟分配角色
	log.Printf("Assigned roles to game %d\n", gameID)
	return nil
}

// createGamePhase 创建游戏阶段
func (s *GameService) createGamePhase(gameID uint, phase string, round int) error {
	// 模拟创建游戏阶段
	log.Printf("Created game phase %s for game %d\n", phase, gameID)
	return nil
}
