import { useState } from 'react';
import { Button } from '@/components/ui/button';
import { Heart } from 'lucide-react';
import { cn } from '@/lib/utils';
import { http } from '@/lib/http';
import { toast } from 'sonner';

interface FavoriteButtonProps {
    playerId: number;
    initialIsFavorite?: boolean;
    className?: string;
    onToggle?: (isFavorite: boolean) => void;
}

export function FavoriteButton({ playerId, initialIsFavorite = false, className, onToggle }: FavoriteButtonProps) {
    const [isFavorite, setIsFavorite] = useState(initialIsFavorite);
    const [loading, setLoading] = useState(false);

    const handleToggle = async (e: React.MouseEvent) => {
        e.preventDefault();
        e.stopPropagation();

        if (loading) return;

        setLoading(true);
        // Optimistic update
        const newState = !isFavorite;
        setIsFavorite(newState);

        try {
            if (newState) {
                await http.post(`/user/favorites/players/${playerId}`);
                toast.success('Added to favorites', { position: 'bottom-center' }); // Simple toast
            } else {
                await http.delete(`/user/favorites/players/${playerId}`);
                toast.success('Removed from favorites', { position: 'bottom-center' });
            }

            if (onToggle) onToggle(newState);
        } catch (err) {
            // Revert state on error
            setIsFavorite(!newState);
            toast.error('Failed to update favorites');
            console.error(err);
        } finally {
            setLoading(false);
        }
    };

    return (
        <Button
            variant="ghost"
            size="icon"
            className={cn(
                "rounded-full transition-all duration-300 hover:bg-pink-500/10 hover:text-pink-500",
                isFavorite ? "text-pink-500 bg-pink-500/10" : "text-muted-foreground",
                className
            )}
            onClick={handleToggle}
            disabled={loading}
        >
            <Heart className={cn("h-5 w-5 transition-all", isFavorite && "fill-current scale-110")} />
            <span className="sr-only">Toggle Favorite</span>
        </Button>
    );
}
