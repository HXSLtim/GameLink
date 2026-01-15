import { useEffect } from 'react';
import { useOrderStore, OrderStatus } from '@/stores';
import { PageContainer } from '@/components/page-container';
import { Card, CardContent, CardFooter, CardHeader } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { ScrollArea } from '@/components/ui/scroll-area';
import { format } from 'date-fns';
import { Clock, CheckCircle, XCircle, Gamepad2, AlertCircle, Ban, Check } from 'lucide-react';
import { useTranslation } from 'react-i18next';
import { toast } from 'sonner';
import { Loader2 } from 'lucide-react';

export default function PlayerOrderListPage() {
    const { playerOrders, fetchPlayerOrders, acceptOrder, rejectOrder, loading } = useOrderStore();
    const { t } = useTranslation();

    useEffect(() => {
        fetchPlayerOrders();
    }, [fetchPlayerOrders]);

    const handleAccept = async (id: string, e: React.MouseEvent) => {
        e.stopPropagation();
        try {
            await acceptOrder(id);
            toast.success(t('order.accepted', { defaultValue: 'Order accepted' }));
        } catch {
            toast.error(t('order.accept_failed', { defaultValue: 'Failed to accept order' }));
        }
    };

    const handleReject = async (id: string, e: React.MouseEvent) => {
        e.stopPropagation();
        try {
            await rejectOrder(id);
            toast.success(t('order.rejected', { defaultValue: 'Order rejected' }));
        } catch {
            toast.error(t('order.reject_failed', { defaultValue: 'Failed to reject order' }));
        }
    };

    const getStatusConfig = (status: OrderStatus) => {
        const configs: Record<string, { label: string; color: string; icon: any }> = {
            [OrderStatus.PENDING]: { label: 'Pending', color: 'bg-yellow-500/10 text-yellow-500 border-yellow-500/20', icon: Clock },
            [OrderStatus.PAID]: { label: 'Paid/Pending', color: 'bg-blue-500/10 text-blue-500 border-blue-500/20', icon: CheckCircle },
            [OrderStatus.ACCEPTED]: { label: 'Accepted', color: 'bg-indigo-500/10 text-indigo-500 border-indigo-500/20', icon: Gamepad2 },
            [OrderStatus.COMPLETED]: { label: 'Completed', color: 'bg-green-500/10 text-green-500 border-green-500/20', icon: CheckCircle },
            [OrderStatus.CANCELLED]: { label: 'Cancelled', color: 'bg-red-500/10 text-red-500 border-red-500/20', icon: XCircle },
            [OrderStatus.REFUNDED]: { label: 'Refunded', color: 'bg-gray-500/10 text-gray-500 border-gray-500/20', icon: AlertCircle },
        };
        return configs[status] || configs[OrderStatus.PENDING];
    };

    return (
        <PageContainer>
            <div className="max-w-5xl mx-auto h-full flex flex-col space-y-6">
                <div className="flex flex-col space-y-2">
                    <h1 className="text-3xl font-bold tracking-tight bg-clip-text text-transparent bg-gradient-to-r from-white to-white/60">
                        {t('player.orders_title', { defaultValue: 'Received Orders' })}
                    </h1>
                    <p className="text-muted-foreground">
                        {t('player.orders_desc', { defaultValue: 'Manage your incoming game requests' })}
                    </p>
                </div>

                <div className="flex-1 min-h-0 bg-muted/30 rounded-xl border border-white/5 p-4">
                    {loading ? (
                        <div className="flex h-full items-center justify-center">
                            <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
                        </div>
                    ) : playerOrders.length === 0 ? (
                        <div className="flex flex-col items-center justify-center py-20 text-center">
                            <div className="p-6 rounded-full bg-muted/30 mb-4">
                                <Gamepad2 className="h-10 w-10 text-muted-foreground/50" />
                            </div>
                            <h3 className="text-lg font-medium">{t('player.no_orders', { defaultValue: 'No orders yet' })}</h3>
                            <p className="text-sm text-muted-foreground mt-1">{t('player.no_orders_desc', { defaultValue: 'Wait for players to book you!' })}</p>
                        </div>
                    ) : (
                        <ScrollArea className="h-[calc(100vh-250px)] pr-4">
                            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-2 xl:grid-cols-3 gap-4 pb-12">
                                {playerOrders.map((order) => {
                                    const config = getStatusConfig(order.status);
                                    const StatusIcon = config.icon;

                                    return (
                                        <Card key={order.id} className="group flex flex-col overflow-hidden border-white/5 bg-sidebar hover:border-primary/50 transition-all duration-300">
                                            <CardHeader className="bg-muted/10 p-4 flex flex-row items-center justify-between space-y-0 border-b border-white/5">
                                                <Badge variant="outline" className="font-mono text-[10px] text-muted-foreground border-white/10">
                                                    #{order.orderNo.slice(-8)}
                                                </Badge>
                                                <Badge variant="outline" className={`${config.color} border gap-1.5 px-2.5 py-0.5 shadow-sm font-medium`}>
                                                    <StatusIcon className="h-3 w-3" />
                                                    {config.label}
                                                </Badge>
                                            </CardHeader>
                                            <CardContent className="p-5 flex-1 space-y-4">
                                                <div className="flex justify-between items-start">
                                                    <div>
                                                        <h4 className="font-semibold">{order.gameName}</h4>
                                                        <p className="text-sm text-muted-foreground">{format(new Date(order.createdAt), 'MMM d, p')}</p>
                                                    </div>
                                                    <div className="text-right">
                                                        <div className="text-lg font-bold">¥{order.amount.toFixed(0)}</div>
                                                        <div className="text-xs text-muted-foreground">{order.quantity} hours</div>
                                                    </div>
                                                </div>

                                                {/* Mock User Details - In real app, fetch user info */}
                                                <div className="flex items-center gap-2 p-2 rounded-lg bg-white/5 text-sm">
                                                    <div className="h-8 w-8 rounded-full bg-primary/20 flex items-center justify-center text-xs">U{order.userId}</div>
                                                    <div>
                                                        <div className="font-medium">User {order.userId}</div>
                                                        <div className="text-xs text-muted-foreground">Level 5 • 4.8★</div>
                                                    </div>
                                                </div>
                                            </CardContent>
                                            <CardFooter className="bg-muted/5 p-3 px-5 flex justify-end gap-2 border-t border-white/5 mt-auto">
                                                {([OrderStatus.PENDING, OrderStatus.PAID] as OrderStatus[]).includes(order.status) ? (
                                                    <>
                                                        <Button
                                                            variant="outline"
                                                            size="sm"
                                                            className="h-8 border-red-500/20 hover:bg-red-500/10 hover:text-red-500"
                                                            onClick={(e) => handleReject(order.id, e)}
                                                        >
                                                            <Ban className="mr-1 h-3 w-3" />
                                                            {t('common.reject', { defaultValue: 'Reject' })}
                                                        </Button>
                                                        <Button
                                                            size="sm"
                                                            className="h-8 bg-green-600 hover:bg-green-700 text-white"
                                                            onClick={(e) => handleAccept(order.id, e)}
                                                        >
                                                            <Check className="mr-1 h-3 w-3" />
                                                            {t('common.accept', { defaultValue: 'Accept' })}
                                                        </Button>
                                                    </>
                                                ) : (
                                                    <Button variant="ghost" size="sm" className="h-8 w-full" disabled>
                                                        {config.label}
                                                    </Button>
                                                )}
                                            </CardFooter>
                                        </Card>
                                    );
                                })}
                            </div>
                        </ScrollArea>
                    )}
                </div>
            </div>
        </PageContainer>
    );
}
