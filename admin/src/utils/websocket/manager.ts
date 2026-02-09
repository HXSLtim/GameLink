/**
 * WebSocket Manager
 * 单例模式管理 WebSocket 连接，支持自动重连、心跳检测
 */
import { logger } from '@/utils/logger';
import type {
  WSMessage,
  WebSocketConfig,
  EventListeners,
  MessageHandler,
} from './types';
import { ConnectionStatus } from './types';
import type { TConnectionStatus, TMessageType } from './types';

/**
 * WebSocket Manager 类
 */
class WebSocketManager {
  private ws: WebSocket | null = null;
  private config: Required<WebSocketConfig>;
  private status: TConnectionStatus = ConnectionStatus.Disconnected;
  private reconnectTimer: ReturnType<typeof setTimeout> | null = null;
  private heartbeatTimer: ReturnType<typeof setInterval> | null = null;
  private reconnectAttempts = 0;
  private messageHandlers: Map<TMessageType, Set<MessageHandler>> = new Map();
  private eventListeners: EventListeners = {};
  private isManualClose = false;

  // 默认配置
  private defaultConfig: Required<WebSocketConfig> = {
    url: '',
    token: null,
    heartbeatInterval: Number(import.meta.env.VITE_WEBSOCKET_HEARTBEAT_INTERVAL) || 54000, // 54 秒心跳 - 对齐后端配置
    reconnect: true,
    reconnectInterval: Number(import.meta.env.VITE_WEBSOCKET_RECONNECT_INTERVAL) || 1000,
    maxReconnectAttempts: Number(import.meta.env.VITE_WEBSOCKET_RECONNECT_ATTEMPTS) || 10,
    debug: import.meta.env.DEV,
  };

  constructor(config?: Partial<WebSocketConfig>) {
    this.config = { ...this.defaultConfig, ...config };
  }

  /**
   * 更新配置
   */
  updateConfig(config: Partial<WebSocketConfig>): void {
    this.config = { ...this.config, ...config };
  }

  /**
   * 获取当前连接状态
   */
  getStatus(): TConnectionStatus {
    return this.status;
  }

  /**
   * 是否已连接
   */
  isConnected(): boolean {
    return this.status === ConnectionStatus.Connected && this.ws?.readyState === WebSocket.OPEN;
  }

  /**
   * 连接 WebSocket
   */
  connect(url?: string, token?: string): void {
    if (url) this.config.url = url;
    if (token !== undefined) this.config.token = token;

    if (!this.config.url) {
      logger.error('WebSocket URL is required');
      return;
    }

    // 如果已连接，先断开
    if (this.ws && this.ws.readyState === WebSocket.OPEN) {
      this.disconnect();
    }

    this.isManualClose = false;
    this.setStatus(ConnectionStatus.Connecting);

    try {
      // 构建 WebSocket URL（添加 token）
      const wsUrl = this.buildWebSocketUrl();
      this.debugLog('Connecting to:', wsUrl);

      this.ws = new WebSocket(wsUrl);

      // 设置事件处理器
      this.ws.onopen = this.handleOpen.bind(this);
      this.ws.onmessage = this.handleMessage.bind(this);
      this.ws.onerror = this.handleError.bind(this);
      this.ws.onclose = this.handleClose.bind(this);
    } catch (error) {
      logger.error('Failed to create WebSocket:', error);
      this.setStatus(ConnectionStatus.Error);
      this.scheduleReconnect();
    }
  }

  /**
   * 断开连接
   */
  disconnect(): void {
    this.isManualClose = true;
    this.clearTimers();

    if (this.ws) {
      this.debugLog('Manual disconnect');
      this.ws.close(1000, 'Client disconnect');
      this.ws = null;
    }

    this.setStatus(ConnectionStatus.Disconnected);
  }

  /**
   * 发送消息
   */
  send(type: string, data?: unknown): boolean {
    if (!this.isConnected()) {
      logger.warn('WebSocket is not connected, cannot send message');
      return false;
    }

    try {
      const message: WSMessage = {
        type,
        timestamp: new Date().toISOString(),
        data,
      };

      this.ws!.send(JSON.stringify(message));
      this.debugLog('Sent:', message);
      return true;
    } catch (error) {
      logger.error('Failed to send WebSocket message:', error);
      return false;
    }
  }

  /**
   * 注册消息处理器
   */
  on(messageType: TMessageType, handler: MessageHandler): () => void {
    if (!this.messageHandlers.has(messageType)) {
      this.messageHandlers.set(messageType, new Set());
    }
    this.messageHandlers.get(messageType)!.add(handler);

    // 返回取消订阅函数
    return () => {
      this.off(messageType, handler);
    };
  }

  /**
   * 取消消息处理器
   */
  off(messageType: TMessageType, handler: MessageHandler): void {
    const handlers = this.messageHandlers.get(messageType);
    if (handlers) {
      handlers.delete(handler);
      if (handlers.size === 0) {
        this.messageHandlers.delete(messageType);
      }
    }
  }

  /**
   * 设置事件监听器
   */
  setEventListeners(listeners: EventListeners): void {
    this.eventListeners = { ...this.eventListeners, ...listeners };
  }

  /**
   * 清除事件监听器
   */
  clearEventListeners(): void {
    this.eventListeners = {};
  }

  // ==================== 私有方法 ====================

