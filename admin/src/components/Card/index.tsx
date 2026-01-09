/**
 * 统一卡片组件
 * 包装 Ant Design Card，提供标准化的样式变体
 *
 * 使用方式:
 * import { Card } from '@/components';
 * <Card cardVariant="borderless" cardPadding="relaxed">内容</Card>
 */
import React, { useMemo } from 'react';
import type { CSSProperties, ReactNode } from 'react';
import { Card as AntCard } from 'antd';
import type { CardProps as AntCardProps } from 'antd';
import { cardPadding } from '@/theme';

/**
 * 卡片变体样式
 */
const variantStyles: Record<string, CSSProperties> = {
  /** 默认 - 有边框 */
  bordered: {},
  /** 无边框 */
  borderless: { border: 'none' },
  /** 阴影 - 无边框但有阴影 */
  elevated: {
    border: 'none',
    boxShadow:
      '0 1px 2px 0 rgba(0, 0, 0, 0.03), 0 1px 6px -1px rgba(0, 0, 0, 0.02), 0 2px 4px 0 rgba(0, 0, 0, 0.02)',
  },
  /** 填充 - 有背景色无边框 */
  filled: {
    border: 'none',
    backgroundColor: 'var(--ant-color-fill-quaternary)',
  },
};

/**
 * 内边距映射
 */
const paddingMap = {
  none: 0,
  compact: cardPadding.compact,
  standard: cardPadding.standard,
  relaxed: cardPadding.relaxed,
} as const;

export type CardVariant = 'bordered' | 'borderless' | 'elevated' | 'filled';
export type CardPadding = keyof typeof paddingMap;

export interface CardProps extends AntCardProps {
  /** 卡片变体 */
  cardVariant?: CardVariant;
  /** 内边距 */
  cardPadding?: CardPadding;
  /** 子元素 */
  children?: ReactNode;
}

/**
 * 统一卡片组件
 * 优化: 使用 React.memo + useMemo 减少重渲染
 */
export const Card: React.FC<CardProps> = React.memo(({
  cardVariant = 'bordered',
  cardPadding: padding = 'standard',
  style,
  children,
  ...props
}) => {
  // 缓存样式计算
  const finalStyle = useMemo<CSSProperties>(() => {
    const variantStyle = variantStyles[cardVariant] || {};
    return { ...variantStyle, ...style };
  }, [cardVariant, style]);

  // 缓存 body 样式
  const bodyStyle = useMemo<CSSProperties>(() => ({
    padding: paddingMap[padding],
  }), [padding]);

  return (
    <AntCard style={finalStyle} styles={{ body: bodyStyle }} {...props}>
      {children}
    </AntCard>
  );
});

/**
 * 统计卡片 - 用于仪表盘
 * 预设为无边框 + 标准内边距
 */
export const StatisticCard: React.FC<CardProps> = React.memo((props) => {
  return <Card cardVariant="borderless" cardPadding="standard" {...props} />;
});

/**
 * 内容卡片 - 用于表单、详情等
 * 预设为有边框 + 宽松内边距
 */
export const ContentCard: React.FC<CardProps> = React.memo((props) => {
  return <Card cardVariant="bordered" cardPadding="relaxed" {...props} />;
});

export default Card;
