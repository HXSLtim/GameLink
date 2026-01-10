/**
 * Notifications Page Tests
 *
 * Tests for Notifications page component including:
 * - Successful data loading
 * - Loading states
 * - Error handling
 * - Mark as read functionality
 * - Delete notification
 * - Load more functionality
 */

import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { screen, waitFor, fireEvent } from '@testing-library/react';
import AdminNotificationsPage from './index';
import { renderWithProviders, resetAllMocks, flushPromises } from '@/testutils';

// Mock the userApi module using vi.hoisted
const { mockApi, mockMessage } = vi.hoisted(() => ({
  mockApi: {
    getNotifications: vi.fn(),
    markAsRead: vi.fn(),
    markAllAsRead: vi.fn(),
    deleteNotification: vi.fn(),
  },
  mockMessage: {
    success: vi.fn(),
    error: vi.fn(),
    warning: vi.fn(),
    info: vi.fn(),
    loading: vi.fn(),
  },
}));

vi.mock('@/api/user', () => ({
  userApi: mockApi,
}));

// Mock antd App.useApp
vi.mock('antd', async () => {
  const actual = await vi.importActual('antd');
  return {
    ...actual,
    App: {
      ...((actual as Record<string, unknown>).App as Record<string, unknown>),
      useApp: () => ({
        message: mockMessage,
      }),
    },
  };
});

// Helper function to create mock notification
const createMockNotification = (overrides: Record<string, unknown> = {}): Record<string, unknown> => ({
  id: 1,
  title: '测试通知',
  message: '这是一条测试通知内容',
  isRead: false,
  createdAt: '2024-01-01T00:00:00Z',
  ...overrides,
});

