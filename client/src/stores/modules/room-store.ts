import { create } from 'zustand';
import { http } from '@/lib/http';

// ============================================================================
// Types
// ============================================================================

export type ChatGroupType = 'public' | 'order' | 'team' | 'lfg' | 'custom';
export type ChatGroupStatus = 'waiting' | 'ready' | 'in_game' | 'finished' | 'canceled';

export interface GameRoom {
    id: number;
    name: string;
    groupType: ChatGroupType;
    roomStatus: ChatGroupStatus;
    gameId: number;
    gameName?: string;
    createdBy: number;
    hostNickname?: string;
    maxMembers: number;
    currentMembers: number;
    isPrivate: boolean;
    description?: string;
    voiceEnabled: boolean;
    voiceRoomId?: string;
    createdAt: string;
    updatedAt: string;
}

export interface RoomMember {
    id: number;
    groupId: number;
    userId: number;
    nickname: string;
    avatarUrl?: string;
    role: 'owner' | 'admin' | 'member';
    isReady: boolean;
    isActive: boolean;
    joinedAt: string;
}

export interface CreateRoomRequest {
    name: string;
    groupType: ChatGroupType;
    gameId: number;
    maxMembers?: number;
    isPrivate?: boolean;
    password?: string;
    description?: string;
    voiceEnabled?: boolean;
}

export interface UpdateRoomRequest {
    name?: string;
    maxMembers?: number;
    isPrivate?: boolean;
    password?: string;
    description?: string;
}

export interface RoomListOptions {
    page?: number;
    pageSize?: number;
    gameId?: number;
    groupType?: ChatGroupType;
    status?: ChatGroupStatus;
}

export interface RoomState {
    // State
    rooms: GameRoom[];
    currentRoom: GameRoom | null;
    members: RoomMember[];
    isLoading: boolean;
    error: string | null;
    pagination: {
        page: number;
        pageSize: number;
        total: number;
    };

    // Actions
    fetchRooms: (options?: RoomListOptions) => Promise<void>;
    fetchRoom: (roomId: number) => Promise<GameRoom>;
    fetchMyRooms: () => Promise<void>;
    createRoom: (data: CreateRoomRequest) => Promise<GameRoom>;
    updateRoom: (roomId: number, data: UpdateRoomRequest) => Promise<void>;
    closeRoom: (roomId: number) => Promise<void>;
    joinRoom: (roomId: number, password?: string) => Promise<void>;
    leaveRoom: (roomId: number) => Promise<void>;
    toggleReady: (roomId: number) => Promise<void>;
    startGame: (roomId: number) => Promise<void>;
    finishGame: (roomId: number) => Promise<void>;
    kickMember: (roomId: number, userId: number) => Promise<void>;
    fetchMembers: (roomId: number) => Promise<void>;
    setCurrentRoom: (room: GameRoom | null) => void;
    reset: () => void;
}

// ============================================================================
// Initial State
// ============================================================================

const initialState = {
    rooms: [],
    currentRoom: null,
    members: [],
    isLoading: false,
    error: null,
    pagination: {
        page: 1,
        pageSize: 20,
        total: 0,
    },
};

// ============================================================================
// Store
// ============================================================================

