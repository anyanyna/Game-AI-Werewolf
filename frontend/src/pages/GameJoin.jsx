import React, { useState } from 'react'
import { useGame } from '../contexts/GameContext'

function GameJoin({ onNavigate }) {
  const [gameCode, setGameCode] = useState('')
  const { user, joinGame, setError, error, clearError, setIsLoading, isLoading } = useGame()

  const handleJoin = (e) => {
    e.preventDefault()
    clearError()
    setIsLoading(true)
    
    // 模拟加入游戏请求
    fetch('http://localhost:8080/api/game/join', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ game_code: gameCode, user_id: user.id }),
    })
      .then((response) => response.json())
      .then((data) => {
        if (data.error) {
          setError(data.error)
        } else {
          joinGame(data.game)
          onNavigate('gameLobby')
        }
      })
      .catch((error) => {
        setError('加入游戏失败，请稍后重试')
        console.error('Error:', error)
      })
      .finally(() => {
        setIsLoading(false)
      })
  }

  return (
    <div className="card">
      <h2>加入游戏</h2>
      <form onSubmit={handleJoin}>
        <div className="form-group">
          <label htmlFor="gameCode">游戏代码</label>
          <input
            type="text"
            id="gameCode"
            value={gameCode}
            onChange={(e) => setGameCode(e.target.value)}
            required
          />
        </div>
        {error && <div className="error">{error}</div>}
        <button type="submit" className="btn" disabled={isLoading}>
          {isLoading ? '加入中...' : '加入游戏'}
        </button>
      </form>
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

export default GameJoin
