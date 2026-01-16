import { useEffect } from 'react';
import { useVipStore } from '@/stores';
import { PageContainer } from '@/components/page-container';
import { Card, CardContent, CardHeader, CardTitle, CardFooter, CardDescription } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { useTranslation } from 'react-i18next';
import { Crown, Check, Zap, Star, Shield, ArrowLeft, Gem } from 'lucide-react';
import { format } from 'date-fns';
import { toast } from 'sonner';
import { useNavigate } from 'react-router-dom';

export default function VipPage() {
    const { t } = useTranslation();
    const navigate = useNavigate();
    const { userVip, levels, fetchVipStatus, claimMonthlyCoupon } = useVipStore();

    const vipUnlocked = userVip?.vipUnlocked ?? false;
    const vipExpireAt = userVip?.vipExpireAt;
    const currentLevelConfig = levels.find(c => c.id === userVip?.vipLevelId);

    useEffect(() => {
        fetchVipStatus();
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, []);

    const handleSubscribe = async (_levelId: number) => {
        try {
            // For now, claim monthly coupon as a placeholder
            await claimMonthlyCoupon();
            toast.success(t('vip.upgrade_success', { defaultValue: 'Successfully upgraded to VIP!' }));
        } catch {
            toast.error(t('vip.upgrade_failed', { defaultValue: 'Subscription failed. Check your wallet balance.' }));
        }
    };

    const FEATURES = [
        "Exclusive Profile Badge",
        "Priority Support",
        "No Service Fees",
        "Ad-free Experience",
        "Monthly Free Boosts (2h)"
    ];

    return (
        <PageContainer>
            <div className="max-w-6xl mx-auto py-8 px-4 space-y-12">
                {/* Header */}
                <div className="flex items-center gap-4 animate-in fade-in slide-in-from-top-2">
                    <Button variant="ghost" size="icon" onClick={() => navigate('/profile')} className="rounded-full">
                        <ArrowLeft className="h-5 w-5" />
                    </Button>
                    <div>
                        <h1 className="text-2xl font-bold tracking-tight">VIP Membership</h1>
                    </div>
                </div>

                {/* Hero Section */}
                <div className="text-center space-y-4 max-w-2xl mx-auto animate-in zoom-in slide-in-from-bottom-4 duration-500">
                    <div className="mx-auto w-20 h-20 bg-gradient-to-br from-yellow-300 to-yellow-600 rounded-full flex items-center justify-center shadow-2xl shadow-yellow-500/20 mb-6">
                        <Crown className="h-10 w-10 text-white" />
                    </div>
                    <h1 className="text-4xl md:text-5xl font-extrabold tracking-tight bg-clip-text text-transparent bg-gradient-to-r from-yellow-500 via-orange-500 to-yellow-500">
                        Unlock Premium Power
                    </h1>
                    <p className="text-xl text-muted-foreground">
                        Elevate your gaming experience with exclusive perks and benefits.
                    </p>
                </div>

                {/* Current Status */}
                {vipUnlocked && (
                    <Card className="bg-gradient-to-r from-yellow-500/10 via-yellow-500/5 to-transparent border-yellow-500/20">
                        <CardContent className="p-6 flex items-center justify-between">
                            <div className="flex items-center gap-4">
                                <div className="h-12 w-12 rounded-full bg-yellow-500/20 flex items-center justify-center">
                                    <Gem className="h-6 w-6 text-yellow-600" />
                                </div>
                                <div>
                                    <div className="font-bold text-lg text-yellow-700 dark:text-yellow-500">Current Plan: {currentLevelConfig?.title || 'VIP'}</div>
                                    <div className="text-sm text-muted-foreground">Expires on {vipExpireAt ? format(new Date(vipExpireAt), 'PPP') : 'Never'}</div>
                                </div>
                            </div>
                            <Badge variant="outline" className="border-yellow-500 text-yellow-600 bg-yellow-500/10 px-3 py-1">ACTIVE</Badge>
                        </CardContent>
                    </Card>
                )}

                {/* Plans Grid */}
                <div className="grid grid-cols-1 md:grid-cols-3 gap-8 pt-8">
                    {/* Free Plan */}
                    <Card className="border-white/5 bg-background/50 backdrop-blur-sm hover:border-primary/20 transition-all">
                        <CardHeader>
                            <CardTitle className="flex items-center gap-2">
                                <Shield className="h-5 w-5 text-muted-foreground" />
                                Standard
                            </CardTitle>
                            <CardDescription>Get started with basic features</CardDescription>
                        </CardHeader>
                        <CardContent className="space-y-4">
                            <div className="text-3xl font-bold">Free</div>
                            <ul className="space-y-2 text-sm text-muted-foreground">
                                <li className="flex items-center gap-2"><Check className="h-4 w-4 text-green-500" /> Basic Profile</li>
                                <li className="flex items-center gap-2"><Check className="h-4 w-4 text-green-500" /> Community Access</li>
                                <li className="flex items-center gap-2"><Check className="h-4 w-4 text-green-500" /> Standard Support</li>
                            </ul>
                        </CardContent>
                        <CardFooter>
                            <Button variant="outline" className="w-full" disabled>Current Plan</Button>
                        </CardFooter>
                    </Card>

                    {/* Pro Plan (Most Popular) */}
                    <Card className="relative border-yellow-500 shadow-2xl shadow-yellow-500/10 bg-gradient-to-b from-background to-yellow-500/5 scale-105 z-10">
                        <div className="absolute top-0 left-1/2 -translate-x-1/2 -translate-y-1/2 bg-gradient-to-r from-yellow-500 to-orange-500 text-white text-xs font-bold px-3 py-1 rounded-full uppercase tracking-wide shadow-lg">
                            Most Popular
                        </div>
                        <CardHeader>
                            <CardTitle className="flex items-center gap-2 text-yellow-500">
                                <Zap className="h-5 w-5 fill-current" />
                                Pro Gamer
                            </CardTitle>
                            <CardDescription>For serious players</CardDescription>
                        </CardHeader>
                        <CardContent className="space-y-4">
                            <div className="flex items-baseline gap-1">
                                <span className="text-4xl font-bold">¥29.9</span>
                                <span className="text-muted-foreground">/mo</span>
                            </div>
                            <ul className="space-y-3 text-sm">
                                {FEATURES.map((feature, i) => (
                                    <li key={i} className="flex items-center gap-2">
                                        <div className="h-5 w-5 rounded-full bg-yellow-500/20 flex items-center justify-center shrink-0">
                                            <Check className="h-3 w-3 text-yellow-600" />
                                        </div>
                                        {feature}
                                    </li>
                                ))}
                            </ul>
                        </CardContent>
                        <CardFooter>
                            <Button className="w-full bg-gradient-to-r from-yellow-500 to-orange-500 hover:from-yellow-600 hover:to-orange-600 text-white font-bold shadow-lg shadow-orange-500/20" onClick={() => handleSubscribe(2)}>
                                Returns 10x Value
                            </Button>
                        </CardFooter>
                    </Card>

                    {/* Elite Plan */}
                    <Card className="border-white/5 bg-background/50 backdrop-blur-sm hover:border-purple-500/20 transition-all">
                        <CardHeader>
                            <CardTitle className="flex items-center gap-2 text-purple-500">
                                <Star className="h-5 w-5" />
                                Elite
                            </CardTitle>
                            <CardDescription>Ultimate status symbol</CardDescription>
                        </CardHeader>
                        <CardContent className="space-y-4">
                            <div className="flex items-baseline gap-1">
                                <span className="text-3xl font-bold">¥299</span>
                                <span className="text-muted-foreground">/yr</span>
                            </div>
                            <ul className="space-y-2 text-sm text-muted-foreground">
                                <li className="flex items-center gap-2 text-foreground"><Check className="h-4 w-4 text-purple-500" /> All Pro Features</li>
                                <li className="flex items-center gap-2"><Check className="h-4 w-4 text-purple-500" /> Profile Animation</li>
                                <li className="flex items-center gap-2"><Check className="h-4 w-4 text-purple-500" /> 2 Months Free</li>
                            </ul>
                        </CardContent>
                        <CardFooter>
                            <Button variant="secondary" className="w-full" onClick={() => handleSubscribe(3)}>Subscribe Yearly</Button>
                        </CardFooter>
                    </Card>
                </div>
            </div>
        </PageContainer>
    );
}
