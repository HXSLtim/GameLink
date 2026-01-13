import { useEffect } from 'react';
import { usePlayerStore } from '@/stores';
import { PageContainer } from '@/components/page-container';
import { Card, CardContent } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { Skeleton } from '@/components/ui/skeleton';
import { Separator } from '@/components/ui/separator';
import { ScrollArea, ScrollBar } from '@/components/ui/scroll-area';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import { Gamepad2, Trophy, Users, ArrowRight, Star, Zap } from 'lucide-react';
import { useNavigate } from 'react-router-dom';

const GAMES = [
    { id: 1, name: 'League of Legends', icon: Trophy, color: 'text-yellow-500', bg: 'bg-yellow-500/10' },
    { id: 2, name: 'Valorant', icon: Zap, color: 'text-red-500', bg: 'bg-red-500/10' },
    { id: 3, name: 'Apex Legends', icon: Gamepad2, color: 'text-orange-500', bg: 'bg-orange-500/10' },
    { id: 4, name: 'Overwatch 2', icon: Users, color: 'text-blue-500', bg: 'bg-blue-500/10' },
];

export default function HomePage() {
    const { featuredPlayers, fetchFeaturedPlayers, loading } = usePlayerStore();
    const navigate = useNavigate();

    useEffect(() => {
        fetchFeaturedPlayers();
    }, []);

    return (
        <PageContainer>
            <div className="space-y-8 pb-8">
                {/* Hero Section */}
                <section className="relative overflow-hidden rounded-xl bg-gradient-to-r from-primary/10 via-primary/5 to-background border p-8 md:p-12 mx-4 mt-4">
                    <div className="relative z-10 max-w-2xl space-y-4">
                        <Badge variant="secondary" className="mb-2">New Season Live</Badge>
                        <h1 className="text-4xl md:text-5xl font-extrabold tracking-tight lg:text-6xl">
                            Find Your Perfect <span className="text-primary">Duo</span>
                        </h1>
                        <p className="text-muted-foreground text-lg md:text-xl max-w-[600px]">
                            Connect with top-tier gamers for coaching, ranking up, or just having fun. Join the community today.
                        </p>
                        <div className="flex flex-wrap gap-4 pt-4">
                            <Button size="lg" onClick={() => navigate('/players')}>
                                Find Players <ArrowRight className="ml-2 h-4 w-4" />
                            </Button>
                            <Button size="lg" variant="outline" onClick={() => navigate('/profile')}>
                                Become a Player
                            </Button>
                        </div>
                    </div>
                    {/* Decorative Background Elements */}
                    <div className="absolute top-0 right-0 -translate-y-1/4 translate-x-1/4 opacity-10 pointer-events-none">
                        <Gamepad2 className="w-96 h-96 text-primary" />
                    </div>
                </section>

                {/* Popular Games */}
                <section className="px-4 space-y-4">
                    <div className="flex items-center justify-between">
                        <h2 className="text-2xl font-bold tracking-tight">Popular Games</h2>
                        <Button variant="ghost" size="sm" onClick={() => navigate('/players')}>View All</Button>
                    </div>
                    <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
                        {GAMES.map((game) => (
                            <Card
                                key={game.id}
                                className="cursor-pointer hover:bg-muted/50 transition-colors border-none shadow-sm bg-muted/20"
                                onClick={() => navigate(`/players?gameId=${game.id}`)}
                            >
                                <CardContent className="p-6 flex flex-col items-center justify-center gap-3 text-center">
                                    <div className={`p-3 rounded-full ${game.bg} ${game.color}`}>
                                        <game.icon className="h-8 w-8" />
                                    </div>
                                    <span className="font-semibold">{game.name}</span>
                                </CardContent>
                            </Card>
                        ))}
                    </div>
                </section>

                <Separator className="mx-4" />

                {/* Featured Players */}
                <section className="px-4 space-y-4">
                    <div className="flex items-center justify-between">
                        <h2 className="text-2xl font-bold tracking-tight">Featured Pros</h2>
                        <Button variant="ghost" size="sm" onClick={() => navigate('/players')}>View All</Button>
                    </div>

                    <div className="relative">
                        <ScrollArea>
                            <div className="flex space-x-4 pb-4">
                                {loading || featuredPlayers.length === 0 ? (
                                    Array(5).fill(0).map((_, i) => (
                                        <Skeleton key={i} className="h-[280px] w-[220px] rounded-xl flex-none" />
                                    ))
                                ) : (
                                    featuredPlayers.map((player) => (
                                        <Card
                                            key={player.id}
                                            className="w-[220px] flex-none cursor-pointer hover:border-primary/50 transition-all group overflow-hidden"
                                            onClick={() => navigate(`/chat/${player.id}`)} // Or player detail modal
                                        >
                                            <div className="h-[140px] w-full bg-muted relative overflow-hidden">
                                                {/* Cover Image Placeholder */}
                                                <div className="absolute inset-0 bg-gradient-to-t from-background/80 to-transparent z-10" />
                                                <Avatar className="absolute bottom-2 left-2 z-20 h-10 w-10 border-2 border-background">
                                                    <AvatarImage src={player.avatar} />
                                                    <AvatarFallback>{player.username[0]}</AvatarFallback>
                                                </Avatar>
                                            </div>
                                            <CardContent className="p-4 space-y-2">
                                                <div>
                                                    <h3 className="font-semibold truncate group-hover:text-primary transition-colors">{player.nickname || player.username}</h3>
                                                    <p className="text-xs text-muted-foreground truncate">{player.gameName}</p>
                                                </div>
                                                <div className="flex items-center gap-1 text-xs font-medium text-yellow-500">
                                                    <Star className="h-3 w-3 fill-current" />
                                                    {player.rating.toFixed(1)}
                                                    <span className="text-muted-foreground ml-1">({player.orderCount} Orders)</span>
                                                </div>
                                                <div className="flex items-center justify-between pt-2">
                                                    <span className="font-bold text-lg">¥{player.price}</span>
                                                    <Badge variant="secondary" className="text-[10px] h-5">PRO</Badge>
                                                </div>
                                            </CardContent>
                                        </Card>
                                    ))
                                )}
                            </div>
                            <ScrollBar orientation="horizontal" />
                        </ScrollArea>
                    </div>
                </section>

                {/* Call to Action or Info Strip */}
                <section className="px-4">
                    <div className="rounded-lg bg-muted p-8 flex flex-col items-center text-center space-y-4">
                        <h3 className="text-2xl font-bold">Why choose GameLink?</h3>
                        <div className="grid grid-cols-1 md:grid-cols-3 gap-8 w-full max-w-4xl pt-4">
                            <div className="space-y-2">
                                <Shield className="h-8 w-8 text-primary mx-auto" />
                                <h4 className="font-semibold">Secure Payments</h4>
                                <p className="text-sm text-muted-foreground">Your funds are protected until the order is complete.</p>
                            </div>
                            <div className="space-y-2">
                                <Users className="h-8 w-8 text-primary mx-auto" />
                                <h4 className="font-semibold">Verified Pros</h4>
                                <p className="text-sm text-muted-foreground">All players undergo strict skill verification.</p>
                            </div>
                            <div className="space-y-2">
                                <Clock className="h-8 w-8 text-primary mx-auto" />
                                <h4 className="font-semibold">24/7 Support</h4>
                                <p className="text-sm text-muted-foreground">We are here to help you anytime, anywhere.</p>
                            </div>
                        </div>
                    </div>
                </section>
            </div>
        </PageContainer>
    );
}

// Importing icons that were missed in the top import if any
import { Shield, Clock } from 'lucide-react';
