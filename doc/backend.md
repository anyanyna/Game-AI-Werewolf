# 后端服务模块详细设计文档

## 1. 模块概述

后端服务模块负责处理游戏的核心逻辑、数据持久化、用户管理等功能。该模块确保游戏规则的正确执行，数据的安全存储，以及与前端的实时通信。

## 2. 功能需求

### 2.1 游戏逻辑服务
- 维护游戏状态和流程控制
- 处理游戏规则和胜负判定
- 管理玩家操作和游戏事件

### 2.2 AI服务
- 处理AI数字人的决策请求
- 管理AI数字人的进化和状态
- 维护AI数字人池

### 2.3 用户服务
- 管理用户注册和登录
- 存储和管理用户数据
- 记录游戏历史和评分

### 2.4 WebSocket通信
- 实现实时数据同步
- 处理玩家操作和游戏状态更新
- 管理连接和消息传递

### 2.5 URL分享服务
- 生成游戏邀请链接
- 处理通过URL加入游戏的请求
- 管理游戏邀请码

## 3. 技术栈

### 3.1 核心技术
- **框架**：Flask 或 Django
- **数据库**：PostgreSQL
- **通信**：Flask-SocketIO 或 Django Channels
- **认证**：JWT
- **缓存**：Redis
- **ORM**：SQLAlchemy 或 Django ORM

### 3.2 第三方库
- **数据验证**：Pydantic
- **密码加密**：bcrypt
- **日志**：logging
- **测试**：pytest
- **部署**：Gunicorn, Nginx

## 4. 系统架构

### 4.1 模块划分
- **API模块**：处理HTTP请求和响应
- **WebSocket模块**：处理实时通信
- **游戏逻辑模块**：处理游戏规则和流程
- **AI模块**：处理AI数字人的决策和进化
- **用户模块**：处理用户管理和认证
- **数据库模块**：处理数据存储和查询

### 4.2 数据流
1. 前端发送请求到API模块
2. API模块处理请求并调用相应的业务逻辑
3. 业务逻辑模块与数据库模块交互
4. 对于实时操作，通过WebSocket模块推送到前端
5. 前端接收数据并更新界面

## 5. 数据库设计

### 5.1 表结构

#### users表
| 字段名 | 数据类型 | 约束 | 描述 |
|--------|---------|------|------|
| `id` | `SERIAL` | `PRIMARY KEY` | 用户ID |
| `username` | `VARCHAR(50)` | `UNIQUE NOT NULL` | 用户名 |
| `email` | `VARCHAR(100)` | `UNIQUE NOT NULL` | 邮箱 |
| `password_hash` | `VARCHAR(255)` | `NOT NULL` | 密码哈希 |
| `created_at` | `TIMESTAMP` | `NOT NULL DEFAULT CURRENT_TIMESTAMP` | 创建时间 |
| `last_login` | `TIMESTAMP` | `NULL` | 最后登录时间 |
| `role` | `VARCHAR(20)` | `NOT NULL DEFAULT 'user'` | 角色 |

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

#### real_person_guesses表
| 字段名 | 数据类型 | 约束 | 描述 |
|--------|---------|------|------|
| `id` | `SERIAL` | `PRIMARY KEY` | 记录ID |
| `game_id` | `INTEGER` | `NOT NULL REFERENCES games(id) ON DELETE CASCADE` | 游戏ID |
| `guesser_id` | `INTEGER` | `NOT NULL REFERENCES users(id) ON DELETE CASCADE` | 猜测者ID |
| `target_id` | `INTEGER` | `NOT NULL REFERENCES game_players(id) ON DELETE CASCADE` | 目标ID |
| `guess` | `BOOLEAN` | `NOT NULL` | 猜测结果（TRUE为真人，FALSE为AI） |
| `is_correct` | `BOOLEAN` | `NOT NULL` | 是否正确 |
| `created_at` | `TIMESTAMP` | `NOT NULL DEFAULT CURRENT_TIMESTAMP` | 创建时间 |

### 5.2 索引设计

