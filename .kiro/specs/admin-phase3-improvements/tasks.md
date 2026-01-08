# Implementation Plan

## Phase 1: Service Layer Foundation

- [x] 1. Set up service layer infrastructure





  - [x] 1.1 Create service utilities (serviceError.ts, serviceResult.ts)


    - Implement `ServiceException` class with code, message, details
    - Implement `ServiceResult<T>` and `BatchResult<T>` interfaces
    - Implement error code constants
    - _Requirements: 1.3_
  - [x] 1.2 Write property test for service error format consistency


    - **Property 1: Service Error Format Consistency**
    - **Validates: Requirements 1.3**
  - [x] 1.3 Create base service class (base.ts)


    - Implement `BaseService` with dependency injection support
    - Implement `handleError` method for error wrapping
    - Implement `wrapAsync` helper for async operations
    - _Requirements: 1.2, 1.5_
  - [x] 1.4 Write property test for service independence from UI


    - **Property 2: Service Independence from UI**
    - **Validates: Requirements 1.2**

- [x] 2. Checkpoint - Ensure all tests pass





  - Ensure all tests pass, ask the user if questions arise.

## Phase 2: Domain Services Implementation

- [x] 3. Implement UserService






  - [x] 3.1 Create UserService interface and implementation

    - Implement CRUD operations (getUsers, getUserById, createUser, updateUser, deleteUser)
    - Implement status and role update methods
    - _Requirements: 2.1_
  - [x] 3.2 Implement user data validation


    - Implement email format validation (RFC 5322 compliant)
    - Implement phone format validation (Chinese mobile format)
    - Implement password strength validation
    - _Requirements: 2.2_

  - [x] 3.3 Write property test for user data validation

    - **Property 4: User Data Validation Completeness**
    - **Validates: Requirements 2.2**
  - [x] 3.4 Implement batch user operations


    - Implement batchUpdateStatus with validation
    - Implement batchUpdateRole with validation
    - Implement batchDelete with validation
    - _Requirements: 2.3_
  - [x] 3.5 Write property test for batch operation results


    - **Property 5: Batch Operation Result Completeness**
    - **Validates: Requirements 2.3**
  - [x] 3.6 Implement user data export


    - Implement exportUsers method returning headers and rows
    - Format data for Excel/CSV compatibility
    - _Requirements: 2.5_
  - [x] 3.7 Write property test for export data format


    - **Property 6: Export Data Format Consistency**
    - **Validates: Requirements 2.5**

- [x] 4. Checkpoint - Ensure all tests pass
  - Ensure all tests pass, ask the user if questions arise.

- [x] 5. Implement OrderService





  - [x] 5.1 Create OrderService interface and implementation


    - Implement query operations (getOrders, getOrderById)
    - Implement cancelOrder and refundOrder methods
    - _Requirements: 3.1_
  - [x] 5.2 Implement order cancellation validation


    - Implement canCancel method checking order status
    - Enforce cancellation rules (only pending/confirmed orders)
    - _Requirements: 3.2_
  - [x] 5.3 Write property test for order cancellation rules





    - **Property 7: Order Cancellation State Validation**
    - **Validates: Requirements 3.2**
  - [x] 5.4 Implement refund calculation
    - Implement calculateRefund method
    - Validate refund amount against original payment
    - Calculate platform fee and player amount
    - _Requirements: 3.3_
  - [x] 5.5 Write property test for refund calculation
    - **Property 8: Refund Calculation Accuracy**
    - **Validates: Requirements 3.3**
  - [x] 5.6 Implement batch order operations
    - Implement batchCancel with parallel processing
    - Implement batchComplete with parallel processing
    - _Requirements: 3.4_
  - [x] 5.7 Implement order statistics computation
    - Implement computeStatistics method
    - Calculate total revenue, order counts, completion rate
    - Implement computeTrend for trend data
    - _Requirements: 3.5_
  - [x] 5.8 Write property test for order statistics
    - **Property 9: Order Statistics Computation Accuracy**
    - **Validates: Requirements 3.5**

- [x] 6. Checkpoint - Ensure all tests pass
  - Ensure all tests pass, ask the user if questions arise.

