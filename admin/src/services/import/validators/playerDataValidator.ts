/**
 * Player Data Validator
 * Validates player import data against business rules
 *
 * @module services/import/validators/playerDataValidator
 */

import { validateUserEmail } from '../templates/userTemplate';
import { parseSkillTags, validateSkillTags } from '../templates/playerTemplate';
import type { FieldError, RowValidationResult } from './userDataValidator';

/**
 * Result of validating all player rows
 */
export interface PlayerDataValidationResult {
  valid: boolean;
  totalRows: number;
  validRows: RowValidationResult[];
  invalidRows: RowValidationResult[];
  duplicateUserEmails: Map<string, number[]>;
  missingUsers: string[];
  existingPlayers: string[];
}

/**
 * Options for player data validation
 */
export interface PlayerDataValidationOptions {
  /** Set of existing user emails in the database */
  existingUserEmails?: Set<string>;
  /** Set of user emails that are already registered as players */
  existingPlayerEmails?: Set<string>;
  /** Whether to check for duplicates within the import data */
  checkInternalDuplicates?: boolean;
  /** Allowed skill tags (if empty, all tags are allowed) */
  allowedSkillTags?: Set<string>;
}

/**
 * Validate a single player data row
 *
 * @param data - Row data to validate
 * @param rowNumber - Row number for error reporting
 * @param options - Validation options
 * @returns RowValidationResult with all validation errors
 */
