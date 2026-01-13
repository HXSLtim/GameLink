import { create } from 'zustand';
import { http } from '@/lib/http';

export const OrderStatus = {
    PENDING: 'pending',
    PAID: 'paid',
    ACCEPTED: 'accepted',
    COMPLETED: 'completed',
    CANCELLED: 'cancelled',
    REFUNDED: 'refunded'
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

export interface OrderActions {
    fetchOrders: () => Promise<void>;
    fetchOrderById: (id: string) => Promise<void>;
    createOrder: (payload: any) => Promise<void>;
    cancelOrder: (id: string) => Promise<void>;
    updateDraft: (draft: Partial<OrderState['orderDraft']>) => void;
    clearDraft: () => void;
}

export const useOrderStore = create<OrderState & OrderActions>((set) => ({
    myOrders: [],
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
        } catch (err: any) {
            set({ loading: false, error: err.message || 'Failed to fetch orders' });
        }
    },

    fetchOrderById: async (id) => {
        set({ loading: true, error: null });
        try {
            const data = await http.get<Order>(`/orders/${id}`);
            set({ currentOrder: data, loading: false });
        } catch (err: any) {
            set({ loading: false, error: err.message || 'Failed to get order details' });
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
        } catch (err: any) {
            set({ loading: false, error: err.message || 'Failed to create order' });
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
        } catch (err: any) {
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
    }
}));
