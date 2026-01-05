/**
 * Dashboard Page Tests
 *
 * Tests for Dashboard page component including:
 * - Successful data loading
 * - Statistics display
 * - Charts rendering
 * - Recent orders display
 * - Top players display
 * - Time range selection
 */

import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { screen, waitFor, fireEvent } from '@testing-library/react';
import Dashboard from './index';
import { renderWithProviders, resetAllMocks, flushPromises } from '@/testutils';

// Mock the adminApi module using vi.hoisted
const { mockApi, mockMessage } = vi.hoisted(() => ({
  mockApi: {
    getDashboardStats: vi.fn(),
    getRevenueTrend: vi.fn(),
    getUserGrowth: vi.fn(),
    getOrders: vi.fn(),
    getTopPlayers: vi.fn(),
  },
  mockMessage: {
    success: vi.fn(),
    error: vi.fn(),
    warning: vi.fn(),
    info: vi.fn(),
    loading: vi.fn(),
  },
}));

vi.mock('@/api/admin', () => ({
  adminApi: mockApi,
}));

// Mock antd App.useApp
vi.mock('antd', async () => {
  const actual = await vi.importActual('antd');
  return {
    ...actual,
    App: {
      ...((actual as Record<string, unknown>).App as Record<string, unknown>),
      useApp: () => ({ message: mockMessage }),
    },
  };
});

// Mock recharts to avoid rendering issues in tests
vi.mock('recharts', () => ({
  ResponsiveContainer: ({ children }: { children: React.ReactNode }) => <div data-testid="responsive-container">{children}</div>,
  LineChart: ({ children }: { children: React.ReactNode }) => <div data-testid="line-chart">{children}</div>,
  Line: () => null,
  XAxis: () => null,
  YAxis: () => null,
  CartesianGrid: () => null,
  Tooltip: () => null,
  Legend: () => null,
  PieChart: ({ children }: { children: React.ReactNode }) => <div data-testid="pie-chart">{children}</div>,
  Pie: () => null,
  Cell: () => null,
}));

// Helper function to create mock stats
const createMockStats = (overrides: Record<string, unknown> = {}): Record<string, unknown> => ({
  totalOrders: 100,
  totalPaidAmountCents: 1000000,
  totalUsers: 500,
  totalPlayers: 50,
  ordersByStatus: {
    pending: 10,
    confirmed: 20,
    in_progress: 30,
    completed: 35,
    canceled: 5,
  },
  paymentsByStatus: {
    pending: 15,
    paid: 80,
    refunded: 5,
  },
  ...overrides,
});

// Helper function to create mock trend data
const createMockTrendData = (): Record<string, unknown>[] => [
  { date: '2024-01-01', value: 1000 },
  { date: '2024-01-02', value: 1500 },
  { date: '2024-01-03', value: 1200 },
];

// Helper function to create mock order
const createMockOrder = (overrides: Record<string, unknown> = {}): Record<string, unknown> => ({
  id: 1,
  orderNo: 'ORD20240101001',
  userId: 100,
  totalPriceCents: 10000,
  status: 'completed',
  createdAt: '2024-01-01T00:00:00Z',
  ...overrides,
});

// Helper function to create mock top player
const createMockTopPlayer = (overrides: Record<string, unknown> = {}): Record<string, unknown> => ({
  playerId: 1,
  nickname: '热门陪玩师',
  ratingAverage: 4.8,
  ratingCount: 100,
  ...overrides,
});

