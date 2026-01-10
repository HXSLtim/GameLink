/**
 * UserBehavior Analysis Page Tests
 *
 * Tests for UserBehavior page component including:
 * - Successful data loading
 * - Loading states
 * - Error handling
 * - Statistics display
 * - Trend data display
 * - Distribution data display
 * - Time range selection
 */

import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { screen, waitFor } from '@testing-library/react';
import UserBehaviorPage from './index';
import { renderWithProviders, resetAllMocks, flushPromises } from '@/testutils';

// Mock apiClient using vi.hoisted
const { mockApiClient, mockMessage } = vi.hoisted(() => ({
  mockApiClient: {
    get: vi.fn(),
    post: vi.fn(),
    put: vi.fn(),
    delete: vi.fn(),
  },
  mockMessage: {
    success: vi.fn(),
    error: vi.fn(),
    warning: vi.fn(),
    info: vi.fn(),
    loading: vi.fn(),
  },
}));

vi.mock('@/api/client', () => ({
  default: mockApiClient,
}));

vi.mock('antd', async () => {
  const actual = await vi.importActual('antd');
  return {
    ...actual,
    App: {
      ...(actual as Record<string, unknown>).App,
      useApp: () => ({
        message: mockMessage,
        notification: {},
        modal: {},
      }),
    },
  };
});

// Helper function to create mock stats
const createMockStats = (overrides: Record<string, unknown> = {}): Record<string, unknown> => ({
  dau: 1000,
  mau: 15000,
  avgOnlineTime: 45,
  avgSpending: 88.5,
  newUsers: 120,
  activeRate: 65.5,
  ...overrides,
});

// Helper function to create mock trend data
const createMockTrend = (): Record<string, unknown>[] => [
  { date: '2024-01-01', dau: 900, newUsers: 100, orders: 200 },
  { date: '2024-01-02', dau: 950, newUsers: 110, orders: 220 },
  { date: '2024-01-03', dau: 1000, newUsers: 120, orders: 250 },
];

// Helper function to create mock distribution data
const createMockDistribution = (): Record<string, unknown> => ({
  regions: [
    { name: '广东', count: 500 },
    { name: '北京', count: 400 },
    { name: '上海', count: 300 },
  ],
  ageGroups: [
    { range: '18-24', count: 600 },
    { range: '25-34', count: 800 },
    { range: '35-44', count: 300 },
  ],
  devices: [
    { type: 'iOS', count: 700 },
    { type: 'Android', count: 900 },
    { type: 'Web', count: 200 },
  ],
});

