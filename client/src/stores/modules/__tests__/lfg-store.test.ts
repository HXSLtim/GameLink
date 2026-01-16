import { describe, it, expect, vi, beforeEach } from 'vitest';
import { useLFGStore, type LFGRequest } from '../lfg-store';

// Mock http module
vi.mock('@/lib/http', () => ({
    http: {
        get: vi.fn(),
        post: vi.fn(),
        delete: vi.fn(),
    },
}));

import { http } from '@/lib/http';

const mockLFGRequest: LFGRequest = {
    id: 1,
    userId: 100,
    userNickname: 'TestUser',
    userAvatarUrl: 'https://example.com/avatar.png',
    gameId: 1,
    gameName: 'League of Legends',
    requestType: 'find_player',
    title: 'Looking for duo partner',
    description: 'Need a support main for ranked',
    requiredPlayers: 1,
    minRank: 'Diamond',
    maxPriceCents: 5000,
    status: 'pending',
    expiresAt: '2024-01-01T01:00:00Z',
    matchedRoomId: undefined,
    createdAt: '2024-01-01T00:00:00Z',
    updatedAt: '2024-01-01T00:00:00Z',
};

describe('lfg-store', () => {
    beforeEach(() => {
        useLFGStore.getState().reset();
        vi.clearAllMocks();
    });

    describe('fetchRequests', () => {
        it('should fetch LFG requests successfully', async () => {
            const mockResponse = {
                items: [mockLFGRequest],
                pagination: { page: 1, pageSize: 20, total: 1 },
            };
            vi.mocked(http.get).mockResolvedValueOnce(mockResponse);

            await useLFGStore.getState().fetchRequests();

            expect(http.get).toHaveBeenCalledWith('/user/lfg?');
            expect(useLFGStore.getState().requests).toEqual([mockLFGRequest]);
            expect(useLFGStore.getState().pagination).toEqual(mockResponse.pagination);
            expect(useLFGStore.getState().isLoading).toBe(false);
        });

        it('should fetch requests with filters', async () => {
            const mockResponse = {
                items: [mockLFGRequest],
                pagination: { page: 1, pageSize: 10, total: 1 },
            };
            vi.mocked(http.get).mockResolvedValueOnce(mockResponse);

            await useLFGStore.getState().fetchRequests({
                page: 1,
                pageSize: 10,
                gameId: 1,
                requestType: 'find_player',
            });

            expect(http.get).toHaveBeenCalledWith(
                '/user/lfg?page=1&pageSize=10&gameId=1&requestType=find_player'
            );
        });

        it('should handle fetch requests error', async () => {
            const error = new Error('Network error');
            vi.mocked(http.get).mockRejectedValueOnce(error);

            await expect(useLFGStore.getState().fetchRequests()).rejects.toThrow('Network error');
            expect(useLFGStore.getState().error).toBe('Network error');
            expect(useLFGStore.getState().isLoading).toBe(false);
        });
    });

    describe('fetchPendingRequests', () => {
        it('should fetch pending requests successfully', async () => {
            const mockResponse = {
                items: [mockLFGRequest],
                pagination: { page: 1, pageSize: 20, total: 1 },
            };
            vi.mocked(http.get).mockResolvedValueOnce(mockResponse);

            await useLFGStore.getState().fetchPendingRequests();

            expect(http.get).toHaveBeenCalledWith('/user/lfg/pending?');
            expect(useLFGStore.getState().requests).toEqual([mockLFGRequest]);
        });

        it('should fetch pending requests with gameId filter', async () => {
            const mockResponse = {
                items: [mockLFGRequest],
                pagination: { page: 1, pageSize: 20, total: 1 },
            };
            vi.mocked(http.get).mockResolvedValueOnce(mockResponse);

            await useLFGStore.getState().fetchPendingRequests(1);

            expect(http.get).toHaveBeenCalledWith('/user/lfg/pending?gameId=1');
        });
    });

    describe('fetchMyRequests', () => {
        it('should fetch my requests successfully', async () => {
            vi.mocked(http.get).mockResolvedValueOnce([mockLFGRequest]);

            await useLFGStore.getState().fetchMyRequests();

            expect(http.get).toHaveBeenCalledWith('/user/lfg/my?');
            expect(useLFGStore.getState().myRequests).toEqual([mockLFGRequest]);
        });

        it('should fetch my requests with status filter', async () => {
            vi.mocked(http.get).mockResolvedValueOnce([mockLFGRequest]);

            await useLFGStore.getState().fetchMyRequests('pending');

            expect(http.get).toHaveBeenCalledWith('/user/lfg/my?status=pending');
        });
    });

    describe('fetchActiveRequest', () => {
        it('should fetch active request successfully', async () => {
            vi.mocked(http.get).mockResolvedValueOnce(mockLFGRequest);

            await useLFGStore.getState().fetchActiveRequest();

            expect(http.get).toHaveBeenCalledWith('/user/lfg/active');
            expect(useLFGStore.getState().activeRequest).toEqual(mockLFGRequest);
        });

        it('should handle 404 (no active request) gracefully', async () => {
            const error = { status: 404, message: 'Not found' };
            vi.mocked(http.get).mockRejectedValueOnce(error);

            await useLFGStore.getState().fetchActiveRequest();

            expect(useLFGStore.getState().activeRequest).toBeNull();
            expect(useLFGStore.getState().error).toBeNull();
        });

        it('should handle other errors', async () => {
            const error = new Error('Server error');
            vi.mocked(http.get).mockRejectedValueOnce(error);

            await expect(useLFGStore.getState().fetchActiveRequest()).rejects.toThrow('Server error');
            expect(useLFGStore.getState().error).toBe('Server error');
        });
    });

    describe('fetchRequest', () => {
        it('should fetch single request successfully', async () => {
            vi.mocked(http.get).mockResolvedValueOnce(mockLFGRequest);

            const result = await useLFGStore.getState().fetchRequest(1);

            expect(http.get).toHaveBeenCalledWith('/user/lfg/1');
            expect(result).toEqual(mockLFGRequest);
        });

        it('should handle fetch request error', async () => {
            const error = new Error('Request not found');
            vi.mocked(http.get).mockRejectedValueOnce(error);

            await expect(useLFGStore.getState().fetchRequest(999)).rejects.toThrow('Request not found');
            expect(useLFGStore.getState().error).toBe('Request not found');
        });
    });

    describe('createRequest', () => {
        it('should create request successfully', async () => {
            vi.mocked(http.post).mockResolvedValueOnce(mockLFGRequest);

            const result = await useLFGStore.getState().createRequest({
                gameId: 1,
                requestType: 'find_player',
                title: 'Looking for duo partner',
                requiredPlayers: 1,
            });

            expect(http.post).toHaveBeenCalledWith('/user/lfg', {
                gameId: 1,
                requestType: 'find_player',
                title: 'Looking for duo partner',
                requiredPlayers: 1,
            });
            expect(result).toEqual(mockLFGRequest);
            expect(useLFGStore.getState().requests).toContainEqual(mockLFGRequest);
            expect(useLFGStore.getState().myRequests).toContainEqual(mockLFGRequest);
            expect(useLFGStore.getState().activeRequest).toEqual(mockLFGRequest);
        });

        it('should handle create request error', async () => {
            const error = new Error('Failed to create');
            vi.mocked(http.post).mockRejectedValueOnce(error);

            await expect(
                useLFGStore.getState().createRequest({
                    gameId: 1,
                    requestType: 'find_player',
                })
            ).rejects.toThrow('Failed to create');
            expect(useLFGStore.getState().error).toBe('Failed to create');
        });
    });

    describe('cancelRequest', () => {
        it('should cancel request successfully', async () => {
            useLFGStore.setState({
                requests: [mockLFGRequest],
                myRequests: [mockLFGRequest],
                activeRequest: mockLFGRequest,
            });
            vi.mocked(http.delete).mockResolvedValueOnce(undefined);

            await useLFGStore.getState().cancelRequest(1);

            expect(http.delete).toHaveBeenCalledWith('/user/lfg/1');
            expect(useLFGStore.getState().requests).toEqual([]);
            expect(useLFGStore.getState().myRequests).toEqual([]);
            expect(useLFGStore.getState().activeRequest).toBeNull();
        });

        it('should not clear activeRequest if different request is canceled', async () => {
            const otherRequest = { ...mockLFGRequest, id: 2 };
            useLFGStore.setState({
                requests: [mockLFGRequest, otherRequest],
                myRequests: [mockLFGRequest, otherRequest],
                activeRequest: otherRequest,
            });
            vi.mocked(http.delete).mockResolvedValueOnce(undefined);

            await useLFGStore.getState().cancelRequest(1);

            expect(useLFGStore.getState().activeRequest).toEqual(otherRequest);
        });

        it('should handle cancel request error', async () => {
            const error = new Error('Failed to cancel');
            vi.mocked(http.delete).mockRejectedValueOnce(error);

            await expect(useLFGStore.getState().cancelRequest(1)).rejects.toThrow('Failed to cancel');
            expect(useLFGStore.getState().error).toBe('Failed to cancel');
        });
    });

    describe('acceptRequest', () => {
        it('should accept request successfully', async () => {
            useLFGStore.setState({ requests: [mockLFGRequest] });
            const mockResponse = { roomId: 123 };
            vi.mocked(http.post).mockResolvedValueOnce(mockResponse);

            const result = await useLFGStore.getState().acceptRequest(1);

            expect(http.post).toHaveBeenCalledWith('/user/lfg/1/accept');
            expect(result).toEqual(mockResponse);
            expect(useLFGStore.getState().requests).toEqual([]);
        });

        it('should handle accept request error', async () => {
            const error = new Error('Request already matched');
            vi.mocked(http.post).mockRejectedValueOnce(error);

            await expect(useLFGStore.getState().acceptRequest(1)).rejects.toThrow('Request already matched');
            expect(useLFGStore.getState().error).toBe('Request already matched');
        });
    });

    describe('findMatches', () => {
        it('should find matches successfully', async () => {
            const matchedRequests = [mockLFGRequest, { ...mockLFGRequest, id: 2 }];
            vi.mocked(http.get).mockResolvedValueOnce(matchedRequests);

            await useLFGStore.getState().findMatches(1);

            expect(http.get).toHaveBeenCalledWith('/user/lfg/1/matches?');
            expect(useLFGStore.getState().matches).toEqual(matchedRequests);
        });

        it('should find matches with limit', async () => {
            vi.mocked(http.get).mockResolvedValueOnce([mockLFGRequest]);

            await useLFGStore.getState().findMatches(1, 5);

            expect(http.get).toHaveBeenCalledWith('/user/lfg/1/matches?limit=5');
        });

        it('should handle find matches error', async () => {
            const error = new Error('Failed to find matches');
            vi.mocked(http.get).mockRejectedValueOnce(error);

            await expect(useLFGStore.getState().findMatches(1)).rejects.toThrow('Failed to find matches');
            expect(useLFGStore.getState().error).toBe('Failed to find matches');
        });
    });

    describe('fetchPendingCount', () => {
        it('should fetch pending count successfully', async () => {
            vi.mocked(http.get).mockResolvedValueOnce({ count: 10 });

            await useLFGStore.getState().fetchPendingCount();

            expect(http.get).toHaveBeenCalledWith('/user/lfg/count?');
            expect(useLFGStore.getState().pendingCount).toBe(10);
        });

        it('should fetch pending count with gameId', async () => {
            vi.mocked(http.get).mockResolvedValueOnce({ count: 5 });

            await useLFGStore.getState().fetchPendingCount(1);

            expect(http.get).toHaveBeenCalledWith('/user/lfg/count?gameId=1');
            expect(useLFGStore.getState().pendingCount).toBe(5);
        });

        it('should handle fetch pending count error silently', async () => {
            const consoleSpy = vi.spyOn(console, 'error').mockImplementation(() => { });
            vi.mocked(http.get).mockRejectedValueOnce(new Error('Network error'));

            await useLFGStore.getState().fetchPendingCount();

            expect(consoleSpy).toHaveBeenCalled();
            consoleSpy.mockRestore();
        });
    });

    describe('reset', () => {
        it('should reset store to initial state', () => {
            useLFGStore.setState({
                requests: [mockLFGRequest],
                myRequests: [mockLFGRequest],
                activeRequest: mockLFGRequest,
                matches: [mockLFGRequest],
                pendingCount: 10,
                isLoading: true,
                error: 'Some error',
            });

            useLFGStore.getState().reset();

            expect(useLFGStore.getState().requests).toEqual([]);
            expect(useLFGStore.getState().myRequests).toEqual([]);
            expect(useLFGStore.getState().activeRequest).toBeNull();
            expect(useLFGStore.getState().matches).toEqual([]);
            expect(useLFGStore.getState().pendingCount).toBe(0);
            expect(useLFGStore.getState().isLoading).toBe(false);
            expect(useLFGStore.getState().error).toBeNull();
        });
    });
});
