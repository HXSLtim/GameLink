import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { cn } from '@/lib/utils';
import {
    PresenceStatus,
    getStatusColor,
    usePresenceStore,
} from '@/stores/modules/presence-store';
import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuItem,
    DropdownMenuSeparator,
    DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import { Button } from '@/components/ui/button';
import { ChevronDown, Circle, Gamepad2, Coffee, Eye, EyeOff, CheckCircle } from 'lucide-react';
import { toast } from 'sonner';

interface PresenceSelectorProps {
    className?: string;
    showLabel?: boolean;
    size?: 'sm' | 'md' | 'lg';
}

interface StatusOption {
    status: PresenceStatus;
    labelKey: string;
    descriptionKey: string;
    icon: React.ReactNode;
}

export function PresenceSelector({
    className,
    showLabel = true,
    size = 'md',
}: PresenceSelectorProps) {
    const { t } = useTranslation();
    const { myPresence, setStatus } = usePresenceStore();
    const [loading, setLoading] = useState(false);
    const [open, setOpen] = useState(false);

    const currentStatus = myPresence?.status || PresenceStatus.OFFLINE;

    // Status options with i18n keys
    const statusOptions: StatusOption[] = [
        {
            status: PresenceStatus.ONLINE,
            labelKey: 'online',
            descriptionKey: 'online_desc',
            icon: <Circle className="h-3 w-3 fill-current" />,
        },
        {
            status: PresenceStatus.ACCEPTING,
            labelKey: 'accepting',
            descriptionKey: 'accepting_desc',
            icon: <CheckCircle className="h-3 w-3" />,
        },
        {
            status: PresenceStatus.IN_GAME,
            labelKey: 'in_game',
            descriptionKey: 'in_game_desc',
            icon: <Gamepad2 className="h-3 w-3" />,
        },
        {
            status: PresenceStatus.RESTING,
            labelKey: 'resting',
            descriptionKey: 'resting_desc',
            icon: <Coffee className="h-3 w-3" />,
        },
        {
            status: PresenceStatus.INVISIBLE,
            labelKey: 'invisible',
            descriptionKey: 'invisible_desc',
            icon: <EyeOff className="h-3 w-3" />,
        },
        {
            status: PresenceStatus.OFFLINE,
            labelKey: 'offline',
            descriptionKey: 'offline_desc',
            icon: <Eye className="h-3 w-3" />,
        },
    ];

    const currentOption = statusOptions.find((opt) => opt.status === currentStatus) || statusOptions[5];
    const color = getStatusColor(currentStatus);

    const sizeClasses = {
        sm: 'h-7 text-xs px-2',
        md: 'h-9 text-sm px-3',
        lg: 'h-11 text-base px-4',
    };

    const dotSizeClasses = {
        sm: 'h-2 w-2',
        md: 'h-2.5 w-2.5',
        lg: 'h-3 w-3',
    };

    const handleStatusChange = async (status: PresenceStatus) => {
        if (status === currentStatus) {
            setOpen(false);
            return;
        }

        setLoading(true);
        try {
            await setStatus(status);
            toast.success(t('presence.status_updated', { status: t(`presence.${status}`) }));
        } catch {
            toast.error(t('presence.status_update_failed'));
        } finally {
            setLoading(false);
            setOpen(false);
        }
    };

    return (
        <DropdownMenu open={open} onOpenChange={setOpen}>
            <DropdownMenuTrigger asChild>
                <Button
                    variant="outline"
                    className={cn(
                        'gap-2 font-normal',
                        sizeClasses[size],
                        className
                    )}
                    disabled={loading}
                >
                    <span
                        className={cn('rounded-full', dotSizeClasses[size])}
                        style={{ backgroundColor: color }}
                    />
                    {showLabel && (
                        <span className="truncate">{t(`presence.${currentOption.labelKey}`)}</span>
                    )}
                    <ChevronDown className="h-3.5 w-3.5 opacity-50" />
                </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end" className="w-48">
                {statusOptions.slice(0, 4).map((option) => (
                    <DropdownMenuItem
                        key={option.status}
                        onClick={() => handleStatusChange(option.status)}
                        className={cn(
                            'gap-2 cursor-pointer',
                            currentStatus === option.status && 'bg-accent'
                        )}
                    >
                        <span
                            className="h-2.5 w-2.5 rounded-full"
                            style={{ backgroundColor: getStatusColor(option.status) }}
                        />
                        <div className="flex-1">
                            <div className="text-sm">{t(`presence.${option.labelKey}`)}</div>
                            <div className="text-xs text-muted-foreground">{t(`presence.${option.descriptionKey}`)}</div>
                        </div>
                        {currentStatus === option.status && (
                            <CheckCircle className="h-3.5 w-3.5 text-primary" />
                        )}
                    </DropdownMenuItem>
                ))}
                <DropdownMenuSeparator />
                {statusOptions.slice(4).map((option) => (
                    <DropdownMenuItem
                        key={option.status}
                        onClick={() => handleStatusChange(option.status)}
                        className={cn(
                            'gap-2 cursor-pointer',
                            currentStatus === option.status && 'bg-accent'
                        )}
                    >
                        <span
                            className="h-2.5 w-2.5 rounded-full"
                            style={{ backgroundColor: getStatusColor(option.status) }}
                        />
                        <div className="flex-1">
                            <div className="text-sm">{t(`presence.${option.labelKey}`)}</div>
                            <div className="text-xs text-muted-foreground">{t(`presence.${option.descriptionKey}`)}</div>
                        </div>
                        {currentStatus === option.status && (
                            <CheckCircle className="h-3.5 w-3.5 text-primary" />
                        )}
                    </DropdownMenuItem>
                ))}
            </DropdownMenuContent>
        </DropdownMenu>
    );
}

// Compact version for mobile/sidebar
interface PresenceSelectorCompactProps {
    className?: string;
}

export function PresenceSelectorCompact({ className }: PresenceSelectorCompactProps) {
    const { t } = useTranslation();
    const { myPresence, setStatus } = usePresenceStore();
    const [loading, setLoading] = useState(false);

    const currentStatus = myPresence?.status || PresenceStatus.OFFLINE;

    const handleQuickToggle = async () => {
        if (loading) return;

        // Quick toggle between online and offline
        const newStatus = currentStatus === PresenceStatus.OFFLINE
            ? PresenceStatus.ONLINE
            : PresenceStatus.OFFLINE;

        setLoading(true);
        try {
            await setStatus(newStatus);
        } catch {
            toast.error(t('presence.status_update_failed'));
        } finally {
            setLoading(false);
        }
    };

    return (
        <button
            onClick={handleQuickToggle}
            disabled={loading}
            className={cn(
                'relative h-4 w-4 rounded-full transition-all',
                loading && 'opacity-50',
                className
            )}
            style={{ backgroundColor: getStatusColor(currentStatus) }}
            title={t('presence.current_status', { status: t(`presence.${currentStatus}`) })}
        >
            <span className="sr-only">{t('presence.toggle_status')}</span>
        </button>
    );
}
