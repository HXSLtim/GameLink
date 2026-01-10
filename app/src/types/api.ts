/**
 * API Response Types
 * 统一的 API 响应类型定义
 */

/**
 * 通用 API 响应结构
 */
export interface ApiResponse<T = any> {
  success: boolean;
  code: number;
  message: string;
  data: T;
  traceId?: string;
}

/**
 * 分页响应
 */
export interface PaginatedResponse<T> {
  items: T[];
  total: number;
  page: number;
  pageSize: number;
}

/**
 * 用户角色
 */
export type UserRole = 'user' | 'player' | 'admin';

/**
 * 账户状态
 */
export type UserStatus = 'active' | 'suspended' | 'banned';

/**
 * VIP 等级
 */
export type VipLevel = 'normal' | 'vip1' | 'vip2' | 'vip3';

/**
 * 登录类型
 */
export type LoginType = 'password' | 'sms' | 'email' | 'oauth';

/**
 * 用户信息
 */
export interface User {
  id: number;
  phone?: string;
  email?: string;
  name: string;
  avatarUrl?: string;
  role: UserRole;
  status: UserStatus;

  // VIP 相关
  vipLevelId?: number;
  vipUnlocked: boolean;
  vipExp: number;
  totalRechargeCents: number;
  vipExpireAt?: string;
  vipLevel?: VipLevelInfo;

  // 钱包
  wallet?: Wallet;
}

/**
 * VIP 等级信息
 */
export interface VipLevelInfo {
  id: number;
  level: string;
  name: string;
  discount: number;
  requiredExp: number;
}

/**
 * 钱包信息
 */
export interface Wallet {
  id: number;
  userId: number;
  balanceCents: number;
  frozenCents: number;
  version: number;
}

/**
 * 陪玩师等级
 */
export type PlayerRank = 'Bronze' | 'Silver' | 'Gold' | 'Platinum' | 'Diamond' | 'Master';

/**
 * 认证状态
 */
export type VerificationStatus = 'pending' | 'verified' | 'rejected';

/**
 * 在线状态
 */
export type OnlineStatus = 'offline' | 'online' | 'busy';

/**
 * 陪玩师信息
 */
export interface Player {
  id: number;
  userId: number;
  nickname: string;
  bio?: string;
  rank: PlayerRank;
  verificationStatus: VerificationStatus;
  ratingAverage: number;
  ratingCount: number;
  hourlyRateCents: number;
  orderCount: number;
  completionRate: number;

  // 用户信息（头像等）
  user?: {
    avatarUrl?: string;
  };

  // 在线状态
  onlineStatus?: OnlineStatus;
}

/**
 * 游戏信息
 */
export interface Game {
  id: number;
  name: string;
  iconUrl?: string;
  description?: string;
  isActive: boolean;
  category?: string;
}

/**
 * 服务项目
 */
export interface ServiceItem {
  id: number;
  gameId: number;
  name: string;
  description?: string;
  duration: number;
  priceCents: number;
  isActive: boolean;

  game?: Game;
}

/**
 * 订单类型
 */
export type OrderType = 'solo' | 'team' | 'gift';

/**
 * 订单状态
 */
export type OrderStatus = 'pending' | 'waiting' | 'in_progress' | 'pending_confirmation' | 'completed' | 'canceled' | 'disputed';

/**
 * 订单信息
 */
export interface Order {
  id: number;
  orderNo: string;
  userId: number;
  playerIds: number[];
  gameId: number;
  serviceItemId?: number;
  type: OrderType;
  status: OrderStatus;

  // 价格信息
  originalAmountCents: number;
  discountCents: number;
  finalAmountCents: number;

  // 时间信息
  duration: number;
  scheduledStartAt?: string;
  startedAt?: string;
  completedAt?: string;

  // 评价
  rating?: number;
  review?: string;

  // 关联数据
  user?: User;
  players?: Player[];
  game?: Game;
  serviceItem?: ServiceItem;
}

/**
 * 收藏的陪玩师
 */
export interface FavoritePlayer {
  id: number;
  playerId: number;
  nickname: string;
  avatarUrl?: string;
  bio?: string;
  rank: PlayerRank;
  ratingAverage: number;
  hourlyRateCents: number;
  createdAt: string;
}

/**
 * 搜索结果类型
 */
export type SearchResultType = 'player' | 'game';

/**
 * 搜索结果项
 */
export interface SearchResultItem {
  id: number;
  type: SearchResultType;
  name: string;
  description?: string;
  imageUrl?: string;
  extra?: Record<string, any>;
}

/**
 * 搜索响应
 */
export interface SearchResponse {
  players: SearchResultItem[];
  games: SearchResultItem[];
  total: number;
}

/**
 * 登录请求
 */
export interface LoginRequest {
  type: LoginType;
  phone?: string;
  email?: string;
  password?: string;
  code?: string;
}

/**
 * 登录响应
 */
export interface LoginResponse {
  token: string;
  refreshToken?: string;
  user: User;
}

/**
 * 微信登录请求
 */
export interface WeChatLoginRequest {
  code: string;
  encryptedData?: string;
  iv?: string;
}

/**
 * 角色切换请求
 */
export interface SwitchRoleRequest {
  role: UserRole;
}

/**
 * 角色切换响应
 */
export interface SwitchRoleResponse {
  token: string;
  currentRole: UserRole;
}

/**
 * 更新在线状态请求
 */
export interface UpdateOnlineStatusRequest {
  status: OnlineStatus;
  online?: boolean; // 兼容旧版本
}

/**
 * 在线状态响应
 */
export interface OnlineStatusResponse {
  status: OnlineStatus;
  online: boolean;
}

/**
 * 优惠券类型
 */
export type CouponType = 'deduct' | 'discount';

/**
 * 优惠券范围
 */
export type CouponScope = 'all' | 'game' | 'service';

/**
 * 优惠券状态
 */
export type CouponStatus = 'available' | 'used' | 'expired';

/**
 * 用户优惠券
 */
export interface UserCoupon {
  id: number;
  userId: number;
  couponId: number;
  status: CouponStatus;
  expiresAt: string;

  coupon: {
    id: number;
    name: string;
    type: CouponType;
    value: number;
    minAmountCents: number;
    scope: CouponScope;
    scopeIds?: number[];
  };
}

/**
 * 充值档位
 */
export interface RechargeTier {
  id: number;
  amountCents: number;
  bonusCents: number;
  originalCents: number;
  discountPercent: number;
  isActive: boolean;
}

/**
 * 充值记录
 */
export interface RechargeRecord {
  id: number;
  userId: number;
  amountCents: number;
  bonusCents: number;
  totalCents: number;
  method: string;
  status: string;
  createdAt: string;
}
