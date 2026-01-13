import { create } from 'zustand';
// import { persist } from 'zustand/middleware'; // Notifications often transient or fetched fresh

export interface AppNotification {
    id: string;
    title: string;
    message: string;
    type: 'system' | 'order' | 'promotion';
    read: boolean;
    createdAt: string;
    actionUrl?: string;
}

export interface NotificationStore {
    notifications: AppNotification[];
    unreadCount: number;
    loading: boolean;

    fetchNotifications: () => Promise<void>;
    markAsRead: (id: string) => Promise<void>;
    markAllAsRead: () => Promise<void>;
    addNotification: (notification: AppNotification) => void; // For real-time push
}

export const useNotificationStore = create<NotificationStore>((set, get) => ({
    notifications: [],
    unreadCount: 0,
    loading: false,

    fetchNotifications: async () => {
        set({ loading: true });
        try {
            // await http.get('/notifications');
            await new Promise(resolve => setTimeout(resolve, 500));

            // Mock data
            const mock: AppNotification[] = [
                { id: 'n1', title: 'Welcome', message: 'Welcome to GameLink!', type: 'system', read: false, createdAt: new Date().toISOString() },
                { id: 'n2', title: 'Discount', message: '50% off on first order', type: 'promotion', read: false, createdAt: new Date().toISOString() }
            ];

            set({
                notifications: mock,
                unreadCount: mock.filter(n => !n.read).length,
                loading: false
            });
        } catch (e) {
            set({ loading: false });
        }
    },

    markAsRead: async (id) => {
        const { notifications } = get();
        const target = notifications.find(n => n.id === id);
        if (!target || target.read) return;

        // Optimistic
        set(state => {
            const updated = state.notifications.map(n => n.id === id ? { ...n, read: true } : n);
            return {
                notifications: updated,
                unreadCount: Math.max(0, state.unreadCount - 1)
            };
        });

        // API call
        // await http.post(`/notifications/${id}/read`);
    },

    markAllAsRead: async () => {
        set(state => ({
            notifications: state.notifications.map(n => ({ ...n, read: true })),
            unreadCount: 0
        }));
        // await http.post('/notifications/read-all');
    },

    addNotification: (notification) => {
        set(state => ({
            notifications: [notification, ...state.notifications],
            unreadCount: state.unreadCount + 1
        }));
    }
}));