describe('Dashboard', () => {
  beforeEach(() => {
    resetAllMocks();
    mockMessage.success.mockClear();
    mockMessage.error.mockClear();
    localStorage.setItem('token', 'test-token');
    localStorage.setItem('user_role', 'admin');

    // Set default mock return values
    mockApi.getDashboardStats.mockResolvedValue({
      data: { data: createMockStats() },
    });
    mockApi.getRevenueTrend.mockResolvedValue({
      data: { data: createMockTrendData() },
    });
    mockApi.getUserGrowth.mockResolvedValue({
      data: { data: createMockTrendData() },
    });
    mockApi.getOrders.mockResolvedValue({
      data: { data: [createMockOrder()] },
    });
    mockApi.getTopPlayers.mockResolvedValue({
      data: { data: [createMockTopPlayer()] },
    });
  });

  afterEach(() => {
    vi.clearAllMocks();
  });

  describe('Successful Data Loading', () => {
    it('should render dashboard page successfully', async () => {
      renderWithProviders(<Dashboard />);

      await waitFor(() => {
        expect(screen.getByText('仪表盘')).toBeInTheDocument();
      });

      expect(mockApi.getDashboardStats).toHaveBeenCalled();
    });

    it('should display page subtitle', async () => {
      renderWithProviders(<Dashboard />);

      await waitFor(() => {
        expect(screen.getByText('系统性能概览')).toBeInTheDocument();
      });
    });

    it('should call all data APIs on mount', async () => {
      renderWithProviders(<Dashboard />);

      await waitFor(() => {
        expect(mockApi.getDashboardStats).toHaveBeenCalled();
        expect(mockApi.getRevenueTrend).toHaveBeenCalled();
        expect(mockApi.getUserGrowth).toHaveBeenCalled();
        expect(mockApi.getOrders).toHaveBeenCalled();
        expect(mockApi.getTopPlayers).toHaveBeenCalled();
      });
    });
  });

  describe('Statistics Cards', () => {
    it('should display total orders card', async () => {
      renderWithProviders(<Dashboard />);

      await waitFor(() => {
        expect(screen.getByText('总订单数')).toBeInTheDocument();
      });
    });

    it('should display total revenue card', async () => {
      renderWithProviders(<Dashboard />);

      await waitFor(() => {
        expect(screen.getByText('交易总额')).toBeInTheDocument();
      });
    });

    it('should display total users card', async () => {
      renderWithProviders(<Dashboard />);

      await waitFor(() => {
        expect(screen.getByText('总用户数')).toBeInTheDocument();
      });
    });

    it('should display total players card', async () => {
      renderWithProviders(<Dashboard />);

      await waitFor(() => {
        expect(screen.getByText('总陪玩数')).toBeInTheDocument();
      });
    });
  });

  describe('Charts Display', () => {
    it('should display order status distribution chart', async () => {
      renderWithProviders(<Dashboard />);

      await waitFor(() => {
        expect(screen.getByText('订单状态分布')).toBeInTheDocument();
      });
    });

    it('should display payment status distribution chart', async () => {
      renderWithProviders(<Dashboard />);

      await waitFor(() => {
        expect(screen.getByText('支付状态分布')).toBeInTheDocument();
      });
    });

    it('should display revenue trend chart', async () => {
      renderWithProviders(<Dashboard />);

      await waitFor(() => {
        expect(screen.getByText('收入趋势 (近7天)')).toBeInTheDocument();
      });
    });

    it('should display user growth chart', async () => {
      renderWithProviders(<Dashboard />);

      await waitFor(() => {
        expect(screen.getByText('用户增长 (近7天)')).toBeInTheDocument();
      });
    });
  });

  describe('Recent Orders Section', () => {
    it('should display recent orders section', async () => {
      renderWithProviders(<Dashboard />);

      await waitFor(() => {
        expect(screen.getByText('最新订单')).toBeInTheDocument();
      });
    });

    it('should display view all link', async () => {
      renderWithProviders(<Dashboard />);

      await waitFor(() => {
        expect(screen.getByText('查看全部')).toBeInTheDocument();
      });
    });
  });

  describe('Top Players Section', () => {
    it('should display top players section', async () => {
      renderWithProviders(<Dashboard />);

      await waitFor(() => {
        expect(screen.getByText('热门陪玩')).toBeInTheDocument();
      });
    });
  });

  describe('Time Range Selection', () => {
    it('should display time range selector', async () => {
      renderWithProviders(<Dashboard />);

      await waitFor(() => {
        expect(screen.getByText('近7天')).toBeInTheDocument();
      });
    });

    it('should have 30 days option', async () => {
      renderWithProviders(<Dashboard />);

      await waitFor(() => {
        expect(screen.getByText('仪表盘')).toBeInTheDocument();
      });

      // Click to open select dropdown
      const select = screen.getByText('近7天');
      fireEvent.mouseDown(select);

      await waitFor(() => {
        expect(screen.getByText('近30天')).toBeInTheDocument();
      });
    });
  });

  describe('Error Handling', () => {
    it('should display error message when API fails', async () => {
      mockApi.getDashboardStats.mockRejectedValue(new Error('Network error'));

      renderWithProviders(<Dashboard />);

      await waitFor(() => {
        expect(mockMessage.error).toHaveBeenCalledWith('加载仪表盘数据失败');
      });
    });
  });

  describe('Empty State', () => {
    it('should handle empty stats gracefully', async () => {
      mockApi.getDashboardStats.mockResolvedValue({
        data: { data: {} },
      });

      renderWithProviders(<Dashboard />);

      await waitFor(() => {
        expect(screen.getByText('仪表盘')).toBeInTheDocument();
      });

      // Should display 0 for missing values
      expect(mockApi.getDashboardStats).toHaveBeenCalled();
    });

    it('should handle empty orders list', async () => {
      mockApi.getOrders.mockResolvedValue({
        data: { data: [] },
      });

      renderWithProviders(<Dashboard />);

      await waitFor(() => {
        expect(screen.getByText('最新订单')).toBeInTheDocument();
      });
    });

    it('should handle empty top players list', async () => {
      mockApi.getTopPlayers.mockResolvedValue({
        data: { data: [] },
      });

      renderWithProviders(<Dashboard />);

      await waitFor(() => {
        expect(screen.getByText('热门陪玩')).toBeInTheDocument();
      });
    });
  });

  describe('Accessibility', () => {
    it('should have proper page title', async () => {
      renderWithProviders(<Dashboard />);

      await waitFor(() => {
        expect(screen.getByText('仪表盘')).toBeInTheDocument();
      });
    });
  });
});
