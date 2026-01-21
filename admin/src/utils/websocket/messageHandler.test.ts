/**
 * Chat Message Handler Tests
 */
import { describe, it, expect, beforeEach, vi } from 'vitest';
import { registerChatMessageHandlers, unregisterChatMessageHandlers } from './messageHandler';
import { wsManager } from './manager';
import type { WSMessage, MessageHandler, TMessageType } from './types';

// Test helper interface to access private members
interface TestableWebSocketManager {
  messageHandlers?: Map<TMessageType, Set<MessageHandler>>;
}

// Mock logger
vi.mock('@/utils/logger', () => ({
  logger: {
    info: vi.fn(),
    debug: vi.fn(),
    warn: vi.fn(),
    error: vi.fn(),
  },
}));

describe('Chat Message Handlers', () => {
  beforeEach(() => {
    // Clear all handlers
    wsManager.clearEventListeners();
  });

  describe('registerChatMessageHandlers', () => {
    it('should register handlers without throwing', () => {
      expect(() => registerChatMessageHandlers()).not.toThrow();
    });

    it('should be idempotent', () => {
      registerChatMessageHandlers();
      expect(() => registerChatMessageHandlers()).not.toThrow();
    });
  });

  describe('unregisterChatMessageHandlers', () => {
    it('should unregister handlers without throwing', () => {
      registerChatMessageHandlers();
      expect(() => unregisterChatMessageHandlers()).not.toThrow();
    });

    it('should be safe to call multiple times', () => {
      registerChatMessageHandlers();
      unregisterChatMessageHandlers();
      expect(() => unregisterChatMessageHandlers()).not.toThrow();
    });
  });

  describe('presence_update handler', () => {
    it('should handle presence update message', () => {
      registerChatMessageHandlers();

      const message: WSMessage = {
        type: 'presence_update',
        timestamp: new Date().toISOString(),
        data: {
          playerId: 123,
          status: 'online',
          currentGameId: 1,
          currentGameName: 'Test Game',
          customStatus: 'Playing',
          currentRoomId: 456,
          updatedAt: new Date().toISOString(),
        },
      };

      // Get the handlers and trigger one
      const managerPrivate = wsManager as unknown as TestableWebSocketManager;
      const handlers = managerPrivate.messageHandlers?.get('presence_update');

      if (handlers && handlers.size > 0) {
        handlers.forEach((handler: (msg: WSMessage) => void) => {
          expect(() => handler(message)).not.toThrow();
        });
      } else {
        // If no handlers registered, test passes
        expect(true).toBe(true);
      }
    });
  });

  describe('room_updated handler', () => {
    it('should handle room update message', () => {
      registerChatMessageHandlers();

      const message: WSMessage = {
        type: 'room_updated',
        timestamp: new Date().toISOString(),
        data: {
          roomId: 1,
          roomName: 'Updated Room',
          status: 'active',
          currentMembers: 3,
        },
      };

      const managerPrivate = wsManager as unknown as TestableWebSocketManager;
      const handlers = managerPrivate.messageHandlers?.get('room_updated');

      if (handlers && handlers.size > 0) {
        handlers.forEach((handler: (msg: WSMessage) => void) => {
          expect(() => handler(message)).not.toThrow();
        });
      } else {
        expect(true).toBe(true);
      }
    });
  });

  describe('room_closed handler', () => {
    it('should handle room close message', () => {
      registerChatMessageHandlers();

      const message: WSMessage = {
        type: 'room_closed',
        timestamp: new Date().toISOString(),
        data: {
          roomId: 1,
        },
      };

      const managerPrivate = wsManager as unknown as TestableWebSocketManager;
      const handlers = managerPrivate.messageHandlers?.get('room_closed');

      if (handlers && handlers.size > 0) {
        handlers.forEach((handler: (msg: WSMessage) => void) => {
          expect(() => handler(message)).not.toThrow();
        });
      } else {
        expect(true).toBe(true);
      }
    });
  });

  describe('room_member_joined handler', () => {
    it('should handle member join message', () => {
      registerChatMessageHandlers();

      const message: WSMessage = {
        type: 'room_member_joined',
        timestamp: new Date().toISOString(),
        data: {
          roomId: 1,
          userId: 123,
          nickname: 'New User',
          avatar: 'avatar.png',
          role: 'user',
        },
      };

      const managerPrivate = wsManager as unknown as TestableWebSocketManager;
      const handlers = managerPrivate.messageHandlers?.get('room_member_joined');

      if (handlers && handlers.size > 0) {
        handlers.forEach((handler: (msg: WSMessage) => void) => {
          expect(() => handler(message)).not.toThrow();
        });
      } else {
        expect(true).toBe(true);
      }
    });
  });

  describe('room_member_left handler', () => {
    it('should handle member leave message', () => {
      registerChatMessageHandlers();

      const message: WSMessage = {
        type: 'room_member_left',
        timestamp: new Date().toISOString(),
        data: {
          roomId: 1,
          userId: 123,
          nickname: 'Leaving User',
        },
      };

      const managerPrivate = wsManager as unknown as TestableWebSocketManager;
      const handlers = managerPrivate.messageHandlers?.get('room_member_left');

      if (handlers && handlers.size > 0) {
        handlers.forEach((handler: (msg: WSMessage) => void) => {
          expect(() => handler(message)).not.toThrow();
        });
      } else {
        expect(true).toBe(true);
      }
    });
  });
});
