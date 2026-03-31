# 蒙面狼人杀 (Werewolf Game)

<p align="center">
  <img src="https://img.shields.io/badge/Go-1.22-blue" alt="Go Version">
  <img src="https://img.shields.io/badge/React-18-blue" alt="React Version">
  <img src="https://img.shields.io/badge/License-All%20Rights%20Reserved-red" alt="License">
  <img src="https://img.shields.io/badge/Status-Active-brightgreen" alt="Status">
</p>

一个支持 AI 玩家的在线狼人杀游戏平台。采用前后端分离架构，后端使用 Go + Gin，前端使用 React + Vite。

**本项目为闭源软件，仅供内部使用。未经授权，禁止复制、传播、修改或使用。**

## ✨ 功能特性

- **用户系统**: 注册、登录、会话保持
- **游戏大厅**: 创建游戏、加入游戏、等待玩家
- **游戏逻辑**: 标准12人局角色分配 (4狼、预言家、女巫、猎人、4村民)
- **实时通信**: WebSocket 实时更新游戏状态
- **AI 对手**: 自动填充 AI 数字人玩家
- **真人猜测**: 玩家可以猜测谁是真人

## 🛠 技术栈

### 后端
| 技术 | 版本 | 用途 |
|------|------|------|
| Go | 1.22+ | 编程语言 |
| Gin | 1.9.1 | Web 框架 |
| GORM | 1.25+ | ORM 框架 |
| SQLite | (glebarez/sqlite) | 数据库 (纯 Go 实现) |
| WebSocket | gorilla/websocket | 实时通信 |

### 前端
| 技术 | 版本 | 用途 |
|------|------|------|
| React | 18.2 | UI 框架 |
| Vite | 5.0 | 构建工具 |
| Vitest | 1.1 | 单元测试 |
| Playwright | - | E2E 测试 |

## 📁 项目结构

```
Game-AI-Werewolf/
├── backend/                 # Go 后端
│   ├── api/                # API 处理器
│   │   └── handlers.go     # 业务逻辑
│   ├── ai/                 # AI 行为树
│   │   ├── ai.go           # AI 接口
│   │   └── behavior_tree.go
│   ├── database/           # 数据库配置
│   │   └── db.go           # SQLite 连接
│   ├── game/               # 游戏服务
│   │   └── game_service.go # 核心逻辑
│   ├── models/             # 数据模型
│   │   └── models.go
│   ├── repository/         # 数据仓储
│   │   ├── game_repository.go
│   │   ├── user_repository.go
│   │   └── ai_repository.go
│   ├── websocket/          # WebSocket 管理
│   │   ├── handler.go
│   │   └── manager.go
│   └── main.go             # 入口文件
│
├── frontend/               # React 前端
│   ├── src/
│   │   ├── contexts/       # React Context
│   │   │   └── GameContext.jsx
│   │   ├── pages/         # 页面组件
│   │   │   ├── Login.jsx
│   │   │   ├── Register.jsx
│   │   │   ├── Home.jsx
│   │   │   ├── GameCreate.jsx
│   │   │   ├── GameJoin.jsx
│   │   │   ├── GameLobby.jsx
│   │   │   └── GameRoom.jsx
│   │   ├── services/       # API 服务
│   │   │   ├── api.js
│   │   │   └── websocket.js
│   │   └── test/           # 测试配置
│   │       └── setup.js
│   ├── package.json
│   └── vite.config.js
│
├── go.mod                  # Go 依赖
├── go.sum
├── package.json            # 前端依赖 (根目录)
└── README.md
```

## 🚀 快速开始

### 前置要求

- **Go**: 1.22+
- **Node.js**: 18+
- **npm** 或 **yarn**

### 安装步骤

1. **克隆项目**
```bash
git clone https://github.com/your-repo/Game-AI-Werewolf.git
cd Game-AI-Werewolf
```

