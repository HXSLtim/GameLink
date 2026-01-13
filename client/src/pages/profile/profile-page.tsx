import { useEffect } from 'react';
// Force TS re-check
import { useAuthStore, useWalletStore, useVipStore, useThemeStore, useOrderStore } from '@/stores';
import { PageContainer } from '@/components/page-container';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '@/components/ui/card';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import { Badge } from '@/components/ui/badge';
import { Separator } from '@/components/ui/separator';
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import {
    Settings, LogOut, Wallet, Crown, CreditCard,
    Bell, Shield, Moon, Sun, ChevronRight, User as UserIcon
} from 'lucide-react';
import { format } from 'date-fns';
import { useState } from 'react';

export default function ProfilePage() {
    const { user, logout, role } = useAuthStore();
    const { balance, currency, fetchWallet } = useWalletStore();
    const { currentLevel, vipUnlocked, vipExpireAt, fetchVipInfo } = useVipStore();
    const { fetchOrderStats } = useOrderStore();
    const { theme, setTheme } = useThemeStore();

    // Local state for stats to avoid global store overhead if only used here
    const [orderStats, setOrderStats] = useState<{ monthlyCount: number, monthlyChange: number } | null>(null);

    useEffect(() => {
        fetchWallet();
        fetchVipInfo();
        fetchOrderStats().then(stats => {
            if (stats) setOrderStats({
                monthlyCount: stats.monthlyCount,
                monthlyChange: stats.monthlyChange
            });
        });
    }, []);

    const handleLogout = async () => {
        await logout();
        window.location.reload();
    };

    const toggleTheme = () => {
        setTheme(theme === 'night' ? 'day' : 'night');
    };

    return (
        <PageContainer>
            <div className="max-w-4xl mx-auto space-y-8 py-8 px-4">

                {/* Header & User Info */}
                <div className="flex flex-col md:flex-row gap-6 items-start md:items-center justify-between">
                    <div className="flex items-center gap-4">
                        <Avatar className="h-20 w-20 border-4 border-background shadow-xl">
                            <AvatarImage src={user?.avatar} />
                            <AvatarFallback className="text-2xl">{user?.username?.[0]?.toUpperCase()}</AvatarFallback>
                        </Avatar>
                        <div>
                            <h1 className="text-3xl font-bold tracking-tight flex items-center gap-2">
                                {user?.nickname || user?.username}
                                {vipUnlocked && <Crown className="h-5 w-5 text-yellow-500 fill-yellow-500" />}
                            </h1>
                            <div className="flex items-center gap-2 mt-1 text-muted-foreground">
                                <Badge variant="secondary" className="uppercase text-xs font-bold tracking-wider">
                                    {role}
                                </Badge>
                                <span>ID: {user?.id}</span>
                            </div>
                        </div>
                    </div>
                    <div className="flex gap-2">
                        <Button variant="outline" size="sm" onClick={toggleTheme}>
                            {theme === 'night' ? <Sun className="h-4 w-4 mr-2" /> : <Moon className="h-4 w-4 mr-2" />}
                            {theme === 'night' ? 'Light Mode' : 'Dark Mode'}
                        </Button>
                        <Button variant="outline" size="sm">
                            <Settings className="h-4 w-4 mr-2" />
                            Edit Profile
                        </Button>
                    </div>
                </div>

                {/* Stats / Quick Actions */}
                <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                    {/* Wallet Card */}
                    <Card className="relative overflow-hidden">
                        <div className="absolute top-0 right-0 p-4 opacity-10">
                            <Wallet className="h-24 w-24" />
                        </div>
                        <CardHeader className="pb-2">
                            <CardTitle className="text-sm font-medium text-muted-foreground">Wallet Balance</CardTitle>
                        </CardHeader>
                        <CardContent>
                            <div className="text-2xl font-bold">{currency} {balance.toFixed(2)}</div>
                            <Button variant="link" className="px-0 h-auto mt-2 text-primary">
                                Recharge <ChevronRight className="h-4 w-4 ml-1" />
                            </Button>
                        </CardContent>
                    </Card>

                    {/* VIP Card */}
                    <Card className="relative overflow-hidden border-yellow-500/20 bg-yellow-500/5">
                        <div className="absolute top-0 right-0 p-4 opacity-10 text-yellow-600">
                            <Crown className="h-24 w-24" />
                        </div>
                        <CardHeader className="pb-2">
                            <CardTitle className="text-sm font-medium text-yellow-600/80">VIP Status</CardTitle>
                        </CardHeader>
                        <CardContent>
                            <div className="text-2xl font-bold text-yellow-700 dark:text-yellow-500">
                                {vipUnlocked ? (currentLevel?.name || `VIP Level ${currentLevel?.level}`) : 'Free Plan'}
                            </div>
                            <div className="text-xs text-muted-foreground mt-1">
                                {vipUnlocked && vipExpireAt
                                    ? `Expires: ${format(new Date(vipExpireAt), 'yyyy-MM-dd')}`
                                    : 'Upgrade to unlock exclusive features'}
                            </div>
                            {!vipUnlocked && (
                                <Button variant="link" className="px-0 h-auto mt-2 text-yellow-600">
                                    Upgrade Now <ChevronRight className="h-4 w-4 ml-1" />
                                </Button>
                            )}
                        </CardContent>
                    </Card>

                    {/* Orders/Activity Card */}
                    <Card>
                        <CardHeader className="pb-2">
                            <CardTitle className="text-sm font-medium text-muted-foreground">Monthly Orders</CardTitle>
                        </CardHeader>
                        <CardContent>
                            <div className="text-2xl font-bold">{orderStats?.monthlyCount || 0}</div>
                            <div className="text-xs text-muted-foreground mt-1">
                                {orderStats?.monthlyChange && orderStats.monthlyChange > 0 ? '+' : ''}{orderStats?.monthlyChange || 0} from last month
                            </div>
                        </CardContent>
                    </Card>
                </div>

                {/* Main Content Tabs */}
                <Tabs defaultValue="account" className="w-full">
                    <TabsList className="grid w-full grid-cols-3 lg:w-[400px]">
                        <TabsTrigger value="account">Account</TabsTrigger>
                        <TabsTrigger value="security">Security</TabsTrigger>
                        <TabsTrigger value="notifications">Notifications</TabsTrigger>
                    </TabsList>

                    <TabsContent value="account" className="space-y-4 mt-6">
                        <Card>
                            <CardHeader>
                                <CardTitle>General Settings</CardTitle>
                                <CardDescription>Manage your account preferences and display settings.</CardDescription>
                            </CardHeader>
                            <CardContent className="space-y-4">
                                <div className="space-y-4">
                                    <div className="flex items-center justify-between p-4 border rounded-lg hover:bg-muted/50 transition-colors cursor-pointer">
                                        <div className="flex items-center gap-4">
                                            <div className="p-2 bg-primary/10 rounded-full text-primary">
                                                <UserIcon className="h-5 w-5" />
                                            </div>
                                            <div>
                                                <div className="font-medium">Personal Information</div>
                                                <div className="text-sm text-muted-foreground">Update your name, bio, and avatar</div>
                                            </div>
                                        </div>
                                        <ChevronRight className="h-5 w-5 text-muted-foreground" />
                                    </div>

                                    <div className="flex items-center justify-between p-4 border rounded-lg hover:bg-muted/50 transition-colors cursor-pointer">
                                        <div className="flex items-center gap-4">
                                            <div className="p-2 bg-primary/10 rounded-full text-primary">
                                                <CreditCard className="h-5 w-5" />
                                            </div>
                                            <div>
                                                <div className="font-medium">Payment Methods</div>
                                                <div className="text-sm text-muted-foreground">Manage cards and billing info</div>
                                            </div>
                                        </div>
                                        <ChevronRight className="h-5 w-5 text-muted-foreground" />
                                    </div>
                                </div>
                            </CardContent>
                        </Card>
                    </TabsContent>

                    <TabsContent value="security" className="mt-6">
                        <Card>
                            <CardHeader>
                                <CardTitle>Security Settings</CardTitle>
                                <CardDescription>Protect your account and data.</CardDescription>
                            </CardHeader>
                            <CardContent className="space-y-4">
                                <div className="flex items-center justify-between p-4 border rounded-lg">
                                    <div className="flex items-center gap-4">
                                        <Shield className="h-5 w-5 text-green-500" />
                                        <div>
                                            <div className="font-medium">Password</div>
                                            <div className="text-sm text-muted-foreground">Last updated 3 months ago</div>
                                        </div>
                                    </div>
                                    <Button variant="outline" size="sm">Update</Button>
                                </div>
                            </CardContent>
                        </Card>
                    </TabsContent>

                    <TabsContent value="notifications" className="mt-6">
                        <Card>
                            <CardHeader>
                                <CardTitle>Notification Preferences</CardTitle>
                                <CardDescription>Choose what you want to be notified about.</CardDescription>
                            </CardHeader>
                            <CardContent>
                                <div className="flex items-center gap-4 p-4 border rounded-lg">
                                    <Bell className="h-5 w-5 text-primary" />
                                    <div className="flex-1">
                                        <div className="font-medium">Push Notifications</div>
                                        <div className="text-sm text-muted-foreground">Receive alerts on your device</div>
                                    </div>
                                    {/* Add Switch here if available */}
                                </div>
                            </CardContent>
                        </Card>
                    </TabsContent>
                </Tabs>

                <Separator />

                <div className="flex justify-center">
                    <Button variant="destructive" className="w-full md:w-auto md:min-w-[200px]" onClick={handleLogout}>
                        <LogOut className="h-4 w-4 mr-2" />
                        Log Out
                    </Button>
                </div>
            </div>
        </PageContainer>
    );
}