| 索引名 | 表名 | 字段 | 类型 | 描述 |
|--------|------|------|------|------|
| `users_username_idx` | `users` | `username` | `BTREE` | 加速用户名查询 |
| `users_email_idx` | `users` | `email` | `BTREE` | 加速邮箱查询 |
| `games_status_idx` | `games` | `status` | `BTREE` | 加速查询不同状态的游戏 |
| `games_game_code_idx` | `games` | `game_code` | `BTREE` | 加速通过游戏邀请码查询游戏 |
| `game_players_game_id_idx` | `game_players` | `game_id` | `BTREE` | 加速查询游戏中的玩家 |
| `game_players_role_code_idx` | `game_players` | `role_code` | `BTREE` | 加速通过角色代号查询玩家 |
| `ai_digital_persons_last_active_idx` | `ai_digital_persons` | `last_active` | `BTREE` | 加速查询活跃的AI数字人 |
| `game_phases_game_id_idx` | `game_phases` | `game_id` | `BTREE` | 加速查询游戏的阶段 |
| `game_phases_current_phase_idx` | `game_phases` | `current_phase` | `BTREE` | 加速查询当前阶段 |
| `night_actions_phase_id_idx` | `night_actions` | `phase_id` | `BTREE` | 加速查询夜晚行动 |
| `day_actions_phase_id_idx` | `day_actions` | `phase_id` | `BTREE` | 加速查询白天行动 |
| `real_person_guesses_game_id_idx` | `real_person_guesses` | `game_id` | `BTREE` | 加速查询游戏中的猜测记录 |

## 6. API接口设计

### 6.1 用户管理接口

#### 注册
- **路径**：`/api/auth/register`
- **方法**：`POST`
- **参数**：
  - `username`：用户名
  - `email`：邮箱
  - `password`：密码
- **返回**：
  - `user_id`：用户ID
  - `username`：用户名
  - `token`：JWT令牌

#### 登录
- **路径**：`/api/auth/login`
- **方法**：`POST`
- **参数**：
  - `email`：邮箱
  - `password`：密码
- **返回**：
  - `user_id`：用户ID
  - `username`：用户名
  - `token`：JWT令牌

#### 获取个人信息
- **路径**：`/api/users/me`
- **方法**：`GET`
- **参数**：
  - `Authorization`：Bearer令牌
- **返回**：
  - `user_id`：用户ID
  - `username`：用户名
  - `email`：邮箱
  - `created_at`：创建时间

### 6.2 游戏管理接口

#### 创建游戏
- **路径**：`/api/games`
- **方法**：`POST`
- **参数**：
  - `Authorization`：Bearer令牌
  - `host_id`：主持人ID
- **返回**：
  - `game_id`：游戏ID
  - `game_code`：游戏邀请码
  - `status`：游戏状态

#### 加入游戏
- **路径**：`/api/games/join`
- **方法**：`POST`
- **参数**：
  - `Authorization`：Bearer令牌
  - `game_code`：游戏邀请码
- **返回**：
  - `game_id`：游戏ID
  - `role`：角色
  - `role_code`：角色数字代号

#### 开始游戏
- **路径**：`/api/games/{game_id}/start`
- **方法**：`POST`
- **参数**：
  - `Authorization`：Bearer令牌
  - `game_id`：游戏ID
- **返回**：
  - `status`：游戏状态
  - `phase`：当前阶段

#### 获取游戏状态
- **路径**：`/api/games/{game_id}/status`
- **方法**：`GET`
- **参数**：
  - `Authorization`：Bearer令牌
  - `game_id`：游戏ID
- **返回**：
  - `status`：游戏状态
  - `phase`：当前阶段
  - `players`：玩家信息
  - `actions`：最近行动

### 6.3 游戏操作接口

#### 夜晚行动
- **路径**：`/api/games/{game_id}/night-action`
- **方法**：`POST`
- **参数**：
  - `Authorization`：Bearer令牌
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
  - `Authorization`：Bearer令牌
  - `game_id`：游戏ID
  - `player_id`：玩家ID
  - `content`：发言内容
