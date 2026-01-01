// Chat Store - Taro App
// Complete chat and messaging state management with WebSocket support

import { create } from 'zustand';
import Taro from '@tarojs/taro';
import type { ChatMessage, ChatRoom, ChatGroup, WSMessage, UploadProgress, ChatMessageType } from '../types';
import { get, post } from '../../api/client';

// WebSocket connection state
interface WebSocketState {
  ws: Taro.SocketTask | null;
  connected: boolean;
  reconnectAttempts: number;
  maxReconnectAttempts: number;
  reconnectDelay: number;
}

// Pending message for retry
interface PendingMessage {
  id: string;
  roomId: number;
  content: string;
  type: ChatMessageType;
  imageUrl?: string;
  timestamp: number;
  attempts: number;
}

// Constants
const MAX_MESSAGES_PER_ROOM = 100; // Memory optimization: limit messages per room
const MAX_SEND_ATTEMPTS = 3;
const BASE_RECONNECT_DELAY = 2000;
const MAX_RECONNECT_DELAY = 30000;

// Chat Store State
interface ChatState {
  // Room data
  rooms: ChatRoom[];
  currentRoom: ChatRoom | null;
  roomMessages: Record<number, ChatMessage[]>;
  loadingRooms: boolean;
  loadingMessages: boolean;

  // WebSocket state
  ws: WebSocketState;

  // Upload state
  uploads: UploadProgress[];

  // Pending messages (for retry)
  pendingMessages: PendingMessage[];

  // Actions - Room management
  fetchRooms: (page?: number, pageSize?: number) => Promise<void>;
  setCurrentRoom: (room: ChatRoom | null) => void;

  // Actions - Message management
  fetchMessages: (roomId: number, page?: number, pageSize?: number) => Promise<void>;
  sendMessage: (roomId: number, content: string, type: ChatMessageType, imageUrl?: string) => Promise<void>;
  sendImageMessage: (roomId: number, filePath: string) => Promise<void>;
  sendVoiceMessage: (roomId: number, filePath: string, duration: number) => Promise<void>;
  retryMessage: (pendingId: string) => Promise<void>;
  deletePendingMessage: (pendingId: string) => void;

  // Actions - Read status
  markAsRead: (roomId: number) => void;
  markMessageAsRead: (roomId: number, messageId: number) => void;

  // Actions - WebSocket
  connectWebSocket: () => void;
  disconnectWebSocket: () => void;
  sendWebSocketMessage: (message: WSMessage) => void;

  // Actions - Upload
  uploadFile: (roomId: number, filePath: string, type: 'image' | 'voice') => Promise<void>;

  // Actions - Cleanup
  clearOldMessages: (roomId: number, keepLast?: number) => void;
  clearAllMessages: () => void;
  reset: () => void;
}

// Initial WebSocket state
const initialWSState: WebSocketState = {
  ws: null,
  connected: false,
  reconnectAttempts: 0,
  maxReconnectAttempts: 10,
  reconnectDelay: BASE_RECONNECT_DELAY,
};

