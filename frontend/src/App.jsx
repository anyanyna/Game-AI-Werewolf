import React, { useState } from 'react'
import { GameProvider } from './contexts/GameContext'
import Login from './pages/Login.jsx'
import Register from './pages/Register.jsx'
import Home from './pages/Home.jsx'
import GameCreate from './pages/GameCreate.jsx'
import GameJoin from './pages/GameJoin.jsx'
import GameLobby from './pages/GameLobby.jsx'
import GameRoom from './pages/GameRoom.jsx'

function App() {
  const [currentPage, setCurrentPage] = useState('login')

  const navigate = (page) => {
    setCurrentPage(page)
  }

  return (
    <GameProvider>
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
    </GameProvider>
  )
}

export default App
