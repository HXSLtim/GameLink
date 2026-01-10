/**
 * User API Service
 * 用户相关 API
 */

import Taro from '@tarojs/taro';
import apiClient from './request';
import type {
  ApiResponse,
  PaginatedResponse,
  User,
  Wallet,
  Order,
  FavoritePlayer,
  UserCoupon,
  RechargeTier,
  RechargeRecord,
  VipLevelInfo,
} from '@/types/api';

/**
 * 交易记录
 */
interface Transaction {
  id: number;
  userId: number;
  type: 'recharge' | 'payment' | 'refund' | 'withdraw' | 'bonus';
  amountCents: number;
  balanceAfterCents: number;
  description: string;
  createdAt: string;
}

/**
 * 收藏响应
 */
interface AddFavoriteResponse {
  id: number;
  playerId: number;
}

/**
 * VIP 信息响应
 */
interface VipInfoResponse {
  vipLevel: VipLevelInfo | null;
  currentExp: number;
  nextLevelExp: number;
  discount: number;
}

/**
 * 获取当前用户信息
 */
export const getCurrentUser = (): Promise<ApiResponse<User>> => {
  return apiClient.get<User>('/user/me');
};

/**
 * 更新用户信息
 */
export const updateProfile = (data: {
  name?: string;
  avatarUrl?: string;
}): Promise<ApiResponse<User>> => {
  return apiClient.put<User>('/user/me', data);
};

/**
 * 获取钱包信息
 */
export const getWallet = (): Promise<ApiResponse<Wallet>> => {
  return apiClient.get<Wallet>('/user/me/wallet');
};

/**
 * 获取交易记录
 */
export const getTransactions = (params?: {
  page?: number;
  pageSize?: number;
}): Promise<ApiResponse<PaginatedResponse<Transaction>>> => {
  return apiClient.get<PaginatedResponse<Transaction>>('/user/me/transactions', params);
};

/**
 * 获取我的订单列表
 */
export const getMyOrders = (params?: {
  page?: number;
  pageSize?: number;
  status?: string;
}): Promise<ApiResponse<PaginatedResponse<Order>>> => {
  return apiClient.get<PaginatedResponse<Order>>('/user/orders', params);
};

/**
 * 获取订单详情
 */
export const getOrderDetail = (orderId: number): Promise<ApiResponse<Order>> => {
  return apiClient.get<Order>(`/user/orders/${orderId}`);
};

/**
 * 创建订单
 */
export const createOrder = (data: {
  gameId: number;
  playerIds: number[];
  serviceItemId?: number;
  type: 'solo' | 'team';
  duration: number;
  scheduledStartAt?: string;
  couponId?: number;
}): Promise<ApiResponse<Order>> => {
  return apiClient.post<Order>('/user/orders', data);
};

/**
 * 取消订单
 */
export const cancelOrder = (orderId: number): Promise<ApiResponse<void>> => {
  return apiClient.post<void>(`/user/orders/${orderId}/cancel`);
};

/**
 * 确认订单完成
 */
export const confirmOrder = (orderId: number, data: {
  rating: number;
  review?: string;
}): Promise<ApiResponse<void>> => {
  return apiClient.post<void>(`/user/orders/${orderId}/confirm`, data);
};

/**
 * 获取收藏列表
 */
export const getFavorites = (params?: {
  page?: number;
  pageSize?: number;
}): Promise<ApiResponse<PaginatedResponse<FavoritePlayer>>> => {
  return apiClient.get<PaginatedResponse<FavoritePlayer>>('/user/favorites/players', params);
};

/**
 * 添加收藏
 */
export const addFavorite = (playerId: number): Promise<ApiResponse<AddFavoriteResponse>> => {
  return apiClient.post<AddFavoriteResponse>(`/user/favorites/players/${playerId}`);
};

/**
 * 取消收藏
 */
export const removeFavorite = (playerId: number): Promise<ApiResponse<void>> => {
  return apiClient.delete<void>(`/user/favorites/players/${playerId}`);
};

/**
 * 检查是否已收藏
 */
export const checkFavorite = (playerId: number): Promise<ApiResponse<{ isFavorite: boolean; playerId: number }>> => {
  return apiClient.get<{ isFavorite: boolean; playerId: number }>(`/user/favorites/players/${playerId}/check`);
};

/**
 * 获取我的优惠券
 */
export const getCoupons = (params?: {
  page?: number;
  pageSize?: number;
  status?: string;
}): Promise<ApiResponse<PaginatedResponse<UserCoupon>>> => {
  return apiClient.get<PaginatedResponse<UserCoupon>>('/user/coupons', params);
};

/**
 * 领取优惠券
 */
export const claimCoupon = (couponId: number): Promise<ApiResponse<void>> => {
  return apiClient.post<void>(`/user/coupons/${couponId}/claim`);
};

/**
 * 获取充值档位
 */
export const getRechargeTiers = (): Promise<ApiResponse<RechargeTier[]>> => {
  return apiClient.get<RechargeTier[]>('/user/recharge/tiers');
};

/**
 * 创建充值订单
 */
export const createRecharge = (tierId: number): Promise<ApiResponse<RechargeRecord>> => {
  return apiClient.post<RechargeRecord>('/user/recharge/create', { tierId });
};

/**
 * 获取充值记录
 */
export const getRechargeRecords = (params?: {
  page?: number;
  pageSize?: number;
}): Promise<ApiResponse<PaginatedResponse<RechargeRecord>>> => {
  return apiClient.get<PaginatedResponse<RechargeRecord>>('/user/recharge/records', params);
};

/**
 * 获取 VIP 信息
 */
export const getVipInfo = (): Promise<ApiResponse<VipInfoResponse>> => {
  return apiClient.get<VipInfoResponse>('/user/vip/info');
};

/**
 * 领取月度优惠券
 */
export const claimMonthlyCoupon = (): Promise<ApiResponse<void>> => {
  return apiClient.post<void>('/user/vip/monthly-coupon');
};

/**
 * 上传图片
 */
export const uploadImage = (filePath: string): Promise<ApiResponse<{ url: string }>> => {
  return new Promise((resolve, reject) => {
    Taro.uploadFile({
      url: `${process.env.TARO_APP_API_BASE_URL}/user/upload`,
      filePath,
      name: 'file',
      header: {
        'Authorization': `Bearer ${Taro.getStorageSync('token')}`,
      },
      success: (res) => {
        const data = JSON.parse(res.data) as ApiResponse<{ url: string }>;
        if (data.success) {
          resolve(data);
        } else {
          reject(new Error(data.message));
        }
      },
      fail: reject,
    });
  });
};
