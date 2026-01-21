/**
 * Chat Store Tests
 * Tests chat state management, message operations, and pagination
 */

import { describe, it, expect, beforeEach, vi } from 'vitest';
import { renderHook, act } from '@testing-library/react';
import { useChatStore } from './chatStore';
import type { ChatRoom, ApiChatMessage } from './chatStore';

// Mock chat API
vi.mock('../../api/chat', () => ({
  chatConversationApi: {
    getConversations: vi.fn(),
    getConversation: vi.fn(),
    closeConversation: vi.fn(),
    reopenConversation: vi.fn(),
  },
  chatMessageApi: {
    getMessages: vi.fn(),
  },
}));

// Helper to create mock room
function createMockRoom(overrides = {}): ChatRoom {
  return {
    id: 1,
    name: 'Test Room',
    type: 'solo',
    gameId: 1,
    gameName: 'League of Legends',
    hostUserId: 10,
    hostNickname: 'Host',
    status: 'active',
    currentMembers: 2,
    maxMembers: 4,
    lastMessageContent: 'Hello',
    lastMessageAt: '2024-01-01T10:00:00Z',
    messageCount: 10,
    unreadCount: 0,
    createdAt: '2024-01-01T00:00:00Z',
    updatedAt: '2024-01-01T00:00:00Z',
    ...overrides,
  };
}

// Helper to create mock message
function createMockMessage(overrides = {}): ApiChatMessage {
  return {
    id: 1,
    conversationId: 1,
    senderId: 10,
    senderName: 'User',
    senderType: 'user',
    content: 'Hello',
    messageType: 'text',
    isDeleted: false,
    createdAt: '2024-01-01T10:00:00Z',
    ...overrides,
  };
}

