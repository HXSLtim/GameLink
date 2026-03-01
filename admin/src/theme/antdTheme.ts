import type { ThemeConfig } from 'antd';

export const GAMELINK_PRIMARY = {
  base: '#1890ff',
  hover: '#40a9ff',
  active: '#096dd9',
  bg: '#e6f7ff',
  bgHover: '#bae7ff',
  border: '#91d5ff',
  borderHover: '#69c0ff',
} as const;

export const gamelinkThemeToken: ThemeConfig['token'] = {
  colorPrimary: GAMELINK_PRIMARY.base,
  colorSuccess: '#52c41a',
  colorWarning: '#faad14',
  colorError: '#ff4d4f',
  colorInfo: '#1677ff',
  colorPrimaryBg: GAMELINK_PRIMARY.bg,
  colorPrimaryBgHover: GAMELINK_PRIMARY.bgHover,
  colorPrimaryBorder: GAMELINK_PRIMARY.border,
  colorPrimaryBorderHover: GAMELINK_PRIMARY.borderHover,
  colorPrimaryHover: GAMELINK_PRIMARY.hover,
  colorPrimaryActive: GAMELINK_PRIMARY.active,
  borderRadius: 8,
  borderRadiusLG: 12,
  borderRadiusSM: 6,
  marginXS: 8,
  marginSM: 12,
  margin: 16,
  marginMD: 20,
  marginLG: 24,
  marginXL: 32,
};

export const gamelinkThemeComponents: ThemeConfig['components'] = {
  Button: {
    colorPrimary: GAMELINK_PRIMARY.base,
    colorPrimaryHover: GAMELINK_PRIMARY.hover,
    colorPrimaryActive: GAMELINK_PRIMARY.active,
    algorithm: true,
  },
  Card: {
    borderRadiusLG: 12,
  },
  Input: {
    borderRadius: 8,
  },
  Select: {
    borderRadius: 8,
  },
};
