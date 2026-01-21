/**
 * Chat Store - Complete Implementation
 * 聊天消息管理状态管理
 *
 * Features:
 * - 聊天室列表管理
 * - 消息历史（分页加载）
 * - 实时消息（WebSocket 集成）
 * - 未读计数
 * - 消息发送队列
 */

import { create } from 'zustand';
import { logger } from '@/utils/logger';
import {
  chatConversationApi,
  chatMessageApi,
  type ChatConversation,
  type ChatMessage as ApiChatMessage,
} from '@/api/chat';
import { useAuthStore } from './authStore';
import { registerChatMessageHandlers, unregisterChatMessageHandlers } from '@/utils/websocket';

// ==================== 类型定义 ====================

/** 聊天室信息（扩展 Conversation） */
export interface ChatRoom extends ChatConversation {
  unreadCount: number;
}

/** 消息分页状态 */
interface MessagePagination {
  page: number;
  pageSize: number;
  total: number;
  hasMore: boolean;
}

/** 消息发送队列项 */
interface PendingMessage {
  tempId: string;
  conversationId: number;
  content: string;
  messageType: 'text' | 'image' | 'voice';
  timestamp: number;
  retryCount: number;
}

/** 聊天状态管理 */
interface ChatState {
  // 聊天室列表
  rooms: ChatRoom[];
  roomsLoading: boolean;
  roomsTotal: number;
  roomsPage: number;
  roomsPageSize: number;
  roomsHasMore: boolean;

  // 当前聊天室
  currentRoomId: number | null;
  currentRoom: ChatRoom | null;

  // 消息列表（按聊天室分组）
  messages: Record<number, ApiChatMessage[]>;
  messagesLoading: boolean;
  messagesPagination: Record<number, MessagePagination>;

  // 未读计数
  totalUnread: number;

  // 发送队列
  pendingMessages: PendingMessage[];
  sendingMessages: boolean;

  // ==================== 聊天室操作 ====================

  /**
   * 获取聊天室列表
   * @param refresh 是否刷新（重置分页）
   * @param params 查询参数
   */
  fetchRooms: (
    refresh?: boolean,
    params?: {
      type?: 'user_order' | 'group_chat';
      status?: 'active' | 'closed';
      keyword?: string;
    }
  ) => Promise<void>;

  /**
   * 加载更多聊天室
   */
  loadMoreRooms: () => Promise<void>;

  /**
   * 设置当前聊天室
   * @param room 聊天室对象或ID
   */
  setCurrentRoom: (room: ChatRoom | number | null) => Promise<void>;

  /**
   * 更新聊天室信息
   * @param roomId 聊天室ID
   * @param updates 更新内容
   */
  updateRoom: (roomId: number, updates: Partial<ChatRoom>) => void;

  /**
   * 关闭聊天室
   * @param roomId 聊天室ID
   * @param reason 关闭原因
   */
  closeRoom: (roomId: number, reason?: string) => Promise<void>;

  /**
   * 重新打开聊天室
   * @param roomId 聊天室ID
   */
  reopenRoom: (roomId: number) => Promise<void>;

  // ==================== 消息操作 ====================

  /**
   * 获取聊天室消息列表
   * @param conversationId 会话ID
   * @param refresh 是否刷新（重置分页）
   */
  fetchMessages: (conversationId: number, refresh?: boolean) => Promise<void>;

  /**
   * 加载更多消息（向上翻页）
   * @param conversationId 会话ID
   */
  loadMoreMessages: (conversationId: number) => Promise<void>;

  /**
   * 添加消息到聊天室
   * @param conversationId 会话ID
   * @param message 消息对象
   */
  addMessage: (conversationId: number, message: ApiChatMessage) => void;

  /**
   * 批量添加消息
   * @param conversationId 会话ID
   * @param newMessages 新消息列表
   * @param prepend 是否添加到开头（用于加载历史消息）
   */
  addMessages: (
    conversationId: number,
    newMessages: ApiChatMessage[],
    prepend?: boolean
  ) => void;

  /**
   * 删除消息
   * @param conversationId 会话ID
   * @param messageId 消息ID
   */
  deleteMessage: (conversationId: number, messageId: number) => void;

  /**
   * 清空聊天室消息
   * @param conversationId 会话ID
   */
  clearMessages: (conversationId: number) => void;

