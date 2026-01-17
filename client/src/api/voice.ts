/**
 * Voice API
 * Handles voice chat, voice channels, voice settings
 */

import { http } from '@/lib/http';
import type {
    VoiceChannel,
    VoiceSession
} from '@/types/api';

export const voiceApi = {
    /**
     * Get voice channel info
     */
    getChannel: (channelId: number) =>
        http.get<VoiceChannel>(`/voice/channel/${channelId}`),

    /**
     * Join voice channel
     */
    joinChannel: (channelId: number) =>
        http.post<{
            sessionId: string;
            token: string;
            serverUrl: string;
        }>(`/voice/channel/${channelId}/join`),

    /**
     * Leave voice channel
     */
    leaveChannel: (channelId: number) =>
        http.post<void>(`/voice/channel/${channelId}/leave`),

    /**
     * Get active voice sessions
     */
    getActiveSessions: () =>
        http.get<VoiceSession[]>('/voice/sessions'),

    /**
     * Mute/unmute user
     */
    toggleMute: (channelId: number, muted: boolean) =>
        http.post<void>(`/voice/channel/${channelId}/mute`, { muted }),

    /**
     * Get voice settings
     */
    getSettings: () =>
        http.get<{
            inputDevice: string;
            outputDevice: string;
            inputVolume: number;
            outputVolume: number;
            echoCancellation: boolean;
            noiseSuppression: boolean;
        }>('/voice/settings'),

    /**
     * Update voice settings
     */
    updateSettings: (settings: {
        inputDevice?: string;
        outputDevice?: string;
        inputVolume?: number;
        outputVolume?: number;
        echoCancellation?: boolean;
        noiseSuppression?: boolean;
    }) =>
        http.put<void>('/voice/settings', settings),
};
