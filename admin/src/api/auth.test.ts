/**
 * Authentication API Tests
 *
 * Tests for authentication endpoints including login, logout,
 * token refresh, and current user retrieval
 *
 * Requirements:
 * - Login with valid credentials should return token and user data
 * - Login with invalid credentials should return 401
 * - Logout should clear tokens
 * - getMe should return current user data
 * - Token refresh should update stored tokens
 */
import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { authApi, type LoginDto, type RegisterDto, type LoginResponse } from './auth';
import apiClient from './client';
import type { AxiosResponse } from 'axios';

// Mock the API client
vi.mock('./client', () => ({
  default: {
    get: vi.fn(),
    post: vi.fn(),
    put: vi.fn(),
    delete: vi.fn(),
    interceptors: {
      request: { use: vi.fn(), handlers: [] },
      response: { use: vi.fn(), handlers: [] },
    },
    defaults: {
      baseURL: '/api/v1',
      timeout: 10000,
      headers: {},
    },
  },
}));

// Mock crypto utilities
vi.mock('@/utils/crypto', () => ({
  encryptRequest: vi.fn((data) => data),
  shouldEncrypt: vi.fn(() => false),
  logger: {
    error: vi.fn(),
    warn: vi.fn(),
  },
}));

describe('Auth API', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    localStorage.clear();
    sessionStorage.clear();
  });

  afterEach(() => {
    localStorage.clear();
    sessionStorage.clear();
  });

  describe('login', () => {
    const mockLoginData: LoginDto = {
      username: 'admin',
      password: 'password123',
    };

    const mockLoginResponse: LoginResponse = {
      token: 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.valid_token',
      user: {
        id: 1,
        username: 'admin',
        email: 'admin@gamelink.com',
        role: 'admin',
      },
    };

    it('should login successfully with valid credentials', async () => {
      const axiosResponse: AxiosResponse = {
        data: {
          success: true,
          code: 200,
          message: 'Login successful',
          data: mockLoginResponse,
        },
        status: 200,
        statusText: 'OK',
        headers: {},
        config: {} as any,
      };

      vi.mocked(apiClient.post).mockResolvedValueOnce(axiosResponse);

      const result = await authApi.login(mockLoginData);

      expect(apiClient.post).toHaveBeenCalledWith('/auth/login', mockLoginData);
      expect(result.data.data).toEqual(mockLoginResponse);
    });

    it('should handle login with missing password', async () => {
      const loginDataWithoutPassword: LoginDto = {
        username: 'admin',
      };

      const axiosResponse: AxiosResponse = {
        data: {
          success: true,
          code: 200,
          message: 'Login successful',
          data: mockLoginResponse,
        },
        status: 200,
        statusText: 'OK',
        headers: {},
        config: {} as any,
      };

      vi.mocked(apiClient.post).mockResolvedValueOnce(axiosResponse);

      const result = await authApi.login(loginDataWithoutPassword);

      expect(apiClient.post).toHaveBeenCalledWith('/auth/login', loginDataWithoutPassword);
      expect(result.data.data).toEqual(mockLoginResponse);
    });

    it('should handle 401 error with invalid credentials', async () => {
      const error = {
        response: {
          status: 401,
          data: {
            success: false,
            code: 401,
            message: 'Invalid username or password',
          },
        },
      };

      vi.mocked(apiClient.post).mockRejectedValueOnce(error);

      await expect(authApi.login(mockLoginData)).rejects.toMatchObject({
        response: {
          status: 401,
          data: {
            success: false,
            message: 'Invalid username or password',
          },
        },
      });

      expect(apiClient.post).toHaveBeenCalledWith('/auth/login', mockLoginData);
    });

    it('should handle 400 validation error', async () => {
      const error = {
        response: {
          status: 400,
          data: {
            success: false,
            code: 400,
            message: 'Username is required',
          },
        },
      };

      vi.mocked(apiClient.post).mockRejectedValueOnce(error);

      await expect(authApi.login({ username: '', password: 'test' })).rejects.toMatchObject({
        response: {
          status: 400,
        },
      });
    });

    it('should handle network errors', async () => {
      const networkError = new Error('Network Error');
      vi.mocked(apiClient.post).mockRejectedValueOnce(networkError);

      await expect(authApi.login(mockLoginData)).rejects.toThrow('Network Error');
    });

    it('should handle timeout errors', async () => {
      const timeoutError = new Error('timeout of 10000ms exceeded');
      vi.mocked(apiClient.post).mockRejectedValueOnce(timeoutError);

      await expect(authApi.login(mockLoginData)).rejects.toThrow();
    });
  });

  describe('register', () => {
    const mockRegisterData: RegisterDto = {
      name: 'Test User',
      email: 'test@example.com',
      phone: '+1234567890',
      password: 'password123',
      confirmPassword: 'password123',
      avatarUrl: 'https://example.com/avatar.jpg',
      referralCode: 'REF123',
    };

    it('should register successfully with valid data', async () => {
      const axiosResponse: AxiosResponse = {
        data: {
          success: true,
          code: 201,
          message: 'Registration successful',
          data: {
            id: 2,
            name: 'Test User',
            email: 'test@example.com',
          },
        },
        status: 201,
        statusText: 'Created',
        headers: {},
        config: {} as any,
      };

      vi.mocked(apiClient.post).mockResolvedValueOnce(axiosResponse);

      const result = await authApi.register(mockRegisterData);

      expect(apiClient.post).toHaveBeenCalledWith('/auth/register', mockRegisterData);
      expect(result.status).toBe(201);
    });

    it('should handle registration with minimal required fields', async () => {
      const minimalData: RegisterDto = {
        name: 'Minimal User',
        email: 'minimal@example.com',
        password: 'password123',
      };

      const axiosResponse: AxiosResponse = {
        data: {
          success: true,
          code: 201,
          message: 'Registration successful',
          data: { id: 3 },
        },
        status: 201,
        statusText: 'Created',
        headers: {},
        config: {} as any,
      };

      vi.mocked(apiClient.post).mockResolvedValueOnce(axiosResponse);

      await authApi.register(minimalData);

      expect(apiClient.post).toHaveBeenCalledWith('/auth/register', minimalData);
    });

    it('should handle 409 conflict error (email already exists)', async () => {
      const error = {
        response: {
          status: 409,
          data: {
            success: false,
            code: 409,
            message: 'Email already exists',
          },
        },
      };

      vi.mocked(apiClient.post).mockRejectedValueOnce(error);

      await expect(authApi.register(mockRegisterData)).rejects.toMatchObject({
        response: {
          status: 409,
          data: {
            message: 'Email already exists',
          },
        },
      });
    });

    it('should handle validation error for weak password', async () => {
      const weakPasswordData = {
        ...mockRegisterData,
        password: '123',
      };

      const error = {
        response: {
          status: 400,
          data: {
            success: false,
            code: 400,
            message: 'Password must be at least 8 characters',
          },
        },
      };

      vi.mocked(apiClient.post).mockRejectedValueOnce(error);

      await expect(authApi.register(weakPasswordData)).rejects.toMatchObject({
        response: {
          status: 400,
        },
      });
    });
  });

  describe('logout', () => {
    it('should logout successfully', async () => {
      // Set tokens in storage
      localStorage.setItem('token', 'test_token');
      sessionStorage.setItem('auth_token', 'test_token');

      const axiosResponse: AxiosResponse = {
        data: {
          success: true,
          code: 200,
          message: 'Logout successful',
          data: null,
        },
        status: 200,
        statusText: 'OK',
        headers: {},
        config: {} as any,
      };

      vi.mocked(apiClient.post).mockResolvedValueOnce(axiosResponse);

      await authApi.logout();

      expect(apiClient.post).toHaveBeenCalledWith('/auth/logout');

      // Note: Actual token clearing is handled by the auth store or component
      // The API call itself doesn't clear tokens
    });

    it('should handle logout error gracefully', async () => {
      const error = {
        response: {
          status: 500,
          data: {
            success: false,
            message: 'Internal server error',
          },
        },
      };

      vi.mocked(apiClient.post).mockRejectedValueOnce(error);

      await expect(authApi.logout()).rejects.toMatchObject({
        response: {
          status: 500,
        },
      });
    });

    it('should handle logout when not authenticated', async () => {
      const error = {
        response: {
          status: 401,
          data: {
            success: false,
            message: 'Not authenticated',
          },
        },
      };

      vi.mocked(apiClient.post).mockRejectedValueOnce(error);

      await expect(authApi.logout()).rejects.toMatchObject({
        response: {
          status: 401,
        },
      });
    });
  });

  describe('getMe', () => {
    const mockCurrentUser: LoginResponse = {
      token: 'current_token',
      user: {
        id: 1,
        username: 'admin',
        email: 'admin@gamelink.com',
        role: 'admin',
      },
    };

    it('should get current user successfully', async () => {
      const axiosResponse: AxiosResponse = {
        data: {
          success: true,
          code: 200,
          message: 'Success',
          data: mockCurrentUser,
        },
        status: 200,
        statusText: 'OK',
        headers: {},
        config: {} as any,
      };

      vi.mocked(apiClient.get).mockResolvedValueOnce(axiosResponse);

      const result = await authApi.getMe();

      expect(apiClient.get).toHaveBeenCalledWith('/auth/me');
      expect(result.data.data).toEqual(mockCurrentUser);
    });

    it('should handle 401 error (not authenticated)', async () => {
      const error = {
        response: {
          status: 401,
          data: {
            success: false,
            code: 401,
            message: 'Not authenticated',
          },
        },
      };

      vi.mocked(apiClient.get).mockRejectedValueOnce(error);

      await expect(authApi.getMe()).rejects.toMatchObject({
        response: {
          status: 401,
        },
      });
    });

    it('should handle token expired scenario', async () => {
      const error = {
        response: {
          status: 401,
          data: {
            success: false,
            code: 401,
            message: 'Token has expired',
          },
        },
      };

      vi.mocked(apiClient.get).mockRejectedValueOnce(error);

      await expect(authApi.getMe()).rejects.toMatchObject({
        response: {
          data: {
            message: 'Token has expired',
          },
        },
      });
    });
  });

  describe('refresh', () => {
    const mockNewToken = 'new_refreshed_token';

    it('should refresh token successfully', async () => {
      const axiosResponse: AxiosResponse = {
        data: {
          success: true,
          code: 200,
          message: 'Token refreshed',
          data: {
            token: mockNewToken,
          },
        },
        status: 200,
        statusText: 'OK',
        headers: {},
        config: {} as any,
      };

      vi.mocked(apiClient.post).mockResolvedValueOnce(axiosResponse);

      const result = await authApi.refresh();

      expect(apiClient.post).toHaveBeenCalledWith('/auth/refresh');
      expect(result.data.data.token).toBe(mockNewToken);
    });

    it('should handle refresh failure with expired token', async () => {
      const error = {
        response: {
          status: 401,
          data: {
            success: false,
            code: 401,
            message: 'Refresh token has expired',
          },
        },
      };

      vi.mocked(apiClient.post).mockRejectedValueOnce(error);

      await expect(authApi.refresh()).rejects.toMatchObject({
        response: {
          status: 401,
          data: {
            message: 'Refresh token has expired',
          },
        },
      });
    });

    it('should handle refresh failure with invalid token', async () => {
      const error = {
        response: {
          status: 401,
          data: {
            success: false,
            code: 401,
            message: 'Invalid refresh token',
          },
        },
      };

      vi.mocked(apiClient.post).mockRejectedValueOnce(error);

      await expect(authApi.refresh()).rejects.toMatchObject({
        response: {
          status: 401,
        },
      });
    });

    it('should handle refresh when no token is available', async () => {
      const error = {
        response: {
          status: 401,
          data: {
            success: false,
            code: 401,
            message: 'No token provided',
          },
        },
      };

      vi.mocked(apiClient.post).mockRejectedValueOnce(error);

      await expect(authApi.refresh()).rejects.toMatchObject({
        response: {
          status: 401,
        },
      });
    });
  });

  describe('Integration scenarios', () => {
    it('should handle complete login flow', async () => {
      const loginData: LoginDto = {
        username: 'admin',
        password: 'password123',
      };

      const loginResponse: AxiosResponse = {
        data: {
          success: true,
          code: 200,
          message: 'Login successful',
          data: {
            token: 'test_token',
            user: {
              id: 1,
              username: 'admin',
              email: 'admin@gamelink.com',
              role: 'admin',
            },
          },
        },
        status: 200,
        statusText: 'OK',
        headers: {},
        config: {} as any,
      };

      vi.mocked(apiClient.post).mockResolvedValueOnce(loginResponse);

      const result = await authApi.login(loginData);

      expect(result.data.data.token).toBe('test_token');
      expect(result.data.data.user.role).toBe('admin');
    });

    it('should handle multiple failed login attempts', async () => {
      const loginData: LoginDto = {
        username: 'admin',
        password: 'wrong_password',
      };

      const error = {
        response: {
          status: 401,
          data: {
            success: false,
            code: 401,
            message: 'Invalid credentials',
          },
        },
      };

      // Mock both calls to reject
      vi.mocked(apiClient.post).mockRejectedValueOnce(error);
      vi.mocked(apiClient.post).mockRejectedValueOnce(error);

      const firstAttempt = authApi.login(loginData);
      const secondAttempt = authApi.login(loginData);

      await expect(firstAttempt).rejects.toMatchObject({
        response: { status: 401 },
      });
      await expect(secondAttempt).rejects.toMatchObject({
        response: { status: 401 },
      });

      expect(apiClient.post).toHaveBeenCalledTimes(2);
    });
  });

  describe('Type safety', () => {
    it('should enforce LoginDto types', () => {
      const validLogin: LoginDto = {
        username: 'admin',
        password: 'password',
      };

      expect(validLogin.username).toBeDefined();
      // Password is optional but we provided it
      expect(validLogin.password).toBeDefined();
    });

    it('should enforce RegisterDto types', () => {
      const validRegister: RegisterDto = {
        name: 'Test User',
        email: 'test@example.com',
        password: 'password123',
      };

      expect(validRegister.name).toBeDefined();
      expect(validRegister.email).toBeDefined();
      expect(validRegister.password).toBeDefined();
      // Optional fields
      expect(validRegister.phone).toBeUndefined();
      expect(validRegister.avatarUrl).toBeUndefined();
    });

    it('should enforce LoginResponse types', () => {
      const loginResponse: LoginResponse = {
        token: 'test_token',
        user: {
          id: 1,
          username: 'admin',
          email: 'admin@example.com',
          role: 'admin',
        },
      };

      expect(loginResponse.token).toBeDefined();
      expect(loginResponse.user.id).toBe(1);
      expect(loginResponse.user.role).toBe('admin');
    });
  });
});
