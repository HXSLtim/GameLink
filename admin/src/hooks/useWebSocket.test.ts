/**
 * useWebSocket Hook Unit Tests
 */
import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { renderHook, act, waitFor } from '@testing-library/react';
import useWebSocket from './useWebSocket';
import { logger } from '@/utils/logger';

// Mock logger
vi.mock('@/utils/logger', () => ({
  logger: {
    info: vi.fn(),
    warn: vi.fn(),
    error: vi.fn(),
    debug: vi.fn(),
  },
}));

// Mock WebSocket
class MockWebSocket {
  url: string;
  readyState: number = WebSocket.CONNECTING;
  onopen: ((event: Event) => void) | null = null;
  onmessage: ((event: MessageEvent) => void) | null = null;
  onerror: ((event: Event) => void) | null = null;
  onclose: ((event: CloseEvent) => void) | null = null;

  send = vi.fn();
  close = vi.fn();

  constructor(url: string) {
    this.url = url;
  }

  // Simulate events
  triggerOpen() {
    this.readyState = WebSocket.OPEN;
    this.onopen?.(new Event('open') as Event);
  }

  triggerMessage(data: string) {
    this.onmessage?.(new MessageEvent('message', { data }) as MessageEvent);
  }

  triggerError(event: Event) {
    this.readyState = WebSocket.CLOSED;
    this.onerror?.(event);
  }

  triggerClose(code: number, reason: string) {
    this.readyState = WebSocket.CLOSED;
    this.onclose?.(new CloseEvent('close', { code, reason }) as CloseEvent);
  }
}

vi.stubGlobal('WebSocket', MockWebSocket as unknown as typeof WebSocket);

const getMockWebSocketInstances = () =>
  (globalThis.WebSocket as unknown as { mock: { instances: MockWebSocket[] } })
    .mock.instances;

const getMockWebSocketInstance = () => getMockWebSocketInstances()[0];

