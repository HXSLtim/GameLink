/**
 * Property-Based Tests for Structure Validator
 *
 * **Feature: admin-phase3-improvements, Property 14: Import Structure Validation**
 * **Validates: Requirements 5.2**
 *
 * Tests that the structure validator correctly validates that all required columns
 * are present and reports missing columns as structural errors.
 */

import { describe, it, expect } from 'vitest';
import * as fc from 'fast-check';
import {
  validateStructure,
  createColumnMapping,
  mapRowToTemplate,
  getMissingRequiredColumns,
  hasAllRequiredColumns,
} from './structureValidator';
import { userImportTemplate } from '../templates/userTemplate';
import { playerImportTemplate } from '../templates/playerTemplate';
import { gameImportTemplate } from '../templates/gameTemplate';
import type { ImportTemplate } from '../templates/types';

/**
 * Get all required column labels (Chinese) from a template
 */
function getRequiredLabels(template: ImportTemplate): string[] {
  return template.columns.filter((col) => col.required).map((col) => col.labelZh);
}

/**
 * Get all column labels (Chinese) from a template
 */
function getAllLabels(template: ImportTemplate): string[] {
  return template.columns.map((col) => col.labelZh);
}

describe('Structure Validator - Property Tests', () => {
  /**
   * **Feature: admin-phase3-improvements, Property 14: Import Structure Validation**
   * **Validates: Requirements 5.2**
   *
   * For any import file, the service SHALL validate that all required columns are
   * present and report missing columns as structural errors.
   */
  describe('Property 14: Import Structure Validation', () => {
    const templates = [
      { name: 'user', template: userImportTemplate },
      { name: 'player', template: playerImportTemplate },
      { name: 'game', template: gameImportTemplate },
    ];

    templates.forEach(({ name, template }) => {
      describe(`${name} template`, () => {
        it('should accept headers with all required columns present', () => {
          const requiredLabels = getRequiredLabels(template);
          const optionalLabels = template.columns
            .filter((col) => !col.required)
            .map((col) => col.labelZh);

          // Generate subsets of optional columns to include
          const optionalSubsetArb = fc.subarray(optionalLabels);

          fc.assert(
            fc.property(optionalSubsetArb, (optionalSubset) => {
              const headers = [...requiredLabels, ...optionalSubset];
              const result = validateStructure(headers, template);

              // Should be valid when all required columns are present
              expect(result.valid).toBe(true);
              expect(result.missingColumns).toHaveLength(0);
            }),
            { numRuns: 50 }
          );
        });

        it('should reject headers with missing required columns', () => {
          const requiredLabels = getRequiredLabels(template);

          // Skip if no required columns
          if (requiredLabels.length === 0) return;

          // Generate proper subsets of required columns (at least one missing)
          const missingSubsetArb = fc
            .subarray(requiredLabels, { minLength: 0, maxLength: requiredLabels.length - 1 })
            .filter((subset) => subset.length < requiredLabels.length);

          fc.assert(
            fc.property(missingSubsetArb, (presentColumns) => {
              const result = validateStructure(presentColumns, template);

              // Should be invalid when required columns are missing
              expect(result.valid).toBe(false);
              expect(result.missingColumns.length).toBeGreaterThan(0);

              // Missing columns should be the ones not in presentColumns
              const missingSet = new Set(result.missingColumns);
              for (const label of requiredLabels) {
                if (!presentColumns.includes(label)) {
                  expect(missingSet.has(label)).toBe(true);
                }
              }
            }),
            { numRuns: 50 }
          );
        });

        it('should report extra columns not in template', () => {
          const allLabels = getAllLabels(template);
          // Create a set of all possible column identifiers (normalized)
          const allIdentifiers = new Set(
            template.columns.flatMap((col) => [
              col.key.toLowerCase().replace(/[\s_-]+/g, ''),
              col.label.toLowerCase().replace(/[\s_-]+/g, ''),
              col.labelZh.toLowerCase().replace(/[\s_-]+/g, ''),
            ])
          );

          // Generate random extra column names that don't match any template column
          const extraColumnArb = fc
            .array(
              fc.string({ minLength: 1, maxLength: 20 }).filter((s) => {
                const normalized = s.toLowerCase().replace(/[\s_-]+/g, '');
                return normalized.length > 0 && !allIdentifiers.has(normalized);
              }),
              { minLength: 1, maxLength: 5 }
            )
            .filter((extras) => extras.length > 0);

          fc.assert(
            fc.property(extraColumnArb, (extraColumns) => {
              const headers = [...allLabels, ...extraColumns];
              const result = validateStructure(headers, template);

              // Should report extra columns
              expect(result.extraColumns.length).toBe(extraColumns.length);
              for (const extra of extraColumns) {
                expect(result.extraColumns).toContain(extra);
              }
            }),
            { numRuns: 50 }
          );
        });

        it('should match columns by key, label, or Chinese label', () => {
          // Test that columns can be matched by any of their identifiers
          const columnArb = fc.constantFrom(...template.columns);

          fc.assert(
            fc.property(columnArb, (column) => {
              // Test matching by key
              const headersByKey = template.columns.map((col) =>
                col.key === column.key ? column.key : col.labelZh
              );
              const resultByKey = validateStructure(headersByKey, template);
              expect(resultByKey.matchedColumns).toContain(column.key);

              // Test matching by label
              const headersByLabel = template.columns.map((col) =>
                col.key === column.key ? column.label : col.labelZh
              );
              const resultByLabel = validateStructure(headersByLabel, template);
              expect(resultByLabel.matchedColumns).toContain(column.key);

              // Test matching by Chinese label
              const headersByLabelZh = template.columns.map((col) => col.labelZh);
              const resultByLabelZh = validateStructure(headersByLabelZh, template);
              expect(resultByLabelZh.matchedColumns).toContain(column.key);
            }),
            { numRuns: 20 }
          );
        });
      });
    });

    it('should handle empty headers array', () => {
      const templateArb = fc.constantFrom(userImportTemplate, playerImportTemplate, gameImportTemplate);

      fc.assert(
        fc.property(templateArb, (template) => {
          const result = validateStructure([], template);

          // Should be invalid if template has required columns
          const requiredCount = template.columns.filter((col) => col.required).length;
          if (requiredCount > 0) {
            expect(result.valid).toBe(false);
            expect(result.missingColumns.length).toBe(requiredCount);
          }
        }),
        { numRuns: 10 }
      );
    });

    it('should be case-insensitive for header matching', () => {
      const templateArb = fc.constantFrom(userImportTemplate, playerImportTemplate, gameImportTemplate);
      const caseTransformArb = fc.constantFrom(
        (s: string) => s.toLowerCase(),
        (s: string) => s.toUpperCase(),
        (s: string) => s.charAt(0).toUpperCase() + s.slice(1).toLowerCase()
      );

      fc.assert(
        fc.property(templateArb, caseTransformArb, (template, transform) => {
          const headers = template.columns.map((col) => transform(col.labelZh));
          const result = validateStructure(headers, template);

          // Should match all columns regardless of case
          expect(result.matchedColumns.length).toBe(template.columns.length);
        }),
        { numRuns: 30 }
      );
    });
  });
});

