/**
 * Chat Message Handlers
 * 处理聊天相关的 WebSocket 消息
 */
import { logger } from '@/utils/logger';
import { useChatStore } from '@/stores/modules/chatStore';
import type { WSMessage } from './types';
import type { TMessageType } from './types';
import { getWebSocketManager } from './manager';

/**
 * 聊天相关消息类型 (参考用)
 */
const _CHAT_MESSAGE_TYPES: TMessageType[] = [
  'presence_update' as TMessageType,
  'presence_batch' as TMessageType,
  'room_created' as TMessageType,
  'room_updated' as TMessageType,
  'room_closed' as TMessageType,
  'room_member_joined' as TMessageType,
  'room_member_left' as TMessageType,
];

/**
 * 取消订阅函数列表
 */
const unsubscribers: Array<() => void> = [];

/**
 * 处理在线状态更新
 */
function handlePresenceUpdate(message: WSMessage): void {
  try {
    const data = message.data as {
      playerId: number;
      status: string;
      currentGameId?: number;
      currentGameName?: string;
      currentRoomId?: number;
      customStatus?: string;
    };

    // 在聊天列表中更新用户状态
    // 这里可以根据需要在 chatStore 中添加对应的状态更新逻辑
    logger.debug('Presence update:', data);

    // 如果用户在某个聊天室，可以更新聊天室成员状态
    if (data.currentRoomId) {
      const chatStore = useChatStore.getState();
      if (chatStore.currentRoomId === data.currentRoomId) {
        // 触发聊天室更新
        // chatStore.updateMemberPresence(data.playerId, data.status);
      }
    }
  } catch (error) {
    logger.error('Failed to handle presence update:', error);
  }
}

/**
 * 处理房间更新
 */
function handleRoomUpdated(message: WSMessage): void {
  try {
    const data = message.data as {
      roomId: number;
      roomName?: string;
      status?: string;
      currentMembers?: number;
    };

    const chatStore = useChatStore.getState();

    // 更新聊天室信息
    if (data.roomId) {
      const updates: Record<string, unknown> = {};
      if (data.roomName) updates.name = data.roomName;
      if (data.status) updates.status = data.status;
      if (data.currentMembers !== undefined) updates.currentMembers = data.currentMembers;

      chatStore.updateRoom(data.roomId, updates);
      logger.debug('Room updated:', data);
    }
  } catch (error) {
    logger.error('Failed to handle room update:', error);
  }
}

/**
 * 处理房间关闭
 */
function handleRoomClosed(message: WSMessage): void {
  try {
    const data = message.data as { roomId: number };

    const chatStore = useChatStore.getState();

    // 更新聊天室状态为已关闭
    chatStore.updateRoom(data.roomId, { status: 'closed' });

    // 如果当前正在查看这个房间，显示提示
    if (chatStore.currentRoomId === data.roomId) {
      logger.info('Current room has been closed:', data.roomId);
    }
  } catch (error) {
    logger.error('Failed to handle room close:', error);
  }
}

/**
 * 处理房间成员加入
 */
function handleRoomMemberJoined(message: WSMessage): void {
  try {
    const data = message.data as {
      roomId: number;
      userId: number;
      nickname: string;
      avatar?: string;
      role: string;
    };

    const chatStore = useChatStore.getState();

    // 如果当前正在查看这个房间，可以添加系统消息
    if (chatStore.currentRoomId === data.roomId) {
      // 可以添加一条系统消息："XXX 加入了聊天室"
      logger.info(`User joined room: ${data.nickname}, roomId: ${data.roomId}`);
    }
  } catch (error) {
    logger.error('Failed to handle room member join:', error);
  }
}

/**
 * 处理房间成员离开
 */
function handleRoomMemberLeft(message: WSMessage): void {
  try {
    const data = message.data as {
      roomId: number;
      userId: number;
      nickname: string;
    };

    const chatStore = useChatStore.getState();

    // 如果当前正在查看这个房间，可以添加系统消息
    if (chatStore.currentRoomId === data.roomId) {
      logger.info(`User left room: ${data.nickname}, roomId: ${data.roomId}`);
    }
  } catch (error) {
    logger.error('Failed to handle room member leave:', error);
  }
}

/**
 * 注册所有聊天消息处理器
 */
export function registerChatMessageHandlers(): void {
  const wsManager = getWebSocketManager();

  // 清除之前的订阅
  unregisterChatMessageHandlers();

  // 注册各个消息类型的处理器
  unsubscribers.push(
    wsManager.on('presence_update' as TMessageType, handlePresenceUpdate),
    wsManager.on('room_updated' as TMessageType, handleRoomUpdated),
    wsManager.on('room_closed' as TMessageType, handleRoomClosed),
    wsManager.on('room_member_joined' as TMessageType, handleRoomMemberJoined),
    wsManager.on('room_member_left' as TMessageType, handleRoomMemberLeft),
  );

  logger.info('Chat message handlers registered');
}

/**
 * 取消所有聊天消息处理器
 */
export function unregisterChatMessageHandlers(): void {
  unsubscribers.forEach((unsubscribe) => {
    try {
      unsubscribe();
    } catch (error) {
      logger.error('Failed to unsubscribe message handler:', error);
    }
  });

  // 清空数组
  unsubscribers.length = 0;

  logger.debug('Chat message handlers unregistered');
}
