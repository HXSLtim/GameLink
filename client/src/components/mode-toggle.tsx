import { Moon, Sun, Laptop } from "lucide-react"
import { Button } from "@/components/ui/button"
import { useThemeStore } from "@/stores"

export function ModeToggle() {
    const { theme, toggleTheme } = useThemeStore()

    return (
        <Button
            variant="ghost"
            size="icon"
            className="h-9 w-9 relative"
            onClick={toggleTheme}
            title={`Current theme: ${theme}. Click to switch.`}
        >
            <Sun className={`h-[1.2rem] w-[1.2rem] transition-all ${theme === 'day' ? 'rotate-0 scale-100' : '-rotate-90 scale-0 absolute'}`} />
            <Moon className={`h-[1.2rem] w-[1.2rem] transition-all ${theme === 'night' ? 'rotate-0 scale-100' : 'rotate-90 scale-0 absolute'}`} />
            <Laptop className={`h-[1.2rem] w-[1.2rem] transition-all ${theme === 'auto' ? 'rotate-0 scale-100' : 'scale-0 absolute'}`} />
            <span className="sr-only">Toggle theme</span>
        </Button>
    )
}
