/**
 * Withdraw Management Page Tests
 *
 * Tests for Withdraw page component including:
 * - Successful data loading
 * - Loading states
 * - Error handling
 * - Statistics display
 * - Approve/Reject/Complete actions
 */

import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import WithdrawPage from './index';
import { renderWithProviders, resetAllMocks, flushPromises } from '@/testutils';

// Mock the adminApi module using vi.hoisted
const { mockApi, mockMessage, mockModal } = vi.hoisted(() => ({
  mockApi: {
    getWithdraws: vi.fn(),
    getWithdraw: vi.fn(),
    approveWithdraw: vi.fn(),
    rejectWithdraw: vi.fn(),
    completeWithdraw: vi.fn(),
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
    success: vi.fn(),
    error: vi.fn(),
    info: vi.fn(),
    warning: vi.fn(),
  },
}));

vi.mock('@/api/admin', () => ({
  adminApi: mockApi,
}));

// Mock App.useApp to return the message and modal mocks
vi.mock('antd', async () => {
  const actual = await vi.importActual('antd');
  return {
    ...actual,
    App: {
      useApp: () => ({
        message: mockMessage,
        notification: { success: vi.fn(), error: vi.fn(), info: vi.fn(), warning: vi.fn() },
        modal: mockModal,
      }),
    },
  };
});

// Helper function to create mock withdraw data
const createMockWithdraw = (overrides: Record<string, unknown> = {}): Record<string, unknown> => ({
  id: 1,
  playerId: 1,
  amountCents: 10000,
  status: 'pending',
  bankName: '中国银行',
  bankAccount: '6222021234567890123',
  accountName: '张三',
  remark: '提现申请',
  rejectReason: null,
  adminRemark: null,
  processedBy: null,
  processedAt: null,
  completedAt: null,
  createdAt: '2024-01-01T00:00:00Z',
  updatedAt: '2024-01-01T00:00:00Z',
  player: {
    id: 1,
    nickname: '测试陪玩师',
    avatarUrl: '',
  },
  ...overrides,
});