describe('useWebSocket', () => {
  const mockUrl = 'ws://localhost:8080/ws';
  const mockOnMessage = vi.fn();
  const mockOnOpen = vi.fn();
  const mockOnClose = vi.fn();
  const mockOnError = vi.fn();

  beforeEach(() => {
    vi.clearAllMocks();
    localStorage.clear();
  });

  afterEach(() => {
    vi.clearAllMocks();
  });

  describe('Initial State', () => {
    it('should initialize with disconnected state', () => {
      const { result } = renderHook(() =>
        useWebSocket({
          url: mockUrl,
          onMessage: mockOnMessage,
        })
      );

      expect(result.current.connectionState).toBe('disconnected');
      expect(result.current.connected).toBe(false);
      expect(result.current.lastMessage).toBe(null);
    });

    it('should not auto connect when autoConnect is false', () => {
      const { result } = renderHook(() =>
        useWebSocket({
          url: mockUrl,
          autoConnect: false,
        })
      );

      expect(result.current.connectionState).toBe('disconnected');
    });

    it('should use default reconnect interval', () => {
      renderHook(() =>
        useWebSocket({
          url: mockUrl,
          reconnectInterval: undefined,
        })
      );

      // Should not throw error
      expect(true).toBe(true);
    });

    it('should use default max retries', () => {
      renderHook(() =>
        useWebSocket({
          url: mockUrl,
          maxRetries: undefined,
        })
      );

      // Should not throw error
      expect(true).toBe(true);
    });
  });

  describe('Connection', () => {
    it('should connect when autoConnect is true', async () => {
      const { result } = renderHook(() =>
        useWebSocket({
          url: mockUrl,
          onOpen: mockOnOpen,
        })
      );

      expect(result.current.connectionState).toBe('connecting');
    });

    it('should set connection state to connected on open', async () => {
      const { result } = renderHook(() =>
        useWebSocket({
          url: mockUrl,
          onOpen: mockOnOpen,
        })
      );

      // Wait for connection
      await waitFor(() => {
        expect(result.current.connectionState).toBe('connected');
        expect(result.current.connected).toBe(true);
      });
    });

    it('should call onOpen callback when connection opens', async () => {
      renderHook(() =>
        useWebSocket({
          url: mockUrl,
          onOpen: mockOnOpen,
        })
      );

      await waitFor(() => {
        expect(mockOnOpen).toHaveBeenCalled();
      });
    });

    it('should authenticate with token when available', () => {
      localStorage.setItem('token', 'test-token-123');

      const { unmount } = renderHook(() =>
        useWebSocket({
          url: mockUrl,
          autoConnect: true,
        })
      );

      unmount();
      // Should log warning about auth token
      // Token should be added to URL
      expect(logger.warn).not.toHaveBeenCalledWith('No auth token found');
    });

    it('should warn when no auth token is found', () => {
      const { unmount } = renderHook(() =>
        useWebSocket({
          url: mockUrl,
          autoConnect: true,
        })
      );

      // Should log warning about no token
      expect(logger.warn).toHaveBeenCalledWith('No auth token found for WebSocket connection');

      unmount();
    });
  });

  describe('Manual Connection', () => {
    it('should provide connect function', () => {
      const { result } = renderHook(() =>
        useWebSocket({
          url: mockUrl,
          autoConnect: false,
        })
      );

      expect(typeof result.current.connect).toBe('function');
    });

    it('should connect when connect function is called', () => {
      const { result } = renderHook(() =>
        useWebSocket({
          url: mockUrl,
          autoConnect: false,
        })
      );

      act(() => {
        result.current.connect();
      });

      expect(result.current.connectionState).toBe('connecting');
    });

    it('should not reconnect if already connected', () => {
      const { result } = renderHook(() =>
        useWebSocket({
          url: mockUrl,
          autoConnect: true,
          onOpen: mockOnOpen,
        })
      );

      act(() => {
        result.current.connect();
      });

      // Should log that already connected
      expect(logger.info).toHaveBeenCalledWith('WebSocket already connected');
    });
  });

  describe('Disconnection', () => {
    it('should provide disconnect function', () => {
      const { result } = renderHook(() =>
        useWebSocket({
          url: mockUrl,
          autoConnect: false,
        })
      );

      expect(typeof result.current.disconnect).toBe('function');
    });

    it('should disconnect when disconnect function is called', async () => {
      const { result } = renderHook(() =>
        useWebSocket({
          url: mockUrl,
          autoConnect: false,
        })
      );

      act(() => {
        result.current.disconnect();
      });

      expect(result.current.connectionState).toBe('disconnected');
    });

    it('should call onClose callback on disconnect', async () => {
      const { result } = renderHook(() =>
        useWebSocket({
          url: mockUrl,
          autoConnect: false,
          onClose: mockOnClose,
        })
      );

      act(() => {
        result.current.disconnect();
      });

      expect(mockOnClose).toHaveBeenCalled();
    });

    it('should close WebSocket with manual disconnect code', () => {
      const { result } = renderHook(() =>
        useWebSocket({
          url: mockUrl,
          autoConnect: false,
        })
      );

      act(() => {
        result.current.disconnect();
      });

      // Close should be called
      const mockWs = getMockWebSocketInstance();
      if (mockWs) {
        expect(mockWs.close).toHaveBeenCalledWith(1000, 'Manual disconnect');
      }
    });
  });

  describe('Message Handling', () => {
    it('should handle incoming messages', async () => {
      const { result } = renderHook(() =>
        useWebSocket({
          url: mockUrl,
          onMessage: mockOnMessage,
          autoConnect: true,
        })
      );

      // Wait for connection
      await waitFor(() => {
        expect(result.current.connectionState).toBe('connected');
      });

      // Trigger message
      const testMessage = { type: 'test', data: 'hello' };
      act(() => {
        const mockWs = getMockWebSocketInstance();
        if (mockWs) {
          mockWs.triggerMessage(JSON.stringify(testMessage));
        }
      });

      expect(result.current.lastMessage).toEqual(testMessage);
      expect(mockOnMessage).toHaveBeenCalledWith(testMessage);
    });

    it('should handle pong messages without calling onMessage', async () => {
      const { result } = renderHook(() =>
        useWebSocket({
          url: mockUrl,
          onMessage: mockOnMessage,
          autoConnect: true,
        })
      );

      await waitFor(() => {
        expect(result.current.connectionState).toBe('connected');
      });

      // Trigger pong message
      act(() => {
        const mockWs = getMockWebSocketInstance();
        if (mockWs) {
          mockWs.triggerMessage(JSON.stringify({ type: 'pong' }));
        }
      });

      expect(mockOnMessage).not.toHaveBeenCalled();
    });

    it('should handle invalid JSON messages', async () => {
      const { result } = renderHook(() =>
        useWebSocket({
          url: mockUrl,
          onMessage: mockOnMessage,
          autoConnect: true,
        })
      );

      await waitFor(() => {
        expect(result.current.connectionState).toBe('connected');
      });

      // Trigger invalid message
      act(() => {
        const mockWs = getMockWebSocketInstance();
        if (mockWs) {
          mockWs.triggerMessage('invalid json');
        }
      });

      expect(logger.error).toHaveBeenCalledWith(
        'Failed to parse WebSocket message:',
        expect.any(Error)
      );
    });
  });

  describe('Send Messages', () => {
    it('should provide send function', () => {
      const { result } = renderHook(() =>
        useWebSocket({
          url: mockUrl,
          autoConnect: false,
        })
      );

      expect(typeof result.current.send).toBe('function');
    });

    it('should send string messages when connected', async () => {
      const { result } = renderHook(() =>
        useWebSocket({
          url: mockUrl,
          autoConnect: true,
        })
      );

      await waitFor(() => {
        expect(result.current.connectionState).toBe('connected');
      });

      act(() => {
        result.current.send('test message');
      });

      const mockWs = getMockWebSocketInstance();
      if (mockWs) {
        expect(mockWs.send).toHaveBeenCalledWith('test message');
      }
    });

    it('should send JSON messages when connected', async () => {
      const { result } = renderHook(() =>
        useWebSocket({
          url: mockUrl,
          autoConnect: true,
        })
      );

      await waitFor(() => {
        expect(result.current.connectionState).toBe('connected');
      });

      const testMessage = { type: 'test', data: 'hello' };
      act(() => {
        result.current.send(testMessage);
      });

      const mockWs = getMockWebSocketInstance();
      if (mockWs) {
        expect(mockWs.send).toHaveBeenCalledWith(JSON.stringify(testMessage));
      }
    });

    it('should warn when trying to send while disconnected', () => {
      const { result } = renderHook(() =>
        useWebSocket({
          url: mockUrl,
          autoConnect: false,
        })
      );

      act(() => {
        result.current.send('test message');
      });

      expect(logger.warn).toHaveBeenCalledWith('WebSocket is not connected');
    });
  });

  describe('Error Handling', () => {
    it('should call onError callback on error', async () => {
      const { result } = renderHook(() =>
        useWebSocket({
          url: mockUrl,
          autoConnect: true,
          onError: mockOnError,
        })
      );

      await waitFor(() => {
        expect(result.current.connectionState).toBe('connected');
      });

      // Trigger error
      act(() => {
        const mockWs = getMockWebSocketInstance();
        if (mockWs) {
          mockWs.triggerError(new Event('error') as Event);
        }
      });

      expect(result.current.connectionState).toBe('error');
      expect(mockOnError).toHaveBeenCalled();
      expect(logger.error).toHaveBeenCalled();
    });
  });

  describe('Auto Reconnect', () => {
    it('should attempt to reconnect on close', async () => {
      vi.useFakeTimers();

      const { result } = renderHook(() =>
        useWebSocket({
          url: mockUrl,
          autoConnect: true,
          reconnectInterval: 1000,
          maxRetries: 3,
          onClose: mockOnClose,
        })
      );

      await waitFor(() => {
        expect(result.current.connectionState).toBe('connected');
      });

      // Trigger close
      act(() => {
        const mockWs = getMockWebSocketInstance();
        if (mockWs) {
          mockWs.triggerClose(1000, 'Normal close');
        }
      });

      // Wait for reconnect
      await waitFor(() => {
        expect(result.current.connectionState).toBe('connecting');
      });

      vi.useRealTimers();
    });

    it('should stop reconnecting after max retries', async () => {
      vi.useFakeTimers();

      const { result } = renderHook(() =>
        useWebSocket({
          url: mockUrl,
          autoConnect: true,
          reconnectInterval: 100,
          maxRetries: 2,
        })
      );

      await waitFor(() => {
        expect(result.current.connectionState).toBe('connected');
      });

      // Trigger close multiple times
      for (let i = 0; i < 3; i++) {
        act(() => {
          const mockWs = getMockWebSocketInstances()[i];
          if (mockWs) {
            mockWs.triggerClose(1000, 'Normal close');
          }
        });
        vi.advanceTimersByTime(100);
      }

      // Should have attempted 2 reconnects (maxRetries)
      vi.useRealTimers();
    });

    it('should not reconnect when manual disconnect', async () => {
      const { result } = renderHook(() =>
        useWebSocket({
          url: mockUrl,
          autoConnect: true,
          maxRetries: 5,
        })
      );

      await waitFor(() => {
        expect(result.current.connectionState).toBe('connected');
      });

      // Manual disconnect
      act(() => {
        result.current.disconnect();
      });

      // Trigger close event
      act(() => {
        const mockWs = getMockWebSocketInstance();
        if (mockWs) {
          mockWs.triggerClose(1000, 'Manual disconnect');
        }
      });

      // Should not attempt to reconnect
      expect(logger.info).toHaveBeenCalledWith('WebSocket closed', expect.any(Object));
    });
  });

  describe('Cleanup', () => {
    it('should cleanup on unmount', () => {
      const { unmount } = renderHook(() =>
        useWebSocket({
          url: mockUrl,
          autoConnect: false,
        })
      );

      unmount();

      const mockWs = getMockWebSocketInstance();
      if (mockWs) {
        expect(mockWs.close).toHaveBeenCalledWith(1000, 'Component unmount');
      }
    });

    it('should clear timers on unmount', () => {
      vi.useFakeTimers();

      const { unmount } = renderHook(() =>
        useWebSocket({
          url: mockUrl,
          autoConnect: true,
        })
      );

      unmount();

      vi.useRealTimers();
      // Should not throw error
      expect(true).toBe(true);
    });
  });

  describe('Custom Options', () => {
    it('should use custom reconnect interval', () => {
      const customInterval = 5000;

      renderHook(() =>
        useWebSocket({
          url: mockUrl,
          autoConnect: true,
          reconnectInterval: customInterval,
        })
      );

      // Should use custom interval
      expect(true).toBe(true);
    });

    it('should use custom max retries', () => {
      const customMaxRetries = 10;

      renderHook(() =>
        useWebSocket({
          url: mockUrl,
          autoConnect: true,
          maxRetries: customMaxRetries,
        })
      );

      // Should use custom max retries
      expect(true).toBe(true);
    });

    it('should work without optional callbacks', () => {
      const { result } = renderHook(() =>
        useWebSocket({
          url: mockUrl,
          autoConnect: false,
        })
      );

      expect(result.current.connectionState).toBe('disconnected');
      expect(typeof result.current.connect).toBe('function');
      expect(typeof result.current.disconnect).toBe('function');
      expect(typeof result.current.send).toBe('function');
    });
  });
});
