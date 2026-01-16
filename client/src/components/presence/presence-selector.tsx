import { useState } from 'react';
import { cn } from '@/lib/utils';
import {
    PresenceStatus,
    getStatusColor,
    getStatusDisplay,
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
    label: string;
    icon: React.ReactNode;
    color: string;
    description?: string;
}

const statusOptions: StatusOption[] = [
    {
        status: PresenceStatus.ONLINE,
        label: '在线',
        icon: <Circle className="h-3 w-3 fill-current" />,
        color: getStatusColor(PresenceStatus.ONLINE),
        description: '显示为在线状态',
    },
    {
        status: PresenceStatus.ACCEPTING,
        label: '接单中',
        icon: <CheckCircle className="h-3 w-3" />,
        color: getStatusColor(PresenceStatus.ACCEPTING),
        description: '正在接受订单',
    },
    {
        status: PresenceStatus.IN_GAME,
        label: '游戏中',
        icon: <Gamepad2 className="h-3 w-3" />,
        color: getStatusColor(PresenceStatus.IN_GAME),
        description: '正在进行游戏',
    },
    {
        status: PresenceStatus.RESTING,
        label: '休息中',
        icon: <Coffee className="h-3 w-3" />,
        color: getStatusColor(PresenceStatus.RESTING),
        description: '暂时休息，不接单',
    },
    {
        status: PresenceStatus.INVISIBLE,
        label: '隐身',
        icon: <EyeOff className="h-3 w-3" />,
        color: getStatusColor(PresenceStatus.INVISIBLE),
        description: '对他人显示为离线',
    },
    {
        status: PresenceStatus.OFFLINE,
        label: '离线',
        icon: <Eye className="h-3 w-3" />,
        color: getStatusColor(PresenceStatus.OFFLINE),
        description: '设为离线状态',
    },
];

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

export function PresenceSelector({
    className,
    showLabel = true,
    size = 'md',
}: PresenceSelectorProps) {
    const { myPresence, setStatus } = usePresenceStore();
    const [loading, setLoading] = useState(false);
    const [open, setOpen] = useState(false);

    const currentStatus = myPresence?.status || PresenceStatus.OFFLINE;
    const currentOption = statusOptions.find((opt) => opt.status === currentStatus) || statusOptions[5];

    const handleStatusChange = async (status: PresenceStatus) => {
        if (status === currentStatus) {
            setOpen(false);
            return;
        }

        setLoading(true);
        try {
            await setStatus(status);
            toast.success(`状态已更新为${getStatusDisplay(status)}`);
        } catch {
            toast.error('状态更新失败');
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
                        style={{ backgroundColor: currentOption.color }}
                    />
                    {showLabel && (
                        <span className="truncate">{currentOption.label}</span>
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
                            style={{ backgroundColor: option.color }}
                        />
                        <span className="flex-1">{option.label}</span>
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
                            style={{ backgroundColor: option.color }}
                        />
                        <span className="flex-1">{option.label}</span>
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
            toast.error('状态更新失败');
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
            title={`当前状态: ${getStatusDisplay(currentStatus)}`}
        >
            <span className="sr-only">切换在线状态</span>
        </button>
    );
}
