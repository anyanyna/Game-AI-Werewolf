import React, { useState, useEffect } from 'react'
import { useGame } from '../contexts/GameContext'
import webSocketService from '../services/websocket'

function GameRoom({ onNavigate }) {
  const { game, user, players, phase, loadPlayers, loadPhase, setError, error, clearError, setIsLoading, isLoading } = useGame()
  const [playerInfo, setPlayerInfo] = useState(null)
  const [actionTarget, setActionTarget] = useState(null)
  const [isActionSubmitted, setIsActionSubmitted] = useState(false)
  const [socketConnected, setSocketConnected] = useState(false)

  useEffect(() => {
    // 加载当前游戏阶段
    if (game) {
      loadCurrentPhase()
      loadPlayersData()
    }
  }, [game])

  useEffect(() => {
    // 连接WebSocket
    if (game && playerInfo) {
      const socket = webSocketService.connect(game.id, playerInfo.id, user.id)
      setSocketConnected(true)

      // 监听游戏状态更新
      webSocketService.on('game_state_update', (data) => {
        console.log('Game state updated:', data)
        loadCurrentPhase()
        loadPlayersData()
      })

      // 监听玩家行动
      webSocketService.on('player_action', (data) => {
        console.log('Player action:', data)
      })

      // 监听阶段结束
      webSocketService.on('phase_end', (data) => {
        console.log('Phase ended:', data)
        loadCurrentPhase()
        loadPlayersData()
      })

      // 监听游戏结束
      webSocketService.on('game_end', (data) => {
        console.log('Game ended:', data)
        alert(`游戏结束！获胜方：${data.winner === 'werewolf' ? '狼人' : '好人'}`)
      })

      // 清理函数
      return () => {
        webSocketService.disconnect()
        setSocketConnected(false)
      }
    }
  }, [game, playerInfo, user.id])

  const loadCurrentPhase = () => {
    clearError()
    setIsLoading(true)
    fetch(`http://localhost:8080/api/phase/current/${game.id}`)
      .then((response) => response.json())
      .then((data) => {
        if (data.error) {
          setError(data.error)
        } else {
          loadPhase(data.phase)
          setIsActionSubmitted(false)
        }
      })
      .catch((error) => {
        console.error('Error:', error)
      })
      .finally(() => {
        setIsLoading(false)
      })
  }

  const loadPlayersData = () => {
    clearError()
    setIsLoading(true)
    fetch(`http://localhost:8080/api/game/players/${game.id}`)
      .then((response) => response.json())
      .then((data) => {
        if (data.error) {
          setError(data.error)
        } else {
          loadPlayers(data.players)
          // 找到当前用户对应的玩家信息
          const currentPlayer = data.players.find(p => p.userID === user.id)
          setPlayerInfo(currentPlayer)
        }
      })
      .catch((error) => {
        console.error('Error:', error)
      })
      .finally(() => {
        setIsLoading(false)
      })
  }

  const handleNightAction = (actionType, targetID) => {
    clearError()
    setIsLoading(true)
    // 模拟夜晚行动请求
    fetch('http://localhost:8080/api/player/action/night', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        game_id: game.id,
        player_id: playerInfo.id,
        action_type: actionType,
        target_id: targetID
      }),
    })
      .then((response) => response.json())
      .then((data) => {
        if (data.error) {
          setError(data.error)
        } else {
          setIsActionSubmitted(true)
          alert('行动已提交')
          // 发送WebSocket消息
          if (webSocketService.isConnected()) {
            webSocketService.emit('night_action', {
              game_id: game.id,
              player_id: playerInfo.id,
              action_type: actionType,
              target_id: targetID
            })
          }
        }
      })
      .catch((error) => {
        setError('行动失败，请稍后重试')
        console.error('Error:', error)
      })
      .finally(() => {
        setIsLoading(false)
      })
  }

  const handleDayAction = (actionType, content, targetID) => {
    clearError()
    setIsLoading(true)
    // 模拟白天行动请求
    fetch('http://localhost:8080/api/player/action/day', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        game_id: game.id,
        player_id: playerInfo.id,
        action_type: actionType,
        content: content,
        target_id: targetID
      }),
    })
      .then((response) => response.json())
      .then((data) => {
        if (data.error) {
          setError(data.error)
        } else {
          alert('行动已提交')
          // 发送WebSocket消息
          if (webSocketService.isConnected()) {
            webSocketService.emit('day_action', {
              game_id: game.id,
              player_id: playerInfo.id,
              action_type: actionType,
              content: content,
              target_id: targetID
            })
          }
        }
      })
      .catch((error) => {
        setError('行动失败，请稍后重试')
        console.error('Error:', error)
      })
      .finally(() => {
        setIsLoading(false)
      })
  }

  const handleVote = (targetID) => {
    clearError()
    setIsLoading(true)
    // 模拟投票请求
    fetch('http://localhost:8080/api/player/action/vote', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        game_id: game.id,
        player_id: playerInfo.id,
        target_id: targetID
      }),
    })
      .then((response) => response.json())
      .then((data) => {
        if (data.error) {
          setError(data.error)
        } else {
          setIsActionSubmitted(true)
          alert('投票已提交')
          // 发送WebSocket消息
          if (webSocketService.isConnected()) {
            webSocketService.emit('vote', {
              game_id: game.id,
              player_id: playerInfo.id,
              target_id: targetID
            })
          }
        }
      })
      .catch((error) => {
        setError('投票失败，请稍后重试')
        console.error('Error:', error)
      })
      .finally(() => {
        setIsLoading(false)
      })
  }

  const handleEndPhase = () => {
    clearError()
    setIsLoading(true)
    // 模拟结束阶段请求
    fetch('http://localhost:8080/api/phase/end', {
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
          // 重新加载游戏阶段
          loadCurrentPhase()
          loadPlayersData()
          // 发送WebSocket消息
          if (webSocketService.isConnected()) {
            webSocketService.emit('end_phase', {
              game_id: game.id
            })
          }
        }
      })
      .catch((error) => {
        setError('结束阶段失败，请稍后重试')
        console.error('Error:', error)
      })
      .finally(() => {
        setIsLoading(false)
      })
  }

  const renderNightPhase = () => {
    if (!playerInfo || playerInfo.status !== 'alive') {
      return <p>你已死亡，等待游戏结束...</p>
    }

    if (isActionSubmitted) {
      return <p>你的行动已提交，等待其他玩家行动...</p>
    }

    switch (playerInfo.role) {
      case 'werewolf':
        return (
          <div>
            <h4>狼人行动</h4>
            <p>请选择要杀害的玩家：</p>
            <div style={{ display: 'grid', gridTemplateColumns: 'repeat(4, 1fr)', gap: '10px' }}>
              {players.filter(p => p.status === 'alive' && p.role !== 'werewolf').map((p) => (
                <button
                  key={p.id}
                  className="btn"
                  onClick={() => handleNightAction('kill', p.id)}
                >
                  玩家 #{p.number}
                </button>
              ))}
            </div>
          </div>
        )
      case 'seer':
        return (
          <div>
            <h4>预言家行动</h4>
            <p>请选择要查验的玩家：</p>
            <div style={{ display: 'grid', gridTemplateColumns: 'repeat(4, 1fr)', gap: '10px' }}>
              {players.filter(p => p.status === 'alive' && p.id !== playerInfo.id).map((p) => (
                <button
                  key={p.id}
                  className="btn"
                  onClick={() => handleNightAction('check', p.id)}
                >
                  玩家 #{p.number}
                </button>
              ))}
            </div>
          </div>
        )
      case 'witch':
        return (
          <div>
            <h4>女巫行动</h4>
            <p>选择你的行动：</p>
            <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '10px' }}>
              <button className="btn" onClick={() => setActionTarget('save')}>
                救 人
              </button>
              <button className="btn" onClick={() => setActionTarget('poison')}>
                毒 人
              </button>
            </div>
            {actionTarget && (
              <div style={{ marginTop: '20px' }}>
                <p>请选择目标：</p>
                <div style={{ display: 'grid', gridTemplateColumns: 'repeat(4, 1fr)', gap: '10px' }}>
                  {players.filter(p => p.status === 'alive').map((p) => (
                    <button
                      key={p.id}
                      className="btn"
                      onClick={() => {
                        handleNightAction(actionTarget, p.id)
                        setActionTarget(null)
                      }}
                    >
                      玩家 #{p.number}
                    </button>
                  ))}
                </div>
              </div>
            )}
          </div>
        )
      case 'hunter':
      case 'villager':
        return <p>夜晚阶段，请等待其他玩家行动...</p>
      default:
        return <p>夜晚阶段，请等待...</p>
    }
  }

  const renderDayPhase = () => {
    if (!playerInfo || playerInfo.status !== 'alive') {
      return <p>你已死亡，等待游戏结束...</p>
    }

    return (
      <div>
        <h4>白天阶段</h4>
        <p>玩家发言中...</p>
        <div style={{ marginTop: '20px' }}>
          <h5>你的发言：</h5>
          <textarea
            style={{ width: '100%', height: '100px', padding: '10px' }}
            placeholder="请输入你的发言..."
            id="speech"
          ></textarea>
          <button 
            className="btn" 
            style={{ marginTop: '10px' }}
            onClick={() => {
              const speech = document.getElementById('speech').value
              handleDayAction('speak', speech, null)
            }}
          >
            提交发言
          </button>
        </div>
      </div>
    )
  }

  const renderVotingPhase = () => {
    if (!playerInfo || playerInfo.status !== 'alive') {
      return <p>你已死亡，等待游戏结束...</p>
    }

    if (isActionSubmitted) {
      return <p>你的投票已提交，等待其他玩家投票...</p>
    }

    return (
      <div>
        <h4>投票阶段</h4>
        <p>请选择要放逐的玩家：</p>
        <div style={{ display: 'grid', gridTemplateColumns: 'repeat(4, 1fr)', gap: '10px' }}>
          {players.filter(p => p.status === 'alive').map((p) => (
            <button
              key={p.id}
              className="btn"
              onClick={() => handleVote(p.id)}
            >
              玩家 #{p.number}
            </button>
          ))}
        </div>
      </div>
    )
  }

  return (
    <div className="card">
      <h2>游戏房间</h2>
      <div style={{ marginBottom: '20px' }}>
        <p><strong>游戏代码：</strong>{game?.gameCode}</p>
        <p><strong>游戏状态：</strong>{game?.status === 'playing' ? '进行中' : '已结束'}</p>
        <p><strong>WebSocket连接：</strong>{socketConnected ? '已连接' : '未连接'}</p>
        {phase && (
          <>
            <p><strong>当前阶段：</strong>{phase.phase === 'night' ? '夜晚' : phase.phase === 'day' ? '白天' : '投票'}</p>
            <p><strong>当前回合：</strong>{phase.round}</p>
          </>
        )}
      </div>

      {playerInfo && (
        <div style={{ marginBottom: '20px' }}>
          <h3>你的角色</h3>
          <p><strong>编号：</strong>{playerInfo.number}</p>
          <p><strong>角色：</strong>{playerInfo.role === 'werewolf' ? '狼人' : playerInfo.role === 'seer' ? '预言家' : playerInfo.role === 'witch' ? '女巫' : playerInfo.role === 'hunter' ? '猎人' : '村民'}</p>
          <p><strong>状态：</strong>{playerInfo.status === 'alive' ? '存活' : '死亡'}</p>
        </div>
      )}

      <h3>游戏操作</h3>
      {phase && (
        <div>
          {phase.phase === 'night' && renderNightPhase()}
          {phase.phase === 'day' && renderDayPhase()}
          {phase.phase === 'voting' && renderVotingPhase()}
          <button 
            className="btn" 
            style={{ marginTop: '20px' }}
            onClick={handleEndPhase}
            disabled={isLoading}
          >
            {isLoading ? '处理中...' : '结束当前阶段'}
          </button>
        </div>
      )}

      {error && <div className="error">{error}</div>}

      <button 
        className="btn" 
        style={{ marginTop: '10px', backgroundColor: '#666' }} 
        onClick={() => onNavigate('home')}
      >
        退出游戏
      </button>
    </div>
  )
}

export default GameRoom
