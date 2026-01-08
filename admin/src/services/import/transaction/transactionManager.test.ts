/**
 * Property-Based Tests for Import Transaction Manager
 *
 * Tests for ImportTransactionManager implementation.
 *
 * @module services/import/transaction/transactionManager.test
 */

import { describe, it, expect, beforeEach } from 'vitest';
import * as fc from 'fast-check';
import {
  ImportTransactionManager,
  generateTransactionId,
  type TransactionManagerConfig,
} from './transactionManager';
import type { ImportType } from '../templates/types';
import type { ImportTransaction } from './types';

/**
 * Create an in-memory storage for testing
 */
function createTestStorage(): Storage {
  const store = new Map<string, string>();
  return {
    getItem: (key: string) => store.get(key) ?? null,
    setItem: (key: string, value: string) => {
      store.set(key, value);
    },
    removeItem: (key: string) => {
      store.delete(key);
    },
    clear: () => {
      store.clear();
    },
    get length() {
      return store.size;
    },
    key: (index: number) => Array.from(store.keys())[index] ?? null,
  };
}

/**
 * Create a test transaction manager with in-memory storage
 */
function createTestManager(config?: Partial<TransactionManagerConfig>): ImportTransactionManager {
  return new ImportTransactionManager({
    storage: createTestStorage(),
    ...config,
  });
}

