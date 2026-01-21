/**
 * User API
 * Handles user profile, settings, preferences
 */

import { http } from '@/lib/http';
import type {
    User,
    UpdateUserRequest,
    UserPreferences
} from '@/types/api';

export const userApi = {
    /**
     * Get current user profile
     */
    getProfile: () =>
        http.get<User>('/user/profile'),

    /**
     * Update user profile
     */
    updateProfile: (data: UpdateUserRequest) =>
        http.put<User>('/user/profile', data),

    /**
     * Upload avatar
     */
    uploadAvatar: (file: File) => {
        const formData = new FormData();
        formData.append('avatar', file);
        return http.post<{ url: string }>('/user/avatar', formData, {
            headers: { 'Content-Type': 'multipart/form-data' }
        });
    },

    /**
     * Get user preferences
     */
    getPreferences: () =>
        http.get<UserPreferences>('/user/preferences'),

    /**
     * Update user preferences
     */
    updatePreferences: (data: Partial<UserPreferences>) =>
        http.put<UserPreferences>('/user/preferences', data),

    /**
     * Get user statistics
     */
    getStats: () =>
        http.get<{
            totalOrders: number;
            totalSpent: number;
            favoriteCount: number;
            reviewCount: number;
        }>('/user/stats'),

    /**
     * Delete account
     */
    deleteAccount: (password: string) =>
        http.delete<void>('/user/account', { data: { password } }),
};
