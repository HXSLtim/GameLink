/**
 * WebSocket Manager Tests
 * Tests for the singleton wsManager
 */
import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { wsManager } from './manager';

// Mock WebSocket
class MockWebSocket {
  static READY_STATE_CONNECTING = 0;
  static READY_STATE_OPEN = 1;
  static READY_STATE_CLOSING = 2;
  static READY_STATE_CLOSED = 3;

  readyState = MockWebSocket.READY_STATE_OPEN;
  url = '';
  onopen: ((event: Event) => void) | null = null;
  onmessage: ((event: MessageEvent) => void) | null = null;
  onerror: ((event: Event) => void) | null = null;
  onclose: ((event: CloseEvent) => void) | null = null;

  constructor(url: string) {
    this.url = url;
  }

  send(_data: string): void {
    if (this.readyState !== MockWebSocket.READY_STATE_OPEN) {
      throw new Error('WebSocket is not open');
    }
  }

  close(code?: number, reason?: string): void {
    this.readyState = MockWebSocket.READY_STATE_CLOSED;
    if (this.onclose) {
      this.onclose({ code: code ?? 1000, reason: reason ?? '', wasClean: true } as CloseEvent);
    }
  }
}

describe('WebSocketManager (wsManager singleton)', () => {
  beforeEach(() => {
    // Disconnect and reset state before each test
    wsManager.disconnect();
    // Mock WebSocket
    vi.stubGlobal('WebSocket', MockWebSocket);
  });

  afterEach(() => {
    wsManager.disconnect();
    vi.unstubAllGlobals();
  });

  describe('initialization', () => {
    it('should have default disconnected state', () => {
      expect(wsManager.getStatus()).toBe('disconnected');
      expect(wsManager.isConnected()).toBe(false);
    });

    it('should update config', () => {
      wsManager.updateConfig({ url: 'ws://localhost:8080', token: 'test-token' });
      // Config update doesn't change status
      expect(wsManager.getStatus()).toBe('disconnected');
    });
  });

  describe('connection', () => {
    it('should connect to WebSocket server', () => {
      wsManager.updateConfig({ url: 'ws://localhost:8080', token: 'test-token' });
      wsManager.connect();

      expect(wsManager.getStatus()).toBe('connecting');
    });

    it('should disconnect manually', () => {
      wsManager.updateConfig({ url: 'ws://localhost:8080', token: 'test-token' });
      wsManager.connect();
      wsManager.disconnect();

      expect(wsManager.getStatus()).toBe('disconnected');
    });

    it('should not send message when not connected', () => {
      const result = wsManager.send('test', { data: 'test' });
      expect(result).toBe(false);
    });
  });

  describe('message handling', () => {
    it('should register message handler', () => {
      const handler = vi.fn();
      const unsubscribe = wsManager.on('system_status' as TMessageType, handler);

      expect(typeof unsubscribe).toBe('function');
    });

    it('should unregister message handler', () => {
      const handler = vi.fn();
      wsManager.on('test' as TMessageType, handler);
      expect(() => wsManager.off('test' as TMessageType, handler)).not.toThrow();
    });
  });

  describe('event listeners', () => {
    it('should set event listeners', () => {
      const onOpen = vi.fn();
      const onError = vi.fn();

      expect(() => wsManager.setEventListeners({ onOpen, onError })).not.toThrow();
    });

    it('should clear event listeners', () => {
      wsManager.setEventListeners({
        onOpen: vi.fn(),
        onError: vi.fn(),
      });

      expect(() => wsManager.clearEventListeners()).not.toThrow();
    });
  });

  describe('status', () => {
    it('should return correct status', () => {
      expect(wsManager.getStatus()).toBe('disconnected');
    });

    it('should return false for isConnected when not connected', () => {
      expect(wsManager.isConnected()).toBe(false);
    });
  });

  describe('reconnect', () => {
    it('should update reconnect config', () => {
      expect(() => wsManager.updateConfig({ reconnect: true })).not.toThrow();
      expect(() => wsManager.updateConfig({ reconnect: false })).not.toThrow();
    });
  });
});
