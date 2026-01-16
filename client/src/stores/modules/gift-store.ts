import { create } from 'zustand';
import { http } from '@/lib/http';
import { getErrorMessage, logError } from '@/lib/error';

// ============ Interfaces ============

export interface GiftType {
    id: number;
    name: string;
    icon: string;
    priceCents: number;
    description?: string;
    category: string;
    sortOrder: number;
    isActive: boolean;
}

export interface Gift {
    id: number;
    giftTypeId: number;
    giftType: GiftType;
    senderId: number;
    senderName: string;
    senderAvatar: string;
    receiverId: number;
    orderId?: number;
    quantity: number;
    totalCents: number;
    message?: string;
    createdAt: string;
}

export interface GiftStats {
    totalReceived: number;
    totalReceivedCents: number;
    todayReceived: number;
    todayReceivedCents: number;
    weeklyReceived: number;
    weeklyReceivedCents: number;
    monthlyReceived: number;
    monthlyReceivedCents: number;
    topGifts: {
        giftType: GiftType;
        count: number;
        totalCents: number;
    }[];
}

export interface SendGiftRequest {
    receiverId: number;
    giftTypeId: number;
    quantity: number;
    message?: string;
    orderId?: number;
}

// ============ State & Actions ============

export interface GiftState {
    giftTypes: GiftType[];
    receivedGifts: Gift[];
    sentGifts: Gift[];
    stats: GiftStats | null;
    loading: boolean;
    error: string | null;
    pagination: {
        page: number;
        pageSize: number;
        total: number;
        hasMore: boolean;
    };
}

export interface GiftActions {
    // Gift types
    fetchGiftTypes: () => Promise<void>;

    // Received gifts (for players)
    fetchReceivedGifts: (refresh?: boolean) => Promise<void>;
    fetchGiftStats: () => Promise<void>;

    // Sent gifts (for users)
    fetchSentGifts: (refresh?: boolean) => Promise<void>;
    sendGift: (request: SendGiftRequest) => Promise<void>;

    // Helpers
    getGiftTypeById: (id: number) => GiftType | undefined;
    calculateGiftTotal: (giftTypeId: number, quantity: number) => number;
}

// ============ Store ============

export const useGiftStore = create<GiftState & GiftActions>((set, get) => ({
    giftTypes: [],
    receivedGifts: [],
    sentGifts: [],
    stats: null,
    loading: false,
    error: null,
    pagination: {
        page: 1,
        pageSize: 20,
        total: 0,
        hasMore: true
    },

    fetchGiftTypes: async () => {
        try {
            const data = await http.get<GiftType[]>('/public/gifts/types');
            set({ giftTypes: data || [] });
        } catch (err) {
            logError('fetchGiftTypes', err);
        }
    },

    fetchReceivedGifts: async (refresh = false) => {
        set({ loading: true, error: null });
        const { pagination, receivedGifts } = get();
        const currentPage = refresh ? 1 : pagination.page;

        try {
            const data = await http.get<{ items: Gift[]; total: number }>('/player/gifts', {
                params: { page: currentPage, pageSize: pagination.pageSize }
            });

            const newGifts = data.items || [];
            const total = data.total || 0;

            set({
                receivedGifts: refresh ? newGifts : [...receivedGifts, ...newGifts],
                pagination: {
                    ...pagination,
                    page: currentPage,
                    total,
                    hasMore: receivedGifts.length + newGifts.length < total
                },
                loading: false
            });
        } catch (err) {
            logError('fetchReceivedGifts', err);
            set({ loading: false, error: getErrorMessage(err, 'Failed to fetch gifts') });
        }
    },

    fetchGiftStats: async () => {
        try {
            const data = await http.get<GiftStats>('/player/gifts/stats');
            set({ stats: data });
        } catch (err) {
            logError('fetchGiftStats', err);
        }
    },

    fetchSentGifts: async (refresh = false) => {
        set({ loading: true, error: null });
        const { pagination, sentGifts } = get();
        const currentPage = refresh ? 1 : pagination.page;

        try {
            const data = await http.get<{ items: Gift[]; total: number }>('/user/gifts/sent', {
                params: { page: currentPage, pageSize: pagination.pageSize }
            });

            const newGifts = data.items || [];
            const total = data.total || 0;

            set({
                sentGifts: refresh ? newGifts : [...sentGifts, ...newGifts],
                pagination: {
                    ...pagination,
                    page: currentPage,
                    total,
                    hasMore: sentGifts.length + newGifts.length < total
                },
                loading: false
            });
        } catch (err) {
            logError('fetchSentGifts', err);
            set({ loading: false, error: getErrorMessage(err, 'Failed to fetch sent gifts') });
        }
    },

    sendGift: async (request) => {
        set({ loading: true, error: null });
        try {
            await http.post('/user/gifts', request);
            set({ loading: false });
        } catch (err) {
            logError('sendGift', err);
            set({ loading: false, error: getErrorMessage(err, 'Failed to send gift') });
            throw err;
        }
    },

    getGiftTypeById: (id) => {
        return get().giftTypes.find(g => g.id === id);
    },

    calculateGiftTotal: (giftTypeId, quantity) => {
        const giftType = get().getGiftTypeById(giftTypeId);
        if (!giftType) return 0;
        return giftType.priceCents * quantity;
    }
}));