- **返回**：
  - `success`：操作是否成功

#### 投票放逐
- **路径**：`/api/games/{game_id}/day-action/vote`
- **方法**：`POST`
- **参数**：
  - `Authorization`：Bearer令牌
  - `game_id`：游戏ID
  - `player_id`：玩家ID
  - `target_id`：投票目标ID
- **返回**：
  - `success`：操作是否成功
  - `vote_count`：当前得票数

#### 真人猜测
- **路径**：`/api/games/{game_id}/guess`
- **方法**：`POST`
- **参数**：
  - `Authorization`：Bearer令牌
  - `game_id`：游戏ID
  - `guesser_id`：猜测者ID
  - `target_id`：目标ID
  - `guess`：猜测结果
- **返回**：
  - `success`：操作是否成功
  - `is_correct`：是否正确

### 6.4 AI管理接口

#### 创建AI数字人
- **路径**：`/api/ai-digital-persons`
- **方法**：`POST`
- **参数**：
  - `Authorization`：Bearer令牌
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
  - `Authorization`：Bearer令牌
  - `personality`：性格（可选）
  - `experience`：阅历（可选）
  - `limit`：返回数量（可选）
- **返回**：
  - `ais`：AI数字人列表
  - `total`：总数量

#### AI决策
- **路径**：`/api/ai-digital-persons/{ai_id}/decision`
- **方法**：`POST`
- **参数**：
  - `Authorization`：Bearer令牌
  - `ai_id`：AI数字人ID
  - `game_state`：游戏状态
  - `role`：角色
  - `action_type`：行动类型
- **返回**：
  - `decision`：决策结果
  - `confidence`：决策信心

## 7. WebSocket接口设计

### 7.1 连接管理
- **路径**：`/socket.io`
- **事件**：
  - `connect`：建立连接
  - `disconnect`：断开连接
  - `error`：连接错误

### 7.2 游戏事件
- **事件**：`game_state_update`
  - **描述**：游戏状态更新
  - **数据**：
    - `game_id`：游戏ID
    - `state`：游戏状态
    - `phase`：当前阶段
    - `players`：玩家信息

- **事件**：`player_action`
  - **描述**：玩家操作
  - **数据**：
    - `game_id`：游戏ID
    - `player_id`：玩家ID
    - `action_type`：行动类型
    - `target_id`：目标ID
    - `content`：操作内容

- **事件**：`ai_action`
  - **描述**：AI数字人操作
  - **数据**：
    - `game_id`：游戏ID
    - `ai_id`：AI数字人ID
    - `action_type`：行动类型
    - `target_id`：目标ID
    - `content`：操作内容

- **事件**：`chat_message`
  - **描述**：聊天消息
  - **数据**：
    - `game_id`：游戏ID
    - `player_id`：玩家ID
    - `message`：消息内容
    - `timestamp`：时间戳

## 8. 业务逻辑实现

### 8.1 游戏逻辑

#### 游戏创建
- 生成唯一的游戏ID和邀请码
- 设置游戏状态为"准备中"
- 记录主持人信息

#### 角色分配
- 随机分配12个角色（4狼人、4神民、4平民）
- 为每个角色分配唯一的数字代号（1-12）
- 记录角色分配结果到game_players表

#### 游戏流程控制
- 管理游戏阶段转换（准备→夜晚→白天→结束）
- 处理夜晚行动和白天行动
- 检查胜负条件

#### 胜负判定
- 狼人胜利：所有神民或所有平民出局
- 好人胜利：所有狼人出局

### 8.2 AI服务

#### AI数字人决策
- 基于AI数字人的性格、阅历、智商参数生成决策
- 考虑当前游戏状态和角色身份
- 实现不同角色的决策逻辑

#### AI数字人进化
- 游戏结束后根据表现更新能力值
- 实现遗忘机制，长时间不参与游戏后能力逐渐减退
- 确保AI能力与真人相当

