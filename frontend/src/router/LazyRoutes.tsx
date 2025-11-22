/**
 * 路由懒加载组件
 */

import { Spin } from '@arco-design/web-react';

/**
 * 路由加载中组件
 */
export const RouteLoading = () => {
  return (
    <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '100vh' }}>
      <Spin size="large" />
    </div>
  );
};

/**
 * 页面加载中组件
 */
export const PageLoading = () => {
  return (
    <div style={{ padding: '40px', textAlign: 'center' }}>
      <Spin />
      <div style={{ marginTop: '16px' }}>加载中...</div>
    </div>
  );
};