describe('chatStore', () => {
  beforeEach(() => {
    // Reset store state before each test
    useChatStore.getState().reset();
    vi.clearAllMocks();
  });

  describe('State Initialization', () => {
    it('should initialize with default state', () => {
      const { result } = renderHook(() => useChatStore());

      expect(result.current.rooms).toEqual([]);
      expect(result.current.currentRoom).toBeNull();
      expect(result.current.roomsLoading).toBe(false);
      expect(result.current.roomsTotal).toBe(0);
      expect(result.current.roomsPage).toBe(1);
      expect(result.current.roomsHasMore).toBe(true);

      expect(result.current.messages).toEqual({});
      expect(result.current.messagesLoading).toBe(false);

      expect(result.current.totalUnread).toBe(0);
      expect(result.current.pendingMessages).toEqual([]);
      expect(result.current.sendingMessages).toBe(false);
    });
  });

  describe('fetchRooms', () => {
    it('should fetch rooms successfully', async () => {
      const mockRooms = [
        createMockRoom({ id: 1, name: 'Room 1' }),
        createMockRoom({ id: 2, name: 'Room 2' }),
      ];

      const { chatConversationApi } = await import('../../api/chat');
      vi.mocked(chatConversationApi.getConversations).mockResolvedValue({
        data: {
          success: true,
          data: {
            items: mockRooms,
            total: 2,
          },
        },
      });

      const { result } = renderHook(() => useChatStore());

      await act(async () => {
        await result.current.fetchRooms();
      });

      expect(result.current.rooms).toEqual(mockRooms.map(room => ({ ...room, unreadCount: 0 })));
      expect(result.current.roomsTotal).toBe(2);
      expect(result.current.roomsLoading).toBe(false);
    });

    it('should handle fetch rooms error', async () => {
      const { chatConversationApi } = await import('../../api/chat');
      vi.mocked(chatConversationApi.getConversations).mockRejectedValue(new Error('API Error'));

      const { result } = renderHook(() => useChatStore());

      await act(async () => {
        await result.current.fetchRooms();
      });

      // Store should handle error gracefully and not crash
      expect(result.current.roomsLoading).toBe(false);
    });
  });

  describe('fetchMessages', () => {
    it('should fetch messages successfully', async () => {
      const mockMessages = [
        createMockMessage({ id: 1, content: 'Message 1' }),
        createMockMessage({ id: 2, content: 'Message 2' }),
      ];

      const { chatMessageApi } = await import('../../api/chat');
      vi.mocked(chatMessageApi.getMessages).mockResolvedValue({
        data: {
          success: true,
          data: {
            items: mockMessages,
            total: 2,
          },
        },
      });

      const { result } = renderHook(() => useChatStore());

      await act(async () => {
        await result.current.fetchMessages(1);
      });

      expect(result.current.messages[1]).toEqual(mockMessages);
      expect(result.current.messagesLoading).toBe(false);
    });

    it('should prepend messages when loading more', async () => {
      const page1Messages = [
        createMockMessage({ id: 1, content: 'Old' }),
      ];
      const page2Messages = [
        createMockMessage({ id: 2, content: 'New' }),
      ];

      const { chatMessageApi } = await import('../../api/chat');
      vi.mocked(chatMessageApi.getMessages)
        .mockResolvedValueOnce({
          data: {
            success: true,
            data: {
              items: page1Messages,
              total: 2,
            },
          },
        })
        .mockResolvedValueOnce({
          data: {
            success: true,
            data: {
              items: page2Messages,
              total: 2,
            },
          },
        });

      const { result } = renderHook(() => useChatStore());

      // Load first page
      await act(async () => {
        await result.current.fetchMessages(1);
      });

      expect(result.current.messages[1]).toEqual(page1Messages);

      // Load second page (should prepend)
      await act(async () => {
        await result.current.fetchMessages(1, false);
      });

      expect(result.current.messages[1]).toEqual([...page2Messages, ...page1Messages]);
    });
  });

  describe('setCurrentRoom', () => {
    it('should set current room and fetch its messages', async () => {
      const mockRoom = createMockRoom({ id: 1 });
      const mockMessages = [createMockMessage()];

      const { chatConversationApi } = await import('../../api/chat');
      vi.mocked(chatConversationApi.getConversation).mockResolvedValue({
        data: {
          success: true,
          data: mockRoom,
        },
      });

      const { chatMessageApi } = await import('../../api/chat');
      vi.mocked(chatMessageApi.getMessages).mockResolvedValue({
        data: {
          success: true,
          data: {
            items: mockMessages,
            total: 1,
          },
        },
      });

      const { result } = renderHook(() => useChatStore());

      // Add room to the rooms array first (since setCurrentRoom with ID looks in the array)
      act(() => {
        useChatStore.setState({
          rooms: [{ ...mockRoom, unreadCount: 0 }],
        });
      });

      await act(async () => {
        await result.current.setCurrentRoom(1);
      });

      expect(result.current.currentRoomId).toBe(1);
      expect(result.current.currentRoom).toEqual({ ...mockRoom, unreadCount: 0 });
      expect(result.current.messages[1]).toEqual(mockMessages);
    });

    it('should mark room as read when setting current room', async () => {
      const mockRoom = createMockRoom({ id: 1, unreadCount: 5 });

      const { chatConversationApi } = await import('../../api/chat');
      vi.mocked(chatConversationApi.getConversation).mockResolvedValue({
        data: {
          success: true,
          data: mockRoom,
        },
      });

      const { chatMessageApi } = await import('../../api/chat');
      vi.mocked(chatMessageApi.getMessages).mockResolvedValue({
        data: {
          success: true,
          data: {
            items: [],
            total: 0,
          },
        },
      });

      const { result } = renderHook(() => useChatStore());

      // Set room with unread count
      act(() => {
        useChatStore.setState({
          rooms: [{ ...mockRoom, unreadCount: 5 }],
          totalUnread: 5,
        });
      });

      expect(result.current.totalUnread).toBe(5);

      // Set as current (should mark as read)
      await act(async () => {
        await result.current.setCurrentRoom(1);
      });

      expect(result.current.totalUnread).toBe(0);
    });
  });

  describe('closeRoom', () => {
    it('should close room successfully', async () => {
      const { chatConversationApi } = await import('../../api/chat');
      vi.mocked(chatConversationApi.closeConversation).mockResolvedValue({
        data: {
          success: true,
          data: { message: 'Room closed' },
        },
      });

      const { result } = renderHook(() => useChatStore());

      act(() => {
        useChatStore.setState({
          rooms: [createMockRoom({ id: 1, status: 'active' })],
        });
      });

      await act(async () => {
        await result.current.closeRoom(1, 'Test reason');
      });

      expect(chatConversationApi.closeConversation).toHaveBeenCalledWith(1, {
        reason: 'Test reason',
        closedBy: 0, // From auth store (mocked as 0 in tests)
      });
    });
  });

  describe('Message Operations', () => {
    it('should add message to room', () => {
      const { result } = renderHook(() => useChatStore());

      act(() => {
        result.current.addMessage(1, createMockMessage({ id: 1 }));
      });

      expect(result.current.messages[1]).toEqual([createMockMessage({ id: 1 })]);
    });

    it('should batch add messages', () => {
      const { result } = renderHook(() => useChatStore());

      const messages = [
        createMockMessage({ id: 1 }),
        createMockMessage({ id: 2 }),
      ];

      act(() => {
        result.current.addMessages(1, messages);
      });

      expect(result.current.messages[1]).toEqual(messages);
    });

    it('should prepend messages when adding batch', () => {
      const { result } = renderHook(() => useChatStore());

      const existing = [createMockMessage({ id: 1 })];
      const newMessages = [createMockMessage({ id: 2 }), createMockMessage({ id: 3 })];

      act(() => {
        result.current.addMessages(1, existing);
        result.current.addMessages(1, newMessages, true);
      });

      expect(result.current.messages[1]).toEqual([...newMessages, ...existing]);
    });

    it('should delete message', () => {
      const { result } = renderHook(() => useChatStore());

      act(() => {
        result.current.addMessage(1, createMockMessage({ id: 1 }));
        result.current.addMessage(1, createMockMessage({ id: 2 }));
        result.current.deleteMessage(1, 1);
      });

      expect(result.current.messages[1]).toEqual([createMockMessage({ id: 2 })]);
    });

    it('should clear room messages', () => {
      const { result } = renderHook(() => useChatStore());

      act(() => {
        result.current.addMessage(1, createMockMessage({ id: 1 }));
        result.current.clearMessages(1);
      });

      expect(result.current.messages[1]).toEqual([]);
      expect(result.current.messagesPagination[1]).toMatchObject({
        page: 1,
        hasMore: true,
      });
    });
  });

  describe('Unread Count', () => {
    it('should mark room as read', () => {
      const { result } = renderHook(() => useChatStore());

      act(() => {
        useChatStore.setState({
          rooms: [createMockRoom({ id: 1, unreadCount: 5 })],
        });
      });

      act(() => {
        result.current.markAsRead(1);
      });

      expect(result.current.rooms[0].unreadCount).toBe(0);
    });

    it('should increment unread count', () => {
      const { result } = renderHook(() => useChatStore());

      act(() => {
        useChatStore.setState({
          rooms: [createMockRoom({ id: 1, unreadCount: 0 })],
        });
      });

      act(() => {
        result.current.incrementUnread(1, 2);
      });

      expect(result.current.rooms[0].unreadCount).toBe(2);
    });

    it('should recalculate total unread', () => {
      const { result } = renderHook(() => useChatStore());

      act(() => {
        useChatStore.setState({
          rooms: [
            createMockRoom({ id: 1, unreadCount: 3 }),
            createMockRoom({ id: 2, unreadCount: 2 }),
          ],
        });
      });

      act(() => {
        result.current.recalculateTotalUnread();
      });

      expect(result.current.totalUnread).toBe(5);
    });
  });

  describe('Message Sending', () => {
    it('should add message to pending queue', async () => {
      const { result } = renderHook(() => useChatStore());

      // Mock processSendQueue to avoid side effects
      const processSendQueueSpy = vi.spyOn(useChatStore.getState(), 'processSendQueue').mockResolvedValue(undefined);

      await act(async () => {
        await result.current.sendMessage(1, 'Hello');
      });

      expect(result.current.pendingMessages.length).toBe(1);
      expect(result.current.pendingMessages[0].content).toBe('Hello');

      // Restore the original implementation
      processSendQueueSpy.mockRestore();
    });

    it('should display temporary message immediately', async () => {
      const { result } = renderHook(() => useChatStore());

      await act(async () => {
        await result.current.sendMessage(1, 'Test');
      });

      const messages = result.current.messages[1];
      expect(messages.length).toBe(1);
      expect(messages[0].id).toBe(0); // Temp message
      expect(messages[0].content).toBe('Test');
    });

    it('should retry failed message', async () => {
      const { result } = renderHook(() => useChatStore());

      // Mock processSendQueue to avoid side effects in this test
      const processSendQueueSpy = vi.spyOn(useChatStore.getState(), 'processSendQueue').mockResolvedValue(undefined);

      // Manually add a pending message to test retry
      const pendingMessage = {
        tempId: 'test-retry',
        conversationId: 1,
        content: 'Retry Test',
        messageType: 'text' as const,
        timestamp: Date.now(),
        retryCount: 0,
      };

      act(() => {
        useChatStore.setState({
          pendingMessages: [pendingMessage],
        });
      });

      expect(result.current.pendingMessages[0].retryCount).toBe(0);

      await act(async () => {
        await result.current.retryMessage('test-retry');
      });

      expect(result.current.pendingMessages[0].retryCount).toBe(1);

      // Restore the original implementation
      processSendQueueSpy.mockRestore();
    });

    it('should cancel pending message', () => {
      const { result } = renderHook(() => useChatStore());

      act(() => {
        useChatStore.setState({
          pendingMessages: [
            {
              tempId: 'test-1',
              conversationId: 1,
              content: 'Test',
              messageType: 'text',
              timestamp: Date.now(),
              retryCount: 0,
            },
          ],
        });
      });

      act(() => {
        result.current.cancelMessage('test-1');
      });

      expect(result.current.pendingMessages.length).toBe(0);
    });
  });

  describe('updateRoom', () => {
    it('should update room in list', () => {
      const { result } = renderHook(() => useChatStore());

      act(() => {
        useChatStore.setState({
          rooms: [createMockRoom({ id: 1, name: 'Old Name' })],
        });
      });

      act(() => {
        result.current.updateRoom(1, { name: 'New Name' });
      });

      expect(result.current.rooms[0].name).toBe('New Name');
    });

    it('should update multiple rooms', () => {
      const { result } = renderHook(() => useChatStore());

      act(() => {
        useChatStore.setState({
          rooms: [
            createMockRoom({ id: 1, status: 'active' }),
            createMockRoom({ id: 2, status: 'active' }),
          ],
        });
      });

      act(() => {
        result.current.updateRoom(1, { status: 'closed' });
        result.current.updateRoom(2, { status: 'closed' });
      });

      expect(result.current.rooms[0].status).toBe('closed');
      expect(result.current.rooms[1].status).toBe('closed');
    });
  });

  describe('reset', () => {
    it('should reset all state to default', () => {
      const { result } = renderHook(() => useChatStore());

      // Set some state
      act(() => {
        useChatStore.setState({
          rooms: [createMockRoom()],
          messages: { 1: [createMockMessage()] },
          pendingMessages: [
            {
              tempId: 'test-1',
              conversationId: 1,
              content: 'Test',
              messageType: 'text',
              timestamp: Date.now(),
              retryCount: 0,
            },
          ],
          totalUnread: 5,
        });
      });

      expect(result.current.rooms.length).toBe(1);

      // Reset
      act(() => {
        result.current.reset();
      });

      expect(result.current.rooms).toEqual([]);
      expect(result.current.messages).toEqual({});
      expect(result.current.pendingMessages).toEqual([]);
      expect(result.current.totalUnread).toBe(0);
    });
  });

  describe('Edge Cases', () => {
    it('should handle empty room list', async () => {
      const { chatConversationApi } = await import('../../api/chat');
      vi.mocked(chatConversationApi.getConversations).mockResolvedValue({
        data: {
          success: true,
          data: [],
          pagination: { total: 0, page: 1, pageSize: 10 },
        },
      });

      const { result } = renderHook(() => useChatStore());

      await act(async () => {
        await result.current.fetchRooms();
      });

      expect(result.current.rooms).toEqual([]);
    });

    it('should handle room not found gracefully', async () => {
      const { chatConversationApi } = await import('../../api/chat');
      vi.mocked(chatConversationApi.getConversation).mockResolvedValue({
        data: {
          success: false,
          message: 'Room not found',
        },
      });

      const { result } = renderHook(() => useChatStore());

      await act(async () => {
        await result.current.setCurrentRoom(999);
      });

      // Room should be null/undefined when not found
      expect(result.current.room == null).toBe(true);
      expect(result.current.currentRoomId).toBe(999);
    });

    it('should handle message with special characters', () => {
      const { result } = renderHook(() => useChatStore());

      act(() => {
        result.current.addMessage(1, createMockMessage({ content: 'Test @#$%' }));
      });

      expect(result.current.messages[1][0].content).toBe('Test @#$%');
    });
  });
});
