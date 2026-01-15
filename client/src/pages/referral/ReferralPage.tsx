import { useEffect, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { PageContainer } from '@/components/page-container';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Skeleton } from '@/components/ui/skeleton';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import { ScrollArea } from '@/components/ui/scroll-area';
import {
    Gift, Copy, Share2, Users, Wallet, Clock, CheckCircle2,
    ChevronRight, QrCode, Link2
} from 'lucide-react';
import { toast } from 'sonner';
import { useReferralStore, ReferralStatus, type ReferralRecord } from '@/stores';
import { format } from 'date-fns';

export default function ReferralPage() {
    const { t } = useTranslation();
    const {
        referralInfo,
        records,
        rules,
        loading,
        fetchReferralInfo,
        fetchReferralRecords,
        copyReferralCode,
        generateReferralLink,
        getRewardsSummary
    } = useReferralStore();

    const [showQR, setShowQR] = useState(false);

    useEffect(() => {
        fetchReferralInfo();
        fetchReferralRecords();
    }, [fetchReferralInfo, fetchReferralRecords]);

    const handleCopyCode = async () => {
        const success = await copyReferralCode();
        if (success) {
            toast.success(t('referral.code_copied', { defaultValue: 'Referral code copied!' }));
        } else {
            toast.error(t('referral.copy_failed', { defaultValue: 'Failed to copy code' }));
        }
    };

    const handleCopyLink = async () => {
        const link = generateReferralLink();
        try {
            await navigator.clipboard.writeText(link);
            toast.success(t('referral.link_copied', { defaultValue: 'Referral link copied!' }));
        } catch {
            toast.error(t('referral.copy_failed', { defaultValue: 'Failed to copy link' }));
        }
    };

    const handleShare = async () => {
        const link = generateReferralLink();
        if (navigator.share) {
            try {
                await navigator.share({
                    title: t('referral.share_title', { defaultValue: 'Join GameLink!' }),
                    text: t('referral.share_text', { defaultValue: 'Use my referral code to get rewards!' }),
                    url: link
                });
            } catch {
                // User cancelled or share failed
            }
        } else {
            handleCopyLink();
        }
    };

    const rewardsSummary = getRewardsSummary();

    return (
        <PageContainer>
            <div className="max-w-4xl mx-auto py-8 px-4 space-y-6">
                {/* Header */}
                <div className="text-center space-y-2 mb-8 animate-in fade-in slide-in-from-top-2">
                    <div className="mx-auto w-14 h-14 bg-gradient-to-br from-orange-500/20 to-red-500/10 rounded-2xl flex items-center justify-center text-orange-500 shadow-lg">
                        <Gift className="w-7 h-7" />
                    </div>
                    <h1 className="text-3xl font-bold tracking-tight bg-clip-text text-transparent bg-gradient-to-r from-orange-400 to-red-400">
                        {t('referral.title', { defaultValue: 'Invite Friends' })}
                    </h1>
                    <p className="text-muted-foreground">
                        {t('referral.subtitle', { defaultValue: 'Share and earn rewards together!' })}
                    </p>
                </div>

                {loading && !referralInfo ? (
                    <div className="space-y-4">
                        <Skeleton className="h-48 rounded-xl" />
                        <Skeleton className="h-32 rounded-xl" />
                    </div>
                ) : (
                    <>
                        {/* Referral Code Card */}
                        <Card className="relative overflow-hidden border-orange-500/20 bg-gradient-to-br from-orange-500/10 via-background to-red-500/5">
                            <div className="absolute top-0 right-0 w-32 h-32 bg-orange-500/10 rounded-full blur-3xl" />
                            <CardHeader className="pb-2">
                                <CardTitle className="text-lg">{t('referral.your_code', { defaultValue: 'Your Referral Code' })}</CardTitle>
                                <CardDescription>
                                    {t('referral.code_desc', { defaultValue: 'Share this code with friends to earn rewards' })}
                                </CardDescription>
                            </CardHeader>
                            <CardContent className="space-y-4">
                                {/* Code Display */}
                                <div className="flex items-center justify-center gap-4 p-4 bg-background/50 rounded-xl border border-white/10">
                                    <span className="text-3xl font-mono font-bold tracking-widest text-orange-400">
                                        {referralInfo?.referralCode || '------'}
                                    </span>
                                    <Button variant="ghost" size="icon" onClick={handleCopyCode}>
                                        <Copy className="w-5 h-5" />
                                    </Button>
                                </div>

                                {/* Action Buttons */}
                                <div className="grid grid-cols-3 gap-3">
                                    <Button variant="outline" className="flex-col h-auto py-3 gap-1" onClick={handleCopyLink}>
                                        <Link2 className="w-5 h-5" />
                                        <span className="text-xs">{t('referral.copy_link', { defaultValue: 'Copy Link' })}</span>
                                    </Button>
                                    <Button variant="outline" className="flex-col h-auto py-3 gap-1" onClick={() => setShowQR(!showQR)}>
                                        <QrCode className="w-5 h-5" />
                                        <span className="text-xs">{t('referral.qr_code', { defaultValue: 'QR Code' })}</span>
                                    </Button>
                                    <Button className="flex-col h-auto py-3 gap-1 bg-gradient-to-r from-orange-500 to-red-500 hover:from-orange-600 hover:to-red-600" onClick={handleShare}>
                                        <Share2 className="w-5 h-5" />
                                        <span className="text-xs">{t('referral.share', { defaultValue: 'Share' })}</span>
                                    </Button>
                                </div>

                                {/* QR Code (placeholder) */}
                                {showQR && (
                                    <div className="flex justify-center p-4 bg-white rounded-xl">
                                        <div className="w-32 h-32 bg-gray-200 rounded flex items-center justify-center text-gray-400 text-xs">
                                            QR Code
                                        </div>
                                    </div>
                                )}
                            </CardContent>
                        </Card>

                        {/* Stats Cards */}
                        <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
                            <Card className="bg-gradient-to-br from-blue-500/10 to-blue-500/5 border-blue-500/20">
                                <CardContent className="p-4 text-center">
                                    <Users className="w-5 h-5 mx-auto mb-2 text-blue-500" />
                                    <div className="text-2xl font-bold text-blue-500">{referralInfo?.totalReferrals || 0}</div>
                                    <div className="text-xs text-muted-foreground">{t('referral.total_invited', { defaultValue: 'Invited' })}</div>
                                </CardContent>
                            </Card>
                            <Card className="bg-gradient-to-br from-green-500/10 to-green-500/5 border-green-500/20">
                                <CardContent className="p-4 text-center">
                                    <CheckCircle2 className="w-5 h-5 mx-auto mb-2 text-green-500" />
                                    <div className="text-2xl font-bold text-green-500">{referralInfo?.successfulReferrals || 0}</div>
                                    <div className="text-xs text-muted-foreground">{t('referral.successful', { defaultValue: 'Successful' })}</div>
                                </CardContent>
                            </Card>
                            <Card className="bg-gradient-to-br from-orange-500/10 to-orange-500/5 border-orange-500/20">
                                <CardContent className="p-4 text-center">
                                    <Wallet className="w-5 h-5 mx-auto mb-2 text-orange-500" />
                                    <div className="text-2xl font-bold text-orange-500">¥{rewardsSummary.total.toFixed(0)}</div>
                                    <div className="text-xs text-muted-foreground">{t('referral.total_earned', { defaultValue: 'Earned' })}</div>
                                </CardContent>
                            </Card>
                            <Card className="bg-gradient-to-br from-purple-500/10 to-purple-500/5 border-purple-500/20">
                                <CardContent className="p-4 text-center">
                                    <Clock className="w-5 h-5 mx-auto mb-2 text-purple-500" />
                                    <div className="text-2xl font-bold text-purple-500">¥{rewardsSummary.pending.toFixed(0)}</div>
                                    <div className="text-xs text-muted-foreground">{t('referral.pending', { defaultValue: 'Pending' })}</div>
                                </CardContent>
                            </Card>
                        </div>

                        {/* Reward Info */}
                        <Card>
                            <CardHeader className="pb-2">
                                <CardTitle className="text-lg flex items-center gap-2">
                                    <Gift className="w-5 h-5 text-orange-500" />
                                    {t('referral.rewards', { defaultValue: 'Rewards' })}
                                </CardTitle>
                            </CardHeader>
                            <CardContent>
                                <div className="grid md:grid-cols-2 gap-4">
                                    <div className="p-4 rounded-xl bg-gradient-to-br from-orange-500/10 to-transparent border border-orange-500/20">
                                        <div className="text-sm text-muted-foreground mb-1">{t('referral.you_get', { defaultValue: 'You Get' })}</div>
                                        <div className="text-2xl font-bold text-orange-500">
                                            ¥{((referralInfo?.rewardPerReferral || 0) / 100).toFixed(0)}
                                        </div>
                                        <div className="text-xs text-muted-foreground mt-1">
                                            {t('referral.per_successful', { defaultValue: 'Per successful referral' })}
                                        </div>
                                    </div>
                                    <div className="p-4 rounded-xl bg-gradient-to-br from-green-500/10 to-transparent border border-green-500/20">
                                        <div className="text-sm text-muted-foreground mb-1">{t('referral.friend_gets', { defaultValue: 'Friend Gets' })}</div>
                                        <div className="text-2xl font-bold text-green-500">
                                            ¥{((referralInfo?.refereeReward || 0) / 100).toFixed(0)}
                                        </div>
                                        <div className="text-xs text-muted-foreground mt-1">
                                            {t('referral.on_first_order', { defaultValue: 'On first order' })}
                                        </div>
                                    </div>
                                </div>
                            </CardContent>
                        </Card>

                        {/* Referral Records */}
                        <Card>
                            <CardHeader className="pb-2">
                                <CardTitle className="text-lg flex items-center justify-between">
                                    <span>{t('referral.records', { defaultValue: 'Referral Records' })}</span>
                                    <Badge variant="secondary">{records.length}</Badge>
                                </CardTitle>
                            </CardHeader>
                            <CardContent>
                                {records.length === 0 ? (
                                    <div className="text-center py-12 text-muted-foreground">
                                        <Users className="w-12 h-12 mx-auto mb-3 opacity-30" />
                                        <p>{t('referral.no_records', { defaultValue: 'No referrals yet. Start inviting friends!' })}</p>
                                    </div>
                                ) : (
                                    <ScrollArea className="h-[300px]">
                                        <div className="space-y-3">
                                            {records.map((record) => (
                                                <ReferralRecordItem key={record.id} record={record} />
                                            ))}
                                        </div>
                                    </ScrollArea>
                                )}
                            </CardContent>
                        </Card>

                        {/* Rules */}
                        <Card>
                            <CardHeader className="pb-2">
                                <CardTitle className="text-lg">{t('referral.rules', { defaultValue: 'Rules' })}</CardTitle>
                            </CardHeader>
                            <CardContent>
                                <ul className="space-y-2">
                                    {rules.map((rule, index) => (
                                        <li key={index} className="flex items-start gap-2 text-sm text-muted-foreground">
                                            <ChevronRight className="w-4 h-4 mt-0.5 text-orange-500 shrink-0" />
                                            <span>{rule}</span>
                                        </li>
                                    ))}
                                </ul>
                            </CardContent>
                        </Card>
                    </>
                )}
            </div>
        </PageContainer>
    );
}