describe('Import Transaction Tracking - Property Tests', () => {
  /**
   * **Feature: admin-phase3-improvements, Property 30: Import Transaction Tracking**
   * **Validates: Requirements 12.1**
   *
   * For any import operation, the transaction manager SHALL track all created
   * record IDs and persist them to storage.
   */
  describe('Property 30: Import Transaction Tracking', () => {
    it('should track all created record IDs for any import operation', () => {
      // Arbitrary for import types
      const importTypeArb = fc.constantFrom<ImportType>('user', 'player', 'game');

      // Arbitrary for record IDs (positive integers)
      const recordIdArb = fc.integer({ min: 1, max: 1000000 });

      // Arbitrary for number of records to create
      const numRecordsArb = fc.integer({ min: 1, max: 100 });

      fc.assert(
        fc.property(
          importTypeArb,
          numRecordsArb,
          fc.array(recordIdArb, { minLength: 1, maxLength: 100 }),
          (importType, _numRecords, recordIds) => {
            const manager = createTestManager();

            // Start a transaction
            const transaction = manager.startTransaction(importType);
            expect(transaction.id).toBeTruthy();
            expect(transaction.type).toBe(importType);
            expect(transaction.status).toBe('pending');
            expect(transaction.createdRecordIds).toHaveLength(0);

            // Record all created IDs
            for (const recordId of recordIds) {
              manager.recordCreated(transaction.id, recordId);
            }

            // Verify all IDs are tracked
            const updatedTransaction = manager.getTransaction(transaction.id);
            expect(updatedTransaction).toBeDefined();
            expect(updatedTransaction!.createdRecordIds).toHaveLength(recordIds.length);
            expect(updatedTransaction!.createdRecordIds).toEqual(recordIds);
            expect(updatedTransaction!.status).toBe('in_progress');
          }
        ),
        { numRuns: 100 }
      );
    });

    it('should persist transactions to storage', () => {
      // Arbitrary for import types
      const importTypeArb = fc.constantFrom<ImportType>('user', 'player', 'game');

      // Arbitrary for record IDs
      const recordIdsArb = fc.array(fc.integer({ min: 1, max: 1000000 }), {
        minLength: 1,
        maxLength: 50,
      });

      fc.assert(
        fc.property(importTypeArb, recordIdsArb, (importType, recordIds) => {
          const storage = createTestStorage();
          const manager = new ImportTransactionManager({ storage });

          // Start a transaction and record IDs
          const transaction = manager.startTransaction(importType);
          for (const recordId of recordIds) {
            manager.recordCreated(transaction.id, recordId);
          }

          // Verify storage contains the transaction
          const storedData = storage.getItem('import_transactions');
          expect(storedData).toBeTruthy();

          const parsed = JSON.parse(storedData!) as Array<[string, ImportTransaction]>;
          expect(parsed.length).toBeGreaterThan(0);

          const storedTransaction = parsed.find(([id]) => id === transaction.id);
          expect(storedTransaction).toBeDefined();
          expect(storedTransaction![1].createdRecordIds).toEqual(recordIds);
        }),
        { numRuns: 100 }
      );
    });

    it('should load persisted transactions on initialization', () => {
      // Arbitrary for import types
      const importTypeArb = fc.constantFrom<ImportType>('user', 'player', 'game');

      // Arbitrary for record IDs
      const recordIdsArb = fc.array(fc.integer({ min: 1, max: 1000000 }), {
        minLength: 1,
        maxLength: 50,
      });

      fc.assert(
        fc.property(importTypeArb, recordIdsArb, (importType, recordIds) => {
          const storage = createTestStorage();

          // Create first manager and add transaction
          const manager1 = new ImportTransactionManager({ storage });
          const transaction = manager1.startTransaction(importType);
          for (const recordId of recordIds) {
            manager1.recordCreated(transaction.id, recordId);
          }

          // Create second manager with same storage (simulates app restart)
          const manager2 = new ImportTransactionManager({ storage });

          // Verify transaction is loaded
          const loadedTransaction = manager2.getTransaction(transaction.id);
          expect(loadedTransaction).toBeDefined();
          expect(loadedTransaction!.type).toBe(importType);
          expect(loadedTransaction!.createdRecordIds).toEqual(recordIds);
        }),
        { numRuns: 100 }
      );
    });

    it('should include all required metadata in transaction', () => {
      // Arbitrary for import types
      const importTypeArb = fc.constantFrom<ImportType>('user', 'player', 'game');

      // Arbitrary for total rows
      const totalRowsArb = fc.integer({ min: 1, max: 10000 });

      fc.assert(
        fc.property(importTypeArb, totalRowsArb, (importType, totalRows) => {
          const manager = createTestManager();

          // Start a transaction with total rows
          const transaction = manager.startTransaction(importType, totalRows);

          // Verify all required metadata is present
          expect(transaction.id).toBeTruthy();
          expect(typeof transaction.id).toBe('string');
          expect(transaction.type).toBe(importType);
          expect(transaction.status).toBe('pending');
          expect(transaction.createdRecordIds).toEqual([]);
          expect(transaction.startedAt).toBeTruthy();
          expect(new Date(transaction.startedAt).getTime()).not.toBeNaN();
          expect(transaction.totalRows).toBe(totalRows);
          expect(transaction.processedRows).toBe(0);
        }),
        { numRuns: 100 }
      );
    });

    it('should generate unique transaction IDs', () => {
      // Generate multiple IDs and verify uniqueness
      const numIdsArb = fc.integer({ min: 10, max: 100 });

      fc.assert(
        fc.property(numIdsArb, (numIds) => {
          const ids = new Set<string>();
          for (let i = 0; i < numIds; i++) {
            ids.add(generateTransactionId());
          }
          // All IDs should be unique
          expect(ids.size).toBe(numIds);
        }),
        { numRuns: 50 }
      );
    });
  });

  /**
   * **Feature: admin-phase3-improvements, Property 31: Rollback Completeness**
   * **Validates: Requirements 12.3, 12.4**
   *
   * For any rollback operation, the transaction manager SHALL attempt to delete
   * all records created during the import and report any failures.
   */
  describe('Property 31: Rollback Completeness', () => {
    it('should attempt to delete all records created during import', async () => {
      // Arbitrary for import types
      const importTypeArb = fc.constantFrom<ImportType>('user', 'player', 'game');

      // Arbitrary for unique record IDs (using Set to ensure uniqueness)
      const recordIdsArb = fc.array(fc.integer({ min: 1, max: 1000000 }), {
        minLength: 1,
        maxLength: 50,
      }).map(ids => [...new Set(ids)]); // Ensure unique IDs

      await fc.assert(
        fc.asyncProperty(importTypeArb, recordIdsArb, async (importType, recordIds) => {
          const deletedRecords: number[] = [];
          const deleteRecord = async (_type: ImportType, recordId: number) => {
            deletedRecords.push(recordId);
          };

          const manager = new ImportTransactionManager({
            storage: createTestStorage(),
            deleteRecord,
          });

          // Start transaction and record all IDs
          const transaction = manager.startTransaction(importType);
          for (const recordId of recordIds) {
            manager.recordCreated(transaction.id, recordId);
          }

          // Perform rollback
          const result = await manager.rollbackTransaction(transaction.id);

          // Verify all records were attempted for deletion
          expect(deletedRecords.length).toBe(recordIds.length);
          expect(result.rolledBackCount).toBe(recordIds.length);
          expect(result.success).toBe(true);
          expect(result.failedRollbacks).toHaveLength(0);

          // Verify records were deleted in reverse order (LIFO)
          const expectedOrder = [...recordIds].reverse();
          expect(deletedRecords).toEqual(expectedOrder);
        }),
        { numRuns: 100 }
      );
    });

    it('should report all failures during rollback', async () => {
      // Arbitrary for import types
      const importTypeArb = fc.constantFrom<ImportType>('user', 'player', 'game');

      // Arbitrary for record IDs with some that will fail
      const recordIdsArb = fc.array(fc.integer({ min: 1, max: 100 }), {
        minLength: 3,
        maxLength: 20,
      }).map(ids => [...new Set(ids)]); // Ensure unique IDs

      // Arbitrary for which indices should fail (subset of indices)
      const failIndicesArb = fc.array(fc.integer({ min: 0, max: 19 }), {
        minLength: 1,
        maxLength: 5,
      }).map(indices => [...new Set(indices)]); // Ensure unique indices

      await fc.assert(
        fc.asyncProperty(
          importTypeArb,
          recordIdsArb,
          failIndicesArb,
          async (importType, recordIds, failIndices) => {
            // Filter fail indices to be within bounds
            const validFailIndices = failIndices.filter(i => i < recordIds.length);
            if (validFailIndices.length === 0) return; // Skip if no valid fail indices

            const failingIds = new Set(validFailIndices.map(i => recordIds[i]));
            const deletedRecords: number[] = [];

            const deleteRecord = async (_type: ImportType, recordId: number) => {
              if (failingIds.has(recordId)) {
                throw new Error(`Failed to delete record ${recordId}`);
              }
              deletedRecords.push(recordId);
            };

            const manager = new ImportTransactionManager({
              storage: createTestStorage(),
              deleteRecord,
            });

            // Start transaction and record all IDs
            const transaction = manager.startTransaction(importType);
            for (const recordId of recordIds) {
              manager.recordCreated(transaction.id, recordId);
            }

            // Perform rollback
            const result = await manager.rollbackTransaction(transaction.id);

            // Verify failure reporting
            expect(result.success).toBe(false);
            expect(result.failedRollbacks.length).toBe(failingIds.size);

            // Verify all failed IDs are reported
            const reportedFailedIds = new Set(result.failedRollbacks.map(f => f.recordId));
            for (const failId of failingIds) {
              expect(reportedFailedIds.has(failId)).toBe(true);
            }

            // Verify each failure has an error message
            for (const failure of result.failedRollbacks) {
              expect(failure.error).toBeTruthy();
              expect(typeof failure.error).toBe('string');
            }

            // Verify successful + failed = total
            expect(result.rolledBackCount + result.failedRollbacks.length).toBe(recordIds.length);
          }
        ),
        { numRuns: 100 }
      );
    });

    it('should mark transaction as rolled_back after rollback attempt', async () => {
      // Arbitrary for import types
      const importTypeArb = fc.constantFrom<ImportType>('user', 'player', 'game');

      // Arbitrary for record IDs
      const recordIdsArb = fc.array(fc.integer({ min: 1, max: 1000000 }), {
        minLength: 1,
        maxLength: 30,
      });

      // Arbitrary for whether some deletes should fail
      const shouldFailArb = fc.boolean();

      await fc.assert(
        fc.asyncProperty(
          importTypeArb,
          recordIdsArb,
          shouldFailArb,
          async (importType, recordIds, shouldFail) => {
            const deleteRecord = async (_type: ImportType, recordId: number) => {
              if (shouldFail && recordId === recordIds[0]) {
                throw new Error('Simulated failure');
              }
            };

            const manager = new ImportTransactionManager({
              storage: createTestStorage(),
              deleteRecord,
            });

            // Start transaction and record all IDs
            const transaction = manager.startTransaction(importType);
            for (const recordId of recordIds) {
              manager.recordCreated(transaction.id, recordId);
            }

            // Perform rollback
            await manager.rollbackTransaction(transaction.id);

            // Verify transaction status is rolled_back regardless of failures
            const updated = manager.getTransaction(transaction.id);
            expect(updated?.status).toBe('rolled_back');
            expect(updated?.completedAt).toBeTruthy();
          }
        ),
        { numRuns: 100 }
      );
    });

    it('should continue rollback even when some deletions fail', async () => {
      // Arbitrary for import types
      const importTypeArb = fc.constantFrom<ImportType>('user', 'player', 'game');

      // Arbitrary for record IDs (at least 5 to have meaningful test)
      const recordIdsArb = fc.array(fc.integer({ min: 1, max: 1000 }), {
        minLength: 5,
        maxLength: 20,
      }).map(ids => [...new Set(ids)]).filter(ids => ids.length >= 5);

      await fc.assert(
        fc.asyncProperty(importTypeArb, recordIdsArb, async (importType, recordIds) => {
          // Make the middle record fail
          const failingId = recordIds[Math.floor(recordIds.length / 2)];
          const deletedRecords: number[] = [];

          const deleteRecord = async (_type: ImportType, recordId: number) => {
            if (recordId === failingId) {
              throw new Error('Simulated failure');
            }
            deletedRecords.push(recordId);
          };

          const manager = new ImportTransactionManager({
            storage: createTestStorage(),
            deleteRecord,
          });

          // Start transaction and record all IDs
          const transaction = manager.startTransaction(importType);
          for (const recordId of recordIds) {
            manager.recordCreated(transaction.id, recordId);
          }

          // Perform rollback
          const result = await manager.rollbackTransaction(transaction.id);

          // Verify rollback continued past the failure
          // All records except the failing one should be deleted
          expect(deletedRecords.length).toBe(recordIds.length - 1);
          expect(result.rolledBackCount).toBe(recordIds.length - 1);
          expect(result.failedRollbacks.length).toBe(1);
          expect(result.failedRollbacks[0].recordId).toBe(failingId);
        }),
        { numRuns: 100 }
      );
    });
  });

  /**
   * **Feature: admin-phase3-improvements, Property 32: Interrupted Import Detection**
   * **Validates: Requirements 12.5**
   *
   * For any import transaction in 'in_progress' or 'pending' status when the
   * application loads, the system SHALL mark it as 'interrupted' and make it
   * available for cleanup.
   */
  describe('Property 32: Interrupted Import Detection', () => {
    it('should mark in_progress transactions as interrupted on cleanup', async () => {
      // Arbitrary for import types
      const importTypeArb = fc.constantFrom<ImportType>('user', 'player', 'game');

      // Arbitrary for record IDs
      const recordIdsArb = fc.array(fc.integer({ min: 1, max: 1000000 }), {
        minLength: 1,
        maxLength: 30,
      });

      await fc.assert(
        fc.asyncProperty(importTypeArb, recordIdsArb, async (importType, recordIds) => {
          const storage = createTestStorage();

          // Create manager and start a transaction (simulating an import in progress)
          const manager1 = new ImportTransactionManager({ storage });
          const transaction = manager1.startTransaction(importType);

          // Record some IDs to put it in 'in_progress' status
          for (const recordId of recordIds) {
            manager1.recordCreated(transaction.id, recordId);
          }

          // Verify it's in_progress
          expect(manager1.getTransaction(transaction.id)?.status).toBe('in_progress');

          // Simulate app restart - create new manager with same storage
          const manager2 = new ImportTransactionManager({ storage });

          // Call cleanupInterrupted (simulating startup detection)
          await manager2.cleanupInterrupted();

          // Verify transaction is now marked as interrupted
          const updated = manager2.getTransaction(transaction.id);
          expect(updated?.status).toBe('interrupted');

          // Verify it's available via getInterruptedTransactions
          const interrupted = manager2.getInterruptedTransactions();
          expect(interrupted.some(t => t.id === transaction.id)).toBe(true);
        }),
        { numRuns: 100 }
      );
    });

    it('should mark pending transactions as interrupted on cleanup', async () => {
      // Arbitrary for import types
      const importTypeArb = fc.constantFrom<ImportType>('user', 'player', 'game');

      // Arbitrary for total rows
      const totalRowsArb = fc.integer({ min: 1, max: 1000 });

      await fc.assert(
        fc.asyncProperty(importTypeArb, totalRowsArb, async (importType, totalRows) => {
          const storage = createTestStorage();

          // Create manager and start a transaction (but don't record any IDs - stays pending)
          const manager1 = new ImportTransactionManager({ storage });
          const transaction = manager1.startTransaction(importType, totalRows);

          // Verify it's pending
          expect(manager1.getTransaction(transaction.id)?.status).toBe('pending');

          // Simulate app restart
          const manager2 = new ImportTransactionManager({ storage });

          // Call cleanupInterrupted
          await manager2.cleanupInterrupted();

          // Verify transaction is now marked as interrupted
          const updated = manager2.getTransaction(transaction.id);
          expect(updated?.status).toBe('interrupted');

          // Verify it's available via getInterruptedTransactions
          const interrupted = manager2.getInterruptedTransactions();
          expect(interrupted.some(t => t.id === transaction.id)).toBe(true);
        }),
        { numRuns: 100 }
      );
    });

    it('should not mark completed transactions as interrupted', async () => {
      // Arbitrary for import types
      const importTypeArb = fc.constantFrom<ImportType>('user', 'player', 'game');

      // Arbitrary for record IDs
      const recordIdsArb = fc.array(fc.integer({ min: 1, max: 1000000 }), {
        minLength: 1,
        maxLength: 20,
      });

      await fc.assert(
        fc.asyncProperty(importTypeArb, recordIdsArb, async (importType, recordIds) => {
          const storage = createTestStorage();

          // Create manager and complete a transaction
          const manager1 = new ImportTransactionManager({ storage });
          const transaction = manager1.startTransaction(importType);

          for (const recordId of recordIds) {
            manager1.recordCreated(transaction.id, recordId);
          }

          // Commit the transaction
          await manager1.commitTransaction(transaction.id);

          // Verify it's completed
          expect(manager1.getTransaction(transaction.id)?.status).toBe('completed');

          // Simulate app restart
          const manager2 = new ImportTransactionManager({ storage });

          // Call cleanupInterrupted
          await manager2.cleanupInterrupted();

          // Verify transaction is still completed (not interrupted)
          const updated = manager2.getTransaction(transaction.id);
          expect(updated?.status).toBe('completed');

          // Verify it's NOT in interrupted transactions
          const interrupted = manager2.getInterruptedTransactions();
          expect(interrupted.some(t => t.id === transaction.id)).toBe(false);
        }),
        { numRuns: 100 }
      );
    });

    it('should not mark rolled_back transactions as interrupted', async () => {
      // Arbitrary for import types
      const importTypeArb = fc.constantFrom<ImportType>('user', 'player', 'game');

      // Arbitrary for record IDs
      const recordIdsArb = fc.array(fc.integer({ min: 1, max: 1000000 }), {
        minLength: 1,
        maxLength: 20,
      });

      await fc.assert(
        fc.asyncProperty(importTypeArb, recordIdsArb, async (importType, recordIds) => {
          const storage = createTestStorage();
          const deleteRecord = async () => {}; // No-op delete

          // Create manager and rollback a transaction
          const manager1 = new ImportTransactionManager({ storage, deleteRecord });
          const transaction = manager1.startTransaction(importType);

          for (const recordId of recordIds) {
            manager1.recordCreated(transaction.id, recordId);
          }

          // Rollback the transaction
          await manager1.rollbackTransaction(transaction.id);

          // Verify it's rolled_back
          expect(manager1.getTransaction(transaction.id)?.status).toBe('rolled_back');

          // Simulate app restart
          const manager2 = new ImportTransactionManager({ storage });

          // Call cleanupInterrupted
          await manager2.cleanupInterrupted();

          // Verify transaction is still rolled_back (not interrupted)
          const updated = manager2.getTransaction(transaction.id);
          expect(updated?.status).toBe('rolled_back');

          // Verify it's NOT in interrupted transactions
          const interrupted = manager2.getInterruptedTransactions();
          expect(interrupted.some(t => t.id === transaction.id)).toBe(false);
        }),
        { numRuns: 100 }
      );
    });

    it('should preserve transaction data when marking as interrupted', async () => {
      // Arbitrary for import types
      const importTypeArb = fc.constantFrom<ImportType>('user', 'player', 'game');

      // Arbitrary for record IDs
      const recordIdsArb = fc.array(fc.integer({ min: 1, max: 1000000 }), {
        minLength: 1,
        maxLength: 30,
      });

      // Arbitrary for total rows
      const totalRowsArb = fc.integer({ min: 1, max: 1000 });

      await fc.assert(
        fc.asyncProperty(
          importTypeArb,
          recordIdsArb,
          totalRowsArb,
          async (importType, recordIds, totalRows) => {
            const storage = createTestStorage();

            // Create manager and start a transaction
            const manager1 = new ImportTransactionManager({ storage });
            const transaction = manager1.startTransaction(importType, totalRows);

            // Record IDs and update progress
            for (let i = 0; i < recordIds.length; i++) {
              manager1.recordCreated(transaction.id, recordIds[i]);
              manager1.updateProgress(transaction.id, i + 1);
            }

            // Capture original data
            const original = manager1.getTransaction(transaction.id)!;
            const originalRecordIds = [...original.createdRecordIds];
            const originalType = original.type;
            const originalTotalRows = original.totalRows;
            const originalProcessedRows = original.processedRows;
            const originalStartedAt = original.startedAt;

            // Simulate app restart
            const manager2 = new ImportTransactionManager({ storage });
            await manager2.cleanupInterrupted();

            // Verify all data is preserved except status
            const updated = manager2.getTransaction(transaction.id)!;
            expect(updated.status).toBe('interrupted');
            expect(updated.type).toBe(originalType);
            expect(updated.createdRecordIds).toEqual(originalRecordIds);
            expect(updated.totalRows).toBe(originalTotalRows);
            expect(updated.processedRows).toBe(originalProcessedRows);
            expect(updated.startedAt).toBe(originalStartedAt);
          }
        ),
        { numRuns: 100 }
      );
    });
  });
});

