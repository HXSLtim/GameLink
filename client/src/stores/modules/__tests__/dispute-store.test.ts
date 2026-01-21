/**
 * Dispute Store Tests
 * Tests for dispute management, messages, and selectors
 */

import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import {
    useDisputeStore,
    DisputeStatus,
    DisputeResolution,
    DisputeInitiatorType,
    DisputeType,
    selectActiveDisputes,
    selectResolvedDisputes,
    selectIsDisputable,
    type OrderDispute,
    type DisputeTemplate,
    type DisputeMessage,
    type DisputeState,
} from '../dispute-store';

// Mock http module
vi.mock('@/lib/http', () => ({
    http: {
        get: vi.fn(),
        post: vi.fn(),
        put: vi.fn(),
        delete: vi.fn(),
    },
}));

// Mock error module
vi.mock('@/lib/error', () => ({
    getErrorMessage: (err: unknown, fallback: string) =>
        err instanceof Error ? err.message : fallback,
    logError: vi.fn(),
}));

import { http } from '@/lib/http';

const mockHttp = http as unknown as {
    get: ReturnType<typeof vi.fn>;
    post: ReturnType<typeof vi.fn>;
    put: ReturnType<typeof vi.fn>;
    delete: ReturnType<typeof vi.fn>;
};

// Mock data factories
const createMockDispute = (overrides: Partial<OrderDispute> = {}): OrderDispute => ({
    id: 1,
    orderId: 100,
    orderNo: 'ORD-001',
    initiatorId: 1,
    initiatorType: DisputeInitiatorType.USER,
    type: DisputeType.SERVICE_QUALITY,
    status: DisputeStatus.PENDING,
    reason: 'Service was not as described',
    evidenceUrls: ['https://example.com/evidence1.jpg'],
    evidenceText: 'Additional details',
    slaDeadline: '2024-01-01T12:00:00Z',
    slaBreached: false,
    slaRemaining: 1800,
    resolution: DisputeResolution.PENDING,
    createdAt: '2024-01-01T00:00:00Z',
    updatedAt: '2024-01-01T00:00:00Z',
    ...overrides,
});

const createMockTemplate = (overrides: Partial<DisputeTemplate> = {}): DisputeTemplate => ({
    id: 1,
    code: 'service_quality',
    name: 'Service Quality Issue',
    initiatorType: DisputeInitiatorType.USER,
    description: 'Report issues with service quality',
    sortOrder: 1,
    isActive: true,
    ...overrides,
});

const createMockMessage = (overrides: Partial<DisputeMessage> = {}): DisputeMessage => ({
    id: 1,
    disputeId: 1,
    senderId: 1,
    senderName: 'User',
    senderRole: 'user',
    content: 'Hello, I have an issue',
    createdAt: '2024-01-01T00:00:00Z',
    ...overrides,
});

