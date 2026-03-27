import React, { useState } from 'react'
import { useGame } from '../contexts/GameContext'

function Login({ onNavigate }) {
  const [username, setUsername] = useState('')
  const [password, setPassword] = useState('')
  const { login, setError, error, clearError, setIsLoading, isLoading } = useGame()

  const handleSubmit = (e) => {
    e.preventDefault()
    clearError()
    setIsLoading(true)
    
    // 模拟登录请求
    fetch('http://localhost:8080/api/user/login', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ username, password }),
    })
      .then((response) => response.json())
      .then((data) => {
        if (data.error) {
          setError(data.error)
        } else {
          login(data.user)
          onNavigate('home')
        }
      })
      .catch((error) => {
        setError('登录失败，请稍后重试')
        console.error('Error:', error)
      })
      .finally(() => {
        setIsLoading(false)
      })
  }

  return (
    <div className="card">
      <h2>登录</h2>
      <form onSubmit={handleSubmit}>
        <div className="form-group">
          <label htmlFor="username">用户名</label>
          <input
            type="text"
            id="username"
            value={username}
            onChange={(e) => setUsername(e.target.value)}
            required
          />
        </div>
        <div className="form-group">
          <label htmlFor="password">密码</label>
          <input
            type="password"
            id="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            required
          />
        </div>
        {error && <div className="error">{error}</div>}
        <button type="submit" className="btn" disabled={isLoading}>
          {isLoading ? '登录中...' : '登录'}
        </button>
      </form>
      <p style={{ marginTop: '15px' }}>
        还没有账号？ <a href="#" onClick={() => onNavigate('register')}>立即注册</a>
      </p>
    </div>
  )
}

export default Login
