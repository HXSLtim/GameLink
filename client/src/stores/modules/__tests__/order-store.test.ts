import { describe, it, expect, vi, beforeEach } from 'vitest';
import { useOrderStore, OrderStatus, type Order, type OrderStats } from '../order-store';
import { http } from '@/lib/http';

// Mock http client
vi.mock('@/lib/http', () => ({
    http: {
        post: vi.fn(),
        get: vi.fn(),
        put: vi.fn(),
    },
}));

// Mock WebSocket service
vi.mock('@/lib/websocket', () => ({
    wsService: {
        connect: vi.fn(),
        on: vi.fn(),
        off: vi.fn(),
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
        (http.get as ReturnType<typeof vi.fn>).mockResolvedValue(mockOrders);

        await useOrderStore.getState().fetchOrders();

        const state = useOrderStore.getState();
        expect(state.myOrders).toEqual(mockOrders);
        expect(state.loading).toBe(false);
        expect(http.get).toHaveBeenCalledWith('/orders');
    });

    it('should handle fetch orders error', async () => {
        const errorMessage = 'Network Error';
        (http.get as ReturnType<typeof vi.fn>).mockRejectedValue(new Error(errorMessage));

        await useOrderStore.getState().fetchOrders();

        const state = useOrderStore.getState();
        expect(state.myOrders).toEqual([]);
        expect(state.loading).toBe(false);
        expect(state.error).toBe(errorMessage);
    });

    it('should create order successfully', async () => {
        const payload = { playerId: 1, gameId: 1, quantity: 2, amount: 100 };
        const mockOrder = { id: '3', ...payload, orderNo: 'ORD-3' };

        (http.post as ReturnType<typeof vi.fn>).mockResolvedValue(mockOrder);

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
        (http.post as ReturnType<typeof vi.fn>).mockRejectedValue(new Error(errorMessage));

        const invalidPayload = { playerId: 1, gameId: 1, quantity: 1, amount: 50 };
        await expect(useOrderStore.getState().createOrder(invalidPayload))
            .rejects.toThrow(errorMessage);

        const state = useOrderStore.getState();
        expect(state.myOrders).toHaveLength(0);
        expect(state.error).toBe(errorMessage);
    });

    it('should cancel order optimistically', async () => {
        // Use Partial type for partial mock state
        const initialOrder = { id: 1, status: OrderStatus.PENDING } as Partial<Order>;
        useOrderStore.setState({
            myOrders: [initialOrder as Order],
            currentOrder: initialOrder as Order
        });

        (http.post as ReturnType<typeof vi.fn>).mockResolvedValue({});

        await useOrderStore.getState().cancelOrder(1);

        const state = useOrderStore.getState();
        expect(state.myOrders[0].status).toBe(OrderStatus.CANCELLED);
        expect(state.currentOrder?.status).toBe(OrderStatus.CANCELLED);
        expect(http.post).toHaveBeenCalledWith('/orders/1/cancel');
    });

    it('should fetch order by id successfully', async () => {
        const mockOrder = { id: 1, orderNo: 'ORD-1', status: OrderStatus.PENDING } as Order;
        (http.get as ReturnType<typeof vi.fn>).mockResolvedValue(mockOrder);

        await useOrderStore.getState().fetchOrderById(1);

        const state = useOrderStore.getState();
        expect(state.currentOrder).toEqual(mockOrder);
        expect(http.get).toHaveBeenCalledWith('/orders/1');
    });

    it('should handle fetch order by id error', async () => {
        const errorMessage = 'Order not found';
        const error = new Error(errorMessage);
        (error as any).response = { data: { message: errorMessage } };
        (http.get as ReturnType<typeof vi.fn>).mockRejectedValue(error);

        await useOrderStore.getState().fetchOrderById(999);

        const state = useOrderStore.getState();
        expect(state.error).toBe(errorMessage);
        expect(state.loading).toBe(false);
    });

    it('should update draft', () => {
        useOrderStore.getState().updateDraft({ playerId: 1 });
        let state = useOrderStore.getState();
        expect(state.orderDraft).toEqual({ playerId: 1 });

        useOrderStore.getState().updateDraft({ gameId: 2 });
        state = useOrderStore.getState();
        expect(state.orderDraft).toEqual({ playerId: 1, gameId: 2 });
    });

    it('should clear draft', () => {
        useOrderStore.setState({ orderDraft: { playerId: 1, gameId: 2 } });
        expect(useOrderStore.getState().orderDraft).not.toBeNull();

        useOrderStore.getState().clearDraft();
        expect(useOrderStore.getState().orderDraft).toBeNull();
    });

    it('should fetch order stats successfully', async () => {
        const mockStats: OrderStats = {
            totalCount: 10,
            monthlyCount: 5,
            monthlyChange: 20,
            pendingCount: 2,
            inProgressCount: 1,
            completedCount: 6,
            canceledCount: 1,
            totalSpentCents: 50000,
            avgOrderAmountCents: 5000,
        };
        (http.get as ReturnType<typeof vi.fn>).mockResolvedValue(mockStats);

        const result = await useOrderStore.getState().fetchOrderStats();

        expect(result).toEqual(mockStats);
        expect(http.get).toHaveBeenCalledWith('/user/orders/stats');
    });

    it('should return null on fetch order stats error', async () => {
        (http.get as ReturnType<typeof vi.fn>).mockRejectedValue(new Error('Network error'));

        const result = await useOrderStore.getState().fetchOrderStats();

        expect(result).toBeNull();
    });

    it('should fetch player orders successfully', async () => {
        const mockPlayerOrders = [
            { id: 1, orderNo: 'ORD-1', status: OrderStatus.PENDING },
            { id: 2, orderNo: 'ORD-2', status: OrderStatus.CONFIRMED },
        ] as Order[];
        (http.get as ReturnType<typeof vi.fn>).mockResolvedValue(mockPlayerOrders);

        await useOrderStore.getState().fetchPlayerOrders();

        const state = useOrderStore.getState();
        expect(state.playerOrders).toEqual(mockPlayerOrders);
        expect(http.get).toHaveBeenCalledWith('/player/orders');
    });

    it('should handle fetch player orders error', async () => {
        const error = new Error('Network error');
        (error as any).response = { data: { message: 'Network error' } };
        (http.get as ReturnType<typeof vi.fn>).mockRejectedValue(error);

        await useOrderStore.getState().fetchPlayerOrders();

        const state = useOrderStore.getState();
        expect(state.error).toBe('Network error');
        expect(state.loading).toBe(false);
    });

    it('should accept order as player', async () => {
        const pendingOrder = { id: 1, status: OrderStatus.CONFIRMED } as Order;
        useOrderStore.setState({ playerOrders: [pendingOrder] });

        (http.put as ReturnType<typeof vi.fn>).mockResolvedValue({});

        await useOrderStore.getState().acceptOrder(1);

        const state = useOrderStore.getState();
        expect(state.playerOrders[0].status).toBe(OrderStatus.IN_PROGRESS);
        expect(http.put).toHaveBeenCalledWith('/player/orders/1/accept');
    });

    it('should reject order as player', async () => {
        const confirmedOrder = { id: 1, status: OrderStatus.CONFIRMED } as Order;
        useOrderStore.setState({ playerOrders: [confirmedOrder] });

        (http.put as ReturnType<typeof vi.fn>).mockResolvedValue({});

        await useOrderStore.getState().rejectOrder(1);

        const state = useOrderStore.getState();
        expect(state.playerOrders[0].status).toBe(OrderStatus.CANCELLED);
        expect(http.put).toHaveBeenCalledWith('/player/orders/1/reject');
    });

    it('should submit dispute successfully', async () => {
        const activeOrder = { id: 1, status: OrderStatus.IN_PROGRESS } as Order;
        useOrderStore.setState({
            currentOrder: activeOrder,
            myOrders: [activeOrder]
        });

        (http.post as ReturnType<typeof vi.fn>).mockResolvedValue({});

        await useOrderStore.getState().submitDispute(1, 'service_issue', 'Poor service');

        const state = useOrderStore.getState();
        expect(state.currentOrder?.status).toBe(OrderStatus.DISPUTED);
        expect(state.myOrders[0].status).toBe(OrderStatus.DISPUTED);
        expect(http.post).toHaveBeenCalledWith('/orders/1/dispute', {
            reason: 'service_issue',
            description: 'Poor service',
        });
    });

    it('should handle submit dispute error', async () => {
        const activeOrder = { id: 1, status: OrderStatus.IN_PROGRESS } as Order;
        useOrderStore.setState({
            currentOrder: activeOrder,
            myOrders: [activeOrder]
        });

        (http.post as ReturnType<typeof vi.fn>).mockRejectedValue(new Error('Dispute failed'));

        await expect(useOrderStore.getState().submitDispute(1, 'issue', 'desc'))
            .rejects.toThrow('Dispute failed');

        const state = useOrderStore.getState();
        expect(state.error).toBe('Dispute failed');
        expect(state.loading).toBe(false);
    });

    it('should update only matching order in currentOrder on dispute', async () => {
        const activeOrder = { id: 1, status: OrderStatus.IN_PROGRESS } as Order;
        const otherOrder = { id: 2, status: OrderStatus.COMPLETED } as Order;
        useOrderStore.setState({
            currentOrder: otherOrder,
            myOrders: [activeOrder, otherOrder]
        });

        (http.post as ReturnType<typeof vi.fn>).mockResolvedValue({});

        await useOrderStore.getState().submitDispute(1, 'issue', 'desc');

        const state = useOrderStore.getState();
        expect(state.currentOrder?.status).toBe(OrderStatus.COMPLETED); // Unchanged
        expect(state.myOrders[0].status).toBe(OrderStatus.DISPUTED);
        expect(state.myOrders[1].status).toBe(OrderStatus.COMPLETED); // Unchanged
    });

    it('should subscribe to order updates', async () => {
        const { wsService } = await import('@/lib/websocket');

        useOrderStore.getState().subscribeToOrderUpdates();

        // Wait for dynamic import to complete
        await new Promise(resolve => setTimeout(resolve, 10));

        expect(wsService.connect).toHaveBeenCalled();
        expect(wsService.on).toHaveBeenCalledWith('order.status_updated', expect.any(Function));
        expect(wsService.on).toHaveBeenCalledWith('order.created', expect.any(Function));
    });

    it('should unsubscribe from order updates', async () => {
        const { wsService } = await import('@/lib/websocket');

        useOrderStore.getState().subscribeToOrderUpdates();
        await new Promise(resolve => setTimeout(resolve, 10));

        useOrderStore.getState().unsubscribeFromOrderUpdates();

        expect(wsService.off).toHaveBeenCalledWith('order.status_updated', expect.any(Function));
        expect(wsService.off).toHaveBeenCalledWith('order.created', expect.any(Function));
    });

    it('should handle order draft with initial null state', () => {
        expect(useOrderStore.getState().orderDraft).toBeNull();

        useOrderStore.getState().updateDraft({ playerId: 1 });
        expect(useOrderStore.getState().orderDraft).toEqual({ playerId: 1 });
    });

    it('should merge draft updates correctly', () => {
        useOrderStore.getState().updateDraft({ playerId: 1, gameId: 2 });
        useOrderStore.getState().updateDraft({ serviceItemId: 3 });
        useOrderStore.getState().updateDraft({ quantity: 2 });

        const state = useOrderStore.getState();
        expect(state.orderDraft).toEqual({
            playerId: 1,
            gameId: 2,
            serviceItemId: 3,
            quantity: 2,
        });
    });

    it('should handle empty orders list on fetch', async () => {
        (http.get as ReturnType<typeof vi.fn>).mockResolvedValue([]);

        await useOrderStore.getState().fetchOrders();

        const state = useOrderStore.getState();
        expect(state.myOrders).toEqual([]);
        expect(state.loading).toBe(false);
    });

    it('should set loading state during fetch', async () => {
        let resolveFetch: (value: Order[]) => void;
        const fetchPromise = new Promise<Order[]>(resolve => {
            resolveFetch = resolve;
        });

        (http.get as ReturnType<typeof vi.fn>).mockReturnValue(fetchPromise);

        const fetchPromise2 = useOrderStore.getState().fetchOrders();

        // Should be loading
        expect(useOrderStore.getState().loading).toBe(true);

        resolveFetch!([]);
        await fetchPromise2;

        // Should be done loading
        expect(useOrderStore.getState().loading).toBe(false);
    });

    it('should not affect other orders when cancelling specific order', async () => {
        const order1 = { id: 1, status: OrderStatus.PENDING } as Order;
        const order2 = { id: 2, status: OrderStatus.CONFIRMED } as Order;
        const order3 = { id: 3, status: OrderStatus.IN_PROGRESS } as Order;

        useOrderStore.setState({
            myOrders: [order1, order2, order3],
            currentOrder: order1,
        });

        (http.post as ReturnType<typeof vi.fn>).mockResolvedValue({});

        await useOrderStore.getState().cancelOrder(1);

        const state = useOrderStore.getState();
        expect(state.myOrders[0].status).toBe(OrderStatus.CANCELLED);
        expect(state.myOrders[1].status).toBe(OrderStatus.CONFIRMED); // Unchanged
        expect(state.myOrders[2].status).toBe(OrderStatus.IN_PROGRESS); // Unchanged
    });

    it('should handle acceptOrder error gracefully', async () => {
        const order = { id: 1, status: OrderStatus.CONFIRMED } as Order;
        useOrderStore.setState({ playerOrders: [order] });

        (http.put as ReturnType<typeof vi.fn>).mockRejectedValue(new Error('Accept failed'));

        await expect(useOrderStore.getState().acceptOrder(1)).rejects.toThrow();
    });

    it('should handle rejectOrder error gracefully', async () => {
        const order = { id: 1, status: OrderStatus.CONFIRMED } as Order;
        useOrderStore.setState({ playerOrders: [order] });

        (http.put as ReturnType<typeof vi.fn>).mockRejectedValue(new Error('Reject failed'));

        await expect(useOrderStore.getState().rejectOrder(1)).rejects.toThrow();
    });
});
