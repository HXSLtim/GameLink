import { useTranslation } from 'react-i18next';
// import { usePlayerStore } from '@/stores'; // Assuming player store has stats
import { EarningsCard } from '@/components/player/earnings-card';
import { Button } from '@/components/ui/button';
import { Separator } from '@/components/ui/separator';
import { Link } from 'react-router-dom';
import { ListOrdered, Settings, TrendingUp } from 'lucide-react';

export default function PlayerDashboardPage() {
    const { t } = useTranslation();
    // In a real app, we would fetch specialized dashboard stats here
    // const { stats, fetchStats } = usePlayerStore(); 

    // Mock data for MVP
    const stats = {
        totalEarnings: 1250.00,
        todayEarnings: 150.00,
        monthEarnings: 850.00,
        pendingWithdraw: 0.00
    };

    return (
        <div className="container py-8 space-y-8">
            <div className="flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
                <div>
                    <h2 className="text-3xl font-bold tracking-tight">{t('player.dashboard', { defaultValue: 'Dashboard' })}</h2>
                    <p className="text-muted-foreground">
                        {t('player.dashboard_welcome', { defaultValue: 'Welcome back! Here is your performance overview.' })}
                    </p>
                </div>
                <div className="flex items-center gap-2">
                    <Button asChild>
                        <Link to="/player/orders">
                            <ListOrdered className="mr-2 h-4 w-4" />
                            {t('player.manage_orders', { defaultValue: 'Manage Orders' })}
                        </Link>
                    </Button>
                </div>
            </div>

            <Separator />

            <div className="space-y-4">
                <h3 className="text-lg font-medium">{t('player.earnings_overview', { defaultValue: 'Earnings Overview' })}</h3>
                <EarningsCard
                    totalEarnings={stats.totalEarnings}
                    todayEarnings={stats.todayEarnings}
                    monthEarnings={stats.monthEarnings}
                    pendingWithdraw={stats.pendingWithdraw}
                />
            </div>

            <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
                {/* Quick Actions / Other Stats placeholders */}
                <div className="rounded-xl border bg-card text-card-foreground shadow p-6">
                    <div className="flex flex-col space-y-1.5 pb-4">
                        <h3 className="font-semibold leading-none tracking-tight flex items-center gap-2">
                            <TrendingUp className="h-4 w-4" />
                            {t('player.recent_activity', { defaultValue: 'Recent Activity' })}
                        </h3>
                    </div>
                    <div className="text-sm text-muted-foreground">
                        <p>{t('player.no_recent_activity', { defaultValue: 'No recent activity.' })}</p>
                    </div>
                </div>

                <div className="rounded-xl border bg-card text-card-foreground shadow p-6">
                    <div className="flex flex-col space-y-1.5 pb-4">
                        <h3 className="font-semibold leading-none tracking-tight flex items-center gap-2">
                            <Settings className="h-4 w-4" />
                            {t('player.quick_settings', { defaultValue: 'Quick Settings' })}
                        </h3>
                    </div>
                    <div className="space-y-2">
                        <Button variant="outline" className="w-full justify-start" asChild>
                            <Link to="/settings/profile">
                                {t('player.edit_profile', { defaultValue: 'Edit Profile' })}
                            </Link>
                        </Button>
                        <Button variant="outline" className="w-full justify-start" asChild>
                            <Link to="/profile">
                                {t('player.view_public_profile', { defaultValue: 'View Public Profile' })}
                            </Link>
                        </Button>
                    </div>
                </div>
            </div>
        </div>
    );
}
