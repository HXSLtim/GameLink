/**
 * Player Management Page Tests
 *
 * Tests for Player page component including:
 * - Successful data loading
 * - Loading states
 * - Error handling
 * - User interactions (filtering, pagination, search)
 * - Player operations (audit, ban/unban, batch operations)
 * - Permission checks
 */

import React from 'react';
import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import PlayerPage from './index';
import { renderWithProviders, resetAllMocks, flushPromises } from '@/testutils';

// Mock the adminApi module using vi.hoisted to avoid hoisting issues
const { mockGetPlayers, mockUpdatePlayerStatus, mockDeletePlayer, mockBatchUpdatePlayerStatus } = vi.hoisted(() => ({
  mockGetPlayers: vi.fn(),
  mockUpdatePlayerStatus: vi.fn(),
  mockDeletePlayer: vi.fn(),
  mockBatchUpdatePlayerStatus: vi.fn(),
}));

vi.mock('@/api/admin', () => ({
  adminApi: {
    getPlayers: mockGetPlayers,
    updatePlayerStatus: mockUpdatePlayerStatus,
    deletePlayer: mockDeletePlayer,
    batchUpdatePlayerStatus: mockBatchUpdatePlayerStatus,
  },
}));

// Export mockApi for use in tests
const mockApi = {
  getPlayers: mockGetPlayers,
  updatePlayerStatus: mockUpdatePlayerStatus,
  deletePlayer: mockDeletePlayer,
  batchUpdatePlayerStatus: mockBatchUpdatePlayerStatus,
};

// Mock export utilities
vi.mock('@/utils/export', () => ({
  exportToCSV: vi.fn(),
  playerExportColumns: [],
}));

