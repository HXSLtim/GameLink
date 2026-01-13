import { useEffect } from "react"
import { useThemeStore } from "@/stores"

type ThemeProviderProps = {
    children: React.ReactNode
    defaultTheme?: string
    storageKey?: string
}

export function ThemeProvider({
    children,
}: ThemeProviderProps) {
    const { theme } = useThemeStore()

    useEffect(() => {
        const root = window.document.documentElement;

        root.classList.remove("light", "dark");

        if (theme === "auto") {
            const mediaQuery = window.matchMedia("(prefers-color-scheme: dark)");

            const applySystemTheme = () => {
                root.classList.remove("light", "dark");
                root.classList.add(mediaQuery.matches ? "dark" : "light");
            };

            applySystemTheme();

            // Listen for system changes
            mediaQuery.addEventListener("change", applySystemTheme);
            return () => mediaQuery.removeEventListener("change", applySystemTheme);
        }

        const rootClass = theme === 'day' ? 'light' : 'dark';
        root.classList.add(rootClass);
    }, [theme]);

    return <>{children}</>;
}
