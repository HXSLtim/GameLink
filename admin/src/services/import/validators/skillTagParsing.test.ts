/**
 * Property-Based Tests for Skill Tag Parsing
 *
 * **Feature: admin-phase3-improvements, Property 21: Skill Tag Parsing**
 * **Validates: Requirements 7.4**
 *
 * Tests that the skill tag parser correctly splits comma-separated values
 * and trims whitespace from each tag.
 */

import { describe, it, expect } from 'vitest';
import * as fc from 'fast-check';
import { parseSkillTags, validateSkillTags } from '../templates/playerTemplate';

describe('Skill Tag Parsing - Property Tests', () => {
  /**
   * **Feature: admin-phase3-improvements, Property 21: Skill Tag Parsing**
   * **Validates: Requirements 7.4**
   *
   * For any skill tags input string, the parser SHALL correctly split
   * comma-separated values and trim whitespace from each tag.
   */
  describe('Property 21: Skill Tag Parsing', () => {
    // Generate valid tag strings (non-empty, no commas, reasonable length)
    const validTagArb = fc
      .string({ minLength: 1, maxLength: 20 })
      .filter((s) => !s.includes(',') && !s.includes('，') && s.trim().length > 0)
      .map((s) => s.trim());

    it('should correctly split tags by English comma', () => {
      fc.assert(
        fc.property(
          fc.array(validTagArb, { minLength: 1, maxLength: 10 }),
          (tags) => {
            const input = tags.join(',');
            const result = parseSkillTags(input);

            // Should have same number of tags
            expect(result.length).toBe(tags.length);

            // Each tag should be present (trimmed)
            for (const tag of tags) {
              expect(result).toContain(tag.trim());
            }
          }
        ),
        { numRuns: 100 }
      );
    });

    it('should correctly split tags by Chinese comma', () => {
      fc.assert(
        fc.property(
          fc.array(validTagArb, { minLength: 1, maxLength: 10 }),
          (tags) => {
            const input = tags.join('，'); // Chinese comma
            const result = parseSkillTags(input);

            // Should have same number of tags
            expect(result.length).toBe(tags.length);

            // Each tag should be present (trimmed)
            for (const tag of tags) {
              expect(result).toContain(tag.trim());
            }
          }
        ),
        { numRuns: 100 }
      );
    });

    it('should correctly split tags with mixed commas', () => {
      fc.assert(
        fc.property(
          fc.array(validTagArb, { minLength: 2, maxLength: 10 }),
          fc.array(fc.constantFrom(',', '，'), { minLength: 1 }),
          (tags, separators) => {
            // Build input with alternating separators
            let input = tags[0];
            for (let i = 1; i < tags.length; i++) {
              const sep = separators[(i - 1) % separators.length];
              input += sep + tags[i];
            }

            const result = parseSkillTags(input);

            // Should have same number of tags
            expect(result.length).toBe(tags.length);

            // Each tag should be present (trimmed)
            for (const tag of tags) {
              expect(result).toContain(tag.trim());
            }
          }
        ),
        { numRuns: 100 }
      );
    });

    it('should trim whitespace from each tag', () => {
      // Generate tags with various whitespace patterns
      const whitespaceArb = fc.constantFrom('', ' ', '  ', '\t', ' \t ');

      fc.assert(
        fc.property(
          fc.array(validTagArb, { minLength: 1, maxLength: 5 }),
          fc.array(whitespaceArb, { minLength: 2 }),
          (tags, whitespaces) => {
            // Add whitespace around each tag
            const paddedTags = tags.map((tag, i) => {
              const before = whitespaces[i % whitespaces.length];
              const after = whitespaces[(i + 1) % whitespaces.length];
              return before + tag + after;
            });

            const input = paddedTags.join(',');
            const result = parseSkillTags(input);

            // Each result should be trimmed
            for (const tag of result) {
              expect(tag).toBe(tag.trim());
              expect(tag.length).toBeGreaterThan(0);
            }

            // Should have same number of non-empty tags
            expect(result.length).toBe(tags.length);
          }
        ),
        { numRuns: 100 }
      );
    });

    it('should filter out empty tags', () => {
      fc.assert(
        fc.property(
          fc.array(validTagArb, { minLength: 1, maxLength: 5 }),
          (tags) => {
            // Create input with some empty segments
            const withEmpty = [...tags, '', '  ', '\t'];
            const input = withEmpty.join(',');
            const result = parseSkillTags(input);

            // Should only contain non-empty tags
            expect(result.length).toBe(tags.length);
            for (const tag of result) {
              expect(tag.length).toBeGreaterThan(0);
            }
          }
        ),
        { numRuns: 100 }
      );
    });

    it('should return empty array for empty or whitespace-only input', () => {
      const emptyInputArb = fc.constantFrom('', ' ', '  ', '\t', '\n', '   \t   ');

      fc.assert(
        fc.property(emptyInputArb, (input) => {
          const result = parseSkillTags(input);
          expect(result).toEqual([]);
        }),
        { numRuns: 20 }
      );
    });

    it('should return empty array for null or undefined input', () => {
      // @ts-expect-error Testing invalid input
      expect(parseSkillTags(null)).toEqual([]);
      // @ts-expect-error Testing invalid input
      expect(parseSkillTags(undefined)).toEqual([]);
    });

    it('should handle single tag without comma', () => {
      fc.assert(
        fc.property(validTagArb, (tag) => {
          const result = parseSkillTags(tag);
          expect(result).toEqual([tag.trim()]);
        }),
        { numRuns: 50 }
      );
    });

    it('should preserve tag order', () => {
      fc.assert(
        fc.property(
          fc.array(validTagArb, { minLength: 2, maxLength: 10 }),
          (tags) => {
            const input = tags.join(',');
            const result = parseSkillTags(input);

            // Order should be preserved
            for (let i = 0; i < tags.length; i++) {
              expect(result[i]).toBe(tags[i].trim());
            }
          }
        ),
        { numRuns: 100 }
      );
    });
  });

  describe('Skill Tag Validation', () => {
    it('should validate tags with valid length', () => {
      const validTagArb = fc
        .string({ minLength: 1, maxLength: 20 })
        .filter((s) => s.trim().length > 0 && s.trim().length <= 20)
        .map((s) => s.trim());

      fc.assert(
        fc.property(
          fc.array(validTagArb, { minLength: 1, maxLength: 10 }),
          (tags) => {
            const result = validateSkillTags(tags);
            expect(result.valid).toBe(true);
            expect(result.invalidTags).toHaveLength(0);
          }
        ),
        { numRuns: 100 }
      );
    });

    it('should reject tags that are too long', () => {
      const longTagArb = fc.string({ minLength: 21, maxLength: 50 });

      fc.assert(
        fc.property(longTagArb, (longTag) => {
          const result = validateSkillTags([longTag]);
          expect(result.valid).toBe(false);
          expect(result.invalidTags).toContain(longTag);
        }),
        { numRuns: 50 }
      );
    });

    it('should reject empty tags', () => {
      const result = validateSkillTags(['']);
      expect(result.valid).toBe(false);
      expect(result.invalidTags).toContain('');
    });
  });
});

