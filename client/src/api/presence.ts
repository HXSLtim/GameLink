/**
 * Presence API
 * Handles user online status, presence updates
 */

import { http } from '@/lib/http';
import type {
    UserPresence
} from '@/types/api';

export const presenceApi = {
    /**
     * Get user presence
     */
    getPresence: (userId: number) =>
        http.get<UserPresence>(`/presence/${userId}`),

    /**
     * Get multiple users' presence
     */
    getBatchPresence: (userIds: number[]) =>
        http.post<Record<number, UserPresence>>('/presence/batch', { userIds }),

    /**
     * Update own presence
     */
    updatePresence: (data: {
        status: 'online' | 'away' | 'busy' | 'offline';
        customStatus?: string;
    }) =>
        http.put<void>('/presence/update', data),

    /**
     * Get online friends
     */
    getOnlineFriends: () =>
        http.get<Array<{
            userId: number;
            username: string;
            avatar: string;
            status: string;
            lastSeen: string;
        }>>('/presence/online-friends'),
};
