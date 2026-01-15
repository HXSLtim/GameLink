import { describe, it, expect, vi, beforeEach } from 'vitest';
import { useOrderStore } from '../order-store';
import { http } from '@/lib/http';

// Mock http client
vi.mock('@/lib/http', () => ({
    http: {
        post: vi.fn(),
        get: vi.fn(),
    },
}));

describe('Order Store', () => {
    beforeEach(() => {
        useOrderStore.setState({
            myOrders: [],
            currentOrder: null,
            loading: false,
            error: null,
        });
        vi.clearAllMocks();
    });

    it('should fetch orders successfully', async () => {
        const mockOrders = [
            { id: '1', orderNo: 'ORD-1', amount: 100 },
            { id: '2', orderNo: 'ORD-2', amount: 200 }
        ];
        (http.get as any).mockResolvedValue(mockOrders);

        await useOrderStore.getState().fetchOrders();

        const state = useOrderStore.getState();
        expect(state.myOrders).toEqual(mockOrders);
        expect(state.loading).toBe(false);
        expect(http.get).toHaveBeenCalledWith('/orders');
    });

    it('should handle fetch orders error', async () => {
        const errorMessage = 'Network Error';
        (http.get as any).mockRejectedValue(new Error(errorMessage));

        await useOrderStore.getState().fetchOrders();

        const state = useOrderStore.getState();
        expect(state.myOrders).toEqual([]);
        expect(state.loading).toBe(false);
        expect(state.error).toBe(errorMessage);
    });

    it('should create order successfully', async () => {
        const payload = { playerId: 1, amount: 100 };
        const mockOrder = { id: '3', ...payload, orderNo: 'ORD-3' };

        (http.post as any).mockResolvedValue(mockOrder);

        const initialState = useOrderStore.getState();
        expect(initialState.myOrders.length).toBe(0);

        await useOrderStore.getState().createOrder(payload);

        const state = useOrderStore.getState();
        expect(state.myOrders).toHaveLength(1);
        expect(state.myOrders[0]).toEqual(mockOrder);
        expect(state.currentOrder).toEqual(mockOrder);
        expect(http.post).toHaveBeenCalledWith('/orders', payload);
    });

    it('should handle create order error', async () => {
        const errorMessage = 'Failed to create';
        (http.post as any).mockRejectedValue(new Error(errorMessage));

        await expect(useOrderStore.getState().createOrder({}))
            .rejects.toThrow(errorMessage);

        const state = useOrderStore.getState();
        expect(state.myOrders).toHaveLength(0);
        expect(state.error).toBe(errorMessage);
    });

    it('should cancel order optimistically', async () => {
        const initialOrder = { id: '1', status: 'pending' };
        useOrderStore.setState({
            myOrders: [initialOrder as any],
            currentOrder: initialOrder as any
        });

        (http.post as any).mockResolvedValue({});

        await useOrderStore.getState().cancelOrder('1');

        const state = useOrderStore.getState();
        expect(state.myOrders[0].status).toBe('cancelled');
        expect(state.currentOrder?.status).toBe('cancelled');
        expect(http.post).toHaveBeenCalledWith('/orders/1/cancel');
    });
});