- [x] 7. Implement PlayerService





  - [x] 7.1 Create PlayerService interface and implementation


    - Implement CRUD operations
    - Implement skill tag management
    - _Requirements: 4.1_
  - [x] 7.2 Implement player verification workflow


    - Implement verifyPlayer method
    - Implement canVerify method with state transition rules
    - Record audit information
    - _Requirements: 4.2_

  - [x] 7.3 Write property test for verification workflow

    - **Property 10: Player Verification Workflow Enforcement**
    - **Validates: Requirements 4.2**
  - [x] 7.4 Implement player earnings calculation

    - Implement calculateEarnings method
    - Apply commission rules correctly
    - _Requirements: 4.3_
  - [x] 7.5 Write property test for earnings calculation

    - **Property 11: Player Earnings Calculation Accuracy**
    - **Validates: Requirements 4.3**
  - [x] 7.6 Implement player statistics computation

    - Implement computeStatistics method
    - Calculate total earnings, order counts, rating averages
    - _Requirements: 4.4_
  - [x] 7.7 Write property test for player statistics

    - **Property 12: Player Statistics Computation Accuracy**
    - **Validates: Requirements 4.4**
  - [x] 7.8 Implement batch player operations

    - Implement batchUpdateStatus with state validation
    - Implement batchDelete
    - _Requirements: 4.5_

- [x] 8. Checkpoint - Ensure all tests pass
  - Ensure all tests pass, ask the user if questions arise.
  - ✅ All 1616 tests passed (66 test files)

## Phase 3: Data Import Framework

- [x] 9. Create import infrastructure





  - [x] 9.1 Create file parsers


    - Implement Excel parser using xlsx library
    - Implement CSV parser
    - Handle file size validation (max 10MB)
    - _Requirements: 5.1_
  - [x] 9.2 Write property test for file format validation


    - **Property 13: File Format Validation**
    - **Validates: Requirements 5.1**
  - [x] 9.3 Create import template definitions


    - Define userImportTemplate with columns
    - Define playerImportTemplate with columns
    - Define gameImportTemplate with columns
    - _Requirements: 6.1, 7.1, 8.1_
  - [x] 9.4 Implement structure validation


    - Validate required columns presence
    - Report missing/extra columns
    - _Requirements: 5.2_
  - [x] 9.5 Write property test for structure validation


    - **Property 14: Import Structure Validation**
    - **Validates: Requirements 5.2**

- [x] 10. Checkpoint - Ensure all tests pass




  - Ensure all tests pass, ask the user if questions arise.

- [x] 11. Implement data validators





  - [x] 11.1 Create user data validator


    - Validate email format and uniqueness
    - Validate phone format and uniqueness
    - Collect all errors per row
    - _Requirements: 6.2_
  - [x] 11.2 Create player data validator


    - Validate user reference exists
    - Validate user not already a player
    - Validate skill tags format
    - _Requirements: 7.2, 7.4_
  - [x] 11.3 Write property test for skill tag parsing


    - **Property 21: Skill Tag Parsing**
    - **Validates: Requirements 7.4**
  - [x] 11.4 Create game data validator


    - Validate game key uniqueness
    - Validate category references
    - _Requirements: 8.2_
  - [x] 11.5 Write property test for duplicate detection


    - **Property 17: Import Duplicate Detection**
    - **Validates: Requirements 6.2, 7.2, 8.2**
  - [x] 11.6 Write property test for data validation completeness


    - **Property 15: Import Data Validation Completeness**
    - **Validates: Requirements 5.3**

- [x] 12. Checkpoint - Ensure all tests pass
  - Ensure all tests pass, ask the user if questions arise.
  - ✅ All 208 service tests passed (10 test files)

