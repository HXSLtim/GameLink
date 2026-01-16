import { create } from 'zustand';
import { http } from '@/lib/http';
import { getErrorMessage, logError } from '@/lib/error';

// ============ Enums ============

export const TeamRole = {
    LEADER: 'leader',
    MEMBER: 'member'
} as const;

export type TeamRole = typeof TeamRole[keyof typeof TeamRole];

export const InviteStatus = {
    PENDING: 'pending',
    ACCEPTED: 'accepted',
    REJECTED: 'rejected',
    EXPIRED: 'expired'
} as const;

export type InviteStatus = typeof InviteStatus[keyof typeof InviteStatus];

// ============ Interfaces ============

export interface TeamMember {
    id: number;
    playerId: number;
    teamId: number;
    role: TeamRole;
    nickname: string;
    avatar: string;
    rating: number;
    joinedAt: string;
}

export interface Team {
    id: number;
    name: string;
    description?: string;
    logo?: string;
    leaderId: number;
    leaderName: string;
    memberCount: number;
    maxMembers: number;
    totalOrders: number;
    totalEarningsCents: number;
    createdAt: string;
    updatedAt: string;
}

export interface TeamInvite {
    id: number;
    teamId: number;
    teamName: string;
    teamLogo?: string;
    inviterId: number;
    inviterName: string;
    inviteeId: number;
    status: InviteStatus;
    message?: string;
    expiresAt: string;
    createdAt: string;
}

export interface CreateTeamRequest {
    name: string;
    description?: string;
    logo?: string;
}

export interface UpdateTeamRequest {
    name?: string;
    description?: string;
    logo?: string;
}

export interface InviteMemberRequest {
    playerId: number;
    message?: string;
}

// ============ State & Actions ============

export interface TeamState {
    myTeam: Team | null;
    members: TeamMember[];
    invites: TeamInvite[];
    receivedInvites: TeamInvite[];
    loading: boolean;
    error: string | null;
}

export interface TeamActions {
    // Team CRUD
    fetchMyTeam: () => Promise<void>;
    createTeam: (request: CreateTeamRequest) => Promise<void>;
    updateTeam: (request: UpdateTeamRequest) => Promise<void>;
    deleteTeam: () => Promise<void>;

    // Members
    fetchMembers: () => Promise<void>;
    kickMember: (memberId: number) => Promise<void>;
    transferLeadership: (memberId: number) => Promise<void>;
    leaveTeam: () => Promise<void>;

    // Invites
    fetchInvites: () => Promise<void>;
    inviteMember: (request: InviteMemberRequest) => Promise<void>;
    cancelInvite: (inviteId: number) => Promise<void>;

    // Received invites
    fetchReceivedInvites: () => Promise<void>;
    acceptInvite: (inviteId: number) => Promise<void>;
    rejectInvite: (inviteId: number) => Promise<void>;

    // Helpers
    isLeader: () => boolean;
    canInvite: () => boolean;
}

// ============ Store ============

