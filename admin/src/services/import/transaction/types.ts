/**
 * Import Transaction Types
 * Defines types for import transaction management and rollback
 *
 * @module services/import/transaction/types
 */

import type { ImportType } from '../templates/types';

/**
 * Transaction status enum
 */
export type TransactionStatus =
  | 'pending'
  | 'in_progress'
  | 'completed'
  | 'rolled_back'
  | 'interrupted';

/**
 * Import transaction record
 */
export interface ImportTransaction {
  /** Unique identifier for this transaction */
  id: string;
  /** Type of import (user, player, game) */
  type: ImportType;
  /** Current status of the transaction */
  status: TransactionStatus;
  /** IDs of records created during this transaction */
  createdRecordIds: number[];
  /** Timestamp when transaction started */
  startedAt: string;
  /** Timestamp when transaction completed/rolled back */
  completedAt?: string;
  /** Total rows being imported */
  totalRows?: number;
  /** Number of rows processed so far */
  processedRows?: number;
}

/**
 * Result of a rollback operation
 */
export interface RollbackResult {
  /** Whether all rollbacks succeeded */
  success: boolean;
  /** Number of records successfully rolled back */
  rolledBackCount: number;
  /** Records that failed to rollback */
  failedRollbacks: Array<{
    recordId: number;
    error: string;
  }>;
}

/**
 * Transaction manager interface
 */
export interface ITransactionManager {
  /**
   * Start a new import transaction
   */
  startTransaction(type: ImportType, totalRows?: number): ImportTransaction;

  /**
   * Record a created record ID in the transaction
   */
  recordCreated(transactionId: string, recordId: number): void;

  /**
   * Update transaction progress
   */
  updateProgress(transactionId: string, processedRows: number): void;

  /**
   * Commit a transaction (mark as completed)
   */
  commitTransaction(transactionId: string): Promise<void>;

  /**
   * Rollback a transaction (delete all created records)
   */
  rollbackTransaction(transactionId: string): Promise<RollbackResult>;

  /**
   * Get a transaction by ID
   */
  getTransaction(transactionId: string): ImportTransaction | undefined;

  /**
   * Get all transactions
   */
  getAllTransactions(): ImportTransaction[];

  /**
   * Get interrupted transactions (pending or in_progress on load)
   */
  getInterruptedTransactions(): ImportTransaction[];

  /**
   * Mark pending/in_progress transactions as interrupted
   */
  cleanupInterrupted(): Promise<void>;

  /**
   * Delete a transaction record
   */
  deleteTransaction(transactionId: string): void;
}

/**
 * Delete function type for different import types
 */
export type DeleteRecordFn = (type: ImportType, recordId: number) => Promise<void>;