#### AI数字人池管理
- 管理AI数字人的创建和分配
- 维护AI数字人的状态和活跃度
- 确保AI数字人的多样性

### 8.3 用户服务

#### 用户认证
- 处理用户注册和登录
- 生成和验证JWT令牌
- 密码加密和验证

#### 用户数据管理
- 存储和管理用户信息
- 记录用户游戏历史
- 计算用户评分和统计数据

#### 游戏记录管理
- 保存游戏历史记录
- 分析游戏数据
- 生成游戏报告

### 8.4 URL分享服务

#### 邀请码生成
- 生成唯一的游戏邀请码
- 关联邀请码与游戏ID
- 支持短链接生成

#### 加入游戏处理
- 解析邀请码
- 验证游戏状态
- 将玩家加入游戏

## 9. 性能优化

### 9.1 数据库优化
- 使用索引加速查询
- 合理使用事务，确保数据一致性
- 批量处理操作，减少数据库访问次数
- 使用连接池，减少连接建立的开销

### 9.2 缓存优化
- 使用Redis缓存热点数据
- 缓存游戏状态和玩家信息
- 实现缓存过期策略

### 9.3 并发处理
- 使用异步处理，提高并发性能
- 实现锁机制，防止并发操作冲突
- 优化WebSocket连接管理

### 9.4 代码优化
- 优化算法，减少时间复杂度
- 减少不必要的计算和IO操作
- 使用异步IO，提高处理效率

## 10. 安全性考虑

### 10.1 认证与授权
- 使用JWT进行身份认证
- 实现基于角色的访问控制
- 验证用户权限，防止越权操作

### 10.2 数据安全
- 密码加密存储
- 敏感数据加密传输
- 防止SQL注入和XSS攻击

### 10.3 游戏安全
- 防作弊机制，检测异常操作
- 验证游戏操作的合法性
- 防止恶意用户破坏游戏平衡

### 10.4 网络安全
- 使用HTTPS加密传输
- 实现速率限制，防止DDoS攻击
- 监控异常访问模式

## 11. 测试策略

### 11.1 单元测试
- 测试API接口
- 测试游戏逻辑
- 测试AI决策逻辑

### 11.2 集成测试
- 测试模块间的集成
- 测试数据库操作
- 测试WebSocket通信

### 11.3 端到端测试
- 测试完整的游戏流程
- 测试多用户同时操作
- 测试异常情况处理

### 11.4 性能测试
- 测试系统在高负载下的表现
- 测试响应时间和吞吐量
- 测试数据库性能

## 12. 部署策略

### 12.1 服务器配置
- **Web服务器**：Nginx
- **应用服务器**：Gunicorn
- **数据库**：PostgreSQL
- **缓存**：Redis

### 12.2 环境配置
- **开发环境**：本地开发配置
- **测试环境**：测试服务器配置
- **生产环境**：生产服务器配置

### 12.3 部署流程
- **CI/CD**：使用GitHub Actions或Jenkins
- **自动化部署**：脚本化部署流程
- **监控**：使用Prometheus和Grafana

### 12.4 扩展性考虑
- **水平扩展**：支持多服务器部署
- **负载均衡**：使用Nginx负载均衡
- **数据库复制**：实现数据库高可用

## 13. 未来扩展

### 13.1 功能扩展
- **更多游戏模式**：支持不同人数和规则的游戏模式
- **社交功能**：添加好友系统、聊天功能
- **赛事系统**：组织线上比赛

### 13.2 技术扩展
- **微服务架构**：将系统拆分为微服务
- **容器化**：使用Docker和Kubernetes
- **云服务**：使用AWS或阿里云等云服务

## 14. 结论

后端服务模块是游戏的核心组成部分，通过详细的数据库设计、API接口设计和业务逻辑实现，确保游戏规则的正确执行和数据的安全存储。该模块的设计考虑了性能优化、安全性和可扩展性，为游戏提供了稳定、高效的后端支持。同时，设计中预留了充分的扩展空间，可以在未来根据技术发展和用户需求进行功能和技术的升级。