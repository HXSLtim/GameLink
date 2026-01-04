/**
 * Review Format Utility Tests
 *
 * Coverage Target: 85%+
 *
 * Test Scenarios:
 * 1. Rating formatting (stars, numbers)
 * 2. Date/time formatting (datetime, date, relative time)
 * 3. Review status text and color getters
 * 4. Report status text and color getters
 * 5. Sensitive word category and severity getters
 * 6. Sort by text getter
 * 7. Text truncation
 * 8. Image count formatting
 * 9. Percentage formatting
 */
import { describe, it, expect, vi } from 'vitest';
import * as fc from 'fast-check';
import {
    formatRatingStars,
    formatRatingNumber,
    formatDateTime,
    formatDate,
    formatRelativeTime,
    getReviewStatusText,
    getReviewStatusColor,
    getReportStatusText,
    getReportStatusColor,
    getSensitiveWordCategoryText,
    getSensitiveWordCategoryColor,
    getSensitiveWordSeverityText,
    getSensitiveWordSeverityColor,
    getSortByText,
    truncateText,
    formatImageCount,
    formatPercentage,
} from './reviewFormat';
import type { ReviewStatus, ReviewReportStatus, SensitiveWordCategory, SensitiveWordSeverity, ReviewSortBy } from '@/types/review';

