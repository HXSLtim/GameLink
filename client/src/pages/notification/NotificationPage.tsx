import { useEffect } from 'react';
import { useNotificationStore } from '@/stores';
import { PageContainer } from '@/components/page-container';
import { Card, CardContent } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { ScrollArea } from '@/components/ui/scroll-area';
import { useTranslation } from 'react-i18next';
import { Bell, CheckCheck, Info, Tag, MessageSquare, ArrowLeft } from 'lucide-react';
import { format } from 'date-fns';
import { useNavigate } from 'react-router-dom';
import { cn } from '@/lib/utils';

export default function NotificationPage() {
    const { t } = useTranslation();
    const navigate = useNavigate();
    const { notifications, unreadCount, fetchNotifications, markAsRead, markAllAsRead } = useNotificationStore();

    useEffect(() => {
        fetchNotifications();
    }, [fetchNotifications]);

    const getIcon = (type: string) => {
        switch (type) {
            case 'order': return <MessageSquare className="h-5 w-5 text-blue-500" />;
            case 'promotion': return <Tag className="h-5 w-5 text-orange-500" />;
            default: return <Info className="h-5 w-5 text-primary" />;
        }
    };

    return (
        <PageContainer>
            <div className="max-w-4xl mx-auto h-full flex flex-col">
                {/* Header */}
                <div className="flex items-center justify-between p-4 bg-background/80 backdrop-blur-md sticky top-0 z-10 border-b">
                    <div className="flex items-center gap-3">
                        <Button variant="ghost" size="icon" onClick={() => navigate(-1)} className="rounded-full">
                            <ArrowLeft className="h-5 w-5" />
                        </Button>
                        <div>
                            <h1 className="text-xl font-bold flex items-center gap-2">
                                {t('notifications.title', { defaultValue: 'Notifications' })}
                                {unreadCount > 0 && (
                                    <Badge className="bg-primary text-primary-foreground h-5 min-w-[20px] px-1.5 rounded-full flex items-center justify-center text-xs">
                                        {unreadCount}
                                    </Badge>
                                )}
                            </h1>
                        </div>
                    </div>
                    {notifications.length > 0 && (
                        <Button variant="ghost" size="sm" onClick={() => markAllAsRead()} className="text-muted-foreground hover:text-primary">
                            <CheckCheck className="h-4 w-4 mr-2" />
                            Mark all read
                        </Button>
                    )}
                </div>

                {/* List */}
                <ScrollArea className="flex-1 px-4 py-4">
                    {notifications.length === 0 ? (
                        <div className="flex flex-col items-center justify-center py-20 text-muted-foreground space-y-4">
                            <div className="p-6 bg-muted/50 rounded-full">
                                <Bell className="h-12 w-12 opacity-20" />
                            </div>
                            <p>No notifications yet</p>
                        </div>
                    ) : (
                        <div className="space-y-3 pb-20">
                            {notifications.map((notification) => (
                                <Card
                                    key={notification.id}
                                    onClick={() => markAsRead(notification.id)}
                                    className={cn(
                                        "cursor-pointer transition-all hover:bg-muted/40 border-l-4",
                                        notification.read ? "border-l-transparent opacity-70 bg-muted/20" : "border-l-primary bg-background shadow-sm"
                                    )}
                                >
                                    <CardContent className="p-4 flex gap-4">
                                        <div className={cn(
                                            "mt-1 p-2 rounded-full h-fit",
                                            notification.read ? "bg-muted" : "bg-primary/10"
                                        )}>
                                            {getIcon(notification.type)}
                                        </div>
                                        <div className="flex-1 space-y-1">
                                            <div className="flex justify-between items-start">
                                                <h4 className={cn("font-medium text-sm", !notification.read && "text-primary font-bold")}>
                                                    {notification.title}
                                                </h4>
                                                <span className="text-[10px] text-muted-foreground shrink-0 ml-2">
                                                    {format(new Date(notification.createdAt), 'MMM d, HH:mm')}
                                                </span>
                                            </div>
                                            <p className="text-sm text-muted-foreground leading-snug">
                                                {notification.message}
                                            </p>
                                        </div>
                                    </CardContent>
                                </Card>
                            ))}
                        </div>
                    )}
                </ScrollArea>
            </div>
        </PageContainer>
    );
}
