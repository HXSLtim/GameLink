/**
 * Game Import Template
 * Defines the structure for importing game data
 *
 * @module services/import/templates/gameTemplate
 */

import type { ImportTemplate } from './types';

/**
 * Game key validation regex (alphanumeric with underscores/hyphens)
 */
const GAME_KEY_REGEX = /^[a-zA-Z0-9_-]+$/;

/**
 * Game import template definition
 *
 * Columns:
 * - key (required): Unique game identifier
 * - name (required): Game display name
 * - category (optional): Game category name
 * - description (optional): Game description
 * - isActive (optional): Whether the game is active, defaults to true
 */
export const gameImportTemplate: ImportTemplate = {
  type: 'game',
  name: 'Game Import Template',
  nameZh: '游戏导入模板',
  description: 'Template for importing game data including key, name, category, and status',
  columns: [
    {
      key: 'key',
      label: 'Game Key',
      labelZh: '游戏标识',
      required: true,
      type: 'string',
      description: 'Unique game identifier (alphanumeric, underscores, hyphens)',
      example: 'honor_of_kings',
      validation: (value) => typeof value === 'string' && GAME_KEY_REGEX.test(value),
    },
    {
      key: 'name',
      label: 'Game Name',
      labelZh: '游戏名称',
      required: true,
      type: 'string',
      description: 'Game display name',
      example: '王者荣耀',
      validation: (value) => typeof value === 'string' && value.trim().length > 0,
    },
    {
      key: 'category',
      label: 'Category',
      labelZh: '分类',
      required: false,
      type: 'string',
      description: 'Game category name',
      example: 'MOBA',
      validation: (value) =>
        value === undefined || value === '' || (typeof value === 'string' && value.length <= 50),
    },
    {
      key: 'description',
      label: 'Description',
      labelZh: '描述',
      required: false,
      type: 'string',
      description: 'Game description',
      example: '5v5 MOBA手游',
      validation: (value) =>
        value === undefined || value === '' || (typeof value === 'string' && value.length <= 500),
    },
    {
      key: 'isActive',
      label: 'Is Active',
      labelZh: '是否启用',
      required: false,
      type: 'boolean',
      defaultValue: true,
      description: 'Whether the game is active (true/false or 是/否)',
      example: '是',
      validation: (value) => {
        if (value === undefined || value === '') return true;
        if (typeof value === 'boolean') return true;
        const strValue = String(value).toLowerCase();
        return ['true', 'false', '是', '否', '1', '0'].includes(strValue);
      },
    },
  ],
  exampleRows: [
    {
      key: 'honor_of_kings',
      name: '王者荣耀',
      category: 'MOBA',
      description: '5v5 MOBA手游',
      isActive: true,
    },
    {
      key: 'league_of_legends',
      name: '英雄联盟',
      category: 'MOBA',
      description: '经典MOBA游戏',
      isActive: true,
    },
    {
      key: 'pubg_mobile',
      name: '和平精英',
      category: 'FPS',
      description: '战术竞技手游',
      isActive: true,
    },
  ],
};

/**
 * Validate game key format
 */
export function validateGameKey(key: string): boolean {
  return GAME_KEY_REGEX.test(key);
}

/**
 * Parse boolean value from various formats
 */
export function parseBoolean(value: unknown): boolean {
  if (typeof value === 'boolean') {
    return value;
  }

  if (value === undefined || value === null || value === '') {
    return true; // Default to true for isActive
  }

  const strValue = String(value).toLowerCase().trim();

  if (['true', '是', '1', 'yes', 'y'].includes(strValue)) {
    return true;
  }

  if (['false', '否', '0', 'no', 'n'].includes(strValue)) {
    return false;
  }

  return true; // Default to true
}

/**
 * Normalize game key to lowercase with underscores
 */
export function normalizeGameKey(key: string): string {
  return key
    .toLowerCase()
    .trim()
    .replace(/\s+/g, '_')
    .replace(/-+/g, '_')
    .replace(/[^a-z0-9_]/g, '');
}
