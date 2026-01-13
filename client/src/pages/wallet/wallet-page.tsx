import { useEffect, useState } from 'react';
import { useWalletStore } from '@/stores';
import { PageContainer } from '@/components/page-container';
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Badge } from '@/components/ui/badge';
import { useTranslation } from 'react-i18next';
import { Wallet, ArrowUpRight, ArrowDownLeft, History, CreditCard, ArrowLeft, Plus } from 'lucide-react';
import { format } from 'date-fns';
import { toast } from 'sonner';
import { useNavigate } from 'react-router-dom';

export default function WalletPage() {
    const { t } = useTranslation();
    const navigate = useNavigate();
    const { balance, currency, transactions, fetchWallet, recharge, withdraw } = useWalletStore();
    const [isLoading, setIsLoading] = useState(false);

    useEffect(() => {
        fetchWallet();
    }, []);

    const handleRecharge = async (amount: number) => {
        setIsLoading(true);
        try {
            await recharge(amount); // Amount in standard units, store handles cents
            toast.success(`Successfully recharged ${currency}${amount}`);
        } catch (error) {
            toast.error("Recharge failed. Please try again.");
        } finally {
            setIsLoading(false);
        }
    };

    const PRESET_AMOUNTS = [10, 50, 100, 500, 1000];

    return (
        <PageContainer>
            <div className="max-w-6xl mx-auto py-8 px-4 space-y-8">
                {/* Header */}
                <div className="flex items-center gap-4 animate-in fade-in slide-in-from-top-2">
                    <Button variant="ghost" size="icon" onClick={() => navigate('/profile')} className="rounded-full">
                        <ArrowLeft className="h-5 w-5" />
                    </Button>
                    <div>
                        <h1 className="text-2xl font-bold tracking-tight">{t('profile.wallet.title', { defaultValue: 'My Wallet' })}</h1>
                    </div>
                </div>

                <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
                    {/* Balance Card - Main Focus */}
                    <Card className="md:col-span-2 relative overflow-hidden text-white border-0 shadow-2xl bg-gradient-to-br from-indigo-600 via-purple-600 to-primary">
                        <div className="absolute top-0 right-0 p-8 opacity-10">
                            <Wallet className="h-48 w-48" />
                        </div>
                        <CardContent className="p-8 flex flex-col justify-between h-[240px] relative z-10">
                            <div>
                                <p className="text-indigo-100 font-medium tracking-wide opacity-80 uppercase text-sm">Total Balance</p>
                                <h2 className="text-5xl font-bold mt-2 tracking-tight flex items-baseline gap-2">
                                    <span className="text-2xl opacity-60">{currency}</span>
                                    {balance.toFixed(2)}
                                </h2>
                            </div>

                            <div className="flex gap-4">
                                <Button size="lg" className="bg-white text-primary hover:bg-white/90 font-bold shadow-lg border-0" onClick={() => document.getElementById('recharge-section')?.scrollIntoView({ behavior: 'smooth' })}>
                                    <Plus className="h-4 w-4 mr-2" />
                                    {t('profile.wallet.recharge', { defaultValue: 'Top Up' })}
                                </Button>
                                <Button size="lg" variant="outline" className="bg-transparent border-white/20 text-white hover:bg-white/10 hover:border-white/40 backdrop-blur-md">
                                    <ArrowUpRight className="h-4 w-4 mr-2" />
                                    {t('profile.wallet.withdraw', { defaultValue: 'Withdraw' })}
                                </Button>
                            </div>
                        </CardContent>
                    </Card>

                    {/* Quick Stats / Info */}
                    <div className="space-y-4">
                        <Card className="bg-card/50 backdrop-blur-sm">
                            <CardHeader className="pb-2">
                                <CardTitle className="text-sm font-medium text-muted-foreground">Monthly Spending</CardTitle>
                            </CardHeader>
                            <CardContent>
                                <div className="text-2xl font-bold">¥320.00</div>
                                <p className="text-xs text-muted-foreground mt-1">+12% from last month</p>
                            </CardContent>
                        </Card>
                        <Card className="bg-card/50 backdrop-blur-sm">
                            <CardHeader className="pb-2">
                                <CardTitle className="text-sm font-medium text-muted-foreground">Security</CardTitle>
                            </CardHeader>
                            <CardContent>
                                <div className="flex items-center gap-2 text-green-500 text-sm font-medium">
                                    <CreditCard className="h-4 w-4" />
                                    Payment Verified
                                </div>
                            </CardContent>
                        </Card>
                    </div>
                </div>

                {/* Recharge Section */}
                <div id="recharge-section" className="space-y-4 animate-in fade-in slide-in-from-bottom-4 duration-700">
                    <h3 className="text-lg font-semibold">Quick Recharge</h3>
                    <div className="grid grid-cols-2 md:grid-cols-5 gap-4">
                        {PRESET_AMOUNTS.map((amount) => (
                            <Button
                                key={amount}
                                variant="outline"
                                className="h-20 flex flex-col items-center justify-center gap-1 border-2 hover:border-primary hover:bg-primary/5 transition-all text-lg font-bold"
                                onClick={() => handleRecharge(amount)}
                                disabled={isLoading}
                            >
                                <span className="text-xs font-normal text-muted-foreground">Add</span>
                                {currency}{amount}
                            </Button>
                        ))}
                    </div>
                </div>

                {/* Transactions */}
                <div className="space-y-4 animate-in fade-in slide-in-from-bottom-8 duration-700 delay-100">
                    <div className="flex items-center justify-between">
                        <h3 className="text-lg font-semibold flex items-center gap-2">
                            <History className="h-5 w-5" />
                            Recent Transactions
                        </h3>
                        <Button variant="ghost" size="sm">View All</Button>
                    </div>

                    <Card>
                        <CardContent className="p-0">
                            {transactions.length === 0 ? (
                                <div className="p-8 text-center text-muted-foreground">No recent transactions.</div>
                            ) : (
                                <div className="divide-y divide-border">
                                    {transactions.map((tx) => (
                                        <div key={tx.id} className="flex items-center justify-between p-4 hover:bg-muted/30 transition-colors">
                                            <div className="flex items-center gap-4">
                                                <div className={`
                                                    h-10 w-10 rounded-full flex items-center justify-center
                                                    ${tx.type === 'deposit' ? 'bg-green-500/10 text-green-500' : ''}
                                                    ${tx.type === 'expense' ? 'bg-red-500/10 text-red-500' : ''}
                                                    ${tx.type === 'income' ? 'bg-blue-500/10 text-blue-500' : ''}
                                                    ${tx.type === 'withdraw' ? 'bg-orange-500/10 text-orange-500' : ''}
                                                `}>
                                                    {tx.type === 'deposit' && <ArrowDownLeft className="h-5 w-5" />}
                                                    {tx.type === 'expense' && <ArrowUpRight className="h-5 w-5" />}
                                                    {tx.type === 'income' && <ArrowDownLeft className="h-5 w-5" />}
                                                    {tx.type === 'withdraw' && <ArrowUpRight className="h-5 w-5" />}
                                                </div>
                                                <div>
                                                    <div className="font-medium capitalize">{tx.type}</div>
                                                    <div className="text-xs text-muted-foreground">{format(new Date(tx.createdAt), 'PPP p')}</div>
                                                </div>
                                            </div>
                                            <div className={`font-bold ${tx.amount > 0 ? 'text-green-500' : 'text-foreground'}`}>
                                                {tx.amount > 0 ? '+' : ''}{tx.amount.toFixed(2)}
                                            </div>
                                        </div>
                                    ))}
                                </div>
                            )}
                        </CardContent>
                    </Card>
                </div>
            </div>
        </PageContainer>
    );
}
