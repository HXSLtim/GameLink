import { useEffect, useState } from 'react';
import { Button } from '@/components/ui/button';
import { X, Download } from 'lucide-react';
import { useTranslation } from 'react-i18next';
import { toast } from 'sonner';

interface BeforeInstallPromptEvent extends Event {
    prompt: () => Promise<void>;
    userChoice: Promise<{ outcome: 'accepted' | 'dismissed'; platform: string }>;
}

export function InstallPrompt() {
    const { t } = useTranslation();
    const [deferredPrompt, setDeferredPrompt] = useState<BeforeInstallPromptEvent | null>(null);
    const [isVisible, setIsVisible] = useState(false);

    useEffect(() => {
        const handler = (e: Event) => {
            // Prevent the mini-infobar from appearing on mobile
            e.preventDefault();
            // Stash the event so it can be triggered later.
            setDeferredPrompt(e as BeforeInstallPromptEvent);
            // Update UI notify the user they can install the PWA
            setIsVisible(true);
        };

        window.addEventListener('beforeinstallprompt', handler);

        return () => {
            window.removeEventListener('beforeinstallprompt', handler);
        };
    }, []);

    const handleInstallClick = async () => {
        setIsVisible(false);

        if (!deferredPrompt) {
            return;
        }

        // Show the install prompt
        await deferredPrompt.prompt();

        // Wait for the user to respond to the prompt
        const { outcome } = await deferredPrompt.userChoice;

        if (outcome === 'accepted') {
            toast.success(t('pwa.install_success', { defaultValue: 'App installed successfully!' }));
        } else {
            // Dismissed
        }

        setDeferredPrompt(null);
    };

    if (!isVisible) {
        return null;
    }

    return (
        <div className="fixed bottom-20 left-4 right-4 z-50 md:bottom-4 md:left-auto md:right-4 md:w-96 animate-in slide-in-from-bottom-5 fade-in duration-500">
            <div className="flex items-center justify-between gap-4 rounded-xl border bg-background/95 p-4 shadow-lg backdrop-blur supports-[backdrop-filter]:bg-background/60">
                <div className="space-y-1">
                    <h4 className="text-sm font-semibold">{t('pwa.install_title', { defaultValue: 'Install App' })}</h4>
                    <p className="text-xs text-muted-foreground">
                        {t('pwa.install_desc', { defaultValue: 'Add to home screen for better experience.' })}
                    </p>
                </div>
                <div className="flex items-center gap-2">
                    <Button size="sm" onClick={handleInstallClick} className="whitespace-nowrap">
                        <Download className="mr-2 h-3 w-3" />
                        {t('common.install', { defaultValue: 'Install' })}
                    </Button>
                    <Button size="icon" variant="ghost" className="h-8 w-8" onClick={() => setIsVisible(false)}>
                        <X className="h-4 w-4" />
                        <span className="sr-only">Close</span>
                    </Button>
                </div>
            </div>
        </div>
    );
}
