import React, { useState } from 'react'
import { useGame } from '../contexts/GameContext'

function GameCreate({ onNavigate }) {
  const { user, createGame, setError, error, clearError, setIsLoading, isLoading } = useGame()

  const handleCreate = async () => {
    clearError()
    setIsLoading(true)
    
    try {
      const result = await createGame(user.id, 12)
      onNavigate('gameLobby')
    } catch (err) {
      setError(err.message || '创建游戏失败，请稍后重试')
    } finally {
      setIsLoading(false)
    }
  }

  return (
    <div className="card">
      <h2>创建游戏</h2>
      <p>点击下方按钮创建一个新的狼人杀游戏：</p>
      {error && <div className="error">{error}</div>}
      <button className="btn" onClick={handleCreate} disabled={isLoading}>
        {isLoading ? '创建中...' : '创建游戏'}
      </button>
      <button 
        className="btn" 
        style={{ marginTop: '10px', backgroundColor: '#666' }} 
        onClick={() => onNavigate('home')}
      >
        返回主页
      </button>
    </div>
  )
}

export default GameCreate