import { create } from 'zustand';
import { http } from '@/lib/http';
import TRTC from 'trtc-js-sdk';

// Development-only logging helper
const devLog = (...args: unknown[]) => {
    if (import.meta.env.DEV) {
        console.log('[Voice]', ...args);
    }
};

const devError = (...args: unknown[]) => {
    if (import.meta.env.DEV) {
        console.error('[Voice]', ...args);
    }
};

// ============================================================================
// Types
// ============================================================================

export interface VoiceToken {
    userSig: string;
    sdkAppId: number;
    userId: string;
    roomId: string;
    expireAt: number;
}

export interface VoiceStatus {
    roomId: number;
    voiceEnabled: boolean;
    voiceRoomId: string;
    provider: string;
    sdkAppId: number;
    startedAt?: string;
    duration: number;
    maxMembers: number;
}

export interface VoiceMember {
    userId: string;
    nickname: string;
    avatarUrl?: string;
    isSpeaking: boolean;
    isMuted: boolean;
    isDeafened: boolean;
}

// TRTC Types
type TRTCClient = ReturnType<typeof TRTC.createClient>;
type TRTCLocalStream = ReturnType<typeof TRTC.createStream>;

export interface VoiceState {
    // 状态
    isInVoice: boolean;
    isMuted: boolean;
    isDeafened: boolean;
    isConnecting: boolean;
    currentRoomId: number | null;
    voiceToken: VoiceToken | null;
    voiceStatus: VoiceStatus | null;
    voiceMembers: VoiceMember[];
    error: string | null;

    // 音频设备
    audioInputDevices: MediaDeviceInfo[];
    audioOutputDevices: MediaDeviceInfo[];
    selectedInputDevice: string | null;
    selectedOutputDevice: string | null;

    // TRTC 客户端
    trtcClient: TRTCClient | null;
    localStream: TRTCLocalStream | null;

    // Actions
    getVoiceToken: (roomId: number) => Promise<VoiceToken>;
    getVoiceStatus: (roomId: number) => Promise<VoiceStatus>;
    startVoice: (roomId: number) => Promise<void>;
    stopVoice: (roomId: number) => Promise<void>;
    joinVoice: (roomId: number) => Promise<void>;
    leaveVoice: () => Promise<void>;
    toggleMute: () => void;
    toggleDeafen: () => void;
    setInputDevice: (deviceId: string) => Promise<void>;
    setOutputDevice: (deviceId: string) => void;
    refreshDevices: () => Promise<void>;
    reset: () => void;
}

// ============================================================================
// Initial State
// ============================================================================

const initialState = {
    isInVoice: false,
    isMuted: false,
    isDeafened: false,
    isConnecting: false,
    currentRoomId: null,
    voiceToken: null,
    voiceStatus: null,
    voiceMembers: [],
    error: null,
    audioInputDevices: [],
    audioOutputDevices: [],
    selectedInputDevice: null,
    selectedOutputDevice: null,
    trtcClient: null,
    localStream: null,
};

// ============================================================================
// Store
// ============================================================================

