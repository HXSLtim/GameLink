import { useState, useEffect } from 'react';

/**
 * Hook to detect media query matches
 */
export function useMediaQuery(query: string): boolean {
    const [matches, setMatches] = useState(() => {
        if (typeof window !== 'undefined') {
            return window.matchMedia(query).matches;
        }
        return false;
    });

    useEffect(() => {
        const mediaQuery = window.matchMedia(query);
        const handler = (event: MediaQueryListEvent) => setMatches(event.matches);

        // Listen for changes
        mediaQuery.addEventListener('change', handler);
        return () => mediaQuery.removeEventListener('change', handler);
    }, [query]);

    return matches;
}

// Preset breakpoints (matching Tailwind defaults)
export function useIsMobile() {
    return useMediaQuery('(max-width: 639px)');
}

export function useIsTablet() {
    return useMediaQuery('(min-width: 640px) and (max-width: 1023px)');
}

export function useIsDesktop() {
    return useMediaQuery('(min-width: 1024px)');
}
