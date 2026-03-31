import React, { createContext, useContext, useState, useCallback, useRef } from 'react'
import webSocketService from '../services/websocket'
import api from '../services/api'

// 创建游戏状态Context
const GameContext = createContext()

// 自定义Hook，用于使用游戏状态
export const useGame = () => {
  const context = useContext(GameContext)
  if (!context) {
    throw new Error('useGame must be used within a GameProvider')
  }
  return context
}

// 游戏状态Provider组件
export const GameProvider = ({ children }) => {
  // 从 localStorage 恢复状态
  const getInitialUser = () => {
    try {
      const saved = localStorage.getItem('werewolf_user')
      return saved ? JSON.parse(saved) : null
    } catch {
      return null
    }
  }

  const getInitialGame = () => {
    try {
      const saved = localStorage.getItem('werewolf_game')
      return saved ? JSON.parse(saved) : null
    } catch {
      return null
    }
  }

  // 状态管理
  const [user, setUser] = useState(getInitialUser)
  const [game, setGame] = useState(getInitialGame)
  const [player, setPlayer] = useState(null)
  const [players, setPlayers] = useState([])
  const [phase, setPhase] = useState(null)
  const [isLoading, setIsLoading] = useState(false)
  const [error, setError] = useState('')
  const [wsConnected, setWsConnected] = useState(false)
  const [gameLogs, setGameLogs] = useState({ phases: [], players: [] })
  const [playerLogs, setPlayerLogs] = useState({ role: '', playerNumber: 0, privateLogs: [] })
  const [isHydrated, setIsHydrated] = useState(false) // 标记是否已从localStorage恢复
  
  const wsRef = useRef(null)

  // 用户状态变化时持久化
  React.useEffect(() => {
    if (user) {
      localStorage.setItem('werewolf_user', JSON.stringify(user))
    } else {
      localStorage.removeItem('werewolf_user')
    }
  }, [user])

  // 游戏状态变化时持久化
  React.useEffect(() => {
    if (game) {
      localStorage.setItem('werewolf_game', JSON.stringify(game))
    } else {
      localStorage.removeItem('werewolf_game')
    }
  }, [game])

  // 从localStorage恢复后，同步游戏状态
  React.useEffect(() => {
    if (user && game && game.id) {
      // 异步刷新游戏数据，确保状态最新
      const syncGame = async () => {
        try {
          const gameResult = await api.game.getInfo(game.id)
          setGame(gameResult.game)
          // 如果游戏正在进行中，也刷新玩家和阶段
          if (gameResult.game.status === 'playing') {
            const playersResult = await api.game.getPlayers(game.id)
            setPlayers(playersResult.players)
            const phaseResult = await api.phase.getCurrent(game.id).catch(() => null)
            if (phaseResult) {
              setPhase(phaseResult.phase)
            }
          }
        } catch (err) {
          console.error('Failed to sync game state:', err)
          // 如果同步失败，清除游戏状态
          setGame(null)
          localStorage.removeItem('werewolf_game')
        } finally {
          setIsHydrated(true)
        }
      }
      syncGame()
    } else {
      setIsHydrated(true)
    }
  }, [user]) // 当用户恢复时触发

  // 登录
  const login = useCallback(async (username, password) => {
    setIsLoading(true)
    setError('')
    try {
      const result = await api.user.login(username, password)
      setUser(result.user)
      return result
    } catch (err) {
      setError(err.message)
      throw err
    } finally {
      setIsLoading(false)
    }
  }, [])

  // 注册
  const register = useCallback(async (username, email, password) => {
    setIsLoading(true)
    setError('')
    try {
      const result = await api.user.register(username, email, password)
      setUser(result.user)
      return result
    } catch (err) {
      setError(err.message)
      throw err
    } finally {
      setIsLoading(false)
    }
  }, [])

  // 登出
  const logout = useCallback(() => {
    if (wsRef.current) {
      webSocketService.disconnect()
      wsRef.current = null
    }
    setUser(null)
    setGame(null)
    setPlayer(null)
    setPlayers([])
    setPhase(null)
    setWsConnected(false)
  }, [])

  // 创建游戏
  const createGame = useCallback(async (hostId, playerCount = 12) => {
    setIsLoading(true)
    setError('')
    try {
      const result = await api.game.create(hostId, playerCount)
      setGame(result.game)
      return result
    } catch (err) {
      setError(err.message)
      throw err
    } finally {
      setIsLoading(false)
    }
  }, [])

  // 加入游戏
  const joinGame = useCallback(async (gameCode, userId) => {
    setIsLoading(true)
    setError('')
    try {
      const result = await api.game.join(gameCode, userId)
      setGame(result.game)
      setPlayer(result.player)
      return result
    } catch (err) {
      setError(err.message)
      throw err
    } finally {
      setIsLoading(false)
    }
  }, [])

  // 开始游戏
  const startGame = useCallback(async (gameId) => {
    setIsLoading(true)
    setError('')
    try {
      const result = await api.game.start(gameId)
      setGame(result.game)
      // 刷新玩家列表
      const playersResult = await api.game.getPlayers(gameId)
      setPlayers(playersResult.players)
      return result
    } catch (err) {
      setError(err.message)
      throw err
    } finally {
      setIsLoading(false)
    }
  }, [])

  // 加载游戏信息
  const loadGameInfo = useCallback(async (gameId) => {
    try {
      const [gameResult, playersResult, phaseResult] = await Promise.all([
        api.game.getInfo(gameId),
        api.game.getPlayers(gameId),
        gameId ? api.phase.getCurrent(gameId).catch(() => null) : null,
      ])
      setGame(gameResult.game)
      setPlayers(playersResult.players)
      if (phaseResult) {
        setPhase(phaseResult.phase)
      }
    } catch (err) {
      console.error('Failed to load game info:', err)
    }
  }, [])

  // 刷新玩家列表
  const refreshPlayers = useCallback(async (gameId, playerId = null) => {
    try {
      // 如果有 playerId，传递给后端用于角色过滤
      const result = playerId 
        ? await api.game.getPlayers(gameId, playerId)
        : await api.game.getPlayers(gameId)
      setPlayers(result.players)
    } catch (err) {
      console.error('Failed to refresh players:', err)
    }
  }, [])

  // 刷新阶段
  const refreshPhase = useCallback(async (gameId) => {
    try {
      const result = await api.phase.getCurrent(gameId)
      setPhase(result.phase)
    } catch (err) {
      console.error('Failed to refresh phase:', err)
    }
  }, [])

  // 连接 WebSocket
  const connectWebSocket = useCallback((gameId, playerId, userId) => {
    return new Promise((resolve, reject) => {
      try {
        wsRef.current = webSocketService.connect(gameId, playerId, userId)
        
        // 监听连接事件
        webSocketService.on('connect', () => {
          console.log('WebSocket connected')
          setWsConnected(true)
          resolve()
        })

        webSocketService.on('disconnect', () => {
          console.log('WebSocket disconnected')
          setWsConnected(false)
        })

        webSocketService.on('error', (err) => {
          console.error('WebSocket error:', err)
          setError(err.message || 'WebSocket error')
        })

        // 监听游戏事件
        webSocketService.on('player_joined', (data) => {
          refreshPlayers(gameId)
        })

        webSocketService.on('phase_changed', (data) => {
          setPhase(data.data.phase)
        })

        webSocketService.on('player_spoke', (data) => {
          // 可以添加发言日志
          console.log('Player spoke:', data)
        })

        webSocketService.on('player_voted', (data) => {
          console.log('Player voted:', data)
        })

        webSocketService.on('game_ended', (data) => {
          setGame(prev => ({ ...prev, status: 'ended', winner: data.data.winner }))
        })

      } catch (err) {
        reject(err)
      }
    })
  }, [refreshPlayers])

  // 断开 WebSocket
  const disconnectWebSocket = useCallback(() => {
    if (wsRef.current) {
      webSocketService.disconnect()
      wsRef.current = null
      setWsConnected(false)
    }
  }, [])

  // 发送 WebSocket 消息
  const sendMessage = useCallback((type, data) => {
    if (webSocketService.isConnected()) {
      webSocketService.send(type, data)
    } else {
      console.warn('WebSocket not connected')
    }
  }, [])

  // 清除错误
  const clearError = useCallback(() => {
    setError('')
  }, [])

  // 加载游戏公共日志
  const loadGameLogs = useCallback(async (gameId) => {
    try {
      const result = await api.log.getGameLogs(gameId)
      setGameLogs(result)
    } catch (err) {
      console.error('Failed to load game logs:', err)
    }
  }, [])

  // 加载玩家私人日志
  const loadPlayerLogs = useCallback(async (gameId, playerId) => {
    try {
      const result = await api.log.getPlayerLogs(gameId, playerId)
      setPlayerLogs(result)
    } catch (err) {
      console.error('Failed to load player logs:', err)
    }
  }, [])

  // 提供给子组件的值
  const value = {
    user,
    game,
    player,
    players,
    phase,
    isLoading,
    isHydrated, // 新增：标记是否已完成状态恢复
    error,
    wsConnected,
    gameLogs,
    playerLogs,
    login,
    register,
    logout,
    createGame,
    joinGame,
    startGame,
    loadGameInfo,
    refreshPlayers,
    refreshPhase,
    connectWebSocket,
    disconnectWebSocket,
    sendMessage,
    clearError,
    setIsLoading,
    setError,
    loadGameLogs,
    loadPlayerLogs
  }

  return (
    <GameContext.Provider value={value}>
      {children}
    </GameContext.Provider>
  )
}

export default GameContext