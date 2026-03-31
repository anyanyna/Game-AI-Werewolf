import React, { useState } from 'react'
import { useGame } from '../contexts/GameContext'

function GameJoin({ onNavigate }) {
  const [gameCode, setGameCode] = useState('')
  const { user, joinGame, setError, error, clearError, setIsLoading, isLoading } = useGame()

  const handleJoin = async (e) => {
    e.preventDefault()
    clearError()
    setIsLoading(true)
    
    try {
      await joinGame(gameCode, user.id)
      onNavigate('gameLobby')
    } catch (err) {
      setError(err.message || '加入游戏失败，请稍后重试')
    } finally {
      setIsLoading(false)
    }
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