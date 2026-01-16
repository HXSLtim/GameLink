import { create } from 'zustand';
import { http } from '@/lib/http';
import { getErrorMessage, logError } from '@/lib/error';
import type { wsService as WsServiceType } from '@/lib/websocket';

// Store WebSocket handler references for proper cleanup
let orderStatusHandler: ((payload: unknown) => void) | null = null;
let orderCreatedHandler: ((payload: unknown) => void) | null = null;
let wsServiceInstance: typeof WsServiceType | null = null;

export const OrderStatus = {
    PENDING: 'pending',
    CONFIRMED: 'confirmed', // Paid, waiting for player
    IN_PROGRESS: 'in_progress', // Player joined
    COMPLETED: 'completed',
    CANCELLED: 'canceled', // Changed to match backend
    REFUNDED: 'refunded',
    DISPUTED: 'disputed'
} as const;

export type OrderStatus = typeof OrderStatus[keyof typeof OrderStatus];

export interface Order {
    id: number;
    orderNo: string;
    playerId: number;
    userId: number;
    gameId: number;
    gameName: string;
    amount: number;
    quantity: number;
    status: OrderStatus;
    createdAt: string;
    scheduledTime?: string;
}

export interface OrderState {
    myOrders: Order[];
    playerOrders: Order[]; // For players to see received orders
    currentOrder: Order | null;
    activeOrders: Order[];

    // Draft
    orderDraft: {
        playerId?: number;
        gameId?: number;
        serviceItemId?: number;
        quantity?: number;
        scheduledTime?: string;
    } | null;

    loading: boolean;
    error: string | null;
}

export interface CreateOrderPayload {
    playerId: number;
    gameId: number;
    serviceItemId?: number;
    quantity: number;
    amount: number;
    note?: string;
    scheduledTime?: string;
}

export interface OrderActions {
    fetchOrders: () => Promise<void>;
    fetchOrderById: (id: number) => Promise<void>;
    createOrder: (payload: CreateOrderPayload) => Promise<void>;
    cancelOrder: (id: number) => Promise<void>;
    updateDraft: (draft: Partial<OrderState['orderDraft']>) => void;
    clearDraft: () => void;

    // New Action for Stats
    fetchOrderStats: () => Promise<OrderStats | null>;

    // Player Actions
    fetchPlayerOrders: () => Promise<void>;
    acceptOrder: (id: number) => Promise<void>;
    rejectOrder: (id: number) => Promise<void>;
    submitDispute: (id: number, reason: string, description: string) => Promise<void>;
    subscribeToOrderUpdates: () => void;
    unsubscribeFromOrderUpdates: () => void;
}

export interface OrderStats {
    totalCount: number;
    monthlyCount: number;
    monthlyChange: number;
    pendingCount: number;
    inProgressCount: number;
    completedCount: number;
    canceledCount: number;
    totalSpentCents: number;
    avgOrderAmountCents: number;
}

