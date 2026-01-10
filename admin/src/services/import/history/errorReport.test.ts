/**
 * Tests for Import Error Report Generator
 *
 * @module services/import/history/errorReport.test
 */

import { describe, it, expect } from 'vitest';
import {
  generateErrorReport,
  hasErrorDetails,
  getErrorSummary,
} from './errorReport';
import { createImportHistoryRecord } from './storage';
import type { ImportHistoryRecord, ImportRowResult } from './types';

/**
 * Helper to read blob content as text
 */
async function readBlobAsText(blob: Blob): Promise<string> {
  return new Promise((resolve, reject) => {
    const reader = new FileReader();
    reader.onload = () => resolve(reader.result as string);
    reader.onerror = reject;
    reader.readAsText(blob);
  });
}

describe('Error Report Generator', () => {
  describe('generateErrorReport', () => {
    it('should generate CSV report with correct headers', () => {
      const record = createImportHistoryRecord('user', 'users.csv', 1024, 1, 'Admin', 3);
      record.rowResults = [
        {
          rowNumber: 2,
          success: false,
          originalData: { name: 'Test', email: 'invalid', phone: '123' },
          errorField: 'email',
          errorMessage: 'Invalid email format',
        },
      ];
      record.status = 'failed';
      record.skippedCount = 1;

      const report = generateErrorReport(record);

      expect(report.blob).toBeInstanceOf(Blob);
      expect(report.mimeType).toBe('text/csv;charset=utf-8');
      expect(report.rowCount).toBe(1);
      expect(report.fileName).toContain('用户导入错误报告');
      expect(report.fileName.endsWith('.csv')).toBe(true);
    });

    it('should include original data columns in report', async () => {
      const record = createImportHistoryRecord('user', 'users.csv', 1024, 1, 'Admin', 1);
      record.rowResults = [
        {
          rowNumber: 2,
          success: false,
          originalData: { name: 'John', email: 'bad@', phone: '13800138000', role: 'user', status: 'active' },
          errorField: 'email',
          errorMessage: 'Invalid email format',
        },
      ];

      const report = generateErrorReport(record);
      const content = await readBlobAsText(report.blob);

      // Should contain original data
      expect(content).toContain('John');
      expect(content).toContain('bad@');
      expect(content).toContain('13800138000');
    });

    it('should include status and error columns', async () => {
      const record = createImportHistoryRecord('user', 'users.csv', 1024, 1, 'Admin', 1);
      record.rowResults = [
        {
          rowNumber: 2,
          success: false,
          originalData: { name: 'Test' },
          errorField: 'email',
          errorMessage: 'Email is required',
        },
      ];

      const report = generateErrorReport(record);
      const content = await readBlobAsText(report.blob);

      // Should contain status columns
      expect(content).toContain('导入状态');
      expect(content).toContain('错误字段');
      expect(content).toContain('错误信息');
      expect(content).toContain('失败');
      expect(content).toContain('email');
      expect(content).toContain('Email is required');
    });

    it('should exclude successful rows by default', async () => {
      const record = createImportHistoryRecord('user', 'users.csv', 1024, 1, 'Admin', 2);
      record.rowResults = [
        {
          rowNumber: 2,
          success: true,
          originalData: { name: 'Good', email: 'good@test.com' },
          createdRecordId: 1,
        },
        {
          rowNumber: 3,
          success: false,
          originalData: { name: 'Bad', email: 'bad' },
          errorField: 'email',
          errorMessage: 'Invalid email',
        },
      ];

      const report = generateErrorReport(record);
      const content = await readBlobAsText(report.blob);

      expect(report.rowCount).toBe(1);
      expect(content).not.toContain('Good');
      expect(content).toContain('Bad');
    });

    it('should include successful rows when option is set', async () => {
      const record = createImportHistoryRecord('user', 'users.csv', 1024, 1, 'Admin', 2);
      record.rowResults = [
        {
          rowNumber: 2,
          success: true,
          originalData: { name: 'Good', email: 'good@test.com' },
          createdRecordId: 1,
        },
        {
          rowNumber: 3,
          success: false,
          originalData: { name: 'Bad', email: 'bad' },
          errorField: 'email',
          errorMessage: 'Invalid email',
        },
      ];

      const report = generateErrorReport(record, { includeSuccessful: true });
      const content = await readBlobAsText(report.blob);

      expect(report.rowCount).toBe(2);
      expect(content).toContain('Good');
      expect(content).toContain('成功');
      expect(content).toContain('Bad');
      expect(content).toContain('失败');
    });

    it('should escape CSV special characters', async () => {
      const record = createImportHistoryRecord('user', 'users.csv', 1024, 1, 'Admin', 1);
      record.rowResults = [
        {
          rowNumber: 2,
          success: false,
          originalData: { name: 'Test, User', email: 'test@test.com' },
          errorField: 'name',
          errorMessage: 'Name contains "special" chars',
        },
      ];

      const report = generateErrorReport(record);
      const content = await readBlobAsText(report.blob);

      // Comma should be escaped with quotes
      expect(content).toContain('"Test, User"');
      // Quotes should be escaped
      expect(content).toContain('""special""');
    });

    it('should generate correct file name for different import types', () => {
      const userRecord = createImportHistoryRecord('user', 'users.csv', 1024, 1, 'Admin', 1);
      userRecord.rowResults = [{ rowNumber: 2, success: false, originalData: {} }];
      
      const playerRecord = createImportHistoryRecord('player', 'players.csv', 1024, 1, 'Admin', 1);
      playerRecord.rowResults = [{ rowNumber: 2, success: false, originalData: {} }];
      
      const gameRecord = createImportHistoryRecord('game', 'games.csv', 1024, 1, 'Admin', 1);
      gameRecord.rowResults = [{ rowNumber: 2, success: false, originalData: {} }];

      expect(generateErrorReport(userRecord).fileName).toContain('用户');
      expect(generateErrorReport(playerRecord).fileName).toContain('陪玩师');
      expect(generateErrorReport(gameRecord).fileName).toContain('游戏');
    });

    it('should use custom file name when provided', () => {
      const record = createImportHistoryRecord('user', 'users.csv', 1024, 1, 'Admin', 1);
      record.rowResults = [{ rowNumber: 2, success: false, originalData: {} }];

      const report = generateErrorReport(record, { fileName: 'custom-report.csv' });

      expect(report.fileName).toBe('custom-report.csv');
    });

    it('should handle empty row results', () => {
      const record = createImportHistoryRecord('user', 'users.csv', 1024, 1, 'Admin', 0);
      record.rowResults = [];

      const report = generateErrorReport(record);

      expect(report.rowCount).toBe(0);
    });

    it('should handle records without row results', () => {
      const record = createImportHistoryRecord('user', 'users.csv', 1024, 1, 'Admin', 5);
      // No rowResults set

      const report = generateErrorReport(record);

      expect(report.rowCount).toBe(0);
    });

    it('should add BOM for Excel compatibility', async () => {
      const record = createImportHistoryRecord('user', 'users.csv', 1024, 1, 'Admin', 1);
      record.rowResults = [
        {
          rowNumber: 2,
          success: false,
          originalData: { name: '中文名' },
          errorMessage: '错误信息',
        },
      ];

      const report = generateErrorReport(record);
      
      // Read as ArrayBuffer to check raw bytes
      const arrayBuffer = await new Promise<ArrayBuffer>((resolve, reject) => {
        const reader = new FileReader();
        reader.onload = () => resolve(reader.result as ArrayBuffer);
        reader.onerror = reject;
        reader.readAsArrayBuffer(report.blob);
      });
      
      const bytes = new Uint8Array(arrayBuffer);
      
      // UTF-8 BOM is EF BB BF (239, 187, 191)
      expect(bytes[0]).toBe(0xEF);
      expect(bytes[1]).toBe(0xBB);
      expect(bytes[2]).toBe(0xBF);
    });
  });

  describe('hasErrorDetails', () => {
    it('should return false for records without row results', () => {
      const record = createImportHistoryRecord('user', 'users.csv', 1024, 1, 'Admin', 5);
      
      expect(hasErrorDetails(record)).toBe(false);
    });

    it('should return false for records with empty row results', () => {
      const record = createImportHistoryRecord('user', 'users.csv', 1024, 1, 'Admin', 0);
      record.rowResults = [];
      
      expect(hasErrorDetails(record)).toBe(false);
    });

    it('should return false for records with only successful rows', () => {
      const record = createImportHistoryRecord('user', 'users.csv', 1024, 1, 'Admin', 2);
      record.rowResults = [
        { rowNumber: 2, success: true, originalData: {}, createdRecordId: 1 },
        { rowNumber: 3, success: true, originalData: {}, createdRecordId: 2 },
      ];
      
      expect(hasErrorDetails(record)).toBe(false);
    });

    it('should return true for records with error rows', () => {
      const record = createImportHistoryRecord('user', 'users.csv', 1024, 1, 'Admin', 2);
      record.rowResults = [
        { rowNumber: 2, success: true, originalData: {}, createdRecordId: 1 },
        { rowNumber: 3, success: false, originalData: {}, errorMessage: 'Error' },
      ];
      
      expect(hasErrorDetails(record)).toBe(true);
    });
  });

  describe('getErrorSummary', () => {
    it('should return correct counts from row results', () => {
      const record = createImportHistoryRecord('user', 'users.csv', 1024, 1, 'Admin', 5);
      record.rowResults = [
        { rowNumber: 2, success: true, originalData: {}, createdRecordId: 1 },
        { rowNumber: 3, success: true, originalData: {}, createdRecordId: 2 },
        { rowNumber: 4, success: false, originalData: {}, errorField: 'email', errorMessage: 'Invalid' },
        { rowNumber: 5, success: false, originalData: {}, errorField: 'email', errorMessage: 'Duplicate' },
        { rowNumber: 6, success: false, originalData: {}, errorField: 'phone', errorMessage: 'Invalid' },
      ];

      const summary = getErrorSummary(record);

      expect(summary.totalRows).toBe(5);
      expect(summary.successCount).toBe(2);
      expect(summary.errorCount).toBe(3);
      expect(summary.errorsByField).toEqual({ email: 2, phone: 1 });
    });

    it('should use record counts when row results not available', () => {
      const record = createImportHistoryRecord('user', 'users.csv', 1024, 1, 'Admin', 10);
      record.importedCount = 7;
      record.skippedCount = 3;

      const summary = getErrorSummary(record);

      expect(summary.totalRows).toBe(10);
      expect(summary.successCount).toBe(7);
      expect(summary.errorCount).toBe(3);
      expect(summary.errorsByField).toEqual({});
    });

    it('should handle errors without field specified', () => {
      const record = createImportHistoryRecord('user', 'users.csv', 1024, 1, 'Admin', 2);
      record.rowResults = [
        { rowNumber: 2, success: false, originalData: {}, errorMessage: 'Unknown error' },
        { rowNumber: 3, success: false, originalData: {}, errorField: 'email', errorMessage: 'Invalid' },
      ];

      const summary = getErrorSummary(record);

      expect(summary.errorCount).toBe(2);
      expect(summary.errorsByField).toEqual({ email: 1 });
    });
  });
});


