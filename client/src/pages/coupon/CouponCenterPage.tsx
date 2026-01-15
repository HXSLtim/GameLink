import { useEffect, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { PageContainer } from '@/components/page-container';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from '@/components/ui/card';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { Badge } from '@/components/ui/badge';
import { Skeleton } from '@/components/ui/skeleton';
import { Ticket, CalendarClock, CheckCircle2, Tag, Percent, Gift, Clock } from 'lucide-react';
import { toast } from 'sonner';
import { useCouponStore, CouponType, CouponStatus, type Coupon, type CouponTemplate } from '@/stores';
import { format } from 'date-fns';

export default function CouponCenterPage() {
    const { t } = useTranslation();
    const {
        myCoupons,
        availableCoupons,
        couponCounts,
        loading,
        fetchMyCoupons,
        fetchAvailableCoupons,
        fetchCouponCounts,
        claimCoupon
    } = useCouponStore();

    const [activeTab, setActiveTab] = useState('market');
    const [claimingId, setClaimingId] = useState<number | null>(null);

    useEffect(() => {
        fetchAvailableCoupons();
        fetchMyCoupons();
        fetchCouponCounts();
    }, [fetchAvailableCoupons, fetchMyCoupons, fetchCouponCounts]);

    const handleClaim = async (templateId: number) => {
        setClaimingId(templateId);
        try {
            await claimCoupon(templateId);
            toast.success(t('coupon.claim_success', { defaultValue: 'Coupon claimed successfully!' }));
        } catch (err: any) {
            const message = err.message || t('coupon.claim_failed', { defaultValue: 'Failed to claim coupon' });
            toast.error(message);
        } finally {
            setClaimingId(null);
        }
    };

    const availableMyCoupons = myCoupons.filter(c => c.status === CouponStatus.AVAILABLE);
    const usedMyCoupons = myCoupons.filter(c => c.status === CouponStatus.USED || c.status === CouponStatus.EXPIRED);

    return (
        <PageContainer>
            <div className="max-w-4xl mx-auto py-8 px-4 space-y-6">
                {/* Header */}
                <div className="text-center space-y-2 mb-8 animate-in fade-in slide-in-from-top-2">
                    <div className="mx-auto w-14 h-14 bg-gradient-to-br from-primary/20 to-primary/5 rounded-2xl flex items-center justify-center text-primary shadow-lg">
                        <Ticket className="w-7 h-7" />
                    </div>
                    <h1 className="text-3xl font-bold tracking-tight bg-clip-text text-transparent bg-gradient-to-r from-white to-white/60">
                        {t('coupon.title', { defaultValue: 'Coupon Center' })}
                    </h1>
                    <p className="text-muted-foreground">
                        {t('coupon.subtitle', { defaultValue: 'Find the best deals and save money.' })}
                    </p>
                </div>

                {/* Stats Cards */}
                <div className="grid grid-cols-3 gap-4 mb-6">
                    <Card className="bg-gradient-to-br from-green-500/10 to-green-500/5 border-green-500/20">
                        <CardContent className="p-4 text-center">
                            <div className="text-2xl font-bold text-green-500">{couponCounts.available}</div>
                            <div className="text-xs text-muted-foreground">{t('coupon.available_count', { defaultValue: 'Available' })}</div>
                        </CardContent>
                    </Card>
                    <Card className="bg-gradient-to-br from-blue-500/10 to-blue-500/5 border-blue-500/20">
                        <CardContent className="p-4 text-center">
                            <div className="text-2xl font-bold text-blue-500">{couponCounts.used}</div>
                            <div className="text-xs text-muted-foreground">{t('coupon.used_count', { defaultValue: 'Used' })}</div>
                        </CardContent>
                    </Card>
                    <Card className="bg-gradient-to-br from-gray-500/10 to-gray-500/5 border-gray-500/20">
                        <CardContent className="p-4 text-center">
                            <div className="text-2xl font-bold text-gray-500">{couponCounts.expired}</div>
                            <div className="text-xs text-muted-foreground">{t('coupon.expired_count', { defaultValue: 'Expired' })}</div>
                        </CardContent>
                    </Card>
                </div>

                {/* Tabs */}
                <Tabs value={activeTab} onValueChange={setActiveTab} className="w-full">
                    <TabsList className="grid w-full grid-cols-3 mb-8 p-1 bg-muted/40 backdrop-blur-md border border-white/5 rounded-full">
                        <TabsTrigger value="market" className="rounded-full data-[state=active]:bg-background/80">
                            {t('coupon.market', { defaultValue: 'Get Coupons' })}
                        </TabsTrigger>
                        <TabsTrigger value="available" className="rounded-full data-[state=active]:bg-background/80">
                            {t('coupon.my_available', { defaultValue: 'My Coupons' })} ({availableMyCoupons.length})
                        </TabsTrigger>
                        <TabsTrigger value="history" className="rounded-full data-[state=active]:bg-background/80">
                            {t('coupon.history', { defaultValue: 'History' })}
                        </TabsTrigger>
                    </TabsList>

                    {/* Market Tab */}
                    <TabsContent value="market" className="space-y-4">
                        {loading && availableCoupons.length === 0 ? (
                            <div className="grid gap-4 md:grid-cols-2">
                                {[1, 2, 3, 4].map(i => (
                                    <Skeleton key={i} className="h-48 rounded-xl" />
                                ))}
                            </div>
                        ) : availableCoupons.length === 0 ? (
                            <EmptyState message={t('coupon.no_available', { defaultValue: 'No coupons available at the moment.' })} />
                        ) : (
                            <div className="grid gap-4 md:grid-cols-2">
                                {availableCoupons.map((template) => (
                                    <CouponTemplateCard
                                        key={template.id}
                                        template={template}
                                        onClaim={() => handleClaim(template.id)}
                                        claiming={claimingId === template.id}
                                    />
                                ))}
                            </div>
                        )}
                    </TabsContent>

                    {/* My Available Coupons Tab */}
                    <TabsContent value="available" className="space-y-4">
                        {loading && availableMyCoupons.length === 0 ? (
                            <div className="grid gap-4 md:grid-cols-2">
                                {[1, 2].map(i => (
                                    <Skeleton key={i} className="h-48 rounded-xl" />
                                ))}
                            </div>
                        ) : availableMyCoupons.length === 0 ? (
                            <EmptyState message={t('coupon.no_owned', { defaultValue: "You haven't claimed any coupons yet." })} />
                        ) : (
                            <div className="grid gap-4 md:grid-cols-2">
                                {availableMyCoupons.map((coupon) => (
                                    <MyCouponCard key={coupon.id} coupon={coupon} />
                                ))}
                            </div>
                        )}
                    </TabsContent>

                    {/* History Tab */}
                    <TabsContent value="history" className="space-y-4">
                        {usedMyCoupons.length === 0 ? (
                            <EmptyState message={t('coupon.no_history', { defaultValue: 'No coupon history yet.' })} />
                        ) : (
                            <div className="grid gap-4 md:grid-cols-2">
                                {usedMyCoupons.map((coupon) => (
                                    <MyCouponCard key={coupon.id} coupon={coupon} isHistory />
                                ))}
                            </div>
                        )}
                    </TabsContent>
                </Tabs>
            </div>
        </PageContainer>
    );
}

function EmptyState({ message }: { message: string }) {
    return (
        <div className="flex flex-col items-center justify-center py-20 text-center animate-in fade-in zoom-in-95">
            <div className="p-6 rounded-full bg-muted/30 mb-4">
                <Ticket className="h-10 w-10 text-muted-foreground/50" />
            </div>
            <p className="text-muted-foreground">{message}</p>
        </div>
    );
}

function CouponTemplateCard({
    template,
    onClaim,
    claiming
}: {
    template: CouponTemplate;
    onClaim: () => void;
    claiming: boolean;
}) {
    const { t } = useTranslation();

    const getTypeIcon = () => {
        switch (template.type) {
            case CouponType.DISCOUNT: return <Percent className="w-4 h-4" />;
            case CouponType.AMOUNT: return <Tag className="w-4 h-4" />;
            case CouponType.FREE: return <Gift className="w-4 h-4" />;
            default: return <Ticket className="w-4 h-4" />;
        }
    };

    const getTypeColor = () => {
        switch (template.type) {
            case CouponType.DISCOUNT: return 'bg-purple-500';
            case CouponType.AMOUNT: return 'bg-blue-500';
            case CouponType.FREE: return 'bg-green-500';
            default: return 'bg-primary';
        }
    };

    const getDiscountDisplay = () => {
        switch (template.type) {
            case CouponType.DISCOUNT:
                return `${Math.round((1 - (template.discountRate || 1)) * 100)}% OFF`;
            case CouponType.AMOUNT:
                return `¥${((template.amountCents || 0) / 100).toFixed(0)}`;
            case CouponType.FREE:
                return t('coupon.free', { defaultValue: 'FREE' });
            default:
                return '';
        }
    };

    const isOutOfStock = template.remainingCount <= 0;

    return (
        <Card className="group relative overflow-hidden border-white/5 bg-background/40 backdrop-blur-md hover:border-primary/50 transition-all duration-300">
            <div className="absolute inset-0 bg-gradient-to-br from-primary/5 via-transparent to-transparent opacity-0 group-hover:opacity-100 transition-opacity" />

            <Badge className={`absolute top-3 right-3 ${getTypeColor()} text-white gap-1`}>
                {getTypeIcon()}
                {template.type === CouponType.DISCOUNT ? t('coupon.type_discount', { defaultValue: 'Discount' }) :
                 template.type === CouponType.AMOUNT ? t('coupon.type_amount', { defaultValue: 'Cash Off' }) :
                 t('coupon.type_free', { defaultValue: 'Free' })}
            </Badge>

            <CardHeader className="pb-2">
                <CardTitle className="text-3xl font-bold text-primary">
                    {getDiscountDisplay()}
                </CardTitle>
                <CardDescription>
                    {template.minOrderCents > 0
                        ? t('coupon.min_spend', { defaultValue: 'Min. spend ¥{{amount}}', amount: (template.minOrderCents / 100).toFixed(0) })
                        : t('coupon.no_minimum', { defaultValue: 'No minimum spend' })}
                </CardDescription>
            </CardHeader>

            <CardContent className="pb-2">
                <h3 className="font-semibold mb-1">{template.name}</h3>
                <div className="flex items-center justify-between text-xs text-muted-foreground mt-3">
                    <div className="flex items-center gap-1">
                        <CalendarClock className="w-3 h-3" />
                        {t('coupon.valid_days', { defaultValue: 'Valid for {{days}} days', days: template.validDays })}
                    </div>
                    <div className="flex items-center gap-1">
                        <Clock className="w-3 h-3" />
                        {t('coupon.remaining', { defaultValue: '{{count}} left', count: template.remainingCount })}
                    </div>
                </div>
            </CardContent>

            <CardFooter className="pt-2 bg-muted/10 border-t border-white/5">
                <Button
                    className="w-full rounded-full"
                    onClick={onClaim}
                    disabled={claiming || isOutOfStock}
                >
                    {claiming ? t('common.loading', { defaultValue: 'Loading...' }) :
                     isOutOfStock ? t('coupon.out_of_stock', { defaultValue: 'Out of Stock' }) :
                     t('common.claim', { defaultValue: 'Claim Now' })}
                </Button>
            </CardFooter>
        </Card>
    );
}

function MyCouponCard({ coupon, isHistory = false }: { coupon: Coupon; isHistory?: boolean }) {
    const { t } = useTranslation();

    const getDiscountDisplay = () => {
        switch (coupon.type) {
            case CouponType.DISCOUNT:
                return `${Math.round((1 - (coupon.discountRate || 1)) * 100)}% OFF`;
            case CouponType.AMOUNT:
                return `¥${((coupon.amountCents || 0) / 100).toFixed(0)}`;
            case CouponType.FREE:
                return t('coupon.free', { defaultValue: 'FREE' });
            default:
                return '';
        }
    };

    const isExpired = coupon.status === CouponStatus.EXPIRED || new Date(coupon.validUntil) < new Date();
    const isUsed = coupon.status === CouponStatus.USED;

    return (
        <Card className={`relative overflow-hidden border-dashed ${isHistory ? 'opacity-60' : ''}`}>
            {(isExpired || isUsed) && (
                <div className="absolute inset-0 bg-background/50 backdrop-blur-[1px] z-10 flex items-center justify-center">
                    <Badge variant="secondary" className="text-lg px-4 py-1">
                        {isUsed ? t('coupon.used', { defaultValue: 'Used' }) : t('coupon.expired', { defaultValue: 'Expired' })}
                    </Badge>
                </div>
            )}

            <CardHeader className="pb-2">
                <div className="flex justify-between items-start">
                    <div>
                        <CardTitle className="text-2xl font-bold text-primary">
                            {getDiscountDisplay()}
                        </CardTitle>
                        <CardDescription className="mt-1">
                            {coupon.minOrderCents > 0
                                ? t('coupon.min_spend', { defaultValue: 'Min. spend ¥{{amount}}', amount: (coupon.minOrderCents / 100).toFixed(0) })
                                : t('coupon.no_minimum', { defaultValue: 'No minimum spend' })}
                        </CardDescription>
                    </div>
                </div>
            </CardHeader>

            <CardContent className="pb-2">
                <h3 className="font-semibold mb-1">{coupon.name}</h3>
                <div className="flex items-center text-xs text-muted-foreground mt-3">
                    <CalendarClock className="w-3 h-3 mr-1" />
                    {t('coupon.expires', { defaultValue: 'Expires: {{date}}', date: format(new Date(coupon.validUntil), 'yyyy-MM-dd') })}
                </div>
            </CardContent>

            {!isHistory && !isExpired && !isUsed && (
                <CardFooter className="pt-2 bg-muted/10 border-t border-white/5 flex justify-between items-center">
                    <span className="flex items-center text-sm font-medium text-green-500">
                        <CheckCircle2 className="w-4 h-4 mr-1" />
                        {t('coupon.owned', { defaultValue: 'Ready to use' })}
                    </span>
                    <Button variant="link" size="sm" className="h-auto p-0 text-primary">
                        {t('coupon.use_now', { defaultValue: 'Use Now' })}
                    </Button>
                </CardFooter>
            )}
        </Card>
    );
}
