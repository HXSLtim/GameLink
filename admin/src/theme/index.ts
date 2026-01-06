/**
 * 设计令牌系统 - 统一导出
 *
 * 使用方式:
 * import { semanticColors, spacing, fontSize } from '@/theme';
 *
 * 或按需导入:
 * import { chartColors } from '@/theme/colors';
 */

// 颜色
export {
  semanticColors,
  chartColors,
  statusColors,
  iconBgColors,
  type SemanticColorKey,
  type ChartColorKey,
  type StatusColorKey,
  type IconBgColorKey,
} from './colors';

// 间距
export {
  BASE_UNIT,
  spacing,
  gap,
  pagePadding,
  cardPadding,
  borderRadius,
  type SpacingKey,
  type GapKey,
  type BorderRadiusKey,
} from './spacing';

// 字体
export {
  fontSize,
  lineHeight,
  fontWeight,
  textStyles,
  type FontSizeKey,
  type FontWeightKey,
  type TextStyleKey,
} from './typography';

/**
 * 响应式断点 (与 Ant Design Grid 一致)
 */
export const breakpoints = {
  xs: 480,
  sm: 576,
  md: 768,
  lg: 992,
  xl: 1200,
  xxl: 1600,
} as const;

/**
 * 媒体查询
 */
export const mediaQuery = {
  mobile: `@media (max-width: ${breakpoints.sm - 1}px)`,
  tablet: `@media (min-width: ${breakpoints.sm}px) and (max-width: ${breakpoints.lg - 1}px)`,
  desktop: `@media (min-width: ${breakpoints.lg}px)`,
  // 最小宽度查询
  minSm: `@media (min-width: ${breakpoints.sm}px)`,
  minMd: `@media (min-width: ${breakpoints.md}px)`,
  minLg: `@media (min-width: ${breakpoints.lg}px)`,
  minXl: `@media (min-width: ${breakpoints.xl}px)`,
} as const;

/**
 * 动画时长
 */
export const duration = {
  fast: 150,
  normal: 300,
  slow: 500,
} as const;

/**
 * 阴影
 */
export const shadow = {
  none: 'none',
  sm: '0 1px 2px 0 rgba(0, 0, 0, 0.05)',
  md: '0 4px 6px -1px rgba(0, 0, 0, 0.1)',
  lg: '0 10px 15px -3px rgba(0, 0, 0, 0.1)',
  xl: '0 20px 25px -5px rgba(0, 0, 0, 0.1)',
} as const;
