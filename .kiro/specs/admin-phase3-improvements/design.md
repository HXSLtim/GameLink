# Design Document: Admin Phase 3 Improvements

## Overview

Phase 3 改进包含两个核心功能模块：

1. **Service 层解耦** - 引入领域服务层架构，将业务逻辑从 UI 组件和 Zustand Store 中抽离，提升代码可维护性和可测试性
2. **数据导入功能** - 提供 Excel/CSV 批量数据导入能力，支持用户、陪玩师、游戏等数据的批量导入

### 设计目标

- 清晰的职责分离：UI → Store → Service → API
- 统一的错误处理和数据转换
- 可测试的业务逻辑
- 用户友好的批量数据导入体验

## Architecture

### 整体架构

```
┌─────────────────────────────────────────────────────────────────┐
│                        UI Components                             │
│  (Pages, Modals, Forms)                                         │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                      Zustand Stores                              │
│  (UI State, Cache, Pagination)                                  │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                     Domain Services                              │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────────┐    │
│  │UserService│  │OrderService│ │PlayerService│ │ImportService│    │
│  └──────────┘  └──────────┘  └──────────┘  └──────────────┘    │
│                                                                  │
│  - Business Logic                                                │
│  - Data Validation                                               │
│  - Data Transformation                                           │
│  - Error Handling                                                │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                        API Layer                                 │
│  (adminApi, userApi, etc.)                                      │
└─────────────────────────────────────────────────────────────────┘
```

### 目录结构

```
admin/src/services/
├── domain/                      # 领域服务
│   ├── base.ts                  # 基础服务类
│   ├── userService.ts           # 用户领域服务
│   ├── orderService.ts          # 订单领域服务
│   ├── playerService.ts         # 陪玩师领域服务
│   └── index.ts                 # 统一导出
├── import/                      # 数据导入服务
│   ├── importService.ts         # 导入核心服务
│   ├── parsers/                 # 文件解析器
│   │   ├── excelParser.ts       # Excel 解析
│   │   ├── csvParser.ts         # CSV 解析
│   │   └── index.ts
│   ├── validators/              # 数据验证器
│   │   ├── userValidator.ts     # 用户数据验证
│   │   ├── playerValidator.ts   # 陪玩师数据验证
│   │   ├── gameValidator.ts     # 游戏数据验证
│   │   └── index.ts
│   ├── templates/               # 导入模板
│   │   ├── userTemplate.ts      # 用户模板定义
│   │   ├── playerTemplate.ts    # 陪玩师模板定义
│   │   ├── gameTemplate.ts      # 游戏模板定义
│   │   └── index.ts
│   └── index.ts
├── utils/                       # 服务工具
│   ├── serviceError.ts          # 统一错误类
│   ├── serviceResult.ts         # 统一结果类
│   └── index.ts
└── index.ts                     # 统一导出
```

## Components and Interfaces

### 1. 基础服务接口

```typescript
// services/utils/serviceError.ts
export interface ServiceError {
  code: string;
  message: string;
  details?: Record<string, unknown>;
  originalError?: Error;
}

export class ServiceException extends Error {
  constructor(
    public code: string,
    message: string,
    public details?: Record<string, unknown>,
    public originalError?: Error
  ) {
    super(message);
    this.name = 'ServiceException';
  }

  toError(): ServiceError {
    return {
      code: this.code,
      message: this.message,
      details: this.details,
      originalError: this.originalError,
    };
  }
}

// services/utils/serviceResult.ts
export interface ServiceResult<T> {
  success: boolean;
  data?: T;
  error?: ServiceError;
}

export interface BatchResult<T> {
  success: boolean;
  total: number;
  succeeded: number;
  failed: number;
  results: Array<{
    index: number;
    success: boolean;
    data?: T;
    error?: ServiceError;
  }>;
}
```

### 2. 基础服务类

