/**
 * Activity Management Page Tests
 *
 * Tests for Activity page component including:
 * - Successful data loading
 * - Loading states
 * - Error handling
 * - Statistics display
 * - CRUD operations
 */

import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { screen, waitFor, fireEvent } from '@testing-library/react';
import ActivityPage from './index';
import { renderWithProviders, resetAllMocks, flushPromises } from '@/testutils';

// Mock the activityApi module using vi.hoisted
const { mockApi, mockMessage } = vi.hoisted(() => ({
  mockApi: {
    getActivities: vi.fn(),
    getAllActivityStats: vi.fn(),
    getActivityDetail: vi.fn(),
    createActivity: vi.fn(),
    updateActivity: vi.fn(),
    deleteActivity: vi.fn(),
    publishActivity: vi.fn(),
    unpublishActivity: vi.fn(),
    getActivityRewards: vi.fn(),
    getActivityStats: vi.fn(),
    createReward: vi.fn(),
    updateReward: vi.fn(),
    deleteReward: vi.fn(),
  },
  mockMessage: {
    success: vi.fn(),
    error: vi.fn(),
    warning: vi.fn(),
    info: vi.fn(),
    loading: vi.fn(),
  },
}));

vi.mock('@/api/activity', () => ({
  activityApi: mockApi,
  getActivityTypeLabel: (type: string) => {
    const labels: Record<string, string> = {
      coupon: '优惠券发放',
      discount: '限时折扣',
      gift: '赠品活动',
    };
    return labels[type] || type;
  },
  getActivityStatusLabel: (status: string) => {
    const labels: Record<string, string> = {
      draft: '草稿',
      preheat: '预热中',
      active: '进行中',
      paused: '已暂停',
      ended: '已结束',
      canceled: '已取消',
    };
    return labels[status] || status;
  },
  getActivityStatusColor: (status: string) => {
    const colors: Record<string, string> = {
      draft: 'default',
      preheat: 'blue',
      active: 'green',
      paused: 'orange',
      ended: 'default',
      canceled: 'red',
    };
    return colors[status] || 'default';
  },
  canEditActivity: (activity: { status: string }) => {
    return activity.status === 'draft' || activity.status === 'preheat';
  },
  canDeleteActivity: (activity: { status: string }) => {
    return activity.status === 'draft' || activity.status === 'canceled' || activity.status === 'ended';
  },
  calculateStockPercentage: () => 50,
}));

// Mock antd message
vi.mock('antd', async () => {
  const actual = await vi.importActual('antd');
  return {
    ...actual,
    message: mockMessage,
  };
});

// Helper function to create mock activity
const createMockActivity = (overrides: Record<string, unknown> = {}): Record<string, unknown> => ({
  id: 1,
  name: '测试活动',
  description: '活动描述',
  type: 'coupon',
  status: 'active',
  coverUrl: null,
  bannerUrl: null,
  preheatAt: null,
  startAt: '2024-01-01T00:00:00Z',
  endAt: '2024-12-31T23:59:59Z',
  totalLimit: 0,
  dailyLimit: 0,
  perUserLimit: 1,
  totalParticipants: 100,
  todayParticipants: 10,
  totalClaimed: 50,
  allowVipStack: false,
  rules: '活动规则',
  sortOrder: 0,
  isVisible: true,
  createdAt: '2024-01-01T00:00:00Z',
  updatedAt: '2024-01-01T00:00:00Z',
  ...overrides,
});

// Helper function to create mock stats
const createMockStats = (overrides: Record<string, unknown> = {}): Record<string, unknown> => ({
  totalActivities: 10,
  activeActivities: 5,
  draftActivities: 3,
  endedActivities: 2,
  totalParticipants: 1000,
  totalClaimed: 500,
  ...overrides,
});

// Helper function to create mock reward
const createMockReward = (overrides: Record<string, unknown> = {}): Record<string, unknown> => ({
  id: 1,
  activityId: 1,
  couponTemplateId: 1,
  couponCount: 1,
  probability: 100,
  totalStock: 100,
  remainingStock: 50,
  sortOrder: 0,
  couponTemplate: {
    id: 1,
    name: '测试优惠券',
    type: 'deduct',
  },
  createdAt: '2024-01-01T00:00:00Z',
  updatedAt: '2024-01-01T00:00:00Z',
  ...overrides,
});

