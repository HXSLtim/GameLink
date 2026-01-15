import { useOnlineStatus } from '@/hooks/use-online-status';
import { WifiOff } from 'lucide-react';

/**
 * Offline banner component - shows when user loses internet connection
 */
export function OfflineBanner() {
    const isOnline = useOnlineStatus();

    if (isOnline) return null;

    return (
        <div className="fixed top-0 left-0 right-0 z-50 bg-destructive text-destructive-foreground py-2 px-4 text-center text-sm font-medium flex items-center justify-center gap-2">
            <WifiOff className="w-4 h-4" />
            <span>网络连接已断开，部分功能可能不可用</span>
        </div>
    );
}
