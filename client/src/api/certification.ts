/**
 * Certification API
 * Handles player certification, real-name verification, skill authentication
 */

import { http } from '@/lib/http';
import type {
    Certification,
    RealNameVerification,
    SkillAuthentication
} from '@/types/api';

export const certificationApi = {
    /**
     * Get player certifications
     */
    getCertifications: () =>
        http.get<Certification[]>('/certification/list'),

    /**
     * Submit real-name verification
     */
    submitRealName: (data: {
        realName: string;
        idCard: string;
        idCardFront: string;
        idCardBack: string;
    }) =>
        http.post<{ verificationId: number }>('/certification/real-name', data),

    /**
     * Get real-name verification status
     */
    getRealNameStatus: () =>
        http.get<RealNameVerification>('/certification/real-name/status'),

    /**
     * Submit skill authentication
     */
    submitSkillAuth: (data: {
        gameId: number;
        rankId: number;
        proofImages: string[];
        description: string;
    }) =>
        http.post<{ authId: number }>('/certification/skill-auth', data),

    /**
     * Get skill authentication status
     */
    getSkillAuthStatus: (gameId: number) =>
        http.get<SkillAuthentication>(`/certification/skill-auth/${gameId}/status`),

    /**
     * Upload certification image
     */
    uploadImage: (file: File, type: 'id-card' | 'skill-proof') => {
        const formData = new FormData();
        formData.append('image', file);
        formData.append('type', type);
        return http.post<{ url: string }>('/certification/upload/image', formData, {
            headers: { 'Content-Type': 'multipart/form-data' }
        });
    },
};