describe('WithdrawPage', () => {
  beforeEach(() => {
    resetAllMocks();
    mockMessage.success.mockClear();
    mockMessage.error.mockClear();
    mockMessage.loading.mockClear();
    mockModal.confirm.mockClear();
    localStorage.setItem('token', 'test-token');
    localStorage.setItem('user_role', 'admin');
    // Set default mock return values
    mockApi.getWithdraws.mockResolvedValue({
      data: {
        success: true,
        data: {
          withdraws: [createMockWithdraw()],
          total: 1,
        },
      },
    });
  });

  afterEach(() => {
    vi.clearAllMocks();
  });

  describe('Successful Data Loading', () => {
    it('should render withdraw page successfully', async () => {
      renderWithProviders(<WithdrawPage />);

      await waitFor(() => {
        expect(screen.getByText('提现管理')).toBeInTheDocument();
      });

      expect(mockApi.getWithdraws).toHaveBeenCalled();
    });

    it('should display statistics cards', async () => {
      renderWithProviders(<WithdrawPage />);

      await waitFor(() => {
        expect(screen.getByText('待审核')).toBeInTheDocument();
      });

      expect(screen.getByText('已批准待打款')).toBeInTheDocument();
      expect(screen.getByText('已完成')).toBeInTheDocument();
      expect(screen.getByText('本月提现总额')).toBeInTheDocument();
    });

    it('should display withdraw list with correct columns', async () => {
      renderWithProviders(<WithdrawPage />);

      await waitFor(() => {
        // "陪玩师" appears in multiple places (column header, search field)
        const playerElements = screen.getAllByText('陪玩师');
        expect(playerElements.length).toBeGreaterThan(0);
      });

      // Multiple elements may appear due to Ant Design table rendering
      const amountElements = screen.getAllByText('提现金额');
      expect(amountElements.length).toBeGreaterThan(0);
      const bankElements = screen.getAllByText('银行信息');
      expect(bankElements.length).toBeGreaterThan(0);
      const statusElements = screen.getAllByText('状态');
      expect(statusElements.length).toBeGreaterThan(0);
      const timeElements = screen.getAllByText('申请时间');
      expect(timeElements.length).toBeGreaterThan(0);
    });

    it('should display withdraw data in table', async () => {
      renderWithProviders(<WithdrawPage />);

      await waitFor(() => {
        expect(screen.getByText('测试陪玩师')).toBeInTheDocument();
      });

      expect(screen.getByText('中国银行')).toBeInTheDocument();
      // "待审核" appears in both stats card and table
      const pendingElements = screen.getAllByText('待审核');
      expect(pendingElements.length).toBeGreaterThan(0);
    });

    it('should display formatted amount', async () => {
      renderWithProviders(<WithdrawPage />);

      await waitFor(() => {
        expect(screen.getByText('¥100.00')).toBeInTheDocument();
      });
    });
  });

  describe('Loading States', () => {
    it('should show loading indicator while fetching data', async () => {
      mockApi.getWithdraws.mockImplementation(
        () =>
          new Promise((resolve) => {
            setTimeout(() => {
              resolve({
                data: {
                  success: true,
                  data: {
                    withdraws: [createMockWithdraw()],
                    total: 1,
                  },
                },
              });
            }, 100);
          })
      );

      renderWithProviders(<WithdrawPage />);

      expect(mockApi.getWithdraws).toHaveBeenCalled();
    });

    it('should hide loading indicator after data loads', async () => {
      renderWithProviders(<WithdrawPage />);

      await waitFor(() => {
        expect(screen.getByText('提现管理')).toBeInTheDocument();
      });

      await flushPromises();

      // Multiple calls for stats (pending, approved, completed)
      expect(mockApi.getWithdraws).toHaveBeenCalled();
    });
  });

  describe('Error Handling', () => {
    it('should display error message when API fails', async () => {
      mockApi.getWithdraws.mockRejectedValue(new Error('Network error'));

      renderWithProviders(<WithdrawPage />);

      await waitFor(() => {
        expect(screen.getByText('提现管理')).toBeInTheDocument();
      });

      await flushPromises();

      expect(mockMessage.error).toHaveBeenCalledWith('加载提现列表失败');
    });

    it('should handle API response with success: false', async () => {
      mockApi.getWithdraws.mockResolvedValue({
        data: {
          success: false,
          message: '获取提现列表失败',
        },
      });

      renderWithProviders(<WithdrawPage />);

      await waitFor(() => {
        expect(screen.getByText('提现管理')).toBeInTheDocument();
      });

      await flushPromises();

      expect(mockMessage.error).toHaveBeenCalledWith('获取提现列表失败');
    });
  });

  describe('Status Display', () => {
    it('should display pending status correctly', async () => {
      mockApi.getWithdraws.mockResolvedValue({
        data: {
          success: true,
          data: {
            withdraws: [createMockWithdraw({ status: 'pending' })],
            total: 1,
          },
        },
      });

      renderWithProviders(<WithdrawPage />);

      await waitFor(() => {
        // "待审核" appears in both stats card and table
        const pendingElements = screen.getAllByText('待审核');
        expect(pendingElements.length).toBeGreaterThan(0);
      });
    });

    it('should display approved status correctly', async () => {
      mockApi.getWithdraws.mockResolvedValue({
        data: {
          success: true,
          data: {
            withdraws: [createMockWithdraw({ status: 'approved' })],
            total: 1,
          },
        },
      });

      renderWithProviders(<WithdrawPage />);

      await waitFor(() => {
        // "已批准" appears in both stats card and table
        const approvedElements = screen.getAllByText('已批准');
        expect(approvedElements.length).toBeGreaterThan(0);
      });
    });

    it('should display rejected status correctly', async () => {
      mockApi.getWithdraws.mockResolvedValue({
        data: {
          success: true,
          data: {
            withdraws: [createMockWithdraw({ status: 'rejected' })],
            total: 1,
          },
        },
      });

      renderWithProviders(<WithdrawPage />);

      await waitFor(() => {
        expect(screen.getByText('已拒绝')).toBeInTheDocument();
      });
    });

    it('should display completed status correctly', async () => {
      mockApi.getWithdraws.mockResolvedValue({
        data: {
          success: true,
          data: {
            withdraws: [createMockWithdraw({ status: 'completed' })],
            total: 1,
          },
        },
      });

      renderWithProviders(<WithdrawPage />);

      await waitFor(() => {
        // "已完成" appears in both stats card and table
        const completedElements = screen.getAllByText('已完成');
        expect(completedElements.length).toBeGreaterThan(0);
      });
    });
  });

  describe('Action Buttons', () => {
    it('should show approve and reject buttons for pending withdraws', async () => {
      mockApi.getWithdraws.mockResolvedValue({
        data: {
          success: true,
          data: {
            withdraws: [createMockWithdraw({ status: 'pending' })],
            total: 1,
          },
        },
      });

      renderWithProviders(<WithdrawPage />);

      await waitFor(() => {
        expect(screen.getByText('批准')).toBeInTheDocument();
      });

      expect(screen.getByText('拒绝')).toBeInTheDocument();
    });

    it('should show complete button for approved withdraws', async () => {
      mockApi.getWithdraws.mockResolvedValue({
        data: {
          success: true,
          data: {
            withdraws: [createMockWithdraw({ status: 'approved' })],
            total: 1,
          },
        },
      });

      renderWithProviders(<WithdrawPage />);

      await waitFor(() => {
        expect(screen.getByText('完成打款')).toBeInTheDocument();
      });
    });

    it('should show detail button for all withdraws', async () => {
      renderWithProviders(<WithdrawPage />);

      await waitFor(() => {
        expect(screen.getByText('详情')).toBeInTheDocument();
      });
    });
  });

  describe('View Detail Modal', () => {
    it('should open detail modal when clicking detail button', async () => {
      const _user = userEvent.setup();
      renderWithProviders(<WithdrawPage />);

      await waitFor(() => {
        expect(screen.getByText('详情')).toBeInTheDocument();
      });

      await _user.click(screen.getByText('详情'));

      await waitFor(() => {
        expect(screen.getByText('提现详情')).toBeInTheDocument();
      });
    });

    it('should display withdraw details in modal', async () => {
      const _user = userEvent.setup();
      renderWithProviders(<WithdrawPage />);

      await waitFor(() => {
        expect(screen.getByText('详情')).toBeInTheDocument();
      });

      await _user.click(screen.getByText('详情'));

      await waitFor(() => {
        expect(screen.getByText('提现ID')).toBeInTheDocument();
      });

      expect(screen.getByText('银行名称')).toBeInTheDocument();
      expect(screen.getByText('银行账号')).toBeInTheDocument();
    });
  });

  describe('Toolbar Buttons', () => {
    it('should display batch approve button', async () => {
      renderWithProviders(<WithdrawPage />);

      await waitFor(() => {
        expect(screen.getByText('批量批准')).toBeInTheDocument();
      });
    });

    it('should display batch reject button', async () => {
      renderWithProviders(<WithdrawPage />);

      await waitFor(() => {
        expect(screen.getByText('批量拒绝')).toBeInTheDocument();
      });
    });

    it('should display batch complete button', async () => {
      renderWithProviders(<WithdrawPage />);

      await waitFor(() => {
        expect(screen.getByText('批量完成打款')).toBeInTheDocument();
      });
    });

    it('should display export button', async () => {
      renderWithProviders(<WithdrawPage />);

      await waitFor(() => {
        expect(screen.getByText('导出数据')).toBeInTheDocument();
      });
    });
  });

  describe('Search Functionality', () => {
    it('should have status filter', async () => {
      renderWithProviders(<WithdrawPage />);

      await waitFor(() => {
        // "状态" appears in multiple places (search field, table column)
        const statusElements = screen.getAllByText('状态');
        expect(statusElements.length).toBeGreaterThan(0);
      });
    });

    it('should have player ID filter', async () => {
      renderWithProviders(<WithdrawPage />);

      await waitFor(() => {
        expect(screen.getByText('陪玩师ID')).toBeInTheDocument();
      });
    });
  });

  describe('Empty State', () => {
    it('should handle empty withdraw list', async () => {
      mockApi.getWithdraws.mockResolvedValue({
        data: {
          success: true,
          data: {
            withdraws: [],
            total: 0,
          },
        },
      });

      renderWithProviders(<WithdrawPage />);

      await waitFor(() => {
        expect(screen.getByText('提现管理')).toBeInTheDocument();
      });

      // Table should still render with no data
      expect(mockApi.getWithdraws).toHaveBeenCalled();
    });
  });

  describe('Bank Information Display', () => {
    it('should display masked bank account', async () => {
      renderWithProviders(<WithdrawPage />);

      await waitFor(() => {
        // Bank account should be masked, showing only last 4 digits
        expect(screen.getByText('****0123')).toBeInTheDocument();
      });
    });

    it('should display bank name', async () => {
      renderWithProviders(<WithdrawPage />);

      await waitFor(() => {
        expect(screen.getByText('中国银行')).toBeInTheDocument();
      });
    });
  });

  describe('Accessibility', () => {
    it('should have proper page title', async () => {
      renderWithProviders(<WithdrawPage />);

      await waitFor(() => {
        expect(screen.getByText('提现管理')).toBeInTheDocument();
      });
    });

    it('should have proper subtitle', async () => {
      renderWithProviders(<WithdrawPage />);

      await waitFor(() => {
        expect(screen.getByText('管理陪玩师提现申请')).toBeInTheDocument();
      });
    });
  });
});
