import { useEffect, useRef } from 'react';
import { BrowserRouter } from 'react-router-dom';
import { App as AntdApp } from 'antd';
import AppRouter from '@/router';
import { AdminProvider } from '@/context/AdminContext';
import { ThemeProvider } from '@/context/ThemeContext';
import { smartInit } from '@/services/init';
import { logger } from '@/utils/logger';
import './App.css';

/**
 * 应用根组件
 */
function App() {
  /**
   * 应用初始化
   * 自动同步路由和权限到后端，为超管分配所有权限
   */
  const initialized = useRef(false);

  useEffect(() => {
    if (initialized.current) return;
    initialized.current = true;

    const runInit = async () => {
      try {
        const result = await smartInit({
          syncMenus: true,
          syncPermissions: true,
          assignSuperAdminPermissions: true,
          verbose: import.meta.env.DEV,
        });

        if (result && import.meta.env.DEV) {
          if (result.success) {
            logger.info('[App] 初始化成功:', {
              菜单: result.menuSync ? `创建${result.menuSync.created}个，更新${result.menuSync.updated}个` : '跳过',
              权限: result.permissionSync ? `创建${result.permissionSync.created}个，更新${result.permissionSync.updated}个` : '跳过',
              超管权限: result.superAdminAssign?.message || '跳过',
              耗时: `${result.duration}ms`,
            });
          } else {
            logger.warn('[App] 初始化存在错误:', result.errors);
          }
        }
      } catch (error) {
        logger.error('[App] 初始化失败:', error);
      }
    };

    runInit();
  }, []);

  return (
    <ThemeProvider>
      <AntdApp>
        <AdminProvider>
          <BrowserRouter>
            <AppRouter />
          </BrowserRouter>
        </AdminProvider>
      </AntdApp>
    </ThemeProvider>
  );
}

export default App;