/**
 * Property-Based Tests for Import Error Report Format
 */
import * as fc from 'fast-check';

/**
 * Arbitrary for ImportType
 */
const importTypeArb = fc.constantFrom('user', 'player', 'game') as fc.Arbitrary<'user' | 'player' | 'game'>;

/**
 * Arbitrary for error field names based on import type
 */
const errorFieldArb = fc.constantFrom(
  'email', 'phone', 'name', 'password', 'role', 'status',
  'userEmail', 'nickname', 'hourlyRate', 'key', 'category'
);

/**
 * Arbitrary for error messages
 */
const errorMessageArb = fc.oneof(
  fc.constant('Field is required'),
  fc.constant('Invalid email format'),
  fc.constant('Invalid phone format'),
  fc.constant('Duplicate value'),
  fc.constant('Value too long'),
  fc.string({ minLength: 5, maxLength: 100 })
);

/**
 * Arbitrary for original data based on import type
 */
function originalDataArb(type: 'user' | 'player' | 'game'): fc.Arbitrary<Record<string, unknown>> {
  switch (type) {
    case 'user':
      return fc.record({
        name: fc.string({ minLength: 1, maxLength: 50 }),
        email: fc.emailAddress(),
        phone: fc.string({ minLength: 11, maxLength: 11 }),
        role: fc.constantFrom('user', 'player', 'admin'),
        status: fc.constantFrom('active', 'banned', 'suspended'),
      });
    case 'player':
      return fc.record({
        userEmail: fc.emailAddress(),
        nickname: fc.option(fc.string({ minLength: 1, maxLength: 30 }), { nil: undefined }),
        bio: fc.option(fc.string({ minLength: 1, maxLength: 200 }), { nil: undefined }),
        hourlyRate: fc.option(fc.integer({ min: 20, max: 200 }), { nil: undefined }),
        mainGame: fc.option(fc.string({ minLength: 1, maxLength: 50 }), { nil: undefined }),
        skillTags: fc.option(fc.string({ minLength: 1, maxLength: 100 }), { nil: undefined }),
      });
    case 'game':
      return fc.record({
        key: fc.string({ minLength: 1, maxLength: 30 }),
        name: fc.string({ minLength: 1, maxLength: 50 }),
        category: fc.option(fc.string({ minLength: 1, maxLength: 30 }), { nil: undefined }),
        description: fc.option(fc.string({ minLength: 1, maxLength: 200 }), { nil: undefined }),
        isActive: fc.option(fc.boolean(), { nil: undefined }),
      });
  }
}

