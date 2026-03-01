/// <reference types="vite/client" />

interface ImportMetaEnv {
  // API Configuration
  readonly VITE_API_BASE_URL: string

  // Security: Encryption Configuration
  readonly VITE_CRYPTO_ENABLED: string
  readonly VITE_CRYPTO_SECRET_KEY: string
  readonly VITE_CRYPTO_IV: string
  readonly VITE_CRYPTO_USE_SIGNATURE: string

  // Application Configuration
  readonly VITE_APP_TITLE: string
  readonly VITE_APP_VERSION: string
  readonly VITE_DEBUG: string

  // Feature Flags
  readonly VITE_ENABLE_PWA: string
  readonly VITE_ENABLE_WEBSOCKET: string
  readonly VITE_WEBSOCKET_RECONNECT_ATTEMPTS: string
  readonly VITE_WEBSOCKET_RECONNECT_INTERVAL: string
  readonly VITE_WEBSOCKET_HEARTBEAT_INTERVAL: string

  // Error Monitoring (Sentry)
  readonly VITE_SENTRY_ENABLED: string
  readonly VITE_SENTRY_DSN: string
  readonly VITE_SENTRY_ENVIRONMENT: string
  readonly VITE_SENTRY_SAMPLE_RATE: string
  readonly VITE_SENTRY_TRACES_SAMPLE_RATE: string

  // UI Configuration
  readonly VITE_DEFAULT_PAGE_SIZE: string
  readonly VITE_MAX_PAGE_SIZE: string
  readonly VITE_DATE_FORMAT: string
  readonly VITE_TIMEZONE: string
}

interface ImportMeta {
  readonly env: ImportMetaEnv
}
