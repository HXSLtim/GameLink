/**
 * Property-Based Tests for UserService
 *
 * Tests user data validation, batch operations, and export functionality
 * using property-based testing with fast-check.
 */

import { describe, it, expect, vi, beforeEach } from 'vitest';
import * as fc from 'fast-check';
import { UserService } from './userService';
import type { CreateUserDto, User } from '../../api/admin';

// Mock the admin API
vi.mock('../../api/admin', () => ({
  adminApi: {
    getUsers: vi.fn(),
    getUser: vi.fn(),
    createUser: vi.fn(),
    updateUser: vi.fn(),
    deleteUser: vi.fn(),
    updateUserStatus: vi.fn(),
    updateUserRole: vi.fn(),
  },
}));

describe('UserService - Property Tests', () => {
  let userService: UserService;

  beforeEach(() => {
    userService = new UserService();
    vi.clearAllMocks();
  });

  /**
   * **Feature: admin-phase3-improvements, Property 4: User Data Validation Completeness**
   * **Validates: Requirements 2.2**
   *
   * For any user data input, the validation function SHALL check email format,
   * phone format, and password strength, returning all validation errors
   * (not just the first one).
   */
  describe('Property 4: User Data Validation Completeness', () => {
    // Arbitrary for valid emails (RFC 5322 simplified)
    const validEmailArb = fc
      .tuple(
        fc.stringMatching(/^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+$/),
        fc.stringMatching(/^[a-zA-Z0-9][a-zA-Z0-9-]*[a-zA-Z0-9]$/),
        fc.stringMatching(/^[a-zA-Z]{2,}$/)
      )
      .map(([local, domain, tld]) => `${local}@${domain}.${tld}`)
      .filter((email) => email.length > 5 && email.length < 100);

    // Arbitrary for invalid emails
    const invalidEmailArb = fc.oneof(
      fc.constant(''),
      fc.constant('invalid'),
      fc.constant('@nodomain.com'),
      fc.constant('spaces in@email.com'),
      fc.constant('double@@at.com'),
      fc.constant('missing.at.sign'),
      fc.stringMatching(/^[^@]+$/).filter((s) => s.length > 0) // No @ symbol
    );

    // Arbitrary for valid Chinese phone numbers
    const validPhoneArb = fc
      .tuple(fc.constantFrom('13', '14', '15', '16', '17', '18', '19'), fc.stringMatching(/^\d{9}$/))
      .map(([prefix, suffix]) => `${prefix}${suffix}`);

    // Arbitrary for invalid phone numbers
    const invalidPhoneArb = fc.oneof(
      fc.constant(''),
      fc.constant('12345678901'), // Invalid prefix (12x)
      fc.constant('1381234567'), // Too short
      fc.constant('138123456789'), // Too long
      fc.constant('abc12345678'), // Contains letters
      fc.stringMatching(/^[02-9]\d{10}$/) // Doesn't start with 1
    );

    // Arbitrary for valid passwords
    const validPasswordArb = fc
      .tuple(
        fc.stringMatching(/^[A-Z]$/),
        fc.stringMatching(/^[a-z]$/),
        fc.stringMatching(/^\d$/),
        fc.constantFrom('!', '@', '#', '$', '%', '^', '&', '*'),
        fc.stringMatching(/^[a-zA-Z0-9]{4,}$/)
      )
      .map(([upper, lower, num, special, rest]) => `${upper}${lower}${num}${special}${rest}`);

    // Arbitrary for weak passwords (missing requirements)
    const weakPasswordArb = fc.oneof(
      fc.constant('short'), // Too short
      fc.constant('alllowercase1!'), // No uppercase
      fc.constant('ALLUPPERCASE1!'), // No lowercase
      fc.constant('NoNumbers!!'), // No numbers
      fc.constant('NoSpecial123'), // No special chars
      fc.constant('') // Empty
    );

    it('should return all validation errors for completely invalid data', () => {
      fc.assert(
        fc.property(invalidEmailArb, invalidPhoneArb, weakPasswordArb, (email, phone, password) => {
          const data: Partial<CreateUserDto> = {
            email,
            phone,
            password,
            name: '', // Invalid name
            role: 'invalid' as 'user',
            status: 'invalid' as 'active',
          };

          const result = userService.validateUserData(data);

          // Should not be valid
          expect(result.valid).toBe(false);

          // Should have multiple errors (not just the first one)
          expect(result.errors.length).toBeGreaterThan(1);

          // Should have errors for different fields
          const errorFields = result.errors.map((e) => e.field);

          // At minimum, should check email, phone, and password
          // (name, role, status may also have errors)
          expect(errorFields.length).toBeGreaterThanOrEqual(3);
        }),
        { numRuns: 100 }
      );
    });

    it('should validate email format correctly', () => {
      fc.assert(
        fc.property(validEmailArb, (email) => {
          const isValid = userService.validateEmail(email);
          expect(isValid).toBe(true);
        }),
        { numRuns: 100 }
      );

      fc.assert(
        fc.property(invalidEmailArb, (email) => {
          const isValid = userService.validateEmail(email);
          expect(isValid).toBe(false);
        }),
        { numRuns: 100 }
      );
    });

    it('should validate Chinese phone format correctly', () => {
      fc.assert(
        fc.property(validPhoneArb, (phone) => {
          const isValid = userService.validatePhone(phone);
          expect(isValid).toBe(true);
        }),
        { numRuns: 100 }
      );

      fc.assert(
        fc.property(invalidPhoneArb, (phone) => {
          const isValid = userService.validatePhone(phone);
          expect(isValid).toBe(false);
        }),
        { numRuns: 100 }
      );
    });

    it('should validate password strength correctly', () => {
      fc.assert(
        fc.property(validPasswordArb, (password) => {
          const result = userService.validatePassword(password);
          expect(result.valid).toBe(true);
          expect(result.errors).toHaveLength(0);
        }),
        { numRuns: 100 }
      );

      fc.assert(
        fc.property(weakPasswordArb, (password) => {
          const result = userService.validatePassword(password);
          expect(result.valid).toBe(false);
          expect(result.errors.length).toBeGreaterThan(0);
        }),
        { numRuns: 100 }
      );
    });

    it('should return all password errors, not just the first one', () => {
      // Password with multiple issues
      const result = userService.validatePassword('abc');

      expect(result.valid).toBe(false);
      // Should have errors for: length, uppercase, number, special char
      expect(result.errors.length).toBeGreaterThanOrEqual(4);
    });

    it('should validate all fields and collect all errors', () => {
      const data: Partial<CreateUserDto> = {
        email: 'invalid-email',
        phone: '12345',
        password: 'weak',
        name: 'a', // Too short
        role: 'invalid-role' as 'user',
        status: 'invalid-status' as 'active',
      };

      const result = userService.validateUserData(data);

      expect(result.valid).toBe(false);

      // Should have errors for all invalid fields
      const errorFields = new Set(result.errors.map((e) => e.field));
      expect(errorFields.has('email')).toBe(true);
      expect(errorFields.has('phone')).toBe(true);
      expect(errorFields.has('password')).toBe(true);
      expect(errorFields.has('name')).toBe(true);
      expect(errorFields.has('role')).toBe(true);
      expect(errorFields.has('status')).toBe(true);
    });

    it('should return valid for correct data', () => {
      // Generate valid names that are not just whitespace
      // Name must have at least 2 non-whitespace characters
      const alphanumericChars = 'abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789';
      const validNameArb = fc
        .array(fc.constantFrom(...alphanumericChars.split('')), { minLength: 2, maxLength: 50 })
        .map((chars) => chars.join(''));

      fc.assert(
        fc.property(
          validEmailArb,
          validPhoneArb,
          validPasswordArb,
          validNameArb,
          fc.constantFrom('user', 'player', 'admin'),
          fc.constantFrom('active', 'banned', 'suspended'),
          (email, phone, password, name, role, status) => {
            const data: Partial<CreateUserDto> = {
              email,
              phone,
              password,
              name,
              role: role as 'user' | 'player' | 'admin',
              status: status as 'active' | 'banned' | 'suspended',
            };

            const result = userService.validateUserData(data);
            expect(result.valid).toBe(true);
            expect(result.errors).toHaveLength(0);
          }
        ),
        { numRuns: 100 }
      );
    });
  });

  /**
   * **Feature: admin-phase3-improvements, Property 5: Batch Operation Result Completeness**
   * **Validates: Requirements 2.3**
   *
   * For any batch operation on users, the result SHALL contain an entry for
   * each input item indicating success or failure with details.
   */
  describe('Property 5: Batch Operation Result Completeness', () => {
    // Arbitrary for user IDs
    const userIdsArb = fc.array(fc.integer({ min: 1, max: 10000 }), {
      minLength: 1,
      maxLength: 20,
    });

    // Arbitrary for valid statuses
    const validStatusArb = fc.constantFrom('active', 'banned', 'suspended');

    // Arbitrary for valid roles
    const validRoleArb = fc.constantFrom('user', 'player', 'admin');

    it('batchUpdateStatus result should have entry for each input user', async () => {
      const { adminApi } = await import('../../api/admin');

      await fc.assert(
        fc.asyncProperty(userIdsArb, validStatusArb, async (userIds, status) => {
          // Mock API to succeed for all calls
          vi.mocked(adminApi.updateUserStatus).mockResolvedValue({
            data: { success: true, data: {} as User },
          } as never);

          const result = await userService.batchUpdateStatus(userIds, status);

          // Result should have entry for each input
          expect(result.total).toBe(userIds.length);
          expect(result.results.length).toBe(userIds.length);

          // Each result should have an index
          const indices = result.results.map((r) => r.index);
          for (let i = 0; i < userIds.length; i++) {
            expect(indices).toContain(i);
          }

          // succeeded + failed should equal total
          expect(result.succeeded + result.failed).toBe(result.total);

          // success flag should reflect whether all succeeded
          expect(result.success).toBe(result.failed === 0);
        }),
        { numRuns: 50 }
      );
    });

    it('batchUpdateRole result should have entry for each input user', async () => {
      const { adminApi } = await import('../../api/admin');

      await fc.assert(
        fc.asyncProperty(userIdsArb, validRoleArb, async (userIds, role) => {
          // Mock API to succeed for all calls
          vi.mocked(adminApi.updateUserRole).mockResolvedValue({
            data: { success: true, data: {} as User },
          } as never);

          const result = await userService.batchUpdateRole(userIds, role);

          // Result should have entry for each input
          expect(result.total).toBe(userIds.length);
          expect(result.results.length).toBe(userIds.length);

          // succeeded + failed should equal total
          expect(result.succeeded + result.failed).toBe(result.total);
        }),
        { numRuns: 50 }
      );
    });

    it('batchDelete result should have entry for each input user', async () => {
      const { adminApi } = await import('../../api/admin');

      await fc.assert(
        fc.asyncProperty(userIdsArb, async (userIds) => {
          // Mock API to succeed for all calls
          vi.mocked(adminApi.deleteUser).mockResolvedValue({
            data: { success: true },
          } as never);

          const result = await userService.batchDelete(userIds);

          // Result should have entry for each input
          expect(result.total).toBe(userIds.length);
          expect(result.results.length).toBe(userIds.length);

          // succeeded + failed should equal total
          expect(result.succeeded + result.failed).toBe(result.total);
        }),
        { numRuns: 50 }
      );
    });

    it('batch operations should report partial failures correctly', async () => {
      const { adminApi } = await import('../../api/admin');

      await fc.assert(
        fc.asyncProperty(
          // Generate unique user IDs to avoid duplicate ID issues
          fc.array(fc.integer({ min: 1, max: 1000 }), { minLength: 3, maxLength: 10 })
            .map(ids => [...new Set(ids)]) // Remove duplicates
            .filter(ids => ids.length >= 3), // Ensure at least 3 unique IDs
          fc.integer({ min: 0 }),
          async (userIds, failIndex) => {
            const actualFailIndex = failIndex % userIds.length;

            // Mock API to fail for one specific index
            vi.mocked(adminApi.updateUserStatus).mockImplementation(async (id) => {
              if (id === userIds[actualFailIndex]) {
                throw new Error('API Error');
              }
              return { data: { success: true, data: {} as User } } as never;
            });

            const result = await userService.batchUpdateStatus(userIds, 'active');

            // Should have exactly one failure
            expect(result.failed).toBe(1);
            expect(result.succeeded).toBe(userIds.length - 1);
            expect(result.success).toBe(false);

            // The failed item should have an error
            const failedResult = result.results.find((r) => r.index === actualFailIndex);
            expect(failedResult).toBeDefined();
            expect(failedResult?.success).toBe(false);
            expect(failedResult?.error).toBeDefined();
            expect(failedResult?.error?.code).toBeDefined();
            expect(failedResult?.error?.message).toBeDefined();
          }
        ),
        { numRuns: 30 }
      );
    });

    it('batch operations with invalid status should fail all items', async () => {
      const invalidStatus = 'invalid_status';
      const userIds = [1, 2, 3];

      const result = await userService.batchUpdateStatus(userIds, invalidStatus);

      // All should fail due to invalid status
      expect(result.success).toBe(false);
      expect(result.failed).toBe(userIds.length);
      expect(result.succeeded).toBe(0);

      // Each result should have an error
      for (const itemResult of result.results) {
        expect(itemResult.success).toBe(false);
        expect(itemResult.error).toBeDefined();
      }
    });

    it('batch operations with empty array should return empty result', async () => {
      const result = await userService.batchUpdateStatus([], 'active');

      expect(result.total).toBe(0);
      expect(result.succeeded).toBe(0);
      expect(result.failed).toBe(0);
      expect(result.results).toHaveLength(0);
      expect(result.success).toBe(true);
    });
  });

  /**
   * **Feature: admin-phase3-improvements, Property 6: Export Data Format Consistency**
   * **Validates: Requirements 2.5**
   *
   * For any user data export, the exported data SHALL contain headers matching
   * the expected columns and rows containing properly formatted values.
   */
  describe('Property 6: Export Data Format Consistency', () => {
    // Arbitrary for user data
    const userArb = fc.record({
      id: fc.integer({ min: 1, max: 100000 }),
      name: fc.string({ minLength: 1, maxLength: 50 }),
      email: fc.emailAddress(),
      phone: fc
        .tuple(fc.constantFrom('13', '14', '15', '16', '17', '18', '19'), fc.stringMatching(/^\d{9}$/))
        .map(([prefix, suffix]) => `${prefix}${suffix}`),
      role: fc.constantFrom('user', 'player', 'admin') as fc.Arbitrary<'user' | 'player' | 'admin'>,
      status: fc.constantFrom('active', 'banned', 'suspended') as fc.Arbitrary<
        'active' | 'banned' | 'suspended'
      >,
      lastLoginAt: fc.option(
        fc.integer({ min: 1577836800000, max: 1924905600000 }).map((ts) => new Date(ts).toISOString()),
        { nil: undefined }
      ),
      createdAt: fc.integer({ min: 1577836800000, max: 1924905600000 }).map((ts) => new Date(ts).toISOString()),
    });

    const usersArb = fc.array(userArb, { minLength: 0, maxLength: 50 });

    it('exported data should have correct headers', () => {
      fc.assert(
        fc.property(usersArb, (users) => {
          const result = userService.exportUsers(users as User[]);

          // Should have expected headers
          expect(result.headers).toEqual([
            'ID',
            '姓名',
            '邮箱',
            '手机号',
            '角色',
            '状态',
            '最后登录时间',
            '创建时间',
          ]);
        }),
        { numRuns: 50 }
      );
    });

    it('exported rows should match input users count', () => {
      fc.assert(
        fc.property(usersArb, (users) => {
          const result = userService.exportUsers(users as User[]);

          // Number of rows should match number of users
          expect(result.rows.length).toBe(users.length);
        }),
        { numRuns: 50 }
      );
    });

    it('each row should have same number of columns as headers', () => {
      fc.assert(
        fc.property(usersArb, (users) => {
          const result = userService.exportUsers(users as User[]);

          // Each row should have same number of columns as headers
          for (const row of result.rows) {
            expect(row.length).toBe(result.headers.length);
          }
        }),
        { numRuns: 50 }
      );
    });

    it('exported data should contain user information', () => {
      fc.assert(
        fc.property(fc.array(userArb, { minLength: 1, maxLength: 10 }), (users) => {
          const result = userService.exportUsers(users as User[]);

          // Each row should contain the user's data
          for (let i = 0; i < users.length; i++) {
            const user = users[i];
            const row = result.rows[i];

            // ID should match
            expect(row[0]).toBe(String(user.id));

            // Name should match
            expect(row[1]).toBe(user.name);

            // Email should match
            expect(row[2]).toBe(user.email);

            // Phone should match
            expect(row[3]).toBe(user.phone);

            // Role should be formatted (Chinese)
            expect(['普通用户', '陪玩师', '管理员']).toContain(row[4]);

            // Status should be formatted (Chinese)
            expect(['正常', '封禁', '暂停']).toContain(row[5]);

            // Last login should be formatted or '-'
            if (user.lastLoginAt) {
              expect(row[6]).not.toBe('-');
            } else {
              expect(row[6]).toBe('-');
            }

            // Created at should be formatted
            expect(row[7]).not.toBe('');
          }
        }),
        { numRuns: 50 }
      );
    });

    it('all row values should be strings', () => {
      fc.assert(
        fc.property(usersArb, (users) => {
          const result = userService.exportUsers(users as User[]);

          // All values should be strings (for CSV/Excel compatibility)
          for (const row of result.rows) {
            for (const cell of row) {
              expect(typeof cell).toBe('string');
            }
          }
        }),
        { numRuns: 50 }
      );
    });

    it('empty user list should return headers with no rows', () => {
      const result = userService.exportUsers([]);

      expect(result.headers.length).toBe(8);
      expect(result.rows.length).toBe(0);
    });
  });
});