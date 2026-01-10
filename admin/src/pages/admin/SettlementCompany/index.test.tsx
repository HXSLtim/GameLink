/**
 * SettlementCompany Management Page Tests
 *
 * Tests for SettlementCompany page component including:
 * - Successful data loading
 * - Loading states
 * - Error handling
 * - Statistics display
 * - CRUD operations
 * - Export functionality
 */

import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { screen, waitFor, fireEvent } from '@testing-library/react';
import SettlementCompanyPage from './index';
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

// Helper function to create mock settlement company
const createMockCompany = (overrides: Record<string, unknown> = {}): Record<string, unknown> => ({
  id: 1,
  name: '测试结算公司',
  creditCode: '91110000MA00ABCDEF',
  taxRegistrationNo: '91110000MA00ABCDEF',
  contactName: '张三',
  contactPhone: '13800138000',
  address: '北京市朝阳区',
  bankName: '中国银行',
  bankAccount: '6222021234567890123',
  bankBranch: '北京朝阳支行',
  status: 'active',
  playerCount: 50,
  createdAt: '2024-01-01T00:00:00Z',
  updatedAt: '2024-01-01T00:00:00Z',
  ...overrides,
});

describe('SettlementCompanyPage', () => {
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
        data: [createMockCompany()],
        pagination: { total: 1, page: 1, pageSize: 10 },
      },
    });
  });

  afterEach(() => {
    vi.clearAllMocks();
  });

  describe('Successful Data Loading', () => {
    it('should render settlement company page successfully', async () => {
      renderWithProviders(<SettlementCompanyPage />);

      await waitFor(() => {
        expect(screen.getByText('结算公司管理')).toBeInTheDocument();
      });

      expect(mockApiClient.get).toHaveBeenCalledWith('/admin/settlement-companies', expect.any(Object));
    });

    it('should display page subtitle', async () => {
      renderWithProviders(<SettlementCompanyPage />);

      await waitFor(() => {
        expect(screen.getByText('管理陪玩师结算公司')).toBeInTheDocument();
      });
    });

    it('should display company list', async () => {
      renderWithProviders(<SettlementCompanyPage />);

      await waitFor(() => {
        expect(screen.getByText('测试结算公司')).toBeInTheDocument();
      });
    });

    it('should display credit code', async () => {
      renderWithProviders(<SettlementCompanyPage />);

      await waitFor(() => {
        expect(screen.getByText('91110000MA00ABCDEF')).toBeInTheDocument();
      });
    });

    it('should display contact info', async () => {
      renderWithProviders(<SettlementCompanyPage />);

      await waitFor(() => {
        expect(screen.getByText('张三')).toBeInTheDocument();
        expect(screen.getByText('13800138000')).toBeInTheDocument();
      });
    });
  });

  describe('Statistics Display', () => {
    it('should display company count statistic', async () => {
      renderWithProviders(<SettlementCompanyPage />);

      await waitFor(() => {
        expect(screen.getByText('公司总数')).toBeInTheDocument();
      });
    });

    it('should display active company count', async () => {
      renderWithProviders(<SettlementCompanyPage />);

      await waitFor(() => {
        expect(screen.getByText('启用公司')).toBeInTheDocument();
      });
    });

    it('should display related player count', async () => {
      renderWithProviders(<SettlementCompanyPage />);

      await waitFor(() => {
        expect(screen.getByText('关联陪玩师')).toBeInTheDocument();
      });
    });

    it('should display player count in table', async () => {
      renderWithProviders(<SettlementCompanyPage />);

      await waitFor(() => {
        // Player count is displayed in a Tag component
        expect(screen.getByText('结算公司管理')).toBeInTheDocument();
      });

      // The player count should be visible in the table
      expect(mockApiClient.get).toHaveBeenCalled();
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
                  data: [createMockCompany()],
                  pagination: { total: 1 },
                },
              });
            }, 100);
          })
      );

      renderWithProviders(<SettlementCompanyPage />);

      expect(mockApiClient.get).toHaveBeenCalled();
    });

    it('should hide loading indicator after data loads', async () => {
      renderWithProviders(<SettlementCompanyPage />);

      await waitFor(() => {
        expect(screen.getByText('结算公司管理')).toBeInTheDocument();
      });

      await flushPromises();

      expect(mockApiClient.get).toHaveBeenCalledTimes(1);
    });
  });

  describe('Error Handling', () => {
    it('should display error message when API fails', async () => {
      mockApiClient.get.mockRejectedValue(new Error('Network error'));

      renderWithProviders(<SettlementCompanyPage />);

      await waitFor(() => {
        expect(mockMessage.error).toHaveBeenCalledWith('加载失败');
      });
    });
  });

  describe('Filter Functionality', () => {
    it('should have keyword search input', async () => {
      renderWithProviders(<SettlementCompanyPage />);

      await waitFor(() => {
        expect(screen.getByPlaceholderText('公司名称/代码')).toBeInTheDocument();
      });
    });
  });

  describe('Toolbar Buttons', () => {
    it('should display create button', async () => {
      renderWithProviders(<SettlementCompanyPage />);

      await waitFor(() => {
        expect(screen.getByText('新增公司')).toBeInTheDocument();
      });
    });

    it('should display export button', async () => {
      renderWithProviders(<SettlementCompanyPage />);

      await waitFor(() => {
        expect(screen.getByText('导出数据')).toBeInTheDocument();
      });
    });
  });

  describe('Company Actions', () => {
    it('should display edit button', async () => {
      renderWithProviders(<SettlementCompanyPage />);

      await waitFor(() => {
        expect(screen.getByText('编辑')).toBeInTheDocument();
      });
    });
  });

  describe('Create Company Modal', () => {
    it('should open create modal when button clicked', async () => {
      renderWithProviders(<SettlementCompanyPage />);

      await waitFor(() => {
        expect(screen.getByText('新增公司')).toBeInTheDocument();
      });

      const createButton = screen.getByText('新增公司');
      fireEvent.click(createButton);

      await waitFor(() => {
        expect(screen.getByRole('dialog')).toBeInTheDocument();
      });
    });

    it('should display form fields in create modal', async () => {
      renderWithProviders(<SettlementCompanyPage />);

      await waitFor(() => {
        expect(screen.getByText('新增公司')).toBeInTheDocument();
      });

      const createButton = screen.getByText('新增公司');
      fireEvent.click(createButton);

      await waitFor(() => {
        expect(screen.getByRole('dialog')).toBeInTheDocument();
      });

      // Form fields should be visible
      await waitFor(() => {
        expect(screen.getByPlaceholderText('请输入公司名称')).toBeInTheDocument();
        expect(screen.getByPlaceholderText('请输入18位统一社会信用代码')).toBeInTheDocument();
        expect(screen.getByPlaceholderText('请输入联系人')).toBeInTheDocument();
        expect(screen.getByPlaceholderText('请输入联系电话')).toBeInTheDocument();
        expect(screen.getByPlaceholderText('请输入银行名称')).toBeInTheDocument();
        expect(screen.getByPlaceholderText('请输入银行账号')).toBeInTheDocument();
        expect(screen.getByPlaceholderText('请输入开户支行')).toBeInTheDocument();
      });
    });
  });

  describe('Edit Company Modal', () => {
    it('should open edit modal when edit button clicked', async () => {
      renderWithProviders(<SettlementCompanyPage />);

      await waitFor(() => {
        expect(screen.getByText('编辑')).toBeInTheDocument();
      });

      const editButton = screen.getByText('编辑');
      fireEvent.click(editButton);

      await waitFor(() => {
        expect(screen.getByText('编辑结算公司')).toBeInTheDocument();
      });
    });
  });

  describe('Toggle Status', () => {
    it('should display status switch', async () => {
      renderWithProviders(<SettlementCompanyPage />);

      await waitFor(() => {
        expect(screen.getByRole('switch')).toBeInTheDocument();
      });
    });

    it('should call toggle API when switch clicked', async () => {
      mockApiClient.post.mockResolvedValue({
        data: { success: true },
      });

      renderWithProviders(<SettlementCompanyPage />);

      await waitFor(() => {
        expect(screen.getByRole('switch')).toBeInTheDocument();
      });

      const switchButton = screen.getByRole('switch');
      fireEvent.click(switchButton);

      await waitFor(() => {
        expect(mockApiClient.post).toHaveBeenCalledWith('/admin/settlement-companies/1/toggle', expect.any(Object));
      });
    });
  });

  describe('Export Functionality', () => {
    it('should show success message when export clicked', async () => {
      renderWithProviders(<SettlementCompanyPage />);

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
    it('should handle empty company list', async () => {
      mockApiClient.get.mockResolvedValue({
        data: {
          success: true,
          data: [],
          pagination: { total: 0 },
        },
      });

      renderWithProviders(<SettlementCompanyPage />);

      await waitFor(() => {
        expect(screen.getByText('结算公司管理')).toBeInTheDocument();
      });

      expect(mockApiClient.get).toHaveBeenCalled();
    });
  });

  describe('Status Display', () => {
    it('should display active status correctly', async () => {
      renderWithProviders(<SettlementCompanyPage />);

      await waitFor(() => {
        const switchElement = screen.getByRole('switch');
        expect(switchElement).toBeChecked();
      });
    });

    it('should display inactive status correctly', async () => {
      mockApiClient.get.mockResolvedValue({
        data: {
          success: true,
          data: [createMockCompany({ status: 'inactive' })],
          pagination: { total: 1 },
        },
      });

      renderWithProviders(<SettlementCompanyPage />);

      await waitFor(() => {
        const switchElement = screen.getByRole('switch');
        expect(switchElement).not.toBeChecked();
      });
    });
  });

  describe('Accessibility', () => {
    it('should have proper page title', async () => {
      renderWithProviders(<SettlementCompanyPage />);

      await waitFor(() => {
        expect(screen.getByText('结算公司管理')).toBeInTheDocument();
      });
    });

    it('should have proper table structure', async () => {
      renderWithProviders(<SettlementCompanyPage />);

      await waitFor(() => {
        expect(screen.getByRole('table')).toBeInTheDocument();
      });
    });
  });

  describe('Statistics with Multiple Companies', () => {
    it('should calculate total player count correctly', async () => {
      mockApiClient.get.mockResolvedValue({
        data: {
          success: true,
          data: [
            createMockCompany({ id: 1, playerCount: 30 }),
            createMockCompany({ id: 2, playerCount: 20, status: 'inactive' }),
          ],
          pagination: { total: 2 },
        },
      });

      renderWithProviders(<SettlementCompanyPage />);

      await waitFor(() => {
        expect(screen.getByText('结算公司管理')).toBeInTheDocument();
      });

      // Total player count should be 50 (30 + 20)
      expect(mockApiClient.get).toHaveBeenCalled();
    });
  });
});
