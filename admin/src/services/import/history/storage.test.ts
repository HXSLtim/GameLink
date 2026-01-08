/**
 * Property-Based Tests for Import History Storage
 *
 * @module services/import/history/storage.test
 */

import { describe, it, expect, beforeEach } from 'vitest';
import * as fc from 'fast-check';
import {
  ImportHistoryStorage,
  createImportHistoryRecord,
  generateImportId,
} from './storage';
import type { ImportHistoryRecord, ImportRowResult } from './types';

/**
 * Create a mock in-memory storage for testing
 */
function createMockStorage(): Storage {
  const store = new Map<string, string>();
  return {
    getItem: (key: string) => store.get(key) ?? null,
    setItem: (key: string, value: string) => { store.set(key, value); },
    removeItem: (key: string) => { store.delete(key); },
    clear: () => { store.clear(); },
    get length() { return store.size; },
    key: (index: number) => Array.from(store.keys())[index] ?? null,
  };
}

/**
 * Arbitrary for ImportType
 */
const importTypeArb = fc.constantFrom('user', 'player', 'game') as fc.Arbitrary<'user' | 'player' | 'game'>;

/**
 * Arbitrary for ImportStatus
 */
const importStatusArb = fc.constantFrom('pending', 'processing', 'completed', 'failed', 'partial') as fc.Arbitrary<'pending' | 'processing' | 'completed' | 'failed' | 'partial'>;

/**
 * Arbitrary for ImportRowResult
 */
const importRowResultArb: fc.Arbitrary<ImportRowResult> = fc.record({
  rowNumber: fc.integer({ min: 1, max: 10000 }),
  success: fc.boolean(),
  originalData: fc.dictionary(fc.string({ minLength: 1, maxLength: 20 }), fc.string()),
  errorMessage: fc.option(fc.string({ minLength: 1, maxLength: 200 }), { nil: undefined }),
  errorField: fc.option(fc.string({ minLength: 1, maxLength: 50 }), { nil: undefined }),
  createdRecordId: fc.option(fc.integer({ min: 1 }), { nil: undefined }),
});

/**
 * Arbitrary for ISO date string using timestamp approach (more reliable)
 */
const minTimestamp = new Date('2020-01-01').getTime();
const maxTimestamp = new Date('2030-01-01').getTime();
const isoDateStringArb = fc.integer({ min: minTimestamp, max: maxTimestamp })
  .map(ts => new Date(ts).toISOString());

/**
 * Arbitrary for optional ISO date string
 */
const optionalIsoDateStringArb = fc.option(isoDateStringArb, { nil: undefined });

/**
 * Arbitrary for ImportHistoryRecord
 */
const importHistoryRecordArb: fc.Arbitrary<ImportHistoryRecord> = fc.record({
  id: fc.string({ minLength: 10, maxLength: 50 }),
  type: importTypeArb,
  fileName: fc.string({ minLength: 1, maxLength: 100 }),
  fileSize: fc.integer({ min: 1, max: 10 * 1024 * 1024 }),
  uploadedBy: fc.integer({ min: 1, max: 10000 }),
  uploadedByName: fc.option(fc.string({ minLength: 1, maxLength: 50 }), { nil: undefined }),
  uploadedAt: isoDateStringArb,
  completedAt: optionalIsoDateStringArb,
  totalRows: fc.integer({ min: 0, max: 10000 }),
  importedCount: fc.integer({ min: 0, max: 10000 }),
  skippedCount: fc.integer({ min: 0, max: 10000 }),
  status: importStatusArb,
  errorSummary: fc.option(fc.string({ minLength: 1, maxLength: 500 }), { nil: undefined }),
  rowResults: fc.option(fc.array(importRowResultArb, { minLength: 0, maxLength: 10 }), { nil: undefined }),
});