describe('Skill Tag Parsing - Unit Tests', () => {
  describe('parseSkillTags', () => {
    it('should parse simple comma-separated tags', () => {
      expect(parseSkillTags('上分,陪玩,教学')).toEqual(['上分', '陪玩', '教学']);
    });

    it('should parse tags with Chinese comma', () => {
      expect(parseSkillTags('上分，陪玩，教学')).toEqual(['上分', '陪玩', '教学']);
    });

    it('should trim whitespace', () => {
      expect(parseSkillTags(' 上分 , 陪玩 , 教学 ')).toEqual(['上分', '陪玩', '教学']);
    });

    it('should filter empty tags', () => {
      expect(parseSkillTags('上分,,陪玩,  ,教学')).toEqual(['上分', '陪玩', '教学']);
    });

    it('should handle single tag', () => {
      expect(parseSkillTags('上分')).toEqual(['上分']);
    });

    it('should return empty array for empty string', () => {
      expect(parseSkillTags('')).toEqual([]);
    });
  });

  describe('validateSkillTags', () => {
    it('should validate valid tags', () => {
      const result = validateSkillTags(['上分', '陪玩', '教学']);
      expect(result.valid).toBe(true);
      expect(result.invalidTags).toHaveLength(0);
    });

    it('should reject empty tags', () => {
      const result = validateSkillTags(['上分', '', '教学']);
      expect(result.valid).toBe(false);
      expect(result.invalidTags).toContain('');
    });

    it('should reject tags longer than 20 characters', () => {
      const longTag = 'a'.repeat(21);
      const result = validateSkillTags(['上分', longTag]);
      expect(result.valid).toBe(false);
      expect(result.invalidTags).toContain(longTag);
    });

    it('should handle empty array', () => {
      const result = validateSkillTags([]);
      expect(result.valid).toBe(true);
      expect(result.invalidTags).toHaveLength(0);
    });
  });
});
