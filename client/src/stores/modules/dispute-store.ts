import { create } from 'zustand';
import { http } from '@/lib/http';

// ============ Enums ============

export const DisputeStatus = {
    PENDING: 'pending',       // 待处理
    ASSIGNED: 'assigned',     // 已指派
    MEDIATING: 'mediating',   // 调解中
    RESOLVED: 'resolved',     // 已解决
    REJECTED: 'rejected',     // 已驳回
    CANCELED: 'canceled'      // 已取消
} as const;

export type DisputeStatus = typeof DisputeStatus[keyof typeof DisputeStatus];

export const DisputeResolution = {
    REFUND: 'refund',         // 全额退款
    PARTIAL: 'partial',       // 部分退款
    REASSIGN: 'reassign',     // 重新指派
    REJECT: 'reject',         // 驳回
    PENDING: 'pending'        // 待决定
} as const;

export type DisputeResolution = typeof DisputeResolution[keyof typeof DisputeResolution];

export const DisputeInitiatorType = {
    USER: 'user',             // 用户发起
    PLAYER: 'player'          // 陪玩师发起
} as const;

export type DisputeInitiatorType = typeof DisputeInitiatorType[keyof typeof DisputeInitiatorType];

export const DisputeType = {
    SERVICE_QUALITY: 'service_quality',           // 服务质量问题
    BAD_ATTITUDE: 'bad_attitude',                 // 态度问题
    INCOMPLETE_SERVICE: 'incomplete_service',     // 未完成服务
    USER_NOT_COOPERATIVE: 'user_not_cooperative', // 用户不配合
    USER_HARASSMENT: 'user_harassment',           // 用户骚扰
    OTHER: 'other'                                // 其他
} as const;

export type DisputeType = typeof DisputeType[keyof typeof DisputeType];

// ============ Interfaces ============

export interface OrderDispute {
    id: number;
    orderId: number;
    orderNo?: string;
    initiatorId: number;
    initiatorType: DisputeInitiatorType;
    type: DisputeType;
    status: DisputeStatus;
    reason: string;
    evidenceUrls: string[];
    evidenceText?: string;

    // SLA info
    slaDeadline?: string;
    slaBreached: boolean;
    slaRemaining?: number;  // seconds

    // Resolution info
    resolution: DisputeResolution;
    resolvedBy?: number;
    resolvedAt?: string;
    resolveRemark?: string;

    // Assigned CS
    assignedServiceId?: number;
    assignedServiceName?: string;

    createdAt: string;
    updatedAt: string;
}

export interface DisputeTemplate {
    id: number;
    code: string;
    name: string;
    initiatorType: DisputeInitiatorType;
    description: string;
    sortOrder: number;
    isActive: boolean;
}

export interface CreateDisputeRequest {
    orderId: number;
    type: DisputeType;
    reason: string;
    evidenceUrls?: string[];
    evidenceText?: string;
}

export interface DisputeMessage {
    id: number;
    disputeId: number;
    senderId: number;
    senderName: string;
    senderRole: 'user' | 'player' | 'cs';
    content: string;
    createdAt: string;
}

// ============ State & Actions ============

export interface DisputeState {
    // My disputes (as user or player)
    myDisputes: OrderDispute[];
    currentDispute: OrderDispute | null;

    // Templates for creating disputes
    templates: DisputeTemplate[];

    // Messages for current dispute
    messages: DisputeMessage[];

    // Pagination
    pagination: {
        page: number;
        pageSize: number;
        total: number;
        hasMore: boolean;
    };

    // Status
    loading: boolean;
    error: string | null;
}

export interface DisputeActions {
    // Fetch disputes
    fetchMyDisputes: (page?: number) => Promise<void>;
    fetchDisputeById: (id: number) => Promise<void>;
    fetchTemplates: (initiatorType?: DisputeInitiatorType) => Promise<void>;

    // Create & manage disputes
    createDispute: (data: CreateDisputeRequest) => Promise<OrderDispute>;
    cancelDispute: (id: number) => Promise<void>;
    addEvidence: (id: number, evidenceUrls: string[], evidenceText?: string) => Promise<void>;

    // Messages
    fetchMessages: (disputeId: number) => Promise<void>;
    sendMessage: (disputeId: number, content: string) => Promise<void>;

    // Check eligibility
    canCreateDispute: (orderId: number) => Promise<boolean>;

    // Utilities
    getStatusLabel: (status: DisputeStatus) => string;
    getTypeLabel: (type: DisputeType) => string;
    getResolutionLabel: (resolution: DisputeResolution) => string;

    // Reset
    reset: () => void;
}

// ============ Label Maps ============

const STATUS_LABELS: Record<DisputeStatus, string> = {
    [DisputeStatus.PENDING]: '待处理',
    [DisputeStatus.ASSIGNED]: '已指派客服',
    [DisputeStatus.MEDIATING]: '调解中',
    [DisputeStatus.RESOLVED]: '已解决',
    [DisputeStatus.REJECTED]: '已驳回',
    [DisputeStatus.CANCELED]: '已取消'
};

const TYPE_LABELS: Record<DisputeType, string> = {
    [DisputeType.SERVICE_QUALITY]: '服务质量问题',
    [DisputeType.BAD_ATTITUDE]: '态度问题',
    [DisputeType.INCOMPLETE_SERVICE]: '未完成服务',
    [DisputeType.USER_NOT_COOPERATIVE]: '用户不配合',
    [DisputeType.USER_HARASSMENT]: '用户骚扰',
    [DisputeType.OTHER]: '其他'
};