export const useChatStore = create<ChatState>((set, get) => ({
  // Initial state
  rooms: [],
  currentRoom: null,
  roomMessages: {},
  loadingRooms: false,
  loadingMessages: false,
  ws: initialWSState,
  uploads: [],
  pendingMessages: [],

  // ============================================
  // Room Management
  // ============================================

  fetchRooms: async (page = 1, pageSize = 20) => {
    const state = get();
    if (state.loadingRooms) return;

    set({ loadingRooms: true });

    try {
      const response = await get<{
        groups: ChatGroup[];
        total: number;
      }>('/chat/groups', {
        showLoading: page === 1,
        loadingText: '加载会话列表...',
      });

      if (response.success && response.data) {
        // Transform ChatGroup to ChatRoom
        const chatRooms: ChatRoom[] = response.data.groups.map((group) => ({
          id: group.id,
          groupId: group.id,
          orderId: group.relatedOrderId,
          name: group.groupName,
          avatar: group.avatarUrl,
          type: group.groupType,
          participants: group.members,
          lastMessage: group.lastMessage,
          unreadCount: group.unreadCount || 0,
          isActive: group.isActive,
          isMuted: false, // TODO: get from member data
          createdAt: group.createdAt,
        }));

        set({
          rooms: page === 1 ? chatRooms : [...state.rooms, ...chatRooms],
          loadingRooms: false,
        });
      } else {
        set({ loadingRooms: false });
        Taro.showToast({
          title: response.message || '加载失败',
          icon: 'none',
        });
      }
    } catch (error: any) {
      console.error('[chatStore] Fetch rooms error:', error);
      set({ loadingRooms: false });
      Taro.showToast({
        title: error.message || '网络错误',
        icon: 'none',
      });
    }
  },

  setCurrentRoom: (room: ChatRoom | null) => {
    set({ currentRoom: room });

    // Mark room as read when opening
    if (room && room.unreadCount > 0) {
      get().markAsRead(room.id);
    }

    // Load messages for the room
    if (room) {
      get().fetchMessages(room.id);
    } else {
      // Clear current messages when leaving room
      set({ roomMessages: {} });
    }
  },

  // ============================================
  // Message Management
  // ============================================

  fetchMessages: async (roomId: number, page = 1, pageSize = 50) => {
    const state = get();
    if (state.loadingMessages) return;

    set({ loadingMessages: true });

    try {
      const response = await get<{
        messages: ChatMessage[];
        total: number;
      }>(`/chat/groups/${roomId}/messages?page=${page}&pageSize=${pageSize}`, {
        showLoading: page === 1,
        loadingText: '加载消息...',
      });

      if (response.success && response.data) {
        const messages = response.data.messages.map((msg) => ({
          ...msg,
          sending: false,
          sendError: false,
          retryCount: 0,
        }));

        set({
          roomMessages: {
            ...state.roomMessages,
            [roomId]: page === 1 ? messages : [...messages, ...state.roomMessages[roomId]],
          },
          loadingMessages: false,
        });
      } else {
        set({ loadingMessages: false });
        Taro.showToast({
          title: response.message || '加载失败',
          icon: 'none',
        });
      }
    } catch (error: any) {
      console.error('[chatStore] Fetch messages error:', error);
      set({ loadingMessages: false });
      Taro.showToast({
        title: error.message || '网络错误',
        icon: 'none',
      });
    }
  },

  sendMessage: async (roomId: number, content: string, type: ChatMessageType, imageUrl?: string) => {
    const state = get();
    const currentUserId = getCurrentUserId();
    if (!currentUserId) {
      Taro.showToast({ title: '请先登录', icon: 'none' });
      return;
    }

    // Create optimistic message
    const tempId = Date.now();
    const optimisticMessage: ChatMessage = {
      id: tempId,
      groupId: roomId,
      senderId: currentUserId,
      content,
      messageType: type,
      imageUrl,
      isDeleted: false,
      auditStatus: 'pending',
      createdAt: new Date().toISOString(),
      sending: true,
      sendError: false,
      retryCount: 0,
    };

    // Add to messages immediately
    set({
      roomMessages: {
        ...state.roomMessages,
        [roomId]: [...(state.roomMessages[roomId] || []), optimisticMessage],
      },
    });

    try {
      const response = await post<ChatMessage>(`/chat/groups/${roomId}/messages`, {
        content,
        messageType: type,
        imageUrl,
      });

      if (response.success && response.data) {
        // Replace optimistic message with real message
        set({
          roomMessages: {
            ...state.roomMessages,
            [roomId]: state.roomMessages[roomId]?.map((msg) =>
              msg.id === tempId ? { ...response.data!, sending: false } : msg
            ) || [response.data!],
          },
        });
      } else {
        throw new Error(response.message || '发送失败');
      }
    } catch (error: any) {
      console.error('[chatStore] Send message error:', error);

      // Mark message as error
      set({
        roomMessages: {
          ...state.roomMessages,
          [roomId]: state.roomMessages[roomId]?.map((msg) =>
            msg.id === tempId ? { ...msg, sending: false, sendError: true } : msg
          ) || [],
        },
      });

      // Add to pending messages for retry
      const pendingId = `pending-${tempId}`;
      const pendingMessage: PendingMessage = {
        id: pendingId,
        roomId,
        content,
        type,
        imageUrl,
        timestamp: Date.now(),
        attempts: 1,
      };

      set({
        pendingMessages: [...state.pendingMessages, pendingMessage],
      });

      Taro.showToast({
        title: '发送失败，可重试',
        icon: 'none',
      });

      // Vibrate to notify user
      try {
        Taro.vibrateShort({ type: 'error' });
      } catch (e) {
        // Ignore vibration errors
      }
    }
  },

  sendImageMessage: async (roomId: number, filePath: string) => {
    await get().uploadFile(roomId, filePath, 'image');
  },

  sendVoiceMessage: async (roomId: number, filePath: string, duration: number) => {
    await get().uploadFile(roomId, filePath, 'voice');
  },

  retryMessage: async (pendingId: string) => {
    const state = get();
    const pendingMsg = state.pendingMessages.find((m) => m.id === pendingId);

    if (!pendingMsg) return;

    // Check retry limit
    if (pendingMsg.attempts >= MAX_SEND_ATTEMPTS) {
      Taro.showToast({
        title: '重试次数已达上限',
        icon: 'none',
      });
      return;
    }

    // Remove from pending and update attempt count
    set({
      pendingMessages: state.pendingMessages.filter((m) => m.id !== pendingId),
    });

    // Update pending message
    const updatedPendingMsg = {
      ...pendingMsg,
      attempts: pendingMsg.attempts + 1,
    };

    // Retry sending
    await get().sendMessage(
      updatedPendingMsg.roomId,
      updatedPendingMsg.content,
      updatedPendingMsg.type,
      updatedPendingMsg.imageUrl
    );
  },

  deletePendingMessage: (pendingId: string) => {
    const state = get();
    set({
      pendingMessages: state.pendingMessages.filter((m) => m.id !== pendingId),
    });
  },

  // ============================================
  // Read Status
  // ============================================

  markAsRead: (roomId: number) => {
    set((state) => ({
      rooms: state.rooms.map((r) =>
        r.id === roomId ? { ...r, unreadCount: 0 } : r
      ),
    }));

    // Notify server via WebSocket
    const wsState = get().ws;
    if (wsState.connected && wsState.ws) {
      const wsMessage: WSMessage = {
        type: 'chat_read',
        timestamp: new Date().toISOString(),
        data: {
          groupId: roomId,
        },
      };
      get().sendWebSocketMessage(wsMessage);
    }
  },

  markMessageAsRead: (roomId: number, messageId: number) => {
    // Update local read status
    set((state) => ({
      roomMessages: {
        ...state.roomMessages,
        [roomId]: state.roomMessages[roomId]?.map((msg) =>
          msg.id === messageId ? { ...msg, readAt: new Date().toISOString() } : msg
        ) || [],
      },
    }));

    // Notify server
    const wsState = get().ws;
    if (wsState.connected && wsState.ws) {
      const wsMessage: WSMessage = {
        type: 'chat_read',
        timestamp: new Date().toISOString(),
        data: {
          groupId: roomId,
          messageId,
        },
      };
      get().sendWebSocketMessage(wsMessage);
    }
  },

  // ============================================
  // WebSocket
  // ============================================

  connectWebSocket: () => {
    const state = get();
    if (state.ws.connected) return;

    const token = Taro.getStorageSync('token');
    if (!token) {
      console.warn('[chatStore] No token, cannot connect WebSocket');
      return;
    }

    // Build WebSocket URL
    const wsUrl = `${process.env.TARO_APP_WS_BASE_URL || 'ws://localhost:8080'}/api/v1/user/ws/chat?token=${token}`;

    console.log('[chatStore] Connecting WebSocket...');

    const wsTask = Taro.connectSocket({
      url: wsUrl,
      header: {
        'Authorization': `Bearer ${token}`,
      },
    });

    wsTask.onOpen(() => {
      console.log('[chatStore] WebSocket connected');
      set({
        ws: {
          ...state.ws,
          ws: wsTask,
          connected: true,
          reconnectAttempts: 0,
          reconnectDelay: BASE_RECONNECT_DELAY,
        },
      });

      // Subscribe to chat updates
      get().sendWebSocketMessage({
        type: 'subscribe',
        timestamp: new Date().toISOString(),
        data: {
          topics: ['chat', 'messages'],
        },
      });
    });

    wsTask.onMessage((msg) => {
      try {
        const wsMessage: WSMessage = JSON.parse(msg.data as string);
        console.log('[chatStore] WebSocket message:', wsMessage);
        handleWebSocketMessage(wsMessage);
      } catch (error) {
        console.error('[chatStore] Failed to parse WebSocket message:', error);
      }
    });

    wsTask.onError((error) => {
      console.error('[chatStore] WebSocket error:', error);
      set({
        ws: {
          ...state.ws,
          connected: false,
        },
      });
    });

    wsTask.onClose(() => {
      console.log('[chatStore] WebSocket closed');
      set({
        ws: {
          ...state.ws,
          connected: false,
        },
      });

      // Attempt to reconnect
      attemptReconnect();
    });

    set({
      ws: {
        ...state.ws,
        ws: wsTask,
      },
    });
  },

  disconnectWebSocket: () => {
    const state = get();
    if (state.ws.ws) {
      state.ws.ws.close();
      set({
        ws: initialWSState,
      });
    }
  },

  sendWebSocketMessage: (message: WSMessage) => {
    const state = get();
    if (state.ws.connected && state.ws.ws) {
      state.ws.ws.send({
        data: JSON.stringify(message),
      });
    } else {
      console.warn('[chatStore] WebSocket not connected, cannot send message');
    }
  },

  // ============================================
  // File Upload
  // ============================================

  uploadFile: async (roomId: number, filePath: string, type: 'image' | 'voice') => {
    const fileId = `upload-${Date.now()}`;
    const fileName = filePath.split('/').pop() || 'unknown';

    // Add upload to state
    set((state) => ({
      uploads: [
        ...state.uploads,
        {
          fileId,
          fileName,
          progress: 0,
          status: 'uploading',
        },
      ],
    }));

    try {
      const token = Taro.getStorageSync('token');
      const uploadUrl = `${process.env.TARO_APP_API_BASE_URL || 'http://localhost:8080/api/v1'}/upload`;

      const uploadTask = Taro.uploadFile({
        url: uploadUrl,
        filePath,
        name: 'file',
        formData: {
          type,
        },
        header: {
          'Authorization': `Bearer ${token}`,
        },
      });

      uploadTask.onProgressUpdate((res) => {
        set((state) => ({
          uploads: state.uploads.map((u) =>
            u.fileId === fileId ? { ...u, progress: res.progress } : u
          ),
        }));
      });

      const result = await uploadTask;

      if (result.statusCode === 200) {
        const data = JSON.parse(result.data);
        if (data.success && data.data) {
          const fileUrl = data.data.url;

          // Update upload state
          set((state) => ({
            uploads: state.uploads.map((u) =>
              u.fileId === fileId
                ? { ...u, progress: 100, status: 'success', url: fileUrl }
                : u
            ),
          }));

          // Send message with uploaded file
          if (type === 'image') {
            await get().sendMessage(roomId, '', 'image', fileUrl);
          } else if (type === 'voice') {
            await get().sendMessage(roomId, fileUrl, 'voice');
          }

          // Remove upload from state after delay
          setTimeout(() => {
            set((state) => ({
              uploads: state.uploads.filter((u) => u.fileId !== fileId),
            }));
          }, 3000);
        } else {
          throw new Error(data.message || '上传失败');
        }
      } else {
        throw new Error('上传失败');
      }
    } catch (error: any) {
      console.error('[chatStore] Upload error:', error);

      // Update upload state to error
      set((state) => ({
        uploads: state.uploads.map((u) =>
          u.fileId === fileId
            ? { ...u, status: 'error', error: error.message || '上传失败' }
            : u
        ),
      }));

      Taro.showToast({
        title: error.message || '上传失败',
        icon: 'none',
      });
    }
  },

  // ============================================
  // Cleanup
  // ============================================

  clearOldMessages: (roomId: number, keepLast = MAX_MESSAGES_PER_ROOM) => {
    set((state) => {
      const messages = state.roomMessages[roomId] || [];
      if (messages.length <= keepLast) return state;

      return {
        roomMessages: {
          ...state.roomMessages,
          [roomId]: messages.slice(-keepLast),
        },
      };
    });
  },

  clearAllMessages: () => {
    set({ roomMessages: {} });
  },

  reset: () => {
    get().disconnectWebSocket();
    set({
      rooms: [],
      currentRoom: null,
      roomMessages: {},
      loadingRooms: false,
      loadingMessages: false,
      uploads: [],
      pendingMessages: [],
    });
  },
}));

