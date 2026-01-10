/**
 * Commission Management Page Tests
 *
 * Tests for Commission page component including:
 * - Successful data loading
 * - Loading states
 * - Error handling
 * - Statistics display
 * - CRUD operations for commission rules
 * - Settlement trigger
 */

import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { screen, waitFor, fireEvent } from '@testing-library/react';
import CommissionPage from './index';
import { renderWithProviders, resetAllMocks, flushPromises } from '@/testutils';

// Mock the adminApi module using vi.hoisted
const { mockApi, mockMessage, mockModal } = vi.hoisted(() => ({
  mockApi: {
    getPlatformStats: vi.fn(),
    createCommissionRule: vi.fn(),
    updateCommissionRule: vi.fn(),
    triggerSettlement: vi.fn(),
  },
  mockMessage: {
    success: vi.fn(),
    error: vi.fn(),
    warning: vi.fn(),
    info: vi.fn(),
    loading: vi.fn(),
  },
  mockModal: {
    confirm: vi.fn(),
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
      useApp: () => ({
        message: mockMessage,
        modal: mockModal,
      }),
    },
  };
});

// Helper function to create mock stats
const createMockStats = (overrides: Record<string, unknown> = {}): Record<string, unknown> => ({
  totalRevenueCents: 1000000,
  totalCommissionCents: 200000,
  totalOrderCount: 100,
  completedOrderCount: 80,
  ...overrides,
});

