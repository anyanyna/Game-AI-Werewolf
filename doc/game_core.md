# 游戏核心逻辑模块详细设计文档

## 1. 模块概述

游戏核心逻辑模块是整个游戏的核心，负责处理游戏规则、流程控制、角色分配等核心功能。该模块确保游戏按照标准狼人杀规则运行，并处理游戏中的各种事件和状态转换。

## 2. 功能需求

### 2.1 游戏配置
- **人数配置**：12人标准局
- **角色配置**：4狼人、4神民（预言家、女巫、猎人、守卫）、4平民
- **角色数字代号**：1-12

### 2.2 游戏流程
1. **准备阶段**：生成游戏ID、分配角色、填充AI数字人
2. **夜晚阶段**：狼人杀人、预言家验人、女巫用药、守卫守护
3. **白天阶段**：宣布夜晚结果、玩家发言、投票放逐
4. **重复阶段**：夜晚-白天循环，直到某一方获胜
5. **结束阶段**：宣布获胜方、展示评分、真人猜测环节

### 2.3 胜负判定
- **狼人胜利**：所有神民或所有平民出局
- **好人胜利**：所有狼人出局

## 3. 数据库设计

### 3.1 表结构

#### games表
| 字段名 | 数据类型 | 约束 | 描述 |
|--------|---------|------|------|
| `id` | `SERIAL` | `PRIMARY KEY` | 游戏ID |
| `status` | `VARCHAR(20)` | `NOT NULL` | 游戏状态（准备中、进行中、已结束） |
| `created_at` | `TIMESTAMP` | `NOT NULL DEFAULT CURRENT_TIMESTAMP` | 创建时间 |
| `started_at` | `TIMESTAMP` | `NULL` | 开始时间 |
| `ended_at` | `TIMESTAMP` | `NULL` | 结束时间 |
| `winner` | `VARCHAR(10)` | `NULL` | 获胜方（狼人/好人） |
| `host_id` | `INTEGER` | `REFERENCES users(id)` | 主持人ID |
| `game_code` | `VARCHAR(20)` | `UNIQUE NOT NULL` | 游戏邀请码 |

#### game_players表
| 字段名 | 数据类型 | 约束 | 描述 |
|--------|---------|------|------|
| `id` | `SERIAL` | `PRIMARY KEY` | 记录ID |
| `game_id` | `INTEGER` | `NOT NULL REFERENCES games(id) ON DELETE CASCADE` | 游戏ID |
| `player_id` | `INTEGER` | `REFERENCES users(id) NULL` | 玩家ID（NULL表示AI） |
| `ai_id` | `INTEGER` | `REFERENCES ai_digital_persons(id) NULL` | AI数字人ID（NULL表示真人） |
| `role` | `VARCHAR(20)` | `NOT NULL` | 角色（狼人、预言家、女巫、猎人、守卫、平民） |
| `role_code` | `INTEGER` | `NOT NULL` | 角色数字代号（1-12） |
| `status` | `VARCHAR(20)` | `NOT NULL DEFAULT 'alive'` | 状态（存活、死亡） |
| `score` | `INTEGER` | `NULL` | 游戏评分 |
| `is_real` | `BOOLEAN` | `NOT NULL` | 是否为真人 |
| `join_time` | `TIMESTAMP` | `NOT NULL DEFAULT CURRENT_TIMESTAMP` | 加入时间 |

#### game_phases表
| 字段名 | 数据类型 | 约束 | 描述 |
|--------|---------|------|------|
| `id` | `SERIAL` | `PRIMARY KEY` | 阶段ID |
| `game_id` | `INTEGER` | `NOT NULL REFERENCES games(id) ON DELETE CASCADE` | 游戏ID |
| `phase_type` | `VARCHAR(20)` | `NOT NULL` | 阶段类型（准备、夜晚、白天、结束） |
| `phase_number` | `INTEGER` | `NOT NULL` | 阶段编号 |
| `start_time` | `TIMESTAMP` | `NOT NULL DEFAULT CURRENT_TIMESTAMP` | 开始时间 |
| `end_time` | `TIMESTAMP` | `NULL` | 结束时间 |
| `current_phase` | `BOOLEAN` | `NOT NULL DEFAULT TRUE` | 是否为当前阶段 |

#### night_actions表
| 字段名 | 数据类型 | 约束 | 描述 |
|--------|---------|------|------|
| `id` | `SERIAL` | `PRIMARY KEY` | 行动ID |
| `game_id` | `INTEGER` | `NOT NULL REFERENCES games(id) ON DELETE CASCADE` | 游戏ID |
| `phase_id` | `INTEGER` | `NOT NULL REFERENCES game_phases(id) ON DELETE CASCADE` | 阶段ID |
| `player_id` | `INTEGER` | `NOT NULL REFERENCES game_players(id) ON DELETE CASCADE` | 玩家ID |
| `action_type` | `VARCHAR(20)` | `NOT NULL` | 行动类型（杀人、验人、用药、守护） |
| `target_id` | `INTEGER` | `REFERENCES game_players(id) NULL` | 目标ID |
| `action_result` | `VARCHAR(50)` | `NULL` | 行动结果 |
| `timestamp` | `TIMESTAMP` | `NOT NULL DEFAULT CURRENT_TIMESTAMP` | 行动时间 |

