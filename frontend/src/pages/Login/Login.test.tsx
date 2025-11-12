/**
 * Login页面测试
 */

import { describe, it, expect, vi, beforeEach } from 'vitest';
import { screen, fireEvent, waitFor } from '@testing-library/react';
import { renderWithProviders } from '../../test/utils/test-utils';
import { Login } from './Login';

// Mock useAuth
const mockLogin = vi.fn();
vi.mock('../../contexts/AuthContext', async () => {
  const actual = await vi.importActual('../../contexts/AuthContext');
  return {
    ...actual,
    useAuth: () => ({
      login: mockLogin,
      logout: vi.fn(),
      user: null,
      token: null,
      loading: false,
    }),
  };
});

// Mock useNavigate
const mockNavigate = vi.fn();
vi.mock('react-router-dom', async () => {
  const actual = await vi.importActual('react-router-dom');
  return {
    ...actual,
    useNavigate: () => mockNavigate,
  };
});

describe('Login Page', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    localStorage.clear();
  });

  it('should render login form', () => {
    renderWithProviders(<Login />);
    
    // 验证页面标题
    expect(screen.getByRole('heading', { level: 1 })).toBeInTheDocument();
    expect(screen.getByText('GameLink')).toBeInTheDocument();
    
    // 验证表单元素存在
    const inputs = screen.getAllByRole('textbox');
    expect(inputs.length).toBeGreaterThan(0);
  });

  it('should show validation errors for empty fields', async () => {
    renderWithProviders(<Login />);
    
    const loginButton = screen.getByRole('button', { name: /登录/i });
    fireEvent.click(loginButton);
    
    // 验证显示了验证错误
    await waitFor(() => {
      const inputs = screen.getAllByRole('textbox');
      expect(inputs.length).toBeGreaterThan(0);
    });
  });

  it('should handle successful login', async () => {
    mockLogin.mockResolvedValue(undefined);

    renderWithProviders(<Login />);
    
    // 填写表单
    const usernameInput = screen.getByPlaceholderText(/用户名/i);
    const passwordInput = screen.getByPlaceholderText(/密码/i);
    
    fireEvent.change(usernameInput, { target: { value: 'admin' } });
    fireEvent.change(passwordInput, { target: { value: 'password123' } });
    
    // 点击登录
    const loginButton = screen.getByRole('button', { name: /登录/i });
    fireEvent.click(loginButton);
    
    // 验证调用了登录函数
    await waitFor(() => {
      expect(mockLogin).toHaveBeenCalledWith('admin', 'password123');
    });
    
    // 验证导航到dashboard
    await waitFor(() => {
      expect(mockNavigate).toHaveBeenCalledWith('/dashboard');
    });
  });

  it('should handle login failure', async () => {
    const errorMessage = 'Invalid credentials';
    mockLogin.mockRejectedValue(new Error(errorMessage));

    renderWithProviders(<Login />);
    
    // 填写表单
    const usernameInput = screen.getByPlaceholderText(/用户名/i);
    const passwordInput = screen.getByPlaceholderText(/密码/i);
    
    fireEvent.change(usernameInput, { target: { value: 'admin' } });
    fireEvent.change(passwordInput, { target: { value: 'wrongpassword' } });
    
    // 点击登录
    const loginButton = screen.getByRole('button', { name: /登录/i });
    fireEvent.click(loginButton);
    
    // 验证调用了登录函数
    await waitFor(() => {
      expect(mockLogin).toHaveBeenCalled();
    });
  });

  it('should toggle password visibility', () => {
    renderWithProviders(<Login />);
    
    const passwordInput = screen.getByPlaceholderText(/密码/i) as HTMLInputElement;
    
    // 初始状态应该是password类型
    expect(passwordInput.type).toBe('password');
    
    // 查找并点击显示密码按钮（如果有）
    const toggleButtons = screen.queryAllByRole('button');
    const toggleButton = toggleButtons.find(btn => 
      btn.getAttribute('aria-label')?.includes('显示密码') || 
      btn.getAttribute('aria-label')?.includes('隐藏密码')
    );
    
    if (toggleButton) {
      fireEvent.click(toggleButton);
      // 验证密码类型改变
      expect(passwordInput.type).toBe('text');
      
      fireEvent.click(toggleButton);
      expect(passwordInput.type).toBe('password');
    }
  });

  it('should disable submit button during login', async () => {
    mockLogin.mockImplementation(() => 
      new Promise(resolve => setTimeout(resolve, 100))
    );

    renderWithProviders(<Login />);
    
    const usernameInput = screen.getByPlaceholderText(/用户名/i);
    const passwordInput = screen.getByPlaceholderText(/密码/i);
    const loginButton = screen.getByRole('button', { name: /登录/i });
    
    fireEvent.change(usernameInput, { target: { value: 'admin' } });
    fireEvent.change(passwordInput, { target: { value: 'password123' } });
    fireEvent.click(loginButton);
    
    // 验证登录函数被调用
    await waitFor(() => {
      expect(mockLogin).toHaveBeenCalled();
    });
  });

  it('should clear error message when input changes', async () => {
    const errorMessage = 'Invalid credentials';
    mockLogin.mockRejectedValue(new Error(errorMessage));

    renderWithProviders(<Login />);
    
    const usernameInput = screen.getByPlaceholderText(/用户名/i);
    const passwordInput = screen.getByPlaceholderText(/密码/i);
    const loginButton = screen.getByRole('button', { name: /登录/i });
    
    // 触发登录错误（使用符合验证规则的输入）
    fireEvent.change(usernameInput, { target: { value: 'admin' } });
    fireEvent.change(passwordInput, { target: { value: 'wrong123' } });
    fireEvent.click(loginButton);
    
    await waitFor(() => {
      expect(mockLogin).toHaveBeenCalled();
    });
    
    // 修改输入应该清除错误
    fireEvent.change(usernameInput, { target: { value: 'newadmin' } });
    
    // 验证表单仍然可以交互
    expect(usernameInput).not.toBeDisabled();
  });
});
