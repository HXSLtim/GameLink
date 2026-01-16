import { useEffect } from 'react';
import { usePlayerStore } from '@/stores';
import { PageContainer } from '@/components/page-container';
import { Card, CardContent } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { Skeleton } from '@/components/ui/skeleton';
import { ScrollArea, ScrollBar } from '@/components/ui/scroll-area';
import { Gamepad2, Trophy, Users, ArrowRight, Star, Zap, Shield, Clock, Bell } from 'lucide-react';
import { useNavigate } from 'react-router-dom';
import { useTranslation } from 'react-i18next';

const GAMES = [
    { id: 1, name: 'League of Legends', icon: Trophy, color: 'text-yellow-500', from: 'from-yellow-500/20', to: 'to-orange-500/20' },
    { id: 2, name: 'Valorant', icon: Zap, color: 'text-red-500', from: 'from-red-500/20', to: 'to-rose-500/20' },
    { id: 3, name: 'Apex Legends', icon: Gamepad2, color: 'text-orange-500', from: 'from-orange-500/20', to: 'to-amber-500/20' },
    { id: 4, name: 'Overwatch 2', icon: Users, color: 'text-blue-500', from: 'from-blue-500/20', to: 'to-cyan-500/20' },
];

export default function HomePage() {
    const { featuredPlayers, fetchFeaturedPlayers, loading } = usePlayerStore();
    const navigate = useNavigate();
    const { t } = useTranslation();

    useEffect(() => {
        fetchFeaturedPlayers();
    }, [fetchFeaturedPlayers]);

    return (
        <PageContainer>
            <div className="space-y-10 pb-12">
                {/* Hero Section */}
                <section className="relative overflow-hidden rounded-3xl bg-gradient-to-br from-violet-600/10 via-background to-indigo-600/10 border border-white/5 p-8 md:p-16 mx-0 mt-4 animate-in fade-in zoom-in-95 duration-500">
                    <div className="absolute inset-0 bg-[url('https://images.unsplash.com/photo-1511512578047-dfb367046420?q=80&w=2671&auto=format&fit=crop')] bg-cover bg-center opacity-5 mix-blend-overlay pointer-events-none" />

                    <div className="relative z-10 max-w-3xl space-y-6">
                        <Badge variant="outline" className="mb-2 bg-primary/10 text-primary border-primary/20 backdrop-blur-sm px-3 py-1">
                            <span className="relative flex h-2 w-2 mr-2">
                                <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-primary opacity-75"></span>
                                <span className="relative inline-flex rounded-full h-2 w-2 bg-primary"></span>
                            </span>
                            {t('nav.new_season')}
                        </Badge>
                        <h1 className="text-4xl md:text-6xl font-extrabold tracking-tight leading-tight">
                            {t('nav.hero_title')} <span className="text-transparent bg-clip-text bg-gradient-to-r from-violet-500 to-indigo-500">{t('nav.hero_highlight')}</span>
                        </h1>
                        <p className="text-muted-foreground text-lg md:text-xl max-w-[600px] leading-relaxed">
                            {t('nav.hero_desc')}
                        </p>
                        <div className="flex flex-wrap gap-4 pt-6">
                            <Button size="lg" className="rounded-full px-8 bg-primary hover:bg-primary/90 shadow-lg shadow-primary/25 transition-all hover:scale-105" onClick={() => navigate('/players')}>
                                {t('nav.find_players')} <ArrowRight className="ml-2 h-4 w-4" />
                            </Button>
                            <Button size="lg" variant="outline" className="rounded-full px-8 backdrop-blur-sm hover:bg-white/10" onClick={() => navigate('/profile')}>
                                {t('nav.become_player')}
                            </Button>
                        </div>
                    </div>
                </section>

                {/* Popular Games */}
                <section className="space-y-6">
                    <div className="flex items-center justify-between px-2">
                        <div className="space-y-1">
                            <h2 className="text-2xl font-bold tracking-tight">{t('nav.popular_games')}</h2>
                            <p className="text-sm text-muted-foreground">{t('nav.home_content.select_game_desc')}</p>
                        </div>
                        <Button variant="ghost" size="sm" onClick={() => navigate('/players')} className="group">
                            {t('nav.view_all')} <ArrowRight className="ml-2 h-4 w-4 transition-transform group-hover:translate-x-1" />
                        </Button>
                    </div>
                    <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
                        {GAMES.map((game, i) => (
                            <Card
                                key={game.id}
                                className="group cursor-pointer hover:border-primary/50 transition-all duration-300 hover:-translate-y-1 hover:shadow-xl bg-background/50 backdrop-blur-sm border-white/5 overflow-hidden"
                                onClick={() => navigate(`/players?gameId=${game.id}`)}
                                style={{ animationDelay: `${i * 100}ms` }}
                            >
                                <div className={`absolute inset-0 bg-gradient-to-br ${game.from} ${game.to} opacity-0 group-hover:opacity-100 transition-opacity duration-300`} />
                                <CardContent className="relative p-6 flex flex-col items-center justify-center gap-4 text-center z-10">
                                    <div className={`p-4 rounded-2xl bg-background/80 shadow-inner ring-1 ring-white/10 group-hover:scale-110 transition-transform duration-300 ${game.color}`}>
                                        <game.icon className="h-8 w-8" />
                                    </div>
                                    <span className="font-semibold text-lg group-hover:text-primary transition-colors">{game.name}</span>
                                </CardContent>
                            </Card>
                        ))}
                    </div>
                </section>

                {/* Featured Players */}
                <section className="space-y-6">
                    <div className="flex items-center justify-between px-2">
                        <div className="space-y-1">
                            <h2 className="text-2xl font-bold tracking-tight">{t('nav.featured_pros')}</h2>
                            <p className="text-sm text-muted-foreground">{t('nav.home_content.top_players_desc')}</p>
                        </div>
                        <div className="flex gap-2">
                            <Button variant="ghost" size="icon" className="relative group" onClick={() => navigate('/notifications')}>
                                <Bell className="h-5 w-5 group-hover:text-primary transition-colors" />
                                <span className="absolute top-2.5 right-2 h-2 w-2 bg-red-500 rounded-full animate-pulse ring-2 ring-background" />
                            </Button>
                            <Button variant="ghost" size="sm" onClick={() => navigate('/players')} className="group">
                                {t('nav.view_all')} <ArrowRight className="ml-2 h-4 w-4 transition-transform group-hover:translate-x-1" />
                            </Button>
                        </div>
                    </div>

                    <div className="relative -mx-4 px-4 md:mx-0 md:px-0">
                        <ScrollArea className="w-full whitespace-nowrap rounded-lg">
                            <div className="flex w-max space-x-4 pb-4">
                                {loading || featuredPlayers.length === 0 ? (
                                    Array(5).fill(0).map((_, i) => (
                                        <Skeleton key={i} className="h-[280px] w-[240px] rounded-2xl flex-none bg-muted/50" />
                                    ))
                                ) : (
                                    featuredPlayers.map((player) => (
                                        <Card
                                            key={player.id}
                                            className="w-[240px] flex-none cursor-pointer border-white/5 hover:border-primary/50 transition-all duration-300 group overflow-hidden bg-background/50 backdrop-blur-sm hover:shadow-lg hover:-translate-y-1"
                                            onClick={() => navigate(`/chat/${player.id}`)}
                                        >
                                            <div className="h-[160px] w-full bg-muted relative overflow-hidden">
                                                <div className="absolute inset-0 bg-gradient-to-t from-background via-transparent to-transparent z-10" />
                                                <img
                                                    src={player.avatar || "https://github.com/shadcn.png"}
                                                    alt={player.nickname}
                                                    className="w-full h-full object-cover transition-transform duration-500 group-hover:scale-105"
                                                />
                                                <div className="absolute top-2 right-2 z-20">
                                                    <Badge variant="secondary" className="bg-background/80 backdrop-blur-sm text-[10px] font-bold shadow-sm">
                                                        TOP 10
                                                    </Badge>
                                                </div>
                                            </div>
                                            <CardContent className="p-4 space-y-3 relative z-20">
                                                <div>
                                                    <h3 className="font-semibold text-base truncate group-hover:text-primary transition-colors">{player.nickname || player.username}</h3>
                                                    <p className="text-xs text-muted-foreground truncate font-medium">{player.gameName}</p>
                                                </div>
                                                <div className="flex items-center justify-between text-xs font-medium">
                                                    <div className="flex items-center gap-1 text-amber-500 bg-amber-500/10 px-2 py-0.5 rounded-full">
                                                        <Star className="h-3 w-3 fill-current" />
                                                        {player.rating.toFixed(1)}
                                                    </div>
                                                    <span className="text-muted-foreground">{player.orderCount} {t('nav.home_content.orders_suffix')}</span>
                                                </div>
                                                <div className="pt-2 flex items-center justify-between border-t border-border/50">
                                                    <span className="font-bold text-lg text-primary">¥{player.price}<span className="text-xs text-muted-foreground font-normal">{t('nav.home_content.per_hour')}</span></span>
                                                    <div className="h-8 w-8 rounded-full bg-primary/10 flex items-center justify-center text-primary group-hover:bg-primary group-hover:text-white transition-colors">
                                                        <ArrowRight className="h-4 w-4" />
                                                    </div>
                                                </div>
                                            </CardContent>
                                        </Card>
                                    ))
                                )}
                            </div>
                            <ScrollBar orientation="horizontal" className="hidden" />
                        </ScrollArea>
                    </div>
                </section>

                {/* Trust Indicators */}
                <section className="py-8">
                    <div className="rounded-3xl bg-gradient-to-b from-muted/50 to-muted/10 border border-white/5 p-8 md:p-12">
                        <div className="text-center space-y-2 mb-10">
                            <h3 className="text-2xl font-bold">{t('nav.why_choose')}</h3>
                            <p className="text-muted-foreground">{t('nav.home_content.features_subtitle')}</p>
                        </div>
                        <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
                            <div className="flex flex-col items-center text-center space-y-3 p-4 rounded-xl hover:bg-white/5 transition-colors">
                                <div className="p-3 rounded-full bg-emerald-500/10 text-emerald-500 mb-2">
                                    <Shield className="h-8 w-8" />
                                </div>
                                <h4 className="font-semibold text-lg">{t('nav.secure_payments')}</h4>
                                <p className="text-sm text-muted-foreground max-w-[250px]">{t('nav.secure_desc')}</p>
                            </div>
                            <div className="flex flex-col items-center text-center space-y-3 p-4 rounded-xl hover:bg-white/5 transition-colors">
                                <div className="p-3 rounded-full bg-blue-500/10 text-blue-500 mb-2">
                                    <Users className="h-8 w-8" />
                                </div>
                                <h4 className="font-semibold text-lg">{t('nav.verified_pros')}</h4>
                                <p className="text-sm text-muted-foreground max-w-[250px]">{t('nav.verified_desc')}</p>
                            </div>
                            <div className="flex flex-col items-center text-center space-y-3 p-4 rounded-xl hover:bg-white/5 transition-colors">
                                <div className="p-3 rounded-full bg-violet-500/10 text-violet-500 mb-2">
                                    <Clock className="h-8 w-8" />
                                </div>
                                <h4 className="font-semibold text-lg">{t('nav.support_247')}</h4>
                                <p className="text-sm text-muted-foreground max-w-[250px]">{t('nav.support_desc')}</p>
                            </div>
                        </div>
                    </div>
                </section>
            </div>
        </PageContainer>
    );
}
