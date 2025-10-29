import React, { createContext, useContext, useState, useCallback, useMemo } from 'react';
import { zhCN, type TranslationKeys } from '../i18n/locales/zh-CN';
import { enUS } from '../i18n/locales/en-US';

export type SupportedLocale = 'zh-CN' | 'en-US';

const LOCALE_STORAGE_KEY = 'gamelink_locale';

interface I18nContextValue {
  locale: SupportedLocale;
  t: TranslationKeys;
  setLocale: (locale: SupportedLocale) => void;
}

const I18nContext = createContext<I18nContextValue | null>(null);

const localeMap: Record<SupportedLocale, TranslationKeys> = {
  'zh-CN': zhCN,
  'en-US': enUS,
};

const detectBrowserLocale = (): SupportedLocale => {
  const browserLang = navigator.language;
  if (browserLang.startsWith('zh')) return 'zh-CN';
  if (browserLang.startsWith('en')) return 'en-US';
  return 'zh-CN';
};

export const I18nProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [locale, setLocaleState] = useState<SupportedLocale>(() => {
    const saved = localStorage.getItem(LOCALE_STORAGE_KEY);
    if (saved === 'zh-CN' || saved === 'en-US') {
      return saved;
    }
    return detectBrowserLocale();
  });

  const setLocale = useCallback((newLocale: SupportedLocale) => {
    setLocaleState(newLocale);
    localStorage.setItem(LOCALE_STORAGE_KEY, newLocale);
  }, []);

  const t = useMemo(() => localeMap[locale], [locale]);

  const value = useMemo<I18nContextValue>(
    () => ({
      locale,
      t,
      setLocale,
    }),
    [locale, t, setLocale],
  );

  return <I18nContext.Provider value={value}>{children}</I18nContext.Provider>;
};

export const useI18n = (): I18nContextValue => {
  const context = useContext(I18nContext);
  if (!context) {
    throw new Error('useI18n must be used within I18nProvider');
  }
  return context;
};
