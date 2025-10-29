import {
  createContext,
  useContext,
  useState,
  useEffect,
  useCallback,
  useMemo,
  useRef,
} from 'react';
import { storage } from '../utils/storage';

export type ThemeMode = 'light' | 'dark' | 'system';
export type EffectiveTheme = 'light' | 'dark';

const THEME_STORAGE_KEY = 'gamelink_theme';
const DARK_THEME_CLASS = 'dark-theme';

interface ThemeContextValue {
  mode: ThemeMode;
  effective: EffectiveTheme;
  setMode: (mode: ThemeMode, x?: number, y?: number) => void;
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
 * Apply theme class to HTML elements with optional ripple effect
 */
const applyThemeClass = (theme: EffectiveTheme, x?: number, y?: number): void => {
  if (typeof document === 'undefined') return;

  const html = document.documentElement;
  const body = document.body;

  // 如果提供了坐标，使用扩散动画
  if (typeof x === 'number' && typeof y === 'number') {
    applyThemeWithRipple(theme, x, y);
    return;
  }

  // 否则直接切换
  if (theme === 'dark') {
    html.classList.add(DARK_THEME_CLASS);
    body.classList.add(DARK_THEME_CLASS);
  } else {
    html.classList.remove(DARK_THEME_CLASS);
    body.classList.remove(DARK_THEME_CLASS);
  }
};

/**
 * Apply theme with ripple animation effect
 * 使用圆形扩散 + 透明度渐变实现平滑过渡
 */
const applyThemeWithRipple = (theme: EffectiveTheme, x: number, y: number): void => {
  const html = document.documentElement;
  const body = document.body;

  // 计算需要的最大半径（确保完全覆盖屏幕）
  const maxDistance = Math.hypot(
    Math.max(x, window.innerWidth - x),
    Math.max(y, window.innerHeight - y),
  );
  const finalSize = maxDistance * 2.5;

  // 创建扩散遮罩 - 使用目标主题的背景色 + 透明度动画
  const ripple = document.createElement('div');
  ripple.style.cssText = `
    position: fixed;
    left: ${x}px;
    top: ${y}px;
    width: 0;
    height: 0;
    border-radius: 50%;
    background-color: ${theme === 'dark' ? 'var(--bg-inverse)' : 'var(--bg-primary)'};
    transform: translate(-50%, -50%) scale(0);
    opacity: 1;
    pointer-events: none;
    z-index: 99999;
    transition: transform 0.8s cubic-bezier(0.4, 0, 0.2, 1),
                opacity 0.8s cubic-bezier(0.4, 0, 0.2, 1);
  `;

  document.body.appendChild(ripple);

  // 强制重排
  ripple.offsetHeight;

  // 触发扩散动画
  requestAnimationFrame(() => {
    ripple.style.width = `${finalSize * 2}px`;
    ripple.style.height = `${finalSize * 2}px`;
    ripple.style.transform = 'translate(-50%, -50%) scale(1)';
  });

  // 关键：等扩散完全覆盖屏幕后再切换主题
  setTimeout(() => {
    // 切换主题类名
    if (theme === 'dark') {
      html.classList.add(DARK_THEME_CLASS);
      body.classList.add(DARK_THEME_CLASS);
    } else {
      html.classList.remove(DARK_THEME_CLASS);
      body.classList.remove(DARK_THEME_CLASS);
    }

    // 主题切换后立即开始淡出遮罩
    ripple.style.opacity = '0';
  }, 600);

  // 淡出动画完成后移除遮罩
  setTimeout(() => {
    ripple.remove();
  }, 1400); // 600ms 扩散 + 800ms 淡出
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

  // 用于标记是否正在进行动画切换
  const isAnimatingRef = useRef(false);

  // Apply theme when mode changes (without ripple for system changes)
  useEffect(() => {
    // 如果正在进行动画切换，跳过立即切换
    if (isAnimatingRef.current) {
      isAnimatingRef.current = false;
      return;
    }

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

  const handleSetMode = useCallback((newMode: ThemeMode, x?: number, y?: number) => {
    // 如果提供了坐标，使用动画切换
    if (typeof x === 'number' && typeof y === 'number') {
      // 标记正在进行动画切换，阻止 useEffect 的立即切换
      isAnimatingRef.current = true;

      const systemTheme = getSystemColorScheme();
      const effectiveTheme = newMode === 'system' ? systemTheme : newMode;

      // 更新 mode（会触发 useEffect，但被 isAnimatingRef 阻止）
      setModeState(newMode);
      storage.setItem(THEME_STORAGE_KEY, newMode);
      setEffective(effectiveTheme);

      // 执行带动画的主题切换
      applyThemeWithRipple(effectiveTheme, x, y);
    } else {
      // 没有坐标，直接切换（无动画）
      setModeState(newMode);
      storage.setItem(THEME_STORAGE_KEY, newMode);
    }
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