describe('Import Transaction Manager - Unit Tests', () => {
  let manager: ImportTransactionManager;
  let storage: Storage;

  beforeEach(() => {
    storage = createTestStorage();
    manager = new ImportTransactionManager({ storage });
  });

  describe('startTransaction', () => {
    it('should create a new transaction with pending status', () => {
      const transaction = manager.startTransaction('user');

      expect(transaction.id).toBeTruthy();
      expect(transaction.type).toBe('user');
      expect(transaction.status).toBe('pending');
      expect(transaction.createdRecordIds).toEqual([]);
      expect(transaction.startedAt).toBeTruthy();
    });

    it('should store total rows if provided', () => {
      const transaction = manager.startTransaction('player', 100);

      expect(transaction.totalRows).toBe(100);
      expect(transaction.processedRows).toBe(0);
    });
  });

  describe('recordCreated', () => {
    it('should add record ID to transaction', () => {
      const transaction = manager.startTransaction('user');
      manager.recordCreated(transaction.id, 123);

      const updated = manager.getTransaction(transaction.id);
      expect(updated?.createdRecordIds).toContain(123);
    });

    it('should update status to in_progress', () => {
      const transaction = manager.startTransaction('user');
      expect(transaction.status).toBe('pending');

      manager.recordCreated(transaction.id, 123);

      const updated = manager.getTransaction(transaction.id);
      expect(updated?.status).toBe('in_progress');
    });

    it('should handle non-existent transaction gracefully', () => {
      // Should not throw
      expect(() => manager.recordCreated('non-existent', 123)).not.toThrow();
    });
  });

  describe('updateProgress', () => {
    it('should update processed rows count', () => {
      const transaction = manager.startTransaction('user', 100);
      manager.updateProgress(transaction.id, 50);

      const updated = manager.getTransaction(transaction.id);
      expect(updated?.processedRows).toBe(50);
    });
  });

  describe('commitTransaction', () => {
    it('should mark transaction as completed', async () => {
      const transaction = manager.startTransaction('user');
      manager.recordCreated(transaction.id, 123);

      await manager.commitTransaction(transaction.id);

      const updated = manager.getTransaction(transaction.id);
      expect(updated?.status).toBe('completed');
      expect(updated?.completedAt).toBeTruthy();
    });
  });

  describe('getTransaction', () => {
    it('should return undefined for non-existent transaction', () => {
      const result = manager.getTransaction('non-existent');
      expect(result).toBeUndefined();
    });

    it('should return transaction by ID', () => {
      const transaction = manager.startTransaction('game');
      const result = manager.getTransaction(transaction.id);

      expect(result).toBeDefined();
      expect(result?.id).toBe(transaction.id);
    });
  });

  describe('getAllTransactions', () => {
    it('should return all transactions', () => {
      manager.startTransaction('user');
      manager.startTransaction('player');
      manager.startTransaction('game');

      const all = manager.getAllTransactions();
      expect(all).toHaveLength(3);
    });
  });

  describe('deleteTransaction', () => {
    it('should remove transaction from storage', () => {
      const transaction = manager.startTransaction('user');
      expect(manager.getTransaction(transaction.id)).toBeDefined();

      manager.deleteTransaction(transaction.id);
      expect(manager.getTransaction(transaction.id)).toBeUndefined();
    });
  });

  describe('rollbackTransaction', () => {
    it('should rollback all created records in reverse order', async () => {
      const deletedRecords: number[] = [];
      const deleteRecord = async (_type: ImportType, recordId: number) => {
        deletedRecords.push(recordId);
      };

      const managerWithDelete = new ImportTransactionManager({
        storage: createTestStorage(),
        deleteRecord,
      });

      const transaction = managerWithDelete.startTransaction('user');
      managerWithDelete.recordCreated(transaction.id, 1);
      managerWithDelete.recordCreated(transaction.id, 2);
      managerWithDelete.recordCreated(transaction.id, 3);

      const result = await managerWithDelete.rollbackTransaction(transaction.id);

      expect(result.success).toBe(true);
      expect(result.rolledBackCount).toBe(3);
      expect(result.failedRollbacks).toHaveLength(0);
      // Should be deleted in reverse order (LIFO)
      expect(deletedRecords).toEqual([3, 2, 1]);
    });

    it('should mark transaction as rolled_back after rollback', async () => {
      const deleteRecord = async () => {};
      const managerWithDelete = new ImportTransactionManager({
        storage: createTestStorage(),
        deleteRecord,
      });

      const transaction = managerWithDelete.startTransaction('user');
      managerWithDelete.recordCreated(transaction.id, 1);

      await managerWithDelete.rollbackTransaction(transaction.id);

      const updated = managerWithDelete.getTransaction(transaction.id);
      expect(updated?.status).toBe('rolled_back');
      expect(updated?.completedAt).toBeTruthy();
    });

    it('should handle delete failures gracefully', async () => {
      const deleteRecord = async (_type: ImportType, recordId: number) => {
        if (recordId === 2) {
          throw new Error('Delete failed');
        }
      };

      const managerWithDelete = new ImportTransactionManager({
        storage: createTestStorage(),
        deleteRecord,
      });

      const transaction = managerWithDelete.startTransaction('user');
      managerWithDelete.recordCreated(transaction.id, 1);
      managerWithDelete.recordCreated(transaction.id, 2);
      managerWithDelete.recordCreated(transaction.id, 3);

      const result = await managerWithDelete.rollbackTransaction(transaction.id);

      expect(result.success).toBe(false);
      expect(result.rolledBackCount).toBe(2); // 3 and 1 succeeded
      expect(result.failedRollbacks).toHaveLength(1);
      expect(result.failedRollbacks[0].recordId).toBe(2);
      expect(result.failedRollbacks[0].error).toBe('Delete failed');
    });

    it('should return failure result for non-existent transaction', async () => {
      const result = await manager.rollbackTransaction('non-existent');

      expect(result.success).toBe(false);
      expect(result.rolledBackCount).toBe(0);
      expect(result.failedRollbacks).toHaveLength(0);
    });

    it('should log all rollback operations', async () => {
      const logEntries: Array<{ message: string; context?: Record<string, unknown> }> = [];
      const mockLogger = {
        debug: (msg: string, ctx?: Record<string, unknown>) => logEntries.push({ message: msg, context: ctx }),
        info: (msg: string, ctx?: Record<string, unknown>) => logEntries.push({ message: msg, context: ctx }),
        warn: (msg: string, ctx?: Record<string, unknown>) => logEntries.push({ message: msg, context: ctx }),
        error: (msg: string, _err?: Error, ctx?: Record<string, unknown>) => logEntries.push({ message: msg, context: ctx }),
        getServiceName: () => 'TestLogger',
        getLogEntries: () => [],
        clearLogEntries: () => {},
      };

      const deleteRecord = async () => {};
      const managerWithLogger = new ImportTransactionManager({
        storage: createTestStorage(),
        deleteRecord,
        logger: mockLogger,
      });

      const transaction = managerWithLogger.startTransaction('user');
      managerWithLogger.recordCreated(transaction.id, 1);
      managerWithLogger.recordCreated(transaction.id, 2);

      await managerWithLogger.rollbackTransaction(transaction.id);

      // Should have logged: start, record created x2, rollback start, record rolled back x2, rollback completed
      const rollbackLogs = logEntries.filter(
        (e) => e.message.toLowerCase().includes('rollback') || e.message.toLowerCase().includes('rolled back')
      );
      expect(rollbackLogs.length).toBeGreaterThanOrEqual(4); // Start, 2 records rolled back, completed
    });
  });
});
