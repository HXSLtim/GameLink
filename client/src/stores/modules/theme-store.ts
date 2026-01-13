import { create } from 'zustand';
import { persist } from 'zustand/middleware';

export interface ThemeState {
    theme: 'day' | 'night' | 'auto'; // 'day' maps to light/kook, 'night' maps to dark/discord
    primaryColor: string;
    fontSize: 'sm' | 'md' | 'lg';
    sidebarCollapsed: boolean;
    soundEnabled: boolean;
}

export interface ThemeActions {
    setTheme: (theme: ThemeState['theme']) => void;
    toggleTheme: () => void;
    setFontSize: (size: ThemeState['fontSize']) => void;
    toggleSidebar: () => void;
    toggleSound: () => void;
}

export const useThemeStore = create<ThemeState & ThemeActions>()(
    persist(
        (set, get) => ({
            // Initial State
            theme: 'night', // Default to Dark/Discord
            primaryColor: 'hsl(var(--primary))',
            fontSize: 'md',
            sidebarCollapsed: false,
            soundEnabled: true,

            // Actions
            setTheme: (theme) => set({ theme }),

            toggleTheme: () => {
                const current = get().theme;
                const cycle: Record<string, ThemeState['theme']> = {
                    'night': 'auto',
                    'auto': 'day',
                    'day': 'night'
                };
                set({ theme: cycle[current] });
            },

            setFontSize: (fontSize) => set({ fontSize }),

            toggleSidebar: () => set((state) => ({ sidebarCollapsed: !state.sidebarCollapsed })),

            toggleSound: () => set((state) => ({ soundEnabled: !state.soundEnabled })),
        }),
        {
            name: 'theme-storage',
        }
    )
);
