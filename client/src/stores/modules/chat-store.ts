import { create } from 'zustand';
import { http } from '@/lib/http';
import { getErrorMessage, logError } from '@/lib/error';

// ============ Enums (Match backend api/internal/model/chat.go) ============

export const ChatGroupType = {
    PUBLIC: 'public',       // 公开群组
    ORDER: 'order',         // 订单房间
    TEAM: 'team',           // 组队房间
    LFG: 'lfg',             // 快速匹配房间
    CUSTOM: 'custom'        // 自定义房间
} as const;

export type ChatGroupType = typeof ChatGroupType[keyof typeof ChatGroupType];

export const ChatGroupStatus = {
    WAITING: 'waiting',       // 等待中
    READY: 'ready',          // 准备就绪
    IN_GAME: 'in_game',     // 游戏中
    FINISHED: 'finished',   // 已结束
    CANCELED: 'canceled'    // 已取消
} as const;

export type ChatGroupStatus = typeof ChatGroupStatus[keyof typeof ChatGroupStatus];

export const ChatMessageType = {
    TEXT: 'text',           // 文本消息
    IMAGE: 'image',         // 图片消息
    FILE: 'file',           // 文件消息
    SYSTEM: 'system',       // 系统消息
    VOICE: 'voice',         // 语音消息（预留）
    EMOJI: 'emoji'          // 表情包（预留）
} as const;

export type ChatMessageType = typeof ChatMessageType[keyof typeof ChatMessageType];

export const ChatMessageAuditStatus = {
    PENDING: 'pending',     // 待审核
    APPROVED: 'approved',   // 已通过
    REJECTED: 'rejected',   // 已拒绝
    DELETED: 'deleted'      // 已删除
} as const;

export type ChatMessageAuditStatus = typeof ChatMessageAuditStatus[keyof typeof ChatMessageAuditStatus];

export const ChatMemberRole = {
    OWNER: 'owner',         // 群主
    ADMIN: 'admin',         // 管理员
    MEMBER: 'member'        // 普通成员
} as const;

export type ChatMemberRole = typeof ChatMemberRole[keyof typeof ChatMemberRole];

// ============ Interfaces (Match backend models) ============

export interface ChatMessage {
    id: number | string;          // Allow string for temporary messages (temp-xxx)
    groupId: number;              // 匹配后端 GroupID
    senderId: number;             // 匹配后端 SenderID
    content: string;
    messageType: ChatMessageType; // 匹配后端 MessageType (more types)
    replyToId?: number;          // 匹配后端 ReplyToID
    imageUrl?: string;           // 匹配后端 ImageURL
    metadata?: Record<string, unknown>; // 匹配后端 Metadata (JSON)
    isDeleted: boolean;          // 匹配后端 IsDeleted
    auditStatus?: ChatMessageAuditStatus; // 匹配后端 AuditStatus
    createdAt: string;
    updatedAt?: string;
    isTemporary?: boolean;       // Flag to identify temporary messages
}

export interface ChatGroup {
    id: number;
    groupName: string;
    groupType: ChatGroupType;
    relatedOrderId?: number;
    createdBy: number;
    maxMembers: number;
    isActive: boolean;
    avatarUrl?: string;
    description?: string;
    currentMembers: number;
    isPrivate: boolean;
    gameStatus?: ChatGroupStatus;
    voiceEnabled: boolean;
    createdAt: string;
    updatedAt: string;
}

export interface ChatGroupMember {
    groupId: number;
    userId: number;
    role: ChatMemberRole;
    nickname?: string;
    joinedAt: string;
    lastReadAt?: string;
    lastReadMessageId?: number;
    isMuted: boolean;
    isActive: boolean;
    isReady: boolean;              // 游戏房间准备状态
}

// Simplified conversation for frontend UI (derived from ChatGroup + last message)
export interface Conversation {
    id: number;                   // Same as groupId
    groupName: string;
    groupType: ChatGroupType;
    participantId: number;        // The other user (for 1:1 chats)
    participantName: string;
    participantAvatar?: string;
    lastMessage: string;
    lastMessageTime: string;
    unreadCount: number;          // Derived from lastReadMessageId vs latest message ID
    online: boolean;              // Participant online status
    isActive: boolean;
    isPrivate: boolean;
}

