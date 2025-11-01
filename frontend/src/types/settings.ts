/**
 * 系统设置类型定义
 */

/**
 * 平台基础设置
 */
export interface PlatformSettings {
  siteName: string;
  siteUrl: string;
  contactEmail: string;
  contactPhone?: string;
  icp?: string;
  description?: string;
  keywords?: string;
  logo?: string;
  favicon?: string;
}

/**
 * 支付设置
 */
export interface PaymentSettings {
  alipayEnabled: boolean;
  alipayAppId?: string;
  alipayPrivateKey?: string;
  wechatEnabled: boolean;
  wechatAppId?: string;
  wechatMchId?: string;
  wechatApiKey?: string;
  balanceEnabled: boolean;
  minRechargeAmount: number;
  maxRechargeAmount: number;
  serviceFeeRate: number; // 平台服务费率（百分比）
}

/**
 * 订单设置
 */
export interface OrderSettings {
  autoConfirmMinutes: number; // 自动确认时间（分钟）
  autoCancelMinutes: number; // 自动取消时间（分钟）
  refundEnabled: boolean;
  refundDays: number; // 退款期限（天）
  minOrderAmount: number; // 最小订单金额（分）
  maxOrderAmount: number; // 最大订单金额（分）
}

/**
 * 用户设置
 */
export interface UserSettings {
  registrationEnabled: boolean;
  emailVerificationRequired: boolean;
  phoneVerificationRequired: boolean;
  minPasswordLength: number;
  defaultUserRole: string;
  sessionTimeout: number; // 会话超时时间（分钟）
}

/**
 * 陪玩师设置
 */
export interface PlayerSettings {
  verificationRequired: boolean;
  minAge: number;
  minHourlyRate: number; // 最低时薪（分）
  maxHourlyRate: number; // 最高时薪（分）
  commissionRate: number; // 平台抽成比例（百分比）
}

/**
 * 系统通知设置
 */
export interface NotificationSettings {
  emailEnabled: boolean;
  smsEnabled: boolean;
  pushEnabled: boolean;
  orderNotifyUser: boolean;
  orderNotifyPlayer: boolean;
  paymentNotifyUser: boolean;
  reviewNotifyPlayer: boolean;
}

/**
 * 公告信息
 */
export interface Announcement {
  id: number;
  title: string;
  content: string;
  type: 'info' | 'warning' | 'error' | 'success';
  startDate?: string;
  endDate?: string;
  enabled: boolean;
  createdAt: string;
  updatedAt: string;
}

/**
 * 创建公告请求
 */
export interface CreateAnnouncementRequest {
  title: string;
  content: string;
  type: 'info' | 'warning' | 'error' | 'success';
  startDate?: string;
  endDate?: string;
  enabled: boolean;
}

/**
 * 更新公告请求
 */
export interface UpdateAnnouncementRequest {
  title?: string;
  content?: string;
  type?: 'info' | 'warning' | 'error' | 'success';
  startDate?: string;
  endDate?: string;
  enabled?: boolean;
}

/**
 * 系统设置（所有设置的集合）
 */
export interface SystemSettings {
  platform: PlatformSettings;
  payment: PaymentSettings;
  order: OrderSettings;
  user: UserSettings;
  player: PlayerSettings;
  notification: NotificationSettings;
}


