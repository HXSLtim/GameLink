/**
 * Discord风格按钮组件
 */
import React, { type ReactNode, type ButtonHTMLAttributes } from 'react';
import './Button.less';

export interface ButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {
  /**
   * 按钮类型
   */
  type?: 'primary' | 'secondary' | 'ghost' | 'danger' | 'link';

  /**
   * 按钮大小
   */
  size?: 'small' | 'medium' | 'large';

  /**
   * 是否禁用
   */
  disabled?: boolean;

  /**
   * 加载状态
   */
  loading?: boolean;

  /**
   * 图标
   */
  icon?: ReactNode;

  /**
   * 块级按钮
   */
  block?: boolean;

  /**
   * 子元素
   */
  children?: ReactNode;

  /**
   * 点击事件
   */
  onClick?: () => void;
}

const Button: React.FC<ButtonProps> = ({
  type = 'primary',
  size = 'medium',
  disabled = false,
  loading = false,
  icon,
  block = false,
  children,
  className = '',
  onClick,
  ...rest
}) => {
  const handleClick = () => {
    if (!disabled && !loading && onClick) {
      onClick();
    }
  };

  const classNames = [
    'discord-button',
    `discord-button--${type}`,
    `discord-button--${size}`,
    loading && 'discord-button--loading',
    disabled && 'discord-button--disabled',
    block && 'discord-button--block',
    className,
  ]
    .filter(Boolean)
    .join(' ');

  return (
    <button
      className={classNames}
      disabled={disabled || loading}
      onClick={handleClick}
      {...rest}
    >
      {loading && (
        <span className="discord-button__loading">
          <span className="discord-button__spinner" />
        </span>
      )}
      {!loading && icon && <span className="discord-button__icon">{icon}</span>}
      {children && <span className="discord-button__text">{children}</span>}
    </button>
  );
};

export default Button;
