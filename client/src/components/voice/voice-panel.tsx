
import { useTranslation } from 'react-i18next';
import {
    Mic, MicOff, Headphones,
    VolumeX, PhoneOff, Settings
} from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuItem,
    DropdownMenuLabel,
    DropdownMenuSeparator,
    DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import { cn } from '@/lib/utils';
import { useVoiceStore } from '@/stores';
import { VoiceMemberItem } from './voice-member-item';

interface VoicePanelProps {
    roomId: number;
    voiceEnabled?: boolean;
    isHost?: boolean;
    onStartVoice?: () => void;
    onStopVoice?: () => void;
}

export function VoicePanel({
    roomId,
    voiceEnabled = false,
    isHost = false,
    onStartVoice,
    onStopVoice,
}: VoicePanelProps) {
    const { t } = useTranslation();
    const {
        isInVoice,
        isMuted,
        isDeafened,
        isConnecting,
        voiceMembers,
        audioInputDevices,
        audioOutputDevices,
        selectedInputDevice,
        selectedOutputDevice,
        joinVoice,
        leaveVoice,
        toggleMute,
        toggleDeafen,
        setInputDevice,
        setOutputDevice,
        refreshDevices,
    } = useVoiceStore();

    const handleJoinVoice = async () => {
        await refreshDevices();
        await joinVoice(roomId);
    };

    const handleLeaveVoice = async () => {
        await leaveVoice();
    };

    // Voice not enabled for this room
    if (!voiceEnabled) {
        return (
            <Card className="border-dashed">
                <CardContent className="py-6 text-center">
                    <Mic className="h-8 w-8 mx-auto mb-2 text-muted-foreground" />
                    <p className="text-sm text-muted-foreground mb-3">
                        {t('voice.notEnabled')}
                    </p>
                    {isHost && onStartVoice && (
                        <Button onClick={onStartVoice} size="sm">
                            {t('voice.enable')}
                        </Button>
                    )}
                </CardContent>
            </Card>
        );
    }

    return (
        <Card>
            <CardHeader className="pb-3">
                <div className="flex items-center justify-between">
                    <CardTitle className="text-base flex items-center gap-2">
                        <Mic className="h-4 w-4 text-green-500" />
                        {t('voice.title')}
                    </CardTitle>
                    {isHost && isInVoice && onStopVoice && (
                        <Button
                            variant="ghost"
                            size="sm"
                            className="text-destructive"
                            onClick={onStopVoice}
                        >
                            {t('voice.disable')}
                        </Button>
                    )}
                </div>
            </CardHeader>
            <CardContent className="space-y-4">
                {/* Voice members */}
                {isInVoice && voiceMembers.length > 0 && (
                    <div className="space-y-2">
                        {voiceMembers.map((member) => (
                            <VoiceMemberItem key={member.userId} member={member} />
                        ))}
                    </div>
                )}

                {/* Join/Leave button */}
                {!isInVoice ? (
                    <Button
                        className="w-full"
                        onClick={handleJoinVoice}
                        disabled={isConnecting}
                    >
                        {isConnecting ? t('voice.connecting') : t('voice.join')}
                    </Button>
                ) : (
                    <div className="space-y-3">
                        {/* Voice controls */}
                        <div className="flex items-center justify-center gap-2">
                            <Button
                                variant={isMuted ? 'destructive' : 'secondary'}
                                size="icon"
                                onClick={toggleMute}
                                className="h-10 w-10"
                            >
                                {isMuted ? (
                                    <MicOff className="h-5 w-5" />
                                ) : (
                                    <Mic className="h-5 w-5" />
                                )}
                            </Button>
                            <Button
                                variant={isDeafened ? 'destructive' : 'secondary'}
                                size="icon"
                                onClick={toggleDeafen}
                                className="h-10 w-10"
                            >
                                {isDeafened ? (
                                    <VolumeX className="h-5 w-5" />
                                ) : (
                                    <Headphones className="h-5 w-5" />
                                )}
                            </Button>
                            <DropdownMenu>
                                <DropdownMenuTrigger asChild>
                                    <Button variant="secondary" size="icon" className="h-10 w-10">
                                        <Settings className="h-5 w-5" />
                                    </Button>
                                </DropdownMenuTrigger>
                                <DropdownMenuContent align="center" className="w-56">
                                    <DropdownMenuLabel>{t('voice.inputDevice')}</DropdownMenuLabel>
                                    {audioInputDevices.map((device) => (
                                        <DropdownMenuItem
                                            key={device.deviceId}
                                            onClick={() => setInputDevice(device.deviceId)}
                                            className={cn(
                                                selectedInputDevice === device.deviceId &&
                                                'bg-accent'
                                            )}
                                        >
                                            {device.label || t('voice.defaultDevice')}
                                        </DropdownMenuItem>
                                    ))}
                                    <DropdownMenuSeparator />
                                    <DropdownMenuLabel>{t('voice.outputDevice')}</DropdownMenuLabel>
                                    {audioOutputDevices.map((device) => (
                                        <DropdownMenuItem
                                            key={device.deviceId}
                                            onClick={() => setOutputDevice(device.deviceId)}
                                            className={cn(
                                                selectedOutputDevice === device.deviceId &&
                                                'bg-accent'
                                            )}
                                        >
                                            {device.label || t('voice.defaultDevice')}
                                        </DropdownMenuItem>
                                    ))}
                                </DropdownMenuContent>
                            </DropdownMenu>
                            <Button
                                variant="destructive"
                                size="icon"
                                onClick={handleLeaveVoice}
                                className="h-10 w-10"
                            >
                                <PhoneOff className="h-5 w-5" />
                            </Button>
                        </div>

                        {/* Status text */}
                        <p className="text-xs text-center text-muted-foreground">
                            {isMuted && isDeafened
                                ? t('voice.mutedAndDeafened')
                                : isMuted
                                    ? t('voice.muted')
                                    : isDeafened
                                        ? t('voice.deafened')
                                        : t('voice.connected')}
                        </p>
                    </div>
                )}
            </CardContent>
        </Card>
    );
}
