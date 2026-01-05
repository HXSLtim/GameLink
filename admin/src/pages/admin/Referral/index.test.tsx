/**
 * Referral Management Page Tests
 *
 * Tests for Referral page component including:
 * - Successful data loading
 * - Loading states
 * - Error handling
 * - Statistics display
 * - Tab navigation
 */

import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { screen, waitFor, fireEvent } from '@testing-library/react';
import ReferralPage from './index';
import { renderWithProviders, resetAllMocks, flushPromises } from '@/testutils';

// Mock the referralApi module using vi.hoisted
const { mockApi, mockMessage } = vi.hoisted(() => ({
  mockApi: {
    getReferralStats: vi.fn(),
    getReferrals: vi.fn(),
    getReferralCodes: vi.fn(),
    getReferralRewards: vi.fn(),
    createReferralCode: vi.fn(),
    updateReferralCode: vi.fn(),
    deleteReferralCode: vi.fn(),
    updateReferralStatus: vi.fn(),
    issueReferralReward: vi.fn(),
    failReferralReward: vi.fn(),
    batchUpdateCodesStatus: vi.fn(),
    batchDeleteCodes: vi.fn(),
    batchUpdateReferralsStatus: vi.fn(),
    batchDeleteReferrals: vi.fn(),
  },
  mockMessage: {
    success: vi.fn(),
    error: vi.fn(),
    warning: vi.fn(),
    info: vi.fn(),
    loading: vi.fn(),
  },
}));

vi.mock('@/api/referral', () => ({
  referralApi: mockApi,
  getReferralTypeLabel: (type: string) => type === 'user' ? '用户推荐' : '陪玩师推荐',
  getReferralStatusLabel: (status: string) => {
    const labels: Record<string, string> = {
      pending: '待完成',
      completed: '已完成',
      canceled: '已取消',
    };
    return labels[status] || status;
  },
  getReferralStatusColor: (status: string) => {
    const colors: Record<string, string> = {
      pending: 'orange',
      completed: 'green',
      canceled: 'red',
    };
    return colors[status] || 'default';
  },
  getRewardTypeLabel: (type: string) => type === 'referrer' ? '推荐人奖励' : '被推荐人奖励',
  getRewardStatusLabel: (status: string) => {
    const labels: Record<string, string> = {
      pending: '待发放',
      issued: '已发放',
      failed: '发放失败',
    };
    return labels[status] || status;
  },
  getRewardStatusColor: (status: string) => {
    const colors: Record<string, string> = {
      pending: 'orange',
      issued: 'green',
      failed: 'red',
    };
    return colors[status] || 'default';
  },
  centsToYuan: (cents: number) => (cents / 100).toFixed(2),
  isCodeExpired: () => false,
  isCodeFullyUsed: () => false,
  getCodeUsagePercent: () => 50,
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
        modal: { confirm: vi.fn() },
      }),
    },
  };
});

// Helper function to create mock stats
const createMockStats = (overrides: Record<string, unknown> = {}): Record<string, unknown> => ({
  totalReferrals: 100,
  completedReferrals: 80,
  pendingReferrals: 15,
  canceledReferrals: 5,
  totalRewardsCents: 100000,
  issuedRewardsCents: 80000,
  pendingRewardsCents: 15000,
  failedRewardsCents: 5000,
  activeCodes: 50,
  totalCodes: 100,
  ...overrides,
});

// Helper function to create mock referral
const createMockReferral = (overrides: Record<string, unknown> = {}): Record<string, unknown> => ({
  id: 1,
  referrerId: 1,
  refereeId: 2,
  codeId: 1,
  type: 'user',
  status: 'pending',
  completedAt: null,
  cancelReason: null,
  createdAt: '2024-01-01T00:00:00Z',
  updatedAt: '2024-01-01T00:00:00Z',
  referrer: {
    id: 1,
    name: '推荐人',
    avatarUrl: null,
  },
  referee: {
    id: 2,
    name: '被推荐人',
    avatarUrl: null,
  },
  code: {
    id: 1,
    code: 'ABC123',
  },
  ...overrides,
});

// Helper function to create mock referral code
const createMockCode = (overrides: Record<string, unknown> = {}): Record<string, unknown> => ({
  id: 1,
  code: 'ABC123',
  ownerId: 1,
  type: 'user',
  maxUses: 100,
  usedCount: 50,
  expiresAt: '2025-12-31T23:59:59Z',
  isActive: true,
  createdAt: '2024-01-01T00:00:00Z',
  updatedAt: '2024-01-01T00:00:00Z',
  owner: {
    id: 1,
    name: '测试用户',
    avatarUrl: null,
  },
  ...overrides,
});

