import { useEffect, useState } from 'react';
import { useOrderStore, OrderStatus } from '@/stores';
import { PageContainer } from '@/components/page-container';
import { Tabs, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Card, CardContent, CardFooter, CardHeader } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { ScrollArea } from '@/components/ui/scroll-area';
import { format } from 'date-fns';
import { Package, Clock, CheckCircle, XCircle, AlertCircle, RefreshCcw, Gamepad2, ChevronRight } from 'lucide-react';
import { useNavigate } from 'react-router-dom';
import { useTranslation } from 'react-i18next';

export default function OrderListPage() {
    const { myOrders, fetchOrders, loading } = useOrderStore();
    const navigate = useNavigate();
    const { t } = useTranslation();
    const [activeTab, setActiveTab] = useState("all");

    useEffect(() => {
        fetchOrders();
    }, []);

    const getStatusConfig = (status: OrderStatus) => {
        const configs: Record<string, { labelKey: string; color: string; icon: any }> = {
            [OrderStatus.PENDING]: { labelKey: 'nav.order_status.pending', color: 'bg-yellow-500/10 text-yellow-500 border-yellow-500/20', icon: Clock },
            [OrderStatus.PAID]: { labelKey: 'nav.order_status.pending', color: 'bg-blue-500/10 text-blue-500 border-blue-500/20', icon: CheckCircle }, // Treating Paid as Pending/Processing
            [OrderStatus.ACCEPTED]: { labelKey: 'nav.order_status.accepted', color: 'bg-indigo-500/10 text-indigo-500 border-indigo-500/20', icon: Gamepad2 },
            [OrderStatus.COMPLETED]: { labelKey: 'nav.order_status.completed', color: 'bg-green-500/10 text-green-500 border-green-500/20', icon: CheckCircle },
            [OrderStatus.CANCELLED]: { labelKey: 'nav.order_status.cancelled', color: 'bg-red-500/10 text-red-500 border-red-500/20', icon: XCircle },
            [OrderStatus.REFUNDED]: { labelKey: 'nav.order_status.cancelled', color: 'bg-gray-500/10 text-gray-500 border-gray-500/20', icon: AlertCircle },
        };
        return configs[status] || configs[OrderStatus.PENDING];
    };

    const filteredOrders = activeTab === "all"
        ? myOrders
        : myOrders.filter(order => {
            if (activeTab === 'active') return ([OrderStatus.PENDING, OrderStatus.PAID, OrderStatus.ACCEPTED] as OrderStatus[]).includes(order.status);
            if (activeTab === 'completed') return order.status === OrderStatus.COMPLETED;
            if (activeTab === 'cancelled') return ([OrderStatus.CANCELLED, OrderStatus.REFUNDED] as OrderStatus[]).includes(order.status);
            return true;
        });

    return (
        <PageContainer>
            <div className="max-w-6xl mx-auto h-full flex flex-col space-y-6">
                <div className="flex flex-col md:flex-row md:items-center justify-between gap-4 animate-in fade-in slide-in-from-top-2">
                    <div className="space-y-1">
                        <h1 className="text-3xl font-bold tracking-tight bg-clip-text text-transparent bg-gradient-to-r from-white to-white/60">
                            {t('nav.orders_title')}
                        </h1>
                        <p className="text-muted-foreground">
                            {t('nav.orders_desc')}
                        </p>
                    </div>
                    <Button variant="outline" size="sm" onClick={() => fetchOrders()} className="self-start md:self-auto gap-2 backdrop-blur-sm bg-background/50">
                        <RefreshCcw className={`h-3.5 w-3.5 ${loading ? 'animate-spin' : ''}`} />
                        Refresh
                    </Button>
                </div>

                <Tabs defaultValue="all" value={activeTab} onValueChange={setActiveTab} className="flex-1 flex flex-col">
                    <TabsList className="grid w-full grid-cols-4 md:w-[400px] p-1 bg-muted/40 backdrop-blur-md border border-white/5 rounded-full">
                        <TabsTrigger value="all" className="rounded-full data-[state=active]:bg-background/80 data-[state=active]:text-foreground data-[state=active]:shadow-sm transition-all">{t('nav.filter.all')}</TabsTrigger>
                        <TabsTrigger value="active" className="rounded-full data-[state=active]:bg-background/80 transition-all">Active</TabsTrigger>
                        <TabsTrigger value="completed" className="rounded-full data-[state=active]:bg-background/80 transition-all">Done</TabsTrigger>
                        <TabsTrigger value="cancelled" className="rounded-full data-[state=active]:bg-background/80 transition-all">Cancel</TabsTrigger>
                    </TabsList>

                    <div className="mt-6 flex-1 min-h-0">
                        {filteredOrders.length === 0 && !loading ? (
                            <div className="flex flex-col items-center justify-center py-20 text-center animate-in fade-in zoom-in-95">
                                <div className="p-6 rounded-full bg-muted/30 mb-4">
                                    <Package className="h-10 w-10 text-muted-foreground/50" />
                                </div>
                                <h3 className="text-lg font-medium">{t('nav.no_players')}</h3>
                                {/* Reusing 'no_players' or generic 'no_data' - ideally should sort out 'no_orders' key, but 'no_players' is okay for now or fallback */}
                                <p className="text-sm text-muted-foreground mt-1 max-w-[250px] mx-auto">You haven't placed any orders in this category yet.</p>
                                <Button className="mt-6 rounded-full px-8" onClick={() => navigate('/players')}>{t('nav.find_players_title')}</Button>
                            </div>
                        ) : (
                            <ScrollArea className="h-[calc(100vh-280px)] pr-4">
                                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 pb-12">
                                    {filteredOrders.map((order, i) => {
                                        const config = getStatusConfig(order.status);
                                        const StatusIcon = config.icon;

                                        return (
                                            <Card
                                                key={order.id}
                                                className="group flex flex-col overflow-hidden border-white/5 bg-background/40 backdrop-blur-md hover:border-primary/50 transition-all duration-300"
                                                style={{ animationDelay: `${i * 50}ms` }}
                                            >
                                                <CardHeader className="bg-muted/10 p-4 flex flex-row items-center justify-between space-y-0 border-b border-white/5">
                                                    <div className="flex items-center gap-2">
                                                        <Badge variant="outline" className="font-mono text-[10px] text-muted-foreground border-white/10">
                                                            #{order.orderNo.slice(-8)}
                                                        </Badge>
                                                        <span className="text-xs text-muted-foreground">
                                                            {format(new Date(order.createdAt), 'MMM d')}
                                                        </span>
                                                    </div>
                                                    <Badge variant="outline" className={`${config.color} border gap-1.5 px-2.5 py-0.5 shadow-sm font-medium`}>
                                                        <StatusIcon className="h-3 w-3" />
                                                        {t(config.labelKey)}
                                                    </Badge>
                                                </CardHeader>
                                                <CardContent className="p-5 flex-1">
                                                    <div className="flex items-start gap-4">
                                                        <div className="h-12 w-12 rounded-xl bg-gradient-to-br from-primary/20 to-primary/5 flex items-center justify-center shrink-0 border border-white/5 group-hover:scale-105 transition-transform">
                                                            <Gamepad2 className="h-6 w-6 text-primary" />
                                                        </div>
                                                        <div className="flex-1 min-w-0 space-y-2">
                                                            <div className="font-semibold text-lg truncate flex items-center gap-2">
                                                                {order.gameName}
                                                            </div>
                                                            <Badge variant="secondary" className="text-[10px] font-normal px-1.5 py-0 h-5">Boosting</Badge>

                                                            <div className="grid grid-cols-2 gap-2 text-sm mt-2">
                                                                <div className="text-muted-foreground">Qty: <span className="text-foreground">{order.quantity}h</span></div>
                                                                <div className="text-muted-foreground text-right">Total: <span className="font-bold text-foreground">¥{order.amount.toFixed(0)}</span></div>
                                                            </div>

                                                            {order.scheduledTime && (
                                                                <div className="flex items-center gap-1.5 text-xs text-blue-400 bg-blue-500/10 px-2 py-1 rounded w-fit mt-1">
                                                                    <Clock className="h-3 w-3" />
                                                                    {format(new Date(order.scheduledTime), 'p')}
                                                                </div>
                                                            )}
                                                        </div>
                                                    </div>
                                                </CardContent>
                                                <CardFooter className="bg-muted/5 p-3 px-5 flex justify-between gap-3 border-t border-white/5 mt-auto">
                                                    <Button variant="ghost" size="sm" className="h-8 hover:bg-white/5 text-xs" onClick={() => navigate(`/orders/${order.id}`)}>
                                                        Details <ChevronRight className="ml-1 h-3 w-3" />
                                                    </Button>
                                                    {([OrderStatus.PENDING] as OrderStatus[]).includes(order.status) && (
                                                        <Button size="sm" className="h-8 rounded-full bg-primary/90 hover:bg-primary shadow-lg shadow-primary/20 text-xs">
                                                            Pay Now
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
                </Tabs>
            </div>
        </PageContainer>
    );
}