describe('Import History Storage - Property Tests', () => {
  let storage: ImportHistoryStorage;
  let mockStorage: Storage;

  beforeEach(async () => {
    mockStorage = createMockStorage();
    storage = new ImportHistoryStorage(mockStorage);
    await storage.clear();
  });

  /**
   * Property 23: Import Metadata Recording
   * For any completed import operation, the metadata record SHALL contain
   * timestamp, user ID, file name, record counts, and status.
   *
   * **Feature: admin-phase3-improvements, Property 23: Import Metadata Recording**
   * **Validates: Requirements 9.1**
   */
  describe('Property 23: Import Metadata Recording', () => {
    it('should preserve all required metadata fields when saving a record', async () => {
      await fc.assert(
        fc.asyncProperty(importHistoryRecordArb, async (record) => {
          // Save the record
          await storage.save(record);

          // Retrieve the record
          const retrieved = await storage.getById(record.id);

          // Verify all required metadata fields are preserved
          expect(retrieved).not.toBeNull();
          if (retrieved) {
            // Timestamp (uploadedAt)
            expect(retrieved.uploadedAt).toBe(record.uploadedAt);
            
            // User ID
            expect(retrieved.uploadedBy).toBe(record.uploadedBy);
            
            // File name
            expect(retrieved.fileName).toBe(record.fileName);
            
            // Record counts
            expect(retrieved.totalRows).toBe(record.totalRows);
            expect(retrieved.importedCount).toBe(record.importedCount);
            expect(retrieved.skippedCount).toBe(record.skippedCount);
            
            // Status
            expect(retrieved.status).toBe(record.status);
            
            // Additional metadata
            expect(retrieved.type).toBe(record.type);
            expect(retrieved.fileSize).toBe(record.fileSize);
          }

          // Clean up for next iteration
          await storage.delete(record.id);
        }),
        { numRuns: 100 }
      );
    });

    it('should preserve optional metadata fields when present', async () => {
      await fc.assert(
        fc.asyncProperty(importHistoryRecordArb, async (record) => {
          await storage.save(record);
          const retrieved = await storage.getById(record.id);

          expect(retrieved).not.toBeNull();
          if (retrieved) {
            // Optional fields should be preserved if present
            expect(retrieved.uploadedByName).toBe(record.uploadedByName);
            expect(retrieved.completedAt).toBe(record.completedAt);
            expect(retrieved.errorSummary).toBe(record.errorSummary);
            
            // Row results should be preserved
            if (record.rowResults) {
              expect(retrieved.rowResults).toEqual(record.rowResults);
            }
          }

          await storage.delete(record.id);
        }),
        { numRuns: 100 }
      );
    });

    it('should create records with all required fields using helper function', () => {
      fc.assert(
        fc.property(
          importTypeArb,
          fc.string({ minLength: 1, maxLength: 100 }),
          fc.integer({ min: 1, max: 10 * 1024 * 1024 }),
          fc.integer({ min: 1, max: 10000 }),
          fc.option(fc.string({ minLength: 1, maxLength: 50 }), { nil: undefined }),
          fc.integer({ min: 0, max: 10000 }),
          (type, fileName, fileSize, uploadedBy, uploadedByName, totalRows) => {
            const record = createImportHistoryRecord(
              type,
              fileName,
              fileSize,
              uploadedBy,
              uploadedByName,
              totalRows
            );

            // Verify all required fields are present
            expect(record.id).toBeDefined();
            expect(record.id.length).toBeGreaterThan(0);
            expect(record.type).toBe(type);
            expect(record.fileName).toBe(fileName);
            expect(record.fileSize).toBe(fileSize);
            expect(record.uploadedBy).toBe(uploadedBy);
            expect(record.uploadedByName).toBe(uploadedByName);
            expect(record.uploadedAt).toBeDefined();
            expect(record.totalRows).toBe(totalRows);
            expect(record.importedCount).toBe(0);
            expect(record.skippedCount).toBe(0);
            expect(record.status).toBe('pending');
          }
        ),
        { numRuns: 100 }
      );
    });

    it('should generate unique IDs for each record', () => {
      fc.assert(
        fc.property(fc.integer({ min: 10, max: 100 }), (count) => {
          const ids = new Set<string>();
          for (let i = 0; i < count; i++) {
            ids.add(generateImportId());
          }
          // All generated IDs should be unique
          return ids.size === count;
        }),
        { numRuns: 50 }
      );
    });

    it('should update records while preserving unchanged fields', async () => {
      await fc.assert(
        fc.asyncProperty(
          importHistoryRecordArb,
          importStatusArb,
          fc.integer({ min: 0, max: 10000 }),
          fc.integer({ min: 0, max: 10000 }),
          async (record, newStatus, newImportedCount, newSkippedCount) => {
            await storage.save(record);

            // Update only specific fields
            await storage.update(record.id, {
              status: newStatus,
              importedCount: newImportedCount,
              skippedCount: newSkippedCount,
            });

            const retrieved = await storage.getById(record.id);

            expect(retrieved).not.toBeNull();
            if (retrieved) {
              // Updated fields should have new values
              expect(retrieved.status).toBe(newStatus);
              expect(retrieved.importedCount).toBe(newImportedCount);
              expect(retrieved.skippedCount).toBe(newSkippedCount);

              // Unchanged fields should be preserved
              expect(retrieved.id).toBe(record.id);
              expect(retrieved.type).toBe(record.type);
              expect(retrieved.fileName).toBe(record.fileName);
              expect(retrieved.fileSize).toBe(record.fileSize);
              expect(retrieved.uploadedBy).toBe(record.uploadedBy);
              expect(retrieved.uploadedAt).toBe(record.uploadedAt);
              expect(retrieved.totalRows).toBe(record.totalRows);
            }

            await storage.delete(record.id);
          }
        ),
        { numRuns: 100 }
      );
    });
  });
});


