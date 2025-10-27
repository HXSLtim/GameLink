/**
 * 国际化工具模块
 * 提供基础的 i18n 支持，未来可集成 react-i18next
 */

export type SupportedLocale = 'zh-CN' | 'en-US';

export const DEFAULT_LOCALE: SupportedLocale = 'zh-CN';

/**
 * 获取当前语言设置
 */
export function getCurrentLocale(): SupportedLocale {
  // 从 localStorage 读取
  const saved = localStorage.getItem('gamelink_locale');
  if (saved === 'zh-CN' || saved === 'en-US') {
    return saved;
  }

  // 从浏览器语言推断
  const browserLang = navigator.language;
  if (browserLang.startsWith('zh')) {
    return 'zh-CN';
  }
  if (browserLang.startsWith('en')) {
    return 'en-US';
  }

  return DEFAULT_LOCALE;
}

/**
 * 设置语言
 */
export function setLocale(locale: SupportedLocale): void {
  localStorage.setItem('gamelink_locale', locale);
  // 触发语言变更事件
  window.dispatchEvent(new CustomEvent('locale-change', { detail: locale }));
}

/**
 * 简单的翻译函数（占位）
 * 未来可替换为 i18next 的 t 函数
 */
export function t(key: string, defaultValue?: string): string {
  // 当前简化实现，直接返回默认值或 key
  return defaultValue || key;
}

/**
 * i18n 工具对象
 */
export const i18n = {
  getCurrentLocale,
  setLocale,
  t,
  DEFAULT_LOCALE,
};