describe('review format utility', () => {
    describe('formatRatingStars', () => {
        it('should format rating to stars', () => {
            expect(formatRatingStars(5)).toBe('★★★★★');
            expect(formatRatingStars(4)).toBe('★★★★☆');
            expect(formatRatingStars(3)).toBe('★★★☆☆');
            expect(formatRatingStars(2)).toBe('★★☆☆☆');
            expect(formatRatingStars(1)).toBe('★☆☆☆☆');
        });

        it('should handle decimal ratings', () => {
            expect(formatRatingStars(4.7)).toBe('★★★★☆');
            expect(formatRatingStars(3.2)).toBe('★★★☆☆');
        });

        it('should handle edge cases', () => {
            expect(formatRatingStars(0)).toBe('☆☆☆☆☆');
            expect(formatRatingStars(6)).toBe('★★★★★'); // Cap at 5
            expect(formatRatingStars(-1)).toBe('☆☆☆☆☆');
        });

        /**
         * Property: Star count should be 5 characters
         */
        it('should always return 5 characters', () => {
            fc.assert(
                fc.property(fc.float({ min: 0, max: 10 }), (rating) => {
                    const result = formatRatingStars(rating);
                    return result.length === 5;
                }),
                { numRuns: 50 }
            );
        });

        /**
         * Property: Only star and empty star characters should be used
         */
        it('should only contain star characters', () => {
            fc.assert(
                fc.property(fc.float({ min: 0, max: 10 }), (rating) => {
                    const result = formatRatingStars(rating);
                    return [...result].every(char => char === '★' || char === '☆');
                }),
                { numRuns: 50 }
            );
        });
    });

    describe('formatRatingNumber', () => {
        it('should format rating to number string', () => {
            expect(formatRatingNumber(4.5)).toBe('4.5分');
            expect(formatRatingNumber(3)).toBe('3.0分');
            expect(formatRatingNumber(5)).toBe('5.0分');
        });

        it('should use custom precision', () => {
            expect(formatRatingNumber(4.567, 2)).toBe('4.57分');
            expect(formatRatingNumber(4.567, 3)).toBe('4.567分');
            expect(formatRatingNumber(4, 0)).toBe('4分');
        });

        it('should handle edge cases', () => {
            expect(formatRatingNumber(0)).toBe('0.0分');
            expect(formatRatingNumber(5)).toBe('5.0分');
        });
    });

    describe('formatDateTime', () => {
        it('should format date time string', () => {
            expect(formatDateTime('2024-01-15T10:30:00Z')).toBe('2024-01-15 10:30:00');
            expect(formatDateTime('2024-12-31T23:59:59Z')).toBe('2024-12-31 23:59:59');
        });

        it('should use custom format', () => {
            expect(formatDateTime('2024-01-15T10:30:00Z', 'YYYY-MM-DD')).toBe('2024-01-15');
            expect(formatDateTime('2024-01-15T10:30:00Z', 'HH:mm:ss')).toBe('10:30:00');
        });

        it('should handle null/undefined', () => {
            expect(formatDateTime(null)).toBe('-');
            expect(formatDateTime(undefined)).toBe('-');
            expect(formatDateTime('')).toBe('-');
        });

        it('should handle invalid dates gracefully', () => {
            expect(formatDateTime('invalid-date')).toMatch(/Invalid/);
        });
    });

    describe('formatDate', () => {
        it('should format date string', () => {
            expect(formatDate('2024-01-15T10:30:00Z')).toBe('2024-01-15');
            expect(formatDate('2024-12-31T23:59:59Z')).toBe('2024-12-31');
        });

        it('should handle null/undefined', () => {
            expect(formatDate(null)).toBe('-');
            expect(formatDate(undefined)).toBe('-');
        });
    });

    describe('formatRelativeTime', () => {
        beforeEach(() => {
            vi.useFakeTimers();
            vi.setSystemTime(new Date('2024-01-15T12:00:00Z'));
        });

        afterEach(() => {
            vi.useRealTimers();
        });

        it('should format very recent time', () => {
            expect(formatRelativeTime('2024-01-15T11:59:30Z')).toBe('刚刚');
        });

        it('should format minutes ago', () => {
            expect(formatRelativeTime('2024-01-15T11:30:00Z')).toBe('30分钟前');
            expect(formatRelativeTime('2024-01-15T11:00:00Z')).toBe('1小时前');
        });

        it('should format hours ago', () => {
            expect(formatRelativeTime('2024-01-15T08:00:00Z')).toBe('4小时前');
            expect(formatRelativeTime('2024-01-15T00:00:00Z')).toBe('12小时前');
        });

        it('should format days ago', () => {
            expect(formatRelativeTime('2024-01-14T12:00:00Z')).toBe('1天前');
            expect(formatRelativeTime('2024-01-10T12:00:00Z')).toBe('5天前');
        });

        it('should format old dates as absolute', () => {
            expect(formatRelativeTime('2024-01-01T12:00:00Z')).toBe('2024-01-01');
            expect(formatRelativeTime('2023-12-01T12:00:00Z')).toBe('2023-12-01');
        });

        it('should handle null/undefined', () => {
            expect(formatRelativeTime(null)).toBe('-');
            expect(formatRelativeTime(undefined)).toBe('-');
        });
    });

    describe('getReviewStatusText', () => {
        it('should return correct text for status', () => {
            const statuses: ReviewStatus[] = ['pending', 'approved', 'rejected', 'flagged'];
            statuses.forEach(status => {
                const text = getReviewStatusText(status);
                expect(text).toBeTruthy();
                expect(typeof text).toBe('string');
            });
        });

        it('should return original status for unknown values', () => {
            expect(getReviewStatusText('unknown' as ReviewStatus)).toBe('unknown');
        });
    });

    describe('getReviewStatusColor', () => {
        it('should return color for status', () => {
            const statuses: ReviewStatus[] = ['pending', 'approved', 'rejected', 'flagged'];
            statuses.forEach(status => {
                const color = getReviewStatusColor(status);
                expect(color).toBeTruthy();
                expect(typeof color).toBe('string');
            });
        });

        it('should return default color for unknown values', () => {
            expect(getReviewStatusColor('unknown' as ReviewStatus)).toBe('default');
        });
    });

    describe('getReportStatusText', () => {
        it('should return correct text for status', () => {
            const statuses: ReviewReportStatus[] = ['pending', 'investigating', 'resolved', 'dismissed'];
            statuses.forEach(status => {
                const text = getReportStatusText(status);
                expect(text).toBeTruthy();
                expect(typeof text).toBe('string');
            });
        });

        it('should return original status for unknown values', () => {
            expect(getReportStatusText('unknown' as ReviewReportStatus)).toBe('unknown');
        });
    });

    describe('getReportStatusColor', () => {
        it('should return color for status', () => {
            const statuses: ReviewReportStatus[] = ['pending', 'investigating', 'resolved', 'dismissed'];
            statuses.forEach(status => {
                const color = getReportStatusColor(status);
                expect(color).toBeTruthy();
                expect(typeof color).toBe('string');
            });
        });

        it('should return default color for unknown values', () => {
            expect(getReportStatusColor('unknown' as ReviewReportStatus)).toBe('default');
        });
    });

    describe('getSensitiveWordCategoryText', () => {
        it('should return correct text for category', () => {
            const categories: SensitiveWordCategory[] = ['politics', 'violence', 'porn', 'spam'];
            categories.forEach(category => {
                const text = getSensitiveWordCategoryText(category);
                expect(text).toBeTruthy();
                expect(typeof text).toBe('string');
            });
        });

        it('should return original category for unknown values', () => {
            expect(getSensitiveWordCategoryText('unknown' as SensitiveWordCategory)).toBe('unknown');
        });
    });

    describe('getSensitiveWordCategoryColor', () => {
        it('should return color for category', () => {
            const categories: SensitiveWordCategory[] = ['politics', 'violence', 'porn', 'spam'];
            categories.forEach(category => {
                const color = getSensitiveWordCategoryColor(category);
                expect(color).toBeTruthy();
                expect(typeof color).toBe('string');
            });
        });

        it('should return default color for unknown values', () => {
            expect(getSensitiveWordCategoryColor('unknown' as SensitiveWordCategory)).toBe('default');
        });
    });

    describe('getSensitiveWordSeverityText', () => {
        it('should return correct text for severity', () => {
            const severities: SensitiveWordSeverity[] = ['low', 'medium', 'high', 'critical'];
            severities.forEach(severity => {
                const text = getSensitiveWordSeverityText(severity);
                expect(text).toBeTruthy();
                expect(typeof text).toBe('string');
            });
        });

        it('should return original severity for unknown values', () => {
            expect(getSensitiveWordSeverityText('unknown' as SensitiveWordSeverity)).toBe('unknown');
        });
    });

    describe('getSensitiveWordSeverityColor', () => {
        it('should return color for severity', () => {
            const severities: SensitiveWordSeverity[] = ['low', 'medium', 'high', 'critical'];
            severities.forEach(severity => {
                const color = getSensitiveWordSeverityColor(severity);
                expect(color).toBeTruthy();
                expect(typeof color).toBe('string');
            });
        });

        it('should return default color for unknown values', () => {
            expect(getSensitiveWordSeverityColor('unknown' as SensitiveWordSeverity)).toBe('default');
        });
    });

    describe('getSortByText', () => {
        it('should return correct text for sort option', () => {
            const sortOptions: ReviewSortBy[] = ['newest', 'oldest', 'rating_high', 'rating_low', 'most_helpful'];
            sortOptions.forEach(sortBy => {
                const text = getSortByText(sortBy);
                expect(text).toBeTruthy();
                expect(typeof text).toBe('string');
            });
        });

        it('should return original sortBy for unknown values', () => {
            expect(getSortByText('unknown' as ReviewSortBy)).toBe('unknown');
        });
    });

    describe('truncateText', () => {
        it('should not truncate short text', () => {
            expect(truncateText('short')).toBe('short');
            expect(truncateText('exactly 50!', 50)).toBe('exactly 50!');
        });

        it('should truncate long text', () => {
            expect(truncateText('this is a very long text that should be truncated', 20)).toBe('this is a very lo...');
            expect(truncateText('1234567890123456789012345678901', 30)).toBe('123456789012345678901234567...');
        });

        it('should use default maxLength', () => {
            const longText = 'x'.repeat(100);
            const result = truncateText(longText);
            expect(result.length).toBe(53); // 50 + '...'
            expect(result).toMatch(/\.\.\.$/);
        });

        it('should handle null/undefined', () => {
            expect(truncateText(null)).toBe('-');
            expect(truncateText(undefined)).toBe('-');
        });

        it('should handle empty string', () => {
            expect(truncateText('')).toBe('-');
        });

        /**
         * Property: Truncated text should not exceed maxLength + 3
         */
        it('should respect max length', () => {
            fc.assert(
                fc.property(fc.string(), fc.integer({ min: 5, max: 100 }), (text, maxLength) => {
                    const result = truncateText(text, maxLength);
                    return result.length <= maxLength + 3 || result === '-';
                }),
                { numRuns: 50 }
            );
        });

        /**
         * Property: Truncated text should end with ...
         */
        it('should add ellipsis when truncated', () => {
            fc.assert(
                fc.property(
                    fc.string({ minLength: 51, maxLength: 200 }),
                    (text) => {
                        const result = truncateText(text, 50);
                        if (text.length > 50) {
                            return result.endsWith('...');
                        }
                        return result === text;
                    }
                ),
                { numRuns: 20 }
            );
        });
    });

    describe('formatImageCount', () => {
        it('should format zero images', () => {
            expect(formatImageCount(0)).toBe('无图片');
        });

        it('should format positive count', () => {
            expect(formatImageCount(1)).toBe('1张图片');
            expect(formatImageCount(5)).toBe('5张图片');
            expect(formatImageCount(100)).toBe('100张图片');
        });

        it('should handle edge cases', () => {
            expect(formatImageCount(-1)).toBe('-1张图片');
        });
    });

    describe('formatPercentage', () => {
        it('should format percentage', () => {
            expect(formatPercentage(50)).toBe('50.0%');
            expect(formatPercentage(75.5)).toBe('75.5%');
            expect(formatPercentage(100)).toBe('100.0%');
        });

        it('should use custom precision', () => {
            expect(formatPercentage(50.123, 0)).toBe('50%');
            expect(formatPercentage(50.123, 2)).toBe('50.12%');
            expect(formatPercentage(50.123, 3)).toBe('50.123%');
        });

        it('should handle edge cases', () => {
            expect(formatPercentage(0)).toBe('0.0%');
            expect(formatPercentage(0.5)).toBe('0.5%');
            expect(formatPercentage(99.99)).toBe('99.99%');
        });

        /**
         * Property: Formatted percentage should end with %
         */
        it('should always end with percent sign', () => {
            fc.assert(
                fc.property(fc.float({ min: 0, max: 100 }), (value) => {
                    const result = formatPercentage(value);
                    return result.endsWith('%');
                }),
                { numRuns: 50 }
            );
        });
    });

    describe('property-based tests', () => {
        /**
         * Property: Rating stars should have correct number of filled stars
         */
        it('should have correct filled star count', () => {
            fc.assert(
                fc.property(fc.integer({ min: 0, max: 5 }), (rating) => {
                    const stars = formatRatingStars(rating);
                    const filledCount = (stars.match(/★/g) || []).length;
                    return filledCount === Math.min(Math.max(rating, 0), 5);
                }),
                { numRuns: 20 }
            );
        });

        /**
         * Property: Percentage should be within 0-100 range
         */
        it('should handle percentage values', () => {
            fc.assert(
                fc.property(fc.float({ min: 0, max: 100 }), fc.integer({ min: 0, max: 5 }),
                    (value, precision) => {
                        const result = formatPercentage(value, precision);
                        const numericValue = parseFloat(result.replace('%', ''));
                        return Math.abs(numericValue - value) < Math.pow(10, -precision);
                    }
                ),
                { numRuns: 30 }
            );
        });
    });

    describe('edge cases', () => {
        it('should handle unicode in truncateText', () => {
            expect(truncateText('你好世界你好世界', 5)).toBe('你好世界好...');
        });

        it('should handle emojis in truncateText', () => {
            expect(truncateText('😀😁😂🤣😃😄😅😆😉😊', 5)).toMatch(/😀😁😂🤣😃/);
        });

        it('should handle special characters', () => {
            expect(formatDateTime('2024-01-15T10:30:00Z')).not.toContain('Invalid');
        });

        it('should handle very long strings', () => {
            const longString = 'x'.repeat(10000);
            expect(() => truncateText(longString)).not.toThrow();
        });
    });

    describe('integration tests', () => {
        it('should format complete review data', () => {
            const rating = 4.5;
            const stars = formatRatingStars(rating);
            const number = formatRatingNumber(rating);

            expect(stars).toBe('★★★★☆');
            expect(number).toBe('4.5分');
        });

        it('should format complete report data', () => {
            const status: ReviewReportStatus = 'pending';
            const text = getReportStatusText(status);
            const color = getReportStatusColor(status);

            expect(text).toBeTruthy();
            expect(color).toBeTruthy();
        });

        it('should format sensitive word data', () => {
            const category: SensitiveWordCategory = 'politics';
            const severity: SensitiveWordSeverity = 'high';

            const categoryText = getSensitiveWordCategoryText(category);
            const categoryColor = getSensitiveWordCategoryColor(category);
            const severityText = getSensitiveWordSeverityText(severity);
            const severityColor = getSensitiveWordSeverityColor(severity);

            expect(categoryText).toBeTruthy();
            expect(categoryColor).toBeTruthy();
            expect(severityText).toBeTruthy();
            expect(severityColor).toBeTruthy();
        });
    });
});
