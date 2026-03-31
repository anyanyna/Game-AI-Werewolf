package game

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"werewolf-game/backend/ai"
	"werewolf-game/backend/models"
	"werewolf-game/backend/repository"

	"gorm.io/gorm"
)

// GameService 游戏服务
type GameService struct {
	gameRepo  *repository.GameRepository
	userRepo  *repository.UserRepository
	aiRepo    *repository.AIDigitalPersonRepository
	aiFactory *ai.AIFactory
	rand      *rand.Rand
}

// NewGameService 创建游戏服务实例
func NewGameService(db *gorm.DB) *GameService {
	return &GameService{
		gameRepo:  repository.NewGameRepository(db),
		userRepo:  repository.NewUserRepository(db),
		aiRepo:    repository.NewAIDigitalPersonRepository(db),
		aiFactory: ai.GetGlobalFactory(),
		rand:      rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// CreateGame 创建新游戏
func (s *GameService) CreateGame(hostID uint, playerCount int) (*models.Game, error) {
	// 生成游戏代码
	gameCode := s.generateGameCode()

	// 确保游戏代码唯一
	_, err := s.gameRepo.GetByGameCode(gameCode)
	if err == nil {
		// 代码已存在，重新生成
		gameCode = s.generateGameCode()
	}

	// 创建游戏
	game := &models.Game{
		Status:    "waiting",
		GameCode:  gameCode,
		HostID:    hostID,
		Winner:    "",
		CreatedAt: time.Now(),
	}

	if err := s.gameRepo.Create(game); err != nil {
		return nil, err
	}

	// 为主人创建玩家记录
	hostPlayer := &models.GamePlayer{
		GameID:            game.ID,
		UserID:            &hostID,
		AIDigitalPersonID: nil,
		Role:              "",
		Status:            "alive",
		IsRealPerson:      true,
		Number:            1,
		CreatedAt:         time.Now(),
	}

	if err := s.gameRepo.AddPlayer(hostPlayer); err != nil {
		return nil, err
	}

	return game, nil
}

// JoinGame 加入游戏
func (s *GameService) JoinGame(gameCode string, userID uint) (*models.Game, *models.GamePlayer, error) {
	// 获取游戏
	game, err := s.gameRepo.GetByGameCode(gameCode)
	if err != nil {
		return nil, nil, fmt.Errorf("game not found: %s", gameCode)
	}

	if game.Status != "waiting" {
		return nil, nil, fmt.Errorf("game already started")
	}

	// 检查是否已加入
	existingPlayer, _ := s.gameRepo.GetPlayerByUserIDAndGameID(userID, game.ID)
	if existingPlayer != nil {
		return game, existingPlayer, nil
	}

	// 获取当前玩家数量
	players, _ := s.gameRepo.GetPlayers(game.ID)
	playerNumber := len(players) + 1

	// 创建玩家记录
	player := &models.GamePlayer{
		GameID:            game.ID,
		UserID:            &userID,
		AIDigitalPersonID: nil,
		Role:              "",
		Status:            "alive",
		IsRealPerson:      true,
		Number:            playerNumber,
		CreatedAt:         time.Now(),
	}

	if err := s.gameRepo.AddPlayer(player); err != nil {
		return nil, nil, err
	}

	return game, player, nil
}

// StartGame 开始游戏
func (s *GameService) StartGame(gameID uint) (*models.Game, error) {
	game, err := s.gameRepo.GetByID(gameID)
	if err != nil {
		return nil, err
	}

	if game.Status != "waiting" {
		return nil, fmt.Errorf("game already started")
	}

	players, err := s.gameRepo.GetPlayers(gameID)
	if err != nil {
		return nil, err
	}

	playerCount := len(players)
	if playerCount < 1 {
		return nil, fmt.Errorf("not enough players")
	}

	// 填充AI数字人到12人
	aiCount := 12 - playerCount
	if aiCount > 0 {
		// 尝试从数据库获取AI数字人，如果没有则创建
		aiPersons, err := s.getOrCreateAIDigitalPersons(aiCount)
		if err != nil {
			return nil, err
		}

		for i, aiPerson := range aiPersons {
			player := &models.GamePlayer{
				GameID:            game.ID,
				UserID:            nil,
				AIDigitalPersonID: &aiPerson.ID,
				Role:              "",
				Status:            "alive",
				IsRealPerson:      false,
				Number:            playerCount + i + 1,
				CreatedAt:         time.Now(),
			}
			if err := s.gameRepo.AddPlayer(player); err != nil {
				return nil, err
			}
		}
	}

	// 重新获取所有玩家
	players, err = s.gameRepo.GetPlayers(gameID)
	if err != nil {
		return nil, err
	}

	// 分配角色
	if err := s.assignRoles(gameID, players); err != nil {
		return nil, err
	}

	// 更新游戏状态
	game.Status = "playing"
	now := time.Now()
	game.StartedAt = &now
	if err := s.gameRepo.Update(game); err != nil {
		return nil, err
	}

	// 创建第一个夜晚阶段
	phase := &models.GamePhase{
		GameID:    gameID,
		Phase:     "night",
		Round:     1,
		StartTime: time.Now(),
	}
	log.Printf("[DEBUG] Creating phase for game %d", gameID)

	// 直接使用 GORM 创建
	result := s.gameRepo.DB().Create(phase)
	if result.Error != nil {
		log.Printf("[DEBUG] CreatePhase ERROR: %v, Phase=%+v", result.Error, phase)
		return nil, result.Error
	}
	log.Printf("[DEBUG] Phase created successfully, ID=%d, GameID=%d", phase.ID, phase.GameID)

	return game, nil
}

// assignRoles 分配角色
func (s *GameService) assignRoles(gameID uint, players []models.GamePlayer) error {
	// 标准12人局角色分配
	roles := []string{
		"werewolf", "werewolf", "werewolf", "werewolf", // 4狼
		"seer",                                         // 预言家
		"witch",                                        // 女巫
		"hunter",                                       // 猎人
		"villager", "villager", "villager", "villager", // 4平民
	}

	// 打乱角色顺序
	s.rand.Shuffle(len(roles), func(i, j int) {
		roles[i], roles[j] = roles[j], roles[i]
	})

	// 分配角色
	for i, player := range players {
		if i < len(roles) {
			player.RoleRaw = roles[i]
			if err := s.gameRepo.UpdatePlayer(&player); err != nil {
				return err
			}
		}
	}

	// 初始化女巫的解药和毒药
	for _, player := range players {
		if player.RoleRaw == "witch" {
			phase, err := s.gameRepo.GetCurrentPhase(gameID)
			if err == nil {
				phase.WitchHasSave = true
				phase.WitchHasPoison = true
				s.gameRepo.UpdatePhase(phase)
			}
			break
		}
	}

	return nil
}

// getOrCreateAIDigitalPersons 获取或创建AI数字人
func (s *GameService) getOrCreateAIDigitalPersons(count int) ([]models.AIDigitalPerson, error) {
	// 尝试从数据库获取
	persons, err := s.aiRepo.GetRandom(count)
	if err != nil || len(persons) < count {
		// 创建新的AI数字人
		persons = make([]models.AIDigitalPerson, count)
		personalities := []string{"outgoing", "introverted", "logical", "intuitive"}

		for i := 0; i < count; i++ {
			person := models.AIDigitalPerson{
				Name:         fmt.Sprintf("AI_%d", time.Now().UnixNano()+int64(i)),
				Personality:  personalities[s.rand.Intn(len(personalities))],
				Experience:   s.rand.Intn(100),
				Intelligence: 0.3 + s.rand.Float64()*0.6, // 0.3-0.9
				CreatedAt:    time.Now(),
				SessionID:    fmt.Sprintf("session_%d", time.Now().UnixNano()+int64(i)),
			}
			if err := s.aiRepo.Create(&person); err != nil {
				return nil, err
			}
			persons[i] = person
		}
	}

	return persons, nil
}

// generateGameCode 生成游戏代码
func (s *GameService) generateGameCode() string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	code := make([]byte, 6)
	for i := range code {
		code[i] = charset[s.rand.Intn(len(charset))]
	}
	return string(code)
}

// GetGame 获取游戏信息
func (s *GameService) GetGame(gameID uint) (*models.Game, error) {
	return s.gameRepo.GetByID(gameID)
}

// GetGameByCode 根据游戏代码获取游戏
func (s *GameService) GetGameByCode(gameCode string) (*models.Game, error) {
	return s.gameRepo.GetByGameCode(gameCode)
}

// GetPlayers 获取游戏玩家（根据查看者角色过滤敏感信息）
func (s *GameService) GetPlayers(gameID uint) ([]models.GamePlayer, error) {
	return s.gameRepo.GetPlayers(gameID)
}

// GetPlayersForViewer 根据查看者角色过滤玩家信息
func (s *GameService) GetPlayersForViewer(gameID uint, viewerPlayerID uint) ([]models.GamePlayer, error) {
	players, err := s.gameRepo.GetPlayers(gameID)
	if err != nil {
		return nil, err
	}

	// 获取查看者的角色
	viewerPlayer, err := s.gameRepo.GetPlayerByID(viewerPlayerID)
	if err != nil {
		// 如果找不到查看者，返回原始数据
		return players, nil
	}

	// 过滤每个玩家的敏感信息
	for i := range players {
		players[i].Role = s.filterRoleForViewer(&players[i], viewerPlayer)
	}

	return players, nil
}

// filterRoleForViewer 根据查看者角色过滤被查看者的角色
func (s *GameService) filterRoleForViewer(target, viewer *models.GamePlayer) string {
	// 自己永远可以看到自己的真实角色
	if target.ID == viewer.ID {
		return target.RoleRaw
	}

	// 死者或游戏结束后所有人都可以看到真实角色
	if target.Status == "dead" {
		return target.RoleRaw
	}

	// 游戏已结束
	game, _ := s.gameRepo.GetByID(target.GameID)
	if game != nil && game.Status == "ended" {
		return target.RoleRaw
	}

	// 查看者是狼人
	if viewer.RoleRaw == "werewolf" {
		// 狼人可以互相看到
		if target.RoleRaw == "werewolf" {
			return "werewolf"
		}
		// 狼人看其他任何人都是隐藏的
		return "unknown"
	}

	// 其他所有情况，角色都是隐藏的
	return "unknown"
}

// GetPlayer 获取玩家信息
func (s *GameService) GetPlayer(playerID uint) (*models.GamePlayer, error) {
	return s.gameRepo.GetPlayerByID(playerID)
}

// ProcessNightAction 处理夜晚行动
func (s *GameService) ProcessNightAction(gameID, playerID uint, actionType string, targetID uint) error {
	player, err := s.gameRepo.GetPlayerByID(playerID)
	if err != nil {
		return err
	}

	// 获取当前阶段
	phase, err := s.gameRepo.GetCurrentPhase(gameID)
	if err != nil {
		return err
	}

	// 记录夜晚行动
	action := &models.NightAction{
		GameID:      gameID,
		GamePhaseID: phase.ID,
		PlayerID:    playerID,
		ActionType:  actionType,
		TargetID:    targetID,
		CreatedAt:   time.Now(),
	}

	if err := s.gameRepo.CreateNightAction(action); err != nil {
		return err
	}

	// 根据角色类型处理行动
	switch player.RoleRaw {
	case "werewolf":
		// 狼人杀人
		phase.KilledPlayer = &targetID
	case "seer":
		// 预言家验人 - 记录查验结果
		targetPlayer, _ := s.gameRepo.GetPlayerByID(targetID)
		if targetPlayer != nil {
			phase.CheckedPlayer = &targetID
			// 根据目标角色返回查验结果
			if targetPlayer.RoleRaw == "werewolf" {
				phase.CheckedResult = "werewolf"
			} else {
				phase.CheckedResult = "good"
			}
		}
	case "witch":
		// 女巫用药 - 检查药水是否还有
		if actionType == "save" && !phase.WitchHasSave {
			return fmt.Errorf("解药已用尽")
		}
		if actionType == "poison" && !phase.WitchHasPoison {
			return fmt.Errorf("毒药已用尽")
		}

		if actionType == "save" {
			phase.SavedPlayer = &targetID
			phase.WitchHasSave = false // 使用了解药
		} else if actionType == "poison" {
			phase.PoisonedPlayer = &targetID
			phase.WitchHasPoison = false // 使用了毒药
		}
	}

	return s.gameRepo.UpdatePhase(phase)
}

// ProcessDayAction 处理白天行动
func (s *GameService) ProcessDayAction(gameID, playerID uint, actionType, content string, targetID *uint) error {
	phase, err := s.gameRepo.GetCurrentPhase(gameID)
	if err != nil {
		return err
	}

	action := &models.DayAction{
		GameID:      gameID,
		GamePhaseID: phase.ID,
		PlayerID:    playerID,
		ActionType:  actionType,
		Content:     content,
		TargetID:    targetID,
		CreatedAt:   time.Now(),
	}

	return s.gameRepo.CreateDayAction(action)
}

// EndPhase 结束当前阶段
func (s *GameService) EndPhase(gameID uint) (*models.GamePhase, error) {
	game, err := s.gameRepo.GetByID(gameID)
	if err != nil {
		return nil, err
	}

	phase, err := s.gameRepo.GetCurrentPhase(gameID)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	phase.EndTime = &now

	// 根据当前阶段类型处理
	if phase.Phase == "night" {
		// 夜晚结束，处理死亡逻辑
		// 检查被毒死的玩家
		if phase.PoisonedPlayer != nil {
			if err := s.killPlayer(*phase.PoisonedPlayer); err != nil {
				return nil, err
			}
		}
		// 检查被狼人杀死的玩家（如果女巫没救）
		if phase.KilledPlayer != nil && phase.SavedPlayer == nil {
			if err := s.killPlayer(*phase.KilledPlayer); err != nil {
				return nil, err
			}
		}

		// 进入白天
		phase.Phase = "day"
		// 触发AI白天发言（等待玩家发言后）
		go func() {
			time.Sleep(3 * time.Second)
			s.triggerAIActions(gameID)
		}()
	} else if phase.Phase == "day" {
		// 白天结束，进入投票阶段
		phase.Phase = "voting"
		phase.StartTime = time.Now()
		phase.EndTime = nil
		// 触发AI投票行动
		go s.triggerAIActions(gameID)
	} else if phase.Phase == "voting" {
		// 投票阶段结束，计算投票结果
		// 获取所有投票
		dayActions, err := s.gameRepo.GetDayActions(gameID, phase.ID)
		if err != nil {
			return nil, err
		}

		// 统计投票
		voteCount := make(map[uint]int)
		var voters []uint
		for _, action := range dayActions {
			if action.ActionType == "vote" && action.TargetID != nil {
				voteCount[*action.TargetID]++
				voters = append(voters, action.PlayerID)
			}
		}

		// 找出最高票数
		var maxVotes int
		var votedPlayer *uint
		for targetID, count := range voteCount {
			if count > maxVotes {
				maxVotes = count
				votedPlayer = &targetID
			}
		}

		// 如果有平票，暂不平票（简化处理：第一个最高票当选）
		// 放逐玩家
		if votedPlayer != nil && maxVotes > 0 {
			// 检查是否是猎人
			votedPlayerInfo, _ := s.gameRepo.GetPlayerByID(*votedPlayer)
			if votedPlayerInfo != nil && votedPlayerInfo.RoleRaw == "hunter" {
				// 猎人被放逐，需要等待猎人选择带走谁
				// 暂时标记，稍后处理
				phase.HunterVotedOut = votedPlayer
			}

			// 执行放逐
			if err := s.killPlayer(*votedPlayer); err != nil {
				return nil, err
			}
			phase.VotedPlayer = votedPlayer
		}

		// 进入下一个夜晚
		phase.Round++
		phase.Phase = "night"
		phase.StartTime = time.Now()
		phase.EndTime = nil
		// 清除上一轮的数据
		phase.KilledPlayer = nil
		phase.SavedPlayer = nil
		phase.PoisonedPlayer = nil
		phase.CheckedPlayer = nil
		phase.CheckedResult = ""
		phase.HunterVotedOut = nil

		// 触发AI夜晚行动（下一轮的夜晚）
		go s.triggerAIActions(gameID)
	}

	if err := s.gameRepo.UpdatePhase(phase); err != nil {
		return nil, err
	}

	// 检查游戏是否结束
	if s.checkGameEnd(gameID) {
		game.Status = "ended"
		s.determineWinner(game)
		s.gameRepo.Update(game)
	}

	return phase, nil
}

// killPlayer 处死玩家
func (s *GameService) killPlayer(playerID uint) error {
	player, err := s.gameRepo.GetPlayerByID(playerID)
	if err != nil {
		return err
	}
	player.Status = "dead"
	return s.gameRepo.UpdatePlayer(player)
}

// checkGameEnd 检查游戏是否结束
func (s *GameService) checkGameEnd(gameID uint) bool {
	players, err := s.gameRepo.GetPlayers(gameID)
	if err != nil {
		return false
	}

	var werewolfCount, goodCount int
	for _, p := range players {
		if p.Status != "alive" {
			continue
		}
		if p.Role == "werewolf" {
			werewolfCount++
		} else {
			goodCount++
		}
	}

	// 狼人全灭 = 好人胜利
	if werewolfCount == 0 {
		return true
	}
	// 狼人数量 >= 好人数 = 狼人胜利
	if werewolfCount >= goodCount {
		return true
	}

	return false
}

// determineWinner 确定获胜方
func (s *GameService) determineWinner(game *models.Game) {
	players, _ := s.gameRepo.GetPlayers(game.ID)

	var werewolfCount, goodCount int
	for _, p := range players {
		if p.Status != "alive" {
			continue
		}
		if p.Role == "werewolf" {
			werewolfCount++
		} else {
			goodCount++
		}
	}

	if werewolfCount == 0 || werewolfCount < goodCount {
		game.Winner = "good"
	} else {
		game.Winner = "werewolf"
	}
}

// GuessRealPerson 真人猜测
func (s *GameService) GuessRealPerson(gameID, userID, targetID uint, isRealPerson bool) (*models.RealPersonGuess, error) {
	guess := &models.RealPersonGuess{
		GameID:       gameID,
		UserID:       userID,
		TargetID:     targetID,
		IsRealPerson: isRealPerson,
		IsCorrect:    false,
		CreatedAt:    time.Now(),
	}

	// 获取被猜测的玩家，判断是否正确
	targetPlayer, err := s.gameRepo.GetPlayerByID(targetID)
	if err == nil {
		guess.IsCorrect = (targetPlayer.IsRealPerson == isRealPerson)
	}

	if err := s.gameRepo.CreateGuess(guess); err != nil {
		return nil, err
	}

	return guess, nil
}

// GetGuesses 获取游戏猜测结果
func (s *GameService) GetGuesses(gameID uint) ([]models.RealPersonGuess, error) {
	return s.gameRepo.GetGuessesByGameID(gameID)
}

// GetCurrentPhase 获取当前阶段
func (s *GameService) GetCurrentPhase(gameID uint) (*models.GamePhase, error) {
	return s.gameRepo.GetCurrentPhase(gameID)
}

// VerifyLogin 验证用户登录
func (s *GameService) VerifyLogin(username, password string) (*models.User, error) {
	user, err := s.userRepo.GetByUsername(username)
	if err != nil {
		return nil, fmt.Errorf("invalid username or password")
	}

	if !s.userRepo.VerifyPassword(user, password) {
		return nil, fmt.Errorf("invalid username or password")
	}

	return user, nil
}

// RegisterUser 注册用户
func (s *GameService) RegisterUser(username, email, password string) (*models.User, error) {
	// 检查用户名是否已存在
	_, err := s.userRepo.GetByUsername(username)
	if err == nil {
		return nil, fmt.Errorf("username already exists")
	}

	// 检查邮箱是否已存在
	_, err = s.userRepo.GetByEmail(email)
	if err == nil {
		return nil, fmt.Errorf("email already exists")
	}

	// 创建用户
	user := &models.User{
		Username:     username,
		Email:        email,
		PasswordHash: password, // TODO: 使用bcrypt加密
		Role:         "user",
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, fmt.Errorf("failed to create user: %v", err)
	}

	return user, nil
}

// GetGameLogs 获取游戏日志（所有阶段、夜晚行动、白天行动）
func (s *GameService) GetGameLogs(gameID uint) (map[string]interface{}, error) {
	// 获取所有阶段
	phases, err := s.gameRepo.GetPhasesByGameID(gameID)
	if err != nil {
		return nil, err
	}

	// 获取所有玩家信息，用于显示玩家编号
	players, err := s.gameRepo.GetPlayers(gameID)
	if err != nil {
		players = []models.GamePlayer{}
	}

	// 构建玩家ID到编号的映射
	playerNumMap := make(map[uint]int)
	for _, p := range players {
		playerNumMap[p.ID] = p.Number
	}

	// 为每个阶段获取夜晚行动和白天行动
	typeLog := make([]map[string]interface{}, 0)
	for _, phase := range phases {
		phaseLog := map[string]interface{}{
			"round":      phase.Round,
			"phase":      phase.Phase,
			"start_time": phase.StartTime.Format("15:04"),
		}

		// 添加夜晚行动
		nightActions, _ := s.gameRepo.GetNightActions(gameID, phase.ID)
		if len(nightActions) > 0 {
			nightLogs := make([]map[string]interface{}, 0)
			for _, action := range nightActions {
				playerNum := playerNumMap[action.PlayerID]
				targetNum := playerNumMap[action.TargetID]
				nightLogs = append(nightLogs, map[string]interface{}{
					"player_num":  playerNum,
					"action_type": action.ActionType,
					"target_num":  targetNum,
				})
			}
			phaseLog["night_actions"] = nightLogs
		}

		// 添加白天行动
		dayActions, _ := s.gameRepo.GetDayActions(gameID, phase.ID)
		if len(dayActions) > 0 {
			dayLogs := make([]map[string]interface{}, 0)
			for _, action := range dayActions {
				playerNum := playerNumMap[action.PlayerID]
				dayLog := map[string]interface{}{
					"player_num":  playerNum,
					"action_type": action.ActionType,
					"content":     action.Content,
				}
				if action.TargetID != nil {
					targetNum := playerNumMap[*action.TargetID]
					dayLog["target_num"] = targetNum
				}
				dayLogs = append(dayLogs, dayLog)
			}
			phaseLog["day_actions"] = dayLogs
		}

		// 添加阶段结果（死亡信息等）
		if phase.KilledPlayer != nil && phase.SavedPlayer == nil {
			phaseLog["killed"] = playerNumMap[*phase.KilledPlayer]
		}
		if phase.PoisonedPlayer != nil {
			phaseLog["poisoned"] = playerNumMap[*phase.PoisonedPlayer]
		}
		if phase.VotedPlayer != nil {
			phaseLog["voted_out"] = playerNumMap[*phase.VotedPlayer]
		}

		typeLog = append(typeLog, phaseLog)
	}

	return map[string]interface{}{
		"phases":  typeLog,
		"players": players,
	}, nil
}

// GetPlayerLogs 获取特定玩家的私人日志（仅夜晚行动）
func (s *GameService) GetPlayerLogs(gameID, playerID uint) (map[string]interface{}, error) {
	// 获取玩家信息
	player, err := s.gameRepo.GetPlayerByID(playerID)
	if err != nil {
		return nil, err
	}

	// 获取所有阶段
	phases, err := s.gameRepo.GetPhasesByGameID(gameID)
	if err != nil {
		return nil, err
	}

	// 获取所有玩家信息
	players, err := s.gameRepo.GetPlayers(gameID)
	if err != nil {
		players = []models.GamePlayer{}
	}

	// 构建玩家ID到编号的映射
	playerNumMap := make(map[uint]int)
	for _, p := range players {
		playerNumMap[p.ID] = p.Number
	}

	// 根据角色类型构建私人日志
	privateLogs := make([]map[string]interface{}, 0)

	for _, phase := range phases {
		// 只有夜晚阶段才有私人日志
		if phase.Phase != "night" {
			continue
		}

		phaseLog := map[string]interface{}{
			"round": phase.Round,
			"phase": "night",
		}

		switch player.Role {
		case "werewolf":
			// 狼人可以看到本轮其他狼人的杀人投票
			nightActions, _ := s.gameRepo.GetNightActions(gameID, phase.ID)
			wolfActions := make([]map[string]interface{}, 0)
			for _, action := range nightActions {
				if action.ActionType == "kill" {
					// 获取行动的玩家信息
					actionPlayer, _ := s.gameRepo.GetPlayerByID(action.PlayerID)
					if actionPlayer != nil && actionPlayer.RoleRaw == "werewolf" {
						wolfActions = append(wolfActions, map[string]interface{}{
							"player_num": playerNumMap[action.PlayerID],
							"target_num": playerNumMap[action.TargetID],
						})
					}
				}
			}
			if len(wolfActions) > 0 {
				phaseLog["wolf_kill_votes"] = wolfActions
			}

		case "seer":
			// 预言家可以看到查验结果
			if phase.CheckedPlayer != nil {
				phaseLog["checked_player"] = playerNumMap[*phase.CheckedPlayer]
				phaseLog["checked_result"] = phase.CheckedResult
			}

		case "witch":
			// 女巫可以看到救人和毒人信息
			if phase.SavedPlayer != nil {
				phaseLog["saved_player"] = playerNumMap[*phase.SavedPlayer]
			}
			if phase.PoisonedPlayer != nil {
				phaseLog["poisoned_player"] = playerNumMap[*phase.PoisonedPlayer]
			}

		case "hunter":
			// 猎人可以看到自己被杀死的信息（如果有夜间死亡）
			// 简化处理
		}

		// 只有当有私人日志时才添加
		if len(phaseLog) > 2 {
			privateLogs = append(privateLogs, phaseLog)
		}
	}

	return map[string]interface{}{
		"role":          player.Role,
		"player_number": player.Number,
		"private_logs":  privateLogs,
	}, nil
}

// triggerAIActions 自动触发AI行动
func (s *GameService) triggerAIActions(gameID uint) {
	log.Printf("[AI] Triggering AI actions for game %d", gameID)

	// 给AI一些准备时间
	time.Sleep(2 * time.Second)

	players, err := s.gameRepo.GetPlayers(gameID)
	if err != nil {
		log.Printf("[AI] Failed to get players: %v", err)
		return
	}

	phase, err := s.gameRepo.GetCurrentPhase(gameID)
	if err != nil {
		log.Printf("[AI] Failed to get phase: %v", err)
		return
	}

	log.Printf("[AI] Current phase: %s, Round: %d, Players: %d", phase.Phase, phase.Round, len(players))

	// 为每个AI玩家触发行动
	for i := range players {
		player := &players[i]

		// 跳过真人玩家和已死亡的玩家
		if player.IsRealPerson || player.Status != "alive" {
			continue
		}

		log.Printf("[AI] Processing AI player %d (role: %s)", player.ID, player.RoleRaw)

		// 获取AI provider
		provider, err := s.aiFactory.GetProvider("llm")
		if err != nil {
			// fallback to behavior tree
			provider, _ = s.aiFactory.GetProvider("behavior_tree")
		}

		if provider == nil {
			log.Printf("[AI] No provider found for player %d", player.ID)
			continue
		}

		// 初始化AI
		if aiPersonID := player.AIDigitalPersonID; aiPersonID != nil {
			aiPerson, _ := s.aiRepo.GetByID(*aiPersonID)
			if aiPerson != nil {
				provider.InitAI(aiPerson)
			}
		}

		// 构建AI上下文
		ctx := &ai.AIGameContext{
			GameID:       gameID,
			Round:        phase.Round,
			Phase:        phase.Phase,
			Player:       player,
			AllPlayers:   convertToAIPlayers(players),
			AlivePlayers: getAliveAIPlayers(players),
			GameInfo:     nil,
			RoleInfo:     &ai.RoleInfo{MyRole: player.RoleRaw},
		}

		// 根据阶段触发不同行动
		switch phase.Phase {
		case "night":
			s.triggerAINightAction(provider, ctx, gameID, player)
		case "day":
			s.triggerAIDayAction(provider, ctx, gameID, player)
		case "voting":
			s.triggerAIVote(provider, ctx, gameID, player)
		}
	}
}

// convertToAIPlayers 转换为AI玩家切片
func convertToAIPlayers(players []models.GamePlayer) []*models.GamePlayer {
	result := make([]*models.GamePlayer, len(players))
	for i := range players {
		result[i] = &players[i]
	}
	return result
}

// getAliveAIPlayers 获取存活AI玩家
func getAliveAIPlayers(players []models.GamePlayer) []*models.GamePlayer {
	var result []*models.GamePlayer
	for i := range players {
		if players[i].Status == "alive" {
			result = append(result, &players[i])
		}
	}
	return result
}

// triggerAINightAction 触发AI夜晚行动
func (s *GameService) triggerAINightAction(provider ai.AIProvider, ctx *ai.AIGameContext, gameID uint, player *models.GamePlayer) {
	actionType, targetID, err := provider.ProcessNightAction(ctx)
	if err != nil || actionType == "no_action" || targetID == 0 {
		return
	}

	// 执行行动
	err = s.ProcessNightAction(gameID, player.ID, actionType, targetID)
	if err != nil {
		log.Printf("[AI] Player %d night action failed: %v", player.ID, err)
	}
}

// triggerAIDayAction 触发AI白天行动
func (s *GameService) triggerAIDayAction(provider ai.AIProvider, ctx *ai.AIGameContext, gameID uint, player *models.GamePlayer) {
	actionType, speech, _, err := provider.ProcessDayAction(ctx)
	if err != nil || actionType != "speak" {
		return
	}

	// 执行发言
	err = s.ProcessDayAction(gameID, player.ID, "speak", speech, nil)
	if err != nil {
		log.Printf("[AI] Player %d day action failed: %v", player.ID, err)
	}
}

// triggerAIVote 触发AI投票
func (s *GameService) triggerAIVote(provider ai.AIProvider, ctx *ai.AIGameContext, gameID uint, player *models.GamePlayer) {
	targetID, err := provider.ProcessVote(ctx)
	if err != nil || targetID == 0 {
		return
	}

	// 执行投票
	err = s.ProcessDayAction(gameID, player.ID, "vote", "", &targetID)
	if err != nil {
		log.Printf("[AI] Player %d vote failed: %v", player.ID, err)
	}
}
