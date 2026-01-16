import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { renderHook, act } from '@testing-library/react';
import { useOnlineStatus } from '../use-online-status';

describe('useOnlineStatus', () => {
    const originalNavigator = window.navigator;

    beforeEach(() => {
        // Reset navigator.onLine to true
        Object.defineProperty(window.navigator, 'onLine', {
            value: true,
            writable: true,
            configurable: true,
        });
    });

    afterEach(() => {
        // Restore original navigator
        Object.defineProperty(window, 'navigator', {
            value: originalNavigator,
            writable: true,
        });
    });

    it('should return true when online', () => {
        Object.defineProperty(window.navigator, 'onLine', { value: true, configurable: true });

        const { result } = renderHook(() => useOnlineStatus());
        expect(result.current).toBe(true);
    });

    it('should return false when offline', () => {
        Object.defineProperty(window.navigator, 'onLine', { value: false, configurable: true });

        const { result } = renderHook(() => useOnlineStatus());
        expect(result.current).toBe(false);
    });

    it('should update when going offline', () => {
        const { result } = renderHook(() => useOnlineStatus());
        expect(result.current).toBe(true);

        act(() => {
            window.dispatchEvent(new Event('offline'));
        });

        expect(result.current).toBe(false);
    });

    it('should update when going online', () => {
        Object.defineProperty(window.navigator, 'onLine', { value: false, configurable: true });

        const { result } = renderHook(() => useOnlineStatus());
        expect(result.current).toBe(false);

        act(() => {
            window.dispatchEvent(new Event('online'));
        });

        expect(result.current).toBe(true);
    });

    it('should add and remove event listeners', () => {
        const addEventListenerSpy = vi.spyOn(window, 'addEventListener');
        const removeEventListenerSpy = vi.spyOn(window, 'removeEventListener');

        const { unmount } = renderHook(() => useOnlineStatus());

        expect(addEventListenerSpy).toHaveBeenCalledWith('online', expect.any(Function));
        expect(addEventListenerSpy).toHaveBeenCalledWith('offline', expect.any(Function));

        unmount();

        expect(removeEventListenerSpy).toHaveBeenCalledWith('online', expect.any(Function));
        expect(removeEventListenerSpy).toHaveBeenCalledWith('offline', expect.any(Function));

        addEventListenerSpy.mockRestore();
        removeEventListenerSpy.mockRestore();
    });

    it('should handle multiple online/offline transitions', () => {
        const { result } = renderHook(() => useOnlineStatus());

        expect(result.current).toBe(true);

        act(() => {
            window.dispatchEvent(new Event('offline'));
        });
        expect(result.current).toBe(false);

        act(() => {
            window.dispatchEvent(new Event('online'));
        });
        expect(result.current).toBe(true);

        act(() => {
            window.dispatchEvent(new Event('offline'));
        });
        expect(result.current).toBe(false);

        act(() => {
            window.dispatchEvent(new Event('online'));
        });
        expect(result.current).toBe(true);
    });
});
