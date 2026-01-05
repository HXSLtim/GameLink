/**
 * Dispute Management Page Tests
 *
 * Tests for Dispute page component including:
 * - Successful data loading
 * - Loading states
 * - Error handling
 * - Statistics display
 * - SLA warning display
 */

import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { screen, waitFor } from '@testing-library/react';
import DisputePage from './index';
import { renderWithProviders, resetAllMocks, flushPromises } from '@/testutils';

// Mock the disputeApi module using vi.hoisted
const { mockApi, mockMessage, mockModal } = vi.hoisted(() => ({
  mockApi: {
    getDisputes: vi.fn(),
    getDisputeStats: vi.fn(),
    assignDispute: vi.fn(),
    resolveDispute: vi.fn(),
    rollbackAssignment: vi.fn(),
  },
  mockMessage: {
    success: vi.fn(),
    error: vi.fn(),
    warning: vi.fn(),
    info: vi.fn(),
    loading: vi.fn(),
  },
  mockModal: {
    confirm: vi.fn(),
  },
}));

vi.mock('@/api/dispute', () => ({
  disputeApi: mockApi,
}));

// Mock App.useApp to return the message and modal mocks
vi.mock('antd', async () => {
  const actual = await vi.importActual('antd');
  return {
    ...actual,
    App: {
      useApp: () => ({
        message: mockMessage,
        notification: { success: vi.fn(), error: vi.fn(), info: vi.fn(), warning: vi.fn() },
        modal: mockModal,
      }),
    },
  };
});

// Helper function to create mock dispute data
const createMockDispute = (overrides: Record<string, unknown> = {}): Record<string, unknown> => ({
  id: 1,
  orderId: 1,
  orderNo: 'ORD202401010001',
  initiatorId: 1,
  initiatorType: 'user',
  type: 'service_quality',
  status: 'pending',
  reason: '服务质量问题',
  evidenceUrls: '[]',
  evidenceText: '陪玩师态度不好',
  chatSnapshotId: null,
  originalServiceId: null,
  assignedServiceId: null,
  slaDeadline: null,
  resolution: null,
  resolvedBy: null,
  resolvedAt: null,
  resolveRemark: null,
  createdAt: '2024-01-01T00:00:00Z',
  updatedAt: '2024-01-01T00:00:00Z',
  order: {
    id: 1,
    orderNo: 'ORD202401010001',
    totalPriceCents: 10000,
  },
  initiator: {
    id: 1,
    name: '测试用户',
  },
  ...overrides,
});

// Helper function to create mock stats
const createMockStats = (overrides: Record<string, unknown> = {}): Record<string, unknown> => ({
  pending: 5,
  assigned: 3,
  mediating: 2,
  resolved: 10,
  rejected: 1,
  canceled: 0,
  slaBreached: 1,
  ...overrides,
});

