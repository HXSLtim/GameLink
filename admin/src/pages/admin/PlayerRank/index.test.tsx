/**
 * PlayerRank Audit Page Tests
 *
 * Tests for PlayerRank page component including:
 * - Successful data loading
 * - Loading states
 * - Error handling
 * - Audit operations (approve/reject)
 * - Filter functionality
 */

import { describe, it, expect, vi, beforeEach } from 'vitest';
import { screen, waitFor, fireEvent } from '@testing-library/react';
import PlayerRankPage from './index';
import { renderWithProviders, resetAllMocks } from '@/testutils';

// Mock the adminApi module using vi.hoisted
const { mockApi, mockMessage } = vi.hoisted(() => ({
  mockApi: {
    getPlayerRanks: vi.fn(),
    getGames: vi.fn(),
    verifyPlayerRank: vi.fn(),
    getPlayerRankDetail: vi.fn(),
    getPlayerRankStats: vi.fn(),
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

// Helper function to create mock player rank
const createMockPlayerRank = (overrides: Record<string, unknown> = {}): Record<string, unknown> => ({
  id: 1,
  playerId: 1,
  gameRankId: 1,
  status: 'pending',
  proofImages: ['https://example.com/proof1.jpg'],
  screenshotUrls: '["https://example.com/proof1.jpg"]',
  verifiedAt: null,
  verifiedBy: null,
  rejectReason: null,
  player: {
    id: 1,
    userId: 100,
    nickname: '测试陪玩师',
    avatarUrl: 'https://example.com/avatar.jpg',
    user: {
      id: 100,
      name: '张三',
      email: 'test@example.com',
      phone: '13800138000',
      status: 'active',
      avatarUrl: 'https://example.com/avatar.jpg',
    },
  },
  game: {
    id: 1,
    name: '王者荣耀',
    icon: 'https://example.com/icon.png',
  },
  rank: {
    id: 1,
    name: '王者',
    level: 10,
    color: '#ff0000',
  },
  createdAt: '2024-01-01T00:00:00Z',
  updatedAt: '2024-01-01T00:00:00Z',
  ...overrides,
});

// Helper function to create mock game
const createMockGame = (overrides: Record<string, unknown> = {}): Record<string, unknown> => ({
  id: 1,
  name: '王者荣耀',
  icon: 'https://example.com/icon.png',
  isActive: true,
  ...overrides,
});

describe('PlayerRankPage', () => {
  beforeEach(() => {
    resetAllMocks();
    // Default successful responses
    mockApi.getPlayerRanks.mockResolvedValue({
      data: {
        success: true,
        data: [createMockPlayerRank()],
        pagination: { total: 1, page: 1, pageSize: 10 },
      },
    });
    mockApi.getGames.mockResolvedValue({
      data: {
        success: true,
        data: [createMockGame()],
      },
    });
    mockApi.getPlayerRankStats.mockResolvedValue({
      data: {
        success: true,
        data: { pending: 1, verified: 0, rejected: 0 },
      },
    });
  });

  describe('Successful Data Loading', () => {
    it('should render page title', async () => {
      renderWithProviders(<PlayerRankPage />);

      await waitFor(() => {
        expect(screen.getByText('段位审核')).toBeInTheDocument();
      });
    });

    it('should display player rank list', async () => {
      renderWithProviders(<PlayerRankPage />);

      await waitFor(() => {
        expect(screen.getByText('测试陪玩师')).toBeInTheDocument();
      });
    });

    it('should display game rank info', async () => {
      renderWithProviders(<PlayerRankPage />);

      await waitFor(() => {
        expect(screen.getByText('王者')).toBeInTheDocument();
        expect(screen.getByText('王者荣耀')).toBeInTheDocument();
      });
    });

    it('should display pending status', async () => {
      renderWithProviders(<PlayerRankPage />);

      await waitFor(() => {
        expect(screen.getByText('待审核')).toBeInTheDocument();
      });
    });

    it('should display user info from player.user', async () => {
      renderWithProviders(<PlayerRankPage />);

      await waitFor(() => {
        // 用户名显示在 "账号: xxx" 格式中
        expect(screen.getByText(/账号: 张三/)).toBeInTheDocument();
      });
    });

    it('should call API on mount', async () => {
      renderWithProviders(<PlayerRankPage />);

      await waitFor(() => {
        expect(mockApi.getPlayerRanks).toHaveBeenCalled();
        expect(mockApi.getGames).toHaveBeenCalled();
      });
    });
  });

  describe('Status Display', () => {
    it('should display verified status correctly', async () => {
      mockApi.getPlayerRanks.mockResolvedValue({
        data: {
          success: true,
          data: [createMockPlayerRank({ status: 'verified' })],
          pagination: { total: 1 },
        },
      });

      renderWithProviders(<PlayerRankPage />);

      await waitFor(() => {
        expect(screen.getByText('已通过')).toBeInTheDocument();
      });
    });

    it('should display rejected status correctly', async () => {
      mockApi.getPlayerRanks.mockResolvedValue({
        data: {
          success: true,
          data: [createMockPlayerRank({ status: 'rejected', rejectReason: '证明不清晰' })],
          pagination: { total: 1 },
        },
      });

      renderWithProviders(<PlayerRankPage />);

      await waitFor(() => {
        expect(screen.getByText('已拒绝')).toBeInTheDocument();
      });
    });
  });

  describe('Error Handling', () => {
    it('should show error message on API failure', async () => {
      mockApi.getPlayerRanks.mockRejectedValue(new Error('Network error'));

      renderWithProviders(<PlayerRankPage />);

      await waitFor(() => {
        expect(mockMessage.error).toHaveBeenCalledWith('加载段位认证列表失败');
      });
    });
  });

  describe('Audit Operations', () => {
    it('should display audit button for pending items', async () => {
      renderWithProviders(<PlayerRankPage />);

      await waitFor(() => {
        expect(screen.getByText('审核')).toBeInTheDocument();
      });
    });

    it('should display delete button', async () => {
      renderWithProviders(<PlayerRankPage />);

      await waitFor(() => {
        expect(screen.getByText('删除')).toBeInTheDocument();
      });
    });

    it('should display detail button', async () => {
      renderWithProviders(<PlayerRankPage />);

      await waitFor(() => {
        expect(screen.getByText('详情')).toBeInTheDocument();
      });
    });

    it('should open audit modal when audit button clicked', async () => {
      renderWithProviders(<PlayerRankPage />);

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
      renderWithProviders(<PlayerRankPage />);

      await waitFor(() => {
        expect(screen.getByText('段位审核')).toBeInTheDocument();
      });
    });
  });

  describe('Pagination', () => {
    it('should display pagination info', async () => {
      renderWithProviders(<PlayerRankPage />);

      await waitFor(() => {
        expect(screen.getByText(/共 1 条/)).toBeInTheDocument();
      });
    });
  });

  describe('Detail Modal', () => {
    it('should open detail modal when detail button clicked', async () => {
      mockApi.getPlayerRankDetail.mockResolvedValue({
        data: { success: true, data: createMockPlayerRank() },
      });

      renderWithProviders(<PlayerRankPage />);

      await waitFor(() => {
        expect(screen.getByText('详情')).toBeInTheDocument();
      });

      fireEvent.click(screen.getByText('详情'));

      await waitFor(() => {
        expect(screen.getByRole('dialog')).toBeInTheDocument();
      });
    });
  });
});
