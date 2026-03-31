import React, { useState, useEffect, useCallback } from 'react'
import { useGame } from '../contexts/GameContext'
import api from '../services/api'
import LogPanel from '../components/LogPanel'

function GameRoom({ onNavigate }) {
  const { 
    game, 
    user, 
    player,
    players, 
    phase,
    loadGameInfo,
    refreshPlayers,
    refreshPhase,
    connectWebSocket,
    disconnectWebSocket,
    sendMessage,
    setError, 
    error, 
    clearError, 
    setIsLoading, 
    isLoading,
    wsConnected,
    gameLogs,
    playerLogs,
    loadGameLogs,
    loadPlayerLogs
  } = useGame()

  const [playerInfo, setPlayerInfo] = useState(null)
  const [actionTarget, setActionTarget] = useState(null)
  const [isActionSubmitted, setIsActionSubmitted] = useState(false)
  const [speechContent, setSpeechContent] = useState('')

  // 加载游戏数据
  useEffect(() => {
    if (game) {
      loadGameInfo(game.id)
    }
  }, [game?.id])

  // 设置当前玩家
  useEffect(() => {
    if (players.length > 0 && user) {
      const currentPlayer = players.find(p => p.user_id === user?.id || p.is_real_person === true)
      if (currentPlayer) {
        console.log('Found current player:', currentPlayer)
        setPlayerInfo(currentPlayer)
        // 加载玩家私人日志
        if (game) {
          loadGameLogs(game.id)
          loadPlayerLogs(game.id, currentPlayer.id)
        }
      }
    }
  }, [players, user, game?.id])

  // 连接 WebSocket
  useEffect(() => {
    if (game && playerInfo) {
      connectWebSocket(game.id, playerInfo.id, user.id)
        .then(() => {
          console.log('WebSocket connected')
        })
        .catch(err => {
          console.error('Failed to connect WebSocket:', err)
        })

      return () => {
        disconnectWebSocket()
      }
    }
  }, [game?.id, playerInfo?.id, user?.id])

  // 刷新日志
  useEffect(() => {
    if (game && playerInfo) {
      loadGameLogs(game.id)
      loadPlayerLogs(game.id, playerInfo.id)
    }
  }, [phase?.phase, phase?.round, game?.id, playerInfo?.id, loadGameLogs, loadPlayerLogs])

  const handleNightAction = async (actionType, targetID) => {
    clearError()
    setIsLoading(true)
    try {
      await api.player.nightAction(game.id, playerInfo.id, actionType, targetID)
      setIsActionSubmitted(true)
      sendMessage('night_action', { action_type: actionType, target_id: targetID })
    } catch (err) {
      setError(err.message || '行动失败，请稍后重试')
    } finally {
      setIsLoading(false)
    }
  }

  const handleDayAction = async () => {
    if (!speechContent.trim()) return
    clearError()
    setIsLoading(true)
    try {
      await api.player.dayAction(game.id, playerInfo.id, 'speak', speechContent, null)
      setIsActionSubmitted(true)  // 标记已发言
      sendMessage('day_action', { content: speechContent })
    } catch (err) {
      setError(err.message || '发言失败，请稍后重试')
    } finally {
      setIsLoading(false)
    }
  }

  const handleVote = async (targetID) => {
    clearError()
    setIsLoading(true)
    try {
      await api.player.vote(game.id, playerInfo.id, targetID)
      setIsActionSubmitted(true)
      sendMessage('vote', { target_id: targetID })
    } catch (err) {
      setError(err.message || '投票失败，请稍后重试')
    } finally {
      setIsLoading(false)
    }
  }

  const handleEndPhase = async () => {
    clearError()
    setIsLoading(true)
    try {
      const result = await api.phase.end(game.id)
      await refreshPhase(game.id)
      await refreshPlayers(game.id, playerInfo?.id)
      // 刷新日志
      if (game && playerInfo) {
        loadGameLogs(game.id)
        loadPlayerLogs(game.id, playerInfo.id)
      }
      sendMessage('end_phase', {})
      setIsActionSubmitted(false)
    } catch (err) {
      setError(err.message || '结束阶段失败，请稍后重试')
    } finally {
      setIsLoading(false)
    }
  }

  // 渲染夜晚阶段
  const renderNightPhase = () => {
    if (!playerInfo || playerInfo.status !== 'alive') {
      return (
        <div style={{
          padding: '30px',
          textAlign: 'center',
          background: '#f5f5f5',
          borderRadius: '10px'
        }}>
          <div style={{ fontSize: '40px', marginBottom: '10px' }}>💀</div>
          <div style={{ fontSize: '18px', color: '#666' }}>
            你已死亡，等待游戏结束...
          </div>
        </div>
      )
    }

    if (isActionSubmitted) {
      return (
        <div style={{
          padding: '30px',
          textAlign: 'center',
          background: '#e3f2fd',
          borderRadius: '10px'
        }}>
          <div style={{ fontSize: '40px', marginBottom: '10px' }}>⏳</div>
          <div style={{ fontSize: '18px', color: '#1976d2' }}>
            你的行动已提交，等待其他玩家行动...
          </div>
        </div>
      )
    }

    const role = playerInfo.role
    const aliveNonWolfPlayers = players.filter(p => p.status === 'alive' && p.role !== 'werewolf')
    const alivePlayers = players.filter(p => p.status === 'alive')

    switch (role) {
      case 'werewolf':
        // 获取其他狼人队友
        const wolfTeammates = players.filter(p => p.status === 'alive' && p.role === 'werewolf' && p.id !== playerInfo.id)
        return (
          <div>
            <div style={{
              padding: '15px',
              background: 'linear-gradient(135deg, #d32f2f 0%, #b71c1c 100%)',
              borderRadius: '10px',
              color: '#fff',
              marginBottom: '15px'
            }}>
              <div style={{ fontSize: '24px', marginBottom: '5px' }}>🐺 狼人行动</div>
              <div style={{ fontSize: '14px', opacity: 0.9 }}>
                选择要杀害的玩家。注意：只有狼人能看到你的选择
              </div>
              {wolfTeammates.length > 0 && (
                <div style={{ 
                  marginTop: '10px', 
                  padding: '10px', 
                  background: 'rgba(255,255,255,0.2)', 
                  borderRadius: '8px',
                  fontSize: '12px'
                }}>
                  <div>🐺 队友：{wolfTeammates.map(p => `#${p.number}`).join(', ')}</div>
                </div>
              )}
            </div>
            <p style={{ marginBottom: '10px', fontWeight: 'bold' }}>请选择要杀害的玩家：</p>
            <div style={{ display: 'grid', gridTemplateColumns: 'repeat(4, 1fr)', gap: '10px' }}>
              {aliveNonWolfPlayers.map((p) => (
                <button
                  key={p.id}
                  className="btn"
                  style={{
                    padding: '15px',
                    background: '#ffcdd2',
                    border: '2px solid #d32f2f'
                  }}
                  onClick={() => handleNightAction('kill', p.id)}
                >
                  <div style={{ fontSize: '20px' }}>{getRoleIcon(p.role)}</div>
                  <div>#{p.number}</div>
                </button>
              ))}
            </div>
          </div>
        )
      case 'seer':
        return (
          <div>
            <div style={{
              padding: '15px',
              background: 'linear-gradient(135deg, #7b1fa2 0%, #4a148c 100%)',
              borderRadius: '10px',
              color: '#fff',
              marginBottom: '15px'
            }}>
              <div style={{ fontSize: '24px', marginBottom: '5px' }}>🔮 预言家行动</div>
              <div style={{ fontSize: '14px', opacity: 0.9 }}>
                查验一名玩家的身份，了解其是好人还是狼人
              </div>
            </div>
            <p style={{ marginBottom: '10px', fontWeight: 'bold' }}>请选择要查验的玩家：</p>
            <div style={{ display: 'grid', gridTemplateColumns: 'repeat(4, 1fr)', gap: '10px' }}>
              {alivePlayers.filter(p => p.id !== playerInfo.id).map((p) => (
                <button
                  key={p.id}
                  className="btn"
                  style={{
                    padding: '15px',
                    background: '#e1bee7',
                    border: '2px solid #7b1fa2'
                  }}
                  onClick={() => handleNightAction('check', p.id)}
                >
                  <div style={{ fontSize: '20px' }}>{getRoleIcon(p.role)}</div>
                  <div>#{p.number}</div>
                </button>
              ))}
            </div>
          </div>
        )
      case 'witch':
        return (
          <div>
            <div style={{
              padding: '15px',
              background: 'linear-gradient(135deg, #388e3c 0%, #1b5e20 100%)',
              borderRadius: '10px',
              color: '#fff',
              marginBottom: '15px'
            }}>
              <div style={{ fontSize: '24px', marginBottom: '5px' }}>🧙‍♀️ 女巫行动</div>
              <div style={{ fontSize: '14px', opacity: 0.9 }}>
                你有一瓶解药和一瓶毒药，使用时请慎重
              </div>
            </div>
            {!actionTarget ? (
              <>
                <p style={{ marginBottom: '15px', fontWeight: 'bold' }}>选择你的行动：</p>
                <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '15px' }}>
                  <button 
                    className="btn" 
                    style={{ 
                      padding: '20px',
                      background: '#c8e6c9',
                      border: '2px solid #388e3c'
                    }}
                    onClick={() => setActionTarget('save')}
                  >
                    <div style={{ fontSize: '30px' }}>💊</div>
                    <div style={{ fontSize: '16px', fontWeight: 'bold' }}>救 人</div>
                    <div style={{ fontSize: '12px', color: '#666' }}>使用解药</div>
                  </button>
                  <button 
                    className="btn" 
                    style={{ 
                      padding: '20px',
                      background: '#ffccbc',
                      border: '2px solid #d84315'
                    }}
                    onClick={() => setActionTarget('poison')}
                  >
                    <div style={{ fontSize: '30px' }}>☠️</div>
                    <div style={{ fontSize: '16px', fontWeight: 'bold' }}>毒 人</div>
                    <div style={{ fontSize: '12px', color: '#666' }}>使用毒药</div>
                  </button>
                </div>
              </>
            ) : (
              <div style={{ marginTop: '20px' }}>
                <p style={{ marginBottom: '10px', fontWeight: 'bold' }}>
                  {actionTarget === 'save' ? '💊 选择要救的玩家：' : '☠️ 选择要毒的玩家：'}
                </p>
                <div style={{ display: 'grid', gridTemplateColumns: 'repeat(4, 1fr)', gap: '10px' }}>
                  {alivePlayers.map((p) => (
                    <button
                      key={p.id}
                      className="btn"
                      style={{
                        padding: '15px',
                        background: actionTarget === 'save' ? '#c8e6c9' : '#ffccbc',
                        border: `2px solid ${actionTarget === 'save' ? '#388e3c' : '#d84315'}`
                      }}
                      onClick={() => {
                        handleNightAction(actionTarget, p.id)
                        setActionTarget(null)
                      }}
                    >
                      <div style={{ fontSize: '20px' }}>{getRoleIcon(p.role)}</div>
                      <div>#{p.number}</div>
                    </button>
                  ))}
                </div>
                <button
                  className="btn"
                  style={{ 
                    marginTop: '15px', 
                    background: '#999',
                    padding: '10px 20px'
                  }}
                  onClick={() => setActionTarget(null)}
                >
                  ← 返回
                </button>
              </div>
            )}
          </div>
        )
      default:
        return (
          <div style={{
            padding: '30px',
            textAlign: 'center',
            background: '#f5f5f5',
            borderRadius: '10px'
          }}>
            <div style={{ fontSize: '40px', marginBottom: '10px' }}>🌙</div>
            <div style={{ fontSize: '18px', color: '#666' }}>
              夜晚阶段，请等待其他玩家行动...
            </div>
          </div>
        )
    }
  }

  // 渲染白天阶段
  const renderDayPhase = () => {
    if (!playerInfo || playerInfo.status !== 'alive') {
      return (
        <div style={{
          padding: '30px',
          textAlign: 'center',
          background: '#f5f5f5',
          borderRadius: '10px'
        }}>
          <div style={{ fontSize: '40px', marginBottom: '10px' }}>💀</div>
          <div style={{ fontSize: '18px', color: '#666' }}>
            你已死亡，等待游戏结束...
          </div>
        </div>
      )
    }

    return (
      <div>
        <div style={{
          padding: '15px',
          background: 'linear-gradient(135deg, #ff9800 0%, #f57c00 100%)',
          borderRadius: '10px',
          color: '#fff',
          marginBottom: '15px'
        }}>
          <div style={{ fontSize: '24px', marginBottom: '5px' }}>☀️ 白天阶段 - 发言</div>
          <div style={{ fontSize: '14px', opacity: 0.9 }}>
            分析局势，发言引导找出狼人
          </div>
        </div>
        <textarea
          style={{ 
            width: '100%', 
            height: '120px', 
            padding: '15px',
            borderRadius: '10px',
            border: '2px solid #ddd',
            fontSize: '14px',
            resize: 'none'
          }}
          placeholder="请输入你的发言... 分析局势、提出怀疑、或者为好人辩护"
          value={speechContent}
          onChange={(e) => setSpeechContent(e.target.value)}
        />
        <div style={{ display: 'flex', gap: '10px', marginTop: '10px' }}>
          <button 
            className="btn" 
            style={{ 
              flex: 1,
              padding: '15px',
              backgroundColor: '#4caf50'
            }}
            onClick={handleDayAction}
            disabled={isLoading || !speechContent.trim()}
          >
            📢 提交发言
          </button>
          <button 
            className="btn" 
            style={{ 
              padding: '15px',
              backgroundColor: '#999'
            }}
            onClick={() => setSpeechContent('')}
            disabled={isLoading || !speechContent}
          >
            🗑️ 清空
          </button>
        </div>
      </div>
    )
  }

  // 渲染投票阶段
  const renderVotingPhase = () => {
    if (!playerInfo || playerInfo.status !== 'alive') {
      return (
        <div style={{
          padding: '30px',
          textAlign: 'center',
          background: '#f5f5f5',
          borderRadius: '10px'
        }}>
          <div style={{ fontSize: '40px', marginBottom: '10px' }}>💀</div>
          <div style={{ fontSize: '18px', color: '#666' }}>
            你已死亡，等待游戏结束...
          </div>
        </div>
      )
    }

    if (isActionSubmitted) {
      return (
        <div style={{
          padding: '30px',
          textAlign: 'center',
          background: '#e3f2fd',
          borderRadius: '10px'
        }}>
          <div style={{ fontSize: '40px', marginBottom: '10px' }}>⏳</div>
          <div style={{ fontSize: '18px', color: '#1976d2' }}>
            你的投票已提交，等待其他玩家投票...
          </div>
        </div>
      )
    }

    return (
      <div>
        <div style={{
          padding: '15px',
          background: 'linear-gradient(135deg, #f44336 0%, #d32f2f 100%)',
          borderRadius: '10px',
          color: '#fff',
          marginBottom: '15px'
        }}>
          <div style={{ fontSize: '24px', marginBottom: '5px' }}>🗳️ 投票阶段</div>
          <div style={{ fontSize: '14px', opacity: 0.9 }}>
            选择要放逐的玩家，注意不要冤枉好人
          </div>
        </div>
        <p style={{ marginBottom: '15px', fontWeight: 'bold' }}>请选择要放逐的玩家：</p>
        <div style={{ display: 'grid', gridTemplateColumns: 'repeat(4, 1fr)', gap: '10px' }}>
          {players.filter(p => p.status === 'alive').map((p) => (
            <button
              key={p.id}
              className="btn"
              style={{
                padding: '15px',
                background: '#ffcdd2',
                border: '2px solid #d32f2f'
              }}
              onClick={() => handleVote(p.id)}
            >
              <div style={{ fontSize: '20px' }}>{getRoleIcon(p.role)}</div>
              <div>#{p.number}</div>
            </button>
          ))}
        </div>
      </div>
    )
  }

  // 获取角色名称
  const getRoleName = (role) => {
    if (role === 'unknown') return '???'
    const roleMap = {
      'werewolf': '狼人',
      'seer': '预言家',
      'witch': '女巫',
      'hunter': '猎人',
      'villager': '村民'
    }
    return roleMap[role] || role
  }

  // 获取角色图标
  const getRoleIcon = (role) => {
    if (role === 'unknown') return '❓'
    const iconMap = {
      'werewolf': '🐺',
      'seer': '🔮',
      'witch': '🧙‍♀️',
      'hunter': '🔫',
      'villager': '👤'
    }
    return iconMap[role] || '❓'
  }

  // 获取角色颜色
  const getRoleColor = (role) => {
    if (role === 'unknown') return '#999'
    const colorMap = {
      'werewolf': '#d32f2f',
      'seer': '#7b1fa2',
      'witch': '#388e3c',
      'hunter': '#f57c00',
      'villager': '#1976d2'
    }
    return colorMap[role] || '#666'
  }

  const getPhaseName = (phase) => {
    const phaseMap = {
      'night': '夜晚',
      'day': '白天',
      'voting': '投票'
    }
    return phaseMap[phase] || phase
  }

  // 渲染增强的玩家卡片
  const renderPlayerCard = (p) => {
    const isDead = p.status !== 'alive'
    const isSelf = p.id === playerInfo?.id
    
    // 获取死亡原因
    let deathReason = ''
    const killedInfo = gameLogs.phases?.find(phase => phase.killed === p.number)
    const poisonedInfo = gameLogs.phases?.find(phase => phase.poisoned === p.number)
    const votedInfo = gameLogs.phases?.find(phase => phase.voted_out === p.number)
    
    if (poisonedInfo) deathReason = '被毒杀'
    else if (killedInfo) deathReason = '被杀害'
    else if (votedInfo) deathReason = '被放逐'

    return (
      <div 
        key={p.id}
        style={{
          padding: '8px',
          background: isSelf ? '#fff3e0' : (isDead ? '#ffebee' : '#e8f5e9'),
          borderRadius: '6px',
          border: isSelf ? '2px solid #ff9800' : '1px solid #ddd',
          textAlign: 'center',
          boxShadow: isSelf ? '0 2px 8px rgba(255, 152, 0, 0.3)' : '0 1px 3px rgba(0,0,0,0.1)',
          cursor: 'pointer',
          position: 'relative'
        }}
      >
        {/* 角色图标 - 存活显示，或是自己死亡后也显示 */}
        {(isSelf || !isDead) && (
          <div style={{ fontSize: '20px', marginBottom: '2px' }}>
            {getRoleIcon(p.role)}
          </div>
        )}
        
        {/* 玩家编号 */}
        <div style={{ 
          fontWeight: 'bold', 
          fontSize: '14px',
          color: isDead ? '#999' : '#333'
        }}>
          #{p.number}
          {isSelf && <span style={{ color: '#ff9800', fontSize: '11px' }}> (你)</span>}
        </div>
        
        {/* 状态 */}
        <div style={{ 
          fontSize: '10px', 
          color: isDead ? '#d32f2f' : '#388e3c',
          marginTop: '2px'
        }}>
          {isDead ? '💀 死亡' : '✅ 存活'}
        </div>
        
        {/* 死亡原因（只显示原因，不显示角色名） */}
        {isDead && deathReason && (
          <div style={{
            fontSize: '10px',
            color: '#666',
            marginTop: '2px'
          }}>
            {deathReason}
          </div>
        )}
        
        {/* 真实玩家标识 */}
        {p.is_real_person && !isDead && (
          <div style={{
            position: 'absolute',
            top: '2px',
            right: '2px',
            fontSize: '9px',
            background: '#2196f3',
            color: '#fff',
            padding: '1px 3px',
            borderRadius: '3px'
          }}>
            真人
          </div>
        )}
      </div>
    )
  }

  // 渲染阶段进度条
  const renderPhaseProgress = () => {
    if (!phase) return null
    
    const phases = ['night', 'day', 'voting']
    const currentIndex = phases.indexOf(phase.phase)
    
    return (
      <div style={{
        marginBottom: '20px',
        padding: '15px',
        background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
        borderRadius: '10px',
        color: '#fff'
      }}>
        <div style={{ 
          display: 'flex', 
          justifyContent: 'space-between', 
          alignItems: 'center',
          marginBottom: '10px'
        }}>
          <div>
            <div style={{ fontSize: '12px', opacity: 0.8 }}>当前阶段</div>
            <div style={{ fontSize: '20px', fontWeight: 'bold' }}>
              {phase.phase === 'night' && '🌙 '}
              {phase.phase === 'day' && '☀️ '}
              {phase.phase === 'voting' && '🗳️ '}
              {getPhaseName(phase.phase)}
            </div>
          </div>
          <div style={{ textAlign: 'right' }}>
            <div style={{ fontSize: '12px', opacity: 0.8 }}>第 {phase.round} 轮</div>
          </div>
        </div>
        
        {/* 阶段进度指示器 */}
        <div style={{
          display: 'flex',
          gap: '10px',
          marginTop: '10px'
        }}>
          {phases.map((p, idx) => (
            <div
              key={p}
              style={{
                flex: 1,
                padding: '8px',
                background: idx === currentIndex ? 'rgba(255,255,255,0.3)' : 'rgba(255,255,255,0.1)',
                borderRadius: '6px',
                textAlign: 'center',
                fontSize: '12px',
                fontWeight: idx <= currentIndex ? 'bold' : 'normal',
                opacity: idx <= currentIndex ? 1 : 0.5
              }}
            >
              {p === 'night' && '🌙 夜晚'}
              {p === 'day' && '☀️ 白天'}
              {p === 'voting' && '🗳️ 投票'}
            </div>
          ))}
        </div>
      </div>
    )
  }

  return (
    <div style={{ display: 'flex', height: '100vh', overflow: 'hidden' }}>
      <div className="card" style={{ flex: 1, padding: '15px', display: 'flex', flexDirection: 'column' }}>
        <h2 style={{ textAlign: 'center', marginBottom: '10px', color: '#333', fontSize: '18px' }}>
          🎭 游戏房间
        </h2>
        
        {/* 游戏信息栏 */}
        <div style={{ 
          display: 'flex', 
          justifyContent: 'space-between',
          alignItems: 'center',
          marginBottom: '10px',
          padding: '8px 12px',
          background: '#f5f5f5',
          borderRadius: '6px'
        }}>
          <div>
            <span style={{ fontWeight: 'bold', fontSize: '12px' }}>游戏代码：</span>
            <span style={{ 
              fontFamily: 'monospace', 
              fontSize: '16px',
              letterSpacing: '2px',
              color: '#1976d2'
            }}>
              {game?.game_code || game?.gameCode}
            </span>
          </div>
          <div style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
            <span style={{
              padding: '3px 8px',
              borderRadius: '10px',
              fontSize: '11px',
              background: game?.status === 'playing' ? '#4caf50' : '#999',
              color: '#fff'
            }}>
              {game?.status === 'playing' ? '进行中' : '已结束'}
            </span>
            <span style={{ 
              fontSize: '11px',
              color: wsConnected ? '#4caf50' : '#f44336'
            }}>
              {wsConnected ? '🟢' : '🔴'}
            </span>
          </div>
        </div>

        {/* 阶段进度条 - 紧凑版 */}
        {phase && (
          <div style={{
            marginBottom: '10px',
            padding: '10px',
            background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
            borderRadius: '8px',
            color: '#fff'
          }}>
            <div style={{ 
              display: 'flex', 
              justifyContent: 'space-between', 
              alignItems: 'center'
            }}>
              <div style={{ fontSize: '16px', fontWeight: 'bold' }}>
                {phase.phase === 'night' && '🌙 '}
                {phase.phase === 'day' && '☀️ '}
                {phase.phase === 'voting' && '🗳️ '}
                {getPhaseName(phase.phase)}
              </div>
              <div style={{ fontSize: '12px' }}>第 {phase.round} 轮</div>
            </div>
            <div style={{
              display: 'flex',
              gap: '8px',
              marginTop: '8px'
            }}>
              {['night', 'day', 'voting'].map((p, idx) => (
                <div
                  key={p}
                  style={{
                    flex: 1,
                    padding: '6px',
                    background: idx === ['night', 'day', 'voting'].indexOf(phase.phase) ? 'rgba(255,255,255,0.3)' : 'rgba(255,255,255,0.1)',
                    borderRadius: '4px',
                    textAlign: 'center',
                    fontSize: '11px'
                  }}
                >
                  {p === 'night' && '🌙 夜'}
                  {p === 'day' && '☀️ 昼'}
                  {p === 'voting' && '🗳️ 投'}
                </div>
              ))}
            </div>
          </div>
        )}

        {/* 玩家列表 - 紧凑布局 */}
        <div style={{ marginBottom: '10px', flex: '0 0 auto' }}>
          <div style={{ 
            marginBottom: '8px', 
            fontSize: '13px',
            color: '#666'
          }}>
            👥 玩家 ({players.filter(p => p.status === 'alive').length}/{players.length})
          </div>
          <div style={{ 
            display: 'grid', 
            gridTemplateColumns: 'repeat(6, 1fr)', 
            gap: '6px'
          }}>
            {players.map(renderPlayerCard)}
          </div>
        </div>

        {/* 游戏操作区域 - 紧凑 */}
        <div style={{ 
          flex: 1,
          background: '#fafafa', 
          borderRadius: '8px',
          padding: '12px',
          overflow: 'auto'
        }}>
          {phase && (
            <div>
              {phase.phase === 'night' && renderNightPhase()}
              {phase.phase === 'day' && renderDayPhase()}
              {phase.phase === 'voting' && renderVotingPhase()}
              
              {/* 根据阶段显示不同的按钮文字 */}
              <button 
                className="btn" 
                style={{ 
                  marginTop: '10px', 
                  backgroundColor: '#FF9800',
                  padding: '10px 20px',
                  fontSize: '13px'
                }}
                onClick={handleEndPhase}
                disabled={isLoading}
              >
                {isLoading ? '处理中...' : 
                  phase.phase === 'night' ? '🌙 结束夜晚 → 进入白天' :
                  phase.phase === 'day' ? '☀️ 结束发言 → 进入投票' :
                  '🗳️ 结束投票 → 进入黑夜'
                }
              </button>
            </div>
          )}
        </div>

        {error && (
          <div style={{
            padding: '8px',
            background: '#ffebee',
            borderRadius: '6px',
            color: '#d32f2f',
            fontSize: '12px',
            marginTop: '8px'
          }}>
            ⚠️ {error}
          </div>
        )}

        <button 
          className="btn" 
          style={{ 
            marginTop: '8px', 
            backgroundColor: '#666',
            padding: '8px 16px',
            fontSize: '12px'
          }} 
          onClick={() => onNavigate('home')}
        >
          🚪 退出
        </button>
      </div>
      <LogPanel />
    </div>
  )
}

export default GameRoom