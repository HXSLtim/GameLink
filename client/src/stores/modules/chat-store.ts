import { create } from 'zustand';
import { http } from '@/lib/http';
import { getErrorMessage, logError } from '@/lib/error';

export interface Message {
    id: number;
    conversationId: number;
    senderId: number;
    content: string;
    type: 'text' | 'image' | 'system';
    createdAt: string;
    read: boolean;
}

export interface Conversation {
    id: number;
    participantId: number;
    participantName: string;
    participantAvatar: string;
    lastMessage: string;
    lastMessageTime: string;
    unreadCount: number;
    online: boolean; // Participant online status
}

export interface ChatState {
    conversations: Conversation[];
    currentConversationId: number | null;
    messages: Record<number, Message[]>; // Keyed by conversationId

    // Status
    totalUnreadCount: number;
    isConnected: boolean;
    loading: boolean;
    error: string | null;
}

export interface ChatActions {
    fetchConversations: () => Promise<void>;
    selectConversation: (conversationId: number) => Promise<void>;
    fetchMessages: (conversationId: number) => Promise<void>;
    sendMessage: (content: string, type?: 'text' | 'image') => Promise<void>;
    receiveMessage: (message: Message) => void;
    markAsRead: (conversationId: number) => Promise<void>;
    setConnected: (connected: boolean) => void;
    reset: () => void;
}

const INITIAL_STATE: ChatState = {
    conversations: [],
    currentConversationId: null,
    messages: {},
    totalUnreadCount: 0,
    isConnected: false,
    loading: false,
    error: null,
};

export const useChatStore = create<ChatState & ChatActions>((set, get) => ({
    ...INITIAL_STATE,

    fetchConversations: async () => {
        set({ loading: true, error: null });
        try {
            const data = await http.get<Conversation[]>('/chat/conversations');
            const conversations = data || [];

            set({
                conversations,
                loading: false,
                totalUnreadCount: conversations.reduce((acc, curr) => acc + curr.unreadCount, 0)
            });
        } catch (err) {
            logError('fetchConversations', err);
            set({ loading: false, error: getErrorMessage(err, 'Failed to fetch conversations') });
        }
    },

    selectConversation: async (conversationId: number) => {
        set({ currentConversationId: conversationId });
        const { messages, markAsRead, fetchMessages } = get();

        // Mark as read immediately when selecting
        void markAsRead(conversationId);

        // Fetch messages if not loaded
        if (!messages[conversationId]) {
            await fetchMessages(conversationId);
        }
    },

    fetchMessages: async (conversationId: number) => {
        set({ loading: true });
        try {
            const data = await http.get<Message[]>(`/chat/conversations/${conversationId}/messages`);

            set((state) => ({
                messages: {
                    ...state.messages,
                    [conversationId]: data || []
                },
                loading: false
            }));
        } catch (err) {
            logError('fetchMessages', err);
            set({ loading: false, error: getErrorMessage(err, 'Failed to fetch messages') });
        }
    },

    sendMessage: async (content, type = 'text') => {
        const { currentConversationId } = get();
        if (!currentConversationId) return;

        // Optimistic update
        const tempId = `temp-${Date.now()}`;
        const newMessage: Message = {
            id: tempId,
            conversationId: currentConversationId,
            senderId: 0, // Will be replaced by server response
            content,
            type,
            createdAt: new Date().toISOString(),
            read: true
        };

        set(state => ({
            messages: {
                ...state.messages,
                [currentConversationId]: [...(state.messages[currentConversationId] || []), newMessage]
            }
        }));

        try {
            const response = await http.post<Message>(
                `/chat/conversations/${currentConversationId}/messages`,
                { content, type }
            );

            // Replace temp message with real one
            set(state => ({
                messages: {
                    ...state.messages,
                    [currentConversationId]: state.messages[currentConversationId]?.map(m =>
                        m.id === tempId ? response : m
                    ) || []
                }
            ));

            // Update conversation's last message
            set(state => ({
                conversations: state.conversations.map(c =>
                    c.id === currentConversationId
                        ? { ...c, lastMessage: content, lastMessageTime: new Date().toISOString() }
                        : c
                )
            }));
        } catch (err) {
            logError('sendMessage', err);
            // Rollback: Remove the temporary message
            set(state => ({
                messages: {
                    ...state.messages,
                    [currentConversationId]: state.messages[currentConversationId]?.filter(m => m.id !== tempId) || []
                },
                error: getErrorMessage(err, 'Failed to send message')
            }));
            throw err;
        }
    },

    receiveMessage: (message) => {
        set(state => {
            const conversationMsgs = state.messages[message.conversationId] || [];
            // Check for duplicates
            if (conversationMsgs.some(m => m.id === message.id)) return state;

            return {
                messages: {
                    ...state.messages,
                    [message.conversationId]: [...conversationMsgs, message]
                },
                // Update conversation's last message
                conversations: state.conversations.map(c =>
                    c.id === message.conversationId
                        ? {
                            ...c,
                            lastMessage: message.content,
                            lastMessageTime: message.createdAt,
                            unreadCount: state.currentConversationId !== message.conversationId
                                ? c.unreadCount + 1
                                : c.unreadCount
                        }
                        : c
                ),
                // Increment unread if not current conversation
                totalUnreadCount: state.currentConversationId !== message.conversationId
                    ? state.totalUnreadCount + 1
                    : state.totalUnreadCount
            };
        });
    },

    markAsRead: async (conversationId) => {
        // Optimistic update
        set(state => ({
            conversations: state.conversations.map(c =>
                c.id === conversationId ? { ...c, unreadCount: 0 } : c
            ),
            totalUnreadCount: state.conversations.reduce((acc, curr) =>
                curr.id === conversationId ? acc : acc + curr.unreadCount,
                0)
        }));

        try {
            await http.post(`/chat/conversations/${conversationId}/read`);
        } catch (err) {
            logError('markAsRead', err);
            // Don't throw - marking as read is not critical
        }
    },

    setConnected: (connected) => {
        set({ isConnected: connected });
    },

    reset: () => set(INITIAL_STATE)
}));

// ============ Selectors ============

export const selectCurrentMessages = (state: ChatState) => {
    if (!state.currentConversationId) return [];
    return state.messages[state.currentConversationId] || [];
};

export const selectCurrentConversation = (state: ChatState) => {
    if (!state.currentConversationId) return null;
    return state.conversations.find(c => c.id === state.currentConversationId) || null;
};

export const selectHasUnread = (state: ChatState) => state.totalUnreadCount > 0;
