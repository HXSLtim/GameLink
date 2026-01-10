/**
 * User Domain Service
 * Encapsulates all user-related business logic
 *
 * @module services/domain/userService
 */

import {
  BaseService,
  type ServiceDependencies,
} from './base';
import {
  ServiceErrorCodes,
  type ServiceResult,
  type BatchResult,
  ServiceResultHelper,
} from '../utils';
import type {
  User,
  CreateUserDto,
  UpdateUserDto,
  UserQueryParams,
} from '@/api/admin';

/**
 * User validation result
 */
export interface UserValidationResult {
  valid: boolean;
  errors: Array<{
    field: string;
    message: string;
  }>;
}

/**
 * Password validation result
 */
export interface PasswordValidationResult {
  valid: boolean;
  errors: string[];
}

/**
 * User export data format
 */
export interface UserExportData {
  headers: string[];
  rows: string[][];
}

/**
 * User Service Interface
 */
export interface IUserService {
  // CRUD Operations
  getUsers(params?: UserQueryParams): Promise<ServiceResult<User[]>>;
  getUserById(id: number): Promise<ServiceResult<User>>;
  createUser(data: CreateUserDto): Promise<ServiceResult<User>>;
  updateUser(id: number, data: UpdateUserDto): Promise<ServiceResult<User>>;
  deleteUser(id: number): Promise<ServiceResult<void>>;

  // Status & Role
  updateUserStatus(id: number, status: string): Promise<ServiceResult<User>>;
  updateUserRole(id: number, role: string): Promise<ServiceResult<User>>;

  // Batch Operations
  batchUpdateStatus(userIds: number[], status: string): Promise<BatchResult<void>>;
  batchUpdateRole(userIds: number[], role: string): Promise<BatchResult<void>>;
  batchDelete(userIds: number[]): Promise<BatchResult<void>>;

  // Validation
  validateUserData(data: Partial<CreateUserDto>): UserValidationResult;
  validateEmail(email: string): boolean;
  validatePhone(phone: string): boolean;
  validatePassword(password: string): PasswordValidationResult;

  // Export
  exportUsers(users: User[]): UserExportData;
}

/**
 * RFC 5322 compliant email regex pattern
 * Simplified version that covers most common email formats
 */
const EMAIL_REGEX = /^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$/;

/**
 * Chinese mobile phone number regex pattern
 * Supports 13x, 14x, 15x, 16x, 17x, 18x, 19x prefixes
 */
const CHINESE_PHONE_REGEX = /^1[3-9]\d{9}$/;

/**
 * Password requirements
 */
const PASSWORD_MIN_LENGTH = 8;
const PASSWORD_REQUIREMENTS = {
  minLength: PASSWORD_MIN_LENGTH,
  requireUppercase: true,
  requireLowercase: true,
  requireNumber: true,
  requireSpecial: true,
};

/**
 * Valid user statuses
 */
const VALID_STATUSES = ['active', 'banned', 'suspended'];

/**
 * Valid user roles
 */
const VALID_ROLES = ['user', 'player', 'admin'];

/**
 * User Service Implementation
 *
 * Provides all user-related business logic including:
 * - CRUD operations
 * - Data validation
 * - Batch operations
 * - Data export
 */
export class UserService extends BaseService implements IUserService {
  constructor(deps: ServiceDependencies = {}) {
    super(deps);
  }

  // ==================== CRUD Operations ====================

  /**
   * Get users with optional filtering
   */
  async getUsers(params?: UserQueryParams): Promise<ServiceResult<User[]>> {
    return this.withLogging('getUsers', { params: params ?? {} }, async () => {
      return this.wrapAsync(async () => {
        const response = await this.api.getUsers(params);
        return response.data.data;
      }, 'Failed to fetch users');
    });
  }

  /**
   * Get a single user by ID
   */
  async getUserById(id: number): Promise<ServiceResult<User>> {
    return this.withLogging('getUserById', { id }, async () => {
      return this.wrapAsync(async () => {
        const response = await this.api.getUser(id);
        return response.data.data;
      }, `Failed to fetch user ${id}`);
    });
  }

  /**
   * Create a new user
   */
  async createUser(data: CreateUserDto): Promise<ServiceResult<User>> {
    // Validate user data first
    const validation = this.validateUserData(data);
    if (!validation.valid) {
      return ServiceResultHelper.failure({
        code: ServiceErrorCodes.VALIDATION_ERROR,
        message: 'User data validation failed',
        details: { errors: validation.errors },
      });
    }

    return this.withLogging('createUser', { email: data.email }, async () => {
      return this.wrapAsync(async () => {
        const response = await this.api.createUser(data);
        return response.data.data;
      }, 'Failed to create user');
    });
  }

