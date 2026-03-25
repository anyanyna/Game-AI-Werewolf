# AI数字人系统模块详细设计文档

## 1. 模块概述

AI数字人系统模块负责管理AI数字人的创建、行为模拟、能力进化等核心功能。该模块确保AI数字人能够以自然、智能的方式参与游戏，为玩家提供良好的游戏体验。

## 2. 功能需求

### 2.1 AI数字人池管理
- 建立AI数字人池，存储和管理AI数字人
- 为每个AI数字人创建独立的人格soul.md文件
- 管理AI数字人的生命周期和状态

### 2.2 AI数字人属性
- **性格**：开朗、内敛、逻辑型、直觉型等
- **阅历**：新手、普通、资深
- **智商**：随机生成（偏向普通真人水平）
- **能力值**：逻辑推理、发言质量、团队协作等
- **游戏历史**：参与的游戏记录

### 2.3 AI行为模拟
- 基于性格、阅历、智商参数生成AI决策
- 实现AI数字人的发言逻辑
- 开发AI数字人的投票策略

### 2.4 AI数字人进化
- 游戏结束后根据表现缓慢进化
- 实现遗忘机制，长时间不参与游戏后能力逐渐减退
- 确保AI能力与真人相当，保持游戏平衡性

### 2.5 独立Session通道
- 为每个AI数字人创建独立的会话
- 维护独立的对话历史和上下文
- 防止大模型认为多个AI是同一个人

## 3. 数据库设计

### 3.1 表结构

#### ai_digital_persons表
| 字段名 | 数据类型 | 约束 | 描述 |
|--------|---------|------|------|
| `id` | `SERIAL` | `PRIMARY KEY` | AI数字人ID |
| `name` | `VARCHAR(50)` | `NOT NULL` | AI名称 |
| `personality` | `VARCHAR(50)` | `NOT NULL` | 性格 |
| `experience` | `VARCHAR(20)` | `NOT NULL` | 阅历 |
| `intelligence` | `INTEGER` | `NOT NULL` | 智商（1-100） |
| `abilities` | `JSONB` | `NOT NULL` | 能力值（JSON格式） |
| `game_history` | `JSONB` | `NOT NULL` | 游戏历史（JSON格式） |
| `last_active` | `TIMESTAMP` | `NOT NULL DEFAULT CURRENT_TIMESTAMP` | 最后活跃时间 |
| `created_at` | `TIMESTAMP` | `NOT NULL DEFAULT CURRENT_TIMESTAMP` | 创建时间 |
| `updated_at` | `TIMESTAMP` | `NOT NULL DEFAULT CURRENT_TIMESTAMP` | 更新时间 |

#### ai_sessions表
| 字段名 | 数据类型 | 约束 | 描述 |
|--------|---------|------|------|
| `id` | `SERIAL` | `PRIMARY KEY` | 会话ID |
| `ai_id` | `INTEGER` | `NOT NULL REFERENCES ai_digital_persons(id) ON DELETE CASCADE` | AI数字人ID |
| `session_id` | `VARCHAR(100)` | `NOT NULL UNIQUE` | 会话标识 |
| `context` | `JSONB` | `NOT NULL` | 会话上下文 |
| `last_used` | `TIMESTAMP` | `NOT NULL DEFAULT CURRENT_TIMESTAMP` | 最后使用时间 |
| `created_at` | `TIMESTAMP` | `NOT NULL DEFAULT CURRENT_TIMESTAMP` | 创建时间 |

#### ai_game_performances表
| 字段名 | 数据类型 | 约束 | 描述 |
|--------|---------|------|------|
| `id` | `SERIAL` | `PRIMARY KEY` | 记录ID |
| `ai_id` | `INTEGER` | `NOT NULL REFERENCES ai_digital_persons(id) ON DELETE CASCADE` | AI数字人ID |
| `game_id` | `INTEGER` | `NOT NULL REFERENCES games(id) ON DELETE CASCADE` | 游戏ID |
| `role` | `VARCHAR(20)` | `NOT NULL` | 角色 |
| `score` | `INTEGER` | `NOT NULL` | 游戏评分 |
| `performance` | `JSONB` | `NOT NULL` | 表现数据（JSON格式） |
| `created_at` | `TIMESTAMP` | `NOT NULL DEFAULT CURRENT_TIMESTAMP` | 创建时间 |

#### ai_evolution_records表
| 字段名 | 数据类型 | 约束 | 描述 |
|--------|---------|------|------|
| `id` | `SERIAL` | `PRIMARY KEY` | 记录ID |
| `ai_id` | `INTEGER` | `NOT NULL REFERENCES ai_digital_persons(id) ON DELETE CASCADE` | AI数字人ID |
| `evolution_type` | `VARCHAR(50)` | `NOT NULL` | 进化类型（能力提升、性格变化等） |
| `old_value` | `JSONB` | `NOT NULL` | 旧值 |
| `new_value` | `JSONB` | `NOT NULL` | 新值 |
| `reason` | `TEXT` | `NULL` | 进化原因 |
| `created_at` | `TIMESTAMP` | `NOT NULL DEFAULT CURRENT_TIMESTAMP` | 创建时间 |

