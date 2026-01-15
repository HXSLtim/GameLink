import { useEffect, useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { useOrderStore, OrderStatus } from '@/stores';
import { PageContainer } from '@/components/page-container';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Avatar, AvatarFallback } from '@/components/ui/avatar';
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from "@/components/ui/dialog";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import {
    CheckCircle, Gamepad2, ArrowLeft, MessageSquare, CreditCard, Ban, Copy, AlertTriangle, ShieldAlert
} from 'lucide-react';
import { toast } from 'sonner';

export default function OrderDetailPage() {
    const { id } = useParams<{ id: string }>();
    const navigate = useNavigate();
    const { currentOrder, fetchOrderById, loading, cancelOrder, submitDispute } = useOrderStore();

    const [disputeOpen, setDisputeOpen] = useState(false);
    const [disputeType, setDisputeType] = useState('');
    const [disputeDesc, setDisputeDesc] = useState('');
    const [submittingDispute, setSubmittingDispute] = useState(false);

    useEffect(() => {
        if (id) {
            fetchOrderById(id);
        }
    }, [id]);

    const handleCopy = (text: string) => {
        navigator.clipboard.writeText(text);
        toast.success("Copied to clipboard");
    };

    const handleCancel = async () => {
        if (!currentOrder) return;
        if (window.confirm("Are you sure you want to cancel this order?")) {
            await cancelOrder(currentOrder.id);
            toast.success("Order cancelled");
        }
    };

    const handleDisputeSubmit = async () => {
        if (!currentOrder || !disputeType || !disputeDesc) {
            toast.error("Please fill in all fields");
            return;
        }

        setSubmittingDispute(true);
        try {
            await submitDispute(currentOrder.id, disputeType, disputeDesc);
            toast.success("Dispute submitted successfully");
            setDisputeOpen(false);
        } catch (error: any) {
            toast.error(error.message || "Failed to submit dispute");
        } finally {
            setSubmittingDispute(false);
        }
    };

    if (loading || !currentOrder) {
        return (
            <div className="flex items-center justify-center h-screen">
                <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary" />
            </div>
        );
    }

    const steps = [
        { status: OrderStatus.PENDING, label: 'Pending Payment', done: true },
        { status: OrderStatus.CONFIRMED, label: 'Confirmed', done: ([OrderStatus.CONFIRMED, OrderStatus.IN_PROGRESS, OrderStatus.COMPLETED, OrderStatus.DISPUTED] as OrderStatus[]).includes(currentOrder.status) },
        { status: OrderStatus.IN_PROGRESS, label: 'In Progress', done: ([OrderStatus.IN_PROGRESS, OrderStatus.COMPLETED, OrderStatus.DISPUTED] as OrderStatus[]).includes(currentOrder.status) },
        { status: OrderStatus.COMPLETED, label: 'Completed', done: currentOrder.status === OrderStatus.COMPLETED }
    ];

    const isCancelled = ([OrderStatus.CANCELED, OrderStatus.REFUNDED] as OrderStatus[]).includes(currentOrder.status);
    const isDisputed = currentOrder.status === OrderStatus.DISPUTED;

    return (
        <PageContainer>
            <div className="max-w-2xl mx-auto py-8 px-4 space-y-6">
                {/* Header */}
                <div className="flex items-center gap-4 animate-in fade-in slide-in-from-top-2">
                    <Button variant="ghost" size="icon" onClick={() => navigate('/orders')} className="rounded-full">
                        <ArrowLeft className="h-5 w-5" />
                    </Button>
                    <div>
                        <h1 className="text-2xl font-bold tracking-tight">Order Details</h1>
                        <p className="text-sm text-muted-foreground flex items-center gap-2">
                            Order #{currentOrder.orderNo}
                            <Copy className="h-3 w-3 cursor-pointer hover:text-primary" onClick={() => handleCopy(currentOrder.orderNo)} />
                        </p>
                    </div>
                    <div className="ml-auto">
                        <Badge variant="outline" className={`
                            px-3 py-1 font-bold uppercase tracking-wider
                            ${currentOrder.status === OrderStatus.COMPLETED ? 'bg-green-500/10 text-green-500 border-green-500/20' : ''}
                            ${currentOrder.status === OrderStatus.PENDING ? 'bg-yellow-500/10 text-yellow-500 border-yellow-500/20' : ''}
                            ${isCancelled ? 'bg-red-500/10 text-red-500 border-red-500/20' : ''}
                            ${isDisputed ? 'bg-orange-500/10 text-orange-500 border-orange-500/20' : ''}
                        `}>
                            {currentOrder.status}
                        </Badge>
                    </div>
                </div>

                {isDisputed && (
                    <Card className="bg-orange-500/10 border-orange-500/20">
                        <CardContent className="flex items-center gap-4 p-4 text-orange-500">
                            <AlertTriangle className="h-5 w-5" />
                            <div>
                                <h4 className="font-bold">Order Under Dispute</h4>
                                <p className="text-sm opacity-90">Our support team is reviewing this order. We will contact you shortly.</p>
                            </div>
                        </CardContent>
                    </Card>
                )}

                {/* Progress Visual */}
                {!isCancelled && !isDisputed && (
                    <Card className="bg-muted/40 border-dashed">
                        <CardContent className="p-6">
                            <div className="relative flex items-center justify-between">
                                {/* Connector Line */}
                                <div className="absolute left-0 top-1/2 w-full h-0.5 bg-muted -z-10" />

                                {steps.map((step, idx) => (
                                    <div key={idx} className="flex flex-col items-center gap-2 bg-background p-2 rounded-lg z-10">
                                        <div className={`
                                            h-8 w-8 rounded-full flex items-center justify-center border-2 transition-colors
                                            ${step.done ? 'bg-primary border-primary text-primary-foreground' : 'bg-muted border-muted-foreground/30 text-muted-foreground'}
                                        `}>
                                            {step.done ? <CheckCircle className="h-4 w-4" /> : <div className="h-2 w-2 rounded-full bg-current" />}
                                        </div>
                                        <span className={`text-[10px] uppercase font-bold ${step.done ? 'text-foreground' : 'text-muted-foreground'}`}>{step.label}</span>
                                    </div>
                                ))}
                            </div>
                        </CardContent>
                    </Card>
                )}

                {/* Main Content */}
                <div className="grid gap-6">
                    {/* Item Info */}
                    <Card>
                        <CardHeader>
                            <CardTitle className="text-base">Service Information</CardTitle>
                        </CardHeader>
                        <CardContent className="flex gap-4">
                            <div className="h-20 w-20 rounded-lg bg-secondary flex items-center justify-center">
                                <Gamepad2 className="h-10 w-10 text-muted-foreground" />
                            </div>
                            <div className="flex-1">
                                <h3 className="font-bold text-lg">{currentOrder.gameName}</h3>
                                <p className="text-muted-foreground text-sm">Service: Gameplay Companion</p>
                                <div className="mt-2 flex gap-4 text-sm">
                                    <div className="flex flex-col">
                                        <span className="text-muted-foreground text-xs">Quantity</span>
                                        <span className="font-medium">{currentOrder.quantity} Hours</span>
                                    </div>
                                    <div className="flex flex-col">
                                        <span className="text-muted-foreground text-xs">Total Amount</span>
                                        <span className="font-medium">¥{currentOrder.amount.toFixed(2)}</span>
                                    </div>
                                </div>
                            </div>
                        </CardContent>
                    </Card>

                    {/* Player Info */}
                    <Card>
                        <CardHeader>
                            <CardTitle className="text-base">Player</CardTitle>
                        </CardHeader>
                        <CardContent className="flex items-center justify-between">
                            <div className="flex items-center gap-3">
                                <Avatar>
                                    <AvatarFallback>P{currentOrder.playerId}</AvatarFallback>
                                </Avatar>
                                <div>
                                    <div className="font-medium">Player #{currentOrder.playerId}</div>
                                    <div className="text-xs text-muted-foreground">Certified Pro</div>
                                </div>
                            </div>
                            <Button variant="outline" size="sm" onClick={() => navigate(`/chat/${currentOrder.playerId}`)}>
                                <MessageSquare className="h-4 w-4 mr-2" />
                                Contact
                            </Button>
                        </CardContent>
                    </Card>
                </div>

                {/* Actions */}
                {!isCancelled && !isDisputed && (
                    <div className="sticky bottom-0 left-0 right-0 p-4 bg-background/80 backdrop-blur-md border-t border-white/5 flex flex-col md:flex-row gap-4 -mx-4 -mb-8 sm:mx-0 sm:mb-0 sm:static sm:bg-transparent sm:border-0">
                        {currentOrder.status === OrderStatus.PENDING && (
                            <Button className="flex-1 font-bold shadow-lg shadow-primary/20" size="lg">
                                <CreditCard className="h-4 w-4 mr-2" />
                                Pay Now
                            </Button>
                        )}
                        {currentOrder.status === OrderStatus.PENDING && (
                            <Button variant="destructive" className="flex-1" size="lg" onClick={handleCancel}>
                                <Ban className="h-4 w-4 mr-2" />
                                Cancel Order
                            </Button>
                        )}
                        {currentOrder.status === OrderStatus.COMPLETED && (
                            <Button className="flex-1" size="lg" onClick={() => navigate(`/players/${currentOrder.playerId}`)}>
                                Book Again
                            </Button>
                        )}

                        {(currentOrder.status === OrderStatus.IN_PROGRESS || currentOrder.status === OrderStatus.COMPLETED) && (
                            <Button variant="outline" className="flex-1 text-destructive hover:bg-destructive/10 border-destructive/20" onClick={() => setDisputeOpen(true)}>
                                <ShieldAlert className="h-4 w-4 mr-2" />
                                Report Issue
                            </Button>
                        )}
                    </div>
                )}

                <Dialog open={disputeOpen} onOpenChange={setDisputeOpen}>
                    <DialogContent>
                        <DialogHeader>
                            <DialogTitle>Report an Issue</DialogTitle>
                            <DialogDescription>
                                Please describe the issue you are experiencing with this order. Our support team will review it.
                            </DialogDescription>
                        </DialogHeader>
                        <div className="grid gap-4 py-4">
                            <div className="grid gap-2">
                                <Label htmlFor="type">Issue Type</Label>
                                <Select onValueChange={setDisputeType}>
                                    <SelectTrigger>
                                        <SelectValue placeholder="Select issue type" />
                                    </SelectTrigger>
                                    <SelectContent>
                                        <SelectItem value="service_quality">Service Quality</SelectItem>
                                        <SelectItem value="bad_attitude">Bad Attitude</SelectItem>
                                        <SelectItem value="incomplete_service">Incomplete Service</SelectItem>
                                        <SelectItem value="other">Other</SelectItem>
                                    </SelectContent>
                                </Select>
                            </div>
                            <div className="grid gap-2">
                                <Label htmlFor="desc">Description</Label>
                                <Textarea
                                    id="desc"
                                    placeholder="Please provide details..."
                                    value={disputeDesc}
                                    onChange={(e) => setDisputeDesc(e.target.value)}
                                />
                            </div>
                        </div>
                        <DialogFooter>
                            <Button variant="outline" onClick={() => setDisputeOpen(false)}>Cancel</Button>
                            <Button onClick={handleDisputeSubmit} disabled={submittingDispute}>
                                {submittingDispute ? 'Submitting...' : 'Submit Dispute'}
                            </Button>
                        </DialogFooter>
                    </DialogContent>
                </Dialog>
            </div>
        </PageContainer>
    );
}
