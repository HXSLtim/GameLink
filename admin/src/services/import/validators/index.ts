/**
 * Import Validators Module
 * Exports all validation utilities for import operations
 *
 * @module services/import/validators
 */

export {
  validateStructure,
  createColumnMapping,
  mapRowToTemplate,
  getMissingRequiredColumns,
  hasAllRequiredColumns,
  type StructureValidationResult,
  type ColumnMapping,
} from './structureValidator';

export {
  validateUserRow,
  validateUserData,
  findDuplicates,
  normalizeUserData,
  type FieldError,
  type RowValidationResult,
  type UserDataValidationResult,
  type UserDataValidationOptions,
} from './userDataValidator';

export {
  validatePlayerRow,
  validatePlayerData,
  findDuplicateUserEmails,
  normalizePlayerData,
  type PlayerDataValidationResult,
  type PlayerDataValidationOptions,
} from './playerDataValidator';

export {
  validateGameRow,
  validateGameData,
  findDuplicateGameKeys,
  normalizeGameData,
  type GameDataValidationResult,
  type GameDataValidationOptions,
} from './gameDataValidator';
