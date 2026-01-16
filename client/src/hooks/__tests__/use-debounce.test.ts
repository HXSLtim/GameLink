import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { renderHook, act } from '@testing-library/react';
import { useDebounce } from '../use-debounce';

describe('useDebounce', () => {
    beforeEach(() => {
        vi.useFakeTimers();
    });

    afterEach(() => {
        vi.useRealTimers();
    });

    it('should return initial value immediately', () => {
        const { result } = renderHook(() => useDebounce('initial', 300));
        expect(result.current).toBe('initial');
    });

    it('should debounce value changes', () => {
        const { result, rerender } = renderHook(
            ({ value, delay }) => useDebounce(value, delay),
            { initialProps: { value: 'initial', delay: 300 } }
        );

        expect(result.current).toBe('initial');

        // Change value
        rerender({ value: 'updated', delay: 300 });

        // Value should not change immediately
        expect(result.current).toBe('initial');

        // Fast forward time
        act(() => {
            vi.advanceTimersByTime(300);
        });

        // Now value should be updated
        expect(result.current).toBe('updated');
    });

    it('should reset timer on rapid value changes', () => {
        const { result, rerender } = renderHook(
            ({ value, delay }) => useDebounce(value, delay),
            { initialProps: { value: 'initial', delay: 300 } }
        );

        // Change value multiple times rapidly
        rerender({ value: 'change1', delay: 300 });
        act(() => {
            vi.advanceTimersByTime(100);
        });

        rerender({ value: 'change2', delay: 300 });
        act(() => {
            vi.advanceTimersByTime(100);
        });

        rerender({ value: 'change3', delay: 300 });

        // Value should still be initial
        expect(result.current).toBe('initial');

        // Fast forward past debounce delay
        act(() => {
            vi.advanceTimersByTime(300);
        });

        // Should have the last value
        expect(result.current).toBe('change3');
    });

    it('should use default delay of 300ms', () => {
        const { result, rerender } = renderHook(
            ({ value }) => useDebounce(value),
            { initialProps: { value: 'initial' } }
        );

        rerender({ value: 'updated' });

        // Should not update before 300ms
        act(() => {
            vi.advanceTimersByTime(299);
        });
        expect(result.current).toBe('initial');

        // Should update after 300ms
        act(() => {
            vi.advanceTimersByTime(1);
        });
        expect(result.current).toBe('updated');
    });

    it('should work with different types', () => {
        // Number
        const { result: numberResult, rerender: rerenderNumber } = renderHook(
            ({ value }) => useDebounce(value, 100),
            { initialProps: { value: 0 } }
        );

        rerenderNumber({ value: 42 });
        act(() => {
            vi.advanceTimersByTime(100);
        });
        expect(numberResult.current).toBe(42);

        // Object
        const { result: objectResult, rerender: rerenderObject } = renderHook(
            ({ value }) => useDebounce(value, 100),
            { initialProps: { value: { name: 'initial' } } }
        );

        rerenderObject({ value: { name: 'updated' } });
        act(() => {
            vi.advanceTimersByTime(100);
        });
        expect(objectResult.current).toEqual({ name: 'updated' });
    });

    it('should handle delay changes', () => {
        const { result, rerender } = renderHook(
            ({ value, delay }) => useDebounce(value, delay),
            { initialProps: { value: 'initial', delay: 300 } }
        );

        // Change both value and delay
        rerender({ value: 'updated', delay: 500 });

        // Should not update at 300ms
        act(() => {
            vi.advanceTimersByTime(300);
        });
        expect(result.current).toBe('initial');

        // Should update at 500ms
        act(() => {
            vi.advanceTimersByTime(200);
        });
        expect(result.current).toBe('updated');
    });
});
