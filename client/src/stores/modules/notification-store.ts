import { create } from 'zustand';
import { http } from '@/lib/http';

// ============ Enums ============

export const NotificationType = {
    SYSTEM: 'system',           // 系统通知
    ORDER: 'order',             // 订单通知
    PROMOTION: 'promotion',     // 促销通知
    CHAT: 'chat',               // 聊天通知
    WALLET: 'wallet'            // 钱包通知
} as const;

export type NotificationType = typeof NotificationType[keyof typeof NotificationType];

// ============ Interfaces ============

export interface AppNotification {
    id: number;
    title: string;
    content: string; // Backend uses 'content', not 'message'
    type: NotificationType;
    read: boolean;
    createdAt: string;
    actionUrl?: string;
    metadata?: Record<string, unknown>;
}

export interface NotificationStore {
    notifications: AppNotification[];
    unreadCount: number;
    loading: boolean;
    error: string | null;
    hasMore: boolean;
    page: number;

    fetchNotifications: (reset?: boolean) => Promise<void>;
    fetchUnreadCount: () => Promise<void>;
    markAsRead: (id: number) => Promise<void>;
    markAllAsRead: () => Promise<void>;
    deleteNotification: (id: number) => Promise<void>;
    addNotification: (notification: AppNotification) => void;
    loadMore: () => Promise<void>;
}

// ============ Store ============

export const useNotificationStore = create<NotificationStore>((set, get) => ({
    notifications: [],
    unreadCount: 0,
    loading: false,
    error: null,
    hasMore: true,
    page: 1,

    fetchNotifications: async (reset = false) => {
        const currentPage = reset ? 1 : get().page;
        set({ loading: true, error: null });

        try {
            const data = await http.get<{
                items: AppNotification[];
                total: number;
                page: number;
                pageSize: number;
            }>('/user/notifications', {
                params: { page: currentPage, pageSize: 20 }
            });

            const items = data.items || [];
            const hasMore = items.length === 20;

            set(state => ({
                notifications: reset ? items : [...state.notifications, ...items],
                unreadCount: items.filter(n => !n.read).length,
                hasMore,
                page: currentPage,
                loading: false
            }));
        } catch (err) {
            set({
                loading: false,
                error: err instanceof Error ? err.message : 'Failed to fetch notifications'
            });
        }
    },

    fetchUnreadCount: async () => {
        try {
            const data = await http.get<{ count: number }>('/user/notifications/unread-count');
            set({ unreadCount: data.count });
        } catch (err) {
            console.error('Failed to fetch unread count:', err);
        }
    },

    markAsRead: async (id) => {
        const { notifications } = get();
        const target = notifications.find(n => n.id === id);
        if (!target || target.read) return;

        // Optimistic update
        set(state => {
            const updated = state.notifications.map(n =>
                n.id === id ? { ...n, read: true } : n
            );
            return {
                notifications: updated,
                unreadCount: Math.max(0, state.unreadCount - 1)
            };
        });

        try {
            await http.put(`/user/notifications/${id}/read`);
        } catch (err) {
            // Rollback on error
            set(state => {
                const updated = state.notifications.map(n =>
                    n.id === id ? { ...n, read: false } : n
                );
                return {
                    notifications: updated,
                    unreadCount: state.unreadCount + 1,
                    error: err instanceof Error ? err.message : 'Failed to mark as read'
                };
            });
        }
    },

    markAllAsRead: async () => {
        const previousNotifications = [...get().notifications];
        const previousUnreadCount = get().unreadCount;

        // Optimistic update
        set(state => ({
            notifications: state.notifications.map(n => ({ ...n, read: true })),
            unreadCount: 0
        }));

        try {
            await http.put('/user/notifications/read-all');
        } catch (err) {
            // Rollback on error
            set({
                notifications: previousNotifications,
                unreadCount: previousUnreadCount,
                error: err instanceof Error ? err.message : 'Failed to mark all as read'
            });
        }
    },

    deleteNotification: async (id) => {
        const previousNotifications = [...get().notifications];
        const target = previousNotifications.find(n => n.id === id);

        // Optimistic update
        set(state => ({
            notifications: state.notifications.filter(n => n.id !== id),
            unreadCount: target && !target.read
                ? Math.max(0, state.unreadCount - 1)
                : state.unreadCount
        }));

        try {
            await http.delete(`/user/notifications/${id}`);
        } catch (err) {
            // Rollback on error
            set({
                notifications: previousNotifications,
                unreadCount: get().unreadCount + (target && !target.read ? 1 : 0),
                error: err instanceof Error ? err.message : 'Failed to delete notification'
            });
        }
    },

    addNotification: (notification) => {
        set(state => ({
            notifications: [notification, ...state.notifications],
            unreadCount: notification.read ? state.unreadCount : state.unreadCount + 1
        }));
    },

    loadMore: async () => {
        const { hasMore, loading, page } = get();
        if (!hasMore || loading) return;

        set({ page: page + 1 });
        await get().fetchNotifications();
    }
}));
