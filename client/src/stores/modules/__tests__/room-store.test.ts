import { describe, it, expect, vi, beforeEach } from 'vitest';
import { useRoomStore, GameRoom, RoomMember } from '../room-store';

// Mock http module
vi.mock('@/lib/http', () => ({
    http: {
        get: vi.fn(),
        post: vi.fn(),
        put: vi.fn(),
        delete: vi.fn(),
    },
}));

import { http } from '@/lib/http';

const mockRoom: GameRoom = {
    id: 1,
    name: 'Test Room',
    groupType: 'team',
    roomStatus: 'waiting',
    gameId: 1,
    gameName: 'League of Legends',
    createdBy: 100,
    hostNickname: 'TestHost',
    maxMembers: 5,
    currentMembers: 2,
    isPrivate: false,
    description: 'A test room',
    voiceEnabled: true,
    voiceRoomId: 'voice-123',
    createdAt: '2024-01-01T00:00:00Z',
    updatedAt: '2024-01-01T00:00:00Z',
};

const mockMember: RoomMember = {
    id: 1,
    groupId: 1,
    userId: 100,
    nickname: 'TestUser',
    avatarUrl: 'https://example.com/avatar.png',
    role: 'owner',
    isReady: true,
    isActive: true,
    joinedAt: '2024-01-01T00:00:00Z',
};

