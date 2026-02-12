/**
 * useWebSocket Hook Unit Tests
 */
import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { renderHook, act, waitFor } from '@testing-library/react';
import useWebSocket from './useWebSocket';
import { logger } from '@/utils/logger';

vi.mock('@/utils/logger', () => ({
  logger: {
    info: vi.fn(),
    warn: vi.fn(),
    error: vi.fn(),
    debug: vi.fn(),
  },
}));

class MockWebSocket {
  static instances: MockWebSocket[] = [];
  static CONNECTING = 0;
  static OPEN = 1;
  static CLOSING = 2;
  static CLOSED = 3;

  url: string;
  readyState = MockWebSocket.CONNECTING;
  onopen: ((event: Event) => void) | null = null;
  onmessage: ((event: MessageEvent) => void) | null = null;
  onerror: ((event: Event) => void) | null = null;
  onclose: ((event: CloseEvent) => void) | null = null;

  send = vi.fn();
  close = vi.fn((code?: number, reason?: string) => {
    this.readyState = MockWebSocket.CLOSED;
    this.onclose?.(
      new CloseEvent('close', {
        code: code ?? 1000,
        reason: reason ?? '',
      })
    );
  });

  constructor(url: string) {
    this.url = url;
    MockWebSocket.instances.push(this);
  }

  static reset() {
    MockWebSocket.instances = [];
  }

  triggerOpen() {
    this.readyState = MockWebSocket.OPEN;
    this.onopen?.(new Event('open'));
  }

  triggerMessage(payload: string) {
    this.onmessage?.(new MessageEvent('message', { data: payload }));
  }

  triggerError(event: Event = new Event('error')) {
    this.readyState = MockWebSocket.CLOSED;
    this.onerror?.(event);
  }

  triggerClose(code = 1000, reason = '') {
    this.readyState = MockWebSocket.CLOSED;
    this.onclose?.(new CloseEvent('close', { code, reason }));
  }
}

vi.stubGlobal('WebSocket', MockWebSocket as unknown as typeof WebSocket);

const getWs = (index = 0): MockWebSocket => {
  const instance = MockWebSocket.instances[index];
  expect(instance).toBeTruthy();
  return instance;
};

