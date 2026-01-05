/**
 * RoutingRule Management Page Tests
 *
 * Tests for RoutingRule page component including:
 * - Successful data loading
 * - Loading states
 * - Error handling
 * - CRUD operations
 * - Rule testing
 * - History viewing
 * - Export functionality
 */

import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { screen, waitFor, fireEvent } from '@testing-library/react';
import RoutingRulePage from './index';
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

// Helper function to create mock routing rule
const createMockRule = (overrides: Record<string, unknown> = {}): Record<string, unknown> => ({
  id: 1,
  name: '测试分流规则',
  description: '规则描述',
  priority: 1,
  status: 'active',
  targetEntityId: 100,
  targetEntityName: '收款主体A',
  conditions: { field: 'order_amount', operator: 'eq', value: 100 },
  createdAt: '2024-01-01T00:00:00Z',
  updatedAt: '2024-01-01T00:00:00Z',
  ...overrides,
});

// Helper function to create mock history
const createMockHistory = (overrides: Record<string, unknown> = {}): Record<string, unknown> => ({
  id: 1,
  ruleId: 1,
  action: '创建规则',
  changes: '新建分流规则',
  operatorId: 1,
  operatorName: '管理员',
  createdAt: '2024-01-01T00:00:00Z',
  ...overrides,
});

