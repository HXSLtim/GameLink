import { useEffect, useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { http } from '@/lib/http';
import { useWalletStore } from '@/stores';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardFooter, CardHeader, CardTitle } from '@/components/ui/card';
import { PaymentMethods, type PaymentMethodType } from '@/components/payment/payment-methods';
import { Loader2, ArrowLeft, Clock, ShieldCheck } from 'lucide-react';
import { toast } from 'sonner';

interface OrderDetails {
    id: string;
    orderNo: string;
    amount: number;
    status: string;
    serviceName: string;
    createdAt: string;
}

export default function PaymentPage() {
    const { orderId } = useParams();
    const navigate = useNavigate();
    const { t } = useTranslation();
    const { fetchWallet, getBalance } = useWalletStore();

    const balance = getBalance() / 100; // Convert cents to yuan

    const [order, setOrder] = useState<OrderDetails | null>(null);
    const [loading, setLoading] = useState(true);
    const [processPayment, setProcessPayment] = useState(false);
    const [paymentMethod, setPaymentMethod] = useState<PaymentMethodType>('balance');

    useEffect(() => {
        const init = async () => {
            try {
                // Fetch wallet to get latest balance
                await fetchWallet();

                // Fetch Order Details
                const orderData = await http.get<OrderDetails>(`/user/orders/${orderId}`);
                setOrder(orderData);

                // If order is not pending, redirect
                if (orderData.status !== 'pending') {
                    toast.info(t('payment.order_not_pending'));
                    navigate(`/orders/${orderId}`);
                }
            } catch (err) {
                console.error(err);
                toast.error(t('payment.fetch_failed'));
                navigate('/orders');
            } finally {
                setLoading(false);
            }
        };

        if (orderId) {
            init();
        }
    }, [orderId, fetchWallet, navigate, t]);

    const handlePay = async () => {
        if (!order) return;

        if (paymentMethod === 'balance' && (balance || 0) < order.amount) {
            toast.error(t('payment.insufficient_balance'));
            // Optionally navigate to recharge page
            // navigate('/wallet/recharge');
            return;
        }

        setProcessPayment(true);
        try {
            await http.post('/user/payments', {
                orderId: order.id,
                method: paymentMethod,
                amount: order.amount
            });

            toast.success(t('payment.success'));
            navigate(`/orders/${orderId}`, { replace: true });
        } catch (err: any) {
            console.error(err);
            toast.error(err.message || t('payment.failed'));
        } finally {
            setProcessPayment(false);
        }
    };

    if (loading) {
        return (
            <div className="flex h-screen items-center justify-center">
                <Loader2 className="h-8 w-8 animate-spin text-primary" />
            </div>
        );
    }

    if (!order) return null;

    return (
        <div className="min-h-screen bg-muted/30 p-4 md:p-8">
            <div className="mx-auto max-w-2xl space-y-6">
                <Button variant="ghost" className="pl-0 gap-2" onClick={() => navigate(-1)}>
                    <ArrowLeft className="h-4 w-4" />
                    {t('common.back', { defaultValue: 'Back' })}
                </Button>

                <div className="grid gap-6">
                    <div className="flex flex-col items-center space-y-2 text-center">
                        <h1 className="text-3xl font-bold tracking-tight">{t('payment.checkout', { defaultValue: 'Checkout' })}</h1>
                        <p className="text-muted-foreground">{t('payment.complete_transaction', { defaultValue: 'Complete your transaction securely' })}</p>
                    </div>

                    <Card>
                        <CardHeader className="pb-4">
                            <CardTitle className="text-base font-medium text-muted-foreground">
                                {t('payment.order_summary', { defaultValue: 'Order Summary' })}
                            </CardTitle>
                        </CardHeader>
                        <CardContent className="grid gap-4">
                            <div className="flex items-center justify-between">
                                <span className="font-medium">{order.serviceName}</span>
                                <span className="font-mono text-lg">¥{order.amount.toFixed(2)}</span>
                            </div>
                            <div className="flex items-center justify-between text-sm text-muted-foreground">
                                <span>{t('payment.order_no', { defaultValue: 'Order No.' })}</span>
                                <span className="font-mono">{order.orderNo}</span>
                            </div>
                            <div className="flex items-center gap-2 rounded-md bg-yellow-500/10 p-3 text-sm text-yellow-600 dark:text-yellow-400">
                                <Clock className="h-4 w-4" />
                                <span>{t('payment.expires_in', { defaultValue: 'pay within 15 mins' })}</span>
                            </div>
                        </CardContent>
                    </Card>

                    <Card>
                        <CardHeader>
                            <CardTitle>{t('payment.payment_method', { defaultValue: 'Payment Method' })}</CardTitle>
                        </CardHeader>
                        <CardContent>
                            <PaymentMethods
                                value={paymentMethod}
                                onChange={setPaymentMethod}
                                balance={balance || 0}
                                amount={order.amount}
                            />
                        </CardContent>
                        <CardFooter className="flex flex-col gap-4 pt-4">
                            <Button
                                className="w-full h-12 text-lg"
                                onClick={handlePay}
                                disabled={processPayment || (paymentMethod === 'balance' && (balance || 0) < order.amount)}
                            >
                                {processPayment && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                                {t('payment.pay_now', { defaultValue: 'Pay' })} ¥{order.amount.toFixed(2)}
                            </Button>

                            <div className="flex items-center justify-center gap-2 text-xs text-muted-foreground">
                                <ShieldCheck className="h-3 w-3" />
                                {t('payment.secure_payment', { defaultValue: 'Secure Payment' })}
                            </div>
                        </CardFooter>
                    </Card>
                </div>
            </div>
        </div>
    );
}
