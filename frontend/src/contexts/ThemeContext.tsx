import { createContext, useContext, useState, useEffect, useCallback, useMemo } from 'react';
import { storage } from '../utils/storage';

export type ThemeMode = 'light' | 'dark' | 'system';
export type EffectiveTheme = 'light' | 'dark';

const THEME_STORAGE_KEY = 'gamelink_theme';
const DARK_THEME_CLASS = 'arco-theme-dark';

interface ThemeContextValue {
  mode: ThemeMode;
  effective: EffectiveTheme;
  setMode: (mode: ThemeMode) => void;
}

interface ThemeProviderProps {
  children: React.ReactNode;
}

const ThemeContext = createContext<ThemeContextValue | null>(null);

/**
 * Get the system's preferred color scheme
 */
const getSystemColorScheme = (): EffectiveTheme => {
  if (typeof window === 'undefined') return 'light';
  return window.matchMedia?.('(prefers-color-scheme: dark)').matches ? 'dark' : 'light';
};

/**
 * Apply theme class to HTML elements
 */
const applyThemeClass = (theme: EffectiveTheme): void => {
  if (typeof document === 'undefined') return;

  const html = document.documentElement;
  const body = document.body;

  if (theme === 'dark') {
    html.classList.add(DARK_THEME_CLASS);
    body.classList.add(DARK_THEME_CLASS);
    body.setAttribute('arco-theme', 'dark');
  } else {
    html.classList.remove(DARK_THEME_CLASS);
    body.classList.remove(DARK_THEME_CLASS);
    body.removeAttribute('arco-theme');
  }
};

/**
 * Theme provider component
 * 
 * @component
 * @description Manages application theme (light/dark/system) and applies appropriate CSS classes
 */
export const ThemeProvider: React.FC<ThemeProviderProps> = ({ children }) => {
  const [mode, setModeState] = useState<ThemeMode>(() => {
    const savedTheme = storage.getItem<ThemeMode>(THEME_STORAGE_KEY);
    return savedTheme || 'system';
  });

  const [effective, setEffective] = useState<EffectiveTheme>(() => {
    return mode === 'system' ? getSystemColorScheme() : mode;
  });

  // Apply theme when mode changes
  useEffect(() => {
    const systemTheme = getSystemColorScheme();
    const effectiveTheme = mode === 'system' ? systemTheme : mode;
    setEffective(effectiveTheme);
    applyThemeClass(effectiveTheme);
  }, [mode]);

  // Listen to system theme changes
  useEffect(() => {
    if (!window.matchMedia) return;

    const mediaQuery = window.matchMedia('(prefers-color-scheme: dark)');
    const handleSystemThemeChange = (event: MediaQueryListEvent) => {
      if (mode === 'system') {
        const newTheme: EffectiveTheme = event.matches ? 'dark' : 'light';
        setEffective(newTheme);
        applyThemeClass(newTheme);
      }
    };

    // Modern browsers
    mediaQuery.addEventListener('change', handleSystemThemeChange);

    return () => {
      mediaQuery.removeEventListener('change', handleSystemThemeChange);
    };
  }, [mode]);

  const handleSetMode = useCallback((newMode: ThemeMode) => {
    setModeState(newMode);
    storage.setItem(THEME_STORAGE_KEY, newMode);
  }, []);

  const contextValue = useMemo<ThemeContextValue>(
    () => ({
      mode,
      effective,
      setMode: handleSetMode,
    }),
    [mode, effective, handleSetMode],
  );

  return <ThemeContext.Provider value={contextValue}>{children}</ThemeContext.Provider>;
};

/**
 * Hook to access theme context
 * 
 * @throws Error if used outside ThemeProvider
 */
export const useTheme = (): ThemeContextValue => {
  const context = useContext(ThemeContext);
  if (!context) {
    throw new Error('useTheme must be used within ThemeProvider');
  }
  return context;
};
