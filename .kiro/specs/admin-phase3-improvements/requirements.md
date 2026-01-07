# Requirements Document

## Introduction

本文档定义了 GameLink 管理后台 Phase 3 改进的需求规范。Phase 3 聚焦于两个核心改进：
1. **Service 层解耦** - 引入领域服务层，将业务逻辑从组件和 Store 中抽离
2. **数据导入功能** - 支持批量导入用户、陪玩师、游戏等数据

这些改进旨在提升代码可维护性、可测试性，并为运营团队提供高效的数据管理工具。

## Glossary

- **Admin System**: GameLink 管理后台系统，基于 React + TypeScript + Ant Design 构建
- **Service Layer**: 业务逻辑层，负责封装复杂业务规则，独立于 UI 和状态管理
- **Domain Service**: 领域服务，按业务领域划分的服务模块（如 UserService、OrderService）
- **Store**: Zustand 状态管理模块，负责管理 UI 状态和缓存
- **Data Import**: 批量数据导入功能，支持 Excel/CSV 文件解析和验证
- **Validation Rule**: 数据验证规则，确保导入数据符合业务约束
- **Import Template**: 导入模板，预定义的 Excel/CSV 格式文件

## Requirements

### Requirement 1: Service Layer Architecture

**User Story:** As a developer, I want a clear separation between business logic and UI components, so that I can maintain and test the codebase more effectively.

#### Acceptance Criteria

1. WHEN a developer creates a new business feature THEN the Admin System SHALL provide a standardized service layer structure under `admin/src/services/domain/`
2. WHEN business logic is executed THEN the Service Layer SHALL handle all API calls, data transformation, and business rule validation independently from UI components
3. WHEN a service method fails THEN the Service Layer SHALL return a standardized error object containing error code, message, and optional details
4. WHEN multiple API calls are required for a single operation THEN the Service Layer SHALL orchestrate these calls and handle partial failures gracefully
5. WHEN a service is instantiated THEN the Service Layer SHALL support dependency injection for easier testing and mocking

### Requirement 2: User Domain Service

**User Story:** As a developer, I want a dedicated UserService that encapsulates all user-related business logic, so that user operations are consistent across the application.

#### Acceptance Criteria

1. WHEN managing users THEN the UserService SHALL provide methods for CRUD operations, status changes, role assignments, and batch operations
2. WHEN validating user data THEN the UserService SHALL enforce business rules including email format, phone format, and password strength requirements
3. WHEN performing batch user operations THEN the UserService SHALL validate all items before execution and return detailed results for each item
4. WHEN querying users THEN the UserService SHALL support filtering by role, status, date range, and keyword search
5. WHEN exporting user data THEN the UserService SHALL generate properly formatted data suitable for Excel/CSV export

### Requirement 3: Order Domain Service

**User Story:** As a developer, I want a dedicated OrderService that encapsulates all order-related business logic, so that order operations follow consistent business rules.

#### Acceptance Criteria

1. WHEN managing orders THEN the OrderService SHALL provide methods for querying, canceling, refunding, and batch operations
2. WHEN canceling an order THEN the OrderService SHALL validate the order status and enforce cancellation rules based on order state
3. WHEN processing a refund THEN the OrderService SHALL calculate the refund amount based on business rules and validate against the original payment
4. WHEN performing batch order operations THEN the OrderService SHALL process items in parallel where possible and return aggregated results
5. WHEN querying order statistics THEN the OrderService SHALL compute metrics including total revenue, order counts by status, and trend data

### Requirement 4: Player Domain Service

**User Story:** As a developer, I want a dedicated PlayerService that encapsulates all player-related business logic, so that player management is consistent and maintainable.

#### Acceptance Criteria

1. WHEN managing players THEN the PlayerService SHALL provide methods for CRUD operations, verification status changes, and skill tag management
2. WHEN verifying a player THEN the PlayerService SHALL enforce verification workflow rules and record audit information
3. WHEN calculating player earnings THEN the PlayerService SHALL apply commission rules and compute net earnings accurately
4. WHEN querying player statistics THEN the PlayerService SHALL compute metrics including total earnings, order counts, and rating averages
5. WHEN batch updating player status THEN the PlayerService SHALL validate each player's current state before applying changes

### Requirement 5: Data Import Framework

