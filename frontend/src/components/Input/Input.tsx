import { InputHTMLAttributes, ReactNode, forwardRef, useState } from 'react';
import styles from './Input.module.less';

export interface InputProps extends Omit<InputHTMLAttributes<HTMLInputElement>, 'size'> {
  /** 输入框尺寸 */
  size?: 'small' | 'medium' | 'large';
  /** 前缀图标 */
  prefix?: ReactNode;
  /** 后缀图标 */
  suffix?: ReactNode;
  /** 是否有错误 */
  error?: boolean;
  /** 错误信息 */
  errorMessage?: string;
  /** 是否允许清空 */
  allowClear?: boolean;
  /** 自定义类名 */
  className?: string;
}

export const Input = forwardRef<HTMLInputElement, InputProps>(
  (
    {
      size = 'medium',
      prefix,
      suffix,
      error = false,
      errorMessage,
      allowClear = false,
      className = '',
      value,
      onChange,
      disabled,
      ...rest
    },
    ref,
  ) => {
    const [internalValue, setInternalValue] = useState(value || '');

    const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
      setInternalValue(e.target.value);
      onChange?.(e);
    };

    const handleClear = () => {
      const event = {
        target: { value: '' },
      } as React.ChangeEvent<HTMLInputElement>;
      setInternalValue('');
      onChange?.(event);
    };

    const showClear = allowClear && internalValue && !disabled;

    const wrapperClassNames = [
      styles.wrapper,
      styles[size],
      error ? styles.error : '',
      disabled ? styles.disabled : '',
      className,
    ]
      .filter(Boolean)
      .join(' ');

    return (
      <div className={styles.container}>
        <div className={wrapperClassNames}>
          {prefix && <span className={styles.prefix}>{prefix}</span>}
          <input
            ref={ref}
            className={styles.input}
            value={value !== undefined ? value : internalValue}
            onChange={handleChange}
            disabled={disabled}
            {...rest}
          />
          {showClear && (
            <button type="button" className={styles.clearBtn} onClick={handleClear}>
              <svg width="16" height="16" viewBox="0 0 16 16" fill="currentColor">
                <path d="M8 0C3.58 0 0 3.58 0 8C0 12.42 3.58 16 8 16C12.42 16 16 12.42 16 8C16 3.58 12.42 0 8 0ZM11.3 10.3C11.5 10.5 11.5 10.9 11.3 11.1C11.2 11.2 11 11.2 10.9 11.2C10.8 11.2 10.6 11.2 10.5 11.1L8 8.6L5.5 11.1C5.4 11.2 5.2 11.2 5.1 11.2C5 11.2 4.8 11.2 4.7 11.1C4.5 10.9 4.5 10.5 4.7 10.3L7.2 7.8L4.7 5.3C4.5 5.1 4.5 4.7 4.7 4.5C4.9 4.3 5.3 4.3 5.5 4.5L8 7L10.5 4.5C10.7 4.3 11.1 4.3 11.3 4.5C11.5 4.7 11.5 5.1 11.3 5.3L8.8 7.8L11.3 10.3Z" />
              </svg>
            </button>
          )}
          {suffix && <span className={styles.suffix}>{suffix}</span>}
        </div>
        {error && errorMessage && <div className={styles.errorMessage}>{errorMessage}</div>}
      </div>
    );
  },
);

Input.displayName = 'Input';

// 密码输入框
export interface PasswordInputProps extends Omit<InputProps, 'type' | 'suffix'> {}

export const PasswordInput = forwardRef<HTMLInputElement, PasswordInputProps>((props, ref) => {
  const [visible, setVisible] = useState(false);

  const toggleVisible = () => {
    setVisible(!visible);
  };

  const eyeIcon = (
    <button
      type="button"
      onClick={toggleVisible}
      className={styles.eyeBtn}
      tabIndex={-1}
      aria-label={visible ? '隐藏密码' : '显示密码'}
    >
      {visible ? (
        <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor">
          <path
            d="M1 12C1 12 5 4 12 4C19 4 23 12 23 12C23 12 19 20 12 20C5 20 1 12 1 12Z"
            strokeWidth="2"
            strokeLinecap="round"
            strokeLinejoin="round"
          />
          <circle cx="12" cy="12" r="3" strokeWidth="2" />
        </svg>
      ) : (
        <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor">
          <path
            d="M17.94 17.94C16.2306 19.243 14.1491 19.9649 12 20C5 20 1 12 1 12C2.24389 9.68192 3.96914 7.65663 6.06 6.06M9.9 4.24C10.5883 4.0789 11.2931 3.99836 12 4C19 4 23 12 23 12C22.393 13.1356 21.6691 14.2047 20.84 15.19M14.12 14.12C13.8454 14.4147 13.5141 14.6512 13.1462 14.8151C12.7782 14.9791 12.3809 15.0673 11.9781 15.0744C11.5753 15.0815 11.1752 15.0074 10.8016 14.8565C10.4281 14.7056 10.0887 14.4811 9.80385 14.1962C9.51897 13.9113 9.29439 13.572 9.14351 13.1984C8.99262 12.8248 8.91853 12.4247 8.92563 12.0219C8.93274 11.6191 9.02091 11.2218 9.18488 10.8538C9.34884 10.4858 9.58525 10.1546 9.88 9.88"
            strokeWidth="2"
            strokeLinecap="round"
            strokeLinejoin="round"
          />
          <line x1="1" y1="1" x2="23" y2="23" strokeWidth="2" strokeLinecap="round" />
        </svg>
      )}
    </button>
  );

  return <Input ref={ref} type={visible ? 'text' : 'password'} suffix={eyeIcon} {...props} />;
});

PasswordInput.displayName = 'PasswordInput';
