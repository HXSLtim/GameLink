import React from 'react';
import styles from './Skeleton.module.less';

export interface SkeletonProps {
  /** 骨架屏类型 */
  variant?: 'text' | 'rect' | 'circle' | 'card';
  /** 宽度 */
  width?: string | number;
  /** 高度 */
  height?: string | number;
  /** 是否显示动画 */
  animated?: boolean;
  /** 自定义类名 */
  className?: string;
  /** 自定义样式 */
  style?: React.CSSProperties;
}

export const Skeleton: React.FC<SkeletonProps> = ({
  variant = 'text',
  width,
  height,
  animated = true,
  className,
  style,
}) => {
  const variantClass = styles[variant];
  const animatedClass = animated ? styles.animated : '';

  const customStyle: React.CSSProperties = {
    ...style,
    width: typeof width === 'number' ? `${width}px` : width,
    height: typeof height === 'number' ? `${height}px` : height,
  };

  return (
    <div
      className={`${styles.skeleton} ${variantClass} ${animatedClass} ${className || ''}`}
      style={customStyle}
    />
  );
};

/**
 * 表格骨架屏
 */
export interface TableSkeletonProps {
  rows?: number;
  columns?: number;
  className?: string;
}

export const TableSkeleton: React.FC<TableSkeletonProps> = ({
  rows = 5,
  columns = 6,
  className,
}) => {
  return (
    <div className={`${styles.tableSkeleton} ${className || ''}`}>
      {/* 表头 */}
      <div className={styles.tableHeader}>
        {Array.from({ length: columns }).map((_, i) => (
          <Skeleton key={`header-${i}`} variant="rect" height={40} />
        ))}
      </div>
      {/* 表格行 */}
      {Array.from({ length: rows }).map((_, rowIndex) => (
        <div key={`row-${rowIndex}`} className={styles.tableRow}>
          {Array.from({ length: columns }).map((_, colIndex) => (
            <Skeleton key={`cell-${rowIndex}-${colIndex}`} variant="rect" height={56} />
          ))}
        </div>
      ))}
    </div>
  );
};

/**
 * 卡片骨架屏
 */
export interface CardSkeletonProps {
  hasImage?: boolean;
  lines?: number;
  className?: string;
}

export const CardSkeleton: React.FC<CardSkeletonProps> = ({
  hasImage = false,
  lines = 3,
  className,
}) => {
  return (
    <div className={`${styles.cardSkeleton} ${className || ''}`}>
      {hasImage && <Skeleton variant="rect" height={200} className={styles.cardImage} />}
      <div className={styles.cardContent}>
        {Array.from({ length: lines }).map((_, i) => (
          <Skeleton
            key={`line-${i}`}
            variant="text"
            width={i === lines - 1 ? '60%' : '100%'}
            className={styles.cardLine}
          />
        ))}
      </div>
    </div>
  );
};

/**
 * 统计卡片骨架屏
 */
export const StatCardSkeleton: React.FC<{ className?: string }> = ({ className }) => {
  return (
    <div className={`${styles.statCardSkeleton} ${className || ''}`}>
      <div className={styles.statIcon}>
        <Skeleton variant="rect" width={56} height={56} />
      </div>
      <div className={styles.statContent}>
        <Skeleton variant="text" width="40%" height={14} />
        <Skeleton variant="text" width="60%" height={32} />
        <Skeleton variant="text" width="50%" height={12} />
      </div>
    </div>
  );
};

/**
 * 列表项骨架屏
 */
export interface ListItemSkeletonProps {
  hasAvatar?: boolean;
  lines?: number;
  className?: string;
}

export const ListItemSkeleton: React.FC<ListItemSkeletonProps> = ({
  hasAvatar = true,
  lines = 2,
  className,
}) => {
  return (
    <div className={`${styles.listItemSkeleton} ${className || ''}`}>
      {hasAvatar && <Skeleton variant="circle" width={40} height={40} />}
      <div className={styles.listItemContent}>
        {Array.from({ length: lines }).map((_, i) => (
          <Skeleton
            key={`line-${i}`}
            variant="text"
            width={i === 0 ? '80%' : '60%'}
            height={i === 0 ? 16 : 14}
          />
        ))}
      </div>
    </div>
  );
};

/**
 * 页面骨架屏 - 用于整个页面的加载状态
 */
export interface PageSkeletonProps {
  className?: string;
}

export const PageSkeleton: React.FC<PageSkeletonProps> = ({ className }) => {
  return (
    <div className={`${styles.pageSkeleton} ${className || ''}`}>
      {/* 页面头部 */}
      <div className={styles.pageHeader}>
        <Skeleton variant="rect" width={100} height={40} />
        <Skeleton variant="text" width={200} height={32} />
      </div>

      {/* 页面内容 */}
      <div className={styles.pageContent}>
        {/* 左侧内容 */}
        <div className={styles.pageColumn}>
          <CardSkeleton hasImage={true} lines={6} />
          <CardSkeleton lines={4} />
        </div>

        {/* 右侧内容 */}
        <div className={styles.pageColumn}>
          <CardSkeleton lines={3} />
          <CardSkeleton lines={4} />
        </div>
      </div>
    </div>
  );
};
