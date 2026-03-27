package models

import (
	"time"
)

// User 用户模型
type User struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	Username     string    `json:"username" gorm:"uniqueIndex"`
	Email        string    `json:"email" gorm:"uniqueIndex"`
	PasswordHash string    `json:"-" gorm:"not null"`
	CreatedAt    time.Time `json:"created_at"`
	LastLogin    *time.Time `json:"last_login"`
	Role         string    `json:"role" gorm:"default:user"`
}

// Game 游戏模型
type Game struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Status    string    `json:"status" gorm:"not null"` // waiting, playing, ended
	GameCode  string    `json:"game_code" gorm:"uniqueIndex;not null"`
	HostID    uint      `json:"host_id"`
	Winner    string    `json:"winner"`
	CreatedAt time.Time `json:"created_at"`
	StartedAt *time.Time `json:"started_at"`
	EndedAt   *time.Time `json:"ended_at"`
}

// GamePlayer 游戏玩家模型
type GamePlayer struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	GameID        uint      `json:"game_id" gorm:"index"`
	UserID        *uint     `json:"user_id"`
	AIDigitalPersonID *uint `json:"ai_digital_person_id"`
	Role          string    `json:"role" gorm:"not null"` // werewolf, villager, seer, witch, hunter
	Status        string    `json:"status" gorm:"not null"` // alive, dead
	IsRealPerson  bool      `json:"is_real_person"`
	CreatedAt     time.Time `json:"created_at"`
	Number        int       `json:"number" gorm:"not null"` // 角色数字代号
}

// AIDigitalPerson AI数字人模型
type AIDigitalPerson struct {
	ID             uint      `json:"id" gorm:"primaryKey"`
	Name           string    `json:"name" gorm:"not null"`
	Personality    string    `json:"personality"`
	Experience     int       `json:"experience" gorm:"default:0"`
	Intelligence   float64   `json:"intelligence" gorm:"default:0.5"`
	CreatedAt      time.Time `json:"created_at"`
	LastPlayedAt   *time.Time `json:"last_played_at"`
	SessionID      string    `json:"session_id" gorm:"uniqueIndex"`
}

// GamePhase 游戏阶段模型
type GamePhase struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	GameID    uint      `json:"game_id" gorm:"index"`
	Phase     string    `json:"phase" gorm:"not null"` // night, day, voting
	Round     int       `json:"round" gorm:"not null"`
	StartTime time.Time `json:"start_time"`
	EndTime   *time.Time `json:"end_time"`
}

// NightAction 夜晚行动模型
type NightAction struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	GameID       uint      `json:"game_id" gorm:"index"`
	GamePhaseID  uint      `json:"game_phase_id" gorm:"index"`
	PlayerID     uint      `json:"player_id" gorm:"index"`
	ActionType   string    `json:"action_type" gorm:"not null"` // kill, save, poison, check
	TargetID     uint      `json:"target_id"`
	Result       string    `json:"result"`
	CreatedAt    time.Time `json:"created_at"`
}

// DayAction 白天行动模型
type DayAction struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	GameID       uint      `json:"game_id" gorm:"index"`
	GamePhaseID  uint      `json:"game_phase_id" gorm:"index"`
	PlayerID     uint      `json:"player_id" gorm:"index"`
	ActionType   string    `json:"action_type" gorm:"not null"` // speak, vote
	Content      string    `json:"content"`
	TargetID     *uint     `json:"target_id"`
	CreatedAt    time.Time `json:"created_at"`
}

// RealPersonGuess 真人猜测模型
type RealPersonGuess struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	GameID       uint      `json:"game_id" gorm:"index"`
	UserID       uint      `json:"user_id" gorm:"index"`
	TargetID     uint      `json:"target_id" gorm:"index"` // 被猜测的玩家ID
	IsRealPerson bool      `json:"is_real_person"` // 猜测是否为真人
	IsCorrect    bool      `json:"is_correct"` // 猜测是否正确
	CreatedAt    time.Time `json:"created_at"`
}