#### day_actions表
| 字段名 | 数据类型 | 约束 | 描述 |
|--------|---------|------|------|
| `id` | `SERIAL` | `PRIMARY KEY` | 行动ID |
| `game_id` | `INTEGER` | `NOT NULL REFERENCES games(id) ON DELETE CASCADE` | 游戏ID |
| `phase_id` | `INTEGER` | `NOT NULL REFERENCES game_phases(id) ON DELETE CASCADE` | 阶段ID |
| `player_id` | `INTEGER` | `NOT NULL REFERENCES game_players(id) ON DELETE CASCADE` | 玩家ID |
| `action_type` | `VARCHAR(20)` | `NOT NULL` | 行动类型（发言、投票） |
| `content` | `TEXT` | `NULL` | 发言内容 |
| `target_id` | `INTEGER` | `REFERENCES game_players(id) NULL` | 投票目标ID |
| `timestamp` | `TIMESTAMP` | `NOT NULL DEFAULT CURRENT_TIMESTAMP` | 行动时间 |

### 3.2 索引设计

| 索引名 | 表名 | 字段 | 类型 | 描述 |
|--------|------|------|------|------|
| `games_status_idx` | `games` | `status` | `BTREE` | 加速查询不同状态的游戏 |
| `games_game_code_idx` | `games` | `game_code` | `BTREE` | 加速通过游戏邀请码查询游戏 |
| `game_players_game_id_idx` | `game_players` | `game_id` | `BTREE` | 加速查询游戏中的玩家 |
| `game_players_role_code_idx` | `game_players` | `role_code` | `BTREE` | 加速通过角色代号查询玩家 |
| `game_players_status_idx` | `game_players` | `status` | `BTREE` | 加速查询存活/死亡的玩家 |
| `game_phases_game_id_idx` | `game_phases` | `game_id` | `BTREE` | 加速查询游戏的阶段 |
| `game_phases_current_phase_idx` | `game_phases` | `current_phase` | `BTREE` | 加速查询当前阶段 |
| `night_actions_phase_id_idx` | `night_actions` | `phase_id` | `BTREE` | 加速查询夜晚行动 |
| `day_actions_phase_id_idx` | `day_actions` | `phase_id` | `BTREE` | 加速查询白天行动 |

## 4. 核心功能实现

### 4.1 游戏创建与准备

#### 游戏创建
- 生成唯一的游戏ID和邀请码
- 设置游戏状态为"准备中"
- 记录主持人信息

#### 角色分配
- 随机分配12个角色（4狼人、4神民、4平民）
- 为每个角色分配唯一的数字代号（1-12）
- 记录角色分配结果到game_players表

#### AI数字人填充
- 当真人玩家不足12人时，从AI数字人池中选择AI数字人填充
- 为每个AI数字人分配角色
- 记录AI数字人信息到game_players表

### 4.2 游戏流程控制

#### 夜晚阶段
1. 狼人选择杀人目标
   - 狼人团队讨论并选择杀人目标
   - 记录狼人杀人行动到night_actions表

2. 预言家验人
   - 预言家选择验人目标
   - 系统返回验人结果（好人/狼人）
   - 记录预言家验人行动到night_actions表

3. 女巫用药
   - 女巫选择使用解药或毒药
   - 记录女巫用药行动到night_actions表

4. 守卫守护
   - 守卫选择守护目标
   - 记录守卫守护行动到night_actions表

#### 白天阶段
1. 宣布夜晚结果
   - 宣布夜晚死亡的玩家
   - 处理猎人死亡时的开枪技能

2. 玩家发言
   - 按照发言顺序，每个玩家依次发言
   - 记录玩家发言内容到day_actions表

3. 投票放逐
   - 每个玩家投票选择放逐目标
   - 记录玩家投票行动到day_actions表
   - 计算投票结果，放逐得票最多的玩家

### 4.3 胜负判定

#### 狼人胜利条件
- 所有神民出局
- 或所有平民出局

#### 好人胜利条件
- 所有狼人出局

#### 游戏结束处理
- 更新游戏状态为"已结束"
- 记录获胜方
- 计算玩家评分
- 进入真人猜测环节

## 5. API接口设计

### 5.1 游戏管理接口

#### 创建游戏
- **路径**：`/api/games`
- **方法**：`POST`
- **参数**：
  - `host_id`：主持人ID
- **返回**：
  - `game_id`：游戏ID
  - `game_code`：游戏邀请码
  - `status`：游戏状态

#### 加入游戏
- **路径**：`/api/games/join`
- **方法**：`POST`
- **参数**：
  - `game_code`：游戏邀请码
  - `player_id`：玩家ID
