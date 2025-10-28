import { useState, FormEvent } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { Button, Input, PasswordInput, Form, FormItem } from 'components';
import styles from './Register.module.less';

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

// 邮箱图标
const EmailIcon = () => (
  <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor">
    <path
      d="M4 4H20C21.1 4 22 4.9 22 6V18C22 19.1 21.1 20 20 20H4C2.9 20 2 19.1 2 18V6C2 4.9 2.9 4 4 4Z"
      strokeWidth="2"
      strokeLinecap="round"
      strokeLinejoin="round"
    />
    <path d="M22 6L12 13L2 6" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" />
  </svg>
);

// 手机图标
const PhoneIcon = () => (
  <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor">
    <path
      d="M22 16.92V19.92C22.0011 20.1985 21.9441 20.4742 21.8325 20.7293C21.7209 20.9845 21.5573 21.2136 21.3521 21.4019C21.1469 21.5901 20.9046 21.7335 20.6407 21.8227C20.3769 21.9119 20.0974 21.9451 19.82 21.92C16.7428 21.5856 13.787 20.5341 11.19 18.85C8.77382 17.3147 6.72533 15.2662 5.18999 12.85C3.49997 10.2412 2.44824 7.27099 2.11999 4.18C2.095 3.90347 2.12787 3.62476 2.21649 3.36162C2.30512 3.09849 2.44756 2.85669 2.63476 2.65162C2.82196 2.44655 3.0498 2.28271 3.30379 2.17052C3.55777 2.05833 3.83233 2.00026 4.10999 2H7.10999C7.5953 1.99522 8.06579 2.16708 8.43376 2.48353C8.80173 2.79999 9.04207 3.23945 9.10999 3.72C9.23662 4.68007 9.47144 5.62273 9.80999 6.53C9.94454 6.88792 9.97366 7.27691 9.8939 7.65088C9.81415 8.02485 9.62886 8.36811 9.35999 8.64L8.08999 9.91C9.51355 12.4135 11.5864 14.4864 14.09 15.91L15.36 14.64C15.6319 14.3711 15.9751 14.1858 16.3491 14.1061C16.7231 14.0263 17.1121 14.0555 17.47 14.19C18.3773 14.5286 19.3199 14.7634 20.28 14.89C20.7658 14.9585 21.2094 15.2032 21.5265 15.5775C21.8437 15.9518 22.0122 16.4296 22 16.92Z"
      strokeWidth="2"
      strokeLinecap="round"
      strokeLinejoin="round"
    />
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
  email: string;
  name: string;
  phone: string;
  password: string;
  confirmPassword: string;
}

interface FormErrors {
  email?: string;
  name?: string;
  phone?: string;
  password?: string;
  confirmPassword?: string;
}

export const Register = () => {
  const navigate = useNavigate();
  const [formValues, setFormValues] = useState<FormValues>({
    email: '',
    name: '',
    phone: '',
    password: '',
    confirmPassword: '',
  });
  const [errors, setErrors] = useState<FormErrors>({});
  const [loading, setLoading] = useState(false);

  const validateForm = (): boolean => {
    const newErrors: FormErrors = {};

    // 邮箱验证
    if (!formValues.email) {
      newErrors.email = '请输入邮箱';
    } else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(formValues.email)) {
      newErrors.email = '请输入有效的邮箱地址';
    }

    // 姓名验证
    if (!formValues.name) {
      newErrors.name = '请输入姓名';
    } else if (formValues.name.length < 2) {
      newErrors.name = '姓名至少2个字符';
    }

    // 手机号验证（可选）
    if (formValues.phone && !/^1[3-9]\d{9}$/.test(formValues.phone)) {
      newErrors.phone = '请输入有效的手机号';
    }

    // 密码验证
    if (!formValues.password) {
      newErrors.password = '请输入密码';
    } else if (formValues.password.length < 6) {
      newErrors.password = '密码至少6个字符';
    }

    // 确认密码验证
    if (!formValues.confirmPassword) {
      newErrors.confirmPassword = '请确认密码';
    } else if (formValues.password !== formValues.confirmPassword) {
      newErrors.confirmPassword = '两次密码不一致';
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
      // 调用注册 API
      const response = await fetch('/api/v1/auth/register', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          email: formValues.email,
          name: formValues.name,
          phone: formValues.phone || undefined, // 可选字段，空值不传
          password: formValues.password,
        }),
      });

      const data = await response.json();

      if (data.success) {
        // 注册成功，跳转到登录页
        navigate('/login');
      } else {
        setErrors({
          confirmPassword: data.message || '注册失败，请稍后重试',
        });
      }
    } catch (error: any) {
      setErrors({
        confirmPassword: error.message || '注册失败，请检查网络连接',
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
    <div className={styles.registerContainer}>
      {/* 动态背景装饰 */}
      <div className={styles.bgDecoration}>
        <div className={styles.circle1}></div>
        <div className={styles.circle2}></div>
        <div className={styles.circle3}></div>
      </div>

      {/* 注册卡片 */}
      <div className={styles.registerCard}>
        {/* 头部 */}
        <div className={styles.header}>
          <h1 className={styles.title}>GameLink</h1>
          <p className={styles.subtitle}>创建新账号</p>
        </div>

        {/* 注册表单 */}
        <Form onSubmit={handleSubmit} className={styles.form}>
          <FormItem label="" error={errors.email}>
            <Input
              size="large"
              type="email"
              prefix={<EmailIcon />}
              placeholder="邮箱"
              value={formValues.email}
              onChange={handleInputChange('email')}
              allowClear
            />
          </FormItem>

          <FormItem label="" error={errors.name}>
            <Input
              size="large"
              prefix={<UserIcon />}
              placeholder="姓名"
              value={formValues.name}
              onChange={handleInputChange('name')}
              allowClear
            />
          </FormItem>

          <FormItem label="" error={errors.phone}>
            <Input
              size="large"
              type="tel"
              prefix={<PhoneIcon />}
              placeholder="手机号（可选）"
              value={formValues.phone}
              onChange={handleInputChange('phone')}
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

          <FormItem label="" error={errors.confirmPassword}>
            <PasswordInput
              size="large"
              prefix={<LockIcon />}
              placeholder="确认密码"
              value={formValues.confirmPassword}
              onChange={handleInputChange('confirmPassword')}
              allowClear
            />
          </FormItem>

          <FormItem label="">
            <Button type="submit" variant="primary" size="large" block loading={loading}>
              注册
            </Button>
          </FormItem>
        </Form>

        {/* 底部 */}
        <div className={styles.footer}>
          <span className={styles.footerText}>已有账号？</span>
          <Link to="/login" className={styles.link}>
            立即登录
          </Link>
        </div>
      </div>

      {/* 版权信息 */}
      <div className={styles.copyright}>© 2025 GameLink. All rights reserved.</div>
    </div>
  );
};