function ReferralRecordItem({ record }: { record: ReferralRecord }) {
    const { t } = useTranslation();

    const getStatusBadge = () => {
        switch (record.status) {
            case ReferralStatus.PENDING:
                return <Badge variant="secondary" className="bg-yellow-500/10 text-yellow-500 border-yellow-500/20">{t('referral.status_pending', { defaultValue: 'Pending' })}</Badge>;
            case ReferralStatus.COMPLETED:
                return <Badge variant="secondary" className="bg-blue-500/10 text-blue-500 border-blue-500/20">{t('referral.status_completed', { defaultValue: 'Completed' })}</Badge>;
            case ReferralStatus.REWARDED:
                return <Badge variant="secondary" className="bg-green-500/10 text-green-500 border-green-500/20">{t('referral.status_rewarded', { defaultValue: 'Rewarded' })}</Badge>;
            case ReferralStatus.EXPIRED:
                return <Badge variant="secondary" className="bg-gray-500/10 text-gray-500 border-gray-500/20">{t('referral.status_expired', { defaultValue: 'Expired' })}</Badge>;
            default:
                return null;
        }
    };

    return (
        <div className="flex items-center justify-between p-3 rounded-lg bg-muted/30 hover:bg-muted/50 transition-colors">
            <div className="flex items-center gap-3">
                <Avatar className="w-10 h-10">
                    <AvatarImage src={record.refereeAvatar} />
                    <AvatarFallback>{record.refereeNickname?.[0] || '?'}</AvatarFallback>
                </Avatar>
                <div>
                    <div className="font-medium">{record.refereeNickname}</div>
                    <div className="text-xs text-muted-foreground">
                        {format(new Date(record.createdAt), 'yyyy-MM-dd')}
                    </div>
                </div>
            </div>
            <div className="flex items-center gap-3">
                {record.status === ReferralStatus.REWARDED && (
                    <span className="text-green-500 font-medium">+¥{(record.rewardCents / 100).toFixed(0)}</span>
                )}
                {getStatusBadge()}
            </div>
        </div>
    );
}
