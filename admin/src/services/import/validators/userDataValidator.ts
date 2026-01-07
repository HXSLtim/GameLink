/**
 * User Data Validator
 * Validates user import data against business rules
 *
 * @module services/import/validators/userDataValidator
 */

import { validateUserEmail, validateUserPhone } from '../templates/userTemplate';

/**
 * Validation error for a specific field
 */
export interface FieldError {
  field: string;
  message: string;
}

/**
 * Result of validating a single row
 */
export interface RowValidationResult {
  rowNumber: number;
  valid: boolean;
  errors: FieldError[];
  data: Record<string, unknown>;
}

/**
 * Result of validating all rows
 */
export interface UserDataValidationResult {
  valid: boolean;
  totalRows: number;
  validRows: RowValidationResult[];
  invalidRows: RowValidationResult[];
  duplicateEmails: Map<string, number[]>;
  duplicatePhones: Map<string, number[]>;
}

/**
 * Options for user data validation
 */
export interface UserDataValidationOptions {
  /** Existing emails in the database (for uniqueness check) */
  existingEmails?: Set<string>;
  /** Existing phones in the database (for uniqueness check) */
  existingPhones?: Set<string>;
  /** Whether to check for duplicates within the import data */
  checkInternalDuplicates?: boolean;
}

/**
 * Valid user roles
 */
const VALID_ROLES = ['user', 'player', 'admin'];

/**
 * Valid user statuses
 */
const VALID_STATUSES = ['active', 'banned', 'suspended'];

/**
 * Validate a single user data row
 *
 * @param data - Row data to validate
 * @param rowNumber - Row number for error reporting
 * @param options - Validation options
 * @returns RowValidationResult with all validation errors
 */
export function validateUserRow(
  data: Record<string, unknown>,
  rowNumber: number,
  options: UserDataValidationOptions = {}
): RowValidationResult {
  const errors: FieldError[] = [];

  // Validate name (required)
  const name = data.name;
  if (name === undefined || name === null || String(name).trim() === '') {
    errors.push({ field: 'name', message: '姓名不能为空' });
  } else if (String(name).trim().length < 2) {
    errors.push({ field: 'name', message: '姓名至少需要2个字符' });
  } else if (String(name).trim().length > 50) {
    errors.push({ field: 'name', message: '姓名不能超过50个字符' });
  }

  // Validate email (required)
  const email = data.email;
  if (email === undefined || email === null || String(email).trim() === '') {
    errors.push({ field: 'email', message: '邮箱不能为空' });
  } else {
    const emailStr = String(email).trim().toLowerCase();
    if (!validateUserEmail(emailStr)) {
      errors.push({ field: 'email', message: '邮箱格式不正确' });
    } else if (options.existingEmails?.has(emailStr)) {
      errors.push({ field: 'email', message: '该邮箱已存在于系统中' });
    }
  }

  // Validate phone (required)
  const phone = data.phone;
  if (phone === undefined || phone === null || String(phone).trim() === '') {
    errors.push({ field: 'phone', message: '手机号不能为空' });
  } else {
    const phoneStr = String(phone).trim();
    if (!validateUserPhone(phoneStr)) {
      errors.push({ field: 'phone', message: '手机号格式不正确，需要是有效的中国手机号' });
    } else if (options.existingPhones?.has(phoneStr)) {
      errors.push({ field: 'phone', message: '该手机号已存在于系统中' });
    }
  }

  // Validate role (optional)
  const role = data.role;
  if (role !== undefined && role !== null && String(role).trim() !== '') {
    const roleStr = String(role).trim().toLowerCase();
    if (!VALID_ROLES.includes(roleStr)) {
      errors.push({
        field: 'role',
        message: `角色无效，有效值为: ${VALID_ROLES.join(', ')}`,
      });
    }
  }

  // Validate status (optional)
  const status = data.status;
  if (status !== undefined && status !== null && String(status).trim() !== '') {
    const statusStr = String(status).trim().toLowerCase();
    if (!VALID_STATUSES.includes(statusStr)) {
      errors.push({
        field: 'status',
        message: `状态无效，有效值为: ${VALID_STATUSES.join(', ')}`,
      });
    }
  }

  return {
    rowNumber,
    valid: errors.length === 0,
    errors,
    data,
  };
}

/**
 * Find duplicate values in import data
 *
 * @param rows - Array of row data
 * @param field - Field to check for duplicates
 * @returns Map of duplicate values to row numbers
 */
export function findDuplicates(
  rows: Array<{ rowNumber: number; data: Record<string, unknown> }>,
  field: string
): Map<string, number[]> {
  const valueToRows = new Map<string, number[]>();

  for (const row of rows) {
    const value = row.data[field];
    if (value !== undefined && value !== null && String(value).trim() !== '') {
      const normalizedValue = String(value).trim().toLowerCase();
      const existing = valueToRows.get(normalizedValue) || [];
      existing.push(row.rowNumber);
      valueToRows.set(normalizedValue, existing);
    }
  }

  // Filter to only duplicates (more than one occurrence)
  const duplicates = new Map<string, number[]>();
  for (const [value, rowNumbers] of valueToRows) {
    if (rowNumbers.length > 1) {
      duplicates.set(value, rowNumbers);
    }
  }

  return duplicates;
}

