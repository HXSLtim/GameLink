/**
 * 通用下拉选项配置
 */

export interface SelectOption {
  label: string;
  value: string;
}

// ==================== 订单状态选项 ====================
export const ORDER_STATUS_OPTIONS: SelectOption[] = [
  { label: '全部状态', value: '' },
  { label: '待处理', value: 'pending' },
  { label: '已确认', value: 'confirmed' },
  { label: '进行中', value: 'in_progress' },
  { label: '已完成', value: 'completed' },
  { label: '已取消', value: 'canceled' },
  { label: '已退款', value: 'refunded' },
];

// ==================== 支付状态选项 ====================
export const PAYMENT_STATUS_OPTIONS: SelectOption[] = [
  { label: '全部状态', value: '' },
  { label: '待支付', value: 'pending' },
  { label: '已支付', value: 'paid' },
  { label: '支付失败', value: 'failed' },
  { label: '已退款', value: 'refunded' },
  { label: '已取消', value: 'cancelled' },
];

// ==================== 支付方式选项 ====================
export const PAYMENT_METHOD_OPTIONS: SelectOption[] = [
  { label: '全部方式', value: '' },
  { label: '支付宝', value: 'alipay' },
  { label: '微信支付', value: 'wechat' },
  { label: '余额支付', value: 'balance' },
];

// ==================== 用户角色选项 ====================
export const USER_ROLE_OPTIONS: SelectOption[] = [
  { label: '全部角色', value: '' },
  { label: '普通用户', value: 'user' },
  { label: '陪玩师', value: 'player' },
  { label: '管理员', value: 'admin' },
];

// ==================== 用户状态选项 ====================
export const USER_STATUS_OPTIONS: SelectOption[] = [
  { label: '全部状态', value: '' },
  { label: '正常', value: 'active' },
  { label: '暂停', value: 'suspended' },
  { label: '封禁', value: 'banned' },
];

// ==================== 认证状态选项 ====================
export const VERIFICATION_STATUS_OPTIONS: SelectOption[] = [
  { label: '全部状态', value: '' },
  { label: '待认证', value: 'pending' },
  { label: '已认证', value: 'verified' },
  { label: '已拒绝', value: 'rejected' },
];

// ==================== 评分选项 ====================
export const RATING_OPTIONS: SelectOption[] = [
  { label: '全部评分', value: '' },
  { label: '1星及以上', value: '1' },
  { label: '2星及以上', value: '2' },
  { label: '3星及以上', value: '3' },
  { label: '4星及以上', value: '4' },
  { label: '5星', value: '5' },
];

export const RATING_SELECT_OPTIONS: SelectOption[] = [
  { label: '⭐ 1星 - 非常差', value: '1' },
  { label: '⭐⭐ 2星 - 较差', value: '2' },
  { label: '⭐⭐⭐ 3星 - 一般', value: '3' },
  { label: '⭐⭐⭐⭐ 4星 - 满意', value: '4' },
  { label: '⭐⭐⭐⭐⭐ 5星 - 非常满意', value: '5' },
];

// ==================== 游戏分类选项 ====================
export const GAME_CATEGORY_OPTIONS: SelectOption[] = [
  { label: '全部分类', value: '' },
  { label: 'MOBA', value: 'moba' },
  { label: 'FPS', value: 'fps' },
  { label: 'RPG', value: 'rpg' },
  { label: '策略', value: 'strategy' },
  { label: '体育', value: 'sports' },
  { label: '竞速', value: 'racing' },
  { label: '益智', value: 'puzzle' },
  { label: '其他', value: 'other' },
];

// ==================== 货币选项 ====================
export const CURRENCY_OPTIONS: SelectOption[] = [
  { label: '人民币', value: 'CNY' },
  { label: '美元', value: 'USD' },
];
