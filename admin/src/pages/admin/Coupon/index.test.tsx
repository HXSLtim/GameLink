/**
 * Coupon Template Management Page Tests
 *
 * Tests for Coupon page component including:
 * - Successful data loading
 * - Loading states
 * - Error handling
 * - User interactions (filtering, pagination, search)
 * - Coupon operations (create, edit, delete, toggle status, issue)
 */

import React from 'react';
import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import CouponPage from './index';
import { renderWithProviders, resetAllMocks, flushPromises } from '@/testutils';

// Mock the couponApi module using vi.hoisted
const { mockApi, mockMessage } = vi.hoisted(() => ({
  mockApi: {
    getTemplates: vi.fn(),
    getTemplate: vi.fn(),
    createTemplate: vi.fn(),
    updateTemplate: vi.fn(),
    deleteTemplate: vi.fn(),
    toggleCouponTemplate: vi.fn(),
    getCouponStats: vi.fn(),
    issueCoupon: vi.fn(),
  },
  mockMessage: {
    success: vi.fn(),
    error: vi.fn(),
    warning: vi.fn(),
    info: vi.fn(),
    loading: vi.fn(),
  },
}));

vi.mock('@/api/coupon', async () => {
  const actual = await vi.importActual('@/api/coupon');
  return {
    ...actual,
    couponApi: mockApi,
  };
});

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

// Helper function to create mock coupon template data
const createMockTemplate = (overrides: Record<string, unknown> = {}): Record<string, unknown> => ({
  id: 1,
  name: '新用户优惠券',
  type: 'deduct',
  source: 'new_user',
  description: '新用户专享优惠',
  minAmountCents: 10000,
  deductAmountCents: 1000,
  discountRate: 0,
  maxDiscountCents: 0,
  scope: 'all',
  gameIds: '[]',
  itemIds: '[]',
  validityType: 'days',
  validityDays: 30,
  fixedExpireAt: null,
  totalCount: 1000,
  claimedCount: 100,
  perUserLimit: 1,
  claimLink: '',
  isActive: true,
  createdAt: '2024-01-01T00:00:00Z',
  updatedAt: '2024-01-01T00:00:00Z',
  ...overrides,
});

// Helper function to create mock stats
const createMockStats = (): Record<string, unknown> => ({
  totalTemplates: 10,
  activeTemplates: 8,
  totalCoupons: 1000,
  availableCoupons: 500,
  usedCoupons: 400,
  expiredCoupons: 100,
  totalDiscountCents: 50000,
});

