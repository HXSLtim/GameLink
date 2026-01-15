import { useNavigate } from 'react-router-dom';
import { Button } from '@/components/ui/button';
import { useTranslation } from 'react-i18next';
import { Ghost, Home } from 'lucide-react';

export default function NotFoundPage() {
    const navigate = useNavigate();
    const { t } = useTranslation();

    return (
        <div className="flex flex-col items-center justify-center min-h-screen bg-background text-foreground p-4 text-center">
            <div className="relative mb-8">
                <div className="absolute inset-0 bg-primary/20 blur-3xl rounded-full" />
                <Ghost className="h-32 w-32 text-primary relative z-10 animate-bounce" />
            </div>

            <h1 className="text-9xl font-black text-primary/20 select-none">404</h1>
            <h2 className="text-2xl font-bold mt-4 mb-2">{t('error.not_found', { defaultValue: 'Page Not Found' })}</h2>
            <p className="text-muted-foreground max-w-md mb-8">
                {t('error.not_found_desc', { defaultValue: "Oops! The page you are looking for seems to have wandered off into the void." })}
            </p>

            <Button size="lg" onClick={() => navigate('/')} className="font-bold gap-2">
                <Home className="h-4 w-4" />
                {t('common.back_home', { defaultValue: 'Back to Home' })}
            </Button>
        </div>
    );
}
