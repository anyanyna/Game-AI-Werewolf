import React, { useState } from 'react'
import { useGame } from '../contexts/GameContext'

function GameCreate({ onNavigate }) {
  const { user, createGame, setError, error, clearError, setIsLoading, isLoading } = useGame()

  const handleCreate = () => {
    clearError()
    setIsLoading(true)
    // 模拟创建游戏请求
    fetch('http://localhost:8080/api/game/create', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ host_id: user.id }),
    })
      .then((response) => response.json())
      .then((data) => {
        if (data.error) {
          setError(data.error)
        } else {
          createGame(data.game)
          onNavigate('gameLobby')
        }
      })
      .catch((error) => {
        setError('创建游戏失败，请稍后重试')
        console.error('Error:', error)
      })
      .finally(() => {
        setIsLoading(false)
      })
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
