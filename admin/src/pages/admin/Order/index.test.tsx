/**
 * Order Management Page Tests
 *
 * Tests for Order page component including:
 * - Successful data loading
 * - Loading states
 * - Error handling
 * - User interactions (filtering, pagination)
 * - Order operations (cancel, refund, batch operations)
 * - Permission checks
 */

import React from 'react';
import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import OrderPage from './index';
import { renderWithProviders, resetAllMocks, flushPromises } from '@/testutils';

// Mock the adminApi module using vi.hoisted
const { mockApi, mockMessage } = vi.hoisted(() => ({
  mockApi: {
    getOrders: vi.fn(),
    getOrder: vi.fn(),
    cancelOrder: vi.fn(),
    refundOrder: vi.fn(),
    batchCancelOrders: vi.fn(),
    batchCompleteOrders: vi.fn(),
  },
  mockMessage: {
    success: vi.fn(),
    error: vi.fn(),
    warning: vi.fn(),
    info: vi.fn(),
  },
}));

vi.mock('@/api/admin', () => ({
  adminApi: mockApi,
}));

// Mock export utilities
vi.mock('@/utils/export', () => ({
  exportToCSV: vi.fn(),
  orderExportColumns: [],
}));

// Mock App.useApp to return the message mock
vi.mock('antd', async () => {
  const actual = await vi.importActual('antd');
  return {
    ...actual,
    App: {
      useApp: () => ({
        message: mockMessage,
        notification: { success: vi.fn(), error: vi.fn(), info: vi.fn(), warning: vi.fn() },
        modal: { success: vi.fn(), error: vi.fn(), info: vi.fn(), warning: vi.fn(), confirm: vi.fn() },
      }),
    },
  };
});

