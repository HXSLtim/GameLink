import React from 'react';
import styles from './Space.module.less';

export type SpaceDirection = 'horizontal' | 'vertical';

export interface SpaceProps {
  size?: 'small' | 'medium' | 'large' | number;
  direction?: SpaceDirection;
  wrap?: boolean;
  align?: 'start' | 'center' | 'end' | 'baseline';
  className?: string;
  style?: React.CSSProperties;
  children?: React.ReactNode;
  split?: React.ReactNode;
}

export const Space: React.FC<SpaceProps> = ({
  size = 'medium',
  direction = 'horizontal',
  wrap = false,
  align = 'center',
  className = '',
  style,
  children,
  split,
}) => {
  const gapValue = typeof size === 'number' ? `${size}px` : undefined;
  const classNames = [
    styles.space,
    styles[direction],
    wrap ? styles.wrap : styles.noWrap,
    styles[`align-${align}`],
    className,
  ]
    .filter(Boolean)
    .join(' ');

  const spaceStyle: React.CSSProperties = {
    ...style,
    gap: gapValue,
  };

  const items = React.Children.toArray(children);

  return (
    <div className={classNames} style={spaceStyle}>
      {split
        ? items.map((child, index) => (
            <React.Fragment key={index}>
              {child}
              {index < items.length - 1 ? <span className={styles.split}>{split}</span> : null}
            </React.Fragment>
          ))
        : items}
    </div>
  );
};