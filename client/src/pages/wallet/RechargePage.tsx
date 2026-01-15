import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { PageContainer } from '@/components/page-container';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { useWalletStore } from '@/stores';
import { PaymentMethods } from '@/components/payment/payment-methods';
import type { PaymentMethodType } from '@/components/payment/payment-methods';
import { toast } from 'sonner';
import { Coins, ShieldCheck } from 'lucide-react';
import { useNavigate } from 'react-router-dom';

const RECHARGE_PLANS = [
    { id: 1, amount: 6, bonus: 0 },
    { id: 2, amount: 30, bonus: 2 },
    { id: 3, amount: 68, bonus: 5 },
    { id: 4, amount: 128, bonus: 12 },
    { id: 5, amount: 328, bonus: 35 },
    { id: 6, amount: 648, bonus: 70 },
];

export default function RechargePage() {
    const { t } = useTranslation();
    const navigate = useNavigate();
    const { recharge, loading } = useWalletStore();

    const [selectedPlanId, setSelectedPlanId] = useState<number | null>(2); // Default to 30
    const [paymentMethod, setPaymentMethod] = useState<PaymentMethodType>('alipay');

    const selectedPlan = RECHARGE_PLANS.find(p => p.id === selectedPlanId);

    const handleRecharge = async () => {
        if (!selectedPlan) return;

        // Only allow wechat or alipay for recharge
        const rechargeMethod = paymentMethod === 'balance' ? 'alipay' : paymentMethod as 'wechat' | 'alipay';

        try {
            await recharge(selectedPlan.amount * 100, rechargeMethod); // Convert to cents
            toast.success(t('wallet.recharge_success', { defaultValue: 'Recharge successful!' }));
            navigate('/wallet');
        } catch {
            toast.error(t('wallet.recharge_failed', { defaultValue: 'Recharge failed' }));
        }
    };

    return (
        <PageContainer>
            <div className="max-w-3xl mx-auto py-8 px-4 space-y-6">
                <div className="text-center space-y-2 mb-8">
                    <div className="mx-auto w-12 h-12 bg-primary/10 rounded-full flex items-center justify-center text-primary">
                        <Coins className="w-6 h-6" />
                    </div>
                    <h1 className="text-3xl font-bold tracking-tight">{t('wallet.recharge', { defaultValue: 'Recharge Balance' })}</h1>
                    <p className="text-muted-foreground">{t('wallet.recharge_desc', { defaultValue: 'Top up your wallet to pay for orders.' })}</p>
                </div>

                <Card>
                    <CardHeader>
                        <CardTitle>{t('wallet.select_amount', { defaultValue: 'Select Amount' })}</CardTitle>
                    </CardHeader>
                    <CardContent>
                        <div className="grid grid-cols-2 md:grid-cols-3 gap-4">
                            {RECHARGE_PLANS.map((plan) => (
                                <div
                                    key={plan.id}
                                    onClick={() => setSelectedPlanId(plan.id)}
                                    className={`
                                        cursor-pointer relative overflow-hidden rounded-xl border-2 p-4 text-center transition-all hover:bg-muted/50
                                        ${selectedPlanId === plan.id ? 'border-primary bg-primary/5' : 'border-muted'}
                                    `}
                                >
                                    <div className="text-2xl font-bold">¥{plan.amount}</div>
                                    {plan.bonus > 0 && (
                                        <div className="text-xs text-green-600 font-medium mt-1">
                                            +{plan.bonus} {t('common.bonus', { defaultValue: 'Bonus' })}
                                        </div>
                                    )}
                                    {selectedPlanId === plan.id && (
                                        <div className="absolute top-0 right-0 p-1 bg-primary rounded-bl-lg">
                                            <ShieldCheck className="w-3 h-3 text-white" />
                                        </div>
                                    )}
                                </div>
                            ))}
                        </div>
                    </CardContent>
                </Card>

                <Card>
                    <CardHeader>
                        <CardTitle>{t('payment.select_method', { defaultValue: 'Payment Method' })}</CardTitle>
                    </CardHeader>
                    <CardContent>
                        <PaymentMethods
                            value={paymentMethod}
                            onChange={setPaymentMethod}
                            balance={9999}
                            amount={selectedPlan?.amount || 0}
                        />
                    </CardContent>
                </Card>

                <div className="fixed bottom-0 left-0 right-0 p-4 border-t bg-background/95 backdrop-blur md:static md:bg-transparent md:border-0 md:p-0">
                    <div className="max-w-3xl mx-auto flex items-center justify-between gap-4">
                        <div className="text-lg">
                            <span className="text-muted-foreground">{t('common.total', { defaultValue: 'Total' })}: </span>
                            <span className="font-bold text-2xl text-primary">¥{selectedPlan?.amount || 0}</span>
                        </div>
                        <Button
                            size="lg"
                            className="w-1/2 md:w-auto md:min-w-[200px]"
                            onClick={handleRecharge}
                            disabled={loading || !selectedPlan}
                        >
                            {loading ? t('common.processing', { defaultValue: 'Processing...' }) : t('wallet.pay_now', { defaultValue: 'Pay Now' })}
                        </Button>
                    </div>
                </div>
            </div>
        </PageContainer>
    );
}
