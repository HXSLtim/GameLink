/**
 * Import Transaction Manager
 * Manages import transactions with persistence and rollback support
 *
 * @module services/import/transaction/transactionManager
 */

import type { ImportType } from '../templates/types';
import type {
  ImportTransaction,
  RollbackResult,
  ITransactionManager,
  DeleteRecordFn,
} from './types';
import type { ServiceLogger } from '../../utils/logger';
import { createServiceLogger } from '../../utils/logger';

/**
 * Storage key for transactions in localStorage
 */
const STORAGE_KEY = 'import_transactions';

/**
 * Generate a unique transaction ID
 */
export function generateTransactionId(): string {
  if (typeof crypto !== 'undefined' && crypto.randomUUID) {
    return crypto.randomUUID();
  }
  // Fallback: timestamp + random string
  const timestamp = Date.now().toString(36);
  const random = Math.random().toString(36).substring(2, 10);
  return `txn-${timestamp}-${random}`;
}

/**
 * Configuration for ImportTransactionManager
 */
export interface TransactionManagerConfig {
  /** Storage instance (defaults to localStorage) */
  storage?: Storage;
  /** Logger instance */
  logger?: ServiceLogger;
  /** Function to delete records by type and ID */
  deleteRecord?: DeleteRecordFn;
  /** Time to keep completed transactions before cleanup (ms) */
  completedRetentionMs?: number;
}

/**
 * Create an in-memory storage for testing or SSR environments
 */
