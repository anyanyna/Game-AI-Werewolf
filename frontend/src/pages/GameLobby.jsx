import React, { useEffect } from 'react'
import { useGame } from '../contexts/GameContext'

function GameLobby({ onNavigate }) {
  const { 
    game, 
    user, 
    players, 
    loadGameInfo, 
    startGame,
    setError, 
    error, 
    clearError, 
    setIsLoading, 
    isLoading 
  } = useGame()

  useEffect(() => {
    if (game) {
      loadGameInfo(game.id)
    }
  }, [game?.id, loadGameInfo])

  const handleStart = async () => {
    clearError()
    setIsLoading(true)
    
    try {
      await startGame(game.id)
      onNavigate('gameRoom')
    } catch (err) {
      setError(err.message || '开始游戏失败，请稍后重试')
    } finally {
      setIsLoading(false)
    }
  }

  const copyGameCode = () => {
    if (game?.game_code) {
      navigator.clipboard.writeText(game.game_code)
      alert('游戏代码已复制到剪贴板')
    }
  }

  return (
    <div className="card">
      <h2>游戏大厅</h2>
      <div style={{ marginBottom: '20px' }}>
        <p><strong>游戏代码：</strong>
          <span 
            style={{ cursor: 'pointer', color: '#2196F3' }} 
            onClick={copyGameCode}
            title="点击复制"
          >
            {game?.game_code || game?.gameCode || '加载中...'}
          </span>
        </p>
        <p><strong>游戏状态：</strong>{game?.status === 'waiting' ? '等待中' : '进行中'}</p>
        <p><strong>玩家数量：</strong>{players.length}/12</p>
      </div>

      <h3>玩家列表</h3>
      <ul style={{ listStyle: 'none', padding: 0 }}>
        {players.map((player) => (
          <li key={player.id} style={{ padding: '5px 0', borderBottom: '1px solid #eee' }}>
            #{player.number} - 
            {player.is_real_person || player.isRealPerson ? (
              <span style={{ color: '#4CAF50' }}>真人玩家</span>
            ) : (
              <span style={{ color: '#9E9E9E' }}>AI数字人</span>
            )}
          </li>
        ))}
      </ul>

      {error && <div className="error">{error}</div>}

      {game?.host_id === user?.id && (
        <button 
          className="btn" 
          onClick={handleStart} 
          disabled={isLoading || players.length < 1}
        >
          {isLoading ? '开始中...' : '开始游戏'}
        </button>
      )}

      <button 
        className="btn" 
        style={{ marginTop: '10px', backgroundColor: '#666' }} 
        onClick={() => onNavigate('home')}
      >
        离开游戏
      </button>
    </div>
  )
}

export default GameLobby