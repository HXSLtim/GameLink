import { useState, useEffect, useCallback } from 'react';
import { useTranslation } from 'react-i18next';
import { PageContainer } from '@/components/page-container';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '@/components/ui/card';
import { Label } from '@/components/ui/label';
import { Switch } from '@/components/ui/switch';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { toast } from 'sonner';
import { Clock, Image as ImageIcon, Mic, Loader2 } from 'lucide-react';
import { usePlayerStore, type TimeSlot } from '@/stores/modules/player-store';

// Day mapping: 0=Sunday, 1=Monday, ..., 6=Saturday
const DAY_KEYS = ['sun', 'mon', 'tue', 'wed', 'thu', 'fri', 'sat'] as const;
type DayKey = typeof DAY_KEYS[number];

export default function EditPlayerProfilePage() {
    const { t } = useTranslation();
    const { myProfile, loading, fetchMyProfile, updateProfile } = usePlayerStore();
    const [saving, setSaving] = useState(false);
    const [initialLoading, setInitialLoading] = useState(true);

    // Availability state derived from serviceTimeSlots
    const [availability, setAvailability] = useState<Record<DayKey, boolean>>({
        sun: false,
        mon: false,
        tue: false,
        wed: false,
        thu: false,
        fri: false,
        sat: false,
    });

    // Album state from gallery
    const [album, setAlbum] = useState<string[]>([]);

    // Load profile on mount
    useEffect(() => {
        const loadProfile = async () => {
            try {
                await fetchMyProfile();
            } finally {
                setInitialLoading(false);
            }
        };
        loadProfile();
    }, [fetchMyProfile]);

    // Sync state when profile loads
    useEffect(() => {
        if (myProfile) {
            // Convert serviceTimeSlots to availability map
            const newAvailability: Record<DayKey, boolean> = {
                sun: false, mon: false, tue: false, wed: false, thu: false, fri: false, sat: false,
            };
            myProfile.serviceTimeSlots?.forEach((slot: TimeSlot) => {
                const dayKey = DAY_KEYS[slot.dayOfWeek];
                if (dayKey) newAvailability[dayKey] = true;
            });
            setAvailability(newAvailability);

            // Set album from gallery
            setAlbum(myProfile.gallery || []);
        }
    }, [myProfile]);

    const handleSave = useCallback(async () => {
        setSaving(true);
        try {
            // Convert availability back to TimeSlots
            const serviceTimeSlots: TimeSlot[] = DAY_KEYS
                .map((day, index) => ({ dayOfWeek: index, day }))
                .filter(({ day }) => availability[day])
                .map(({ dayOfWeek }) => ({
                    dayOfWeek,
                    startTime: '09:00',
                    endTime: '22:00',
                }));

            await updateProfile({
                serviceTimeSlots,
                gallery: album,
            });
            toast.success(t('player.edit.success', { defaultValue: 'Profile updated successfully!' }));
        } catch {
            toast.error(t('player.edit.error', { defaultValue: 'Failed to update profile' }));
        } finally {
            setSaving(false);
        }
    }, [availability, album, updateProfile, t]);

    const toggleDay = (day: DayKey) => {
        setAvailability(prev => ({ ...prev, [day]: !prev[day] }));
    };

    const handleDeletePhoto = (index: number) => {
        setAlbum(prev => prev.filter((_, i) => i !== index));
    };

    // Show loading state while fetching profile
    if (initialLoading || loading) {
        return (
            <PageContainer>
                <div className="flex items-center justify-center min-h-[400px]">
                    <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
                </div>
            </PageContainer>
        );
    }

    // Show error if no profile found
    if (!myProfile) {
        return (
            <PageContainer>
                <div className="flex flex-col items-center justify-center min-h-[400px] space-y-4">
                    <p className="text-muted-foreground">{t('player.edit.no_profile', { defaultValue: 'Player profile not found. Please apply to become a player first.' })}</p>
                </div>
            </PageContainer>
        );
    }

    return (
        <PageContainer>
            <div className="max-w-4xl mx-auto py-8 px-4 space-y-6">
                <div className="flex justify-between items-center">
                    <div>
                        <h1 className="text-3xl font-bold tracking-tight">{t('player.edit.title', { defaultValue: 'Player Settings' })}</h1>
                        <p className="text-muted-foreground">{t('player.edit.subtitle', { defaultValue: 'Manage your service availability and media.' })}</p>
                    </div>
                    <Button onClick={handleSave} disabled={saving}>
                        {saving ? (
                            <>
                                <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                                {t('common.saving', { defaultValue: 'Saving...' })}
                            </>
                        ) : t('common.save', { defaultValue: 'Save Changes' })}
                    </Button>
                </div>

                <Tabs defaultValue="availability" className="space-y-4">
                    <TabsList>
                        <TabsTrigger value="availability" className="flex items-center gap-2">
                            <Clock className="w-4 h-4" />
                            {t('player.edit.availability', { defaultValue: 'Availability' })}
                        </TabsTrigger>
                        <TabsTrigger value="media" className="flex items-center gap-2">
                            <ImageIcon className="w-4 h-4" />
                            {t('player.edit.media', { defaultValue: 'Media & Album' })}
                        </TabsTrigger>
                        <TabsTrigger value="voice" className="flex items-center gap-2">
                            <Mic className="w-4 h-4" />
                            {t('player.edit.voice', { defaultValue: 'Voice Intro' })}
                        </TabsTrigger>
                    </TabsList>

                    <TabsContent value="availability">
                        <Card>
                            <CardHeader>
                                <CardTitle>{t('player.edit.weekly_schedule', { defaultValue: 'Weekly Schedule' })}</CardTitle>
                                <CardDescription>{t('player.edit.schedule_desc', { defaultValue: 'Toggle the days you are available to take orders.' })}</CardDescription>
                            </CardHeader>
                            <CardContent className="space-y-6">
                                {Object.entries(availability).map(([day, isAvailable]) => (
                                    <div key={day} className="flex items-center justify-between border-b pb-4 last:border-0 last:pb-0">
                                        <div className="space-y-0.5">
                                            <Label className="text-base capitalize">{day}</Label>
                                            <p className="text-sm text-muted-foreground">
                                                {isAvailable
                                                    ? t('player.edit.open_all_day', { defaultValue: 'Available all day' })
                                                    : t('player.edit.unavailable', { defaultValue: 'Unavailable' })
                                                }
                                            </p>
                                        </div>
                                        <Switch
                                            checked={isAvailable}
                                            onCheckedChange={() => toggleDay(day as keyof typeof availability)}
                                        />
                                    </div>
                                ))}
                            </CardContent>
                        </Card>
                    </TabsContent>

                    <TabsContent value="media">
                        <Card>
                            <CardHeader>
                                <CardTitle>{t('player.edit.album', { defaultValue: 'Photo Album' })}</CardTitle>
                                <CardDescription>{t('player.edit.album_desc', { defaultValue: 'Showcase your best moments. Max 9 photos.' })}</CardDescription>
                            </CardHeader>
                            <CardContent>
                                <div className="grid grid-cols-3 md:grid-cols-4 gap-4">
                                    {album.map((src, index) => (
                                        <div key={index} className="aspect-square relative rounded-lg overflow-hidden border bg-muted group">
                                            <img src={src} alt={`Album ${index}`} className="w-full h-full object-cover transition-transform group-hover:scale-105" />
                                            <Button
                                                variant="destructive"
                                                size="icon"
                                                className="absolute top-2 right-2 h-6 w-6 opacity-0 group-hover:opacity-100 transition-opacity"
                                                onClick={() => handleDeletePhoto(index)}
                                            >
                                                <span className="sr-only">Delete</span>
                                                &times;
                                            </Button>
                                        </div>
                                    ))}
                                    <div className="aspect-square rounded-lg border-2 border-dashed flex flex-col items-center justify-center cursor-pointer hover:bg-muted/50 transition-colors">
                                        <ImageIcon className="h-8 w-8 text-muted-foreground mb-2" />
                                        <span className="text-xs text-muted-foreground">{t('common.upload', { defaultValue: 'Upload' })}</span>
                                    </div>
                                </div>
                            </CardContent>
                        </Card>
                    </TabsContent>

                    <TabsContent value="voice">
                        <Card>
                            <CardHeader>
                                <CardTitle>{t('player.edit.voice_intro', { defaultValue: 'Voice Introduction' })}</CardTitle>
                                <CardDescription>{t('player.edit.voice_desc', { defaultValue: 'Record a greeting for your profile.' })}</CardDescription>
                            </CardHeader>
                            <CardContent className="flex flex-col items-center justify-center py-10 space-y-4">
                                <div className="w-16 h-16 rounded-full bg-primary/10 flex items-center justify-center text-primary">
                                    <Mic className="h-8 w-8" />
                                </div>
                                <Button variant="outline">
                                    {t('player.edit.start_recording', { defaultValue: 'Click to Record' })}
                                </Button>
                                <p className="text-xs text-muted-foreground">Max 30 seconds</p>
                            </CardContent>
                        </Card>
                    </TabsContent>
                </Tabs>
            </div>
        </PageContainer>
    );
}
