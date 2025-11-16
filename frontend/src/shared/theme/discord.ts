/**
 * Discord/KOOK 风格主题配置
 */

export const discordTheme = {
  dark: {
    // 背景颜色
    background: {
      primary: '#36393f',      // Discord经典深灰背景
      secondary: '#2f3136',    // 卡片/面板背景
      tertiary: '#202225',     // 侧边栏背景
      elevated: '#18191c',     // 弹出层/模态框
    },

    // 文字颜色
    text: {
      primary: '#dcddde',      // 主要文字
      secondary: '#b9bbbe',    // 次要文字
      muted: '#72767d',        // 辅助文字
      link: '#00aff4',         // 链接颜色
    },

    // 品牌色
    brand: {
      primary: '#5865f2',      // Discord蓝 - 主要操作按钮
      success: '#3ba55d',      // 成功/在线状态
      warning: '#faa61a',      // 警告
      danger: '#ed4245',       // 错误/删除
    },

    // KOOK绿色备选
    kook: {
      primary: '#6DD230',      // KOOK主题绿
      secondary: '#53C41A',    // KOOK深绿
    },

    // 状态颜色
    status: {
      online: '#3ba55d',       // 在线
      idle: '#faa61a',         // 忙碌
      dnd: '#ed4245',          // 游戏中
      offline: '#747f8d',      // 离线
    },

    // 功能色
    functional: {
      hoverOverlay: 'rgba(255, 255, 255, 0.06)',
      activeOverlay: 'rgba(255, 255, 255, 0.1)',
      divider: 'rgba(255, 255, 255, 0.08)',
    },
  },

  light: {
    background: {
      primary: '#ffffff',
      secondary: '#f2f3f5',
      tertiary: '#e3e5e8',
      elevated: '#ffffff',
    },

    text: {
      primary: '#2e3338',
      secondary: '#4e5058',
      muted: '#747f8d',
      link: '#00aff4',
    },

    brand: {
      primary: '#5865f2',
      success: '#3ba55d',
      warning: '#faa61a',
      danger: '#ed4245',
    },

    kook: {
      primary: '#6DD230',
      secondary: '#53C41A',
    },

    status: {
      online: '#3ba55d',
      idle: '#faa61a',
      dnd: '#ed4245',
      offline: '#747f8d',
    },

    functional: {
      hoverOverlay: 'rgba(0, 0, 0, 0.04)',
      activeOverlay: 'rgba(0, 0, 0, 0.08)',
      divider: 'rgba(0, 0, 0, 0.08)',
    },
  },
};

// 导出CSS变量
export const getCSSVariables = (mode: 'dark' | 'light' = 'dark') => {
  const theme = discordTheme[mode];

  return {
    // 背景色
    '--primary-bg': theme.background.primary,
    '--secondary-bg': theme.background.secondary,
    '--tertiary-bg': theme.background.tertiary,
    '--elevated-bg': theme.background.elevated,

    // 文字色
    '--text-primary': theme.text.primary,
    '--text-secondary': theme.text.secondary,
    '--text-muted': theme.text.muted,
    '--text-link': theme.text.link,

    // 品牌色
    '--brand-primary': theme.brand.primary,
    '--brand-success': theme.brand.success,
    '--brand-warning': theme.brand.warning,
    '--brand-danger': theme.brand.danger,

    // KOOK绿
    '--kook-primary': theme.kook.primary,
    '--kook-secondary': theme.kook.secondary,

    // 状态色
    '--status-online': theme.status.online,
    '--status-idle': theme.status.idle,
    '--status-dnd': theme.status.dnd,
    '--status-offline': theme.status.offline,

    // 功能色
    '--hover-overlay': theme.functional.hoverOverlay,
    '--active-overlay': theme.functional.activeOverlay,
    '--divider': theme.functional.divider,

    // 字体
    '--font-display': '"Noto Sans SC", "PingFang SC", "Microsoft YaHei", sans-serif',
    '--font-mono': '"JetBrains Mono", "Consolas", monospace',

    // 字号
    '--text-xs': '12px',
    '--text-sm': '14px',
    '--text-base': '16px',
    '--text-lg': '18px',
    '--text-xl': '20px',
    '--text-2xl': '24px',
    '--text-3xl': '30px',

    // 间距
    '--space-1': '4px',
    '--space-2': '8px',
    '--space-3': '12px',
    '--space-4': '16px',
    '--space-5': '20px',
    '--space-6': '24px',
    '--space-8': '32px',
    '--space-10': '40px',
    '--space-12': '48px',
    '--space-16': '64px',

    // 圆角
    '--radius-sm': '4px',
    '--radius-md': '8px',
    '--radius-lg': '12px',
    '--radius-xl': '16px',
    '--radius-full': '9999px',

    // 阴影
    '--shadow-sm': '0 1px 2px 0 rgba(0, 0, 0, 0.05)',
    '--shadow-md': '0 4px 6px -1px rgba(0, 0, 0, 0.1)',
    '--shadow-lg': '0 10px 15px -3px rgba(0, 0, 0, 0.1)',
    '--shadow-xl': '0 20px 25px -5px rgba(0, 0, 0, 0.1)',
    '--shadow-2xl': '0 25px 50px -12px rgba(0, 0, 0, 0.25)',
  };
};
