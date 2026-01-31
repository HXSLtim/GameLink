/**
 * WebSocket 连接管理
 */

import { ref, onUnmounted } from 'vue'
import { useUserStore } from '@/store'

export interface WsMessage {
  type: 'chat_message' | 'order_status' | 'order_new' | 'notification' | 'pong'
  timestamp: number
  data: unknown
}

export type WsEventHandler = (message: WsMessage) => void

export function useWebSocket() {
  const userStore = useUserStore()
  
  const isConnected = ref(false)
  const isConnecting = ref(false)
  
  let ws: UniApp.SocketTask | null = null
  let heartbeatTimer: ReturnType<typeof setInterval> | null = null
  let reconnectTimer: ReturnType<typeof setTimeout> | null = null
  let reconnectAttempts = 0
  const maxReconnectAttempts = 5
  
  const handlers = new Map<string, Set<WsEventHandler>>()
  
  /**
   * 连接 WebSocket
   */
  function connect() {
    if (isConnected.value || isConnecting.value) return
    if (!userStore.token) return
    
    isConnecting.value = true
    
    // #ifdef H5
    const wsUrl = `${location.protocol === 'https:' ? 'wss:' : 'ws:'}//${location.host}/api/v1/ws?token=${userStore.token}`
    // #endif
    
    // #ifndef H5
    const baseUrl = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8081'
    const wsProtocol = baseUrl.startsWith('https') ? 'wss:' : 'ws:'
    const wsHost = baseUrl.replace(/^https?:\/\//, '').replace(/\/api\/v1$/, '')
    const wsUrl = `${wsProtocol}//${wsHost}/api/v1/ws?token=${userStore.token}`
    // #endif
    
    ws = uni.connectSocket({
      url: wsUrl,
      success: () => console.log('WebSocket connecting...'),
    })
    
    ws.onOpen(() => {
      console.log('WebSocket connected')
      isConnected.value = true
      isConnecting.value = false
      reconnectAttempts = 0
      startHeartbeat()
    })
    
    ws.onMessage((res) => {
      try {
        const message = JSON.parse(res.data as string) as WsMessage
        emit(message.type, message)
        emit('*', message) // 通配符监听所有消息
      } catch (e) {
        console.error('Parse WS message error:', e)
      }
    })
    
    ws.onClose(() => {
      console.log('WebSocket closed')
      isConnected.value = false
      isConnecting.value = false
      stopHeartbeat()
      scheduleReconnect()
    })
    
    ws.onError((err) => {
      console.error('WebSocket error:', err)
      isConnecting.value = false
    })
  }
  
  /**
   * 断开连接
   */
  function disconnect() {
    stopHeartbeat()
    clearReconnect()
    
    if (ws) {
      ws.close({})
      ws = null
    }
    
    isConnected.value = false
    isConnecting.value = false
  }
  
  /**
   * 发送消息
   */
  function send(data: Record<string, unknown>) {
    if (!ws || !isConnected.value) {
      console.warn('WebSocket not connected')
      return false
    }
    
    ws.send({
      data: JSON.stringify(data),
      success: () => {},
      fail: (err) => console.error('Send message failed:', err),
    })
    
    return true
  }
  
  /**
   * 监听消息
   */
  function on(type: string, handler: WsEventHandler) {
    if (!handlers.has(type)) {
      handlers.set(type, new Set())
    }
    handlers.get(type)!.add(handler)
  }
  
  /**
   * 取消监听
   */
  function off(type: string, handler: WsEventHandler) {
    handlers.get(type)?.delete(handler)
  }
  
  /**
   * 触发事件
   */
  function emit(type: string, message: WsMessage) {
    handlers.get(type)?.forEach(handler => handler(message))
  }
  
  /**
   * 心跳
   */
  function startHeartbeat() {
    heartbeatTimer = setInterval(() => {
      send({ type: 'ping' })
    }, 30000)
  }
  
  function stopHeartbeat() {
    if (heartbeatTimer) {
      clearInterval(heartbeatTimer)
      heartbeatTimer = null
    }
  }
  
  /**
   * 重连
   */
  function scheduleReconnect() {
    if (reconnectTimer) return
    if (reconnectAttempts >= maxReconnectAttempts) {
      console.log('Max reconnect attempts reached')
      return
    }
    
    const delay = Math.min(1000 * Math.pow(2, reconnectAttempts), 30000)
    reconnectAttempts++
    
    reconnectTimer = setTimeout(() => {
      reconnectTimer = null
      connect()
    }, delay)
  }
  
  function clearReconnect() {
    if (reconnectTimer) {
      clearTimeout(reconnectTimer)
      reconnectTimer = null
    }
    reconnectAttempts = 0
  }
  
  // 组件卸载时清理
  onUnmounted(() => {
    disconnect()
    handlers.clear()
  })
  
  return {
    isConnected,
    isConnecting,
    connect,
    disconnect,
    send,
    on,
    off,
  }
}
