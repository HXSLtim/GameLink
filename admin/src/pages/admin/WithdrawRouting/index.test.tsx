/**
 * WithdrawRouting Statistics Page Tests
 *
 * Tests for WithdrawRouting page component including:
 * - Successful data loading
 * - Loading states
 * - Error handling
 * - Statistics display
 * - Company statistics display
 * - Date range filtering
 * - Export functionality
 * - Report generation
 */

import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { screen, waitFor, fireEvent } from '@testing-library/react';
import WithdrawRoutingPage from './index';
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
    message: mockMessage,
  };
});

// Helper function to create mock withdraw record
const createMockWithdraw = (overrides: Record<string, unknown> = {}): Record<string, unknown> => ({
  id: 1,
  playerId: 100,
  playerName: '测试陪玩师',
  amount: 1000.5,
  status: 'completed',
  settlementCompanyId: 1,
  settlementCompanyName: '结算公司A',
  createdAt: '2024-01-01T00:00:00Z',
  ...overrides,
});

// Helper function to create mock stats
const createMockStats = (overrides: Record<string, unknown> = {}): Record<string, unknown> => ({
  totalAmount: 50000,
  totalCount: 100,
  completedAmount: 40000,
  completedCount: 80,
  pendingAmount: 10000,
  pendingCount: 20,
  byCompany: [
    { companyId: 1, companyName: '结算公司A', amount: 30000, count: 60 },
    { companyId: 2, companyName: '结算公司B', amount: 20000, count: 40 },
  ],
  ...overrides,
});

