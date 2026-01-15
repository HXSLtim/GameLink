import { create } from 'zustand';
import { http } from '@/lib/http';

// ============ Enums ============

export const ActivityType = {
    DISCOUNT: 'discount',       // 折扣活动
    RECHARGE: 'recharge',       // 充值活动
    NEW_USER: 'new_user',       // 新用户活动
    FESTIVAL: 'festival',       // 节日活动
    FLASH: 'flash'              // 限时秒杀
} as const;

export type ActivityType = typeof ActivityType[keyof typeof ActivityType];

export const ActivityStatus = {
    DRAFT: 'draft',             // 草稿
    SCHEDULED: 'scheduled',     // 已排期
    ACTIVE: 'active',           // 进行中
    ENDED: 'ended',             // 已结束
    CANCELED: 'canceled'        // 已取消
} as const;

export type ActivityStatus = typeof ActivityStatus[keyof typeof ActivityStatus];

// ============ Interfaces ============

export interface ActivityRule {
    type: 'min_order' | 'first_order' | 'specific_game' | 'specific_player';
    value: any;
    description: string;
}

export interface ActivityReward {
    type: 'coupon' | 'points' | 'discount' | 'gift';
    value: any;
    description: string;
}

export interface Activity {
    id: number;
    name: string;
    description: string;
    type: ActivityType;
    status: ActivityStatus;

    // 展示
    bannerUrl: string;
    detailUrl?: string;

    // 时间
    startAt: string;
    endAt: string;

    // 规则
    rules: ActivityRule[];

    // 参与限制
    maxParticipants?: number;
    currentParticipants: number;
    userLimit: number;           // 每人参与次数限制

    // 奖励
    rewards: ActivityReward[];
}

export interface ActivityParticipation {
    activityId: number;
    userId: number;
    participatedAt: string;
    participationCount: number;
    rewardsReceived: ActivityReward[];
}

// ============ State & Actions ============

export interface ActivityState {
    activities: Activity[];
    currentActivity: Activity | null;
    myParticipations: ActivityParticipation[];
    loading: boolean;
    error: string | null;
}

export interface ActivityActions {
    fetchActivities: (type?: ActivityType) => Promise<void>;
    fetchActivityDetail: (id: number) => Promise<void>;
    fetchMyParticipation: (activityId: number) => Promise<ActivityParticipation | null>;
    joinActivity: (id: number) => Promise<void>;
    getActiveActivities: () => Activity[];
    canParticipate: (activity: Activity, participation?: ActivityParticipation) => { can: boolean; reason?: string };
}

// ============ Store ============

export const useActivityStore = create<ActivityState & ActivityActions>((set, get) => ({
    activities: [],
    currentActivity: null,
    myParticipations: [],
    loading: false,
    error: null,

    fetchActivities: async (type) => {
        set({ loading: true, error: null });
        try {
            const params: Record<string, any> = { status: 'active' };
            if (type) params.type = type;

            const data = await http.get<Activity[]>('/activities', { params });
            set({ activities: data, loading: false });
        } catch (err: any) {
            set({ loading: false, error: err.message || 'Failed to fetch activities' });
        }
    },

    fetchActivityDetail: async (id) => {
        set({ loading: true, error: null });
        try {
            const data = await http.get<Activity>(`/activities/${id}`);
            set({ currentActivity: data, loading: false });
        } catch (err: any) {
            set({ loading: false, error: err.message || 'Failed to fetch activity detail' });
        }
    },

    fetchMyParticipation: async (activityId) => {
        try {
            const data = await http.get<ActivityParticipation>(`/activities/${activityId}/participation`);

            // 更新本地参与记录
            set((state) => ({
                myParticipations: [
                    ...state.myParticipations.filter(p => p.activityId !== activityId),
                    data
                ]
            }));

            return data;
        } catch (err: any) {
            // 404 表示未参与，不是错误
            if (err.status === 404) return null;
            console.error('Failed to fetch participation:', err);
            return null;
        }
    },

    joinActivity: async (id) => {
        set({ loading: true, error: null });
        try {
            await http.post(`/activities/${id}/join`);

            // 刷新活动详情和参与状态
            await Promise.all([
                get().fetchActivityDetail(id),
                get().fetchMyParticipation(id)
            ]);

            set({ loading: false });
        } catch (err: any) {
            set({ loading: false, error: err.message || 'Failed to join activity' });
            throw err;
        }
    },

    getActiveActivities: () => {
        const { activities } = get();
        const now = new Date();
        return activities.filter(activity => {
            if (activity.status !== ActivityStatus.ACTIVE) return false;
            const startAt = new Date(activity.startAt);
            const endAt = new Date(activity.endAt);
            return now >= startAt && now <= endAt;
        });
    },

    canParticipate: (activity, participation) => {
        const now = new Date();
        const startAt = new Date(activity.startAt);
        const endAt = new Date(activity.endAt);

        // 检查活动状态
        if (activity.status !== ActivityStatus.ACTIVE) {
            return { can: false, reason: '活动未开始或已结束' };
        }

        // 检查时间
        if (now < startAt) {
            return { can: false, reason: '活动尚未开始' };
        }
        if (now > endAt) {
            return { can: false, reason: '活动已结束' };
        }

        // 检查参与人数上限
        if (activity.maxParticipants && activity.currentParticipants >= activity.maxParticipants) {
            return { can: false, reason: '参与人数已满' };
        }

        // 检查个人参与次数
        if (participation && participation.participationCount >= activity.userLimit) {
            return { can: false, reason: '已达到参与次数上限' };
        }

        return { can: true };
    }
}));
