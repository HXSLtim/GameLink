import { renderHook, act, waitFor } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { ThemeProvider, useTheme } from './ThemeContext';

const wrapper = ({ children }: { children: React.ReactNode }) => (
  <ThemeProvider>{children}</ThemeProvider>
);

describe('ThemeContext', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    localStorage.clear();
  });

  it('should throw error when useTheme is used outside ThemeProvider', () => {
    // Suppress console.error for this test
    const consoleSpy = vi.spyOn(console, 'error').mockImplementation(() => {});

    expect(() => {
      renderHook(() => useTheme());
    }).toThrow('useTheme must be used within ThemeProvider');

    consoleSpy.mockRestore();
  });

  it('should initialize with system theme mode', () => {
    const { result } = renderHook(() => useTheme(), { wrapper });

    expect(result.current.mode).toBe('system');
    expect(['light', 'dark']).toContain(result.current.effective);
  });

  it('should load theme from localStorage', () => {
    localStorage.setItem('gamelink_theme', 'dark');

    const { result } = renderHook(() => useTheme(), { wrapper });

    expect(result.current.mode).toBe('dark');
    expect(result.current.effective).toBe('dark');
  });

  it('should save theme to localStorage when mode changes', () => {
    const { result } = renderHook(() => useTheme(), { wrapper });

    act(() => {
      result.current.setMode('dark');
    });

    expect(localStorage.getItem('gamelink_theme')).toBe('dark');
    expect(result.current.mode).toBe('dark');
    expect(result.current.effective).toBe('dark');
  });

  it('should change to light mode', () => {
    const { result } = renderHook(() => useTheme(), { wrapper });

    act(() => {
      result.current.setMode('light');
    });

    expect(result.current.mode).toBe('light');
    expect(result.current.effective).toBe('light');
  });

  it('should change to dark mode', () => {
    const { result } = renderHook(() => useTheme(), { wrapper });

    act(() => {
      result.current.setMode('dark');
    });

    expect(result.current.mode).toBe('dark');
    expect(result.current.effective).toBe('dark');
  });

  it('should change back to system mode', () => {
    const { result } = renderHook(() => useTheme(), { wrapper });

    act(() => {
      result.current.setMode('dark');
    });

    expect(result.current.mode).toBe('dark');

    act(() => {
      result.current.setMode('system');
    });

    expect(result.current.mode).toBe('system');
    // Effective should be based on system preference
    expect(['light', 'dark']).toContain(result.current.effective);
  });

  it('should apply dark theme class to document', async () => {
    const { result } = renderHook(() => useTheme(), { wrapper });

    act(() => {
      result.current.setMode('dark');
    });

    // 等待DOM更新
    await waitFor(() => {
      expect(document.documentElement.classList.contains('dark-theme')).toBe(true);
      expect(document.body.classList.contains('dark-theme')).toBe(true);
    });
  });

  it('should remove dark theme class when switching to light', async () => {
    const { result } = renderHook(() => useTheme(), { wrapper });

    act(() => {
      result.current.setMode('dark');
    });

    await waitFor(() => {
      expect(document.documentElement.classList.contains('dark-theme')).toBe(true);
    });

    act(() => {
      result.current.setMode('light');
    });

    // 等待DOM更新
    await waitFor(() => {
      expect(document.documentElement.classList.contains('dark-theme')).toBe(false);
      expect(document.body.classList.contains('dark-theme')).toBe(false);
    });
  });

  it('should persist theme preference across rerenders', () => {
    const { result, rerender } = renderHook(() => useTheme(), { wrapper });

    act(() => {
      result.current.setMode('dark');
    });

    expect(result.current.mode).toBe('dark');

    rerender();

    expect(result.current.mode).toBe('dark');
    expect(localStorage.getItem('gamelink_theme')).toBe('dark');
  });
});
