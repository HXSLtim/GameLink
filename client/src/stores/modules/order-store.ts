import { create } from 'zustand';
import { http } from '@/lib/http';

export const OrderStatus = {
    PENDING: 'pending',
    CONFIRMED: 'confirmed', // Paid, waiting for player
    IN_PROGRESS: 'in_progress', // Player joined
    COMPLETED: 'completed',
    CANCELED: 'canceled', // Single 'l' as per doc
    REFUNDED: 'refunded',
    DISPUTED: 'disputed'
} as const;

export type OrderStatus = typeof OrderStatus[keyof typeof OrderStatus];

export interface Order {
    id: string; // or number
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
    fetchOrderById: (id: string) => Promise<void>;
    createOrder: (payload: CreateOrderPayload) => Promise<void>;
    cancelOrder: (id: string) => Promise<void>;
    updateDraft: (draft: Partial<OrderState['orderDraft']>) => void;
    clearDraft: () => void;

    // New Action for Stats
    fetchOrderStats: () => Promise<OrderStats | null>;

    // Player Actions
    fetchPlayerOrders: () => Promise<void>;
    acceptOrder: (id: string) => Promise<void>;
    rejectOrder: (id: string) => Promise<void>;
    submitDispute: (id: string, reason: string, description: string) => Promise<void>;
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
            set({ loading: false, error: (err as Error).message || 'Failed to fetch orders' });
        }
    },

    fetchOrderById: async (id) => {
        set({ loading: true, error: null });
        try {
            const data = await http.get<Order>(`/orders/${id}`);
            set({ currentOrder: data, loading: false });
        } catch (err) {
            set({ loading: false, error: (err as Error).message || 'Failed to get order details' });
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
            set({ loading: false, error: (err as Error).message || 'Failed to create order' });
            throw err;
        }
    },

    cancelOrder: async (id) => {
        try {
            await http.post(`/orders/${id}/cancel`);
            // Optimistic update
            set((state) => ({
                myOrders: state.myOrders.map(o => o.id === id ? { ...o, status: OrderStatus.CANCELED } : o),
                currentOrder: state.currentOrder?.id === id ? { ...state.currentOrder, status: OrderStatus.CANCELED } : state.currentOrder
            }));
        } catch (err) {
            console.error(err);
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
            set({ loading: false, error: (err as Error).message || 'Failed to fetch order stats' });
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
            set({ loading: false, error: (err as Error).message || 'Failed to fetch player orders' });
        }
    },

    acceptOrder: async (id: string) => {
        await http.put(`/player/orders/${id}/accept`);
        set((state) => ({
            playerOrders: state.playerOrders.map(o => o.id === id ? { ...o, status: OrderStatus.IN_PROGRESS } : o)
        }));
    },

    rejectOrder: async (id: string) => {
        await http.put(`/player/orders/${id}/reject`);
        set((state) => ({
            playerOrders: state.playerOrders.map(o => o.id === id ? { ...o, status: OrderStatus.CANCELED } : o)
        }));
    },

    submitDispute: async (id: string, reason: string, description: string) => {
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

    // WebSocket Integration
    subscribeToOrderUpdates: () => {
        import('@/lib/websocket').then(({ wsService }) => {
            wsService.connect();

            wsService.on('order.status_updated', (payload: { orderId: string, status: OrderStatus }) => {
                const { orderId, status } = payload;
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
            });

            wsService.on('order.created', (newOrder: Order) => {
                set((state) => ({
                    myOrders: [newOrder, ...state.myOrders],
                    playerOrders: [newOrder, ...state.playerOrders] // Naive update, ideally filtered
                }));
            });
        });
    },

    unsubscribeFromOrderUpdates: () => {
        import('@/lib/websocket').then(() => {
            // In a real app with strict handler references, we'd need to store the functions. 
            // For now assuming the service or store lifecycle manages this simplistically 
            // or we accept potential memory leak if unsubscribe isn't perfect in this iteration.
            // A better way is to move WS logic to a dedicated hook or Context.
        });
    }
}));