- [x] 13. Implement ImportService core





  - [x] 13.1 Create ImportService interface and implementation


    - Implement parseFile method
    - Implement validateStructure method
    - Implement getTemplate and downloadTemplate methods
    - _Requirements: 5.1, 5.2, 5.3_
  - [x] 13.2 Implement user import

    - Implement importUsers method
    - Generate secure temporary passwords
    - Set default values for optional fields
    - _Requirements: 6.3, 6.4_
  - [x] 13.3 Write property test for password generation


    - **Property 18: Password Generation Security**
    - **Validates: Requirements 6.3**
  - [x] 13.4 Implement player import

    - Implement importPlayers method
    - Set initial verification status to pending
    - Parse and validate skill tags
    - _Requirements: 7.3_
  - [x] 13.5 Write property test for player import initial state


    - **Property 20: Player Import Initial State**
    - **Validates: Requirements 7.3**
  - [x] 13.6 Implement game import

    - Implement importGames method
    - Apply default values (isActive=true, sortOrder=0)
    - Handle duplicate key options (skip/update/fail)
    - _Requirements: 8.3, 8.4_
  - [x] 13.7 Write property test for game import defaults

    - **Property 22: Game Import Defaults**
    - **Validates: Requirements 8.3**
  - [x] 13.8 Write property test for import summary accuracy


    - **Property 16: Import Summary Accuracy**
    - **Validates: Requirements 5.5, 6.5**

- [x] 14. Checkpoint - Ensure all tests pass
  - Ensure all tests pass, ask the user if questions arise.
  - ✅ All 220 service tests passed (11 test files)

## Phase 4: Import History and UI Integration

- [x] 15. Implement import history







  - [x] 15.1 Create import history storage


    - Store import metadata (timestamp, user, file, counts, status)
    - Store row-by-row results for failed imports
    - _Requirements: 9.1_
  - [x] 15.2 Write property test for metadata recording







    - **Property 23: Import Metadata Recording**
    - **Validates: Requirements 9.1**
  - [x] 15.3 Implement history query and details


    - Implement getImportHistory with filtering
    - Implement getImportDetails
    - _Requirements: 9.2, 9.3_
  - [x] 15.4 Write property test for error detail preservation


    - **Property 19: Error Detail Preservation**
    - **Validates: Requirements 6.4, 9.3, 9.4**
  - [x] 15.5 Implement error report download


    - Generate report with original data + status + errors
    - Support Excel/CSV format
    - _Requirements: 9.5_
  - [x] 15.6 Write property test for report format


    - **Property 24: Import Result Report Format**
    - **Validates: Requirements 9.5**

- [x] 16. Checkpoint - Ensure all tests pass





  - Ensure all tests pass, ask the user if questions arise.

- [x] 17. Create UI components for import





  - [x] 17.1 Create ImportModal component


    - File upload with drag-and-drop
    - Template download button
    - Preview table for parsed data
    - Error display with row numbers
    - _Requirements: 5.4_
  - [x] 17.2 Create ImportHistoryTable component

    - Paginated list of past imports
    - Filter by type and date range
    - Link to import details
    - _Requirements: 9.2_
  - [x] 17.3 Integrate import into User management page


    - Add import button to toolbar
    - Connect to ImportModal
    - Refresh list after import
    - _Requirements: 6.5_

  - [x] 17.4 Integrate import into Player management page

    - Add import button to toolbar
    - Connect to ImportModal
    - Refresh list after import
    - _Requirements: 7.5_

  - [x] 17.5 Integrate import into Game management page

    - Add import button to toolbar
    - Connect to ImportModal
    - Refresh list after import
    - _Requirements: 8.5_

- [x] 18. Checkpoint - Ensure all tests pass

  - Ensure all tests pass, ask the user if questions arise.
  - ✅ All 1746 tests passed (73 test files)

## Phase 5: Store Integration and Refactoring

- [x] 19. Implement observability infrastructure





  - [x] 19.1 Create ServiceLogger interface and implementation


    - Implement debug, info, warn, error methods
    - Implement parameter sanitization for sensitive data
    - Support optional external error tracking integration
    - _Requirements: 10.1, 10.3_

  - [x] 19.2 Create PerformanceMonitor

    - Implement startTimer and recordMetric methods
    - Implement slow operation detection (>3s threshold)
    - Store metrics for debugging
    - _Requirements: 10.2, 10.5_
  - [x] 19.3 Write property test for service logging


    - **Property 25: Service Method Logging**
    - **Validates: Requirements 10.1, 10.2**
  - [x] 19.4 Write property test for slow operation warning


    - **Property 26: Slow Operation Warning**
    - **Validates: Requirements 10.5**
  - [x] 19.5 Integrate logging into BaseService


    - Add withLogging wrapper method
    - Log batch operation progress at 10% intervals
    - _Requirements: 10.4_

