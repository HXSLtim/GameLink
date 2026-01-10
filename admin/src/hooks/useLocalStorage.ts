import { useState, useEffect, useCallback } from 'react';

/**
 * localStorage Hook
 * 用于持久化用户偏好设置
 *
 * @param key localStorage 键名
 * @param initialValue 初始值
 * @returns [value, setValue, removeValue]
 */
export function useLocalStorage<T>(
  key: string,
  initialValue: T
): [T, (value: T | ((prev: T) => T)) => void, () => void] {
  // 获取初始值
  const readValue = useCallback((): T => {
    if (typeof window === 'undefined') {
      return initialValue;
    }

    try {
      const item = window.localStorage.getItem(key);
      return item ? (JSON.parse(item) as T) : initialValue;
    } catch (error) {
      console.warn(`Error reading localStorage key "${key}":`, error);
      return initialValue;
    }
  }, [initialValue, key]);

  const [storedValue, setStoredValue] = useState<T>(readValue);

  // 设置值
  const setValue = useCallback(
    (value: T | ((prev: T) => T)) => {
      try {
        const valueToStore = value instanceof Function ? value(storedValue) : value;
        setStoredValue(valueToStore);
        if (typeof window !== 'undefined') {
          window.localStorage.setItem(key, JSON.stringify(valueToStore));
          // 触发自定义事件，用于跨标签页同步
          window.dispatchEvent(new StorageEvent('storage', { key, newValue: JSON.stringify(valueToStore) }));
        }
      } catch (error) {
        console.warn(`Error setting localStorage key "${key}":`, error);
      }
    },
    [key, storedValue]
  );

  // 移除值
  const removeValue = useCallback(() => {
    try {
      if (typeof window !== 'undefined') {
        window.localStorage.removeItem(key);
        setStoredValue(initialValue);
      }
    } catch (error) {
      console.warn(`Error removing localStorage key "${key}":`, error);
    }
  }, [initialValue, key]);

  // 监听其他标签页的变化
  useEffect(() => {
    const handleStorageChange = (event: StorageEvent) => {
      if (event.key === key && event.newValue !== null) {
        try {
          setStoredValue(JSON.parse(event.newValue) as T);
        } catch {
          // 忽略解析错误
        }
      }
    };

    window.addEventListener('storage', handleStorageChange);
    return () => window.removeEventListener('storage', handleStorageChange);
  }, [key]);

  return [storedValue, setValue, removeValue];
}

/**
 * 仪表盘偏好设置类型
 */
export interface DashboardPreferences {
  /** 折叠的区块 */
  collapsedSections: string[];
  /** 隐藏的小部件 */
  hiddenWidgets: string[];
  /** 小部件顺序 */
  widgetOrder?: string[];
}

/**
 * 仪表盘偏好 Hook
 */
export function useDashboardPreferences() {
  const [preferences, setPreferences, resetPreferences] = useLocalStorage<DashboardPreferences>(
    'gamelink_dashboard_preferences',
    {
      collapsedSections: [],
      hiddenWidgets: [],
      widgetOrder: undefined,
    }
  );

  const toggleSection = useCallback(
    (sectionId: string) => {
      setPreferences((prev) => ({
        ...prev,
        collapsedSections: prev.collapsedSections.includes(sectionId)
          ? prev.collapsedSections.filter((id) => id !== sectionId)
          : [...prev.collapsedSections, sectionId],
      }));
    },
    [setPreferences]
  );

  const toggleWidget = useCallback(
    (widgetId: string) => {
      setPreferences((prev) => ({
        ...prev,
        hiddenWidgets: prev.hiddenWidgets.includes(widgetId)
          ? prev.hiddenWidgets.filter((id) => id !== widgetId)
          : [...prev.hiddenWidgets, widgetId],
      }));
    },
    [setPreferences]
  );

  const isSectionCollapsed = useCallback(
    (sectionId: string) => preferences.collapsedSections.includes(sectionId),
    [preferences.collapsedSections]
  );

  const isWidgetHidden = useCallback(
    (widgetId: string) => preferences.hiddenWidgets.includes(widgetId),
    [preferences.hiddenWidgets]
  );

  return {
    preferences,
    setPreferences,
    resetPreferences,
    toggleSection,
    toggleWidget,
    isSectionCollapsed,
    isWidgetHidden,
  };
}

export default useLocalStorage;
