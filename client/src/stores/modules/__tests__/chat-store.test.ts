/**
 * Chat Store Tests
 * Tests for chat conversations, messages, and real-time communication
 */

import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { useChatStore } from '../chat-store';

// Mock http module
vi.mock('@/lib/http', () => ({
    http: {
        get: vi.fn(),
        post: vi.fn(),
        put: vi.fn(),
        delete: vi.fn(),
    },
}));

// Mock import.meta.env
vi.stubGlobal('import', {
    meta: {
        env: {
            VITE_USE_MOCK: 'false',
        },
    },
});

import { http } from '@/lib/http';

const mockHttp = http as unknown as {
    get: ReturnType<typeof vi.fn>;
    post: ReturnType<typeof vi.fn>;
    put: ReturnType<typeof vi.fn>;
    delete: ReturnType<typeof vi.fn>;
};

describe('Chat Store', () => {
    beforeEach(() => {
        vi.clearAllMocks();
        // Reset store state
        useChatStore.setState({
            conversations: [],
            currentConversationId: null,
            messages: {},
            totalUnreadCount: 0,
            isConnected: false,
            loading: false,
            error: null,
        });
    });

    afterEach(() => {
        vi.restoreAllMocks();
    });

    describe('Initial State', () => {
        it('should have correct initial state', () => {
            const state = useChatStore.getState();

            expect(state.conversations).toEqual([]);
            expect(state.currentConversationId).toBeNull();
            expect(state.messages).toEqual({});
            expect(state.totalUnreadCount).toBe(0);
            expect(state.isConnected).toBe(false);
            expect(state.loading).toBe(false);
            expect(state.error).toBeNull();
        });
    });

    describe('fetchConversations', () => {
        it('should fetch conversations successfully', async () => {
            const mockConversations = [
                {
                    id: '1',
                    participantId: 101,
                    participantName: 'Player1',
                    participantAvatar: 'avatar1.jpg',
                    lastMessage: 'Hello!',
                    lastMessageTime: '2024-01-01T00:00:00Z',
                    unreadCount: 2,
                    online: true,
                },
                {
                    id: '2',
                    participantId: 102,
                    participantName: 'Player2',
                    participantAvatar: 'avatar2.jpg',
                    lastMessage: 'GG!',
                    lastMessageTime: '2024-01-02T00:00:00Z',
                    unreadCount: 0,
                    online: false,
                },
            ];

            mockHttp.get.mockResolvedValueOnce(mockConversations);

            await useChatStore.getState().fetchConversations();

            const state = useChatStore.getState();
            expect(state.conversations).toHaveLength(2);
            expect(state.conversations[0].participantName).toBe('Player1');
            expect(state.totalUnreadCount).toBe(2);
            expect(state.loading).toBe(false);
        });

        it('should handle fetch error', async () => {
            mockHttp.get.mockRejectedValueOnce(new Error('Network error'));

            await useChatStore.getState().fetchConversations();

            const state = useChatStore.getState();
            expect(state.loading).toBe(false);
            expect(state.error).toBe('Network error');
        });
    });

    describe('selectConversation', () => {
        it('should select conversation and mark as read', async () => {
            useChatStore.setState({
                conversations: [
                    {
                        id: '1',
                        participantId: 101,
                        participantName: 'Player1',
                        participantAvatar: 'avatar1.jpg',
                        lastMessage: 'Hello!',
                        lastMessageTime: '2024-01-01',
                        unreadCount: 3,
                        online: true,
                    },
                ],
                totalUnreadCount: 3,
            });

            await useChatStore.getState().selectConversation('1');

            const state = useChatStore.getState();
            expect(state.currentConversationId).toBe('1');
            expect(state.conversations[0].unreadCount).toBe(0);
            expect(state.totalUnreadCount).toBe(0);
        });

        it('should load messages for new conversation', async () => {
            useChatStore.setState({
                conversations: [
                    {
                        id: '1',
                        participantId: 101,
                        participantName: 'Player1',
                        participantAvatar: 'avatar1.jpg',
                        lastMessage: 'Hello!',
                        lastMessageTime: '2024-01-01',
                        unreadCount: 0,
                        online: true,
                    },
                ],
                messages: {},
            });

            await useChatStore.getState().selectConversation('1');

            const state = useChatStore.getState();
            expect(state.currentConversationId).toBe('1');
            // Messages should be loaded (mock data in store)
            expect(state.messages['1']).toBeDefined();
        });

        it('should not reload messages if already loaded', async () => {
            const existingMessages = [
                {
                    id: 'm1',
                    conversationId: '1',
                    senderId: 101,
                    content: 'Existing message',
                    type: 'text' as const,
                    createdAt: '2024-01-01',
                    read: true,
                },
            ];

            useChatStore.setState({
                conversations: [
                    {
                        id: '1',
                        participantId: 101,
                        participantName: 'Player1',
                        participantAvatar: 'avatar1.jpg',
                        lastMessage: 'Hello!',
                        lastMessageTime: '2024-01-01',
                        unreadCount: 0,
                        online: true,
                    },
                ],
                messages: { '1': existingMessages },
            });

            await useChatStore.getState().selectConversation('1');

            const state = useChatStore.getState();
            expect(state.messages['1']).toEqual(existingMessages);
        });
    });

    describe('sendMessage', () => {
        it('should send message with optimistic update', async () => {
            const mockResponse = {
                id: 1,
                groupId: 1,
                senderId: 100,
                content: 'Hello!',
                messageType: 'text',
                createdAt: '2024-01-01T00:00:00Z',
                isDeleted: false
            };
            mockHttp.post.mockResolvedValueOnce(mockResponse);

            useChatStore.setState({
                currentConversationId: 1,
                messages: { 1: [] },
                conversations: [{
                    id: 1,
                    groupName: 'Test',
                    groupType: 'public',
                    participantId: 101,
                    participantName: 'Player1',
                    participantAvatar: 'avatar1.jpg',
                    lastMessage: '',
                    lastMessageTime: '',
                    unreadCount: 0,
                    online: true,
                    isActive: true,
                    isPrivate: false,
                }],
            });

            await useChatStore.getState().sendMessage('Hello!');

            const state = useChatStore.getState();
            expect(state.messages[1]).toHaveLength(1);
            expect(state.messages[1][0].content).toBe('Hello!');
            expect(state.messages[1][0].messageType).toBe('text');
            expect(mockHttp.post).toHaveBeenCalledWith('/chat/groups/1/messages', { content: 'Hello!', messageType: 'text' });
        });

        it('should send image message', async () => {
            const mockResponse = {
                id: 2,
                groupId: 1,
                senderId: 100,
                content: 'image-url.jpg',
                messageType: 'image',
                createdAt: '2024-01-01T00:00:00Z',
                isDeleted: false
            };
            mockHttp.post.mockResolvedValueOnce(mockResponse);

            useChatStore.setState({
                currentConversationId: 1,
                messages: { 1: [] },
                conversations: [{
                    id: 1,
                    groupName: 'Test',
                    groupType: 'public',
                    participantId: 101,
                    participantName: 'Player1',
                    participantAvatar: 'avatar1.jpg',
                    lastMessage: '',
                    lastMessageTime: '',
                    unreadCount: 0,
                    online: true,
                    isActive: true,
                    isPrivate: false,
                }],
            });

            await useChatStore.getState().sendMessage('image-url.jpg', 'image');

            const state = useChatStore.getState();
            expect(state.messages[1][0].messageType).toBe('image');
        });

        it('should not send if no conversation selected', async () => {
            useChatStore.setState({
                currentConversationId: null,
                messages: {},
            });

            await useChatStore.getState().sendMessage('Hello!');

            const state = useChatStore.getState();
            expect(Object.keys(state.messages)).toHaveLength(0);
            expect(mockHttp.post).not.toHaveBeenCalled();
        });

        it('should create messages array if not exists', async () => {
            const mockResponse = {
                id: 3,
                groupId: 1,
                senderId: 100,
                content: 'Hello!',
                messageType: 'text',
                createdAt: '2024-01-01T00:00:00Z',
                isDeleted: false
            };
            mockHttp.post.mockResolvedValueOnce(mockResponse);

            useChatStore.setState({
                currentConversationId: 1,
                messages: {},
                conversations: [{
                    id: 1,
                    groupName: 'Test',
                    groupType: 'public',
                    participantId: 101,
                    participantName: 'Player1',
                    participantAvatar: 'avatar1.jpg',
                    lastMessage: '',
                    lastMessageTime: '',
                    unreadCount: 0,
                    online: true,
                    isActive: true,
                    isPrivate: false,
                }],
            });

            await useChatStore.getState().sendMessage('Hello!');

            const state = useChatStore.getState();
            expect(state.messages[1]).toBeDefined();
            expect(state.messages[1]).toHaveLength(1);
        });

        it('should rollback on send error', async () => {
            mockHttp.post.mockRejectedValueOnce(new Error('Network error'));

            useChatStore.setState({
                currentConversationId: 1,
                messages: { 1: [] },
                conversations: [],
            });

            await expect(useChatStore.getState().sendMessage('Hello!')).rejects.toThrow('Network error');

            const state = useChatStore.getState();
            expect(state.messages[1]).toHaveLength(0);
            expect(state.error).toBe('Network error');
        });
    });

    describe('receiveMessage', () => {
        it('should add received message to conversation', () => {
            useChatStore.setState({
                currentConversationId: 1,
                messages: { 1: [] },
                conversations: [{ id: 1, groupName: 'Test', groupType: 'public', participantId: 101, participantName: 'Player1', lastMessage: '', lastMessageTime: '', unreadCount: 0, online: true, isActive: true, isPrivate: false }],
                totalUnreadCount: 0,
            });

            const newMessage = {
                id: 1,
                groupId: 1,
                senderId: 101,
                content: 'New message',
                messageType: 'text' as const,
                createdAt: '2024-01-01',
                isDeleted: false,
            };

            useChatStore.getState().receiveMessage(newMessage);

            const state = useChatStore.getState();
            expect(state.messages[1]).toHaveLength(1);
            expect(state.messages[1][0].content).toBe('New message');
        });

        it('should increment unread count for non-current conversation', () => {
            useChatStore.setState({
                currentConversationId: 1,
                messages: { 2: [] },
                conversations: [{ id: 2, groupName: 'Test', groupType: 'public', participantId: 102, participantName: 'Player2', lastMessage: '', lastMessageTime: '', unreadCount: 0, online: true, isActive: true, isPrivate: false }],
                totalUnreadCount: 0,
            });

            const newMessage = {
                id: 1,
                groupId: 2, // Different from current
                senderId: 102,
                content: 'New message',
                messageType: 'text' as const,
                createdAt: '2024-01-01',
                isDeleted: false,
            };

            useChatStore.getState().receiveMessage(newMessage);

            expect(useChatStore.getState().totalUnreadCount).toBe(1);
        });

        it('should not increment unread count for current conversation', () => {
            useChatStore.setState({
                currentConversationId: 1,
                messages: { 1: [] },
                conversations: [{ id: 1, groupName: 'Test', groupType: 'public', participantId: 101, participantName: 'Player1', lastMessage: '', lastMessageTime: '', unreadCount: 0, online: true, isActive: true, isPrivate: false }],
                totalUnreadCount: 0,
            });

            const newMessage = {
                id: 1,
                groupId: 1, // Same as current
                senderId: 101,
                content: 'New message',
                messageType: 'text' as const,
                createdAt: '2024-01-01',
                isDeleted: false,
            };

            useChatStore.getState().receiveMessage(newMessage);

            expect(useChatStore.getState().totalUnreadCount).toBe(0);
        });

        it('should not add duplicate messages', () => {
            const existingMessage = {
                id: 1,
                groupId: 1,
                senderId: 101,
                content: 'Existing',
                messageType: 'text' as const,
                createdAt: '2024-01-01',
                isDeleted: false,
            };

            useChatStore.setState({
                currentConversationId: 1,
                messages: { 1: [existingMessage] },
                conversations: [],
            });

            // Try to add same message again
            useChatStore.getState().receiveMessage(existingMessage);

            expect(useChatStore.getState().messages[1]).toHaveLength(1);
        });

        it('should create messages array for new conversation', () => {
            useChatStore.setState({
                currentConversationId: null,
                messages: {},
                conversations: [{ id: 3, groupName: 'Test', groupType: 'public', participantId: 103, participantName: 'Player3', lastMessage: '', lastMessageTime: '', unreadCount: 0, online: true, isActive: true, isPrivate: false }],
                totalUnreadCount: 0,
            });

            const newMessage = {
                id: 1,
                groupId: 3,
                senderId: 103,
                content: 'New conversation',
                messageType: 'text' as const,
                createdAt: '2024-01-01',
                isDeleted: false,
            };

            useChatStore.getState().receiveMessage(newMessage);

            const state = useChatStore.getState();
            expect(state.messages[3]).toBeDefined();
            expect(state.messages[3]).toHaveLength(1);
        });
    });

    describe('markAsRead', () => {
        it('should mark conversation as read', async () => {
            useChatStore.setState({
                conversations: [
                    {
                        id: '1',
                        participantId: 101,
                        participantName: 'Player1',
                        participantAvatar: 'avatar1.jpg',
                        lastMessage: 'Hello!',
                        lastMessageTime: '2024-01-01',
                        unreadCount: 5,
                        online: true,
                    },
                    {
                        id: '2',
                        participantId: 102,
                        participantName: 'Player2',
                        participantAvatar: 'avatar2.jpg',
                        lastMessage: 'Hi!',
                        lastMessageTime: '2024-01-02',
                        unreadCount: 3,
                        online: false,
                    },
                ],
                totalUnreadCount: 8,
            });

            await useChatStore.getState().markAsRead('1');

            const state = useChatStore.getState();
            expect(state.conversations[0].unreadCount).toBe(0);
            expect(state.conversations[1].unreadCount).toBe(3);
            expect(state.totalUnreadCount).toBe(3);
        });

        it('should handle already read conversation', async () => {
            useChatStore.setState({
                conversations: [
                    {
                        id: '1',
                        participantId: 101,
                        participantName: 'Player1',
                        participantAvatar: 'avatar1.jpg',
                        lastMessage: 'Hello!',
                        lastMessageTime: '2024-01-01',
                        unreadCount: 0,
                        online: true,
                    },
                ],
                totalUnreadCount: 0,
            });

            await useChatStore.getState().markAsRead('1');

            const state = useChatStore.getState();
            expect(state.conversations[0].unreadCount).toBe(0);
            expect(state.totalUnreadCount).toBe(0);
        });
    });
});
