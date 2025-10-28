import { useState, FormEvent } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { Button, Input, PasswordInput, Form, FormItem } from 'components';
import { useAuth } from 'contexts/AuthContext';
import styles from './Login.module.less';

// 用户图标
const UserIcon = () => (
  <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor">
    <path
      d="M20 21V19C20 17.9391 19.5786 16.9217 18.8284 16.1716C18.0783 15.4214 17.0609 15 16 15H8C6.93913 15 5.92172 15.4214 5.17157 16.1716C4.42143 16.9217 4 17.9391 4 19V21"
      strokeWidth="2"
      strokeLinecap="round"
      strokeLinejoin="round"
    />
    <circle cx="12" cy="7" r="4" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" />
  </svg>
);

// 锁图标
const LockIcon = () => (
  <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor">
    <rect
      x="3"
      y="11"
      width="18"
      height="11"
      rx="2"
      ry="2"
      strokeWidth="2"
      strokeLinecap="round"
      strokeLinejoin="round"
    />
    <path
      d="M7 11V7C7 5.67392 7.52678 4.40215 8.46447 3.46447C9.40215 2.52678 10.6739 2 12 2C13.3261 2 14.5979 2.52678 15.5355 3.46447C16.4732 4.40215 17 5.67392 17 7V11"
      strokeWidth="2"
      strokeLinecap="round"
      strokeLinejoin="round"
    />
  </svg>
);

interface FormValues {
  username: string;
  password: string;
}

interface FormErrors {
  username?: string;
  password?: string;
}

export const Login = () => {
  const navigate = useNavigate();
  const { login } = useAuth();
  const [formValues, setFormValues] = useState<FormValues>({
    username: '',
    password: '',
  });
  const [errors, setErrors] = useState<FormErrors>({});
  const [loading, setLoading] = useState(false);

  const validateForm = (): boolean => {
    const newErrors: FormErrors = {};

    if (!formValues.username) {
      newErrors.username = '请输入用户名';
    } else if (formValues.username.length < 3) {
      newErrors.username = '用户名至少3个字符';
    }

    if (!formValues.password) {
      newErrors.password = '请输入密码';
    } else if (formValues.password.length < 6) {
      newErrors.password = '密码至少6个字符';
    }

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleSubmit = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();

    if (!validateForm()) {
      return;
    }

    setLoading(true);
    try {
      await login(formValues.username, formValues.password);
      // 登录成功，跳转到仪表盘
      navigate('/dashboard');
    } catch (error: any) {
      setErrors({
        password: error.message || '登录失败，请检查用户名和密码',
      });
    } finally {
      setLoading(false);
    }
  };

  const handleInputChange =
    (field: keyof FormValues) => (e: React.ChangeEvent<HTMLInputElement>) => {
      setFormValues({
        ...formValues,
        [field]: e.target.value,
      });
      // 清除该字段的错误
      if (errors[field]) {
        setErrors({
          ...errors,
          [field]: undefined,
        });
      }
    };

  return (
    <div className={styles.loginContainer}>
      {/* 动态背景装饰 */}
      <div className={styles.bgDecoration}>
        <div className={styles.circle1}></div>
        <div className={styles.circle2}></div>
        <div className={styles.circle3}></div>
      </div>

      {/* 登录卡片 */}
      <div className={styles.loginCard}>
        {/* 头部 */}
        <div className={styles.header}>
          <h1 className={styles.title}>GameLink</h1>
          <p className={styles.subtitle}>欢迎回来</p>
        </div>

        {/* 登录表单 */}
        <Form onSubmit={handleSubmit} className={styles.form}>
          <FormItem label="" error={errors.username}>
            <Input
              size="large"
              prefix={<UserIcon />}
              placeholder="用户名"
              value={formValues.username}
              onChange={handleInputChange('username')}
              allowClear
            />
          </FormItem>

          <FormItem label="" error={errors.password}>
            <PasswordInput
              size="large"
              prefix={<LockIcon />}
              placeholder="密码"
              value={formValues.password}
              onChange={handleInputChange('password')}
              allowClear
            />
          </FormItem>

          <FormItem label="">
            <Button type="submit" variant="primary" size="large" block loading={loading}>
              登录
            </Button>
          </FormItem>
        </Form>

        {/* 底部 */}
        <div className={styles.footer}>
          <span className={styles.footerText}>还没有账号？</span>
          <Link to="/register" className={styles.link}>
            立即注册
          </Link>
        </div>
      </div>

      {/* 版权信息 */}
      <div className={styles.copyright}>© 2025 GameLink. All rights reserved.</div>
    </div>
  );
};
