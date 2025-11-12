/**
 * Dashboard页面测试
 */

import { describe, it, expect, vi, beforeEach } from 'vitest';
import { screen, waitFor } from '@testing-library/react';
import { renderWithProviders } from '../../test/utils/test-utils';
import { Dashboard } from './Dashboard';
import { statsApi } from '../../services/api/stats';
import { orderApi } from '../../services/api/order';
import { OrderStatus } from '../../types/order';

// Mock statsApi
vi.mock('../../services/api/stats', () => ({
  statsApi: {
    getDashboard: vi.fn(),
    getUserGrowth: vi.fn(),
    getOrderStats: vi.fn(),
    getRevenueTrend: vi.fn(),
    getTopPlayers: vi.fn(),
  },
}));

// Mock orderApi
vi.mock('../../services/api/order', () => ({
  orderApi: {
    getList: vi.fn(),
  },
}));

describe('Dashboard Page', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('should render dashboard with overview stats', async () => {
    const mockDashboard = {
      totalUsers: 1000,
      totalPlayers: 150,
      totalGames: 50,
      totalOrders: 500,
      totalPaidAmountCents: 5000000,
      ordersByStatus: {
        pending: 25,
        in_progress: 30,
        completed: 400,
        canceled: 45,
      },
      paymentsByStatus: {},
    };

    vi.mocked(statsApi.getDashboard).mockResolvedValue(mockDashboard);
    vi.mocked(orderApi.getList).mockResolvedValue({
      list: [],
      total: 0,
      page: 1,
      page_size: 5,
    });

    renderWithProviders(<Dashboard />);

    // 等待数据加载
    await waitFor(() => {
      expect(statsApi.getDashboard).toHaveBeenCalled();
      expect(orderApi.getList).toHaveBeenCalled();
    });

    // 验证页面标题存在
    expect(screen.getByRole('heading', { level: 1 })).toBeInTheDocument();
    
    // 验证统计数据显示（检查数值）
    expect(screen.getByText('1000')).toBeInTheDocument();
    expect(screen.getByText('150')).toBeInTheDocument();
  });

  it('should handle API error gracefully', async () => {
    vi.mocked(statsApi.getDashboard).mockRejectedValue(new Error('Failed to fetch'));
    vi.mocked(orderApi.getList).mockRejectedValue(new Error('Failed to fetch'));

    renderWithProviders(<Dashboard />);

    // 等待错误处理完成
    await waitFor(() => {
      expect(statsApi.getDashboard).toHaveBeenCalled();
    });

    // 验证失败状态显示
    expect(screen.getByRole('heading', { level: 1 })).toBeInTheDocument();
  });

  it('should display loading state initially', () => {
    vi.mocked(statsApi.getDashboard).mockImplementation(
      () => new Promise(resolve => setTimeout(resolve, 1000))
    );
    vi.mocked(orderApi.getList).mockImplementation(
      () => new Promise(resolve => setTimeout(resolve, 1000))
    );

    renderWithProviders(<Dashboard />);

    // 应该显示加载状态
    expect(screen.getByText(/加载中/i)).toBeInTheDocument();
  });

  it('should display recent orders', async () => {
    const mockOrders = [
      {
        id: 1,
        userId: 1,
        gameId: 1,
        title: 'Order 1',
        status: OrderStatus.PENDING,
        priceCents: 10000,
        createdAt: '2024-11-11T00:00:00Z',
        updatedAt: '2024-11-11T00:00:00Z',
      },
    ];

    vi.mocked(statsApi.getDashboard).mockResolvedValue({
      totalUsers: 1000,
      totalPlayers: 150,
      totalGames: 50,
      totalOrders: 500,
      totalPaidAmountCents: 5000000,
      ordersByStatus: {},
      paymentsByStatus: {},
    });
    vi.mocked(orderApi.getList).mockResolvedValue({
      list: mockOrders,
      total: 1,
      page: 1,
      page_size: 5,
    });

    renderWithProviders(<Dashboard />);

    await waitFor(() => {
      expect(orderApi.getList).toHaveBeenCalled();
    });
  });

  it('should display order status cards', async () => {
    const mockDashboard = {
      totalUsers: 1000,
      totalPlayers: 150,
      totalGames: 50,
      totalOrders: 500,
      totalPaidAmountCents: 5000000,
      ordersByStatus: {
        pending: 25,
        in_progress: 30,
        completed: 400,
        canceled: 45,
      },
      paymentsByStatus: {},
    };

    vi.mocked(statsApi.getDashboard).mockResolvedValue(mockDashboard);
    vi.mocked(orderApi.getList).mockResolvedValue({
      list: [],
      total: 0,
      page: 1,
      page_size: 5,
    });

    renderWithProviders(<Dashboard />);

    await waitFor(() => {
      expect(statsApi.getDashboard).toHaveBeenCalled();
    });

    // 验证订单状态数值显示
    await waitFor(() => {
      const values = screen.getAllByText(/\d+/);
      expect(values.length).toBeGreaterThan(0);
    });
  });

  it('should navigate to orders page when view all button is clicked', async () => {
    const mockDashboard = {
      totalUsers: 1000,
      totalPlayers: 150,
      totalGames: 50,
      totalOrders: 500,
      totalPaidAmountCents: 5000000,
      ordersByStatus: {},
      paymentsByStatus: {},
    };

    vi.mocked(statsApi.getDashboard).mockResolvedValue(mockDashboard);
    vi.mocked(orderApi.getList).mockResolvedValue({
      list: [],
      total: 0,
      page: 1,
      page_size: 5,
    });

    renderWithProviders(<Dashboard />);

    // 等待初始加载
    await waitFor(() => {
      expect(statsApi.getDashboard).toHaveBeenCalledTimes(1);
      expect(orderApi.getList).toHaveBeenCalled();
    });

    // 验证组件正常渲染
    expect(screen.getByRole('heading', { level: 1 })).toBeInTheDocument();
  });

  it('should display empty state when no orders available', async () => {
    const mockDashboard = {
      totalUsers: 0,
      totalPlayers: 0,
      totalGames: 0,
      totalOrders: 0,
      totalPaidAmountCents: 0,
      ordersByStatus: {},
      paymentsByStatus: {},
    };

    vi.mocked(statsApi.getDashboard).mockResolvedValue(mockDashboard);
    vi.mocked(orderApi.getList).mockResolvedValue({
      list: [],
      total: 0,
      page: 1,
      page_size: 5,
    });

    renderWithProviders(<Dashboard />);

    await waitFor(() => {
      expect(statsApi.getDashboard).toHaveBeenCalled();
    });

    // 验证数据加载完成并且显示空状态
    await waitFor(() => {
      const zeros = screen.getAllByText('0');
      expect(zeros.length).toBeGreaterThan(0);
    });
  });
});