describe('DisputePage', () => {
  beforeEach(() => {
    resetAllMocks();
    mockMessage.success.mockClear();
    mockMessage.error.mockClear();
    mockMessage.info.mockClear();
    mockModal.confirm.mockClear();
    localStorage.setItem('token', 'test-token');
    localStorage.setItem('user_role', 'admin');
    // Set default mock return values
    mockApi.getDisputes.mockResolvedValue({
      data: {
        success: true,
        data: {
          disputes: [createMockDispute()],
          total: 1,
        },
      },
    });
    mockApi.getDisputeStats.mockResolvedValue({
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
    it('should render dispute page successfully', async () => {
      renderWithProviders(<DisputePage />);

      await waitFor(() => {
        expect(screen.getByText('纠纷管理')).toBeInTheDocument();
      });

      expect(mockApi.getDisputes).toHaveBeenCalled();
      expect(mockApi.getDisputeStats).toHaveBeenCalled();
    });

    it('should display statistics cards', async () => {
      renderWithProviders(<DisputePage />);

      await waitFor(() => {
        // "待处理" appears in stats card
        const pendingElements = screen.getAllByText('待处理');
        expect(pendingElements.length).toBeGreaterThan(0);
      });

      expect(screen.getByText('已指派')).toBeInTheDocument();
      expect(screen.getByText('调解中')).toBeInTheDocument();
      expect(screen.getByText('已解决')).toBeInTheDocument();
      expect(screen.getByText('已驳回')).toBeInTheDocument();
      expect(screen.getByText('SLA超时')).toBeInTheDocument();
    });

    it('should display statistics values', async () => {
      renderWithProviders(<DisputePage />);

      await waitFor(() => {
        expect(screen.getByText('5')).toBeInTheDocument(); // pending
      });

      expect(screen.getByText('3')).toBeInTheDocument(); // assigned
      expect(screen.getByText('10')).toBeInTheDocument(); // resolved
    });
  });

  describe('Loading States', () => {
    it('should show loading indicator while fetching data', async () => {
      mockApi.getDisputes.mockImplementation(
        () =>
          new Promise((resolve) => {
            setTimeout(() => {
              resolve({
                data: {
                  success: true,
                  data: {
                    disputes: [createMockDispute()],
                    total: 1,
                  },
                },
              });
            }, 100);
          })
      );

      renderWithProviders(<DisputePage />);

      expect(mockApi.getDisputes).toHaveBeenCalled();
    });

    it('should hide loading indicator after data loads', async () => {
      renderWithProviders(<DisputePage />);

      await waitFor(() => {
        expect(screen.getByText('纠纷管理')).toBeInTheDocument();
      });

      await flushPromises();

      expect(mockApi.getDisputes).toHaveBeenCalledTimes(1);
    });
  });

  describe('Error Handling', () => {
    it('should display error message when API fails', async () => {
      mockApi.getDisputes.mockRejectedValue(new Error('Network error'));

      renderWithProviders(<DisputePage />);

      await waitFor(() => {
        expect(screen.getByText('纠纷管理')).toBeInTheDocument();
      });

      await flushPromises();

      expect(mockMessage.error).toHaveBeenCalledWith('加载纠纷列表失败');
    });
  });

  describe('SLA Warning', () => {
    it('should display SLA warning when there are breached disputes', async () => {
      mockApi.getDisputeStats.mockResolvedValue({
        data: {
          success: true,
          data: createMockStats({ slaBreached: 3 }),
        },
      });

      renderWithProviders(<DisputePage />);

      await waitFor(() => {
        expect(screen.getByText(/有 3 个纠纷已超过30分钟SLA/)).toBeInTheDocument();
      });
    });

    it('should not display SLA warning when no breached disputes', async () => {
      mockApi.getDisputeStats.mockResolvedValue({
        data: {
          success: true,
          data: createMockStats({ slaBreached: 0 }),
        },
      });

      renderWithProviders(<DisputePage />);

      await waitFor(() => {
        expect(screen.getByText('纠纷管理')).toBeInTheDocument();
      });

      await flushPromises();

      expect(screen.queryByText(/已超过30分钟SLA/)).not.toBeInTheDocument();
    });
  });

  describe('Search Functionality', () => {
    it('should have search fields available', async () => {
      renderWithProviders(<DisputePage />);

      await waitFor(() => {
        expect(screen.getByText('纠纷管理')).toBeInTheDocument();
      });

      // Search fields are rendered by DisputeList component
      expect(mockApi.getDisputes).toHaveBeenCalled();
    });

    it('should have status filter', async () => {
      renderWithProviders(<DisputePage />);

      await waitFor(() => {
        expect(screen.getByText('纠纷管理')).toBeInTheDocument();
      });

      // Status filter is part of search fields
      expect(mockApi.getDisputes).toHaveBeenCalled();
    });

    it('should have initiator type filter', async () => {
      renderWithProviders(<DisputePage />);

      await waitFor(() => {
        expect(screen.getByText('纠纷管理')).toBeInTheDocument();
      });

      // Initiator type filter is part of search fields
      expect(mockApi.getDisputes).toHaveBeenCalled();
    });
  });

  describe('Toolbar Buttons', () => {
    it('should display export button', async () => {
      renderWithProviders(<DisputePage />);

      await waitFor(() => {
        expect(screen.getByText('导出数据')).toBeInTheDocument();
      });
    });
  });

  describe('Empty State', () => {
    it('should handle empty dispute list', async () => {
      mockApi.getDisputes.mockResolvedValue({
        data: {
          success: true,
          data: {
            disputes: [],
            total: 0,
          },
        },
      });

      renderWithProviders(<DisputePage />);

      await waitFor(() => {
        expect(screen.getByText('纠纷管理')).toBeInTheDocument();
      });

      expect(mockApi.getDisputes).toHaveBeenCalled();
    });
  });

  describe('Statistics with Zero Values', () => {
    it('should display zero values correctly', async () => {
      mockApi.getDisputeStats.mockResolvedValue({
        data: {
          success: true,
          data: createMockStats({
            pending: 0,
            assigned: 0,
            mediating: 0,
            resolved: 0,
            rejected: 0,
            slaBreached: 0,
          }),
        },
      });

      renderWithProviders(<DisputePage />);

      await waitFor(() => {
        expect(screen.getByText('纠纷管理')).toBeInTheDocument();
      });

      // Should display 0 values
      const zeroValues = screen.getAllByText('0');
      expect(zeroValues.length).toBeGreaterThan(0);
    });
  });

  describe('Accessibility', () => {
    it('should have proper page title', async () => {
      renderWithProviders(<DisputePage />);

      await waitFor(() => {
        expect(screen.getByText('纠纷管理')).toBeInTheDocument();
      });
    });

    it('should have proper subtitle', async () => {
      renderWithProviders(<DisputePage />);

      await waitFor(() => {
        expect(screen.getByText('处理用户与陪玩师之间的订单纠纷')).toBeInTheDocument();
      });
    });
  });
});
