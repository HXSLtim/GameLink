/**
 * Game Data Validator
 * Validates game import data against business rules
 *
 * @module services/import/validators/gameDataValidator
 */

import { validateGameKey, parseBoolean } from '../templates/gameTemplate';
import type { FieldError, RowValidationResult } from './userDataValidator';

/**
 * Result of validating all game rows
 */
export interface GameDataValidationResult {
  valid: boolean;
  totalRows: number;
  validRows: RowValidationResult[];
  invalidRows: RowValidationResult[];
  duplicateKeys: Map<string, number[]>;
  invalidCategories: string[];
}

/**
 * Options for game data validation
 */
export interface GameDataValidationOptions {
  /** Set of existing game keys in the database */
  existingGameKeys?: Set<string>;
  /** Set of valid category names/IDs */
  validCategories?: Set<string>;
  /** Whether to check for duplicates within the import data */
  checkInternalDuplicates?: boolean;
  /** How to handle existing keys: 'skip', 'update', or 'fail' */
  duplicateKeyHandling?: 'skip' | 'update' | 'fail';
}

/**
 * Validate a single game data row
 *
 * @param data - Row data to validate
 * @param rowNumber - Row number for error reporting
 * @param options - Validation options
 * @returns RowValidationResult with all validation errors
 */