describe('room-store', () => {
    beforeEach(() => {
        // Reset store state
        useRoomStore.getState().reset();
        vi.clearAllMocks();
    });

    describe('fetchRooms', () => {
        it('should fetch rooms successfully', async () => {
            const mockResponse = {
                items: [mockRoom],
                pagination: { page: 1, pageSize: 20, total: 1 },
            };
            vi.mocked(http.get).mockResolvedValueOnce(mockResponse);

            await useRoomStore.getState().fetchRooms();

            expect(http.get).toHaveBeenCalledWith('/user/rooms?');
            expect(useRoomStore.getState().rooms).toEqual([mockRoom]);
            expect(useRoomStore.getState().pagination).toEqual(mockResponse.pagination);
            expect(useRoomStore.getState().isLoading).toBe(false);
        });

        it('should fetch rooms with filters', async () => {
            const mockResponse = {
                items: [mockRoom],
                pagination: { page: 1, pageSize: 10, total: 1 },
            };
            vi.mocked(http.get).mockResolvedValueOnce(mockResponse);

            await useRoomStore.getState().fetchRooms({
                page: 1,
                pageSize: 10,
                gameId: 1,
                groupType: 'team',
                status: 'waiting',
            });

            expect(http.get).toHaveBeenCalledWith(
                '/user/rooms?page=1&pageSize=10&gameId=1&groupType=team&status=waiting'
            );
        });

        it('should handle fetch rooms error', async () => {
            const error = new Error('Network error');
            vi.mocked(http.get).mockRejectedValueOnce(error);

            await expect(useRoomStore.getState().fetchRooms()).rejects.toThrow('Network error');
            expect(useRoomStore.getState().error).toBe('Network error');
            expect(useRoomStore.getState().isLoading).toBe(false);
        });
    });

    describe('fetchRoom', () => {
        it('should fetch single room successfully', async () => {
            vi.mocked(http.get).mockResolvedValueOnce(mockRoom);

            const result = await useRoomStore.getState().fetchRoom(1);

            expect(http.get).toHaveBeenCalledWith('/user/rooms/1');
            expect(result).toEqual(mockRoom);
            expect(useRoomStore.getState().currentRoom).toEqual(mockRoom);
            expect(useRoomStore.getState().isLoading).toBe(false);
        });

        it('should handle fetch room error', async () => {
            const error = new Error('Room not found');
            vi.mocked(http.get).mockRejectedValueOnce(error);

            await expect(useRoomStore.getState().fetchRoom(999)).rejects.toThrow('Room not found');
            expect(useRoomStore.getState().error).toBe('Room not found');
        });
    });

    describe('fetchMyRooms', () => {
        it('should fetch my rooms successfully', async () => {
            vi.mocked(http.get).mockResolvedValueOnce([mockRoom]);

            await useRoomStore.getState().fetchMyRooms();

            expect(http.get).toHaveBeenCalledWith('/user/rooms/my');
            expect(useRoomStore.getState().rooms).toEqual([mockRoom]);
        });
    });

    describe('createRoom', () => {
        it('should create room successfully', async () => {
            vi.mocked(http.post).mockResolvedValueOnce(mockRoom);

            const result = await useRoomStore.getState().createRoom({
                name: 'Test Room',
                groupType: 'team',
                gameId: 1,
                maxMembers: 5,
            });

            expect(http.post).toHaveBeenCalledWith('/user/rooms', {
                name: 'Test Room',
                groupType: 'team',
                gameId: 1,
                maxMembers: 5,
            });
            expect(result).toEqual(mockRoom);
            expect(useRoomStore.getState().rooms).toContainEqual(mockRoom);
            expect(useRoomStore.getState().currentRoom).toEqual(mockRoom);
        });

        it('should handle create room error', async () => {
            const error = new Error('Failed to create');
            vi.mocked(http.post).mockRejectedValueOnce(error);

            await expect(
                useRoomStore.getState().createRoom({
                    name: 'Test',
                    groupType: 'team',
                    gameId: 1,
                })
            ).rejects.toThrow('Failed to create');
            expect(useRoomStore.getState().error).toBe('Failed to create');
        });
    });

    describe('updateRoom', () => {
        it('should update room successfully', async () => {
            vi.mocked(http.put).mockResolvedValueOnce(undefined);
            vi.mocked(http.get).mockResolvedValueOnce({ ...mockRoom, name: 'Updated Room' });

            await useRoomStore.getState().updateRoom(1, { name: 'Updated Room' });

            expect(http.put).toHaveBeenCalledWith('/user/rooms/1', { name: 'Updated Room' });
            expect(useRoomStore.getState().currentRoom?.name).toBe('Updated Room');
        });
    });

    describe('closeRoom', () => {
        it('should close room successfully', async () => {
            // Set initial state with room
            useRoomStore.setState({ rooms: [mockRoom], currentRoom: mockRoom });
            vi.mocked(http.delete).mockResolvedValueOnce(undefined);

            await useRoomStore.getState().closeRoom(1);

            expect(http.delete).toHaveBeenCalledWith('/user/rooms/1');
            expect(useRoomStore.getState().rooms).toEqual([]);
            expect(useRoomStore.getState().currentRoom).toBeNull();
        });

        it('should not clear currentRoom if different room is closed', async () => {
            const otherRoom = { ...mockRoom, id: 2 };
            useRoomStore.setState({ rooms: [mockRoom, otherRoom], currentRoom: otherRoom });
            vi.mocked(http.delete).mockResolvedValueOnce(undefined);

            await useRoomStore.getState().closeRoom(1);

            expect(useRoomStore.getState().currentRoom).toEqual(otherRoom);
        });
    });

    describe('joinRoom', () => {
        it('should join room successfully', async () => {
            vi.mocked(http.post).mockResolvedValueOnce(undefined);
            vi.mocked(http.get)
                .mockResolvedValueOnce(mockRoom) // fetchRoom
                .mockResolvedValueOnce([mockMember]); // fetchMembers

            await useRoomStore.getState().joinRoom(1);

            expect(http.post).toHaveBeenCalledWith('/user/rooms/1/join', { password: undefined });
            expect(useRoomStore.getState().currentRoom).toEqual(mockRoom);
            expect(useRoomStore.getState().members).toEqual([mockMember]);
        });

        it('should join private room with password', async () => {
            vi.mocked(http.post).mockResolvedValueOnce(undefined);
            vi.mocked(http.get)
                .mockResolvedValueOnce(mockRoom)
                .mockResolvedValueOnce([mockMember]);

            await useRoomStore.getState().joinRoom(1, 'secret123');

            expect(http.post).toHaveBeenCalledWith('/user/rooms/1/join', { password: 'secret123' });
        });

        it('should handle join room error', async () => {
            const error = new Error('Wrong password');
            vi.mocked(http.post).mockRejectedValueOnce(error);

            await expect(useRoomStore.getState().joinRoom(1, 'wrong')).rejects.toThrow('Wrong password');
            expect(useRoomStore.getState().error).toBe('Wrong password');
        });
    });

    describe('leaveRoom', () => {
        it('should leave room successfully', async () => {
            useRoomStore.setState({ currentRoom: mockRoom, members: [mockMember] });
            vi.mocked(http.post).mockResolvedValueOnce(undefined);

            await useRoomStore.getState().leaveRoom(1);

            expect(http.post).toHaveBeenCalledWith('/user/rooms/1/leave');
            expect(useRoomStore.getState().currentRoom).toBeNull();
            expect(useRoomStore.getState().members).toEqual([]);
        });
    });

    describe('toggleReady', () => {
        it('should toggle ready status successfully', async () => {
            vi.mocked(http.post).mockResolvedValueOnce(undefined);
            vi.mocked(http.get).mockResolvedValueOnce([{ ...mockMember, isReady: false }]);

            await useRoomStore.getState().toggleReady(1);

            expect(http.post).toHaveBeenCalledWith('/user/rooms/1/ready');
        });
    });

    describe('startGame', () => {
        it('should start game successfully', async () => {
            vi.mocked(http.post).mockResolvedValueOnce(undefined);
            vi.mocked(http.get).mockResolvedValueOnce({ ...mockRoom, roomStatus: 'in_game' });

            await useRoomStore.getState().startGame(1);

            expect(http.post).toHaveBeenCalledWith('/user/rooms/1/start');
            expect(useRoomStore.getState().currentRoom?.roomStatus).toBe('in_game');
        });
    });

    describe('finishGame', () => {
        it('should finish game successfully', async () => {
            vi.mocked(http.post).mockResolvedValueOnce(undefined);
            vi.mocked(http.get).mockResolvedValueOnce({ ...mockRoom, roomStatus: 'finished' });

            await useRoomStore.getState().finishGame(1);

            expect(http.post).toHaveBeenCalledWith('/user/rooms/1/finish');
            expect(useRoomStore.getState().currentRoom?.roomStatus).toBe('finished');
        });
    });

    describe('kickMember', () => {
        it('should kick member successfully', async () => {
            vi.mocked(http.post).mockResolvedValueOnce(undefined);
            vi.mocked(http.get)
                .mockResolvedValueOnce([]) // fetchMembers
                .mockResolvedValueOnce({ ...mockRoom, currentMembers: 1 }); // fetchRoom

            await useRoomStore.getState().kickMember(1, 200);

            expect(http.post).toHaveBeenCalledWith('/user/rooms/1/kick/200');
        });
    });

    describe('fetchMembers', () => {
        it('should fetch members successfully', async () => {
            vi.mocked(http.get).mockResolvedValueOnce([mockMember]);

            await useRoomStore.getState().fetchMembers(1);

            expect(http.get).toHaveBeenCalledWith('/user/rooms/1/members');
            expect(useRoomStore.getState().members).toEqual([mockMember]);
        });
    });

    describe('setCurrentRoom', () => {
        it('should set current room', () => {
            useRoomStore.getState().setCurrentRoom(mockRoom);
            expect(useRoomStore.getState().currentRoom).toEqual(mockRoom);
        });

        it('should clear current room', () => {
            useRoomStore.setState({ currentRoom: mockRoom });
            useRoomStore.getState().setCurrentRoom(null);
            expect(useRoomStore.getState().currentRoom).toBeNull();
        });
    });

    describe('reset', () => {
        it('should reset store to initial state', () => {
            useRoomStore.setState({
                rooms: [mockRoom],
                currentRoom: mockRoom,
                members: [mockMember],
                isLoading: true,
                error: 'Some error',
            });

            useRoomStore.getState().reset();

            expect(useRoomStore.getState().rooms).toEqual([]);
            expect(useRoomStore.getState().currentRoom).toBeNull();
            expect(useRoomStore.getState().members).toEqual([]);
            expect(useRoomStore.getState().isLoading).toBe(false);
            expect(useRoomStore.getState().error).toBeNull();
        });
    });
});