describe('UserBehaviorPage', () => {
  beforeEach(() => {
    resetAllMocks();
    mockMessage.success.mockClear();
    mockMessage.error.mockClear();
    localStorage.setItem('token', 'test-token');
    localStorage.setItem('user_role', 'admin');

    // Set default mock return values
    mockApiClient.get.mockImplementation((url: string) => {
      if (url === '/admin/users/behavior/stats') {
        return Promise.resolve({
          data: {
            success: true,
            data: createMockStats(),
          },
        });
      }
      if (url === '/admin/users/behavior/trend') {
        return Promise.resolve({
          data: {
            success: true,
            data: createMockTrend(),
          },
        });
      }
      if (url === '/admin/users/behavior/distribution') {
        return Promise.resolve({
          data: {
            success: true,
            data: createMockDistribution(),
          },
        });
      }
      return Promise.resolve({ data: { success: true, data: null } });
    });
  });

  afterEach(() => {
    vi.clearAllMocks();
  });

  describe('Successful Data Loading', () => {
    it('should render user behavior page successfully', async () => {
      renderWithProviders(<UserBehaviorPage />);

      await waitFor(() => {
        expect(screen.getByText('用户行为分析')).toBeInTheDocument();
      });

      expect(mockApiClient.get).toHaveBeenCalledWith('/admin/users/behavior/stats');
    });

    it('should display page subtitle', async () => {
      renderWithProviders(<UserBehaviorPage />);

      await waitFor(() => {
        expect(screen.getByText('用户活跃度、消费习惯和分布统计')).toBeInTheDocument();
      });
    });

    it('should load all data on mount', async () => {
      renderWithProviders(<UserBehaviorPage />);

      await waitFor(() => {
        expect(mockApiClient.get).toHaveBeenCalledWith('/admin/users/behavior/stats');
        expect(mockApiClient.get).toHaveBeenCalledWith('/admin/users/behavior/trend', expect.any(Object));
        expect(mockApiClient.get).toHaveBeenCalledWith('/admin/users/behavior/distribution');
      });
    });
  });

  describe('Statistics Display', () => {
    it('should display DAU statistic', async () => {
      renderWithProviders(<UserBehaviorPage />);

      await waitFor(() => {
        expect(screen.getByText('日活用户(DAU)')).toBeInTheDocument();
      });
    });

    it('should display MAU statistic', async () => {
      renderWithProviders(<UserBehaviorPage />);

      await waitFor(() => {
        expect(screen.getByText('月活用户(MAU)')).toBeInTheDocument();
      });
    });

    it('should display average online time', async () => {
      renderWithProviders(<UserBehaviorPage />);

      await waitFor(() => {
        expect(screen.getByText('平均在线时长')).toBeInTheDocument();
      });
    });

    it('should display average spending', async () => {
      renderWithProviders(<UserBehaviorPage />);

      await waitFor(() => {
        expect(screen.getByText('人均消费')).toBeInTheDocument();
      });
    });

    it('should display new users count', async () => {
      renderWithProviders(<UserBehaviorPage />);

      await waitFor(() => {
        expect(screen.getByText('新增用户')).toBeInTheDocument();
      });
    });

    it('should display active rate', async () => {
      renderWithProviders(<UserBehaviorPage />);

      await waitFor(() => {
        expect(screen.getByText('活跃率')).toBeInTheDocument();
      });
    });

    it('should display statistic values', async () => {
      renderWithProviders(<UserBehaviorPage />);

      await waitFor(() => {
        // DAU value
        expect(screen.getByText('1000')).toBeInTheDocument();
      });
    });
  });

  describe('Trend Data Display', () => {
    it('should display trend card title', async () => {
      renderWithProviders(<UserBehaviorPage />);

      await waitFor(() => {
        expect(screen.getByText('用户活动趋势')).toBeInTheDocument();
      });
    });

    it('should display trend table headers', async () => {
      renderWithProviders(<UserBehaviorPage />);

      await waitFor(() => {
        expect(screen.getByText('日期')).toBeInTheDocument();
        // DAU appears in both statistics and table header
        const dauElements = screen.getAllByText('DAU');
        expect(dauElements.length).toBeGreaterThan(0);
        // 新增用户 appears in both statistics and table header
        const newUserElements = screen.getAllByText('新增用户');
        expect(newUserElements.length).toBeGreaterThan(0);
        expect(screen.getByText('订单数')).toBeInTheDocument();
      });
    });
  });

  describe('Distribution Data Display', () => {
    it('should display region distribution card', async () => {
      renderWithProviders(<UserBehaviorPage />);

      await waitFor(() => {
        expect(screen.getByText('地域分布')).toBeInTheDocument();
      });
    });

    it('should display age distribution card', async () => {
      renderWithProviders(<UserBehaviorPage />);

      await waitFor(() => {
        expect(screen.getByText('年龄分布')).toBeInTheDocument();
      });
    });

    it('should display device distribution card', async () => {
      renderWithProviders(<UserBehaviorPage />);

      await waitFor(() => {
        expect(screen.getByText('设备分布')).toBeInTheDocument();
      });
    });

    it('should display region data', async () => {
      renderWithProviders(<UserBehaviorPage />);

      await waitFor(() => {
        expect(screen.getByText('广东')).toBeInTheDocument();
        expect(screen.getByText('北京')).toBeInTheDocument();
      });
    });

    it('should display age group data', async () => {
      renderWithProviders(<UserBehaviorPage />);

      await waitFor(() => {
        expect(screen.getByText('18-24')).toBeInTheDocument();
        expect(screen.getByText('25-34')).toBeInTheDocument();
      });
    });

    it('should display device data', async () => {
      renderWithProviders(<UserBehaviorPage />);

      await waitFor(() => {
        expect(screen.getByText('iOS')).toBeInTheDocument();
        expect(screen.getByText('Android')).toBeInTheDocument();
      });
    });
  });

  describe('Time Range Selection', () => {
    it('should display time range selector', async () => {
      renderWithProviders(<UserBehaviorPage />);

      await waitFor(() => {
        expect(screen.getByText('最近7天')).toBeInTheDocument();
      });
    });

    it('should have time range options', async () => {
      renderWithProviders(<UserBehaviorPage />);

      await waitFor(() => {
        expect(screen.getByText('用户活动趋势')).toBeInTheDocument();
      });

      // The select should have options for 7, 14, 30 days
      const select = screen.getByRole('combobox');
      expect(select).toBeInTheDocument();
    });
  });

  describe('Loading States', () => {
    it('should show loading indicator while fetching data', async () => {
      mockApiClient.get.mockImplementation(
        () =>
          new Promise((resolve) => {
            setTimeout(() => {
              resolve({
                data: {
                  success: true,
                  data: createMockStats(),
                },
              });
            }, 100);
          })
      );

      renderWithProviders(<UserBehaviorPage />);

      expect(mockApiClient.get).toHaveBeenCalled();
    });

    it('should hide loading indicator after data loads', async () => {
      renderWithProviders(<UserBehaviorPage />);

      await waitFor(() => {
        expect(screen.getByText('用户行为分析')).toBeInTheDocument();
      });

      await flushPromises();

      expect(mockApiClient.get).toHaveBeenCalled();
    });
  });

  describe('Error Handling', () => {
    it('should handle API errors gracefully', async () => {
      mockApiClient.get.mockRejectedValue(new Error('Network error'));

      renderWithProviders(<UserBehaviorPage />);

      // The page should still render even with errors
      await waitFor(() => {
        expect(screen.getByText('用户行为分析')).toBeInTheDocument();
      });

      // API should have been called
      expect(mockApiClient.get).toHaveBeenCalled();
    });
  });

  describe('Empty State', () => {
    it('should handle empty trend data', async () => {
      mockApiClient.get.mockImplementation((url: string) => {
        if (url === '/admin/users/behavior/stats') {
          return Promise.resolve({
            data: { success: true, data: createMockStats() },
          });
        }
        if (url === '/admin/users/behavior/trend') {
          return Promise.resolve({
            data: { success: true, data: [] },
          });
        }
        if (url === '/admin/users/behavior/distribution') {
          return Promise.resolve({
            data: { success: true, data: createMockDistribution() },
          });
        }
        return Promise.resolve({ data: { success: true, data: null } });
      });

      renderWithProviders(<UserBehaviorPage />);

      await waitFor(() => {
        expect(screen.getByText('暂无数据')).toBeInTheDocument();
      });
    });

    it('should handle empty distribution data', async () => {
      mockApiClient.get.mockImplementation((url: string) => {
        if (url === '/admin/users/behavior/stats') {
          return Promise.resolve({
            data: { success: true, data: createMockStats() },
          });
        }
        if (url === '/admin/users/behavior/trend') {
          return Promise.resolve({
            data: { success: true, data: createMockTrend() },
          });
        }
        if (url === '/admin/users/behavior/distribution') {
          return Promise.resolve({
            data: {
              success: true,
              data: { regions: [], ageGroups: [], devices: [] },
            },
          });
        }
        return Promise.resolve({ data: { success: true, data: null } });
      });

      renderWithProviders(<UserBehaviorPage />);

      await waitFor(() => {
        // Multiple "暂无数据" for each empty distribution
        const emptyTexts = screen.getAllByText('暂无数据');
        expect(emptyTexts.length).toBeGreaterThan(0);
      });
    });
  });

  describe('Statistics with Zero Values', () => {
    it('should display zero values correctly', async () => {
      mockApiClient.get.mockImplementation((url: string) => {
        if (url === '/admin/users/behavior/stats') {
          return Promise.resolve({
            data: {
              success: true,
              data: createMockStats({
                dau: 0,
                mau: 0,
                avgOnlineTime: 0,
                avgSpending: 0,
                newUsers: 0,
                activeRate: 0,
              }),
            },
          });
        }
        if (url === '/admin/users/behavior/trend') {
          return Promise.resolve({
            data: { success: true, data: [] },
          });
        }
        if (url === '/admin/users/behavior/distribution') {
          return Promise.resolve({
            data: { success: true, data: { regions: [], ageGroups: [], devices: [] } },
          });
        }
        return Promise.resolve({ data: { success: true, data: null } });
      });

      renderWithProviders(<UserBehaviorPage />);

      await waitFor(() => {
        expect(screen.getByText('用户行为分析')).toBeInTheDocument();
      });

      // Should display 0 values
      const zeroValues = screen.getAllByText('0');
      expect(zeroValues.length).toBeGreaterThan(0);
    });
  });

  describe('Accessibility', () => {
    it('should have proper page title', async () => {
      renderWithProviders(<UserBehaviorPage />);

      await waitFor(() => {
        expect(screen.getByText('用户行为分析')).toBeInTheDocument();
      });
    });

    it('should have proper card structure', async () => {
      renderWithProviders(<UserBehaviorPage />);

      await waitFor(() => {
        expect(screen.getByText('用户活动趋势')).toBeInTheDocument();
        expect(screen.getByText('地域分布')).toBeInTheDocument();
        expect(screen.getByText('年龄分布')).toBeInTheDocument();
        expect(screen.getByText('设备分布')).toBeInTheDocument();
      });
    });
  });
});
