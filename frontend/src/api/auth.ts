import apiClient from './client';
import type { ApiResponse } from './admin';

export interface LoginDto {
    username: string;
    password?: string;
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
    register: (data: any) => apiClient.post('/auth/register', data),
    logout: () => apiClient.post('/auth/logout'),
    getMe: () => apiClient.get<ApiResponse<LoginResponse>>('/auth/me'),
    refresh: () => apiClient.post('/auth/refresh'),
};