export const useRoomStore = create<RoomState>((set, get) => ({
    ...initialState,

    // Fetch room list
    fetchRooms: async (options?: RoomListOptions) => {
        set({ isLoading: true, error: null });
        try {
            const params = new URLSearchParams();
            if (options?.page) params.append('page', String(options.page));
            if (options?.pageSize) params.append('pageSize', String(options.pageSize));
            if (options?.gameId) params.append('gameId', String(options.gameId));
            if (options?.groupType) params.append('groupType', options.groupType);
            if (options?.status) params.append('status', options.status);

            const response = await http.get<{
                items: GameRoom[];
                pagination: { page: number; pageSize: number; total: number };
            }>(`/user/rooms?${params.toString()}`);

            set({
                rooms: response.items || [],
                pagination: response.pagination || initialState.pagination,
                isLoading: false,
            });
        } catch (error: unknown) {
            const message = error instanceof Error ? error.message : 'Failed to fetch rooms';
            set({ error: message, isLoading: false });
            throw error;
        }
    },

    // Fetch single room
    fetchRoom: async (roomId: number) => {
        set({ isLoading: true, error: null });
        try {
            const response = await http.get<GameRoom>(`/user/rooms/${roomId}`);
            set({ currentRoom: response, isLoading: false });
            return response;
        } catch (error: unknown) {
            const message = error instanceof Error ? error.message : 'Failed to fetch room';
            set({ error: message, isLoading: false });
            throw error;
        }
    },

    // Fetch my rooms
    fetchMyRooms: async () => {
        set({ isLoading: true, error: null });
        try {
            const response = await http.get<GameRoom[]>('/user/rooms/my');
            set({ rooms: response || [], isLoading: false });
        } catch (error: unknown) {
            const message = error instanceof Error ? error.message : 'Failed to fetch my rooms';
            set({ error: message, isLoading: false });
            throw error;
        }
    },

    // Create room
    createRoom: async (data: CreateRoomRequest) => {
        set({ isLoading: true, error: null });
        try {
            const response = await http.post<GameRoom>('/user/rooms', data);
            set((state) => ({
                rooms: [response, ...state.rooms],
                currentRoom: response,
                isLoading: false,
            }));
            return response;
        } catch (error: unknown) {
            const message = error instanceof Error ? error.message : 'Failed to create room';
            set({ error: message, isLoading: false });
            throw error;
        }
    },

    // Update room
    updateRoom: async (roomId: number, data: UpdateRoomRequest) => {
        try {
            await http.put(`/user/rooms/${roomId}`, data);
            // Refresh room data
            await get().fetchRoom(roomId);
        } catch (error: unknown) {
            const message = error instanceof Error ? error.message : 'Failed to update room';
            set({ error: message });
            throw error;
        }
    },

    // Close room
    closeRoom: async (roomId: number) => {
        try {
            await http.delete(`/user/rooms/${roomId}`);
            set((state) => ({
                rooms: state.rooms.filter((r) => r.id !== roomId),
                currentRoom:
                    state.currentRoom?.id === roomId ? null : state.currentRoom,
            }));
        } catch (error: unknown) {
            const message = error instanceof Error ? error.message : 'Failed to close room';
            set({ error: message });
            throw error;
        }
    },

    // Join room
    joinRoom: async (roomId: number, password?: string) => {
        set({ isLoading: true, error: null });
        try {
            await http.post(`/user/rooms/${roomId}/join`, { password });
            // Refresh room data
            await get().fetchRoom(roomId);
            await get().fetchMembers(roomId);
            set({ isLoading: false });
        } catch (error: unknown) {
            const message = error instanceof Error ? error.message : 'Failed to join room';
            set({ error: message, isLoading: false });
            throw error;
        }
    },

    // Leave room
    leaveRoom: async (roomId: number) => {
        try {
            await http.post(`/user/rooms/${roomId}/leave`);
            set((state) => ({
                currentRoom:
                    state.currentRoom?.id === roomId ? null : state.currentRoom,
                members: [],
            }));
        } catch (error: unknown) {
            const message = error instanceof Error ? error.message : 'Failed to leave room';
            set({ error: message });
            throw error;
        }
    },

    // Toggle ready status
    toggleReady: async (roomId: number) => {
        try {
            await http.post(`/user/rooms/${roomId}/ready`);
            // Refresh members
            await get().fetchMembers(roomId);
        } catch (error: unknown) {
            const message = error instanceof Error ? error.message : 'Failed to toggle ready';
            set({ error: message });
            throw error;
        }
    },

    // Start game (host only)
    startGame: async (roomId: number) => {
        try {
            await http.post(`/user/rooms/${roomId}/start`);
            await get().fetchRoom(roomId);
        } catch (error: unknown) {
            const message = error instanceof Error ? error.message : 'Failed to start game';
            set({ error: message });
            throw error;
        }
    },

    // Finish game (host only)
    finishGame: async (roomId: number) => {
        try {
            await http.post(`/user/rooms/${roomId}/finish`);
            await get().fetchRoom(roomId);
        } catch (error: unknown) {
            const message = error instanceof Error ? error.message : 'Failed to finish game';
            set({ error: message });
            throw error;
        }
    },

    // Kick member (host only)
    kickMember: async (roomId: number, userId: number) => {
        try {
            await http.post(`/user/rooms/${roomId}/kick/${userId}`);
            await get().fetchMembers(roomId);
            await get().fetchRoom(roomId);
        } catch (error: unknown) {
            const message = error instanceof Error ? error.message : 'Failed to kick member';
            set({ error: message });
            throw error;
        }
    },

    // Fetch room members
    fetchMembers: async (roomId: number) => {
        try {
            const response = await http.get<RoomMember[]>(
                `/user/rooms/${roomId}/members`
            );
            set({ members: response || [] });
        } catch (error: unknown) {
            const message = error instanceof Error ? error.message : 'Failed to fetch members';
            set({ error: message });
            throw error;
        }
    },

    // Set current room
    setCurrentRoom: (room: GameRoom | null) => {
        set({ currentRoom: room });
    },

    // Reset state
    reset: () => {
        set(initialState);
    },
}));

// ============================================================================
// Selectors
// ============================================================================

export const selectRooms = (state: RoomState) => state.rooms;
export const selectCurrentRoom = (state: RoomState) => state.currentRoom;
export const selectMembers = (state: RoomState) => state.members;
export const selectIsLoading = (state: RoomState) => state.isLoading;
export const selectRoomError = (state: RoomState) => state.error;
