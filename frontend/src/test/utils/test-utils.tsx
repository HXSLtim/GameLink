/**
 * 测试工具函数
 * 提供常用的测试辅助功能
 */

import { render, RenderOptions } from '@testing-library/react';
import { ReactElement, ReactNode } from 'react';
import { MemoryRouter } from 'react-router-dom';
import { AuthProvider } from '../../contexts/AuthContext';
import { ThemeProvider } from '../../contexts/ThemeContext';
import { UserRole, UserStatus } from '../../types/user';
import type { CurrentUser } from '../../types/auth';

/**
 * Mock用户数据
 */
export const mockUsers = {
  admin: {
    id: 1,
    name: 'Admin User',
    username: 'admin',
    email: 'admin@gamelink.com',
    role: UserRole.ADMIN,
    status: UserStatus.ACTIVE,
    createdAt: '2024-01-01T00:00:00Z',
    updatedAt: '2024-01-01T00:00:00Z',
  } as CurrentUser,
  
  player: {
    id: 2,
    name: 'Player User',
    username: 'player',
    email: 'player@gamelink.com',
    role: UserRole.PLAYER,
    status: UserStatus.ACTIVE,
    createdAt: '2024-01-01T00:00:00Z',
    updatedAt: '2024-01-01T00:00:00Z',
  } as CurrentUser,
  
  user: {
    id: 3,
    name: 'Normal User',
    username: 'user',
    email: 'user@gamelink.com',
    role: UserRole.USER,
    status: UserStatus.ACTIVE,
    createdAt: '2024-01-01T00:00:00Z',
    updatedAt: '2024-01-01T00:00:00Z',
  } as CurrentUser,
};

/**
 * 所有Providers的包装器选项
 */
interface AllProvidersProps {
  children: ReactNode;
  initialEntries?: string[];
  initialIndex?: number;
}

/**
 * 包含所有必要Providers的包装组件
 */
export function AllProviders({ children, initialEntries = ['/'], initialIndex = 0 }: AllProvidersProps) {
  return (
    <MemoryRouter initialEntries={initialEntries} initialIndex={initialIndex}>
      <ThemeProvider>
        <AuthProvider>
          {children}
        </AuthProvider>
      </ThemeProvider>
    </MemoryRouter>
  );
}

/**
 * 自定义渲染选项
 */
interface CustomRenderOptions extends Omit<RenderOptions, 'wrapper'> {
  route?: string;
  initialEntries?: string[];
  initialIndex?: number;
  withRouter?: boolean;
  withAuth?: boolean;
  withTheme?: boolean;
}

/**
 * 自定义渲染函数，自动包含常用Providers
 */
export function renderWithProviders(
  ui: ReactElement,
  {
    route = '/',
    initialEntries,
    initialIndex = 0,
    withRouter = true,
    withAuth = true,
    withTheme = true,
    ...renderOptions
  }: CustomRenderOptions = {}
) {
  const entries = initialEntries || [route];
  
  let Wrapper = ({ children }: { children: ReactNode }) => <>{children}</>;
  
  // 按需添加Providers
  if (withRouter) {
    const RouterWrapper = Wrapper;
    Wrapper = ({ children }) => (
      <MemoryRouter initialEntries={entries} initialIndex={initialIndex}>
        <RouterWrapper>{children}</RouterWrapper>
      </MemoryRouter>
    );
  }
  
  if (withTheme) {
    const ThemeWrapper = Wrapper;
    Wrapper = ({ children }) => (
      <ThemeProvider>
        <ThemeWrapper>{children}</ThemeWrapper>
      </ThemeProvider>
    );
  }
  
  if (withAuth) {
    const AuthWrapper = Wrapper;
    Wrapper = ({ children }) => (
      <AuthProvider>
        <AuthWrapper>{children}</AuthWrapper>
      </AuthProvider>
    );
  }
  
  return render(ui, { wrapper: Wrapper, ...renderOptions });
}

/**
 * 等待异步操作
 */
export const waitFor = (ms: number) => new Promise(resolve => setTimeout(resolve, ms));

/**
 * 创建Mock API响应
 */
export function createMockApiResponse<T>(data: T, delay = 0) {
  return new Promise<T>((resolve) => {
    setTimeout(() => resolve(data), delay);
  });
}

/**
 * 创建Mock API错误
 */
export function createMockApiError(message: string, code?: string, delay = 0) {
  return new Promise((_, reject) => {
    setTimeout(() => {
      const error = new Error(message) as any;
      error.code = code;
      reject(error);
    }, delay);
  });
}

/**
 * Mock localStorage
 */
export function mockLocalStorage() {
  const storage: Record<string, string> = {};
  
  return {
    getItem: (key: string) => storage[key] || null,
    setItem: (key: string, value: string) => { storage[key] = value; },
    removeItem: (key: string) => { delete storage[key]; },
    clear: () => { Object.keys(storage).forEach(key => delete storage[key]); },
    get length() { return Object.keys(storage).length; },
    key: (index: number) => Object.keys(storage)[index] || null,
  };
}

/**
 * 常用的测试数据
 */
export const testData = {
  // 列表数据
  pagination: {
    page: 1,
    pageSize: 20,
    total: 100,
  },
  
  // 日期
  dates: {
    now: '2024-11-11T00:00:00Z',
    yesterday: '2024-11-10T00:00:00Z',
    lastWeek: '2024-11-04T00:00:00Z',
    lastMonth: '2024-10-11T00:00:00Z',
  },
  
  // 订单状态
  orderStatuses: ['pending', 'confirmed', 'in_progress', 'completed', 'cancelled'],
  
  // 支付状态
  paymentStatuses: ['pending', 'success', 'failed', 'refunded'],
};

/**
 * 导出所有工具
 */
export * from '@testing-library/react';
export { renderWithProviders as render };