- **返回**：
  - `game_id`：游戏ID
  - `role`：角色
  - `role_code`：角色数字代号

#### 开始游戏
- **路径**：`/api/games/{game_id}/start`
- **方法**：`POST`
- **参数**：
  - `game_id`：游戏ID
- **返回**：
  - `status`：游戏状态
  - `phase`：当前阶段

### 5.2 游戏操作接口

#### 夜晚行动
- **路径**：`/api/games/{game_id}/night-action`
- **方法**：`POST`
- **参数**：
  - `game_id`：游戏ID
  - `player_id`：玩家ID
  - `action_type`：行动类型
  - `target_id`：目标ID
- **返回**：
  - `success`：操作是否成功
  - `result`：行动结果

#### 白天发言
- **路径**：`/api/games/{game_id}/day-action/speak`
- **方法**：`POST`
- **参数**：
  - `game_id`：游戏ID
  - `player_id`：玩家ID
  - `content`：发言内容
- **返回**：
  - `success`：操作是否成功

#### 投票放逐
- **路径**：`/api/games/{game_id}/day-action/vote`
- **方法**：`POST`
- **参数**：
  - `game_id`：游戏ID
  - `player_id`：玩家ID
  - `target_id`：投票目标ID
- **返回**：
  - `success`：操作是否成功
  - `vote_count`：当前得票数

### 5.3 游戏状态接口

#### 获取游戏状态
- **路径**：`/api/games/{game_id}/status`
- **方法**：`GET`
- **参数**：
  - `game_id`：游戏ID
- **返回**：
  - `status`：游戏状态
  - `phase`：当前阶段
  - `players`：玩家信息
  - `actions`：最近行动

#### 获取游戏历史
- **路径**：`/api/games/{game_id}/history`
- **方法**：`GET`
- **参数**：
  - `game_id`：游戏ID
- **返回**：
  - `phases`：游戏阶段历史
  - `actions`：所有行动记录

## 6. 数据流转

### 6.1 游戏创建流程
1. 主持人创建游戏
2. 系统生成游戏ID和邀请码
3. 系统设置游戏状态为"准备中"
4. 玩家通过邀请码加入游戏
5. 当玩家人数达到1人时，系统填充AI数字人至12人
6. 主持人开始游戏
7. 系统分配角色并开始第一夜

### 6.2 夜晚行动流程
1. 系统进入夜晚阶段
2. 狼人选择杀人目标
3. 预言家选择验人目标
4. 女巫选择是否使用解药或毒药
5. 守卫选择守护目标
6. 系统处理夜晚行动结果
7. 系统进入白天阶段

### 6.3 白天行动流程
1. 系统宣布夜晚结果
2. 玩家按顺序发言
3. 玩家投票放逐
4. 系统处理投票结果
5. 系统检查胜负条件
6. 如果游戏未结束，进入下一夜

## 7. 异常处理

### 7.1 游戏状态异常
- **处理**：验证游戏状态，确保操作符合当前阶段
- **返回**：返回错误信息，提示当前游戏状态不允许该操作

### 7.2 玩家操作异常
- **处理**：验证玩家身份和权限，确保玩家只能执行自己的操作
- **返回**：返回错误信息，提示玩家无权限执行该操作

### 7.3 网络异常
- **处理**：实现断线重连机制，恢复游戏状态
- **返回**：重连成功后，同步游戏状态

## 8. 性能优化

### 8.1 数据库优化
- 使用索引加速查询
- 合理使用事务，确保数据一致性
- 批量处理操作，减少数据库访问次数

### 8.2 业务逻辑优化
- 预计算游戏状态，减少实时计算
- 缓存游戏数据，提高响应速度
- 异步处理非实时操作

### 8.3 并发处理
- 使用锁机制，防止并发操作冲突
- 实现乐观锁，提高并发性能
- 合理设计事务隔离级别

## 9. 测试策略

### 9.1 单元测试
- 测试游戏规则逻辑
- 测试角色分配算法
- 测试胜负判定逻辑

### 9.2 集成测试
- 测试游戏流程完整性
- 测试API接口功能
- 测试数据库操作

### 9.3 性能测试
- 测试游戏创建和加入的响应时间
- 测试多玩家同时操作的性能
- 测试游戏状态更新的实时性

## 10. 未来扩展

### 10.1 功能扩展
- 支持不同人数的游戏模式（6人、9人等）
- 支持自定义角色配置
- 支持特殊规则和主题模式

### 10.2 技术扩展
- 实现游戏记录回放功能
- 支持游戏数据分析和统计
- 集成AI辅助决策系统

## 11. 结论

游戏核心逻辑模块是整个游戏的基础，通过详细的数据库设计和业务逻辑实现，确保游戏按照标准狼人杀规则运行。该模块的设计考虑了扩展性和性能优化，为未来的功能扩展和技术升级预留了空间。