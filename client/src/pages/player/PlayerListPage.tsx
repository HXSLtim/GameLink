import { useEffect, useState } from "react";
import { usePlayerStore } from "@/stores";
import { PageContainer } from "@/components/page-container";
import { Card, CardContent, CardFooter } from "@/components/ui/card";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Star, Gamepad2, ShoppingBag, Joystick, Filter } from "lucide-react";
import { Skeleton } from "@/components/ui/skeleton";
import { useTranslation } from "react-i18next";
import { useNavigate } from "react-router-dom";

export default function PlayerListPage() {
    const { players, loading, fetchPlayers, filters, setFilters, pagination } = usePlayerStore();
    const { t } = useTranslation();
    const navigate = useNavigate();
    const [observerRef, setObserverRef] = useState<HTMLDivElement | null>(null);

    // Initial fetch
    useEffect(() => {
        fetchPlayers(true);
    }, []);

    // Infinite scroll observer
    useEffect(() => {
        const observer = new IntersectionObserver(
            (entries) => {
                if (entries[0].isIntersecting && pagination.hasMore && !loading) {
                    fetchPlayers(false);
                }
            },
            { threshold: 0.1 }
        );

        if (observerRef) {
            observer.observe(observerRef);
        }

        return () => {
            if (observerRef) {
                observer.unobserve(observerRef);
            }
        };
    }, [observerRef, pagination.hasMore, loading, fetchPlayers]);

    return (
        <PageContainer>
            <div className="space-y-6">
                {/* Header Section */}
                <div className="flex flex-col md:flex-row md:items-center justify-between gap-4 animate-in fade-in slide-in-from-top-2">
                    <div className="space-y-1">
                        <h1 className="text-3xl font-bold tracking-tight bg-clip-text text-transparent bg-gradient-to-r from-white to-white/60">
                            {t('nav.find_players_title')}
                        </h1>
                        <p className="text-muted-foreground">
                            {t('nav.find_players_desc')}
                        </p>
                    </div>
                </div>

                {/* Filter Bar */}
                <div className="flex items-center gap-2 overflow-x-auto pb-4 scrollbar-hide animate-in fade-in slide-in-from-left-2 duration-500 delay-100">
                    <div className="flex items-center gap-2 p-1 bg-muted/40 backdrop-blur-md rounded-full border border-white/5">
                        <div className="pl-3 pr-2 text-muted-foreground">
                            <Filter className="w-4 h-4" />
                        </div>
                        <Badge
                            variant={filters.sortBy === 'rating' ? 'secondary' : 'outline'}
                            className={`cursor-pointer rounded-full px-4 py-1.5 transition-all ${filters.sortBy === 'rating' ? 'bg-primary/20 text-primary hover:bg-primary/30' : 'hover:bg-white/10 border-transparent text-muted-foreground'}`}
                            onClick={() => setFilters({ sortBy: 'rating' })}
                        >
                            {t('nav.filter.rating')}
                        </Badge>
                        <Badge
                            variant={filters.sortBy === 'orders' ? 'secondary' : 'outline'}
                            className={`cursor-pointer rounded-full px-4 py-1.5 transition-all ${filters.sortBy === 'orders' ? 'bg-primary/20 text-primary hover:bg-primary/30' : 'hover:bg-white/10 border-transparent text-muted-foreground'}`}
                            onClick={() => setFilters({ sortBy: 'orders' })}
                        >
                            {t('nav.filter.orders')}
                        </Badge>
                        <Badge
                            variant={filters.sortBy === 'price' ? 'secondary' : 'outline'}
                            className={`cursor-pointer rounded-full px-4 py-1.5 transition-all ${filters.sortBy === 'price' ? 'bg-primary/20 text-primary hover:bg-primary/30' : 'hover:bg-white/10 border-transparent text-muted-foreground'}`}
                            onClick={() => setFilters({ sortBy: 'price' })}
                        >
                            {t('nav.filter.price')}
                        </Badge>
                    </div>
                </div>

                {loading && players.length === 0 ? (
                    <div className="grid gap-6 grid-cols-2 md:grid-cols-3 lg:grid-cols-4">
                        {Array.from({ length: 8 }).map((_, i) => (
                            <Card key={i} className="overflow-hidden border-border/40 bg-muted/20">
                                <Skeleton className="h-[200px] w-full rounded-none" />
                                <CardContent className="p-4 space-y-3">
                                    <Skeleton className="h-4 w-3/4" />
                                    <Skeleton className="h-4 w-1/2" />
                                </CardContent>
                            </Card>
                        ))}
                    </div>
                ) : (
                    <div className="grid gap-6 grid-cols-2 md:grid-cols-3 lg:grid-cols-4 pb-20">
                        {players.map((player, index) => (
                            <Card
                                key={player.id}
                                className="group relative overflow-hidden border-white/5 bg-background/40 backdrop-blur-md hover:border-primary/50 hover:shadow-2xl hover:-translate-y-1 transition-all duration-300 cursor-pointer"
                                style={{ animationDelay: `${index * 50}ms` }}
                                onClick={() => navigate(`/players/${player.id}`)}
                            >
                                <div className="absolute inset-0 bg-gradient-to-br from-primary/5 via-transparent to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-500" />

                                {/* Image Section */}
                                <div className="relative h-[220px] overflow-hidden bg-muted">
                                    <div className="absolute inset-0 bg-gradient-to-t from-background/90 via-background/20 to-transparent z-10" />
                                    <img
                                        src={`https://api.dicebear.com/7.x/shapes/svg?seed=${player.username}`} // Placeholder cover style
                                        alt="Cover"
                                        className="w-full h-full object-cover transition-transform duration-700 group-hover:scale-110"
                                    />

                                    {/* Status Badge */}
                                    <div className="absolute top-3 right-3 z-20">
                                        <Badge variant={player.online ? "default" : "secondary"} className={`backdrop-blur-md shadow-sm ${player.online ? 'bg-green-500 hover:bg-green-600' : 'bg-black/50 text-white border-white/10'}`}>
                                            {player.online ? t('nav.status.online') : t('nav.status.offline')}
                                        </Badge>
                                    </div>

                                    {/* Avatar & Name Overlay */}
                                    <div className="absolute bottom-3 left-3 z-20 flex items-end gap-3 w-[calc(100%-24px)]">
                                        <div className="relative">
                                            <Avatar className="w-12 h-12 border-2 border-background shadow-lg ring-2 ring-white/5 group-hover:ring-primary/50 transition-all">
                                                <AvatarImage src={player.avatar} />
                                                <AvatarFallback className="bg-primary/20 text-primary">{player.nickname?.[0] || player.username?.[0] || '?'}</AvatarFallback>
                                            </Avatar>
                                            <span className={`absolute bottom-0 right-0 w-3 h-3 rounded-full border-2 border-background ${player.online ? 'bg-green-500' : 'bg-gray-500'}`} />
                                        </div>
                                        <div className="flex-1 min-w-0 mb-1">
                                            <div className="font-bold text-base leading-tight truncate text-white drop-shadow-md group-hover:text-primary transition-colors">
                                                {player.nickname}
                                            </div>
                                            <div className="flex items-center gap-1.5 text-xs text-white/80 mt-0.5">
                                                <Joystick className="w-3 h-3" />
                                                <span className="truncate">{player.gameName || t('player.all_games')}</span>
                                            </div>
                                        </div>
                                    </div>
                                </div>

                                <CardContent className="p-4 grid gap-3 relative z-10">
                                    <div className="flex items-center justify-between text-sm">
                                        <div className="flex items-center gap-1.5 bg-yellow-500/10 text-yellow-500 px-2 py-1 rounded-md font-medium">
                                            <Star className="w-3.5 h-3.5 fill-current" />
                                            {player.rating.toFixed(1)}
                                        </div>
                                        <div className="text-muted-foreground flex items-center gap-1 text-xs">
                                            <ShoppingBag className="w-3.5 h-3.5" />
                                            {player.orderCount} {t('nav.filter.orders')}
                                        </div>
                                    </div>

                                    <div className="space-y-2">
                                        <div className="flex flex-wrap gap-1.5 h-6 overflow-hidden">
                                            {player.tags.slice(0, 3).map(tag => (
                                                <span key={tag} className="text-[10px] px-2 py-0.5 rounded-full bg-secondary/50 text-secondary-foreground border border-white/5">
                                                    {tag}
                                                </span>
                                            ))}
                                        </div>
                                    </div>
                                </CardContent>

                                <CardFooter className="p-4 pt-0 flex bg-transparent items-center justify-between relative z-10">
                                    <div className="flex flex-col">
                                        <span className="text-xs text-muted-foreground">{t('player.starting_at')}</span>
                                        <span className="text-xl font-bold text-primary">¥{player.price}<span className="text-xs font-normal text-muted-foreground">{t('nav.home_content.per_hour')}</span></span>
                                    </div>
                                    <Button size="sm" className="rounded-full shadow-lg shadow-primary/20 group-hover:bg-primary group-hover:text-primary-foreground transition-all">
                                        {t('nav.action.book_now')}
                                    </Button>
                                </CardFooter>
                            </Card>
                        ))}
                    </div>
                )}

                {/* Infinite Scroll Loader */}
                {!loading && players.length > 0 && pagination.hasMore && (
                    <div ref={setObserverRef} className="flex justify-center py-4 w-full">
                        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary"></div>
                    </div>
                )}

                {loading && players.length > 0 && (
                    <div className="flex justify-center py-4 w-full">
                        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary"></div>
                    </div>
                )}

                {!loading && players.length === 0 && (
                    <div className="flex flex-col items-center justify-center py-24 text-center space-y-4 animate-in fade-in zoom-in-95">
                        <div className="p-6 rounded-full bg-muted/30">
                            <Gamepad2 className="w-12 h-12 text-muted-foreground/50" />
                        </div>
                        <div className="space-y-2">
                            <h3 className="text-xl font-semibold">{t('nav.no_players')}</h3>
                            <p className="text-muted-foreground max-w-sm mx-auto">{t('player.no_results_desc')}</p>
                        </div>
                    </div>
                )}
            </div>
        </PageContainer>
    );
}
