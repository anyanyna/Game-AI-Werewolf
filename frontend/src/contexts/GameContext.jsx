import React, { createContext, useContext, useState, useEffect } from 'react'

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
  // 状态管理
  const [user, setUser] = useState(null)
  const [game, setGame] = useState(null)
  const [players, setPlayers] = useState([])
  const [phase, setPhase] = useState(null)
  const [isLoading, setIsLoading] = useState(false)
  const [error, setError] = useState('')

  // 登录
  const login = (userData) => {
    setUser(userData)
  }

  // 注册
  const register = (userData) => {
    setUser(userData)
  }

  // 登出
  const logout = () => {
    setUser(null)
    setGame(null)
    setPlayers([])
    setPhase(null)
  }

  // 创建游戏
  const createGame = (gameData) => {
    setGame(gameData)
  }

  // 加入游戏
  const joinGame = (gameData) => {
    setGame(gameData)
  }

  // 加载游戏玩家
  const loadPlayers = (playersData) => {
    setPlayers(playersData)
  }

  // 加载游戏阶段
  const loadPhase = (phaseData) => {
    setPhase(phaseData)
  }

  // 清除错误
  const clearError = () => {
    setError('')
  }

  // 提供给子组件的值
  const value = {
    user,
    game,
    players,
    phase,
    isLoading,
    error,
    login,
    register,
    logout,
    createGame,
    joinGame,
    loadPlayers,
    loadPhase,
    clearError,
    setIsLoading,
    setError
  }

  return (
    <GameContext.Provider value={value}>
      {children}
    </GameContext.Provider>
  )
}

export default GameContext
