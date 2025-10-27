import { useState, useCallback } from 'react';
import {
  Button,
  Card,
  Form,
  Input,
  Message,
  Typography,
  Checkbox,
  Divider,
  Link,
} from '@arco-design/web-react';
import {
  IconUser,
  IconLock,
  IconApps,
  IconRight,
} from '@arco-design/web-react/icon';
import type { Location } from 'react-router-dom';
import { useNavigate, useLocation } from 'react-router-dom';
import { authService } from '../services/auth';
import { useAuth } from '../contexts/AuthContext';
import styles from './Login.module.less';

interface LoginFormValues {
  username: string;
  password: string;
  remember?: boolean;
}

interface LocationState {
  from?: Location;
}

const INITIAL_VALUES: LoginFormValues = {
  username: 'admin',
  password: 'admin123',
  remember: true,
};

/**
 * Login page component
 *
 * @component
 * @description Modern login interface with enhanced UX and visual design
 */
export const Login = () => {
  const [form] = Form.useForm<LoginFormValues>();
  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();
  const location = useLocation() as Location<LocationState>;
  const { login } = useAuth();

  const handleSubmit = useCallback(async () => {
    setLoading(true);
    try {
      const values = await form.validate();
      const response = await authService.login(values);
      login(response.token);
      Message.success({
        content: '登录成功，欢迎回来！',
        duration: 2000,
      });

      const redirectPath = location.state?.from?.pathname || '/';
      // 延迟跳转，让用户看到成功提示
      setTimeout(() => {
        navigate(redirectPath, { replace: true });
      }, 500);
    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : '登录失败，请检查账号密码';
      Message.error(errorMessage);
      setLoading(false);
    }
  }, [form, login, location.state, navigate]);

  const handleKeyPress = useCallback(
    (event: React.KeyboardEvent) => {
      if (event.key === 'Enter') {
        handleSubmit();
      }
    },
    [handleSubmit],
  );

  return (
    <div className={styles.loginPage}>
      {/* 背景装饰 */}
      <div className={styles.background}>
        <div className={styles.gradientOrb1} />
        <div className={styles.gradientOrb2} />
        <div className={styles.gradientOrb3} />
      </div>

      {/* 登录卡片容器 */}
      <div className={styles.container}>
        {/* Logo 和标题区域 */}
        <div className={styles.header}>
          <div className={styles.logoContainer}>
            <IconApps className={styles.logo} />
          </div>
          <Typography.Title heading={3} className={styles.title}>
            GameLink 管理系统
          </Typography.Title>
          <Typography.Text className={styles.subtitle}>
            欢迎回来，请登录您的账户
          </Typography.Text>
        </div>

        {/* 登录表单卡片 */}
        <Card className={styles.card} bordered={false}>
          <Form
            form={form}
            layout="vertical"
            onSubmit={handleSubmit}
            autoComplete="off"
            initialValues={INITIAL_VALUES}
          >
            <Form.Item
              field="username"
              rules={[
                { required: true, message: '请输入用户名' },
                { minLength: 3, message: '用户名至少3个字符' },
              ]}
            >
              <Input
                prefix={<IconUser />}
                placeholder="请输入用户名"
                size="large"
                onKeyPress={handleKeyPress}
              />
            </Form.Item>

            <Form.Item
              field="password"
              rules={[
                { required: true, message: '请输入密码' },
                { minLength: 6, message: '密码至少6个字符' },
              ]}
            >
              <Input.Password
                prefix={<IconLock />}
                placeholder="请输入密码"
                size="large"
                onKeyPress={handleKeyPress}
              />
            </Form.Item>

            <Form.Item>
              <div className={styles.formOptions}>
                <Form.Item field="remember" noStyle>
                  <Checkbox>记住我</Checkbox>
                </Form.Item>
                <Link className={styles.forgotPassword}>忘记密码？</Link>
              </div>
            </Form.Item>

            <Form.Item>
              <Button
                type="primary"
                htmlType="submit"
                size="large"
                long
                loading={loading}
                className={styles.loginButton}
                icon={<IconRight />}
                iconOnly={false}
              >
                {loading ? '登录中...' : '立即登录'}
              </Button>
            </Form.Item>

            <Divider orientation="center" className={styles.divider}>
              开发环境
            </Divider>

            <div className={styles.devInfo}>
              <Typography.Text type="secondary" className={styles.hint}>
                <span className={styles.hintIcon}>💡</span>
                演示账号：<strong>admin</strong> / 密码：<strong>admin123</strong>
              </Typography.Text>
            </div>
          </Form>
        </Card>

        {/* 页脚信息 */}
        <div className={styles.footer}>
          <Typography.Text type="secondary" className={styles.footerText}>
            © 2024 GameLink. All rights reserved.
          </Typography.Text>
        </div>
      </div>
    </div>
  );
};