export const useOrderStore = create<OrderState & OrderActions>((set) => ({
    myOrders: [],
    playerOrders: [],
    currentOrder: null,
    activeOrders: [],
    orderDraft: null,
    loading: false,
    error: null,

    fetchOrders: async () => {
        set({ loading: true, error: null });
        try {
            const data = await http.get<Order[]>('/orders');
            set({ myOrders: data, loading: false });
        } catch (err) {
            logError('fetchOrders', err);
            set({ loading: false, error: getErrorMessage(err, 'Failed to fetch orders') });
        }
    },

    fetchOrderById: async (id) => {
        set({ loading: true, error: null });
        try {
            const data = await http.get<Order>(`/orders/${id}`);
            set({ currentOrder: data, loading: false });
        } catch (err) {
            logError('fetchOrderById', err);
            set({ loading: false, error: getErrorMessage(err, 'Failed to get order details') });
        }
    },

    createOrder: async (payload) => {
        set({ loading: true, error: null });
        try {
            const data = await http.post<Order>('/orders', payload);
            set((state) => ({
                myOrders: [data, ...state.myOrders],
                currentOrder: data,
                loading: false
            }));
        } catch (err) {
            logError('createOrder', err);
            set({ loading: false, error: getErrorMessage(err, 'Failed to create order') });
            throw err;
        }
    },

    cancelOrder: async (id) => {
        try {
            await http.post(`/orders/${id}/cancel`);
            // Optimistic update
            set((state) => ({
                myOrders: state.myOrders.map(o => o.id === id ? { ...o, status: OrderStatus.CANCELLED } : o),
                currentOrder: state.currentOrder?.id === id ? { ...state.currentOrder, status: OrderStatus.CANCELLED } : state.currentOrder
            }));
        } catch (err) {
            logError('cancelOrder', err);
            throw err;
        }
    },

    updateDraft: (draft) => {
        set((state) => ({
            orderDraft: { ...(state.orderDraft || {}), ...draft }
        }));
    },

    clearDraft: () => {
        set({ orderDraft: null });
    },

    fetchOrderStats: async () => {
        set({ loading: true, error: null });
        try {
            const data = await http.get<OrderStats>('/user/orders/stats');
            set({ loading: false });
            return data;
        } catch (err) {
            logError('fetchOrderStats', err);
            set({ loading: false, error: getErrorMessage(err, 'Failed to fetch order stats') });
            return null;
        }
    },

    // Player Side Actions
    fetchPlayerOrders: async () => {
        set({ loading: true, error: null });
        try {
            const data = await http.get<Order[]>('/player/orders'); // Report says /player/orders
            set({ playerOrders: data, loading: false });
        } catch (err) {
            logError('fetchPlayerOrders', err);
            set({ loading: false, error: getErrorMessage(err, 'Failed to fetch player orders') });
        }
    },

    acceptOrder: async (id: number) => {
        await http.put(`/player/orders/${id}/accept`);
        set((state) => ({
            playerOrders: state.playerOrders.map(o => o.id === id ? { ...o, status: OrderStatus.IN_PROGRESS } : o)
        }));
    },

    rejectOrder: async (id: number) => {
        await http.put(`/player/orders/${id}/reject`);
        set((state) => ({
            playerOrders: state.playerOrders.map(o => o.id === id ? { ...o, status: OrderStatus.CANCELLED } : o)
        }));
    },

    submitDispute: async (id: number, reason: string, description: string) => {
        set({ loading: true, error: null });
        try {
            await http.post(`/orders/${id}/dispute`, { reason, description });
            set((state) => ({
                currentOrder: state.currentOrder?.id === id ? { ...state.currentOrder, status: OrderStatus.DISPUTED } : state.currentOrder,
                myOrders: state.myOrders.map(o => o.id === id ? { ...o, status: OrderStatus.DISPUTED } : o),
                loading: false
            }));
        } catch (err) {
            set({ loading: false, error: (err as Error).message || 'Failed to submit dispute' });
            throw err;
        }
    },

    // WebSocket Integration with proper cleanup
    subscribeToOrderUpdates: () => {
        import('@/lib/websocket').then(({ wsService }) => {
            wsServiceInstance = wsService;
            wsService.connect();

            // Remove existing handlers first to prevent duplicates
            if (orderStatusHandler) {
                wsService.off('order.status_updated', orderStatusHandler);
            }
            if (orderCreatedHandler) {
                wsService.off('order.created', orderCreatedHandler);
            }

            // Create and store new handlers
            orderStatusHandler = (payload: unknown) => {
                const { orderId, status } = payload as { orderId: number, status: OrderStatus };
                set((state) => {
                    const updatedMyOrders = state.myOrders.map(o => o.id === orderId ? { ...o, status } : o);
                    const updatedCurrent = state.currentOrder?.id === orderId ? { ...state.currentOrder, status } : state.currentOrder;
                    const updatedPlayerOrders = state.playerOrders.map(o => o.id === orderId ? { ...o, status } : o);

                    return {
                        myOrders: updatedMyOrders,
                        currentOrder: updatedCurrent,
                        playerOrders: updatedPlayerOrders
                    };
                });
            };

            orderCreatedHandler = (data: unknown) => {
                const newOrder = data as Order;
                set((state) => ({
                    myOrders: [newOrder, ...state.myOrders],
                    playerOrders: [newOrder, ...state.playerOrders]
                }));
            };

            // Register handlers
            wsService.on('order.status_updated', orderStatusHandler);
            wsService.on('order.created', orderCreatedHandler);
        });
    },

    unsubscribeFromOrderUpdates: () => {
        // Properly clean up WebSocket handlers to prevent memory leaks
        if (wsServiceInstance) {
            if (orderStatusHandler) {
                wsServiceInstance.off('order.status_updated', orderStatusHandler);
                orderStatusHandler = null;
            }
            if (orderCreatedHandler) {
                wsServiceInstance.off('order.created', orderCreatedHandler);
                orderCreatedHandler = null;
            }
        }
    }
}));