```typescript
// services/domain/base.ts
export interface ServiceDependencies {
  api: typeof adminApi;
  cache?: CacheService;
  logger?: LoggerService;
}

export abstract class BaseService {
  protected api: typeof adminApi;
  protected cache?: CacheService;
  protected logger?: LoggerService;

  constructor(deps: ServiceDependencies) {
    this.api = deps.api;
    this.cache = deps.cache;
    this.logger = deps.logger;
  }

  protected handleError(error: unknown, context: string): ServiceException {
    if (error instanceof ServiceException) {
      return error;
    }
    
    const message = error instanceof Error ? error.message : 'Unknown error';
    return new ServiceException(
      'SERVICE_ERROR',
      `${context}: ${message}`,
      undefined,
      error instanceof Error ? error : undefined
    );
  }

  protected async wrapAsync<T>(
    operation: () => Promise<T>,
    context: string
  ): Promise<ServiceResult<T>> {
    try {
      const data = await operation();
      return { success: true, data };
    } catch (error) {
      const serviceError = this.handleError(error, context);
      return { success: false, error: serviceError.toError() };
    }
  }
}
```

### 3. UserService 接口

```typescript
// services/domain/userService.ts
export interface UserValidationResult {
  valid: boolean;
  errors: Array<{
    field: string;
    message: string;
  }>;
}

export interface UserExportData {
  headers: string[];
  rows: string[][];
}

export interface UserService {
  // CRUD Operations
  getUsers(params: UserQueryParams): Promise<ServiceResult<User[]>>;
  getUserById(id: number): Promise<ServiceResult<User>>;
  createUser(data: CreateUserDto): Promise<ServiceResult<User>>;
  updateUser(id: number, data: UpdateUserDto): Promise<ServiceResult<User>>;
  deleteUser(id: number): Promise<ServiceResult<void>>;

  // Status & Role
  updateUserStatus(id: number, status: string): Promise<ServiceResult<User>>;
  updateUserRole(id: number, role: string): Promise<ServiceResult<User>>;

  // Batch Operations
  batchUpdateStatus(userIds: number[], status: string): Promise<BatchResult<void>>;
  batchUpdateRole(userIds: number[], role: string): Promise<BatchResult<void>>;
  batchDelete(userIds: number[]): Promise<BatchResult<void>>;

  // Validation
  validateUserData(data: Partial<CreateUserDto>): UserValidationResult;
  validateEmail(email: string): boolean;
  validatePhone(phone: string): boolean;
  validatePassword(password: string): { valid: boolean; errors: string[] };

  // Export
  exportUsers(users: User[]): UserExportData;
}
```

### 4. OrderService 接口

```typescript
// services/domain/orderService.ts
export interface OrderStatistics {
  totalOrders: number;
  totalRevenue: number;
  ordersByStatus: Record<string, number>;
  averageOrderValue: number;
  completionRate: number;
}

export interface RefundCalculation {
  originalAmount: number;
  refundAmount: number;
  platformFee: number;
  playerAmount: number;
}

export interface OrderService {
  // Query
  getOrders(params: OrderQueryParams): Promise<ServiceResult<Order[]>>;
  getOrderById(id: number): Promise<ServiceResult<Order>>;

  // Operations
  cancelOrder(id: number, reason?: string): Promise<ServiceResult<Order>>;
  refundOrder(id: number, data: RefundRequest): Promise<ServiceResult<Order>>;

  // Batch Operations
  batchCancel(orderIds: number[], reason?: string): Promise<BatchResult<void>>;
  batchComplete(orderIds: number[]): Promise<BatchResult<void>>;

  // Validation
  canCancel(order: Order): { allowed: boolean; reason?: string };
  canRefund(order: Order): { allowed: boolean; reason?: string };
  calculateRefund(order: Order, requestedAmount: number): RefundCalculation;

  // Statistics
  computeStatistics(orders: Order[]): OrderStatistics;
  computeTrend(orders: Order[], days: number): TrendData[];
}
```

