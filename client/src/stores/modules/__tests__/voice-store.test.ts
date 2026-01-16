import { describe, it, expect, vi, beforeEach } from 'vitest';

// Mock TRTC SDK before importing voice-store
const mockLocalStream = {
    initialize: vi.fn().mockResolvedValue(undefined),
    close: vi.fn(),
    muteAudio: vi.fn(),
    unmuteAudio: vi.fn(),
    switchDevice: vi.fn().mockResolvedValue(undefined),
    play: vi.fn(),
};

const mockClient = {
    on: vi.fn(),
    join: vi.fn().mockResolvedValue(undefined),
    leave: vi.fn().mockResolvedValue(undefined),
    publish: vi.fn().mockResolvedValue(undefined),
    subscribe: vi.fn().mockResolvedValue(undefined),
};

vi.mock('trtc-js-sdk', () => ({
    default: {
        createClient: vi.fn(() => mockClient),
        createStream: vi.fn(() => mockLocalStream),
    },
}));

// Mock http module
vi.mock('@/lib/http', () => ({
    http: {
        get: vi.fn(),
        post: vi.fn(),
    },
}));

import {
    useVoiceStore,
    type VoiceToken,
    type VoiceStatus,
    selectIsInVoice,
    selectIsMuted,
    selectIsDeafened,
    selectVoiceMembers,
    selectVoiceStatus,
} from '../voice-store';

// Mock navigator.mediaDevices
const mockEnumerateDevices = vi.fn();
Object.defineProperty(global.navigator, 'mediaDevices', {
    value: {
        enumerateDevices: mockEnumerateDevices,
    },
    writable: true,
});

import { http } from '@/lib/http';

const mockVoiceToken: VoiceToken = {
    userSig: 'test-user-sig',
    sdkAppId: 1400000000,
    userId: 'user-100',
    roomId: 'room-1',
    expireAt: Date.now() + 3600000,
};

const mockVoiceStatus: VoiceStatus = {
    roomId: 1,
    voiceEnabled: true,
    voiceRoomId: 'voice-room-1',
    provider: 'trtc',
    sdkAppId: 1400000000,
    startedAt: '2024-01-01T00:00:00Z',
    duration: 0,
    maxMembers: 10,
};