// ============================================
// Helper Functions
// ============================================

/**
 * Get current user ID from storage
 */
function getCurrentUserId(): number | null {
  const userInfo = Taro.getStorageSync('userInfo');
  return userInfo?.id || null;
}

/**
 * Handle incoming WebSocket message
 */
function handleWebSocketMessage(message: WSMessage) {
  const state = useChatStore.getState();

  switch (message.type) {
    case 'chat_message':
      handleNewMessage(message.data as any);
      break;

    case 'chat_typing':
      handleTypingIndicator(message.data as any);
      break;

    case 'chat_read':
      handleReadReceipt(message.data as any);
      break;

    case 'ping':
      // Respond with pong
      state.sendWebSocketMessage({
        type: 'pong',
        timestamp: new Date().toISOString(),
      });
      break;

    default:
      console.log('[chatStore] Unknown WebSocket message type:', message.type);
  }
}

/**
 * Handle new chat message from WebSocket
 */
function handleNewMessage(data: any) {
  const { message, groupId } = data;

  if (!message || !groupId) return;

  // Add message to room
  useChatStore.setState((state) => ({
    roomMessages: {
      ...state.roomMessages,
      [groupId]: [...(state.roomMessages[groupId] || []), message],
    },
  }));

  // Update room's last message and unread count
  useChatStore.setState((state) => ({
    rooms: state.rooms.map((room) =>
      room.id === groupId
        ? {
            ...room,
            lastMessage: message,
            unreadCount: room.id === state.currentRoom?.id ? 0 : room.unreadCount + 1,
          }
        : room
    ),
  }));

  // Clear old messages if exceeding limit
  const messages = useChatStore.getState().roomMessages[groupId];
  if (messages && messages.length > MAX_MESSAGES_PER_ROOM) {
    useChatStore.getState().clearOldMessages(groupId);
  }

  // Notify user if not in current room
  const currentState = useChatStore.getState();
  if (currentState.currentRoom?.id !== groupId) {
    Taro.showToast({
      title: '收到新消息',
      icon: 'none',
      duration: 1500,
    });

    // Vibrate
    try {
      Taro.vibrateShort();
    } catch (e) {
      // Ignore vibration errors
    }

    // Play sound (optional)
    // Taro.playBackgroundAudio({ ... });
  }
}

