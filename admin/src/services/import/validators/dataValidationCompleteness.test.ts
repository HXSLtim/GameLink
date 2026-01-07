/**
 * Property-Based Tests for Import Data Validation Completeness
 *
 * **Feature: admin-phase3-improvements, Property 15: Import Data Validation Completeness**
 * **Validates: Requirements 5.3**
 *
 * Tests that the validation checks all fields against business rules and collects
 * ALL validation errors for each row (not stopping at the first error).
 */

import { describe, it, expect } from 'vitest';
import * as fc from 'fast-check';
import { validateUserRow, validateUserData } from './userDataValidator';
import { validatePlayerRow } from './playerDataValidator';
import { validateGameRow } from './gameDataValidator';

describe('Import Data Validation Completeness - Property Tests', () => {
  /**
   * **Feature: admin-phase3-improvements, Property 15: Import Data Validation Completeness**
   * **Validates: Requirements 5.3**
   *
   * For any import data row, the validation SHALL check all fields against business
   * rules and collect ALL validation errors for the row (not stopping at the first error).
   */
  describe('Property 15: Import Data Validation Completeness', () => {
    describe('User Data Validation Completeness', () => {
      it('should collect all validation errors for a row with multiple invalid fields', () => {
        // Generate data with multiple invalid fields
        const invalidDataArb = fc.record({
          name: fc.constantFrom('', ' ', 'a'), // Invalid: empty or too short
          email: fc.constantFrom('', 'invalid', 'no-at-sign'), // Invalid: not email format
          phone: fc.constantFrom('', '123', 'abc'), // Invalid: not Chinese phone format
          role: fc.constantFrom('invalid_role', 'superuser'), // Invalid: not in allowed list
          status: fc.constantFrom('invalid_status', 'deleted'), // Invalid: not in allowed list
        });

        fc.assert(
          fc.property(invalidDataArb, (data) => {
            const result = validateUserRow(data, 2);

            // Should be invalid
            expect(result.valid).toBe(false);

            // Should have multiple errors (at least 3: name, email, phone are required)
            expect(result.errors.length).toBeGreaterThanOrEqual(3);

            // Should have errors for each invalid field
            const errorFields = result.errors.map((e) => e.field);
            expect(errorFields).toContain('name');
            expect(errorFields).toContain('email');
            expect(errorFields).toContain('phone');
          }),
          { numRuns: 20 }
        );
      });

      it('should report all errors even when first field is invalid', () => {
        fc.assert(
          fc.property(
            fc.integer({ min: 2, max: 100 }),
            (rowNumber) => {
              // Data with all fields invalid
              const data = {
                name: '', // Invalid
                email: 'not-an-email', // Invalid
                phone: '123', // Invalid
                role: 'superadmin', // Invalid
                status: 'deleted', // Invalid
              };

              const result = validateUserRow(data, rowNumber);

              // Should collect errors for all invalid fields
              expect(result.errors.length).toBeGreaterThanOrEqual(5);

              // Verify each field has an error
              const errorFields = new Set(result.errors.map((e) => e.field));
              expect(errorFields.has('name')).toBe(true);
              expect(errorFields.has('email')).toBe(true);
              expect(errorFields.has('phone')).toBe(true);
              expect(errorFields.has('role')).toBe(true);
              expect(errorFields.has('status')).toBe(true);
            }
          ),
          { numRuns: 10 }
        );
      });

      it('should return no errors for valid data', () => {
        const validDataArb = fc.record({
          name: fc.string({ minLength: 2, maxLength: 50 }).filter((s) => s.trim().length >= 2),
          email: fc.integer({ min: 1, max: 9999 }).map((n) => `user${n}@example.com`),
          phone: fc.integer({ min: 0, max: 99999999 }).map((n) => `138${String(n).padStart(8, '0')}`),
          role: fc.constantFrom('user', 'player', 'admin'),
          status: fc.constantFrom('active', 'banned', 'suspended'),
        });

        fc.assert(
          fc.property(validDataArb, (data) => {
            const result = validateUserRow(data, 2);

            expect(result.valid).toBe(true);
            expect(result.errors).toHaveLength(0);
          }),
          { numRuns: 20 }
        );
      });
    });

    describe('Player Data Validation Completeness', () => {
      it('should collect all validation errors for a row with multiple invalid fields', () => {
        const invalidDataArb = fc.record({
          userEmail: fc.constantFrom('', 'invalid', 'no-at-sign'), // Invalid
          nickname: fc.string({ minLength: 51, maxLength: 60 }), // Too long
          bio: fc.string({ minLength: 501, maxLength: 600 }), // Too long
          hourlyRate: fc.constantFrom(-1, -100, 'abc'), // Invalid: negative or not a number
          skillTags: fc.constant('a'.repeat(25) + ',' + 'b'.repeat(25)), // Tags too long
        });

        fc.assert(
          fc.property(invalidDataArb, (data) => {
            const result = validatePlayerRow(data, 2);

            // Should be invalid
            expect(result.valid).toBe(false);

            // Should have multiple errors
            expect(result.errors.length).toBeGreaterThanOrEqual(1);

            // Should have error for userEmail at minimum
            const errorFields = result.errors.map((e) => e.field);
            expect(errorFields).toContain('userEmail');
          }),
          { numRuns: 20 }
        );
      });

      it('should validate all fields even when userEmail is invalid', () => {
        const data = {
          userEmail: 'invalid-email',
          nickname: 'a'.repeat(60), // Too long
          bio: 'a'.repeat(600), // Too long
          hourlyRate: -50, // Negative
          skillTags: 'a'.repeat(25), // Tag too long
        };

        const result = validatePlayerRow(data, 2);

        // Should have errors for multiple fields
        expect(result.errors.length).toBeGreaterThanOrEqual(4);

        const errorFields = new Set(result.errors.map((e) => e.field));
        expect(errorFields.has('userEmail')).toBe(true);
        expect(errorFields.has('nickname')).toBe(true);
        expect(errorFields.has('bio')).toBe(true);
        expect(errorFields.has('hourlyRate')).toBe(true);
      });
    });

    describe('Game Data Validation Completeness', () => {
      it('should collect all validation errors for a row with multiple invalid fields', () => {
        const invalidDataArb = fc.record({
          key: fc.constantFrom('', 'invalid key!', 'has spaces'), // Invalid: empty or invalid chars
          name: fc.constantFrom('', ' '), // Invalid: empty
          category: fc.string({ minLength: 51, maxLength: 60 }), // Too long
          description: fc.string({ minLength: 501, maxLength: 600 }), // Too long
          isActive: fc.constantFrom('maybe', 'unknown'), // Invalid boolean
        });

        fc.assert(
          fc.property(invalidDataArb, (data) => {
            const result = validateGameRow(data, 2);

            // Should be invalid
            expect(result.valid).toBe(false);

            // Should have multiple errors
            expect(result.errors.length).toBeGreaterThanOrEqual(2);

            // Should have errors for key and name at minimum
            const errorFields = result.errors.map((e) => e.field);
            expect(errorFields).toContain('key');
            expect(errorFields).toContain('name');
          }),
          { numRuns: 20 }
        );
      });

      it('should validate all fields even when key is invalid', () => {
        const data = {
          key: 'invalid key!', // Invalid chars
          name: '', // Empty
          category: 'a'.repeat(60), // Too long
          description: 'a'.repeat(600), // Too long
          isActive: 'maybe', // Invalid boolean
        };

        const result = validateGameRow(data, 2);

        // Should have errors for multiple fields
        expect(result.errors.length).toBeGreaterThanOrEqual(4);

        const errorFields = new Set(result.errors.map((e) => e.field));
        expect(errorFields.has('key')).toBe(true);
        expect(errorFields.has('name')).toBe(true);
        expect(errorFields.has('category')).toBe(true);
        expect(errorFields.has('description')).toBe(true);
      });
    });

    describe('Batch Validation Completeness', () => {
      it('should validate all rows and collect all errors', () => {
        // Create rows with varying numbers of errors
        const rows = [
          { rowNumber: 2, data: { name: '', email: '', phone: '' } }, // 3 errors
          { rowNumber: 3, data: { name: 'Valid', email: 'invalid', phone: '13800138000' } }, // 1 error
          { rowNumber: 4, data: { name: 'Valid', email: 'valid@example.com', phone: '13800138000' } }, // 0 errors
          { rowNumber: 5, data: { name: '', email: 'invalid', phone: 'abc' } }, // 3 errors
        ];

        const result = validateUserData(rows, { checkInternalDuplicates: false });

        // Should have 3 invalid rows
        expect(result.invalidRows.length).toBe(3);

        // Should have 1 valid row
        expect(result.validRows.length).toBe(1);

        // Each invalid row should have all its errors collected
        const row2 = result.invalidRows.find((r) => r.rowNumber === 2);
        expect(row2?.errors.length).toBeGreaterThanOrEqual(3);

        const row3 = result.invalidRows.find((r) => r.rowNumber === 3);
        expect(row3?.errors.length).toBeGreaterThanOrEqual(1);

        const row5 = result.invalidRows.find((r) => r.rowNumber === 5);
        expect(row5?.errors.length).toBeGreaterThanOrEqual(3);
      });

      it('should preserve row numbers in error results', () => {
        fc.assert(
          fc.property(
            fc.array(fc.integer({ min: 2, max: 1000 }), { minLength: 1, maxLength: 10 }),
            (rowNumbers) => {
              const uniqueRowNumbers = [...new Set(rowNumbers)];
              const rows = uniqueRowNumbers.map((rowNumber) => ({
                rowNumber,
                data: { name: '', email: '', phone: '' }, // All invalid
              }));

              const result = validateUserData(rows, { checkInternalDuplicates: false });

              // All rows should be invalid
              expect(result.invalidRows.length).toBe(uniqueRowNumbers.length);

              // Each result should have the correct row number
              for (const row of result.invalidRows) {
                expect(uniqueRowNumbers).toContain(row.rowNumber);
              }
            }
          ),
          { numRuns: 20 }
        );
      });
    });
  });
});