/**
 * Validate all user data rows
 *
 * @param rows - Array of row data with row numbers
 * @param options - Validation options
 * @returns UserDataValidationResult with all validation results
 */
export function validateUserData(
  rows: Array<{ rowNumber: number; data: Record<string, unknown> }>,
  options: UserDataValidationOptions = {}
): UserDataValidationResult {
  const validRows: RowValidationResult[] = [];
  const invalidRows: RowValidationResult[] = [];

  // First pass: validate each row individually
  for (const row of rows) {
    const result = validateUserRow(row.data, row.rowNumber, options);
    if (result.valid) {
      validRows.push(result);
    } else {
      invalidRows.push(result);
    }
  }

  // Check for internal duplicates if enabled
  let duplicateEmails = new Map<string, number[]>();
  let duplicatePhones = new Map<string, number[]>();

  if (options.checkInternalDuplicates !== false) {
    duplicateEmails = findDuplicates(rows, 'email');
    duplicatePhones = findDuplicates(rows, 'phone');

    // Add duplicate errors to affected rows
    for (const [email, rowNumbers] of duplicateEmails) {
      for (const rowNum of rowNumbers) {
        // Find the row in validRows or invalidRows
        const validIndex = validRows.findIndex((r) => r.rowNumber === rowNum);
        if (validIndex !== -1) {
          const row = validRows[validIndex];
          row.errors.push({
            field: 'email',
            message: `邮箱 "${email}" 在导入数据中重复出现 (行: ${rowNumbers.join(', ')})`,
          });
          row.valid = false;
          // Move to invalid rows
          validRows.splice(validIndex, 1);
          invalidRows.push(row);
        } else {
          // Already in invalid rows, just add the error
          const invalidRow = invalidRows.find((r) => r.rowNumber === rowNum);
          if (invalidRow) {
            // Check if this error already exists
            const hasError = invalidRow.errors.some(
              (e) => e.field === 'email' && e.message.includes('重复出现')
            );
            if (!hasError) {
              invalidRow.errors.push({
                field: 'email',
                message: `邮箱 "${email}" 在导入数据中重复出现 (行: ${rowNumbers.join(', ')})`,
              });
            }
          }
        }
      }
    }

    for (const [phone, rowNumbers] of duplicatePhones) {
      for (const rowNum of rowNumbers) {
        // Find the row in validRows or invalidRows
        const validIndex = validRows.findIndex((r) => r.rowNumber === rowNum);
        if (validIndex !== -1) {
          const row = validRows[validIndex];
          row.errors.push({
            field: 'phone',
            message: `手机号 "${phone}" 在导入数据中重复出现 (行: ${rowNumbers.join(', ')})`,
          });
          row.valid = false;
          // Move to invalid rows
          validRows.splice(validIndex, 1);
          invalidRows.push(row);
        } else {
          // Already in invalid rows, just add the error
          const invalidRow = invalidRows.find((r) => r.rowNumber === rowNum);
          if (invalidRow) {
            // Check if this error already exists
            const hasError = invalidRow.errors.some(
              (e) => e.field === 'phone' && e.message.includes('重复出现')
            );
            if (!hasError) {
              invalidRow.errors.push({
                field: 'phone',
                message: `手机号 "${phone}" 在导入数据中重复出现 (行: ${rowNumbers.join(', ')})`,
              });
            }
          }
        }
      }
    }
  }

  // Sort invalid rows by row number
  invalidRows.sort((a, b) => a.rowNumber - b.rowNumber);
  validRows.sort((a, b) => a.rowNumber - b.rowNumber);

  return {
    valid: invalidRows.length === 0,
    totalRows: rows.length,
    validRows,
    invalidRows,
    duplicateEmails,
    duplicatePhones,
  };
}

/**
 * Normalize user data for import
 * Applies default values and normalizes field values
 *
 * @param data - Raw row data
 * @returns Normalized data ready for import
 */
export function normalizeUserData(data: Record<string, unknown>): Record<string, unknown> {
  const normalized: Record<string, unknown> = { ...data };

  // Normalize name
  if (normalized.name) {
    normalized.name = String(normalized.name).trim();
  }

  // Normalize email (lowercase)
  if (normalized.email) {
    normalized.email = String(normalized.email).trim().toLowerCase();
  }

  // Normalize phone
  if (normalized.phone) {
    normalized.phone = String(normalized.phone).trim();
  }

  // Normalize role (lowercase, default to 'user')
  if (normalized.role && String(normalized.role).trim() !== '') {
    normalized.role = String(normalized.role).trim().toLowerCase();
  } else {
    normalized.role = 'user';
  }

  // Normalize status (lowercase, default to 'active')
  if (normalized.status && String(normalized.status).trim() !== '') {
    normalized.status = String(normalized.status).trim().toLowerCase();
  } else {
    normalized.status = 'active';
  }

  return normalized;
}
