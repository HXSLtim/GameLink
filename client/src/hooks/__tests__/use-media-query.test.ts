import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { renderHook, act } from '@testing-library/react';
import { useMediaQuery, useIsMobile, useIsTablet, useIsDesktop } from '../use-media-query';

describe('useMediaQuery', () => {
    let mockMatchMedia: ReturnType<typeof vi.fn>;
    let mockAddEventListener: ReturnType<typeof vi.fn>;
    let mockRemoveEventListener: ReturnType<typeof vi.fn>;
    let changeHandler: ((event: MediaQueryListEvent) => void) | null = null;

    beforeEach(() => {
        mockAddEventListener = vi.fn((event, handler) => {
            if (event === 'change') {
                changeHandler = handler;
            }
        });
        mockRemoveEventListener = vi.fn();

        mockMatchMedia = vi.fn((query: string) => ({
            matches: false,
            media: query,
            onchange: null,
            addEventListener: mockAddEventListener,
            removeEventListener: mockRemoveEventListener,
            addListener: vi.fn(),
            removeListener: vi.fn(),
            dispatchEvent: vi.fn(),
        }));

        Object.defineProperty(window, 'matchMedia', {
            writable: true,
            value: mockMatchMedia,
        });
    });

    afterEach(() => {
        changeHandler = null;
        vi.clearAllMocks();
    });

    it('should return false when query does not match', () => {
        mockMatchMedia.mockReturnValue({
            matches: false,
            addEventListener: mockAddEventListener,
            removeEventListener: mockRemoveEventListener,
        });

        const { result } = renderHook(() => useMediaQuery('(min-width: 768px)'));
        expect(result.current).toBe(false);
    });

    it('should return true when query matches', () => {
        mockMatchMedia.mockReturnValue({
            matches: true,
            addEventListener: mockAddEventListener,
            removeEventListener: mockRemoveEventListener,
        });

        const { result } = renderHook(() => useMediaQuery('(min-width: 768px)'));
        expect(result.current).toBe(true);
    });

    it('should update when media query changes', () => {
        mockMatchMedia.mockReturnValue({
            matches: false,
            addEventListener: mockAddEventListener,
            removeEventListener: mockRemoveEventListener,
        });

        const { result } = renderHook(() => useMediaQuery('(min-width: 768px)'));
        expect(result.current).toBe(false);

        // Simulate media query change
        act(() => {
            if (changeHandler) {
                changeHandler({ matches: true } as MediaQueryListEvent);
            }
        });

        expect(result.current).toBe(true);
    });

    it('should add and remove event listeners', () => {
        const { unmount } = renderHook(() => useMediaQuery('(min-width: 768px)'));

        expect(mockAddEventListener).toHaveBeenCalledWith('change', expect.any(Function));

        unmount();

        expect(mockRemoveEventListener).toHaveBeenCalledWith('change', expect.any(Function));
    });

    it('should update listener when query changes', () => {
        const { rerender } = renderHook(
            ({ query }) => useMediaQuery(query),
            { initialProps: { query: '(min-width: 768px)' } }
        );

        expect(mockMatchMedia).toHaveBeenCalledWith('(min-width: 768px)');

        rerender({ query: '(min-width: 1024px)' });

        expect(mockMatchMedia).toHaveBeenCalledWith('(min-width: 1024px)');
    });
});

describe('useIsMobile', () => {
    let mockMatchMedia: ReturnType<typeof vi.fn>;

    beforeEach(() => {
        mockMatchMedia = vi.fn(() => ({
            matches: false,
            addEventListener: vi.fn(),
            removeEventListener: vi.fn(),
        }));

        Object.defineProperty(window, 'matchMedia', {
            writable: true,
            value: mockMatchMedia,
        });
    });

    it('should use correct mobile breakpoint query', () => {
        renderHook(() => useIsMobile());
        expect(mockMatchMedia).toHaveBeenCalledWith('(max-width: 639px)');
    });

    it('should return true on mobile viewport', () => {
        mockMatchMedia.mockReturnValue({
            matches: true,
            addEventListener: vi.fn(),
            removeEventListener: vi.fn(),
        });

        const { result } = renderHook(() => useIsMobile());
        expect(result.current).toBe(true);
    });
});

describe('useIsTablet', () => {
    let mockMatchMedia: ReturnType<typeof vi.fn>;

    beforeEach(() => {
        mockMatchMedia = vi.fn(() => ({
            matches: false,
            addEventListener: vi.fn(),
            removeEventListener: vi.fn(),
        }));

        Object.defineProperty(window, 'matchMedia', {
            writable: true,
            value: mockMatchMedia,
        });
    });

    it('should use correct tablet breakpoint query', () => {
        renderHook(() => useIsTablet());
        expect(mockMatchMedia).toHaveBeenCalledWith('(min-width: 640px) and (max-width: 1023px)');
    });

    it('should return true on tablet viewport', () => {
        mockMatchMedia.mockReturnValue({
            matches: true,
            addEventListener: vi.fn(),
            removeEventListener: vi.fn(),
        });

        const { result } = renderHook(() => useIsTablet());
        expect(result.current).toBe(true);
    });
});

describe('useIsDesktop', () => {
    let mockMatchMedia: ReturnType<typeof vi.fn>;

    beforeEach(() => {
        mockMatchMedia = vi.fn(() => ({
            matches: false,
            addEventListener: vi.fn(),
            removeEventListener: vi.fn(),
        }));

        Object.defineProperty(window, 'matchMedia', {
            writable: true,
            value: mockMatchMedia,
        });
    });

    it('should use correct desktop breakpoint query', () => {
        renderHook(() => useIsDesktop());
        expect(mockMatchMedia).toHaveBeenCalledWith('(min-width: 1024px)');
    });

    it('should return true on desktop viewport', () => {
        mockMatchMedia.mockReturnValue({
            matches: true,
            addEventListener: vi.fn(),
            removeEventListener: vi.fn(),
        });

        const { result } = renderHook(() => useIsDesktop());
        expect(result.current).toBe(true);
    });
});
