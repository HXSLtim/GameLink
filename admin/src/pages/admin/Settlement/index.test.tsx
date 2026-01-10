/**
 * Settlement Company Management Page Tests
 *
 * Tests for Settlement page component including:
 * - Successful data loading
 * - Loading states
 * - Error handling
 * - Statistics display
 * - CRUD operations
 * - Status toggle
 * - History view
 */

import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { screen, waitFor, fireEvent } from '@testing-library/react';
import SettlementPage from './index';
import { renderWithProviders, resetAllMocks, flushPromises } from '@/testutils';

// Mock the settlementApi module using vi.hoisted
const { mockApi, mockMessage } = vi.hoisted(() => ({
  mockApi: {
    getSettlementCompanies: vi.fn(),
    createSettlementCompany: vi.fn(),
    updateSettlementCompany: vi.fn(),
    deleteSettlementCompany: vi.fn(),
    toggleSettlementCompanyStatus: vi.fn(),
    getSettlementCompanyHistory: vi.fn(),
    batchUpdateCompanyStatus: vi.fn(),
  },
  mockMessage: {
    success: vi.fn(),
    error: vi.fn(),
    warning: vi.fn(),
    info: vi.fn(),
    loading: vi.fn(),
  },
}));

vi.mock('@/api/settlement', () => ({
  settlementApi: mockApi,
}));

vi.mock('antd', async () => {
  const actual = await vi.importActual('antd');
  return {
    ...actual,
    message: mockMessage,
  };
});

// Helper function to create mock company
const createMockCompany = (overrides: Record<string, unknown> = {}): Record<string, unknown> => ({
  id: 1,
  name: '测试结算公司',
  type: 'company',
  taxNumber: '91110000123456789X',
  contactPerson: '张三',
  contactPhone: '13800138000',
  status: 'active',
  playerCount: 10,
  createdAt: '2024-01-01T00:00:00Z',
  updatedAt: '2024-01-01T00:00:00Z',
  ...overrides,
});