function createMemoryStorage(): Storage {
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
 * Import Transaction Manager Implementation
 *
 * Features:
 * - Tracks all records created during an import
 * - Persists transactions to localStorage
 * - Supports rollback of created records
 * - Detects and handles interrupted imports
 */
export class ImportTransactionManager implements ITransactionManager {
  private transactions = new Map<string, ImportTransaction>();
  private storage: Storage;
  private logger: ServiceLogger;
  private deleteRecord: DeleteRecordFn;
  private completedRetentionMs: number;

  constructor(config: TransactionManagerConfig = {}) {
    this.storage =
      config.storage ??
      (typeof localStorage !== 'undefined' ? localStorage : createMemoryStorage());
    this.logger = config.logger ?? createServiceLogger('ImportTransactionManager');
    this.deleteRecord = config.deleteRecord ?? this.defaultDeleteRecord.bind(this);
    this.completedRetentionMs = config.completedRetentionMs ?? 60000; // 1 minute default

    // Load persisted transactions on initialization
    this.loadPersistedTransactions();
  }

  /**
   * Start a new import transaction
   */
  startTransaction(type: ImportType, totalRows?: number): ImportTransaction {
    const transaction: ImportTransaction = {
      id: generateTransactionId(),
      type,
      status: 'pending',
      createdRecordIds: [],
      startedAt: new Date().toISOString(),
      totalRows,
      processedRows: 0,
    };

    this.transactions.set(transaction.id, transaction);
    this.persistTransactions();

    this.logger.info('Import transaction started', {
      transactionId: transaction.id,
      type,
      totalRows,
    });

    return transaction;
  }

  /**
   * Record a created record ID in the transaction
   */
  recordCreated(transactionId: string, recordId: number): void {
    const transaction = this.transactions.get(transactionId);
    if (transaction) {
      transaction.createdRecordIds.push(recordId);
      transaction.status = 'in_progress';
      this.persistTransactions();

      this.logger.debug('Record created in transaction', {
        transactionId,
        recordId,
        totalCreated: transaction.createdRecordIds.length,
      });
    }
  }

  /**
   * Update transaction progress
   */
  updateProgress(transactionId: string, processedRows: number): void {
    const transaction = this.transactions.get(transactionId);
    if (transaction) {
      transaction.processedRows = processedRows;
      this.persistTransactions();
    }
  }

  /**
   * Commit a transaction (mark as completed)
   */
  async commitTransaction(transactionId: string): Promise<void> {
    const transaction = this.transactions.get(transactionId);
    if (transaction) {
      transaction.status = 'completed';
      transaction.completedAt = new Date().toISOString();
      this.persistTransactions();

      this.logger.info('Import transaction committed', {
        transactionId,
        recordCount: transaction.createdRecordIds.length,
      });

      // Schedule cleanup after retention period
      setTimeout(() => {
        this.transactions.delete(transactionId);
        this.persistTransactions();
      }, this.completedRetentionMs);
    }
  }

  /**
   * Rollback a transaction (delete all created records)
   */
  async rollbackTransaction(transactionId: string): Promise<RollbackResult> {
    const transaction = this.transactions.get(transactionId);
    if (!transaction) {
      return { success: false, rolledBackCount: 0, failedRollbacks: [] };
    }

    this.logger.info('Starting rollback', {
      transactionId,
      recordCount: transaction.createdRecordIds.length,
      type: transaction.type,
    });

    const result: RollbackResult = {
      success: true,
      rolledBackCount: 0,
      failedRollbacks: [],
    };

    // Rollback in reverse order (LIFO)
    const recordIds = [...transaction.createdRecordIds].reverse();

    for (const recordId of recordIds) {
      try {
        await this.deleteRecord(transaction.type, recordId);
        result.rolledBackCount++;

        this.logger.debug('Record rolled back', {
          transactionId,
          recordId,
          type: transaction.type,
        });
      } catch (error) {
        result.success = false;
        result.failedRollbacks.push({
          recordId,
          error: error instanceof Error ? error.message : 'Unknown error',
        });

        this.logger.error('Rollback failed for record', error as Error, {
          transactionId,
          recordId,
          type: transaction.type,
        });
      }
    }

    // Update transaction status
    transaction.status = 'rolled_back';
    transaction.completedAt = new Date().toISOString();
    this.persistTransactions();

    this.logger.info('Rollback completed', {
      transactionId,
      rolledBackCount: result.rolledBackCount,
      failedCount: result.failedRollbacks.length,
    });

    return result;
  }

  /**
   * Get a transaction by ID
   */
  getTransaction(transactionId: string): ImportTransaction | undefined {
    return this.transactions.get(transactionId);
  }

  /**
   * Get all transactions
   */
  getAllTransactions(): ImportTransaction[] {
    return Array.from(this.transactions.values());
  }

  /**
   * Get interrupted transactions (pending or in_progress)
   */
  getInterruptedTransactions(): ImportTransaction[] {
    return Array.from(this.transactions.values()).filter(
      (t) => t.status === 'interrupted'
    );
  }

  /**
   * Mark pending/in_progress transactions as interrupted
   */
  async cleanupInterrupted(): Promise<void> {
    const activeTransactions = Array.from(this.transactions.values()).filter(
      (t) => t.status === 'in_progress' || t.status === 'pending'
    );

    for (const transaction of activeTransactions) {
      transaction.status = 'interrupted';
      this.logger.warn('Found interrupted transaction', {
        transactionId: transaction.id,
        type: transaction.type,
        recordCount: transaction.createdRecordIds.length,
        processedRows: transaction.processedRows,
        totalRows: transaction.totalRows,
      });
    }

    if (activeTransactions.length > 0) {
      this.persistTransactions();
    }
  }

  /**
   * Delete a transaction record
   */
  deleteTransaction(transactionId: string): void {
    this.transactions.delete(transactionId);
    this.persistTransactions();
  }

  /**
   * Persist transactions to storage
   */
  private persistTransactions(): void {
    try {
      const data = Array.from(this.transactions.entries());
      this.storage.setItem(STORAGE_KEY, JSON.stringify(data));
    } catch (error) {
      this.logger.error('Failed to persist transactions', error as Error);
    }
  }

  /**
   * Load persisted transactions from storage
   */
  private loadPersistedTransactions(): void {
    try {
      const data = this.storage.getItem(STORAGE_KEY);
      if (data) {
        const entries = JSON.parse(data) as Array<[string, ImportTransaction]>;
        this.transactions = new Map(entries);
        this.logger.debug('Loaded persisted transactions', {
          count: this.transactions.size,
        });
      }
    } catch (error) {
      this.logger.error('Failed to load persisted transactions', error as Error);
      // Start fresh if loading fails
      this.transactions = new Map();
    }
  }

  /**
   * Default delete record function (no-op, should be overridden)
   */
  private async defaultDeleteRecord(type: ImportType, recordId: number): Promise<void> {
    this.logger.warn('Default deleteRecord called - no actual deletion performed', {
      type,
      recordId,
    });
    // This is a no-op; actual implementation should be provided via config
  }
}

/**
 * Create an ImportTransactionManager instance
 */
export function createTransactionManager(
  config?: TransactionManagerConfig
): ImportTransactionManager {
  return new ImportTransactionManager(config);
}