export interface ChatState {
    conversations: Conversation[];
    currentConversationId: number | null;
    messages: Record<number, ChatMessage[]>; // Keyed by groupId

    // Status
    totalUnreadCount: number;
    isConnected: boolean;
    loading: boolean;
    error: string | null;
}

export interface ChatActions {
    fetchConversations: () => Promise<void>;
    selectConversation: (groupId: number) => Promise<void>;
    fetchMessages: (groupId: number) => Promise<void>;
    sendMessage: (content: string, type?: ChatMessageType) => Promise<void>;
    receiveMessage: (message: ChatMessage) => void;
    markAsRead: (groupId: number) => Promise<void>;
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

    selectConversation: async (groupId: number) => {
        set({ currentConversationId: groupId });
        const { messages, markAsRead, fetchMessages } = get();

        // Mark as read immediately when selecting
        void markAsRead(groupId);

        // Fetch messages if not loaded
        if (!messages[groupId]) {
            await fetchMessages(groupId);
        }
    },

    fetchMessages: async (groupId: number) => {
        set({ loading: true });
        try {
            const data = await http.get<ChatMessage[]>(`/chat/groups/${groupId}/messages`);

            set((state) => ({
                messages: {
                    ...state.messages,
                    [groupId]: data || []
                },
                loading: false
            }));
        } catch (err) {
            logError('fetchMessages', err);
            set({ loading: false, error: getErrorMessage(err, 'Failed to fetch messages') });
        }
    },

    sendMessage: async (content, type = ChatMessageType.TEXT) => {
        const { currentConversationId } = get();
        if (!currentConversationId) return;

        // Create temporary message with string ID and isTemporary flag
        const tempId = `temp-${Date.now()}`;
        const newMessage: ChatMessage = {
            id: tempId,
            groupId: currentConversationId,
            senderId: 0, // Will be replaced by server response
            content,
            messageType: type,
            createdAt: new Date().toISOString(),
            isDeleted: false,
            isTemporary: true
        };

        set(state => ({
            messages: {
                ...state.messages,
                [currentConversationId]: [...(state.messages[currentConversationId] || []), newMessage]
            }
        }));

        try {
            const response = await http.post<ChatMessage>(
                `/chat/groups/${currentConversationId}/messages`,
                { content, messageType: type }
            );

            // Replace temp message with real one using isTemporary flag for safer matching
            set(state => ({
                messages: {
                    ...state.messages,
                    [currentConversationId]: state.messages[currentConversationId]?.map(m =>
                        (m.id === tempId || m.isTemporary) && m.content === content ? { ...response, isTemporary: false } : m
                    ) || []
                }
            }));

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
            const groupMsgs = state.messages[message.groupId] || [];
            // Check for duplicates
            if (groupMsgs.some(m => m.id === message.id)) return state;

            return {
                messages: {
                    ...state.messages,
                    [message.groupId]: [...groupMsgs, message]
                },
                // Update conversation's last message
                conversations: state.conversations.map(c =>
                    c.id === message.groupId
                        ? {
                            ...c,
                            lastMessage: message.content,
                            lastMessageTime: message.createdAt,
                            unreadCount: state.currentConversationId !== message.groupId
                                ? c.unreadCount + 1
                                : c.unreadCount
                        }
                        : c
                ),
                // Increment unread if not current conversation
                totalUnreadCount: state.currentConversationId !== message.groupId
                    ? state.totalUnreadCount + 1
                    : state.totalUnreadCount
            };
        });
    },

    markAsRead: async (groupId) => {
        // Optimistic update
        set(state => ({
            conversations: state.conversations.map(c =>
                c.id === groupId ? { ...c, unreadCount: 0 } : c
            ),
            totalUnreadCount: state.conversations.reduce((acc, curr) =>
                curr.id === groupId ? acc : acc + curr.unreadCount,
                0
            )
        }));

        try {
            await http.post(`/chat/groups/${groupId}/read`);
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
