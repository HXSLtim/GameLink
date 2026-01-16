import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { Users, Lock, Mic, Gamepad2 } from 'lucide-react';
import { Card, CardContent } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import { Input } from '@/components/ui/input';
import {
    Dialog,
    DialogContent,
    DialogDescription,
    DialogFooter,
    DialogHeader,
    DialogTitle,
} from '@/components/ui/dialog';
import { cn } from '@/lib/utils';
import type { GameRoom, ChatGroupStatus } from '@/stores/modules/room-store';

interface RoomCardProps {
    room: GameRoom;
    currentUserId?: number;
    onJoin?: (roomId: number, password?: string) => Promise<void>;
    isJoining?: boolean;
}

const statusConfig: Record<ChatGroupStatus, { label: string; variant: 'default' | 'secondary' | 'destructive' | 'outline'; className: string }> = {
    waiting: { label: 'room.status.waiting', variant: 'default', className: 'bg-amber-500 hover:bg-amber-500' },
    ready: { label: 'room.status.ready', variant: 'default', className: 'bg-green-500 hover:bg-green-500' },
    in_game: { label: 'room.status.in_game', variant: 'default', className: 'bg-blue-500 hover:bg-blue-500' },
    finished: { label: 'room.status.finished', variant: 'secondary', className: '' },
    canceled: { label: 'room.status.canceled', variant: 'destructive', className: '' },
};

