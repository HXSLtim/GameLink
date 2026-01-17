/**
 * Authentication API
 * Handles login, register, logout, token refresh
 */

import { http } from '@/lib/http';
import type {
    LoginRequest,
    LoginResponse,
    RegisterRequest,
    RegisterResponse,
    RefreshResponse,
    MeResponse
} from '@/types/api';

export const authApi = {
    /**
     * User login
     */
    login: (data: LoginRequest) =>
        http.post<LoginResponse>('/auth/login', data),

    /**
     * User registration
     */
    register: (data: RegisterRequest) =>
        http.post<RegisterResponse>('/auth/register', data),

    /**
     * User logout
     */
    logout: () =>
        http.post<void>('/auth/logout'),

    /**
     * Refresh access token
     */
    refresh: (refreshToken: string) =>
        http.post<RefreshResponse>('/auth/refresh', { refreshToken }),

    /**
     * Get current user info
     */
    me: () =>
        http.get<MeResponse>('/auth/me'),

    /**
     * Request password reset
     */
    forgotPassword: (email: string) =>
        http.post<void>('/auth/forgot-password', { email }),

    /**
     * Reset password with code
     */
    resetPassword: (code: string, newPassword: string) =>
        http.post<void>('/auth/reset-password', { code, newPassword }),

    /**
     * Change password (authenticated)
     */
    changePassword: (oldPassword: string, newPassword: string) =>
        http.put<void>('/auth/change-password', { oldPassword, newPassword }),

    /**
     * Send SMS verification code
     */
    sendSmsCode: (phone: string) =>
        http.post<void>('/auth/send-sms-code', { phone }),

    /**
     * Send email verification code
     */
    sendEmailCode: (email: string) =>
        http.post<void>('/auth/send-email-code', { email }),
};