- [x] 20. Implement concurrency control
  - [x] 20.1 Create ConcurrencyController
    - Implement withDeduplication for duplicate prevention
    - Implement processWithConcurrency with configurable limits
    - Implement Semaphore for concurrency limiting
    - _Requirements: 11.1, 11.2_
  - [x] 20.2 Write property test for concurrency limit
    - **Property 27: Batch Concurrency Limit**
    - **Validates: Requirements 11.1**
  - [x] 20.3 Write property test for duplicate prevention
    - **Property 28: Duplicate Operation Prevention**
    - **Validates: Requirements 11.2**
  - [x] 20.4 Implement retry with exponential backoff
    - Implement withRetry method
    - Implement isRetryableError detection
    - Configure retry attempts and delays
    - _Requirements: 11.3_
  - [x] 20.5 Write property test for retry behavior
    - **Property 29: Retry with Exponential Backoff**
    - **Validates: Requirements 11.3**
  - [x] 20.6 Implement batch chunking
    - Implement chunkArray utility
    - Configure default chunk size (50 items)
    - _Requirements: 11.4_

- [x] 21. Checkpoint - Ensure all tests pass





  - Ensure all tests pass, ask the user if questions arise.

- [x] 22. Implement import transaction management







  - [x] 22.1 Create ImportTransactionManager


    - Implement startTransaction, recordCreated, commitTransaction
    - Implement transaction persistence to localStorage
    - _Requirements: 12.1_
  - [x] 22.2 Write property test for transaction tracking


    - **Property 30: Import Transaction Tracking**
    - **Validates: Requirements 12.1**
  - [x] 22.3 Implement rollback functionality




    - Implement rollbackTransaction method
    - Delete records in reverse order
    - Log all rollback operations
    - _Requirements: 12.3, 12.4_
  - [x] 22.4 Write property test for rollback completeness


    - **Property 31: Rollback Completeness**
    - **Validates: Requirements 12.3, 12.4**
  - [x] 22.5 Implement interrupted import handling


    - Implement cleanupInterrupted method
    - Load persisted transactions on startup
    - Mark interrupted transactions
    - _Requirements: 12.5_
  - [x] 22.6 Write property test for interrupted import detection


    - **Property 32: Interrupted Import Detection**
    - **Validates: Requirements 12.5**
  - [x] 22.7 Update ImportService with transaction support


    - Implement importWithTransaction method
    - Add rollback prompt callback
    - Implement getInterruptedImports and resumeOrCleanup
    - _Requirements: 12.2_

- [x] 23. Checkpoint - Ensure all tests pass
  - Ensure all tests pass, ask the user if questions arise.
  - ✅ All 1840 tests passed (76 test files)

- [x] 24. Refactor stores to use services








  - [x] 24.1 Refactor userStore to use UserService

    - Replace direct API calls with service calls
    - Use service validation methods
    - Handle service errors consistently
    - _Requirements: 1.2_

  - [x] 24.2 Refactor orderStore to use OrderService

    - Replace direct API calls with service calls
    - Use service calculation methods
    - Handle service errors consistently
    - _Requirements: 1.2_


  - [x] 24.3 Refactor playerStore to use PlayerService





    - Replace direct API calls with service calls
    - Use service validation methods
    - Handle service errors consistently

    - _Requirements: 1.2_


  - [x] 24.4 Write property test for multi-API orchestration









    - **Property 3: Multi-API Orchestration Graceful Handling**
    - **Validates: Requirements 1.4**

- [x] 25. Final Checkpoint - Ensure all tests pass





  - Ensure all tests pass, ask the user if questions arise.

- [x] 26. Update service exports and documentation





  - [x] 26.1 Create unified service exports


    - Export all services from services/index.ts
    - Export all types and interfaces
    - _Requirements: 1.1_


  - [x] 26.2 Update stores documentation
    - Document service layer usage patterns
    - Update BEST_PRACTICES.md

    - _Requirements: 13.1_
  - [x] 26.3 Verify test coverage meets 80% threshold

    - Run coverage report
    - Add tests for uncovered code paths
    - _Requirements: 13.5_
