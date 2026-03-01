type ErrorContext = Record<string, unknown>;

interface SentryLike {
  captureException: (error: unknown, context?: { tags?: Record<string, string>; extra?: ErrorContext }) => void;
}

declare global {
  interface Window {
    Sentry?: SentryLike;
  }
}

const RUNTIME_MONITORING_FLAG = "__gamelink_runtime_monitoring_initialized__";

function isMonitoringEnabled(): boolean {
  return import.meta.env.VITE_ENABLE_ERROR_MONITORING === "true";
}

function getReleaseTag(): string {
  return import.meta.env.VITE_SENTRY_RELEASE || import.meta.env.VITE_APP_VERSION || "unknown";
}

function reportToSentry(error: unknown, context?: ErrorContext): void {
  window.Sentry?.captureException(error, {
    tags: {
      app: "gamelink-app",
      release: getReleaseTag(),
      env: import.meta.env.VITE_SENTRY_ENVIRONMENT || import.meta.env.MODE || "unknown",
    },
    extra: context,
  });
}

export function reportRuntimeError(error: unknown, context?: ErrorContext): void {
  if (!isMonitoringEnabled()) {
    return;
  }

  if (window.Sentry && import.meta.env.VITE_SENTRY_DSN) {
    reportToSentry(error, context);
    return;
  }

  if (import.meta.env.DEV) {
    console.error("[runtime-monitoring:fallback]", error, context);
  }
}

export function setupGlobalErrorMonitoring(): void {
  if (typeof window === "undefined") {
    return;
  }

  const globalState = window as typeof window & Record<string, unknown>;
  if (globalState[RUNTIME_MONITORING_FLAG]) {
    return;
  }

  globalState[RUNTIME_MONITORING_FLAG] = true;

  window.addEventListener("error", (event) => {
    reportRuntimeError(event.error || event.message, {
      type: "window.error",
      filename: event.filename,
      lineno: event.lineno,
      colno: event.colno,
    });
  });

  window.addEventListener("unhandledrejection", (event) => {
    reportRuntimeError(event.reason, {
      type: "window.unhandledrejection",
    });
  });
}