describe('voice-store', () => {
    beforeEach(() => {
        useVoiceStore.getState().reset();
        vi.clearAllMocks();
    });

    // ========== getVoiceToken Tests ==========

    describe('getVoiceToken', () => {
        it('should fetch voice token successfully', async () => {
            vi.mocked(http.get).mockResolvedValueOnce(mockVoiceToken);

            const result = await useVoiceStore.getState().getVoiceToken(1);

            expect(http.get).toHaveBeenCalledWith('/user/rooms/1/voice/token');
            expect(result).toEqual(mockVoiceToken);
            expect(useVoiceStore.getState().voiceToken).toEqual(mockVoiceToken);
            expect(useVoiceStore.getState().error).toBeNull();
        });

        it('should handle fetch token error', async () => {
            const error = new Error('Token fetch failed');
            vi.mocked(http.get).mockRejectedValueOnce(error);

            await expect(useVoiceStore.getState().getVoiceToken(1)).rejects.toThrow('Token fetch failed');
            expect(useVoiceStore.getState().error).toBe('Token fetch failed');
        });

        it('should use default error message when error has no message', async () => {
            vi.mocked(http.get).mockRejectedValueOnce({});

            await expect(useVoiceStore.getState().getVoiceToken(1)).rejects.toBeDefined();
            expect(useVoiceStore.getState().error).toBe('获取语音Token失败');
        });
    });

    // ========== getVoiceStatus Tests ==========

    describe('getVoiceStatus', () => {
        it('should fetch voice status successfully', async () => {
            vi.mocked(http.get).mockResolvedValueOnce(mockVoiceStatus);

            const result = await useVoiceStore.getState().getVoiceStatus(1);

            expect(http.get).toHaveBeenCalledWith('/user/rooms/1/voice/status');
            expect(result).toEqual(mockVoiceStatus);
            expect(useVoiceStore.getState().voiceStatus).toEqual(mockVoiceStatus);
            expect(useVoiceStore.getState().error).toBeNull();
        });

        it('should handle fetch status error', async () => {
            const error = new Error('Status fetch failed');
            vi.mocked(http.get).mockRejectedValueOnce(error);

            await expect(useVoiceStore.getState().getVoiceStatus(1)).rejects.toThrow('Status fetch failed');
            expect(useVoiceStore.getState().error).toBe('Status fetch failed');
        });
    });

    // ========== startVoice Tests ==========

    describe('startVoice', () => {
        it('should start voice successfully', async () => {
            vi.mocked(http.post).mockResolvedValueOnce(undefined);
            vi.mocked(http.get).mockResolvedValueOnce(mockVoiceStatus);

            await useVoiceStore.getState().startVoice(1);

            expect(http.post).toHaveBeenCalledWith('/user/rooms/1/voice/start');
            expect(http.get).toHaveBeenCalledWith('/user/rooms/1/voice/status');
            expect(useVoiceStore.getState().voiceStatus).toEqual(mockVoiceStatus);
        });

        it('should handle start voice error', async () => {
            const error = new Error('Failed to start voice');
            vi.mocked(http.post).mockRejectedValueOnce(error);

            await expect(useVoiceStore.getState().startVoice(1)).rejects.toThrow('Failed to start voice');
            expect(useVoiceStore.getState().error).toBe('Failed to start voice');
        });
    });

    // ========== stopVoice Tests ==========

    describe('stopVoice', () => {
        it('should stop voice successfully', async () => {
            vi.mocked(http.post).mockResolvedValueOnce(undefined);
            vi.mocked(http.get).mockResolvedValueOnce({ ...mockVoiceStatus, voiceEnabled: false });

            await useVoiceStore.getState().stopVoice(1);

            expect(http.post).toHaveBeenCalledWith('/user/rooms/1/voice/stop');
            expect(http.get).toHaveBeenCalledWith('/user/rooms/1/voice/status');
        });

        it('should handle stop voice error', async () => {
            const error = new Error('Failed to stop voice');
            vi.mocked(http.post).mockRejectedValueOnce(error);

            await expect(useVoiceStore.getState().stopVoice(1)).rejects.toThrow('Failed to stop voice');
            expect(useVoiceStore.getState().error).toBe('Failed to stop voice');
        });
    });

    // ========== joinVoice Tests ==========

    describe('joinVoice', () => {
        it('should join voice room successfully', async () => {
            vi.mocked(http.get).mockResolvedValueOnce(mockVoiceToken);

            await useVoiceStore.getState().joinVoice(1);

            expect(http.get).toHaveBeenCalledWith('/user/rooms/1/voice/token');
            expect(useVoiceStore.getState().isInVoice).toBe(true);
            expect(useVoiceStore.getState().isConnecting).toBe(false);
            expect(useVoiceStore.getState().currentRoomId).toBe(1);
            expect(useVoiceStore.getState().voiceToken).toEqual(mockVoiceToken);
        });

        it('should throw error if already in voice room', async () => {
            useVoiceStore.setState({ isInVoice: true });

            await expect(useVoiceStore.getState().joinVoice(1)).rejects.toThrow('已在语音房间中');
        });

        it('should set isConnecting during join', async () => {
            let connectingState = false;
            vi.mocked(http.get).mockImplementationOnce(async () => {
                connectingState = useVoiceStore.getState().isConnecting;
                return mockVoiceToken;
            });

            await useVoiceStore.getState().joinVoice(1);

            expect(connectingState).toBe(true);
            expect(useVoiceStore.getState().isConnecting).toBe(false);
        });

        it('should handle join voice error', async () => {
            const error = new Error('Join failed');
            vi.mocked(http.get).mockRejectedValueOnce(error);

            await expect(useVoiceStore.getState().joinVoice(1)).rejects.toThrow('Join failed');
            expect(useVoiceStore.getState().error).toBe('Join failed');
            expect(useVoiceStore.getState().isConnecting).toBe(false);
            expect(useVoiceStore.getState().isInVoice).toBe(false);
        });
    });

    // ========== leaveVoice Tests ==========

    describe('leaveVoice', () => {
        it('should leave voice room successfully', async () => {
            useVoiceStore.setState({
                isInVoice: true,
                currentRoomId: 1,
                voiceToken: mockVoiceToken,
                isMuted: true,
                isDeafened: true,
            });

            await useVoiceStore.getState().leaveVoice();

            expect(useVoiceStore.getState().isInVoice).toBe(false);
            expect(useVoiceStore.getState().currentRoomId).toBeNull();
            expect(useVoiceStore.getState().voiceToken).toBeNull();
            expect(useVoiceStore.getState().isMuted).toBe(false);
            expect(useVoiceStore.getState().isDeafened).toBe(false);
        });

        it('should do nothing if not in voice room', async () => {
            useVoiceStore.setState({ isInVoice: false });

            await useVoiceStore.getState().leaveVoice();

            expect(useVoiceStore.getState().isInVoice).toBe(false);
        });
    });

    // ========== toggleMute Tests ==========

    describe('toggleMute', () => {
        it('should toggle mute from false to true', () => {
            useVoiceStore.setState({ isMuted: false });

            useVoiceStore.getState().toggleMute();

            expect(useVoiceStore.getState().isMuted).toBe(true);
        });

        it('should toggle mute from true to false', () => {
            useVoiceStore.setState({ isMuted: true });

            useVoiceStore.getState().toggleMute();

            expect(useVoiceStore.getState().isMuted).toBe(false);
        });
    });

    // ========== toggleDeafen Tests ==========

    describe('toggleDeafen', () => {
        it('should toggle deafen from false to true', () => {
            useVoiceStore.setState({ isDeafened: false });

            useVoiceStore.getState().toggleDeafen();

            expect(useVoiceStore.getState().isDeafened).toBe(true);
        });

        it('should toggle deafen from true to false', () => {
            useVoiceStore.setState({ isDeafened: true });

            useVoiceStore.getState().toggleDeafen();

            expect(useVoiceStore.getState().isDeafened).toBe(false);
        });
    });

    // ========== Device Selection Tests ==========

    describe('setInputDevice', () => {
        it('should set selected input device', () => {
            useVoiceStore.getState().setInputDevice('device-1');

            expect(useVoiceStore.getState().selectedInputDevice).toBe('device-1');
        });
    });

    describe('setOutputDevice', () => {
        it('should set selected output device', () => {
            useVoiceStore.getState().setOutputDevice('device-2');

            expect(useVoiceStore.getState().selectedOutputDevice).toBe('device-2');
        });
    });

    describe('refreshDevices', () => {
        it('should refresh audio devices list', async () => {
            const mockDevices = [
                { deviceId: 'input-1', kind: 'audioinput', label: 'Mic 1' },
                { deviceId: 'input-2', kind: 'audioinput', label: 'Mic 2' },
                { deviceId: 'output-1', kind: 'audiooutput', label: 'Speaker 1' },
                { deviceId: 'video-1', kind: 'videoinput', label: 'Camera 1' },
            ] as MediaDeviceInfo[];
            mockEnumerateDevices.mockResolvedValueOnce(mockDevices);

            await useVoiceStore.getState().refreshDevices();

            expect(useVoiceStore.getState().audioInputDevices).toHaveLength(2);
            expect(useVoiceStore.getState().audioOutputDevices).toHaveLength(1);
        });

        it('should handle device enumeration error', async () => {
            const consoleSpy = vi.spyOn(console, 'error').mockImplementation(() => { });
            mockEnumerateDevices.mockRejectedValueOnce(new Error('Permission denied'));

            await useVoiceStore.getState().refreshDevices();

            expect(consoleSpy).toHaveBeenCalled();
            consoleSpy.mockRestore();
        });
    });

    // ========== Reset Tests ==========

    describe('reset', () => {
        it('should reset store to initial state', () => {
            useVoiceStore.setState({
                isInVoice: true,
                isMuted: true,
                isDeafened: true,
                isConnecting: true,
                currentRoomId: 1,
                voiceToken: mockVoiceToken,
                voiceStatus: mockVoiceStatus,
                voiceMembers: [{ userId: 'user-1', nickname: 'Test', isSpeaking: false, isMuted: false, isDeafened: false }],
                error: 'Some error',
                selectedInputDevice: 'device-1',
                selectedOutputDevice: 'device-2',
            });

            useVoiceStore.getState().reset();

            expect(useVoiceStore.getState().isInVoice).toBe(false);
            expect(useVoiceStore.getState().isMuted).toBe(false);
            expect(useVoiceStore.getState().isDeafened).toBe(false);
            expect(useVoiceStore.getState().isConnecting).toBe(false);
            expect(useVoiceStore.getState().currentRoomId).toBeNull();
            expect(useVoiceStore.getState().voiceToken).toBeNull();
            expect(useVoiceStore.getState().voiceStatus).toBeNull();
            expect(useVoiceStore.getState().voiceMembers).toEqual([]);
            expect(useVoiceStore.getState().error).toBeNull();
            expect(useVoiceStore.getState().selectedInputDevice).toBeNull();
            expect(useVoiceStore.getState().selectedOutputDevice).toBeNull();
        });
    });

    // ========== Selectors Tests ==========

    describe('selectors', () => {
        it('selectIsInVoice should return isInVoice state', () => {
            const state = { ...useVoiceStore.getState(), isInVoice: true };
            expect(selectIsInVoice(state)).toBe(true);
        });

        it('selectIsMuted should return isMuted state', () => {
            const state = { ...useVoiceStore.getState(), isMuted: true };
            expect(selectIsMuted(state)).toBe(true);
        });

        it('selectIsDeafened should return isDeafened state', () => {
            const state = { ...useVoiceStore.getState(), isDeafened: true };
            expect(selectIsDeafened(state)).toBe(true);
        });

        it('selectVoiceMembers should return voiceMembers state', () => {
            const members = [{ userId: 'user-1', nickname: 'Test', isSpeaking: false, isMuted: false, isDeafened: false }];
            const state = { ...useVoiceStore.getState(), voiceMembers: members };
            expect(selectVoiceMembers(state)).toEqual(members);
        });

        it('selectVoiceStatus should return voiceStatus state', () => {
            const state = { ...useVoiceStore.getState(), voiceStatus: mockVoiceStatus };
            expect(selectVoiceStatus(state)).toEqual(mockVoiceStatus);
        });
    });
});
