/**
 * Player Import Template
 * Defines the structure for importing player (陪玩师) data
 *
 * @module services/import/templates/playerTemplate
 */

import type { ImportTemplate } from './types';

/**
 * Email validation regex (RFC 5322 simplified)
 */
const EMAIL_REGEX = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;

/**
 * Player import template definition
 *
 * Columns:
 * - userEmail (required): Email of the existing user to register as player
 * - nickname (optional): Player's display nickname
 * - bio (optional): Player's biography/introduction
 * - hourlyRate (optional): Hourly rate in yuan (will be converted to cents)
 * - mainGame (optional): Main game name or ID
 * - skillTags (optional): Comma-separated skill tags
 */
export const playerImportTemplate: ImportTemplate = {
  type: 'player',
  name: 'Player Import Template',
  nameZh: '陪玩师导入模板',
  description:
    'Template for importing player data. Users must exist before importing as players.',
  columns: [
    {
      key: 'userEmail',
      label: 'User Email',
      labelZh: '用户邮箱',
      required: true,
      type: 'email',
      description: 'Email of the existing user to register as player',
      example: 'player@example.com',
      validation: (value) => typeof value === 'string' && EMAIL_REGEX.test(value),
    },
    {
      key: 'nickname',
      label: 'Nickname',
      labelZh: '昵称',
      required: false,
      type: 'string',
      description: "Player's display nickname",
      example: '游戏达人',
      validation: (value) =>
        value === undefined || value === '' || (typeof value === 'string' && value.length <= 50),
    },
    {
      key: 'bio',
      label: 'Bio',
      labelZh: '简介',
      required: false,
      type: 'string',
      description: "Player's biography/introduction",
      example: '专业陪玩，技术过硬',
      validation: (value) =>
        value === undefined || value === '' || (typeof value === 'string' && value.length <= 500),
    },
    {
      key: 'hourlyRate',
      label: 'Hourly Rate (Yuan)',
      labelZh: '时薪(元)',
      required: false,
      type: 'number',
      defaultValue: 0,
      description: 'Hourly rate in yuan (will be converted to cents internally)',
      example: '30',
      validation: (value) => {
        if (value === undefined || value === '') return true;
        const num = typeof value === 'number' ? value : parseFloat(String(value));
        return !isNaN(num) && num >= 0;
      },
    },
    {
      key: 'mainGame',
      label: 'Main Game',
      labelZh: '主游戏',
      required: false,
      type: 'string',
      description: 'Main game name or ID',
      example: '王者荣耀',
      validation: (value) =>
        value === undefined || value === '' || typeof value === 'string' || typeof value === 'number',
    },
    {
      key: 'skillTags',
      label: 'Skill Tags',
      labelZh: '技能标签',
      required: false,
      type: 'string',
      description: 'Comma-separated skill tags',
      example: '上分,陪玩,教学',
      validation: (value) => value === undefined || value === '' || typeof value === 'string',
    },
  ],
  exampleRows: [
    {
      userEmail: 'player1@example.com',
      nickname: '游戏达人',
      bio: '专业陪玩，技术过硬',
      hourlyRate: 30,
      mainGame: '王者荣耀',
      skillTags: '上分,陪玩,教学',
    },
    {
      userEmail: 'player2@example.com',
      nickname: '电竞小王子',
      bio: '多年游戏经验',
      hourlyRate: 50,
      mainGame: '英雄联盟',
      skillTags: '上分,开黑',
    },
    {
      userEmail: 'player3@example.com',
      nickname: '休闲玩家',
      bio: '',
      hourlyRate: 20,
      mainGame: '和平精英',
      skillTags: '陪玩',
    },
  ],
};

/**
 * Parse skill tags from comma-separated string
 * Handles both Chinese and English commas
 */
export function parseSkillTags(tagsString: string): string[] {
  if (!tagsString || typeof tagsString !== 'string') {
    return [];
  }

  // Split by comma (both Chinese and English) and trim whitespace
  return tagsString
    .split(/[,，]/)
    .map((tag) => tag.trim())
    .filter((tag) => tag.length > 0);
}

/**
 * Convert hourly rate from yuan to cents
 */
export function yuanToCents(yuan: number): number {
  return Math.round(yuan * 100);
}

/**
 * Validate skill tags format
 */
export function validateSkillTags(tags: string[]): { valid: boolean; invalidTags: string[] } {
  const invalidTags: string[] = [];

  for (const tag of tags) {
    // Tags should be non-empty and not too long
    if (tag.length === 0 || tag.length > 20) {
      invalidTags.push(tag);
    }
  }

  return {
    valid: invalidTags.length === 0,
    invalidTags,
  };
}