### 3.2 索引设计

| 索引名 | 表名 | 字段 | 类型 | 描述 |
|--------|------|------|------|------|
| `ai_digital_persons_last_active_idx` | `ai_digital_persons` | `last_active` | `BTREE` | 加速查询活跃的AI数字人 |
| `ai_digital_persons_personality_idx` | `ai_digital_persons` | `personality` | `BTREE` | 加速按性格查询AI数字人 |
| `ai_digital_persons_experience_idx` | `ai_digital_persons` | `experience` | `BTREE` | 加速按阅历查询AI数字人 |
| `ai_sessions_ai_id_idx` | `ai_sessions` | `ai_id` | `BTREE` | 加速查询AI数字人的会话 |
| `ai_sessions_session_id_idx` | `ai_sessions` | `session_id` | `BTREE` | 加速通过会话ID查询 |
| `ai_game_performances_ai_id_idx` | `ai_game_performances` | `ai_id` | `BTREE` | 加速查询AI数字人的游戏表现 |
| `ai_game_performances_game_id_idx` | `ai_game_performances` | `game_id` | `BTREE` | 加速查询游戏中的AI表现 |
| `ai_evolution_records_ai_id_idx` | `ai_evolution_records` | `ai_id` | `BTREE` | 加速查询AI数字人的进化记录 |

## 4. 核心功能实现

### 4.1 AI数字人池管理

#### AI数字人创建
- 生成AI数字人的基本属性（性格、阅历、智商）
- 初始化能力值和游戏历史
- 创建独立的会话通道
- 生成人格soul.md文件

#### AI数字人分配
- 根据游戏需求从AI数字人池中选择合适的AI
- 考虑AI的活跃度、能力水平和性格特点
- 确保AI数字人的多样性

#### AI数字人状态管理
- 跟踪AI数字人的活跃状态
- 定期清理不活跃的AI数字人
- 更新AI数字人的能力和状态

### 4.2 AI行为模拟

#### 决策机制
- 基于AI数字人的性格、阅历、智商参数生成决策
- 考虑当前游戏状态和角色身份
- 实现不同角色的决策逻辑（狼人、神民、平民）

#### 发言逻辑
- 狼人：伪装、误导、保护同伴
- 预言家：提供信息、引导投票
- 女巫：隐藏身份、使用药水
- 猎人：谨慎发言、准备开枪
- 守卫：保护重要目标
- 平民：分析信息、跟随神民

#### 投票策略
- 基于分析结果选择投票目标
- 考虑团队利益和个人判断
- 模拟不同性格的投票风格

### 4.3 AI数字人进化

#### 能力进化
- 游戏结束后根据表现更新能力值
- 表现好的AI数字人能力提升
- 表现差的AI数字人能力可能下降
- 进化速度缓慢，符合人类学习曲线

#### 遗忘机制
- 长时间不参与游戏的AI数字人能力逐渐减退
- 记忆衰减模拟人类遗忘过程
- 重新参与游戏后能力会逐渐恢复

#### 平衡机制
- 监控AI数字人的能力水平
- 确保AI能力与真人相当
- 防止AI过强或过弱影响游戏平衡

### 4.4 独立Session通道

#### 会话管理
- 为每个AI数字人创建独立的会话
- 维护会话上下文和对话历史
- 确保会话的隔离性和独立性

#### 上下文管理
- 为每个AI数字人维护独立的上下文
- 记录对话历史和游戏状态
- 确保AI数字人能够基于历史信息做出决策

#### 会话刷新
- 定期刷新会话，保持会话活跃
- 处理会话过期和重新创建
- 确保会话的稳定性和可靠性

## 5. API接口设计

### 5.1 AI数字人管理接口

#### 创建AI数字人
- **路径**：`/api/ai-digital-persons`
- **方法**：`POST`
- **参数**：
  - `name`：AI名称
  - `personality`：性格
  - `experience`：阅历
  - `intelligence`：智商
- **返回**：
  - `ai_id`：AI数字人ID
  - `status`：创建状态

#### 获取AI数字人列表
- **路径**：`/api/ai-digital-persons`
- **方法**：`GET`
- **参数**：
  - `personality`：性格（可选）
  - `experience`：阅历（可选）
  - `limit`：返回数量（可选）
- **返回**：
  - `ais`：AI数字人列表
  - `total`：总数量

#### 获取AI数字人详情
- **路径**：`/api/ai-digital-persons/{ai_id}`
- **方法**：`GET`
- **参数**：
  - `ai_id`：AI数字人ID
