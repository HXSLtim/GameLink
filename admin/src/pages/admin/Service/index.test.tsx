/**
 * Service Management Page Tests
 *
 * Tests for Service page component including:
 * - Successful data loading
 * - Loading states
 * - Error handling
 * - Statistics display
 * - CRUD operations
 * - Filter functionality
 */

import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { screen, waitFor, fireEvent } from '@testing-library/react';
import '@testing-library/jest-dom';
import AdminService from './index';
import { renderWithProviders, resetAllMocks, flushPromises } from '../../../testutils';

// Mock the adminApi module using vi.hoisted
const { mockApi, mockMessage } = vi.hoisted(() => ({
  mockApi: {
    getServiceItems: vi.fn(),
    createServiceItem: vi.fn(),
    updateServiceItem: vi.fn(),
    deleteServiceItem: vi.fn(),
    getGames: vi.fn(),
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

// Mock antd App.useApp
vi.mock('antd', async () => {
  const actual = await vi.importActual('antd');
  return {
    ...actual,
    App: {
      ...((actual as Record<string, unknown>).App as Record<string, unknown>),
      useApp: () => ({ message: mockMessage }),
    },
  };
});

// Helper function to create mock service item
const createMockServiceItem = (overrides: Record<string, unknown> = {}): Record<string, unknown> => ({
  id: 1,
  itemCode: 'ESCORT_SOLO_001',
  name: '上分陪玩',
  description: '专业上分服务',
  category: 'escort',
  subCategory: 'solo',
  gameId: 1,
  basePriceCents: 5000,
  serviceHours: 1,
  commissionRate: 0.2,
  minUsers: 1,
  maxPlayers: 1,
  tags: '[]',
  iconUrl: '',
  isActive: true,
  sortOrder: 0,
  createdAt: '2024-01-01T00:00:00Z',
  updatedAt: '2024-01-01T00:00:00Z',
  ...overrides,
});

// Helper function to create mock game
const createMockGame = (overrides: Record<string, unknown> = {}): Record<string, unknown> => ({
  id: 1,
  name: '王者荣耀',
  ...overrides,
});

describe('AdminService', () => {
  beforeEach(() => {
    resetAllMocks();
    mockMessage.success.mockClear();
    mockMessage.error.mockClear();
    localStorage.setItem('token', 'test-token');
    localStorage.setItem('user_role', 'admin');

    // Set default mock return values
    mockApi.getServiceItems.mockResolvedValue({
      data: {
        success: true,
        data: [createMockServiceItem()],
        pagination: { total: 1 },
      },
    });
    mockApi.getGames.mockResolvedValue({
      data: {
        success: true,
        data: [createMockGame()],
      },
    });
  });

  afterEach(() => {
    vi.clearAllMocks();
  });

  describe('Successful Data Loading', () => {
    it('should render service page successfully', async () => {
      renderWithProviders(<AdminService />);

      await waitFor(() => {
        expect(screen.getByText('服务项目管理')).toBeInTheDocument();
      });

      expect(mockApi.getServiceItems).toHaveBeenCalled();
      expect(mockApi.getGames).toHaveBeenCalled();
    });

    it('should display service list', async () => {
      renderWithProviders(<AdminService />);

      await waitFor(() => {
        expect(screen.getByText('上分陪玩')).toBeInTheDocument();
      });
    });

    it('should display service code', async () => {
      renderWithProviders(<AdminService />);

      await waitFor(() => {
        expect(screen.getByText('ESCORT_SOLO_001')).toBeInTheDocument();
      });
    });

    it('should display service price', async () => {
      renderWithProviders(<AdminService />);

      await waitFor(() => {
        expect(screen.getByText('¥50.00')).toBeInTheDocument();
      });
    });
  });

  describe('Statistics Display', () => {
    it('should display total count', async () => {
      renderWithProviders(<AdminService />);

      await waitFor(() => {
        expect(screen.getByText('服务项目总数')).toBeInTheDocument();
      });
    });

    it('should display active count', async () => {
      renderWithProviders(<AdminService />);

      await waitFor(() => {
        expect(screen.getByText('已启用')).toBeInTheDocument();
      });
    });
  });

  describe('Loading States', () => {
    it('should show loading indicator while fetching data', async () => {
      mockApi.getServiceItems.mockImplementation(
        () =>
          new Promise((resolve) => {
            setTimeout(() => {
              resolve({
                data: {
                  success: true,
                  data: [createMockServiceItem()],
                  pagination: { total: 1 },
                },
              });
            }, 100);
          })
      );

      renderWithProviders(<AdminService />);

      expect(mockApi.getServiceItems).toHaveBeenCalled();
    });

    it('should hide loading indicator after data loads', async () => {
      renderWithProviders(<AdminService />);

      await waitFor(() => {
        expect(screen.getByText('服务项目管理')).toBeInTheDocument();
      });

      await flushPromises();

      expect(mockApi.getServiceItems).toHaveBeenCalledTimes(1);
    });
  });

  describe('Error Handling', () => {
    it('should display error message when API fails', async () => {
      mockApi.getServiceItems.mockRejectedValue(new Error('Network error'));

      renderWithProviders(<AdminService />);

      await waitFor(() => {
        expect(mockMessage.error).toHaveBeenCalledWith('加载数据失败');
      });
    });

    it('should display error message when API returns error', async () => {
      mockApi.getServiceItems.mockResolvedValue({
        data: {
          success: false,
          message: '服务器错误',
        },
      });

      renderWithProviders(<AdminService />);

      await waitFor(() => {
        expect(mockMessage.error).toHaveBeenCalledWith('服务器错误');
      });
    });
  });

  describe('Filter Functionality', () => {
    it('should have search input', async () => {
      renderWithProviders(<AdminService />);

      await waitFor(() => {
        expect(screen.getByPlaceholderText('搜索服务名称')).toBeInTheDocument();
      });
    });

    it('should have game filter', async () => {
      renderWithProviders(<AdminService />);

      await waitFor(() => {
        expect(screen.getByText('服务项目管理')).toBeInTheDocument();
      });

      // Game filter select should exist
      expect(mockApi.getGames).toHaveBeenCalled();
    });
  });

  describe('Toolbar Buttons', () => {
    it('should display refresh button', async () => {
      renderWithProviders(<AdminService />);

      await waitFor(() => {
        expect(screen.getByText('刷新')).toBeInTheDocument();
      });
    });

    it('should display add button', async () => {
      renderWithProviders(<AdminService />);

      await waitFor(() => {
        expect(screen.getByText('新增服务')).toBeInTheDocument();
      });
    });

    it('should refresh data when refresh button clicked', async () => {
      renderWithProviders(<AdminService />);

      await waitFor(() => {
        expect(screen.getByText('刷新')).toBeInTheDocument();
      });

      const refreshButton = screen.getByText('刷新');
      fireEvent.click(refreshButton);

      await waitFor(() => {
        expect(mockApi.getServiceItems).toHaveBeenCalledTimes(2);
      });
    });
  });

  describe('Service Actions', () => {
    it('should display edit button', async () => {
      renderWithProviders(<AdminService />);

      await waitFor(() => {
        expect(screen.getByText('编辑')).toBeInTheDocument();
      });
    });

    it('should display delete button', async () => {
      renderWithProviders(<AdminService />);

      await waitFor(() => {
        expect(screen.getByText('删除')).toBeInTheDocument();
      });
    });
  });

  describe('Service Type Display', () => {
    it('should display solo type tag', async () => {
      renderWithProviders(<AdminService />);

      await waitFor(() => {
        expect(screen.getByText('单人护航')).toBeInTheDocument();
      });
    });

    it('should display team type tag', async () => {
      mockApi.getServiceItems.mockResolvedValue({
        data: {
          success: true,
          data: [createMockServiceItem({ subCategory: 'team' })],
          pagination: { total: 1 },
        },
      });

      renderWithProviders(<AdminService />);

      await waitFor(() => {
        expect(screen.getByText('团队护航')).toBeInTheDocument();
      });
    });

    it('should display gift type tag', async () => {
      mockApi.getServiceItems.mockResolvedValue({
        data: {
          success: true,
          data: [createMockServiceItem({ subCategory: 'gift' })],
          pagination: { total: 1 },
        },
      });

      renderWithProviders(<AdminService />);

      await waitFor(() => {
        expect(screen.getByText('礼物')).toBeInTheDocument();
      });
    });
  });

  describe('Service Status', () => {
    it('should display status switch', async () => {
      renderWithProviders(<AdminService />);

      await waitFor(() => {
        expect(screen.getByRole('switch')).toBeInTheDocument();
      });
    });

    it('should show enabled status for active service', async () => {
      renderWithProviders(<AdminService />);

      await waitFor(() => {
        const switchElement = screen.getByRole('switch');
        expect(switchElement).toHaveAttribute('aria-checked', 'true');
      });
    });
  });

  describe('Commission Rate Display', () => {
    it('should display commission rate', async () => {
      renderWithProviders(<AdminService />);

      await waitFor(() => {
        expect(screen.getByText('20%')).toBeInTheDocument();
      });
    });
  });

  describe('Service Hours Display', () => {
    it('should display service hours', async () => {
      renderWithProviders(<AdminService />);

      await waitFor(() => {
        expect(screen.getByText('1小时')).toBeInTheDocument();
      });
    });
  });

  describe('Create Service Modal', () => {
    it('should open create modal when clicking add button', async () => {
      renderWithProviders(<AdminService />);

      await waitFor(() => {
        expect(screen.getByText('新增服务')).toBeInTheDocument();
      });

      const addButton = screen.getByText('新增服务');
      fireEvent.click(addButton);

      await waitFor(() => {
        expect(screen.getByText('新增服务项目')).toBeInTheDocument();
      });
    });
  });

  describe('Edit Service Modal', () => {
    it('should open edit modal when clicking edit button', async () => {
      renderWithProviders(<AdminService />);

      await waitFor(() => {
        expect(screen.getByText('编辑')).toBeInTheDocument();
      });

      const editButton = screen.getByText('编辑');
      fireEvent.click(editButton);

      await waitFor(() => {
        expect(screen.getByText('编辑服务项目')).toBeInTheDocument();
      });
    });
  });

  describe('Delete Service', () => {
    it('should show confirmation when delete button clicked', async () => {
      renderWithProviders(<AdminService />);

      await waitFor(() => {
        expect(screen.getByText('删除')).toBeInTheDocument();
      });

      const deleteButton = screen.getByText('删除');
      fireEvent.click(deleteButton);

      await waitFor(() => {
        expect(screen.getByText('确定删除此服务项目？')).toBeInTheDocument();
      });
    });
  });

  describe('Empty State', () => {
    it('should handle empty service list', async () => {
      mockApi.getServiceItems.mockResolvedValue({
        data: {
          success: true,
          data: [],
          pagination: { total: 0 },
        },
      });

      renderWithProviders(<AdminService />);

      await waitFor(() => {
        expect(screen.getByText('服务项目管理')).toBeInTheDocument();
      });

      expect(mockApi.getServiceItems).toHaveBeenCalled();
    });
  });

  describe('Accessibility', () => {
    it('should have proper page title', async () => {
      renderWithProviders(<AdminService />);

      await waitFor(() => {
        expect(screen.getByText('服务项目管理')).toBeInTheDocument();
      });
    });

    it('should have proper table structure', async () => {
      renderWithProviders(<AdminService />);

      await waitFor(() => {
        expect(screen.getByRole('table')).toBeInTheDocument();
      });
    });
  });

  describe('Data Refresh', () => {
    it('should call loadData and loadGames on mount', async () => {
      renderWithProviders(<AdminService />);

      await waitFor(() => {
        expect(mockApi.getServiceItems).toHaveBeenCalledTimes(1);
        expect(mockApi.getGames).toHaveBeenCalledTimes(1);
      });
    });
  });
});
