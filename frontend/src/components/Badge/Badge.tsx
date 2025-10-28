import React from 'react';
import styles from './Badge.module.less';

export interface BadgeProps {
  count?: number;
  dot?: boolean;
  showZero?: boolean;
  overflowCount?: number;
  children?: React.ReactNode;
  className?: string;
  style?: React.CSSProperties;
}

export const Badge: React.FC<BadgeProps> = ({
  count = 0,
  dot = false,
  showZero = false,
  overflowCount = 99,
  children,
  className,
  style,
}) => {
  const displayCount = count > overflowCount ? `${overflowCount}+` : count;
  const shouldShow = dot || count > 0 || (showZero && count === 0);

  if (!children) {
    // 独立徽章模式
    if (!shouldShow) return null;

    return (
      <span className={`${styles.badgeStandalone} ${className || ''}`} style={style}>
        {dot ? (
          <span className={styles.dot} />
        ) : (
          <span className={styles.count}>{displayCount}</span>
        )}
      </span>
    );
  }

  // 包裹模式
  return (
    <span className={`${styles.badgeWrapper} ${className || ''}`} style={style}>
      {children}
      {shouldShow && (
        <span className={styles.badge}>
          {dot ? (
            <span className={styles.dot} />
          ) : (
            <span className={styles.count}>{displayCount}</span>
          )}
        </span>
      )}
    </span>
  );
};
