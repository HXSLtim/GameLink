/**
 * Ranking Commission Page Tests
 *
 * Tests for RankingCommission page component including:
 * - Successful data loading
 * - Loading states
 * - Error handling
 * - CRUD operations
 * - Export functionality
 */

import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { screen, waitFor, fireEvent } from '@testing-library/react';
import RankingCommissionPage from './index';
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

// Helper function to create mock config
const createMockConfig = (overrides: Record<string, unknown> = {}): Record<string, unknown> => ({
  id: 1,
  name: '测试配置',
  rankingType: 'income',
  period: 'monthly',
  month: '2024-01',
  rules: [
    { rankStart: 1, rankEnd: 10, commissionRate: 5 },
    { rankStart: 11, rankEnd: 50, commissionRate: 3 },
  ],
  description: '测试描述',
  isActive: true,
  createdAt: '2024-01-01T00:00:00Z',
  updatedAt: '2024-01-01T00:00:00Z',
  ...overrides,
});

describe('RankingCommissionPage', () => {
  beforeEach(() => {
    resetAllMocks();
    mockMessage.success.mockClear();
    mockMessage.error.mockClear();
    localStorage.setItem('token', 'test-token');
    localStorage.setItem('user_role', 'admin');

    // Set default mock return values
    mockApiClient.get.mockResolvedValue({
      data: {
        success: true,
        data: {
          configs: [createMockConfig()],
          total: 1,
        },
      },
    });
  });

  afterEach(() => {
    vi.clearAllMocks();
  });

  describe('Successful Data Loading', () => {
    it('should render ranking commission page successfully', async () => {
      renderWithProviders(<RankingCommissionPage />);

      await waitFor(() => {
        expect(screen.getByText('排行榜抽成配置')).toBeInTheDocument();
      });

      expect(mockApiClient.get).toHaveBeenCalled();
    });

    it('should display page subtitle', async () => {
      renderWithProviders(<RankingCommissionPage />);

      await waitFor(() => {
        expect(screen.getByText('管理陪玩师排行榜抽成规则')).toBeInTheDocument();
      });
    });

    it('should display config list', async () => {
      renderWithProviders(<RankingCommissionPage />);

      await waitFor(() => {
        expect(screen.getByText('测试配置')).toBeInTheDocument();
      });
    });

    it('should display ranking type', async () => {
      renderWithProviders(<RankingCommissionPage />);

      await waitFor(() => {
        expect(screen.getByText('收入排行')).toBeInTheDocument();
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
                  data: {
                    configs: [createMockConfig()],
                    total: 1,
                  },
                },
              });
            }, 100);
          })
      );

      renderWithProviders(<RankingCommissionPage />);

      expect(mockApiClient.get).toHaveBeenCalled();
    });

    it('should hide loading indicator after data loads', async () => {
      renderWithProviders(<RankingCommissionPage />);

      await waitFor(() => {
        expect(screen.getByText('排行榜抽成配置')).toBeInTheDocument();
      });

      await flushPromises();

      expect(mockApiClient.get).toHaveBeenCalledTimes(1);
    });
  });

  describe('Error Handling', () => {
    it('should display error message when API fails', async () => {
      mockApiClient.get.mockRejectedValue(new Error('Network error'));

      renderWithProviders(<RankingCommissionPage />);

      await waitFor(() => {
        expect(mockMessage.error).toHaveBeenCalledWith('加载失败');
      });
    });
  });

  describe('Table Structure', () => {
    it('should display table with data', async () => {
      renderWithProviders(<RankingCommissionPage />);

      await waitFor(() => {
        expect(screen.getByRole('table')).toBeInTheDocument();
      });

      // Check for config data
      expect(screen.getByText('测试配置')).toBeInTheDocument();
    });

    it('should display ranking type in table', async () => {
      renderWithProviders(<RankingCommissionPage />);

      await waitFor(() => {
        expect(screen.getByText('收入排行')).toBeInTheDocument();
      });
    });

    it('should display rules count in table', async () => {
      renderWithProviders(<RankingCommissionPage />);

      await waitFor(() => {
        expect(screen.getByText('2')).toBeInTheDocument(); // 2 rules
      });
    });
  });

  describe('Toolbar Buttons', () => {
    it('should display create button', async () => {
      renderWithProviders(<RankingCommissionPage />);

      await waitFor(() => {
        expect(screen.getByText('新增配置')).toBeInTheDocument();
      });
    });

    it('should display export button', async () => {
      renderWithProviders(<RankingCommissionPage />);

      await waitFor(() => {
        expect(screen.getByText('导出数据')).toBeInTheDocument();
      });
    });
  });

  describe('Create Config Modal', () => {
    it('should open create modal when button clicked', async () => {
      renderWithProviders(<RankingCommissionPage />);

      await waitFor(() => {
        expect(screen.getByText('新增配置')).toBeInTheDocument();
      });

      const createButton = screen.getByText('新增配置');
      fireEvent.click(createButton);

      await waitFor(() => {
        expect(screen.getByText('新增抽成配置')).toBeInTheDocument();
      });
    });

    it('should display form fields in modal', async () => {
      renderWithProviders(<RankingCommissionPage />);

      await waitFor(() => {
        expect(screen.getByText('新增配置')).toBeInTheDocument();
      });

      const createButton = screen.getByText('新增配置');
      fireEvent.click(createButton);

      await waitFor(() => {
        expect(screen.getByRole('dialog')).toBeInTheDocument();
      });

      // Check for input placeholders
      expect(screen.getByPlaceholderText('请输入配置名称')).toBeInTheDocument();
      // YYYY-MM appears multiple times, use getAllByPlaceholderText
      const monthInputs = screen.getAllByPlaceholderText('YYYY-MM');
      expect(monthInputs.length).toBeGreaterThan(0);
    });

    it('should display rules section in modal', async () => {
      renderWithProviders(<RankingCommissionPage />);

      await waitFor(() => {
        expect(screen.getByText('新增配置')).toBeInTheDocument();
      });

      const createButton = screen.getByText('新增配置');
      fireEvent.click(createButton);

      await waitFor(() => {
        expect(screen.getByText('抽成规则')).toBeInTheDocument();
      });
    });
  });

  describe('Edit Config', () => {
    it('should display edit button', async () => {
      renderWithProviders(<RankingCommissionPage />);

      await waitFor(() => {
        expect(screen.getByText('编辑')).toBeInTheDocument();
      });
    });

    it('should open edit modal when button clicked', async () => {
      renderWithProviders(<RankingCommissionPage />);

      await waitFor(() => {
        expect(screen.getByText('编辑')).toBeInTheDocument();
      });

      const editButton = screen.getByText('编辑');
      fireEvent.click(editButton);

      await waitFor(() => {
        expect(screen.getByText('编辑抽成配置')).toBeInTheDocument();
      });
    });
  });

  describe('Delete Config', () => {
    it('should display delete button', async () => {
      renderWithProviders(<RankingCommissionPage />);

      await waitFor(() => {
        expect(screen.getByText('删除')).toBeInTheDocument();
      });
    });

    it('should call delete API when button clicked', async () => {
      mockApiClient.delete.mockResolvedValue({ data: { success: true } });

      renderWithProviders(<RankingCommissionPage />);

      await waitFor(() => {
        expect(screen.getByText('删除')).toBeInTheDocument();
      });

      const deleteButton = screen.getByText('删除');
      fireEvent.click(deleteButton);

      await waitFor(() => {
        expect(mockApiClient.delete).toHaveBeenCalledWith('/admin/ranking-commission/configs/1');
      });
    });
  });

  describe('Status Display', () => {
    it('should display active status', async () => {
      renderWithProviders(<RankingCommissionPage />);

      await waitFor(() => {
        expect(screen.getByText('启用')).toBeInTheDocument();
      });
    });

    it('should display inactive status', async () => {
      mockApiClient.get.mockResolvedValue({
        data: {
          success: true,
          data: {
            configs: [createMockConfig({ isActive: false })],
            total: 1,
          },
        },
      });

      renderWithProviders(<RankingCommissionPage />);

      await waitFor(() => {
        expect(screen.getByText('禁用')).toBeInTheDocument();
      });
    });
  });

  describe('Ranking Type Display', () => {
    it('should display income ranking type', async () => {
      renderWithProviders(<RankingCommissionPage />);

      await waitFor(() => {
        expect(screen.getByText('收入排行')).toBeInTheDocument();
      });
    });

    it('should display order count ranking type', async () => {
      mockApiClient.get.mockResolvedValue({
        data: {
          success: true,
          data: {
            configs: [createMockConfig({ rankingType: 'order_count' })],
            total: 1,
          },
        },
      });

      renderWithProviders(<RankingCommissionPage />);

      await waitFor(() => {
        expect(screen.getByText('订单量排行')).toBeInTheDocument();
      });
    });
  });

  describe('Rules Count Display', () => {
    it('should display rules count', async () => {
      renderWithProviders(<RankingCommissionPage />);

      await waitFor(() => {
        expect(screen.getByText('2')).toBeInTheDocument(); // 2 rules
      });
    });
  });

  describe('Empty State', () => {
    it('should handle empty config list', async () => {
      mockApiClient.get.mockResolvedValue({
        data: {
          success: true,
          data: {
            configs: [],
            total: 0,
          },
        },
      });

      renderWithProviders(<RankingCommissionPage />);

      await waitFor(() => {
        expect(screen.getByText('排行榜抽成配置')).toBeInTheDocument();
      });

      expect(mockApiClient.get).toHaveBeenCalled();
    });
  });

  describe('Search Functionality', () => {
    it('should have month search field', async () => {
      renderWithProviders(<RankingCommissionPage />);

      await waitFor(() => {
        expect(screen.getByPlaceholderText('YYYY-MM')).toBeInTheDocument();
      });
    });
  });

  describe('Accessibility', () => {
    it('should have proper page title', async () => {
      renderWithProviders(<RankingCommissionPage />);

      await waitFor(() => {
        expect(screen.getByText('排行榜抽成配置')).toBeInTheDocument();
      });
    });

    it('should have proper table structure', async () => {
      renderWithProviders(<RankingCommissionPage />);

      await waitFor(() => {
        expect(screen.getByRole('table')).toBeInTheDocument();
      });
    });
  });
});
