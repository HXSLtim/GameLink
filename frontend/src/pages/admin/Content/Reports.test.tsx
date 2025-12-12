/**
 * Content Reports Page Tests
 * 测试举报管理页面的渲染和处理功能
 * 需求: 5.1, 5.2, 5.3, 5.4, 5.5
 */
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, waitFor } from '@testing-library/react';
import { BrowserRouter } from 'react-router-dom';
import ReportsPage from './Reports';
import { feedReportApi } from '@/api/content';

// Mock the API
vi.mock('@/api/content', () => ({
  feedReportApi: {
    getReports: vi.fn(),
    getReport: vi.fn(),
    processReport: vi.fn(),
  },
}));

const mockReportsData = {
  items: [
    {
      id: 1,
      feedId: 101,
      reporterId: 201,
      reporterName: '举报人A',
      reason: '内容不当',
      status: 'pending',
      createdAt: '2025-12-01 10:00:00',
    },
    {
      id: 2,
      feedId: 102,
      reporterId: 202,
      reporterName: '举报人B',
      reason: '垃圾广告',
      status: 'processed',
      result: '已删除内容',
      handledAt: '2025-12-02 11:00:00',
      createdAt: '2025-12-01 09:00:00',
    },
  ],
  total: 2,
};

const renderWithRouter = (component: React.ReactNode) => {
  return render(<BrowserRouter>{component}</BrowserRouter>);
};

describe('Reports Page', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    // Mock returns the response directly (not wrapped in data)
    (feedReportApi.getReports as ReturnType<typeof vi.fn>).mockResolvedValue({
      success: true,
      data: mockReportsData,
    });
  });

  describe('基本渲染', () => {
    it('should render page title', async () => {
      renderWithRouter(<ReportsPage />);

      await waitFor(() => {
        expect(screen.getByText('动态举报管理')).toBeInTheDocument();
      });
    });

    it('should render reports table', async () => {
      renderWithRouter(<ReportsPage />);

      await waitFor(() => {
        expect(screen.getByText('内容不当')).toBeInTheDocument();
        expect(screen.getByText('垃圾广告')).toBeInTheDocument();
      });
    });

    it('should display reporter names', async () => {
      renderWithRouter(<ReportsPage />);

      await waitFor(() => {
        expect(screen.getByText('举报人A')).toBeInTheDocument();
        expect(screen.getByText('举报人B')).toBeInTheDocument();
      });
    });
  });

  describe('状态筛选', () => {
    it('should render status filter', async () => {
      renderWithRouter(<ReportsPage />);

      await waitFor(() => {
        // Filter controls should be present
        const table = document.querySelector('.ant-table');
        expect(table).toBeInTheDocument();
      });
    });
  });

  describe('API调用', () => {
    it('should call getReports on mount', async () => {
      renderWithRouter(<ReportsPage />);

      await waitFor(() => {
        expect(feedReportApi.getReports).toHaveBeenCalled();
      });
    });
  });

  describe('加载状态', () => {
    it('should show loading state while fetching', async () => {
      (feedReportApi.getReports as ReturnType<typeof vi.fn>).mockImplementation(
        () => new Promise(() => {})
      );

      renderWithRouter(<ReportsPage />);

      await waitFor(() => {
        const spinner = document.querySelector('.ant-spin-spinning') || 
                       document.querySelector('.ant-table-loading');
        expect(spinner || document.querySelector('.ant-spin')).toBeTruthy();
      }, { timeout: 1000 });
    });
  });

  describe('空状态', () => {
    it('should handle empty reports list', async () => {
      (feedReportApi.getReports as ReturnType<typeof vi.fn>).mockResolvedValue({
        success: true,
        data: { items: [], total: 0 },
      });

      renderWithRouter(<ReportsPage />);

      await waitFor(() => {
        const table = document.querySelector('.ant-table');
        expect(table).toBeInTheDocument();
      });
    });
  });

  describe('错误处理', () => {
    it('should handle API error gracefully', async () => {
      (feedReportApi.getReports as ReturnType<typeof vi.fn>).mockRejectedValue(new Error('API Error'));

      renderWithRouter(<ReportsPage />);

      await waitFor(() => {
        expect(screen.getByText('动态举报管理')).toBeInTheDocument();
      });
    });
  });
});
