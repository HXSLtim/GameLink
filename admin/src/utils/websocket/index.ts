/**
 * WebSocket Utilities
 * 导出所有 WebSocket 相关工具
 */

// Manager
export { WebSocketManager, wsManager, initWebSocket, disconnectWebSocket, getWebSocketManager } from './manager';

// Types
export type {
  WSMessage,
  TMessageType,
  SystemStatus,
  OnlineUsers,
  OrderQueue,
  Alert,
  PresenceUpdate,
  RoomEvent,
  RoomMemberEvent,
  LFGEvent,
  WebSocketConfig,
  MessageHandler,
  EventListeners,
} from './types';

export { MessageType, ConnectionStatus } from './types';

// Message Handler
export { registerChatMessageHandlers, unregisterChatMessageHandlers } from './messageHandler';
