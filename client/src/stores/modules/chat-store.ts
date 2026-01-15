import { create } from 'zustand';
import { http } from '@/lib/http';

export interface Message {
    id: string;
    conversationId: string;
    senderId: number;
    content: string;
    type: 'text' | 'image' | 'system';
    createdAt: string;
    read: boolean;
}

export interface Conversation {
    id: string;
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
    currentConversationId: string | null;
    messages: Record<string, Message[]>; // Keyed by conversationId

    // Status
    totalUnreadCount: number;
    isConnected: boolean;
    loading: boolean;
    error: string | null;
}

export interface ChatActions {
    fetchConversations: () => Promise<void>;
    selectConversation: (conversationId: string) => Promise<void>;
    sendMessage: (content: string, type?: 'text' | 'image') => Promise<void>;
    receiveMessage: (message: Message) => void;
    markAsRead: (conversationId: string) => Promise<void>;
}

export const useChatStore = create<ChatState & ChatActions>((set, get) => ({
    conversations: [],
    currentConversationId: null,
    messages: {},
    totalUnreadCount: 0,
    isConnected: false,
    loading: false,
    error: null,

    fetchConversations: async () => {
        set({ loading: true, error: null });
        try {
            const USE_MOCK = import.meta.env.VITE_USE_MOCK === 'true';

            if (USE_MOCK) {
                // Mock data for development
                const mockConversations: Conversation[] = [
                    {
                        id: '1',
                        participantId: 101,
                        participantName: "欢乐使者",
                        participantAvatar: "",
                        lastMessage: "Hello! Are you available?",
                        lastMessageTime: new Date().toISOString(),
                        unreadCount: 2,
                        online: true,
                    },
                    {
                        id: '2',
                        participantId: 102,
                        participantName: "游戏高手",
                        participantAvatar: "",
                        lastMessage: "GG well played!",
                        lastMessageTime: new Date(Date.now() - 3600000).toISOString(),
                        unreadCount: 0,
                        online: false,
                    }
                ];

                set({
                    conversations: mockConversations,
                    loading: false,
                    totalUnreadCount: mockConversations.reduce((acc, curr) => acc + curr.unreadCount, 0)
                });
            } else {
                // Real API call
                // Assuming the API returns Conversation[] directly or wrapped
                // Ideally this should be typed properly with generic Http response
                const data = await http.get<Conversation[]>('/chat/conversations');

                set({
                    conversations: data,
                    loading: false,
                    totalUnreadCount: data.reduce((acc, curr) => acc + curr.unreadCount, 0)
                });
            }
        } catch (err) {
            set({ loading: false, error: err instanceof Error ? err.message : 'Failed to fetch conversations' });
        }
    },

    selectConversation: async (conversationId: string) => {
        set({ currentConversationId: conversationId });
        const { messages, markAsRead } = get();

        // Mark as read immediately when selecting
        void markAsRead(conversationId);

        // Fetch messages if not loaded
        if (!messages[conversationId]) {
            set({ loading: true });
            try {
                // Mock API call
                // const history = await http.get<Message[]>(`/chat/conversations/${conversationId}/messages`);

                const mockHandler = conversationId === '1' ? [
                    { id: 'm1', conversationId, senderId: 101, content: 'Hi there!', type: 'text', createdAt: new Date(Date.now() - 60000).toISOString(), read: true },
                    { id: 'm2', conversationId, senderId: 999, content: 'Hello! Yes I am.', type: 'text', createdAt: new Date().toISOString(), read: true }, // 999 is self
                ] : [];

                set((state) => ({
                    messages: {
                        ...state.messages,
                        [conversationId]: mockHandler as Message[] // Cast for mock
                    },
                    loading: false
                }));

            } catch (err) {
                console.error(err);
                set({ loading: false });
            }
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
            senderId: 0, // 0 usually means 'me' in optimistic UI before real ID
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
            // await http.post(`/chat/conversations/${currentConversationId}/messages`, { content, type });
            // In real app, replace tempId with real ID from response

            // Simulating API call for now since we are in mock mode mostly
            // If VITE_USE_MOCK is false, we would await the real call
            const USE_MOCK = import.meta.env.VITE_USE_MOCK === 'true';
            if (!USE_MOCK) {
                // TODO: Uncomment when backend is ready
                // const response = await http.post<Message>(`/chat/conversations/${currentConversationId}/messages`, { content, type });
                // set(state => ({
                //     messages: {
                //         ...state.messages,
                //         [currentConversationId]: state.messages[currentConversationId]?.map(m =>
                //             m.id === tempId ? { ...m, id: response.id } : m
                //         )
                //     }
                // }));
            }

        } catch (err) {
            console.error("Failed to send message", err);
            // Rollback: Remove the temporary message
            set(state => ({
                messages: {
                    ...state.messages,
                    [currentConversationId]: state.messages[currentConversationId]?.filter(m => m.id !== tempId) || []
                },
                error: 'Failed to send message. Please try again.'
            }));
        }
    },

    receiveMessage: (message) => {
        set(state => {
            const conversationMsgs = state.messages[message.conversationId] || [];
            // Check for dups
            if (conversationMsgs.some(m => m.id === message.id)) return state;

            return {
                messages: {
                    ...state.messages,
                    [message.conversationId]: [...conversationMsgs, message]
                },
                // Increment unread if not current conversation
                totalUnreadCount: state.currentConversationId !== message.conversationId
                    ? state.totalUnreadCount + 1
                    : state.totalUnreadCount
            };
        });
    },

    markAsRead: async (conversationId) => {
        // Optimistic
        set(state => ({
            conversations: state.conversations.map(c =>
                c.id === conversationId ? { ...c, unreadCount: 0 } : c
            ),
            totalUnreadCount: state.conversations.reduce((acc, curr) =>
                curr.id === conversationId ? acc : acc + curr.unreadCount,
                0)
        }));

        try {
            // await http.post(`/chat/conversations/${conversationId}/read`);
        } catch (e) {
            console.error(e);
        }
    }
}));
