const CACHE_CLEANUP_VERSION = 'v1';
const CACHE_CLEANUP_MARK = `gamelink_admin_sw_cleanup_${CACHE_CLEANUP_VERSION}`;

const LEGACY_CACHE_PATTERNS = [
  /^workbox/i,
  /precache/i,
  /^api-cache$/i,
  /^image-cache$/i,
  /^static-resources-cache$/i,
  /^vite-plugin-pwa/i,
];

function isLegacyPwaCache(cacheName: string): boolean {
  return LEGACY_CACHE_PATTERNS.some((pattern) => pattern.test(cacheName));
}

export async function cleanupLegacyServiceWorker(): Promise<void> {
  if (typeof window === 'undefined') {
    return;
  }

  if (window.localStorage.getItem(CACHE_CLEANUP_MARK)) {
    return;
  }

  try {
    if ('serviceWorker' in navigator && navigator.serviceWorker.getRegistrations) {
      const registrations = await navigator.serviceWorker.getRegistrations();
      await Promise.all(registrations.map((registration) => registration.unregister()));
    }

    if ('caches' in window && caches.keys) {
      const cacheNames = await caches.keys();
      const legacyCacheNames = cacheNames.filter(isLegacyPwaCache);
      await Promise.all(legacyCacheNames.map((cacheName) => caches.delete(cacheName)));
    }

    window.localStorage.setItem(CACHE_CLEANUP_MARK, '1');
  } catch {
    // Ignore cleanup failures: avoid blocking app startup.
  }
}
