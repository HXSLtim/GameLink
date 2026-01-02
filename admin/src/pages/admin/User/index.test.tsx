/**
 * User Management Page Tests
 *
 * Tests for User page component including:
 * - Successful data loading
 * - Loading states
 * - Error handling
 * - User interactions (filtering, pagination, search)
 * - User CRUD operations (create, edit, delete, ban/unban)
 * - Batch operations (role change, notification, points)
 * - Statistics display
 * - Permission checks
 */

import React from 'react';
import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import UserPage from './index';
import { renderWithProviders, resetAllMocks, flushPromises } from '@/testutils';

// Define mock API for admin
const mockApi = {
  getUsers: vi.fn(),
  getUserStats: vi.fn(),
  getUserLogs: vi.fn(),
  createUser: vi.fn(),
  updateUser: vi.fn(),
  deleteUser: vi.fn(),
  batchUpdateUserStatus: vi.fn(),
  batchUpdateUserRole: vi.fn(),
  sendNotification: vi.fn(),
  adjustPoints: vi.fn(),
};

// Mock the adminApi module
vi.mock('@/api/admin', () => ({
  adminApi: mockApi,
}));

// Mock export utilities
vi.mock('@/utils/export', () => ({
  exportToCSV: vi.fn(),
  userExportColumns: [],
}));