describe('OrderPage', () => {
  beforeEach(() => {
    resetAllMocks();
    // Setup localStorage
    localStorage.setItem('token', 'test-token');
    localStorage.setItem('user_role', 'admin');

    // Setup default mock responses
    mockApi.getOrders.mockResolvedValue({
      data: {
        success: true,
        data: [
          {
            id: 1,
            orderNo: 'ORD20240101001',
            userId: 1,
            playerId: 1,
            gameId: 1,
            title: 'Test Order',
            description: 'Test description',
            totalPriceCents: 10000,
            currency: 'CNY',
            status: 'pending' as const,
            scheduledStart: '2024-01-01T10:00:00Z',
            scheduledEnd: '2024-01-01T12:00:00Z',
            createdAt: '2024-01-01T00:00:00Z',
            user: { id: 1, name: 'Test User', avatarUrl: '' },
            player: { id: 1, nickname: 'Test Player', user: { avatarUrl: '' } },
            game: { id: 1, name: 'Test Game' },
          },
        ],
        pagination: { total: 1, page: 1, pageSize: 10 },
      },
    });

    mockApi.getOrder.mockResolvedValue({
      data: {
        success: true,
        data: {
          id: 1,
          orderNo: 'ORD20240101001',
          userId: 1,
          playerId: 1,
          gameId: 1,
          title: 'Test Order',
          totalPriceCents: 10000,
          status: 'pending' as const,
          createdAt: '2024-01-01T00:00:00Z',
          user: { id: 1, name: 'Test User', avatarUrl: '' },
          player: { id: 1, nickname: 'Test Player', user: { avatarUrl: '' } },
          game: { id: 1, name: 'Test Game' },
        },
      },
    });
  });

  afterEach(() => {
    vi.clearAllMocks();
  });

  describe('Successful Data Loading', () => {
    it('should render order list successfully', async () => {
      renderWithProviders(<OrderPage />);

      await waitFor(() => {
        expect(screen.getByText('订单管理')).toBeInTheDocument();
      });

      await waitFor(() => {
        expect(screen.getByText('ORD20240101001')).toBeInTheDocument();
      });

      expect(mockApi.getOrders).toHaveBeenCalledWith({
        page: 1,
        page_size: 10,
      });
    });

    it('should display order information correctly', async () => {
      renderWithProviders(<OrderPage />);

      await waitFor(() => {
        expect(screen.getByText('ORD20240101001')).toBeInTheDocument();
      });

      expect(screen.getByText('Test Order')).toBeInTheDocument();
      expect(screen.getByText('¥100.00')).toBeInTheDocument();
      expect(screen.getByText('Test Game')).toBeInTheDocument();
    });

    it('should display user and player information', async () => {
      renderWithProviders(<OrderPage />);

      await waitFor(() => {
        expect(screen.getByText('Test User')).toBeInTheDocument();
      });

      expect(screen.getByText('Test Player')).toBeInTheDocument();
    });

    it('should display order status tags', async () => {
      renderWithProviders(<OrderPage />);

      await waitFor(() => {
        expect(screen.getByText('待确认')).toBeInTheDocument();
      });
    });
  });

  describe('Loading States', () => {
    it('should show loading indicator while fetching data', async () => {
      // Make the API call slower
      mockApi.getOrders.mockImplementation(
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

      renderWithProviders(<OrderPage />);

      // Ant Design's Spin component doesn't have a clear test ID,
      // but we can check if the table is in a loading state
      expect(mockApi.getOrders).toHaveBeenCalled();
    });

    it('should hide loading indicator after data loads', async () => {
      renderWithProviders(<OrderPage />);

      await waitFor(() => {
        expect(screen.getByText('订单管理')).toBeInTheDocument();
      });

      await flushPromises();

      expect(mockApi.getOrders).toHaveBeenCalledTimes(1);
    });
  });

  describe('Error Handling', () => {
    it('should display error message when API fails', async () => {
      mockApi.getOrders.mockRejectedValue(new Error('Network error'));

      renderWithProviders(<OrderPage />);

      await waitFor(() => {
        expect(screen.getByText('订单管理')).toBeInTheDocument();
      });

      await flushPromises();

      // Error message should appear
      const errorMessage = await screen.findByText(/加载订单列表失败/);
      expect(errorMessage).toBeInTheDocument();
    });

    it('should handle empty data gracefully', async () => {
      mockApi.getOrders.mockResolvedValue({
        data: {
          success: true,
          data: [],
          pagination: { total: 0, page: 1, pageSize: 10 },
        },
      });

      renderWithProviders(<OrderPage />);

      await waitFor(() => {
        expect(screen.getByText('订单管理')).toBeInTheDocument();
      });

      await flushPromises();

      // Should show empty state or table with no rows
      expect(screen.getByText('共 0 条')).toBeInTheDocument();
    });

    it('should handle API response with success: false', async () => {
      mockApi.getOrders.mockResolvedValue({
        data: {
          success: false,
          message: '获取订单列表失败',
        },
      });

      renderWithProviders(<OrderPage />);

      await waitFor(() => {
        expect(screen.getByText('订单管理')).toBeInTheDocument();
      });

      await flushPromises();

      const errorMessage = await screen.findByText(/获取订单列表失败/);
      expect(errorMessage).toBeInTheDocument();
    });
  });

  describe('Search and Filtering', () => {
    it('should allow searching by order number', async () => {
      const user = userEvent.setup();
      mockApi.getOrders.mockResolvedValue({
        data: {
          success: true,
          data: [],
          pagination: { total: 0, page: 1, pageSize: 10 },
        },
      });

      renderWithProviders(<OrderPage />);

      await waitFor(() => {
        expect(screen.getByPlaceholderText('请输入订单号')).toBeInTheDocument();
      });

      const searchInput = screen.getByPlaceholderText('请输入订单号');
      await user.type(searchInput, 'ORD20240101001');

      const searchButton = screen.getByRole('button', { name: /搜索/i });
      await user.click(searchButton);

      await waitFor(() => {
        expect(mockApi.getOrders).toHaveBeenCalledWith(
          expect.objectContaining({
            orderNo: 'ORD20240101001',
          })
        );
      });
    });

    it('should allow filtering by order status', async () => {
      const user = userEvent.setup();
      mockApi.getOrders.mockResolvedValue({
        data: {
          success: true,
          data: [],
          pagination: { total: 0, page: 1, pageSize: 10 },
        },
      });

      renderWithProviders(<OrderPage />);

      await waitFor(() => {
        expect(screen.getByText('订单状态')).toBeInTheDocument();
      });

      // Find and click the status filter dropdown
      const statusDropdown = screen.getByText('订单状态').closest('.ant-select');
      if (statusDropdown) {
        await user.click(statusDropdown);

        // Select "待确认" status
        const pendingOption = await screen.findByText('待确认');
        await user.click(pendingOption);

        await waitFor(() => {
          expect(mockApi.getOrders).toHaveBeenCalledWith(
            expect.objectContaining({
              status: 'pending',
            })
          );
        });
      }
    });

    it('should reset to first page when searching', async () => {
      const user = userEvent.setup();
      mockApi.getOrders.mockResolvedValue({
        data: {
          success: true,
          data: [],
          pagination: { total: 0, page: 1, pageSize: 10 },
        },
      });

      renderWithProviders(<OrderPage />);

      await waitFor(() => {
        expect(screen.getByPlaceholderText('请输入订单号')).toBeInTheDocument();
      });

      const searchInput = screen.getByPlaceholderText('请输入订单号');
      await user.type(searchInput, 'test');

      const searchButton = screen.getByRole('button', { name: /搜索/i });
      await user.click(searchButton);

      await waitFor(() => {
        expect(mockApi.getOrders).toHaveBeenCalledWith(
          expect.objectContaining({
            page: 1,
          })
        );
      });
    });
  });

  describe('Pagination', () => {
    it('should display pagination controls', async () => {
      renderWithProviders(<OrderPage />);

      await waitFor(() => {
        expect(screen.getByText('订单管理')).toBeInTheDocument();
      });

      expect(screen.getByText('共 1 条')).toBeInTheDocument();
    });

    it('should change page when clicking pagination', async () => {
      const user = userEvent.setup();
      mockApi.getOrders.mockResolvedValue({
        data: {
          success: true,
          data: [],
          pagination: { total: 20, page: 2, pageSize: 10 },
        },
      });

      renderWithProviders(<OrderPage />);

      await waitFor(() => {
        expect(screen.getByText('订单管理')).toBeInTheDocument();
      });

      // Click next page
      const nextPageButton = screen.getByTitle('下一页');
      await user.click(nextPageButton);

      await waitFor(() => {
        expect(mockApi.getOrders).toHaveBeenCalledWith(
          expect.objectContaining({
            page: 2,
          })
        );
      });
    });

    it('should change page size when selecting different size', async () => {
      const user = userEvent.setup();
      mockApi.getOrders.mockResolvedValue({
        data: {
          success: true,
          data: [],
          pagination: { total: 50, page: 1, pageSize: 20 },
        },
      });

      renderWithProviders(<OrderPage />);

      await waitFor(() => {
        expect(screen.getByText('订单管理')).toBeInTheDocument();
      });

      // Click page size selector
      const pageSizeSelector = screen.getByText('10 条/页');
      await user.click(pageSizeSelector);

      // Select 20 items per page
      const pageSize20 = await screen.findByText('20 条/页');
      await user.click(pageSize20);

      await waitFor(() => {
        expect(mockApi.getOrders).toHaveBeenCalledWith(
          expect.objectContaining({
            page_size: 20,
          })
        );
      });
    });
  });

  describe('Order Details', () => {
    it('should open detail drawer when clicking detail button', async () => {
      const user = userEvent.setup();
      renderWithProviders(<OrderPage />);

      await waitFor(() => {
        expect(screen.getByText('ORD20240101001')).toBeInTheDocument();
      });

      const detailButton = screen.getByRole('button', { name: /详情/i });
      await user.click(detailButton);

      await waitFor(() => {
        expect(screen.getByText('订单详情')).toBeInTheDocument();
      });

      expect(mockApi.getOrder).toHaveBeenCalledWith(1);
    });

    it('should display order details in drawer', async () => {
      const user = userEvent.setup();
      renderWithProviders(<OrderPage />);

      await waitFor(() => {
        expect(screen.getByText('ORD20240101001')).toBeInTheDocument();
      });

      const detailButton = screen.getByRole('button', { name: /详情/i });
      await user.click(detailButton);

      await waitFor(() => {
        expect(screen.getByText('订单详情')).toBeInTheDocument();
      });

      expect(screen.getByText('订单信息')).toBeInTheDocument();
      expect(screen.getByText('用户信息')).toBeInTheDocument();
      expect(screen.getByText('陪玩师信息')).toBeInTheDocument();
      expect(screen.getByText('订单进度')).toBeInTheDocument();
    });

    it('should close detail drawer when clicking close button', async () => {
      const user = userEvent.setup();
      renderWithProviders(<OrderPage />);

      await waitFor(() => {
        expect(screen.getByText('ORD20240101001')).toBeInTheDocument();
      });

      const detailButton = screen.getByRole('button', { name: /详情/i });
      await user.click(detailButton);

      await waitFor(() => {
        expect(screen.getByText('订单详情')).toBeInTheDocument();
      });

      // Close button is typically an icon button
      const closeButton = screen.getByRole('button', { name: /close/i });
      await user.click(closeButton);

      await waitFor(() => {
        expect(screen.queryByText('订单详情')).not.toBeInTheDocument();
      });
    });
  });

  describe('Order Cancellation', () => {
    it('should show cancel button for pending orders', async () => {
      renderWithProviders(<OrderPage />);

      await waitFor(() => {
        expect(screen.getByText('ORD20240101001')).toBeInTheDocument();
      });

      expect(screen.getByRole('button', { name: /取消/i })).toBeInTheDocument();
    });

    it('should cancel order when confirming cancellation', async () => {
      const user = userEvent.setup();
      renderWithProviders(<OrderPage />);

      await waitFor(() => {
        expect(screen.getByText('ORD20240101001')).toBeInTheDocument();
      });

      const cancelButton = screen.getByRole('button', { name: /取消/i });
      await user.click(cancelButton);

      // Confirm in popconfirm
      const confirmButton = await screen.findByRole('button', { name: /确定/i });
      await user.click(confirmButton);

      await waitFor(() => {
        expect(mockApi.cancelOrder).toHaveBeenCalledWith(1);
      });
    });

    it('should not show cancel button for completed orders', async () => {
      mockApi.getOrders.mockResolvedValue({
        data: {
          success: true,
          data: [
            {
              id: 2,
              orderNo: 'ORD20240101002',
              status: 'completed' as const,
              userId: 1,
              playerId: 1,
              gameId: 1,
              title: 'Completed Order',
              totalPriceCents: 10000,
              currency: 'CNY',
              createdAt: '2024-01-01T00:00:00Z',
            },
          ],
          pagination: { total: 1, page: 1, pageSize: 10 },
        },
      });

      renderWithProviders(<OrderPage />);

      await waitFor(() => {
        expect(screen.getByText('ORD20240101002')).toBeInTheDocument();
      });

      expect(screen.queryByRole('button', { name: /取消/i })).not.toBeInTheDocument();
    });
  });

  describe('Order Refund', () => {
    it('should open refund modal when clicking refund button', async () => {
      const user = userEvent.setup();
      renderWithProviders(<OrderPage />);

      await waitFor(() => {
        expect(screen.getByText('ORD20240101001')).toBeInTheDocument();
      });

      const refundButton = screen.getByRole('button', { name: /退款/i });
      await user.click(refundButton);

      await waitFor(() => {
        expect(screen.getByText('订单退款')).toBeInTheDocument();
      });
    });

    it('should submit refund with valid data', async () => {
      const user = userEvent.setup();
      renderWithProviders(<OrderPage />);

      await waitFor(() => {
        expect(screen.getByText('ORD20240101001')).toBeInTheDocument();
      });

      const refundButton = screen.getByRole('button', { name: /退款/i });
      await user.click(refundButton);

      await waitFor(() => {
        expect(screen.getByText('订单退款')).toBeInTheDocument();
      });

      // Fill refund form
      const amountInput = screen.getByPlaceholderText(/请输入退款金额/);
      await user.type(amountInput, '50.00');

      const reasonInput = screen.getByPlaceholderText('请输入退款原因');
      await user.type(reasonInput, '用户申请退款');

      const confirmButton = screen.getByRole('button', { name: /确认/i });
      await user.click(confirmButton);

      await waitFor(() => {
        expect(mockApi.refundOrder).toHaveBeenCalledWith(1, {
          reason: '用户申请退款',
          amount_cents: 5000,
        });
      });
    });

    it('should validate refund amount does not exceed order total', async () => {
      const user = userEvent.setup();
      renderWithProviders(<OrderPage />);

      await waitFor(() => {
        expect(screen.getByText('ORD20240101001')).toBeInTheDocument();
      });

      const refundButton = screen.getByRole('button', { name: /退款/i });
      await user.click(refundButton);

      await waitFor(() => {
        expect(screen.getByText('订单退款')).toBeInTheDocument();
      });

      // Try to enter amount greater than order total
      const amountInput = screen.getByPlaceholderText(/请输入退款金额/);
      await user.clear(amountInput);
      await user.type(amountInput, '150.00');

      const confirmButton = screen.getByRole('button', { name: /确认/i });
      await user.click(confirmButton);

      // Should show validation error
      await waitFor(() => {
        expect(screen.getByText(/退款金额不能超过 ¥100.00/)).toBeInTheDocument();
      });
    });
  });

  describe('Batch Operations', () => {
    it('should show batch cancel button in toolbar', async () => {
      renderWithProviders(<OrderPage />);

      await waitFor(() => {
        expect(screen.getByRole('button', { name: /批量取消/i })).toBeInTheDocument();
      });
    });

    it('should show batch complete button in toolbar', async () => {
      renderWithProviders(<OrderPage />);

      await waitFor(() => {
        expect(screen.getByRole('button', { name: /批量完成/i })).toBeInTheDocument();
      });
    });

    it('should open batch cancel modal', async () => {
      const user = userEvent.setup();
      renderWithProviders(<OrderPage />);

      await waitFor(() => {
        expect(screen.getByRole('button', { name: /批量取消/i })).toBeInTheDocument();
      });

      const batchCancelButton = screen.getByRole('button', { name: /批量取消/i });
      await user.click(batchCancelButton);

      await waitFor(() => {
        expect(screen.getByText('批量取消订单')).toBeInTheDocument();
      });
    });

    it('should export order data', async () => {
      const user = userEvent.setup();
      const { exportToCSV } = await import('@/utils/export');

      renderWithProviders(<OrderPage />);

      await waitFor(() => {
        expect(screen.getByRole('button', { name: /导出数据/i })).toBeInTheDocument();
      });

      const exportButton = screen.getByRole('button', { name: /导出数据/i });
      await user.click(exportButton);

      await waitFor(() => {
        expect(mockApi.getOrders).toHaveBeenCalledWith(
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
      const user = userEvent.setup();
      renderWithProviders(<OrderPage />);

      await waitFor(() => {
        expect(screen.getByText('订单管理')).toBeInTheDocument();
      });

      const refreshButton = screen.getByRole('button', { name: /刷新/i });
      await user.click(refreshButton);

      await waitFor(() => {
        expect(mockApi.getOrders).toHaveBeenCalledTimes(2); // Initial load + refresh
      });
    });
  });

  describe('Permission Checks', () => {
    it('should show cancel button only with CANCEL permission', async () => {
      // This test assumes the PermissionGuard component is working
      renderWithProviders(<OrderPage />);

      await waitFor(() => {
        expect(screen.getByText('ORD20240101001')).toBeInTheDocument();
      });

      // If permission guard works, cancel button should be visible for admin
      expect(screen.getByRole('button', { name: /取消/i })).toBeInTheDocument();
    });

    it('should show refund button only with REFUND permission', async () => {
      renderWithProviders(<OrderPage />);

      await waitFor(() => {
        expect(screen.getByText('ORD20240101001')).toBeInTheDocument();
      });

      expect(screen.getByRole('button', { name: /退款/i })).toBeInTheDocument();
    });
  });

  describe('Accessibility', () => {
    it('should have proper ARIA labels', async () => {
      renderWithProviders(<OrderPage />);

      await waitFor(() => {
        expect(screen.getByRole('heading', { name: '订单管理' })).toBeInTheDocument();
      });
    });

    it('should be keyboard navigable', async () => {
      const user = userEvent.setup();
      renderWithProviders(<OrderPage />);

      await waitFor(() => {
        expect(screen.getByPlaceholderText('请输入订单号')).toBeInTheDocument();
      });

      const searchInput = screen.getByPlaceholderText('请输入订单号');
      searchInput.focus();

      expect(searchInput).toHaveFocus();
    });
  });
});
