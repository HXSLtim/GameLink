/**
 * 间距比例尺
 * 基于 4px 基准单位的间距系统
 *
 * 使用方式:
 * import { spacing, gap } from '@/theme';
 * style={{ padding: spacing.md, gap: gap.sm }}
 */

/**
 * 基础间距单位 (px)
 */
export const BASE_UNIT = 4;

/**
 * 间距比例尺
 * 基于 4px 基准，遵循 4-8-12-16-24-32-48 的常见设计系统
 */
export const spacing = {
  /** 0px */
  none: 0,
  /** 4px - 极小间距 */
  xs: BASE_UNIT, // 4
  /** 8px - 小间距 */
  sm: BASE_UNIT * 2, // 8
  /** 12px - 中小间距 */
  md: BASE_UNIT * 3, // 12
  /** 16px - 中间距 */
  lg: BASE_UNIT * 4, // 16
  /** 24px - 大间距 */
  xl: BASE_UNIT * 6, // 24
  /** 32px - 超大间距 */
  xxl: BASE_UNIT * 8, // 32
  /** 48px - 巨大间距 */
  xxxl: BASE_UNIT * 12, // 48
} as const;

/**
 * 组件间隙 (用于 Grid gutter, Flex gap 等)
 */
export const gap = {
  /** 8px - 紧凑 */
  compact: [8, 8] as [number, number],
  /** 16px - 标准 */
  standard: [16, 16] as [number, number],
  /** 24px - 宽松 */
  relaxed: [24, 24] as [number, number],
} as const;

/**
 * 页面内边距
 */
export const pagePadding = {
  /** 移动端 */
  mobile: spacing.lg, // 16
  /** 平板 */
  tablet: spacing.xl, // 24
  /** 桌面 */
  desktop: spacing.xl, // 24
} as const;

/**
 * 卡片内边距
 */
export const cardPadding = {
  /** 紧凑 */
  compact: spacing.md, // 12
  /** 标准 */
  standard: spacing.lg, // 16
  /** 宽松 */
  relaxed: spacing.xl, // 24
} as const;

/**
 * 边框圆角
 */
export const borderRadius = {
  /** 0px - 无圆角 */
  none: 0,
  /** 2px - 极小 */
  xs: 2,
  /** 4px - 小 */
  sm: 4,
  /** 6px - 中 (Ant Design 默认) */
  md: 6,
  /** 8px - 大 */
  lg: 8,
  /** 12px - 超大 */
  xl: 12,
  /** 50% - 圆形 */
  full: '50%',
} as const;

export type SpacingKey = keyof typeof spacing;
export type GapKey = keyof typeof gap;
export type BorderRadiusKey = keyof typeof borderRadius;
