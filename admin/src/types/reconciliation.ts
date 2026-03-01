/**
 * 对账管理类型定义和常量
 */

import type { ReconciliationType, ReconciliationStatus } from '@/api/reconciliation';

/**
 * 对账状态文本映射
 */
export const RECONCILIATION_STATUS_TEXT: Record<ReconciliationStatus, string> = {
  pending: '待处理',
  progress: '进行中',
  success: '已完成',
  failed: '失败',
  exception: '异常',
};

/**
 * 对账状态颜色映射
 */
export const RECONCILIATION_STATUS_COLOR: Record<ReconciliationStatus, string> = {
  pending: 'default',
  progress: 'processing',
  success: 'success',
  failed: 'error',
  exception: 'warning',
};

/**
 * 对账类型文本映射
 */
export const RECONCILIATION_TYPE_TEXT: Record<ReconciliationType, string> = {
  payment: '支付对账',
  internal: '内部对账',
  bank: '银行对账',
  manual: '手动对账',
};

/**
 * 对账类型颜色映射
 */
export const RECONCILIATION_TYPE_COLOR: Record<ReconciliationType, string> = {
  payment: 'blue',
  internal: 'green',
  bank: 'purple',
  manual: 'orange',
};
