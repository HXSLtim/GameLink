import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { PageContainer } from '@/components/page-container';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '@/components/ui/card';
import { Label } from '@/components/ui/label';
import { Switch } from '@/components/ui/switch';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { toast } from 'sonner';
import { Clock, Image as ImageIcon, Mic } from 'lucide-react';

export default function EditPlayerProfilePage() {
    const { t } = useTranslation();
    const [loading, setLoading] = useState(false);

    // Mock Availability State (True = Available)
    const [availability, setAvailability] = useState({
        mon: true,
        tue: true,
        wed: true,
        thu: true,
        fri: true,
        sat: true,
        sun: true,
    });

    const [album, setAlbum] = useState([
        'https://api.dicebear.com/7.x/adventurer/svg?seed=Felix',
        'https://api.dicebear.com/7.x/adventurer/svg?seed=Coco',
    ]);

    const handleSave = () => {
        setLoading(true);
        setTimeout(() => {
            setLoading(false);
            toast.success(t('player.edit.success', { defaultValue: 'Profile updated successfully!' }));
        }, 1000);
    };

    const toggleDay = (day: keyof typeof availability) => {
        setAvailability(prev => ({ ...prev, [day]: !prev[day] }));
    };

    const handleDeletePhoto = (index: number) => {
        setAlbum(prev => prev.filter((_, i) => i !== index));
    };

    return (
        <PageContainer>
            <div className="max-w-4xl mx-auto py-8 px-4 space-y-6">
                <div className="flex justify-between items-center">
                    <div>
                        <h1 className="text-3xl font-bold tracking-tight">{t('player.edit.title', { defaultValue: 'Player Settings' })}</h1>
                        <p className="text-muted-foreground">{t('player.edit.subtitle', { defaultValue: 'Manage your service availability and media.' })}</p>
                    </div>
                    <Button onClick={handleSave} disabled={loading}>
                        {loading ? t('common.saving', { defaultValue: 'Saving...' }) : t('common.save', { defaultValue: 'Save Changes' })}
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
