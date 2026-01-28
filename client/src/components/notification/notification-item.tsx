import { Button } from '@/components/ui/button';
import {
    Bell,
    ShoppingBag,
    MessageSquare,
    Wallet,
    Gift,
    CheckCircle,
    Clock,
    ChevronRight,
} from 'lucide-react';
import { cn } from '@/lib/utils';
import { formatDistanceToNow, parseISO } from 'date-fns';
import { zhCN } from 'date-fns/locale';
import type { Notification } from '@/types';

interface NotificationItemProps {
    notification: Notification;
    /** 点击通知回调 */
    onClick?: () => void;
    /** 标记已读回调 */
    onMarkRead?: () => void;
    /** 变体 */
    variant?: 'default' | 'compact';
    /** 自定义类名 */
    className?: string;
}

// 通知类型配置
const NOTIFICATION_CONFIG: Record<
    Notification['type'],
    {
        icon: typeof Bell;
        color: string;
        bgColor: string;
    }
> = {
    system: {
        icon: Bell,
        color: 'text-blue-500',
        bgColor: 'bg-blue-50',
    },
    order: {
        icon: ShoppingBag,
        color: 'text-green-500',
        bgColor: 'bg-green-50',
    },
    chat: {
        icon: MessageSquare,
        color: 'text-purple-500',
        bgColor: 'bg-purple-50',
    },
    wallet: {
        icon: Wallet,
        color: 'text-amber-500',
        bgColor: 'bg-amber-50',
    },
    promotion: {
        icon: Gift,
        color: 'text-pink-500',
        bgColor: 'bg-pink-50',
    },
};

export function NotificationItem({
    notification,
    onClick,
    onMarkRead,
    variant = 'default',
    className,
}: NotificationItemProps) {
    const config = NOTIFICATION_CONFIG[notification.type] || NOTIFICATION_CONFIG.system;
    const Icon = config.icon;

    const formatTime = (dateString: string) => {
        try {
            return formatDistanceToNow(parseISO(dateString), {
                addSuffix: true,
                locale: zhCN,
            });
        } catch {
            return dateString;
        }
    };

    const handleClick = () => {
        if (!notification.isRead && onMarkRead) {
            onMarkRead();
        }
        onClick?.();
    };

    if (variant === 'compact') {
        return (
            <div
                className={cn(
                    'flex items-center gap-3 p-3 rounded-lg cursor-pointer transition-colors',
                    notification.isRead
                        ? 'bg-background hover:bg-muted/50'
                        : 'bg-primary/5 hover:bg-primary/10',
                    className
                )}
                onClick={handleClick}
            >
                <div
                    className={cn(
                        'flex-shrink-0 w-8 h-8 rounded-full flex items-center justify-center',
                        config.bgColor
                    )}
                >
                    <Icon className={cn('h-4 w-4', config.color)} />
                </div>
                <div className="flex-1 min-w-0">
                    <div className="font-medium truncate text-sm">{notification.title}</div>
                    <div className="text-xs text-muted-foreground truncate">
                        {formatTime(notification.createdAt)}
                    </div>
                </div>
                {!notification.isRead && (
                    <div className="w-2 h-2 rounded-full bg-primary flex-shrink-0" />
                )}
            </div>
        );
    }

    return (
        <div
            className={cn(
                'relative p-4 rounded-lg border cursor-pointer transition-all',
                notification.isRead
                    ? 'bg-background hover:bg-muted/30 border-border'
                    : 'bg-primary/5 hover:bg-primary/10 border-primary/20',
                className
            )}
            onClick={handleClick}
        >
            {/* 未读指示器 */}
            {!notification.isRead && (
                <div className="absolute top-4 right-4 w-2 h-2 rounded-full bg-primary" />
            )}

            <div className="flex gap-4">
                {/* 图标 */}
                <div
                    className={cn(
                        'flex-shrink-0 w-10 h-10 rounded-full flex items-center justify-center',
                        config.bgColor
                    )}
                >
                    <Icon className={cn('h-5 w-5', config.color)} />
                </div>

                {/* 内容 */}
                <div className="flex-1 min-w-0">
                    <div className="flex items-start justify-between gap-2">
                        <h4 className={cn(
                            'font-medium',
                            !notification.isRead && 'text-foreground'
                        )}>
                            {notification.title}
                        </h4>
                    </div>

                    <p className="text-sm text-muted-foreground mt-1 line-clamp-2">
                        {notification.content}
                    </p>

                    <div className="flex items-center justify-between mt-3">
                        <div className="flex items-center gap-1 text-xs text-muted-foreground">
                            <Clock className="h-3 w-3" />
                            <span>{formatTime(notification.createdAt)}</span>
                        </div>

                        <div className="flex items-center gap-2">
                            {!notification.isRead && onMarkRead && (
                                <Button
                                    variant="ghost"
                                    size="sm"
                                    className="h-7 text-xs"
                                    onClick={(e) => {
                                        e.stopPropagation();
                                        onMarkRead();
                                    }}
                                >
                                    <CheckCircle className="h-3 w-3 mr-1" />
                                    标记已读
                                </Button>
                            )}
                            {onClick && (
                                <Button variant="ghost" size="sm" className="h-7 text-xs">
                                    查看
                                    <ChevronRight className="h-3 w-3 ml-1" />
                                </Button>
                            )}
                        </div>
                    </div>
                </div>
            </div>
        </div>
    );
}

/**
 * 通知列表空状态
 */
export function NotificationEmpty() {
    return (
        <div className="flex flex-col items-center justify-center py-12 text-center">
            <div className="w-16 h-16 rounded-full bg-muted flex items-center justify-center mb-4">
                <Bell className="h-8 w-8 text-muted-foreground" />
            </div>
            <h3 className="font-medium text-lg mb-1">暂无通知</h3>
            <p className="text-sm text-muted-foreground">新的通知将会显示在这里</p>
        </div>
    );
}
