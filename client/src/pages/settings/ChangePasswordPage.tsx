import { useState } from 'react';
import { useAuthStore } from '@/stores';
import { PageContainer } from '@/components/page-container';
import { Card, CardContent, CardFooter } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { useTranslation } from 'react-i18next';
import { ArrowLeft, Save, Loader2, KeyRound, Lock, ShieldCheck } from 'lucide-react';
import { useNavigate } from 'react-router-dom';
import { toast } from 'sonner';

export default function ChangePasswordPage() {
    const { t } = useTranslation();
    const navigate = useNavigate();
    const { changePassword } = useAuthStore();

    const [isLoading, setIsLoading] = useState(false);
    const [formData, setFormData] = useState({
        oldPassword: '',
        newPassword: '',
        confirmPassword: ''
    });

    const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const { name, value } = e.target;
        setFormData(prev => ({ ...prev, [name]: value }));
    };

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();

        if (formData.newPassword !== formData.confirmPassword) {
            toast.error(t('settings.password_mismatch'));
            return;
        }

        if (formData.newPassword.length < 6) {
            toast.error(t('auth.password_min_length'));
            return;
        }

        setIsLoading(true);
        try {
            await changePassword({
                oldPassword: formData.oldPassword,
                newPassword: formData.newPassword
            });
            toast.success(t('settings.update_success'));
            setTimeout(() => navigate('/profile'), 500);
        } catch (error) {
            toast.error(error instanceof Error ? error.message : 'Failed to change password');
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
                        <h1 className="text-2xl font-bold tracking-tight">{t('settings.security')}</h1>
                        <p className="text-sm text-muted-foreground">{t('settings.security_desc')}</p>
                    </div>
                </div>

                <div className="grid gap-6 animate-in fade-in slide-in-from-bottom-4">
                    <Card className="border-l-4 border-l-yellow-500 bg-yellow-500/5 border-t-0 border-b-0 border-r-0 shadow-sm">
                        <CardContent className="pt-6 flex gap-4">
                            <ShieldCheck className="h-6 w-6 text-yellow-600 shrink-0" />
                            <div className="text-sm text-muted-foreground">
                                <p className="font-medium text-foreground mb-1">{t('settings.security_tip')}</p>
                                {t('settings.tip_content')}
                            </div>
                        </CardContent>
                    </Card>

                    <Card className="border-0 shadow-lg bg-card/50 backdrop-blur-sm">
                        <form onSubmit={handleSubmit}>
                            <CardContent className="space-y-4 pt-6">
                                <div className="grid gap-2">
                                    <Label htmlFor="oldPassword">{t('settings.current_password')}</Label>
                                    <div className="relative">
                                        <KeyRound className="absolute left-3 top-3 h-4 w-4 text-muted-foreground" />
                                        <Input
                                            id="oldPassword"
                                            name="oldPassword"
                                            type="password"
                                            className="pl-9"
                                            value={formData.oldPassword}
                                            onChange={handleChange}
                                            placeholder={t('settings.current_password')}
                                        />
                                    </div>
                                </div>

                                <div className="grid gap-2">
                                    <Label htmlFor="newPassword">{t('settings.new_password')}</Label>
                                    <div className="relative">
                                        <Lock className="absolute left-3 top-3 h-4 w-4 text-muted-foreground" />
                                        <Input
                                            id="newPassword"
                                            name="newPassword"
                                            type="password"
                                            className="pl-9"
                                            value={formData.newPassword}
                                            onChange={handleChange}
                                            placeholder={t('settings.new_password')}
                                        />
                                    </div>
                                </div>

                                <div className="grid gap-2">
                                    <Label htmlFor="confirmPassword">{t('settings.confirm_password')}</Label>
                                    <div className="relative">
                                        <Lock className="absolute left-3 top-3 h-4 w-4 text-muted-foreground" />
                                        <Input
                                            id="confirmPassword"
                                            name="confirmPassword"
                                            type="password"
                                            className="pl-9"
                                            value={formData.confirmPassword}
                                            onChange={handleChange}
                                            placeholder={t('settings.confirm_password')}
                                        />
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
            </div>
        </PageContainer>
    );
}
