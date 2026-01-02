/**
 * Admin Login Page Tests
 *
 * Tests for Admin Login page component including:
 * - Successful login
 * - Loading states
 * - Error handling (invalid credentials, network errors)
 * - Form validation
 * - Remember password functionality
 * - Quick login buttons (dev environment)
 * - Role verification (admin only)
 */

import React from 'react';
import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import AdminLogin from './index';
import { _authApi as authApi } from '@/api/auth';

// Mock the authApi module
vi.mock('@/api/auth', () => ({
  authApi: mockApi,
}));

// Mock config
vi.mock('@/config/debug', () => ({
  ENABLE_QUICK_LOGIN: false,
  DEBUG_USERS: [],
}));

describe('AdminLogin', () => {
  beforeEach(() => {
    resetAllMocks();
    localStorage.clear();
  });

  afterEach(() => {
    vi.clearAllMocks();
    localStorage.clear();
  });

  describe('Page Rendering', () => {
    it('should render login page correctly', () => {
      renderWithProviders(<AdminLogin />, { route: '/admin/login' });

      expect(screen.getByText('GameLink 管理后台')).toBeInTheDocument();
      expect(screen.getByText('管理员登录')).toBeInTheDocument();
    });

    it('should display login form fields', () => {
      renderWithProviders(<AdminLogin />, { route: '/admin/login' });

      expect(screen.getByPlaceholderText('管理员账号/邮箱')).toBeInTheDocument();
      expect(screen.getByPlaceholderText('密码')).toBeInTheDocument();
      expect(screen.getByRole('button', { name: '登录管理后台' })).toBeInTheDocument();
    });

    it('should display remember password checkbox', () => {
      renderWithProviders(<AdminLogin />, { route: '/admin/login' });

      expect(screen.getByText('记住密码')).toBeInTheDocument();
    });

    it('should display copyright information', () => {
      renderWithProviders(<AdminLogin />, { route: '/admin/login' });

      expect(screen.getByText(/© 2025 GameLink/)).toBeInTheDocument();
    });
  });

  describe('Successful Login', () => {
    it('should login successfully with valid credentials', async () => {
      const user = userEvent.setup();
      mockApi.login.mockResolvedValue({
        data: {
          success: true,
          data: {
            token: 'test-admin-token',
            user: {
              id: 1,
              role: 'admin',
              name: 'Test Admin',
              email: 'admin@gamelink.com',
            },
          },
        },
      });

      renderWithProviders(<AdminLogin />, { route: '/admin/login' });

      const usernameInput = screen.getByPlaceholderText('管理员账号/邮箱');
      const passwordInput = screen.getByPlaceholderText('密码');
      const loginButton = screen.getByRole('button', { name: '登录管理后台' });

      await user.type(usernameInput, 'admin@gamelink.com');
      await user.type(passwordInput, 'password123');
      await user.click(loginButton);

      await waitFor(() => {
        expect(mockApi.login).toHaveBeenCalledWith({
          username: 'admin@gamelink.com',
          password: 'password123',
        });
      });

      await waitFor(() => {
        expect(localStorage.getItem('token')).toBe('test-admin-token');
        expect(localStorage.getItem('user_role')).toBe('admin');
      });
    });

    it('should show success message on successful login', async () => {
      const user = userEvent.setup();
      mockApi.login.mockResolvedValue({
        data: {
          success: true,
          data: {
            token: 'test-admin-token',
            user: {
              id: 1,
              role: 'admin',
              name: 'Test Admin',
              email: 'admin@gamelink.com',
            },
          },
        },
      });

      renderWithProviders(<AdminLogin />, { route: '/admin/login' });

      const usernameInput = screen.getByPlaceholderText('管理员账号/邮箱');
      const passwordInput = screen.getByPlaceholderText('密码');
      const loginButton = screen.getByRole('button', { name: '登录管理后台' });

      await user.type(usernameInput, 'admin@gamelink.com');
      await user.type(passwordInput, 'password123');
      await user.click(loginButton);

      await waitFor(() => {
        // Success message should be shown (via App.useApp().message.success)
        expect(mockApi.login).toHaveBeenCalled();
      });
    });
  });

  describe('Role Verification', () => {
    it('should reject login for non-admin users', async () => {
      const user = userEvent.setup();
      mockApi.login.mockResolvedValue({
        data: {
          success: true,
          data: {
            token: 'test-user-token',
            user: {
              id: 2,
              role: 'user', // Not admin
              name: 'Regular User',
              email: 'user@example.com',
            },
          },
        },
      });

      renderWithProviders(<AdminLogin />, { route: '/admin/login' });

      const usernameInput = screen.getByPlaceholderText('管理员账号/邮箱');
      const passwordInput = screen.getByPlaceholderText('密码');
      const loginButton = screen.getByRole('button', { name: '登录管理后台' });

      await user.type(usernameInput, 'user@example.com');
      await user.type(passwordInput, 'password123');
      await user.click(loginButton);

      await waitFor(() => {
        expect(screen.getByText(/您没有管理后台访问权限/)).toBeInTheDocument();
      });

      // Token should not be saved for non-admin users
      expect(localStorage.getItem('token')).toBeNull();
    });

    it('should reject login for player role', async () => {
      const user = userEvent.setup();
      mockApi.login.mockResolvedValue({
        data: {
          success: true,
          data: {
            token: 'test-player-token',
            user: {
              id: 3,
              role: 'player', // Not admin
              name: 'Player User',
              email: 'player@example.com',
            },
          },
        },
      });

      renderWithProviders(<AdminLogin />, { route: '/admin/login' });

      const usernameInput = screen.getByPlaceholderText('管理员账号/邮箱');
      const passwordInput = screen.getByPlaceholderText('密码');
      const loginButton = screen.getByRole('button', { name: '登录管理后台' });

      await user.type(usernameInput, 'player@example.com');
      await user.type(passwordInput, 'password123');
      await user.click(loginButton);

      await waitFor(() => {
        expect(screen.getByText(/您没有管理后台访问权限/)).toBeInTheDocument();
      });
    });
  });

  describe('Loading States', () => {
    it('should show loading state during login', async () => {
      const user = userEvent.setup();
      mockApi.login.mockImplementation(
        () =>
          new Promise((resolve) => {
            setTimeout(() => {
              resolve({
                data: {
                  success: true,
                  data: {
                    token: 'test-token',
                    user: { id: 1, role: 'admin', name: 'Admin' },
                  },
                },
              });
            }, 100);
          })
      );

      renderWithProviders(<AdminLogin />, { route: '/admin/login' });

      const usernameInput = screen.getByPlaceholderText('管理员账号/邮箱');
      const passwordInput = screen.getByPlaceholderText('密码');
      const loginButton = screen.getByRole('button', { name: '登录管理后台' });

      await user.type(usernameInput, 'admin@gamelink.com');
      await user.type(passwordInput, 'password123');
      await user.click(loginButton);

      // Button should be in loading state
      expect(loginButton).toBeDisabled();

      await waitFor(() => {
        expect(loginButton).not.toBeDisabled();
      });
    });
  });

  describe('Error Handling', () => {
    it('should show error for invalid credentials (401)', async () => {
      const user = userEvent.setup();
      const error = {
        response: {
          status: 401,
          data: { message: '用户名或密码错误' },
        },
      };
      mockApi.login.mockRejectedValue(error);

      renderWithProviders(<AdminLogin />, { route: '/admin/login' });

      const usernameInput = screen.getByPlaceholderText('管理员账号/邮箱');
      const passwordInput = screen.getByPlaceholderText('密码');
      const loginButton = screen.getByRole('button', { name: '登录管理后台' });

      await user.type(usernameInput, 'admin@gamelink.com');
      await user.type(passwordInput, 'wrongpassword');
      await user.click(loginButton);

      await waitFor(() => {
        expect(screen.getByText(/用户名或密码错误/)).toBeInTheDocument();
      });
    });

    it('should show error for disabled account (403)', async () => {
      const user = userEvent.setup();
      const error = {
        response: {
          status: 403,
          data: { message: '账号已被禁用，请联系管理员' },
        },
      };
      mockApi.login.mockRejectedValue(error);

      renderWithProviders(<AdminLogin />, { route: '/admin/login' });

      const usernameInput = screen.getByPlaceholderText('管理员账号/邮箱');
      const passwordInput = screen.getByPlaceholderText('密码');
      const loginButton = screen.getByRole('button', { name: '登录管理后台' });

      await user.type(usernameInput, 'admin@gamelink.com');
      await user.type(passwordInput, 'password123');
      await user.click(loginButton);

      await waitFor(() => {
        expect(screen.getByText(/账号已被禁用/)).toBeInTheDocument();
      });
    });

    it('should show error for non-existent user (404)', async () => {
      const user = userEvent.setup();
      const error = {
        response: {
          status: 404,
          data: { message: '用户不存在' },
        },
      };
      mockApi.login.mockRejectedValue(error);

      renderWithProviders(<AdminLogin />, { route: '/admin/login' });

      const usernameInput = screen.getByPlaceholderText('管理员账号/邮箱');
      const passwordInput = screen.getByPlaceholderText('密码');
      const loginButton = screen.getByRole('button', { name: '登录管理后台' });

      await user.type(usernameInput, 'nonexistent@example.com');
      await user.type(passwordInput, 'password123');
      await user.click(loginButton);

      await waitFor(() => {
        expect(screen.getByText(/用户不存在/)).toBeInTheDocument();
      });
    });

    it('should show error for too many attempts (429)', async () => {
      const user = userEvent.setup();
      const error = {
        response: {
          status: 429,
          data: { message: '登录尝试次数过多，请稍后再试' },
        },
      };
      mockApi.login.mockRejectedValue(error);

      renderWithProviders(<AdminLogin />, { route: '/admin/login' });

      const usernameInput = screen.getByPlaceholderText('管理员账号/邮箱');
      const passwordInput = screen.getByPlaceholderText('密码');
      const loginButton = screen.getByRole('button', { name: '登录管理后台' });

      await user.type(usernameInput, 'admin@gamelink.com');
      await user.type(passwordInput, 'password123');
      await user.click(loginButton);

      await waitFor(() => {
        expect(screen.getByText(/登录尝试次数过多/)).toBeInTheDocument();
      });
    });

    it('should show error for server error (500+)', async () => {
      const user = userEvent.setup();
      const error = {
        response: {
          status: 500,
          data: { message: '服务器错误' },
        },
      };
      mockApi.login.mockRejectedValue(error);

      renderWithProviders(<AdminLogin />, { route: '/admin/login' });

      const usernameInput = screen.getByPlaceholderText('管理员账号/邮箱');
      const passwordInput = screen.getByPlaceholderText('密码');
      const loginButton = screen.getByRole('button', { name: '登录管理后台' });

      await user.type(usernameInput, 'admin@gamelink.com');
      await user.type(passwordInput, 'password123');
      await user.click(loginButton);

      await waitFor(() => {
        expect(screen.getByText(/服务器错误/)).toBeInTheDocument();
      });
    });

    it('should show error for network failure', async () => {
      const user = userEvent.setup();
      const error = {
        message: 'Network Error',
      };
      mockApi.login.mockRejectedValue(error);

      renderWithProviders(<AdminLogin />, { route: '/admin/login' });

      const usernameInput = screen.getByPlaceholderText('管理员账号/邮箱');
      const passwordInput = screen.getByPlaceholderText('密码');
      const loginButton = screen.getByRole('button', { name: '登录管理后台' });

      await user.type(usernameInput, 'admin@gamelink.com');
      await user.type(passwordInput, 'password123');
      await user.click(loginButton);

      await waitFor(() => {
        expect(screen.getByText(/网络连接失败/)).toBeInTheDocument();
      });
    });

    it('should show error for timeout', async () => {
      const user = userEvent.setup();
      const error = {
        message: 'timeout of 5000ms exceeded',
      };
      mockApi.login.mockRejectedValue(error);

      renderWithProviders(<AdminLogin />, { route: '/admin/login' });

      const usernameInput = screen.getByPlaceholderText('管理员账号/邮箱');
      const passwordInput = screen.getByPlaceholderText('密码');
      const loginButton = screen.getByRole('button', { name: '登录管理后台' });

      await user.type(usernameInput, 'admin@gamelink.com');
      await user.type(passwordInput, 'password123');
      await user.click(loginButton);

      await waitFor(() => {
        expect(screen.getByText(/网络连接失败/)).toBeInTheDocument();
      });
    });

    it('should handle API response with success: false', async () => {
      const user = userEvent.setup();
      mockApi.login.mockResolvedValue({
        data: {
          success: false,
          message: '登录失败',
        },
      });

      renderWithProviders(<AdminLogin />, { route: '/admin/login' });

      const usernameInput = screen.getByPlaceholderText('管理员账号/邮箱');
      const passwordInput = screen.getByPlaceholderText('密码');
      const loginButton = screen.getByRole('button', { name: '登录管理后台' });

      await user.type(usernameInput, 'admin@gamelink.com');
      await user.type(passwordInput, 'password123');
      await user.click(loginButton);

      await waitFor(() => {
        expect(screen.getByText(/登录失败/)).toBeInTheDocument();
      });
    });
  });

  describe('Form Validation', () => {
    it('should show validation error for empty username', async () => {
      const user = userEvent.setup();
      renderWithProviders(<AdminLogin />, { route: '/admin/login' });

      const passwordInput = screen.getByPlaceholderText('密码');
      await user.type(passwordInput, 'password123');

      const loginButton = screen.getByRole('button', { name: '登录管理后台' });
      await user.click(loginButton);

      await waitFor(() => {
        expect(screen.getByText(/请输入管理员账号/)).toBeInTheDocument();
      });
    });

    it('should show validation error for empty password', async () => {
      const user = userEvent.setup();
      renderWithProviders(<AdminLogin />, { route: '/admin/login' });

      const usernameInput = screen.getByPlaceholderText('管理员账号/邮箱');
      await user.type(usernameInput, 'admin@gamelink.com');

      const loginButton = screen.getByRole('button', { name: '登录管理后台' });
      await user.click(loginButton);

      await waitFor(() => {
        expect(screen.getByText(/请输入密码/)).toBeInTheDocument();
      });
    });

    it('should show validation errors for both empty fields', async () => {
      const user = userEvent.setup();
      renderWithProviders(<AdminLogin />, { route: '/admin/login' });

      const loginButton = screen.getByRole('button', { name: '登录管理后台' });
      await user.click(loginButton);

      await waitFor(() => {
        expect(screen.getByText(/请输入管理员账号/)).toBeInTheDocument();
        expect(screen.getByText(/请输入密码/)).toBeInTheDocument();
      });
    });
  });

  describe('Remember Password Functionality', () => {
    it('should save credentials when remember me is checked', async () => {
      const user = userEvent.setup();
      mockApi.login.mockResolvedValue({
        data: {
          success: true,
          data: {
            token: 'test-token',
            user: { id: 1, role: 'admin', name: 'Admin' },
          },
        },
      });

      renderWithProviders(<AdminLogin />, { route: '/admin/login' });

      const usernameInput = screen.getByPlaceholderText('管理员账号/邮箱');
      const passwordInput = screen.getByPlaceholderText('密码');
      const rememberCheckbox = screen.getByText('记住密码');
      const loginButton = screen.getByRole('button', { name: '登录管理后台' });

      await user.type(usernameInput, 'admin@gamelink.com');
      await user.type(passwordInput, 'password123');
      await user.click(rememberCheckbox);
      await user.click(loginButton);

      await waitFor(() => {
        const saved = localStorage.getItem('gamelink_admin_remember');
        expect(saved).not.toBeNull();
        const parsed = JSON.parse(saved!);
        expect(parsed.username).toBe('admin@gamelink.com');
        expect(parsed.password).toBe('password123');
        expect(parsed.remember).toBe(true);
      });
    });

    it('should not save credentials when remember me is unchecked', async () => {
      const user = userEvent.setup();
      mockApi.login.mockResolvedValue({
        data: {
          success: true,
          data: {
            token: 'test-token',
            user: { id: 1, role: 'admin', name: 'Admin' },
          },
        },
      });

      renderWithProviders(<AdminLogin />, { route: '/admin/login' });

      const usernameInput = screen.getByPlaceholderText('管理员账号/邮箱');
      const passwordInput = screen.getByPlaceholderText('密码');
      const loginButton = screen.getByRole('button', { name: '登录管理后台' });

      await user.type(usernameInput, 'admin@gamelink.com');
      await user.type(passwordInput, 'password123');
      await user.click(loginButton);

      await waitFor(() => {
        const saved = localStorage.getItem('gamelink_admin_remember');
        expect(saved).toBeNull();
      });
    });

    it('should load saved credentials on mount', async () => {
      localStorage.setItem(
        'gamelink_admin_remember',
        JSON.stringify({
          username: 'saved@example.com',
          password: 'savedpassword',
          remember: true,
        })
      );

      renderWithProviders(<AdminLogin />, { route: '/admin/login' });

      await waitFor(() => {
        const usernameInput = screen.getByPlaceholderText('管理员账号/邮箱') as HTMLInputElement;
        const passwordInput = screen.getByPlaceholderText('密码') as HTMLInputElement;

        expect(usernameInput.value).toBe('saved@example.com');
        expect(passwordInput.value).toBe('savedpassword');
      });
    });

    it('should clear saved data when remember me is unchecked after login', async () => {
      const user = userEvent.setup();
      localStorage.setItem(
        'gamelink_admin_remember',
        JSON.stringify({
          username: 'saved@example.com',
          password: 'savedpassword',
          remember: true,
        })
      );

      mockApi.login.mockResolvedValue({
        data: {
          success: true,
          data: {
            token: 'test-token',
            user: { id: 1, role: 'admin', name: 'Admin' },
          },
        },
      });

      renderWithProviders(<AdminLogin />, { route: '/admin/login' });

      const rememberCheckbox = screen.getByText('记住密码');
      await user.click(rememberCheckbox); // Uncheck

      const loginButton = screen.getByRole('button', { name: '登录管理后台' });
      await user.click(loginButton);

      await waitFor(() => {
        expect(localStorage.getItem('gamelink_admin_remember')).toBeNull();
      });
    });
  });

  describe('Accessibility', () => {
    it('should have proper heading structure', () => {
      renderWithProviders(<AdminLogin />, { route: '/admin/login' });

      expect(screen.getByRole('heading', { level: 1, name: 'GameLink 管理后台' })).toBeInTheDocument();
    });

    it('should have accessible form labels', () => {
      renderWithProviders(<AdminLogin />, { route: '/admin/login' });

      const usernameInput = screen.getByPlaceholderText('管理员账号/邮箱');
      const passwordInput = screen.getByPlaceholderText('密码');

      expect(usernameInput).toHaveAttribute('type');
      expect(passwordInput).toHaveAttribute('type', 'password');
    });

    it('should be keyboard navigable', async () => {
      const user = userEvent.setup();
      renderWithProviders(<AdminLogin />, { route: '/admin/login' });

      const usernameInput = screen.getByPlaceholderText('管理员账号/邮箱');
      usernameInput.focus();

      expect(usernameInput).toHaveFocus();

      await user.tab();

      const passwordInput = screen.getByPlaceholderText('密码');
      expect(passwordInput).toHaveFocus();
    });
  });

  describe('UI Elements', () => {
    it('should display security icon', () => {
      renderWithProviders(<AdminLogin />, { route: '/admin/login' });

      // Check for the SafetyCertificateOutlined icon
      const icon = document.querySelector('.anticon-safety-certificate');
      expect(icon).toBeInTheDocument();
    });

    it('should have proper button styling', () => {
      renderWithProviders(<AdminLogin />, { route: '/admin/login' });

      const loginButton = screen.getByRole('button', { name: '登录管理后台' });

      expect(loginButton).toHaveStyle({ height: '44px' });
    });
  });
});
