import { useEffect } from "react";
import { usePlayerStore } from "@/stores";
import { PageContainer, PageHeader } from "@/components/page-container";
import { Card, CardContent, CardFooter, CardHeader } from "@/components/ui/card";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Star, Gamepad2, ShoppingBag } from "lucide-react";
import { Skeleton } from "@/components/ui/skeleton";

export default function PlayerListPage() {
    const {
        players,
        loading,
        fetchPlayers,
        filters,
        setFilters
    } = usePlayerStore();

    useEffect(() => {
        fetchPlayers(true); // Initial fetch
    }, []); // eslint-disable-line react-hooks/exhaustive-deps

    return (
        <PageContainer>
            <PageHeader
                title="Find Players"
                description="Discover talented players to game with."
            />

            {/* Filters (Basic placeholder) */}
            <div className="flex gap-2 mb-6 overflow-x-auto pb-2">
                <Badge
                    variant={filters.sortBy === 'rating' ? 'default' : 'outline'}
                    className="cursor-pointer"
                    onClick={() => setFilters({ sortBy: 'rating' })}
                >
                    Highest Rated
                </Badge>
                <Badge
                    variant={filters.sortBy === 'orders' ? 'default' : 'outline'}
                    className="cursor-pointer"
                    onClick={() => setFilters({ sortBy: 'orders' })}
                >
                    Most Orders
                </Badge>
                <Badge
                    variant={filters.sortBy === 'price' ? 'default' : 'outline'}
                    className="cursor-pointer"
                    onClick={() => setFilters({ sortBy: 'price' })}
                >
                    Price
                </Badge>
            </div>

            {loading && players.length === 0 ? (
                <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
                    {Array.from({ length: 8 }).map((_, i) => (
                        <Card key={i} className="overflow-hidden">
                            <CardHeader className="p-0">
                                <Skeleton className="h-48 w-full" />
                            </CardHeader>
                            <CardContent className="p-4 space-y-2">
                                <Skeleton className="h-4 w-1/2" />
                                <Skeleton className="h-4 w-1/4" />
                            </CardContent>
                        </Card>
                    ))}
                </div>
            ) : (
                <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
                    {players.map((player) => (
                        <Card key={player.id} className="overflow-hidden hover:shadow-lg transition-shadow border-border/50">
                            <div className="relative h-48 bg-muted">
                                {/* Cover Image Placeholder */}
                                <div className="absolute inset-0 bg-gradient-to-t from-black/60 to-transparent z-10" />
                                <img
                                    src={`https://api.dicebear.com/7.x/shapes/svg?seed=${player.username}`}
                                    alt="Cover"
                                    className="w-full h-full object-cover"
                                />
                                <div className="absolute bottom-3 left-3 z-20 flex items-center gap-2">
                                    <Avatar className="w-10 h-10 border-2 border-background">
                                        <AvatarImage src={player.avatar} />
                                        <AvatarFallback>{player.nickname[0]}</AvatarFallback>
                                    </Avatar>
                                    <div className="text-white">
                                        <div className="font-bold leading-tight">{player.nickname}</div>
                                        <div className="text-xs opacity-80 flex items-center gap-1">
                                            <span className={`w-2 h-2 rounded-full ${player.online ? 'bg-green-500' : 'bg-gray-400'}`} />
                                            {player.online ? 'Online' : 'Offline'}
                                        </div>
                                    </div>
                                </div>
                            </div>

                            <CardContent className="p-4 grid gap-2">
                                <div className="flex items-center justify-between text-sm">
                                    <div className="flex items-center text-yellow-500 font-medium">
                                        <Star className="w-4 h-4 mr-1 fill-current" />
                                        {player.rating.toFixed(1)}
                                    </div>
                                    <div className="text-muted-foreground flex items-center">
                                        <ShoppingBag className="w-3 h-3 mr-1" />
                                        {player.orderCount} Orders
                                    </div>
                                </div>

                                <div className="flex items-center justify-between mt-2">
                                    <div className="flex items-center text-sm text-primary font-medium bg-primary/10 px-2 py-1 rounded">
                                        <Gamepad2 className="w-3 h-3 mr-1" />
                                        {player.gameName}
                                    </div>
                                    <div className="text-lg font-bold text-foreground">
                                        ${player.price} <span className="text-xs text-muted-foreground font-normal">/hr</span>
                                    </div>
                                </div>

                                <div className="flex flex-wrap gap-1 mt-2">
                                    {player.tags.slice(0, 3).map(tag => (
                                        <span key={tag} className="text-[10px] px-1.5 py-0.5 rounded-full bg-secondary text-secondary-foreground border border-transparent">
                                            {tag}
                                        </span>
                                    ))}
                                </div>
                            </CardContent>

                            <CardFooter className="p-3 bg-muted/30">
                                <Button className="w-full h-9 text-sm" variant="default">
                                    Book Now
                                </Button>
                            </CardFooter>
                        </Card>
                    ))}
                </div>
            )}

            {!loading && players.length === 0 && (
                <div className="text-center py-20 text-muted-foreground">
                    No players found matching your criteria.
                </div>
            )}
        </PageContainer>
    );
}
