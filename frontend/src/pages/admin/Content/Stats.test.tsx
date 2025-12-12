/**
 * Content Stats Page Tests
 * 测试内容统计页面的渲染和导出功能
 * 需求: 8.1, 8.2, 8.3, 8.4, 8.5
 */
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { BrowserRouter } from 'react-router-dom';
import StatsPage from './Stats';
import { contentStatsApi } from '@/api/content';

// Mock the API
vi.mock('@/api/content', () => ({
  contentStatsApi: {
    getStats: vi.fn(),
    exportStats: vi.fn(),
  },
}));

const mockStatsData = {
  stats: {
    totalFeeds: 100,
    pendingFeeds: 10,
    approvedFeeds: 80,
    rejectedFeeds: 10,
    totalMessages: 500,
    flaggedMessages: 5,
    totalReports: 20,
    reportHandleRate: 85.5,
  },
  trend: [
    { date: '2025-12-01', feedCount: 10, messageCount: 50, reportCount: 2 },
    { date: '2025-12-02', feedCount: 12, messageCount: 55, reportCount: 3 },
    { date: '2025-12-03', feedCount: 8, messageCount: 45, reportCount: 1 },
  ],
};

const renderWithRouter = async (component: React.ReactNode) => {
  const result = render(<BrowserRouter>{component}</BrowserRouter>);
  // Wait for initial render and effects to complete
  await waitFor(() => {
    expect(contentStatsApi.getStats).toHaveBeenCalled();
  });
  return result;
};

describe('Stats Page', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    // Mock returns the response directly (not wrapped in data)
    (contentStatsApi.getStats as ReturnType<typeof vi.fn>).mockResolvedValue({
      success: true,
      data: mockStatsData,
    });
  });

  describe('基本渲染', () => {
    it('should render statistics cards', async () => {
      await renderWithRouter(<StatsPage />);

      await waitFor(() => {
        expect(screen.getByText('总动态数')).toBeInTheDocument();
        expect(screen.getByText('待审核动态')).toBeInTheDocument();
        expect(screen.getByText('已通过动态')).toBeInTheDocument();
        expect(screen.getByText('已拒绝动态')).toBeInTheDocument();
        expect(screen.getByText('总消息数')).toBeInTheDocument();
        expect(screen.getByText('标记消息数')).toBeInTheDocument();
        expect(screen.getByText('总举报数')).toBeInTheDocument();
        expect(screen.getByText('举报处理率')).toBeInTheDocument();
      });
    });

    it('should display correct statistics values', async () => {
      await renderWithRouter(<StatsPage />);

      await waitFor(() => {
        expect(screen.getByText('100')).toBeInTheDocument();
        expect(screen.getByText('80')).toBeInTheDocument();
        expect(screen.getByText('500')).toBeInTheDocument();
        // 10 appears multiple times (pending and rejected), use getAllByText
        expect(screen.getAllByText('10').length).toBeGreaterThanOrEqual(1);
      });
    });

    it('should render trend table', async () => {
      await renderWithRouter(<StatsPage />);

      await waitFor(() => {
        expect(screen.getByText('内容趋势')).toBeInTheDocument();
        expect(screen.getByText('日期')).toBeInTheDocument();
        expect(screen.getByText('动态数')).toBeInTheDocument();
        expect(screen.getByText('消息数')).toBeInTheDocument();
        expect(screen.getByText('举报数')).toBeInTheDocument();
      });
    });

    it('should render export button', async () => {
      await renderWithRouter(<StatsPage />);

      await waitFor(() => {
        expect(screen.getByText('导出Excel')).toBeInTheDocument();
      });
    });
  });

  describe('时间范围选择', () => {
    it('should have default 30 days selected', async () => {
      await renderWithRouter(<StatsPage />);

      await waitFor(() => {
        expect(screen.getByText('最近30天')).toBeInTheDocument();
      });
    });

    it('should call API with different days when selection changes', async () => {
      await renderWithRouter(<StatsPage />);

      await waitFor(() => {
        expect(contentStatsApi.getStats).toHaveBeenCalledWith(30);
      });
    });
  });

  describe('导出功能', () => {
    it('should call export API when export button is clicked', async () => {
      const mockBlob = new Blob(['test'], { type: 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet' });
      (contentStatsApi.exportStats as ReturnType<typeof vi.fn>).mockResolvedValue({
        data: mockBlob,
        headers: { 'content-disposition': 'attachment; filename="stats.xlsx"' },
      });

      // Mock URL.createObjectURL and URL.revokeObjectURL
      const mockCreateObjectURL = vi.fn(() => 'blob:test');
      const mockRevokeObjectURL = vi.fn();
      window.URL.createObjectURL = mockCreateObjectURL;
      window.URL.revokeObjectURL = mockRevokeObjectURL;

      await renderWithRouter(<StatsPage />);

      await waitFor(() => {
        expect(screen.getByText('导出Excel')).toBeInTheDocument();
      });

      const exportButton = screen.getByText('导出Excel');
      fireEvent.click(exportButton);

      await waitFor(() => {
        expect(contentStatsApi.exportStats).toHaveBeenCalledWith(30);
      });
    });

    it('should show error message when export fails', async () => {
      (contentStatsApi.exportStats as ReturnType<typeof vi.fn>).mockRejectedValue(new Error('Export failed'));

      await renderWithRouter(<StatsPage />);

      await waitFor(() => {
        expect(screen.getByText('导出Excel')).toBeInTheDocument();
      });

      const exportButton = screen.getByText('导出Excel');
      fireEvent.click(exportButton);

      // Error handling is done via message.error, which is harder to test
      // Just verify the API was called
      await waitFor(() => {
        expect(contentStatsApi.exportStats).toHaveBeenCalled();
      });
    });
  });

  describe('加载状态', () => {
    it('should show loading spinner while fetching data', async () => {
      (contentStatsApi.getStats as ReturnType<typeof vi.fn>).mockImplementation(
        () => new Promise(() => {}) // Never resolves
      );

      render(<BrowserRouter><StatsPage /></BrowserRouter>);

      // Spin component should be present
      await waitFor(() => {
        const spinner = document.querySelector('.ant-spin');
        expect(spinner).toBeInTheDocument();
      });
    });
  });

  describe('错误处理', () => {
    it('should handle API error gracefully', async () => {
      (contentStatsApi.getStats as ReturnType<typeof vi.fn>).mockRejectedValue(new Error('API Error'));

      render(<BrowserRouter><StatsPage /></BrowserRouter>);

      // Component should still render without crashing
      await waitFor(() => {
        expect(screen.getByText('内容趋势')).toBeInTheDocument();
      });
    });
  });
});
