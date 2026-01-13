import { useTranslation } from 'react-i18next';
import { Globe } from 'lucide-react';
import { Button } from '@/components/ui/button';
import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuItem,
    DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import { cn } from '@/lib/utils';

interface LanguageSwitcherProps {
    className?: string;
}

export function LanguageSwitcher({ className }: LanguageSwitcherProps) {
    const { i18n } = useTranslation();

    const changeLanguage = (lng: string) => {
        i18n.changeLanguage(lng);
    };

    const currentLanguage = i18n.language;

    return (
        <div className={cn("z-50", className)}>
            <DropdownMenu>
                <DropdownMenuTrigger asChild>
                    <Button variant="ghost" size="icon" className="h-9 w-9 bg-background/50 backdrop-blur-sm hover:bg-background/80 transition-colors rounded-full">
                        <Globe className="h-[1.2rem] w-[1.2rem]" />
                        <span className="sr-only">Toggle language</span>
                    </Button>
                </DropdownMenuTrigger>
                <DropdownMenuContent align="end">
                    <DropdownMenuItem
                        onClick={() => changeLanguage('zh-CN')}
                        className={cn(currentLanguage === 'zh-CN' && "bg-accent text-accent-foreground")}
                    >
                        简体中文
                    </DropdownMenuItem>
                    <DropdownMenuItem
                        onClick={() => changeLanguage('en-US')}
                        className={cn(currentLanguage === 'en-US' && "bg-accent text-accent-foreground")}
                    >
                        English
                    </DropdownMenuItem>
                </DropdownMenuContent>
            </DropdownMenu>
        </div>
    );
}
