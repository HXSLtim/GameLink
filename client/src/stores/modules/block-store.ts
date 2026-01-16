import { create } from 'zustand';
import { http } from '@/lib/http';
import { getErrorMessage, logError } from '@/lib/error';

// ============ Interfaces ============

export interface BlockedUser {
    id: number;
    blockedUserId: number;
    blockedUser: {
        id: number;
        username: string;
        nickname: string;
        avatar: string;
    };
    reason?: string;
    createdAt: string;
}

export interface BlockUserRequest {
    userId: number;
    reason?: string;
}

// ============ State & Actions ============

export interface BlockState {
    blockedUsers: BlockedUser[];
    blockedUserIds: Set<number>;
    loading: boolean;
    error: string | null;
}

export interface BlockActions {
    fetchBlockedUsers: () => Promise<void>;
    blockUser: (request: BlockUserRequest) => Promise<void>;
    unblockUser: (userId: number) => Promise<void>;
    isBlocked: (userId: number) => boolean;
    checkBlocked: (userId: number) => Promise<boolean>;
}

// ============ Store ============

export const useBlockStore = create<BlockState & BlockActions>((set, get) => ({
    blockedUsers: [],
    blockedUserIds: new Set(),
    loading: false,
    error: null,

    fetchBlockedUsers: async () => {
        set({ loading: true, error: null });
        try {
            const data = await http.get<{ items: BlockedUser[] }>('/user/blocks');
            const blockedUsers = data.items || [];
            const blockedUserIds = new Set(blockedUsers.map(b => b.blockedUserId));
            set({ blockedUsers, blockedUserIds, loading: false });
        } catch (err) {
            logError('fetchBlockedUsers', err);
            set({ loading: false, error: getErrorMessage(err, 'Failed to fetch blocked users') });
        }
    },

    blockUser: async (request) => {
        set({ loading: true, error: null });
        try {
            await http.post('/user/blocks', request);
            // Refresh the list
            await get().fetchBlockedUsers();
        } catch (err) {
            logError('blockUser', err);
            set({ loading: false, error: getErrorMessage(err, 'Failed to block user') });
            throw err;
        }
    },

    unblockUser: async (userId) => {
        const previousBlockedUsers = [...get().blockedUsers];
        const previousBlockedUserIds = new Set(get().blockedUserIds);

        // Optimistic update
        set((state) => ({
            blockedUsers: state.blockedUsers.filter(b => b.blockedUserId !== userId),
            blockedUserIds: new Set([...state.blockedUserIds].filter(id => id !== userId))
        }));

        try {
            await http.delete(`/user/blocks/${userId}`);
        } catch (err) {
            // Rollback on error
            logError('unblockUser', err);
            set({
                blockedUsers: previousBlockedUsers,
                blockedUserIds: previousBlockedUserIds,
                error: getErrorMessage(err, 'Failed to unblock user')
            });
            throw err;
        }
    },

    isBlocked: (userId) => {
        return get().blockedUserIds.has(userId);
    },

    checkBlocked: async (userId) => {
        try {
            const data = await http.get<{ blocked: boolean }>(`/user/blocks/check/${userId}`);
            return data.blocked;
        } catch (err) {
            logError('checkBlocked', err);
            return false;
        }
    }
}));
