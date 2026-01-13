import { useTranslation } from 'react-i18next';
import { Button } from '@/components/ui/button';
import { cn } from '@/lib/utils';

interface LanguageSwitcherProps {
    className?: string;
}

export function LanguageSwitcher({ className }: LanguageSwitcherProps) {
    const { i18n } = useTranslation();

    const toggleLanguage = () => {
        const nextLang = i18n.language === 'zh-CN' ? 'en-US' : 'zh-CN';
        i18n.changeLanguage(nextLang);
    };

    const isZh = i18n.language === 'zh-CN';

    return (
        <Button
            variant="ghost"
            size="icon"
            className={cn("h-9 w-9 bg-background/50 backdrop-blur-sm hover:bg-background/80 transition-colors rounded-full z-50", className)}
            onClick={toggleLanguage}
            title={isZh ? 'Switch to English' : '切换到简体中文'}
        >
            <span className={cn("absolute text-xs font-bold transition-all", isZh ? "rotate-0 scale-100" : "-rotate-90 scale-0")}>
                中
            </span>
            <span className={cn("absolute text-xs font-bold transition-all", !isZh ? "rotate-0 scale-100" : "rotate-90 scale-0")}>
                En
            </span>
            <span className="sr-only">Toggle language</span>
        </Button>
    );
}
