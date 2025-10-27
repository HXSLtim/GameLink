import { render, screen, fireEvent, waitFor } from '../test/testUtils';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { Login } from './Login';
import * as authService from '../services/auth';

// Mock auth service
vi.mock('../services/auth', () => ({
  authService: {
    login: vi.fn(),
    me: vi.fn(),
  },
}));

// Mock useNavigate and useLocation
const mockNavigate = vi.fn();
vi.mock('react-router-dom', async () => {
  const actual = await vi.importActual('react-router-dom');
  return {
    ...actual,
    useNavigate: () => mockNavigate,
    useLocation: () => ({ state: null }),
  };
});

describe('Login', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('should render login form', () => {
    render(<Login />);

    expect(screen.getByText('GameLink 管理系统')).toBeInTheDocument();
    expect(screen.getByText('欢迎回来，请登录您的账户')).toBeInTheDocument();
    expect(screen.getByPlaceholderText('请输入用户名')).toBeInTheDocument();
    expect(screen.getByPlaceholderText('请输入密码')).toBeInTheDocument();
    expect(screen.getByRole('button', { name: /立即登录/i })).toBeInTheDocument();
  });

  it('should display validation errors for empty fields', async () => {
    render(<Login />);

    // Clear form first
    const usernameInput = screen.getByDisplayValue('admin');
    const passwordInput = screen.getByDisplayValue('admin123');

    fireEvent.change(usernameInput, { target: { value: '' } });
    fireEvent.change(passwordInput, { target: { value: '' } });

    const loginButton = screen.getByRole('button', { name: /立即登录/i });
    fireEvent.click(loginButton);

    await waitFor(() => {
      expect(screen.getByText('请输入用户名')).toBeInTheDocument();
      expect(screen.getByText('请输入密码')).toBeInTheDocument();
    });
  });

  it('should display validation error for short username', async () => {
    render(<Login />);

    const usernameInput = screen.getByPlaceholderText('请输入用户名');
    fireEvent.change(usernameInput, { target: { value: 'ab' } });

    const loginButton = screen.getByRole('button', { name: /立即登录/i });
    fireEvent.click(loginButton);

    await waitFor(() => {
      expect(screen.getByText('用户名至少3个字符')).toBeInTheDocument();
    });
  });

  it('should display validation error for short password', async () => {
    render(<Login />);

    const usernameInput = screen.getByPlaceholderText('请输入用户名');
    const passwordInput = screen.getByPlaceholderText('请输入密码');

    fireEvent.change(usernameInput, { target: { value: 'admin' } });
    fireEvent.change(passwordInput, { target: { value: '12345' } });

    const loginButton = screen.getByRole('button', { name: /立即登录/i });
    fireEvent.click(loginButton);

    await waitFor(() => {
      expect(screen.getByText('密码至少6个字符')).toBeInTheDocument();
    });
  });

  it('should call login service on form submit', async () => {
    const mockLogin = vi.mocked(authService.authService.login);
    mockLogin.mockResolvedValue({ token: 'test-token' });

    render(<Login />);

    const usernameInput = screen.getByPlaceholderText('请输入用户名');
    const passwordInput = screen.getByPlaceholderText('请输入密码');

    fireEvent.change(usernameInput, { target: { value: 'testuser' } });
    fireEvent.change(passwordInput, { target: { value: 'testpass123' } });

    const loginButton = screen.getByRole('button', { name: /立即登录/i });
    fireEvent.click(loginButton);

    await waitFor(() => {
      expect(mockLogin).toHaveBeenCalledWith({
        username: 'testuser',
        password: 'testpass123',
        remember: true,
      });
    });
  });

  it('should show remember me checkbox', () => {
    render(<Login />);

    const rememberCheckbox = screen.getByRole('checkbox', { name: /记住我/i });
    expect(rememberCheckbox).toBeInTheDocument();
    expect(rememberCheckbox).toBeChecked(); // Default is true
  });

  it('should show forgot password link', () => {
    render(<Login />);

    expect(screen.getByText('忘记密码？')).toBeInTheDocument();
  });

  it('should display dev environment hint', () => {
    render(<Login />);

    expect(screen.getByText(/演示账号/i)).toBeInTheDocument();
    expect(screen.getByText(/admin123/i)).toBeInTheDocument();
  });
});


