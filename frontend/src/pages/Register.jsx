import React, { useState } from 'react'
import { useGame } from '../contexts/GameContext'

function Register({ onNavigate }) {
  const [username, setUsername] = useState('')
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const { register, setError, error, clearError, setIsLoading, isLoading } = useGame()

  const handleSubmit = (e) => {
    e.preventDefault()
    clearError()
    setIsLoading(true)
    
    // 模拟注册请求
    fetch('http://localhost:8080/api/user/register', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ username, email, password }),
    })
      .then((response) => response.json())
      .then((data) => {
        if (data.error) {
          setError(data.error)
        } else {
          register(data.user)
          onNavigate('home')
        }
      })
      .catch((error) => {
        setError('注册失败，请稍后重试')
        console.error('Error:', error)
      })
      .finally(() => {
        setIsLoading(false)
      })
  }

  return (
    <div className="card">
      <h2>注册</h2>
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
          <label htmlFor="email">邮箱</label>
          <input
            type="email"
            id="email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
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
            minLength={6}
          />
        </div>
        {error && <div className="error">{error}</div>}
        <button type="submit" className="btn" disabled={isLoading}>
          {isLoading ? '注册中...' : '注册'}
        </button>
      </form>
      <p style={{ marginTop: '15px' }}>
        已有账号？ <a href="#" onClick={() => onNavigate('login')}>立即登录</a>
      </p>
    </div>
  )
}

export default Register
