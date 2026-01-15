import { describe, it, expect, vi, beforeEach } from 'vitest';
import { HttpClient } from './http';
import { useAuthStore } from '@/stores';

// Correct approach:
const { mockAxiosInstance } = vi.hoisted(() => {
    const instance = {
        interceptors: {
            request: { use: vi.fn() },
            response: { use: vi.fn() },
        },
        get: vi.fn(),
        post: vi.fn(),
        put: vi.fn(),
        delete: vi.fn(),
        patch: vi.fn(),
    };
    return { mockAxiosInstance: instance };
});

vi.mock('axios', () => ({
    default: {
        create: vi.fn(() => mockAxiosInstance),
    },
}));

// Mock store
vi.mock('@/stores', () => ({
    useAuthStore: {
        getState: vi.fn(),
    },
}));

describe('HttpClient', () => {
    let httpClient: HttpClient;
    let requestInterceptor: any;
    let responseInterceptor: any;

    beforeEach(() => {
        vi.clearAllMocks();

        // Capture interceptors
        // We reuse the hoisted instance
        (mockAxiosInstance.interceptors.request.use as any).mockImplementation((success: any, fail: any) => {
            requestInterceptor = { success, fail };
        });
        (mockAxiosInstance.interceptors.response.use as any).mockImplementation((success: any, fail: any) => {
            responseInterceptor = { success, fail };
        });

        // Create new instance for each test
        httpClient = new HttpClient();
    });

    it('should setup interceptors on initialization', () => {
        expect(mockAxiosInstance.interceptors.request.use).toHaveBeenCalled();
        expect(mockAxiosInstance.interceptors.response.use).toHaveBeenCalled();
    });

    describe('Interceptors', () => {
        it('should add Authorization header if token exists', async () => {
            const token = 'test-token';
            (useAuthStore.getState as any).mockReturnValue({ token });

            const config = { headers: {} };
            const result = requestInterceptor.success(config);

            expect(result.headers.Authorization).toBe(`Bearer ${token}`);
        });

        it('should not add Authorization header if no token', async () => {
            (useAuthStore.getState as any).mockReturnValue({ token: null });

            const config = { headers: {} };
            const result = requestInterceptor.success(config);

            expect(result.headers.Authorization).toBeUndefined();
        });

        it('should unwrap API success response', async () => {
            const apiResponse = {
                data: {
                    success: true,
                    code: 200,
                    message: 'ok',
                    data: { id: 1, name: 'test' }
                }
            };

            // Call the success response interceptor
            const result = responseInterceptor.success(apiResponse);
            // The interceptor just returns the response, the unwrap happens in methods
            expect(result).toBe(apiResponse);
        });
    });

    describe('Methods', () => {
        it('methods should unwrap response data', async () => {
            const mockData = { id: 123 };
            const apiResponse = {
                data: {
                    success: true,
                    code: 200,
                    message: 'Success',
                    data: mockData
                }
            };

            mockAxiosInstance.get.mockResolvedValue(apiResponse);

            const result = await httpClient.get('/test');
            expect(result).toEqual(mockData);
        });

        it('should return raw data if not matching standard format', async () => {
            const rawData = { some: 'data' };
            mockAxiosInstance.get.mockResolvedValue({ data: rawData });

            const result = await httpClient.get('/test');
            expect(result).toEqual(rawData);
        });
    });
});
