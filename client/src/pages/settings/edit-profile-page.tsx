import { useState } from 'react';
import { useAuthStore } from '@/stores';
import { PageContainer } from '@/components/page-container';
import { Card, CardContent, CardFooter } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label'; // Assuming Label exists or I use standard label
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import { useTranslation } from 'react-i18next';
import { ArrowLeft, Save, Loader2, Camera, User, Mail } from 'lucide-react';
import { useNavigate } from 'react-router-dom';
import { toast } from 'sonner';

export default function EditProfilePage() {
    const { t } = useTranslation();
    const navigate = useNavigate();
    const { user, updateProfile } = useAuthStore();

    const [isLoading, setIsLoading] = useState(false);
    const [formData, setFormData] = useState({
        nickname: user?.nickname || '',
        username: user?.username || '',
        email: user?.email || '',
        avatar: user?.avatar || ''
    });

    const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const { name, value } = e.target;
        setFormData(prev => ({ ...prev, [name]: value }));
    };

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setIsLoading(true);
        try {
            await updateProfile({
                nickname: formData.nickname,
                email: formData.email,
                avatar: formData.avatar
            });
            toast.success(t('settings.update_success'));
            setTimeout(() => navigate('/profile'), 500);
        } catch (error: any) {
            toast.error(error.message || 'Failed to update profile');
        } finally {
            setIsLoading(false);
        }
    };

    return (
        <PageContainer>
            <div className="max-w-2xl mx-auto py-8 px-4 space-y-6">
                <div className="flex items-center gap-4 animate-in fade-in slide-in-from-top-2">
                    <Button variant="ghost" size="icon" onClick={() => navigate('/profile')} className="rounded-full">
                        <ArrowLeft className="h-5 w-5" />
                    </Button>
                    <div>
                        <h1 className="text-2xl font-bold tracking-tight">{t('settings.personal_info')}</h1>
                        <p className="text-sm text-muted-foreground">{t('settings.personal_info_desc')}</p>
                    </div>
                </div>

                <Card className="border-0 shadow-lg bg-card/50 backdrop-blur-sm animate-in fade-in slide-in-from-bottom-4">
                    <form onSubmit={handleSubmit}>
                        <CardContent className="space-y-6 pt-6">
                            {/* Avatar Section */}
                            <div className="flex flex-col items-center gap-4 py-4">
                                <div className="relative group cursor-pointer">
                                    <Avatar className="h-32 w-32 border-4 border-background shadow-xl">
                                        <AvatarImage src={formData.avatar} />
                                        <AvatarFallback className="text-4xl">{formData.nickname?.[0] || 'U'}</AvatarFallback>
                                    </Avatar>
                                    <div className="absolute bottom-0 right-0 p-2 bg-primary text-primary-foreground rounded-full shadow-lg transform group-hover:scale-110 transition-transform">
                                        <Camera className="h-5 w-5" />
                                    </div>
                                </div>
                                <div className="w-full max-w-sm">
                                    <Label className="text-xs text-muted-foreground mb-1 block">{t('settings.avatar_url')}</Label>
                                    <Input
                                        name="avatar"
                                        value={formData.avatar}
                                        onChange={handleChange}
                                        placeholder="https://..."
                                        className="text-center bg-background/50"
                                    />
                                </div>
                            </div>

                            <div className="space-y-4">
                                <div className="grid gap-2">
                                    <Label htmlFor="nickname" className="flex items-center gap-2">
                                        <User className="h-4 w-4 text-muted-foreground" />
                                        {t('auth.nickname')}
                                    </Label>
                                    <Input
                                        id="nickname"
                                        name="nickname"
                                        value={formData.nickname}
                                        onChange={handleChange}
                                        placeholder={t('auth.nickname')}
                                    />
                                </div>

                                <div className="grid gap-2">
                                    <Label htmlFor="email" className="flex items-center gap-2">
                                        <Mail className="h-4 w-4 text-muted-foreground" />
                                        {t('auth.email')}
                                    </Label>
                                    <Input
                                        id="email"
                                        name="email"
                                        type="email"
                                        value={formData.email}
                                        onChange={handleChange}
                                        placeholder="name@example.com"
                                    />
                                </div>

                                <div className="grid gap-2">
                                    <Label className="text-muted-foreground">{t('auth.username')}</Label>
                                    <Input
                                        value={formData.username}
                                        disabled
                                        className="bg-muted opacity-50"
                                    />
                                    <p className="text-[10px] text-muted-foreground">{t('settings.username_locked')}</p>
                                </div>
                            </div>
                        </CardContent>
                        <CardFooter className="flex justify-end gap-3 pb-6">
                            <Button type="button" variant="ghost" onClick={() => navigate('/profile')}>
                                {t('settings.cancel')}
                            </Button>
                            <Button type="submit" disabled={isLoading} className="min-w-[120px]">
                                {isLoading ? <Loader2 className="h-4 w-4 animate-spin mr-2" /> : <Save className="h-4 w-4 mr-2" />}
                                {t('settings.save')}
                            </Button>
                        </CardFooter>
                    </form>
                </Card>
            </div>
        </PageContainer>
    );
}
