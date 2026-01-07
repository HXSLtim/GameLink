/**
 * Import Templates Module
 * Exports all import template definitions
 *
 * @module services/import/templates
 */

// Types
export * from './types';

// Templates
export {
  userImportTemplate,
  validateUserEmail,
  validateUserPhone,
  validateUserRole,
  validateUserStatus,
} from './userTemplate';

export {
  playerImportTemplate,
  parseSkillTags,
  yuanToCents,
  validateSkillTags,
} from './playerTemplate';

export {
  gameImportTemplate,
  validateGameKey,
  parseBoolean,
  normalizeGameKey,
} from './gameTemplate';

import type { ImportTemplate, ImportType } from './types';
import { userImportTemplate } from './userTemplate';
import { playerImportTemplate } from './playerTemplate';
import { gameImportTemplate } from './gameTemplate';

/**
 * Map of import types to their templates
 */
export const importTemplates: Record<ImportType, ImportTemplate> = {
  user: userImportTemplate,
  player: playerImportTemplate,
  game: gameImportTemplate,
};

/**
 * Get template by import type
 */
export function getTemplate(type: ImportType): ImportTemplate {
  return importTemplates[type];
}

/**
 * Get all available import types
 */
export function getAvailableImportTypes(): ImportType[] {
  return Object.keys(importTemplates) as ImportType[];
}

/**
 * Check if an import type is valid
 */
export function isValidImportType(type: string): type is ImportType {
  return type in importTemplates;
}
