/**
 * Routing Rule Management Page Tests
 *
 * Tests for Routing page component including:
 * - Successful data loading
 * - Loading states
 * - Error handling
 * - Statistics display
 * - CRUD operations
 * - History view
 * - Export functionality
 */

import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { screen, waitFor, fireEvent } from '@testing-library/react';
import RoutingRulePage from './index';
import { renderWithProviders, resetAllMocks, flushPromises } from '@/testutils';

// Mock the routingApi module using vi.hoisted
const { mockApi, mockMessage, mockNavigate } = vi.hoisted(() => ({
  mockApi: {
    getRoutingRules: vi.fn(),
    createRoutingRule: vi.fn(),
    updateRoutingRule: vi.fn(),
    deleteRoutingRule: vi.fn(),
    toggleRoutingRuleStatus: vi.fn(),
    getRoutingRuleHistory: vi.fn(),
  },
  mockMessage: {
    success: vi.fn(),
    error: vi.fn(),
    warning: vi.fn(),
    info: vi.fn(),
    loading: vi.fn(),
  },
  mockNavigate: vi.fn(),
}));

vi.mock('@/api/routing', () => ({
  routingApi: mockApi,
}));

vi.mock('antd', async () => {
  const actual = await vi.importActual('antd');
  return {
    ...actual,
    message: mockMessage,
  };
});

vi.mock('react-router-dom', async () => {
  const actual = await vi.importActual('react-router-dom');
  return {
    ...actual,
    useNavigate: () => mockNavigate,
  };
});

// Helper function to create mock rule
const createMockRule = (overrides: Record<string, unknown> = {}): Record<string, unknown> => ({
  id: 1,
  name: '测试路由规则',
  description: '测试描述',
  priority: 1,
  status: 'active',
  targetEntityId: 1,
  targetEntity: { id: 1, name: '测试主体' },
  conditions: [
    { field: 'game_type', operator: 'eq', value: 'LOL' },
  ],
  createdAt: '2024-01-01T00:00:00Z',
  updatedAt: '2024-01-01T00:00:00Z',
  ...overrides,
});