describe('Structure Validator - Unit Tests', () => {
  describe('validateStructure', () => {
    it('should validate user template with all columns', () => {
      const headers = ['姓名', '邮箱', '手机号', '角色', '状态'];
      const result = validateStructure(headers, userImportTemplate);

      expect(result.valid).toBe(true);
      expect(result.missingColumns).toHaveLength(0);
      expect(result.matchedColumns).toHaveLength(5);
    });

    it('should validate user template with only required columns', () => {
      const headers = ['姓名', '邮箱', '手机号'];
      const result = validateStructure(headers, userImportTemplate);

      expect(result.valid).toBe(true);
      expect(result.missingColumns).toHaveLength(0);
      expect(result.matchedColumns).toHaveLength(3);
    });

    it('should fail validation when required column is missing', () => {
      const headers = ['姓名', '邮箱']; // Missing 手机号
      const result = validateStructure(headers, userImportTemplate);

      expect(result.valid).toBe(false);
      expect(result.missingColumns).toContain('手机号');
    });

    it('should report extra columns', () => {
      const headers = ['姓名', '邮箱', '手机号', '未知列'];
      const result = validateStructure(headers, userImportTemplate);

      expect(result.valid).toBe(true);
      expect(result.extraColumns).toContain('未知列');
    });
  });

  describe('createColumnMapping', () => {
    it('should create correct mapping for user template', () => {
      const headers = ['姓名', '邮箱', '手机号'];
      const mapping = createColumnMapping(headers, userImportTemplate);

      expect(mapping.size).toBe(3);
      expect(mapping.get(0)?.templateKey).toBe('name');
      expect(mapping.get(1)?.templateKey).toBe('email');
      expect(mapping.get(2)?.templateKey).toBe('phone');
    });

    it('should handle mixed key and label headers', () => {
      const headers = ['name', '邮箱', 'phone'];
      const mapping = createColumnMapping(headers, userImportTemplate);

      expect(mapping.size).toBe(3);
      expect(mapping.get(0)?.templateKey).toBe('name');
      expect(mapping.get(1)?.templateKey).toBe('email');
      expect(mapping.get(2)?.templateKey).toBe('phone');
    });
  });

  describe('mapRowToTemplate', () => {
    it('should map row data to template keys', () => {
      const headers = ['姓名', '邮箱', '手机号'];
      const mapping = createColumnMapping(headers, userImportTemplate);
      const row = {
        姓名: '张三',
        邮箱: 'test@example.com',
        手机号: '13800138000',
      };

      const mapped = mapRowToTemplate(row, headers, mapping, userImportTemplate);

      expect(mapped.name).toBe('张三');
      expect(mapped.email).toBe('test@example.com');
      expect(mapped.phone).toBe('13800138000');
    });

    it('should apply default values for missing optional columns', () => {
      const headers = ['姓名', '邮箱', '手机号'];
      const mapping = createColumnMapping(headers, userImportTemplate);
      const row = {
        姓名: '张三',
        邮箱: 'test@example.com',
        手机号: '13800138000',
      };

      const mapped = mapRowToTemplate(row, headers, mapping, userImportTemplate);

      expect(mapped.role).toBe('user'); // Default value
      expect(mapped.status).toBe('active'); // Default value
    });
  });

  describe('getMissingRequiredColumns', () => {
    it('should return empty array when all required columns present', () => {
      const headers = ['姓名', '邮箱', '手机号'];
      const missing = getMissingRequiredColumns(headers, userImportTemplate);

      expect(missing).toHaveLength(0);
    });

    it('should return missing required columns', () => {
      const headers = ['姓名'];
      const missing = getMissingRequiredColumns(headers, userImportTemplate);

      expect(missing.length).toBe(2);
      expect(missing.map((c) => c.key)).toContain('email');
      expect(missing.map((c) => c.key)).toContain('phone');
    });
  });

  describe('hasAllRequiredColumns', () => {
    it('should return true when all required columns present', () => {
      const headers = ['姓名', '邮箱', '手机号'];
      expect(hasAllRequiredColumns(headers, userImportTemplate)).toBe(true);
    });

    it('should return false when required columns missing', () => {
      const headers = ['姓名'];
      expect(hasAllRequiredColumns(headers, userImportTemplate)).toBe(false);
    });
  });
});
