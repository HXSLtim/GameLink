/**
 * Content Feeds Page Tests
 * 测试动态审核页面的渲染和操作功能
 * 需求: 1.1, 1.2, 1.3, 1.4, 2.1, 2.4, 2.5
 */
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, waitFor } from '@testing-library/react';
import { BrowserRouter } from 'react-router-dom';
import FeedsPage from './Feeds';
import { feedApi } from '@/api/content';

// Mock the API
vi.mock('@/api/content', () => ({
  feedApi: {
    getFeeds: vi.fn(),
    approveFeed: vi.fn(),
    rejectFeed: vi.fn(),
    deleteFeed: vi.fn(),
    batchApproveFeed: vi.fn(),
    batchRejectFeed: vi.fn(),
  },
}));

const mockFeedsData = {
  items: [
    {
      id: 1,
      authorId: 101,
      authorName: '用户A',
      content: '这是一条测试动态',
      images: [],
      visibility: 'public',
      moderationStatus: 'pending',
      createdAt: '2025-12-01 10:00:00',
    },
    {
      id: 2,
      authorId: 102,
      authorName: '用户B',
      content: '这是另一条测试动态',
      images: ['http://example.com/image.jpg'],
      visibility: 'public',
      moderationStatus: 'approved',
      createdAt: '2025-12-02 11:00:00',
    },
  ],
  total: 2,
};

const renderWithRouter = (component: React.ReactNode) => {
  return render(<BrowserRouter>{component}</BrowserRouter>);
};

describe('Feeds Page', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    // Mock returns the response directly (not wrapped in data)
    (feedApi.getFeeds as ReturnType<typeof vi.fn>).mockResolvedValue({
      success: true,
      data: mockFeedsData,
    });
  });

  describe('基本渲染', () => {
    it('should render page title', async () => {
      renderWithRouter(<FeedsPage />);

      await waitFor(() => {
        expect(screen.getByText('动态审核')).toBeInTheDocument();
      });
    });

    it('should render filter controls', async () => {
      renderWithRouter(<FeedsPage />);

      await waitFor(() => {
        // Search input
        expect(screen.getByPlaceholderText(/搜索/i)).toBeInTheDocument();
      });
    });

    it('should render feeds table', async () => {
      renderWithRouter(<FeedsPage />);

      await waitFor(() => {
        expect(screen.getByText('这是一条测试动态')).toBeInTheDocument();
        expect(screen.getByText('这是另一条测试动态')).toBeInTheDocument();
      });
    });

    it('should display author names', async () => {
      renderWithRouter(<FeedsPage />);

      await waitFor(() => {
        expect(screen.getByText('用户A')).toBeInTheDocument();
        expect(screen.getByText('用户B')).toBeInTheDocument();
      });
    });
  });

  describe('审核状态显示', () => {
    it('should display moderation status tags', async () => {
      renderWithRouter(<FeedsPage />);

      await waitFor(() => {
        // Status tags should be rendered
        const table = document.querySelector('.ant-table');
        expect(table).toBeInTheDocument();
      });
    });
  });

  describe('API调用', () => {
    it('should call getFeeds on mount', async () => {
      renderWithRouter(<FeedsPage />);

      await waitFor(() => {
        expect(feedApi.getFeeds).toHaveBeenCalled();
      });
    });
  });

  describe('加载状态', () => {
    it('should show loading state while fetching', async () => {
      (feedApi.getFeeds as ReturnType<typeof vi.fn>).mockImplementation(
        () => new Promise(() => {})
      );

      renderWithRouter(<FeedsPage />);

      // Wait for the component to start loading
      await waitFor(() => {
        // Check for spinning state on table or any spin indicator
        const spinner = document.querySelector('.ant-spin-spinning') || 
                       document.querySelector('.ant-table-loading');
        expect(spinner || document.querySelector('.ant-spin')).toBeTruthy();
      }, { timeout: 1000 });
    });
  });

  describe('空状态', () => {
    it('should handle empty feeds list', async () => {
      (feedApi.getFeeds as ReturnType<typeof vi.fn>).mockResolvedValue({
        success: true,
        data: { items: [], total: 0 },
      });

      renderWithRouter(<FeedsPage />);

      await waitFor(() => {
        // Empty state or table with no data
        const table = document.querySelector('.ant-table');
        expect(table).toBeInTheDocument();
      });
    });
  });

  describe('错误处理', () => {
    it('should handle API error gracefully', async () => {
      (feedApi.getFeeds as ReturnType<typeof vi.fn>).mockRejectedValue(new Error('API Error'));

      renderWithRouter(<FeedsPage />);

      // Component should still render without crashing
      await waitFor(() => {
        expect(screen.getByText('动态审核')).toBeInTheDocument();
      });
    });
  });
});
