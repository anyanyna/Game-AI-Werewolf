// WebSocket 服务 - 使用原生 WebSocket API
// 与后端 gorilla/websocket 兼容

class WebSocketService {
  constructor() {
    this.socket = null
    this.listeners = {}
    this.reconnectAttempts = 0
    this.maxReconnectAttempts = 5
    this.reconnectDelay = 1000
  }

  // 连接到 WebSocket 服务器
  connect(gameID, playerID, userID) {
    // 构建 WebSocket URL 和查询参数
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
    const host = 'localhost:8080'
    const queryParams = `?game_id=${gameID}&player_id=${playerID}&user_id=${userID || ''}`
    const wsUrl = `${protocol}//${host}/ws${queryParams}`

    console.log('Connecting to WebSocket:', wsUrl)

    try {
      this.socket = new WebSocket(wsUrl)

      // 连接打开
      this.socket.onopen = (event) => {
        console.log('WebSocket connected')
        this.reconnectAttempts = 0
        this.emit('connect', event)
      }

      // 接收消息
      this.socket.onmessage = (event) => {
        // 后端可能发送多条消息（用\n分隔）
        const messages = event.data.split('\n')
        messages.forEach(msg => {
          if (msg.trim()) {
            try {
              const data = JSON.parse(msg)
              this.emit(data.type, data)
              // 同时触发 all 事件
              this.emit('all', data)
            } catch (e) {
              console.error('Failed to parse message:', e, msg)
            }
          }
        })
      }

      // 连接关闭
      this.socket.onclose = (event) => {
        console.log('WebSocket disconnected:', event.code, event.reason)
        this.emit('disconnect', event)
        
        // 自动重连
        if (this.reconnectAttempts < this.maxReconnectAttempts && !event.wasClean) {
          this.reconnectAttempts++
          console.log(`Reconnecting... (${this.reconnectAttempts}/${this.maxReconnectAttempts})`)
          setTimeout(() => {
            this.connect(gameID, playerID, userID)
          }, this.reconnectDelay * this.reconnectAttempts)
        }
      }

      // 错误处理
      this.socket.onerror = (error) => {
        console.error('WebSocket error:', error)
        this.emit('error', error)
      }

      return this.socket
    } catch (error) {
      console.error('Failed to create WebSocket:', error)
      throw error
    }
  }

  // 断开连接
  disconnect() {
    if (this.socket) {
      this.socket.close(1000, 'Client disconnect')
      this.socket = null
      this.listeners = {}
      this.reconnectAttempts = 0
    }
  }

  // 发送消息
  send(type, data = {}) {
    if (this.socket && this.socket.readyState === WebSocket.OPEN) {
      const message = JSON.stringify({
        type,
        ...data
      })
      this.socket.send(message)
      return true
    } else {
      console.warn('WebSocket not connected, cannot send message')
      return false
    }
  }

  // 发送消息（别名，兼容旧API）
  emit(event, data) {
    // 内部方法，不对外暴露
  }

  // 监听事件
  on(event, callback) {
    if (!this.listeners[event]) {
      this.listeners[event] = []
    }
    this.listeners[event].push(callback)
  }

  // 移除事件监听
  off(event, callback) {
    if (this.listeners[event]) {
      if (callback) {
        this.listeners[event] = this.listeners[event].filter(cb => cb !== callback)
      } else {
        delete this.listeners[event]
      }
    }
  }

  // 检查连接状态
  isConnected() {
    return this.socket && this.socket.readyState === WebSocket.OPEN
  }

  // 获取连接状态描述
  getState() {
    if (!this.socket) return 'CLOSED'
    switch (this.socket.readyState) {
      case WebSocket.CONNECTING: return 'CONNECTING'
      case WebSocket.OPEN: return 'OPEN'
      case WebSocket.CLOSING: return 'CLOSING'
      case WebSocket.CLOSED: return 'CLOSED'
      default: return 'UNKNOWN'
    }
  }
}

// 导出单例
const webSocketService = new WebSocketService()
export default webSocketService