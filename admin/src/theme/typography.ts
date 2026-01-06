/**
 * 字体样式系统
 * 统一管理字体大小、行高、字重
 *
 * 使用方式:
 * import { fontSize, fontWeight, textStyles } from '@/theme';
 */

/**
 * 字体大小 (px)
 * 基于 Ant Design 的字体比例尺
 */
export const fontSize = {
  /** 12px - 辅助文字、标签 */
  xs: 12,
  /** 14px - 正文 (Ant Design 默认) */
  sm: 14,
  /** 16px - 大正文 */
  md: 16,
  /** 18px - 小标题 */
  lg: 18,
  /** 20px - 标题 */
  xl: 20,
  /** 24px - 大标题 */
  xxl: 24,
  /** 30px - 页面标题 */
  xxxl: 30,
} as const;

/**
 * 行高
 */
export const lineHeight = {
  /** 1 - 紧凑 */
  tight: 1,
  /** 1.4 - 标准 */
  normal: 1.4,
  /** 1.5715 - Ant Design 默认 */
  base: 1.5715,
  /** 1.8 - 宽松 */
  relaxed: 1.8,
} as const;

/**
 * 字重
 */
export const fontWeight = {
  /** 400 - 正常 */
  normal: 400,
  /** 500 - 中等 */
  medium: 500,
  /** 600 - 半粗 */
  semibold: 600,
  /** 700 - 粗体 */
  bold: 700,
} as const;

/**
 * 预设文字样式
 * 可直接展开到 style 属性
 */
export const textStyles = {
  /** 页面大标题 */
  pageTitle: {
    fontSize: fontSize.xxl,
    fontWeight: fontWeight.semibold,
    lineHeight: lineHeight.normal,
  },

  /** 卡片/区块标题 */
  sectionTitle: {
    fontSize: fontSize.lg,
    fontWeight: fontWeight.semibold,
    lineHeight: lineHeight.normal,
  },

  /** 正文 */
  body: {
    fontSize: fontSize.sm,
    fontWeight: fontWeight.normal,
    lineHeight: lineHeight.base,
  },

  /** 小正文 */
  bodySmall: {
    fontSize: fontSize.xs,
    fontWeight: fontWeight.normal,
    lineHeight: lineHeight.base,
  },

  /** 标签/辅助文字 */
  caption: {
    fontSize: fontSize.xs,
    fontWeight: fontWeight.normal,
    lineHeight: lineHeight.normal,
  },

  /** 统计数字 */
  statistic: {
    fontSize: fontSize.xxl,
    fontWeight: fontWeight.semibold,
    lineHeight: lineHeight.tight,
  },

  /** 按钮文字 */
  button: {
    fontSize: fontSize.sm,
    fontWeight: fontWeight.medium,
    lineHeight: lineHeight.normal,
  },
} as const;

export type FontSizeKey = keyof typeof fontSize;
export type FontWeightKey = keyof typeof fontWeight;
export type TextStyleKey = keyof typeof textStyles;