  /**
   * Update an existing user
   */
  async updateUser(id: number, data: UpdateUserDto): Promise<ServiceResult<User>> {
    // Validate user data (password is optional for updates)
    const validation = this.validateUserData({
      ...data,
      password: data.password || 'ValidPass1!', // Skip password validation if not provided
    });
    if (!validation.valid) {
      return ServiceResultHelper.failure({
        code: ServiceErrorCodes.VALIDATION_ERROR,
        message: 'User data validation failed',
        details: { errors: validation.errors },
      });
    }

    return this.withLogging('updateUser', { id, email: data.email }, async () => {
      return this.wrapAsync(async () => {
        const response = await this.api.updateUser(id, data);
        return response.data.data;
      }, `Failed to update user ${id}`);
    });
  }

  /**
   * Delete a user
   */
  async deleteUser(id: number): Promise<ServiceResult<void>> {
    return this.withLogging('deleteUser', { id }, async () => {
      return this.wrapAsync(async () => {
        await this.api.deleteUser(id);
      }, `Failed to delete user ${id}`);
    });
  }

  // ==================== Status & Role ====================

  /**
   * Update user status
   */
  async updateUserStatus(id: number, status: string): Promise<ServiceResult<User>> {
    if (!VALID_STATUSES.includes(status)) {
      return ServiceResultHelper.failure({
        code: ServiceErrorCodes.USER_INVALID_STATUS,
        message: `Invalid status: ${status}. Valid statuses are: ${VALID_STATUSES.join(', ')}`,
      });
    }

    return this.withLogging('updateUserStatus', { id, status }, async () => {
      return this.wrapAsync(async () => {
        const response = await this.api.updateUserStatus(id, status);
        return response.data.data;
      }, `Failed to update user ${id} status`);
    });
  }

  /**
   * Update user role
   */
  async updateUserRole(id: number, role: string): Promise<ServiceResult<User>> {
    if (!VALID_ROLES.includes(role)) {
      return ServiceResultHelper.failure({
        code: ServiceErrorCodes.VALIDATION_ERROR,
        message: `Invalid role: ${role}. Valid roles are: ${VALID_ROLES.join(', ')}`,
      });
    }

    return this.withLogging('updateUserRole', { id, role }, async () => {
      return this.wrapAsync(async () => {
        const response = await this.api.updateUserRole(id, role);
        return response.data.data;
      }, `Failed to update user ${id} role`);
    });
  }

  // ==================== Batch Operations ====================

  /**
   * Batch update user status
   */
  async batchUpdateStatus(
    userIds: number[],
    status: string
  ): Promise<BatchResult<void>> {
    if (!VALID_STATUSES.includes(status)) {
      return {
        success: false,
        total: userIds.length,
        succeeded: 0,
        failed: userIds.length,
        results: userIds.map((_, index) => ({
          index,
          success: false,
          error: {
            code: ServiceErrorCodes.USER_INVALID_STATUS,
            message: `Invalid status: ${status}`,
          },
        })),
      };
    }

    if (userIds.length === 0) {
      return ServiceResultHelper.emptyBatch(0);
    }

    return this.withLogging(
      'batchUpdateStatus',
      { userIds, status },
      async () => {
        return this.executeBatch(
          userIds,
          async (userId) => {
            await this.api.updateUserStatus(userId, status);
          },
          'batchUpdateStatus'
        );
      }
    );
  }

  /**
   * Batch update user role
   */
  async batchUpdateRole(
    userIds: number[],
    role: string
  ): Promise<BatchResult<void>> {
    if (!VALID_ROLES.includes(role)) {
      return {
        success: false,
        total: userIds.length,
        succeeded: 0,
        failed: userIds.length,
        results: userIds.map((_, index) => ({
          index,
          success: false,
          error: {
            code: ServiceErrorCodes.VALIDATION_ERROR,
            message: `Invalid role: ${role}`,
          },
        })),
      };
    }

    if (userIds.length === 0) {
      return ServiceResultHelper.emptyBatch(0);
    }

    return this.withLogging(
      'batchUpdateRole',
      { userIds, role },
      async () => {
        return this.executeBatch(
          userIds,
          async (userId) => {
            await this.api.updateUserRole(userId, role);
          },
          'batchUpdateRole'
        );
      }
    );
  }

  /**
   * Batch delete users
   */
  async batchDelete(userIds: number[]): Promise<BatchResult<void>> {
    if (userIds.length === 0) {
      return ServiceResultHelper.emptyBatch(0);
    }

    return this.withLogging('batchDelete', { userIds }, async () => {
      return this.executeBatch(
        userIds,
        async (userId) => {
          await this.api.deleteUser(userId);
        },
        'batchDelete'
      );
    });
  }

  // ==================== Validation ====================

