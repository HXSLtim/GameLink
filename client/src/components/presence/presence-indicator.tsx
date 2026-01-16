import { cn } from '@/lib/utils';
import { useTranslation } from 'react-i18next';
import {
    PresenceStatus,
    getStatusColor,
    isOnlineStatus,
} from '@/stores/modules/presence-store';

interface PresenceIndicatorProps {
    status: PresenceStatus;
    currentGameName?: string;
    customStatus?: string;
    size?: 'sm' | 'md' | 'lg';
    showLabel?: boolean;
    showGameName?: boolean;
    className?: string;
}

const sizeClasses = {
    sm: 'h-2 w-2',
    md: 'h-3 w-3',
    lg: 'h-4 w-4',
};

const labelSizeClasses = {
    sm: 'text-xs',
    md: 'text-sm',
    lg: 'text-base',
};

export function PresenceIndicator({
    status,
    currentGameName,
    customStatus,
    size = 'md',
    showLabel = false,
    showGameName = true,
    className,
}: PresenceIndicatorProps) {
    const { t } = useTranslation();
    const color = getStatusColor(status);
    const isOnline = isOnlineStatus(status);

    // Build display text
    let displayText: string;
    if (status === 'in_game' && currentGameName && showGameName) {
        displayText = t('presence.playing', { game: currentGameName });
    } else if (customStatus) {
        displayText = customStatus;
    } else {
        displayText = t(`presence.${status}`, { defaultValue: status });
    }

    return (
        <div className={cn('flex items-center gap-1.5', className)}>
            {/* Status dot */}
            <span
                className={cn(
                    'rounded-full flex-shrink-0',
                    sizeClasses[size],
                    isOnline && 'animate-pulse'
                )}
                style={{ backgroundColor: color }}
            />

            {/* Label */}
            {showLabel && (
                <span
                    className={cn(
                        'text-muted-foreground truncate',
                        labelSizeClasses[size]
                    )}
                >
                    {displayText}
                </span>
            )}
        </div>
    );
}

// Overlay version for avatars
interface PresenceIndicatorOverlayProps {
    status: PresenceStatus;
    size?: 'sm' | 'md' | 'lg';
    position?: 'bottom-right' | 'bottom-left' | 'top-right' | 'top-left';
    className?: string;
}

const positionClasses = {
    'bottom-right': 'bottom-0 right-0',
    'bottom-left': 'bottom-0 left-0',
    'top-right': 'top-0 right-0',
    'top-left': 'top-0 left-0',
};

const overlaySizeClasses = {
    sm: 'h-2.5 w-2.5 border',
    md: 'h-3.5 w-3.5 border-2',
    lg: 'h-4.5 w-4.5 border-2',
};

export function PresenceIndicatorOverlay({
    status,
    size = 'md',
    position = 'bottom-right',
    className,
}: PresenceIndicatorOverlayProps) {
    const color = getStatusColor(status);
    const isOnline = isOnlineStatus(status);

    return (
        <span
            className={cn(
                'absolute rounded-full border-background',
                overlaySizeClasses[size],
                positionClasses[position],
                isOnline && 'animate-pulse',
                className
            )}
            style={{ backgroundColor: color }}
        />
    );
}

// Badge version for lists
interface PresenceBadgeProps {
    status: PresenceStatus;
    currentGameName?: string;
    className?: string;
}

export function PresenceBadge({
    status,
    currentGameName,
    className,
}: PresenceBadgeProps) {
    const { t } = useTranslation();
    const color = getStatusColor(status);

    let displayText: string;
    if (status === 'in_game' && currentGameName) {
        displayText = t('presence.playing', { game: currentGameName });
    } else {
        displayText = t(`presence.${status}`, { defaultValue: status });
    }

    return (
        <span
            className={cn(
                'inline-flex items-center gap-1.5 px-2 py-0.5 rounded-full text-xs font-medium',
                className
            )}
            style={{
                backgroundColor: `${color}20`,
                color: color,
            }}
        >
            <span
                className="h-1.5 w-1.5 rounded-full"
                style={{ backgroundColor: color }}
            />
            {displayText}
        </span>
    );
}
