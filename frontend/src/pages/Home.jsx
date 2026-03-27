import React from 'react'
import { useGame } from '../contexts/GameContext'

function Home({ onNavigate }) {
  const { user, logout } = useGame()

  return (
    <div className="card">
      <h2>欢迎，{user?.username}！</h2>
      <p>选择您要进行的操作：</p>
      
      <div style={{ display: 'flex', gap: '20px', marginTop: '20px' }}>
        <button className="btn" onClick={() => onNavigate('gameCreate')}>
          创建游戏
        </button>
        <button className="btn" onClick={() => onNavigate('gameJoin')}>
          加入游戏
        </button>
        <button 
          className="btn" 
          style={{ backgroundColor: '#f44336' }}
          onClick={logout}
        >
          退出登录
        </button>
      </div>
    </div>
  )
}

export default Home