const RESOLUTION_LABELS: Record<DisputeResolution, string> = {
    [DisputeResolution.REFUND]: '全额退款',
    [DisputeResolution.PARTIAL]: '部分退款',
    [DisputeResolution.REASSIGN]: '重新指派',
    [DisputeResolution.REJECT]: '驳回申诉',
    [DisputeResolution.PENDING]: '待处理'
};

// ============ Initial State ============

const INITIAL_STATE: DisputeState = {
    myDisputes: [],
    currentDispute: null,
    templates: [],
    messages: [],
    pagination: {
        page: 1,
        pageSize: 20,
        total: 0,
        hasMore: false
    },
    loading: false,
    error: null
};

// ============ Store ============

export const useDisputeStore = create<DisputeState & DisputeActions>((set, get) => ({
    ...INITIAL_STATE,

    fetchMyDisputes: async (page = 1) => {
        set({ loading: true, error: null });
        try {
            const { pagination } = get();
            const data = await http.get<{
                items: OrderDispute[];
                total: number;
            }>('/user/disputes', {
                params: { page, pageSize: pagination.pageSize }
            });

            const items = data.items || [];
            const total = data.total || 0;

            set({
                myDisputes: page === 1 ? items : [...get().myDisputes, ...items],
                pagination: {
                    ...pagination,
                    page,
                    total,
                    hasMore: get().myDisputes.length + items.length < total
                },
                loading: false
            });
        } catch (err) {
            set({
                loading: false,
                error: err instanceof Error ? err.message : 'Failed to fetch disputes'
            });
        }
    },

    fetchDisputeById: async (id) => {
        set({ loading: true, error: null });
        try {
            const data = await http.get<OrderDispute>(`/user/disputes/${id}`);
            set({ currentDispute: data, loading: false });
        } catch (err) {
            set({
                loading: false,
                error: err instanceof Error ? err.message : 'Failed to fetch dispute'
            });
        }
    },

    fetchTemplates: async (initiatorType) => {
        try {
            const params = initiatorType ? { initiatorType } : {};
            const data = await http.get<DisputeTemplate[]>('/disputes/templates', { params });
            set({ templates: data || [] });
        } catch (err) {
            console.error('Failed to fetch dispute templates:', err);
        }
    },

    createDispute: async (data) => {
        set({ loading: true, error: null });
        try {
            const dispute = await http.post<OrderDispute>('/user/disputes', data);

            // Add to list
            set((state) => ({
                myDisputes: [dispute, ...state.myDisputes],
                currentDispute: dispute,
                loading: false
            }));

            return dispute;
        } catch (err) {
            set({
                loading: false,
                error: err instanceof Error ? err.message : 'Failed to create dispute'
            });
            throw err;
        }
    },

    cancelDispute: async (id) => {
        set({ loading: true, error: null });
        try {
            await http.post(`/user/disputes/${id}/cancel`);

            // Update in list
            set((state) => ({
                myDisputes: state.myDisputes.map(d =>
                    d.id === id ? { ...d, status: DisputeStatus.CANCELED } : d
                ),
                currentDispute: state.currentDispute?.id === id
                    ? { ...state.currentDispute, status: DisputeStatus.CANCELED }
                    : state.currentDispute,
                loading: false
            }));
        } catch (err) {
            set({
                loading: false,
                error: err instanceof Error ? err.message : 'Failed to cancel dispute'
            });
            throw err;
        }
    },

    addEvidence: async (id, evidenceUrls, evidenceText) => {
        set({ loading: true, error: null });
        try {
            const updated = await http.post<OrderDispute>(
                `/user/disputes/${id}/evidence`,
                { evidenceUrls, evidenceText }
            );

            set((state) => ({
                currentDispute: state.currentDispute?.id === id ? updated : state.currentDispute,
                myDisputes: state.myDisputes.map(d => d.id === id ? updated : d),
                loading: false
            }));
        } catch (err) {
            set({
                loading: false,
                error: err instanceof Error ? err.message : 'Failed to add evidence'
            });
            throw err;
        }
    },

    fetchMessages: async (disputeId) => {
        try {
            const data = await http.get<DisputeMessage[]>(
                `/user/disputes/${disputeId}/messages`
            );
            set({ messages: data || [] });
        } catch (err) {
            console.error('Failed to fetch dispute messages:', err);
        }
    },

    sendMessage: async (disputeId, content) => {
        try {
            const message = await http.post<DisputeMessage>(
                `/user/disputes/${disputeId}/messages`,
                { content }
            );

            set((state) => ({
                messages: [...state.messages, message]
            }));
        } catch (err) {
            console.error('Failed to send message:', err);
            throw err;
        }
    },

    canCreateDispute: async (orderId) => {
        try {
            const data = await http.get<{ canDispute: boolean; reason?: string }>(
                `/orders/${orderId}/can-dispute`
            );
            return data.canDispute;
        } catch {
            return false;
        }
    },

    getStatusLabel: (status) => STATUS_LABELS[status] || status,
    getTypeLabel: (type) => TYPE_LABELS[type] || type,
    getResolutionLabel: (resolution) => RESOLUTION_LABELS[resolution] || resolution,

    reset: () => set(INITIAL_STATE)
}));

// ============ Selectors ============

export const selectActiveDisputes = (state: DisputeState) =>
    state.myDisputes.filter(d =>
        d.status === DisputeStatus.PENDING ||
        d.status === DisputeStatus.ASSIGNED ||
        d.status === DisputeStatus.MEDIATING
    );

export const selectResolvedDisputes = (state: DisputeState) =>
    state.myDisputes.filter(d =>
        d.status === DisputeStatus.RESOLVED ||
        d.status === DisputeStatus.REJECTED ||
        d.status === DisputeStatus.CANCELED
    );

export const selectIsDisputable = (status: DisputeStatus) =>
    status === DisputeStatus.PENDING || status === DisputeStatus.ASSIGNED;
