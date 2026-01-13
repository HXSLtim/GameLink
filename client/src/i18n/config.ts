import i18n from 'i18next';
import { initReactI18next } from 'react-i18next';
import LanguageDetector from 'i18next-browser-languagedetector';
import Backend from 'i18next-http-backend';

import zhCN from './locales/zh-CN.json';
import enUS from './locales/en-US.json';

i18n
    // load translation using http -> see /public/locales (i.e. https://github.com/i18next/react-i18next/tree/master/example/react/public/locales)
    // learn more: https://github.com/i18next/i18next-http-backend
    // want your translations to be loaded from a professional CDN? => https://github.com/locize/react-tutorial#step-2---use-the-locize-cdn
    .use(Backend)
    // detect user language
    // learn more: https://github.com/i18next/i18next-browser-languagedetector
    .use(LanguageDetector)
    // pass the i18n instance to react-i18next.
    .use(initReactI18next)
    // init i18next
    // for all options read: https://www.i18next.com/overview/configuration-options
    .init({
        resources: {
            'zh-CN': { translation: zhCN },
            'en-US': { translation: enUS }
        },
        fallbackLng: 'zh-CN',
        debug: import.meta.env.DEV,

        interpolation: {
            escapeValue: false, // not needed for react as it escapes by default
        },

        // Check for user's language preference from settings store logic could be added here or in effect
    });

export default i18n;
