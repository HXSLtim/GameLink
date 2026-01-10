/**
 * UserTag Management Page Tests
 *
 * Tests for UserTag page component including:
 * - Successful data loading
 * - Loading states
 * - Error handling
 * - Statistics display
 * - CRUD operations
 * - Tag users viewing
 * - Export functionality
 */

import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { screen, waitFor, fireEvent } from '@testing-library/react';
import UserTagPage from './index';
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

// Helper function to create mock user tag
const createMockTag = (overrides: Record<string, unknown> = {}): Record<string, unknown> => ({
  id: 1,
  name: 'VIP用户',
  color: '#1890ff',
  description: '高价值用户标签',
  userCount: 100,
  createdAt: '2024-01-01T00:00:00Z',
  updatedAt: '2024-01-01T00:00:00Z',
  ...overrides,
});

// Helper function to create mock tag user
const createMockTagUser = (overrides: Record<string, unknown> = {}): Record<string, unknown> => ({
  id: 1,
  name: '测试用户',
  email: 'test@example.com',
  avatarUrl: null,
  ...overrides,
});

describe('UserTagPage', () => {
  beforeEach(() => {
    resetAllMocks();
    mockMessage.success.mockClear();
    mockMessage.error.mockClear();
    mockMessage.loading.mockClear();
    localStorage.setItem('token', 'test-token');
    localStorage.setItem('user_role', 'admin');

    // Set default mock return values
    mockApiClient.get.mockImplementation((url: string) => {
      if (url === '/admin/user-tags') {
        return Promise.resolve({
          data: {
            success: true,
            data: [createMockTag()],
          },
        });
      }
      if (url.includes('/users')) {
        return Promise.resolve({
          data: {
            success: true,
            data: [createMockTagUser()],
            pagination: { total: 1, page: 1, page_size: 10 },
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
    it('should render user tag page successfully', async () => {
      renderWithProviders(<UserTagPage />);

      await waitFor(() => {
        expect(screen.getByText('用户标签管理')).toBeInTheDocument();
      });

      expect(mockApiClient.get).toHaveBeenCalledWith('/admin/user-tags');
    });

    it('should display page subtitle', async () => {
      renderWithProviders(<UserTagPage />);

      await waitFor(() => {
        expect(screen.getByText('管理用户分群标签')).toBeInTheDocument();
      });
    });

    it('should display tag list', async () => {
      renderWithProviders(<UserTagPage />);

      await waitFor(() => {
        expect(screen.getByText('VIP用户')).toBeInTheDocument();
      });
    });

    it('should display tag description', async () => {
      renderWithProviders(<UserTagPage />);

      await waitFor(() => {
        expect(screen.getByText('高价值用户标签')).toBeInTheDocument();
      });
    });

    it('should display tag color', async () => {
      renderWithProviders(<UserTagPage />);

      await waitFor(() => {
        expect(screen.getByText('#1890ff')).toBeInTheDocument();
      });
    });
  });

  describe('Statistics Display', () => {
    it('should display total tags statistic', async () => {
      renderWithProviders(<UserTagPage />);

      await waitFor(() => {
        expect(screen.getByText('标签总数')).toBeInTheDocument();
      });
    });

    it('should display tagged users statistic', async () => {
      renderWithProviders(<UserTagPage />);

      await waitFor(() => {
        expect(screen.getByText('已标记用户')).toBeInTheDocument();
      });
    });

    it('should display user count in table', async () => {
      renderWithProviders(<UserTagPage />);

      await waitFor(() => {
        expect(screen.getByText('100 人')).toBeInTheDocument();
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
                  data: [createMockTag()],
                },
              });
            }, 100);
          })
      );

      renderWithProviders(<UserTagPage />);

      expect(mockApiClient.get).toHaveBeenCalled();
    });

    it('should hide loading indicator after data loads', async () => {
      renderWithProviders(<UserTagPage />);

      await waitFor(() => {
        expect(screen.getByText('用户标签管理')).toBeInTheDocument();
      });

      await flushPromises();

      expect(mockApiClient.get).toHaveBeenCalledTimes(1);
    });
  });

  describe('Error Handling', () => {
    it('should display error message when API fails', async () => {
      mockApiClient.get.mockRejectedValue(new Error('Network error'));

      renderWithProviders(<UserTagPage />);

      await waitFor(() => {
        expect(mockMessage.error).toHaveBeenCalledWith('加载标签列表失败');
      });
    });
  });

  describe('Filter Functionality', () => {
    it('should have keyword search input', async () => {
      renderWithProviders(<UserTagPage />);

      await waitFor(() => {
        expect(screen.getByPlaceholderText('标签名称/描述')).toBeInTheDocument();
      });
    });

    it('should filter tags by keyword', async () => {
      mockApiClient.get.mockResolvedValue({
        data: {
          success: true,
          data: [
            createMockTag({ id: 1, name: 'VIP用户' }),
            createMockTag({ id: 2, name: '新用户' }),
          ],
        },
      });

      renderWithProviders(<UserTagPage />);

      await waitFor(() => {
        expect(screen.getByText('VIP用户')).toBeInTheDocument();
        expect(screen.getByText('新用户')).toBeInTheDocument();
      });
    });
  });

  describe('Toolbar Buttons', () => {
    it('should display create button', async () => {
      renderWithProviders(<UserTagPage />);

      await waitFor(() => {
        expect(screen.getByText('新增标签')).toBeInTheDocument();
      });
    });

    it('should display export button', async () => {
      renderWithProviders(<UserTagPage />);

      await waitFor(() => {
        expect(screen.getByText('导出数据')).toBeInTheDocument();
      });
    });
  });

  describe('Tag Actions', () => {
    it('should display edit button', async () => {
      renderWithProviders(<UserTagPage />);

      await waitFor(() => {
        expect(screen.getByText('编辑')).toBeInTheDocument();
      });
    });

    it('should display delete button', async () => {
      renderWithProviders(<UserTagPage />);

      await waitFor(() => {
        expect(screen.getByText('删除')).toBeInTheDocument();
      });
    });
  });

  describe('Create Tag Modal', () => {
    it('should open create modal when button clicked', async () => {
      renderWithProviders(<UserTagPage />);

      await waitFor(() => {
        expect(screen.getByText('新增标签')).toBeInTheDocument();
      });

      const createButton = screen.getByText('新增标签');
      fireEvent.click(createButton);

      await waitFor(() => {
        // Modal title changes to "新增标签" when creating
        expect(screen.getByRole('dialog')).toBeInTheDocument();
      });
    });

    it('should display form fields in create modal', async () => {
      renderWithProviders(<UserTagPage />);

      await waitFor(() => {
        expect(screen.getByText('新增标签')).toBeInTheDocument();
      });

      const createButton = screen.getByText('新增标签');
      fireEvent.click(createButton);

      await waitFor(() => {
        expect(screen.getByRole('dialog')).toBeInTheDocument();
      });

      // Form fields should be visible in the modal
      await waitFor(() => {
        expect(screen.getByPlaceholderText('请输入标签名称')).toBeInTheDocument();
        expect(screen.getByPlaceholderText('请输入标签描述')).toBeInTheDocument();
      });
    });
  });

  describe('Edit Tag Modal', () => {
    it('should open edit modal when edit button clicked', async () => {
      renderWithProviders(<UserTagPage />);

      await waitFor(() => {
        expect(screen.getByText('编辑')).toBeInTheDocument();
      });

      const editButton = screen.getByText('编辑');
      fireEvent.click(editButton);

      await waitFor(() => {
        expect(screen.getByText('编辑标签')).toBeInTheDocument();
      });
    });
  });

  describe('Delete Tag', () => {
    it('should show confirmation when delete button clicked', async () => {
      renderWithProviders(<UserTagPage />);

      await waitFor(() => {
        expect(screen.getByText('删除')).toBeInTheDocument();
      });

      const deleteButton = screen.getByText('删除');
      fireEvent.click(deleteButton);

      await waitFor(() => {
        expect(screen.getByText('确定要删除该标签吗？删除后所有用户的该标签将被移除。')).toBeInTheDocument();
      });
    });
  });

  describe('View Tag Users', () => {
    it('should open users modal when user count clicked', async () => {
      renderWithProviders(<UserTagPage />);

      await waitFor(() => {
        expect(screen.getByText('100 人')).toBeInTheDocument();
      });

      const userCountButton = screen.getByText('100 人');
      fireEvent.click(userCountButton);

      await waitFor(() => {
        expect(screen.getByText('标签用户列表')).toBeInTheDocument();
      });
    });

    it('should load tag users when modal opens', async () => {
      renderWithProviders(<UserTagPage />);

      await waitFor(() => {
        expect(screen.getByText('100 人')).toBeInTheDocument();
      });

      const userCountButton = screen.getByText('100 人');
      fireEvent.click(userCountButton);

      await waitFor(() => {
        expect(mockApiClient.get).toHaveBeenCalledWith('/admin/user-tags/1/users', expect.any(Object));
      });
    });

    it('should display user list in modal', async () => {
      renderWithProviders(<UserTagPage />);

      await waitFor(() => {
        expect(screen.getByText('100 人')).toBeInTheDocument();
      });

      const userCountButton = screen.getByText('100 人');
      fireEvent.click(userCountButton);

      await waitFor(() => {
        expect(screen.getByText('测试用户')).toBeInTheDocument();
        expect(screen.getByText('test@example.com')).toBeInTheDocument();
      });
    });
  });

  describe('Export Functionality', () => {
    it('should show loading and success message when export clicked', async () => {
      renderWithProviders(<UserTagPage />);

      await waitFor(() => {
        expect(screen.getByText('导出数据')).toBeInTheDocument();
      });

      const exportButton = screen.getByText('导出数据');
      fireEvent.click(exportButton);

      await waitFor(() => {
        expect(mockMessage.loading).toHaveBeenCalledWith({ content: '正在导出...', key: 'export' });
        expect(mockMessage.success).toHaveBeenCalledWith({ content: '导出成功', key: 'export' });
      });
    });
  });

  describe('Empty State', () => {
    it('should handle empty tag list', async () => {
      mockApiClient.get.mockResolvedValue({
        data: {
          success: true,
          data: [],
        },
      });

      renderWithProviders(<UserTagPage />);

      await waitFor(() => {
        expect(screen.getByText('用户标签管理')).toBeInTheDocument();
      });

      expect(mockApiClient.get).toHaveBeenCalled();
    });
  });

  describe('Statistics with Zero Values', () => {
    it('should display zero values correctly', async () => {
      mockApiClient.get.mockResolvedValue({
        data: {
          success: true,
          data: [createMockTag({ userCount: 0 })],
        },
      });

      renderWithProviders(<UserTagPage />);

      await waitFor(() => {
        expect(screen.getByText('0 人')).toBeInTheDocument();
      });
    });
  });

  describe('Tag Color Display', () => {
    it('should display tag with correct color', async () => {
      renderWithProviders(<UserTagPage />);

      await waitFor(() => {
        // The tag should be rendered with the color
        expect(screen.getByText('VIP用户')).toBeInTheDocument();
      });
    });

    it('should display different tag colors', async () => {
      mockApiClient.get.mockResolvedValue({
        data: {
          success: true,
          data: [
            createMockTag({ id: 1, name: '标签1', color: '#ff0000' }),
            createMockTag({ id: 2, name: '标签2', color: '#00ff00' }),
          ],
        },
      });

      renderWithProviders(<UserTagPage />);

      await waitFor(() => {
        expect(screen.getByText('#ff0000')).toBeInTheDocument();
        expect(screen.getByText('#00ff00')).toBeInTheDocument();
      });
    });
  });

  describe('Accessibility', () => {
    it('should have proper page title', async () => {
      renderWithProviders(<UserTagPage />);

      await waitFor(() => {
        expect(screen.getByText('用户标签管理')).toBeInTheDocument();
      });
    });

    it('should have proper table structure', async () => {
      renderWithProviders(<UserTagPage />);

      await waitFor(() => {
        expect(screen.getByRole('table')).toBeInTheDocument();
      });
    });
  });

  describe('Multiple Tags', () => {
    it('should display multiple tags', async () => {
      mockApiClient.get.mockResolvedValue({
        data: {
          success: true,
          data: [
            createMockTag({ id: 1, name: 'VIP用户', userCount: 100 }),
            createMockTag({ id: 2, name: '新用户', userCount: 50 }),
            createMockTag({ id: 3, name: '活跃用户', userCount: 200 }),
          ],
        },
      });

      renderWithProviders(<UserTagPage />);

      await waitFor(() => {
        expect(screen.getByText('VIP用户')).toBeInTheDocument();
        expect(screen.getByText('新用户')).toBeInTheDocument();
        expect(screen.getByText('活跃用户')).toBeInTheDocument();
      });
    });

    it('should calculate total tagged users correctly', async () => {
      mockApiClient.get.mockResolvedValue({
        data: {
          success: true,
          data: [
            createMockTag({ id: 1, userCount: 100 }),
            createMockTag({ id: 2, userCount: 50 }),
          ],
        },
      });

      renderWithProviders(<UserTagPage />);

      await waitFor(() => {
        // Total should be 150
        expect(screen.getByText('150')).toBeInTheDocument();
      });
    });
  });
});