  /**
   * 构建 WebSocket URL
   */
  private buildWebSocketUrl(): string {
    const { url, token } = this.config;

    // 如果 URL 已包含协议，直接使用
    if (url.startsWith('ws://') || url.startsWith('wss://')) {
      // 添加 token 作为查询参数
      const separator = url.includes('?') ? '&' : '?';
      return token ? `${url}${separator}token=${token}` : url;
    }

    // 否则根据当前协议构建
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const wsUrl = `${protocol}//${url}`;

    // 添加 token
    const separator = wsUrl.includes('?') ? '&' : '?';
    return token ? `${wsUrl}${separator}token=${token}` : wsUrl;
  }

  /**
   * 设置状态
   */
  private setStatus(status: TConnectionStatus): void {
    if (this.status !== status) {
      this.debugLog('Status changed:', this.status, '->', status);
      this.status = status;
    }
  }

  /**
   * 清除定时器
   */
  private clearTimers(): void {
    if (this.reconnectTimer) {
      clearTimeout(this.reconnectTimer);
      this.reconnectTimer = null;
    }
    if (this.heartbeatTimer) {
      clearInterval(this.heartbeatTimer);
      this.heartbeatTimer = null;
    }
  }

  /**
   * 调度重连
   */
  private scheduleReconnect(): void {
    if (this.isManualClose) {
      return;
    }

    if (!this.config.reconnect) {
      this.debugLog('Reconnect disabled, not scheduling reconnect');
      return;
    }

    if (this.reconnectAttempts >= this.config.maxReconnectAttempts) {
      logger.error('Max reconnect attempts reached, giving up');
      this.setStatus(ConnectionStatus.Error);
      return;
    }

    // 指数退避重连
    const delay = Math.min(
      this.config.reconnectInterval * Math.pow(2, this.reconnectAttempts),
      30000 // 最大 30 秒
    );

    this.debugLog(`Scheduling reconnect in ${delay}ms (attempt ${this.reconnectAttempts + 1})`);
    this.setStatus(ConnectionStatus.Reconnecting);

    this.eventListeners.onReconnecting?.(this.reconnectAttempts + 1);

    this.reconnectTimer = setTimeout(() => {
      this.reconnectAttempts++;
      this.connect();
    }, delay);
  }

  /**
   * 启动心跳
   */
  private startHeartbeat(): void {
    this.clearTimers();

    this.heartbeatTimer = setInterval(() => {
      if (this.isConnected()) {
        this.send('ping');
      }
    }, this.config.heartbeatInterval);
  }

  /**
   * 处理连接打开
   */
  private handleOpen(event: Event): void {
    this.debugLog('WebSocket connected');
    this.setStatus(ConnectionStatus.Connected);
    this.reconnectAttempts = 0;
    this.startHeartbeat();

    // 如果重连成功，触发回调
    if (this.eventListeners.onReconnected && this.reconnectAttempts > 0) {
      this.eventListeners.onReconnected();
    }

    this.eventListeners.onOpen?.(event);
  }

  /**
   * 处理消息接收
   */
  private handleMessage(event: MessageEvent): void {
    try {
      const message: WSMessage = JSON.parse(event.data);
      this.debugLog('Received:', message);

      // 处理 pong 消息
      if (message.type === 'pong') {
        return;
      }

      // 触发特定类型的处理器
      const handlers = this.messageHandlers.get(message.type as TMessageType);
      if (handlers) {
        handlers.forEach((handler) => {
          try {
            handler(message);
          } catch (error) {
            logger.error(`Error in message handler for ${message.type}:`, error);
          }
        });
      }

      // 触发全局消息处理器
      this.eventListeners.onMessage?.(message);
    } catch (error) {
      logger.error('Failed to parse WebSocket message:', error);
    }
  }

  /**
   * 处理错误
   */
  private handleError(event: Event): void {
    this.debugLog('WebSocket error:', event);
    this.setStatus(ConnectionStatus.Error);
    this.eventListeners.onError?.(event);
  }

  /**
   * 处理连接关闭
   */
  private handleClose(event: CloseEvent): void {
    this.debugLog('WebSocket closed:', event.code, event.reason);
    this.clearTimers();

    if (!this.isManualClose) {
      this.setStatus(ConnectionStatus.Disconnected);
      this.scheduleReconnect();
    }

    this.eventListeners.onClose?.(event);
  }

  /**
   * 调试日志
   */
  private debugLog(...args: unknown[]): void {
    if (this.config.debug) {
      logger.info('[WebSocket]', ...args);
    }
  }
}

// ==================== 单例导出 ====================

/**
 * WebSocket Manager 单例
 */
const wsManagerInstance = new WebSocketManager();

export { WebSocketManager };
export const wsManager = wsManagerInstance;

/**
 * 初始化 WebSocket 连接
 * @param url WebSocket 服务器地址
 * @param token 认证 token
 */
export function initWebSocket(url: string, token: string): void {
  wsManagerInstance.updateConfig({ url, token });
  wsManagerInstance.connect();
}

/**
 * 断开 WebSocket 连接
 */
export function disconnectWebSocket(): void {
  wsManagerInstance.disconnect();
}

/**
 * 获取 WebSocket Manager 实例
 */
export function getWebSocketManager(): WebSocketManager {
  return wsManagerInstance;
}

export default WebSocketManager;
