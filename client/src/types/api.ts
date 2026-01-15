/**
 * API Response Types for GameLink Client
 */

// Auth API responses
export interface LoginResponse {
    token: string;
    refreshToken: string;
    user: {
        id: string;
        username: string;
        avatar: string;
        email?: string;
        nickname?: string;
    };
    role?: string;
    permissions?: string[];
}

export interface RegisterResponse extends LoginResponse {}

export interface RefreshResponse {
    token: string;
    refreshToken?: string;
    user?: {
        id: string;
        username: string;
        avatar: string;
        email?: string;
        nickname?: string;
    };
}

export interface MeResponse {
    user?: {
        id: string;
        username: string;
        avatar: string;
        email?: string;
        nickname?: string;
        role?: string;
    };
    id?: string;
    username?: string;
    avatar?: string;
    role?: string;
}

// Error response
export interface ApiError {
    response?: {
        data?: {
            message?: string;
        };
    };
    message?: string;
}

// Helper to extract error message
export function getErrorMessage(err: unknown): string {
    if (err && typeof err === 'object') {
        const apiErr = err as ApiError;
        return apiErr.response?.data?.message || apiErr.message || 'An error occurred';
    }
    return 'An error occurred';
}
