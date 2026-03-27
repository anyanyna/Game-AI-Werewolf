import React, { useEffect } from 'react'
import { useGame } from '../contexts/GameContext'

function GameLobby({ onNavigate }) {
  const { game, user, players, loadPlayers, setError, error, clearError, setIsLoading, isLoading } = useGame()

  useEffect(() => {
    // 加载游戏玩家列表
    if (game) {
      clearError()
      setIsLoading(true)
      fetch(`http://localhost:8080/api/game/players/${game.id}`)
        .then((response) => response.json())
        .then((data) => {
          if (data.error) {
            setError(data.error)
          } else {
            loadPlayers(data.players)
          }
        })
        .catch((error) => {
          console.error('Error:', error)
        })
        .finally(() => {
          setIsLoading(false)
        })
    }
  }, [game, loadPlayers, setError, clearError, setIsLoading])

  const handleStart = () => {
    clearError()
    setIsLoading(true)
    // 模拟开始游戏请求
    fetch('http://localhost:8080/api/game/start', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ game_id: game.id }),
    })
      .then((response) => response.json())
      .then((data) => {
        if (data.error) {
          setError(data.error)
        } else {
          onNavigate('gameRoom')
        }
      })
      .catch((error) => {
        setError('开始游戏失败，请稍后重试')
        console.error('Error:', error)
      })
      .finally(() => {
        setIsLoading(false)
      })
  }

  return (
    <div className="card">
      <h2>游戏大厅</h2>
      <div style={{ marginBottom: '20px' }}>
        <p><strong>游戏代码：</strong>{game?.gameCode}</p>
        <p><strong>游戏状态：</strong>{game?.status === 'waiting' ? '等待中' : '进行中'}</p>
        <p><strong>玩家数量：</strong>{players.length}/12</p>
      </div>

      <h3>玩家列表</h3>
      <ul style={{ listStyle: 'none', padding: 0 }}>
        {players.map((player) => (
          <li key={player.id} style={{ padding: '5px 0', borderBottom: '1px solid #eee' }}>
            {player.isRealPerson ? (
              <span>真人玩家 #{player.number}</span>
            ) : (
              <span>AI数字人 #{player.number}</span>
            )}
          </li>
        ))}
      </ul>

      {error && <div className="error">{error}</div>}

      {game?.hostID === user?.id && (
        <button className="btn" onClick={handleStart} disabled={isLoading}>
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
