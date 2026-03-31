import React, { useState, useEffect, useRef } from 'react'
import { GameProvider, useGame } from './contexts/GameContext'
import Login from './pages/Login.jsx'
import Register from './pages/Register.jsx'
import Home from './pages/Home.jsx'
import GameCreate from './pages/GameCreate.jsx'
import GameJoin from './pages/GameJoin.jsx'
import GameLobby from './pages/GameLobby.jsx'
import GameRoom from './pages/GameRoom.jsx'

// 内部组件，需要访问 GameContext
function AppContent({ onNavigate }) {
  const { user, game, isHydrated } = useGame()
  const [currentPage, setCurrentPage] = useState('login')
  const targetPageRef = useRef(null)

  // 根据登录状态和游戏状态决定初始页面
  useEffect(() => {
    // 如果还没有从localStorage恢复完成，等待
    if (!isHydrated) {
      return
    }

    // 如果有目标页面（用户手动导航），优先使用
    if (targetPageRef.current) {
      setCurrentPage(targetPageRef.current)
      targetPageRef.current = null
      return
    }
    
    // 否则根据用户/游戏状态自动决定
    if (user) {
      if (game && game.status === 'playing') {
        setCurrentPage('gameRoom')
      } else if (game && game.status === 'waiting') {
        setCurrentPage('gameLobby')
      } else {
        setCurrentPage('home')
      }
    } else {
      setCurrentPage('login')
    }
  }, [user, game, isHydrated])

  const navigate = (page) => {
    targetPageRef.current = page
    setCurrentPage(page)
    onNavigate(page)
  }

  return (
    <div className="app">
      <div className="header">
        <h1>蒙面狼人杀</h1>
      </div>
      
      <div className="container">
        {currentPage === 'login' && (
          <Login onNavigate={navigate} />
        )}
        
        {currentPage === 'register' && (
          <Register onNavigate={navigate} />
        )}
        
        {currentPage === 'home' && (
          <Home onNavigate={navigate} />
        )}
        
        {currentPage === 'gameCreate' && (
          <GameCreate onNavigate={navigate} />
        )}
        
        {currentPage === 'gameJoin' && (
          <GameJoin onNavigate={navigate} />
        )}
        
        {currentPage === 'gameLobby' && (
          <GameLobby onNavigate={navigate} />
        )}
        
        {currentPage === 'gameRoom' && (
          <GameRoom onNavigate={navigate} />
        )}
      </div>
      
      <div className="footer">
        <p>© 2026 蒙面狼人杀</p>
      </div>
    </div>
  )
}

function App() {
  const [targetPage, setTargetPage] = useState('login')

  const navigate = (page) => {
    setTargetPage(page)
  }

  return (
    <GameProvider>
      <AppContent onNavigate={navigate} />
    </GameProvider>
  )
}

export default App
