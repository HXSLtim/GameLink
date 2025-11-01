import React, { InputHTMLAttributes } from 'react';
import styles from './Checkbox.module.less';

export interface CheckboxProps extends Omit<InputHTMLAttributes<HTMLInputElement>, 'type'> {
  /**
   * 是否选中
   */
  checked?: boolean;
  /**
   * 是否为部分选中状态（半选）
   */
  indeterminate?: boolean;
  /**
   * 是否禁用
   */
  disabled?: boolean;
  /**
   * 变化时的回调
   */
  onChange?: (e: React.ChangeEvent<HTMLInputElement>) => void;
  /**
   * 子元素（标签文本）
   */
  children?: React.ReactNode;
  /**
   * 自定义类名
   */
  className?: string;
}

/**
 * Checkbox 复选框组件
 */
export const Checkbox: React.FC<CheckboxProps> = ({
  checked = false,
  indeterminate = false,
  disabled = false,
  onChange,
  children,
  className = '',
  ...rest
}) => {
  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (!disabled && onChange) {
      onChange(e);
    }
  };

  const classNames = [
    styles.checkbox,
    checked ? styles.checked : '',
    indeterminate ? styles.indeterminate : '',
    disabled ? styles.disabled : '',
    className,
  ]
    .filter(Boolean)
    .join(' ');

  return (
    <label className={classNames}>
      <input
        type="checkbox"
        className={styles.input}
        checked={checked}
        disabled={disabled}
        onChange={handleChange}
        {...rest}
      />
      <span className={styles.box}>
        {indeterminate ? (
          <span className={styles.indeterminateLine} />
        ) : checked ? (
          <svg className={styles.checkIcon} viewBox="0 0 16 16" fill="none">
            <path
              d="M3 8L6.5 11.5L13 5"
              stroke="currentColor"
              strokeWidth="2"
              strokeLinecap="round"
              strokeLinejoin="round"
            />
          </svg>
        ) : null}
      </span>
      {children && <span className={styles.label}>{children}</span>}
    </label>
  );
};


