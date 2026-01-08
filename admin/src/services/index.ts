/**
 * Service Layer - Unified Exports
 * 服务层统一导出
 *
 * This module provides a centralized export point for all service layer components:
 * - Domain Services: Business logic for users, orders, players
 * - Import Services: Data import functionality for Excel/CSV files
 * - Service Utilities: Error handling, logging, performance monitoring, concurrency control
 *
 * @module services
 *
 * @example
 * // Import domain services
 * import { userService, orderService, playerService } from '@/services';
 *
 * // Import service utilities
 * import { ServiceException, ServiceErrorCodes } from '@/services';
 *
 * // Import import services
 * import { importService, parseFile } from '@/services';
 */

// ============================================================================
// Init Services
// ============================================================================
export { initApp, smartInit, forceInit, type InitConfig, type InitResult } from './init';

// ============================================================================
// Service Utilities
// ============================================================================
export * from './utils';

// ============================================================================
// Domain Services
// ============================================================================
export * from './domain';

// ============================================================================
// Import Services
// ============================================================================
export * from './import';

// ============================================================================
// Re-export commonly used types for convenience
// ============================================================================

// Service Result Types
export type {
  ServiceResult,
  BatchResult,
  BatchItemResult,
} from './utils/serviceResult';

// Service Error Types
export type { ServiceError, ServiceErrorCode } from './utils/serviceError';

// Domain Service Types
export type {
  IUserService,
  UserValidationResult,
  PasswordValidationResult,
  UserExportData,
} from './domain/userService';

export type {
  IOrderService,
  CancellationCheckResult,
  RefundCheckResult,
  RefundCalculation,
  OrderStatistics,
} from './domain/orderService';

export type {
  IPlayerService,
  VerificationCheckResult,
  EarningsCalculation,
  PlayerStatistics,
  SkillTagValidationResult,
  VerificationStatus,
} from './domain/playerService';

// Import Service Types
export type {
  IImportService,
  ParsedRow,
  ImportPreview,
  ImportResult,
  ImportOptions,
  DuplicateKeyHandling,
} from './import/importService';

// Import Template Types
export type {
  ImportType,
  ImportTemplate,
  ImportColumn,
  ColumnType,
} from './import/templates/types';

// Import History Types
export type {
  ImportStatus,
  ImportRowResult,
  ImportHistoryRecord,
  ImportHistoryQueryParams,
  ImportHistoryPage,
  IImportHistoryStorage,
  ImportHistoryFilter,
  PaginationOptions,
  ImportHistoryDetails,
  IImportHistoryService,
  ReportFormat,
  ErrorReportOptions,
  ErrorReportResult,
} from './import/history';

// Import Transaction Types
export type {
  TransactionStatus,
  ImportTransaction,
  RollbackResult,
  ITransactionManager,
  DeleteRecordFn,
} from './import/transaction/types';

// Parser Types
export type {
  SupportedFileType,
  ParseResult,
  FileValidationResult,
  FileParser,
} from './import/parsers/types';

// Concurrency Types
export type { ConcurrencyConfig, ProcessResult } from './utils/concurrency';

// Logger Types
export type {
  LogEntry,
  LoggerConfig,
  ExternalErrorTracker,
} from './utils/logger';

// Performance Types
export type { PerformanceMetrics, PerformanceMonitorConfig } from './utils/performance';