  // ==================== 未读计数 ====================

  /**
   * 标记聊天室为已读
   * @param conversationId 会话ID
   */
  markAsRead: (conversationId: number) => void;

  /**
   * 增加聊天室未读数
   * @param conversationId 会话ID
   * @param count 增加数量
   */
  incrementUnread: (conversationId: number, count?: number) => void;

  /**
   * 重新计算总未读数
   */
  recalculateTotalUnread: () => void;

  // ==================== 消息发送 ====================

  /**
   * 发送消息
   * @param conversationId 会话ID
   * @param content 消息内容
   * @param messageType 消息类型
   */
  sendMessage: (
    conversationId: number,
    content: string,
    messageType?: 'text' | 'image' | 'voice'
  ) => Promise<void>;

  /**
   * 重试发送失败的消息
   * @param tempId 临时消息ID
   */
  retryMessage: (tempId: string) => Promise<void>;

  /**
   * 取消发送消息
   * @param tempId 临时消息ID
   */
  cancelMessage: (tempId: string) => void;

  /**
   * 处理发送队列
   */
  processSendQueue: () => Promise<void>;

  // ==================== WebSocket 集成 ====================

  /**
   * 初始化 WebSocket 监听
   * 用于接收实时消息
   */
  initWebSocket: () => void;

  /**
   * 清理 WebSocket 监听
   */
  cleanupWebSocket: () => void;

  // ==================== 工具方法 ====================

  /**
   * 重置所有状态
   */
  reset: () => void;

  /**
   * 获取聊天室最后一条消息
   * @param conversationId 会话ID
   */
  getLastMessage: (conversationId: number) => ApiChatMessage | undefined;
}

// ==================== 常量定义 ====================

const MESSAGES_PAGE_SIZE = 50; // 每次加载消息数量
const ROOMS_PAGE_SIZE = 20; // 每次加载聊天室数量
const MAX_RETRY_COUNT = 3; // 最大重试次数

// ==================== Store 实现 ====================