describe('Dispute Store', () => {
    beforeEach(() => {
        vi.clearAllMocks();
        // Reset store state
        useDisputeStore.getState().reset();
    });

    afterEach(() => {
        vi.restoreAllMocks();
    });

    describe('Initial State', () => {
        it('should have correct initial state', () => {
            const state = useDisputeStore.getState();

            expect(state.myDisputes).toEqual([]);
            expect(state.currentDispute).toBeNull();
            expect(state.templates).toEqual([]);
            expect(state.messages).toEqual([]);
            expect(state.pagination).toEqual({
                page: 1,
                pageSize: 20,
                total: 0,
                hasMore: false,
            });
            expect(state.loading).toBe(false);
            expect(state.error).toBeNull();
        });
    });

    describe('fetchMyDisputes', () => {
        it('should fetch disputes successfully', async () => {
            const mockDisputes = [
                createMockDispute({ id: 1 }),
                createMockDispute({ id: 2 }),
            ];
            mockHttp.get.mockResolvedValueOnce({
                items: mockDisputes,
                total: 2,
            });

            await useDisputeStore.getState().fetchMyDisputes();

            const state = useDisputeStore.getState();
            expect(state.myDisputes).toHaveLength(2);
            expect(state.pagination.total).toBe(2);
            expect(state.loading).toBe(false);
            expect(state.error).toBeNull();
            expect(mockHttp.get).toHaveBeenCalledWith('/user/disputes', {
                params: { page: 1, pageSize: 20 },
            });
        });

        it('should handle fetch disputes error', async () => {
            mockHttp.get.mockRejectedValueOnce(new Error('Network error'));

            await useDisputeStore.getState().fetchMyDisputes();

            const state = useDisputeStore.getState();
            expect(state.myDisputes).toEqual([]);
            expect(state.loading).toBe(false);
            expect(state.error).toBe('Network error');
        });

        it('should append disputes on subsequent pages', async () => {
            const firstPage = [createMockDispute({ id: 1 })];
            const secondPage = [createMockDispute({ id: 2 })];

            mockHttp.get.mockResolvedValueOnce({ items: firstPage, total: 2 });
            await useDisputeStore.getState().fetchMyDisputes(1);

            mockHttp.get.mockResolvedValueOnce({ items: secondPage, total: 2 });
            await useDisputeStore.getState().fetchMyDisputes(2);

            const state = useDisputeStore.getState();
            expect(state.myDisputes).toHaveLength(2);
            expect(state.pagination.page).toBe(2);
        });

        it('should replace disputes on page 1', async () => {
            useDisputeStore.setState({
                myDisputes: [createMockDispute({ id: 99 })],
            });

            const newDisputes = [createMockDispute({ id: 1 })];
            mockHttp.get.mockResolvedValueOnce({ items: newDisputes, total: 1 });

            await useDisputeStore.getState().fetchMyDisputes(1);

            const state = useDisputeStore.getState();
            expect(state.myDisputes).toHaveLength(1);
            expect(state.myDisputes[0].id).toBe(1);
        });

        it('should handle empty response', async () => {
            mockHttp.get.mockResolvedValueOnce({ items: null, total: 0 });

            await useDisputeStore.getState().fetchMyDisputes();

            const state = useDisputeStore.getState();
            expect(state.myDisputes).toEqual([]);
            expect(state.pagination.total).toBe(0);
        });
    });

    describe('fetchDisputeById', () => {
        it('should fetch dispute by id successfully', async () => {
            const mockDispute = createMockDispute({ id: 1 });
            mockHttp.get.mockResolvedValueOnce(mockDispute);

            await useDisputeStore.getState().fetchDisputeById(1);

            const state = useDisputeStore.getState();
            expect(state.currentDispute).toEqual(mockDispute);
            expect(state.loading).toBe(false);
            expect(mockHttp.get).toHaveBeenCalledWith('/user/disputes/1');
        });

        it('should handle fetch dispute error', async () => {
            mockHttp.get.mockRejectedValueOnce(new Error('Not found'));

            await useDisputeStore.getState().fetchDisputeById(999);

            const state = useDisputeStore.getState();
            expect(state.currentDispute).toBeNull();
            expect(state.error).toBe('Not found');
        });
    });

    describe('fetchTemplates', () => {
        it('should fetch templates successfully', async () => {
            const mockTemplates = [
                createMockTemplate({ id: 1 }),
                createMockTemplate({ id: 2, code: 'bad_attitude' }),
            ];
            mockHttp.get.mockResolvedValueOnce(mockTemplates);

            await useDisputeStore.getState().fetchTemplates();

            const state = useDisputeStore.getState();
            expect(state.templates).toHaveLength(2);
            expect(mockHttp.get).toHaveBeenCalledWith('/disputes/templates', { params: {} });
        });

        it('should fetch templates with initiator type filter', async () => {
            mockHttp.get.mockResolvedValueOnce([createMockTemplate()]);

            await useDisputeStore.getState().fetchTemplates(DisputeInitiatorType.USER);

            expect(mockHttp.get).toHaveBeenCalledWith('/disputes/templates', {
                params: { initiatorType: 'user' },
            });
        });

        it('should handle fetch templates error silently', async () => {
            const consoleSpy = vi.spyOn(console, 'error').mockImplementation(() => {});
            mockHttp.get.mockRejectedValueOnce(new Error('Failed'));

            await useDisputeStore.getState().fetchTemplates();

            const state = useDisputeStore.getState();
            expect(state.templates).toEqual([]);
            expect(consoleSpy).toHaveBeenCalled();
            consoleSpy.mockRestore();
        });
    });

    describe('createDispute', () => {
        it('should create dispute successfully', async () => {
            const mockDispute = createMockDispute({ id: 1 });
            mockHttp.post.mockResolvedValueOnce(mockDispute);

            const result = await useDisputeStore.getState().createDispute({
                orderId: 100,
                type: DisputeType.SERVICE_QUALITY,
                reason: 'Service issue',
            });

            const state = useDisputeStore.getState();
            expect(result).toEqual(mockDispute);
            expect(state.myDisputes).toContainEqual(mockDispute);
            expect(state.currentDispute).toEqual(mockDispute);
            expect(state.loading).toBe(false);
            expect(mockHttp.post).toHaveBeenCalledWith('/user/disputes', {
                orderId: 100,
                type: DisputeType.SERVICE_QUALITY,
                reason: 'Service issue',
            });
        });

        it('should handle create dispute error', async () => {
            mockHttp.post.mockRejectedValueOnce(new Error('Failed to create'));

            await expect(
                useDisputeStore.getState().createDispute({
                    orderId: 100,
                    type: DisputeType.SERVICE_QUALITY,
                    reason: 'Service issue',
                })
            ).rejects.toThrow('Failed to create');

            const state = useDisputeStore.getState();
            expect(state.error).toBe('Failed to create');
            expect(state.loading).toBe(false);
        });

        it('should add new dispute to beginning of list', async () => {
            const existingDispute = createMockDispute({ id: 1 });
            useDisputeStore.setState({ myDisputes: [existingDispute] });

            const newDispute = createMockDispute({ id: 2 });
            mockHttp.post.mockResolvedValueOnce(newDispute);

            await useDisputeStore.getState().createDispute({
                orderId: 100,
                type: DisputeType.SERVICE_QUALITY,
                reason: 'New issue',
            });

            const state = useDisputeStore.getState();
            expect(state.myDisputes[0].id).toBe(2);
            expect(state.myDisputes[1].id).toBe(1);
        });
    });

    describe('cancelDispute', () => {
        it('should cancel dispute successfully', async () => {
            const dispute = createMockDispute({ id: 1, status: DisputeStatus.PENDING });
            useDisputeStore.setState({
                myDisputes: [dispute],
                currentDispute: dispute,
            });
            mockHttp.post.mockResolvedValueOnce({});

            await useDisputeStore.getState().cancelDispute(1);

            const state = useDisputeStore.getState();
            expect(state.myDisputes[0].status).toBe(DisputeStatus.CANCELED);
            expect(state.currentDispute?.status).toBe(DisputeStatus.CANCELED);
            expect(state.loading).toBe(false);
            expect(mockHttp.post).toHaveBeenCalledWith('/user/disputes/1/cancel');
        });

        it('should handle cancel dispute error', async () => {
            const dispute = createMockDispute({ id: 1 });
            useDisputeStore.setState({ myDisputes: [dispute] });
            mockHttp.post.mockRejectedValueOnce(new Error('Cannot cancel'));

            await expect(useDisputeStore.getState().cancelDispute(1)).rejects.toThrow('Cannot cancel');

            const state = useDisputeStore.getState();
            expect(state.error).toBe('Cannot cancel');
        });

        it('should not update currentDispute if different id', async () => {
            const dispute1 = createMockDispute({ id: 1 });
            const dispute2 = createMockDispute({ id: 2 });
            useDisputeStore.setState({
                myDisputes: [dispute1, dispute2],
                currentDispute: dispute2,
            });
            mockHttp.post.mockResolvedValueOnce({});

            await useDisputeStore.getState().cancelDispute(1);

            const state = useDisputeStore.getState();
            expect(state.myDisputes[0].status).toBe(DisputeStatus.CANCELED);
            expect(state.currentDispute?.status).toBe(DisputeStatus.PENDING);
        });
    });

    describe('addEvidence', () => {
        it('should add evidence successfully', async () => {
            const dispute = createMockDispute({ id: 1, evidenceUrls: [] });
            const updatedDispute = createMockDispute({
                id: 1,
                evidenceUrls: ['new-evidence.jpg'],
                evidenceText: 'New evidence text',
            });
            useDisputeStore.setState({
                myDisputes: [dispute],
                currentDispute: dispute,
            });
            mockHttp.post.mockResolvedValueOnce(updatedDispute);

            await useDisputeStore.getState().addEvidence(1, ['new-evidence.jpg'], 'New evidence text');

            const state = useDisputeStore.getState();
            expect(state.currentDispute?.evidenceUrls).toContain('new-evidence.jpg');
            expect(state.myDisputes[0].evidenceUrls).toContain('new-evidence.jpg');
            expect(mockHttp.post).toHaveBeenCalledWith('/user/disputes/1/evidence', {
                evidenceUrls: ['new-evidence.jpg'],
                evidenceText: 'New evidence text',
            });
        });

        it('should handle add evidence error', async () => {
            const dispute = createMockDispute({ id: 1 });
            useDisputeStore.setState({ currentDispute: dispute, myDisputes: [dispute] });
            mockHttp.post.mockRejectedValueOnce(new Error('Upload failed'));

            await expect(
                useDisputeStore.getState().addEvidence(1, ['file.jpg'])
            ).rejects.toThrow('Upload failed');

            const state = useDisputeStore.getState();
            expect(state.error).toBe('Upload failed');
        });
    });

    describe('fetchMessages', () => {
        it('should fetch messages successfully', async () => {
            const mockMessages = [
                createMockMessage({ id: 1, content: 'First message' }),
                createMockMessage({ id: 2, content: 'Second message' }),
            ];
            mockHttp.get.mockResolvedValueOnce(mockMessages);

            await useDisputeStore.getState().fetchMessages(1);

            const state = useDisputeStore.getState();
            expect(state.messages).toHaveLength(2);
            expect(state.messages[0].content).toBe('First message');
            expect(mockHttp.get).toHaveBeenCalledWith('/user/disputes/1/messages');
        });

        it('should handle fetch messages error silently', async () => {
            const consoleSpy = vi.spyOn(console, 'error').mockImplementation(() => {});
            mockHttp.get.mockRejectedValueOnce(new Error('Failed'));

            await useDisputeStore.getState().fetchMessages(1);

            const state = useDisputeStore.getState();
            expect(state.messages).toEqual([]);
            expect(consoleSpy).toHaveBeenCalled();
            consoleSpy.mockRestore();
        });

        it('should handle null response', async () => {
            mockHttp.get.mockResolvedValueOnce(null);

            await useDisputeStore.getState().fetchMessages(1);

            const state = useDisputeStore.getState();
            expect(state.messages).toEqual([]);
        });
    });

    describe('sendMessage', () => {
        it('should send message successfully', async () => {
            const newMessage = createMockMessage({ id: 3, content: 'New message' });
            mockHttp.post.mockResolvedValueOnce(newMessage);
            useDisputeStore.setState({ messages: [] });

            await useDisputeStore.getState().sendMessage(1, 'New message');

            const state = useDisputeStore.getState();
            expect(state.messages).toHaveLength(1);
            expect(state.messages[0].content).toBe('New message');
            expect(mockHttp.post).toHaveBeenCalledWith('/user/disputes/1/messages', {
                content: 'New message',
            });
        });

        it('should append message to existing messages', async () => {
            const existingMessage = createMockMessage({ id: 1 });
            const newMessage = createMockMessage({ id: 2, content: 'New' });
            useDisputeStore.setState({ messages: [existingMessage] });
            mockHttp.post.mockResolvedValueOnce(newMessage);

            await useDisputeStore.getState().sendMessage(1, 'New');

            const state = useDisputeStore.getState();
            expect(state.messages).toHaveLength(2);
        });

        it('should handle send message error', async () => {
            const consoleSpy = vi.spyOn(console, 'error').mockImplementation(() => {});
            mockHttp.post.mockRejectedValueOnce(new Error('Send failed'));

            await expect(useDisputeStore.getState().sendMessage(1, 'Test')).rejects.toThrow('Send failed');
            expect(consoleSpy).toHaveBeenCalled();
            consoleSpy.mockRestore();
        });
    });

    describe('canCreateDispute', () => {
        it('should return true when dispute is allowed', async () => {
            mockHttp.get.mockResolvedValueOnce({ canDispute: true });

            const result = await useDisputeStore.getState().canCreateDispute(100);

            expect(result).toBe(true);
            expect(mockHttp.get).toHaveBeenCalledWith('/orders/100/can-dispute');
        });

        it('should return false when dispute is not allowed', async () => {
            mockHttp.get.mockResolvedValueOnce({ canDispute: false, reason: 'Order too old' });

            const result = await useDisputeStore.getState().canCreateDispute(100);

            expect(result).toBe(false);
        });

        it('should return false on error', async () => {
            mockHttp.get.mockRejectedValueOnce(new Error('Network error'));

            const result = await useDisputeStore.getState().canCreateDispute(100);

            expect(result).toBe(false);
        });
    });

    describe('Label Utilities', () => {
        it('should return correct status labels', () => {
            const { getStatusLabel } = useDisputeStore.getState();

            expect(getStatusLabel(DisputeStatus.PENDING)).toBe('待处理');
            expect(getStatusLabel(DisputeStatus.ASSIGNED)).toBe('已指派客服');
            expect(getStatusLabel(DisputeStatus.MEDIATING)).toBe('调解中');
            expect(getStatusLabel(DisputeStatus.RESOLVED)).toBe('已解决');
            expect(getStatusLabel(DisputeStatus.REJECTED)).toBe('已驳回');
            expect(getStatusLabel(DisputeStatus.CANCELED)).toBe('已取消');
        });

        it('should return status as fallback for unknown status', () => {
            const { getStatusLabel } = useDisputeStore.getState();
            expect(getStatusLabel('unknown' as DisputeStatus)).toBe('unknown');
        });

        it('should return correct type labels', () => {
            const { getTypeLabel } = useDisputeStore.getState();

            expect(getTypeLabel(DisputeType.SERVICE_QUALITY)).toBe('服务质量问题');
            expect(getTypeLabel(DisputeType.BAD_ATTITUDE)).toBe('态度问题');
            expect(getTypeLabel(DisputeType.INCOMPLETE_SERVICE)).toBe('未完成服务');
            expect(getTypeLabel(DisputeType.USER_NOT_COOPERATIVE)).toBe('用户不配合');
            expect(getTypeLabel(DisputeType.USER_HARASSMENT)).toBe('用户骚扰');
            expect(getTypeLabel(DisputeType.OTHER)).toBe('其他');
        });

        it('should return correct resolution labels', () => {
            const { getResolutionLabel } = useDisputeStore.getState();

            expect(getResolutionLabel(DisputeResolution.REFUND)).toBe('全额退款');
            expect(getResolutionLabel(DisputeResolution.PARTIAL)).toBe('部分退款');
            expect(getResolutionLabel(DisputeResolution.REASSIGN)).toBe('重新指派');
            expect(getResolutionLabel(DisputeResolution.REJECT)).toBe('驳回申诉');
            expect(getResolutionLabel(DisputeResolution.PENDING)).toBe('待处理');
        });
    });

    describe('reset', () => {
        it('should reset store to initial state', () => {
            useDisputeStore.setState({
                myDisputes: [createMockDispute()],
                currentDispute: createMockDispute(),
                templates: [createMockTemplate()],
                messages: [createMockMessage()],
                pagination: { page: 5, pageSize: 20, total: 100, hasMore: true },
                loading: true,
                error: 'Some error',
            });

            useDisputeStore.getState().reset();

            const state = useDisputeStore.getState();
            expect(state.myDisputes).toEqual([]);
            expect(state.currentDispute).toBeNull();
            expect(state.templates).toEqual([]);
            expect(state.messages).toEqual([]);
            expect(state.pagination).toEqual({
                page: 1,
                pageSize: 20,
                total: 0,
                hasMore: false,
            });
            expect(state.loading).toBe(false);
            expect(state.error).toBeNull();
        });
    });

    describe('Selectors', () => {
        describe('selectActiveDisputes', () => {
            it('should return only active disputes', () => {
                const disputes = [
                    createMockDispute({ id: 1, status: DisputeStatus.PENDING }),
                    createMockDispute({ id: 2, status: DisputeStatus.ASSIGNED }),
                    createMockDispute({ id: 3, status: DisputeStatus.MEDIATING }),
                    createMockDispute({ id: 4, status: DisputeStatus.RESOLVED }),
                    createMockDispute({ id: 5, status: DisputeStatus.REJECTED }),
                    createMockDispute({ id: 6, status: DisputeStatus.CANCELED }),
                ];

                const state = { myDisputes: disputes } as Pick<DisputeState, 'myDisputes'>;
                const active = selectActiveDisputes(state);

                expect(active).toHaveLength(3);
                expect(active.map(d => d.id)).toEqual([1, 2, 3]);
            });

            it('should return empty array when no active disputes', () => {
                const disputes = [
                    createMockDispute({ id: 1, status: DisputeStatus.RESOLVED }),
                    createMockDispute({ id: 2, status: DisputeStatus.CANCELED }),
                ];

                const state = { myDisputes: disputes } as Pick<DisputeState, 'myDisputes'>;
                const active = selectActiveDisputes(state);

                expect(active).toHaveLength(0);
            });
        });

        describe('selectResolvedDisputes', () => {
            it('should return only resolved/rejected/canceled disputes', () => {
                const disputes = [
                    createMockDispute({ id: 1, status: DisputeStatus.PENDING }),
                    createMockDispute({ id: 2, status: DisputeStatus.RESOLVED }),
                    createMockDispute({ id: 3, status: DisputeStatus.REJECTED }),
                    createMockDispute({ id: 4, status: DisputeStatus.CANCELED }),
                ];

                const state = { myDisputes: disputes } as Pick<DisputeState, 'myDisputes'>;
                const resolved = selectResolvedDisputes(state);

                expect(resolved).toHaveLength(3);
                expect(resolved.map(d => d.id)).toEqual([2, 3, 4]);
            });
        });

        describe('selectIsDisputable', () => {
            it('should return true for pending status', () => {
                expect(selectIsDisputable(DisputeStatus.PENDING)).toBe(true);
            });

            it('should return true for assigned status', () => {
                expect(selectIsDisputable(DisputeStatus.ASSIGNED)).toBe(true);
            });

            it('should return false for other statuses', () => {
                expect(selectIsDisputable(DisputeStatus.MEDIATING)).toBe(false);
                expect(selectIsDisputable(DisputeStatus.RESOLVED)).toBe(false);
                expect(selectIsDisputable(DisputeStatus.REJECTED)).toBe(false);
                expect(selectIsDisputable(DisputeStatus.CANCELED)).toBe(false);
            });
        });
    });
});