/**
 * Handle typing indicator from WebSocket
 */
function handleTypingIndicator(data: any) {
  const { groupId, userId, userName } = data;

  // Could be used to show "typing..." indicator in UI
  console.log(`[chatStore] User ${userName} is typing in group ${groupId}`);
}

/**
 * Handle read receipt from WebSocket
 */
function handleReadReceipt(data: any) {
  const { groupId, userId, messageId } = data;

  // Update message read status
  useChatStore.setState((state) => ({
    roomMessages: {
      ...state.roomMessages,
      [groupId]: state.roomMessages[groupId]?.map((msg) =>
        msg.id === messageId ? { ...msg, readAt: new Date().toISOString() } : msg
      ) || [],
    },
  }));
}

/**
 * Attempt to reconnect WebSocket
 */
function attemptReconnect() {
  const currentState = useChatStore.getState();

  if (currentState.ws.reconnectAttempts >= currentState.ws.maxReconnectAttempts) {
    console.error('[chatStore] Max reconnection attempts reached');
    Taro.showModal({
      title: '连接已断开',
      content: '聊天服务连接已断开，请刷新页面重试',
      showCancel: false,
    });
    return;
  }

  const delay = Math.min(
    currentState.ws.reconnectDelay * Math.pow(2, currentState.ws.reconnectAttempts),
    MAX_RECONNECT_DELAY
  );

  console.log(`[chatStore] Reconnecting in ${delay}ms (attempt ${currentState.ws.reconnectAttempts + 1})`);

  setTimeout(() => {
    const state = useChatStore.getState();
    useChatStore.setState({
      ws: {
        ...state.ws,
        reconnectAttempts: state.ws.reconnectAttempts + 1,
      },
    });
    useChatStore.getState().connectWebSocket();
  }, delay);
}