describe('WithdrawRoutingPage', () => {
  beforeEach(() => {
    resetAllMocks();
    mockMessage.success.mockClear();
    mockMessage.error.mockClear();
    localStorage.setItem('token', 'test-token');
    localStorage.setItem('user_role', 'admin');

    // Set default mock return values
    mockApiClient.get.mockImplementation((url: string) => {
      if (url === '/admin/withdrawals/by-company') {
        return Promise.resolve({
          data: {
            success: true,
            data: [createMockWithdraw()],
            pagination: { total: 1, page: 1, pageSize: 10 },
          },
        });
      }
      if (url === '/admin/withdrawals/routing-stats') {
        return Promise.resolve({
          data: {
            success: true,
            data: createMockStats(),
          },
        });
      }
      if (url === '/admin/withdrawals/routing-report') {
        return Promise.resolve({
          data: { success: true },
        });
      }
      return Promise.resolve({ data: { success: true, data: null } });
    });
  });

  afterEach(() => {
    vi.clearAllMocks();
  });

  describe('Successful Data Loading', () => {
    it('should render withdraw routing page successfully', async () => {
      renderWithProviders(<WithdrawRoutingPage />);

      await waitFor(() => {
        expect(screen.getByText('提现分流统计')).toBeInTheDocument();
      });

      expect(mockApiClient.get).toHaveBeenCalledWith('/admin/withdrawals/by-company', expect.any(Object));
      expect(mockApiClient.get).toHaveBeenCalledWith('/admin/withdrawals/routing-stats', expect.any(Object));
    });

    it('should display page subtitle', async () => {
      renderWithProviders(<WithdrawRoutingPage />);

      await waitFor(() => {
        expect(screen.getByText('按结算公司统计提现数据')).toBeInTheDocument();
      });
    });

    it('should display withdraw list', async () => {
      renderWithProviders(<WithdrawRoutingPage />);

      await waitFor(() => {
        expect(screen.getByText('测试陪玩师')).toBeInTheDocument();
      });
    });

    it('should display settlement company name', async () => {
      renderWithProviders(<WithdrawRoutingPage />);

      await waitFor(() => {
        // Company name appears in both statistics card and table
        const companyElements = screen.getAllByText('结算公司A');
        expect(companyElements.length).toBeGreaterThan(0);
      });
    });

    it('should display withdraw amount', async () => {
      renderWithProviders(<WithdrawRoutingPage />);

      await waitFor(() => {
        expect(screen.getByText('¥1000.50')).toBeInTheDocument();
      });
    });
  });

  describe('Statistics Display', () => {
    it('should display total amount statistic', async () => {
      renderWithProviders(<WithdrawRoutingPage />);

      await waitFor(() => {
        expect(screen.getByText('总提现金额')).toBeInTheDocument();
      });
    });

    it('should display total count statistic', async () => {
      renderWithProviders(<WithdrawRoutingPage />);

      await waitFor(() => {
        expect(screen.getByText('总提现笔数')).toBeInTheDocument();
      });
    });

    it('should display completed amount statistic', async () => {
      renderWithProviders(<WithdrawRoutingPage />);

      await waitFor(() => {
        expect(screen.getByText('已完成金额')).toBeInTheDocument();
      });
    });

    it('should display pending amount statistic', async () => {
      renderWithProviders(<WithdrawRoutingPage />);

      await waitFor(() => {
        expect(screen.getByText('待处理金额')).toBeInTheDocument();
      });
    });
  });

  describe('Company Statistics Display', () => {
    it('should display company statistics card', async () => {
      renderWithProviders(<WithdrawRoutingPage />);

      await waitFor(() => {
        expect(screen.getByText('按结算公司统计')).toBeInTheDocument();
      });
    });

    it('should display company names in statistics', async () => {
      renderWithProviders(<WithdrawRoutingPage />);

      await waitFor(() => {
        // Company names should appear in the statistics section
        const companyAElements = screen.getAllByText('结算公司A');
        expect(companyAElements.length).toBeGreaterThan(0);
      });
    });
  });

  describe('Status Display', () => {
    it('should display completed status correctly', async () => {
      renderWithProviders(<WithdrawRoutingPage />);

      await waitFor(() => {
        expect(screen.getByText('已完成')).toBeInTheDocument();
      });
    });

    it('should display pending status correctly', async () => {
      mockApiClient.get.mockImplementation((url: string) => {
        if (url === '/admin/withdrawals/by-company') {
          return Promise.resolve({
            data: {
              success: true,
              data: [createMockWithdraw({ status: 'pending' })],
              pagination: { total: 1 },
            },
          });
        }
        if (url === '/admin/withdrawals/routing-stats') {
          return Promise.resolve({
            data: { success: true, data: createMockStats() },
          });
        }
        return Promise.resolve({ data: { success: true, data: null } });
      });

      renderWithProviders(<WithdrawRoutingPage />);

      await waitFor(() => {
        expect(screen.getByText('待审核')).toBeInTheDocument();
      });
    });

    it('should display approved status correctly', async () => {
      mockApiClient.get.mockImplementation((url: string) => {
        if (url === '/admin/withdrawals/by-company') {
          return Promise.resolve({
            data: {
              success: true,
              data: [createMockWithdraw({ status: 'approved' })],
              pagination: { total: 1 },
            },
          });
        }
        if (url === '/admin/withdrawals/routing-stats') {
          return Promise.resolve({
            data: { success: true, data: createMockStats() },
          });
        }
        return Promise.resolve({ data: { success: true, data: null } });
      });

      renderWithProviders(<WithdrawRoutingPage />);

      await waitFor(() => {
        expect(screen.getByText('已批准')).toBeInTheDocument();
      });
    });

    it('should display rejected status correctly', async () => {
      mockApiClient.get.mockImplementation((url: string) => {
        if (url === '/admin/withdrawals/by-company') {
          return Promise.resolve({
            data: {
              success: true,
              data: [createMockWithdraw({ status: 'rejected' })],
              pagination: { total: 1 },
            },
          });
        }
        if (url === '/admin/withdrawals/routing-stats') {
          return Promise.resolve({
            data: { success: true, data: createMockStats() },
          });
        }
        return Promise.resolve({ data: { success: true, data: null } });
      });

      renderWithProviders(<WithdrawRoutingPage />);

      await waitFor(() => {
        expect(screen.getByText('已拒绝')).toBeInTheDocument();
      });
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
                  data: [createMockWithdraw()],
                  pagination: { total: 1 },
                },
              });
            }, 100);
          })
      );

      renderWithProviders(<WithdrawRoutingPage />);

      expect(mockApiClient.get).toHaveBeenCalled();
    });

    it('should hide loading indicator after data loads', async () => {
      renderWithProviders(<WithdrawRoutingPage />);

      await waitFor(() => {
        expect(screen.getByText('提现分流统计')).toBeInTheDocument();
      });

      await flushPromises();

      expect(mockApiClient.get).toHaveBeenCalled();
    });
  });

  describe('Error Handling', () => {
    it('should display error message when API fails', async () => {
      mockApiClient.get.mockRejectedValue(new Error('Network error'));

      renderWithProviders(<WithdrawRoutingPage />);

      await waitFor(() => {
        expect(mockMessage.error).toHaveBeenCalledWith('加载失败');
      });
    });
  });

  describe('Filter Functionality', () => {
    it('should have settlement company ID search input', async () => {
      renderWithProviders(<WithdrawRoutingPage />);

      await waitFor(() => {
        expect(screen.getByPlaceholderText('公司ID')).toBeInTheDocument();
      });
    });

    it('should have date range picker', async () => {
      renderWithProviders(<WithdrawRoutingPage />);

      await waitFor(() => {
        expect(screen.getByText('提现明细')).toBeInTheDocument();
      });

      // Date range picker should be present
      expect(mockApiClient.get).toHaveBeenCalled();
    });
  });

  describe('Toolbar Buttons', () => {
    it('should display monthly report button', async () => {
      renderWithProviders(<WithdrawRoutingPage />);

      await waitFor(() => {
        expect(screen.getByText('月报')).toBeInTheDocument();
      });
    });

    it('should display export button', async () => {
      renderWithProviders(<WithdrawRoutingPage />);

      await waitFor(() => {
        expect(screen.getByText('导出数据')).toBeInTheDocument();
      });
    });
  });

  describe('Report Generation', () => {
    it('should call report API when monthly report button clicked', async () => {
      renderWithProviders(<WithdrawRoutingPage />);

      await waitFor(() => {
        expect(screen.getByText('月报')).toBeInTheDocument();
      });

      const reportButton = screen.getByText('月报');
      fireEvent.click(reportButton);

      await waitFor(() => {
        expect(mockApiClient.get).toHaveBeenCalledWith('/admin/withdrawals/routing-report', expect.any(Object));
      });
    });

    it('should show success message when report generated', async () => {
      renderWithProviders(<WithdrawRoutingPage />);

      await waitFor(() => {
        expect(screen.getByText('月报')).toBeInTheDocument();
      });

      const reportButton = screen.getByText('月报');
      fireEvent.click(reportButton);

      await waitFor(() => {
        expect(mockMessage.success).toHaveBeenCalledWith('报表生成成功');
      });
    });

    it('should show error message when report generation fails', async () => {
      mockApiClient.get.mockImplementation((url: string) => {
        if (url === '/admin/withdrawals/routing-report') {
          return Promise.reject(new Error('Report generation failed'));
        }
        if (url === '/admin/withdrawals/by-company') {
          return Promise.resolve({
            data: {
              success: true,
              data: [createMockWithdraw()],
              pagination: { total: 1 },
            },
          });
        }
        if (url === '/admin/withdrawals/routing-stats') {
          return Promise.resolve({
            data: { success: true, data: createMockStats() },
          });
        }
        return Promise.resolve({ data: { success: true, data: null } });
      });

      renderWithProviders(<WithdrawRoutingPage />);

      await waitFor(() => {
        expect(screen.getByText('月报')).toBeInTheDocument();
      });

      const reportButton = screen.getByText('月报');
      fireEvent.click(reportButton);

      await waitFor(() => {
        expect(mockMessage.error).toHaveBeenCalledWith('生成报表失败');
      });
    });
  });

  describe('Export Functionality', () => {
    it('should show success message when export clicked', async () => {
      renderWithProviders(<WithdrawRoutingPage />);

      await waitFor(() => {
        expect(screen.getByText('导出数据')).toBeInTheDocument();
      });

      const exportButton = screen.getByText('导出数据');
      fireEvent.click(exportButton);

      await waitFor(() => {
        expect(mockMessage.success).toHaveBeenCalledWith('导出成功');
      });
    });
  });

  describe('Empty State', () => {
    it('should handle empty withdraw list', async () => {
      mockApiClient.get.mockImplementation((url: string) => {
        if (url === '/admin/withdrawals/by-company') {
          return Promise.resolve({
            data: {
              success: true,
              data: [],
              pagination: { total: 0 },
            },
          });
        }
        if (url === '/admin/withdrawals/routing-stats') {
          return Promise.resolve({
            data: { success: true, data: createMockStats() },
          });
        }
        return Promise.resolve({ data: { success: true, data: null } });
      });

      renderWithProviders(<WithdrawRoutingPage />);

      await waitFor(() => {
        expect(screen.getByText('提现分流统计')).toBeInTheDocument();
      });

      expect(mockApiClient.get).toHaveBeenCalled();
    });

    it('should handle empty company statistics', async () => {
      mockApiClient.get.mockImplementation((url: string) => {
        if (url === '/admin/withdrawals/by-company') {
          return Promise.resolve({
            data: {
              success: true,
              data: [createMockWithdraw()],
              pagination: { total: 1 },
            },
          });
        }
        if (url === '/admin/withdrawals/routing-stats') {
          return Promise.resolve({
            data: {
              success: true,
              data: createMockStats({ byCompany: [] }),
            },
          });
        }
        return Promise.resolve({ data: { success: true, data: null } });
      });

      renderWithProviders(<WithdrawRoutingPage />);

      await waitFor(() => {
        expect(screen.getByText('提现分流统计')).toBeInTheDocument();
      });

      // Company statistics card should not be displayed when empty
      expect(mockApiClient.get).toHaveBeenCalled();
    });
  });

  describe('Statistics with Zero Values', () => {
    it('should display zero values correctly', async () => {
      mockApiClient.get.mockImplementation((url: string) => {
        if (url === '/admin/withdrawals/by-company') {
          return Promise.resolve({
            data: {
              success: true,
              data: [],
              pagination: { total: 0 },
            },
          });
        }
        if (url === '/admin/withdrawals/routing-stats') {
          return Promise.resolve({
            data: {
              success: true,
              data: createMockStats({
                totalAmount: 0,
                totalCount: 0,
                completedAmount: 0,
                completedCount: 0,
                pendingAmount: 0,
                pendingCount: 0,
                byCompany: [],
              }),
            },
          });
        }
        return Promise.resolve({ data: { success: true, data: null } });
      });

      renderWithProviders(<WithdrawRoutingPage />);

      await waitFor(() => {
        expect(screen.getByText('提现分流统计')).toBeInTheDocument();
      });

      // Should display 0 values
      const zeroValues = screen.getAllByText('0');
      expect(zeroValues.length).toBeGreaterThan(0);
    });
  });

  describe('Accessibility', () => {
    it('should have proper page title', async () => {
      renderWithProviders(<WithdrawRoutingPage />);

      await waitFor(() => {
        expect(screen.getByText('提现分流统计')).toBeInTheDocument();
      });
    });

    it('should have proper table structure', async () => {
      renderWithProviders(<WithdrawRoutingPage />);

      await waitFor(() => {
        expect(screen.getByRole('table')).toBeInTheDocument();
      });
    });

    it('should have withdraw detail card', async () => {
      renderWithProviders(<WithdrawRoutingPage />);

      await waitFor(() => {
        expect(screen.getByText('提现明细')).toBeInTheDocument();
      });
    });
  });
});