  /**
   * Validate user data comprehensively
   * Returns ALL validation errors, not just the first one
   */
  validateUserData(data: Partial<CreateUserDto>): UserValidationResult {
    const errors: Array<{ field: string; message: string }> = [];

    // Validate email
    if (data.email !== undefined) {
      if (!data.email || data.email.trim() === '') {
        errors.push({ field: 'email', message: 'Email is required' });
      } else if (!this.validateEmail(data.email)) {
        errors.push({ field: 'email', message: 'Invalid email format' });
      }
    }

    // Validate phone
    if (data.phone !== undefined) {
      if (!data.phone || data.phone.trim() === '') {
        errors.push({ field: 'phone', message: 'Phone is required' });
      } else if (!this.validatePhone(data.phone)) {
        errors.push({ field: 'phone', message: 'Invalid phone format. Must be a valid Chinese mobile number' });
      }
    }

    // Validate password
    if (data.password !== undefined) {
      const passwordValidation = this.validatePassword(data.password);
      if (!passwordValidation.valid) {
        passwordValidation.errors.forEach((error) => {
          errors.push({ field: 'password', message: error });
        });
      }
    }

    // Validate name
    if (data.name !== undefined) {
      if (!data.name || data.name.trim() === '') {
        errors.push({ field: 'name', message: 'Name is required' });
      } else if (data.name.length < 2) {
        errors.push({ field: 'name', message: 'Name must be at least 2 characters' });
      } else if (data.name.length > 50) {
        errors.push({ field: 'name', message: 'Name must be at most 50 characters' });
      }
    }

    // Validate role
    if (data.role !== undefined) {
      if (!VALID_ROLES.includes(data.role)) {
        errors.push({ field: 'role', message: `Invalid role. Valid roles are: ${VALID_ROLES.join(', ')}` });
      }
    }

    // Validate status
    if (data.status !== undefined) {
      if (!VALID_STATUSES.includes(data.status)) {
        errors.push({ field: 'status', message: `Invalid status. Valid statuses are: ${VALID_STATUSES.join(', ')}` });
      }
    }

    return {
      valid: errors.length === 0,
      errors,
    };
  }

  /**
   * Validate email format (RFC 5322 compliant)
   */
  validateEmail(email: string): boolean {
    if (!email || typeof email !== 'string') {
      return false;
    }
    return EMAIL_REGEX.test(email.trim());
  }

  /**
   * Validate Chinese mobile phone format
   */
  validatePhone(phone: string): boolean {
    if (!phone || typeof phone !== 'string') {
      return false;
    }
    return CHINESE_PHONE_REGEX.test(phone.trim());
  }

  /**
   * Validate password strength
   * Returns all validation errors
   */
  validatePassword(password: string): PasswordValidationResult {
    const errors: string[] = [];

    if (!password || typeof password !== 'string') {
      return { valid: false, errors: ['Password is required'] };
    }

    if (password.length < PASSWORD_REQUIREMENTS.minLength) {
      errors.push(`Password must be at least ${PASSWORD_REQUIREMENTS.minLength} characters`);
    }

    if (PASSWORD_REQUIREMENTS.requireUppercase && !/[A-Z]/.test(password)) {
      errors.push('Password must contain at least one uppercase letter');
    }

    if (PASSWORD_REQUIREMENTS.requireLowercase && !/[a-z]/.test(password)) {
      errors.push('Password must contain at least one lowercase letter');
    }

    if (PASSWORD_REQUIREMENTS.requireNumber && !/\d/.test(password)) {
      errors.push('Password must contain at least one number');
    }

    if (PASSWORD_REQUIREMENTS.requireSpecial && !/[!@#$%^&*()_+\-=[\]{};':"\\|,.<>/?]/.test(password)) {
      errors.push('Password must contain at least one special character');
    }

    return {
      valid: errors.length === 0,
      errors,
    };
  }

  // ==================== Export ====================

  /**
   * Export users to a format suitable for Excel/CSV
   */
  exportUsers(users: User[]): UserExportData {
    const headers = [
      'ID',
      '姓名',
      '邮箱',
      '手机号',
      '角色',
      '状态',
      '最后登录时间',
      '创建时间',
    ];

    const rows = users.map((user) => [
      String(user.id),
      user.name,
      user.email,
      user.phone,
      this.formatRole(user.role),
      this.formatStatus(user.status),
      user.lastLoginAt ? this.formatDate(user.lastLoginAt) : '-',
      this.formatDate(user.createdAt),
    ]);

    return { headers, rows };
  }

  // ==================== Private Helpers ====================

  /**
   * Format role for display
   */
  private formatRole(role: string): string {
    const roleMap: Record<string, string> = {
      user: '普通用户',
      player: '陪玩师',
      admin: '管理员',
    };
    return roleMap[role] || role;
  }

  /**
   * Format status for display
   */
  private formatStatus(status: string): string {
    const statusMap: Record<string, string> = {
      active: '正常',
      banned: '封禁',
      suspended: '暂停',
    };
    return statusMap[status] || status;
  }

  /**
   * Format date for export
   */
  private formatDate(dateStr: string): string {
    try {
      const date = new Date(dateStr);
      return date.toLocaleString('zh-CN', {
        year: 'numeric',
        month: '2-digit',
        day: '2-digit',
        hour: '2-digit',
        minute: '2-digit',
      });
    } catch {
      return dateStr;
    }
  }
}

/**
 * Default UserService instance
 */
export const userService = new UserService();
