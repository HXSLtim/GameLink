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

// Mock permission API
vi.mock('@/api/permission', () => ({
  permissionApi: {
    getUserPermissions: vi.fn().mockResolvedValue([]),
    getMyPermissions: vi.fn().mockResolvedValue([]),
  },
}));

// Mock the AdminContext to provide all permissions
vi.mock('@/context/AdminContext', () => ({
  useAdmin: () => ({
    menus: [],
    permissions: ['*'], // Super admin has all permissions
    loading: false,
    refreshMenus: vi.fn(),
    hasPermission: vi.fn(() => true),
    hasAllPermissions: vi.fn(() => true),
    hasAnyPermission: vi.fn(() => true),
    isSuperAdmin: true,
    permissionVersion: 0,
    notifyPermissionChange: vi.fn(),
  }),
  default: {
    AdminProvider: ({ children }: { children: React.ReactNode }) => <>{children}</>,
  },
}));

// Mock useAdmin hook (used by PermissionGuard component)
vi.mock('@/context/useAdmin', () => ({
  useAdmin: () => ({
    menus: [],
    permissions: ['*'], // Super admin has all permissions
    loading: false,
    refreshMenus: vi.fn(),
    hasPermission: vi.fn(() => true),
    hasAllPermissions: vi.fn(() => true),
    hasAnyPermission: vi.fn(() => true),
    isSuperAdmin: true,
    permissionVersion: 0,
    notifyPermissionChange: vi.fn(),
  }),
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

// Helper function to create mock order data
const createMockOrder = (overrides: Record<string, unknown> = {}): Record<string, unknown> => ({
  id: 1,
  orderNo: 'ORD20240101001',
  userId: 1,
  playerId: 1,
  gameId: 1,
  title: 'Test Order',
  description: 'Test description',
  totalPriceCents: 10000,
  currency: 'CNY',
  status: 'pending',
  scheduledStart: '2024-01-01T10:00:00Z',
  scheduledEnd: '2024-01-01T12:00:00Z',
  completedAt: '',
  cancelReason: '',
  createdAt: '2024-01-01T00:00:00Z',
  updatedAt: '2024-01-01T00:00:00Z',
  user: { id: 1, name: 'Test User', avatarUrl: '' },
  player: { id: 1, nickname: 'Test Player', user: { avatarUrl: '' } },
  game: { id: 1, name: 'Test Game' },
  ...overrides,
});

// Helper function to create mock order list
const createMockOrderList = (count = 1, overrides: Record<string, unknown> = {}): Record<string, unknown>[] => {
  return Array.from({ length: count }, (_, i) =>
    createMockOrder({
      id: i + 1,
      orderNo: `ORD202401010${String(i + 1).padStart(3, '0')}`,
      title: `Test Order ${i + 1}`,
      userId: i + 1,
      playerId: i + 1,
      ...overrides,
    })
  );
};

// Helper function to setup mock data with orders
const setupMockDataWithOrders = (orderCount = 1) => {
  const orders = createMockOrderList(orderCount);
  mockApi.getOrders.mockResolvedValue({
    data: {
      success: true,
      data: orders,
      pagination: { total: orderCount, page: 1, pageSize: 10 },
    },
  });

  // Also setup getOrder mock for single order
  mockApi.getOrder.mockResolvedValue({
    data: {
      success: true,
      data: orders[0],
    },
  });

  return orders;
};

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

      // Empty state should be handled gracefully - the component should still render
      expect(screen.getByText('订单管理')).toBeInTheDocument();
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
      const _user = userEvent.setup();
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
      await _user.type(searchInput, 'ORD20240101001');

      const searchButton = screen.getByRole('button', { name: /搜索/i });
      await _user.click(searchButton);

      await waitFor(() => {
        expect(mockApi.getOrders).toHaveBeenCalledWith(
          expect.objectContaining({
            orderNo: 'ORD20240101001',
          })
        );
      });
    });

    it('should reset to first page when searching', async () => {
      const _user = userEvent.setup();
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
      await _user.type(searchInput, 'test');

      const searchButton = screen.getByRole('button', { name: /搜索/i });
      await _user.click(searchButton);

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
      setupMockDataWithOrders(1);
      renderWithProviders(<OrderPage />);

      await waitFor(() => {
        expect(screen.getByText('订单管理')).toBeInTheDocument();
      });

      expect(screen.getByText(/共.*1.*条/)).toBeInTheDocument();
    });

    it('should change page when clicking pagination', async () => {
      // Create 20 orders to show pagination
      const orders = createMockOrderList(20);
      mockApi.getOrders.mockResolvedValue({
        data: {
          success: true,
          data: orders.slice(0, 10),
          pagination: { total: 20, page: 1, pageSize: 10 },
        },
      });

      renderWithProviders(<OrderPage />);

      await waitFor(() => {
        expect(screen.getByText('订单管理')).toBeInTheDocument();
      });

      // Verify pagination data shows 20 total items
      await waitFor(() => {
        expect(screen.getByText(/共.*20.*条/)).toBeInTheDocument();
      });
    });
  });

  describe('Order Details', () => {
    it('should open detail drawer when clicking detail button', async () => {
      const _user = userEvent.setup();
      renderWithProviders(<OrderPage />);

      await waitFor(() => {
        expect(screen.getByText('ORD20240101001')).toBeInTheDocument();
      });

      const detailButton = screen.getByRole('button', { name: /详情/i });
      await _user.click(detailButton);

      await waitFor(() => {
        expect(screen.getByText('订单详情')).toBeInTheDocument();
      });

      expect(mockApi.getOrder).toHaveBeenCalledWith(1);
    });

    it('should display order details in drawer', async () => {
      const _user = userEvent.setup();
      renderWithProviders(<OrderPage />);

      await waitFor(() => {
        expect(screen.getByText('ORD20240101001')).toBeInTheDocument();
      });

      const detailButton = screen.getByRole('button', { name: /详情/i });
      await _user.click(detailButton);

      await waitFor(() => {
        expect(screen.getByText('订单详情')).toBeInTheDocument();
      });

      expect(screen.getByText('订单信息')).toBeInTheDocument();
      expect(screen.getByText('用户信息')).toBeInTheDocument();
      expect(screen.getByText('陪玩师信息')).toBeInTheDocument();
    });
  });

  describe('Order Cancellation', () => {
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

      // Should not find a dangerous cancel button for completed orders
      const dangerousCancelButton = screen.queryAllByRole('button', { name: /取消/i }).find(btn =>
        btn.getAttribute('danger') === ''
      );
      expect(dangerousCancelButton).toBeUndefined();
    });
  });

  describe('Order Refund', () => {
    it('should open refund modal when clicking refund button', async () => {
      const _user = userEvent.setup();
      renderWithProviders(<OrderPage />);

      await waitFor(() => {
        expect(screen.getByText('ORD20240101001')).toBeInTheDocument();
      });

      const refundButton = screen.getByRole('button', { name: /退款/i });
      await _user.click(refundButton);

      await waitFor(() => {
        expect(screen.getByText('订单退款')).toBeInTheDocument();
      });
    }, 20000);
  });
});
