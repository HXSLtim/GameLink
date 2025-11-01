import React, { useEffect, useState } from 'react';
import { Button, Card, Input, Space, Tag, RouteCache } from 'components';
import { Outlet, useLocation, useNavigate } from 'react-router-dom';

/**
 * RouteCache 演示主页面
 * - 提供切换缓存开关与子页导航
 * - 使用 RouteCache 包裹 Outlet，使子页在路由跳转间保持状态
 */
export const CacheDemo: React.FC = () => {
  const [enabled, setEnabled] = useState<boolean>(true);
  const navigate = useNavigate();
  const location = useLocation();

  const go = (sub: 'a' | 'b') => navigate(`/cache-demo/${sub}`);

  return (
    <div style={{ padding: 24 }}>
      <Space size={12} direction="vertical" style={{ width: '100%' }}>
        <Card title="RouteCache 演示" extra={<Tag color={enabled ? 'green' : 'red'}>{enabled ? '缓存开启' : '缓存关闭'}</Tag>}>
          <Space size={12} wrap>
            <Button onClick={() => setEnabled((v) => !v)}>
              {enabled ? '关闭缓存' : '开启缓存'}
            </Button>
            <Button onClick={() => go('a')} disabled={location.pathname === '/cache-demo/a'}>
              跳转到页面 A
            </Button>
            <Button onClick={() => go('b')} disabled={location.pathname === '/cache-demo/b'}>
              跳转到页面 B
            </Button>
          </Space>
          <div style={{ marginTop: 12, color: 'var(--text-secondary)' }}>
            在开启缓存时，页面内的计时器与输入内容会在 A/B 跳转间保持；关闭缓存将导致子页重新挂载。
          </div>
        </Card>

        <RouteCache enabled={enabled} maxCache={5} cacheRoutes={["/cache-demo/a", "/cache-demo/b"]}>
          <div style={{ marginTop: 8 }}>
            <Outlet />
          </div>
        </RouteCache>
      </Space>
    </div>
  );
};

/**
 * 展示状态持久化的通用子组件
 */
const StatefulBox: React.FC<{ label: string }> = ({ label }) => {
  const [text, setText] = useState<string>('');
  const [counter, setCounter] = useState<number>(0);
  const [mountedAt] = useState<string>(() => new Date().toLocaleTimeString());

  useEffect(() => {
    const t = setInterval(() => setCounter((c) => c + 1), 1000);
    return () => clearInterval(t);
  }, []);

  return (
    <Card title={`${label}（挂载时间：${mountedAt}）`}>
      <Space size={12} direction="vertical" style={{ width: '100%' }}>
        <div>
          计时器：<Tag color="blue">{counter}s</Tag>
        </div>
        <Input
          value={text}
          placeholder={`在 ${label} 输入一些内容，然后切换到另一页再切回来，观察是否保留。`}
          onChange={(e) => setText(e.target.value)}
        />
        <div>
          当前输入镜像：<Tag>{text || '（空）'}</Tag>
        </div>
      </Space>
    </Card>
  );
};

export const CachePageA: React.FC = () => {
  return <StatefulBox label="页面 A" />;
};

export const CachePageB: React.FC = () => {
  return <StatefulBox label="页面 B" />;
};