2. **启动后端**
```bash
cd backend
go mod tidy
go build -o server.exe main.go
./server.exe
# 后端运行在 http://localhost:8080
```

3. **启动前端**
```bash
cd frontend
npm install
npm run dev
# 前端运行在 http://localhost:5173
```

4. **访问游戏**
打开浏览器访问 http://localhost:5173

## 🔌 API 接口

### 用户接口
| 方法 | 路径 | 描述 |
|------|------|------|
| POST | `/api/user/register` | 用户注册 |
| POST | `/api/user/login` | 用户登录 |
| GET | `/api/user/info/:user_id` | 获取用户信息 |

### 游戏接口
| 方法 | 路径 | 描述 |
|------|------|------|
| POST | `/api/game/create` | 创建游戏 |
| POST | `/api/game/join` | 加入游戏 |
| POST | `/api/game/start` | 开始游戏 |
| GET | `/api/game/info/:game_id` | 获取游戏信息 |
| GET | `/api/game/players/:game_id` | 获取玩家列表 |

### 阶段接口
| 方法 | 路径 | 描述 |
|------|------|------|
| GET | `/api/phase/current/:game_id` | 获取当前阶段 |
| POST | `/api/phase/end` | 结束当前阶段 |

### 玩家操作接口
| 方法 | 路径 | 描述 |
|------|------|------|
| POST | `/api/player/action/night` | 夜晚行动 |
| POST | `/api/player/action/day` | 白天发言 |
| POST | `/api/player/action/vote` | 投票 |

## 🧪 测试

### 后端测试
```bash
cd backend
go test -v ./...
```

### 前端测试
```bash
cd frontend
npm run test:run
```

### 测试覆盖率
| 测试类型 | 测试文件 | 测试数量 |
|----------|----------|----------|
| 后端 API | `backend/api/handlers_test.go` | 22 |
| 前端单元 | `frontend/src/gameLogic.test.js` | 26 |
| 前端组件 | `frontend/src/pages/GameRoom.test.jsx` | 9 |
| **总计** | | **57** |

## 🎮 游戏规则

### 角色配置 (12人局)
- **狼人** (4人): 夜间击杀村民
- **预言家** (1人): 夜间查验玩家身份
- **女巫** (1人): 夜间救人或毒人
- **猎人** (1人): 被放逐时带走一人
- **村民** (4人): 白天投票找出狼人

### 游戏流程
1. **夜晚阶段**: 狼人杀人、预言家查验、女巫救人/毒人
2. **白天阶段**: 玩家发言讨论
3. **投票阶段**: 票选可疑玩家放逐
4. 循环直到一方胜利

### 胜利条件
- **狼人胜利**: 狼人数量等于村民数量
- **村民胜利**: 所有狼人被放逐

## 📝 开发指南

### 数据库
- 使用 SQLite (纯 Go 实现，无需 CGO)
- 数据库文件: `backend/werewolf.db`
- 测试数据库: `backend/test.db`

### 环境变量
```go
// backend/main.go
port := 8080  // 服务端口
```

### 添加新功能
1. 在 `backend/api/handlers.go` 添加新的 handler
2. 在 `backend/game/game_service.go` 添加业务逻辑
3. 在 `frontend/src/services/api.js` 添加前端调用
4. 编写测试用例覆盖新功能

## 🤝 贡献指南

本项目为闭源项目，如需贡献代码或报告问题，请联系项目维护者。

## 📄 许可证

**版权所有 © 2024 蒙面狼人杀团队**

本项目为闭源软件。保留所有权利。

未经明确授权，禁止以下行为：
- 复制、修改、分发本软件
- 使用本软件进行商业活动
- 复制相关文档或代码
- 反向工程或逆向分析

如需授权或获取更多信息，请联系项目维护者。

## 📧 联系方式

如有问题，请联系项目维护者。

---

<p align="center">© 2024 蒙面狼人杀团队 - 版权所有</p>