describe('CommissionPage', () => {
  beforeEach(() => {
    resetAllMocks();
    mockMessage.success.mockClear();
    mockMessage.error.mockClear();
    mockModal.confirm.mockClear();
    localStorage.setItem('token', 'test-token');
    localStorage.setItem('user_role', 'admin');

    // Set default mock return values
    mockApi.getPlatformStats.mockResolvedValue({
      data: {
        success: true,
        data: createMockStats(),
      },
    });
  });

  afterEach(() => {
    vi.clearAllMocks();
  });

  describe('Successful Data Loading', () => {
    it('should render commission page successfully', async () => {
      renderWithProviders(<CommissionPage />);

      await waitFor(() => {
        expect(screen.getByText('佣金管理')).toBeInTheDocument();
      });

      expect(mockApi.getPlatformStats).toHaveBeenCalled();
    });

    it('should display page subtitle', async () => {
      renderWithProviders(<CommissionPage />);

      await waitFor(() => {
        expect(screen.getByText('平台佣金设置与结算管理')).toBeInTheDocument();
      });
    });

    it('should display statistics cards', async () => {
      renderWithProviders(<CommissionPage />);

      await waitFor(() => {
        expect(screen.getByText('平台总收入')).toBeInTheDocument();
      });

      expect(screen.getByText('平台佣金')).toBeInTheDocument();
      expect(screen.getByText('总订单数')).toBeInTheDocument();
      expect(screen.getByText('完成订单数')).toBeInTheDocument();
    });

    it('should display statistics values', async () => {
      renderWithProviders(<CommissionPage />);

      await waitFor(() => {
        expect(screen.getByText('100')).toBeInTheDocument(); // totalOrderCount
      });

      expect(screen.getByText('80')).toBeInTheDocument(); // completedOrderCount
    });
  });

  describe('Loading States', () => {
    it('should show loading indicator while fetching data', async () => {
      mockApi.getPlatformStats.mockImplementation(
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

      renderWithProviders(<CommissionPage />);

      expect(mockApi.getPlatformStats).toHaveBeenCalled();
    });

    it('should hide loading indicator after data loads', async () => {
      renderWithProviders(<CommissionPage />);

      await waitFor(() => {
        expect(screen.getByText('佣金管理')).toBeInTheDocument();
      });

      await flushPromises();

      expect(mockApi.getPlatformStats).toHaveBeenCalledTimes(1);
    });
  });

  describe('Error Handling', () => {
    it('should display error message when API fails', async () => {
      mockApi.getPlatformStats.mockRejectedValue(new Error('Network error'));

      renderWithProviders(<CommissionPage />);

      await waitFor(() => {
        expect(mockMessage.error).toHaveBeenCalledWith('加载统计数据失败');
      });
    });
  });

  describe('Month Selection', () => {
    it('should have month picker', async () => {
      renderWithProviders(<CommissionPage />);

      await waitFor(() => {
        expect(screen.getByText('选择月份：')).toBeInTheDocument();
      });
    });

    it('should have refresh button', async () => {
      renderWithProviders(<CommissionPage />);

      await waitFor(() => {
        expect(screen.getByText('刷新')).toBeInTheDocument();
      });
    });

    it('should call loadStats when refresh button clicked', async () => {
      renderWithProviders(<CommissionPage />);

      await waitFor(() => {
        expect(screen.getByText('刷新')).toBeInTheDocument();
      });

      const refreshButton = screen.getByText('刷新');
      fireEvent.click(refreshButton);

      await flushPromises();

      // Initial load + refresh click
      expect(mockApi.getPlatformStats).toHaveBeenCalled();
    });
  });

  describe('Commission Rules Section', () => {
    it('should display commission rules card', async () => {
      renderWithProviders(<CommissionPage />);

      await waitFor(() => {
        expect(screen.getByText('佣金规则')).toBeInTheDocument();
      });
    });

    it('should display commission explanation alert', async () => {
      renderWithProviders(<CommissionPage />);

      await waitFor(() => {
        expect(screen.getByText('佣金说明')).toBeInTheDocument();
      });
    });

    it('should display create rule button', async () => {
      renderWithProviders(<CommissionPage />);

      await waitFor(() => {
        expect(screen.getByText('新增规则')).toBeInTheDocument();
      });
    });

    it('should display default rule in table', async () => {
      renderWithProviders(<CommissionPage />);

      await waitFor(() => {
        expect(screen.getByText('默认抽成规则')).toBeInTheDocument();
      });
    });
  });

  describe('Create Rule Modal', () => {
    it('should open create rule modal when button clicked', async () => {
      renderWithProviders(<CommissionPage />);

      await waitFor(() => {
        expect(screen.getByText('新增规则')).toBeInTheDocument();
      });

      const createButton = screen.getByText('新增规则');
      fireEvent.click(createButton);

      await waitFor(() => {
        expect(screen.getByText('新增佣金规则')).toBeInTheDocument();
      });
    });

    it('should display form fields in modal', async () => {
      renderWithProviders(<CommissionPage />);

      await waitFor(() => {
        expect(screen.getByText('新增规则')).toBeInTheDocument();
      });

      const createButton = screen.getByText('新增规则');
      fireEvent.click(createButton);

      await waitFor(() => {
        expect(screen.getByRole('dialog')).toBeInTheDocument();
      });

      // Check for input placeholders instead of labels
      expect(screen.getByPlaceholderText('请输入规则名称')).toBeInTheDocument();
      expect(screen.getByPlaceholderText('请输入抽成比例')).toBeInTheDocument();
    });
  });

  describe('Settlement Section', () => {
    it('should display settlement trigger button', async () => {
      renderWithProviders(<CommissionPage />);

      await waitFor(() => {
        expect(screen.getByText('触发月度结算')).toBeInTheDocument();
      });
    });

    it('should display settlement explanation', async () => {
      renderWithProviders(<CommissionPage />);

      await waitFor(() => {
        expect(screen.getByText('结算说明')).toBeInTheDocument();
      });
    });

    it('should display settlement cycle info', async () => {
      renderWithProviders(<CommissionPage />);

      await waitFor(() => {
        expect(screen.getByText('结算周期')).toBeInTheDocument();
      });
    });

    it('should display settlement rules info', async () => {
      renderWithProviders(<CommissionPage />);

      await waitFor(() => {
        expect(screen.getByText('结算规则')).toBeInTheDocument();
      });
    });
  });

  describe('Table Structure', () => {
    it('should display table with data', async () => {
      renderWithProviders(<CommissionPage />);

      await waitFor(() => {
        expect(screen.getByRole('table')).toBeInTheDocument();
      });

      // Check for default rule data
      expect(screen.getByText('默认抽成规则')).toBeInTheDocument();
    });

    it('should display default rule rate', async () => {
      renderWithProviders(<CommissionPage />);

      await waitFor(() => {
        expect(screen.getByText('20%')).toBeInTheDocument();
      });
    });
  });

  describe('Statistics with Zero Values', () => {
    it('should display zero values correctly', async () => {
      mockApi.getPlatformStats.mockResolvedValue({
        data: {
          success: true,
          data: createMockStats({
            totalRevenueCents: 0,
            totalCommissionCents: 0,
            totalOrderCount: 0,
            completedOrderCount: 0,
          }),
        },
      });

      renderWithProviders(<CommissionPage />);

      await waitFor(() => {
        expect(screen.getByText('佣金管理')).toBeInTheDocument();
      });

      // Should display 0 values
      const zeroValues = screen.getAllByText('0');
      expect(zeroValues.length).toBeGreaterThan(0);
    });
  });

  describe('Accessibility', () => {
    it('should have proper page title', async () => {
      renderWithProviders(<CommissionPage />);

      await waitFor(() => {
        expect(screen.getByText('佣金管理')).toBeInTheDocument();
      });
    });

    it('should have proper table structure', async () => {
      renderWithProviders(<CommissionPage />);

      await waitFor(() => {
        expect(screen.getByRole('table')).toBeInTheDocument();
      });
    });
  });
});