export function RoomCard({ room, currentUserId, onJoin, isJoining }: RoomCardProps) {
    const { t } = useTranslation();
    const navigate = useNavigate();
    const [showPasswordDialog, setShowPasswordDialog] = useState(false);
    const [password, setPassword] = useState('');
    const [passwordError, setPasswordError] = useState('');

    const status = statusConfig[room.roomStatus] || statusConfig.waiting;
    const isFull = room.currentMembers >= room.maxMembers;
    const isHost = currentUserId === room.createdBy;
    const canJoin = room.roomStatus === 'waiting' && !isFull && !isHost;

    const handleCardClick = () => {
        if (isHost) {
            navigate(`/rooms/${room.id}`);
        }
    };

    const handleJoinClick = async (e: React.MouseEvent) => {
        e.stopPropagation();
        if (!onJoin) {
            navigate(`/rooms/${room.id}`);
            return;
        }

        if (room.isPrivate) {
            setShowPasswordDialog(true);
        } else {
            try {
                await onJoin(room.id);
                navigate(`/rooms/${room.id}`);
            } catch (error) {
                console.error('Failed to join room:', error);
            }
        }
    };

    const handlePasswordSubmit = async () => {
        if (!password.trim()) {
            setPasswordError(t('room.errors.passwordRequired'));
            return;
        }
        setPasswordError('');

        try {
            await onJoin?.(room.id, password);
            setShowPasswordDialog(false);
            setPassword('');
            navigate(`/rooms/${room.id}`);
        } catch (error: unknown) {
            setPasswordError(error instanceof Error ? error.message : t('room.errors.wrongPassword'));
        }
    };

    const getButtonContent = () => {
        if (isHost) return { text: t('room.enter'), disabled: false };
        if (isFull) return { text: t('room.full'), disabled: true };
        if (room.roomStatus !== 'waiting') return { text: t('room.status.' + room.roomStatus), disabled: true };
        return { text: t('room.join'), disabled: false };
    };

    const buttonContent = getButtonContent();

    return (
        <>
            <Card
                className={cn(
                    'transition-all hover:shadow-md',
                    'border border-border/50 bg-card/80 backdrop-blur-sm',
                    isHost && 'cursor-pointer hover:scale-[1.02]',
                    !canJoin && !isHost && 'opacity-75'
                )}
                onClick={handleCardClick}
            >
                <CardContent className="p-4">
                    {/* Header */}
                    <div className="flex items-start justify-between mb-3">
                        <div className="flex-1 min-w-0">
                            <div className="flex items-center gap-2 mb-1">
                                <h3 className="font-semibold text-foreground truncate">
                                    {room.name}
                                </h3>
                                {room.isPrivate && (
                                    <Lock className="h-3.5 w-3.5 text-muted-foreground flex-shrink-0" />
                                )}
                                {room.voiceEnabled && (
                                    <Mic className="h-3.5 w-3.5 text-green-500 flex-shrink-0" />
                                )}
                            </div>
                            {room.gameName && (
                                <div className="flex items-center gap-1.5 text-sm text-muted-foreground">
                                    <Gamepad2 className="h-3.5 w-3.5" />
                                    <span className="truncate">{room.gameName}</span>
                                </div>
                            )}
                        </div>
                        <Badge variant={status.variant} className={cn('text-xs text-white', status.className)}>
                            {t(status.label)}
                        </Badge>
                    </div>

                    {/* Description */}
                    {room.description && (
                        <p className="text-sm text-muted-foreground mb-3 line-clamp-2">
                            {room.description}
                        </p>
                    )}

                    {/* Footer */}
                    <div className="flex items-center justify-between">
                        {/* Host & Members */}
                        <div className="flex items-center gap-3">
                            <div className="flex items-center gap-2">
                                <Avatar className="h-6 w-6">
                                    <AvatarImage src={undefined} />
                                    <AvatarFallback className="text-xs">
                                        {room.hostNickname?.charAt(0) || 'H'}
                                    </AvatarFallback>
                                </Avatar>
                                <span className="text-sm text-muted-foreground truncate max-w-[80px]">
                                    {room.hostNickname || t('room.host')}
                                </span>
                            </div>
                            <div className="flex items-center gap-1">
                                <Users className="h-4 w-4 text-muted-foreground" />
                                <span className={cn('text-sm font-medium', isFull ? 'text-red-500' : 'text-muted-foreground')}>
                                    {room.currentMembers}/{room.maxMembers}
                                </span>
                            </div>
                        </div>

                        {/* Join Button */}
                        <Button
                            size="sm"
                            variant={buttonContent.disabled ? 'secondary' : 'default'}
                            disabled={buttonContent.disabled || isJoining}
                            onClick={handleJoinClick}
                        >
                            {isJoining ? t('room.joining') : buttonContent.text}
                        </Button>
                    </div>
                </CardContent>
            </Card>

            {/* Password Dialog */}
            <Dialog open={showPasswordDialog} onOpenChange={setShowPasswordDialog}>
                <DialogContent className="sm:max-w-md">
                    <DialogHeader>
                        <DialogTitle>{t('room.enterPassword')}</DialogTitle>
                        <DialogDescription>
                            {t('room.enterPasswordDescription')}
                        </DialogDescription>
                    </DialogHeader>
                    <div className="space-y-4">
                        <Input
                            type="password"
                            placeholder={t('room.passwordPlaceholder')}
                            value={password}
                            onChange={(e) => setPassword(e.target.value)}
                            onKeyDown={(e) => e.key === 'Enter' && handlePasswordSubmit()}
                        />
                        {passwordError && (
                            <p className="text-sm text-destructive">{passwordError}</p>
                        )}
                    </div>
                    <DialogFooter>
                        <Button variant="outline" onClick={() => setShowPasswordDialog(false)}>
                            {t('common.cancel')}
                        </Button>
                        <Button onClick={handlePasswordSubmit} disabled={isJoining}>
                            {isJoining ? t('room.joining') : t('room.join')}
                        </Button>
                    </DialogFooter>
                </DialogContent>
            </Dialog>
        </>
    );
}

export function RoomCardSkeleton() {
    return (
        <Card className="border border-border/50 bg-card/80">
            <CardContent className="p-4">
                <div className="flex items-start justify-between mb-3">
                    <div className="flex-1">
                        <div className="h-5 w-32 bg-muted rounded animate-pulse mb-2" />
                        <div className="h-4 w-24 bg-muted rounded animate-pulse" />
                    </div>
                    <div className="h-5 w-16 bg-muted rounded-full animate-pulse" />
                </div>
                <div className="h-4 w-full bg-muted rounded animate-pulse mb-3" />
                <div className="flex items-center justify-between">
                    <div className="flex items-center gap-3">
                        <div className="flex items-center gap-2">
                            <div className="h-6 w-6 bg-muted rounded-full animate-pulse" />
                            <div className="h-4 w-16 bg-muted rounded animate-pulse" />
                        </div>
                        <div className="h-4 w-12 bg-muted rounded animate-pulse" />
                    </div>
                    <div className="h-8 w-16 bg-muted rounded animate-pulse" />
                </div>
            </CardContent>
        </Card>
    );
}