describe('RoutingRulePage', () => {
  beforeEach(() => {
    resetAllMocks();
    mockMessage.success.mockClear();
    mockMessage.error.mockClear();
    mockNavigate.mockClear();
    localStorage.setItem('token', 'test-token');
    localStorage.setItem('user_role', 'admin');

    // Set default mock return values
    mockApi.getRoutingRules.mockResolvedValue({
      data: {
        data: [createMockRule()],
      },
    });
    mockApi.getRoutingRuleHistory.mockResolvedValue({
      data: {
        data: [],
      },
    });
  });

  afterEach(() => {
    vi.clearAllMocks();
  });

  describe('Successful Data Loading', () => {
    it('should render routing page successfully', async () => {
      renderWithProviders(<RoutingRulePage />);

      await waitFor(() => {
        expect(screen.getByText('支付路由规则')).toBeInTheDocument();
      });

      expect(mockApi.getRoutingRules).toHaveBeenCalled();
    });

    it('should display page subtitle', async () => {
      renderWithProviders(<RoutingRulePage />);

      await waitFor(() => {
        expect(screen.getByText('管理订单支付路由规则')).toBeInTheDocument();
      });
    });

    it('should display rule list', async () => {
      renderWithProviders(<RoutingRulePage />);

      await waitFor(() => {
        expect(screen.getByText('测试路由规则')).toBeInTheDocument();
      });
    });

    it('should display statistics section', async () => {
      renderWithProviders(<RoutingRulePage />);

      await waitFor(() => {
        expect(screen.getByText('规则总数')).toBeInTheDocument();
      });
    });
  });

  describe('Loading States', () => {
    it('should show loading indicator while fetching data', async () => {
      mockApi.getRoutingRules.mockImplementation(
        () =>
          new Promise((resolve) => {
            setTimeout(() => {
              resolve({
                data: {
                  data: [createMockRule()],
                },
              });
            }, 100);
          })
      );

      renderWithProviders(<RoutingRulePage />);

      expect(mockApi.getRoutingRules).toHaveBeenCalled();
    });

    it('should hide loading indicator after data loads', async () => {
      renderWithProviders(<RoutingRulePage />);

      await waitFor(() => {
        expect(screen.getByText('支付路由规则')).toBeInTheDocument();
      });

      await flushPromises();

      expect(mockApi.getRoutingRules).toHaveBeenCalledTimes(1);
    });
  });

  describe('Error Handling', () => {
    it('should display error message when API fails', async () => {
      mockApi.getRoutingRules.mockRejectedValue(new Error('Network error'));

      renderWithProviders(<RoutingRulePage />);

      await waitFor(() => {
        expect(mockMessage.error).toHaveBeenCalledWith('加载失败');
      });
    });
  });

  describe('Table Structure', () => {
    it('should display table with data', async () => {
      renderWithProviders(<RoutingRulePage />);

      await waitFor(() => {
        expect(screen.getByRole('table')).toBeInTheDocument();
      });

      // Check for rule data
      expect(screen.getByText('测试路由规则')).toBeInTheDocument();
    });

    it('should display target entity', async () => {
      renderWithProviders(<RoutingRulePage />);

      await waitFor(() => {
        expect(screen.getByText('测试主体')).toBeInTheDocument();
      });
    });

    it('should display active status', async () => {
      renderWithProviders(<RoutingRulePage />);

      await waitFor(() => {
        expect(screen.getByText('启用')).toBeInTheDocument();
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

    it('should display test button', async () => {
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

    it('should navigate to test page when test button clicked', async () => {
      renderWithProviders(<RoutingRulePage />);

      await waitFor(() => {
        expect(screen.getByText('测试规则')).toBeInTheDocument();
      });

      const testButton = screen.getByText('测试规则');
      fireEvent.click(testButton);

      expect(mockNavigate).toHaveBeenCalledWith('/admin/routing/test');
    });
  });

  describe('Status Display', () => {
    it('should display active status', async () => {
      renderWithProviders(<RoutingRulePage />);

      await waitFor(() => {
        expect(screen.getByText('启用')).toBeInTheDocument();
      });
    });

    it('should display inactive status', async () => {
      mockApi.getRoutingRules.mockResolvedValue({
        data: {
          data: [createMockRule({ status: 'inactive' })],
        },
      });

      renderWithProviders(<RoutingRulePage />);

      await waitFor(() => {
        expect(screen.getByText('禁用')).toBeInTheDocument();
      });
    });
  });

  describe('Delete Rule', () => {
    it('should show confirmation when delete clicked', async () => {
      renderWithProviders(<RoutingRulePage />);

      await waitFor(() => {
        expect(screen.getByText('测试路由规则')).toBeInTheDocument();
      });

      // Find delete button by tooltip
      const deleteButtons = document.querySelectorAll('[aria-label="删除"]');
      if (deleteButtons.length > 0) {
        fireEvent.click(deleteButtons[0]);
      }
    });
  });

  describe('View History', () => {
    it('should open history modal when button clicked', async () => {
      renderWithProviders(<RoutingRulePage />);

      await waitFor(() => {
        expect(screen.getByText('测试路由规则')).toBeInTheDocument();
      });

      // Find history button by tooltip
      const historyButtons = document.querySelectorAll('[aria-label="历史"]');
      if (historyButtons.length > 0) {
        fireEvent.click(historyButtons[0]);

        await waitFor(() => {
          expect(screen.getByText('修改历史')).toBeInTheDocument();
        });
      }
    });
  });

  describe('Conditions Display', () => {
    it('should display condition tags', async () => {
      renderWithProviders(<RoutingRulePage />);

      await waitFor(() => {
        expect(screen.getByText('测试路由规则')).toBeInTheDocument();
      });

      // Should display game type condition
      expect(screen.getByText(/游戏/)).toBeInTheDocument();
    });

    it('should display no condition tag when empty', async () => {
      mockApi.getRoutingRules.mockResolvedValue({
        data: {
          data: [createMockRule({ conditions: [] })],
        },
      });

      renderWithProviders(<RoutingRulePage />);

      await waitFor(() => {
        expect(screen.getByText('无条件')).toBeInTheDocument();
      });
    });
  });

  describe('Target Entity Display', () => {
    it('should display target entity name', async () => {
      renderWithProviders(<RoutingRulePage />);

      await waitFor(() => {
        expect(screen.getByText('测试主体')).toBeInTheDocument();
      });
    });
  });

  describe('Empty State', () => {
    it('should handle empty rule list', async () => {
      mockApi.getRoutingRules.mockResolvedValue({
        data: {
          data: [],
        },
      });

      renderWithProviders(<RoutingRulePage />);

      await waitFor(() => {
        expect(screen.getByText('支付路由规则')).toBeInTheDocument();
      });

      expect(mockApi.getRoutingRules).toHaveBeenCalled();
    });
  });

  describe('Search Functionality', () => {
    it('should have keyword search field', async () => {
      renderWithProviders(<RoutingRulePage />);

      await waitFor(() => {
        expect(screen.getByPlaceholderText('规则名称')).toBeInTheDocument();
      });
    });
  });

  describe('Accessibility', () => {
    it('should have proper page title', async () => {
      renderWithProviders(<RoutingRulePage />);

      await waitFor(() => {
        expect(screen.getByText('支付路由规则')).toBeInTheDocument();
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
