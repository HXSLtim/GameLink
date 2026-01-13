import { Moon, Sun, Laptop } from "lucide-react"
import { useTranslation } from "react-i18next"
import { Button } from "@/components/ui/button"
import { useThemeStore } from "@/stores"

export function ModeToggle() {
    const { theme, toggleTheme } = useThemeStore()
    const { t } = useTranslation()

    const themeName = {
        day: t('app.theme.light'),
        night: t('app.theme.dark'),
        auto: t('app.theme.system')
    }[theme];

    return (
        <Button
            variant="ghost"
            size="icon"
            className="h-9 w-9 relative rounded-full bg-background/60 backdrop-blur-md border border-border/40 hover:bg-accent hover:text-accent-foreground transition-all"
            onClick={toggleTheme}
            title={`${t('app.theme.toggle')} (${themeName})`}
        >
            <Sun className={`h-[1.2rem] w-[1.2rem] transition-all ${theme === 'day' ? 'rotate-0 scale-100' : '-rotate-90 scale-0 absolute'}`} />
            <Moon className={`h-[1.2rem] w-[1.2rem] transition-all ${theme === 'night' ? 'rotate-0 scale-100' : 'rotate-90 scale-0 absolute'}`} />
            <Laptop className={`h-[1.2rem] w-[1.2rem] transition-all ${theme === 'auto' ? 'rotate-0 scale-100' : 'scale-0 absolute'}`} />
            <span className="sr-only">{t('app.theme.toggle')}</span>
        </Button>
    )
}