### 5. PlayerService 接口

```typescript
// services/domain/playerService.ts
export interface PlayerStatistics {
  totalEarnings: number;
  monthlyEarnings: number;
  completedOrders: number;
  averageRating: number;
  ratingCount: number;
}

export interface EarningsCalculation {
  grossAmount: number;
  commissionRate: number;
  commissionAmount: number;
  netAmount: number;
}

export interface PlayerService {
  // CRUD
  getPlayers(params: PlayerQueryParams): Promise<ServiceResult<Player[]>>;
  getPlayerById(id: number): Promise<ServiceResult<Player>>;
  createPlayer(data: CreatePlayerDto): Promise<ServiceResult<Player>>;
  updatePlayer(id: number, data: UpdatePlayerDto): Promise<ServiceResult<Player>>;
  deletePlayer(id: number): Promise<ServiceResult<void>>;

  // Verification
  verifyPlayer(id: number, status: string, remark?: string): Promise<ServiceResult<Player>>;
  canVerify(player: Player, newStatus: string): { allowed: boolean; reason?: string };

  // Skill Tags
  updateSkillTags(id: number, tags: string[]): Promise<ServiceResult<void>>;
  validateSkillTags(tags: string[]): { valid: boolean; invalidTags: string[] };

  // Batch Operations
  batchUpdateStatus(playerIds: number[], status: string): Promise<BatchResult<void>>;
  batchDelete(playerIds: number[]): Promise<BatchResult<void>>;

  // Earnings
  calculateEarnings(order: Order, commissionRule?: CommissionRule): EarningsCalculation;
  computeStatistics(playerId: number, orders: Order[]): PlayerStatistics;
}
```

### 6. ImportService 接口

```typescript
// services/import/importService.ts
export type ImportType = 'user' | 'player' | 'game';

export interface ImportTemplate {
  type: ImportType;
  columns: Array<{
    key: string;
    label: string;
    required: boolean;
    type: 'string' | 'number' | 'boolean' | 'date';
    validation?: (value: unknown) => boolean;
  }>;
}

export interface ParsedRow {
  rowNumber: number;
  data: Record<string, unknown>;
  errors: Array<{
    field: string;
    message: string;
  }>;
  isValid: boolean;
}

export interface ImportPreview {
  totalRows: number;
  validRows: ParsedRow[];
  invalidRows: ParsedRow[];
  structureErrors: string[];
}

export interface ImportResult {
  success: boolean;
  totalRows: number;
  importedCount: number;
  skippedCount: number;
  errors: Array<{
    rowNumber: number;
    field?: string;
    message: string;
  }>;
}

export interface ImportHistory {
  id: string;
  type: ImportType;
  fileName: string;
  uploadedBy: number;
  uploadedAt: string;
  totalRows: number;
  importedCount: number;
  skippedCount: number;
  status: 'pending' | 'processing' | 'completed' | 'failed';
  errorDetails?: string;
}

export interface ImportService {
  // Template
  getTemplate(type: ImportType): ImportTemplate;
  downloadTemplate(type: ImportType): Blob;

  // Parsing
  parseFile(file: File, type: ImportType): Promise<ImportPreview>;
  validateStructure(headers: string[], template: ImportTemplate): string[];

  // Import
  importData(type: ImportType, rows: ParsedRow[]): Promise<ImportResult>;
  importUsers(rows: ParsedRow[]): Promise<ImportResult>;
  importPlayers(rows: ParsedRow[]): Promise<ImportResult>;
  importGames(rows: ParsedRow[]): Promise<ImportResult>;

  // History
  getImportHistory(params: { type?: ImportType; page?: number }): Promise<ImportHistory[]>;
  getImportDetails(id: string): Promise<ImportHistory>;
  downloadErrorReport(id: string): Blob;
}
```

## Data Models

### Service Error Codes

