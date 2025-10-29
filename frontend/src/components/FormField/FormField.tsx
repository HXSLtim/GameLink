import React, { ReactNode } from 'react';
import styles from './FormField.module.less';

export interface FormFieldProps {
  label: string;
  required?: boolean;
  error?: string;
  children: ReactNode;
  className?: string;
}

/**
 * 通用表单字段组件
 * 包含标签、必填标识、错误提示
 */
export const FormField: React.FC<FormFieldProps> = ({
  label,
  required = false,
  error,
  children,
  className = '',
}) => {
  return (
    <div className={`${styles.formField} ${className}`}>
      <label className={styles.label}>
        {label}
        {required && <span className={styles.required}>*</span>}
      </label>
      {children}
      {error && <div className={styles.error}>{error}</div>}
    </div>
  );
};