/**
 * Arbitrary for ImportRowResult with errors
 */
function errorRowResultArb(type: 'user' | 'player' | 'game'): fc.Arbitrary<ImportRowResult> {
  return fc.record({
    rowNumber: fc.integer({ min: 2, max: 10000 }),
    success: fc.constant(false),
    originalData: originalDataArb(type),
    errorMessage: errorMessageArb,
    errorField: errorFieldArb,
    createdRecordId: fc.constant(undefined),
  });
}

/**
 * Arbitrary for ImportRowResult with success
 */
function successRowResultArb(type: 'user' | 'player' | 'game'): fc.Arbitrary<ImportRowResult> {
  return fc.record({
    rowNumber: fc.integer({ min: 2, max: 10000 }),
    success: fc.constant(true),
    originalData: originalDataArb(type),
    errorMessage: fc.constant(undefined),
    errorField: fc.constant(undefined),
    createdRecordId: fc.integer({ min: 1, max: 100000 }),
  });
}

/**
 * Property 24: Import Result Report Format
 * For any import result download, the report SHALL contain all original data columns
 * plus status and error columns.
 *
 * **Feature: admin-phase3-improvements, Property 24: Import Result Report Format**
 * **Validates: Requirements 9.5**
 */
describe('Property 24: Import Result Report Format', () => {
  /**
   * Helper to read blob content as text
   */
  async function readBlobAsText(blob: Blob): Promise<string> {
    return new Promise((resolve, reject) => {
      const reader = new FileReader();
      reader.onload = () => resolve(reader.result as string);
      reader.onerror = reject;
      reader.readAsText(blob);
    });
  }

  /**
   * Get expected column headers for a type
   */
  function getExpectedHeaders(type: 'user' | 'player' | 'game'): string[] {
    const typeHeaders: Record<string, string[]> = {
      user: ['姓名', '邮箱', '手机号', '角色', '状态'],
      player: ['用户邮箱', '昵称', '简介', '时薪(元)', '主游戏', '技能标签'],
      game: ['游戏标识', '游戏名称', '分类', '描述', '是否启用'],
    };
    return typeHeaders[type];
  }

  it('should contain all original data column headers for any import type', async () => {
    await fc.assert(
      fc.asyncProperty(
        importTypeArb,
        fc.array(fc.oneof(
          fc.record({
            rowNumber: fc.integer({ min: 2, max: 100 }),
            success: fc.boolean(),
            originalData: fc.dictionary(fc.string({ minLength: 1, maxLength: 20 }), fc.string()),
            errorMessage: fc.option(fc.string(), { nil: undefined }),
            errorField: fc.option(fc.string(), { nil: undefined }),
            createdRecordId: fc.option(fc.integer({ min: 1 }), { nil: undefined }),
          })
        ), { minLength: 1, maxLength: 5 }),
        async (type, rowResults) => {
          const record = createImportHistoryRecord(type, 'test.csv', 1024, 1, 'Admin', rowResults.length);
          record.rowResults = rowResults;
          record.status = 'partial';

          const report = generateErrorReport(record, { includeSuccessful: true });
          const content = await readBlobAsText(report.blob);

          // Get expected headers for this type
          const expectedHeaders = getExpectedHeaders(type);

          // Verify all original data column headers are present
          for (const header of expectedHeaders) {
            expect(content).toContain(header);
          }
        }
      ),
      { numRuns: 50 }
    );
  });

  it('should contain status column in report', async () => {
    await fc.assert(
      fc.asyncProperty(
        importTypeArb,
        async (type) => {
          const record = createImportHistoryRecord(type, 'test.csv', 1024, 1, 'Admin', 1);
          record.rowResults = [{
            rowNumber: 2,
            success: false,
            originalData: {},
            errorMessage: 'Test error',
          }];

          const report = generateErrorReport(record);
          const content = await readBlobAsText(report.blob);

          // Status column header must be present
          expect(content).toContain('导入状态');
        }
      ),
      { numRuns: 30 }
    );
  });

  it('should contain error field column in report', async () => {
    await fc.assert(
      fc.asyncProperty(
        importTypeArb,
        async (type) => {
          const record = createImportHistoryRecord(type, 'test.csv', 1024, 1, 'Admin', 1);
          record.rowResults = [{
            rowNumber: 2,
            success: false,
            originalData: {},
            errorField: 'email',
            errorMessage: 'Invalid',
          }];

          const report = generateErrorReport(record);
          const content = await readBlobAsText(report.blob);

          // Error field column header must be present
          expect(content).toContain('错误字段');
        }
      ),
      { numRuns: 30 }
    );
  });

  it('should contain error message column in report', async () => {
    await fc.assert(
      fc.asyncProperty(
        importTypeArb,
        async (type) => {
          const record = createImportHistoryRecord(type, 'test.csv', 1024, 1, 'Admin', 1);
          record.rowResults = [{
            rowNumber: 2,
            success: false,
            originalData: {},
            errorMessage: 'Test error message',
          }];

          const report = generateErrorReport(record);
          const content = await readBlobAsText(report.blob);

          // Error message column header must be present
          expect(content).toContain('错误信息');
        }
      ),
      { numRuns: 30 }
    );
  });

  it('should preserve original data values in report for error rows', async () => {
    await fc.assert(
      fc.asyncProperty(
        importTypeArb.chain((type) => 
          fc.tuple(
            fc.constant(type),
            fc.array(errorRowResultArb(type), { minLength: 1, maxLength: 3 })
          )
        ),
        async ([type, errorRows]) => {
          const record = createImportHistoryRecord(type, 'test.csv', 1024, 1, 'Admin', errorRows.length);
          record.rowResults = errorRows;
          record.status = 'failed';

          const report = generateErrorReport(record);
          const content = await readBlobAsText(report.blob);

          // Each error row's original data values should be in the report
          for (const row of errorRows) {
            for (const value of Object.values(row.originalData)) {
              if (value !== undefined && value !== null && String(value).length > 0) {
                // The value should appear in the content (possibly escaped)
                const stringValue = String(value);
                // For simple values without special chars, they should appear directly
                if (!stringValue.includes(',') && !stringValue.includes('"') && !stringValue.includes('\n')) {
                  expect(content).toContain(stringValue);
                }
              }
            }
          }
        }
      ),
      { numRuns: 50 }
    );
  });

  it('should preserve error messages in report', async () => {
    await fc.assert(
      fc.asyncProperty(
        importTypeArb.chain((type) => 
          fc.tuple(
            fc.constant(type),
            fc.array(errorRowResultArb(type), { minLength: 1, maxLength: 3 })
          )
        ),
        async ([type, errorRows]) => {
          const record = createImportHistoryRecord(type, 'test.csv', 1024, 1, 'Admin', errorRows.length);
          record.rowResults = errorRows;
          record.status = 'failed';

          const report = generateErrorReport(record);
          const content = await readBlobAsText(report.blob);

          // Each error message should appear in the report
          for (const row of errorRows) {
            if (row.errorMessage) {
              const errorMsg = row.errorMessage;
              // For simple messages without special chars
              if (!errorMsg.includes(',') && !errorMsg.includes('"') && !errorMsg.includes('\n')) {
                expect(content).toContain(errorMsg);
              }
            }
          }
        }
      ),
      { numRuns: 50 }
    );
  });

  it('should show success status for successful rows when included', async () => {
    await fc.assert(
      fc.asyncProperty(
        importTypeArb.chain((type) => 
          fc.tuple(
            fc.constant(type),
            fc.array(successRowResultArb(type), { minLength: 1, maxLength: 3 })
          )
        ),
        async ([type, successRows]) => {
          const record = createImportHistoryRecord(type, 'test.csv', 1024, 1, 'Admin', successRows.length);
          record.rowResults = successRows;
          record.status = 'completed';
          record.importedCount = successRows.length;

          const report = generateErrorReport(record, { includeSuccessful: true });
          const content = await readBlobAsText(report.blob);

          // Success status should appear for each successful row
          const successCount = (content.match(/成功/g) || []).length;
          expect(successCount).toBeGreaterThanOrEqual(successRows.length);
        }
      ),
      { numRuns: 50 }
    );
  });

  it('should show failure status for error rows', async () => {
    await fc.assert(
      fc.asyncProperty(
        importTypeArb.chain((type) => 
          fc.tuple(
            fc.constant(type),
            fc.array(errorRowResultArb(type), { minLength: 1, maxLength: 3 })
          )
        ),
        async ([type, errorRows]) => {
          const record = createImportHistoryRecord(type, 'test.csv', 1024, 1, 'Admin', errorRows.length);
          record.rowResults = errorRows;
          record.status = 'failed';
          record.skippedCount = errorRows.length;

          const report = generateErrorReport(record);
          const content = await readBlobAsText(report.blob);

          // Failure status should appear for each error row
          const failureCount = (content.match(/失败/g) || []).length;
          expect(failureCount).toBeGreaterThanOrEqual(errorRows.length);
        }
      ),
      { numRuns: 50 }
    );
  });

  it('should have correct number of data rows in report', async () => {
    await fc.assert(
      fc.asyncProperty(
        importTypeArb,
        fc.integer({ min: 1, max: 10 }),
        fc.integer({ min: 0, max: 10 }),
        async (type, errorCount, successCount) => {
          const errorRows: ImportRowResult[] = [];
          const successRows: ImportRowResult[] = [];

          for (let i = 0; i < errorCount; i++) {
            errorRows.push({
              rowNumber: i + 2,
              success: false,
              originalData: {},
              errorMessage: `Error ${i}`,
            });
          }

          for (let i = 0; i < successCount; i++) {
            successRows.push({
              rowNumber: errorCount + i + 2,
              success: true,
              originalData: {},
              createdRecordId: i + 1,
            });
          }

          const record = createImportHistoryRecord(type, 'test.csv', 1024, 1, 'Admin', errorCount + successCount);
          record.rowResults = [...errorRows, ...successRows];
          record.status = 'partial';
          record.importedCount = successCount;
          record.skippedCount = errorCount;

          // Without includeSuccessful, only error rows
          const errorOnlyReport = generateErrorReport(record);
          expect(errorOnlyReport.rowCount).toBe(errorCount);

          // With includeSuccessful, all rows
          const allRowsReport = generateErrorReport(record, { includeSuccessful: true });
          expect(allRowsReport.rowCount).toBe(errorCount + successCount);
        }
      ),
      { numRuns: 50 }
    );
  });
});
