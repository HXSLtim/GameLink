// ============================================
// Common Types
// ============================================

export interface UserInfo {
  id: number;
  name: string;
  email?: string;
  phone?: string;
  avatar?: string;
  role: string;
  permissions: string[];
  createdAt: string;
  updatedAt: string;
}

export interface LoginRequest {
  username: string;
  password: string;
}

export interface LoginResponse {
  token: string;
  user: UserInfo;
}

export interface MenuItem {
  id: number;
  key: string;
  label: string;
  icon?: string;
  path?: string;
  parentId?: number | null;
  sort: number;
  permission?: string;
  children?: MenuItem[];
}

// User types
export interface User {
  id: number;
  name: string;
  email: string;
  phone: string;
  status: 'active' | 'inactive' | 'blocked';
  role: string;
  createdAt: string;
}

// Order types
export interface Order {
  id: number;
  orderNo: string;
  userId: number;
  playerIds: number[];
  status: 'pending' | 'in_progress' | 'completed' | 'cancelled' | 'refunded';
  amount: number;
  duration: number;
  createdAt: string;
}

// Player types
export interface Player {
  id: number;
  userId: number;
  nickname: string;
  avatar: string;
  rank: number;
  level: number;
  pricePerHour: number;
  status: 'available' | 'busy' | 'offline';
  tags: string[];
}

// Chat types
export interface ChatMessage {
  id: number;
  orderId: number;
  senderId: number;
  senderType: 'user' | 'player' | 'admin';
  content: string;
  type: 'text' | 'image' | 'voice';
  createdAt: string;
  read: boolean;
}

export interface ChatRoom {
  id: number;
  orderId: number;
  participants: number[];
  lastMessage?: ChatMessage;
  unreadCount: number;
}