describe('RoutingRulePage', () => {
  beforeEach(() => {
    resetAllMocks();
    mockMessage.success.mockClear();
    mockMessage.error.mockClear();
    localStorage.setItem('token', 'test-token');
    localStorage.setItem('user_role', 'admin');

    // Set default mock return values
    mockApiClient.get.mockImplementation((url: string) => {
      if (url === '/admin/routing-rules') {
        return Promise.resolve({
          data: {
            success: true,
            data: [createMockRule()],
            pagination: { total: 1, page: 1, pageSize: 10 },
          },
        });
      }
      if (url.includes('/history')) {
        return Promise.resolve({
          data: {
            success: true,
            data: [createMockHistory()],
          },
        });
      }
      return Promise.resolve({ data: { success: true, data: null } });
    });
  });

  afterEach(() => {
    vi.clearAllMocks();
  });

  describe('Successful Data Loading', () => {
    it('should render routing rule page successfully', async () => {
      renderWithProviders(<RoutingRulePage />);

      await waitFor(() => {
        expect(screen.getByText('分流规则管理')).toBeInTheDocument();
      });

      expect(mockApiClient.get).toHaveBeenCalledWith('/admin/routing-rules', expect.any(Object));
    });

    it('should display page subtitle', async () => {
      renderWithProviders(<RoutingRulePage />);

      await waitFor(() => {
        expect(screen.getByText('管理订单分流规则')).toBeInTheDocument();
      });
    });

    it('should display rule list', async () => {
      renderWithProviders(<RoutingRulePage />);

      await waitFor(() => {
        expect(screen.getByText('测试分流规则')).toBeInTheDocument();
      });
    });

    it('should display target entity name', async () => {
      renderWithProviders(<RoutingRulePage />);

      await waitFor(() => {
        expect(screen.getByText('收款主体A')).toBeInTheDocument();
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
                  data: [createMockRule()],
                  pagination: { total: 1 },
                },
              });
            }, 100);
          })
      );

      renderWithProviders(<RoutingRulePage />);

      expect(mockApiClient.get).toHaveBeenCalled();
    });

    it('should hide loading indicator after data loads', async () => {
      renderWithProviders(<RoutingRulePage />);

      await waitFor(() => {
        expect(screen.getByText('分流规则管理')).toBeInTheDocument();
      });

      await flushPromises();

      expect(mockApiClient.get).toHaveBeenCalledTimes(1);
    });
  });

  describe('Error Handling', () => {
    it('should display error message when API fails', async () => {
      mockApiClient.get.mockRejectedValue(new Error('Network error'));

      renderWithProviders(<RoutingRulePage />);

      await waitFor(() => {
        expect(mockMessage.error).toHaveBeenCalledWith('加载失败');
      });
    });
  });

  describe('Filter Functionality', () => {
    it('should have keyword search input', async () => {
      renderWithProviders(<RoutingRulePage />);

      await waitFor(() => {
        expect(screen.getByPlaceholderText('规则名称')).toBeInTheDocument();
      });
    });
  });

  describe('Toolbar Buttons', () => {
    it('should display create button', async () => {
      renderWithProviders(<RoutingRulePage />);

      await waitFor(() => {
        expect(screen.getByText('新增规则')).toBeInTheDocument();
      });
    });

    it('should display test rule button', async () => {
      renderWithProviders(<RoutingRulePage />);

      await waitFor(() => {
        expect(screen.getByText('测试规则')).toBeInTheDocument();
      });
    });

    it('should display export button', async () => {
      renderWithProviders(<RoutingRulePage />);

      await waitFor(() => {
        expect(screen.getByText('导出数据')).toBeInTheDocument();
      });
    });
  });

  describe('Rule Actions', () => {
    it('should display edit button', async () => {
      renderWithProviders(<RoutingRulePage />);

      await waitFor(() => {
        expect(screen.getByText('编辑')).toBeInTheDocument();
      });
    });

    it('should display history button', async () => {
      renderWithProviders(<RoutingRulePage />);

      await waitFor(() => {
        expect(screen.getByText('历史')).toBeInTheDocument();
      });
    });

    it('should display delete button', async () => {
      renderWithProviders(<RoutingRulePage />);

      await waitFor(() => {
        expect(screen.getByText('删除')).toBeInTheDocument();
      });
    });
  });

  describe('Create Rule Modal', () => {
    it('should open create modal when button clicked', async () => {
      renderWithProviders(<RoutingRulePage />);

      await waitFor(() => {
        expect(screen.getByText('新增规则')).toBeInTheDocument();
      });

      const createButton = screen.getByText('新增规则');
      fireEvent.click(createButton);

      await waitFor(() => {
        expect(screen.getByText('新增分流规则')).toBeInTheDocument();
      });
    });

    it('should display form fields in create modal', async () => {
      renderWithProviders(<RoutingRulePage />);

      await waitFor(() => {
        expect(screen.getByText('新增规则')).toBeInTheDocument();
      });

      const createButton = screen.getByText('新增规则');
      fireEvent.click(createButton);

      await waitFor(() => {
        expect(screen.getByRole('dialog')).toBeInTheDocument();
      });

      // Check form fields by placeholder (only Input fields, not Select)
      await waitFor(() => {
        expect(screen.getByPlaceholderText('请输入规则名称')).toBeInTheDocument();
        expect(screen.getByPlaceholderText('请输入描述')).toBeInTheDocument();
      });
    });
  });

  describe('Edit Rule Modal', () => {
    it('should open edit modal when edit button clicked', async () => {
      renderWithProviders(<RoutingRulePage />);

      await waitFor(() => {
        expect(screen.getByText('编辑')).toBeInTheDocument();
      });

      const editButton = screen.getByText('编辑');
      fireEvent.click(editButton);

      await waitFor(() => {
        expect(screen.getByText('编辑分流规则')).toBeInTheDocument();
      });
    });
  });

  describe('View History Modal', () => {
    it('should open history modal when history button clicked', async () => {
      renderWithProviders(<RoutingRulePage />);

      await waitFor(() => {
        expect(screen.getByText('历史')).toBeInTheDocument();
      });

      const historyButton = screen.getByText('历史');
      fireEvent.click(historyButton);

      await waitFor(() => {
        expect(screen.getByText('修改历史')).toBeInTheDocument();
      });
    });

    it('should load history data when modal opens', async () => {
      renderWithProviders(<RoutingRulePage />);

      await waitFor(() => {
        expect(screen.getByText('历史')).toBeInTheDocument();
      });

      const historyButton = screen.getByText('历史');
      fireEvent.click(historyButton);

      await waitFor(() => {
        expect(mockApiClient.get).toHaveBeenCalledWith('/admin/routing-rules/1/history');
      });
    });
  });

  describe('Test Rule Modal', () => {
    it('should open test modal when test button clicked', async () => {
      renderWithProviders(<RoutingRulePage />);

      await waitFor(() => {
        expect(screen.getByText('测试规则')).toBeInTheDocument();
      });

      const testButton = screen.getByText('测试规则');
      fireEvent.click(testButton);

      await waitFor(() => {
        expect(screen.getByRole('dialog')).toBeInTheDocument();
      });
    });

    it('should display test form fields', async () => {
      renderWithProviders(<RoutingRulePage />);

      await waitFor(() => {
        expect(screen.getByText('测试规则')).toBeInTheDocument();
      });

      const testButton = screen.getByText('测试规则');
      fireEvent.click(testButton);

      await waitFor(() => {
        expect(screen.getByRole('dialog')).toBeInTheDocument();
      });

      // Form fields should be visible
      await waitFor(() => {
        expect(screen.getByPlaceholderText('请输入订单金额')).toBeInTheDocument();
        expect(screen.getByPlaceholderText('请输入游戏ID')).toBeInTheDocument();
      });
    });
  });

  describe('Delete Rule', () => {
    it('should call delete API when delete button clicked', async () => {
      mockApiClient.delete.mockResolvedValue({
        data: { success: true },
      });

      renderWithProviders(<RoutingRulePage />);

      await waitFor(() => {
        expect(screen.getByText('删除')).toBeInTheDocument();
      });

      const deleteButton = screen.getByText('删除');
      fireEvent.click(deleteButton);

      await waitFor(() => {
        expect(mockApiClient.delete).toHaveBeenCalledWith('/admin/routing-rules/1');
      });
    });
  });

  describe('Toggle Status', () => {
    it('should display status switch', async () => {
      renderWithProviders(<RoutingRulePage />);

      await waitFor(() => {
        expect(screen.getByRole('switch')).toBeInTheDocument();
      });
    });

    it('should call toggle API when switch clicked', async () => {
      mockApiClient.post.mockResolvedValue({
        data: { success: true },
      });

      renderWithProviders(<RoutingRulePage />);

      await waitFor(() => {
        expect(screen.getByRole('switch')).toBeInTheDocument();
      });

      const switchButton = screen.getByRole('switch');
      fireEvent.click(switchButton);

      await waitFor(() => {
        expect(mockApiClient.post).toHaveBeenCalledWith('/admin/routing-rules/1/toggle', expect.any(Object));
      });
    });
  });

  describe('Export Functionality', () => {
    it('should show success message when export clicked', async () => {
      renderWithProviders(<RoutingRulePage />);

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
    it('should handle empty rule list', async () => {
      mockApiClient.get.mockResolvedValue({
        data: {
          success: true,
          data: [],
          pagination: { total: 0 },
        },
      });

      renderWithProviders(<RoutingRulePage />);

      await waitFor(() => {
        expect(screen.getByText('分流规则管理')).toBeInTheDocument();
      });

      expect(mockApiClient.get).toHaveBeenCalled();
    });
  });

  describe('Accessibility', () => {
    it('should have proper page title', async () => {
      renderWithProviders(<RoutingRulePage />);

      await waitFor(() => {
        expect(screen.getByText('分流规则管理')).toBeInTheDocument();
      });
    });

    it('should have proper table structure', async () => {
      renderWithProviders(<RoutingRulePage />);

      await waitFor(() => {
        expect(screen.getByRole('table')).toBeInTheDocument();
      });
    });
  });
});