**User Story:** As an admin, I want to import data from Excel/CSV files, so that I can efficiently onboard large amounts of data without manual entry.

#### Acceptance Criteria

1. WHEN uploading an import file THEN the Admin System SHALL accept Excel (.xlsx, .xls) and CSV (.csv) file formats up to 10MB in size
2. WHEN parsing an import file THEN the Admin System SHALL validate the file structure against the expected template and report structural errors
3. WHEN validating import data THEN the Admin System SHALL check each row against business rules and collect all validation errors before processing
4. WHEN displaying validation results THEN the Admin System SHALL show a preview of valid rows and a detailed list of errors with row numbers and field names
5. WHEN the user confirms import THEN the Admin System SHALL process valid rows and skip invalid rows, providing a final summary of imported and skipped records

### Requirement 6: User Data Import

**User Story:** As an admin, I want to import user data from Excel/CSV files, so that I can quickly onboard multiple users at once.

#### Acceptance Criteria

1. WHEN downloading the user import template THEN the Admin System SHALL provide an Excel file with columns for name, email, phone, role, and status
2. WHEN validating user import data THEN the Admin System SHALL check for duplicate emails and phones against existing database records
3. WHEN importing users THEN the Admin System SHALL generate secure temporary passwords for new users and optionally send welcome emails
4. WHEN an import row has validation errors THEN the Admin System SHALL display the specific field errors and allow the user to download an error report
5. WHEN the import completes THEN the Admin System SHALL display a summary showing total rows, successful imports, and failed rows with reasons

### Requirement 7: Player Data Import

**User Story:** As an admin, I want to import player data from Excel/CSV files, so that I can efficiently register multiple players.

#### Acceptance Criteria

1. WHEN downloading the player import template THEN the Admin System SHALL provide an Excel file with columns for user reference, nickname, bio, hourly rate, main game, and skill tags
2. WHEN validating player import data THEN the Admin System SHALL verify that referenced users exist and are not already registered as players
3. WHEN importing players THEN the Admin System SHALL set the initial verification status to pending and create associated records
4. WHEN skill tags are provided THEN the Admin System SHALL parse comma-separated tags and validate against allowed tag values
5. WHEN the import completes THEN the Admin System SHALL display a summary and provide options to view imported players or download error report

### Requirement 8: Game Data Import

**User Story:** As an admin, I want to import game data from Excel/CSV files, so that I can quickly add multiple games to the platform.

#### Acceptance Criteria

1. WHEN downloading the game import template THEN the Admin System SHALL provide an Excel file with columns for key, name, category, description, and active status
2. WHEN validating game import data THEN the Admin System SHALL check for duplicate game keys and validate category references
3. WHEN importing games THEN the Admin System SHALL create game records with default sort order and active status
4. WHEN a game key already exists THEN the Admin System SHALL offer options to skip, update, or fail the import row
5. WHEN the import completes THEN the Admin System SHALL refresh the game list and display the import summary

### Requirement 9: Import History and Audit

**User Story:** As an admin, I want to view the history of data imports, so that I can track what data was imported and by whom.

#### Acceptance Criteria

1. WHEN an import operation completes THEN the Admin System SHALL record the import metadata including timestamp, user, file name, record counts, and status
2. WHEN viewing import history THEN the Admin System SHALL display a paginated list of past imports with filtering by type and date range
3. WHEN viewing import details THEN the Admin System SHALL show the original file name, row-by-row results, and any error messages
4. WHEN an import fails partially THEN the Admin System SHALL preserve the error details for later review and potential retry
5. WHEN downloading import results THEN the Admin System SHALL generate a report file containing the original data plus status and error columns

### Requirement 10: Service Layer Testing

**User Story:** As a developer, I want comprehensive tests for the service layer, so that I can ensure business logic correctness and prevent regressions.

#### Acceptance Criteria

1. WHEN writing service tests THEN the Admin System SHALL provide mock utilities for API calls and external dependencies
2. WHEN testing service methods THEN the test suite SHALL cover success cases, error cases, and edge cases for each method
3. WHEN testing batch operations THEN the test suite SHALL verify partial success scenarios and error aggregation
4. WHEN testing data validation THEN the test suite SHALL verify all validation rules are correctly enforced
5. WHEN running the test suite THEN the service layer tests SHALL achieve at least 80% code coverage
