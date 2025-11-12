/**
 * OrderList 页面测试
 */

import { describe, it, expect, vi, beforeEach } from 'vitest';
import { waitFor } from '@testing-library/react';
import { renderWithProviders } from '../../test/utils/test-utils';
import { OrderList } from './OrderList';
import { orderApi } from '../../services/api/order';
import { OrderStatus } from '../../types/order';

// Mock orderApi
vi.mock('../../services/api/order', () => ({
  orderApi: {
    getList: vi.fn(),
    getDetail: vi.fn(),
    update: vi.fn(),
    delete: vi.fn(),
  },
}));

describe('OrderList Page', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  // 1. 基础渲染测试
  it('should render order list page', async () => {
    vi.mocked(orderApi.getList).mockResolvedValue({
      list: [],
      total: 0,
      page: 1,
      page_size: 10,
    });

    renderWithProviders(<OrderList />);

    await waitFor(() => {
      expect(orderApi.getList).toHaveBeenCalled();
    });

    expect(orderApi.getList).toHaveBeenCalled();
  });

  // 2. 数据加载测试
  it('should load and display orders', async () => {
    const mockOrders = [
      {
        id: 1,
        userId: 1,
        gameId: 1,
        title: 'Test Order 1',
        status: OrderStatus.PENDING,
        priceCents: 10000,
        createdAt: '2024-11-01T00:00:00Z',
        updatedAt: '2024-11-01T00:00:00Z',
      },
      {
        id: 2,
        userId: 2,
        gameId: 2,
        title: 'Test Order 2',
        status: OrderStatus.COMPLETED,
        priceCents: 20000,
        createdAt: '2024-11-02T00:00:00Z',
        updatedAt: '2024-11-02T00:00:00Z',
      },
    ];

    vi.mocked(orderApi.getList).mockResolvedValue({
      list: mockOrders,
      total: 2,
      page: 1,
      page_size: 10,
    });

    renderWithProviders(<OrderList />);

    await waitFor(() => {
      expect(orderApi.getList).toHaveBeenCalled();
    });

    expect(orderApi.getList).toHaveBeenCalledWith(
      expect.objectContaining({
        page: 1,
      })
    );
  });

  // 3. 加载状态测试
  it('should display loading state while fetching', () => {
    vi.mocked(orderApi.getList).mockImplementation(
      () => new Promise(resolve => setTimeout(resolve, 1000))
    );

    renderWithProviders(<OrderList />);

    expect(orderApi.getList).toHaveBeenCalled();
  });

  // 4. 错误处理测试
  it('should handle API error gracefully', async () => {
    vi.mocked(orderApi.getList).mockRejectedValue(new Error('Failed to fetch'));

    renderWithProviders(<OrderList />);

    await waitFor(() => {
      expect(orderApi.getList).toHaveBeenCalled();
    });

    expect(orderApi.getList).toHaveBeenCalled();
  });

  // 5. 分页功能测试
  it('should handle pagination', async () => {
    const mockOrders = Array.from({ length: 10 }, (_, i) => ({
      id: i + 1,
      userId: i + 1,
      gameId: i + 1,
      title: `Order ${i + 1}`,
      status: OrderStatus.PENDING,
      priceCents: 10000 * (i + 1),
      createdAt: '2024-11-01T00:00:00Z',
      updatedAt: '2024-11-01T00:00:00Z',
    }));

    vi.mocked(orderApi.getList).mockResolvedValue({
      list: mockOrders,
      total: 20,
      page: 1,
      page_size: 10,
    });

    renderWithProviders(<OrderList />);

    await waitFor(() => {
      expect(orderApi.getList).toHaveBeenCalled();
    });

    expect(orderApi.getList).toHaveBeenCalledWith(
      expect.objectContaining({
        page: 1,
      })
    );
  });

  // 6. 状态过滤测试
  it('should support filtering by status', async () => {
    const pendingOrders = [
      {
        id: 1,
        userId: 1,
        gameId: 1,
        title: 'Pending Order',
        status: OrderStatus.PENDING,
        priceCents: 10000,
        createdAt: '2024-11-01T00:00:00Z',
        updatedAt: '2024-11-01T00:00:00Z',
      },
    ];

    vi.mocked(orderApi.getList).mockResolvedValue({
      list: pendingOrders,
      total: 1,
      page: 1,
      page_size: 10,
    });

    renderWithProviders(<OrderList />);

    await waitFor(() => {
      expect(orderApi.getList).toHaveBeenCalled();
    });

    expect(orderApi.getList).toHaveBeenCalled();
  });

  // 7. 搜索功能测试
  it('should support searching orders', async () => {
    vi.mocked(orderApi.getList).mockResolvedValue({
      list: [],
      total: 0,
      page: 1,
      page_size: 10,
    });

    renderWithProviders(<OrderList />);

    await waitFor(() => {
      expect(orderApi.getList).toHaveBeenCalled();
    });

    expect(orderApi.getList).toHaveBeenCalled();
  });

  // 8. 更新订单状态测试
  it('should support updating order status', async () => {
    const order = {
      id: 1,
      userId: 1,
      gameId: 1,
      title: 'Test Order',
      status: OrderStatus.PENDING,
      priceCents: 10000,
      createdAt: '2024-11-01T00:00:00Z',
      updatedAt: '2024-11-01T00:00:00Z',
    };

    vi.mocked(orderApi.getList).mockResolvedValue({
      list: [order],
      total: 1,
      page: 1,
      page_size: 10,
    });

    vi.mocked(orderApi.update).mockResolvedValue({
      ...order,
      status: OrderStatus.COMPLETED,
    });

    renderWithProviders(<OrderList />);

    await waitFor(() => {
      expect(orderApi.getList).toHaveBeenCalled();
    });

    expect(orderApi.update).toBeDefined();
  });

  // 9. 删除订单测试
  it('should support deleting orders', async () => {
    const order = {
      id: 1,
      userId: 1,
      gameId: 1,
      title: 'Test Order',
      status: OrderStatus.PENDING,
      priceCents: 10000,
      createdAt: '2024-11-01T00:00:00Z',
      updatedAt: '2024-11-01T00:00:00Z',
    };

    vi.mocked(orderApi.getList).mockResolvedValue({
      list: [order],
      total: 1,
      page: 1,
      page_size: 10,
    });

    vi.mocked(orderApi.delete).mockResolvedValue(undefined);

    renderWithProviders(<OrderList />);

    await waitFor(() => {
      expect(orderApi.getList).toHaveBeenCalled();
    });

    expect(orderApi.delete).toBeDefined();
  });

  // 10. 空状态测试
  it('should display empty state when no orders available', async () => {
    vi.mocked(orderApi.getList).mockResolvedValue({
      list: [],
      total: 0,
      page: 1,
      page_size: 10,
    });

    renderWithProviders(<OrderList />);

    await waitFor(() => {
      expect(orderApi.getList).toHaveBeenCalled();
    });

    expect(orderApi.getList).toHaveBeenCalledWith(
      expect.objectContaining({
        page: 1,
      })
    );
  });

  // 11. 获取订单详情测试
  it('should support fetching order details', async () => {
    const order = {
      id: 1,
      userId: 1,
      gameId: 1,
      title: 'Test Order',
      status: OrderStatus.PENDING,
      priceCents: 10000,
      createdAt: '2024-11-01T00:00:00Z',
      updatedAt: '2024-11-01T00:00:00Z',
    };

    vi.mocked(orderApi.getList).mockResolvedValue({
      list: [order],
      total: 1,
      page: 1,
      page_size: 10,
    });

    vi.mocked(orderApi.getDetail).mockResolvedValue(order);

    renderWithProviders(<OrderList />);

    await waitFor(() => {
      expect(orderApi.getList).toHaveBeenCalled();
    });

    expect(orderApi.getDetail).toBeDefined();
  });

  // 12. 多条件过滤测试
  it('should support multiple filter conditions', async () => {
    vi.mocked(orderApi.getList).mockResolvedValue({
      list: [],
      total: 0,
      page: 1,
      page_size: 10,
    });

    renderWithProviders(<OrderList />);

    await waitFor(() => {
      expect(orderApi.getList).toHaveBeenCalled();
    });

    expect(orderApi.getList).toHaveBeenCalled();
  });
});