export const useTeamStore = create<TeamState & TeamActions>((set, get) => ({
    myTeam: null,
    members: [],
    invites: [],
    receivedInvites: [],
    loading: false,
    error: null,

    // ========== Team CRUD ==========

    fetchMyTeam: async () => {
        set({ loading: true, error: null });
        try {
            const data = await http.get<Team>('/player/team');
            set({ myTeam: data, loading: false });
        } catch (err) {
            // 404 means no team, not an error
            if (err && typeof err === 'object' && 'status' in err && err.status === 404) {
                set({ myTeam: null, loading: false });
                return;
            }
            logError('fetchMyTeam', err);
            set({ loading: false, error: getErrorMessage(err, 'Failed to fetch team') });
        }
    },

    createTeam: async (request) => {
        set({ loading: true, error: null });
        try {
            const data = await http.post<Team>('/player/team', request);
            set({ myTeam: data, loading: false });
        } catch (err) {
            logError('createTeam', err);
            set({ loading: false, error: getErrorMessage(err, 'Failed to create team') });
            throw err;
        }
    },

    updateTeam: async (request) => {
        set({ loading: true, error: null });
        try {
            const data = await http.put<Team>('/player/team', request);
            set({ myTeam: data, loading: false });
        } catch (err) {
            logError('updateTeam', err);
            set({ loading: false, error: getErrorMessage(err, 'Failed to update team') });
            throw err;
        }
    },

    deleteTeam: async () => {
        set({ loading: true, error: null });
        try {
            await http.delete('/player/team');
            set({ myTeam: null, members: [], invites: [], loading: false });
        } catch (err) {
            logError('deleteTeam', err);
            set({ loading: false, error: getErrorMessage(err, 'Failed to delete team') });
            throw err;
        }
    },

    // ========== Members ==========

    fetchMembers: async () => {
        set({ loading: true, error: null });
        try {
            const data = await http.get<{ items: TeamMember[] }>('/player/team/members');
            set({ members: data.items || [], loading: false });
        } catch (err) {
            logError('fetchMembers', err);
            set({ loading: false, error: getErrorMessage(err, 'Failed to fetch members') });
        }
    },

    kickMember: async (memberId) => {
        set({ loading: true, error: null });
        try {
            await http.post(`/player/team/members/${memberId}/kick`);
            // Refresh members list
            await get().fetchMembers();
            // Update team member count
            set((state) => ({
                myTeam: state.myTeam
                    ? { ...state.myTeam, memberCount: state.myTeam.memberCount - 1 }
                    : null,
                loading: false
            }));
        } catch (err) {
            logError('kickMember', err);
            set({ loading: false, error: getErrorMessage(err, 'Failed to kick member') });
            throw err;
        }
    },

    transferLeadership: async (memberId) => {
        set({ loading: true, error: null });
        try {
            await http.post('/player/team/transfer', { newLeaderId: memberId });
            // Refresh team and members
            await Promise.all([get().fetchMyTeam(), get().fetchMembers()]);
        } catch (err) {
            logError('transferLeadership', err);
            set({ loading: false, error: getErrorMessage(err, 'Failed to transfer leadership') });
            throw err;
        }
    },

    leaveTeam: async () => {
        set({ loading: true, error: null });
        try {
            await http.post('/player/team/leave');
            set({ myTeam: null, members: [], invites: [], loading: false });
        } catch (err) {
            logError('leaveTeam', err);
            set({ loading: false, error: getErrorMessage(err, 'Failed to leave team') });
            throw err;
        }
    },

    // ========== Invites ==========

    fetchInvites: async () => {
        try {
            const data = await http.get<{ items: TeamInvite[] }>('/player/team/invites');
            set({ invites: data.items || [] });
        } catch (err) {
            logError('fetchInvites', err);
        }
    },

    inviteMember: async (request) => {
        set({ loading: true, error: null });
        try {
            await http.post('/player/team/invites', request);
            await get().fetchInvites();
            set({ loading: false });
        } catch (err) {
            logError('inviteMember', err);
            set({ loading: false, error: getErrorMessage(err, 'Failed to invite member') });
            throw err;
        }
    },

    cancelInvite: async (inviteId) => {
        try {
            await http.delete(`/player/team/invites/${inviteId}`);
            set((state) => ({
                invites: state.invites.filter(i => i.id !== inviteId)
            }));
        } catch (err) {
            logError('cancelInvite', err);
            throw err;
        }
    },

    // ========== Received Invites ==========

    fetchReceivedInvites: async () => {
        try {
            const data = await http.get<{ items: TeamInvite[] }>('/player/team/invites/received');
            set({ receivedInvites: data.items || [] });
        } catch (err) {
            logError('fetchReceivedInvites', err);
        }
    },

    acceptInvite: async (inviteId) => {
        set({ loading: true, error: null });
        try {
            await http.post(`/player/team/invites/${inviteId}/accept`);
            // Refresh team data
            await get().fetchMyTeam();
            set((state) => ({
                receivedInvites: state.receivedInvites.filter(i => i.id !== inviteId),
                loading: false
            }));
        } catch (err) {
            logError('acceptInvite', err);
            set({ loading: false, error: getErrorMessage(err, 'Failed to accept invite') });
            throw err;
        }
    },

    rejectInvite: async (inviteId) => {
        try {
            await http.post(`/player/team/invites/${inviteId}/reject`);
            set((state) => ({
                receivedInvites: state.receivedInvites.filter(i => i.id !== inviteId)
            }));
        } catch (err) {
            logError('rejectInvite', err);
            throw err;
        }
    },

    // ========== Helpers ==========

    isLeader: () => {
        const { myTeam, members } = get();
        if (!myTeam) return false;
        // Find current user in members and check role
        const currentMember = members.find(m => m.playerId === myTeam.leaderId);
        return currentMember?.role === TeamRole.LEADER;
    },

    canInvite: () => {
        const { myTeam } = get();
        if (!myTeam) return false;
        return myTeam.memberCount < myTeam.maxMembers;
    }
}));