```typescript
export const ServiceErrorCodes = {
  // General
  UNKNOWN_ERROR: 'UNKNOWN_ERROR',
  VALIDATION_ERROR: 'VALIDATION_ERROR',
  NOT_FOUND: 'NOT_FOUND',
  UNAUTHORIZED: 'UNAUTHORIZED',
  FORBIDDEN: 'FORBIDDEN',

  // User
  USER_EMAIL_EXISTS: 'USER_EMAIL_EXISTS',
  USER_PHONE_EXISTS: 'USER_PHONE_EXISTS',
  USER_INVALID_STATUS: 'USER_INVALID_STATUS',

  // Order
  ORDER_CANNOT_CANCEL: 'ORDER_CANNOT_CANCEL',
  ORDER_CANNOT_REFUND: 'ORDER_CANNOT_REFUND',
  ORDER_INVALID_REFUND_AMOUNT: 'ORDER_INVALID_REFUND_AMOUNT',

  // Player
  PLAYER_INVALID_VERIFICATION: 'PLAYER_INVALID_VERIFICATION',
  PLAYER_ALREADY_EXISTS: 'PLAYER_ALREADY_EXISTS',

  // Import
  IMPORT_INVALID_FILE: 'IMPORT_INVALID_FILE',
  IMPORT_FILE_TOO_LARGE: 'IMPORT_FILE_TOO_LARGE',
  IMPORT_INVALID_STRUCTURE: 'IMPORT_INVALID_STRUCTURE',
  IMPORT_VALIDATION_FAILED: 'IMPORT_VALIDATION_FAILED',
} as const;
```

### Import Templates

```typescript
// User Import Template
export const userImportTemplate: ImportTemplate = {
  type: 'user',
  columns: [
    { key: 'name', label: '姓名', required: true, type: 'string' },
    { key: 'email', label: '邮箱', required: true, type: 'string' },
    { key: 'phone', label: '手机号', required: true, type: 'string' },
    { key: 'role', label: '角色', required: false, type: 'string' },
    { key: 'status', label: '状态', required: false, type: 'string' },
  ],
};

// Player Import Template
export const playerImportTemplate: ImportTemplate = {
  type: 'player',
  columns: [
    { key: 'userEmail', label: '用户邮箱', required: true, type: 'string' },
    { key: 'nickname', label: '昵称', required: false, type: 'string' },
    { key: 'bio', label: '简介', required: false, type: 'string' },
    { key: 'hourlyRate', label: '时薪(元)', required: false, type: 'number' },
    { key: 'mainGame', label: '主游戏', required: false, type: 'string' },
    { key: 'skillTags', label: '技能标签', required: false, type: 'string' },
  ],
};

// Game Import Template
export const gameImportTemplate: ImportTemplate = {
  type: 'game',
  columns: [
    { key: 'key', label: '游戏标识', required: true, type: 'string' },
    { key: 'name', label: '游戏名称', required: true, type: 'string' },
    { key: 'category', label: '分类', required: false, type: 'string' },
    { key: 'description', label: '描述', required: false, type: 'string' },
    { key: 'isActive', label: '是否启用', required: false, type: 'boolean' },
  ],
};
```

## Correctness Properties

*A property is a characteristic or behavior that should hold true across all valid executions of a system-essentially, a formal statement about what the system should do. Properties serve as the bridge between human-readable specifications and machine-verifiable correctness guarantees.*

### Property 1: Service Error Format Consistency
*For any* service method that throws an error, the error object SHALL contain code, message, and optionally details fields in a consistent format.
**Validates: Requirements 1.3**

### Property 2: Service Independence from UI
*For any* service module, the module SHALL NOT import any React, Ant Design, or other UI-specific dependencies.
**Validates: Requirements 1.2**

### Property 3: Multi-API Orchestration Graceful Handling
*For any* service method that calls multiple APIs, if one API fails, the service SHALL return a result that indicates partial success and includes details of which operations succeeded and which failed.
**Validates: Requirements 1.4**

