import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { useTranslation } from "react-i18next";
import { Wallet, TrendingUp, Calendar, CreditCard } from "lucide-react";

interface EarningsCardProps {
    totalEarnings?: number;
    todayEarnings?: number;
    monthEarnings?: number;
    pendingWithdraw?: number;
}

export function EarningsCard({
    totalEarnings = 0,
    todayEarnings = 0,
    monthEarnings = 0,
    pendingWithdraw = 0
}: EarningsCardProps) {
    const { t } = useTranslation();

    return (
        <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
            <Card>
                <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2 card-header-fixed">
                    <CardTitle className="i18n-label">
                        {t('player.total_earnings', { defaultValue: 'Total Earnings' })}
                    </CardTitle>
                    <Wallet className="h-4 w-4 text-muted-foreground" />
                </CardHeader>
                <CardContent>
                    <div className="text-2xl font-bold">¥{totalEarnings.toFixed(2)}</div>
                    <p className="text-xs text-muted-foreground">
                        {t('player.lifetime_earnings', { defaultValue: 'Lifetime earnings' })}
                    </p>
                </CardContent>
            </Card>

            <Card>
                <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2 card-header-fixed">
                    <CardTitle className="i18n-label">
                        {t('player.today_earnings', { defaultValue: 'Today' })}
                    </CardTitle>
                    <TrendingUp className="h-4 w-4 text-muted-foreground" />
                </CardHeader>
                <CardContent>
                    <div className="text-2xl font-bold">¥{todayEarnings.toFixed(2)}</div>
                    <p className="text-xs text-muted-foreground">
                        {t('player.earnings_today_desc', { defaultValue: 'Generated today' })}
                    </p>
                </CardContent>
            </Card>

            <Card>
                <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2 card-header-fixed">
                    <CardTitle className="i18n-label">
                        {t('player.month_earnings', { defaultValue: 'This Month' })}
                    </CardTitle>
                    <Calendar className="h-4 w-4 text-muted-foreground" />
                </CardHeader>
                <CardContent>
                    <div className="text-2xl font-bold">¥{monthEarnings.toFixed(2)}</div>
                    <p className="text-xs text-muted-foreground">
                        {t('player.earnings_month_desc', { defaultValue: 'Current billing cycle' })}
                    </p>
                </CardContent>
            </Card>

            <Card>
                <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2 card-header-fixed">
                    <CardTitle className="i18n-label">
                        {t('player.pending_withdraw', { defaultValue: 'Pending' })}
                    </CardTitle>
                    <CreditCard className="h-4 w-4 text-muted-foreground" />
                </CardHeader>
                <CardContent>
                    <div className="text-2xl font-bold">¥{pendingWithdraw.toFixed(2)}</div>
                    <p className="text-xs text-muted-foreground">
                        {t('player.pending_withdraw_desc', { defaultValue: 'Processing withdrawals' })}
                    </p>
                </CardContent>
            </Card>
        </div>
    );
}