describe('PlayerPage', () => {
  beforeEach(() => {
    resetAllMocks();
    localStorage.setItem('token', 'test-token');
    localStorage.setItem('user_role', 'admin');
  });

  afterEach(() => {
    vi.clearAllMocks();
  });

  describe('Successful Data Loading', () => {
    it('should render player list successfully', async () => {
      renderWithProviders(<PlayerPage />);

      await waitFor(() => {
        expect(screen.getByText('陪玩师管理')).toBeInTheDocument();
      });

      await waitFor(() => {
        expect(screen.getByText('Test Player')).toBeInTheDocument();
      });

      expect(mockApi.getPlayers).toHaveBeenCalledWith({
        page: 1,
        page_size: 10,
      });
    });

    it('should display player information correctly', async () => {
      renderWithProviders(<PlayerPage />);

      await waitFor(() => {
        expect(screen.getByText('Test Player')).toBeInTheDocument();
      });

      expect(screen.getByText('Diamond')).toBeInTheDocument();
      expect(screen.getByText('¥50.00')).toBeInTheDocument();
      expect(screen.getByText('Test Game')).toBeInTheDocument();
    });

    it('should display player rating', async () => {
      renderWithProviders(<PlayerPage />);

      await waitFor(() => {
        expect(screen.getByText('4.5')).toBeInTheDocument();
      });

      expect(screen.getByText('(10)')).toBeInTheDocument();
    });

    it('should display player skill tags', async () => {
      renderWithProviders(<PlayerPage />);

      await waitFor(() => {
        expect(screen.getByText('friendly')).toBeInTheDocument();
      });

      expect(screen.getByText('skilled')).toBeInTheDocument();
    });

    it('should display player verification status', async () => {
      renderWithProviders(<PlayerPage />);

      await waitFor(() => {
        expect(screen.getByText('已通过')).toBeInTheDocument();
      });
    });
  });

  describe('Loading States', () => {
    it('should show loading indicator while fetching data', async () => {
      mockApi.getPlayers.mockImplementation(
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

      renderWithProviders(<PlayerPage />);

      expect(mockApi.getPlayers).toHaveBeenCalled();
    });

    it('should hide loading indicator after data loads', async () => {
      renderWithProviders(<PlayerPage />);

      await waitFor(() => {
        expect(screen.getByText('陪玩师管理')).toBeInTheDocument();
      });

      await flushPromises();

      expect(mockApi.getPlayers).toHaveBeenCalledTimes(1);
    });
  });

  describe('Error Handling', () => {
    it('should display error message when API fails', async () => {
      mockApi.getPlayers.mockRejectedValue(new Error('Network error'));

      renderWithProviders(<PlayerPage />);

      await waitFor(() => {
        expect(screen.getByText('陪玩师管理')).toBeInTheDocument();
      });

      await flushPromises();

      const errorMessage = await screen.findByText(/加载陪玩师列表失败/);
      expect(errorMessage).toBeInTheDocument();
    });

    it('should handle empty data gracefully', async () => {
      mockApi.getPlayers.mockResolvedValue({
        data: {
          success: true,
          data: [],
          pagination: { total: 0, page: 1, pageSize: 10 },
        },
      });

      renderWithProviders(<PlayerPage />);

      await waitFor(() => {
        expect(screen.getByText('陪玩师管理')).toBeInTheDocument();
      });

      await flushPromises();

      expect(screen.getByText('共 0 条')).toBeInTheDocument();
    });

    it('should handle API response with success: false', async () => {
      mockApi.getPlayers.mockResolvedValue({
        data: {
          success: false,
          message: '获取陪玩师列表失败',
        },
      });

      renderWithProviders(<PlayerPage />);

      await waitFor(() => {
        expect(screen.getByText('陪玩师管理')).toBeInTheDocument();
      });

      await flushPromises();

      const errorMessage = await screen.findByText(/获取陪玩师列表失败/);
      expect(errorMessage).toBeInTheDocument();
    });
  });

  describe('Search and Filtering', () => {
    it('should allow searching by keyword', async () => {
      const _user = userEvent.setup();
      mockApi.getPlayers.mockResolvedValue({
        data: {
          success: true,
          data: [],
          pagination: { total: 0, page: 1, pageSize: 10 },
        },
      });

      renderWithProviders(<PlayerPage />);

      await waitFor(() => {
        expect(screen.getByPlaceholderText('名称/ID')).toBeInTheDocument();
      });

      const searchInput = screen.getByPlaceholderText('名称/ID');
      await user.type(searchInput, 'Test Player');

      const searchButton = screen.getByRole('button', { name: /搜索/i });
      await user.click(searchButton);

      await waitFor(() => {
        expect(mockApi.getPlayers).toHaveBeenCalledWith(
          expect.objectContaining({
            keyword: 'Test Player',
          })
        );
      });
    });

    it('should allow filtering by status', async () => {
      const _user = userEvent.setup();
      mockApi.getPlayers.mockResolvedValue({
        data: {
          success: true,
          data: [],
          pagination: { total: 0, page: 1, pageSize: 10 },
        },
      });

      renderWithProviders(<PlayerPage />);

      await waitFor(() => {
        expect(screen.getByText('状态')).toBeInTheDocument();
      });

      const statusDropdown = screen.getByText('状态').closest('.ant-select');
      if (statusDropdown) {
        await user.click(statusDropdown);

        const verifiedOption = await screen.findByText('已通过');
        await user.click(verifiedOption);

        await waitFor(() => {
          expect(mockApi.getPlayers).toHaveBeenCalledWith(
            expect.objectContaining({
              status: 'verified',
            })
          );
        });
      }
    });

    it('should reset to first page when searching', async () => {
      const _user = userEvent.setup();
      mockApi.getPlayers.mockResolvedValue({
        data: {
          success: true,
          data: [],
          pagination: { total: 0, page: 1, pageSize: 10 },
        },
      });

      renderWithProviders(<PlayerPage />);

      await waitFor(() => {
        expect(screen.getByPlaceholderText('名称/ID')).toBeInTheDocument();
      });

      const searchInput = screen.getByPlaceholderText('名称/ID');
      await user.type(searchInput, 'test');

      const searchButton = screen.getByRole('button', { name: /搜索/i });
      await user.click(searchButton);

      await waitFor(() => {
        expect(mockApi.getPlayers).toHaveBeenCalledWith(
          expect.objectContaining({
            page: 1,
          })
        );
      });
    });
  });

  describe('Pagination', () => {
    it('should display pagination controls', async () => {
      renderWithProviders(<PlayerPage />);

      await waitFor(() => {
        expect(screen.getByText('陪玩师管理')).toBeInTheDocument();
      });

      expect(screen.getByText('共 1 条')).toBeInTheDocument();
    });

    it('should change page when clicking pagination', async () => {
      const _user = userEvent.setup();
      mockApi.getPlayers.mockResolvedValue({
        data: {
          success: true,
          data: [],
          pagination: { total: 20, page: 2, pageSize: 10 },
        },
      });

      renderWithProviders(<PlayerPage />);

      await waitFor(() => {
        expect(screen.getByText('陪玩师管理')).toBeInTheDocument();
      });

      const nextPageButton = screen.getByTitle('下一页');
      await user.click(nextPageButton);

      await waitFor(() => {
        expect(mockApi.getPlayers).toHaveBeenCalledWith(
          expect.objectContaining({
            page: 2,
          })
        );
      });
    });

    it('should change page size when selecting different size', async () => {
      const _user = userEvent.setup();
      mockApi.getPlayers.mockResolvedValue({
        data: {
          success: true,
          data: [],
          pagination: { total: 50, page: 1, pageSize: 20 },
        },
      });

      renderWithProviders(<PlayerPage />);

      await waitFor(() => {
        expect(screen.getByText('陪玩师管理')).toBeInTheDocument();
      });

      const pageSizeSelector = screen.getByText('10 条/页');
      await user.click(pageSizeSelector);

      const pageSize20 = await screen.findByText('20 条/页');
      await user.click(pageSize20);

      await waitFor(() => {
        expect(mockApi.getPlayers).toHaveBeenCalledWith(
          expect.objectContaining({
            page_size: 20,
          })
        );
      });
    });
  });

  describe('Player Details', () => {
    it('should open detail drawer when clicking detail button', async () => {
      const _user = userEvent.setup();
      renderWithProviders(<PlayerPage />);

      await waitFor(() => {
        expect(screen.getByText('Test Player')).toBeInTheDocument();
      });

      const detailButton = screen.getByRole('button', { name: /详情/i });
      await user.click(detailButton);

      await waitFor(() => {
        expect(screen.getByText('陪玩师详情')).toBeInTheDocument();
      });
    });

    it('should display player statistics in drawer', async () => {
      const _user = userEvent.setup();
      renderWithProviders(<PlayerPage />);

      await waitFor(() => {
        expect(screen.getByText('Test Player')).toBeInTheDocument();
      });

      const detailButton = screen.getByRole('button', { name: /详情/i });
      await user.click(detailButton);

      await waitFor(() => {
        expect(screen.getByText('评分')).toBeInTheDocument();
        expect(screen.getByText('评价数')).toBeInTheDocument();
        expect(screen.getByText('时薪')).toBeInTheDocument();
      });
    });

    it('should display player basic information', async () => {
      const _user = userEvent.setup();
      renderWithProviders(<PlayerPage />);

      await waitFor(() => {
        expect(screen.getByText('Test Player')).toBeInTheDocument();
      });

      const detailButton = screen.getByRole('button', { name: /详情/i });
      await user.click(detailButton);

      await waitFor(() => {
        expect(screen.getByText('基本信息')).toBeInTheDocument();
      });

      expect(screen.getByText('Test Game')).toBeInTheDocument();
      expect(screen.getByText('Test bio')).toBeInTheDocument();
    });
  });

  describe('Player Audit', () => {
    it('should show audit button for pending players', async () => {
      mockApi.getPlayers.mockResolvedValue({
        data: {
          success: true,
          data: [
            {
              id: 2,
              userId: 2,
              nickname: 'Pending Player',
              verificationStatus: 'pending' as const,
              hourlyRateCents: 5000,
              ratingAverage: 0,
              ratingCount: 0,
              skillTags: [],
              bio: 'Pending',
              mainGame: { id: 1, name: 'Test Game' },
              createdAt: '2024-01-01T00:00:00Z',
              user: { id: 2, name: 'Pending User', avatarUrl: '' },
            },
          ],
          pagination: { total: 1, page: 1, pageSize: 10 },
        },
      });

      renderWithProviders(<PlayerPage />);

      await waitFor(() => {
        expect(screen.getByText('Pending Player')).toBeInTheDocument();
      });

      expect(screen.getByRole('button', { name: /审核/i })).toBeInTheDocument();
    });

    it('should open audit modal when clicking audit button', async () => {
      const _user = userEvent.setup();
      mockApi.getPlayers.mockResolvedValue({
        data: {
          success: true,
          data: [
            {
              id: 2,
              userId: 2,
              nickname: 'Pending Player',
              verificationStatus: 'pending' as const,
              hourlyRateCents: 5000,
              ratingAverage: 0,
              ratingCount: 0,
              skillTags: [],
              bio: 'Pending',
              mainGame: { id: 1, name: 'Test Game' },
              createdAt: '2024-01-01T00:00:00Z',
              user: { id: 2, name: 'Pending User', avatarUrl: '' },
            },
          ],
          pagination: { total: 1, page: 1, pageSize: 10 },
        },
      });

      renderWithProviders(<PlayerPage />);

      await waitFor(() => {
        expect(screen.getByText('Pending Player')).toBeInTheDocument();
      });

      const auditButton = screen.getByRole('button', { name: /审核/i });
      await user.click(auditButton);

      await waitFor(() => {
        expect(screen.getByText('审核陪玩师申请')).toBeInTheDocument();
      });
    });

    it('should approve player when clicking approve button', async () => {
      const _user = userEvent.setup();
      mockApi.getPlayers.mockResolvedValue({
        data: {
          success: true,
          data: [
            {
              id: 2,
              userId: 2,
              nickname: 'Pending Player',
              verificationStatus: 'pending' as const,
              hourlyRateCents: 5000,
              ratingAverage: 0,
              ratingCount: 0,
              skillTags: [],
              bio: 'Pending',
              mainGame: { id: 1, name: 'Test Game' },
              createdAt: '2024-01-01T00:00:00Z',
              user: { id: 2, name: 'Pending User', avatarUrl: '' },
            },
          ],
          pagination: { total: 1, page: 1, pageSize: 10 },
        },
      });

      renderWithProviders(<PlayerPage />);

      await waitFor(() => {
        expect(screen.getByText('Pending Player')).toBeInTheDocument();
      });

      const auditButton = screen.getByRole('button', { name: /审核/i });
      await user.click(auditButton);

      await waitFor(() => {
        expect(screen.getByText('审核陪玩师申请')).toBeInTheDocument();
      });

      const approveButton = screen.getByRole('button', { name: /通过/i });
      await user.click(approveButton);

      await waitFor(() => {
        expect(mockApi.updatePlayerVerification).toHaveBeenCalledWith(2, 'verified', '');
      });
    });

    it('should reject player when clicking reject button', async () => {
      const _user = userEvent.setup();
      mockApi.getPlayers.mockResolvedValue({
        data: {
          success: true,
          data: [
            {
              id: 2,
              userId: 2,
              nickname: 'Pending Player',
              verificationStatus: 'pending' as const,
              hourlyRateCents: 5000,
              ratingAverage: 0,
              ratingCount: 0,
              skillTags: [],
              bio: 'Pending',
              mainGame: { id: 1, name: 'Test Game' },
              createdAt: '2024-01-01T00:00:00Z',
              user: { id: 2, name: 'Pending User', avatarUrl: '' },
            },
          ],
          pagination: { total: 1, page: 1, pageSize: 10 },
        },
      });

      renderWithProviders(<PlayerPage />);

      await waitFor(() => {
        expect(screen.getByText('Pending Player')).toBeInTheDocument();
      });

      const auditButton = screen.getByRole('button', { name: /审核/i });
      await user.click(auditButton);

      await waitFor(() => {
        expect(screen.getByText('审核陪玩师申请')).toBeInTheDocument();
      });

      const rejectButton = screen.getByRole('button', { name: /拒绝/i });
      await user.click(rejectButton);

      await waitFor(() => {
        expect(mockApi.updatePlayerVerification).toHaveBeenCalledWith(2, 'rejected', '');
      });
    });
  });

  describe('Player Ban/Unban', () => {
    it('should show ban button for verified players', async () => {
      renderWithProviders(<PlayerPage />);

      await waitFor(() => {
        expect(screen.getByText('Test Player')).toBeInTheDocument();
      });

      expect(screen.getByRole('button', { name: /封禁/i })).toBeInTheDocument();
    });

    it('should ban player when confirming ban action', async () => {
      const _user = userEvent.setup();
      renderWithProviders(<PlayerPage />);

      await waitFor(() => {
        expect(screen.getByText('Test Player')).toBeInTheDocument();
      });

      const banButton = screen.getByRole('button', { name: /封禁/i });
      await user.click(banButton);

      const confirmButton = await screen.findByRole('button', { name: /确定/i });
      await user.click(confirmButton);

      await waitFor(() => {
        expect(mockApi.updatePlayerVerification).toHaveBeenCalledWith(1, 'rejected');
      });
    });

    it('should show unban button for rejected players', async () => {
      mockApi.getPlayers.mockResolvedValue({
        data: {
          success: true,
          data: [
            {
              id: 2,
              userId: 2,
              nickname: 'Banned Player',
              verificationStatus: 'rejected' as const,
              hourlyRateCents: 5000,
              ratingAverage: 0,
              ratingCount: 0,
              skillTags: [],
              bio: 'Banned',
              mainGame: { id: 1, name: 'Test Game' },
              createdAt: '2024-01-01T00:00:00Z',
              user: { id: 2, name: 'Banned User', avatarUrl: '' },
            },
          ],
          pagination: { total: 1, page: 1, pageSize: 10 },
        },
      });

      renderWithProviders(<PlayerPage />);

      await waitFor(() => {
        expect(screen.getByText('Banned Player')).toBeInTheDocument();
      });

      expect(screen.getByRole('button', { name: /解封/i })).toBeInTheDocument();
    });
  });

  describe('Batch Operations', () => {
    it('should show batch modify status button', async () => {
      renderWithProviders(<PlayerPage />);

      await waitFor(() => {
        expect(screen.getByRole('button', { name: /批量修改状态/i })).toBeInTheDocument();
      });
    });

    it('should show batch delete button', async () => {
      renderWithProviders(<PlayerPage />);

      await waitFor(() => {
        expect(screen.getByRole('button', { name: /批量删除/i })).toBeInTheDocument();
      });
    });

    it('should export player data', async () => {
      const _user = userEvent.setup();
      const { exportToCSV } = await import('@/utils/export');

      renderWithProviders(<PlayerPage />);

      await waitFor(() => {
        expect(screen.getByRole('button', { name: /导出数据/i })).toBeInTheDocument();
      });

      const exportButton = screen.getByRole('button', { name: /导出数据/i });
      await user.click(exportButton);

      await waitFor(() => {
        expect(mockApi.getPlayers).toHaveBeenCalledWith(
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
      renderWithProviders(<PlayerPage />);

      await waitFor(() => {
        expect(screen.getByText('陪玩师管理')).toBeInTheDocument();
      });

      const refreshButton = screen.getByRole('button', { name: /刷新/i });
      await user.click(refreshButton);

      await waitFor(() => {
        expect(mockApi.getPlayers).toHaveBeenCalledTimes(2);
      });
    });
  });

  describe('Accessibility', () => {
    it('should have proper ARIA labels', async () => {
      renderWithProviders(<PlayerPage />);

      await waitFor(() => {
        expect(screen.getByRole('heading', { name: '陪玩师管理' })).toBeInTheDocument();
      });
    });

    it('should be keyboard navigable', async () => {
      const _user = userEvent.setup();
      renderWithProviders(<PlayerPage />);

      await waitFor(() => {
        expect(screen.getByPlaceholderText('名称/ID')).toBeInTheDocument();
      });

      const searchInput = screen.getByPlaceholderText('名称/ID');
      searchInput.focus();

      expect(searchInput).toHaveFocus();
    });
  });
});
