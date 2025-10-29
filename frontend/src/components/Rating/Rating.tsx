import React from 'react';
import { StarFilledIcon, StarOutlinedIcon } from '../Icons/icons';
import styles from './Rating.module.less';

export interface RatingProps {
  /** 评分值 (1-5) */
  value: number;
  /** 是否显示数字 */
  showNumber?: boolean;
  /** 是否显示文本描述 */
  showText?: boolean;
  /** 星星大小 */
  size?: number;
  /** 自定义类名 */
  className?: string;
}

/** 评分文本映射 */
const RATING_TEXT: Record<number, string> = {
  1: '非常差',
  2: '较差',
  3: '一般',
  4: '满意',
  5: '非常满意',
};

/**
 * 根据评分获取主题颜色类名
 */
const getRatingColorClass = (rating: number): string => {
  if (rating >= 4.5) return styles.ratingExcellent;
  if (rating >= 3.5) return styles.ratingGood;
  if (rating >= 2.5) return styles.ratingAverage;
  return styles.ratingPoor;
};

/**
 * 评分组件
 * 
 * @component
 * @example
 * ```tsx
 * <Rating value={4} showNumber showText />
 * <Rating value={3.5} size={16} />
 * ```
 */
export const Rating: React.FC<RatingProps> = ({
  value,
  showNumber = false,
  showText = false,
  size = 16,
  className,
}) => {
  const roundedValue = Math.round(value);
  const text = RATING_TEXT[roundedValue] || `${value}星`;
  const colorClass = getRatingColorClass(value);

  return (
    <span className={`${styles.rating} ${colorClass} ${className || ''}`}>
      {Array.from({ length: 5 }).map((_, index) => {
        const isFilled = index < roundedValue;
        const StarComponent = isFilled ? StarFilledIcon : StarOutlinedIcon;
        
        return (
          <StarComponent
            key={index}
            size={size}
            className={`${styles.star} ${isFilled ? styles.filled : styles.outlined}`}
          />
        );
      })}
      {showNumber && <span className={styles.number}>{value}</span>}
      {showText && <span className={styles.text}>{text}</span>}
    </span>
  );
};

/**
 * 简单评分显示（只显示星星和数字）
 */
export const SimpleRating: React.FC<{ value: number; size?: number }> = ({ value, size = 14 }) => {
  const roundedValue = Math.round(value);
  const colorClass = getRatingColorClass(value);

  return (
    <span className={`${styles.simpleRating} ${colorClass}`}>
      {Array.from({ length: 5 }).map((_, index) => {
        const isFilled = index < roundedValue;
        const StarComponent = isFilled ? StarFilledIcon : StarOutlinedIcon;
        
        return (
          <StarComponent 
            key={index} 
            size={size} 
            className={`${styles.star} ${isFilled ? styles.filled : styles.outlined}`}
          />
        );
      })}
      <span className={styles.number}>{value}</span>
    </span>
  );
};

/**
 * 获取评分文本
 */
export const getRatingText = (rating: number): string => {
  const rounded = Math.round(rating);
  return RATING_TEXT[rounded] || `${rating}星`;
};

/**
 * 获取评分颜色（CSS 变量名）
 */
export const getRatingColor = (rating: number): string => {
  if (rating >= 4.5) return 'var(--rating-excellent)';
  if (rating >= 3.5) return 'var(--rating-good)';
  if (rating >= 2.5) return 'var(--rating-average)';
  return 'var(--rating-poor)';
};