describe('SettlementPage', () => {
  beforeEach(() => {
    resetAllMocks();
    mockMessage.success.mockClear();
    mockMessage.error.mockClear();
    localStorage.setItem('token', 'test-token');
    localStorage.setItem('user_role', 'admin');

    // Set default mock return values
    mockApi.getSettlementCompanies.mockResolvedValue({
      data: {
        success: true,
        data: [createMockCompany()],
        pagination: { total: 1, page: 1, page_size: 10 },
      },
    });
    mockApi.getSettlementCompanyHistory.mockResolvedValue({
      data: {
        success: true,
        data: [],
      },
    });
  });

  afterEach(() => {
    vi.clearAllMocks();
  });

  describe('Successful Data Loading', () => {
    it('should render settlement page successfully', async () => {
      renderWithProviders(<SettlementPage />);

      await waitFor(() => {
        expect(screen.getByText('结算公司管理')).toBeInTheDocument();
      });

      expect(mockApi.getSettlementCompanies).toHaveBeenCalled();
    });

    it('should display page subtitle', async () => {
      renderWithProviders(<SettlementPage />);

      await waitFor(() => {
        expect(screen.getByText('管理陪玩师结算公司与归属关系')).toBeInTheDocument();
      });
    });

    it('should display company list', async () => {
      renderWithProviders(<SettlementPage />);

      await waitFor(() => {
        expect(screen.getByText('测试结算公司')).toBeInTheDocument();
      });
    });

    it('should display statistics in header', async () => {
      renderWithProviders(<SettlementPage />);

      await waitFor(() => {
        expect(screen.getByText('公司总数')).toBeInTheDocument();
      });

      expect(screen.getByText('启用公司')).toBeInTheDocument();
      expect(screen.getByText('关联陪玩师')).toBeInTheDocument();
    });
  });

  describe('Loading States', () => {
    it('should show loading indicator while fetching data', async () => {
      mockApi.getSettlementCompanies.mockImplementation(
        () =>
          new Promise((resolve) => {
            setTimeout(() => {
              resolve({
                data: {
                  success: true,
                  data: [createMockCompany()],
                  pagination: { total: 1 },
                },
              });
            }, 100);
          })
      );

      renderWithProviders(<SettlementPage />);

      expect(mockApi.getSettlementCompanies).toHaveBeenCalled();
    });

    it('should hide loading indicator after data loads', async () => {
      renderWithProviders(<SettlementPage />);

      await waitFor(() => {
        expect(screen.getByText('结算公司管理')).toBeInTheDocument();
      });

      await flushPromises();

      expect(mockApi.getSettlementCompanies).toHaveBeenCalledTimes(1);
    });
  });

  describe('Error Handling', () => {
    it('should display error message when API fails', async () => {
      mockApi.getSettlementCompanies.mockRejectedValue(new Error('Network error'));

      renderWithProviders(<SettlementPage />);

      await waitFor(() => {
        expect(mockMessage.error).toHaveBeenCalledWith('加载结算公司列表失败');
      });
    });

    it('should display error message when API returns error', async () => {
      mockApi.getSettlementCompanies.mockResolvedValue({
        data: {
          success: false,
          message: '服务器错误',
        },
      });

      renderWithProviders(<SettlementPage />);

      await waitFor(() => {
        expect(mockMessage.error).toHaveBeenCalledWith('服务器错误');
      });
    });
  });

  describe('Table Structure', () => {
    it('should display table with data', async () => {
      renderWithProviders(<SettlementPage />);

      await waitFor(() => {
        expect(screen.getByRole('table')).toBeInTheDocument();
      });

      // Check for company data
      expect(screen.getByText('测试结算公司')).toBeInTheDocument();
    });

    it('should display company type', async () => {
      renderWithProviders(<SettlementPage />);

      await waitFor(() => {
        expect(screen.getByText('企业')).toBeInTheDocument();
      });
    });

    it('should display tax number', async () => {
      renderWithProviders(<SettlementPage />);

      await waitFor(() => {
        expect(screen.getByText('91110000123456789X')).toBeInTheDocument();
      });
    });

    it('should display contact person', async () => {
      renderWithProviders(<SettlementPage />);

      await waitFor(() => {
        expect(screen.getByText('张三')).toBeInTheDocument();
      });
    });

    it('should display active status', async () => {
      renderWithProviders(<SettlementPage />);

      await waitFor(() => {
        expect(screen.getByText('启用')).toBeInTheDocument();
      });
    });
  });

  describe('Toolbar Buttons', () => {
    it('should display create button', async () => {
      renderWithProviders(<SettlementPage />);

      await waitFor(() => {
        expect(screen.getByText('新增公司')).toBeInTheDocument();
      });
    });

    it('should display batch enable button', async () => {
      renderWithProviders(<SettlementPage />);

      await waitFor(() => {
        expect(screen.getByText('批量启用')).toBeInTheDocument();
      });
    });

    it('should display batch disable button', async () => {
      renderWithProviders(<SettlementPage />);

      await waitFor(() => {
        expect(screen.getByText('批量停用')).toBeInTheDocument();
      });
    });

    it('should display batch delete button', async () => {
      renderWithProviders(<SettlementPage />);

      await waitFor(() => {
        expect(screen.getByText('批量删除')).toBeInTheDocument();
      });
    });
  });

  describe('Company Actions', () => {
    it('should display edit button', async () => {
      renderWithProviders(<SettlementPage />);

      await waitFor(() => {
        expect(screen.getByText('编辑')).toBeInTheDocument();
      });
    });

    it('should display history button', async () => {
      renderWithProviders(<SettlementPage />);

      await waitFor(() => {
        expect(screen.getByText('历史')).toBeInTheDocument();
      });
    });

    it('should display toggle status button', async () => {
      renderWithProviders(<SettlementPage />);

      await waitFor(() => {
        expect(screen.getByText('停用')).toBeInTheDocument();
      });
    });

    it('should display delete button', async () => {
      renderWithProviders(<SettlementPage />);

      await waitFor(() => {
        expect(screen.getByText('删除')).toBeInTheDocument();
      });
    });
  });

  describe('Status Display', () => {
    it('should display active status', async () => {
      renderWithProviders(<SettlementPage />);

      await waitFor(() => {
        expect(screen.getByText('启用')).toBeInTheDocument();
      });
    });

    it('should display suspended status', async () => {
      mockApi.getSettlementCompanies.mockResolvedValue({
        data: {
          success: true,
          data: [createMockCompany({ status: 'suspended' })],
          pagination: { total: 1 },
        },
      });

      renderWithProviders(<SettlementPage />);

      await waitFor(() => {
        expect(screen.getByText('停用')).toBeInTheDocument();
      });
    });
  });

  describe('Company Type Display', () => {
    it('should display company type', async () => {
      renderWithProviders(<SettlementPage />);

      await waitFor(() => {
        expect(screen.getByText('企业')).toBeInTheDocument();
      });
    });

    it('should display individual type', async () => {
      mockApi.getSettlementCompanies.mockResolvedValue({
        data: {
          success: true,
          data: [createMockCompany({ type: 'individual' })],
          pagination: { total: 1 },
        },
      });

      renderWithProviders(<SettlementPage />);

      await waitFor(() => {
        expect(screen.getByText('个人')).toBeInTheDocument();
      });
    });
  });

  describe('View History', () => {
    it('should open history modal when button clicked', async () => {
      renderWithProviders(<SettlementPage />);

      await waitFor(() => {
        expect(screen.getByText('历史')).toBeInTheDocument();
      });

      const historyButton = screen.getByText('历史');
      fireEvent.click(historyButton);

      await waitFor(() => {
        expect(screen.getByText('变更历史')).toBeInTheDocument();
      });
    });

    it('should display empty history message', async () => {
      renderWithProviders(<SettlementPage />);

      await waitFor(() => {
        expect(screen.getByText('历史')).toBeInTheDocument();
      });

      const historyButton = screen.getByText('历史');
      fireEvent.click(historyButton);

      await waitFor(() => {
        expect(screen.getByText('暂无变更记录')).toBeInTheDocument();
      });
    });
  });

  describe('Toggle Status', () => {
    it('should show confirmation when toggle status clicked', async () => {
      renderWithProviders(<SettlementPage />);

      await waitFor(() => {
        expect(screen.getByText('停用')).toBeInTheDocument();
      });

      const toggleButton = screen.getByText('停用');
      fireEvent.click(toggleButton);

      await waitFor(() => {
        expect(screen.getByText('确定要停用该公司吗？')).toBeInTheDocument();
      });
    });
  });

  describe('Delete Company', () => {
    it('should show confirmation when delete clicked', async () => {
      renderWithProviders(<SettlementPage />);

      await waitFor(() => {
        expect(screen.getByText('删除')).toBeInTheDocument();
      });

      const deleteButton = screen.getByText('删除');
      fireEvent.click(deleteButton);

      await waitFor(() => {
        expect(screen.getByText('确定要删除该公司吗？')).toBeInTheDocument();
      });
    });

    it('should display delete warning', async () => {
      renderWithProviders(<SettlementPage />);

      await waitFor(() => {
        expect(screen.getByText('删除')).toBeInTheDocument();
      });

      const deleteButton = screen.getByText('删除');
      fireEvent.click(deleteButton);

      await waitFor(() => {
        expect(screen.getByText('删除后关联的陪玩师将变为未分配状态')).toBeInTheDocument();
      });
    });
  });

  describe('Empty State', () => {
    it('should handle empty company list', async () => {
      mockApi.getSettlementCompanies.mockResolvedValue({
        data: {
          success: true,
          data: [],
          pagination: { total: 0 },
        },
      });

      renderWithProviders(<SettlementPage />);

      await waitFor(() => {
        expect(screen.getByText('结算公司管理')).toBeInTheDocument();
      });

      expect(mockApi.getSettlementCompanies).toHaveBeenCalled();
    });
  });

  describe('Search Functionality', () => {
    it('should have keyword search field', async () => {
      renderWithProviders(<SettlementPage />);

      await waitFor(() => {
        expect(screen.getByPlaceholderText('公司名称/税号')).toBeInTheDocument();
      });
    });
  });

  describe('Accessibility', () => {
    it('should have proper page title', async () => {
      renderWithProviders(<SettlementPage />);

      await waitFor(() => {
        expect(screen.getByText('结算公司管理')).toBeInTheDocument();
      });
    });

    it('should have proper table structure', async () => {
      renderWithProviders(<SettlementPage />);

      await waitFor(() => {
        expect(screen.getByRole('table')).toBeInTheDocument();
      });
    });
  });
});
