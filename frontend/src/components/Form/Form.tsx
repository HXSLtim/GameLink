import { ReactNode, FormHTMLAttributes, createContext, useContext } from 'react';
import styles from './Form.module.less';

export interface FormProps extends FormHTMLAttributes<HTMLFormElement> {
  /** 表单内容 */
  children: ReactNode;
  /** 表单布局 */
  layout?: 'horizontal' | 'vertical';
  /** 自定义类名 */
  className?: string;
}

interface FormContextValue {
  layout: 'horizontal' | 'vertical';
}

const FormContext = createContext<FormContextValue>({ layout: 'vertical' });

export const Form: React.FC<FormProps> = ({
  children,
  layout = 'vertical',
  className = '',
  ...rest
}) => {
  const classNames = [styles.form, styles[layout], className].filter(Boolean).join(' ');

  return (
    <FormContext.Provider value={{ layout }}>
      <form className={classNames} {...rest}>
        {children}
      </form>
    </FormContext.Provider>
  );
};

// FormItem 组件
export interface FormItemProps {
  /** 表单项内容 */
  children: ReactNode;
  /** 标签 */
  label?: ReactNode;
  /** 字段名 */
  name?: string;
  /** 是否必填 */
  required?: boolean;
  /** 错误信息 */
  error?: string;
  /** 提示信息 */
  help?: string;
  /** 自定义类名 */
  className?: string;
}

export const FormItem: React.FC<FormItemProps> = ({
  children,
  label,
  required = false,
  error,
  help,
  className = '',
}) => {
  const { layout } = useContext(FormContext);

  const classNames = [
    styles.formItem,
    styles[`${layout}Item`],
    error ? styles.error : '',
    className,
  ]
    .filter(Boolean)
    .join(' ');

  return (
    <div className={classNames}>
      {label && (
        <label className={styles.label}>
          {required && <span className={styles.required}>*</span>}
          {label}
        </label>
      )}
      <div className={styles.control}>{children}</div>
      {error && <div className={styles.errorMessage}>{error}</div>}
      {!error && help && <div className={styles.helpMessage}>{help}</div>}
    </div>
  );
};
