import { useState, useEffect, useCallback } from 'react';
import { useTranslation } from 'react-i18next';
import { Search, Gamepad2, ChevronDown, X } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import {
    Dialog,
    DialogContent,
    DialogHeader,
    DialogTitle,
} from '@/components/ui/dialog';
import { ScrollArea } from '@/components/ui/scroll-area';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import { cn } from '@/lib/utils';
import { http } from '@/lib/http';

export interface Game {
    id: number;
    name: string;
    icon: string;
    category?: string;
}

interface GameSelectorProps {
    value?: number;
    onChange: (gameId: number, game: Game) => void;
    placeholder?: string;
    disabled?: boolean;
    error?: string;
}

interface GameListResponse {
    games: Game[];
    total: number;
    page: number;
    pageSize: number;
}

export function GameSelector({
    value,
    onChange,
    placeholder,
    disabled = false,
    error,
}: GameSelectorProps) {
    const { t } = useTranslation();
    const [open, setOpen] = useState(false);
    const [search, setSearch] = useState('');
    const [games, setGames] = useState<Game[]>([]);
    const [isLoading, setIsLoading] = useState(false);
    const [selectedGame, setSelectedGame] = useState<Game | null>(null);

    // Fetch games from API
    const fetchGames = useCallback(async (keyword?: string) => {
        setIsLoading(true);
        try {
            const params = new URLSearchParams();
            params.append('pageSize', '50');
            if (keyword) {
                params.append('keyword', keyword);
            }
            const response = await http.get<GameListResponse>(
                `/public/games?${params.toString()}`
            );
            setGames(response.games || []);
        } catch (err) {
            console.error('Failed to fetch games:', err);
            setGames([]);
        } finally {
            setIsLoading(false);
        }
    }, []);

    // Fetch games on mount and when dialog opens
    useEffect(() => {
        if (open) {
            fetchGames();
        }
    }, [open, fetchGames]);

    // Fetch selected game info if value is provided but no selectedGame
    useEffect(() => {
        if (value && !selectedGame) {
            http.get<Game>(`/public/games/${value}`)
                .then((game) => setSelectedGame(game))
                .catch(() => {});
        }
    }, [value, selectedGame]);

    // Debounced search
    useEffect(() => {
        if (!open) return;
        const timer = setTimeout(() => {
            fetchGames(search || undefined);
        }, 300);
        return () => clearTimeout(timer);
    }, [search, open, fetchGames]);

    const handleSelect = (game: Game) => {
        setSelectedGame(game);
        onChange(game.id, game);
        setOpen(false);
        setSearch('');
    };

    const handleClear = (e: React.MouseEvent) => {
        e.stopPropagation();
        setSelectedGame(null);
        onChange(0, { id: 0, name: '', icon: '' });
    };

    return (
        <>
            <Button
                type="button"
                variant="outline"
                role="combobox"
                aria-expanded={open}
                disabled={disabled}
                onClick={() => setOpen(true)}
                className={cn(
                    'w-full justify-between font-normal',
                    !selectedGame && 'text-muted-foreground',
                    error && 'border-destructive'
                )}
            >
                <div className="flex items-center gap-2 truncate">
                    {selectedGame ? (
                        <>
                            <Avatar className="h-5 w-5">
                                <AvatarImage src={selectedGame.icon} />
                                <AvatarFallback>
                                    <Gamepad2 className="h-3 w-3" />
                                </AvatarFallback>
                            </Avatar>
                            <span className="truncate">{selectedGame.name}</span>
                        </>
                    ) : (
                        <>
                            <Gamepad2 className="h-4 w-4" />
                            <span>{placeholder || t('game.selectGame')}</span>
                        </>
                    )}
                </div>
                <div className="flex items-center gap-1">
                    {selectedGame && (
                        <X
                            className="h-4 w-4 opacity-50 hover:opacity-100"
                            onClick={handleClear}
                        />
                    )}
                    <ChevronDown className="h-4 w-4 opacity-50" />
                </div>
            </Button>

            <Dialog open={open} onOpenChange={setOpen}>
                <DialogContent className="sm:max-w-md">
                    <DialogHeader>
                        <DialogTitle>{t('game.selectGame')}</DialogTitle>
                    </DialogHeader>
                    <div className="space-y-4">
                        {/* Search Input */}
                        <div className="relative">
                            <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
                            <Input
                                placeholder={t('game.searchPlaceholder')}
                                value={search}
                                onChange={(e) => setSearch(e.target.value)}
                                className="pl-9"
                                autoFocus
                            />
                        </div>

                        {/* Game List */}
                        <ScrollArea className="h-[300px]">
                            {isLoading ? (
                                <div className="space-y-2 p-1">
                                    {[1, 2, 3, 4, 5].map((i) => (
                                        <div
                                            key={i}
                                            className="flex items-center gap-3 p-2 rounded-md"
                                        >
                                            <div className="h-8 w-8 rounded bg-muted animate-pulse" />
                                            <div className="h-4 w-32 bg-muted rounded animate-pulse" />
                                        </div>
                                    ))}
                                </div>
                            ) : games.length === 0 ? (
                                <div className="text-center py-8 text-muted-foreground">
                                    {search ? t('game.noResults') : t('game.noGames')}
                                </div>
                            ) : (
                                <div className="space-y-1 p-1">
                                    {games.map((game) => (
                                        <button
                                            key={game.id}
                                            type="button"
                                            onClick={() => handleSelect(game)}
                                            className={cn(
                                                'w-full flex items-center gap-3 p-2 rounded-md',
                                                'hover:bg-accent transition-colors text-left',
                                                selectedGame?.id === game.id && 'bg-accent'
                                            )}
                                        >
                                            <Avatar className="h-8 w-8">
                                                <AvatarImage src={game.icon} />
                                                <AvatarFallback>
                                                    <Gamepad2 className="h-4 w-4" />
                                                </AvatarFallback>
                                            </Avatar>
                                            <div className="flex-1 min-w-0">
                                                <div className="font-medium truncate">
                                                    {game.name}
                                                </div>
                                                {game.category && (
                                                    <div className="text-xs text-muted-foreground">
                                                        {game.category}
                                                    </div>
                                                )}
                                            </div>
                                        </button>
                                    ))}
                                </div>
                            )}
                        </ScrollArea>
                    </div>
                </DialogContent>
            </Dialog>
        </>
    );
}