export function validatePlayerRow(
  data: Record<string, unknown>,
  rowNumber: number,
  options: PlayerDataValidationOptions = {}
): RowValidationResult {
  const errors: FieldError[] = [];

  // Validate userEmail (required)
  const userEmail = data.userEmail;
  if (userEmail === undefined || userEmail === null || String(userEmail).trim() === '') {
    errors.push({ field: 'userEmail', message: '用户邮箱不能为空' });
  } else {
    const emailStr = String(userEmail).trim().toLowerCase();
    if (!validateUserEmail(emailStr)) {
      errors.push({ field: 'userEmail', message: '用户邮箱格式不正确' });
    } else {
      // Check if user exists
      if (options.existingUserEmails && !options.existingUserEmails.has(emailStr)) {
        errors.push({ field: 'userEmail', message: '该用户邮箱不存在于系统中' });
      }
      // Check if user is already a player
      if (options.existingPlayerEmails?.has(emailStr)) {
        errors.push({ field: 'userEmail', message: '该用户已经是陪玩师' });
      }
    }
  }

  // Validate nickname (optional, max 50 chars)
  const nickname = data.nickname;
  if (nickname !== undefined && nickname !== null && String(nickname).trim() !== '') {
    const nicknameStr = String(nickname).trim();
    if (nicknameStr.length > 50) {
      errors.push({ field: 'nickname', message: '昵称不能超过50个字符' });
    }
  }

  // Validate bio (optional, max 500 chars)
  const bio = data.bio;
  if (bio !== undefined && bio !== null && String(bio).trim() !== '') {
    const bioStr = String(bio).trim();
    if (bioStr.length > 500) {
      errors.push({ field: 'bio', message: '简介不能超过500个字符' });
    }
  }

  // Validate hourlyRate (optional, must be non-negative number)
  const hourlyRate = data.hourlyRate;
  if (hourlyRate !== undefined && hourlyRate !== null && String(hourlyRate).trim() !== '') {
    const rateNum = typeof hourlyRate === 'number' ? hourlyRate : parseFloat(String(hourlyRate));
    if (isNaN(rateNum)) {
      errors.push({ field: 'hourlyRate', message: '时薪必须是有效的数字' });
    } else if (rateNum < 0) {
      errors.push({ field: 'hourlyRate', message: '时薪不能为负数' });
    } else if (rateNum > 10000) {
      errors.push({ field: 'hourlyRate', message: '时薪不能超过10000元' });
    }
  }

  // Validate mainGame (optional)
  const mainGame = data.mainGame;
  if (mainGame !== undefined && mainGame !== null && String(mainGame).trim() !== '') {
    const gameStr = String(mainGame).trim();
    if (gameStr.length > 100) {
      errors.push({ field: 'mainGame', message: '主游戏名称不能超过100个字符' });
    }
  }

  // Validate skillTags (optional)
  const skillTags = data.skillTags;
  if (skillTags !== undefined && skillTags !== null && String(skillTags).trim() !== '') {
    const tagsStr = String(skillTags).trim();
    const parsedTags = parseSkillTags(tagsStr);
    const tagValidation = validateSkillTags(parsedTags);

    if (!tagValidation.valid) {
      errors.push({
        field: 'skillTags',
        message: `技能标签格式无效: ${tagValidation.invalidTags.join(', ')}`,
      });
    }

    // Check against allowed tags if provided
    if (options.allowedSkillTags && options.allowedSkillTags.size > 0) {
      const invalidTags = parsedTags.filter((tag) => !options.allowedSkillTags!.has(tag));
      if (invalidTags.length > 0) {
        errors.push({
          field: 'skillTags',
          message: `不允许的技能标签: ${invalidTags.join(', ')}`,
        });
      }
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
 * Find duplicate user emails in import data
 *
 * @param rows - Array of row data
 * @returns Map of duplicate emails to row numbers
 */
export function findDuplicateUserEmails(
  rows: Array<{ rowNumber: number; data: Record<string, unknown> }>
): Map<string, number[]> {
  const emailToRows = new Map<string, number[]>();

  for (const row of rows) {
    const email = row.data.userEmail;
    if (email !== undefined && email !== null && String(email).trim() !== '') {
      const normalizedEmail = String(email).trim().toLowerCase();
      const existing = emailToRows.get(normalizedEmail) || [];
      existing.push(row.rowNumber);
      emailToRows.set(normalizedEmail, existing);
    }
  }

  // Filter to only duplicates (more than one occurrence)
  const duplicates = new Map<string, number[]>();
  for (const [email, rowNumbers] of emailToRows) {
    if (rowNumbers.length > 1) {
      duplicates.set(email, rowNumbers);
    }
  }

  return duplicates;
}

/**
 * Validate all player data rows
 *
 * @param rows - Array of row data with row numbers
 * @param options - Validation options
 * @returns PlayerDataValidationResult with all validation results
 */
export function validatePlayerData(
  rows: Array<{ rowNumber: number; data: Record<string, unknown> }>,
  options: PlayerDataValidationOptions = {}
): PlayerDataValidationResult {
  const validRows: RowValidationResult[] = [];
  const invalidRows: RowValidationResult[] = [];
  const missingUsers: string[] = [];
  const existingPlayers: string[] = [];

  // First pass: validate each row individually
  for (const row of rows) {
    const result = validatePlayerRow(row.data, row.rowNumber, options);

    // Track missing users and existing players
    const email = row.data.userEmail;
    if (email && String(email).trim() !== '') {
      const emailStr = String(email).trim().toLowerCase();
      if (options.existingUserEmails && !options.existingUserEmails.has(emailStr)) {
        if (!missingUsers.includes(emailStr)) {
          missingUsers.push(emailStr);
        }
      }
      if (options.existingPlayerEmails?.has(emailStr)) {
        if (!existingPlayers.includes(emailStr)) {
          existingPlayers.push(emailStr);
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
  let duplicateUserEmails = new Map<string, number[]>();

  if (options.checkInternalDuplicates !== false) {
    duplicateUserEmails = findDuplicateUserEmails(rows);

    // Add duplicate errors to affected rows
    for (const [email, rowNumbers] of duplicateUserEmails) {
      for (const rowNum of rowNumbers) {
        // Find the row in validRows or invalidRows
        const validIndex = validRows.findIndex((r) => r.rowNumber === rowNum);
        if (validIndex !== -1) {
          const row = validRows[validIndex];
          row.errors.push({
            field: 'userEmail',
            message: `用户邮箱 "${email}" 在导入数据中重复出现 (行: ${rowNumbers.join(', ')})`,
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
              (e) => e.field === 'userEmail' && e.message.includes('重复出现')
            );
            if (!hasError) {
              invalidRow.errors.push({
                field: 'userEmail',
                message: `用户邮箱 "${email}" 在导入数据中重复出现 (行: ${rowNumbers.join(', ')})`,
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
    duplicateUserEmails,
    missingUsers,
    existingPlayers,
  };
}

/**
 * Normalize player data for import
 * Applies default values and normalizes field values
 *
 * @param data - Raw row data
 * @returns Normalized data ready for import
 */
export function normalizePlayerData(data: Record<string, unknown>): Record<string, unknown> {
  const normalized: Record<string, unknown> = { ...data };

  // Normalize userEmail (lowercase)
  if (normalized.userEmail) {
    normalized.userEmail = String(normalized.userEmail).trim().toLowerCase();
  }

  // Normalize nickname
  if (normalized.nickname && String(normalized.nickname).trim() !== '') {
    normalized.nickname = String(normalized.nickname).trim();
  } else {
    normalized.nickname = '';
  }

  // Normalize bio
  if (normalized.bio && String(normalized.bio).trim() !== '') {
    normalized.bio = String(normalized.bio).trim();
  } else {
    normalized.bio = '';
  }

  // Normalize hourlyRate (convert to number, default to 0)
  if (normalized.hourlyRate !== undefined && normalized.hourlyRate !== null && String(normalized.hourlyRate).trim() !== '') {
    const rate = typeof normalized.hourlyRate === 'number'
      ? normalized.hourlyRate
      : parseFloat(String(normalized.hourlyRate));
    normalized.hourlyRate = isNaN(rate) ? 0 : rate;
  } else {
    normalized.hourlyRate = 0;
  }

  // Normalize mainGame
  if (normalized.mainGame && String(normalized.mainGame).trim() !== '') {
    normalized.mainGame = String(normalized.mainGame).trim();
  } else {
    normalized.mainGame = '';
  }

  // Normalize skillTags (parse and rejoin)
  if (normalized.skillTags && String(normalized.skillTags).trim() !== '') {
    const tags = parseSkillTags(String(normalized.skillTags));
    normalized.skillTags = tags;
    normalized.skillTagsString = tags.join(',');
  } else {
    normalized.skillTags = [];
    normalized.skillTagsString = '';
  }

  // Set initial verification status to pending
  normalized.verificationStatus = 'pending';

  return normalized;
}
