/**
 * Review Management Page Tests
 */

import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { screen, waitFor, fireEvent } from '@testing-library/react';
import ReviewList from './index';
import { renderWithProviders, resetAllMocks } from '@/testutils';

// Mock the reviewApi module using vi.hoisted
const { mockApi, mockMessage } = vi.hoisted(() => ({
  mockApi: {
    getReviews: vi.fn(),
    approveReview: vi.fn(),
    deleteReview: vi.fn(),
    getReviewLogs: vi.fn(),
    batchApproveReviews: vi.fn(),
  },
  mockMessage: {
    success: vi.fn(),
    error: vi.fn(),
    warning: vi.fn(),
    info: vi.fn(),
    loading: vi.fn(),
  },
}));

vi.mock('@/api/review', () => ({
  reviewApi: mockApi,
}));

// Mock antd message
vi.mock('antd', async () => {
  const actual = await vi.importActual('antd');
  return {
    ...actual,
    message: mockMessage,
  };
});

// Mock usePermissions hook
vi.mock('@/hooks/usePermission', () => ({
  usePermissions: () => ({
    canApprove: true,
    canReject: true,
    canDelete: true,
    canViewLogs: true,
  }),
}));

// Helper function to create mock review
const createMockReview = (overrides: Record<string, unknown> = {}): Record<string, unknown> => ({
  id: 1,
  orderId: 1001,
  reviewerId: 100,
  reviewerName: '测试用户',
  playerId: 200,
  playerName: '测试陪玩师',
  rating: 5,
  comment: '服务很好，非常满意！',
  images: [],
  status: 'pending',
  isReported: false,
  createdAt: '2024-01-01T00:00:00Z',
  updatedAt: '2024-01-01T00:00:00Z',
  ...overrides,
});

describe('ReviewList', () => {
  beforeEach(() => {
    resetAllMocks();
    mockMessage.success.mockClear();
    mockMessage.error.mockClear();
    localStorage.setItem('token', 'test-token');
    localStorage.setItem('user_role', 'admin');

    mockApi.getReviews.mockResolvedValue({
      data: {
        success: true,
        data: [createMockReview()],
        pagination: { total: 1 },
      },
    });
    mockApi.getReviewLogs.mockResolvedValue({
      data: { success: true, data: [] },
    });
  });

  afterEach(() => {
    vi.clearAllMocks();
  });

  describe('Successful Data Loading', () => {
    it('should render review page successfully', async () => {
      renderWithProviders(<ReviewList />);
      await waitFor(() => {
        expect(screen.getByText('评价列表')).toBeInTheDocument();
      });
      expect(mockApi.getReviews).toHaveBeenCalled();
    });

    it('should display reviewer name', async () => {
      renderWithProviders(<ReviewList />);
      await waitFor(() => {
        expect(screen.getByText('测试用户')).toBeInTheDocument();
      });
    });

    it('should display player name', async () => {
      renderWithProviders(<ReviewList />);
      await waitFor(() => {
        expect(screen.getByText('测试陪玩师')).toBeInTheDocument();
      });
    });
  });

  describe('Loading States', () => {
    it('should call API on mount', async () => {
      renderWithProviders(<ReviewList />);
      await waitFor(() => {
        expect(mockApi.getReviews).toHaveBeenCalledTimes(1);
      });
    });
  });

  describe('Error Handling', () => {
    it('should display error message when API fails', async () => {
      mockApi.getReviews.mockRejectedValue(new Error('Network error'));
      renderWithProviders(<ReviewList />);
      await waitFor(() => {
        expect(mockMessage.error).toHaveBeenCalledWith('获取评价列表失败');
      });
    });
  });

  describe('Filter Functionality', () => {
    it('should have search input', async () => {
      renderWithProviders(<ReviewList />);
      await waitFor(() => {
        expect(screen.getByPlaceholderText('搜索订单ID/评价内容')).toBeInTheDocument();
      });
    });

    it('should have search button', async () => {
      renderWithProviders(<ReviewList />);
      await waitFor(() => {
        expect(screen.getByText('搜索')).toBeInTheDocument();
      });
    });

    it('should have reset button', async () => {
      renderWithProviders(<ReviewList />);
      await waitFor(() => {
        expect(screen.getByText('重置')).toBeInTheDocument();
      });
    });
  });

  describe('Review Status Display', () => {
    it('should display pending status tag', async () => {
      renderWithProviders(<ReviewList />);
      await waitFor(() => {
        expect(screen.getByText('待审核')).toBeInTheDocument();
      });
    });

    it('should display approved status tag', async () => {
      mockApi.getReviews.mockResolvedValue({
        data: {
          success: true,
          data: [createMockReview({ status: 'approved' })],
          pagination: { total: 1 },
        },
      });
      renderWithProviders(<ReviewList />);
      await waitFor(() => {
        expect(screen.getByText('已通过')).toBeInTheDocument();
      });
    });
  });

  describe('Review Actions', () => {
    it('should display detail button', async () => {
      renderWithProviders(<ReviewList />);
      await waitFor(() => {
        expect(screen.getByText('详情')).toBeInTheDocument();
      });
    });

    it('should display approve button', async () => {
      renderWithProviders(<ReviewList />);
      await waitFor(() => {
        expect(screen.getByText('批准')).toBeInTheDocument();
      });
    });

    it('should display delete button', async () => {
      renderWithProviders(<ReviewList />);
      await waitFor(() => {
        expect(screen.getByText('删除')).toBeInTheDocument();
      });
    });
  });

  describe('View Detail', () => {
    it('should open detail drawer when clicking detail button', async () => {
      renderWithProviders(<ReviewList />);
      await waitFor(() => {
        expect(screen.getByText('详情')).toBeInTheDocument();
      });
      const detailButton = screen.getByText('详情');
      fireEvent.click(detailButton);
      await waitFor(() => {
        expect(screen.getByText('评价详情')).toBeInTheDocument();
      });
    });
  });

  describe('Accessibility', () => {
    it('should have proper page title', async () => {
      renderWithProviders(<ReviewList />);
      await waitFor(() => {
        expect(screen.getByText('评价列表')).toBeInTheDocument();
      });
    });

    it('should have proper table structure', async () => {
      renderWithProviders(<ReviewList />);
      await waitFor(() => {
        expect(screen.getByRole('table')).toBeInTheDocument();
      });
    });
  });

  describe('Table Structure', () => {
    it('should render table with data', async () => {
      renderWithProviders(<ReviewList />);
      await waitFor(() => {
        expect(screen.getByRole('table')).toBeInTheDocument();
      });
      // Check that data is rendered
      expect(screen.getByText('测试用户')).toBeInTheDocument();
    });

    it('should display review comment', async () => {
      renderWithProviders(<ReviewList />);
      await waitFor(() => {
        expect(screen.getByText('服务很好，非常满意！')).toBeInTheDocument();
      });
    });
  });

  describe('Empty State', () => {
    it('should handle empty review list', async () => {
      mockApi.getReviews.mockResolvedValue({
        data: {
          success: true,
          data: [],
          pagination: { total: 0 },
        },
      });
      renderWithProviders(<ReviewList />);
      await waitFor(() => {
        expect(screen.getByText('评价列表')).toBeInTheDocument();
      });
      expect(mockApi.getReviews).toHaveBeenCalled();
    });
  });
});