export const useChatStore = create<ChatState>((set, get) => ({
  // ==================== 初始状态 ====================

  // 聊天室列表
  rooms: [],
  roomsLoading: false,
  roomsTotal: 0,
  roomsPage: 1,
  roomsPageSize: ROOMS_PAGE_SIZE,
  roomsHasMore: true,

  // 当前聊天室
  currentRoomId: null,
  currentRoom: null,

  // 消息列表
  messages: {},
  messagesLoading: false,
  messagesPagination: {},

  // 未读计数
  totalUnread: 0,

  // 发送队列
  pendingMessages: [],
  sendingMessages: false,

  // ==================== 聊天室操作实现 ====================

  /**
   * 获取聊天室列表
   */
  fetchRooms: async (refresh = false, params) => {
    const state = get();

    // 如果正在加载，不重复请求
    if (state.roomsLoading) {
      return;
    }

    // 重置分页
    if (refresh) {
      set({
        roomsPage: 1,
        roomsHasMore: true,
      });
    }

    // 没有更多数据
    if (!refresh && !state.roomsHasMore) {
      return;
    }

    set({ roomsLoading: true });

    try {
      const { roomsPage, roomsPageSize } = get();
      const response = await chatConversationApi.getConversations({
        page: roomsPage,
        pageSize: roomsPageSize,
        ...params,
      });

      if (response.data && response.data.data && response.data.data.items) {
        const fetchedRooms = response.data.data.items.map(
          (room) =>
            ({
              ...room,
              unreadCount: 0, // 初始化未读数为0
            }) as ChatRoom
        );

        set({
          rooms:
            refresh || roomsPage === 1
              ? fetchedRooms
              : [...state.rooms, ...fetchedRooms],
          roomsTotal: response.data.data.total,
          roomsPage: roomsPage + 1,
          roomsHasMore: (refresh ? 0 : state.rooms.length) + fetchedRooms.length < response.data.data.total,
        });

        // 重新计算未读数
        get().recalculateTotalUnread();
      }
    } catch (error) {
      logger.error('Failed to fetch chat rooms:', error);
    } finally {
      set({ roomsLoading: false });
    }
  },

  /**
   * 加载更多聊天室
   */
  loadMoreRooms: async () => {
    await get().fetchRooms(false);
  },

  /**
   * 设置当前聊天室
   */
  setCurrentRoom: async (room) => {
    let roomId: number | null = null;
    let roomData: ChatRoom | null = null;

    if (typeof room === 'number') {
      roomId = room;
      // 从列表中查找
      roomData = get().rooms.find((r) => r.id === roomId) || null;
    } else if (room) {
      roomId = room.id;
      roomData = room;
    }

    set({
      currentRoomId: roomId,
      currentRoom: roomData,
    });

    // 标记为已读
    if (roomId) {
      get().markAsRead(roomId);

      // 如果消息未加载，则加载消息
      const state = get();
      if (!state.messages[roomId] || state.messages[roomId].length === 0) {
        await get().fetchMessages(roomId, true);
      }
    }
  },

  /**
   * 更新聊天室信息
   */
  updateRoom: (roomId, updates) => {
    set((state) => {
      const roomIndex = state.rooms.findIndex((r) => r.id === roomId);
      const newRooms = [...state.rooms];
      const newCurrentRoom =
        state.currentRoom?.id === roomId
          ? { ...state.currentRoom, ...updates }
          : state.currentRoom;

      if (roomIndex !== -1) {
        newRooms[roomIndex] = { ...newRooms[roomIndex], ...updates };
      }

      return {
        rooms: newRooms,
        currentRoom: newCurrentRoom,
      };
    });
  },

  /**
   * 关闭聊天室
   */
  closeRoom: async (roomId, reason) => {
    try {
      // 获取当前管理员 ID
      const userInfo = useAuthStore.getState().userInfo;
      const closedBy = userInfo?.id ?? 0;

      await chatConversationApi.closeConversation(roomId, {
        reason,
        closedBy,
      });
      get().updateRoom(roomId, { status: 'closed' } as Partial<ChatRoom>);
    } catch (error) {
      logger.error('Failed to close room:', error);
      throw error;
    }
  },

  /**
   * 重新打开聊天室
   */
  reopenRoom: async (roomId) => {
    try {
      await chatConversationApi.reopenConversation(roomId);
      get().updateRoom(roomId, { status: 'active' } as Partial<ChatRoom>);
    } catch (error) {
      logger.error('Failed to reopen room:', error);
      throw error;
    }
  },

  // ==================== 消息操作实现 ====================

  /**
   * 获取聊天室消息列表
   */
  fetchMessages: async (conversationId, refresh = false) => {
    const state = get();

    // 防止重复加载
    if (state.messagesLoading) {
      return;
    }

    // 初始化分页状态
    const currentPagination = state.messagesPagination[conversationId];
    if (refresh || !currentPagination) {
      set({
        messagesPagination: {
          ...state.messagesPagination,
          [conversationId]: {
            page: 1,
            pageSize: MESSAGES_PAGE_SIZE,
            total: 0,
            hasMore: true,
          },
        },
      });
    }

    const pagination = get().messagesPagination[conversationId];

    // 没有更多数据
    if (!refresh && !pagination.hasMore) {
      return;
    }

    set({ messagesLoading: true });

    try {
      const response = await chatMessageApi.getMessages(conversationId, {
        page: pagination.page,
        pageSize: pagination.pageSize,
      });

      if (response.data && response.data.data && response.data.data.items) {
        const fetchedMessages = response.data.data.items;
        const existingMessages = state.messages[conversationId] || [];

        set({
          messages: {
            ...state.messages,
            [conversationId]:
              refresh || pagination.page === 1
                ? fetchedMessages
                : [...fetchedMessages, ...existingMessages],
          },
          messagesPagination: {
            ...state.messagesPagination,
            [conversationId]: {
              page: pagination.page + 1,
              pageSize: pagination.pageSize,
              total: response.data.data.total,
              hasMore:
                (refresh ? 0 : existingMessages.length) + fetchedMessages.length <
                response.data.data.total,
            },
          },
        });
      }
    } catch (error) {
      logger.error('Failed to fetch messages:', error);
    } finally {
      set({ messagesLoading: false });
    }
  },

  /**
   * 加载更多消息
   */
  loadMoreMessages: async (conversationId) => {
    await get().fetchMessages(conversationId, false);
  },

  /**
   * 添加单条消息
   */
  addMessage: (conversationId, message) => {
    set((state) => {
      const roomMessages = state.messages[conversationId] || [];
      const newMessages = [...roomMessages, message];

      // 更新聊天室最后消息信息
      const roomIndex = state.rooms.findIndex((r) => r.id === conversationId);
      const newRooms = [...state.rooms];
      if (roomIndex !== -1) {
        newRooms[roomIndex] = {
          ...newRooms[roomIndex],
          lastMessageContent: message.content,
          lastMessageAt: message.createdAt,
          messageCount: newRooms[roomIndex].messageCount + 1,
        };
      }

      return {
        messages: {
          ...state.messages,
          [conversationId]: newMessages,
        },
        rooms: newRooms,
      };
    });
  },

  /**
   * 批量添加消息
   */
  addMessages: (conversationId, newMessages, prepend = false) => {
    set((state) => {
      const roomMessages = state.messages[conversationId] || [];
      const updatedMessages = prepend
        ? [...newMessages, ...roomMessages]
        : [...roomMessages, ...newMessages];

      return {
        messages: {
          ...state.messages,
          [conversationId]: updatedMessages,
        },
      };
    });
  },

  /**
   * 删除消息
   */
  deleteMessage: (conversationId, messageId) => {
    set((state) => {
      const roomMessages = state.messages[conversationId] || [];
      return {
        messages: {
          ...state.messages,
          [conversationId]: roomMessages.filter((m) => m.id !== messageId),
        },
      };
    });
  },

  /**
   * 清空聊天室消息
   */
  clearMessages: (conversationId) => {
    set((state) => ({
      messages: {
        ...state.messages,
        [conversationId]: [],
      },
      messagesPagination: {
        ...state.messagesPagination,
        [conversationId]: {
          page: 1,
          pageSize: MESSAGES_PAGE_SIZE,
          total: 0,
          hasMore: true,
        },
      },
    }));
  },

  // ==================== 未读计数实现 ====================

  /**
   * 标记聊天室为已读
   */
  markAsRead: (conversationId) => {
    set((state) => {
      const roomIndex = state.rooms.findIndex((r) => r.id === conversationId);
      if (roomIndex !== -1 && state.rooms[roomIndex].unreadCount > 0) {
        const newRooms = [...state.rooms];
        newRooms[roomIndex] = { ...newRooms[roomIndex], unreadCount: 0 };
        return { rooms: newRooms };
      }
      return {};
    });
    get().recalculateTotalUnread();
  },

  /**
   * 增加未读数
   */
  incrementUnread: (conversationId, count = 1) => {
    set((state) => {
      const roomIndex = state.rooms.findIndex((r) => r.id === conversationId);
      if (roomIndex !== -1) {
        const newRooms = [...state.rooms];
        newRooms[roomIndex] = {
          ...newRooms[roomIndex],
          unreadCount: newRooms[roomIndex].unreadCount + count,
        };
        return { rooms: newRooms };
      }
      return {};
    });
    get().recalculateTotalUnread();
  },

  /**
   * 重新计算总未读数
   */
  recalculateTotalUnread: () => {
    const { rooms } = get();
    const total = rooms.reduce((sum, room) => sum + room.unreadCount, 0);
    set({ totalUnread: total });
  },

  // ==================== 消息发送实现 ====================

  /**
   * 发送消息
   */
  sendMessage: async (conversationId, content, messageType = 'text') => {
    const tempId = `temp-${Date.now()}-${Math.random()}`;

    // 添加到待发送队列
    const pendingMessage: PendingMessage = {
      tempId,
      conversationId,
      content,
      messageType,
      timestamp: Date.now(),
      retryCount: 0,
    };

    set((state) => ({
      pendingMessages: [...state.pendingMessages, pendingMessage],
    }));

    // 获取当前用户信息
    const userInfo = useAuthStore.getState().userInfo;

    // 立即显示临时消息
    const tempChatMessage: ApiChatMessage = {
      id: 0,
      conversationId,
      senderId: userInfo?.id ?? 0,
      senderName: userInfo?.name ?? '管理员',
      senderType: 'user', // 管理员发送消息使用 user 类型
      content,
      messageType,
      isDeleted: false,
      createdAt: new Date().toISOString(),
    };

    get().addMessage(conversationId, tempChatMessage);

    // 处理发送队列
    await get().processSendQueue();
  },

  /**
   * 重试发送失败的消息
   */
  retryMessage: async (tempId) => {
    set((state) => ({
      pendingMessages: state.pendingMessages.map((msg) =>
        msg.tempId === tempId ? { ...msg, retryCount: msg.retryCount + 1 } : msg
      ),
    }));
    await get().processSendQueue();
  },

  /**
   * 取消发送消息
   */
  cancelMessage: (tempId) => {
    set((state) => {
      const index = state.pendingMessages.findIndex((m) => m.tempId === tempId);
      if (index !== -1) {
        const msg = state.pendingMessages[index];
        const newPending = [...state.pendingMessages];
        newPending.splice(index, 1);

        // 移除临时消息
        const roomMessages = state.messages[msg.conversationId] || [];
        return {
          pendingMessages: newPending,
          messages: {
            ...state.messages,
            [msg.conversationId]: roomMessages.filter((m) => m.id !== 0),
          },
        };
      }
      return {};
    });
  },

  /**
   * 处理发送队列
   */
  processSendQueue: async () => {
    const state = get();

    // 正在发送或队列为空
    if (state.sendingMessages || state.pendingMessages.length === 0) {
      return;
    }

    set({ sendingMessages: true });

    try {
      const firstPending = state.pendingMessages[0];
      if (!firstPending) {
        return;
      }

      // 注意: 当前管理员面板不直接发送消息
      // 如需实现此功能，需在后端添加管理员发送消息的 API
      // 参考: POST /api/v1/admin/chat/conversations/:id/messages
      //
      // 实现示例:
      // await chatMessageApi.sendMessage(firstPending.conversationId, {
      //   content: firstPending.content,
      //   messageType: firstPending.messageType,
      // });

      // 模拟发送成功（开发测试用）
      await new Promise((resolve) => setTimeout(resolve, 300));

      // 移除已发送的消息
      set((state) => ({
        pendingMessages: state.pendingMessages.slice(1),
      }));

      // 继续处理下一条
      if (get().pendingMessages.length > 0) {
        await get().processSendQueue();
      }
    } catch (error) {
      logger.error('Failed to send message:', error);

      // 检查是否需要重试
      const firstPending = get().pendingMessages[0];
      if (firstPending && firstPending.retryCount >= MAX_RETRY_COUNT) {
        // 超过最大重试次数，移除消息
        set((state) => ({
          pendingMessages: state.pendingMessages.slice(1),
        }));
      }
    } finally {
      set({ sendingMessages: false });
    }
  },

  // ==================== WebSocket 集成实现 ====================

  /**
   * 初始化 WebSocket 监听
   */
  initWebSocket: () => {
    // 注册聊天消息处理器
    registerChatMessageHandlers();

    // WebSocket 连接在 App.tsx 中统一管理
    // 这里只注册消息处理器，不负责连接管理
    logger.info('Chat WebSocket message handlers registered');
  },

  /**
   * 清理 WebSocket 监听
   */
  cleanupWebSocket: () => {
    // 取消聊天消息处理器
    unregisterChatMessageHandlers();

    logger.info('Chat WebSocket message handlers unregistered');
  },

  // ==================== 工具方法实现 ====================

  /**
   * 重置所有状态
   */
  reset: () => {
    set({
      rooms: [],
      roomsLoading: false,
      roomsTotal: 0,
      roomsPage: 1,
      roomsHasMore: true,

      currentRoomId: null,
      currentRoom: null,

      messages: {},
      messagesLoading: false,
      messagesPagination: {},

      totalUnread: 0,

      pendingMessages: [],
      sendingMessages: false,
    });
  },

  /**
   * 获取聊天室最后一条消息
   */
  getLastMessage: (conversationId) => {
    const { messages } = get();
    const roomMessages = messages[conversationId];
    return roomMessages && roomMessages.length > 0
      ? roomMessages[roomMessages.length - 1]
      : undefined;
  },
}));

// ==================== 选择器 ====================

/**
 * 获取聊天室消息列表
 */
export const selectRoomMessages = (conversationId: number) =>
  useChatStore.getState().messages[conversationId] || [];

/**
 * 获取聊天室未读数
 */
export const selectRoomUnread = (conversationId: number) => {
  const room = useChatStore
    .getState()
    .rooms.find((r) => r.id === conversationId);
  return room?.unreadCount || 0;
};

/**
 * 获取消息分页状态
 */
export const selectMessagesPagination = (conversationId: number) =>
  useChatStore.getState().messagesPagination[conversationId];