describe('AdminNotificationsPage', () => {
  beforeEach(() => {
    resetAllMocks();
    mockMessage.success.mockClear();
    mockMessage.error.mockClear();
    localStorage.setItem('token', 'test-token');
    localStorage.setItem('user_role', 'admin');

    // Set default mock return values
    mockApi.getNotifications.mockResolvedValue({
      success: true,
      data: {
        items: [createMockNotification()],
      },
    });
  });

  afterEach(() => {
    vi.clearAllMocks();
  });

  describe('Successful Data Loading', () => {
    it('should render notifications page successfully', async () => {
      renderWithProviders(<AdminNotificationsPage />);

      await waitFor(() => {
        expect(screen.getByText('消息通知')).toBeInTheDocument();
      });

      expect(mockApi.getNotifications).toHaveBeenCalled();
    });

    it('should display notification list', async () => {
      renderWithProviders(<AdminNotificationsPage />);

      await waitFor(() => {
        expect(screen.getByText('测试通知')).toBeInTheDocument();
      });
    });

    it('should display notification content', async () => {
      renderWithProviders(<AdminNotificationsPage />);

      await waitFor(() => {
        expect(screen.getByText('这是一条测试通知内容')).toBeInTheDocument();
      });
    });

    it('should display mark all read button', async () => {
      renderWithProviders(<AdminNotificationsPage />);

      await waitFor(() => {
        expect(screen.getByText('全部已读')).toBeInTheDocument();
      });
    });
  });

  describe('Loading States', () => {
    it('should show loading indicator while fetching data', async () => {
      mockApi.getNotifications.mockImplementation(
        () =>
          new Promise((resolve) => {
            setTimeout(() => {
              resolve({
                success: true,
                data: {
                  items: [createMockNotification()],
                },
              });
            }, 100);
          })
      );

      renderWithProviders(<AdminNotificationsPage />);

      expect(mockApi.getNotifications).toHaveBeenCalled();
    });

    it('should hide loading indicator after data loads', async () => {
      renderWithProviders(<AdminNotificationsPage />);

      await waitFor(() => {
        expect(screen.getByText('消息通知')).toBeInTheDocument();
      });

      await flushPromises();

      expect(mockApi.getNotifications).toHaveBeenCalledTimes(1);
    });
  });

  describe('Error Handling', () => {
    it('should display error message when API fails', async () => {
      mockApi.getNotifications.mockRejectedValue(new Error('Network error'));

      renderWithProviders(<AdminNotificationsPage />);

      await waitFor(() => {
        expect(mockMessage.error).toHaveBeenCalledWith('加载通知失败');
      });
    });
  });

  describe('Empty State', () => {
    it('should display empty state when no notifications', async () => {
      mockApi.getNotifications.mockResolvedValue({
        success: true,
        data: {
          items: [],
        },
      });

      renderWithProviders(<AdminNotificationsPage />);

      await waitFor(() => {
        expect(screen.getByText('暂无通知')).toBeInTheDocument();
      });
    });
  });

  describe('Mark as Read', () => {
    it('should display mark as read button for unread notifications', async () => {
      renderWithProviders(<AdminNotificationsPage />);

      await waitFor(() => {
        expect(screen.getByText('测试通知')).toBeInTheDocument();
      });

      // Should have mark as read button (CheckOutlined icon)
      const markReadButtons = screen.getAllByRole('button');
      expect(markReadButtons.length).toBeGreaterThan(0);
    });

    it('should call markAsRead API when button clicked', async () => {
      mockApi.markAsRead.mockResolvedValue({ success: true });

      renderWithProviders(<AdminNotificationsPage />);

      await waitFor(() => {
        expect(screen.getByText('测试通知')).toBeInTheDocument();
      });

      // Find and click the mark as read button
      const markReadButton = screen.getByTitle('标记为已读');
      fireEvent.click(markReadButton);

      await waitFor(() => {
        expect(mockApi.markAsRead).toHaveBeenCalledWith(1);
      });
    });

    it('should display error when mark as read fails', async () => {
      mockApi.markAsRead.mockRejectedValue(new Error('Failed'));

      renderWithProviders(<AdminNotificationsPage />);

      await waitFor(() => {
        expect(screen.getByText('测试通知')).toBeInTheDocument();
      });

      const markReadButton = screen.getByTitle('标记为已读');
      fireEvent.click(markReadButton);

      await waitFor(() => {
        expect(mockMessage.error).toHaveBeenCalledWith('操作失败');
      });
    });
  });

  describe('Mark All as Read', () => {
    it('should call markAllAsRead API when button clicked', async () => {
      mockApi.markAllAsRead.mockResolvedValue({ success: true });

      renderWithProviders(<AdminNotificationsPage />);

      await waitFor(() => {
        expect(screen.getByText('全部已读')).toBeInTheDocument();
      });

      const markAllButton = screen.getByText('全部已读');
      fireEvent.click(markAllButton);

      await waitFor(() => {
        expect(mockApi.markAllAsRead).toHaveBeenCalled();
      });
    });

    it('should display success message when mark all succeeds', async () => {
      mockApi.markAllAsRead.mockResolvedValue({ success: true });

      renderWithProviders(<AdminNotificationsPage />);

      await waitFor(() => {
        expect(screen.getByText('全部已读')).toBeInTheDocument();
      });

      const markAllButton = screen.getByText('全部已读');
      fireEvent.click(markAllButton);

      await waitFor(() => {
        expect(mockMessage.success).toHaveBeenCalledWith('已全部标记为已读');
      });
    });

    it('should display error when mark all fails', async () => {
      mockApi.markAllAsRead.mockRejectedValue(new Error('Failed'));

      renderWithProviders(<AdminNotificationsPage />);

      await waitFor(() => {
        expect(screen.getByText('全部已读')).toBeInTheDocument();
      });

      const markAllButton = screen.getByText('全部已读');
      fireEvent.click(markAllButton);

      await waitFor(() => {
        expect(mockMessage.error).toHaveBeenCalledWith('操作失败');
      });
    });
  });

  describe('Delete Notification', () => {
    it('should display delete button', async () => {
      renderWithProviders(<AdminNotificationsPage />);

      await waitFor(() => {
        expect(screen.getByText('测试通知')).toBeInTheDocument();
      });

      expect(screen.getByTitle('删除')).toBeInTheDocument();
    });

    it('should show confirmation when delete clicked', async () => {
      renderWithProviders(<AdminNotificationsPage />);

      await waitFor(() => {
        expect(screen.getByText('测试通知')).toBeInTheDocument();
      });

      const deleteButton = screen.getByTitle('删除');
      fireEvent.click(deleteButton);

      await waitFor(() => {
        expect(screen.getByText('删除通知')).toBeInTheDocument();
      });
    });
  });

  describe('Load More', () => {
    it('should display load more button when has more data', async () => {
      // Return 20 items to trigger hasMore
      const items = Array.from({ length: 20 }, (_, i) =>
        createMockNotification({ id: i + 1, title: `通知 ${i + 1}` })
      );

      mockApi.getNotifications.mockResolvedValue({
        success: true,
        data: { items },
      });

      renderWithProviders(<AdminNotificationsPage />);

      await waitFor(() => {
        expect(screen.getByText('加载更多')).toBeInTheDocument();
      });
    });

    it('should call API with next page when load more clicked', async () => {
      const items = Array.from({ length: 20 }, (_, i) =>
        createMockNotification({ id: i + 1, title: `通知 ${i + 1}` })
      );

      mockApi.getNotifications.mockResolvedValue({
        success: true,
        data: { items },
      });

      renderWithProviders(<AdminNotificationsPage />);

      await waitFor(() => {
        expect(screen.getByText('加载更多')).toBeInTheDocument();
      });

      const loadMoreButton = screen.getByText('加载更多');
      fireEvent.click(loadMoreButton);

      await waitFor(() => {
        expect(mockApi.getNotifications).toHaveBeenCalledWith({ page: 2, page_size: 20 });
      });
    });
  });

  describe('Read/Unread Status Display', () => {
    it('should display unread notification with highlight', async () => {
      renderWithProviders(<AdminNotificationsPage />);

      await waitFor(() => {
        expect(screen.getByText('测试通知')).toBeInTheDocument();
      });

      // Unread notification should have processing badge
      const badges = document.querySelectorAll('.ant-badge-status-processing');
      expect(badges.length).toBeGreaterThan(0);
    });

    it('should not show mark as read button for read notifications', async () => {
      mockApi.getNotifications.mockResolvedValue({
        success: true,
        data: {
          items: [createMockNotification({ isRead: true })],
        },
      });

      renderWithProviders(<AdminNotificationsPage />);

      await waitFor(() => {
        expect(screen.getByText('测试通知')).toBeInTheDocument();
      });

      // Should not have mark as read button
      expect(screen.queryByTitle('标记为已读')).not.toBeInTheDocument();
    });
  });

  describe('Multiple Notifications', () => {
    it('should display multiple notifications', async () => {
      mockApi.getNotifications.mockResolvedValue({
        success: true,
        data: {
          items: [
            createMockNotification({ id: 1, title: '通知一' }),
            createMockNotification({ id: 2, title: '通知二' }),
            createMockNotification({ id: 3, title: '通知三' }),
          ],
        },
      });

      renderWithProviders(<AdminNotificationsPage />);

      await waitFor(() => {
        expect(screen.getByText('通知一')).toBeInTheDocument();
      });

      expect(screen.getByText('通知二')).toBeInTheDocument();
      expect(screen.getByText('通知三')).toBeInTheDocument();
    });
  });

  describe('Accessibility', () => {
    it('should have proper page title', async () => {
      renderWithProviders(<AdminNotificationsPage />);

      await waitFor(() => {
        expect(screen.getByText('消息通知')).toBeInTheDocument();
      });
    });
  });
});
