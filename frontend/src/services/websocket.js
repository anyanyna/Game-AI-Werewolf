import io from 'socket.io-client'

class WebSocketService {
  constructor() {
    this.socket = null
    this.listeners = {}
  }

  connect(gameID, playerID, userID) {
    // 连接到WebSocket服务器
    this.socket = io('http://localhost:8080', {
      query: {
        game_id: gameID,
        player_id: playerID,
        user_id: userID
      }
    })

    // 监听连接事件
    this.socket.on('connect', () => {
      console.log('WebSocket connected')
    })

    // 监听断开连接事件
    this.socket.on('disconnect', () => {
      console.log('WebSocket disconnected')
    })

    // 监听错误事件
    this.socket.on('error', (error) => {
      console.error('WebSocket error:', error)
    })

    return this.socket
  }

  disconnect() {
    if (this.socket) {
      this.socket.disconnect()
      this.socket = null
    }
  }

  on(event, callback) {
    if (this.socket) {
      this.socket.on(event, callback)
      // 保存监听器，以便后续可以移除
      if (!this.listeners[event]) {
        this.listeners[event] = []
      }
      this.listeners[event].push(callback)
    }
  }

  off(event, callback) {
    if (this.socket) {
      this.socket.off(event, callback)
      // 从监听器列表中移除
      if (this.listeners[event]) {
        this.listeners[event] = this.listeners[event].filter(cb => cb !== callback)
      }
    }
  }

  emit(event, data) {
    if (this.socket) {
      this.socket.emit(event, data)
    }
  }

  isConnected() {
    return this.socket && this.socket.connected
  }
}

// 导出单例
const webSocketService = new WebSocketService()
export default webSocketService
