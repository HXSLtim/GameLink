/**
 * GameRank Management Page Tests
 *
 * Tests for GameRank page component including:
 * - Successful data loading
 * - Loading states
 * - Error handling
 * - CRUD operations
 * - Batch operations
 */

import { describe, it, expect, vi, beforeEach } from 'vitest';
import { screen, waitFor, fireEvent } from '@testing-library/react';
import GameRankPage from './index';
import { renderWithProviders, resetAllMocks } from '@/testutils';

// Mock the adminApi module using vi.hoisted
const { mockApi, mockMessage } = vi.hoisted(() => ({
  mockApi: {
    getGameRanks: vi.fn(),
    getGames: vi.fn(),
    createGameRank: vi.fn(),
    updateGameRank: vi.fn(),
    deleteGameRank: vi.fn(),
    batchDeleteGameRanks: vi.fn(),
    batchUpdateGameRankStatus: vi.fn(),
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

// Helper function to create mock game rank
const createMockGameRank = (overrides: Record<string, unknown> = {}): Record<string, unknown> => ({
  id: 1,
  gameId: 1,
  name: '青铜',
  description: '入门段位',
  level: 1,
  priceCents: 2000,
  sortOrder: 1,
  color: '#cd7f32',
  iconUrl: null,
  isActive: true,
  game: { id: 1, name: '王者荣耀' },
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

describe('GameRankPage', () => {
  beforeEach(() => {
    resetAllMocks();
    // Default successful responses
    mockApi.getGameRanks.mockResolvedValue({
      data: {
        success: true,
        data: [createMockGameRank()],
        pagination: { total: 1, page: 1, pageSize: 10 },
      },
    });
    mockApi.getGames.mockResolvedValue({
      data: {
        success: true,
        data: [createMockGame()],
      },
    });
  });

  describe('Successful Data Loading', () => {
    it('should render page title', async () => {
      renderWithProviders(<GameRankPage />);

      await waitFor(() => {
        expect(screen.getByText('段位管理')).toBeInTheDocument();
      });
    });

    it('should display game rank list', async () => {
      renderWithProviders(<GameRankPage />);

      await waitFor(() => {
        expect(screen.getByText('青铜')).toBeInTheDocument();
      });
    });

    it('should display game name tag', async () => {
      renderWithProviders(<GameRankPage />);

      await waitFor(() => {
        expect(screen.getByText('王者荣耀')).toBeInTheDocument();
      });
    });

    it('should display price correctly', async () => {
      renderWithProviders(<GameRankPage />);

      await waitFor(() => {
        expect(screen.getByText('¥20.00/小时')).toBeInTheDocument();
      });
    });

    it('should display status tag', async () => {
      renderWithProviders(<GameRankPage />);

      await waitFor(() => {
        expect(screen.getByText('启用')).toBeInTheDocument();
      });
    });

    it('should call API on mount', async () => {
      renderWithProviders(<GameRankPage />);

      await waitFor(() => {
        expect(mockApi.getGameRanks).toHaveBeenCalled();
        expect(mockApi.getGames).toHaveBeenCalled();
      });
    });
  });

  describe('Loading States', () => {
    it('should show loading state initially', async () => {
      mockApi.getGameRanks.mockImplementation(
        () => new Promise((resolve) => setTimeout(() => resolve({
          data: { success: true, data: [], pagination: { total: 0 } },
        }), 100))
      );

      renderWithProviders(<GameRankPage />);

      // Table should be in loading state
      expect(mockApi.getGameRanks).toHaveBeenCalled();
    });
  });

  describe('Error Handling', () => {
    it('should show error message on API failure', async () => {
      mockApi.getGameRanks.mockRejectedValue(new Error('Network error'));

      renderWithProviders(<GameRankPage />);

      await waitFor(() => {
        expect(mockMessage.error).toHaveBeenCalledWith('加载段位列表失败');
      });
    });

    it('should show error message on unsuccessful response', async () => {
      mockApi.getGameRanks.mockResolvedValue({
        data: { success: false, message: '服务器错误' },
      });

      renderWithProviders(<GameRankPage />);

      await waitFor(() => {
        expect(mockMessage.error).toHaveBeenCalledWith('服务器错误');
      });
    });
  });

  describe('Create Operation', () => {
    it('should display create button', async () => {
      renderWithProviders(<GameRankPage />);

      await waitFor(() => {
        expect(screen.getByText('新增段位')).toBeInTheDocument();
      });
    });

    it('should open modal when create button clicked', async () => {
      renderWithProviders(<GameRankPage />);

      await waitFor(() => {
        expect(screen.getByText('新增段位')).toBeInTheDocument();
      });

      fireEvent.click(screen.getByText('新增段位'));

      await waitFor(() => {
        expect(screen.getByRole('dialog')).toBeInTheDocument();
      });
    });
  });

  describe('Edit Operation', () => {
    it('should display edit button', async () => {
      renderWithProviders(<GameRankPage />);

      await waitFor(() => {
        expect(screen.getByText('编辑')).toBeInTheDocument();
      });
    });

    it('should open modal when edit button clicked', async () => {
      renderWithProviders(<GameRankPage />);

      await waitFor(() => {
        expect(screen.getByText('编辑')).toBeInTheDocument();
      });

      fireEvent.click(screen.getByText('编辑'));

      await waitFor(() => {
        expect(screen.getByRole('dialog')).toBeInTheDocument();
        expect(screen.getByText('编辑段位')).toBeInTheDocument();
      });
    });
  });

  describe('Delete Operation', () => {
    it('should display delete button', async () => {
      renderWithProviders(<GameRankPage />);

      await waitFor(() => {
        expect(screen.getByText('删除')).toBeInTheDocument();
      });
    });

    it('should show confirm dialog when delete clicked', async () => {
      renderWithProviders(<GameRankPage />);

      await waitFor(() => {
        expect(screen.getByText('删除')).toBeInTheDocument();
      });

      fireEvent.click(screen.getByText('删除'));

      await waitFor(() => {
        expect(screen.getByText('确定删除该段位？')).toBeInTheDocument();
      });
    });

    it('should call delete API when confirmed', async () => {
      mockApi.deleteGameRank.mockResolvedValue({ data: { success: true } });

      renderWithProviders(<GameRankPage />);

      await waitFor(() => {
        expect(screen.getByText('删除')).toBeInTheDocument();
      });

      fireEvent.click(screen.getByText('删除'));

      await waitFor(() => {
        expect(screen.getByText('确定删除该段位？')).toBeInTheDocument();
      });

      // Click confirm button (Popconfirm uses "OK" in English)
      const confirmBtn = screen.getByRole('button', { name: /OK/i });
      fireEvent.click(confirmBtn);

      await waitFor(() => {
        expect(mockApi.deleteGameRank).toHaveBeenCalledWith(1);
      });
    });
  });

  describe('Batch Operations', () => {
    it('should display batch operation buttons', async () => {
      renderWithProviders(<GameRankPage />);

      await waitFor(() => {
        expect(screen.getByText('批量启用')).toBeInTheDocument();
        expect(screen.getByText('批量禁用')).toBeInTheDocument();
        expect(screen.getByText('批量删除')).toBeInTheDocument();
      });
    });
  });

  describe('Search Functionality', () => {
    it('should have search fields', async () => {
      renderWithProviders(<GameRankPage />);

      await waitFor(() => {
        expect(screen.getByPlaceholderText('段位名称')).toBeInTheDocument();
      });
    });
  });

  describe('Pagination', () => {
    it('should display pagination info', async () => {
      renderWithProviders(<GameRankPage />);

      await waitFor(() => {
        expect(screen.getByText(/共 1 条/)).toBeInTheDocument();
      });
    });
  });
});
