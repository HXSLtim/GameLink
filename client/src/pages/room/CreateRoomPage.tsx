import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { ArrowLeft, Mic } from 'lucide-react';
import { PageContainer, PageHeader } from '@/components/page-container';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Textarea } from '@/components/ui/textarea';
import { Switch } from '@/components/ui/switch';
import {
    Select,
    SelectContent,
    SelectItem,
    SelectTrigger,
    SelectValue,
} from '@/components/ui/select';
import { Card, CardContent } from '@/components/ui/card';
import { useRoomStore } from '@/stores';
import type { ChatGroupType } from '@/stores/modules/room-store';

export default function CreateRoomPage() {
    const { t } = useTranslation();
    const navigate = useNavigate();
    const { createRoom, isLoading } = useRoomStore();

    const [formData, setFormData] = useState({
        name: '',
        groupType: 'team' as ChatGroupType,
        gameId: '',
        maxMembers: '5',
        isPrivate: false,
        password: '',
        description: '',
        voiceEnabled: false,
    });

    const [errors, setErrors] = useState<Record<string, string>>({});

    const validate = () => {
        const newErrors: Record<string, string> = {};
        if (!formData.name.trim()) {
            newErrors.name = t('room.create.errors.nameRequired');
        }
        if (!formData.gameId) {
            newErrors.gameId = t('room.create.errors.gameRequired');
        }
        if (formData.isPrivate && !formData.password.trim()) {
            newErrors.password = t('room.create.errors.passwordRequired');
        }
        setErrors(newErrors);
        return Object.keys(newErrors).length === 0;
    };

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        if (!validate()) return;

        try {
            const room = await createRoom({
                name: formData.name.trim(),
                groupType: formData.groupType,
                gameId: Number(formData.gameId),
                maxMembers: Number(formData.maxMembers),
                isPrivate: formData.isPrivate,
                password: formData.password,
                description: formData.description.trim(),
                voiceEnabled: formData.voiceEnabled,
            });
            navigate(`/rooms/${room.id}`);
        } catch (error) {
            console.error('Failed to create room:', error);
        }
    };

    return (
        <PageContainer>
            <div className="flex items-center gap-4 mb-6">
                <Button variant="ghost" size="icon" onClick={() => navigate('/rooms')}>
                    <ArrowLeft className="h-5 w-5" />
                </Button>
                <PageHeader
                    title={t('room.create')}
                    description={t('room.description')}
                />
            </div>

            <Card className="max-w-2xl">
                <CardContent className="pt-6">
                    <form onSubmit={handleSubmit} className="space-y-6">
                        {/* Room Name */}
                        <div className="space-y-2">
                            <Label htmlFor="name">{t('room.createForm.name')}</Label>
                            <Input
                                id="name"
                                value={formData.name}
                                onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                                placeholder={t('room.createForm.namePlaceholder')}
                                maxLength={64}
                            />
                            {errors.name && (
                                <p className="text-sm text-destructive">{errors.name}</p>
                            )}
                        </div>

                        {/* Room Type */}
                        <div className="space-y-2">
                            <Label>{t('room.createForm.type')}</Label>
                            <Select
                                value={formData.groupType}
                                onValueChange={(value) =>
                                    setFormData({ ...formData, groupType: value as ChatGroupType })
                                }
                            >
                                <SelectTrigger>
                                    <SelectValue />
                                </SelectTrigger>
                                <SelectContent>
                                    <SelectItem value="team">{t('room.type.team')}</SelectItem>
                                    <SelectItem value="lfg">{t('room.type.lfg')}</SelectItem>
                                    <SelectItem value="custom">{t('room.type.custom')}</SelectItem>
                                </SelectContent>
                            </Select>
                        </div>

                        {/* Game Selection */}
                        <div className="space-y-2">
                            <Label htmlFor="gameId">{t('room.createForm.game')}</Label>
                            <Input
                                id="gameId"
                                type="number"
                                value={formData.gameId}
                                onChange={(e) => setFormData({ ...formData, gameId: e.target.value })}
                                placeholder={t('room.createForm.gamePlaceholder')}
                            />
                            {errors.gameId && (
                                <p className="text-sm text-destructive">{errors.gameId}</p>
                            )}
                        </div>

                        {/* Max Members */}
                        <div className="space-y-2">
                            <Label>{t('room.create.maxMembers')}</Label>
                            <Select
                                value={formData.maxMembers}
                                onValueChange={(value) =>
                                    setFormData({ ...formData, maxMembers: value })
                                }
                            >
                                <SelectTrigger>
                                    <SelectValue />
                                </SelectTrigger>
                                <SelectContent>
                                    {[2, 3, 4, 5, 6, 8, 10, 15, 20].map((n) => (
                                        <SelectItem key={n} value={String(n)}>
                                            {n} {t('room.create.people')}
                                        </SelectItem>
                                    ))}
                                </SelectContent>
                            </Select>
                        </div>

                        {/* Description */}
                        <div className="space-y-2">
                            <Label htmlFor="description">{t('room.create.descriptionLabel')}</Label>
                            <Textarea
                                id="description"
                                value={formData.description}
                                onChange={(e) =>
                                    setFormData({ ...formData, description: e.target.value })
                                }
                                placeholder={t('room.create.descriptionPlaceholder')}
                                rows={3}
                                maxLength={256}
                            />
                        </div>

                        {/* Private Room */}
                        <div className="flex items-center justify-between">
                            <div>
                                <Label>{t('room.create.private')}</Label>
                                <p className="text-sm text-muted-foreground">
                                    {t('room.create.privateDescription')}
                                </p>
                            </div>
                            <Switch
                                checked={formData.isPrivate}
                                onCheckedChange={(checked) =>
                                    setFormData({ ...formData, isPrivate: checked })
                                }
                            />
                        </div>

                        {/* Password (if private) */}
                        {formData.isPrivate && (
                            <div className="space-y-2">
                                <Label htmlFor="password">{t('room.create.password')}</Label>
                                <Input
                                    id="password"
                                    type="password"
                                    value={formData.password}
                                    onChange={(e) =>
                                        setFormData({ ...formData, password: e.target.value })
                                    }
                                    placeholder={t('room.create.passwordPlaceholder')}
                                    maxLength={32}
                                />
                                {errors.password && (
                                    <p className="text-sm text-destructive">{errors.password}</p>
                                )}
                            </div>
                        )}

                        {/* Voice Enabled */}
                        <div className="flex items-center justify-between">
                            <div className="flex items-center gap-2">
                                <Mic className="h-4 w-4 text-muted-foreground" />
                                <div>
                                    <Label>{t('room.create.voice')}</Label>
                                    <p className="text-sm text-muted-foreground">
                                        {t('room.create.voiceDescription')}
                                    </p>
                                </div>
                            </div>
                            <Switch
                                checked={formData.voiceEnabled}
                                onCheckedChange={(checked) =>
                                    setFormData({ ...formData, voiceEnabled: checked })
                                }
                            />
                        </div>

                        {/* Submit */}
                        <div className="flex gap-3 pt-4">
                            <Button
                                type="button"
                                variant="outline"
                                onClick={() => navigate('/rooms')}
                                className="flex-1"
                            >
                                {t('common.cancel')}
                            </Button>
                            <Button type="submit" disabled={isLoading} className="flex-1">
                                {isLoading ? t('room.createForm.creating') : t('room.createForm.submit')}
                            </Button>
                        </div>
                    </form>
                </CardContent>
            </Card>
        </PageContainer>
    );
}
