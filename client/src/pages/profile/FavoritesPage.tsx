import { useEffect, useState } from 'react';
import { PageContainer } from '@/components/page-container';
import { useTranslation } from 'react-i18next';
import { http } from '@/lib/http';
import { Button } from '@/components/ui/button';
import { Card, CardContent } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { FavoriteButton } from '@/components/player/favorite-button';
import { Loader2, Heart, Star, Gamepad2 } from 'lucide-react';
import { useNavigate } from 'react-router-dom';
import type { Player } from '@/stores/modules/player-store';

export default function FavoritesPage() {
    const { t } = useTranslation();
    const navigate = useNavigate();
    const [players, setPlayers] = useState<Player[]>([]);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        const fetchFavorites = async () => {
            try {
                const data = await http.get<Player[]>('/user/favorites/players');
                setPlayers(data);
            } catch (err) {
                console.error("Failed to fetch favorites", err);
            } finally {
                setLoading(false);
            }
        };

        fetchFavorites();
    }, []);

    const handleRemove = (id: number) => {
        setPlayers(prev => prev.filter(p => p.id !== id));
    };

    return (
        <PageContainer>
            <div className="max-w-6xl mx-auto h-full flex flex-col space-y-6">
                <div className="flex flex-col space-y-2">
                    <h1 className="text-3xl font-bold tracking-tight bg-clip-text text-transparent bg-gradient-to-r from-pink-400 to-purple-400">
                        {t('nav.favorites', { defaultValue: 'My Favorites' })}
                    </h1>
                    <p className="text-muted-foreground">
                        {t('favorites.desc', { defaultValue: 'Your favorite players are here.' })}
                    </p>
                </div>

                <div className="flex-1 min-h-0 bg-muted/30 rounded-xl border border-white/5 p-4">
                    {loading ? (
                        <div className="flex h-full items-center justify-center" data-testid="loading-spinner">
                            <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
                        </div>
                    ) : players.length === 0 ? (
                        <div className="flex flex-col items-center justify-center py-20 text-center animate-in fade-in zoom-in-95" data-testid="empty-state">
                            <div className="p-6 rounded-full bg-pink-500/10 mb-4">
                                <Heart className="h-10 w-10 text-pink-500" />
                            </div>
                            <h3 className="text-lg font-medium">{t('favorites.empty', { defaultValue: 'No favorites yet' })}</h3>
                            <p className="text-sm text-muted-foreground mt-1 max-w-[250px] mx-auto">
                                {t('favorites.empty_desc', { defaultValue: 'Go find some awesome players!' })}
                            </p>
                            <Button className="mt-6 rounded-full px-8 bg-pink-600 hover:bg-pink-700" onClick={() => navigate('/players')}>
                                {t('nav.find_players_title', { defaultValue: 'Find Players' })}
                            </Button>
                        </div>
                    ) : (
                        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
                            {players.map((player) => (
                                <Card key={player.id} className="group overflow-hidden border-white/5 bg-sidebar hover:border-pink-500/30 transition-all duration-300">
                                    <div className="relative aspect-[4/3] bg-muted">
                                        <img
                                            src={player.avatar || `https://api.dicebear.com/7.x/avataaars/svg?seed=${player.username}`}
                                            alt={player.nickname}
                                            className="h-full w-full object-cover transition-transform duration-500 group-hover:scale-110"
                                        />
                                        <div className="absolute inset-0 bg-gradient-to-t from-black/80 via-black/20 to-transparent opacity-80" />

                                        <div className="absolute top-3 right-3 z-10">
                                            <FavoriteButton
                                                playerId={player.id}
                                                initialIsFavorite={true}
                                                onToggle={(isFav) => !isFav && handleRemove(player.id)}
                                                className="bg-black/20 hover:bg-black/40 backdrop-blur-md"
                                            />
                                        </div>

                                        <div className="absolute bottom-3 left-3 right-3 text-white">
                                            <div className="flex items-center justify-between">
                                                <h3 className="font-bold text-lg truncate">{player.nickname}</h3>
                                                <div className="flex items-center text-yellow-400 text-sm font-bold bg-black/40 px-1.5 py-0.5 rounded backdrop-blur-md">
                                                    <Star className="h-3 w-3 fill-yellow-400 mr-1" />
                                                    {player.rating?.toFixed(1) || '-'}
                                                </div>
                                            </div>
                                            <div className="flex items-center text-xs text-white/80 mt-1">
                                                <Gamepad2 className="h-3 w-3 mr-1" />
                                                Valorant
                                            </div>
                                        </div>
                                    </div>
                                    <CardContent className="p-4 space-y-4">
                                        <div className="flex flex-wrap gap-1.5 h-6 overflow-hidden">
                                            {player.tags?.slice(0, 3).map(tag => (
                                                <Badge key={tag} variant="secondary" className="text-[10px] px-1.5 py-0 h-5 bg-white/5 hover:bg-white/10">{tag}</Badge>
                                            ))}
                                        </div>
                                        <div className="flex items-center justify-between pt-2">
                                            <div className="flex items-baseline gap-1">
                                                <span className="text-lg font-bold text-primary">¥{player.price}</span>
                                                <span className="text-xs text-muted-foreground">/hr</span>
                                            </div>
                                            <Button size="sm" className="h-8 rounded-full px-4" onClick={() => navigate(`/players/${player.id}`)}>
                                                {t('player.book', { defaultValue: 'Book' })}
                                            </Button>
                                        </div>
                                    </CardContent>
                                </Card>
                            ))}
                        </div>
                    )}
                </div>
            </div>
        </PageContainer>
    );
}