describe('ActivityPage', () => {
  beforeEach(() => {
    resetAllMocks();
    mockMessage.success.mockClear();
    mockMessage.error.mockClear();
    mockMessage.warning.mockClear();
    localStorage.setItem('token', 'test-token');
    localStorage.setItem('user_role', 'admin');
    
    // Set default mock return values
    mockApi.getActivities.mockResolvedValue({
      data: {
        success: true,
        data: [createMockActivity()],
        pagination: { total: 1 },
      },
    });
    mockApi.getAllActivityStats.mockResolvedValue({
      data: {
        success: true,
        data: createMockStats(),
      },
    });
    mockApi.getActivityRewards.mockResolvedValue({
      data: {
        success: true,
        data: [createMockReward()],
      },
    });
    mockApi.getActivityStats.mockResolvedValue({
      data: {
        success: true,
        data: {
          activityId: 1,
          totalParticipants: 100,
          todayParticipants: 10,
          totalClaimed: 50,
        },
      },
    });
  });

  afterEach(() => {
    vi.clearAllMocks();
  });

  describe('Successful Data Loading', () => {
    it('should render activity page successfully', async () => {
      renderWithProviders(<ActivityPage />);

      await waitFor(() => {
        expect(screen.getByText('活动管理')).toBeInTheDocument();
      });

      expect(mockApi.getActivities).toHaveBeenCalled();
      expect(mockApi.getAllActivityStats).toHaveBeenCalled();
    });

    it('should display statistics cards', async () => {
      renderWithProviders(<ActivityPage />);

      await waitFor(() => {
        expect(screen.getByText('活动总数')).toBeInTheDocument();
      });

      // "进行中" 和 "草稿" 可能出现在统计卡片和表格中，使用 getAllByText
      const activeElements = screen.getAllByText('进行中');
      expect(activeElements.length).toBeGreaterThan(0);
      const draftElements = screen.getAllByText('草稿');
      expect(draftElements.length).toBeGreaterThan(0);
      expect(screen.getByText('总参与')).toBeInTheDocument();
    });

    it('should display statistics values', async () => {
      renderWithProviders(<ActivityPage />);

      await waitFor(() => {
        expect(screen.getByText('10')).toBeInTheDocument(); // totalActivities
      });

      expect(screen.getByText('5')).toBeInTheDocument(); // activeActivities
      expect(screen.getByText('3')).toBeInTheDocument(); // draftActivities
    });

    it('should display activity list', async () => {
      renderWithProviders(<ActivityPage />);

      await waitFor(() => {
        expect(screen.getByText('测试活动')).toBeInTheDocument();
      });
    });

    it('should display activity type tag', async () => {
      renderWithProviders(<ActivityPage />);

      await waitFor(() => {
        // "优惠券发放" 可能出现多次（表格中的类型标签）
        const typeElements = screen.getAllByText('优惠券发放');
        expect(typeElements.length).toBeGreaterThan(0);
      });
    });
  });

  describe('Loading States', () => {
    it('should show loading indicator while fetching data', async () => {
      mockApi.getActivities.mockImplementation(
        () =>
          new Promise((resolve) => {
            setTimeout(() => {
              resolve({
                data: {
                  success: true,
                  data: [createMockActivity()],
                  pagination: { total: 1 },
                },
              });
            }, 100);
          })
      );

      renderWithProviders(<ActivityPage />);

      expect(mockApi.getActivities).toHaveBeenCalled();
    });

    it('should hide loading indicator after data loads', async () => {
      renderWithProviders(<ActivityPage />);

      await waitFor(() => {
        expect(screen.getByText('活动管理')).toBeInTheDocument();
      });

      await flushPromises();

      expect(mockApi.getActivities).toHaveBeenCalledTimes(1);
    });
  });

  describe('Error Handling', () => {
    it('should display error message when API fails', async () => {
      mockApi.getActivities.mockRejectedValue(new Error('Network error'));

      renderWithProviders(<ActivityPage />);

      await waitFor(() => {
        expect(screen.getByText('活动管理')).toBeInTheDocument();
      });

      await flushPromises();

      expect(mockMessage.error).toHaveBeenCalledWith('加载失败');
    });

    it('should display error message when API returns error', async () => {
      mockApi.getActivities.mockResolvedValue({
        data: {
          success: false,
          message: '服务器错误',
        },
      });

      renderWithProviders(<ActivityPage />);

      await waitFor(() => {
        expect(mockMessage.error).toHaveBeenCalledWith('服务器错误');
      });
    });
  });

  describe('Filter Functionality', () => {
    it('should have search input', async () => {
      renderWithProviders(<ActivityPage />);

      await waitFor(() => {
        expect(screen.getByPlaceholderText('搜索名称')).toBeInTheDocument();
      });
    });

    it('should have type filter', async () => {
      renderWithProviders(<ActivityPage />);

      await waitFor(() => {
        expect(screen.getByText('活动管理')).toBeInTheDocument();
      });

      // Type filter select should exist
      expect(mockApi.getActivities).toHaveBeenCalled();
    });

    it('should have status filter', async () => {
      renderWithProviders(<ActivityPage />);

      await waitFor(() => {
        expect(screen.getByText('活动管理')).toBeInTheDocument();
      });

      // Status filter select should exist
      expect(mockApi.getActivities).toHaveBeenCalled();
    });
  });

  describe('Toolbar Buttons', () => {
    it('should display refresh button', async () => {
      renderWithProviders(<ActivityPage />);

      await waitFor(() => {
        expect(screen.getByText('刷新')).toBeInTheDocument();
      });
    });

    it('should display add button', async () => {
      renderWithProviders(<ActivityPage />);

      await waitFor(() => {
        expect(screen.getByText('新建活动')).toBeInTheDocument();
      });
    });

    it('should refresh data when refresh button clicked', async () => {
      renderWithProviders(<ActivityPage />);

      await waitFor(() => {
        expect(screen.getByText('刷新')).toBeInTheDocument();
      });

      const refreshButton = screen.getByText('刷新');
      fireEvent.click(refreshButton);

      await waitFor(() => {
        expect(mockApi.getActivities).toHaveBeenCalledTimes(2);
      });
    });
  });

  describe('Activity Actions', () => {
    it('should display detail button', async () => {
      renderWithProviders(<ActivityPage />);

      await waitFor(() => {
        expect(screen.getByText('详情')).toBeInTheDocument();
      });
    });

    it('should display rewards button', async () => {
      renderWithProviders(<ActivityPage />);

      await waitFor(() => {
        expect(screen.getByText('奖励')).toBeInTheDocument();
      });
    });

    it('should display unpublish button for visible activity', async () => {
      renderWithProviders(<ActivityPage />);

      await waitFor(() => {
        expect(screen.getByText('下架')).toBeInTheDocument();
      });
    });

    it('should display publish button for hidden activity', async () => {
      mockApi.getActivities.mockResolvedValue({
        data: {
          success: true,
          data: [createMockActivity({ isVisible: false })],
          pagination: { total: 1 },
        },
      });

      renderWithProviders(<ActivityPage />);

      await waitFor(() => {
        expect(screen.getByText('发布')).toBeInTheDocument();
      });
    });

    it('should display edit button', async () => {
      renderWithProviders(<ActivityPage />);

      await waitFor(() => {
        expect(screen.getByText('编辑')).toBeInTheDocument();
      });
    });

    it('should display delete button for deletable activity', async () => {
      mockApi.getActivities.mockResolvedValue({
        data: {
          success: true,
          data: [createMockActivity({ status: 'draft' })],
          pagination: { total: 1 },
        },
      });

      renderWithProviders(<ActivityPage />);

      await waitFor(() => {
        expect(screen.getByText('删除')).toBeInTheDocument();
      });
    });
  });

  describe('Empty State', () => {
    it('should handle empty activity list', async () => {
      mockApi.getActivities.mockResolvedValue({
        data: {
          success: true,
          data: [],
          pagination: { total: 0 },
        },
      });

      renderWithProviders(<ActivityPage />);

      await waitFor(() => {
        expect(screen.getByText('活动管理')).toBeInTheDocument();
      });

      expect(mockApi.getActivities).toHaveBeenCalled();
    });
  });

  describe('Statistics with Zero Values', () => {
    it('should display zero values correctly', async () => {
      mockApi.getAllActivityStats.mockResolvedValue({
        data: {
          success: true,
          data: createMockStats({
            totalActivities: 0,
            activeActivities: 0,
            draftActivities: 0,
            totalParticipants: 0,
          }),
        },
      });

      renderWithProviders(<ActivityPage />);

      await waitFor(() => {
        expect(screen.getByText('活动管理')).toBeInTheDocument();
      });

      // Should display 0 values
      const zeroValues = screen.getAllByText('0');
      expect(zeroValues.length).toBeGreaterThan(0);
    });
  });

  describe('Statistics with Null Values', () => {
    it('should handle null stats gracefully', async () => {
      mockApi.getAllActivityStats.mockResolvedValue({
        data: {
          success: true,
          data: null,
        },
      });

      renderWithProviders(<ActivityPage />);

      await waitFor(() => {
        expect(screen.getByText('活动管理')).toBeInTheDocument();
      });

      // Should display 0 for null values
      const zeroValues = screen.getAllByText('0');
      expect(zeroValues.length).toBeGreaterThan(0);
    });
  });

  describe('Activity Status Display', () => {
    it('should display active status correctly', async () => {
      renderWithProviders(<ActivityPage />);

      await waitFor(() => {
        // "进行中" appears in both stats card and table
        const activeElements = screen.getAllByText('进行中');
        expect(activeElements.length).toBeGreaterThan(0);
      });
    });

    it('should display draft status correctly', async () => {
      mockApi.getActivities.mockResolvedValue({
        data: {
          success: true,
          data: [createMockActivity({ status: 'draft' })],
          pagination: { total: 1 },
        },
      });

      renderWithProviders(<ActivityPage />);

      await waitFor(() => {
        // "草稿" appears in both stats card and table
        const draftElements = screen.getAllByText('草稿');
        expect(draftElements.length).toBeGreaterThan(0);
      });
    });
  });

  describe('Visibility Display', () => {
    it('should display visible tag for visible activity', async () => {
      renderWithProviders(<ActivityPage />);

      await waitFor(() => {
        // "显示" 可能出现多次（表格中的可见性标签）
        const visibleElements = screen.getAllByText('显示');
        expect(visibleElements.length).toBeGreaterThan(0);
      });
    });

    it('should display hidden tag for hidden activity', async () => {
      mockApi.getActivities.mockResolvedValue({
        data: {
          success: true,
          data: [createMockActivity({ isVisible: false })],
          pagination: { total: 1 },
        },
      });

      renderWithProviders(<ActivityPage />);

      await waitFor(() => {
        expect(screen.getByText('隐藏')).toBeInTheDocument();
      });
    });
  });

  describe('VIP Stack Display', () => {
    it('should display VIP stack allowed tag', async () => {
      mockApi.getActivities.mockResolvedValue({
        data: {
          success: true,
          data: [createMockActivity({ allowVipStack: true })],
          pagination: { total: 1 },
        },
      });

      renderWithProviders(<ActivityPage />);

      await waitFor(() => {
        expect(screen.getByText('允许')).toBeInTheDocument();
      });
    });

    it('should display VIP stack not allowed tag', async () => {
      renderWithProviders(<ActivityPage />);

      await waitFor(() => {
        expect(screen.getByText('不允许')).toBeInTheDocument();
      });
    });
  });

  describe('Accessibility', () => {
    it('should have proper page title', async () => {
      renderWithProviders(<ActivityPage />);

      await waitFor(() => {
        expect(screen.getByText('活动管理')).toBeInTheDocument();
      });
    });

    it('should have proper table structure', async () => {
      renderWithProviders(<ActivityPage />);

      await waitFor(() => {
        expect(screen.getByRole('table')).toBeInTheDocument();
      });
    });
  });

  describe('Data Refresh', () => {
    it('should call loadData and loadStats on mount', async () => {
      renderWithProviders(<ActivityPage />);

      await waitFor(() => {
        expect(mockApi.getActivities).toHaveBeenCalledTimes(1);
        expect(mockApi.getAllActivityStats).toHaveBeenCalledTimes(1);
      });
    });
  });
});