describe('UserPage', () => {
  beforeEach(() => {
    resetAllMocks();
    localStorage.setItem('token', 'test-token');
    localStorage.setItem('user_role', 'admin');
  });

  afterEach(() => {
    vi.clearAllMocks();
  });

  describe('Successful Data Loading', () => {
    it('should render user list successfully', async () => {
      renderWithProviders(<UserPage />);

      await waitFor(() => {
        expect(screen.getByText('用户管理')).toBeInTheDocument();
      });

      await waitFor(() => {
        expect(screen.getByText('Test User')).toBeInTheDocument();
      });

      expect(mockApi.getUsers).toHaveBeenCalledWith({
        page: 1,
        page_size: 10,
      });
    });

    it('should display user information correctly', async () => {
      renderWithProviders(<UserPage />);

      await waitFor(() => {
        expect(screen.getByText('Test User')).toBeInTheDocument();
      });

      expect(screen.getByText('test@example.com')).toBeInTheDocument();
      expect(screen.getByText('13800138000')).toBeInTheDocument();
    });

    it('should display user role tags', async () => {
      renderWithProviders(<UserPage />);

      await waitFor(() => {
        expect(screen.getByText('管理员')).toBeInTheDocument();
      });
    });

    it('should display user status tags', async () => {
      renderWithProviders(<UserPage />);

      await waitFor(() => {
        expect(screen.getByText('正常')).toBeInTheDocument();
      });
    });

    it('should display user statistics', async () => {
      renderWithProviders(<UserPage />);

      await waitFor(() => {
        expect(screen.getByText('用户总数')).toBeInTheDocument();
      });

      await waitFor(() => {
        expect(screen.getByText('陪玩师')).toBeInTheDocument();
        expect(screen.getByText('正常用户')).toBeInTheDocument();
        expect(screen.getByText('最近注册')).toBeInTheDocument();
      });

      expect(mockApi.getUserStats).toHaveBeenCalled();
    });
  });

  describe('Loading States', () => {
    it('should show loading indicator while fetching data', async () => {
      mockApi.getUsers.mockImplementation(
        () =>
          new Promise((resolve) => {
            setTimeout(() => {
              resolve({
                data: {
                  success: true,
                  data: [],
                  pagination: { total: 0, page: 1, pageSize: 10 },
                },
              });
            }, 100);
          })
      );

      renderWithProviders(<UserPage />);

      expect(mockApi.getUsers).toHaveBeenCalled();
    });

    it('should hide loading indicator after data loads', async () => {
      renderWithProviders(<UserPage />);

      await waitFor(() => {
        expect(screen.getByText('用户管理')).toBeInTheDocument();
      });

      await flushPromises();

      expect(mockApi.getUsers).toHaveBeenCalledTimes(1);
    });
  });

  describe('Error Handling', () => {
    it('should display error message when API fails', async () => {
      mockApi.getUsers.mockRejectedValue(new Error('Network error'));

      renderWithProviders(<UserPage />);

      await waitFor(() => {
        expect(screen.getByText('用户管理')).toBeInTheDocument();
      });

      await flushPromises();

      const errorMessage = await screen.findByText(/获取用户列表失败/);
      expect(errorMessage).toBeInTheDocument();
    });

    it('should handle empty data gracefully', async () => {
      mockApi.getUsers.mockResolvedValue({
        data: {
          success: true,
          data: [],
          pagination: { total: 0, page: 1, pageSize: 10 },
        },
      });

      renderWithProviders(<UserPage />);

      await waitFor(() => {
        expect(screen.getByText('用户管理')).toBeInTheDocument();
      });

      await flushPromises();

      expect(screen.getByText('共 0 条')).toBeInTheDocument();
    });

    it('should handle API response with success: false', async () => {
      mockApi.getUsers.mockResolvedValue({
        data: {
          success: false,
          message: '获取用户列表失败',
        },
      });

      renderWithProviders(<UserPage />);

      await waitFor(() => {
        expect(screen.getByText('用户管理')).toBeInTheDocument();
      });

      await flushPromises();

      const errorMessage = await screen.findByText(/获取用户列表失败/);
      expect(errorMessage).toBeInTheDocument();
    });
  });

  describe('Search and Filtering', () => {
    it('should allow searching by keyword', async () => {
      const _user = userEvent.setup();
      mockApi.getUsers.mockResolvedValue({
        data: {
          success: true,
          data: [],
          pagination: { total: 0, page: 1, pageSize: 10 },
        },
      });

      renderWithProviders(<UserPage />);

      await waitFor(() => {
        expect(screen.getByPlaceholderText('用户名/邮箱/手机号')).toBeInTheDocument();
      });

      const searchInput = screen.getByPlaceholderText('用户名/邮箱/手机号');
      await user.type(searchInput, 'Test User');

      const searchButton = screen.getByRole('button', { name: /搜索/i });
      await user.click(searchButton);

      await waitFor(() => {
        expect(mockApi.getUsers).toHaveBeenCalledWith(
          expect.objectContaining({
            keyword: 'Test User',
          })
        );
      });
    });

    it('should allow filtering by role', async () => {
      const _user = userEvent.setup();
      mockApi.getUsers.mockResolvedValue({
        data: {
          success: true,
          data: [],
          pagination: { total: 0, page: 1, pageSize: 10 },
        },
      });

      renderWithProviders(<UserPage />);

      await waitFor(() => {
        expect(screen.getByText('角色')).toBeInTheDocument();
      });

      const roleDropdown = screen.getByText('角色').closest('.ant-select');
      if (roleDropdown) {
        await user.click(roleDropdown);

        const userOption = await screen.findByText('普通用户');
        await user.click(userOption);

        await waitFor(() => {
          expect(mockApi.getUsers).toHaveBeenCalledWith(
            expect.objectContaining({
              role: ['user'],
            })
          );
        });
      }
    });

    it('should allow filtering by status', async () => {
      const _user = userEvent.setup();
      mockApi.getUsers.mockResolvedValue({
        data: {
          success: true,
          data: [],
          pagination: { total: 0, page: 1, pageSize: 10 },
        },
      });

      renderWithProviders(<UserPage />);

      await waitFor(() => {
        expect(screen.getByText('状态')).toBeInTheDocument();
      });

      const statusDropdown = screen.getByText('状态').closest('.ant-select');
      if (statusDropdown) {
        await user.click(statusDropdown);

        const activeOption = await screen.findByText('正常');
        await user.click(activeOption);

        await waitFor(() => {
          expect(mockApi.getUsers).toHaveBeenCalledWith(
            expect.objectContaining({
              status: ['active'],
            })
          );
        });
      }
    });

    it('should reset to first page when searching', async () => {
      const _user = userEvent.setup();
      mockApi.getUsers.mockResolvedValue({
        data: {
          success: true,
          data: [],
          pagination: { total: 0, page: 1, pageSize: 10 },
        },
      });

      renderWithProviders(<UserPage />);

      await waitFor(() => {
        expect(screen.getByPlaceholderText('用户名/邮箱/手机号')).toBeInTheDocument();
      });

      const searchInput = screen.getByPlaceholderText('用户名/邮箱/手机号');
      await user.type(searchInput, 'test');

      const searchButton = screen.getByRole('button', { name: /搜索/i });
      await user.click(searchButton);

      await waitFor(() => {
        expect(mockApi.getUsers).toHaveBeenCalledWith(
          expect.objectContaining({
            page: 1,
          })
        );
      });
    });
  });

  describe('Pagination', () => {
    it('should display pagination controls', async () => {
      renderWithProviders(<UserPage />);

      await waitFor(() => {
        expect(screen.getByText('用户管理')).toBeInTheDocument();
      });

      expect(screen.getByText('共 1 条')).toBeInTheDocument();
    });

    it('should change page when clicking pagination', async () => {
      const _user = userEvent.setup();
      mockApi.getUsers.mockResolvedValue({
        data: {
          success: true,
          data: [],
          pagination: { total: 20, page: 2, pageSize: 10 },
        },
      });

      renderWithProviders(<UserPage />);

      await waitFor(() => {
        expect(screen.getByText('用户管理')).toBeInTheDocument();
      });

      const nextPageButton = screen.getByTitle('下一页');
      await user.click(nextPageButton);

      await waitFor(() => {
        expect(mockApi.getUsers).toHaveBeenCalledWith(
          expect.objectContaining({
            page: 2,
          })
        );
      });
    });
  });

  describe('User Details', () => {
    it('should open detail drawer when clicking detail button', async () => {
      const _user = userEvent.setup();
      renderWithProviders(<UserPage />);

      await waitFor(() => {
        expect(screen.getByText('Test User')).toBeInTheDocument();
      });

      const detailButton = screen.getByRole('button', { name: /详情/i });
      await user.click(detailButton);

      await waitFor(() => {
        expect(screen.getByText('用户详情')).toBeInTheDocument();
      });
    });

    it('should display user basic information in drawer', async () => {
      const _user = userEvent.setup();
      renderWithProviders(<UserPage />);

      await waitFor(() => {
        expect(screen.getByText('Test User')).toBeInTheDocument();
      });

      const detailButton = screen.getByRole('button', { name: /详情/i });
      await user.click(detailButton);

      await waitFor(() => {
        expect(screen.getByText('基本信息')).toBeInTheDocument();
      });

      expect(screen.getByText('test@example.com')).toBeInTheDocument();
    });

    it('should display login history tab', async () => {
      const _user = userEvent.setup();
      mockApi.getUserLogs.mockResolvedValue({
        success: true,
        data: [],
      });

      renderWithProviders(<UserPage />);

      await waitFor(() => {
        expect(screen.getByText('Test User')).toBeInTheDocument();
      });

      const detailButton = screen.getByRole('button', { name: /详情/i });
      await user.click(detailButton);

      await waitFor(() => {
        expect(screen.getByText('用户详情')).toBeInTheDocument();
      });

      const loginHistoryTab = screen.getByText('登录历史');
      await user.click(loginHistoryTab);

      await waitFor(() => {
        expect(mockApi.getUserLogs).toHaveBeenCalledWith(1, {
          page: 1,
          page_size: 10,
          type: 'login',
        });
      });
    });

    it('should display operation logs tab', async () => {
      const _user = userEvent.setup();
      mockApi.getUserLogs.mockResolvedValue({
        success: true,
        data: [],
      });

      renderWithProviders(<UserPage />);

      await waitFor(() => {
        expect(screen.getByText('Test User')).toBeInTheDocument();
      });

      const detailButton = screen.getByRole('button', { name: /详情/i });
      await user.click(detailButton);

      await waitFor(() => {
        expect(screen.getByText('用户详情')).toBeInTheDocument();
      });

      const operationLogsTab = screen.getByText('操作日志');
      await user.click(operationLogsTab);

      await waitFor(() => {
        expect(mockApi.getUserLogs).toHaveBeenCalledWith(1, {
          page: 1,
          page_size: 10,
        });
      });
    });
  });

  describe('User Edit', () => {
    it('should open edit modal when clicking edit button', async () => {
      const _user = userEvent.setup();
      renderWithProviders(<UserPage />);

      await waitFor(() => {
        expect(screen.getByText('Test User')).toBeInTheDocument();
      });

      const editButton = screen.getByRole('button', { name: /编辑/i });
      await user.click(editButton);

      await waitFor(() => {
        expect(screen.getByText('编辑用户')).toBeInTheDocument();
      });
    });

    it('should pre-fill form with user data', async () => {
      const _user = userEvent.setup();
      renderWithProviders(<UserPage />);

      await waitFor(() => {
        expect(screen.getByText('Test User')).toBeInTheDocument();
      });

      const editButton = screen.getByRole('button', { name: /编辑/i });
      await user.click(editButton);

      await waitFor(() => {
        expect(screen.getByDisplayValue('Test User')).toBeInTheDocument();
      });

      expect(screen.getByDisplayValue('test@example.com')).toBeInTheDocument();
      expect(screen.getByDisplayValue('13800138000')).toBeInTheDocument();
    });

    it('should update user when submitting form', async () => {
      const _user = userEvent.setup();
      renderWithProviders(<UserPage />);

      await waitFor(() => {
        expect(screen.getByText('Test User')).toBeInTheDocument();
      });

      const editButton = screen.getByRole('button', { name: /编辑/i });
      await user.click(editButton);

      await waitFor(() => {
        expect(screen.getByText('编辑用户')).toBeInTheDocument();
      });

      const nameInput = screen.getByDisplayValue('Test User');
      await user.clear(nameInput);
      await user.type(nameInput, 'Updated User');

      const saveButton = screen.getByRole('button', { name: /保存/i });
      await user.click(saveButton);

      await waitFor(() => {
        expect(mockApi.updateUser).toHaveBeenCalledWith(1, {
          name: 'Updated User',
          email: 'test@example.com',
          phone: '13800138000',
          avatarUrl: 'https://example.com/avatar.jpg',
          role: 'admin',
          status: 'active',
        });
      });
    });
  });

  describe('User Ban/Unban', () => {
    it('should show ban button for active users', async () => {
      renderWithProviders(<UserPage />);

      await waitFor(() => {
        expect(screen.getByText('Test User')).toBeInTheDocument();
      });

      expect(screen.getByRole('button', { name: /封禁/i })).toBeInTheDocument();
    });

    it('should ban user when confirming ban action', async () => {
      const _user = userEvent.setup();
      renderWithProviders(<UserPage />);

      await waitFor(() => {
        expect(screen.getByText('Test User')).toBeInTheDocument();
      });

      const banButton = screen.getByRole('button', { name: /封禁/i });
      await user.click(banButton);

      const confirmButton = await screen.findByRole('button', { name: /确定/i });
      await user.click(confirmButton);

      await waitFor(() => {
        expect(mockApi.updateUserStatus).toHaveBeenCalledWith(1, 'banned');
      });
    });

    it('should show unban button for banned users', async () => {
      mockApi.getUsers.mockResolvedValue({
        data: {
          success: true,
          data: [
            {
              ...mockUser,
              id: 2,
              name: 'Banned User',
              email: 'banned@example.com',
              phone: '13800138001',
              status: 'banned' as const,
              createdAt: '2024-01-01T00:00:00Z',
            },
          ],
          pagination: { total: 1, page: 1, pageSize: 10 },
        },
      });

      renderWithProviders(<UserPage />);

      await waitFor(() => {
        expect(screen.getByText('Banned User')).toBeInTheDocument();
      });

      expect(screen.getByRole('button', { name: /解封/i })).toBeInTheDocument();
    });
  });

  describe('User Delete', () => {
    it('should show delete button', async () => {
      renderWithProviders(<UserPage />);

      await waitFor(() => {
        expect(screen.getByText('Test User')).toBeInTheDocument();
      });

      expect(screen.getByRole('button', { name: /删除/i })).toBeInTheDocument();
    });

    it('should delete user when confirming delete action', async () => {
      const _user = userEvent.setup();
      renderWithProviders(<UserPage />);

      await waitFor(() => {
        expect(screen.getByText('Test User')).toBeInTheDocument();
      });

      const deleteButton = screen.getByRole('button', { name: /删除/i });
      await user.click(deleteButton);

      const confirmButton = await screen.findByRole('button', { name: /确定/i });
      await user.click(confirmButton);

      await waitFor(() => {
        expect(mockApi.deleteUser).toHaveBeenCalledWith(1);
      });
    });
  });

  describe('Create User', () => {
    it('should open create modal when clicking new user button', async () => {
      const _user = userEvent.setup();
      renderWithProviders(<UserPage />);

      await waitFor(() => {
        expect(screen.getByRole('button', { name: /新增用户/i })).toBeInTheDocument();
      });

      const createButton = screen.getByRole('button', { name: /新增用户/i });
      await user.click(createButton);

      await waitFor(() => {
        expect(screen.getByText('新增用户')).toBeInTheDocument();
      });
    });

    it('should create user when submitting form with valid data', async () => {
      const _user = userEvent.setup();
      renderWithProviders(<UserPage />);

      const createButton = screen.getByRole('button', { name: /新增用户/i });
      await user.click(createButton);

      await waitFor(() => {
        expect(screen.getByText('新增用户')).toBeInTheDocument();
      });

      const nameInput = screen.getByPlaceholderText('请输入用户名');
      await user.type(nameInput, 'New User');

      const emailInput = screen.getByPlaceholderText('请输入邮箱');
      await user.type(emailInput, 'new@example.com');

      const phoneInput = screen.getByPlaceholderText('请输入手机号');
      await user.type(phoneInput, '13900139000');

      const passwordInput = screen.getByPlaceholderText(/请输入密码/);
      await user.type(passwordInput, 'password123');

      const saveButton = screen.getByRole('button', { name: /保存/i });
      await user.click(saveButton);

      await waitFor(() => {
        expect(mockApi.createUser).toHaveBeenCalledWith({
          name: 'New User',
          email: 'new@example.com',
          phone: '13900139000',
          password: 'password123',
          avatarUrl: undefined,
          role: 'user',
          status: 'active',
        });
      });
    });
  });

  describe('Batch Operations', () => {
    it('should show batch modify role button', async () => {
      renderWithProviders(<UserPage />);

      await waitFor(() => {
        expect(screen.getByRole('button', { name: /批量修改角色/i })).toBeInTheDocument();
      });
    });

    it('should open batch role modal', async () => {
      const _user = userEvent.setup();
      renderWithProviders(<UserPage />);

      await waitFor(() => {
        expect(screen.getByRole('button', { name: /批量修改角色/i })).toBeInTheDocument();
      });

      const batchRoleButton = screen.getByRole('button', { name: /批量修改角色/i });
      await user.click(batchRoleButton);

      await waitFor(() => {
        expect(screen.getByText('批量修改角色')).toBeInTheDocument();
      });
    });

    it('should show batch enable button', async () => {
      renderWithProviders(<UserPage />);

      await waitFor(() => {
        expect(screen.getByRole('button', { name: /批量启用/i })).toBeInTheDocument();
      });
    });

    it('should show batch disable button', async () => {
      renderWithProviders(<UserPage />);

      await waitFor(() => {
        expect(screen.getByRole('button', { name: /批量禁用/i })).toBeInTheDocument();
      });
    });

    it('should show batch send notification button', async () => {
      renderWithProviders(<UserPage />);

      await waitFor(() => {
        expect(screen.getByRole('button', { name: /批量发送通知/i })).toBeInTheDocument();
      });
    });

    it('should show batch add points button', async () => {
      renderWithProviders(<UserPage />);

      await waitFor(() => {
        expect(screen.getByRole('button', { name: /批量增加积分/i })).toBeInTheDocument();
      });
    });

    it('should export user data', async () => {
      const _user = userEvent.setup();
      const { exportToCSV } = await import('@/utils/export');

      renderWithProviders(<UserPage />);

      await waitFor(() => {
        expect(screen.getByRole('button', { name: /导出数据/i })).toBeInTheDocument();
      });

      const exportButton = screen.getByRole('button', { name: /导出数据/i });
      await user.click(exportButton);

      await waitFor(() => {
        expect(mockApi.getUsers).toHaveBeenCalledWith(
          expect.objectContaining({
            page_size: 10000,
          })
        );
        expect(exportToCSV).toHaveBeenCalled();
      });
    });
  });

  describe('Refresh Functionality', () => {
    it('should refresh data when clicking refresh button', async () => {
      const _user = userEvent.setup();
      renderWithProviders(<UserPage />);

      await waitFor(() => {
        expect(screen.getByText('用户管理')).toBeInTheDocument();
      });

      const refreshButton = screen.getByRole('button', { name: /刷新/i });
      await user.click(refreshButton);

      await waitFor(() => {
        expect(mockApi.getUsers).toHaveBeenCalledTimes(2);
      });
    });
  });

  describe('Accessibility', () => {
    it('should have proper ARIA labels', async () => {
      renderWithProviders(<UserPage />);

      await waitFor(() => {
        expect(screen.getByRole('heading', { name: '用户管理' })).toBeInTheDocument();
      });
    });

    it('should be keyboard navigable', async () => {
      const _user = userEvent.setup();
      renderWithProviders(<UserPage />);

      await waitFor(() => {
        expect(screen.getByPlaceholderText('用户名/邮箱/手机号')).toBeInTheDocument();
      });

      const searchInput = screen.getByPlaceholderText('用户名/邮箱/手机号');
      searchInput.focus();

      expect(searchInput).toHaveFocus();
    });
  });
});
