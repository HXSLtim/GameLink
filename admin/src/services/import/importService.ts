/**
 * Import Service
 * Core service for handling data imports from Excel/CSV files
 *
 * @module services/import/importService
 */

import { BaseService, type ServiceDependencies } from '../domain/base';
import type { ImportTemplate, ImportType } from './templates/types';
import { getTemplate } from './templates';
import { parseFile } from './parsers';
import type { ParseResult } from './parsers/types';
import {
  validateStructure,
  createColumnMapping,
  mapRowToTemplate,
  type StructureValidationResult,
} from './validators/structureValidator';
import {
  validateUserData,
  normalizeUserData,
} from './validators/userDataValidator';
import {
  validatePlayerData,
  normalizePlayerData,
} from './validators/playerDataValidator';
import {
  validateGameData,
  normalizeGameData,
} from './validators/gameDataValidator';
import { yuanToCents } from './templates/playerTemplate';
import type { CreateUserDto, CreatePlayerDto, CreateGameDto } from '@/api/admin';

/**
 * Parsed row with validation status
 */
export interface ParsedRow {
  rowNumber: number;
  data: Record<string, unknown>;
  errors: Array<{ field: string; message: string }>;
  isValid: boolean;
}

/**
 * Import preview result
 */
export interface ImportPreview {
  totalRows: number;
  validRows: ParsedRow[];
  invalidRows: ParsedRow[];
  structureErrors: string[];
}

/**
 * Import result summary
 */
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

/**
 * Options for duplicate key handling in game import
 */
export type DuplicateKeyHandling = 'skip' | 'update' | 'fail';

/**
 * Import options
 */
export interface ImportOptions {
  /** How to handle duplicate keys (for game import) */
  duplicateKeyHandling?: DuplicateKeyHandling;
  /** Progress callback */
  onProgress?: (completed: number, total: number) => void;
}

/**
 * Import Service Interface
 */
export interface IImportService {
  // Template methods
  getTemplate(type: ImportType): ImportTemplate;
  downloadTemplate(type: ImportType): Blob;

  // Parsing methods
  parseFile(file: File, type: ImportType): Promise<ImportPreview>;
  validateStructure(headers: string[], template: ImportTemplate): StructureValidationResult;

  // Import methods
  importUsers(rows: ParsedRow[], options?: ImportOptions): Promise<ImportResult>;
  importPlayers(rows: ParsedRow[], options?: ImportOptions): Promise<ImportResult>;
  importGames(rows: ParsedRow[], options?: ImportOptions): Promise<ImportResult>;
}


/**
 * Password generation configuration
 */
const PASSWORD_CONFIG = {
  length: 12,
  uppercase: 'ABCDEFGHIJKLMNOPQRSTUVWXYZ',
  lowercase: 'abcdefghijklmnopqrstuvwxyz',
  numbers: '0123456789',
  special: '!@#$%^&*',
};

/**
 * Generate a secure temporary password
 * Password meets requirements: length >= 8, contains uppercase, lowercase, number, and special character
 */
export function generateSecurePassword(): string {
  const { length, uppercase, lowercase, numbers, special } = PASSWORD_CONFIG;

  // Ensure at least one of each required character type
  const requiredChars = [
    uppercase[Math.floor(Math.random() * uppercase.length)],
    lowercase[Math.floor(Math.random() * lowercase.length)],
    numbers[Math.floor(Math.random() * numbers.length)],
    special[Math.floor(Math.random() * special.length)],
  ];

  // Fill remaining length with random characters from all pools
  const allChars = uppercase + lowercase + numbers + special;
  const remainingLength = length - requiredChars.length;
  const randomChars: string[] = [];

  for (let i = 0; i < remainingLength; i++) {
    randomChars.push(allChars[Math.floor(Math.random() * allChars.length)]);
  }

  // Combine and shuffle
  const allPasswordChars = [...requiredChars, ...randomChars];
  for (let i = allPasswordChars.length - 1; i > 0; i--) {
    const j = Math.floor(Math.random() * (i + 1));
    [allPasswordChars[i], allPasswordChars[j]] = [allPasswordChars[j], allPasswordChars[i]];
  }

  return allPasswordChars.join('');
}

/**
 * Validate password meets security requirements
 */
