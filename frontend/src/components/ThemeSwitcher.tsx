import { useMemo, useCallback } from 'react';
import { Select, Badge } from '@arco-design/web-react';
import { IconMoon, IconSun, IconDesktop } from '@arco-design/web-react/icon';
import type { ThemeMode } from '../contexts/ThemeContext';
import { useTheme } from '../contexts/ThemeContext';
import styles from './ThemeSwitcher.module.less';

interface ThemeOption {
  label: React.ReactNode;
  value: ThemeMode;
}

const ICON_STYLE = { verticalAlign: -2, marginRight: 6 };

/**
 * Theme switcher component
 *
 * @component
 * @description Allows users to switch between light, dark, and system theme modes
 */
export const ThemeSwitcher = () => {
  const { mode, effective, setMode } = useTheme();

  const options = useMemo<ThemeOption[]>(
    () => [
      {
        value: 'system',
        label: (
          <span>
            <IconDesktop style={ICON_STYLE} /> 跟随系统
          </span>
        ),
      },
      {
        value: 'light',
        label: (
          <span>
            <IconSun style={ICON_STYLE} /> 亮色
          </span>
        ),
      },
      {
        value: 'dark',
        label: (
          <span>
            <IconMoon style={ICON_STYLE} /> 暗色
          </span>
        ),
      },
    ],
    [],
  );

  const handleModeChange = useCallback(
    (value: ThemeMode) => {
      setMode(value);
    },
    [setMode],
  );

  const badgeStatus = effective === 'dark' ? 'processing' : 'success';
  const badgeText = effective === 'dark' ? 'Dark' : 'Light';

  return (
    <div className={styles.container}>
      <Badge status={badgeStatus} text={badgeText} />
      <Select
        size="small"
        value={mode}
        className={styles.select}
        onChange={handleModeChange}
        options={options}
        triggerProps={{ position: 'bl' }}
      />
    </div>
  );
};