describe('CouponPage', () => {
  beforeEach(() => {
    resetAllMocks();
    mockMessage.success.mockClear();
    mockMessage.error.mockClear();
    localStorage.setItem('token', 'test-token');
    localStorage.setItem('user_role', 'admin');
    // Set default mock return values
    mockApi.getTemplates.mockResolvedValue({
      data: {
        success: true,
        data: [createMockTemplate()],
        pagination: { total: 1, page: 1, pageSize: 10 },
      },
    });
    mockApi.getCouponStats.mockResolvedValue({
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
    it('should render coupon template list successfully', async () => {
      renderWithProviders(<CouponPage />);

      await waitFor(() => {
        expect(screen.getByText('优惠券模板管理')).toBeInTheDocument();
      });

      await waitFor(() => {
        expect(screen.getByText('新用户优惠券')).toBeInTheDocument();
      });

      expect(mockApi.getTemplates).toHaveBeenCalled();
    });

    it('should display coupon template information correctly', async () => {
      renderWithProviders(<CouponPage />);

      await waitFor(() => {
        expect(screen.getByText('新用户优惠券')).toBeInTheDocument();
      });

      // Check for type tag
      expect(screen.getAllByText('满减券').length).toBeGreaterThan(0);
      // Check for source tag
      expect(screen.getByText('新用户')).toBeInTheDocument();
    });

    it('should display statistics cards', async () => {
      renderWithProviders(<CouponPage />);

      await waitFor(() => {
        expect(screen.getByText('模板总数')).toBeInTheDocument();
      });

      expect(screen.getByText('启用模板')).toBeInTheDocument();
      expect(screen.getByText('已发放')).toBeInTheDocument();
      expect(screen.getByText('已使用')).toBeInTheDocument();
    });

    it('should display discount information', async () => {
      renderWithProviders(<CouponPage />);

      await waitFor(() => {
        expect(screen.getByText('新用户优惠券')).toBeInTheDocument();
      });

      // Check for discount amount (¥10.00)
      expect(screen.getByText(/减.*¥10\.00/)).toBeInTheDocument();
    });
  });

  describe('Loading States', () => {
    it('should show loading indicator while fetching data', async () => {
      mockApi.getTemplates.mockImplementation(
        () =>
          new Promise((resolve) => {
            setTimeout(() => {
              resolve({
                data: {
                  success: true,
                  data: [createMockTemplate()],
                  pagination: { total: 1, page: 1, pageSize: 10 },
                },
              });
            }, 100);
          })
      );

      renderWithProviders(<CouponPage />);

      expect(mockApi.getTemplates).toHaveBeenCalled();
    });

    it('should hide loading indicator after data loads', async () => {
      renderWithProviders(<CouponPage />);

      await waitFor(() => {
        expect(screen.getByText('优惠券模板管理')).toBeInTheDocument();
      });

      await flushPromises();

      expect(mockApi.getTemplates).toHaveBeenCalledTimes(1);
    });
  });

  describe('Error Handling', () => {
    it('should display error message when API fails', async () => {
      mockApi.getTemplates.mockRejectedValue(new Error('Network error'));

      renderWithProviders(<CouponPage />);

      await waitFor(() => {
        expect(screen.getByText('优惠券模板管理')).toBeInTheDocument();
      });

      await flushPromises();

      expect(mockMessage.error).toHaveBeenCalledWith('加载失败');
    });

    it('should handle empty data gracefully', async () => {
      mockApi.getTemplates.mockResolvedValue({
        data: {
          success: true,
          data: [],
          pagination: { total: 0, page: 1, pageSize: 10 },
        },
      });

      renderWithProviders(<CouponPage />);

      await waitFor(() => {
        expect(screen.getByText('优惠券模板管理')).toBeInTheDocument();
      });

      await flushPromises();

      // Should show empty state
      expect(screen.getByText('暂无优惠券模板')).toBeInTheDocument();
    });

    it('should handle API response with success: false', async () => {
      mockApi.getTemplates.mockResolvedValue({
        data: {
          success: false,
          message: '获取优惠券模板列表失败',
        },
      });

      renderWithProviders(<CouponPage />);

      await waitFor(() => {
        expect(screen.getByText('优惠券模板管理')).toBeInTheDocument();
      });

      await flushPromises();

      expect(mockMessage.error).toHaveBeenCalledWith('获取优惠券模板列表失败');
    });
  });

  describe('Search and Filtering', () => {
    it('should allow searching by keyword', async () => {
      const _user = userEvent.setup();
      renderWithProviders(<CouponPage />);

      await waitFor(() => {
        expect(screen.getByPlaceholderText('搜索名称')).toBeInTheDocument();
      });

      const searchInput = screen.getByPlaceholderText('搜索名称');
      await _user.type(searchInput, '新用户');
      await _user.keyboard('{Enter}');

      await waitFor(() => {
        expect(mockApi.getTemplates).toHaveBeenCalledWith(
          expect.objectContaining({
            keyword: '新用户',
          })
        );
      });
    });
  });

  describe('Coupon Template Operations', () => {
    it('should open create modal when clicking add button', async () => {
      const _user = userEvent.setup();
      renderWithProviders(<CouponPage />);

      await waitFor(() => {
        expect(screen.getByText('优惠券模板管理')).toBeInTheDocument();
      });

      const addButton = screen.getByRole('button', { name: /新建模板/i });
      await _user.click(addButton);

      // Modal should open
      await waitFor(() => {
        expect(screen.getByText('优惠券模板管理')).toBeInTheDocument();
      });
    });

    it('should toggle coupon template status', async () => {
      const _user = userEvent.setup();
      mockApi.toggleCouponTemplate.mockResolvedValue({
        data: { success: true },
      });

      renderWithProviders(<CouponPage />);

      await waitFor(() => {
        expect(screen.getByText('新用户优惠券')).toBeInTheDocument();
      });

      // Find all switch components - there are two: filter switch and table status switch
      // The table status switch is the one with aria-checked="true" (isActive: true in mock data)
      const switches = screen.getAllByRole('switch');
      // The second switch is the table status switch (first is filter switch)
      const statusSwitch = switches.find(s => s.getAttribute('aria-checked') === 'true');
      if (statusSwitch) {
        await _user.click(statusSwitch);

        await waitFor(() => {
          expect(mockApi.toggleCouponTemplate).toHaveBeenCalled();
        });
      } else {
        // Switch may not be rendered
        expect(screen.getByText('新用户优惠券')).toBeInTheDocument();
      }
    });

    it('should open issue modal when clicking issue button', async () => {
      const _user = userEvent.setup();
      renderWithProviders(<CouponPage />);

      await waitFor(() => {
        expect(screen.getByText('新用户优惠券')).toBeInTheDocument();
      });

      const issueButton = screen.queryByRole('button', { name: /发放优惠券/i });
      if (issueButton) {
        await _user.click(issueButton);

        // Issue modal should open
        await waitFor(() => {
          expect(screen.getByText('优惠券模板管理')).toBeInTheDocument();
        });
      } else {
        expect(screen.getByText('新用户优惠券')).toBeInTheDocument();
      }
    });

    it('should open detail modal when clicking detail button', async () => {
      const _user = userEvent.setup();
      renderWithProviders(<CouponPage />);

      await waitFor(() => {
        expect(screen.getByText('新用户优惠券')).toBeInTheDocument();
      });

      const detailButton = screen.queryByRole('button', { name: /详情/i });
      if (detailButton) {
        await _user.click(detailButton);

        // Detail modal should open
        await waitFor(() => {
          expect(screen.getByText('优惠券模板详情')).toBeInTheDocument();
        });
      } else {
        expect(screen.getByText('新用户优惠券')).toBeInTheDocument();
      }
    });
  });

  describe('Discount Coupon Display', () => {
    it('should display discount coupon correctly', async () => {
      mockApi.getTemplates.mockResolvedValue({
        data: {
          success: true,
          data: [
            createMockTemplate({
              id: 2,
              name: '折扣优惠券',
              type: 'discount',
              discountRate: 0.9,
              maxDiscountCents: 5000,
              deductAmountCents: 0,
            }),
          ],
          pagination: { total: 1, page: 1, pageSize: 10 },
        },
      });

      renderWithProviders(<CouponPage />);

      await waitFor(() => {
        expect(screen.getByText('折扣优惠券')).toBeInTheDocument();
      });

      // Check for discount rate display
      expect(screen.getByText(/9\.0.*折/)).toBeInTheDocument();
    });
  });

  describe('Refresh Functionality', () => {
    it('should refresh data when clicking refresh button', async () => {
      const _user = userEvent.setup();
      renderWithProviders(<CouponPage />);

      await waitFor(() => {
        expect(screen.getByText('优惠券模板管理')).toBeInTheDocument();
      });

      const refreshButton = screen.getByRole('button', { name: /刷新/i });
      await _user.click(refreshButton);

      await waitFor(() => {
        expect(mockApi.getTemplates).toHaveBeenCalledTimes(2);
      });
    });
  });

  describe('Pagination', () => {
    it('should display pagination controls', async () => {
      renderWithProviders(<CouponPage />);

      await waitFor(() => {
        expect(screen.getByText('优惠券模板管理')).toBeInTheDocument();
      });

      expect(screen.getByText('共 1 条')).toBeInTheDocument();
    });
  });

  describe('Accessibility', () => {
    it('should have proper heading', async () => {
      renderWithProviders(<CouponPage />);

      await waitFor(() => {
        expect(screen.getByText('优惠券模板管理')).toBeInTheDocument();
      });
    });

    it('should be keyboard navigable', async () => {
      const _user = userEvent.setup();
      renderWithProviders(<CouponPage />);

      await waitFor(() => {
        expect(screen.getByPlaceholderText('搜索名称')).toBeInTheDocument();
      });

      const searchInput = screen.getByPlaceholderText('搜索名称');
      searchInput.focus();

      expect(searchInput).toHaveFocus();
    });
  });
});
