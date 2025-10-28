/**
 * 订单状态枚举
 */
export enum OrderStatus {
  /** 待支付 */
  PENDING_PAYMENT = 'pending_payment',
  /** 待接单 */
  PENDING_ACCEPT = 'pending_accept',
  /** 待审核 */
  PENDING_REVIEW = 'pending_review',
  /** 审核中 */
  IN_REVIEW = 'in_review',
  /** 审核通过 */
  REVIEW_APPROVED = 'review_approved',
  /** 审核拒绝 */
  REVIEW_REJECTED = 'review_rejected',
  /** 进行中 */
  IN_PROGRESS = 'in_progress',
  /** 已完成 */
  COMPLETED = 'completed',
  /** 已取消 */
  CANCELLED = 'cancelled',
  /** 已退款 */
  REFUNDED = 'refunded',
}

/**
 * 订单审核状态
 */
export enum ReviewStatus {
  /** 待审核 */
  PENDING = 'pending',
  /** 审核中 */
  IN_REVIEW = 'in_review',
  /** 已通过 */
  APPROVED = 'approved',
  /** 已拒绝 */
  REJECTED = 'rejected',
}

/**
 * 游戏类型
 */
export enum GameType {
  /** 王者荣耀 */
  HONOR_OF_KINGS = 'honor_of_kings',
  /** 英雄联盟 */
  LEAGUE_OF_LEGENDS = 'league_of_legends',
  /** 和平精英 */
  PEACEKEEPER_ELITE = 'peacekeeper_elite',
  /** 原神 */
  GENSHIN_IMPACT = 'genshin_impact',
  /** 其他 */
  OTHER = 'other',
}

/**
 * 服务类型
 */
export enum ServiceType {
  /** 陪玩 */
  ACCOMPANY = 'accompany',
  /** 代练 */
  BOOST = 'boost',
  /** 上分 */
  RANK_UP = 'rank_up',
  /** 娱乐 */
  ENTERTAINMENT = 'entertainment',
}

/**
 * 订单操作类型
 */
export enum OrderActionType {
  /** 创建订单 */
  CREATE = 'create',
  /** 支付订单 */
  PAY = 'pay',
  /** 接单 */
  ACCEPT = 'accept',
  /** 开始服务 */
  START = 'start',
  /** 提交审核 */
  SUBMIT_REVIEW = 'submit_review',
  /** 审核通过 */
  APPROVE = 'approve',
  /** 审核拒绝 */
  REJECT = 'reject',
  /** 完成订单 */
  COMPLETE = 'complete',
  /** 取消订单 */
  CANCEL = 'cancel',
  /** 申请退款 */
  REQUEST_REFUND = 'request_refund',
  /** 退款 */
  REFUND = 'refund',
}

/**
 * 订单操作日志
 */
export interface OrderLog {
  id: string;
  orderId: string;
  action: OrderActionType;
  operator: string;
  operatorId: string;
  operatorRole: 'user' | 'player' | 'admin' | 'system';
  content: string;
  statusBefore?: OrderStatus;
  statusAfter?: OrderStatus;
  metadata?: Record<string, any>;
  createdAt: string;
}

/**
 * 订单审核记录
 */
export interface OrderReview {
  id: string;
  orderId: string;
  reviewer: string;
  reviewerId: string;
  status: ReviewStatus;
  result?: 'approved' | 'rejected';
  reason?: string;
  screenshots?: string[];
  createdAt: string;
  updatedAt: string;
}

/**
 * 陪玩者信息
 */
export interface PlayerInfo {
  id: string;
  username: string;
  avatar?: string;
  rating: number;
  completedOrders: number;
  level: number;
  tags: string[];
}

/**
 * 用户信息
 */
export interface UserInfo {
  id: string;
  username: string;
  avatar?: string;
  phone?: string;
}

/**
 * 订单数据结构
 */
export interface Order {
  id: string;
  orderNo: string;
  user: UserInfo;
  player?: PlayerInfo;
  gameType: GameType;
  serviceType: ServiceType;
  status: OrderStatus;
  reviewStatus?: ReviewStatus;
  price: number;
  duration: number;
  description: string;
  requirements?: string;
  screenshots?: string[];
  reviewNote?: string;
  createdAt: string;
  updatedAt: string;
  paidAt?: string;
  acceptedAt?: string;
  startedAt?: string;
  completedAt?: string;
  cancelledAt?: string;
}

/**
 * 订单详情（包含操作历史）
 */
export interface OrderDetail extends Order {
  logs: OrderLog[];
  reviews: OrderReview[];
}

/**
 * 订单查询参数
 */
export interface OrderQueryParams {
  page?: number;
  pageSize?: number;
  keyword?: string;
  status?: OrderStatus;
  reviewStatus?: ReviewStatus;
  gameType?: GameType;
  serviceType?: ServiceType;
  startDate?: string;
  endDate?: string;
  userId?: string;
  playerId?: string;
}

/**
 * 订单统计数据
 */
export interface OrderStatistics {
  total: number;
  pendingPayment: number;
  pendingReview: number;
  inProgress: number;
  completed: number;
  cancelled: number;
  todayOrders: number;
  todayRevenue: number;
}
