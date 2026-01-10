/**
 * Player API Service
 * 陪玩师相关 API
 */

import apiClient from './request';
import type {
  ApiResponse,
  PaginatedResponse,
  Player,
  Order,
  Wallet,
  OnlineStatusResponse,
  UpdateOnlineStatusRequest,
} from '@/types/api';

/**
 * 交易记录
 */
interface PlayerTransaction {
  id: number;
  playerId: number;
  type: 'earning' | 'withdraw' | 'commission' | 'bonus';
  amountCents: number;
  balanceAfterCents: number;
  description: string;
  orderId?: number;
  createdAt: string;
}

/**
 * 提现记录
 */
interface WithdrawRecord {
  id: number;
  playerId: number;
  amountCents: number;
  feeCents: number;
  actualAmountCents: number;
  method: string;
  status: 'pending' | 'processing' | 'completed' | 'rejected';
  accountInfo: Record<string, string>;
  createdAt: string;
  completedAt?: string;
}

/**
 * 评价记录
 */
interface PlayerReview {
  id: number;
  orderId: number;
  userId: number;
  playerId: number;
  rating: number;
  content?: string;
  createdAt: string;
  user?: {
    name: string;
    avatarUrl?: string;
  };
}

/**
 * 收入明细
 */
interface IncomeDetail {
  id: number;
  orderId: number;
  orderNo: string;
  grossAmountCents: number;
  commissionCents: number;
  netAmountCents: number;
  status: 'frozen' | 'available' | 'withdrawn';
  createdAt: string;
  availableAt?: string;
}

/**
 * 收入统计响应
 */
interface EarningsResponse {
  totalCents: number;
  frozenCents: number;
  availableCents: number;
  orderCount: number;
  completionRate: number;
  averageRating: number;
}

/**
 * 提现响应
 */
interface WithdrawResponse {
  id: number;
  amountCents: number;
  feeCents: number;
  actualAmountCents: number;
  status: string;
}

/**
 * 获取当前陪玩师资料
 */
export const getMyPlayerProfile = (): Promise<ApiResponse<Player>> => {
  return apiClient.get<Player>('/player/me');
};

/**
 * 更新陪玩师资料
 */
export const updatePlayerProfile = (data: {
  nickname?: string;
  bio?: string;
  hourlyRateCents?: number;
}): Promise<ApiResponse<Player>> => {
  return apiClient.put<Player>('/player/me', data);
};

/**
 * 获取在线状态
 */
export const getOnlineStatus = (): Promise<ApiResponse<OnlineStatusResponse>> => {
  return apiClient.get<OnlineStatusResponse>('/player/online-status');
};

/**
 * 更新在线状态
 */
export const updateOnlineStatus = (data: UpdateOnlineStatusRequest): Promise<ApiResponse<OnlineStatusResponse>> => {
  return apiClient.put<OnlineStatusResponse>('/player/online-status', data);
};

/**
 * 获取我的订单列表（陪玩师端）
 */
export const getMyOrders = (params?: {
  page?: number;
  pageSize?: number;
  status?: string;
}): Promise<ApiResponse<PaginatedResponse<Order>>> => {
  return apiClient.get<PaginatedResponse<Order>>('/player/orders', params);
};

/**
 * 接单
 */
export const acceptOrder = (orderId: number): Promise<ApiResponse<Order>> => {
  return apiClient.post<Order>(`/player/orders/${orderId}/accept`);
};

/**
 * 拒单
 */
export const rejectOrder = (orderId: number, reason?: string): Promise<ApiResponse<void>> => {
  return apiClient.post<void>(`/player/orders/${orderId}/reject`, { reason });
};

/**
 * 开始服务
 */
export const startOrder = (orderId: number): Promise<ApiResponse<Order>> => {
  return apiClient.post<Order>(`/player/orders/${orderId}/start`);
};

/**
 * 完成服务
 */
export const completeOrder = (orderId: number): Promise<ApiResponse<Order>> => {
  return apiClient.post<Order>(`/player/orders/${orderId}/complete`);
};

/**
 * 获取收入统计
 */
export const getEarnings = (params?: {
  period?: 'today' | 'week' | 'month' | 'all';
}): Promise<ApiResponse<EarningsResponse>> => {
  return apiClient.get<EarningsResponse>('/player/earnings', params);
};

/**
 * 获取钱包信息
 */
export const getWallet = (): Promise<ApiResponse<Wallet>> => {
  return apiClient.get<Wallet>('/player/me/wallet');
};

/**
 * 获取交易记录
 */
export const getTransactions = (params?: {
  page?: number;
  pageSize?: number;
}): Promise<ApiResponse<PaginatedResponse<PlayerTransaction>>> => {
  return apiClient.get<PaginatedResponse<PlayerTransaction>>('/player/me/transactions', params);
};

/**
 * 提现申请
 */
export const withdraw = (data: {
  amountCents: number;
  method: string;
  accountInfo: Record<string, string>;
}): Promise<ApiResponse<WithdrawResponse>> => {
  return apiClient.post<WithdrawResponse>('/player/me/withdraw', data);
};

/**
 * 获取提现记录
 */
export const getWithdrawRecords = (params?: {
  page?: number;
  pageSize?: number;
}): Promise<ApiResponse<PaginatedResponse<WithdrawRecord>>> => {
  return apiClient.get<PaginatedResponse<WithdrawRecord>>('/player/me/withdraw-records', params);
};

/**
 * 获取评价列表
 */
export const getReviews = (params?: {
  page?: number;
  pageSize?: number;
}): Promise<ApiResponse<PaginatedResponse<PlayerReview>>> => {
  return apiClient.get<PaginatedResponse<PlayerReview>>('/player/me/reviews', params);
};

/**
 * 获取排行榜
 */
export const getRanking = (params?: {
  type?: 'daily' | 'weekly' | 'monthly';
  gameId?: number;
}): Promise<ApiResponse<Player[]>> => {
  return apiClient.get<Player[]>('/public/ranking', params);
};

/**
 * 获取收入明细
 */
export const getIncomeDetails = (params?: {
  page?: number;
  pageSize?: number;
  startDate?: string;
  endDate?: string;
}): Promise<ApiResponse<PaginatedResponse<IncomeDetail>>> => {
  return apiClient.get<PaginatedResponse<IncomeDetail>>('/player/me/income-details', params);
};
