/**
 * Notification Store Tests
 * Tests for notification management, read status, and pagination
 */

import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { useNotificationStore, NotificationType } from '../notification-store';

// Mock http module
vi.mock('@/lib/http', () => ({
    http: {
        get: vi.fn(),
        post: vi.fn(),
        put: vi.fn(),
        delete: vi.fn(),
    },
}));

import { http } from '@/lib/http';

const mockHttp = http as unknown as {
    get: ReturnType<typeof vi.fn>;
    post: ReturnType<typeof vi.fn>;
    put: ReturnType<typeof vi.fn>;
    delete: ReturnType<typeof vi.fn>;
};

describe('Notification Store', () => {
    beforeEach(() => {
        vi.clearAllMocks();
        // Reset store state
        useNotificationStore.setState({
            notifications: [],
            unreadCount: 0,
            loading: false,
            error: null,
            hasMore: true,
            page: 1,
        });
    });

    afterEach(() => {
        vi.restoreAllMocks();
    });

    describe('Initial State', () => {
        it('should have correct initial state', () => {
            const state = useNotificationStore.getState();

            expect(state.notifications).toEqual([]);
            expect(state.unreadCount).toBe(0);
            expect(state.loading).toBe(false);
            expect(state.error).toBeNull();
            expect(state.hasMore).toBe(true);
            expect(state.page).toBe(1);
        });
    });

    describe('fetchNotifications', () => {
        it('should fetch notifications successfully', async () => {
            const mockNotifications = {
                items: [
                    {
                        id: '1',
                        title: '订单完成',
                        message: '您的订单已完成',
                        type: NotificationType.ORDER,
                        read: false,
                        createdAt: '2024-01-01T00:00:00Z',
                    },
                    {
                        id: '2',
                        title: '系统通知',
                        message: '系统维护通知',
                        type: NotificationType.SYSTEM,
                        read: true,
                        createdAt: '2024-01-02T00:00:00Z',
                    },
                ],
                total: 2,
                page: 1,
                pageSize: 20,
            };

            mockHttp.get.mockResolvedValueOnce(mockNotifications);

            await useNotificationStore.getState().fetchNotifications(true);

            const state = useNotificationStore.getState();
            expect(state.notifications).toHaveLength(2);
            expect(state.notifications[0].title).toBe('订单完成');
            expect(state.unreadCount).toBe(1); // Only 1 unread
            expect(state.loading).toBe(false);
            expect(state.error).toBeNull();
        });

        it('should reset notifications when refresh is true', async () => {
            // Set initial notifications
            useNotificationStore.setState({
                notifications: [
                    {
                        id: 'old',
                        title: 'Old',
                        message: 'Old message',
                        type: NotificationType.SYSTEM,
                        read: true,
                        createdAt: '2024-01-01',
                    },
                ],
                page: 3,
            });

            mockHttp.get.mockResolvedValueOnce({
                items: [
                    {
                        id: 'new',
                        title: 'New',
                        message: 'New message',
                        type: NotificationType.SYSTEM,
                        read: false,
                        createdAt: '2024-01-02',
                    },
                ],
                total: 1,
                page: 1,
                pageSize: 20,
            });

            await useNotificationStore.getState().fetchNotifications(true);

            const state = useNotificationStore.getState();
            expect(state.notifications).toHaveLength(1);
            expect(state.notifications[0].id).toBe('new');
        });

        it('should append notifications when refresh is false', async () => {
            useNotificationStore.setState({
                notifications: [
                    {
                        id: '1',
                        title: 'First',
                        message: 'First message',
                        type: NotificationType.SYSTEM,
                        read: true,
                        createdAt: '2024-01-01',
                    },
                ],
                page: 1,
            });

            mockHttp.get.mockResolvedValueOnce({
                items: [
                    {
                        id: '2',
                        title: 'Second',
                        message: 'Second message',
                        type: NotificationType.SYSTEM,
                        read: false,
                        createdAt: '2024-01-02',
                    },
                ],
                total: 2,
                page: 1,
                pageSize: 20,
            });

            await useNotificationStore.getState().fetchNotifications(false);

            const state = useNotificationStore.getState();
            expect(state.notifications).toHaveLength(2);
        });

        it('should handle fetch error', async () => {
            mockHttp.get.mockRejectedValueOnce(new Error('Network error'));

            await useNotificationStore.getState().fetchNotifications();

            const state = useNotificationStore.getState();
            expect(state.loading).toBe(false);
            expect(state.error).toBe('Network error');
        });

        it('should set hasMore based on items length', async () => {
            // Less than 20 items means no more
            mockHttp.get.mockResolvedValueOnce({
                items: Array(10).fill({
                    id: '1',
                    title: 'Test',
                    message: 'Test',
                    type: NotificationType.SYSTEM,
                    read: true,
                    createdAt: '2024-01-01',
                }),
                total: 10,
                page: 1,
                pageSize: 20,
            });

            await useNotificationStore.getState().fetchNotifications(true);

            expect(useNotificationStore.getState().hasMore).toBe(false);
        });
    });

    describe('fetchUnreadCount', () => {
        it('should fetch unread count', async () => {
            mockHttp.get.mockResolvedValueOnce({ count: 5 });

            await useNotificationStore.getState().fetchUnreadCount();

            expect(useNotificationStore.getState().unreadCount).toBe(5);
        });

        it('should handle error silently', async () => {
            const consoleSpy = vi.spyOn(console, 'error').mockImplementation(() => { });
            mockHttp.get.mockRejectedValueOnce(new Error('API error'));

            await useNotificationStore.getState().fetchUnreadCount();

            expect(consoleSpy).toHaveBeenCalled();
            consoleSpy.mockRestore();
        });
    });

    describe('markAsRead', () => {
        it('should mark notification as read (optimistic update)', async () => {
            useNotificationStore.setState({
                notifications: [
                    {
                        id: '1',
                        title: 'Test',
                        message: 'Test message',
                        type: NotificationType.ORDER,
                        read: false,
                        createdAt: '2024-01-01',
                    },
                ],
                unreadCount: 1,
            });

            mockHttp.put.mockResolvedValueOnce({});

            await useNotificationStore.getState().markAsRead('1');

            const state = useNotificationStore.getState();
            expect(state.notifications[0].read).toBe(true);
            expect(state.unreadCount).toBe(0);
        });

        it('should not update already read notification', async () => {
            useNotificationStore.setState({
                notifications: [
                    {
                        id: '1',
                        title: 'Test',
                        message: 'Test message',
                        type: NotificationType.ORDER,
                        read: true,
                        createdAt: '2024-01-01',
                    },
                ],
                unreadCount: 0,
            });

            await useNotificationStore.getState().markAsRead('1');

            expect(mockHttp.put).not.toHaveBeenCalled();
        });

        it('should rollback on error', async () => {
            useNotificationStore.setState({
                notifications: [
                    {
                        id: '1',
                        title: 'Test',
                        message: 'Test message',
                        type: NotificationType.ORDER,
                        read: false,
                        createdAt: '2024-01-01',
                    },
                ],
                unreadCount: 1,
            });

            mockHttp.put.mockRejectedValueOnce(new Error('API error'));

            await useNotificationStore.getState().markAsRead('1');

            const state = useNotificationStore.getState();
            expect(state.notifications[0].read).toBe(false);
            expect(state.unreadCount).toBe(1);
            expect(state.error).toBe('API error');
        });
    });

    describe('markAllAsRead', () => {
        it('should mark all notifications as read', async () => {
            useNotificationStore.setState({
                notifications: [
                    {
                        id: '1',
                        title: 'Test 1',
                        message: 'Message 1',
                        type: NotificationType.ORDER,
                        read: false,
                        createdAt: '2024-01-01',
                    },
                    {
                        id: '2',
                        title: 'Test 2',
                        message: 'Message 2',
                        type: NotificationType.SYSTEM,
                        read: false,
                        createdAt: '2024-01-02',
                    },
                ],
                unreadCount: 2,
            });

            mockHttp.put.mockResolvedValueOnce({});

            await useNotificationStore.getState().markAllAsRead();

            const state = useNotificationStore.getState();
            expect(state.notifications.every(n => n.read)).toBe(true);
            expect(state.unreadCount).toBe(0);
        });

        it('should rollback on error', async () => {
            const originalNotifications = [
                {
                    id: '1',
                    title: 'Test',
                    message: 'Message',
                    type: NotificationType.ORDER,
                    read: false,
                    createdAt: '2024-01-01',
                },
            ];

            useNotificationStore.setState({
                notifications: originalNotifications,
                unreadCount: 1,
            });

            mockHttp.put.mockRejectedValueOnce(new Error('API error'));

            await useNotificationStore.getState().markAllAsRead();

            const state = useNotificationStore.getState();
            expect(state.notifications[0].read).toBe(false);
            expect(state.unreadCount).toBe(1);
        });
    });

    describe('deleteNotification', () => {
        it('should delete notification (optimistic update)', async () => {
            useNotificationStore.setState({
                notifications: [
                    {
                        id: '1',
                        title: 'Test',
                        message: 'Message',
                        type: NotificationType.ORDER,
                        read: false,
                        createdAt: '2024-01-01',
                    },
                ],
                unreadCount: 1,
            });

            mockHttp.delete.mockResolvedValueOnce({});

            await useNotificationStore.getState().deleteNotification('1');

            const state = useNotificationStore.getState();
            expect(state.notifications).toHaveLength(0);
            expect(state.unreadCount).toBe(0);
        });

        it('should not decrement unread count for read notification', async () => {
            useNotificationStore.setState({
                notifications: [
                    {
                        id: '1',
                        title: 'Test',
                        message: 'Message',
                        type: NotificationType.ORDER,
                        read: true,
                        createdAt: '2024-01-01',
                    },
                ],
                unreadCount: 0,
            });

            mockHttp.delete.mockResolvedValueOnce({});

            await useNotificationStore.getState().deleteNotification('1');

            expect(useNotificationStore.getState().unreadCount).toBe(0);
        });

        it('should rollback on error', async () => {
            const notification = {
                id: '1',
                title: 'Test',
                message: 'Message',
                type: NotificationType.ORDER,
                read: false,
                createdAt: '2024-01-01',
            };

            useNotificationStore.setState({
                notifications: [notification],
                unreadCount: 1,
            });

            mockHttp.delete.mockRejectedValueOnce(new Error('API error'));

            await useNotificationStore.getState().deleteNotification('1');

            const state = useNotificationStore.getState();
            expect(state.notifications).toHaveLength(1);
        });
    });

    describe('addNotification', () => {
        it('should add notification to the beginning', () => {
            useNotificationStore.setState({
                notifications: [
                    {
                        id: '1',
                        title: 'Old',
                        message: 'Old message',
                        type: NotificationType.SYSTEM,
                        read: true,
                        createdAt: '2024-01-01',
                    },
                ],
                unreadCount: 0,
            });

            const newNotification = {
                id: '2',
                title: 'New',
                message: 'New message',
                type: NotificationType.ORDER,
                read: false,
                createdAt: '2024-01-02',
            };

            useNotificationStore.getState().addNotification(newNotification);

            const state = useNotificationStore.getState();
            expect(state.notifications).toHaveLength(2);
            expect(state.notifications[0].id).toBe('2');
            expect(state.unreadCount).toBe(1);
        });

        it('should not increment unread count for read notification', () => {
            useNotificationStore.setState({
                notifications: [],
                unreadCount: 0,
            });

            const readNotification = {
                id: '1',
                title: 'Read',
                message: 'Already read',
                type: NotificationType.SYSTEM,
                read: true,
                createdAt: '2024-01-01',
            };

            useNotificationStore.getState().addNotification(readNotification);

            expect(useNotificationStore.getState().unreadCount).toBe(0);
        });
    });

    describe('loadMore', () => {
        it('should load more notifications', async () => {
            useNotificationStore.setState({
                notifications: [
                    {
                        id: '1',
                        title: 'First',
                        message: 'First message',
                        type: NotificationType.SYSTEM,
                        read: true,
                        createdAt: '2024-01-01',
                    },
                ],
                page: 1,
                hasMore: true,
                loading: false,
            });

            mockHttp.get.mockResolvedValueOnce({
                items: [
                    {
                        id: '2',
                        title: 'Second',
                        message: 'Second message',
                        type: NotificationType.SYSTEM,
                        read: true,
                        createdAt: '2024-01-02',
                    },
                ],
                total: 2,
                page: 2,
                pageSize: 20,
            });

            await useNotificationStore.getState().loadMore();

            const state = useNotificationStore.getState();
            expect(state.page).toBe(2);
            expect(state.notifications).toHaveLength(2);
        });

        it('should not load more when hasMore is false', async () => {
            useNotificationStore.setState({
                hasMore: false,
                loading: false,
            });

            await useNotificationStore.getState().loadMore();

            expect(mockHttp.get).not.toHaveBeenCalled();
        });

        it('should not load more when already loading', async () => {
            useNotificationStore.setState({
                hasMore: true,
                loading: true,
            });

            await useNotificationStore.getState().loadMore();

            expect(mockHttp.get).not.toHaveBeenCalled();
        });
    });
});
