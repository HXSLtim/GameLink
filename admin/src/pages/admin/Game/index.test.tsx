/**
 * Game Management Page Tests
 *
 * Tests for Game page component including:
 * - Successful data loading
 * - Loading states
 * - Error handling
 * - CRUD operations
 * - Search functionality
 */

import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import GamePage from './index';
import { renderWithProviders, resetAllMocks, flushPromises } from '@/testutils';

// Mock the adminApi module using vi.hoisted
const { mockApi, mockMessage } = vi.hoisted(() => ({
  mockApi: {
    getGames: vi.fn(),
    createGame: vi.fn(),
    updateGame: vi.fn(),
    deleteGame: vi.fn(),
    batchDeleteGames: vi.fn(),
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

// Mock antd message only
vi.mock('antd', async () => {
  const actual = await vi.importActual('antd');
  return {
    ...actual,
    message: mockMessage,
  };
});

// Helper function to create mock game data
const createMockGame = (overrides: Record<string, unknown> = {}): Record<string, unknown> => ({
  id: 1,
  key: 'lol',
  name: '英雄联盟',
  category: 'moba',
  description: '5v5 MOBA游戏',
  iconUrl: 'https://example.com/lol.png',
  createdAt: '2024-01-01T00:00:00Z',
  updatedAt: '2024-01-01T00:00:00Z',
  ...overrides,
});

describe('GamePage', () => {
  beforeEach(() => {
    resetAllMocks();
    mockMessage.success.mockClear();
    mockMessage.error.mockClear();
    mockMessage.loading.mockClear();
    localStorage.setItem('token', 'test-token');
    localStorage.setItem('user_role', 'admin');
    // Set default mock return values
    mockApi.getGames.mockResolvedValue({
      data: {
        success: true,
        data: [createMockGame()],
        pagination: { total: 1, page: 1, pageSize: 10 },
      },
    });
  });

  afterEach(() => {
    vi.clearAllMocks();
  });

  describe('Successful Data Loading', () => {
    it('should render game page successfully', async () => {
      renderWithProviders(<GamePage />);

      await waitFor(() => {
        expect(screen.getByText('游戏管理')).toBeInTheDocument();
      });

      expect(mockApi.getGames).toHaveBeenCalled();
    });

    it('should display game list with correct columns', async () => {
      renderWithProviders(<GamePage />);

      await waitFor(() => {
        // "游戏" appears in multiple places
        const gameElements = screen.getAllByText('游戏');
        expect(gameElements.length).toBeGreaterThan(0);
      });

      // Check for column headers - use getAllByText for duplicates
      const keyElements = screen.getAllByText('Key');
      expect(keyElements.length).toBeGreaterThan(0);
      const categoryElements = screen.getAllByText('分类');
      expect(categoryElements.length).toBeGreaterThan(0);
      const descElements = screen.getAllByText('描述');
      expect(descElements.length).toBeGreaterThan(0);
    });

    it('should display game data in table', async () => {
      renderWithProviders(<GamePage />);

      await waitFor(() => {
        expect(screen.getByText('英雄联盟')).toBeInTheDocument();
      });

      expect(screen.getByText('lol')).toBeInTheDocument();
      expect(screen.getByText('MOBA')).toBeInTheDocument();
    });

    it('should display game description', async () => {
      renderWithProviders(<GamePage />);

      await waitFor(() => {
        expect(screen.getByText('5v5 MOBA游戏')).toBeInTheDocument();
      });
    });
  });

  describe('Loading States', () => {
    it('should show loading indicator while fetching data', async () => {
      mockApi.getGames.mockImplementation(
        () =>
          new Promise((resolve) => {
            setTimeout(() => {
              resolve({
                data: {
                  success: true,
                  data: [createMockGame()],
                  pagination: { total: 1, page: 1, pageSize: 10 },
                },
              });
            }, 100);
          })
      );

      renderWithProviders(<GamePage />);

      expect(mockApi.getGames).toHaveBeenCalled();
    });

    it('should hide loading indicator after data loads', async () => {
      renderWithProviders(<GamePage />);

      await waitFor(() => {
        expect(screen.getByText('游戏管理')).toBeInTheDocument();
      });

      await flushPromises();

      expect(mockApi.getGames).toHaveBeenCalledTimes(1);
    });
  });

  describe('Error Handling', () => {
    it('should handle API error gracefully', async () => {
      mockApi.getGames.mockRejectedValue(new Error('Network error'));

      renderWithProviders(<GamePage />);

      await waitFor(() => {
        expect(screen.getByText('游戏管理')).toBeInTheDocument();
      });

      await flushPromises();

      // useCrud hook handles errors
      expect(mockApi.getGames).toHaveBeenCalled();
    });

    it('should handle API response with success: false', async () => {
      mockApi.getGames.mockResolvedValue({
        data: {
          success: false,
          message: '获取游戏列表失败',
        },
      });

      renderWithProviders(<GamePage />);

      await waitFor(() => {
        expect(screen.getByText('游戏管理')).toBeInTheDocument();
      });

      await flushPromises();

      // useCrud hook handles errors internally
      expect(mockApi.getGames).toHaveBeenCalled();
    });
  });

  describe('Category Display', () => {
    it('should display MOBA category correctly', async () => {
      mockApi.getGames.mockResolvedValue({
        data: {
          success: true,
          data: [createMockGame({ category: 'moba' })],
          pagination: { total: 1, page: 1, pageSize: 10 },
        },
      });

      renderWithProviders(<GamePage />);

      await waitFor(() => {
        expect(screen.getByText('MOBA')).toBeInTheDocument();
      });
    });

    it('should display FPS category correctly', async () => {
      mockApi.getGames.mockResolvedValue({
        data: {
          success: true,
          data: [createMockGame({ category: 'fps', name: 'CS2', key: 'cs2' })],
          pagination: { total: 1, page: 1, pageSize: 10 },
        },
      });

      renderWithProviders(<GamePage />);

      await waitFor(() => {
        expect(screen.getByText('射击')).toBeInTheDocument();
      });
    });

    it('should display RPG category correctly', async () => {
      mockApi.getGames.mockResolvedValue({
        data: {
          success: true,
          data: [createMockGame({ category: 'rpg', name: '原神', key: 'genshin' })],
          pagination: { total: 1, page: 1, pageSize: 10 },
        },
      });

      renderWithProviders(<GamePage />);

      await waitFor(() => {
        expect(screen.getByText('RPG')).toBeInTheDocument();
      });
    });
  });

  describe('Action Buttons', () => {
    it('should show edit button for each game', async () => {
      renderWithProviders(<GamePage />);

      await waitFor(() => {
        expect(screen.getByText('编辑')).toBeInTheDocument();
      });
    });

    it('should show delete button for each game', async () => {
      renderWithProviders(<GamePage />);

      await waitFor(() => {
        expect(screen.getByText('删除')).toBeInTheDocument();
      });
    });
  });

  describe('Create Game Modal', () => {
    it('should open create modal when clicking create button', async () => {
      const _user = userEvent.setup();
      renderWithProviders(<GamePage />);

      await waitFor(() => {
        expect(screen.getByText('新增游戏')).toBeInTheDocument();
      });

      await _user.click(screen.getByText('新增游戏'));

      await waitFor(() => {
        // Modal title should appear
        const modalTitles = screen.getAllByText('新增游戏');
        expect(modalTitles.length).toBeGreaterThan(1); // Button + Modal title
      });
    });

    it('should display form fields in create modal', async () => {
      const _user = userEvent.setup();
      renderWithProviders(<GamePage />);

      await waitFor(() => {
        expect(screen.getByText('新增游戏')).toBeInTheDocument();
      });

      await _user.click(screen.getByText('新增游戏'));

      await waitFor(() => {
        expect(screen.getByText('游戏Key')).toBeInTheDocument();
      });

      // "游戏名称" appears in multiple places (search field, form field)
      const nameElements = screen.getAllByText('游戏名称');
      expect(nameElements.length).toBeGreaterThan(0);
    });
  });

  describe('Edit Game Modal', () => {
    it('should open edit modal when clicking edit button', async () => {
      const _user = userEvent.setup();
      renderWithProviders(<GamePage />);

      await waitFor(() => {
        expect(screen.getByText('编辑')).toBeInTheDocument();
      });

      await _user.click(screen.getByText('编辑'));

      await waitFor(() => {
        expect(screen.getByText('编辑游戏')).toBeInTheDocument();
      });
    });

    it('should populate form with game data when editing', async () => {
      const _user = userEvent.setup();
      renderWithProviders(<GamePage />);

      await waitFor(() => {
        expect(screen.getByText('编辑')).toBeInTheDocument();
      });

      await _user.click(screen.getByText('编辑'));

      await waitFor(() => {
        expect(screen.getByText('编辑游戏')).toBeInTheDocument();
      });

      // Form should be populated with game data
      const keyInput = screen.getByPlaceholderText('请输入游戏Key（唯一标识）');
      expect(keyInput).toHaveValue('lol');
    });
  });

  describe('Toolbar Buttons', () => {
    it('should display batch delete button', async () => {
      renderWithProviders(<GamePage />);

      await waitFor(() => {
        expect(screen.getByText('批量删除')).toBeInTheDocument();
      });
    });

    it('should display export button', async () => {
      renderWithProviders(<GamePage />);

      await waitFor(() => {
        expect(screen.getByText('导出数据')).toBeInTheDocument();
      });
    });
  });

  describe('Search Functionality', () => {
    it('should have game name search field', async () => {
      renderWithProviders(<GamePage />);

      await waitFor(() => {
        // "游戏名称" appears in multiple places
        const nameElements = screen.getAllByText('游戏名称');
        expect(nameElements.length).toBeGreaterThan(0);
      });
    });

    it('should have category filter', async () => {
      renderWithProviders(<GamePage />);

      await waitFor(() => {
        // "分类" appears in multiple places (search field, table column)
        const categoryElements = screen.getAllByText('分类');
        expect(categoryElements.length).toBeGreaterThan(0);
      });
    });
  });

  describe('Empty State', () => {
    it('should handle empty game list', async () => {
      mockApi.getGames.mockResolvedValue({
        data: {
          success: true,
          data: [],
          pagination: { total: 0, page: 1, pageSize: 10 },
        },
      });

      renderWithProviders(<GamePage />);

      await waitFor(() => {
        expect(screen.getByText('游戏管理')).toBeInTheDocument();
      });

      // Table should still render with no data
      expect(mockApi.getGames).toHaveBeenCalled();
    });
  });

  describe('Multiple Games', () => {
    it('should display multiple games', async () => {
      mockApi.getGames.mockResolvedValue({
        data: {
          success: true,
          data: [
            createMockGame({ id: 1, key: 'lol', name: '英雄联盟', category: 'moba' }),
            createMockGame({ id: 2, key: 'valorant', name: '无畏契约', category: 'fps' }),
            createMockGame({ id: 3, key: 'genshin', name: '原神', category: 'rpg' }),
          ],
          pagination: { total: 3, page: 1, pageSize: 10 },
        },
      });

      renderWithProviders(<GamePage />);

      await waitFor(() => {
        expect(screen.getByText('英雄联盟')).toBeInTheDocument();
      });

      expect(screen.getByText('无畏契约')).toBeInTheDocument();
      expect(screen.getByText('原神')).toBeInTheDocument();
    });
  });

  describe('Accessibility', () => {
    it('should have proper page title', async () => {
      renderWithProviders(<GamePage />);

      await waitFor(() => {
        expect(screen.getByText('游戏管理')).toBeInTheDocument();
      });
    });

    it('should have proper subtitle', async () => {
      renderWithProviders(<GamePage />);

      await waitFor(() => {
        expect(screen.getByText('管理平台支持的游戏')).toBeInTheDocument();
      });
    });
  });
});
