import apiClient from './client';
import type { ApiResponse } from './admin';

export interface LoginDto {
    username: string;
    password?: string;
}

export interface RegisterDto {
    name: string;  // 后端期望 name 字段，不是 username
    email: string;
    phone?: string;
    password: string;
    confirmPassword?: string;
    avatarUrl?: string;
    referralCode?: string;
}

export interface LoginResponse {
    token: string;
    user: {
        id: number;
        username: string;
        email: string;
        role: string;
    };
}

export const authApi = {
    login: (data: LoginDto) => apiClient.post<ApiResponse<LoginResponse>>('/auth/login', data),
    register: (data: RegisterDto) => apiClient.post('/auth/register', data),
    logout: () => apiClient.post('/auth/logout'),
    getMe: () => apiClient.get<ApiResponse<LoginResponse>>('/auth/me'),
    refresh: () => apiClient.post('/auth/refresh'),
};