### Property 4: User Data Validation Completeness
*For any* user data input, the validation function SHALL check email format, phone format, and password strength, returning all validation errors (not just the first one).
**Validates: Requirements 2.2**

### Property 5: Batch Operation Result Completeness
*For any* batch operation on users, orders, or players, the result SHALL contain an entry for each input item indicating success or failure with details.
**Validates: Requirements 2.3, 3.4, 4.5**

### Property 6: Export Data Format Consistency
*For any* user data export, the exported data SHALL contain headers matching the expected columns and rows containing properly formatted values.
**Validates: Requirements 2.5**

### Property 7: Order Cancellation State Validation
*For any* order cancellation request, the service SHALL only allow cancellation if the order status is in a cancellable state (pending, confirmed) and SHALL reject with a clear reason otherwise.
**Validates: Requirements 3.2**

### Property 8: Refund Calculation Accuracy
*For any* refund calculation, the refund amount SHALL NOT exceed the original payment amount, and the sum of refund amount, platform fee, and player amount SHALL equal the original amount.
**Validates: Requirements 3.3**

### Property 9: Order Statistics Computation Accuracy
*For any* set of orders, the computed statistics SHALL accurately reflect total revenue (sum of all order amounts), order counts by status, and completion rate.
**Validates: Requirements 3.5**

### Property 10: Player Verification Workflow Enforcement
*For any* player verification status change, the service SHALL enforce valid state transitions (pending → verified/rejected) and SHALL reject invalid transitions.
**Validates: Requirements 4.2**

### Property 11: Player Earnings Calculation Accuracy
*For any* player earnings calculation, the net earnings SHALL equal gross amount minus commission, and commission SHALL be calculated using the applicable commission rate.
**Validates: Requirements 4.3**

### Property 12: Player Statistics Computation Accuracy
*For any* player statistics computation, the total earnings SHALL equal the sum of net earnings from all completed orders, and rating average SHALL be correctly computed.
**Validates: Requirements 4.4**

### Property 13: File Format Validation
*For any* import file upload, the service SHALL accept only .xlsx, .xls, and .csv files up to 10MB, rejecting others with appropriate error messages.
**Validates: Requirements 5.1**

### Property 14: Import Structure Validation
*For any* import file, the service SHALL validate that all required columns are present and report missing columns as structural errors.
**Validates: Requirements 5.2**

### Property 15: Import Data Validation Completeness
*For any* import data row, the validation SHALL check all fields against business rules and collect ALL validation errors for the row (not stopping at the first error).
**Validates: Requirements 5.3**

### Property 16: Import Summary Accuracy
*For any* import operation, the summary SHALL accurately report total rows, successful imports, and failed rows, where total = successful + failed.
**Validates: Requirements 5.5, 6.5**

### Property 17: Import Duplicate Detection
*For any* import data containing duplicates (emails for users, user references for players, keys for games), the validation SHALL detect and report all duplicates.
**Validates: Requirements 6.2, 7.2, 8.2**

### Property 18: Password Generation Security
*For any* generated temporary password, the password SHALL meet minimum security requirements (length ≥ 8, contains uppercase, lowercase, number, and special character).
**Validates: Requirements 6.3**

### Property 19: Error Detail Preservation
*For any* import row with validation errors, the error details SHALL include the row number, field name, and specific error message, and these details SHALL be preserved for later retrieval.
**Validates: Requirements 6.4, 9.3, 9.4**

### Property 20: Player Import Initial State
*For any* successfully imported player, the initial verification status SHALL be 'pending'.
**Validates: Requirements 7.3**

### Property 21: Skill Tag Parsing
*For any* skill tags input string, the parser SHALL correctly split comma-separated values and trim whitespace from each tag.
**Validates: Requirements 7.4**

### Property 22: Game Import Defaults
*For any* imported game without explicit isActive value, the default SHALL be true, and without explicit sortOrder, the default SHALL be 0.
**Validates: Requirements 8.3**

