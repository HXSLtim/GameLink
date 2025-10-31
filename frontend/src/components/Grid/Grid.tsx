import React from 'react';
import styles from './Grid.module.less';

export type GridJustify = 'start' | 'center' | 'end' | 'space-between' | 'space-around';
export type GridAlign = 'top' | 'middle' | 'bottom';

export interface RowProps {
  gutter?: number | [number, number];
  justify?: GridJustify;
  align?: GridAlign;
  wrap?: boolean;
  className?: string;
  style?: React.CSSProperties;
  children?: React.ReactNode;
}

export const Row: React.FC<RowProps> = ({
  gutter = 0,
  justify = 'start',
  align = 'top',
  wrap = true,
  className = '',
  style,
  children,
}) => {
  const [horizontal, vertical] = Array.isArray(gutter) ? gutter : [gutter, 0];

  const classNames = [
    styles.row,
    styles[`justify-${justify}`],
    styles[`align-${align}`],
    wrap ? styles.wrap : styles.noWrap,
    className,
  ]
    .filter(Boolean)
    .join(' ');

  const rowStyle: React.CSSProperties = {
    ...style,
    marginLeft: horizontal ? -(horizontal / 2) : undefined,
    marginRight: horizontal ? -(horizontal / 2) : undefined,
    rowGap: vertical || undefined,
  };

  return <div className={classNames} style={rowStyle}>{children}</div>;
};

export interface ResponsiveColSpan {
  span?: number;
  offset?: number;
}

export interface ColProps extends ResponsiveColSpan {
  xs?: ResponsiveColSpan;
  sm?: ResponsiveColSpan;
  md?: ResponsiveColSpan;
  lg?: ResponsiveColSpan;
  xl?: ResponsiveColSpan;
  className?: string;
  style?: React.CSSProperties;
  children?: React.ReactNode;
}

const calcPercent = (span?: number) => (span ? (span / 24) * 100 : undefined);

export const Col: React.FC<ColProps> = ({
  span,
  offset,
  xs,
  sm,
  md,
  lg,
  xl,
  className = '',
  style,
  children,
}) => {
  const classNames = [styles.col, className].filter(Boolean).join(' ');

  const baseWidth = calcPercent(span);
  const baseMarginLeft = calcPercent(offset);

  const colStyle: React.CSSProperties = {
    ...style,
    flexBasis: baseWidth ? `${baseWidth}%` : undefined,
    maxWidth: baseWidth ? `${baseWidth}%` : undefined,
    marginLeft: baseMarginLeft ? `${baseMarginLeft}%` : undefined,
  };

  return (
    <div
      className={classNames}
      style={colStyle}
      data-xs-span={xs?.span}
      data-sm-span={sm?.span}
      data-md-span={md?.span}
      data-lg-span={lg?.span}
      data-xl-span={xl?.span}
      data-xs-offset={xs?.offset}
      data-sm-offset={sm?.offset}
      data-md-offset={md?.offset}
      data-lg-offset={lg?.offset}
      data-xl-offset={xl?.offset}
    >
      {children}
    </div>
  );
};

export interface GridProps {
  children?: React.ReactNode;
}

export const Grid: React.FC<GridProps> = ({ children }) => {
  return <div className={styles.grid}>{children}</div>;
};