- **返回**：
  - `ai`：AI数字人详情
  - `abilities`：能力值
  - `game_history`：游戏历史

### 5.2 AI行为接口

#### AI决策
- **路径**：`/api/ai-digital-persons/{ai_id}/decision`
- **方法**：`POST`
- **参数**：
  - `ai_id`：AI数字人ID
  - `game_state`：游戏状态
  - `role`：角色
  - `action_type`：行动类型
- **返回**：
  - `decision`：决策结果
  - `confidence`：决策信心

#### AI发言
- **路径**：`/api/ai-digital-persons/{ai_id}/speak`
- **方法**：`POST`
- **参数**：
  - `ai_id`：AI数字人ID
  - `game_state`：游戏状态
  - `role`：角色
  - `context`：发言上下文
- **返回**：
  - `content`：发言内容
  - `tone`：发言语气

#### AI投票
- **路径**：`/api/ai-digital-persons/{ai_id}/vote`
- **方法**：`POST`
- **参数**：
  - `ai_id`：AI数字人ID
  - `game_state`：游戏状态
  - `role`：角色
  - `candidates`：投票候选人
- **返回**：
  - `target_id`：投票目标ID
  - `reason`：投票原因

### 5.3 AI进化接口

#### 更新AI能力
- **路径**：`/api/ai-digital-persons/{ai_id}/evolve`
- **方法**：`POST`
- **参数**：
  - `ai_id`：AI数字人ID
  - `game_performance`：游戏表现
  - `role`：角色
- **返回**：
  - `success`：操作是否成功
  - `evolution`：进化结果

#### 获取AI进化记录
- **路径**：`/api/ai-digital-persons/{ai_id}/evolution-records`
- **方法**：`GET`
- **参数**：
  - `ai_id`：AI数字人ID
  - `limit`：返回数量（可选）
- **返回**：
  - `records`：进化记录列表

## 6. 数据流转

### 6.1 AI数字人创建流程
1. 系统生成AI数字人的基本属性
2. 初始化能力值和游戏历史
3. 创建独立的会话通道
4. 生成人格soul.md文件
5. 将AI数字人加入AI数字人池

### 6.2 AI数字人参与游戏流程
1. 游戏需要AI数字人填充时，从AI数字人池中选择合适的AI
2. 为AI数字人分配角色
3. 游戏过程中，AI数字人根据游戏状态做出决策
4. 游戏结束后，AI数字人根据表现进化
5. AI数字人返回AI数字人池，等待下一次分配

### 6.3 AI数字人进化流程
1. 游戏结束后，系统评估AI数字人的表现
2. 根据表现计算能力变化
3. 更新AI数字人的能力值和游戏历史
4. 记录进化过程
5. 应用遗忘机制（如果AI数字人长时间不活动）

## 7. 异常处理

### 7.1 AI数字人不存在
- **处理**：验证AI数字人ID的有效性
- **返回**：返回错误信息，提示AI数字人不存在

### 7.2 会话过期
- **处理**：检测会话状态，过期时重新创建会话
- **返回**：重新创建会话后继续操作

### 7.3 决策失败
- **处理**：捕获决策过程中的异常，使用默认决策
- **返回**：返回默认决策结果

## 8. 性能优化

### 8.1 数据库优化
- 使用索引加速查询
- 合理使用JSONB类型存储复杂数据
- 批量处理AI数字人的更新

### 8.2 会话管理优化
- 实现会话池，减少会话创建和销毁的开销
- 使用缓存存储活跃会话
- 定期清理过期会话

### 8.3 AI决策优化
- 缓存AI决策结果
- 批量处理AI请求
- 使用轻量级模型进行初步决策

## 9. 测试策略

### 9.1 单元测试
- 测试AI数字人创建和管理
- 测试AI决策逻辑
- 测试AI进化机制

### 9.2 集成测试
- 测试AI数字人参与游戏的完整流程
- 测试AI数字人与游戏核心逻辑的集成
- 测试AI数字人之间的交互

### 9.3 性能测试
- 测试AI数字人决策的响应时间
- 测试多AI数字人同时决策的性能
- 测试AI数字人进化的计算开销

## 10. 未来扩展

### 10.1 功能扩展
- 支持AI数字人的个性化定制
- 实现AI数字人之间的社交互动
- 开发AI数字人的技能树系统

### 10.2 技术扩展
- 集成更先进的AI模型
- 实现AI数字人的情感识别和表达
- 开发AI数字人的自主学习能力

## 11. 结论

AI数字人系统模块是游戏的重要组成部分，通过详细的数据库设计和业务逻辑实现，确保AI数字人能够以自然、智能的方式参与游戏。该模块的设计考虑了AI数字人的行为模拟、能力进化和会话管理，为玩家提供了良好的游戏体验。同时，设计中预留了充分的扩展空间，可以在未来根据技术发展和用户需求进行功能和技术的升级。