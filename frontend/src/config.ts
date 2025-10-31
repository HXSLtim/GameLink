export const API_BASE = '/api/v1';
export const STORAGE_KEYS = {
  token: 'gamelink_token',
  user: 'gamelink_user',
  theme: 'gamelink_theme',
  language: 'gamelink_language',
};

// 特性开关配置
export const FEATURE_FLAGS = {
  showcase: {
    // 是否默认展开组件使用示例
    expandExamplesByDefault: true,
    // 是否启用独立的 /showcase 路由（生产默认关闭）
    enableShowcaseRoute: !import.meta.env.PROD,
    // 是否在主布局中启用 /components 子路由与菜单项（生产默认关闭）
    enableComponentsRoute: !import.meta.env.PROD,
  },
  cacheDemo: {
    // 是否启用 /cache-demo 演示路由（生产默认关闭）
    enableCacheDemoRoute: !import.meta.env.PROD,
  },
};
