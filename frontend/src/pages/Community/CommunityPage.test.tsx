import { describe, it, expect, vi, beforeEach } from 'vitest';
import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import { CommunityPage } from './CommunityPage';
import { feedApi } from '../../services/api/feed';
import { notificationApi } from '../../services/api/notification';

vi.mock('../../services/api/feed', () => ({
  feedApi: {
    list: vi.fn(),
    create: vi.fn(),
    report: vi.fn(),
  },
}));

vi.mock('../../services/api/notification', () => ({
  notificationApi: {
    list: vi.fn(),
    markRead: vi.fn(),
    unreadCount: vi.fn(),
  },
}));

describe('CommunityPage', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    vi.mocked(feedApi.list).mockResolvedValue({ items: [], nextCursor: undefined });
    vi.mocked(notificationApi.unreadCount).mockResolvedValue({ unread: 0 });
  });

  it('loads feeds on mount', async () => {
    render(<CommunityPage />);
    await waitFor(() => {
      expect(feedApi.list).toHaveBeenCalled();
    });
  });

  it('opens publish modal', async () => {
    render(<CommunityPage />);
    const publishButton = screen.getByText('发布动态');
    fireEvent.click(publishButton);
    expect(await screen.findByText('动态内容')).toBeInTheDocument();
  });

  it('opens notification modal and fetches notifications', async () => {
    vi.mocked(notificationApi.list).mockResolvedValue({ items: [], page: 1, pageSize: 10, total: 0, unreadCount: 0 });
    render(<CommunityPage />);
    const button = await screen.findByText('通知中心');
    fireEvent.click(button);
    await waitFor(() => {
      expect(notificationApi.list).toHaveBeenCalled();
    });
  });
});
