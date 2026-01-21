import { useEffect, useState } from 'react';
import { useLocation, useNavigate } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { PageContainer } from '@/components/page-container';
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Separator } from '@/components/ui/separator';
import { Label } from '@/components/ui/label';
import { Input } from '@/components/ui/input';
import { useOrderStore } from '@/stores';
import { ArrowLeft, CreditCard, Loader2, Minus, Plus, ShieldCheck } from 'lucide-react';
import { toast } from 'sonner';

export default function CreateOrderPage() {
    const navigate = useNavigate();
    const location = useLocation();
    const { t } = useTranslation();
    const { createOrder, loading } = useOrderStore();

    // Parse query params for initial data
    const searchParams = new URLSearchParams(location.search);
    const playerId = searchParams.get('playerId');
    const serviceId = searchParams.get('serviceId');
    const priceStr = searchParams.get('price');
    const gameName = searchParams.get('gameName') || 'Game Service';

    const pricePerUnit = priceStr ? parseFloat(priceStr) : 10;
    const [quantity, setQuantity] = useState(1);
    const [note, setNote] = useState('');

    useEffect(() => {
        if (!playerId) {
            toast.error(t('orders.error_missing_params', { defaultValue: 'Missing required order parameters' }));
            navigate('/');
        }
    }, [playerId, navigate, t]);

    const totalAmount = pricePerUnit * quantity;

    const handleQuantityChange = (delta: number) => {
        setQuantity(q => Math.max(1, q + delta));
    };

    const handleCreateOrder = async () => {
        if (!playerId) return;

        try {
            // Create Order
            const orderPayload = {
                playerId: parseInt(playerId),
                serviceItemId: serviceId ? parseInt(serviceId) : undefined,
                gameId: 1, // TODO: Get from query params or player data
                quantity,
                amount: totalAmount,
                note
            };

            await createOrder(orderPayload);

            // Get the created order from store and navigate to payment
            const { currentOrder } = useOrderStore.getState();
            if (currentOrder?.id) {
                toast.success(t('orders.create_success', { defaultValue: 'Order created! Redirecting to payment...' }));
                navigate(`/payment/${currentOrder.id}`);
            } else {
                // Fallback to orders list if no order ID
                navigate('/orders');
            }

        } catch (error) {
            console.error(error);
            toast.error(error instanceof Error ? error.message : t('orders.create_failed', { defaultValue: 'Failed to create order' }));
        }
    };

    if (!playerId) return null;

    return (
        <PageContainer>
            <div className="max-w-3xl mx-auto py-8 px-4 space-y-6">
                <Button variant="ghost" className="pl-0 hover:bg-transparent" onClick={() => navigate(-1)}>
                    <ArrowLeft className="h-4 w-4 mr-2" />
                    {t('common.back', { defaultValue: 'Back' })}
                </Button>

                <Card className="border-0 shadow-xl bg-card/60 backdrop-blur-xl animate-in fade-in slide-in-from-bottom-4">
                    <CardHeader>
                        <CardTitle className="text-2xl">{t('orders.review_order', { defaultValue: 'Review Order' })}</CardTitle>
                        <CardDescription>{t('orders.review_order_desc', { defaultValue: 'Please review your order details before payment.' })}</CardDescription>
                    </CardHeader>
                    <CardContent className="space-y-6">
                        {/* Service Details */}
                        <div className="flex justify-between items-start p-4 bg-muted/30 rounded-lg">
                            <div>
                                <h3 className="font-semibold text-lg">{gameName}</h3>
                                <p className="text-muted-foreground text-sm">Player ID: {playerId}</p>
                            </div>
                            <div className="text-right">
                                <div className="font-mono text-xl text-primary font-bold">${pricePerUnit}/hr</div>
                                <div className="text-xs text-muted-foreground">Base Price</div>
                            </div>
                        </div>

                        <Separator />

                        {/* Configuration */}
                        <div className="grid gap-6 md:grid-cols-2">
                            <div className="space-y-2">
                                <Label>{t('orders.quantity', { defaultValue: 'Quantity (Hours)' })}</Label>
                                <div className="flex items-center gap-4">
                                    <Button variant="outline" size="icon" onClick={() => handleQuantityChange(-1)} disabled={quantity <= 1}>
                                        <Minus className="h-4 w-4" />
                                    </Button>
                                    <span className="text-xl font-mono w-12 text-center">{quantity}</span>
                                    <Button variant="outline" size="icon" onClick={() => handleQuantityChange(1)}>
                                        <Plus className="h-4 w-4" />
                                    </Button>
                                </div>
                            </div>
                            <div className="space-y-2">
                                <Label>{t('orders.note', { defaultValue: 'Note to Player (Optional)' })}</Label>
                                <Input
                                    placeholder={t('orders.note_placeholder', { defaultValue: 'E.g., I want to play ranked...' })}
                                    value={note}
                                    onChange={(e) => setNote(e.target.value)}
                                />
                            </div>
                        </div>

                        {/* Summary */}
                        <div className="bg-primary/5 p-6 rounded-xl space-y-3">
                            <div className="flex justify-between text-sm">
                                <span className="text-muted-foreground">{t('orders.subtotal', { defaultValue: 'Subtotal' })}</span>
                                <span>${(totalAmount ?? 0).toFixed(2)}</span>
                            </div>
                            <div className="flex justify-between text-sm">
                                <span className="text-muted-foreground">{t('orders.service_fee', { defaultValue: 'Service Fee' })}</span>
                                <span>$0.00</span>
                            </div>
                            <Separator className="bg-primary/10" />
                            <div className="flex justify-between items-end pt-2">
                                <span className="font-semibold">{t('orders.total', { defaultValue: 'Total' })}</span>
                                <span className="text-3xl font-bold text-primary">${(totalAmount ?? 0).toFixed(2)}</span>
                            </div>
                        </div>
                    </CardContent>
                    <CardFooter className="flex flex-col gap-4">
                        <Button
                            className="w-full h-12 text-lg shadow-lg shadow-primary/25"
                            onClick={handleCreateOrder}
                            disabled={loading}
                        >
                            {loading ? (
                                <Loader2 className="mr-2 h-5 w-5 animate-spin" />
                            ) : (
                                <CreditCard className="mr-2 h-5 w-5" />
                            )}
                            {loading ? t('orders.processing', { defaultValue: 'Processing...' }) : t('orders.confirm_pay', { defaultValue: 'Confirm & Pay' })}
                        </Button>
                        <div className="flex items-center justify-center gap-2 text-xs text-muted-foreground">
                            <ShieldCheck className="h-3 w-3" />
                            {t('orders.secure_payment', { defaultValue: 'Secure payment powered by GameLink' })}
                        </div>
                    </CardFooter>
                </Card>
            </div>
        </PageContainer>
    );
}
