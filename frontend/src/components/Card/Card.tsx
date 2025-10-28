import { ReactNode } from 'react';
import styles from './Card.module.less';

export interface CardProps {
  /** 卡片内容 */
  children: ReactNode;
  /** 标题 */
  title?: ReactNode;
  /** 额外内容（显示在右上角） */
  extra?: ReactNode;
  /** 是否有边框 */
  bordered?: boolean;
  /** 是否可悬停 */
  hoverable?: boolean;
  /** 自定义类名 */
  className?: string;
  /** 自定义样式 */
  style?: React.CSSProperties;
}

export const Card: React.FC<CardProps> = ({
  children,
  title,
  extra,
  bordered = true,
  hoverable = false,
  className = '',
  style,
}) => {
  const classNames = [
    styles.card,
    bordered ? styles.bordered : '',
    hoverable ? styles.hoverable : '',
    className,
  ]
    .filter(Boolean)
    .join(' ');

  return (
    <div className={classNames} style={style}>
      {(title || extra) && (
        <div className={styles.header}>
          {title && <div className={styles.title}>{title}</div>}
          {extra && <div className={styles.extra}>{extra}</div>}
        </div>
      )}
      <div className={styles.body}>{children}</div>
    </div>
  );
};
