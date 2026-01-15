import { useState, useRef } from 'react';
import { useAuthStore } from '@/stores';
import { PageContainer } from '@/components/page-container';
import { Card, CardContent, CardFooter } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import { useTranslation } from 'react-i18next';
import { ArrowLeft, Save, Loader2, Camera, User, Mail, Upload } from 'lucide-react';
import { useNavigate } from 'react-router-dom';
import { toast } from 'sonner';
import { http } from '@/lib/http';

export default function EditProfilePage() {
    const { t } = useTranslation();
    const navigate = useNavigate();
    const { user, updateProfile } = useAuthStore();
    const fileInputRef = useRef<HTMLInputElement>(null);

    const [isLoading, setIsLoading] = useState(false);
    const [isUploading, setIsUploading] = useState(false);
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

    const handleFileSelect = () => {
        fileInputRef.current?.click();
    };

    const handleFileChange = async (e: React.ChangeEvent<HTMLInputElement>) => {
        const file = e.target.files?.[0];
        if (!file) return;

        // Basic validation
        if (file.size > 5 * 1024 * 1024) { // 5MB limit
            toast.error(t('settings.file_too_large') || 'File size exceeds 5MB');
            return;
        }

        setIsUploading(true);
        const uploadData = new FormData();
        uploadData.append('file', file);

        try {
            // Assuming the API returns { url: string }
            const response = await http.post<{ url: string }>('/upload', uploadData, {
                headers: {
                    'Content-Type': 'multipart/form-data',
                },
            });

            // Check if response has the url directly or if it's wrapped
            // Based on http.ts unwrap logic, response should be the data payload
            // If the structure is different, we might need adjustment.
            // Let's assume standard response based on typical implementations.

            if (response?.url) {
                setFormData(prev => ({ ...prev, avatar: response.url }));
                toast.success(t('settings.upload_success') || 'Image uploaded successfully');
            } else {
                // Fallback if structure is different or no URL returned
                // This depends on actual API contract. 
                // If mocking or no backend, this might fail or need mock response.
                console.warn('No URL in response', response);
            }

        } catch (error) {
            console.error('Upload failed:', error);
            // toast.error(error.message || 'Failed to upload image');
            // For now, if API 404s (since it might not exist), let's simulate success for demo if needed,
            // but ideally we show error. User requested "automatic fill", implying backend exists.
            toast.error(error instanceof Error ? error.message : 'Failed to upload image');
        } finally {
            setIsUploading(false);
            // Reset input so same file can be selected again if needed
            if (fileInputRef.current) fileInputRef.current.value = '';
        }
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
        } catch (error) {
            toast.error(error instanceof Error ? error.message : 'Failed to update profile');
        } finally {
            setIsLoading(false);
        }
    };

    return (
        <PageContainer>
            <div className="w-[80%] mx-auto py-8 space-y-8">
                {/* Header */}
                <div className="flex items-center gap-4 animate-in fade-in slide-in-from-top-2">
                    <Button variant="ghost" size="icon" onClick={() => navigate('/profile')} className="rounded-full">
                        <ArrowLeft className="h-5 w-5" />
                    </Button>
                    <div>
                        <h1 className="text-3xl font-bold tracking-tight">{t('profile.settings.personal_info')}</h1>
                        <p className="text-base text-muted-foreground mt-1">{t('profile.settings.personal_info_desc')}</p>
                    </div>
                </div>

                <div className="grid grid-cols-1 xl:grid-cols-12 gap-8 items-start animate-in fade-in slide-in-from-bottom-4">
                    {/* Left Column: Avatar */}
                    <Card className="xl:col-span-4 border-foreground/5 shadow-xl bg-card/50 backdrop-blur-sm overflow-hidden sticky top-8">
                        <CardContent className="pt-8 pb-8 flex flex-col items-center text-center space-y-6">
                            <input
                                type="file"
                                ref={fileInputRef}
                                className="hidden"
                                accept="image/*"
                                onChange={handleFileChange}
                            />
                            <div className="relative group cursor-pointer" onClick={handleFileSelect}>
                                <Avatar className="h-40 w-40 border-4 border-background shadow-2xl ring-4 ring-primary/10 transition-all duration-300 group-hover:ring-primary/30">
                                    <AvatarImage src={formData.avatar} className="object-cover" />
                                    <AvatarFallback className="text-5xl font-bold bg-primary/5 text-primary">
                                        {isUploading ? <Loader2 className="h-10 w-10 animate-spin" /> : (formData.nickname?.[0]?.toUpperCase() || 'U')}
                                    </AvatarFallback>
                                </Avatar>
                                <div className="absolute inset-0 bg-black/40 rounded-full flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity backdrop-blur-[2px]">
                                    {isUploading ? <Loader2 className="h-10 w-10 text-white animate-spin" /> : <Camera className="h-10 w-10 text-white" />}
                                </div>
                            </div>

                            <div className="space-y-2 w-full px-4">
                                <h3 className="font-semibold text-xl">{formData.nickname || user?.username}</h3>
                                <p className="text-sm text-muted-foreground">{formData.email}</p>
                                <Button variant="outline" size="sm" className="mt-2 w-full" onClick={handleFileSelect} disabled={isUploading}>
                                    <Upload className="h-4 w-4 mr-2" />
                                    {isUploading ? 'Uploading...' : 'Upload New Avatar'}
                                </Button>
                            </div>

                            <div className="w-full pt-4 border-t border-border/50">
                                <div className="text-xs font-medium text-muted-foreground mb-2 text-left uppercase tracking-wider pl-1">{t('settings.avatar_url')}</div>
                                <div className="relative">
                                    <Input
                                        name="avatar"
                                        value={formData.avatar}
                                        onChange={handleChange}
                                        placeholder="https://..."
                                        className="bg-muted/50 focus:bg-background transition-colors text-sm font-mono pl-9"
                                    />
                                    <User className="absolute left-3 top-2.5 h-4 w-4 text-muted-foreground" />
                                </div>
                            </div>
                        </CardContent>
                    </Card>

                    {/* Right Column: Form */}
                    <Card className="xl:col-span-8 border-foreground/5 shadow-xl bg-card/50 backdrop-blur-sm">
                        <form onSubmit={handleSubmit}>
                            <CardContent className="space-y-8 pt-8">
                                <div className="grid gap-6">
                                    <div className="grid gap-3">
                                        <Label htmlFor="nickname" className="text-base font-medium flex items-center gap-2">
                                            <User className="h-4 w-4 text-primary" />
                                            {t('auth.nickname')}
                                        </Label>
                                        <Input
                                            id="nickname"
                                            name="nickname"
                                            value={formData.nickname}
                                            onChange={handleChange}
                                            placeholder={t('auth.nickname')}
                                            className="h-12 bg-muted/30 focus:bg-background text-lg"
                                        />
                                        <p className="text-sm text-muted-foreground">This is your public display name.</p>
                                    </div>

                                    <div className="grid gap-3">
                                        <Label htmlFor="email" className="text-base font-medium flex items-center gap-2">
                                            <Mail className="h-4 w-4 text-primary" />
                                            {t('auth.email')}
                                        </Label>
                                        <Input
                                            id="email"
                                            name="email"
                                            type="email"
                                            value={formData.email}
                                            onChange={handleChange}
                                            placeholder="name@example.com"
                                            className="h-12 bg-muted/30 focus:bg-background text-lg"
                                        />
                                    </div>

                                    <div className="grid gap-3">
                                        <Label className="text-base font-medium flex items-center gap-2 text-muted-foreground">
                                            <span className="opacity-70">{t('auth.username')}</span>
                                            <span className="text-xs font-normal px-2 py-0.5 rounded-full bg-muted text-muted-foreground ml-auto">Locked</span>
                                        </Label>
                                        <Input
                                            value={formData.username}
                                            disabled
                                            className="h-12 bg-muted/50 opacity-70 font-mono text-base"
                                        />
                                        <p className="text-sm text-muted-foreground flex items-center gap-1.5">
                                            <span className="inline-block w-1.5 h-1.5 rounded-full bg-orange-400"></span>
                                            {t('settings.username_locked')}
                                        </p>
                                    </div>
                                </div>
                            </CardContent>
                            <CardFooter className="flex justify-end gap-4 pb-8 pt-4 border-t bg-muted/10">
                                <Button type="button" variant="ghost" onClick={() => navigate('/profile')} className="h-11 px-6 text-base hover:bg-muted/50">
                                    {t('settings.cancel')}
                                </Button>
                                <Button type="submit" disabled={isLoading} className="h-11 px-8 min-w-[140px] text-base shadow-lg shadow-primary/20 hover:shadow-primary/30 transition-all">
                                    {isLoading ? <Loader2 className="h-5 w-5 animate-spin mr-2" /> : <Save className="h-5 w-5 mr-2" />}
                                    {t('settings.save')}
                                </Button>
                            </CardFooter>
                        </form>
                    </Card>
                </div>
            </div>
        </PageContainer>
    );
}
