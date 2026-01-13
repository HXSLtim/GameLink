import { useEffect, useState } from 'react';
import { useOrderStore, OrderStatus } from '@/stores';
import { PageContainer } from '@/components/page-container';
import { Tabs, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Card, CardContent, CardFooter, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { ScrollArea } from '@/components/ui/scroll-area';
import { format } from 'date-fns';
import { Package, Clock, CheckCircle, XCircle, AlertCircle, RefreshCcw } from 'lucide-react';
import { useNavigate } from 'react-router-dom';

const statusMap: Record<string, { label: string; color: string; icon: any }> = {
    [OrderStatus.PENDING]: { label: 'Pending', color: 'bg-yellow-500/10 text-yellow-500', icon: Clock },
    [OrderStatus.PAID]: { label: 'Paid', color: 'bg-blue-500/10 text-blue-500', icon: CheckCircle },
    [OrderStatus.ACCEPTED]: { label: 'Accepted', color: 'bg-indigo-500/10 text-indigo-500', icon: CheckCircle },
    [OrderStatus.COMPLETED]: { label: 'Completed', color: 'bg-green-500/10 text-green-500', icon: CheckCircle },
    [OrderStatus.CANCELLED]: { label: 'Cancelled', color: 'bg-red-500/10 text-red-500', icon: XCircle },
    [OrderStatus.REFUNDED]: { label: 'Refunded', color: 'bg-gray-500/10 text-gray-500', icon: AlertCircle },
};

export default function OrderListPage() {
    const { myOrders, fetchOrders, loading } = useOrderStore();
    const navigate = useNavigate();
    const [activeTab, setActiveTab] = useState("all");

    useEffect(() => {
        fetchOrders();
    }, []);

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
            <div className="max-w-4xl mx-auto py-8 px-4 h-full flex flex-col">
                <div className="flex items-center justify-between mb-6">
                    <div>
                        <h1 className="text-3xl font-bold tracking-tight">My Orders</h1>
                        <p className="text-muted-foreground mt-1">Manage and track your service orders</p>
                    </div>
                    <Button variant="outline" size="icon" onClick={() => fetchOrders()}>
                        <RefreshCcw className={`h-4 w-4 ${loading ? 'animate-spin' : ''}`} />
                    </Button>
                </div>

                <Tabs defaultValue="all" className="flex-1 flex flex-col" onValueChange={setActiveTab}>
                    <TabsList className="grid w-full grid-cols-4 lg:w-[400px]">
                        <TabsTrigger value="all">All</TabsTrigger>
                        <TabsTrigger value="active">Active</TabsTrigger>
                        <TabsTrigger value="completed">Done</TabsTrigger>
                        <TabsTrigger value="cancelled">Cancel</TabsTrigger>
                    </TabsList>

                    <div className="mt-4 flex-1 min-h-0">
                        {filteredOrders.length === 0 && !loading ? (
                            <div className="flex flex-col items-center justify-center h-[300px] border border-dashed rounded-lg text-center p-8">
                                <Package className="h-10 w-10 text-muted-foreground mb-4" />
                                <h3 className="text-lg font-medium">No orders found</h3>
                                <p className="text-sm text-muted-foreground mt-1">You haven't placed any orders in this category yet.</p>
                                <Button className="mt-4" onClick={() => navigate('/players')}>Find Players</Button>
                            </div>
                        ) : (
                            <ScrollArea className="h-full">
                                <div className="space-y-4 pb-4">
                                    {filteredOrders.map((order) => {
                                        const StatusIcon = statusMap[order.status]?.icon || AlertCircle;
                                        return (
                                            <Card key={order.id} className="overflow-hidden cursor-pointer hover:border-primary/50 transition-colors">
                                                <CardHeader className="bg-muted/30 p-4 flex flex-row items-center justify-between space-y-0">
                                                    <div className="space-y-1">
                                                        <CardTitle className="text-sm font-medium">Order #{order.orderNo}</CardTitle>
                                                        <div className="text-xs text-muted-foreground">
                                                            {format(new Date(order.createdAt), 'PPP p')}
                                                        </div>
                                                    </div>
                                                    <Badge variant="outline" className={`${statusMap[order.status]?.color} border-0 font-medium flex items-center gap-1`}>
                                                        <StatusIcon className="h-3 w-3" />
                                                        {statusMap[order.status]?.label}
                                                    </Badge>
                                                </CardHeader>
                                                <CardContent className="p-4">
                                                    <div className="flex items-center gap-4">
                                                        <div className="h-16 w-16 bg-muted rounded-md flex items-center justify-center shrink-0">
                                                            {/* Placeholder for game image */}
                                                            <Package className="h-8 w-8 text-muted-foreground/50" />
                                                        </div>
                                                        <div className="flex-1 min-w-0">
                                                            <div className="font-semibold truncate">{order.gameName} Boosting</div>
                                                            <div className="text-sm text-muted-foreground mt-1">
                                                                Quantity: x{order.quantity}
                                                            </div>
                                                            {order.scheduledTime && (
                                                                <div className="text-xs text-blue-500 mt-1 flex items-center gap-1">
                                                                    <Clock className="h-3 w-3" />
                                                                    Scheduled: {format(new Date(order.scheduledTime), 'PP p')}
                                                                </div>
                                                            )}
                                                        </div>
                                                        <div className="text-right">
                                                            <div className="font-bold text-lg">¥{order.amount.toFixed(2)}</div>
                                                        </div>
                                                    </div>
                                                </CardContent>
                                                <CardFooter className="bg-muted/10 p-2 px-4 flex justify-end gap-2">
                                                    <Button variant="ghost" size="sm">View Details</Button>
                                                    {([OrderStatus.PENDING] as OrderStatus[]).includes(order.status) && (
                                                        <Button size="sm">Pay Now</Button>
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
