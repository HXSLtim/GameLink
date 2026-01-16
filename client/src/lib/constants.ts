/**
 * Application-wide constants
 * Extracted from various files to avoid magic numbers
 */

// UI Timing Constants
export const SCROLL_DELAY_MS = 50;
export const BOOKING_REDIRECT_DELAY_MS = 1500;
export const FILTER_DEBOUNCE_MS = 500;

// WebSocket Constants
export const HEARTBEAT_INTERVAL_MS = 30000;
export const HEARTBEAT_TIMEOUT_MS = 10000;
export const MAX_RECONNECT_ATTEMPTS = 5;

// Auth Constants
export const TOKEN_REFRESH_BUFFER_MS = 60000;

// Pagination Constants
export const DEFAULT_PAGE_SIZE = 20;
export const MAX_PAGE_SIZE = 100;

// Cache Constants
export const CACHE_TTL_MS = 5 * 60 * 1000; // 5 minutes