export const useVoiceStore = create<VoiceState>((set, get) => ({
    ...initialState,

    // 获取语音Token
    getVoiceToken: async (roomId: number) => {
        try {
            const response = await http.get<VoiceToken>(
                `/user/rooms/${roomId}/voice/token`
            );
            set({ voiceToken: response, error: null });
            return response;
        } catch (error: unknown) {
            const message = error instanceof Error ? error.message : '获取语音Token失败';
            set({ error: message });
            throw error;
        }
    },

    // 获取语音状态
    getVoiceStatus: async (roomId: number) => {
        try {
            const response = await http.get<VoiceStatus>(
                `/user/rooms/${roomId}/voice/status`
            );
            set({ voiceStatus: response, error: null });
            return response;
        } catch (error: unknown) {
            const message = error instanceof Error ? error.message : '获取语音状态失败';
            set({ error: message });
            throw error;
        }
    },

    // 开启语音（房主）
    startVoice: async (roomId: number) => {
        try {
            await http.post(`/user/rooms/${roomId}/voice/start`);
            await get().getVoiceStatus(roomId);
        } catch (error: unknown) {
            const message = error instanceof Error ? error.message : '开启语音失败';
            set({ error: message });
            throw error;
        }
    },

    // 关闭语音（房主）
    stopVoice: async (roomId: number) => {
        try {
            await http.post(`/user/rooms/${roomId}/voice/stop`);
            await get().getVoiceStatus(roomId);
        } catch (error: unknown) {
            const message = error instanceof Error ? error.message : '关闭语音失败';
            set({ error: message });
            throw error;
        }
    },

    // 加入语音房间
    joinVoice: async (roomId: number) => {
        const state = get();
        if (state.isInVoice) {
            throw new Error('已在语音房间中');
        }

        set({ isConnecting: true, error: null });

        try {
            // 1. 获取Token
            const token = await get().getVoiceToken(roomId);

            // 2. 初始化TRTC客户端
            const client = TRTC.createClient({
                mode: 'rtc',
                sdkAppId: token.sdkAppId,
                userId: token.userId,
                userSig: token.userSig,
            });

            // 设置事件监听
            client.on('stream-added', (event) => {
                const remoteStream = event.stream;
                const remoteUserId = remoteStream.getUserId();
                devLog('Remote stream added:', remoteUserId);
                client.subscribe(remoteStream);
            });

            client.on('stream-subscribed', (event) => {
                const remoteStream = event.stream;
                remoteStream.play('remote-audio-container');
                devLog('Remote stream subscribed');
            });

            client.on('stream-removed', (event) => {
                const remoteStream = event.stream;
                devLog('Remote stream removed:', remoteStream.getUserId());
            });

            client.on('peer-join', (event) => {
                devLog('Peer joined:', event.userId);
            });

            client.on('peer-leave', (event) => {
                devLog('Peer left:', event.userId);
            });

            // 3. 创建本地音频流
            const localStream = TRTC.createStream({
                audio: true,
                video: false,
            });
            await localStream.initialize();

            // 4. 加入房间
            await client.join({ roomId: token.roomId });

            // 5. 发布本地流
            await client.publish(localStream);

            devLog('Successfully joined voice room:', token.roomId);

            set({
                isInVoice: true,
                isConnecting: false,
                currentRoomId: roomId,
                voiceToken: token,
                trtcClient: client,
                localStream: localStream,
            });
        } catch (error: unknown) {
            const message = error instanceof Error ? error.message : '加入语音失败';
            set({
                isConnecting: false,
                error: message,
            });
            throw error;
        }
    },

    // 离开语音房间
    leaveVoice: async () => {
        const state = get();
        if (!state.isInVoice) {
            return;
        }

        try {
            // 停止本地流
            if (state.localStream) {
                state.localStream.close();
            }

            // 离开房间
            if (state.trtcClient) {
                await state.trtcClient.leave();
            }

            devLog('Left voice room');

            set({
                isInVoice: false,
                currentRoomId: null,
                voiceToken: null,
                trtcClient: null,
                localStream: null,
                voiceMembers: [],
                isMuted: false,
                isDeafened: false,
            });
        } catch (error: unknown) {
            const message = error instanceof Error ? error.message : '离开语音失败';
            set({ error: message });
            throw error;
        }
    },

    // 切换静音
    toggleMute: () => {
        const state = get();
        const newMuted = !state.isMuted;

        if (state.localStream) {
            if (newMuted) {
                state.localStream.muteAudio();
            } else {
                state.localStream.unmuteAudio();
            }
        }

        devLog('Mute:', newMuted);
        set({ isMuted: newMuted });
    },

    // 切换耳机静音
    toggleDeafen: () => {
        const state = get();
        const newDeafened = !state.isDeafened;

        // 实际实现需要控制远程流的音量
        devLog('Deafen:', newDeafened);
        set({ isDeafened: newDeafened });
    },

    // 设置输入设备
    setInputDevice: async (deviceId: string) => {
        const state = get();
        if (state.localStream) {
            await state.localStream.switchDevice('audio', deviceId);
        }
        devLog('Set input device:', deviceId);
        set({ selectedInputDevice: deviceId });
    },

    // 设置输出设备
    setOutputDevice: (deviceId: string) => {
        devLog('Set output device:', deviceId);
        set({ selectedOutputDevice: deviceId });
    },

    // 刷新设备列表
    refreshDevices: async () => {
        try {
            const devices = await navigator.mediaDevices.enumerateDevices();
            const audioInputs = devices.filter(
                (d) => d.kind === 'audioinput'
            );
            const audioOutputs = devices.filter(
                (d) => d.kind === 'audiooutput'
            );

            set({
                audioInputDevices: audioInputs,
                audioOutputDevices: audioOutputs,
            });
        } catch (error) {
            devError('Failed to enumerate devices:', error);
        }
    },

    // 重置状态
    reset: () => {
        const state = get();
        // 清理TRTC资源
        if (state.localStream) {
            state.localStream.close();
        }
        if (state.trtcClient) {
            state.trtcClient.leave().catch(() => {
                // 忽略离开房间时的错误
            });
        }
        set(initialState);
    },
}));

// ============================================================================
// Selectors
// ============================================================================

export const selectIsInVoice = (state: VoiceState) => state.isInVoice;
export const selectIsMuted = (state: VoiceState) => state.isMuted;
export const selectIsDeafened = (state: VoiceState) => state.isDeafened;
export const selectVoiceMembers = (state: VoiceState) => state.voiceMembers;
export const selectVoiceStatus = (state: VoiceState) => state.voiceStatus;
