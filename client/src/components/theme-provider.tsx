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
        const root = window.document.documentElement

        root.classList.remove("light", "dark")

        if (theme === "auto") {
            const systemTheme = window.matchMedia("(prefers-color-scheme: dark)")
                .matches
                ? "dark"
                : "light"

            root.classList.add(systemTheme)
            return
        }

        // 'day' maps to light theme tokens, 'night' maps to dark theme tokens
        // For Tailwind class compat, we add 'light' or 'dark' to root
        const rootClass = theme === 'day' ? 'light' : 'dark'
        root.classList.add(rootClass)

    }, [theme])

    return <>{children}</>
}
