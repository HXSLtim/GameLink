/**
 * Discord风格卡片组件
 */
import React, { type ReactNode } from 'react';
import './Card.less';

export interface CardProps {
  /**
   * 卡片内容
   */
  children: ReactNode;

  /**
   * 卡片标题
   */
  title?: string;

  /**
   * 卡片副标题
   */
  subtitle?: string;

  /**
   * 头部右侧额外内容
   */
  extra?: ReactNode;

  /**
   * 是否可悬浮（hover效果）
   */
  hoverable?: boolean;

  /**
   * 是否可点击
   */
  clickable?: boolean;

  /**
   * 自定义类名
   */
  className?: string;

  /**
   * 点击事件
   */
  onClick?: () => void;

  /**
   * 卡片内边距大小
   */
  padding?: 'none' | 'small' | 'medium' | 'large';

  /**
   * 是否显示边框
   */
  bordered?: boolean;
}

const Card: React.FC<CardProps> = ({
  children,
  title,
  subtitle,
  extra,
  hoverable = false,
  clickable = false,
  className = '',
  onClick,
  padding = 'medium',
  bordered = true,
}) => {
  const classNames = [
    'discord-card',
    hoverable && 'discord-card--hoverable',
    clickable && 'discord-card--clickable',
    bordered && 'discord-card--bordered',
    `discord-card--padding-${padding}`,
    className,
  ]
    .filter(Boolean)
    .join(' ');

  const hasHeader = title || subtitle || extra;

  return (
    <div className={classNames} onClick={onClick}>
      {hasHeader && (
        <div className="discord-card__header">
          <div className="discord-card__header-content">
            {title && <h3 className="discord-card__title">{title}</h3>}
            {subtitle && <p className="discord-card__subtitle">{subtitle}</p>}
          </div>
          {extra && <div className="discord-card__extra">{extra}</div>}
        </div>
      )}
      <div className="discord-card__body">{children}</div>
    </div>
  );
};

export default Card;
