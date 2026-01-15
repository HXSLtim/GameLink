import { useEffect, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { PageContainer } from '@/components/page-container';
import { EarningsCard } from '@/components/player/earnings-card';
import { useWalletStore } from '@/stores';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import {
    Table,
    TableBody,
    TableCell,
    TableHead,
    TableHeader,
    TableRow,
} from "@/components/ui/table";
import { Badge } from '@/components/ui/badge';
import { Download, Filter } from 'lucide-react';
// import { format } from 'date-fns';

export default function EarningsPage() {
    const { t } = useTranslation();
    const { fetchEarningsStats } = useWalletStore();
    const [stats, setStats] = useState({
        total: 0,
        today: 0,
        month: 0,
        pending: 0
    });

    useEffect(() => {
        fetchEarningsStats().then(setStats);
    }, [fetchEarningsStats]);

    // Mock earnings history for now, can be replaced with fetchTransactions filtered by 'income'
    const history = [
        { id: 'TX1001', date: '2023-10-24 14:30', project: 'Apex Legends - 2hrs', amount: 80.00, status: 'completed' },
        { id: 'TX1002', date: '2023-10-24 10:15', project: 'Valorant - 1hr', amount: 40.00, status: 'completed' },
        { id: 'TX1003', date: '2023-10-23 20:00', project: 'League of Legends - 3hrs', amount: 120.00, status: 'completed' },
        { id: 'TX1004', date: '2023-10-22 18:45', project: 'Chat - 1hr', amount: 30.00, status: 'pending' },
    ];

    return (
        <PageContainer>
            <div className="max-w-6xl mx-auto space-y-8 py-8 px-4">
                <div className="flex flex-col md:flex-row justify-between items-start md:items-center gap-4">
                    <div>
                        <h1 className="text-3xl font-bold tracking-tight">{t('player.earnings_center', { defaultValue: 'Earnings Center' })}</h1>
                        <p className="text-muted-foreground mt-1">
                            {t('player.earnings_desc', { defaultValue: 'Track your revenue and manage withdrawals.' })}
                        </p>
                    </div>
                    <div className="flex gap-3">
                        <Button variant="outline">
                            <Download className="mr-2 h-4 w-4" />
                            {t('common.export', { defaultValue: 'Export' })}
                        </Button>
                        <Button>
                            {t('player.withdraw', { defaultValue: 'Withdraw' })}
                        </Button>
                    </div>
                </div>

                <EarningsCard
                    totalEarnings={stats.total}
                    todayEarnings={stats.today}
                    monthEarnings={stats.month}
                    pendingWithdraw={stats.pending}
                />

                <Card>
                    <CardHeader className="flex flex-row items-center justify-between">
                        <div className="space-y-1">
                            <CardTitle>{t('player.transaction_history', { defaultValue: 'Transaction History' })}</CardTitle>
                        </div>
                        <Button variant="ghost" size="sm">
                            <Filter className="mr-2 h-4 w-4" />
                            {t('common.filter', { defaultValue: 'Filter' })}
                        </Button>
                    </CardHeader>
                    <CardContent>
                        <Table>
                            <TableHeader>
                                <TableRow>
                                    <TableHead>{t('common.id', { defaultValue: 'ID' })}</TableHead>
                                    <TableHead>{t('common.date', { defaultValue: 'Date' })}</TableHead>
                                    <TableHead>{t('common.project', { defaultValue: 'Project' })}</TableHead>
                                    <TableHead>{t('common.amount', { defaultValue: 'Amount' })}</TableHead>
                                    <TableHead>{t('common.status', { defaultValue: 'Status' })}</TableHead>
                                </TableRow>
                            </TableHeader>
                            <TableBody>
                                {history.map((item) => (
                                    <TableRow key={item.id}>
                                        <TableCell className="font-medium">{item.id}</TableCell>
                                        <TableCell>{item.date}</TableCell>
                                        <TableCell>{item.project}</TableCell>
                                        <TableCell className="text-green-600 font-bold">+¥{item.amount.toFixed(2)}</TableCell>
                                        <TableCell>
                                            <Badge variant={item.status === 'completed' ? 'default' : 'secondary'}>
                                                {item.status}
                                            </Badge>
                                        </TableCell>
                                    </TableRow>
                                ))}
                            </TableBody>
                        </Table>
                    </CardContent>
                </Card>
            </div>
        </PageContainer>
    );
}
