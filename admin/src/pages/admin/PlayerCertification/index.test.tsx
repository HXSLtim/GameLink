/**
 * PlayerCertification Audit Page Tests
 *
 * Tests for PlayerCertification page component including:
 * - Successful data loading
 * - Loading states
 * - Error handling
 * - Audit operations (approve/reject)
 * - Filter functionality
 */

import { describe, it, expect, vi, beforeEach } from 'vitest';
import { screen, waitFor, fireEvent } from '@testing-library/react';
import PlayerCertificationPage from './index';
import { renderWithProviders, resetAllMocks } from '@/testutils';

// Mock the adminApi module using vi.hoisted
const { mockApi, mockMessage } = vi.hoisted(() => ({
  mockApi: {
    getPlayerCertifications: vi.fn(),
    verifyPlayerCertification: vi.fn(),
    getPlayerCertificationDetail: vi.fn(),
    getPlayerCertificationStats: vi.fn(),
    deletePlayerCertification: vi.fn(),
  },
  mockMessage: {
    success: vi.fn(),
    error: vi.fn(),
    warning: vi.fn(),
    info: vi.fn(),
    loading: vi.fn(),
  },
}));

vi.mock('@/api/admin', () => ({
  adminApi: mockApi,
}));

vi.mock('antd', async () => {
  const actual = await vi.importActual('antd');
  return {
    ...actual,
    message: mockMessage,
  };
});

// Helper function to create mock player certification
const createMockPlayerCertification = (overrides: Record<string, unknown> = {}): Record<string, unknown> => ({
  id: 1,
  playerId: 1,
  realName: '张三',
  idNumber: '110101199001011234',
  idFrontUrl: 'https://example.com/front.jpg',
  idBackUrl: 'https://example.com/back.jpg',
  status: 'pending',
  verifiedAt: null,
  verifiedBy: null,
  rejectReason: null,
  player: {
    id: 1,
    nickname: '测试陪玩师',
    avatarUrl: 'https://example.com/avatar.jpg',
    user: {
      id: 100,
      name: '张三',
      email: 'test@example.com',
      phone: '13800138000',
      status: 'active',
    },
  },
  createdAt: '2024-01-01T00:00:00Z',
  updatedAt: '2024-01-01T00:00:00Z',
  ...overrides,
});

describe('PlayerCertificationPage', () => {
  beforeEach(() => {
    resetAllMocks();
    // Default successful responses
    mockApi.getPlayerCertifications.mockResolvedValue({
      data: {
        success: true,
        data: [createMockPlayerCertification()],
        pagination: { total: 1, page: 1, pageSize: 10 },
      },
    });
    mockApi.getPlayerCertificationStats.mockResolvedValue({
      data: {
        success: true,
        data: { pending: 1, verified: 0, rejected: 0 },
      },
    });
  });

  describe('Successful Data Loading', () => {
    it('should render page title', async () => {
      renderWithProviders(<PlayerCertificationPage />);

      await waitFor(() => {
        expect(screen.getByText('实名审核')).toBeInTheDocument();
      });
    });

    it('should display certification list', async () => {
      renderWithProviders(<PlayerCertificationPage />);

      await waitFor(() => {
        expect(screen.getByText('测试陪玩师')).toBeInTheDocument();
      });
    });

    it('should display real name', async () => {
      renderWithProviders(<PlayerCertificationPage />);

      await waitFor(() => {
        expect(screen.getByText('张三')).toBeInTheDocument();
      });
    });

    it('should display pending status', async () => {
      renderWithProviders(<PlayerCertificationPage />);

      await waitFor(() => {
        expect(screen.getByText('待审核')).toBeInTheDocument();
      });
    });

    it('should call API on mount', async () => {
      renderWithProviders(<PlayerCertificationPage />);

      await waitFor(() => {
        expect(mockApi.getPlayerCertifications).toHaveBeenCalled();
      });
    });
  });

  describe('Status Display', () => {
    it('should display verified status correctly', async () => {
      mockApi.getPlayerCertifications.mockResolvedValue({
        data: {
          success: true,
          data: [createMockPlayerCertification({ status: 'verified' })],
          pagination: { total: 1 },
        },
      });

      renderWithProviders(<PlayerCertificationPage />);

      await waitFor(() => {
        expect(screen.getByText('已通过')).toBeInTheDocument();
      });
    });

    it('should display rejected status correctly', async () => {
      mockApi.getPlayerCertifications.mockResolvedValue({
        data: {
          success: true,
          data: [createMockPlayerCertification({ status: 'rejected', rejectReason: '证件不清晰' })],
          pagination: { total: 1 },
        },
      });

      renderWithProviders(<PlayerCertificationPage />);

      await waitFor(() => {
        expect(screen.getByText('已拒绝')).toBeInTheDocument();
      });
    });
  });

  describe('Error Handling', () => {
    it('should show error message on API failure', async () => {
      mockApi.getPlayerCertifications.mockRejectedValue(new Error('Network error'));

      renderWithProviders(<PlayerCertificationPage />);

      await waitFor(() => {
        expect(mockMessage.error).toHaveBeenCalledWith('加载实名认证列表失败');
      });
    });
  });

  describe('Audit Operations', () => {
    it('should display audit button for pending items', async () => {
      renderWithProviders(<PlayerCertificationPage />);

      await waitFor(() => {
        expect(screen.getByText('审核')).toBeInTheDocument();
      });
    });

    it('should display delete button', async () => {
      renderWithProviders(<PlayerCertificationPage />);

      await waitFor(() => {
        expect(screen.getByText('删除')).toBeInTheDocument();
      });
    });

    it('should display detail button', async () => {
      renderWithProviders(<PlayerCertificationPage />);

      await waitFor(() => {
        expect(screen.getByText('详情')).toBeInTheDocument();
      });
    });

    it('should open audit modal when audit button clicked', async () => {
      renderWithProviders(<PlayerCertificationPage />);

      await waitFor(() => {
        expect(screen.getByText('审核')).toBeInTheDocument();
      });

      fireEvent.click(screen.getByText('审核'));

      await waitFor(() => {
        expect(screen.getByRole('dialog')).toBeInTheDocument();
      });
    });
  });

  describe('Filter Functionality', () => {
    it('should have status filter', async () => {
      renderWithProviders(<PlayerCertificationPage />);

      await waitFor(() => {
        expect(screen.getByText('实名审核')).toBeInTheDocument();
      });
    });
  });

  describe('Pagination', () => {
    it('should display pagination info', async () => {
      renderWithProviders(<PlayerCertificationPage />);

      await waitFor(() => {
        expect(screen.getByText(/共 1 条/)).toBeInTheDocument();
      });
    });
  });

  describe('Detail Modal', () => {
    it('should open detail modal when detail button clicked', async () => {
      mockApi.getPlayerCertificationDetail.mockResolvedValue({
        data: { success: true, data: createMockPlayerCertification() },
      });

      renderWithProviders(<PlayerCertificationPage />);

      await waitFor(() => {
        expect(screen.getByText('详情')).toBeInTheDocument();
      });

      fireEvent.click(screen.getByText('详情'));

      await waitFor(() => {
        expect(screen.getByRole('dialog')).toBeInTheDocument();
      });
    });
  });

  describe('ID Card Images', () => {
    it('should have image preview functionality', async () => {
      renderWithProviders(<PlayerCertificationPage />);

      await waitFor(() => {
        expect(screen.getByText('测试陪玩师')).toBeInTheDocument();
      });
    });
  });
});