describe('Import Data Validation Completeness - Unit Tests', () => {
  describe('User validation collects all errors', () => {
    it('should collect name, email, and phone errors together', () => {
      const result = validateUserRow(
        {
          name: '',
          email: 'not-email',
          phone: '123',
        },
        2
      );

      expect(result.valid).toBe(false);
      expect(result.errors.length).toBeGreaterThanOrEqual(3);

      const fields = result.errors.map((e) => e.field);
      expect(fields).toContain('name');
      expect(fields).toContain('email');
      expect(fields).toContain('phone');
    });

    it('should include role and status errors when invalid', () => {
      const result = validateUserRow(
        {
          name: 'Valid Name',
          email: 'valid@example.com',
          phone: '13800138000',
          role: 'superuser',
          status: 'deleted',
        },
        2
      );

      expect(result.valid).toBe(false);
      expect(result.errors.length).toBe(2);

      const fields = result.errors.map((e) => e.field);
      expect(fields).toContain('role');
      expect(fields).toContain('status');
    });
  });

  describe('Player validation collects all errors', () => {
    it('should collect all field errors together', () => {
      const result = validatePlayerRow(
        {
          userEmail: 'invalid',
          nickname: 'a'.repeat(60),
          bio: 'a'.repeat(600),
          hourlyRate: -10,
        },
        2
      );

      expect(result.valid).toBe(false);
      expect(result.errors.length).toBeGreaterThanOrEqual(4);

      const fields = result.errors.map((e) => e.field);
      expect(fields).toContain('userEmail');
      expect(fields).toContain('nickname');
      expect(fields).toContain('bio');
      expect(fields).toContain('hourlyRate');
    });
  });

  describe('Game validation collects all errors', () => {
    it('should collect all field errors together', () => {
      const result = validateGameRow(
        {
          key: 'invalid key!',
          name: '',
          category: 'a'.repeat(60),
          description: 'a'.repeat(600),
        },
        2
      );

      expect(result.valid).toBe(false);
      expect(result.errors.length).toBeGreaterThanOrEqual(4);

      const fields = result.errors.map((e) => e.field);
      expect(fields).toContain('key');
      expect(fields).toContain('name');
      expect(fields).toContain('category');
      expect(fields).toContain('description');
    });
  });

  describe('Batch validation', () => {
    it('should validate all rows independently', () => {
      const rows = [
        { rowNumber: 2, data: { name: 'User 1', email: 'a@example.com', phone: '13800138000' } },
        { rowNumber: 3, data: { name: '', email: '', phone: '' } },
        { rowNumber: 4, data: { name: 'User 3', email: 'c@example.com', phone: '13800138002' } },
      ];

      const result = validateUserData(rows, { checkInternalDuplicates: false });

      expect(result.validRows.length).toBe(2);
      expect(result.invalidRows.length).toBe(1);
      expect(result.invalidRows[0].rowNumber).toBe(3);
    });
  });
});