/**
 * Property 19: Error Detail Preservation
 * For any import row with validation errors, the error details SHALL include
 * the row number, field name, and specific error message, and these details
 * SHALL be preserved for later retrieval.
 *
 * **Feature: admin-phase3-improvements, Property 19: Error Detail Preservation**
 * **Validates: Requirements 6.4, 9.3, 9.4**
 */
describe('Property 19: Error Detail Preservation', () => {
  let storage: ImportHistoryStorage;
  let mockStorage: Storage;

  beforeEach(async () => {
    mockStorage = createMockStorage();
    storage = new ImportHistoryStorage(mockStorage);
    await storage.clear();
  });

  /**
   * Arbitrary for error field names
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
    fc.constant('Invalid format'),
    fc.string({ minLength: 5, maxLength: 200 })
  );

  /**
   * Arbitrary for ImportRowResult with errors
   */
  const errorRowResultArb: fc.Arbitrary<ImportRowResult> = fc.record({
    rowNumber: fc.integer({ min: 2, max: 10000 }), // Row 2+ (after header)
    success: fc.constant(false),
    originalData: fc.dictionary(
      fc.constantFrom('email', 'phone', 'name', 'role', 'status'),
      fc.string({ minLength: 1, maxLength: 100 })
    ),
    errorMessage: errorMessageArb,
    errorField: errorFieldArb,
    createdRecordId: fc.constant(undefined),
  });

  /**
   * Arbitrary for ImportRowResult with success
   */
  const successRowResultArb: fc.Arbitrary<ImportRowResult> = fc.record({
    rowNumber: fc.integer({ min: 2, max: 10000 }),
    success: fc.constant(true),
    originalData: fc.dictionary(
      fc.constantFrom('email', 'phone', 'name', 'role', 'status'),
      fc.string({ minLength: 1, maxLength: 100 })
    ),
    errorMessage: fc.constant(undefined),
    errorField: fc.constant(undefined),
    createdRecordId: fc.integer({ min: 1, max: 100000 }),
  });

  /**
   * Arbitrary for mixed row results (some success, some errors)
   */
  const mixedRowResultsArb = fc.array(
    fc.oneof(errorRowResultArb, successRowResultArb),
    { minLength: 1, maxLength: 20 }
  );

  it('should preserve row number for all error rows', async () => {
    await fc.assert(
      fc.asyncProperty(
        importTypeArb,
        fc.string({ minLength: 1, maxLength: 100 }),
        fc.integer({ min: 1, max: 10 * 1024 * 1024 }),
        fc.integer({ min: 1, max: 10000 }),
        mixedRowResultsArb,
        async (type, fileName, fileSize, uploadedBy, rowResults) => {
          const record = createImportHistoryRecord(type, fileName, fileSize, uploadedBy, undefined, rowResults.length);
          record.rowResults = rowResults;
          record.status = 'partial';
          record.skippedCount = rowResults.filter(r => !r.success).length;
          record.importedCount = rowResults.filter(r => r.success).length;

          await storage.save(record);
          const retrieved = await storage.getById(record.id);

          expect(retrieved).not.toBeNull();
          expect(retrieved!.rowResults).toBeDefined();
          expect(retrieved!.rowResults!.length).toBe(rowResults.length);

          // Verify all row numbers are preserved
          for (let i = 0; i < rowResults.length; i++) {
            expect(retrieved!.rowResults![i].rowNumber).toBe(rowResults[i].rowNumber);
          }

          await storage.delete(record.id);
        }
      ),
      { numRuns: 100 }
    );
  });

  it('should preserve error field name for all error rows', async () => {
    await fc.assert(
      fc.asyncProperty(
        importTypeArb,
        fc.array(errorRowResultArb, { minLength: 1, maxLength: 10 }),
        async (type, errorRows) => {
          const record = createImportHistoryRecord(type, 'test.csv', 1024, 1, undefined, errorRows.length);
          record.rowResults = errorRows;
          record.status = 'failed';
          record.skippedCount = errorRows.length;

          await storage.save(record);
          const retrieved = await storage.getById(record.id);

          expect(retrieved).not.toBeNull();
          expect(retrieved!.rowResults).toBeDefined();

          // Verify all error field names are preserved
          for (let i = 0; i < errorRows.length; i++) {
            expect(retrieved!.rowResults![i].errorField).toBe(errorRows[i].errorField);
          }

          await storage.delete(record.id);
        }
      ),
      { numRuns: 100 }
    );
  });

  it('should preserve specific error message for all error rows', async () => {
    await fc.assert(
      fc.asyncProperty(
        importTypeArb,
        fc.array(errorRowResultArb, { minLength: 1, maxLength: 10 }),
        async (type, errorRows) => {
          const record = createImportHistoryRecord(type, 'test.csv', 1024, 1, undefined, errorRows.length);
          record.rowResults = errorRows;
          record.status = 'failed';
          record.skippedCount = errorRows.length;

          await storage.save(record);
          const retrieved = await storage.getById(record.id);

          expect(retrieved).not.toBeNull();
          expect(retrieved!.rowResults).toBeDefined();

          // Verify all error messages are preserved
          for (let i = 0; i < errorRows.length; i++) {
            expect(retrieved!.rowResults![i].errorMessage).toBe(errorRows[i].errorMessage);
          }

          await storage.delete(record.id);
        }
      ),
      { numRuns: 100 }
    );
  });

  it('should preserve original data for error rows to enable error report generation', async () => {
    await fc.assert(
      fc.asyncProperty(
        importTypeArb,
        fc.array(errorRowResultArb, { minLength: 1, maxLength: 10 }),
        async (type, errorRows) => {
          const record = createImportHistoryRecord(type, 'test.csv', 1024, 1, undefined, errorRows.length);
          record.rowResults = errorRows;
          record.status = 'failed';
          record.skippedCount = errorRows.length;

          await storage.save(record);
          const retrieved = await storage.getById(record.id);

          expect(retrieved).not.toBeNull();
          expect(retrieved!.rowResults).toBeDefined();

          // Verify original data is preserved for each error row
          for (let i = 0; i < errorRows.length; i++) {
            expect(retrieved!.rowResults![i].originalData).toEqual(errorRows[i].originalData);
          }

          await storage.delete(record.id);
        }
      ),
      { numRuns: 100 }
    );
  });

  it('should preserve error details through update operations', async () => {
    await fc.assert(
      fc.asyncProperty(
        importTypeArb,
        fc.array(errorRowResultArb, { minLength: 1, maxLength: 10 }),
        importStatusArb,
        async (type, errorRows, newStatus) => {
          const record = createImportHistoryRecord(type, 'test.csv', 1024, 1, undefined, errorRows.length);
          record.rowResults = errorRows;
          record.status = 'processing';
          record.skippedCount = errorRows.length;

          await storage.save(record);

          // Update status but not rowResults
          await storage.update(record.id, { status: newStatus });

          const retrieved = await storage.getById(record.id);

          expect(retrieved).not.toBeNull();
          expect(retrieved!.status).toBe(newStatus);
          
          // Error details should still be preserved
          expect(retrieved!.rowResults).toBeDefined();
          expect(retrieved!.rowResults!.length).toBe(errorRows.length);
          
          for (let i = 0; i < errorRows.length; i++) {
            expect(retrieved!.rowResults![i].rowNumber).toBe(errorRows[i].rowNumber);
            expect(retrieved!.rowResults![i].errorField).toBe(errorRows[i].errorField);
            expect(retrieved!.rowResults![i].errorMessage).toBe(errorRows[i].errorMessage);
            expect(retrieved!.rowResults![i].originalData).toEqual(errorRows[i].originalData);
          }

          await storage.delete(record.id);
        }
      ),
      { numRuns: 100 }
    );
  });

  it('should preserve complete error details for later retrieval via query', async () => {
    await fc.assert(
      fc.asyncProperty(
        importTypeArb,
        fc.array(errorRowResultArb, { minLength: 1, maxLength: 5 }),
        async (type, errorRows) => {
          const record = createImportHistoryRecord(type, 'test.csv', 1024, 1, undefined, errorRows.length);
          record.rowResults = errorRows;
          record.status = 'failed';
          record.skippedCount = errorRows.length;

          await storage.save(record);

          // Query to retrieve the record
          const queryResult = await storage.query({ type, status: 'failed' });

          expect(queryResult.records.length).toBeGreaterThan(0);
          
          const foundRecord = queryResult.records.find(r => r.id === record.id);
          expect(foundRecord).toBeDefined();
          expect(foundRecord!.rowResults).toBeDefined();
          expect(foundRecord!.rowResults!.length).toBe(errorRows.length);

          // Verify all error details are preserved in query results
          for (let i = 0; i < errorRows.length; i++) {
            const retrievedRow = foundRecord!.rowResults![i];
            const originalRow = errorRows[i];
            
            // Row number
            expect(retrievedRow.rowNumber).toBe(originalRow.rowNumber);
            // Field name
            expect(retrievedRow.errorField).toBe(originalRow.errorField);
            // Error message
            expect(retrievedRow.errorMessage).toBe(originalRow.errorMessage);
            // Original data
            expect(retrievedRow.originalData).toEqual(originalRow.originalData);
          }

          await storage.delete(record.id);
        }
      ),
      { numRuns: 100 }
    );
  });

  it('should handle records with both success and error rows', async () => {
    await fc.assert(
      fc.asyncProperty(
        importTypeArb,
        fc.array(successRowResultArb, { minLength: 1, maxLength: 5 }),
        fc.array(errorRowResultArb, { minLength: 1, maxLength: 5 }),
        async (type, successRows, errorRows) => {
          const allRows = [...successRows, ...errorRows];
          const record = createImportHistoryRecord(type, 'test.csv', 1024, 1, undefined, allRows.length);
          record.rowResults = allRows;
          record.status = 'partial';
          record.importedCount = successRows.length;
          record.skippedCount = errorRows.length;

          await storage.save(record);
          const retrieved = await storage.getById(record.id);

          expect(retrieved).not.toBeNull();
          expect(retrieved!.rowResults).toBeDefined();
          expect(retrieved!.rowResults!.length).toBe(allRows.length);

          // Verify error rows have all required error details
          const retrievedErrorRows = retrieved!.rowResults!.filter(r => !r.success);
          expect(retrievedErrorRows.length).toBe(errorRows.length);

          for (const errorRow of retrievedErrorRows) {
            // Each error row must have row number
            expect(errorRow.rowNumber).toBeGreaterThanOrEqual(2);
            // Each error row must have error field
            expect(errorRow.errorField).toBeDefined();
            // Each error row must have error message
            expect(errorRow.errorMessage).toBeDefined();
            // Each error row must have original data
            expect(errorRow.originalData).toBeDefined();
          }

          await storage.delete(record.id);
        }
      ),
      { numRuns: 100 }
    );
  });
});
