/**
 * WebSocket Hook
 * 提供 WebSocket 连接管理、自动重连、心跳检测等功能
 */
import { useState, useEffect, useRef, useCallback } from 'react';
import type {
  WSMessage,
  WSConnectionState,
  UseWebSocketOptions,
  UseWebSocketReturn,
} from '@/types/monitor';

const DEFAULT_RECONNECT_INTERVAL = 3000; // 3 seconds
const DEFAULT_MAX_RETRIES = 5;
const PING_INTERVAL = 30000; // 30 seconds

/**
 * WebSocket Hook
 * @param options - WebSocket 配置选项
 * @returns WebSocket 连接状态和操作方法
 */
export function useWebSocket(options: UseWebSocketOptions): UseWebSocketReturn {
  const {
    url,
    onMessage,
    onError,
    onOpen,
    onClose,
    reconnectInterval = DEFAULT_RECONNECT_INTERVAL,
    maxRetries = DEFAULT_MAX_RETRIES,
    autoConnect = true,
  } = options;

  const [connectionState, setConnectionState] = useState<WSConnectionState>('disconnected');
  const [lastMessage, setLastMessage] = useState<WSMessage | null>(null);

  const wsRef = useRef<WebSocket | null>(null);
  const retryCountRef = useRef(0);
  const reconnectTimerRef = useRef<ReturnType<typeof setTimeout> | null>(null);
  const pingTimerRef = useRef<ReturnType<typeof setInterval> | null>(null);
  const shouldReconnectRef = useRef(true);

  /**
   * 获取带有认证的 WebSocket URL
   */
  const getAuthenticatedUrl = useCallback((): string => {
    const token = localStorage.getItem('token');
    if (!token) {
      console.warn('No auth token found for WebSocket connection');
      return url;
    }

    // 将 token 添加到 URL 查询参数
    const separator = url.includes('?') ? '&' : '?';
    return `${url}${separator}token=${encodeURIComponent(token)}`;
  }, [url]);

  /**
   * 清理定时器
   */
  const clearTimers = useCallback(() => {
    if (reconnectTimerRef.current) {
      clearTimeout(reconnectTimerRef.current);
      reconnectTimerRef.current = null;
    }
    if (pingTimerRef.current) {
      clearInterval(pingTimerRef.current);
      pingTimerRef.current = null;
    }
  }, []);

  /**
   * 启动心跳检测
   */
  const startPingInterval = useCallback(() => {
    if (pingTimerRef.current) {
      clearInterval(pingTimerRef.current);
    }

    pingTimerRef.current = setInterval(() => {
      if (wsRef.current?.readyState === WebSocket.OPEN) {
        const pingMessage: WSMessage = {
          type: 'ping',
          timestamp: new Date().toISOString(),
        };
        wsRef.current.send(JSON.stringify(pingMessage));
      }
    }, PING_INTERVAL);
  }, []);

  /**
   * 重连函数 - 提取到 connect 函数外部
   */
  const reconnect = useCallback(() => {
    if (wsRef.current?.readyState === WebSocket.OPEN) {
      console.log('WebSocket already connected');
      return;
    }

    clearTimers();
    shouldReconnectRef.current = true;
    setConnectionState('connecting');

    try {
      const authenticatedUrl = getAuthenticatedUrl();
      wsRef.current = new WebSocket(authenticatedUrl);

      // 创建事件处理器
      const handleOpen = () => {
        console.log('WebSocket connected');
        setConnectionState('connected');
        retryCountRef.current = 0;
        startPingInterval();
        onOpen?.();
      };

      const handleMessage = (event: MessageEvent) => {
        try {
          const message: WSMessage = JSON.parse(event.data);
          setLastMessage(message);

          // 处理 pong 消息
          if (message.type === 'pong') {
            return;
          }

          onMessage?.(message);
        } catch (error) {
          console.error('Failed to parse WebSocket message:', error);
        }
      };

      const handleError = (event: Event) => {
        console.error('WebSocket error:', event);
        setConnectionState('error');
        onError?.(event);
      };

      const handleClose = (event: CloseEvent) => {
        console.log('WebSocket closed:', event.code, event.reason);
        setConnectionState('disconnected');
        clearTimers();
        onClose?.();

        // 自动重连
        if (shouldReconnectRef.current && retryCountRef.current < maxRetries) {
          retryCountRef.current++;
          console.log(`Reconnecting... (attempt ${retryCountRef.current}/${maxRetries})`);

          reconnectTimerRef.current = setTimeout(() => {
            if (autoConnect) {
              reconnect();
            }
          }, reconnectInterval * retryCountRef.current);
        }
      };

      // 绑定事件处理器
      if (wsRef.current) {
        wsRef.current.onopen = handleOpen;
        wsRef.current.onmessage = handleMessage;
        wsRef.current.onerror = handleError;
        wsRef.current.onclose = handleClose;
      }
    } catch (error) {
      console.error('Failed to create WebSocket:', error);
      setConnectionState('error');
    }
  }, [
    getAuthenticatedUrl,
    clearTimers,
    startPingInterval,
    onOpen,
    onMessage,
    onError,
    onClose,
    maxRetries,
    reconnectInterval,
    autoConnect,
  ]);

  /**
   * 连接 WebSocket - 简化为调用 reconnect
   */
  const connect = useCallback(() => {
    reconnect();
  }, [reconnect]);

  /**
   * 断开 WebSocket 连接
   */
  const disconnect = useCallback(() => {
    shouldReconnectRef.current = false;
    clearTimers();

    if (wsRef.current) {
      wsRef.current.close(1000, 'Manual disconnect');
      wsRef.current = null;
    }

    setConnectionState('disconnected');
  }, [clearTimers]);

  /**
   * 发送消息
   */
  const send = useCallback((data: unknown) => {
    if (wsRef.current?.readyState === WebSocket.OPEN) {
      const message = typeof data === 'string' ? data : JSON.stringify(data);
      wsRef.current.send(message);
    } else {
      console.warn('WebSocket is not connected');
    }
  }, []);

  // 自动连接
  useEffect(() => {
    if (autoConnect) {
      connect();
    }

    return () => {
      shouldReconnectRef.current = false;
      clearTimers();
      if (wsRef.current) {
        wsRef.current.close(1000, 'Component unmount');
      }
    };
  }, [autoConnect, connect, clearTimers]);

  return {
    connected: connectionState === 'connected',
    connectionState,
    send,
    connect,
    disconnect,
    lastMessage,
  };
}

export default useWebSocket;
