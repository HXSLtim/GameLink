import { useNavigate } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { Users, Clock, Gamepad2, Trophy, DollarSign } from 'lucide-react';
import { Card, CardContent } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import { cn } from '@/lib/utils';
import type { LFGRequest, LFGRequestType, LFGRequestStatus } from '@/stores/modules/lfg-store';

interface LFGRequestCardProps {
    request: LFGRequest;
    currentUserId?: number;
    onAccept?: (requestId: number) => void;
    onCancel?: (requestId: number) => void;
    isAccepting?: boolean;
}

const typeConfig: Record<LFGRequestType, { label: string; color: string }> = {
    find_player: { label: 'lfg.type.findPlayer', color: 'bg-blue-500' },
    find_team: { label: 'lfg.type.findTeam', color: 'bg-purple-500' },
};

const statusConfig: Record<LFGRequestStatus, { label: string; color: string }> = {
    pending: { label: 'lfg.status.pending', color: 'bg-green-500' },
    matched: { label: 'lfg.status.matched', color: 'bg-blue-500' },
    expired: { label: 'lfg.status.expired', color: 'bg-gray-500' },
    canceled: { label: 'lfg.status.canceled', color: 'bg-red-500' },
};

export function LFGRequestCard({
    request,
    currentUserId,
    onAccept,
    onCancel,
    isAccepting = false,
}: LFGRequestCardProps) {
    const { t } = useTranslation();
    const navigate = useNavigate();

    const typeInfo = typeConfig[request.requestType] || typeConfig.find_player;
    const statusInfo = statusConfig[request.status] || statusConfig.pending;
    const isOwner = currentUserId === request.userId;
    const canAccept = request.status === 'pending' && !isOwner;
    const canCancel = request.status === 'pending' && isOwner;

    // Calculate time remaining
    const expiresAt = new Date(request.expiresAt);
    const now = new Date();
    const timeRemaining = Math.max(0, Math.floor((expiresAt.getTime() - now.getTime()) / 60000));

    const formatPrice = (cents: number) => {
        return `¥${(cents / 100).toFixed(0)}`;
    };

    return (
        <Card
            className={cn(
                'transition-all hover:shadow-md',
                'border border-border/50 bg-card/80 backdrop-blur-sm',
                request.status !== 'pending' && 'opacity-75'
            )}
        >
            <CardContent className="p-4">
                {/* Header */}
                <div className="flex items-start justify-between mb-3">
                    <div className="flex items-center gap-3">
                        <Avatar className="h-10 w-10">
                            <AvatarImage src={request.userAvatarUrl} />
                            <AvatarFallback>
                                {request.userNickname?.charAt(0) || 'U'}
                            </AvatarFallback>
                        </Avatar>
                        <div>
                            <div className="flex items-center gap-2">
                                <span className="font-medium text-foreground">
                                    {request.userNickname}
                                </span>
                                {isOwner && (
                                    <Badge variant="outline" className="text-xs">
                                        {t('lfg.you')}
                                    </Badge>
                                )}
                            </div>
                            <div className="flex items-center gap-2 text-sm text-muted-foreground">
                                <Badge
                                    variant="secondary"
                                    className={cn('text-xs text-white', typeInfo.color)}
                                >
                                    {t(typeInfo.label)}
                                </Badge>
                            </div>
                        </div>
                    </div>
                    <Badge
                        variant="secondary"
                        className={cn('text-xs text-white', statusInfo.color)}
                    >
                        {t(statusInfo.label)}
                    </Badge>
                </div>

                {/* Title & Description */}
                {request.title && (
                    <h3 className="font-semibold text-foreground mb-1">{request.title}</h3>
                )}
                {request.description && (
                    <p className="text-sm text-muted-foreground mb-3 line-clamp-2">
                        {request.description}
                    </p>
                )}

                {/* Info Grid */}
                <div className="grid grid-cols-2 gap-2 mb-3">
                    {/* Game */}
                    {request.gameName && (
                        <div className="flex items-center gap-1.5 text-sm text-muted-foreground">
                            <Gamepad2 className="h-4 w-4" />
                            <span className="truncate">{request.gameName}</span>
                        </div>
                    )}

                    {/* Required Players */}
                    <div className="flex items-center gap-1.5 text-sm text-muted-foreground">
                        <Users className="h-4 w-4" />
                        <span>
                            {t('lfg.requiredPlayers', { count: request.requiredPlayers })}
                        </span>
                    </div>

                    {/* Min Rank */}
                    {request.minRank && (
                        <div className="flex items-center gap-1.5 text-sm text-muted-foreground">
                            <Trophy className="h-4 w-4" />
                            <span>{request.minRank}</span>
                        </div>
                    )}

                    {/* Max Price */}
                    {request.maxPriceCents && request.maxPriceCents > 0 && (
                        <div className="flex items-center gap-1.5 text-sm text-muted-foreground">
                            <DollarSign className="h-4 w-4" />
                            <span>{t('lfg.maxPrice', { price: formatPrice(request.maxPriceCents) })}</span>
                        </div>
                    )}
                </div>

                {/* Footer */}
                <div className="flex items-center justify-between">
                    {/* Time remaining */}
                    {request.status === 'pending' && (
                        <div className="flex items-center gap-1.5 text-sm">
                            <Clock className="h-4 w-4 text-muted-foreground" />
                            <span
                                className={cn(
                                    timeRemaining <= 5 ? 'text-red-500' : 'text-muted-foreground'
                                )}
                            >
                                {timeRemaining > 0
                                    ? t('lfg.timeRemaining', { minutes: timeRemaining })
                                    : t('lfg.expiringSoon')}
                            </span>
                        </div>
                    )}
                    {request.status !== 'pending' && <div />}

                    {/* Actions */}
                    <div className="flex items-center gap-2">
                        {canCancel && onCancel && (
                            <Button
                                variant="outline"
                                size="sm"
                                onClick={() => onCancel(request.id)}
                            >
                                {t('lfg.cancel')}
                            </Button>
                        )}
                        {canAccept && onAccept && (
                            <Button
                                size="sm"
                                onClick={() => onAccept(request.id)}
                                disabled={isAccepting}
                            >
                                {isAccepting ? t('lfg.accepting') : t('lfg.accept')}
                            </Button>
                        )}
                        {request.status === 'matched' && request.matchedRoomId && (
                            <Button
                                size="sm"
                                onClick={() => navigate(`/rooms/${request.matchedRoomId}`)}
                            >
                                {t('lfg.goToRoom')}
                            </Button>
                        )}
                    </div>
                </div>
            </CardContent>
        </Card>
    );
}

export function LFGRequestCardSkeleton() {
    return (
        <Card className="border border-border/50 bg-card/80">
            <CardContent className="p-4">
                <div className="flex items-start justify-between mb-3">
                    <div className="flex items-center gap-3">
                        <div className="h-10 w-10 rounded-full bg-muted animate-pulse" />
                        <div>
                            <div className="h-4 w-24 bg-muted rounded animate-pulse mb-1" />
                            <div className="h-5 w-20 bg-muted rounded animate-pulse" />
                        </div>
                    </div>
                    <div className="h-5 w-16 bg-muted rounded animate-pulse" />
                </div>
                <div className="h-5 w-48 bg-muted rounded animate-pulse mb-2" />
                <div className="h-4 w-full bg-muted rounded animate-pulse mb-3" />
                <div className="grid grid-cols-2 gap-2 mb-3">
                    <div className="h-4 w-24 bg-muted rounded animate-pulse" />
                    <div className="h-4 w-20 bg-muted rounded animate-pulse" />
                </div>
                <div className="flex items-center justify-between">
                    <div className="h-4 w-24 bg-muted rounded animate-pulse" />
                    <div className="h-8 w-20 bg-muted rounded animate-pulse" />
                </div>
            </CardContent>
        </Card>
    );
}
