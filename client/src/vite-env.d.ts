/// <reference types="vite/client" />

interface ImportMetaEnv {
    readonly VITE_API_BASE_URL: string;
    readonly VITE_WS_URL: string;
    readonly VITE_APP_ENV: 'development' | 'staging' | 'production';
    readonly VITE_ENABLE_MOCK: string;
    readonly VITE_SENTRY_DSN?: string;
    readonly VITE_ANALYTICS_ID?: string;
}

interface ImportMeta {
    readonly env: ImportMetaEnv;
}
