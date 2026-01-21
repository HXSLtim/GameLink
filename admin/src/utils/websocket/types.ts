/**
 * WebSocket Types
 * 与后端 api/internal/ws/message.go 对应
 */

/**
 * WebSocket 消息结构
 */
export interface WSMessage<T = unknown> {
  type: string;
  timestamp: string;
  data?: T;
}

/**
 * 消息类型常量 (使用 const object 代替 enum 以兼容 erasableSyntaxOnly)
 */
export const MessageType = {
  // System status
  SystemStatus: 'system_status',
  OnlineUsers: 'online_users',
  OrderQueue: 'order_queue',
  Alert: 'alert',

  // Client messages
  Ping: 'ping',
  Pong: 'pong',
  Subscribe: 'subscribe',

  // Presence (Discord/Kook style)
  PresenceUpdate: 'presence_update',
  PresenceSubscribe: 'presence_subscribe',
  PresenceBatch: 'presence_batch',

  // Room events
  RoomCreated: 'room_created',
  RoomUpdated: 'room_updated',
  RoomClosed: 'room_closed',
  RoomMemberJoined: 'room_member_joined',
  RoomMemberLeft: 'room_member_left',
  RoomMemberReady: 'room_member_ready',
  RoomStarted: 'room_started',
  RoomFinished: 'room_finished',

  // LFG events
  LFGNew: 'lfg_new',
  LFGMatched: 'lfg_matched',
  LFGExpired: 'lfg_expired',
  LFGCanceled: 'lfg_canceled',

  // Voice events
  VoiceStarted: 'voice_started',
  VoiceStopped: 'voice_stopped',
  VoiceMemberJoined: 'voice_member_joined',
  VoiceMemberLeft: 'voice_member_left',
  VoiceMemberMuted: 'voice_member_muted',
} as const;

/**
 * MessageType 类型
 */
export type TMessageType = typeof MessageType[keyof typeof MessageType];

/**
 * 系统状态数据
 */
export interface SystemStatus {
  cpuUsage: number;
  memoryUsage: number;
  memoryTotal: number;
  memoryUsed: number;
  goroutines: number;
  dbConnections: {
    active: number;
    idle: number;
    max: number;
  };
  uptime: number;
  requestsPerSec: number;
  status: 'healthy' | 'degraded' | 'critical';
}

/**
 * 在线用户数据
 */
export interface OnlineUsers {
  total: number;
  peak: number;
  byRole: Record<string, number>;
  updatedAt: string;
}

/**
 * 订单队列数据
 */
export interface OrderQueue {
  pending: number;
  processing: number;
  completed: number;
  processingSpeed: number;
  averageWaitTime: number;
  hasBacklog: boolean;
}

/**
 * 警告数据
 */
export interface Alert {
  id: string;
  level: 'high' | 'medium' | 'low';
  type: 'system' | 'business' | 'security';
  title: string;
  message: string;
  source: string;
  createdAt: string;
  isRead: boolean;
}

/**
 * 在线状态更新
 */
export interface PresenceUpdate {
  playerId: number;
  status: string;
  currentGameId?: number;
  currentGameName?: string;
  customStatus?: string;
  currentRoomId?: number;
  updatedAt: string;
}

/**
 * 房间事件
 */
export interface RoomEvent {
  roomId: number;
  roomName: string;
  roomType: string;
  gameId: number;
  gameName?: string;
  hostUserId: number;
  status: string;
  currentMembers: number;
  maxMembers: number;
}

/**
 * 房间成员事件
 */
export interface RoomMemberEvent {
  roomId: number;
  userId: number;
  nickname: string;
  avatar?: string;
  role: string;
  isReady: boolean;
}

/**
 * LFG 事件
 */
export interface LFGEvent {
  requestId: number;
  userId: number;
  userNickname?: string;
  gameId: number;
  gameName?: string;
  requestType: string;
  title: string;
  description?: string;
  requiredPlayers: number;
  status: string;
  matchedRoomId?: number;
}

/**
 * WebSocket 连接状态
 */
export const ConnectionStatus = {
  Connecting: 'connecting',
  Connected: 'connected',
  Disconnected: 'disconnected',
  Reconnecting: 'reconnecting',
  Error: 'error',
} as const;

/**
 * ConnectionStatus 类型
 */
export type TConnectionStatus = typeof ConnectionStatus[keyof typeof ConnectionStatus];

/**
 * WebSocket 配置
 */
export interface WebSocketConfig {
  url: string;
  token: string | null;
  heartbeatInterval?: number;
  reconnect?: boolean;
  reconnectInterval?: number;
  maxReconnectAttempts?: number;
  debug?: boolean;
}

/**
 * 消息处理器类型
 */
export type MessageHandler = <T = unknown>(message: WSMessage<T>) => void;

/**
 * 事件监听器映射
 */
export interface EventListeners {
  onOpen?: (event: Event) => void;
  onMessage?: (message: WSMessage) => void;
  onError?: (error: Event) => void;
  onClose?: (event: CloseEvent) => void;
  onReconnecting?: (attempt: number) => void;
  onReconnected?: () => void;
}