### Property 23: Import Metadata Recording
*For any* completed import operation, the metadata record SHALL contain timestamp, user ID, file name, record counts, and status.
**Validates: Requirements 9.1**

### Property 24: Import Result Report Format
*For any* import result download, the report SHALL contain all original data columns plus status and error columns.
**Validates: Requirements 9.5**

## Error Handling

### Service Layer Error Handling Strategy

```typescript
// 1. API errors are caught and wrapped in ServiceException
try {
  const response = await this.api.getUsers(params);
  return { success: true, data: response.data.data };
} catch (error) {
  throw this.handleError(error, 'Failed to fetch users');
}

// 2. Validation errors are returned as ServiceResult with error
const validation = this.validateUserData(data);
if (!validation.valid) {
  return {
    success: false,
    error: {
      code: ServiceErrorCodes.VALIDATION_ERROR,
      message: 'User data validation failed',
      details: { errors: validation.errors },
    },
  };
}

// 3. Batch operations collect all errors
const results: BatchResult<void> = {
  success: true,
  total: items.length,
  succeeded: 0,
  failed: 0,
  results: [],
};

for (const [index, item] of items.entries()) {
  try {
    await this.processItem(item);
    results.succeeded++;
    results.results.push({ index, success: true });
  } catch (error) {
    results.failed++;
    results.success = false;
    results.results.push({
      index,
      success: false,
      error: this.handleError(error, `Item ${index}`).toError(),
    });
  }
}
```

### Import Error Handling

```typescript
// Structure validation errors
if (missingColumns.length > 0) {
  return {
    totalRows: 0,
    validRows: [],
    invalidRows: [],
    structureErrors: missingColumns.map(col => `Missing required column: ${col}`),
  };
}

// Row validation errors - collect all errors per row
for (const [index, row] of rows.entries()) {
  const errors: Array<{ field: string; message: string }> = [];
  
  for (const column of template.columns) {
    const value = row[column.key];
    if (column.required && !value) {
      errors.push({ field: column.key, message: `${column.label} is required` });
    }
    if (column.validation && !column.validation(value)) {
      errors.push({ field: column.key, message: `${column.label} is invalid` });
    }
  }
  
  parsedRows.push({
    rowNumber: index + 2, // +2 for header row and 0-index
    data: row,
    errors,
    isValid: errors.length === 0,
  });
}
```

## Testing Strategy

### Unit Testing

使用 Vitest 进行单元测试，测试覆盖：

1. **Service 方法测试**
   - 成功场景
   - 错误场景
   - 边界条件

2. **验证逻辑测试**
   - 邮箱格式验证
   - 手机号格式验证
   - 密码强度验证
   - 业务规则验证

3. **数据转换测试**
   - 导出数据格式
   - 导入数据解析

### Property-Based Testing

使用 fast-check 进行属性测试，验证：

1. **错误格式一致性** - 生成随机错误场景，验证错误对象格式
2. **批量操作结果完整性** - 生成随机批量数据，验证结果数量匹配
3. **验证规则正确性** - 生成随机输入数据，验证验证结果
4. **计算准确性** - 生成随机订单/收益数据，验证计算结果

### Mock Utilities

```typescript
// services/__mocks__/apiMock.ts
export const createMockApi = () => ({
  getUsers: vi.fn(),
  createUser: vi.fn(),
  updateUser: vi.fn(),
  deleteUser: vi.fn(),
  // ... other methods
});

// Usage in tests
const mockApi = createMockApi();
const userService = new UserServiceImpl({ api: mockApi });

mockApi.getUsers.mockResolvedValue({
  data: { success: true, data: [mockUser] },
});

const result = await userService.getUsers({});
expect(result.success).toBe(true);
```

### Test Coverage Requirements

- Service 层代码覆盖率 ≥ 80%
- 所有公共方法必须有测试
- 所有验证规则必须有测试
- 所有错误场景必须有测试