export function validateGameRow(
  data: Record<string, unknown>,
  rowNumber: number,
  options: GameDataValidationOptions = {}
): RowValidationResult {
  const errors: FieldError[] = [];

  // Validate key (required)
  const key = data.key;
  if (key === undefined || key === null || String(key).trim() === '') {
    errors.push({ field: 'key', message: '游戏标识不能为空' });
  } else {
    const keyStr = String(key).trim();
    if (!validateGameKey(keyStr)) {
      errors.push({
        field: 'key',
        message: '游戏标识格式不正确，只能包含字母、数字、下划线和连字符',
      });
    } else if (keyStr.length > 50) {
      errors.push({ field: 'key', message: '游戏标识不能超过50个字符' });
    } else {
      // Check if key already exists (based on handling mode)
      const normalizedKey = keyStr.toLowerCase();
      if (options.existingGameKeys?.has(normalizedKey)) {
        const handling = options.duplicateKeyHandling || 'fail';
        if (handling === 'fail') {
          errors.push({ field: 'key', message: '该游戏标识已存在于系统中' });
        }
        // For 'skip' and 'update', we don't add an error
      }
    }
  }

  // Validate name (required)
  const name = data.name;
  if (name === undefined || name === null || String(name).trim() === '') {
    errors.push({ field: 'name', message: '游戏名称不能为空' });
  } else {
    const nameStr = String(name).trim();
    if (nameStr.length > 100) {
      errors.push({ field: 'name', message: '游戏名称不能超过100个字符' });
    }
  }

  // Validate category (optional)
  const category = data.category;
  if (category !== undefined && category !== null && String(category).trim() !== '') {
    const categoryStr = String(category).trim();
    if (categoryStr.length > 50) {
      errors.push({ field: 'category', message: '分类名称不能超过50个字符' });
    }
    // Check against valid categories if provided
    if (options.validCategories && options.validCategories.size > 0) {
      const normalizedCategory = categoryStr.toLowerCase();
      if (!options.validCategories.has(normalizedCategory)) {
        errors.push({
          field: 'category',
          message: `分类 "${categoryStr}" 不存在于系统中`,
        });
      }
    }
  }

  // Validate description (optional)
  const description = data.description;
  if (description !== undefined && description !== null && String(description).trim() !== '') {
    const descStr = String(description).trim();
    if (descStr.length > 500) {
      errors.push({ field: 'description', message: '描述不能超过500个字符' });
    }
  }

  // Validate isActive (optional, must be boolean-like)
  const isActive = data.isActive;
  if (isActive !== undefined && isActive !== null && String(isActive).trim() !== '') {
    const strValue = String(isActive).toLowerCase().trim();
    const validBooleanValues = ['true', 'false', '是', '否', '1', '0', 'yes', 'no', 'y', 'n'];
    if (typeof isActive !== 'boolean' && !validBooleanValues.includes(strValue)) {
      errors.push({
        field: 'isActive',
        message: '是否启用必须是有效的布尔值 (true/false, 是/否, 1/0)',
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
 * Find duplicate game keys in import data
 *
 * @param rows - Array of row data
 * @returns Map of duplicate keys to row numbers
 */
export function findDuplicateGameKeys(
  rows: Array<{ rowNumber: number; data: Record<string, unknown> }>
): Map<string, number[]> {
  const keyToRows = new Map<string, number[]>();

  for (const row of rows) {
    const key = row.data.key;
    if (key !== undefined && key !== null && String(key).trim() !== '') {
      const normalizedKey = String(key).trim().toLowerCase();
      const existing = keyToRows.get(normalizedKey) || [];
      existing.push(row.rowNumber);
      keyToRows.set(normalizedKey, existing);
    }
  }

  // Filter to only duplicates (more than one occurrence)
  const duplicates = new Map<string, number[]>();
  for (const [key, rowNumbers] of keyToRows) {
    if (rowNumbers.length > 1) {
      duplicates.set(key, rowNumbers);
    }
  }

  return duplicates;
}

/**
 * Validate all game data rows
 *
 * @param rows - Array of row data with row numbers
 * @param options - Validation options
 * @returns GameDataValidationResult with all validation results
 */
export function validateGameData(
  rows: Array<{ rowNumber: number; data: Record<string, unknown> }>,
  options: GameDataValidationOptions = {}
): GameDataValidationResult {
  const validRows: RowValidationResult[] = [];
  const invalidRows: RowValidationResult[] = [];
  const invalidCategories: string[] = [];

  // First pass: validate each row individually
  for (const row of rows) {
    const result = validateGameRow(row.data, row.rowNumber, options);

    // Track invalid categories
    const category = row.data.category;
    if (category && String(category).trim() !== '') {
      const categoryStr = String(category).trim();
      if (options.validCategories && options.validCategories.size > 0) {
        const normalizedCategory = categoryStr.toLowerCase();
        if (!options.validCategories.has(normalizedCategory)) {
          if (!invalidCategories.includes(categoryStr)) {
            invalidCategories.push(categoryStr);
          }
        }
      }
    }

    if (result.valid) {
      validRows.push(result);
    } else {
      invalidRows.push(result);
    }
  }

  // Check for internal duplicates if enabled
  let duplicateKeys = new Map<string, number[]>();

  if (options.checkInternalDuplicates !== false) {
    duplicateKeys = findDuplicateGameKeys(rows);

    // Add duplicate errors to affected rows
    for (const [key, rowNumbers] of duplicateKeys) {
      for (const rowNum of rowNumbers) {
        // Find the row in validRows or invalidRows
        const validIndex = validRows.findIndex((r) => r.rowNumber === rowNum);
        if (validIndex !== -1) {
          const row = validRows[validIndex];
          row.errors.push({
            field: 'key',
            message: `游戏标识 "${key}" 在导入数据中重复出现 (行: ${rowNumbers.join(', ')})`,
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
              (e) => e.field === 'key' && e.message.includes('重复出现')
            );
            if (!hasError) {
              invalidRow.errors.push({
                field: 'key',
                message: `游戏标识 "${key}" 在导入数据中重复出现 (行: ${rowNumbers.join(', ')})`,
              });
            }
          }
        }
      }
    }
  }

  // Sort rows by row number
  invalidRows.sort((a, b) => a.rowNumber - b.rowNumber);
  validRows.sort((a, b) => a.rowNumber - b.rowNumber);

  return {
    valid: invalidRows.length === 0,
    totalRows: rows.length,
    validRows,
    invalidRows,
    duplicateKeys,
    invalidCategories,
  };
}

/**
 * Normalize game data for import
 * Applies default values and normalizes field values
 *
 * @param data - Raw row data
 * @returns Normalized data ready for import
 */
export function normalizeGameData(data: Record<string, unknown>): Record<string, unknown> {
  const normalized: Record<string, unknown> = { ...data };

  // Normalize key (lowercase)
  if (normalized.key) {
    normalized.key = String(normalized.key).trim().toLowerCase();
  }

  // Normalize name
  if (normalized.name) {
    normalized.name = String(normalized.name).trim();
  }

  // Normalize category
  if (normalized.category && String(normalized.category).trim() !== '') {
    normalized.category = String(normalized.category).trim();
  } else {
    normalized.category = '';
  }

  // Normalize description
  if (normalized.description && String(normalized.description).trim() !== '') {
    normalized.description = String(normalized.description).trim();
  } else {
    normalized.description = '';
  }

  // Normalize isActive (default to true)
  normalized.isActive = parseBoolean(normalized.isActive);

  // Set default sortOrder
  if (normalized.sortOrder === undefined || normalized.sortOrder === null) {
    normalized.sortOrder = 0;
  } else {
    const sortOrder = typeof normalized.sortOrder === 'number'
      ? normalized.sortOrder
      : parseInt(String(normalized.sortOrder), 10);
    normalized.sortOrder = isNaN(sortOrder) ? 0 : sortOrder;
  }

  return normalized;
}
