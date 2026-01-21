import { useEffect, useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import {
    ArrowLeft,
    Users,
    Lock,
    Gamepad2,
    Play,
    Square,
    LogOut,
    Settings,
    Check,
} from 'lucide-react';
import { PageContainer } from '@/components/page-container';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import {
    AlertDialog,
    AlertDialogAction,
    AlertDialogCancel,
    AlertDialogContent,
    AlertDialogDescription,
    AlertDialogFooter,
    AlertDialogHeader,
    AlertDialogTitle,
    AlertDialogTrigger,
} from '@/components/ui/alert-dialog';
import {
    Dialog,
    DialogContent,
    DialogDescription,
    DialogFooter,
    DialogHeader,
    DialogTitle,
} from '@/components/ui/dialog';
import { RoomMemberList, RoomMemberListSkeleton } from '@/components/room';
import { VoicePanel } from '@/components/voice';
import { useRoomStore, useAuthStore, useVoiceStore } from '@/stores';
import { cn } from '@/lib/utils';
import type { ChatGroupStatus } from '@/stores/modules/room-store';

const statusConfig: Record<ChatGroupStatus, { label: string; color: string }> = {
    waiting: { label: 'room.status.waiting', color: 'bg-green-500' },
    ready: { label: 'room.status.ready', color: 'bg-blue-500' },
    in_game: { label: 'room.status.inGame', color: 'bg-orange-500' },
    finished: { label: 'room.status.finished', color: 'bg-gray-500' },
    canceled: { label: 'room.status.canceled', color: 'bg-red-500' },
};

export default function RoomDetailPage() {
    const { t } = useTranslation();
    const { id } = useParams<{ id: string }>();
    const navigate = useNavigate();
    const roomId = Number(id);

    const { user } = useAuthStore();
    const {
        currentRoom,
        members,
        isLoading,
        fetchRoom,
        fetchMembers,
        joinRoom,
        leaveRoom,
        toggleReady,
        startGame,
        finishGame,
        kickMember,
        closeRoom,
    } = useRoomStore();
    const { startVoice, stopVoice } = useVoiceStore();

    const [isJoining, setIsJoining] = useState(false);
    const [password, setPassword] = useState('');
    const [showPasswordDialog, setShowPasswordDialog] = useState(false);

    useEffect(() => {
        if (roomId) {
            fetchRoom(roomId);
            fetchMembers(roomId);
        }
    }, [roomId, fetchRoom, fetchMembers]);

    const isHost = currentRoom?.createdBy === Number(user?.id);
    const isMember = members.some((m) => m.userId === Number(user?.id) && m.isActive);
    const currentMember = members.find((m) => m.userId === Number(user?.id));
    const status = currentRoom ? statusConfig[currentRoom.roomStatus] : statusConfig.waiting;

    const handleJoin = async () => {
        if (!roomId) return;

        // If room is private and no password entered yet, show dialog
        if (currentRoom?.isPrivate && !password) {
            setShowPasswordDialog(true);
            return;
        }

        setIsJoining(true);
        try {
            await joinRoom(roomId, password);
            setPassword('');
            setShowPasswordDialog(false);
        } finally {
            setIsJoining(false);
        }
    };

    const handlePasswordSubmit = async () => {
        if (!roomId || !password.trim()) return;
        setIsJoining(true);
        try {
            await joinRoom(roomId, password);
            setPassword('');
            setShowPasswordDialog(false);
        } finally {
            setIsJoining(false);
        }
    };

    const handleLeave = async () => {
        if (!roomId) return;
        await leaveRoom(roomId);
        navigate('/rooms');
    };

    const handleToggleReady = async () => {
        if (!roomId) return;
        await toggleReady(roomId);
    };

    const handleStartGame = async () => {
        if (!roomId) return;
        await startGame(roomId);
    };

    const handleFinishGame = async () => {
        if (!roomId) return;
        await finishGame(roomId);
    };

    const handleKick = async (userId: number) => {
        if (!roomId) return;
        await kickMember(roomId, userId);
    };

    const handleClose = async () => {
        if (!roomId) return;
        await closeRoom(roomId);
        navigate('/rooms');
    };

    const handleStartVoice = async () => {
        if (!roomId) return;
        await startVoice(roomId);
        await fetchRoom(roomId);
    };

    const handleStopVoice = async () => {
        if (!roomId) return;
        await stopVoice(roomId);
        await fetchRoom(roomId);
    };

    if (isLoading || !currentRoom) {
        return (
            <PageContainer>
                <div className="space-y-4">
                    <div className="h-8 w-48 bg-muted rounded animate-pulse" />
                    <div className="h-32 bg-muted rounded animate-pulse" />
                    <RoomMemberListSkeleton />
                </div>
            </PageContainer>
        );
    }

    return (
        <PageContainer>
            {/* Header */}
            <div className="flex items-center gap-4 mb-6">
                <Button variant="ghost" size="icon" onClick={() => navigate('/rooms')}>
                    <ArrowLeft className="h-5 w-5" />
                </Button>
                <div className="flex-1">
                    <div className="flex items-center gap-2">
                        <h1 className="text-xl font-bold text-foreground">{currentRoom.name}</h1>
                        {currentRoom.isPrivate && <Lock className="h-4 w-4 text-muted-foreground" />}
                    </div>
                    {currentRoom.gameName && (
                        <div className="flex items-center gap-1.5 text-sm text-muted-foreground">
                            <Gamepad2 className="h-4 w-4" />
                            <span>{currentRoom.gameName}</span>
                        </div>
                    )}
                </div>
                <Badge variant="secondary" className={cn('text-white', status.color)}>
                    {t(status.label)}
                </Badge>
            </div>

            <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
                {/* Main Content */}
                <div className="lg:col-span-2 space-y-6">
                    {/* Room Info */}
                    <Card>
                        <CardHeader>
                            <CardTitle className="text-base">{t('room.info')}</CardTitle>
                        </CardHeader>
                        <CardContent>
                            {currentRoom.description && (
                                <p className="text-sm text-muted-foreground mb-4">
                                    {currentRoom.description}
                                </p>
                            )}
                            <div className="flex items-center gap-4 text-sm">
                                <div className="flex items-center gap-1.5">
                                    <Users className="h-4 w-4 text-muted-foreground" />
                                    <span>
                                        {currentRoom.currentMembers}/{currentRoom.maxMembers}
                                    </span>
                                </div>
                            </div>
                        </CardContent>
                    </Card>

                    {/* Members */}
                    <Card>
                        <CardHeader>
                            <CardTitle className="text-base">{t('room.members')}</CardTitle>
                        </CardHeader>
                        <CardContent>
                            <RoomMemberList
                                members={members}
                                hostUserId={currentRoom.createdBy}
                                currentUserId={Number(user?.id)}
                                isHost={isHost}
                                onKick={handleKick}
                            />
                        </CardContent>
                    </Card>
                </div>

                {/* Sidebar */}
                <div className="space-y-6">
                    {/* Actions */}
                    <Card>
                        <CardHeader>
                            <CardTitle className="text-base">{t('room.actions')}</CardTitle>
                        </CardHeader>
                        <CardContent className="space-y-3">
                            {!isMember ? (
                                /* Join button */
                                <Button className="w-full" onClick={handleJoin} disabled={isJoining}>
                                    {isJoining ? t('room.joining') : t('room.actions.join')}
                                </Button>
                            ) : (
                                <>
                                    {/* Ready toggle */}
                                    {currentRoom.roomStatus === 'waiting' && !isHost && (
                                        <Button
                                            variant={currentMember?.isReady ? 'secondary' : 'default'}
                                            className="w-full"
                                            onClick={handleToggleReady}
                                        >
                                            <Check className="h-4 w-4 mr-2" />
                                            {currentMember?.isReady ? t('room.actions.cancelReady') : t('room.actions.ready')}
                                        </Button>
                                    )}

                                    {/* Host actions */}
                                    {isHost && (
                                        <>
                                            {currentRoom.roomStatus === 'waiting' && (
                                                <Button className="w-full" onClick={handleStartGame}>
                                                    <Play className="h-4 w-4 mr-2" />
                                                    {t('room.actions.start')}
                                                </Button>
                                            )}
                                            {currentRoom.roomStatus === 'in_game' && (
                                                <Button className="w-full" onClick={handleFinishGame}>
                                                    <Square className="h-4 w-4 mr-2" />
                                                    {t('room.actions.finish')}
                                                </Button>
                                            )}
                                            <Button
                                                variant="outline"
                                                className="w-full"
                                                onClick={() => navigate(`/rooms/${roomId}/edit`)}
                                            >
                                                <Settings className="h-4 w-4 mr-2" />
                                                {t('room.detail.settings')}
                                            </Button>
                                        </>
                                    )}

                                    {/* Leave/Close */}
                                    {isHost ? (
                                        <AlertDialog>
                                            <AlertDialogTrigger asChild>
                                                <Button variant="destructive" className="w-full">
                                                    {t('room.actions.close')}
                                                </Button>
                                            </AlertDialogTrigger>
                                            <AlertDialogContent>
                                                <AlertDialogHeader>
                                                    <AlertDialogTitle>{t('room.dialogs.closeConfirm.title')}</AlertDialogTitle>
                                                    <AlertDialogDescription>
                                                        {t('room.dialogs.closeConfirm.description')}
                                                    </AlertDialogDescription>
                                                </AlertDialogHeader>
                                                <AlertDialogFooter>
                                                    <AlertDialogCancel>{t('common.cancel')}</AlertDialogCancel>
                                                    <AlertDialogAction onClick={handleClose}>
                                                        {t('common.confirm')}
                                                    </AlertDialogAction>
                                                </AlertDialogFooter>
                                            </AlertDialogContent>
                                        </AlertDialog>
                                    ) : (
                                        <Button variant="outline" className="w-full" onClick={handleLeave}>
                                            <LogOut className="h-4 w-4 mr-2" />
                                            {t('room.actions.leave')}
                                        </Button>
                                    )}
                                </>
                            )}
                        </CardContent>
                    </Card>

                    {/* Voice Panel */}
                    {isMember && (
                        <VoicePanel
                            roomId={roomId}
                            voiceEnabled={currentRoom.voiceEnabled}
                            isHost={isHost}
                            onStartVoice={handleStartVoice}
                            onStopVoice={handleStopVoice}
                        />
                    )}
                </div>
            </div>

            {/* Password Dialog for Private Rooms */}
            <Dialog open={showPasswordDialog} onOpenChange={setShowPasswordDialog}>
                <DialogContent>
                    <DialogHeader>
                        <DialogTitle>{t('room.dialogs.enterPassword.title')}</DialogTitle>
                        <DialogDescription>
                            {t('room.dialogs.enterPassword.description')}
                        </DialogDescription>
                    </DialogHeader>
                    <div className="space-y-4 py-4">
                        <div className="space-y-2">
                            <Label htmlFor="room-password">{t('room.password')}</Label>
                            <Input
                                id="room-password"
                                type="password"
                                value={password}
                                onChange={(e) => setPassword(e.target.value)}
                                placeholder={t('room.dialogs.enterPassword.placeholder')}
                                onKeyDown={(e) => {
                                    if (e.key === 'Enter') {
                                        handlePasswordSubmit();
                                    }
                                }}
                            />
                        </div>
                    </div>
                    <DialogFooter>
                        <Button
                            variant="outline"
                            onClick={() => {
                                setShowPasswordDialog(false);
                                setPassword('');
                            }}
                        >
                            {t('common.cancel')}
                        </Button>
                        <Button onClick={handlePasswordSubmit} disabled={isJoining || !password.trim()}>
                            {isJoining ? t('room.joining') : t('room.actions.join')}
                        </Button>
                    </DialogFooter>
                </DialogContent>
            </Dialog>
        </PageContainer>
    );
}