export function validatePasswordSecurity(password: string): {
  valid: boolean;
  errors: string[];
} {
  const errors: string[] = [];

  if (password.length < 8) {
    errors.push('Password must be at least 8 characters');
  }
  if (!/[A-Z]/.test(password)) {
    errors.push('Password must contain at least one uppercase letter');
  }
  if (!/[a-z]/.test(password)) {
    errors.push('Password must contain at least one lowercase letter');
  }
  if (!/[0-9]/.test(password)) {
    errors.push('Password must contain at least one number');
  }
  if (!/[!@#$%^&*]/.test(password)) {
    errors.push('Password must contain at least one special character (!@#$%^&*)');
  }

  return {
    valid: errors.length === 0,
    errors,
  };
}

/**
 * Import Service Implementation
 */
export class ImportService extends BaseService implements IImportService {
  constructor(deps: ServiceDependencies = {}) {
    super(deps);
  }

  /**
   * Get import template by type
   */
  getTemplate(type: ImportType): ImportTemplate {
    return getTemplate(type);
  }

  /**
   * Download import template as Excel file
   */
  downloadTemplate(type: ImportType): Blob {
    const template = this.getTemplate(type);

    // Create CSV content with headers and example rows
    const headers = template.columns.map((col) => col.labelZh);
    const rows: string[][] = [headers];

    // Add example rows if available
    if (template.exampleRows) {
      for (const exampleRow of template.exampleRows) {
        const row = template.columns.map((col) => {
          const value = exampleRow[col.key];
          return value !== undefined && value !== null ? String(value) : '';
        });
        rows.push(row);
      }
    }

    // Convert to CSV
    const csvContent = rows
      .map((row) =>
        row.map((cell) => {
          // Escape quotes and wrap in quotes if contains comma or quote
          if (cell.includes(',') || cell.includes('"') || cell.includes('\n')) {
            return `"${cell.replace(/"/g, '""')}"`;
          }
          return cell;
        }).join(',')
      )
      .join('\n');

    // Add BOM for Excel compatibility with Chinese characters
    const bom = '\uFEFF';
    return new Blob([bom + csvContent], { type: 'text/csv;charset=utf-8' });
  }

  /**
   * Parse and validate an import file
   */
  async parseFile(file: File, type: ImportType): Promise<ImportPreview> {
    const template = this.getTemplate(type);

    // Parse the file
    const parseResult: ParseResult = await parseFile(file);

    if (!parseResult.success) {
      return {
        totalRows: 0,
        validRows: [],
        invalidRows: [],
        structureErrors: [parseResult.error || 'Failed to parse file'],
      };
    }

    // Validate structure
    const structureResult = this.validateStructure(parseResult.headers, template);

    if (!structureResult.valid) {
      return {
        totalRows: parseResult.totalRows,
        validRows: [],
        invalidRows: [],
        structureErrors: structureResult.errors,
      };
    }

    // Create column mapping
    const mapping = createColumnMapping(parseResult.headers, template);

    // Map and validate each row
    const parsedRows: ParsedRow[] = [];

    for (let i = 0; i < parseResult.rows.length; i++) {
      const rawRow = parseResult.rows[i];
      const rowNumber = i + 2; // +2 for header row and 1-based indexing

      // Map row to template structure
      const mappedData = mapRowToTemplate(rawRow, parseResult.headers, mapping, template);

      parsedRows.push({
        rowNumber,
        data: mappedData,
        errors: [],
        isValid: true,
      });
    }

    // Validate data based on type
    const validationResult = this.validateDataByType(type, parsedRows);

    return {
      totalRows: parseResult.totalRows,
      validRows: validationResult.validRows,
      invalidRows: validationResult.invalidRows,
      structureErrors: structureResult.errors.filter((e) => !e.includes('将被忽略')),
    };
  }

  /**
   * Validate structure of headers against template
   */
  validateStructure(headers: string[], template: ImportTemplate): StructureValidationResult {
    return validateStructure(headers, template);
  }

  /**
   * Validate data based on import type
   */
  private validateDataByType(
    type: ImportType,
    rows: ParsedRow[]
  ): { validRows: ParsedRow[]; invalidRows: ParsedRow[] } {
    const rowsWithNumbers = rows.map((r) => ({
      rowNumber: r.rowNumber,
      data: r.data,
    }));

    switch (type) {
      case 'user': {
        const result = validateUserData(rowsWithNumbers, { checkInternalDuplicates: true });
        return {
          validRows: result.validRows.map((r) => ({
            rowNumber: r.rowNumber,
            data: r.data,
            errors: r.errors,
            isValid: r.valid,
          })),
          invalidRows: result.invalidRows.map((r) => ({
            rowNumber: r.rowNumber,
            data: r.data,
            errors: r.errors,
            isValid: r.valid,
          })),
        };
      }
      case 'player': {
        const result = validatePlayerData(rowsWithNumbers, { checkInternalDuplicates: true });
        return {
          validRows: result.validRows.map((r) => ({
            rowNumber: r.rowNumber,
            data: r.data,
            errors: r.errors,
            isValid: r.valid,
          })),
          invalidRows: result.invalidRows.map((r) => ({
            rowNumber: r.rowNumber,
            data: r.data,
            errors: r.errors,
            isValid: r.valid,
          })),
        };
      }
      case 'game': {
        const result = validateGameData(rowsWithNumbers, { checkInternalDuplicates: true });
        return {
          validRows: result.validRows.map((r) => ({
            rowNumber: r.rowNumber,
            data: r.data,
            errors: r.errors,
            isValid: r.valid,
          })),
          invalidRows: result.invalidRows.map((r) => ({
            rowNumber: r.rowNumber,
            data: r.data,
            errors: r.errors,
            isValid: r.valid,
          })),
        };
      }
      default:
        return { validRows: rows, invalidRows: [] };
    }
  }


  /**
   * Import users from parsed rows
   * Generates secure temporary passwords for new users
   */
  async importUsers(rows: ParsedRow[], options?: ImportOptions): Promise<ImportResult> {
    const result: ImportResult = {
      success: true,
      totalRows: rows.length,
      importedCount: 0,
      skippedCount: 0,
      errors: [],
    };

    const validRows = rows.filter((r) => r.isValid);
    const invalidRows = rows.filter((r) => !r.isValid);

    // Add errors from invalid rows
    for (const row of invalidRows) {
      result.skippedCount++;
      for (const error of row.errors) {
        result.errors.push({
          rowNumber: row.rowNumber,
          field: error.field,
          message: error.message,
        });
      }
    }

    // Process valid rows
    for (let i = 0; i < validRows.length; i++) {
      const row = validRows[i];

      try {
        // Normalize data
        const normalizedData = normalizeUserData(row.data);

        // Generate secure password
        const password = generateSecurePassword();

        // Create user DTO
        const createDto: CreateUserDto = {
          name: String(normalizedData.name),
          email: String(normalizedData.email),
          phone: String(normalizedData.phone),
          password,
          role: (normalizedData.role as 'user' | 'player' | 'admin') || 'user',
          status: (normalizedData.status as 'active' | 'banned' | 'suspended') || 'active',
        };

        // Call API to create user
        await this.api.createUser(createDto);
        result.importedCount++;

        // Report progress
        if (options?.onProgress) {
          options.onProgress(i + 1, validRows.length);
        }
      } catch (error) {
        result.success = false;
        result.skippedCount++;

        const errorMessage =
          error instanceof Error ? error.message : 'Unknown error occurred';
        result.errors.push({
          rowNumber: row.rowNumber,
          message: errorMessage,
        });
      }
    }

    // Update success status
    result.success = result.errors.length === 0;

    return result;
  }

  /**
   * Import players from parsed rows
   * Sets initial verification status to pending
   */
  async importPlayers(rows: ParsedRow[], options?: ImportOptions): Promise<ImportResult> {
    const result: ImportResult = {
      success: true,
      totalRows: rows.length,
      importedCount: 0,
      skippedCount: 0,
      errors: [],
    };

    const validRows = rows.filter((r) => r.isValid);
    const invalidRows = rows.filter((r) => !r.isValid);

    // Add errors from invalid rows
    for (const row of invalidRows) {
      result.skippedCount++;
      for (const error of row.errors) {
        result.errors.push({
          rowNumber: row.rowNumber,
          field: error.field,
          message: error.message,
        });
      }
    }

    // First, we need to look up user IDs by email
    // This would typically be done via an API call
    // For now, we'll assume the user lookup is handled elsewhere

    // Process valid rows
    for (let i = 0; i < validRows.length; i++) {
      const row = validRows[i];

      try {
        // Normalize data
        const normalizedData = normalizePlayerData(row.data);

        // Look up user by email to get userId
        const userEmail = String(normalizedData.userEmail);
        const userResponse = await this.api.getUsers({ keyword: userEmail, page_size: 1 });
        const users = userResponse.data?.data || [];
        const user = users.find(
          (u) => u.email.toLowerCase() === userEmail.toLowerCase()
        );

        if (!user) {
          result.skippedCount++;
          result.errors.push({
            rowNumber: row.rowNumber,
            field: 'userEmail',
            message: `User with email "${userEmail}" not found`,
          });
          continue;
        }

        // Convert hourly rate from yuan to cents
        const hourlyRateCents = normalizedData.hourlyRate
          ? yuanToCents(Number(normalizedData.hourlyRate))
          : 0;

        // Create player DTO with initial verification status as pending
        const createDto: CreatePlayerDto = {
          userId: user.id,
          nickname: normalizedData.nickname ? String(normalizedData.nickname) : undefined,
          bio: normalizedData.bio ? String(normalizedData.bio) : undefined,
          hourlyRateCents,
          verificationStatus: 'pending', // Always set to pending for new imports
        };

        // Call API to create player
        await this.api.createPlayer(createDto);
        result.importedCount++;

        // Report progress
        if (options?.onProgress) {
          options.onProgress(i + 1, validRows.length);
        }
      } catch (error) {
        result.success = false;
        result.skippedCount++;

        const errorMessage =
          error instanceof Error ? error.message : 'Unknown error occurred';
        result.errors.push({
          rowNumber: row.rowNumber,
          message: errorMessage,
        });
      }
    }

    // Update success status
    result.success = result.errors.length === 0;

    return result;
  }

  /**
   * Import games from parsed rows
   * Applies default values (isActive=true, sortOrder=0)
   * Handles duplicate key options (skip/update/fail)
   */
  async importGames(rows: ParsedRow[], options?: ImportOptions): Promise<ImportResult> {
    const result: ImportResult = {
      success: true,
      totalRows: rows.length,
      importedCount: 0,
      skippedCount: 0,
      errors: [],
    };

    const duplicateHandling = options?.duplicateKeyHandling || 'fail';
    const validRows = rows.filter((r) => r.isValid);
    const invalidRows = rows.filter((r) => !r.isValid);

    // Add errors from invalid rows
    for (const row of invalidRows) {
      result.skippedCount++;
      for (const error of row.errors) {
        result.errors.push({
          rowNumber: row.rowNumber,
          field: error.field,
          message: error.message,
        });
      }
    }

    // Get existing games to check for duplicates
    let existingGames: Array<{ id: number; key: string }> = [];
    try {
      const gamesResponse = await this.api.getGames({ page_size: 1000 });
      existingGames = (gamesResponse.data?.data || []).map((g) => ({
        id: g.id,
        key: g.key.toLowerCase(),
      }));
    } catch {
      // If we can't fetch existing games, proceed without duplicate check
      this.logger.warn('Could not fetch existing games for duplicate check');
    }

    const existingKeyMap = new Map(existingGames.map((g) => [g.key, g.id]));

    // Process valid rows
    for (let i = 0; i < validRows.length; i++) {
      const row = validRows[i];

      try {
        // Normalize data
        const normalizedData = normalizeGameData(row.data);
        const gameKey = String(normalizedData.key).toLowerCase();

        // Check for duplicate key
        const existingGameId = existingKeyMap.get(gameKey);

        if (existingGameId !== undefined) {
          switch (duplicateHandling) {
            case 'skip':
              result.skippedCount++;
              result.errors.push({
                rowNumber: row.rowNumber,
                field: 'key',
                message: `Game key "${gameKey}" already exists, skipped`,
              });
              continue;

            case 'update':
              // Update existing game
              await this.api.updateGame(existingGameId, {
                name: String(normalizedData.name),
                category: normalizedData.category
                  ? String(normalizedData.category)
                  : undefined,
                description: normalizedData.description
                  ? String(normalizedData.description)
                  : undefined,
                isActive: normalizedData.isActive as boolean,
                sortOrder: normalizedData.sortOrder as number,
              });
              result.importedCount++;
              break;

            case 'fail':
            default:
              result.skippedCount++;
              result.errors.push({
                rowNumber: row.rowNumber,
                field: 'key',
                message: `Game key "${gameKey}" already exists`,
              });
              continue;
          }
        } else {
          // Create new game with defaults
          const createDto: CreateGameDto = {
            key: gameKey,
            name: String(normalizedData.name),
            category: normalizedData.category
              ? String(normalizedData.category)
              : undefined,
            description: normalizedData.description
              ? String(normalizedData.description)
              : undefined,
            isActive: normalizedData.isActive !== undefined
              ? (normalizedData.isActive as boolean)
              : true, // Default to true
            sortOrder: normalizedData.sortOrder !== undefined
              ? (normalizedData.sortOrder as number)
              : 0, // Default to 0
          };

          // Call API to create game
          await this.api.createGame(createDto);
          result.importedCount++;

          // Add to existing map to prevent duplicates within same import
          existingKeyMap.set(gameKey, -1);
        }

        // Report progress
        if (options?.onProgress) {
          options.onProgress(i + 1, validRows.length);
        }
      } catch (error) {
        result.success = false;
        result.skippedCount++;

        const errorMessage =
          error instanceof Error ? error.message : 'Unknown error occurred';
        result.errors.push({
          rowNumber: row.rowNumber,
          message: errorMessage,
        });
      }
    }

    // Update success status
    result.success = result.errors.length === 0;

    return result;
  }
}

// Export singleton instance
export const importService = new ImportService();
