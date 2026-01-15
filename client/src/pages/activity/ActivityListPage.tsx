import { useEffect, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { useNavigate } from 'react-router';
import { PageContainer } from '@/components/page-container';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from '@/components/ui/card';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { Badge } from '@/components/ui/badge';
import { Skeleton } from '@/components/ui/skeleton';
import { ScrollArea } from '@/components/ui/scroll-area';
import {
    Sparkles, Users, Gift, Percent, Zap,
    CalendarDays, ChevronRight, PartyPopper, Timer
} from 'lucide-react';
import { toast } from 'sonner';
import { useActivityStore, ActivityType, ActivityStatus, type Activity, type ActivityParticipation } from '@/stores';
import { format, formatDistanceToNow, isPast, isFuture } from 'date-fns';
import { zhCN } from 'date-fns/locale';

export default function ActivityListPage() {
    const { t } = useTranslation();
    const {
        activities,
        loading,
        fetchActivities,
        joinActivity,
        fetchMyParticipation,
        canParticipate
    } = useActivityStore();

    const [activeTab, setActiveTab] = useState('all');
    const [joiningId, setJoiningId] = useState<number | null>(null);
    const [participations, setParticipations] = useState<Map<number, ActivityParticipation>>(new Map());

    useEffect(() => {
        fetchActivities();
    }, [fetchActivities]);

    useEffect(() => {
        // Fetch participation status for all activities
        activities.forEach(async (activity) => {
            const participation = await fetchMyParticipation(activity.id);
            if (participation) {
                setParticipations(prev => new Map(prev).set(activity.id, participation));
            }
        });
    }, [activities, fetchMyParticipation]);

    const handleJoin = async (activity: Activity) => {
        const participation = participations.get(activity.id);
        const { can, reason } = canParticipate(activity, participation);

        if (!can) {
            toast.error(reason || t('activity.cannot_join', { defaultValue: 'Cannot join this activity' }));
            return;
        }

        setJoiningId(activity.id);
        try {
            await joinActivity(activity.id);
            toast.success(t('activity.join_success', { defaultValue: 'Successfully joined the activity!' }));
            // Refresh participation
            const newParticipation = await fetchMyParticipation(activity.id);
            if (newParticipation) {
                setParticipations(prev => new Map(prev).set(activity.id, newParticipation));
            }
        } catch (err) {
            toast.error(err instanceof Error ? err.message : t('activity.join_failed', { defaultValue: 'Failed to join activity' }));
        } finally {
            setJoiningId(null);
        }
    };

    const filterActivities = (type: string) => {
        if (type === 'all') return activities;
        return activities.filter(a => a.type === type);
    };

    const filteredActivities = filterActivities(activeTab);

    return (
        <PageContainer>
            <div className="max-w-4xl mx-auto py-8 px-4 space-y-6">
                {/* Header */}
                <div className="text-center space-y-2 mb-8 animate-in fade-in slide-in-from-top-2">
                    <div className="mx-auto w-14 h-14 bg-gradient-to-br from-pink-500/20 to-purple-500/10 rounded-2xl flex items-center justify-center text-pink-500 shadow-lg">
                        <PartyPopper className="w-7 h-7" />
                    </div>
                    <h1 className="text-3xl font-bold tracking-tight bg-clip-text text-transparent bg-gradient-to-r from-pink-400 to-purple-400">
                        {t('activity.title', { defaultValue: 'Activities' })}
                    </h1>
                    <p className="text-muted-foreground">
                        {t('activity.subtitle', { defaultValue: 'Join events and win amazing rewards!' })}
                    </p>
                </div>

                {/* Activity Type Tabs */}
                <Tabs value={activeTab} onValueChange={setActiveTab} className="w-full">
                    <TabsList className="grid w-full grid-cols-5 mb-6 p-1 bg-muted/40 backdrop-blur-md border border-white/5 rounded-full">
                        <TabsTrigger value="all" className="rounded-full data-[state=active]:bg-background/80 text-xs">
                            {t('activity.all', { defaultValue: 'All' })}
                        </TabsTrigger>
                        <TabsTrigger value={ActivityType.DISCOUNT} className="rounded-full data-[state=active]:bg-background/80 text-xs">
                            <Percent className="w-3 h-3 mr-1" />
                            {t('activity.discount', { defaultValue: 'Discount' })}
                        </TabsTrigger>
                        <TabsTrigger value={ActivityType.RECHARGE} className="rounded-full data-[state=active]:bg-background/80 text-xs">
                            <Gift className="w-3 h-3 mr-1" />
                            {t('activity.recharge', { defaultValue: 'Recharge' })}
                        </TabsTrigger>
                        <TabsTrigger value={ActivityType.NEW_USER} className="rounded-full data-[state=active]:bg-background/80 text-xs">
                            <Sparkles className="w-3 h-3 mr-1" />
                            {t('activity.new_user', { defaultValue: 'New User' })}
                        </TabsTrigger>
                        <TabsTrigger value={ActivityType.FLASH} className="rounded-full data-[state=active]:bg-background/80 text-xs">
                            <Zap className="w-3 h-3 mr-1" />
                            {t('activity.flash', { defaultValue: 'Flash' })}
                        </TabsTrigger>
                    </TabsList>

                    <TabsContent value={activeTab} className="mt-0">
                        {loading && activities.length === 0 ? (
                            <div className="space-y-4">
                                {[1, 2, 3].map(i => (
                                    <Skeleton key={i} className="h-64 rounded-xl" />
                                ))}
                            </div>
                        ) : filteredActivities.length === 0 ? (
                            <EmptyState />
                        ) : (
                            <ScrollArea className="h-[calc(100vh-320px)]">
                                <div className="space-y-4 pr-4">
                                    {filteredActivities.map((activity) => (
                                        <ActivityCard
                                            key={activity.id}
                                            activity={activity}
                                            participation={participations.get(activity.id)}
                                            onJoin={() => handleJoin(activity)}
                                            joining={joiningId === activity.id}
                                            canParticipate={canParticipate}
                                        />
                                    ))}
                                </div>
                            </ScrollArea>
                        )}
                    </TabsContent>
                </Tabs>
            </div>
        </PageContainer>
    );
}

function EmptyState() {
    const { t } = useTranslation();
    return (
        <div className="flex flex-col items-center justify-center py-20 text-center animate-in fade-in zoom-in-95">
            <div className="p-6 rounded-full bg-muted/30 mb-4">
                <PartyPopper className="h-10 w-10 text-muted-foreground/50" />
            </div>
            <p className="text-muted-foreground">
                {t('activity.no_activities', { defaultValue: 'No activities available at the moment.' })}
            </p>
        </div>
    );
}

interface ActivityCardProps {
    activity: Activity;
    participation?: ActivityParticipation;
    onJoin: () => void;
    joining: boolean;
    canParticipate: (activity: Activity, participation?: ActivityParticipation) => { can: boolean; reason?: string };
}

function ActivityCard({ activity, participation, onJoin, joining, canParticipate }: ActivityCardProps) {
    const { t } = useTranslation();
    const navigate = useNavigate();

    const getTypeIcon = () => {
        switch (activity.type) {
            case ActivityType.DISCOUNT: return <Percent className="w-4 h-4" />;
            case ActivityType.RECHARGE: return <Gift className="w-4 h-4" />;
            case ActivityType.NEW_USER: return <Sparkles className="w-4 h-4" />;
            case ActivityType.FESTIVAL: return <PartyPopper className="w-4 h-4" />;
            case ActivityType.FLASH: return <Zap className="w-4 h-4" />;
            default: return <Gift className="w-4 h-4" />;
        }
    };

    const getTypeColor = () => {
        switch (activity.type) {
            case ActivityType.DISCOUNT: return 'bg-purple-500';
            case ActivityType.RECHARGE: return 'bg-green-500';
            case ActivityType.NEW_USER: return 'bg-blue-500';
            case ActivityType.FESTIVAL: return 'bg-pink-500';
            case ActivityType.FLASH: return 'bg-orange-500';
            default: return 'bg-primary';
        }
    };

    const getStatusBadge = () => {
        const startAt = new Date(activity.startAt);
        const endAt = new Date(activity.endAt);

        if (activity.status === ActivityStatus.ENDED || isPast(endAt)) {
            return <Badge variant="secondary" className="bg-gray-500/10 text-gray-500">{t('activity.ended', { defaultValue: 'Ended' })}</Badge>;
        }
        if (activity.status === ActivityStatus.SCHEDULED || isFuture(startAt)) {
            return <Badge variant="secondary" className="bg-blue-500/10 text-blue-500">{t('activity.upcoming', { defaultValue: 'Upcoming' })}</Badge>;
        }
        if (activity.status === ActivityStatus.ACTIVE) {
            return <Badge variant="secondary" className="bg-green-500/10 text-green-500 animate-pulse">{t('activity.ongoing', { defaultValue: 'Ongoing' })}</Badge>;
        }
        return null;
    };

    const getTimeDisplay = () => {
        const startAt = new Date(activity.startAt);
        const endAt = new Date(activity.endAt);

        if (isFuture(startAt)) {
            return t('activity.starts_in', {
                defaultValue: 'Starts {{time}}',
                time: formatDistanceToNow(startAt, { addSuffix: true, locale: zhCN })
            });
        }
        if (isPast(endAt)) {
            return t('activity.ended_at', {
                defaultValue: 'Ended {{date}}',
                date: format(endAt, 'MM-dd HH:mm')
            });
        }
        return t('activity.ends_in', {
            defaultValue: 'Ends {{time}}',
            time: formatDistanceToNow(endAt, { addSuffix: true, locale: zhCN })
        });
    };

    const { can, reason } = canParticipate(activity, participation);
    const hasParticipated = participation && participation.participationCount > 0;

    return (
        <Card className="group relative overflow-hidden border-white/5 bg-background/40 backdrop-blur-md hover:border-primary/30 transition-all duration-300">
            {/* Banner */}
            {activity.bannerUrl && (
                <div className="relative h-40 overflow-hidden">
                    <img
                        src={activity.bannerUrl}
                        alt={activity.name}
                        className="w-full h-full object-cover group-hover:scale-105 transition-transform duration-500"
                    />
                    <div className="absolute inset-0 bg-gradient-to-t from-background via-background/50 to-transparent" />

                    {/* Type Badge */}
                    <Badge className={`absolute top-3 left-3 ${getTypeColor()} text-white gap-1`}>
                        {getTypeIcon()}
                        {activity.type === ActivityType.DISCOUNT ? t('activity.type_discount', { defaultValue: 'Discount' }) :
                         activity.type === ActivityType.RECHARGE ? t('activity.type_recharge', { defaultValue: 'Recharge' }) :
                         activity.type === ActivityType.NEW_USER ? t('activity.type_new_user', { defaultValue: 'New User' }) :
                         activity.type === ActivityType.FESTIVAL ? t('activity.type_festival', { defaultValue: 'Festival' }) :
                         t('activity.type_flash', { defaultValue: 'Flash Sale' })}
                    </Badge>

                    {/* Status Badge */}
                    <div className="absolute top-3 right-3">
                        {getStatusBadge()}
                    </div>
                </div>
            )}

            <CardHeader className="pb-2">
                <div className="flex items-start justify-between">
                    <div>
                        <CardTitle className="text-xl">{activity.name}</CardTitle>
                        <CardDescription className="mt-1 line-clamp-2">
                            {activity.description}
                        </CardDescription>
                    </div>
                </div>
            </CardHeader>

            <CardContent className="pb-2">
                {/* Time & Participants Info */}
                <div className="flex flex-wrap gap-4 text-sm text-muted-foreground">
                    <div className="flex items-center gap-1">
                        <Timer className="w-4 h-4" />
                        <span>{getTimeDisplay()}</span>
                    </div>
                    <div className="flex items-center gap-1">
                        <CalendarDays className="w-4 h-4" />
                        <span>{format(new Date(activity.startAt), 'MM/dd')} - {format(new Date(activity.endAt), 'MM/dd')}</span>
                    </div>
                    {activity.maxParticipants && (
                        <div className="flex items-center gap-1">
                            <Users className="w-4 h-4" />
                            <span>{activity.currentParticipants}/{activity.maxParticipants}</span>
                        </div>
                    )}
                </div>

                {/* Rewards Preview */}
                {activity.rewards.length > 0 && (
                    <div className="mt-3 flex flex-wrap gap-2">
                        {activity.rewards.slice(0, 3).map((reward, index) => (
                            <Badge key={index} variant="outline" className="bg-primary/5 border-primary/20">
                                <Gift className="w-3 h-3 mr-1" />
                                {reward.description}
                            </Badge>
                        ))}
                        {activity.rewards.length > 3 && (
                            <Badge variant="outline" className="bg-muted/50">
                                +{activity.rewards.length - 3}
                            </Badge>
                        )}
                    </div>
                )}

                {/* Participation Status */}
                {hasParticipated && (
                    <div className="mt-3 p-2 rounded-lg bg-green-500/10 border border-green-500/20">
                        <div className="flex items-center gap-2 text-sm text-green-500">
                            <Sparkles className="w-4 h-4" />
                            <span>
                                {t('activity.participated', {
                                    defaultValue: 'Participated {{count}} time(s)',
                                    count: participation?.participationCount || 0
                                })}
                            </span>
                        </div>
                    </div>
                )}
            </CardContent>

            <CardFooter className="pt-2 border-t border-white/5 flex justify-between items-center">
                <Button
                    variant="ghost"
                    size="sm"
                    className="text-muted-foreground hover:text-foreground"
                    onClick={() => navigate(`/activity/${activity.id}`)}
                >
                    {t('activity.view_details', { defaultValue: 'View Details' })}
                    <ChevronRight className="w-4 h-4 ml-1" />
                </Button>

                <Button
                    onClick={onJoin}
                    disabled={joining || !can}
                    className={can ? 'bg-gradient-to-r from-pink-500 to-purple-500 hover:from-pink-600 hover:to-purple-600' : ''}
                >
                    {joining ? t('common.loading', { defaultValue: 'Loading...' }) :
                     !can ? (reason || t('activity.cannot_join', { defaultValue: 'Cannot Join' })) :
                     hasParticipated ? t('activity.join_again', { defaultValue: 'Join Again' }) :
                     t('activity.join_now', { defaultValue: 'Join Now' })}
                </Button>
            </CardFooter>
        </Card>
    );
}
