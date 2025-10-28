import { ButtonHTMLAttributes, ReactNode } from 'react';
import styles from './Button.module.less';

export interface ButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {
  /** 按钮内容 */
  children: ReactNode;
  /** 按钮类型 */
  variant?: 'primary' | 'secondary' | 'text';
  /** 按钮尺寸 */
  size?: 'small' | 'medium' | 'large';
  /** 是否为块级按钮 */
  block?: boolean;
  /** 加载状态 */
  loading?: boolean;
  /** 禁用状态 */
  disabled?: boolean;
  /** 图标 */
  icon?: ReactNode;
  /** 自定义类名 */
  className?: string;
}

export const Button: React.FC<ButtonProps> = ({
  children,
  variant = 'primary',
  size = 'medium',
  block = false,
  loading = false,
  disabled = false,
  icon,
  className = '',
  ...rest
}) => {
  const classNames = [
    styles.button,
    styles[variant],
    styles[size],
    block ? styles.block : '',
    loading ? styles.loading : '',
    className,
  ]
    .filter(Boolean)
    .join(' ');

  return (
    <button className={classNames} disabled={disabled || loading} {...rest}>
      {loading && <span className={styles.spinner}></span>}
      {!loading && icon && <span className={styles.icon}>{icon}</span>}
      <span className={styles.content}>{children}</span>
    </button>
  );
};
