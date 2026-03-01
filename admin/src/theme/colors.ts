/**
 * 语义化颜色令牌
 * 统一管理应用中的颜色，避免魔法数字
 *
 * 使用方式:
 * import { semanticColors, chartColors, brandColors } from '@/theme';
 */

/**
 * GameLink 品牌色
 * 统一走 Ant Design 主题变量，避免出现多套主色源。
 */
export const brandColors = {
  primary: 'var(--ant-color-primary)',
  primaryLight: 'var(--ant-color-primary-hover)',
  primaryDark: 'var(--ant-color-primary-active)',
  primaryGradient: 'linear-gradient(135deg, var(--ant-color-primary) 0%, var(--ant-color-primary-active) 100%)',

  gold: 'var(--ant-color-warning)',
  goldGradient: 'linear-gradient(135deg, var(--ant-color-warning) 0%, #d97706 100%)',

  purple: '#722ed1',
  cyan: 'var(--ant-color-info)',
} as const;

/**
 * 语义化颜色 - 基于 Ant Design 主题变量
 * 优先使用 CSS 变量，确保主题切换时自动适配
 */
export const semanticColors = {
  // 主色调（使用 GameLink 品牌色）
  primary: brandColors.primary,
  primaryHover: brandColors.primaryLight,
  primaryActive: brandColors.primaryDark,
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
    'var(--ant-color-primary)', // 蓝色主色
    'var(--ant-color-info)', // 青色信息色
    'var(--ant-color-warning)', // 黄色告警色
    '#ff8042', // 橙色
    '#8884d8', // 紫色
    'var(--ant-color-success)', // 绿色成功色
    '#ffc658', // 金色
    'var(--ant-color-error)', // 红色错误色
  ],

  // 语义化图表颜色（跟随主题变量）
  revenue: 'var(--ant-color-warning)',
  users: 'var(--ant-color-primary)',
  orders: 'var(--ant-color-primary)',
  players: '#722ed1',
  games: 'var(--ant-color-info)',
  vip: 'var(--ant-color-warning)',
} as const;

/**
 * 状态颜色映射
 * 用于订单状态、支付状态等标签
 */
export const statusColors = {
  // 订单状态
  pending: 'var(--ant-color-warning)',
  confirmed: 'var(--ant-color-primary)',
  in_progress: 'var(--ant-color-primary)',
  completed: 'var(--ant-color-success)',
  canceled: 'var(--ant-color-text-tertiary)',
  refunded: 'var(--ant-color-error)',

  // 支付状态
  paid: 'var(--ant-color-success)',
  failed: 'var(--ant-color-error)',
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
  cyan: 'var(--ant-color-info)',
  magenta: '#eb2f96',
  volcano: '#fa541c',
} as const;

export type SemanticColorKey = keyof typeof semanticColors;
export type ChartColorKey = keyof typeof chartColors;
export type StatusColorKey = keyof typeof statusColors;
export type IconBgColorKey = keyof typeof iconBgColors;
