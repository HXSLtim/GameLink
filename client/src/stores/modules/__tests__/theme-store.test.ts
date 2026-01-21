/**
 * Theme Store Tests
 * Tests for theme management and UI preferences
 */

import { describe, it, expect, beforeEach } from 'vitest';
import { useThemeStore } from '../theme-store';

describe('Theme Store', () => {
    beforeEach(() => {
        // Reset store to initial state
        useThemeStore.setState({
            theme: 'night',
            primaryColor: 'hsl(var(--primary))',
            fontSize: 'md',
            sidebarCollapsed: false,
            soundEnabled: true,
        });
    });

    describe('Initial State', () => {
        it('should have correct initial state', () => {
            const state = useThemeStore.getState();

            expect(state.theme).toBe('night');
            expect(state.fontSize).toBe('md');
            expect(state.sidebarCollapsed).toBe(false);
            expect(state.soundEnabled).toBe(true);
        });
    });

    describe('Theme Management', () => {
        it('should set theme to day', () => {
            useThemeStore.getState().setTheme('day');

            expect(useThemeStore.getState().theme).toBe('day');
        });

        it('should set theme to night', () => {
            useThemeStore.getState().setTheme('night');

            expect(useThemeStore.getState().theme).toBe('night');
        });

        it('should set theme to auto', () => {
            useThemeStore.getState().setTheme('auto');

            expect(useThemeStore.getState().theme).toBe('auto');
        });

        it('should toggle theme from night to auto', () => {
            useThemeStore.setState({ theme: 'night' });

            useThemeStore.getState().toggleTheme();

            expect(useThemeStore.getState().theme).toBe('auto');
        });

        it('should toggle theme from auto to day', () => {
            useThemeStore.setState({ theme: 'auto' });

            useThemeStore.getState().toggleTheme();

            expect(useThemeStore.getState().theme).toBe('day');
        });

        it('should toggle theme from day to night', () => {
            useThemeStore.setState({ theme: 'day' });

            useThemeStore.getState().toggleTheme();

            expect(useThemeStore.getState().theme).toBe('night');
        });

        it('should cycle through all themes', () => {
            const store = useThemeStore.getState();

            // night -> auto
            expect(store.theme).toBe('night');
            useThemeStore.getState().toggleTheme();
            expect(useThemeStore.getState().theme).toBe('auto');

            // auto -> day
            useThemeStore.getState().toggleTheme();
            expect(useThemeStore.getState().theme).toBe('day');

            // day -> night
            useThemeStore.getState().toggleTheme();
            expect(useThemeStore.getState().theme).toBe('night');

            // night -> auto again (cycle continues)
            useThemeStore.getState().toggleTheme();
            expect(useThemeStore.getState().theme).toBe('auto');
        });
    });

    describe('Font Size Management', () => {
        it('should set font size to small', () => {
            useThemeStore.getState().setFontSize('sm');

            expect(useThemeStore.getState().fontSize).toBe('sm');
        });

        it('should set font size to medium', () => {
            useThemeStore.getState().setFontSize('md');

            expect(useThemeStore.getState().fontSize).toBe('md');
        });

        it('should set font size to large', () => {
            useThemeStore.getState().setFontSize('lg');

            expect(useThemeStore.getState().fontSize).toBe('lg');
        });

        it('should change font size from small to medium', () => {
            useThemeStore.setState({ fontSize: 'sm' });

            useThemeStore.getState().setFontSize('md');

            expect(useThemeStore.getState().fontSize).toBe('md');
        });

        it('should change font size from medium to large', () => {
            useThemeStore.setState({ fontSize: 'md' });

            useThemeStore.getState().setFontSize('lg');

            expect(useThemeStore.getState().fontSize).toBe('lg');
        });
    });

    describe('Sidebar Management', () => {
        it('should toggle sidebar from expanded to collapsed', () => {
            useThemeStore.setState({ sidebarCollapsed: false });

            useThemeStore.getState().toggleSidebar();

            expect(useThemeStore.getState().sidebarCollapsed).toBe(true);
        });

        it('should toggle sidebar from collapsed to expanded', () => {
            useThemeStore.setState({ sidebarCollapsed: true });

            useThemeStore.getState().toggleSidebar();

            expect(useThemeStore.getState().sidebarCollapsed).toBe(false);
        });

        it('should toggle sidebar multiple times', () => {
            const store = useThemeStore.getState();

            expect(store.sidebarCollapsed).toBe(false);

            useThemeStore.getState().toggleSidebar();
            expect(useThemeStore.getState().sidebarCollapsed).toBe(true);

            useThemeStore.getState().toggleSidebar();
            expect(useThemeStore.getState().sidebarCollapsed).toBe(false);

            useThemeStore.getState().toggleSidebar();
            expect(useThemeStore.getState().sidebarCollapsed).toBe(true);
        });
    });

    describe('Sound Management', () => {
        it('should toggle sound from enabled to disabled', () => {
            useThemeStore.setState({ soundEnabled: true });

            useThemeStore.getState().toggleSound();

            expect(useThemeStore.getState().soundEnabled).toBe(false);
        });

        it('should toggle sound from disabled to enabled', () => {
            useThemeStore.setState({ soundEnabled: false });

            useThemeStore.getState().toggleSound();

            expect(useThemeStore.getState().soundEnabled).toBe(true);
        });

        it('should toggle sound multiple times', () => {
            const store = useThemeStore.getState();

            expect(store.soundEnabled).toBe(true);

            useThemeStore.getState().toggleSound();
            expect(useThemeStore.getState().soundEnabled).toBe(false);

            useThemeStore.getState().toggleSound();
            expect(useThemeStore.getState().soundEnabled).toBe(true);

            useThemeStore.getState().toggleSound();
            expect(useThemeStore.getState().soundEnabled).toBe(false);
        });
    });

    describe('Combined Actions', () => {
        it('should handle multiple theme changes', () => {
            useThemeStore.getState().setTheme('day');
            expect(useThemeStore.getState().theme).toBe('day');

            useThemeStore.getState().setTheme('auto');
            expect(useThemeStore.getState().theme).toBe('auto');

            useThemeStore.getState().setTheme('night');
            expect(useThemeStore.getState().theme).toBe('night');
        });

        it('should handle multiple font size changes', () => {
            useThemeStore.getState().setFontSize('sm');
            expect(useThemeStore.getState().fontSize).toBe('sm');

            useThemeStore.getState().setFontSize('lg');
            expect(useThemeStore.getState().fontSize).toBe('lg');

            useThemeStore.getState().setFontSize('md');
            expect(useThemeStore.getState().fontSize).toBe('md');
        });

        it('should maintain independent state for different settings', () => {
            useThemeStore.setState({
                theme: 'day',
                fontSize: 'lg',
                sidebarCollapsed: true,
                soundEnabled: false,
            });

            const state = useThemeStore.getState();

            expect(state.theme).toBe('day');
            expect(state.fontSize).toBe('lg');
            expect(state.sidebarCollapsed).toBe(true);
            expect(state.soundEnabled).toBe(false);

            // Change one setting should not affect others
            useThemeStore.getState().toggleTheme();

            expect(useThemeStore.getState().theme).not.toBe('day');
            expect(useThemeStore.getState().fontSize).toBe('lg');
            expect(useThemeStore.getState().sidebarCollapsed).toBe(true);
            expect(useThemeStore.getState().soundEnabled).toBe(false);
        });
    });

    describe('Edge Cases', () => {
        it('should handle rapid theme toggles', () => {
            useThemeStore.setState({ theme: 'night' });

            useThemeStore.getState().toggleTheme();
            useThemeStore.getState().toggleTheme();
            useThemeStore.getState().toggleTheme();
            useThemeStore.getState().toggleTheme();

            // night -> auto -> day -> night -> auto
            expect(useThemeStore.getState().theme).toBe('auto');
        });

        it('should handle rapid sidebar toggles', () => {
            useThemeStore.setState({ sidebarCollapsed: false });

            // 10 toggles: false -> true -> false -> ... (even number = false)
            for (let i = 0; i < 10; i++) {
                useThemeStore.getState().toggleSidebar();
            }

            expect(useThemeStore.getState().sidebarCollapsed).toBe(false);
        });

        it('should handle all theme values', () => {
            ['day', 'night', 'auto'].forEach(themeValue => {
                useThemeStore.getState().setTheme(themeValue as any);
                expect(useThemeStore.getState().theme).toBe(themeValue);
            });
        });

        it('should handle all font size values', () => {
            ['sm', 'md', 'lg'].forEach(size => {
                useThemeStore.getState().setFontSize(size as any);
                expect(useThemeStore.getState().fontSize).toBe(size);
            });
        });
    });
});
