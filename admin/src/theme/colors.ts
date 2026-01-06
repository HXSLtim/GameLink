/**
 * 语义化颜色令牌
 * 统一管理应用中的颜色，避免魔法数字
 *
 * 使用方式:
 * import { semanticColors, chartColors } from '@/theme';
 */

/**
 * 语义化颜色 - 基于 Ant Design 主题变量
 * 优先使用 CSS 变量，确保主题切换时自动适配
 */
export const semanticColors = {
  // 主色调
  primary: 'var(--ant-color-primary)',
  primaryHover: 'var(--ant-color-primary-hover)',
  primaryActive: 'var(--ant-color-primary-active)',
  primaryBg: 'var(--ant-color-primary-bg)',

  // 功能色
  success: 'var(--ant-color-success)',
  warning: 'var(--ant-color-warning)',
  error: 'var(--ant-color-error)',
  info: 'var(--ant-color-info)',

  // 中性色
  text: 'var(--ant-color-text)',
  textSecondary: 'var(--ant-color-text-secondary)',
  textTertiary: 'var(--ant-color-text-tertiary)',
  textDisabled: 'var(--ant-color-text-disabled)',

  // 背景色
  bgContainer: 'var(--ant-color-bg-container)',
  bgElevated: 'var(--ant-color-bg-elevated)',
  bgLayout: 'var(--ant-color-bg-layout)',

  // 边框色
  border: 'var(--ant-color-border)',
  borderSecondary: 'var(--ant-color-border-secondary)',

  // 分割线
  split: 'var(--ant-color-split)',
} as const;

/**
 * 图表颜色调色板
 * 用于饼图、折线图等数据可视化
 */
export const chartColors = {
  palette: [
    '#0088FE', // 蓝色
    '#00C49F', // 青色
    '#FFBB28', // 黄色
    '#FF8042', // 橙色
    '#8884d8', // 紫色
    '#82ca9d', // 绿色
    '#ffc658', // 金色
    '#ff7c7c', // 红色
  ],

  // 语义化图表颜色
  revenue: '#faad14', // 收入 - 金色
  users: '#1677ff', // 用户 - 蓝色
  orders: '#52c41a', // 订单 - 绿色
  players: '#722ed1', // 陪玩 - 紫色
} as const;

/**
 * 状态颜色映射
 * 用于订单状态、支付状态等标签
 */
export const statusColors = {
  // 订单状态
  pending: '#faad14', // 待确认 - 金色
  confirmed: '#1677ff', // 已确认 - 蓝色
  in_progress: '#1677ff', // 进行中 - 蓝色
  completed: '#52c41a', // 已完成 - 绿色
  canceled: '#8c8c8c', // 已取消 - 灰色
  refunded: '#ff4d4f', // 已退款 - 红色

  // 支付状态
  paid: '#52c41a', // 已支付 - 绿色
  failed: '#ff4d4f', // 支付失败 - 红色
} as const;

/**
 * 图标背景色
 * 用于 StatCard 等组件的图标背景
 */
export const iconBgColors = {
  primary: 'var(--ant-color-primary)',
  success: 'var(--ant-color-success)',
  warning: 'var(--ant-color-warning)',
  error: 'var(--ant-color-error)',
  purple: '#722ed1',
  cyan: '#13c2c2',
  magenta: '#eb2f96',
  volcano: '#fa541c',
} as const;

export type SemanticColorKey = keyof typeof semanticColors;
export type ChartColorKey = keyof typeof chartColors;
export type StatusColorKey = keyof typeof statusColors;
export type IconBgColorKey = keyof typeof iconBgColors;