describe('useWebSocket', () => {
  const mockUrl = 'ws://localhost:8080/ws';
  const mockOnMessage = vi.fn();
  const mockOnOpen = vi.fn();
  const mockOnClose = vi.fn();
  const mockOnError = vi.fn();

  beforeEach(() => {
    vi.clearAllMocks();
    vi.useRealTimers();
    localStorage.clear();
    MockWebSocket.reset();
  });

  afterEach(() => {
    vi.useRealTimers();
  });

  it('should initialize with disconnected state when autoConnect is false', () => {
    const { result } = renderHook(() =>
      useWebSocket({
        url: mockUrl,
        autoConnect: false,
      })
    );

    expect(result.current.connectionState).toBe('disconnected');
    expect(result.current.connected).toBe(false);
    expect(result.current.lastMessage).toBe(null);
    expect(MockWebSocket.instances).toHaveLength(0);
  });

  it('should connect automatically when autoConnect is true', async () => {
    const { result } = renderHook(() =>
      useWebSocket({
        url: mockUrl,
        autoConnect: true,
      })
    );

    await waitFor(() => {
      expect(MockWebSocket.instances).toHaveLength(1);
      expect(result.current.connectionState).toBe('connecting');
    });
  });

  it('should append token to url when token exists', async () => {
    localStorage.setItem('token', 'abc-123');

    renderHook(() =>
      useWebSocket({
        url: mockUrl,
        autoConnect: true,
      })
    );

    await waitFor(() => {
      expect(MockWebSocket.instances).toHaveLength(1);
    });

    expect(getWs().url).toContain('token=abc-123');
    expect(logger.warn).not.toHaveBeenCalledWith('No auth token found for WebSocket connection');
  });

  it('should warn when token does not exist', async () => {
    renderHook(() =>
      useWebSocket({
        url: mockUrl,
        autoConnect: true,
      })
    );

    await waitFor(() => {
      expect(MockWebSocket.instances).toHaveLength(1);
    });

    expect(logger.warn).toHaveBeenCalledWith('No auth token found for WebSocket connection');
  });

  it('should update to connected and call onOpen when socket opens', async () => {
    const { result } = renderHook(() =>
      useWebSocket({
        url: mockUrl,
        autoConnect: true,
        onOpen: mockOnOpen,
      })
    );

    await waitFor(() => {
      expect(MockWebSocket.instances).toHaveLength(1);
    });

    act(() => {
      getWs().triggerOpen();
    });

    await waitFor(() => {
      expect(result.current.connectionState).toBe('connected');
      expect(result.current.connected).toBe(true);
    });

    expect(mockOnOpen).toHaveBeenCalledTimes(1);
  });

  it('should connect when connect is called manually', async () => {
    const { result } = renderHook(() =>
      useWebSocket({
        url: mockUrl,
        autoConnect: false,
      })
    );

    act(() => {
      result.current.connect();
    });

    await waitFor(() => {
      expect(MockWebSocket.instances).toHaveLength(1);
      expect(result.current.connectionState).toBe('connecting');
    });
  });

  it('should handle incoming message and set lastMessage', async () => {
    const { result } = renderHook(() =>
      useWebSocket({
        url: mockUrl,
        autoConnect: true,
        onMessage: mockOnMessage,
      })
    );

    await waitFor(() => {
      expect(MockWebSocket.instances).toHaveLength(1);
    });

    act(() => {
      getWs().triggerOpen();
    });

    const payload = {
      type: 'system_status',
      timestamp: new Date().toISOString(),
      data: { cpuUsage: 45 },
    };

    act(() => {
      getWs().triggerMessage(JSON.stringify(payload));
    });

    expect(result.current.lastMessage).toEqual(payload);
    expect(mockOnMessage).toHaveBeenCalledWith(payload);
  });

  it('should ignore pong messages for onMessage callback', async () => {
    renderHook(() =>
      useWebSocket({
        url: mockUrl,
        autoConnect: true,
        onMessage: mockOnMessage,
      })
    );

    await waitFor(() => {
      expect(MockWebSocket.instances).toHaveLength(1);
    });

    act(() => {
      getWs().triggerOpen();
      getWs().triggerMessage(JSON.stringify({ type: 'pong', timestamp: new Date().toISOString() }));
    });

    expect(mockOnMessage).not.toHaveBeenCalled();
  });

  it('should log parse error for invalid message', async () => {
    renderHook(() =>
      useWebSocket({
        url: mockUrl,
        autoConnect: true,
        onMessage: mockOnMessage,
      })
    );

    await waitFor(() => {
      expect(MockWebSocket.instances).toHaveLength(1);
    });

    act(() => {
      getWs().triggerOpen();
      getWs().triggerMessage('invalid-json');
    });

    expect(logger.error).toHaveBeenCalledWith(
      'Failed to parse WebSocket message:',
      expect.any(Error)
    );
  });

  it('should send string and json payload when connected', async () => {
    const { result } = renderHook(() =>
      useWebSocket({
        url: mockUrl,
        autoConnect: true,
      })
    );

    await waitFor(() => {
      expect(MockWebSocket.instances).toHaveLength(1);
    });

    act(() => {
      getWs().triggerOpen();
    });

    act(() => {
      result.current.send('hello');
      result.current.send({ type: 'ping' });
    });

    expect(getWs().send).toHaveBeenCalledWith('hello');
    expect(getWs().send).toHaveBeenCalledWith(JSON.stringify({ type: 'ping' }));
  });

  it('should warn when sending while disconnected', () => {
    const { result } = renderHook(() =>
      useWebSocket({
        url: mockUrl,
        autoConnect: false,
      })
    );

    act(() => {
      result.current.send('hello');
    });

    expect(logger.warn).toHaveBeenCalledWith('WebSocket is not connected');
  });

  it('should set error state and call onError when socket errors', async () => {
    const { result } = renderHook(() =>
      useWebSocket({
        url: mockUrl,
        autoConnect: true,
        onError: mockOnError,
      })
    );

    await waitFor(() => {
      expect(MockWebSocket.instances).toHaveLength(1);
    });

    act(() => {
      getWs().triggerOpen();
      getWs().triggerError(new Event('error'));
    });

    expect(result.current.connectionState).toBe('error');
    expect(mockOnError).toHaveBeenCalled();
  });

  it('should disconnect manually and call close with manual reason', async () => {
    const { result } = renderHook(() =>
      useWebSocket({
        url: mockUrl,
        autoConnect: true,
        onClose: mockOnClose,
      })
    );

    await waitFor(() => {
      expect(MockWebSocket.instances).toHaveLength(1);
    });

    act(() => {
      getWs().triggerOpen();
      result.current.disconnect();
    });

    expect(getWs().close).toHaveBeenCalledWith(1000, 'Manual disconnect');
    expect(result.current.connectionState).toBe('disconnected');
    expect(mockOnClose).toHaveBeenCalled();
  });

  it('should reconnect automatically after close', async () => {
    vi.useFakeTimers();

    const { result } = renderHook(() =>
      useWebSocket({
        url: mockUrl,
        autoConnect: true,
        reconnectInterval: 100,
        maxRetries: 2,
        onClose: mockOnClose,
      })
    );

    await act(async () => {
      await Promise.resolve();
    });
    expect(MockWebSocket.instances).toHaveLength(1);

    act(() => {
      getWs(0).triggerOpen();
      getWs(0).triggerClose(1006, 'abnormal');
    });

    expect(result.current.connectionState).toBe('disconnected');
    expect(mockOnClose).toHaveBeenCalled();

    act(() => {
      vi.advanceTimersByTime(100);
    });

    expect(MockWebSocket.instances).toHaveLength(2);
    expect(result.current.connectionState).toBe('connecting');
  });

  it('should stop reconnecting after reaching max retries', async () => {
    vi.useFakeTimers();

    renderHook(() =>
      useWebSocket({
        url: mockUrl,
        autoConnect: true,
        reconnectInterval: 100,
        maxRetries: 1,
      })
    );

    await act(async () => {
      await Promise.resolve();
    });
    expect(MockWebSocket.instances).toHaveLength(1);

    act(() => {
      getWs(0).triggerOpen();
      getWs(0).triggerClose(1006, 'abnormal');
      vi.advanceTimersByTime(100);
    });

    expect(MockWebSocket.instances).toHaveLength(2);

    act(() => {
      getWs(1).triggerClose(1006, 'abnormal');
      vi.advanceTimersByTime(500);
    });

    expect(MockWebSocket.instances).toHaveLength(2);
  });

  it('should cleanup websocket on unmount', async () => {
    const { unmount } = renderHook(() =>
      useWebSocket({
        url: mockUrl,
        autoConnect: true,
      })
    );

    await waitFor(() => {
      expect(MockWebSocket.instances).toHaveLength(1);
    });

    unmount();

    expect(getWs().close).toHaveBeenCalledWith(1000, 'Component unmount');
  });
});
