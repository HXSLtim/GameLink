import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { PageContainer } from '@/components/page-container';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from '@/components/ui/card';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { Ticket, CalendarClock, CheckCircle2 } from 'lucide-react';
import { toast } from 'sonner';

interface Coupon {
    id: string;
    code: string;
    discount: number;
    type: 'percent' | 'fixed';
    minSpend: number;
    title: string;
    desc: string;
    expiresAt: string;
    status: 'available' | 'claimed' | 'used' | 'expired';
}

const MOCK_COUPONS: Coupon[] = [
    {
        id: '1',
        code: 'NEWPLAYER',
        discount: 20,
        type: 'percent',
        minSpend: 50,
        title: 'New Player Discount',
        desc: 'Get 20% off your first order',
        expiresAt: '2026-12-31',
        status: 'available',
    },
    {
        id: '2',
        code: 'SAVE10',
        discount: 10,
        type: 'fixed',
        minSpend: 100,
        title: 'Save ¥10',
        desc: 'Save ¥10 on orders over ¥100',
        expiresAt: '2026-06-30',
        status: 'available',
    },
    {
        id: '3',
        code: 'WELCOME',
        discount: 5,
        type: 'fixed',
        minSpend: 0,
        title: 'Welcome Gift',
        desc: 'Free ¥5 coupon',
        expiresAt: '2026-02-28',
        status: 'claimed',
    }
];

export default function CouponCenterPage() {
    const { t } = useTranslation();
    const [coupons, setCoupons] = useState<Coupon[]>(MOCK_COUPONS);

    const handleClaim = (id: string) => {
        setCoupons(prev => prev.map(c => c.id === id ? { ...c, status: 'claimed' } : c));
        toast.success(t('coupon.claim_success', { defaultValue: 'Coupon claimed successfully!' }));
    };

    const marketCoupons = coupons.filter(c => c.status === 'available');
    const myCoupons = coupons.filter(c => ['claimed', 'used'].includes(c.status));

    return (
        <PageContainer>
            <div className="max-w-3xl mx-auto py-8 px-4 space-y-6">
                <div className="text-center space-y-2 mb-8">
                    <div className="mx-auto w-12 h-12 bg-primary/10 rounded-full flex items-center justify-center text-primary">
                        <Ticket className="w-6 h-6" />
                    </div>
                    <h1 className="text-3xl font-bold tracking-tight">{t('coupon.title', { defaultValue: 'Coupon Center' })}</h1>
                    <p className="text-muted-foreground">{t('coupon.subtitle', { defaultValue: 'Find the best deals and save money.' })}</p>
                </div>

                <Tabs defaultValue="market" className="w-full">
                    <TabsList className="grid w-full grid-cols-2 mb-8">
                        <TabsTrigger value="market">{t('coupon.market', { defaultValue: 'Coupon Market' })}</TabsTrigger>
                        <TabsTrigger value="mine">{t('coupon.mine', { defaultValue: 'My Coupons' })}</TabsTrigger>
                    </TabsList>

                    <TabsContent value="market" className="space-y-4">
                        {marketCoupons.length === 0 ? (
                            <div className="text-center py-20 text-muted-foreground">
                                {t('coupon.no_available', { defaultValue: 'No coupons available at the moment.' })}
                            </div>
                        ) : (
                            <div className="grid gap-4 md:grid-cols-2">
                                {marketCoupons.map((coupon) => (
                                    <CouponCard key={coupon.id} coupon={coupon} onAction={() => handleClaim(coupon.id)} actionLabel={t('common.claim', { defaultValue: 'Claim' })} />
                                ))}
                            </div>
                        )}
                    </TabsContent>

                    <TabsContent value="mine" className="space-y-4">
                        {myCoupons.length === 0 ? (
                            <div className="text-center py-20 text-muted-foreground">
                                {t('coupon.no_owned', { defaultValue: 'You haven\'t claimed any coupons yet.' })}
                            </div>
                        ) : (
                            <div className="grid gap-4 md:grid-cols-2">
                                {myCoupons.map((coupon) => (
                                    <CouponCard key={coupon.id} coupon={coupon} isOwned />
                                ))}
                            </div>
                        )}
                    </TabsContent>
                </Tabs>
            </div>
        </PageContainer>
    );
}

function CouponCard({ coupon, onAction, actionLabel, isOwned = false }: { coupon: Coupon, onAction?: () => void, actionLabel?: string, isOwned?: boolean }) {
    const { t } = useTranslation();

    return (
        <Card className="flex flex-col overflow-hidden relative border-dashed">
            <div className={`absolute top-0 right-0 p-1 px-3 rounded-bl-lg text-xs font-bold text-white ${coupon.type === 'percent' ? 'bg-purple-500' : 'bg-blue-500'}`}>
                {coupon.type === 'percent' ? 'Discount' : 'Cash Off'}
            </div>
            <CardHeader className="pb-2">
                <div className="flex justify-between items-start">
                    <div>
                        <CardTitle className="text-2xl font-bold text-primary">
                            {coupon.type === 'percent' ? `${coupon.discount}% OFF` : `¥${coupon.discount}`}
                        </CardTitle>
                        <CardDescription className="mt-1">
                            {t('coupon.min_spend', { defaultValue: 'Min. spend ¥{{amount}}', amount: coupon.minSpend })}
                        </CardDescription>
                    </div>
                </div>
            </CardHeader>
            <CardContent className="flex-1 pb-2">
                <h3 className="font-semibold mb-1">{coupon.title}</h3>
                <p className="text-sm text-muted-foreground">{coupon.desc}</p>
                <div className="flex items-center text-xs text-muted-foreground mt-4">
                    <CalendarClock className="w-3 h-3 mr-1" />
                    {t('coupon.expires', { defaultValue: 'Expires: {{date}}', date: coupon.expiresAt })}
                </div>
            </CardContent>
            {onAction && !isOwned && (
                <CardFooter className="pt-2 bg-muted/20">
                    <Button className="w-full" onClick={onAction} variant="secondary">
                        {actionLabel}
                    </Button>
                </CardFooter>
            )}
            {isOwned && (
                <CardFooter className="pt-2 bg-muted/20 flex justify-between items-center text-sm font-medium text-green-600">
                    <span className="flex items-center">
                        <CheckCircle2 className="w-4 h-4 mr-1" />
                        {coupon.status === 'used' ? t('coupon.used', { defaultValue: 'Used' }) : t('coupon.owned', { defaultValue: 'Acquired' })}
                    </span>
                    <Button variant="link" size="sm" className="h-auto p-0 text-primary">
                        {t('coupon.use_now', { defaultValue: 'Use Now' })}
                    </Button>
                </CardFooter>
            )}
        </Card>
    );
}
