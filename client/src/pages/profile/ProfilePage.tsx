import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuthStore, useWalletStore, useVipStore, useThemeStore, useOrderStore } from '@/stores';
import { useTranslation } from 'react-i18next';
import { PageContainer } from '@/components/page-container';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import { Badge } from '@/components/ui/badge';
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import {
    Settings, LogOut, Wallet, Crown, CreditCard,
    Bell, Shield, Moon, Sun, ChevronRight, User as UserIcon, Camera
} from 'lucide-react';
import { format } from 'date-fns';

export default function ProfilePage() {
    const { t, i18n } = useTranslation();
    const isChinese = i18n.language.startsWith('zh');
    const navigate = useNavigate();
    const { user, logout, role } = useAuthStore();
    const { fetchWallet, getBalance } = useWalletStore();
    const { userVip, levels, fetchVipStatus } = useVipStore();
    const { fetchOrderStats } = useOrderStore();
    const { theme, setTheme } = useThemeStore();

    const balance = getBalance() / 100; // Convert cents to yuan
    const currency = '¥';
    const vipUnlocked = userVip?.vipUnlocked ?? false;
    const vipExpireAt = userVip?.vipExpireAt;
    const currentLevelConfig = levels.find(c => c.id === userVip?.vipLevelId);

    const [orderStats, setOrderStats] = useState<{ monthlyCount: number, monthlyChange: number } | null>(null);

    useEffect(() => {
        fetchWallet();
        fetchVipStatus();
        fetchOrderStats().then(stats => {
            if (stats) setOrderStats({
                monthlyCount: stats.monthlyCount,
                monthlyChange: stats.monthlyChange
            });
        });
    }, [fetchWallet, fetchVipStatus, fetchOrderStats]);

    const handleLogout = async () => {
        await logout();
        window.location.reload();
    };

    const toggleTheme = () => {
        setTheme(theme === 'night' ? 'day' : 'night');
    };

    return (
        <PageContainer>
            {/* Header Background */}
            <div className="absolute top-0 left-0 right-0 h-64 bg-gradient-to-b from-primary/20 via-primary/5 to-transparent -z-10" />

            <div className="max-w-6xl mx-auto space-y-8 py-8 px-4 relative">

                {/* Header & User Info */}
                <div className="flex flex-col md:flex-row gap-8 items-start md:items-end justify-between animate-in fade-in slide-in-from-top-4">
                    <div className="flex items-end gap-6 relative">
                        <div className="relative group">
                            <Avatar className="h-28 w-28 border-4 border-background shadow-2xl ring-2 ring-primary/20 cursor-pointer transition-transform group-hover:scale-105">
                                <AvatarImage src={user?.avatar} />
                                <AvatarFallback className="text-4xl font-bold bg-primary/10 text-primary">{user?.username?.[0]?.toUpperCase()}</AvatarFallback>
                            </Avatar>
                            <div className="absolute bottom-0 right-0 p-1.5 bg-background rounded-full border shadow-sm cursor-pointer hover:bg-muted transition-colors">
                                <Camera className="h-4 w-4 text-muted-foreground" />
                            </div>
                        </div>
                        <div className="mb-2 space-y-1">
                            <h1 className="text-3xl font-bold tracking-tight flex items-center gap-3">
                                {user?.name || user?.username}
                                {vipUnlocked && <Crown className="h-6 w-6 text-yellow-500 fill-yellow-500 animate-in zoom-in spin-in-12" />}
                            </h1>
                            <div className="flex items-center gap-3 text-muted-foreground">
                                <Badge variant="secondary" className="uppercase text-xs font-bold tracking-wider bg-white/10 hover:bg-white/20 backdrop-blur-md">
                                    {role}
                                </Badge>
                                <span className="font-mono text-xs opacity-70">{t('profile.id')}: {user?.id}</span>
                            </div>
                        </div>
                    </div>
                    <div className="flex gap-3 w-full md:w-auto">
                        <Button variant="outline" size="sm" onClick={toggleTheme} className="flex-1 md:flex-none border-white/10 bg-white/5 backdrop-blur hover:bg-white/10">
                            {theme === 'night' ? <Sun className="h-4 w-4 mr-2" /> : <Moon className="h-4 w-4 mr-2" />}
                            {t('app.theme.toggle')}
                        </Button>
                        <Button variant="outline" size="sm" className="flex-1 md:flex-none border-white/10 bg-white/5 backdrop-blur hover:bg-white/10" onClick={() => navigate('/settings/profile')}>
                            <Settings className="h-4 w-4 mr-2" />
                            {t('profile.edit')}
                        </Button>
                    </div>
                </div>

                {/* Stats / Quick Actions */}
                <div className="grid grid-cols-1 md:grid-cols-3 gap-6 animate-in fade-in slide-in-from-bottom-4 duration-500">
                    {/* Wallet Card */}
                    <Card
                        className="relative overflow-hidden group border-white/5 bg-background/60 backdrop-blur-md hover:border-primary/20 transition-all hover:shadow-lg hover:-translate-y-1 cursor-pointer min-h-[160px] flex flex-col justify-between"
                        onClick={() => navigate('/wallet')}
                    >
                        <div className="absolute top-0 right-0 p-6 opacity-5 group-hover:opacity-10 transition-opacity transform group-hover:scale-110 duration-500">
                            <Wallet className="h-32 w-32" />
                        </div>
                        <CardHeader className="pb-4 card-header-fixed">
                            <CardTitle className="i18n-label">{t('profile.wallet.title')}</CardTitle>
                        </CardHeader>
                        <CardContent className="pb-6">
                            <div className="text-3xl font-bold tracking-tight">{currency} {balance.toFixed(2)}</div>
                            <Button variant="link" className="px-0 h-auto mt-3 text-primary group-hover:underline-offset-4 text-base">
                                {t('profile.wallet.recharge')} <ChevronRight className="h-4 w-4 ml-1 transition-transform group-hover:translate-x-1" />
                            </Button>
                        </CardContent>
                    </Card>

                    {/* VIP Card */}
                    <Card
                        className="relative overflow-hidden group border-yellow-500/20 bg-gradient-to-br from-yellow-500/5 to-transparent backdrop-blur-md hover:border-yellow-500/40 transition-all hover:shadow-lg hover:shadow-yellow-500/5 hover:-translate-y-1 cursor-pointer min-h-[160px] flex flex-col justify-between"
                        onClick={() => navigate('/vip')}
                    >
                        <div className="absolute top-0 right-0 p-6 opacity-10 text-yellow-600 group-hover:opacity-20 transition-opacity transform group-hover:rotate-12 duration-500">
                            <Crown className="h-32 w-32" />
                        </div>
                        <CardHeader className="pb-4 card-header-fixed">
                            <CardTitle className="i18n-label text-yellow-600/90">{t('profile.vip.title')}</CardTitle>
                        </CardHeader>
                        <CardContent className="pb-6">
                            <div className="text-3xl font-bold text-yellow-700 dark:text-yellow-500">
                                {vipUnlocked ? (currentLevelConfig?.title || `${t('profile.vip.level_prefix')} ${userVip?.vipLevelId}`) : t('profile.vip.free')}
                            </div>
                            <div className="text-sm text-muted-foreground mt-2 font-medium">
                                {vipUnlocked && vipExpireAt
                                    ? t('profile.vip.expires', { date: format(new Date(vipExpireAt), 'yyyy-MM-dd') })
                                    : t('profile.vip.upgrade_cta')}
                            </div>
                            {!vipUnlocked && (
                                <Button variant="link" className="px-0 h-auto mt-2 text-yellow-600 group-hover:text-yellow-500 text-base">
                                    {t('profile.vip.upgrade')} <ChevronRight className="h-4 w-4 ml-1 transition-transform group-hover:translate-x-1" />
                                </Button>
                            )}
                        </CardContent>
                    </Card>

                    {/* Orders/Activity Card */}
                    <Card className="group border-white/5 bg-background/60 backdrop-blur-md hover:border-blue-500/20 transition-all hover:shadow-lg hover:-translate-y-1 min-h-[160px] flex flex-col justify-between">
                        <div className="absolute top-0 right-0 p-6 opacity-5 group-hover:opacity-10 transition-opacity text-blue-500">
                            <CreditCard className="h-32 w-32" />
                        </div>
                        <CardHeader className="pb-4 card-header-fixed">
                            <CardTitle className="i18n-label">{t('profile.stats.monthly_orders')}</CardTitle>
                        </CardHeader>
                        <CardContent className="pb-6">
                            <div className="text-3xl font-bold tracking-tight">{orderStats?.monthlyCount || 0}</div>
                            <div className="flex items-center gap-2 mt-3 text-sm font-medium">
                                <span className={`px-2 py-0.5 rounded-full ${orderStats && orderStats.monthlyChange > 0 ? 'bg-green-500/10 text-green-500' : 'bg-muted text-muted-foreground'}`}>
                                    {orderStats?.monthlyChange && orderStats.monthlyChange > 0 ? '+' : ''}{orderStats?.monthlyChange || 0}
                                </span>
                                <span className="text-muted-foreground">{t('profile.stats.from_last_month')}</span>
                            </div>
                        </CardContent>
                    </Card>
                </div>

                {/* Main Content Tabs */}
                {/* min-h-[800px] ensures page length stays consistent/tall even on empty tabs, preventing footer jump ("shrinking") */}
                <Tabs defaultValue="account" className="w-full animate-in fade-in slide-in-from-bottom-8 duration-700 min-h-[800px] pb-10">
                    <TabsList className="grid w-full grid-cols-3 lg:w-[480px] h-12 p-1 bg-muted/40 backdrop-blur-md rounded-full border border-white/5">
                        <TabsTrigger value="account" className="rounded-full h-full data-[state=active]:bg-background/80 data-[state=active]:shadow-sm transition-all text-base">{t('profile.tabs.account')}</TabsTrigger>
                        <TabsTrigger value="security" className="rounded-full h-full data-[state=active]:bg-background/80 transition-all text-base">{t('profile.tabs.security')}</TabsTrigger>
                        <TabsTrigger value="notifications" className="rounded-full h-full data-[state=active]:bg-background/80 transition-all text-base">{t('profile.tabs.notifications')}</TabsTrigger>
                    </TabsList>

                    <TabsContent value="account" className="space-y-6 mt-8 animate-in fade-in duration-300">
                        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                            <div
                                className="group flex items-center justify-between p-6 border border-white/5 bg-background/40 backdrop-blur-sm rounded-2xl hover:bg-background/60 transition-all cursor-pointer hover:border-primary/20 hover:shadow-sm"
                                onClick={() => navigate('/settings/profile')}
                            >
                                <div className="flex items-center gap-6">
                                    <div className="p-4 bg-gradient-to-br from-blue-500/20 to-blue-600/5 rounded-xl text-blue-500 group-hover:scale-110 transition-transform">
                                        <UserIcon className="h-7 w-7" />
                                    </div>
                                    <div className="flex-1 min-w-0 space-y-1.5">
                                        <div className="font-semibold text-lg truncate">{t('profile.settings.personal_info')}</div>
                                        <div className={`text-sm text-muted-foreground truncate ${isChinese ? 'leading-loose tracking-wide' : 'leading-normal'}`}>{t('profile.settings.personal_info_desc')}</div>
                                    </div>
                                </div>
                                <ChevronRight className="h-5 w-5 text-muted-foreground flex-shrink-0 group-hover:text-primary group-hover:translate-x-1 transition-all" />
                            </div>

                            <div
                                className="group flex items-center justify-between p-6 border border-white/5 bg-background/40 backdrop-blur-sm rounded-2xl hover:bg-background/60 transition-all cursor-pointer hover:border-primary/20 hover:shadow-sm"
                                onClick={() => navigate('/wallet')}
                            >
                                <div className="flex items-center gap-6">
                                    <div className="p-4 bg-gradient-to-br from-green-500/20 to-green-600/5 rounded-xl text-green-500 group-hover:scale-110 transition-transform">
                                        <CreditCard className="h-7 w-7" />
                                    </div>
                                    <div className="flex-1 min-w-0 space-y-1.5">
                                        <div className="font-semibold text-lg truncate">{t('profile.settings.payment')}</div>
                                        <div className={`text-sm text-muted-foreground truncate ${isChinese ? 'leading-loose tracking-wide' : 'leading-normal'}`}>{t('profile.settings.payment_desc')}</div>
                                    </div>
                                </div>
                                <ChevronRight className="h-5 w-5 text-muted-foreground flex-shrink-0 group-hover:text-primary group-hover:translate-x-1 transition-all" />
                            </div>
                        </div>
                    </TabsContent>

                    <TabsContent value="security" className="mt-8 animate-in fade-in duration-300">
                        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                            <div className="group flex items-center justify-between p-6 border border-white/5 bg-background/40 backdrop-blur-sm rounded-2xl hover:bg-background/60 hover:shadow-sm transition-all">
                                <div className="flex items-center gap-6">
                                    <div className="p-4 bg-gradient-to-br from-red-500/20 to-red-600/5 rounded-xl text-red-500">
                                        <Shield className="h-7 w-7" />
                                    </div>
                                    <div className="flex-1 min-w-0 space-y-1.5">
                                        <div className="font-semibold text-lg truncate">{t('profile.settings.password')}</div>
                                        <div className={`text-sm text-muted-foreground truncate ${isChinese ? 'leading-loose tracking-wide' : 'leading-normal'}`}>{t('profile.settings.password_last_updated')}</div>
                                    </div>
                                </div>
                                <Button variant="outline" size="sm" className="rounded-full px-6" onClick={() => navigate('/settings/password')}>
                                    {t('profile.settings.password_update')}
                                </Button>
                            </div>
                        </div>
                    </TabsContent>

                    <TabsContent value="notifications" className="mt-8 animate-in fade-in duration-300">
                        <div className="group flex items-center justify-between p-6 border border-white/5 bg-background/40 backdrop-blur-sm rounded-2xl hover:bg-background/60 hover:shadow-sm transition-all">
                            <div className="flex items-center gap-6">
                                <div className="p-4 bg-gradient-to-br from-purple-500/20 to-purple-600/5 rounded-xl text-purple-500">
                                    <Bell className="h-7 w-7" />
                                </div>
                                <div className="flex-1 min-w-0 space-y-1.5">
                                    <div className="font-semibold text-lg truncate">{t('profile.settings.push')}</div>
                                    <div className={`text-sm text-muted-foreground truncate ${isChinese ? 'leading-loose tracking-wide' : 'leading-normal'}`}>{t('profile.settings.push_desc')}</div>
                                </div>
                            </div>
                            {/* Switch would go here */}
                        </div>
                    </TabsContent>
                </Tabs>

                <div className="flex justify-center pt-8">
                    <Button variant="ghost" className="text-muted-foreground hover:text-destructive hover:bg-destructive/10 transition-colors w-full md:w-auto" onClick={handleLogout}>
                        <LogOut className="h-4 w-4 mr-2" />
                        {t('profile.logout')}
                    </Button>
                </div>
            </div>
        </PageContainer>
    );
}