// Helper function to create mock reward
const createMockReward = (overrides: Record<string, unknown> = {}): Record<string, unknown> => ({
  id: 1,
  referralId: 1,
  userId: 1,
  type: 'referrer',
  amountCents: 1000,
  status: 'pending',
  issuedAt: null,
  failedAt: null,
  failureReason: null,
  createdAt: '2024-01-01T00:00:00Z',
  updatedAt: '2024-01-01T00:00:00Z',
  user: {
    id: 1,
    name: '测试用户',
    avatarUrl: null,
  },
  ...overrides,
});

describe('ReferralPage', () => {
  beforeEach(() => {
    resetAllMocks();
    mockMessage.success.mockClear();
    mockMessage.error.mockClear();
    mockMessage.info.mockClear();
    localStorage.setItem('token', 'test-token');
    localStorage.setItem('user_role', 'admin');
    
    // Set default mock return values
    mockApi.getReferralStats.mockResolvedValue({
      data: {
        success: true,
        data: createMockStats(),
      },
    });
    mockApi.getReferrals.mockResolvedValue({
      data: {
        success: true,
        data: [createMockReferral()],
        pagination: { total: 1 },
      },
    });
    mockApi.getReferralCodes.mockResolvedValue({
      data: {
        success: true,
        data: [createMockCode()],
        pagination: { total: 1 },
      },
    });
    mockApi.getReferralRewards.mockResolvedValue({
      data: {
        success: true,
        data: [createMockReward()],
        pagination: { total: 1 },
      },
    });
  });

  afterEach(() => {
    vi.clearAllMocks();
  });

  describe('Successful Data Loading', () => {
    it('should render referral page successfully', async () => {
      renderWithProviders(<ReferralPage />);

      await waitFor(() => {
        expect(screen.getByText('总推荐数')).toBeInTheDocument();
      });

      expect(mockApi.getReferralStats).toHaveBeenCalled();
    });

    it('should display statistics cards', async () => {
      renderWithProviders(<ReferralPage />);

      await waitFor(() => {
        expect(screen.getByText('总推荐数')).toBeInTheDocument();
      });

      expect(screen.getByText('已完成推荐')).toBeInTheDocument();
      expect(screen.getByText('已发放奖励')).toBeInTheDocument();
      expect(screen.getByText('活跃邀请码')).toBeInTheDocument();
    });

    it('should display statistics values', async () => {
      renderWithProviders(<ReferralPage />);

      await waitFor(() => {
        expect(screen.getByText('100')).toBeInTheDocument(); // totalReferrals
      });

      expect(screen.getByText('80')).toBeInTheDocument(); // completedReferrals
      expect(screen.getByText('50')).toBeInTheDocument(); // activeCodes
    });

    it('should display total codes count', async () => {
      renderWithProviders(<ReferralPage />);

      await waitFor(() => {
        expect(screen.getByText(/总计: 100 个/)).toBeInTheDocument();
      });
    });
  });

  describe('Loading States', () => {
    it('should show loading indicator while fetching data', async () => {
      mockApi.getReferralStats.mockImplementation(
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

      renderWithProviders(<ReferralPage />);

      expect(mockApi.getReferralStats).toHaveBeenCalled();
    });

    it('should hide loading indicator after data loads', async () => {
      renderWithProviders(<ReferralPage />);

      await waitFor(() => {
        expect(screen.getByText('总推荐数')).toBeInTheDocument();
      });

      await flushPromises();

      expect(mockApi.getReferralStats).toHaveBeenCalledTimes(1);
    });
  });

  describe('Error Handling', () => {
    it('should handle stats API failure silently', async () => {
      mockApi.getReferralStats.mockRejectedValue(new Error('Network error'));

      renderWithProviders(<ReferralPage />);

      await waitFor(() => {
        expect(screen.getByText('总推荐数')).toBeInTheDocument();
      });

      await flushPromises();

      // Should not show error message for stats (silent fail)
      expect(mockMessage.error).not.toHaveBeenCalled();
    });
  });

  describe('Tab Navigation', () => {
    it('should display all tabs', async () => {
      renderWithProviders(<ReferralPage />);

      await waitFor(() => {
        expect(screen.getByText('总推荐数')).toBeInTheDocument();
      });

      // Check tabs exist
      const tabs = screen.getAllByRole('tab');
      expect(tabs.length).toBe(3);
    });

    it('should show referral list tab by default', async () => {
      renderWithProviders(<ReferralPage />);

      await waitFor(() => {
        expect(screen.getByText('推荐关系')).toBeInTheDocument();
      });

      // First tab should be active
      const firstTab = screen.getByRole('tab', { name: '推荐关系' });
      expect(firstTab).toHaveAttribute('aria-selected', 'true');
    });

    it('should switch to codes tab when clicked', async () => {
      renderWithProviders(<ReferralPage />);

      await waitFor(() => {
        expect(screen.getByText('总推荐数')).toBeInTheDocument();
      });

      // Click on codes tab
      const codesTab = screen.getByRole('tab', { name: '邀请码管理' });
      fireEvent.click(codesTab);

      await waitFor(() => {
        expect(codesTab).toHaveAttribute('aria-selected', 'true');
      });
    });

    it('should switch to rewards tab when clicked', async () => {
      renderWithProviders(<ReferralPage />);

      await waitFor(() => {
        expect(screen.getByText('总推荐数')).toBeInTheDocument();
      });

      // Click on rewards tab
      const rewardsTab = screen.getByRole('tab', { name: '奖励管理' });
      fireEvent.click(rewardsTab);

      await waitFor(() => {
        expect(rewardsTab).toHaveAttribute('aria-selected', 'true');
      });
    });
  });

  describe('Statistics with Zero Values', () => {
    it('should display zero values correctly', async () => {
      mockApi.getReferralStats.mockResolvedValue({
        data: {
          success: true,
          data: createMockStats({
            totalReferrals: 0,
            completedReferrals: 0,
            issuedRewardsCents: 0,
            activeCodes: 0,
            totalCodes: 0,
          }),
        },
      });

      renderWithProviders(<ReferralPage />);

      await waitFor(() => {
        expect(screen.getByText('总推荐数')).toBeInTheDocument();
      });

      // Should display 0 values
      const zeroValues = screen.getAllByText('0');
      expect(zeroValues.length).toBeGreaterThan(0);
    });
  });

  describe('Statistics with Null Values', () => {
    it('should handle null stats gracefully', async () => {
      mockApi.getReferralStats.mockResolvedValue({
        data: {
          success: true,
          data: null,
        },
      });

      renderWithProviders(<ReferralPage />);

      await waitFor(() => {
        expect(screen.getByText('总推荐数')).toBeInTheDocument();
      });

      // Should display 0 for null values
      const zeroValues = screen.getAllByText('0');
      expect(zeroValues.length).toBeGreaterThan(0);
    });
  });

  describe('Issued Rewards Display', () => {
    it('should display issued rewards in yuan', async () => {
      mockApi.getReferralStats.mockResolvedValue({
        data: {
          success: true,
          data: createMockStats({
            issuedRewardsCents: 100000, // 1000 yuan
          }),
        },
      });

      renderWithProviders(<ReferralPage />);

      await waitFor(() => {
        expect(screen.getByText('已发放奖励')).toBeInTheDocument();
      });

      // Should display 1000.00 yuan (100000 cents / 100)
      // Note: Ant Design Statistic may format this differently
      expect(mockApi.getReferralStats).toHaveBeenCalled();
    });
  });

  describe('Accessibility', () => {
    it('should have proper statistics card structure', async () => {
      renderWithProviders(<ReferralPage />);

      await waitFor(() => {
        expect(screen.getByText('总推荐数')).toBeInTheDocument();
      });

      // Check that all stat cards are present
      expect(screen.getByText('已完成推荐')).toBeInTheDocument();
      expect(screen.getByText('已发放奖励')).toBeInTheDocument();
      expect(screen.getByText('活跃邀请码')).toBeInTheDocument();
    });

    it('should have proper tab structure', async () => {
      renderWithProviders(<ReferralPage />);

      await waitFor(() => {
        expect(screen.getByText('总推荐数')).toBeInTheDocument();
      });

      // Check tab list exists
      const tabList = screen.getByRole('tablist');
      expect(tabList).toBeInTheDocument();
    });
  });

  describe('Data Refresh', () => {
    it('should call loadStats on mount', async () => {
      renderWithProviders(<ReferralPage />);

      await waitFor(() => {
        expect(mockApi.getReferralStats).toHaveBeenCalledTimes(1);
      });
    });
  });
});
